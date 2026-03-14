package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/logging"
)

type BranchInfo struct {
	BranchName   string
	IsCheckedOut bool
}

type GitBranch struct {
	isRepoUnborn    bool // meaning this is a newly init repo, no commit on any branch yet
	currentCheckOut BranchInfo
	allBranches     []BranchInfo // this refer to all local branch
	remoteBranches  []BranchInfo
	logging         *logging.GittiLogging
	gitProcessLock  *GitProcessLock
	FfMerge         bool // determine when merge it was fast forward or not, fast forward will not have the merge commit and non fast forward will have one
}

// ----------------------------------
//
//	Initialize the git branch handler with shared dependencies
//
// ----------------------------------
func InitGitBranch(gitProcessLock *GitProcessLock, ffMerge bool, logging *logging.GittiLogging) *GitBranch {
	gitBranch := GitBranch{
		isRepoUnborn:   false,
		gitProcessLock: gitProcessLock,
		logging:        logging,
		FfMerge:        ffMerge,
	}
	return &gitBranch
}

// ----------------------------------
//
//	Return current branch
//
// ----------------------------------
func (gb *GitBranch) CurrentCheckOut() BranchInfo {
	return gb.currentCheckOut
}

// ----------------------------------
//
//	Return  allbranch
//
// ----------------------------------
func (gb *GitBranch) AllBranches() []BranchInfo {
	copied := make([]BranchInfo, len(gb.allBranches))
	copy(copied, gb.allBranches)
	return copied
}

// ----------------------------------
//
//	Return  remoteBranches
//
// ----------------------------------
func (gb *GitBranch) RemoteBranches() []BranchInfo {
	copied := make([]BranchInfo, len(gb.remoteBranches))
	copy(copied, gb.remoteBranches)
	return copied
}

// ----------------------------------
//
//	Return is repo unborn
//
// ----------------------------------
func (gb *GitBranch) IsRepoUnborn() bool {
	return gb.isRepoUnborn
}

// ----------------------------------
//
//		Retrieve Branches Info
//	 * Passive, this should only be trigger by system
//
// ----------------------------------
func (gb *GitBranch) GetLatestBranchesInfo() {
	gitArgs := []string{"branch"}
	allBranches := []BranchInfo{}

	gb.isRepoUnborn = false

	branchCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	gitOutput, err := branchCmdExecutor.Output()
	if err != nil {
		gb.logging.RegisterNewLog(logging.GET_LATEST_BRANCH_INFO_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.GET_LATEST_BRANCH_INFO_OPS, err.Error()), true)
		return
	}

	gitBranches := processGeneralGitOpsOutputIntoStringArray(gitOutput)

	gb.allBranches = make([]BranchInfo, 0, max(0, len(gitBranches)-1))
	// meaning this was a newly init repo with a uncommited branch
	if len(gitBranches) < 1 {
		gitArgs := []string{"symbolic-ref", "--short", "HEAD"}
		branchCmdExecutor = executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
		gitOutput, err := branchCmdExecutor.Output()
		if err != nil {
			gb.logging.RegisterNewLog(logging.GET_LATEST_BRANCH_INFO_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.GET_LATEST_BRANCH_INFO_OPS, err.Error()), true)
			return
		}
		gitBranches = processGeneralGitOpsOutputIntoStringArray(gitOutput)
		gb.currentCheckOut = BranchInfo{
			BranchName:   gitBranches[0],
			IsCheckedOut: true,
		}
		gb.isRepoUnborn = true
	} else {
		for _, branch := range gitBranches {
			branch = strings.TrimSpace(branch)

			if strings.HasPrefix(branch, "*") {
				branch = strings.TrimSpace(strings.TrimPrefix(branch, "*"))
				gb.currentCheckOut = BranchInfo{
					BranchName:   branch,
					IsCheckedOut: true,
				}
			} else {
				allBranches = append(allBranches, BranchInfo{
					BranchName:   branch,
					IsCheckedOut: false,
				})
			}
		}
	}

	gb.allBranches = allBranches
}

