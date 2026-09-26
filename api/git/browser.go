package git

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gohyuhan/gitti/executor"
)

// the reasons a web page URL cannot be built
var (
	ErrNoRemote              = errors.New("no remote")
	ErrUnsupportedRemoteHost = errors.New("unsupported remote host (only github, gitlab and bitbucket)")
	ErrBranchNotPushed       = errors.New("branch is not on the remote, push it first")
	ErrCommitNotOnRemote     = errors.New("commit is not on the remote, push it first")
	ErrCommitVerification    = errors.New("cannot verify whether the commit is on the remote")
	ErrAmbiguousPushTarget   = errors.New("branch has no clear push target")
	ErrAmbiguousRemoteBranch = errors.New("cannot determine which remote branch belongs to this local branch")
)

// the longest wait for a remote to list its branches, before the last fetch is used instead
const remoteBranchListTimeout = 10 * time.Second

const remoteCommitProofTimeout = 15 * time.Second

// the ssh endpoints on port 443, whose web pages are on the main host
var sshPort443HostToWebHost = map[string]string{
	"ssh.github.com":    "github.com",
	"altssh.gitlab.com": "gitlab.com",
}

// ------------------------------------
//
//	Return the web page URL of a commit. A commit hash names the same commit on every host, so
//	every remote that has the commit gives a correct page. The order: the remote the checked-out
//	branch tracks, the remote git pushes it to, then the other remotes. A fork can lack the
//	commits that are only on origin, and origin lacks the commits pushed only to the fork.
//
// ------------------------------------
func GetCommitWebPageURL(commitHash string, checkedOutBranchName string, repoPath string) (string, error) {
	cmdExecutor := executor.NewCmdExecutor(repoPath)
	remoteNames := getRemoteNames(cmdExecutor)
	var candidateRemotes []commitWebCandidate
	if upstreamRemoteName, _, hasUpstream := getBranchUpstream(cmdExecutor, checkedOutBranchName); hasUpstream {
		candidateRemotes = append(candidateRemotes, commitWebCandidate{upstreamRemoteName, false})
	} else if slices.Contains(remoteNames, "origin") {
		candidateRemotes = append(candidateRemotes, commitWebCandidate{"origin", false})
	}
	candidateRemotes = append(candidateRemotes, commitWebCandidate{getPushRemoteName(cmdExecutor, checkedOutBranchName, remoteNames), true})
	for _, remoteName := range remoteNames {
		candidate := commitWebCandidate{remoteName, false}
		if !slices.Contains(candidateRemotes, candidate) {
			candidateRemotes = append(candidateRemotes, candidate)
		}
	}

	var firstErr error
	for remoteIndex, candidate := range candidateRemotes {
		webPageURL, err := getCommitWebPageURLOnRemote(cmdExecutor, repoPath, commitHash, candidate)
		if err == nil {
			return webPageURL, nil
		}
		if remoteIndex == 0 {
			// the first remote is the one the user works with: its unsupported host is the answer,
			// not the page of another host
			if errors.Is(err, ErrUnsupportedRemoteHost) {
				return "", err
			}
			firstErr = err
		}
	}
	if firstErr == nil {
		return "", ErrNoRemote
	}
	return "", firstErr
}

type commitWebCandidate struct {
	remoteName string
	usePushURL bool
}

