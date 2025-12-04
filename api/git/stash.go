package git

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gohyuhan/gitti/executor"
)

type StashInfo struct {
	Id      string
	Message string
}

type GitStash struct {
	allStash       []StashInfo
	errorLog       []error
	gitProcessLock *GitProcessLock
}

func InitGitStash(gitProcessLock *GitProcessLock) *GitStash {
	gitStash := &GitStash{
		allStash:       []StashInfo{},
		errorLog:       []error{},
		gitProcessLock: gitProcessLock,
	}

	return gitStash
}

func (gs *GitStash) AllStash() []StashInfo {
	copied := make([]StashInfo, len(gs.allStash))
	copy(copied, gs.allStash)
	return copied
}

// ----------------------------------
//
//	Get Latest Info For Stash
//
// ----------------------------------
func (gs *GitStash) GetLatestStashInfo() {
	gitArgs := []string{"stash", "list", "--format=%gd %s"}
	stashInfoCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	stashInfoOutput, stashInfoErr := stashInfoCmdExecutor.Output()
	if stashInfoErr != nil {
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH INFO RETRIEVE ERROR]: %w", stashInfoErr))
	}

	parsedStashInfo := strings.Split(string(stashInfoOutput), "\n")
	if len(parsedStashInfo) < 1 {
		return
	}

	var stashInfoArray []StashInfo
	for _, stashInfo := range parsedStashInfo {
		parsedInfo := strings.SplitN(stashInfo, " ", 2)
		if len(parsedInfo) < 2 {
			continue
		}
		stashInfoArray = append(stashInfoArray, StashInfo{
			Id:      strings.TrimSpace(parsedInfo[0]),
			Message: strings.TrimSpace(parsedInfo[1]),
		})
	}

	gs.allStash = stashInfoArray
}

// ----------------------------------
//
//	Related to Git Stash All including untracked ( both index and worktree except ignored )
//
// ----------------------------------
func (gs *GitStash) GitStashAll(message string) ([]string, int) {
	if !gs.gitProcessLock.CanProceedWithGitOps() {
		return []string{gs.gitProcessLock.OtherProcessRunningWarning()}, -1
	}
	defer gs.gitProcessLock.ReleaseGitOpsLock()

	gitArgs := []string{"stash", "push", "-u", "-m", message}

	stashAllCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	stashAllOutput, stashAllErr := stashAllCmdExecutor.CombinedOutput()
	stashAllOutputStringArray := processGeneralGitOpsOutputIntoStringArray(stashAllOutput)
	if stashAllErr != nil {
		if exitErr, ok := stashAllErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			return stashAllOutputStringArray, status
		}
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH ALL ERROR]: %w", stashAllErr))
		return stashAllOutputStringArray, -1
	}

	return stashAllOutputStringArray, 0
}

// ----------------------------------
//
// # Stash File changes
//
// ----------------------------------
func (gs *GitStash) GitStashFile(filePathName string, message string) ([]string, int) {
	if !gs.gitProcessLock.CanProceedWithGitOps() {
		return []string{gs.gitProcessLock.OtherProcessRunningWarning()}, -1
	}
	defer gs.gitProcessLock.ReleaseGitOpsLock()

	// Parse renamed/copied file format (old -> new)
	actualFileName := filePathName
	if strings.Contains(filePathName, "->") {
		parts := strings.Split(filePathName, "->")
		if len(parts) >= 2 {
			actualFileName = strings.TrimSpace(parts[1])
		}
	}

	var gitArgs []string
	if message == "" {
		gitArgs = []string{"stash", "push", "-u", actualFileName}
	} else {
		gitArgs = []string{"stash", "push", "-u", "-m", message, actualFileName}
	}

	stashCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	stashOutput, stashErr := stashCmdExecutor.CombinedOutput()
	stashOutputStringArray := processGeneralGitOpsOutputIntoStringArray(stashOutput)
	if stashErr != nil {
		if exitErr, ok := stashErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			return stashOutputStringArray, status
		}
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH ERROR]: %w", stashErr))
		return stashOutputStringArray, -1
	}

	return stashOutputStringArray, 0
}

// ----------------------------------
//
// # Git stash apply
//
// ----------------------------------
func (gs *GitStash) GitStashApply(stashId string) ([]string, int) {
	if !gs.gitProcessLock.CanProceedWithGitOps() {
		return []string{gs.gitProcessLock.OtherProcessRunningWarning()}, -1
	}
	defer gs.gitProcessLock.ReleaseGitOpsLock()

	gitArgs := []string{"stash", "apply", stashId}

	stashApplyCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	stashApplyOutput, stashApplyErr := stashApplyCmdExecutor.CombinedOutput()
	stashApplyOutputStringArray := processGeneralGitOpsOutputIntoStringArray(stashApplyOutput)
	if stashApplyErr != nil {
		if exitErr, ok := stashApplyErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			return stashApplyOutputStringArray, status
		}
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH APPLY ERROR]: %w", stashApplyErr))
		return stashApplyOutputStringArray, -1
	}

	return stashApplyOutputStringArray, 0
}

