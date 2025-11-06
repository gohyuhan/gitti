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
	CurrentCheckOut BranchInfo
	AllBranches     []BranchInfo
	ErrorLog        []error
}

func InitGitBranch() {
	gitBranch := GitBranch{}
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
	_, err := cmd.Output()
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
	_, err := cmd.Output()
	if err != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT UNSTASH ERROR]: %w", err))
	}
}
