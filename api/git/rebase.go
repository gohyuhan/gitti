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

type GitRebase struct {
	gitRebaseOutput   []string
	gitRebaseOutputMu sync.RWMutex
	gitProcessLock    *GitProcessLock
	updateChannel     chan string
	logging           *logging.GittiLogging
}

// ----------------------------------
//
//	Initialize the git rebase handler with shared dependencies
//
// ----------------------------------
func InitGitRebase(updateChannel chan string, gitProcessLock *GitProcessLock, logging *logging.GittiLogging) *GitRebase {
	gitRebase := &GitRebase{
		gitRebaseOutput: []string{},
		updateChannel:   updateChannel,
		gitProcessLock:  gitProcessLock,
		logging:         logging,
	}

	return gitRebase
}

// --------------------------------
//
// return the git rebase output
//
// --------------------------------
func (gr *GitRebase) GetGitRebaseOutput() []string {
	gr.gitRebaseOutputMu.RLock()
	defer gr.gitRebaseOutputMu.RUnlock()

	copied := make([]string, len(gr.gitRebaseOutput))
	copy(copied, gr.gitRebaseOutput)
	return copied
}

// --------------------------------
//
// # Git Rebase (this rebase is different from the above rebase which default to its track upstream branch)
// * this rebase will be a more precise control on what remote (or local) and branch it will rebase on
//
// --------------------------------
func (gr *GitRebase) GitRebase(ctx context.Context, remote string, branchName string) int {
	if !gr.gitProcessLock.CanProceedWithGitOps() {
		return -1
	}
	defer func() {
		gr.gitProcessLock.ReleaseGitOpsLock()
	}()

	gr.ClearGitRebaseOutput()

	var gitArgs []string
	if remote == "" {
		gitArgs = []string{"rebase", branchName, "--progress", "--autostash", "--no-edit"}
	} else {
		gitArgs = []string{"pull", remote, branchName, "--progress", "--rebase", "--autostash", "--no-edit"}
	}
	cmdExecutor := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, true)

	// Combine stderr into stdout
	stdout, err := cmdExecutor.StdoutPipe()
	if err != nil {
		gr.logging.RegisterNewLog(logging.REBASE_OPS, "", logging.ERROR, fmt.Sprintf("[PIPE ERROR]: %s", err.Error()), false)
		return -1
	}
	cmdExecutor.Stderr = cmdExecutor.Stdout

	// Start the process
	if err := cmdExecutor.Start(); err != nil {
		gr.logging.RegisterNewLog(logging.REBASE_OPS, "", logging.ERROR, fmt.Sprintf("[START ERROR]: %s", err.Error()), false)
		return -1
	} else {
		gr.logging.RegisterNewLog(logging.REBASE_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
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
				gr.gitRebaseOutputMu.Lock()
				updatedCursorIndex, updatedGitRebaseOutput := handleProgressOutputStream(cursorIndex, scanner, gr.gitRebaseOutput)
				gr.gitRebaseOutput = updatedGitRebaseOutput
				cursorIndex = updatedCursorIndex
				gr.gitRebaseOutputMu.Unlock()
				if time.Since(lastSent) >= STREAMUPDATETHROTTLEMS*time.Millisecond {
					select {
					case gr.updateChannel <- GIT_REBASE_OUTPUT_UPDATE:
						lastSent = time.Now()
					default:
					}
				}
			}
		}
		// trigger an update once it ends
		gr.updateChannel <- GIT_REBASE_OUTPUT_UPDATE
	}()

	waitErr := cmdExecutor.Wait()
	wg.Wait()

	if ctx.Err() != nil {
		gr.logging.RegisterNewLog(logging.REBASE_OPS, strings.Join(gitArgs, " "), logging.WARN, fmt.Sprintf("[%s CANCELLED]: %s", logging.REBASE_OPS, ctx.Err().Error()), true)
		return -1
	}

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gr.logging.RegisterNewLog(logging.REBASE_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.REBASE_OPS, waitErr.Error()), true)
			return status
		}
		gr.logging.RegisterNewLog(logging.REBASE_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.REBASE_OPS, waitErr.Error()), true)
		return -1
	}
	return 0
}

// GitRebaseWithSigning constructs the git rebase command arguments for execution in the terminal.
// This allows for interactive signing (e.g., grG passphrase) by suspending the UI.
func (gr *GitRebase) GitRebaseWithSigning(remote string, branchName string) []string {
	var gitArgs []string
	if remote == "" {
		gitArgs = []string{"rebase", branchName, "--progress", "--autostash", "--no-edit"}
	} else {
		gitArgs = []string{"pull", remote, branchName, "--progress", "--rebase", "--autostash", "--no-edit"}
	}
	return gitArgs
}

// --------------------------------
//
// # Clear the Git Process Output
//
// --------------------------------
func (gr *GitRebase) ClearGitRebaseOutput() {
	gr.gitRebaseOutputMu.Lock()
	defer gr.gitRebaseOutputMu.Unlock()
	gr.gitRebaseOutput = []string{}
}