// ------------------------------------
//
//	Return the web page URL of a commit on one remote. Check its current branch and tag tips;
//	remote-tracking refs can be stale after a push or branch deletion.
//
// ------------------------------------
func getCommitWebPageURLOnRemote(cmdExecutor *executor.CmdExecutor, repoPath string, commitHash string, candidate commitWebCandidate) (string, error) {
	remoteURL, err := getRemoteURL(cmdExecutor, candidate.remoteName, candidate.usePushURL)
	if err != nil {
		return "", err
	}
	webPageURL, err := buildWebPageURL(remoteURL, WEBPAGECOMMIT, commitHash)
	if err != nil {
		return "", err
	}

	remoteTips, isRemoteReached := getRemoteTips(cmdExecutor, remoteURL, true)
	if isRemoteReached {
		var unknownTips []string
		seenTips := make(map[string]bool)
		for _, remoteTip := range remoteTips {
			if isAncestor(cmdExecutor, commitHash, remoteTip) {
				return webPageURL, nil
			}
			if !seenTips[remoteTip] && !isLocalCommitObject(cmdExecutor, remoteTip) {
				unknownTips = append(unknownTips, remoteTip)
				seenTips[remoteTip] = true
			}
		}
		if len(unknownTips) > 0 {
			containsCommit, proofErr := checkUnknownRemoteTips(repoPath, remoteURL, commitHash, unknownTips)
			if containsCommit {
				return webPageURL, nil
			}
			if proofErr != nil {
				return "", fmt.Errorf("%w: %s", ErrCommitVerification, candidate.remoteName)
			}
		}
	} else {
		// An offline remote cannot be checked live; use the last fetched refs.
		if canUseTrackingRef(cmdExecutor, candidate.remoteName, remoteURL, candidate.usePushURL) {
			gitArgs := []string{"for-each-ref", "--contains", commitHash, "--format=%(refname)", "refs/remotes/" + candidate.remoteName + "/"}
			output, _ := cmdExecutor.RunGitCmd(gitArgs, false).Output()
			if len(output) > 0 {
				return webPageURL, nil
			}
		}
	}
	return "", fmt.Errorf("%w: %s", ErrCommitNotOnRemote, candidate.remoteName)
}

func isLocalCommitObject(cmdExecutor *executor.CmdExecutor, objectID string) bool {
	return cmdExecutor.RunGitCmd([]string{"cat-file", "-e", objectID + "^{commit}"}, false).Run() == nil
}

// Fetch unknown advertised tips into a short-lived shared clone. This lets Git prove ancestry
// without changing the user's refs, object store, or FETCH_HEAD.
func checkUnknownRemoteTips(repoPath string, remoteURL string, commitHash string, tips []string) (bool, error) {
	tempRoot, err := os.MkdirTemp("", "gitti-commit-lookup-")
	if err != nil {
		return false, err
	}
	defer os.RemoveAll(tempRoot)
	proofRepo := filepath.Join(tempRoot, "repo.git")
	proofContext, cancel := context.WithTimeout(context.Background(), remoteCommitProofTimeout)
	defer cancel()

	localExecutor := executor.NewCmdExecutor(repoPath)
	cloneCmd := localExecutor.RunGitCmdWithContext(proofContext, []string{"clone", "--bare", "--shared", "--quiet", "--", repoPath, proofRepo}, false)
	cloneCmd.WaitDelay = time.Second
	if err := cloneCmd.Run(); err != nil {
		return false, err
	}

	proofExecutor := executor.NewCmdExecutor(proofRepo)
	fetchArgs := append([]string{"fetch", "--quiet", "--no-tags", "--filter=blob:none", remoteURL}, tips...)
	fetchCmd := proofExecutor.RunGitCmdWithContext(proofContext, fetchArgs, false)
	fetchCmd.WaitDelay = time.Second
	if err := fetchCmd.Run(); err != nil {
		return false, err
	}
	for _, tip := range tips {
		if isAncestor(proofExecutor, commitHash, tip) {
			return true, nil
		}
	}
	return false, nil
}

