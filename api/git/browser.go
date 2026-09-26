package git

import (
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	"github.com/gohyuhan/gitti/executor"
)

// the reasons a web page URL cannot be built
var (
	ErrNoRemote              = errors.New("no remote")
	ErrUnsupportedRemoteHost = errors.New("unsupported remote host (only github, gitlab and bitbucket)")
	ErrBranchNotPushed       = errors.New("branch has no upstream, push it first")
	ErrCommitNotOnRemote     = errors.New("commit is not on the remote, push it first")
)

// the ssh endpoints on port 443, whose web pages are on the main host
var sshPort443HostToWebHost = map[string]string{
	"ssh.github.com":    "github.com",
	"altssh.gitlab.com": "gitlab.com",
}

// ------------------------------------
//
//	Return the web page URL of a commit. The remote is the upstream remote of the checked-out
//	branch, else "origin". The commit must be on a branch of that remote, or the page does not exist.
//
// ------------------------------------
func GetCommitWebPageURL(commitHash string, checkedOutBranchName string) (string, error) {
	remoteName, _, hasUpstream := getBranchUpstream(checkedOutBranchName)
	if !hasUpstream {
		remoteName = "origin"
	}

	webPageURL, err := getRemoteWebPageURL(remoteName, WEBPAGECOMMIT, commitHash)
	if err != nil {
		return "", err
	}

	gitArgs := []string{"branch", "-r", "--contains", commitHash, "--list", remoteName + "/*"}
	remoteBranchesOutput, remoteBranchesErr := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false).Output()
	if remoteBranchesErr != nil || strings.TrimSpace(string(remoteBranchesOutput)) == "" {
		return "", fmt.Errorf("%w: %s", ErrCommitNotOnRemote, remoteName)
	}
	return webPageURL, nil
}

// ------------------------------------
//
//	Return the web page URL of a local branch (WEBPAGEBRANCH) or of a new pull request for it
//	(WEBPAGEPULLREQUEST). The URL uses the remote and the branch name of the upstream, so a
//	local branch "new" that tracks "origin/old" opens "old".
//
// ------------------------------------
func GetBranchWebPageURL(branchName string, pageType string) (string, error) {
	remoteName, remoteBranchName, hasUpstream := getBranchUpstream(branchName)
	if !hasUpstream {
		return "", fmt.Errorf("%w: %s", ErrBranchNotPushed, branchName)
	}
	return getRemoteWebPageURL(remoteName, pageType, remoteBranchName)
}

// ------------------------------------
//
//	Return the remote and the remote-side branch name that a local branch tracks. hasUpstream
//	is false when the branch has no upstream, tracks a local branch, or its remote branch is gone.
//
// ------------------------------------
func getBranchUpstream(branchName string) (remoteName string, remoteBranchName string, hasUpstream bool) {
	// @{upstream} fails when the branch has no upstream (branch.<name>.remote or .merge not set)
	// and when the remote-tracking branch is gone (deleted on the remote and pruned)
	gitArgs := []string{"rev-parse", "--verify", "--quiet", branchName + "@{upstream}"}
	if executor.GittiCmdExecutor.RunGitCmd(gitArgs, false).Run() != nil {
		return "", "", false
	}

	// read the remote and the branch as separate config values instead of splitting "origin/x",
	// because a remote name can hold "/"
	remoteName = getGitConfigValue("branch." + branchName + ".remote")
	// "." is the remote of a branch that tracks another local branch
	if remoteName == "." {
		return "", "", false
	}
	mergeRef := getGitConfigValue("branch." + branchName + ".merge")
	return remoteName, strings.TrimPrefix(mergeRef, "refs/heads/"), true
}

// ------------------------------------
//
//	Return the value of a git config key, or "" when the key is not set
//
// ------------------------------------
func getGitConfigValue(key string) string {
	gitArgs := []string{"config", "--get", key}
	configOutput, configErr := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false).Output()
	if configErr != nil {
		return ""
	}
	return strings.TrimSpace(string(configOutput))
}

// ------------------------------------
//
//	Read the URL of a remote and build the web page URL from it
//
// ------------------------------------
func getRemoteWebPageURL(remoteName string, pageType string, commitHashOrBranchName string) (string, error) {
	// read the URL fresh: the cached remote list has no order, and its push URL can differ from the fetch URL
	gitArgs := []string{"remote", "get-url", remoteName}
	remoteURLOutput, remoteURLErr := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false).Output()
	if remoteURLErr != nil {
		return "", fmt.Errorf("%w: %s", ErrNoRemote, remoteName)
	}
	return buildWebPageURL(strings.TrimSpace(string(remoteURLOutput)), pageType, commitHashOrBranchName)
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
	for index, segment := range branchNamePathSegments {
		branchNamePathSegments[index] = url.PathEscape(segment)
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
		pullRequestPath = "/compare/" + branchNameAsPathSegment + "?expand=1"
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
	default:
		return webBaseURL + pullRequestPath, nil
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
	sshConfigOutput, sshErr := exec.Command("ssh", "-G", "--", host).Output()
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
