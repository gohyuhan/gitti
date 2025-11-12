package git

import (
	"fmt"
	"sync"

	"gitti/cmd"
)

type StashInfo struct {
	Id      string
	Message string
}

type GitStash struct {
	allStash      []StashInfo
	errorLog      []error
	gitStashMutex sync.Mutex
}

// ----------------------------------
//
//	Related to Git Stash including untacked ( except ignored )
//
// ----------------------------------
func (gs *GitStash) GitStashAll() {
	gitArgs := []string{"stash", "--u"}

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	_, err := cmd.CombinedOutput()
	if err != nil {
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT STASH ERROR]: %w", err))
	}
}

// ----------------------------------
//
//	Related to Git UnStash all
//
// ----------------------------------
func (gs *GitStash) GitUnstashAll() {
	gitArgs := []string{"stash", "pop"}

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	_, err := cmd.CombinedOutput()
	if err != nil {
		gs.errorLog = append(gs.errorLog, fmt.Errorf("[GIT UNSTASH ERROR]: %w", err))
	}
}