// ------------------------------------
//
//	Return the web page URL of a local branch or a new pull request for it. Prefer a published
//	copy whose tip belongs to the local branch. This distinguishes a renamed branch tracking
//	its old name from a feature branch tracking a base branch.
//
// ------------------------------------
func GetBranchWebPageURL(branchName string, pageType string, repoPath string) (string, error) {
	cmdExecutor := executor.NewCmdExecutor(repoPath)
	pushRemoteName := getPushRemoteName(cmdExecutor, branchName, getRemoteNames(cmdExecutor))
	pushBranchName, hasPushRefspec, err := getConfiguredPushBranch(cmdExecutor, pushRemoteName, branchName)
	if err != nil {
		return "", err
	}
	candidates := []branchWebCandidate{{
		remoteName:  pushRemoteName,
		branchName:  pushBranchName,
		trackingRef: "refs/remotes/" + pushRemoteName + "/" + pushBranchName,
		usePushURL:  true,
		eligible:    true,
	}}
	if upstreamRemoteName, upstreamBranchName, hasUpstream := getBranchUpstream(cmdExecutor, branchName); hasUpstream && !hasPushRefspec {
		renamedFromUpstream := upstreamBranchName != branchName && wasRenamedFrom(cmdExecutor, branchName, upstreamBranchName)
		upstreamIsPushTarget := upstreamRemoteName == pushRemoteName && pushesToUpstream(cmdExecutor)
		upstreamIsCandidate := upstreamBranchName == branchName || upstreamIsPushTarget || renamedFromUpstream
		if (upstreamRemoteName != pushRemoteName || upstreamBranchName != pushBranchName) &&
			(upstreamIsCandidate || upstreamRemoteName == pushRemoteName) {
			candidates = append(candidates, branchWebCandidate{
				remoteName:        upstreamRemoteName,
				branchName:        upstreamBranchName,
				trackingRef:       branchName + "@{upstream}",
				usePushURL:        upstreamIsPushTarget,
				isPushTarget:      upstreamIsPushTarget,
				preferIfUncertain: renamedFromUpstream,
				eligible:          upstreamIsCandidate,
			})
		}
	}

	var firstErr error
	tipLookups := make(map[string]remoteTipLookup)
	for candidateIndex := range candidates {
		candidate := &candidates[candidateIndex]
		remoteURL, err := getRemoteURL(cmdExecutor, candidate.remoteName, candidate.usePushURL)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		candidate.remoteURL = remoteURL
		candidate.webPageURL, candidate.urlErr = buildWebPageURL(remoteURL, pageType, candidate.branchName)
		lookup, found := tipLookups[remoteURL]
		if !found {
			lookup.tips, lookup.reached = getRemoteTips(cmdExecutor, remoteURL, false)
			tipLookups[remoteURL] = lookup
		}
		if lookup.reached {
			candidate.tip, candidate.exists = lookup.tips[candidate.branchName]
		} else if canUseTrackingRef(cmdExecutor, candidate.remoteName, remoteURL, candidate.usePushURL) {
			// Offline fallback uses the last fetched ref; a reachable remote always wins over stale refs.
			gitArgs := []string{"rev-parse", "--verify", "--quiet", candidate.trackingRef}
			output, err := cmdExecutor.RunGitCmd(gitArgs, false).Output()
			if err == nil {
				candidate.tip, candidate.exists = strings.TrimSpace(string(output)), true
			}
		}
	}

	ancestryExecutor, cleanup, err := branchAncestryExecutor(cmdExecutor, repoPath, candidates)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrAmbiguousRemoteBranch, branchName)
	}
	defer cleanup()

	selectedCandidate := -1
	if candidates[0].exists {
		selectedCandidate = 0
	}
	if len(candidates) > 1 && candidates[1].exists {
		if selectedCandidate < 0 && candidates[1].eligible {
			if candidates[1].preferIfUncertain && !isBranchHistoryRelated(ancestryExecutor, branchName, candidates[1].tip) {
				return "", fmt.Errorf("%w: %s", ErrAmbiguousRemoteBranch, branchName)
			}
			selectedCandidate = 1
		} else if selectedCandidate >= 0 && candidates[1].preferIfUncertain && !candidates[1].isPushTarget && candidates[0].tip == candidates[1].tip {
			return "", fmt.Errorf("%w: %s", ErrAmbiguousRemoteBranch, branchName)
		} else if selectedCandidate >= 0 && candidates[1].eligible && shouldUseUpstream(ancestryExecutor, branchName, candidates[0], candidates[1]) {
			selectedCandidate = 1
		} else if selectedCandidate >= 0 && candidates[1].preferIfUncertain && !isBranchHistoryRelated(ancestryExecutor, branchName, candidates[1].tip) {
			return "", fmt.Errorf("%w: %s", ErrAmbiguousRemoteBranch, branchName)
		} else if selectedCandidate >= 0 && !candidates[1].eligible && hasConflictingUpstream(ancestryExecutor, branchName, candidates[0], candidates[1]) {
			return "", fmt.Errorf("%w: %s", ErrAmbiguousRemoteBranch, branchName)
		}
	}
	if selectedCandidate >= 0 {
		candidate := candidates[selectedCandidate]
		return candidate.webPageURL, candidate.urlErr
	}
	if firstErr != nil {
		return "", firstErr
	}
	return "", fmt.Errorf("%w: %s", ErrBranchNotPushed, branchName)
}