// ----------------------------------
//
//	Set The Global Default Branch Name when git init
//
// ----------------------------------
func SetGitInitDefaultBranch(branchName string, cwd string) {
	gitArgs := []string{"config", "--global", "init.defaultBranch", branchName}

	cmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	_ = cmdExecutor.Run()
}

// ----------------------------------
//
//	Related to Create New Branch ( only create, remain at current branch )
//
// ----------------------------------
func (gb *GitBranch) GitCreateNewBranch(branchName string) {
	if !gb.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gb.gitProcessLock.ReleaseGitOpsLock()

	gitArgs := []string{"branch", branchName}

	if gb.isRepoUnborn {
		gitArgs = []string{"branch", "-M", branchName}
	}

	cmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	err := cmdExecutor.Run()
	gb.logging.RegisterNewLog(logging.CREATE_NEW_BRANCH_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	if err != nil {
		gb.logging.RegisterNewLog(logging.CREATE_NEW_BRANCH_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.CREATE_NEW_BRANCH_OPS, err.Error()), true)
		return
	}
}

// ----------------------------------
//
//	Related to Create New Branch and Move All Changes to new Branch ( create, then switch to new branch )
//
// ----------------------------------
func (gb *GitBranch) GitCreateNewBranchAndSwitch(branchName string) {
	if !gb.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gb.gitProcessLock.ReleaseGitOpsLock()

	createAndSwitchBranchGitArgs := []string{"checkout", "-b", branchName}
	createAndSwitchBranchCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(createAndSwitchBranchGitArgs, false)
	createAndSwitchBranchErr := createAndSwitchBranchCmdExecutor.Run()
	gb.logging.RegisterNewLog(logging.CREATE_NEW_BRANCH_AND_SWITCH_OPS, strings.Join(createAndSwitchBranchGitArgs, " "), logging.INFO, "", true)
	if createAndSwitchBranchErr != nil {
		gb.logging.RegisterNewLog(logging.CREATE_NEW_BRANCH_AND_SWITCH_OPS, strings.Join(createAndSwitchBranchGitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.CREATE_NEW_BRANCH_AND_SWITCH_OPS, createAndSwitchBranchErr.Error()), true)
		return
	}
}

// ----------------------------------
//
//	Related to Create New Branch based on a remote branch
//
// ----------------------------------
func (gb *GitBranch) GitCreateNewBranchBasedOnRemote(remoteName string, branchName string) ([]string, bool) {
	success := false
	if !gb.gitProcessLock.CanProceedWithGitOps() {
		return []string{gb.gitProcessLock.OtherProcessRunningWarning()}, success
	}
	defer gb.gitProcessLock.ReleaseGitOpsLock()

	gitFetch(gb.logging, false)

	remoteBranchName := fmt.Sprintf("%s/%s", remoteName, branchName)

	createBranchBasedOnRemoteGitArgs := []string{"branch", branchName, remoteBranchName}
	createBranchBasedOnRemoteCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(createBranchBasedOnRemoteGitArgs, false)
	gb.logging.RegisterNewLog(logging.CREATE_NEW_BRANCH_BASED_ON_REMOTE_BRANCH_OPS, strings.Join(createBranchBasedOnRemoteGitArgs, " "), logging.INFO, "", true)
	createBranchBasedOnRemoteOutput, createBranchBasedOnRemoteErr := createBranchBasedOnRemoteCmdExecutor.CombinedOutput()

	parsedCreateBranchBasedOnRemoteOutput := processGeneralGitOpsOutputIntoStringArray(createBranchBasedOnRemoteOutput)

	if createBranchBasedOnRemoteErr != nil {
		gb.logging.RegisterNewLog(logging.CREATE_NEW_BRANCH_BASED_ON_REMOTE_BRANCH_OPS, strings.Join(createBranchBasedOnRemoteGitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.CREATE_NEW_BRANCH_BASED_ON_REMOTE_BRANCH_OPS, createBranchBasedOnRemoteErr.Error()), true)
	} else {
		success = true
	}

	return parsedCreateBranchBasedOnRemoteOutput, success
}

