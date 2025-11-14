package git

import (
	"fmt"
	"strings"

	"gitti/cmd"
)

type BranchInfo struct {
	BranchName   string
	IsCheckedOut bool
}

type GitBranch struct {
	isRepoUnborn    bool // meaning this is a newly init repo, no commit on any branch yet
	currentCheckOut BranchInfo
	allBranches     []BranchInfo
	errorLog        []error
	gitProcessLock  *GitProcessLock
}

func InitGitBranch(gitProcessLock *GitProcessLock) *GitBranch {
	gitBranch := GitBranch{
		isRepoUnborn:   false,
		gitProcessLock: gitProcessLock,
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
func (gb *GitBranch) GetLatestBranchesinfo() {
	gitArgs := []string{"branch"}
	allBranches := []BranchInfo{}

	bCmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
	gitOutput, err := bCmd.Output()
	if err != nil {
		gb.errorLog = append(gb.errorLog, fmt.Errorf("[GIT BRANCHES ERROR]: %w", err))
	}

	gitBranches := strings.Split(strings.TrimSpace(string(gitOutput)), "\n")

	gb.allBranches = make([]BranchInfo, 0, max(0, len(gitBranches)-1))
	// meaning this was a newly init repo with a uncommited branch
	if len(gitBranches) < 1 {
		gitArgs := []string{"symbolic-ref", "--short", "HEAD"}
		bCmd = cmd.GittiCmd.RunGitCmd(gitArgs, false)
		gitOutput, err := bCmd.Output()
		if err != nil {
			gb.errorLog = append(gb.errorLog, fmt.Errorf("[GIT BRANCHES ERROR]: %w", err))
		}
		gitBranches = strings.Split(strings.TrimSpace(string(gitOutput)), "\n")
		gb.currentCheckOut = BranchInfo{
			BranchName:   gitBranches[0],
			IsCheckedOut: true,
		}
		gb.isRepoUnborn = true
	}

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

	gb.allBranches = allBranches
}

// ----------------------------------
//
//	Set The Global Default Branch Name when git init
//
// ----------------------------------
func SetGitInitDefaultBranch(branchName string, cwd string) {
	gitArgs := []string{"config", "--global", "init.defaultBranch", branchName}

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
	_ = cmd.Run()
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

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
	_, err := cmd.CombinedOutput()
	if err != nil {
		gb.errorLog = append(gb.errorLog, fmt.Errorf("[GIT CREATE BRANCH ERROR]: %w", err))
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

	stashChangesGitArgs := []string{"stash", "push", "--all"}
	stashChangesCmd := cmd.GittiCmd.RunGitCmd(stashChangesGitArgs, false)
	_, stashChangesErr := stashChangesCmd.CombinedOutput()
	if stashChangesErr != nil {
		gb.errorLog = append(gb.errorLog, fmt.Errorf("[GIT STASH CHANGES ERROR]: %w", stashChangesErr))
		return
	}

	createBranchGitArgs := []string{"checkout", "-b", branchName}
	createBranchCmd := cmd.GittiCmd.RunGitCmd(createBranchGitArgs, false)
	_, createBranchErr := createBranchCmd.CombinedOutput()
	if createBranchErr != nil {
		gb.errorLog = append(gb.errorLog, fmt.Errorf("[GIT CREATE AND SWITCH BRANCH ERROR]: %w", createBranchErr))
		return
	}

	unstashChangesGitArgs := []string{"stash", "pop"}
	unstashChangesCmd := cmd.GittiCmd.RunGitCmd(unstashChangesGitArgs, false)
	_, unstashChangesErr := unstashChangesCmd.CombinedOutput()
	if unstashChangesErr != nil {
		gb.errorLog = append(gb.errorLog, fmt.Errorf("[GIT UNSTASH CHANGE ERROR]: %w", unstashChangesErr))
		return
	}
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

	stashChangesGitArgs := []string{"stash", "push", "--all"}
	stashChangesCmd := cmd.GittiCmd.RunGitCmd(stashChangesGitArgs, false)
	stashChangesOutput, stashChangesErr := stashChangesCmd.CombinedOutput()
	gitOpsOutput = append(gitOpsOutput, processGeneralGitOpsOutputIntoStringArray(stashChangesOutput)...)
	if stashChangesErr != nil {
		gb.errorLog = append(gb.errorLog, fmt.Errorf("[GIT STASH CHANGES ERROR]: %w", stashChangesErr))
		return gitOpsOutput, false
	}

	switchBranchGitArgs := []string{"checkout", "-B", branchName}
	switchBranchCmd := cmd.GittiCmd.RunGitCmd(switchBranchGitArgs, false)
	switchBranchOutput, switchBranchErr := switchBranchCmd.CombinedOutput()
	gitOpsOutput = append(gitOpsOutput, processGeneralGitOpsOutputIntoStringArray(switchBranchOutput)...)
	if switchBranchErr != nil {
		gb.errorLog = append(gb.errorLog, fmt.Errorf("[GIT SWITCH BRANCH ERROR]: %w", switchBranchErr))
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

	stashChangesGitArgs := []string{"stash", "push", "--all"}
	stashChangesCmd := cmd.GittiCmd.RunGitCmd(stashChangesGitArgs, false)
	stashChangesOutput, stashChangesErr := stashChangesCmd.CombinedOutput()
	gitOpsOutput = append(gitOpsOutput, processGeneralGitOpsOutputIntoStringArray(stashChangesOutput)...)
	if stashChangesErr != nil {
		gb.errorLog = append(gb.errorLog, fmt.Errorf("[GIT STASH CHANGES ERROR]: %w", stashChangesErr))
		return gitOpsOutput, false
	}

	switchBranchGitArgs := []string{"checkout", "-B", branchName}
	switchBranchCmd := cmd.GittiCmd.RunGitCmd(switchBranchGitArgs, false)
	switchBranchOutput, switchBranchErr := switchBranchCmd.CombinedOutput()
	gitOpsOutput = append(gitOpsOutput, processGeneralGitOpsOutputIntoStringArray(switchBranchOutput)...)
	if switchBranchErr != nil {
		gb.errorLog = append(gb.errorLog, fmt.Errorf("[GIT SWITCH BRANCH ERROR]: %w", switchBranchErr))
		return gitOpsOutput, false
	}

	unstashChangesGitArgs := []string{"stash", "pop"}
	unstashChangesCmd := cmd.GittiCmd.RunGitCmd(unstashChangesGitArgs, false)
	unstashChangesOutput, unstashChangesErr := unstashChangesCmd.CombinedOutput()
	gitOpsOutput = append(gitOpsOutput, processGeneralGitOpsOutputIntoStringArray(unstashChangesOutput)...)
	if unstashChangesErr != nil {
		gb.errorLog = append(gb.errorLog, fmt.Errorf("[GIT UNSTASH CHANGE ERROR]: %w", unstashChangesErr))
		return gitOpsOutput, false
	}

	return gitOpsOutput, true
}

//test
