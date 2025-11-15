package git

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"

	"gitti/executor"
)

type GitPull struct {
	errorLog            []error
	gitPullOutput       []string
	gitPullProcessMutex sync.Mutex
	gitPullProcessCmd   *exec.Cmd
	gitProcessLock      *GitProcessLock
	updateChannel       chan string
}

func InitGitPull(updateChannel chan string, gitProcessLock *GitProcessLock) *GitPull {
	gitPull := &GitPull{
		errorLog:       []error{},
		gitPullOutput:  []string{},
		updateChannel:  updateChannel,
		gitProcessLock: gitProcessLock,
	}
	gitPull.gitPullProcessCmd = nil

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
func (gp *GitPull) GitPull(pullType string) int {
	if !gp.gitProcessLock.CanProceedWithGitOps() {
		return -1
	}
	defer func() {
		gp.gitPullProcessMutex.Lock()
		gp.gitPullProcessReset()
		gp.gitPullProcessMutex.Unlock()
	}()

	gp.gitPullProcessMutex.Lock()
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

	cmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitPullArgs, true)

	gp.gitPullProcessCmd = cmdExecutor

	// Combine stderr into stdout
	stdout, err := cmdExecutor.StdoutPipe()
	if err != nil {
		gp.gitPullProcessMutex.Unlock()
		gp.errorLog = append(gp.errorLog, fmt.Errorf("[PIPE ERROR]: %w", err))
		return -1
	}
	cmdExecutor.Stderr = cmdExecutor.Stdout

	// Start the process while still holding the mutex
	if err := cmdExecutor.Start(); err != nil {
		gp.gitPullProcessMutex.Unlock()
		gp.errorLog = append(gp.errorLog, fmt.Errorf("[START ERROR]: %w", err))
		return -1
	}

	// Process is now started and can be killed safely
	gp.gitPullProcessMutex.Unlock()

	// Stream combined output
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		scanner.Split(splitOnCarriageReturnOrNewline)
		cursorIndex := 0
		for scanner.Scan() {
			updatedCursorIndex, updatedGitPullOutput := handleProgressOutputStream(cursorIndex, scanner, gp.gitPullOutput)
			gp.gitPullOutput = updatedGitPullOutput
			cursorIndex = updatedCursorIndex
			gp.updateChannel <- GIT_REMOTE_PUSH_OUTPUT_UPDATE
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

// --------------------------------
//
// # Kill the git pull process
//
// --------------------------------
func (gp *GitPull) KillGitPullCmd() {
	gp.gitPullProcessMutex.Lock()
	defer gp.gitPullProcessMutex.Unlock()

	if gp.gitPullProcessCmd != nil && gp.gitPullProcessCmd.Process != nil {
		_ = gp.gitPullProcessCmd.Process.Kill()
	}
}

// --------------------------------
//
// # Reset the git pull process after done
//
// --------------------------------
func (gp *GitPull) gitPullProcessReset() {
	gp.gitPullProcessCmd = nil
	gp.gitProcessLock.ReleaseGitOpsLock()
}
