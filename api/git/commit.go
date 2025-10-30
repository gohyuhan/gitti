package git

import (
	"fmt"
	"os"
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
	gitArgs := []string{"add"}
	for _, files := range GITFILES.FilesStatus {
		if files.SelectedForStage {
			gitArgs = append(gitArgs, files.FileName)
		}
	}
	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gc.RepoPath
	_, err := cmd.Output()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", err))
	}
}

func (gc *GitCommit) GitCommit(message string, description string) {
	// need revision back later
	gitArgs := []string{"commit", "-m", message}
	if len(description) > 1 {
		gitArgs = append(gitArgs, []string{"-m", description}...)
	}
	for _, files := range GITFILES.FilesStatus {
		if files.SelectedForStage {
			gitArgs = append(gitArgs, files.FileName)
		}
	}
	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gc.RepoPath
	_, err := cmd.Output()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", err))
	}

	// run a git reset after a git commit is run, this is to unstage any changes made by pre-commit ( might not be applicable to most user but we still need to handle this )
	// so that those changes will be reflected on the modifed files panels and user can see the content of the modification
	// will no be any effcet if it was commited successfully
	cmd = exec.Command("git", []string{"reset"}...)
	_, err = cmd.Output()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", err))
	}
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

func GitInit(repoPath string) {
	gitArgs := []string{"init"}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = repoPath
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("[GIT INIT ERROR]: %v", err)
		os.Exit(1)
	}
}