// ----------------------------------
//
// # Git stash pop
//
// ----------------------------------
func (gs *GitStash) GitStashPop(stashId string) ([]string, int) {
	if !gs.gitProcessLock.CanProceedWithGitOps() {
		return []string{gs.gitProcessLock.OtherProcessRunningWarning()}, -1
	}
	defer gs.gitProcessLock.ReleaseGitOpsLock()

	gitArgs := []string{"stash", "pop", stashId}

	stashPopCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	stashPopOutput, stashPopErr := stashPopCmdExecutor.CombinedOutput()
	stashPopOutputStringArray := processGeneralGitOpsOutputIntoStringArray(stashPopOutput)
	if stashPopErr != nil {
		if exitErr, ok := stashPopErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			return stashPopOutputStringArray, status
		}
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH POP ERROR]: %w", stashPopErr))
		return stashPopOutputStringArray, -1
	}

	return stashPopOutputStringArray, 0
}

// ----------------------------------
//
// # Git stash drop
//
// ----------------------------------
func (gs *GitStash) GitStashDrop(stashId string) ([]string, int) {
	if !gs.gitProcessLock.CanProceedWithGitOps() {
		return []string{gs.gitProcessLock.OtherProcessRunningWarning()}, -1
	}
	defer gs.gitProcessLock.ReleaseGitOpsLock()

	gitArgs := []string{"stash", "drop", stashId}

	stashDropCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	stashDropOutput, stashDropErr := stashDropCmdExecutor.CombinedOutput()
	stashDropOutputStringArray := processGeneralGitOpsOutputIntoStringArray(stashDropOutput)
	if stashDropErr != nil {
		if exitErr, ok := stashDropErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			return stashDropOutputStringArray, status
		}
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH DROP ERROR]: %w", stashDropErr))
		return stashDropOutputStringArray, -1
	}

	return stashDropOutputStringArray, 0
}

// ----------------------------------
//
// # Git stash detail
//
// ----------------------------------
func (gs *GitStash) GitStashDetail(ctx context.Context, stashId string) []string {
	var parsedDetail []string

	// Use -p flag for small stashes to show patch details
	var gitArgs []string
	isSmall, err := gs.isStashSmall(ctx, stashId)
	if err != nil {
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[STASH DETAIL OPERATION CANCELLED DUE TO CONTEXT SWITCHING]: %w", ctx.Err()))
		return parsedDetail
	}
	if isSmall {
		gitArgs = []string{"stash", "show", "-p", "-u", stashId}
	} else {
		gitArgs = []string{"stash", "show", "-u", stashId}
	}

	detailCmdExecutor := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, true)
	stashDetailOutput, detailCmdErr := detailCmdExecutor.Output()
	if detailCmdErr != nil {
		if ctx.Err() != nil {
			// This catches context.Canceled
			gs.errorLog = append(gs.errorLog, fmt.Errorf("[STASH DETAIL OPERATION CANCELLED DUE TO CONTEXT SWITCHING]: %w", ctx.Err()))
			return parsedDetail
		}
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH DETAIL ERROR]: %w", detailCmdErr))
		return parsedDetail
	}
	parsedDetail = processGeneralGitOpsOutputIntoStringArray(stashDetailOutput)
	return parsedDetail
}

// ----------------------------------
//
// # Helper to determine if stash is small
//
// ----------------------------------
func (gs *GitStash) isStashSmall(ctx context.Context, stashId string) (bool, error) {
	// Fast early-exit: use numstat which shows all files (tracked + untracked)
	// Stop reading after threshold to avoid processing millions of files
	const fileThreshold = 25
	fileCount := 0

	gitArgs := []string{"stash", "show", "-u", "--name-only", stashId}
	showCmdExecutor := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, false)
	showOutput, showErr := showCmdExecutor.StdoutPipe()

	if showErr != nil {
		if ctx.Err() != nil {
			// This catches context.Canceled
			gs.errorLog = append(gs.errorLog, fmt.Errorf("[DETERMINE STASHOPERATION CANCELLED DUE TO CONTEXT SWITCHING]: %w", ctx.Err()))
			return false, ctx.Err()
		}
		return false, showErr
	}

	// Start the process
	if err := showCmdExecutor.Start(); err != nil {
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[DETERMINE STASH SIZE START ERROR]: %w", err))
		return false, err
	}

	defer func() {
		if err := showCmdExecutor.Wait(); err != nil {
			// Only log if it's not a context cancellation
			if ctx.Err() == nil {
				gs.errorLog = append(gs.errorLog, fmt.Errorf("[DETERMINE STASH SIZE WAIT ERROR]: %w", err))
			}
		}
	}()

	scanner := bufio.NewScanner(showOutput)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return false, fmt.Errorf("[DETERMINE STASHOPERATION CANCELLED DUE TO CONTEXT SWITCHING]: %w", ctx.Err())
		default:
			fileCount++
			if fileCount > fileThreshold {
				return false, nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		if ctx.Err() != nil {
			return false, ctx.Err()
		}
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[SCANNER ERROR]: %w", err))
		return false, err
	}

	return true, nil
}
