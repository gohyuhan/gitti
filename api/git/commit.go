package git

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gohyuhan/gitti/executor"
)

type GitCommit struct {
	errorLog              []error
	gitCommitOutput       []string
	gitCommitOutputMu     sync.RWMutex
	gitRemotePushOutput   []string
	gitRemotePushOutputMu sync.RWMutex
	updateChannel         chan string
	gitProcessLock        *GitProcessLock
	remote                []GitRemote
}

type LatestCommitMsgAndDesc struct {
	Message     string
	Description string
}

func InitGitCommit(updateChannel chan string, gitProcessLock *GitProcessLock) *GitCommit {
	gitCommit := GitCommit{
		gitCommitOutput:     []string{},
		gitRemotePushOutput: []string{},
		updateChannel:       updateChannel,
		gitProcessLock:      gitProcessLock,
	}

	return &gitCommit
}

// ----------------------------------
//
//	Return git commit output
//
// ----------------------------------
func (gc *GitCommit) GitCommitOutput() []string {
	gc.gitCommitOutputMu.RLock()
	defer gc.gitCommitOutputMu.RUnlock()

	copied := make([]string, len(gc.gitCommitOutput))
	copy(copied, gc.gitCommitOutput)
	return copied
}

// ----------------------------------
//
//	Return git remote push output
//
// ----------------------------------
func (gc *GitCommit) GitRemotePushOutput() []string {
	gc.gitRemotePushOutputMu.RLock()
	defer gc.gitRemotePushOutputMu.RUnlock()

	copied := make([]string, len(gc.gitRemotePushOutput))
	copy(copied, gc.gitRemotePushOutput)
	return copied
}

// ----------------------------------
//
//	Related to Git Commit
//
// ----------------------------------
func (gc *GitCommit) GitCommit(ctx context.Context, message, description string, isAmendCommit bool) int {
	if !gc.gitProcessLock.CanProceedWithGitOps() {
		return -1
	}

	defer func() {
		gc.gitProcessLock.ReleaseGitOpsLock()
	}()

	gc.ClearGitCommitOutput()
	gitArgs := []string{"commit", "-m", message}
	if isAmendCommit {
		gitArgs = []string{"commit", "--amend", "-m", message}
	}
	if len(description) > 0 {
		gitArgs = append(gitArgs, "-m", description)
	}

	commitCmd := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, true)

	// Combine stderr into stdout
	stdout, err := commitCmd.StdoutPipe()
	if err != nil {
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[PIPE ERROR]: %w", err))
		return -1
	}
	commitCmd.Stderr = commitCmd.Stdout

	// Start the process
	if err := commitCmd.Start(); err != nil {
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[START ERROR]: %w", err))
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
		lastSent := time.Time{}
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return // Stop immediately on cancel
			default:
				gc.gitCommitOutputMu.Lock()
				updatedCursorIndex, updatedGitCommitOutput := handleProgressOutputStream(cursorIndex, scanner, gc.gitCommitOutput)
				gc.gitCommitOutput = updatedGitCommitOutput
				cursorIndex = updatedCursorIndex
				gc.gitCommitOutputMu.Unlock()
				if time.Since(lastSent) >= STREAMUPDATETHROTTLEMS*time.Millisecond {
					if isAmendCommit {
						select {
						case gc.updateChannel <- GIT_AMEND_COMMIT_OUTPUT_UPDATE:
							lastSent = time.Now()
						default:
						}
					} else {
						select {
						case gc.updateChannel <- GIT_COMMIT_OUTPUT_UPDATE:
							lastSent = time.Now()
						default:
						}
					}
				}
			}
		}
		// trigger an update once it ends
		gc.updateChannel <- GIT_COMMIT_OUTPUT_UPDATE
	}()

	waitErr := commitCmd.Wait()
	wg.Wait()

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gc.errorLog = append(gc.errorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", waitErr))
			return status
		}
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[UNEXPECTED ERROR]: %w", waitErr))
		return -1
	}

	return 0
}

func (gc *GitCommit) ClearGitCommitOutput() {
	gc.gitCommitOutputMu.Lock()
	defer gc.gitCommitOutputMu.Unlock()
	gc.gitCommitOutput = []string{}
}