func pushesToUpstream(cmdExecutor *executor.CmdExecutor) bool {
	pushDefault := getGitConfigValue(cmdExecutor, "push.default")
	return pushDefault == "upstream" || pushDefault == "tracking"
}

// remote.<name>.push takes precedence over push.default. Resolve the refspecs that push this
// branch, and reject an ambiguous destination instead of opening another branch's page.
func getConfiguredPushBranch(cmdExecutor *executor.CmdExecutor, remoteName string, branchName string) (string, bool, error) {
	gitArgs := []string{"config", "--get-all", "remote." + remoteName + ".push"}
	output, err := cmdExecutor.RunGitCmd(gitArgs, false).Output()
	if err != nil || len(output) == 0 {
		return branchName, false, nil
	}

	currentBranch := strings.TrimSpace(stringOutput(cmdExecutor, "symbolic-ref", "--quiet", "--short", "HEAD"))
	pushBranch := ""
	for refspec := range strings.SplitSeq(strings.TrimSpace(string(output)), "\n") {
		if refspec == "" {
			continue
		}
		remoteBranch, matches, supported := pushBranchFromRefspec(refspec, branchName, currentBranch)
		if matches && (!supported || pushBranch != "" && pushBranch != remoteBranch) {
			return "", true, fmt.Errorf("%w: %s", ErrAmbiguousPushTarget, branchName)
		}
		if matches {
			pushBranch = remoteBranch
		}
	}
	if pushBranch != "" {
		return pushBranch, true, nil
	}
	return "", true, fmt.Errorf("%w: no refspec for %s", ErrAmbiguousPushTarget, branchName)
}

func pushBranchFromRefspec(refspec string, branchName string, currentBranch string) (remoteBranch string, matches bool, supported bool) {
	refspec = strings.TrimPrefix(refspec, "+")
	if refspec == ":" {
		return branchName, true, true
	}
	source, destination, hasDestination := strings.Cut(refspec, ":")
	if !hasDestination {
		destination = source
	}
	if source == "HEAD" {
		if branchName != currentBranch {
			return "", false, true
		}
		if !hasDestination {
			destination = branchName
		}
	} else {
		source = strings.TrimPrefix(source, "refs/heads/")
		if source == branchName {
			// exact source
		} else if strings.Count(source, "*") == 1 && strings.Count(destination, "*") == 1 {
			prefix, suffix, _ := strings.Cut(source, "*")
			if !strings.HasPrefix(branchName, prefix) || !strings.HasSuffix(branchName, suffix) || len(branchName) < len(prefix)+len(suffix) {
				return "", false, true
			}
			middle := branchName[len(prefix) : len(branchName)-len(suffix)]
			destination = strings.Replace(destination, "*", middle, 1)
		} else {
			return "", false, true
		}
	}
	if strings.HasPrefix(destination, "refs/") && !strings.HasPrefix(destination, "refs/heads/") {
		return "", true, false
	}
	destination = strings.TrimPrefix(destination, "refs/heads/")
	if destination == "" || strings.Contains(destination, "*") {
		return "", true, false
	}
	return destination, true, true
}

func stringOutput(cmdExecutor *executor.CmdExecutor, gitArgs ...string) string {
	output, _ := cmdExecutor.RunGitCmd(gitArgs, false).Output()
	return string(output)
}

// A differently named upstream is often a base branch used for pulls. A branch rename leaves
// an explicit reflog entry, which distinguishes a published old name from that base branch.
func wasRenamedFrom(cmdExecutor *executor.CmdExecutor, branchName string, upstreamBranchName string) bool {
	gitArgs := []string{"reflog", "show", "--format=%gs", "refs/heads/" + branchName}
	output, err := cmdExecutor.RunGitCmd(gitArgs, false).Output()
	if err != nil {
		return false
	}
	previousNames := map[string]bool{branchName: true}
	for line := range strings.SplitSeq(string(output), "\n") {
		rename, isRename := strings.CutPrefix(line, "Branch: renamed refs/heads/")
		if !isRename {
			continue
		}
		oldName, newName, found := strings.Cut(rename, " to refs/heads/")
		if found && previousNames[newName] {
			previousNames[oldName] = true
		}
	}
	return previousNames[upstreamBranchName]
}

