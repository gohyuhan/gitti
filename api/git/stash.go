package git

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/logging"
)

type StashInfo struct {
	Id      string
	Message string
}

type GitStash struct {
	allStash       []StashInfo
	gitProcessLock *GitProcessLock
	logging        *logging.GittiLogging
}

// ----------------------------------
//
//	Initialize the git stash handler with shared dependencies
//
// ----------------------------------
func InitGitStash(gitProcessLock *GitProcessLock, logging *logging.GittiLogging) *GitStash {
	gitStash := &GitStash{
		allStash:       []StashInfo{},
		gitProcessLock: gitProcessLock,
		logging:        logging,
	}

	return gitStash
}

// ----------------------------------
//
//	Return all stash entries
//
// ----------------------------------
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
		gs.logging.RegisterNewLog(logging.GET_STASH_INFO, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.GET_STASH_INFO, stashInfoErr.Error()), true)
		return
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
	gs.logging.RegisterNewLog(logging.STASH_ALL, strings.Join(gitArgs, " "), logging.INFO, "", true)

	var stashAllOutputStringArray []string
	if stashAllErr != nil {
		gs.logging.RegisterNewLog(logging.STASH_ALL, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.STASH_ALL, stashAllErr.Error()), true)
		if exitErr, ok := stashAllErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			return stashAllOutputStringArray, status
		}
		return stashAllOutputStringArray, -1
	}

	stashAllOutputStringArray = processGeneralGitOpsOutputIntoStringArray(stashAllOutput)

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
	gs.logging.RegisterNewLog(logging.STASH_FILE, strings.Join(gitArgs, " "), logging.INFO, "", true)

	var stashOutputStringArray []string
	if stashErr != nil {
		gs.logging.RegisterNewLog(logging.STASH_FILE, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.STASH_FILE, stashErr.Error()), true)
		if exitErr, ok := stashErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			return stashOutputStringArray, status
		}
		return stashOutputStringArray, -1
	}

	stashOutputStringArray = processGeneralGitOpsOutputIntoStringArray(stashOutput)

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
	gs.logging.RegisterNewLog(logging.STASH_APPLY, strings.Join(gitArgs, " "), logging.INFO, "", true)

	var stashApplyOutputStringArray []string
	if stashApplyErr != nil {
		gs.logging.RegisterNewLog(logging.STASH_APPLY, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.STASH_APPLY, stashApplyErr.Error()), true)
		if exitErr, ok := stashApplyErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			return stashApplyOutputStringArray, status
		}
		return stashApplyOutputStringArray, -1
	}
	stashApplyOutputStringArray = processGeneralGitOpsOutputIntoStringArray(stashApplyOutput)

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
	gs.logging.RegisterNewLog(logging.STASH_POP, strings.Join(gitArgs, " "), logging.INFO, "", true)

	var stashPopOutputStringArray []string
	if stashPopErr != nil {
		gs.logging.RegisterNewLog(logging.STASH_POP, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.STASH_POP, stashPopErr.Error()), true)
		if exitErr, ok := stashPopErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			return stashPopOutputStringArray, status
		}
		return stashPopOutputStringArray, -1
	}
	stashPopOutputStringArray = processGeneralGitOpsOutputIntoStringArray(stashPopOutput)

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
	gs.logging.RegisterNewLog(logging.STASH_DROP, strings.Join(gitArgs, " "), logging.INFO, "", true)

	var stashDropOutputStringArray []string
	if stashDropErr != nil {
		gs.logging.RegisterNewLog(logging.STASH_DROP, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.STASH_DROP, stashDropErr.Error()), true)
		if exitErr, ok := stashDropErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			return stashDropOutputStringArray, status
		}
		return stashDropOutputStringArray, -1
	}
	stashDropOutputStringArray = processGeneralGitOpsOutputIntoStringArray(stashDropOutput)

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
			gs.logging.RegisterNewLog(logging.STASH_DETAIL, strings.Join(gitArgs, " "), logging.WARN, fmt.Sprintf("[%s CANCELLED]", logging.STASH_DETAIL), true)
			return parsedDetail
		}
		gs.logging.RegisterNewLog(logging.STASH_DETAIL, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.STASH_DETAIL, detailCmdErr.Error()), true)
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
			return false, ctx.Err()
		}
		return false, showErr
	}

	// Start the process
	if err := showCmdExecutor.Start(); err != nil {
		return false, err
	}

	defer func() {
		showCmdExecutor.Wait()
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
		return false, err
	}

	return true, nil
}
