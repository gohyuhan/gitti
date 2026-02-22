package git

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/logging"
)

type GitPull struct {
	gitPullOutput   []string
	gitPullOutputMu sync.RWMutex
	gitProcessLock  *GitProcessLock
	updateChannel   chan string
	logging         *logging.GittiLogging
}

// ----------------------------------
//
//	Initialize the git pull handler with shared dependencies
//
// ----------------------------------
func InitGitPull(updateChannel chan string, gitProcessLock *GitProcessLock, logging *logging.GittiLogging) *GitPull {
	gitPull := &GitPull{
		gitPullOutput:  []string{},
		updateChannel:  updateChannel,
		gitProcessLock: gitProcessLock,
		logging:        logging,
	}

	return gitPull
}

// --------------------------------
//
// return the git pull output
//
// --------------------------------
func (gp *GitPull) GetGitPullOutput() []string {
	gp.gitPullOutputMu.RLock()
	defer gp.gitPullOutputMu.RUnlock()

	copied := make([]string, len(gp.gitPullOutput))
	copy(copied, gp.gitPullOutput)
	return copied
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
		gp.logging.RegisterNewLog(logging.PULL_OPS, "", logging.ERROR, fmt.Sprintf("[PIPE ERROR]: %s", err.Error()), false)
		return -1
	}
	cmdExecutor.Stderr = cmdExecutor.Stdout

	// Start the process
	if err := cmdExecutor.Start(); err != nil {
		gp.logging.RegisterNewLog(logging.PULL_OPS, "", logging.ERROR, fmt.Sprintf("[START ERROR]: %s", err.Error()), false)
		return -1
	} else {
		gp.logging.RegisterNewLog(logging.PULL_OPS, strings.Join(gitPullArgs, " "), logging.INFO, "", true)
	}

	// Stream combined output
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		scanner.Split(splitOnCarriageReturnOrNewline)
		cursorIndex := 0
		lastSent := time.Time{}
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return // Stop immediately on cancel
			default:
				gp.gitPullOutputMu.Lock()
				updatedCursorIndex, updatedGitPullOutput := handleProgressOutputStream(cursorIndex, scanner, gp.gitPullOutput)
				gp.gitPullOutput = updatedGitPullOutput
				cursorIndex = updatedCursorIndex
				gp.gitPullOutputMu.Unlock()
				if time.Since(lastSent) >= STREAMUPDATETHROTTLEMS*time.Millisecond {
					select {
					case gp.updateChannel <- GIT_PULL_OUTPUT_UPDATE:
						lastSent = time.Now()
					default:
					}
				}
			}
		}
		// trigger an update once it ends
		gp.updateChannel <- GIT_PULL_OUTPUT_UPDATE
	}()

	waitErr := cmdExecutor.Wait()
	wg.Wait()

	if ctx.Err() != nil {
		gp.logging.RegisterNewLog(logging.PULL_OPS, strings.Join(gitPullArgs, " "), logging.WARN, fmt.Sprintf("[%s CANCELLED]: %s", logging.PULL_OPS, ctx.Err().Error()), true)
		return -1
	}

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gp.logging.RegisterNewLog(logging.PULL_OPS, strings.Join(gitPullArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.PULL_OPS, waitErr.Error()), true)
			return status
		}
		gp.logging.RegisterNewLog(logging.PULL_OPS, strings.Join(gitPullArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.PULL_OPS, waitErr.Error()), true)
		return -1
	}
	return 0
}

// GitPullWithSigning constructs the git pull command arguments for execution in the terminal.
// This allows for interactive signing (e.g., GPG passphrase) by suspending the UI.
func (gp *GitPull) GitPullWithSigning(pullType string) []string {
	var gitPullArgs []string
	switch pullType {
	case GITPULL:
		gitPullArgs = []string{"pull", "--progress", "--no-edit"}
	case GITPULLREBASE:
		gitPullArgs = []string{"pull", "--progress", "--rebase", "--autostash", "--no-edit"}
	case GITPULLMERGE:
		gitPullArgs = []string{"pull", "--progress", "--no-rebase", "--no-edit"}
	}
	return gitPullArgs
}

// --------------------------------
//
// # Clear the Git Process Output
//
// --------------------------------
func (gp *GitPull) ClearGitPullOutput() {
	gp.gitPullOutputMu.Lock()
	defer gp.gitPullOutputMu.Unlock()
	gp.gitPullOutput = []string{}
}