type branchWebCandidate struct {
	remoteName        string
	remoteURL         string
	branchName        string
	trackingRef       string
	usePushURL        bool
	isPushTarget      bool
	preferIfUncertain bool
	eligible          bool
	webPageURL        string
	urlErr            error
	tip               string
	exists            bool
}

// A remote tip may not have been fetched yet. Resolve its history in a temporary clone
// before comparing candidates, leaving the user's repository untouched.
func branchAncestryExecutor(cmdExecutor *executor.CmdExecutor, repoPath string, candidates []branchWebCandidate) (*executor.CmdExecutor, func(), error) {
	if len(candidates) < 2 || !candidates[1].exists || (!candidates[0].exists && !candidates[1].preferIfUncertain) {
		return cmdExecutor, func() {}, nil
	}
	var missing []branchWebCandidate
	for _, candidate := range candidates {
		if candidate.exists && !isLocalCommitObject(cmdExecutor, candidate.tip) {
			missing = append(missing, candidate)
		}
	}
	if len(missing) == 0 {
		return cmdExecutor, func() {}, nil
	}

	root, err := os.MkdirTemp("", "gitti-branch-lookup-")
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { _ = os.RemoveAll(root) }
	proofRepo := filepath.Join(root, "repo.git")
	ctx, cancel := context.WithTimeout(context.Background(), remoteCommitProofTimeout)
	defer cancel()
	cloneCmd := cmdExecutor.RunGitCmdWithContext(ctx, []string{"clone", "--bare", "--shared", "--quiet", "--", repoPath, proofRepo}, false)
	cloneCmd.WaitDelay = time.Second
	if err := cloneCmd.Run(); err != nil {
		cleanup()
		return nil, nil, err
	}
	proofExecutor := executor.NewCmdExecutor(proofRepo)
	for _, candidate := range missing {
		fetchCmd := proofExecutor.RunGitCmdWithContext(ctx, []string{"fetch", "--quiet", "--no-tags", "--filter=blob:none", candidate.remoteURL, candidate.tip}, false)
		fetchCmd.WaitDelay = time.Second
		if err := fetchCmd.Run(); err != nil {
			cleanup()
			return nil, nil, err
		}
	}
	return proofExecutor, cleanup, nil
}

type remoteTipLookup struct {
	tips    map[string]string
	reached bool
}

func hasConflictingUpstream(cmdExecutor *executor.CmdExecutor, branchName string, push branchWebCandidate, upstream branchWebCandidate) bool {
	if !isBranchHistoryRelated(cmdExecutor, branchName, upstream.tip) {
		return false
	}
	if push.tip == upstream.tip {
		return true
	}
	localRef := "refs/heads/" + branchName
	localTip := strings.TrimSpace(stringOutput(cmdExecutor, "rev-parse", "--verify", "--quiet", localRef))
	if push.tip == localTip {
		return !isAncestor(cmdExecutor, upstream.tip, localRef)
	}
	// A published feature is clear when both remote tips are in its local history and
	// the same-name tip is newer than the tracked base. Other relationships leave
	// the branch identity unresolved.
	return !isAncestor(cmdExecutor, upstream.tip, push.tip) || !isAncestor(cmdExecutor, push.tip, localRef)
}

func isBranchHistoryRelated(cmdExecutor *executor.CmdExecutor, branchName string, tip string) bool {
	localRef := "refs/heads/" + branchName
	return isAncestor(cmdExecutor, tip, localRef) || isAncestor(cmdExecutor, localRef, tip)
}

