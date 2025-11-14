package git

import (
	"fmt"
	"strings"

	"gitti/cmd"
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
	return gs.allStash
}

// ----------------------------------
//
//	Get Latest Info For Stash
//
// ----------------------------------
func (gs *GitStash) GetLatestStashInfo() {
	gitArgs := []string{"stash", "list", "--format=%gd %s"}
	stashInfocmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
	stashInfoOutput, stashInfoErr := stashInfocmd.Output()
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

	stashAllCmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
	_, stashAllErr := stashAllCmd.CombinedOutput()
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

	stashCmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
	_, stashErr := stashCmd.CombinedOutput()
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

	stashApplyCmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
	_, stashApplyErr := stashApplyCmd.CombinedOutput()
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

	stashPopCmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
	_, stashPopErr := stashPopCmd.CombinedOutput()
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

	stashDropCmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
	_, stashDropErr := stashDropCmd.CombinedOutput()
	if stashDropErr != nil {
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH DROP ERROR]: %w", stashDropErr))
	}
}