// ----------------------------------
//
//	Related to Switch Branch ( Does not bring the changes over )
//
// ----------------------------------
func (gb *GitBranch) GitSwitchBranch(branchName string) ([]string, bool) {
	if !gb.gitProcessLock.CanProceedWithGitOps() {
		return []string{gb.gitProcessLock.OtherProcessRunningWarning()}, false
	}
	defer gb.gitProcessLock.ReleaseGitOpsLock()

	var gitOpsOutput []string

	gitArgs := []string{"stash", "push", "-u"}
	stashChangesCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	stashChangesOutput, stashChangesErr := stashChangesCmdExecutor.CombinedOutput()
	gb.logging.RegisterNewLog(logging.STASH_ALL_FILE_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	gitOpsOutput = append(gitOpsOutput, processGeneralGitOpsOutputIntoStringArray(stashChangesOutput)...)
	if stashChangesErr != nil {
		gb.logging.RegisterNewLog(logging.STASH_ALL_FILE_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.STASH_ALL_FILE_OPS, stashChangesErr.Error()), true)
		return gitOpsOutput, false
	}

	gitArgs = []string{"checkout", branchName}
	switchBranchCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	switchBranchOutput, switchBranchErr := switchBranchCmdExecutor.CombinedOutput()
	gb.logging.RegisterNewLog(logging.SWITCH_BRANCH_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	gitOpsOutput = append(gitOpsOutput, processGeneralGitOpsOutputIntoStringArray(switchBranchOutput)...)
	if switchBranchErr != nil {
		gb.logging.RegisterNewLog(logging.SWITCH_BRANCH_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.SWITCH_BRANCH_OPS, switchBranchErr.Error()), true)
		return gitOpsOutput, false
	}
	return gitOpsOutput, true
}

// ----------------------------------
//
//		Related to Create New Branch based on commit hash ( only create, remain at current branch )
//	 * used to create branch based on commit hash on reflog
//
// ----------------------------------
func (gb *GitBranch) GitCreateNewBranchBasedOnCommitHash(branchName string, commitHash string) {
	if !gb.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gb.gitProcessLock.ReleaseGitOpsLock()

	gitArgs := []string{"branch", branchName, commitHash}

	cmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	err := cmdExecutor.Run()
	gb.logging.RegisterNewLog(logging.CREATE_NEW_BRANCH_BASED_ON_COMMIT_HASH_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	if err != nil {
		gb.logging.RegisterNewLog(logging.CREATE_NEW_BRANCH_BASED_ON_COMMIT_HASH_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.CREATE_NEW_BRANCH_BASED_ON_COMMIT_HASH_OPS, err.Error()), true)
		return
	}
}

// ----------------------------------
//
//	Related to Switch Branch with the changes ( bring the changes over )
//
// ----------------------------------
func (gb *GitBranch) GitSwitchBranchWithChanges(branchName string) ([]string, bool) {
	if !gb.gitProcessLock.CanProceedWithGitOps() {
		return []string{gb.gitProcessLock.OtherProcessRunningWarning()}, false
	}
	defer gb.gitProcessLock.ReleaseGitOpsLock()
	var gitOpsOutput []string

	switchBranchGitArgs := []string{"checkout", branchName}
	switchBranchCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(switchBranchGitArgs, false)
	switchBranchOutput, switchBranchErr := switchBranchCmdExecutor.CombinedOutput()
	gb.logging.RegisterNewLog(logging.SWITCH_BRANCH_OPS, strings.Join(switchBranchGitArgs, " "), logging.INFO, "", true)

	gitOpsOutput = append(gitOpsOutput, processGeneralGitOpsOutputIntoStringArray(switchBranchOutput)...)

	if switchBranchErr != nil {
		gb.logging.RegisterNewLog(logging.SWITCH_BRANCH_OPS, strings.Join(switchBranchGitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.SWITCH_BRANCH_OPS, switchBranchErr.Error()), true)
		return gitOpsOutput, false
	}
	return gitOpsOutput, true
}

