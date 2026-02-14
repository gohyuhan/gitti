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

type GitCommit struct {
	gitCommitOutput       []string
	gitCommitOutputMu     sync.RWMutex
	gitRemotePushOutput   []string
	gitRemotePushOutputMu sync.RWMutex
	updateChannel         chan string
	gitProcessLock        *GitProcessLock
	logging               *logging.GittiLogging
}

type LatestCommitMsgAndDesc struct {
	Message     string
	Description string
}

func InitGitCommit(updateChannel chan string, gitProcessLock *GitProcessLock, logging *logging.GittiLogging) *GitCommit {
	gitCommit := GitCommit{
		gitCommitOutput:     []string{},
		gitRemotePushOutput: []string{},
		updateChannel:       updateChannel,
		gitProcessLock:      gitProcessLock,
		logging:             logging,
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
		gc.logging.RegisterNewLog(logging.COMMIT_OPS, "", logging.ERROR, fmt.Sprintf("[PIPE ERROR]: %s", err.Error()), false)
		return -1
	}
	commitCmd.Stderr = commitCmd.Stdout

	// Start the process
	if err := commitCmd.Start(); err != nil {
		gc.logging.RegisterNewLog(logging.COMMIT_OPS, "", logging.ERROR, fmt.Sprintf("[START ERROR]: %s", err.Error()), false)
		return -1
	} else {
		gc.logging.RegisterNewLog(logging.COMMIT_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
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
		if isAmendCommit {
			gc.updateChannel <- GIT_AMEND_COMMIT_OUTPUT_UPDATE
		} else {
			gc.updateChannel <- GIT_COMMIT_OUTPUT_UPDATE
		}
	}()

	waitErr := commitCmd.Wait()
	wg.Wait()

	if ctx.Err() != nil {
		gc.logging.RegisterNewLog(logging.COMMIT_OPS, strings.Join(gitArgs, " "), logging.WARN, fmt.Sprintf("[%s CANCELLED]", logging.COMMIT_OPS), true)
		return -1
	}

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gc.logging.RegisterNewLog(logging.COMMIT_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.COMMIT_OPS, waitErr.Error()), true)
			return status
		}
		gc.logging.RegisterNewLog(logging.COMMIT_OPS, "", logging.ERROR, fmt.Sprintf("[UNEXPECTED ERROR]: %s", waitErr.Error()), false)
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
	gc.logging.RegisterNewLog(logging.GET_LATEST_COMMIT_MSG_AND_DESC_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	if cmdErr != nil {
		gc.logging.RegisterNewLog(logging.GET_LATEST_COMMIT_MSG_AND_DESC_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.GET_LATEST_COMMIT_MSG_AND_DESC_OPS, cmdErr.Error()), true)
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

	// Combine stderr into stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		gc.logging.RegisterNewLog(logging.GIT_PUSH_OPS, "", logging.ERROR, fmt.Sprintf("[PIPE ERROR]: %s", err.Error()), false)
		return -1
	}
	cmd.Stderr = cmd.Stdout

	// Start the process
	if err := cmd.Start(); err != nil {
		gc.logging.RegisterNewLog(logging.GIT_PUSH_OPS, "", logging.ERROR, fmt.Sprintf("[START ERROR]: %s", err.Error()), false)
		return -1
	} else {
		gc.logging.RegisterNewLog(logging.GIT_PUSH_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
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

	if ctx.Err() != nil {
		gc.logging.RegisterNewLog(logging.GIT_PUSH_OPS, strings.Join(gitArgs, " "), logging.WARN, fmt.Sprintf("[%s CANCELLED]", logging.GIT_PUSH_OPS), true)
		return -1
	}

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gc.logging.RegisterNewLog(logging.GIT_PUSH_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.GIT_PUSH_OPS, exitErr.Error()), true)
			return status
		}
		gc.logging.RegisterNewLog(logging.GIT_PUSH_OPS, "", logging.ERROR, fmt.Sprintf("[UNEXPECTED ERROR]: %s", waitErr.Error()), false)
		return -1
	}

	return 0
}

func (gc *GitCommit) ClearGitRemotePushOutput() {
	gc.gitRemotePushOutputMu.Lock()
	defer gc.gitRemotePushOutputMu.Unlock()
	gc.gitRemotePushOutput = []string{}
}

// ----------------------------------
//
//	Related to Git Commit RESET (apply to the latest commit only)
//
// ----------------------------------
func (gc *GitCommit) GitResetLatestCommit(resetType string) {
	if !gc.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer func() {
		gc.gitProcessLock.ReleaseGitOpsLock()
	}()

	var gitArgs []string

	switch resetType {
	case RESETSOFT:
		gitArgs = []string{"reset", "--soft", "HEAD~1"}
	case RESETHARD:
		gitArgs = []string{"reset", "--hard", "HEAD~1"}
	case RESETMIXED:
		gitArgs = []string{"reset", "--mixed", "HEAD~1"}
	default:
		// we default to reset mixed as default option for reset
		gitArgs = []string{"reset", "--mixed", "HEAD~1"}
	}

	commitLatestResetCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	err := commitLatestResetCmdExecutor.Run()
	gc.logging.RegisterNewLog(logging.GIT_RESET_LATEST_COMMIT_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	if err != nil {
		gc.logging.RegisterNewLog(logging.GIT_RESET_LATEST_COMMIT_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.GIT_RESET_LATEST_COMMIT_OPS, err.Error()), true)
		return
	}
}

// ----------------------------------
//
//	Related to Git Commit RESET (apply to selected commit [using commit hash])
//
// ----------------------------------
func (gc *GitCommit) GitResetToSelectedCommit(resetType string, commitHash string) {
	if !gc.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer func() {
		gc.gitProcessLock.ReleaseGitOpsLock()
	}()

	var gitArgs []string

	switch resetType {
	case RESETSOFT:
		gitArgs = []string{"reset", "--soft", commitHash}
	case RESETHARD:
		gitArgs = []string{"reset", "--hard", commitHash}
	case RESETMIXED:
		gitArgs = []string{"reset", "--mixed", commitHash}
	default:
		// we default to reset mixed as default option for reset
		gitArgs = []string{"reset", "--mixed", commitHash}
	}

	resetToSelectedCommitCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	err := resetToSelectedCommitCmdExecutor.Run()
	gc.logging.RegisterNewLog(logging.GIT_RESET_TO_SELECTED_COMMIT_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	if err != nil {
		gc.logging.RegisterNewLog(logging.GIT_RESET_TO_SELECTED_COMMIT_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.GIT_RESET_TO_SELECTED_COMMIT_OPS, err.Error()), true)
		return
	}
}
