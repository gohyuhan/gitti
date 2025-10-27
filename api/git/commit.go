package git

import (
	"fmt"
	"os/exec"
)

var GITCOMMIT *GitCommit

type GitCommit struct {
	RepoPath string
	ErrorLog []error
}

func InitGitCommit(repoPath string) {
	gitCommit := GitCommit{
		RepoPath: repoPath,
	}
	GITCOMMIT = &gitCommit
}

func (gc *GitCommit) GitInit() {
	gitArgs := []string{"init"}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gc.RepoPath
	_, err := cmd.Output()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", err))
	}
}

func (gc *GitCommit) GitFetch() {
	gitArgs := []string{"fetch"}
	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gc.RepoPath
	_, err := cmd.Output()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", err))
	}
}

func (gc *GitCommit) GitStage() {

}

func (gc *GitCommit) GitCommit() {

}

func (gc *GitCommit) GitPull() {

}

func (gc *GitCommit) GitPush() {
	// gitArgs := []string{"pull`"}
	// cmd := exec.Command("git", gitArgs...)
	// cmd.Dir = gc.RepoPath
	// gitOutput, err := cmd.Output()
	// if err != nil {
	// 	gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", err))
	// }
}