func shouldUseUpstream(cmdExecutor *executor.CmdExecutor, branchName string, push branchWebCandidate, upstream branchWebCandidate) bool {
	if upstream.isPushTarget {
		return true
	}
	localRef := "refs/heads/" + branchName
	pushTipIsLocalHistory := isAncestor(cmdExecutor, push.tip, localRef)
	upstreamTipIsLocalHistory := isAncestor(cmdExecutor, upstream.tip, localRef)
	if !pushTipIsLocalHistory && upstreamTipIsLocalHistory {
		return true
	}
	if !upstream.preferIfUncertain {
		return false
	}
	if isAncestor(cmdExecutor, localRef, upstream.tip) {
		return true
	}
	if !pushTipIsLocalHistory {
		return true
	}
	return upstreamTipIsLocalHistory && isAncestor(cmdExecutor, push.tip, upstream.tip)
}

func isAncestor(cmdExecutor *executor.CmdExecutor, ancestor string, descendant string) bool {
	if ancestor == "" || descendant == "" {
		return false
	}
	gitArgs := []string{"merge-base", "--is-ancestor", ancestor, descendant}
	return cmdExecutor.RunGitCmd(gitArgs, false).Run() == nil
}

// ------------------------------------
//
//	Return the upstream of a local branch: branch.<name>.remote and branch.<name>.merge. They are
//	read as separate config values instead of splitting "origin/x", because a remote name can
//	hold "/". A branch that tracks another local branch (remote ".") has no upstream on a remote.
//
// ------------------------------------
func getBranchUpstream(cmdExecutor *executor.CmdExecutor, branchName string) (remoteName string, remoteBranchName string, hasUpstream bool) {
	remoteName = getGitConfigValue(cmdExecutor, "branch."+branchName+".remote")
	mergeRef := getGitConfigValue(cmdExecutor, "branch."+branchName+".merge")
	if remoteName == "" || remoteName == "." || mergeRef == "" {
		return "", "", false
	}
	return remoteName, strings.TrimPrefix(mergeRef, "refs/heads/"), true
}

// ------------------------------------
//
//	Return the remote that git pushes a branch to, in the order of `git help config`:
//	branch.<name>.pushRemote, remote.pushDefault, branch.<name>.remote, then the only remote when
//	there is one (as `git push` does), else "origin". "." is this repository, not a remote.
//
// ------------------------------------
func getPushRemoteName(cmdExecutor *executor.CmdExecutor, branchName string, remoteNames []string) string {
	for _, configKey := range []string{"branch." + branchName + ".pushRemote", "remote.pushDefault", "branch." + branchName + ".remote"} {
		if remoteName := getGitConfigValue(cmdExecutor, configKey); remoteName != "" && remoteName != "." {
			return remoteName
		}
	}
	if len(remoteNames) == 1 {
		return remoteNames[0]
	}
	return "origin"
}

// ------------------------------------
//
//	Return the names of the remotes, read fresh. A remote name holds no whitespace.
//
// ------------------------------------
func getRemoteNames(cmdExecutor *executor.CmdExecutor) []string {
	gitArgs := []string{"remote"}
	remoteListOutput, remoteListErr := cmdExecutor.RunGitCmd(gitArgs, false).Output()
	if remoteListErr != nil {
		return nil
	}
	return strings.Fields(string(remoteListOutput))
}

// ------------------------------------
//
//	Return the branch tips currently advertised by a remote. Commit checks also include tags;
//	annotated tags have a peeled commit entry with a ^{} suffix.
//
// ------------------------------------
func getRemoteTips(cmdExecutor *executor.CmdExecutor, remoteURL string, includeTags bool) (remoteTips map[string]string, isRemoteReached bool) {
	lookupContext, cancel := context.WithTimeout(context.Background(), remoteBranchListTimeout)
	defer cancel()

	// a pattern, not --heads (deprecated in newer git) or --branches (missing in older git)
	gitArgs := []string{"ls-remote", remoteURL, "refs/heads/*"}
	if includeTags {
		gitArgs = append(gitArgs, "refs/tags/*")
	}
	lsRemoteCmd := cmdExecutor.RunGitCmdWithContext(lookupContext, gitArgs, false)
	// the ssh process that git starts can outlive a killed git and keep the output open
	lsRemoteCmd.WaitDelay = time.Second
	lsRemoteOutput, lsRemoteErr := lsRemoteCmd.Output()
	if lsRemoteErr != nil {
		return nil, false
	}

	remoteTips = map[string]string{}
	for line := range strings.SplitSeq(string(lsRemoteOutput), "\n") {
		// each line is "<commit hash>\t<ref>"; the pattern matches the end of a ref, so a ref such as
		// "refs/pull/1/refs/heads/x" comes too and is skipped
		branchTip, remoteRef, _ := strings.Cut(strings.TrimSpace(line), "\t")
		if remoteBranchName, isBranch := strings.CutPrefix(remoteRef, "refs/heads/"); isBranch {
			remoteTips[remoteBranchName] = branchTip
		} else if includeTags && strings.HasPrefix(remoteRef, "refs/tags/") {
			remoteTips[remoteRef] = branchTip
		}
	}
	return remoteTips, true
}

