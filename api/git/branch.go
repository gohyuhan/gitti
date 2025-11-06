package git

import (
	"fmt"
	"strings"

	"gitti/cmd"
)

var GITBRANCH *GitBranch

type BranchInfo struct {
	BranchName   string
	IsCheckedOut bool
}

type GitBranch struct {
	IsRepoUnborn    bool // meaning this is a newly init repo, no commit on any branch yet
	CurrentCheckOut BranchInfo
	AllBranches     []BranchInfo
	ErrorLog        []error
}

func InitGitBranch() {
	gitBranch := GitBranch{
		IsRepoUnborn: false,
	}
	GITBRANCH = &gitBranch
}

// ----------------------------------
//
//	Related to Branches Info
//
// ----------------------------------
func (gb *GitBranch) GetLatestBranchesinfo() {
	gitArgs := []string{"branch"}
	allBranches := []BranchInfo{}

	bCmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	gitOutput, err := bCmd.Output()
	if err != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT BRANCHES ERROR]: %w", err))
	}

	gitBranches := strings.Split(strings.TrimSpace(string(gitOutput)), "\n")
	gb.AllBranches = make([]BranchInfo, 0, max(0, len(gitBranches)-1))
	// meaning this was a newly init repo with a uncommited branch
	if len(gitBranches) < 1 {
		gitArgs := []string{"symbolic-ref", "--short", "HEAD"}
		bCmd = cmd.GittiCmd.RunGitCmd(gitArgs)
		gitOutput, err := bCmd.Output()
		if err != nil {
			gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT BRANCHES ERROR]: %w", err))
		}
		gitBranches = strings.Split(strings.TrimSpace(string(gitOutput)), "\n")
		gb.CurrentCheckOut = BranchInfo{
			BranchName:   gitBranches[0],
			IsCheckedOut: true,
		}
		gb.IsRepoUnborn = true
	}

	for _, branch := range gitBranches {
		branch = strings.TrimSpace(branch)

		if strings.HasPrefix(branch, "*") {
			branch = strings.TrimSpace(strings.TrimPrefix(branch, "*"))
			gb.CurrentCheckOut = BranchInfo{
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

	gb.AllBranches = allBranches
}

// ----------------------------------
//
//	Set The Global Default Branch Name when git init
//
// ----------------------------------
func SetGitInitDefaultBranch(branchName string, cwd string) {
	gitArgs := []string{"config", "--global", "init.defaultBranch", branchName}

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	_ = cmd.Run()
}

// ----------------------------------
//
//	Related to Git Stash
//
// ----------------------------------
func (gb *GitBranch) GitStash() {
	gitArgs := []string{"stash", "--all"}

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	_, err := cmd.CombinedOutput()
	if err != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT STASH ERROR]: %w", err))
	}
}

// ----------------------------------
//
//	Related to Git UnStash
//
// ----------------------------------
func (gb *GitBranch) GitUnstash() {
	gitArgs := []string{"stash", "pop"}

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	_, err := cmd.CombinedOutput()
	if err != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT UNSTASH ERROR]: %w", err))
	}
}

// ----------------------------------
//
//	Related to Create New Branch ( only create, remain at current branch )
//
// ----------------------------------
func (gb *GitBranch) GitCreateNewBranch(branchName string) {
	gitArgs := []string{"branch", branchName}

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	_, err := cmd.CombinedOutput()
	if err != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT CREATE BRANCH ERROR]: %w", err))
	}
}

// ----------------------------------
//
//	Related to Create New Branch and Move All Changes to new Branch ( create, then switch to new branch )
//
// ----------------------------------
func (gb *GitBranch) GitCreateNewBranchAndSwitch(branchName string) {
	stashChangesGitArgs := []string{"stash", "push", "--all"}
	stashChangesCmd := cmd.GittiCmd.RunGitCmd(stashChangesGitArgs)
	_, stashChangesErr := stashChangesCmd.CombinedOutput()
	if stashChangesErr != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT STASH CHANGES ERROR]: %w", stashChangesErr))
		return
	}

	createBranchGitArgs := []string{"switch", "-c", branchName}
	createBranchCmd := cmd.GittiCmd.RunGitCmd(createBranchGitArgs)
	_, createBranchErr := createBranchCmd.CombinedOutput()
	if createBranchErr != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT CREATE AND SWITCH BRANCH ERROR]: %w", createBranchErr))
		return
	}

	unstashChangesGitArgs := []string{"stash", "pop"}
	unstashChangesCmd := cmd.GittiCmd.RunGitCmd(unstashChangesGitArgs)
	_, unstashChangesErr := unstashChangesCmd.CombinedOutput()
	if unstashChangesErr != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT UNSTASH CHANGE ERROR]: %w", unstashChangesErr))
		return
	}
}

// ----------------------------------
//
//	Related to Switch Branch ( Does not bring the changes over )
//
// ----------------------------------
func (gb *GitBranch) GitSwitchBranch(branchName string) {
	stashChangesGitArgs := []string{"stash", "push", "--all"}
	stashChangesCmd := cmd.GittiCmd.RunGitCmd(stashChangesGitArgs)
	_, stashChangesErr := stashChangesCmd.CombinedOutput()
	if stashChangesErr != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT STASH CHANGES ERROR]: %w", stashChangesErr))
		return
	}

	switchBranchGitArgs := []string{"switch", branchName}
	switchBranchCmd := cmd.GittiCmd.RunGitCmd(switchBranchGitArgs)
	_, switchBranchErr := switchBranchCmd.CombinedOutput()
	if switchBranchErr != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT SWITCH BRANCH ERROR]: %w", switchBranchErr))
		return
	}
}

// ----------------------------------
//
//	Related to Switch Branch with the changes ( bring the changes over )
//
// ----------------------------------
func (gb *GitBranch) GitSwitchBranchWithChanges(branchName string) {
	stashChangesGitArgs := []string{"stash", "push", "--all"}
	stashChangesCmd := cmd.GittiCmd.RunGitCmd(stashChangesGitArgs)
	_, stashChangesErr := stashChangesCmd.CombinedOutput()
	if stashChangesErr != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT STASH CHANGES ERROR]: %w", stashChangesErr))
		return
	}

	switchBranchGitArgs := []string{"switch", branchName}
	switchBranchCmd := cmd.GittiCmd.RunGitCmd(switchBranchGitArgs)
	_, switchBranchErr := switchBranchCmd.CombinedOutput()
	if switchBranchErr != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT SWITCH BRANCH ERROR]: %w", switchBranchErr))
		return
	}

	unstashChangesGitArgs := []string{"stash", "pop"}
	unstashChangesCmd := cmd.GittiCmd.RunGitCmd(unstashChangesGitArgs)
	_, unstashChangesErr := unstashChangesCmd.CombinedOutput()
	if unstashChangesErr != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT UNSTASH CHANGE ERROR]: %w", unstashChangesErr))
		return
	}
}
