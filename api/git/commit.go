package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

var GITCOMMIT *GitCommit

type GitCommit struct {
	RepoPath         string
	ErrorLog         []error
	GitCommitProcess *exec.Cmd
	GitCommitOutput  []string
	UpdateChannel    chan string
}

func InitGitCommit(repoPath string, updateChannel chan string) {
	gitCommit := GitCommit{
		RepoPath:         repoPath,
		GitCommitProcess: nil,
		GitCommitOutput:  []string{},
		UpdateChannel:    updateChannel,
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

func (gc *GitCommit) GitCommit(message, description string) int {
	gitArgs := []string{"commit", "-m", message}
	if len(description) > 0 {
		gitArgs = append(gitArgs, "-m", description)
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gc.RepoPath

	// Combine stderr into stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[PIPE ERROR]: %w", err))
		return -1
	}
	cmd.Stderr = cmd.Stdout

	gc.GitCommitProcess = cmd
	defer func() {
		// ensure cleanup even if Start or Wait fails
		gc.GitCommitProcess = nil
	}()

	if err := cmd.Start(); err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[START ERROR]: %w", err))
		return -1
	}

	// Stream combined output
	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Text()
			gc.GitCommitOutput = append(gc.GitCommitOutput, line)
			gc.UpdateChannel <- GIT_COMMIT_OUTPUT_UPDATE
		}
	}()

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", err))
			return status
		}
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[UNEXPECTED ERROR]: %w", err))
		return -1
	}

	return 0
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

func (gc *GitCommit) ClearGitCommitOutput() {
	gc.GitCommitOutput = []string{}
}

func (gc *GitCommit) KillCommit() {
	if gc.GitCommitProcess != nil && gc.GitCommitProcess.Process != nil {
		_ = gc.GitCommitProcess.Process.Kill()
		gc.GitCommitProcess = nil
	}
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