// ------------------------------------
//
//	Return the value of a git config key, or "" when the key is not set
//
// ------------------------------------
func getGitConfigValue(cmdExecutor *executor.CmdExecutor, key string) string {
	gitArgs := []string{"config", "--get", key}
	configOutput, configErr := cmdExecutor.RunGitCmd(gitArgs, false).Output()
	if configErr != nil {
		return ""
	}
	return strings.TrimSpace(string(configOutput))
}

// ------------------------------------
//
//	Read a remote URL, using its push URL for a pushed branch.
//
// ------------------------------------
func getRemoteURL(cmdExecutor *executor.CmdExecutor, remoteName string, push bool) (string, error) {
	gitArgs := []string{"remote", "get-url"}
	if push {
		gitArgs = append(gitArgs, "--push")
	}
	gitArgs = append(gitArgs, remoteName)
	remoteURLOutput, remoteURLErr := cmdExecutor.RunGitCmd(gitArgs, false).Output()
	if remoteURLErr != nil {
		return "", fmt.Errorf("%w: %s", ErrNoRemote, remoteName)
	}
	return strings.TrimSpace(string(remoteURLOutput)), nil
}

// A fetched tracking ref says nothing about a separate push URL when the remote is offline.
func canUseTrackingRef(cmdExecutor *executor.CmdExecutor, remoteName string, remoteURL string, isPushURL bool) bool {
	if !isPushURL {
		return true
	}
	fetchURL, err := getRemoteURL(cmdExecutor, remoteName, false)
	return err == nil && fetchURL == remoteURL
}

// ------------------------------------
//
//	Build the web page URL of a commit (WEBPAGECOMMIT), a branch (WEBPAGEBRANCH), or a new pull
//	request (WEBPAGEPULLREQUEST) from a remote URL. commitHashOrBranchName is the full commit hash
//	for WEBPAGECOMMIT, else the branch name on the remote.
//
// ------------------------------------
func buildWebPageURL(remoteURL string, pageType string, commitHashOrBranchName string) (string, error) {
	host, webBaseURL, err := parseRemoteWebBaseURL(remoteURL)
	if err != nil {
		return "", err
	}

	// a branch name can hold URL characters such as "#", "%" and "?".
	// As a path it keeps "/" (like `gh browse`); as one path segment "/" becomes %2F (like `gh pr create --web`)
	branchNamePathSegments := strings.Split(commitHashOrBranchName, "/")
	for segmentIndex, segment := range branchNamePathSegments {
		branchNamePathSegments[segmentIndex] = url.PathEscape(segment)
	}
	branchNameAsPath := strings.Join(branchNamePathSegments, "/")
	branchNameAsPathSegment := url.PathEscape(commitHashOrBranchName)
	branchNameAsQueryValue := url.QueryEscape(commitHashOrBranchName)

	// a self-hosted "gitlab.example.com" matches by substring
	var commitPath, branchPath, pullRequestPath string
	switch {
	case strings.Contains(host, "github"):
		commitPath = "/commit/" + commitHashOrBranchName
		branchPath = "/tree/" + branchNameAsPath
		pullRequestPath = "/compare/" + branchNameAsPathSegment + "?quick_pull=1"
	case strings.Contains(host, "gitlab"):
		commitPath = "/-/commit/" + commitHashOrBranchName
		branchPath = "/-/tree/" + branchNameAsPath
		pullRequestPath = "/-/merge_requests/new?merge_request%5Bsource_branch%5D=" + branchNameAsQueryValue
	case strings.Contains(host, "bitbucket"):
		commitPath = "/commits/" + commitHashOrBranchName
		branchPath = "/src/" + branchNameAsPath
		pullRequestPath = "/pull-requests/new?source=" + branchNameAsQueryValue + "&t=1"
	default:
		return "", fmt.Errorf("%w: %s", ErrUnsupportedRemoteHost, host)
	}

	switch pageType {
	case WEBPAGECOMMIT:
		return webBaseURL + commitPath, nil
	case WEBPAGEBRANCH:
		return webBaseURL + branchPath, nil
	case WEBPAGEPULLREQUEST:
		return webBaseURL + pullRequestPath, nil
	default:
		return "", fmt.Errorf("unsupported web page type: %s", pageType)
	}
}

