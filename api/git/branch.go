package git

import (
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
	allBranches     []BranchInfo
	logging         *logging.GittiLogging
	gitProcessLock  *GitProcessLock
}

func InitGitBranch(gitProcessLock *GitProcessLock, logging *logging.GittiLogging) *GitBranch {
	gitBranch := GitBranch{
		isRepoUnborn:   false,
		gitProcessLock: gitProcessLock,
		logging:        logging,
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
	if len(gitBranches) == 1 && gitBranches[0] == "" {
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
	createBranchBasedOnRemoteOutput, createBranchBasedOnRemoteErr := createBranchBasedOnRemoteCmdExecutor.CombinedOutput()
	gb.logging.RegisterNewLog(logging.CREATE_NEW_BRANCH_BASED_ON_REMOTE_BRANCH_OPS, strings.Join(createBranchBasedOnRemoteGitArgs, " "), logging.INFO, "", true)

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
