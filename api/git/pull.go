package git

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"sync"

	"github.com/gohyuhan/gitti/executor"
)

type GitPull struct {
	errorLog       []error
	gitPullOutput  []string
	gitProcessLock *GitProcessLock
	updateChannel  chan string
}

func InitGitPull(updateChannel chan string, gitProcessLock *GitProcessLock) *GitPull {
	gitPull := &GitPull{
		errorLog:       []error{},
		gitPullOutput:  []string{},
		updateChannel:  updateChannel,
		gitProcessLock: gitProcessLock,
	}

	return gitPull
}

// --------------------------------
//
// return the git pull output
//
// --------------------------------
func (gp *GitPull) GetGitPullOutput() []string {
	return gp.gitPullOutput
}

// --------------------------------
//
// # Git Pull and will operate differently based on the user selection type
//
// --------------------------------
func (gp *GitPull) GitPull(ctx context.Context, pullType string) int {
	if !gp.gitProcessLock.CanProceedWithGitOps() {
		return -1
	}
	defer func() {
		gp.gitProcessLock.ReleaseGitOpsLock()
	}()

	gp.ClearGitPullOutput()
	var gitPullArgs []string
	switch pullType {
	case GITPULL:
		gitPullArgs = []string{"pull", "--progress", "--no-edit"}
	case GITPULLREBASE:
		gitPullArgs = []string{"pull", "--progress", "--rebase", "--autostash", "--no-edit"}
	case GITPULLMERGE:
		gitPullArgs = []string{"pull", "--progress", "--no-rebase", "--no-edit"}
	}

	cmdExecutor := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitPullArgs, true)

	// Combine stderr into stdout
	stdout, err := cmdExecutor.StdoutPipe()
	if err != nil {
		gp.errorLog = append(gp.errorLog, fmt.Errorf("[PIPE ERROR]: %w", err))
		return -1
	}
	cmdExecutor.Stderr = cmdExecutor.Stdout

	// Start the process
	if err := cmdExecutor.Start(); err != nil {
		gp.errorLog = append(gp.errorLog, fmt.Errorf("[START ERROR]: %w", err))
		return -1
	}

	// Stream combined output
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		scanner.Split(splitOnCarriageReturnOrNewline)
		cursorIndex := 0
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return // Stop immediately on cancel
			default:
				updatedCursorIndex, updatedGitPullOutput := handleProgressOutputStream(cursorIndex, scanner, gp.gitPullOutput)
				gp.gitPullOutput = updatedGitPullOutput
				cursorIndex = updatedCursorIndex
				gp.updateChannel <- GIT_PULL_OUTPUT_UPDATE
			}
		}
	}()

	waitErr := cmdExecutor.Wait()
	wg.Wait()

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gp.errorLog = append(gp.errorLog, fmt.Errorf("[GIT PULL ERROR]: %w", waitErr))
			return status
		}
		gp.errorLog = append(gp.errorLog, fmt.Errorf("[UNEXPECTED ERROR]: %w", waitErr))
		return -1
	}
	return 0
}

// --------------------------------
//
// # Clear the Git Process Output
//
// --------------------------------
func (gp *GitPull) ClearGitPullOutput() {
	gp.gitPullOutput = []string{}
}