// ----------------------------------
//
//	Related to Git Commit (Amend)
//
// ----------------------------------
func (gc *GitCommit) GetLatestCommitMsgAndDesc() LatestCommitMsgAndDesc {
	gitArgs := []string{"log", "-1", "--pretty=format:%s%n%b", "HEAD"}
	latestCommitCmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	commitMsgAndDesc, cmdErr := latestCommitCmd.Output()
	if cmdErr != nil {
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[GET LATEST COMMIT INFO ERROR]: %w", cmdErr))
		return LatestCommitMsgAndDesc{}
	}

	parsed := strings.SplitN(string(commitMsgAndDesc), "\n", 2)
	title := parsed[0]
	description := ""
	if len(parsed) > 1 {
		description = parsed[1]
	}

	return LatestCommitMsgAndDesc{
		Message:     title,
		Description: description,
	}
}

// ----------------------------------
//
//	Related to Git Push
//
// ----------------------------------
func (gc *GitCommit) GitPush(ctx context.Context, originName string, pushType string, currentCheckOutBranch string) int {
	if !gc.gitProcessLock.CanProceedWithGitOps() {
		return -1
	}
	defer func() {
		gc.gitProcessLock.ReleaseGitOpsLock()
	}()

	gc.ClearGitRemotePushOutput()

	// check if the checkoutbranch has upstream if not include "-u" flag
	_, hasUpstream := hasUpStream()
	var gitArgs []string
	if !hasUpstream {
		gitArgs = []string{"push", "-u"}
	} else {
		gitArgs = []string{"push"}
	}
	switch pushType {
	case FORCEPUSHSAFE:
		gitArgs = append(gitArgs, []string{"--progress", "--force-with-lease", originName}...)
	case FORCEPUSHDANGEROUS:
		gitArgs = append(gitArgs, []string{"--progress", "--force", originName}...)
	default:
		gitArgs = append(gitArgs, []string{"--progress", originName}...)
	}

	// include the current checkout branch name at the end if there was no upstream so that git know which branch to push
	if !hasUpstream {
		gitArgs = append(gitArgs, currentCheckOutBranch)
	}

	cmd := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, true)
	// Disable interactive prompts for credentials
	cmd.Env = append(os.Environ(), "GIT_ASKPASS=true", "GIT_TERMINAL_PROMPT=0")

	// Combine stderr into stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[PIPE ERROR]: %w", err))
		return -1
	}
	cmd.Stderr = cmd.Stdout

	// Start the process
	if err := cmd.Start(); err != nil {
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[START ERROR]: %w", err))
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
		lastSent := time.Time{}
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return // Stop immediately on cancel
			default:
				gc.gitRemotePushOutputMu.Lock()
				updatedCursorIndex, updatedGitRemotePushOutput := handleProgressOutputStream(cursorIndex, scanner, gc.gitRemotePushOutput)
				gc.gitRemotePushOutput = updatedGitRemotePushOutput
				cursorIndex = updatedCursorIndex
				gc.gitRemotePushOutputMu.Unlock()
				if time.Since(lastSent) >= STREAMUPDATETHROTTLEMS*time.Millisecond {
					select {
					case gc.updateChannel <- GIT_REMOTE_PUSH_OUTPUT_UPDATE:
						lastSent = time.Now()
					default:
					}
				}
			}
		}
		// trigger an update once it ends
		gc.updateChannel <- GIT_REMOTE_PUSH_OUTPUT_UPDATE
	}()

	waitErr := cmd.Wait()
	wg.Wait()

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gc.errorLog = append(gc.errorLog, fmt.Errorf("[GIT PUSH ERROR]: %w", waitErr))
			return status
		}
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[UNEXPECTED ERROR]: %w", waitErr))
		return -1
	}
	return 0
}

func (gc *GitCommit) ClearGitRemotePushOutput() {
	gc.gitRemotePushOutputMu.Lock()
	defer gc.gitRemotePushOutputMu.Unlock()
	gc.gitRemotePushOutput = []string{}
}