// ------------------------------------
//
//	Turn a remote URL into the web base URL of the repo ("https://host/owner/repo") and return it
//	with the lowercase host name. Accepted forms:
//	- scp-like "[user@]host:owner/repo(.git)"
//	- "ssh://[user@]host[:port]/owner/repo(.git)"; the port is the ssh port, so it is removed
//	- "git://host[:port]/owner/repo(.git)"; the port is removed
//	- "http(s)://[user@]host[:port]/owner/repo(.git)"; the scheme and the port are kept
//	The host of an ssh form goes through resolveSSHHost. The full path is kept, so GitLab subgroups work.
//
// ------------------------------------
func parseRemoteWebBaseURL(remoteURL string) (host string, webBaseURL string, err error) {
	// the errors never show remoteURL: it can hold a password or a token ("https://user:token@host/...")
	scheme := "https"
	var hostWithPort, repoPath string
	if strings.Contains(remoteURL, "://") {
		parsedURL, parseErr := url.Parse(remoteURL)
		if parseErr != nil {
			return "", "", ErrUnsupportedRemoteHost
		}
		host = parsedURL.Hostname()
		switch parsedURL.Scheme {
		case "http", "https":
			scheme = parsedURL.Scheme
			hostWithPort = parsedURL.Host
		case "ssh", "git+ssh", "ssh+git":
			host = resolveSSHHost(host)
			hostWithPort = host
		default:
			hostWithPort = host
		}
		repoPath = parsedURL.Path
	} else {
		// like git, a "/" before the first ":" makes it a local path, not the scp-like form
		userAndHost, path, found := strings.Cut(remoteURL, ":")
		if !found || strings.Contains(userAndHost, "/") {
			return "", "", ErrUnsupportedRemoteHost
		}
		host = resolveSSHHost(userAndHost[strings.LastIndex(userAndHost, "@")+1:])
		hostWithPort = host
		repoPath = path
	}

	repoPath = strings.TrimSuffix(strings.Trim(repoPath, "/"), ".git")
	if host == "" || repoPath == "" {
		return "", "", fmt.Errorf("%w: %s", ErrUnsupportedRemoteHost, host)
	}
	return strings.ToLower(host), scheme + "://" + hostWithPort + "/" + repoPath, nil
}

// ------------------------------------
//
//	Return the real host of an ssh remote. The host can be an alias from ~/.ssh/config ("Host
//	github-work" with "HostName github.com") or an ssh endpoint on port 443 (ssh.github.com,
//	altssh.gitlab.com). Return the host unchanged when ssh is not installed or fails.
//
// ------------------------------------
func resolveSSHHost(host string) string {
	// `ssh -G` prints the config for the host and exits without connecting;
	// "--" so a host that starts with "-" is never read as an option
	sshConfigContext, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	sshConfigOutput, sshErr := exec.CommandContext(sshConfigContext, "ssh", "-G", "--", host).Output()
	if sshErr == nil {
		for line := range strings.SplitSeq(string(sshConfigOutput), "\n") {
			if hostName, found := strings.CutPrefix(strings.TrimSpace(line), "hostname "); found && hostName != "" {
				host = hostName
				break
			}
		}
	}
	if webHost, isPort443Host := sshPort443HostToWebHost[strings.ToLower(host)]; isPort443Host {
		return webHost
	}
	return host
}
