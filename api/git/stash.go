package git

import (
	"fmt"
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
//	Related to Git Stash including untracked ( except ignored )
//
// ----------------------------------
func (gs *GitStash) GitStashAll(message string) {
	if !gs.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gs.gitProcessLock.ReleaseGitOpsLock()

	gitArgs := []string{"stash", "push", "-u", "-m", message}

	stashAllCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	_, stashAllErr := stashAllCmdExecutor.CombinedOutput()
	if stashAllErr != nil {
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH ALL ERROR]: %w", stashAllErr))
	}
}

// ----------------------------------
//
// # Stash File changes
//
// ----------------------------------
func (gs *GitStash) GitStashFile(filePathName string, message string) {
	if !gs.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gs.gitProcessLock.ReleaseGitOpsLock()

	var gitArgs []string
	if message == "" {
		gitArgs = []string{"stash", "push", "-u", filePathName}
	} else {
		gitArgs = []string{"stash", "push", "-u", "-m", message, filePathName}
	}

	stashCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	_, stashErr := stashCmdExecutor.CombinedOutput()
	if stashErr != nil {
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH ERROR]: %w", stashErr))
	}
}

// ----------------------------------
//
// # Git stash apply
//
// ----------------------------------
func (gs *GitStash) GitStashApply(stashId string) {
	if !gs.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gs.gitProcessLock.ReleaseGitOpsLock()

	gitArgs := []string{"stash", "apply", stashId}

	stashApplyCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	_, stashApplyErr := stashApplyCmdExecutor.CombinedOutput()
	if stashApplyErr != nil {
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH APPLY ERROR]: %w", stashApplyErr))
	}
}

// ----------------------------------
//
// # Git stash pop
//
// ----------------------------------
func (gs *GitStash) GitStashPop(stashId string) {
	if !gs.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gs.gitProcessLock.ReleaseGitOpsLock()

	gitArgs := []string{"stash", "pop", stashId}

	stashPopCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	_, stashPopErr := stashPopCmdExecutor.CombinedOutput()
	if stashPopErr != nil {
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH POP ERROR]: %w", stashPopErr))
	}
}

// ----------------------------------
//
// # Git stash drop
//
// ----------------------------------
func (gs *GitStash) GitStashDrop(stashId string) {
	if !gs.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gs.gitProcessLock.ReleaseGitOpsLock()

	gitArgs := []string{"stash", "drop", stashId}

	stashDropCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	_, stashDropErr := stashDropCmdExecutor.CombinedOutput()
	if stashDropErr != nil {
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH DROP ERROR]: %w", stashDropErr))
	}
}

// ----------------------------------
//
// # Git stash see deatil
//
// ----------------------------------
func (gs *GitStash) GitStashDetail(stashId string) []string {
	var parsedDetail []string
	gitArgs := []string{"stash", "show", "-p", "-u", stashId}

	detailCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, true)
	stashDetailOutput, detailCmdErr := detailCmdExecutor.Output()
	if detailCmdErr != nil {
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH DETAIL ERROR]: %w", detailCmdErr))
		return parsedDetail
	}
	parsedDetail = strings.Split(strings.TrimSpace(string(stashDetailOutput)), "\n")
	return parsedDetail
}
