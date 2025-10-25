package git

import (
	"fmt"
	"os/exec"
	"strings"
)

var GITBRANCH *GitBranch

type BranchInfo struct {
	BranchName   string
	IsCheckedOut bool
}

type GitBranch struct {
	RepoPath        string
	CurrentCheckOut BranchInfo
	AllBranches     []BranchInfo
	ErrorLog        []error
}

func InitGitBranch(repoPath string) {
	gitBranch := GitBranch{
		RepoPath: repoPath,
	}
	GITBRANCH = &gitBranch
}

func (gb *GitBranch) GetLatestBranchesinfo() {
	gitArgs := []string{"branch"}
	allBranches := []BranchInfo{}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gb.RepoPath
	gitOutput, err := cmd.Output()
	if err != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT BRANCHES ERROR]: %w", err))
	}

	gitBranches := strings.Split(strings.TrimSpace(string(gitOutput)), "\n")
	gb.AllBranches = make([]BranchInfo, 0, max(0, len(gitBranches)-1))
	// meaning this was a newly init repo with a uncommited branch
	if len(gitBranches) < 1 {
		gitArgs := []string{"symbolic-ref", "--short", "HEAD"}
		cmd := exec.Command("git", gitArgs...)
		cmd.Dir = gb.RepoPath
		gitOutput, err := cmd.Output()
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