// ----------------------------------
//
//	Related to delete branch in local
//
// ----------------------------------
func (gb *GitBranch) DeleteLocalBranch(branchName string) ([]string, bool) {
	if !gb.gitProcessLock.CanProceedWithGitOps() {
		return []string{gb.gitProcessLock.OtherProcessRunningWarning()}, false
	}
	defer gb.gitProcessLock.ReleaseGitOpsLock()
	gitArgs := []string{"branch", "-D", branchName}
	branchDeleteExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	branchDeleteOutput, branchDeleteErr := branchDeleteExecutor.CombinedOutput()
	gb.logging.RegisterNewLog(logging.DELETE_LOCAL_BRANCH_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)

	gitOpsOutput := processGeneralGitOpsOutputIntoStringArray(branchDeleteOutput)

	if branchDeleteErr != nil {
		gb.logging.RegisterNewLog(logging.DELETE_LOCAL_BRANCH_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.DELETE_LOCAL_BRANCH_OPS, branchDeleteErr.Error()), true)
		return gitOpsOutput, false
	}
	return gitOpsOutput, true
}

// ----------------------------------
//
//		Related to get remote branch
//	 * this run passively and will not be triggered by user manually, this will be trigger after passive and manual git fetch
//
// ----------------------------------
func (gb *GitBranch) GetLatestRemoteBranchesInfo() {
	gitArgs := []string{"branch", "-r"}
	remoteBranchExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	remoteBranchOutput, remoteBranchErr := remoteBranchExecutor.CombinedOutput()

	parsedRemoteBranchOutput := processGeneralGitOpsOutputIntoStringArray(remoteBranchOutput)

	var remoteBranches []BranchInfo
	if remoteBranchErr != nil {
		gb.logging.RegisterNewLog(logging.RETRIEVE_LATEST_REMOTE_BRANCH_INFO, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.RETRIEVE_LATEST_REMOTE_BRANCH_INFO, remoteBranchErr.Error()), true)
		return
	}

	for _, parsedRemote := range parsedRemoteBranchOutput {
		if strings.Contains(parsedRemote, "/HEAD") {
			continue
		}
		remoteBranch := BranchInfo{
			BranchName:   strings.TrimSpace(parsedRemote),
			IsCheckedOut: false,
		}

		remoteBranches = append(remoteBranches, remoteBranch)
	}

	gb.remoteBranches = remoteBranches
}

// ----------------------------------
//
//	Related to get merge branch
//
// ----------------------------------
func (gb *GitBranch) GitMerge(ctx context.Context, branchesName []string) ([]string, bool) {
	if !gb.gitProcessLock.CanProceedWithGitOps() {
		return []string{gb.gitProcessLock.OtherProcessRunningWarning()}, false
	}
	defer gb.gitProcessLock.ReleaseGitOpsLock()
	var gitArgs []string
	if gb.FfMerge {
		gitArgs = []string{"merge", "--ff"}
	} else {
		gitArgs = []string{"merge", "--no-ff"}
	}
	gitArgs = append(gitArgs, branchesName...)

	mergeExecutor := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, false)
	gb.logging.RegisterNewLog(logging.MERGE_BRANCH_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	mergeOutput, mergeErr := mergeExecutor.CombinedOutput()

	gitMergeOpsOutput := processGeneralGitOpsOutputIntoStringArray(mergeOutput)

	if mergeErr != nil {
		gb.logging.RegisterNewLog(logging.MERGE_BRANCH_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.MERGE_BRANCH_OPS, mergeErr.Error()), true)
		return gitMergeOpsOutput, false
	}
	return gitMergeOpsOutput, true
}

// ----------------------------------
//
//	GitMergeWithSigning constructs a git merge command for terminal execution when signing is required.
//	When commit signing is enabled, gitti UI is suspended and the commit is executed directly in the terminal,
//	allowing the user to interact with the signing prompt (e.g., GPG passphrase).
//
// ----------------------------------
func (gb *GitBranch) GitMergeWithSigning(branchesName []string) []string {
	var gitArgs []string
	if gb.FfMerge {
		gitArgs = []string{"merge", "--ff"}
	} else {
		gitArgs = []string{"merge", "--no-ff"}
	}
	gitArgs = append(gitArgs, branchesName...)

	return gitArgs
}
