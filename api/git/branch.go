package git

import (
	"fmt"
	"os/exec"
	"strings"
)

var GITBRANCH *GitBranch

type GitBranch struct {
	RepoPath        string
	CurrentCheckOut string
	AllBranches     []string
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
	allBranches := []string{}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gb.RepoPath
	gitOutput, err := cmd.Output()
	if err != nil {
		gb.ErrorLog = append(gb.ErrorLog, fmt.Errorf("[GIT BRANCHES ERROR]: %w", err))
	}

	gitBranches := strings.Split(strings.TrimSpace(string(gitOutput)), "\n")
	gb.AllBranches = make([]string, 0, max(0, len(gitBranches)-1))
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
		gb.CurrentCheckOut = gitBranches[0]
	}

	for _, branch := range gitBranches {
		branch = strings.TrimSpace(branch)

		if strings.HasPrefix(branch, "*") {
			branch = strings.TrimSpace(strings.TrimPrefix(branch, "*"))
			gb.CurrentCheckOut = branch
		} else {
			allBranches = append(allBranches, branch)
		}
	}

	gb.AllBranches = allBranches
}
