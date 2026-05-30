package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/logging"
)

type WorktreeInfo struct {
	WorktreePath        string
	IsMain              bool
	IsInCurrentWorktree bool
}

type GitWorktree struct {
	updateChannel       chan string
	allWorktree         []WorktreeInfo
	logging             *logging.GittiLogging
	gitProcessLock      *GitProcessLock
	currentWorktreePath string
}

// ------------------------------------
//
//	Initialize the git worktree handler with shared dependencies
//
// ------------------------------------
func InitGitWorktree(updateChannel chan string, gitProcessLock *GitProcessLock, logging *logging.GittiLogging, absoluteWorktreePath string) *GitWorktree {
	gitWorktree := GitWorktree{
		updateChannel:       updateChannel,
		gitProcessLock:      gitProcessLock,
		logging:             logging,
		currentWorktreePath: absoluteWorktreePath,
	}
	return &gitWorktree
}

// ------------------------------------
//
//	Return a defensive copy of the latest cached worktree infos
//
// ------------------------------------
func (gwt *GitWorktree) AllWorktree() []WorktreeInfo {
	copied := make([]WorktreeInfo, len(gwt.allWorktree))
	copy(copied, gwt.allWorktree)
	return copied
}

// ------------------------------------
//
//	Refresh the cached worktree infos by parsing `git worktree list --porcelain`,
//	resolving submodule worktree paths and dropping any unusable records
//
// ------------------------------------
func (gwt *GitWorktree) GetLatestWorktreeInfos() {
	gitArgs := []string{"worktree", "list", "--porcelain"}

	getLatestGitWoktreeCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	getLatestGitWorktreeCmdOutput, getLatestGitWorktreeCmdErr := getLatestGitWoktreeCmdExecutor.Output()

	if getLatestGitWorktreeCmdErr != nil {
		gwt.logging.RegisterNewLog(logging.GET_LATEST_WORKTREE_INFO, "", logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.GET_LATEST_WORKTREE_INFO, getLatestGitWorktreeCmdErr.Error()), true)
		return
	}

	var latestWorktreeInfo []WorktreeInfo
	// Each record is a block of newline-separated attribute lines, records are blank-line separated
	parsedOutput := processOutputIntoStringArrayWithCustomSeperator(getLatestGitWorktreeCmdOutput, WORKTREE_RECORD_SEPARATOR)

	if len(parsedOutput) < 1 {
		gwt.allWorktree = latestWorktreeInfo
		return
	}

	currentWorktreePath := filepath.Clean(gwt.currentWorktreePath)
	for index, worktreeInlineInfo := range parsedOutput {
		// Split a single record into its attribute lines (worktree/HEAD/branch/...)
		parsedWorktreeInfoArray := processOutputIntoStringArrayWithCustomSeperator([]byte(worktreeInlineInfo), WORKTREE_FIELD_SEPARATOR)
		if len(parsedWorktreeInfoArray) < 1 {
			continue
		}

		// The `worktree <path>` attribute is always the first line of a record
		firstLine := strings.TrimSpace(parsedWorktreeInfoArray[0])
		if !strings.HasPrefix(firstLine, "worktree ") {
			gwt.logging.RegisterNewLog(logging.GET_LATEST_WORKTREE_INFO, "", logging.WARN, "skip malformed worktree record", false)
			continue
		}

		actualWorktreePath := strings.TrimSpace(strings.TrimPrefix(firstLine, "worktree "))
		if actualWorktreePath == "" {
			continue
		}

		// A submodule worktree may be reported as its `.git/modules/...` git dir, resolve it back to the real checkout
		if isSubmoduleGitDirPath(actualWorktreePath) {
			processedPath, processError := processSubmoduleActualDirPath(actualWorktreePath)
			if processError != nil {
				gwt.logging.RegisterNewLog(logging.GET_LATEST_WORKTREE_INFO, "", logging.WARN, fmt.Sprintf("skip unresolvable submodule worktree path: %s", processError.Error()), false)
				continue
			}
			actualWorktreePath = processedPath
		}

		normalizedWorktreePath := filepath.Clean(actualWorktreePath)
		// Drop bare/pruned/non-checked-out entries that have no real working tree on disk
		if !isUsableWorktreePath(normalizedWorktreePath) {
			gwt.logging.RegisterNewLog(logging.GET_LATEST_WORKTREE_INFO, "", logging.WARN, fmt.Sprintf("skip unusable worktree path: %s", normalizedWorktreePath), false)
			continue
		}

		worktree := WorktreeInfo{
			WorktreePath: normalizedWorktreePath,
			// git always lists the main worktree first
			IsMain:              index == 0,
			IsInCurrentWorktree: normalizedWorktreePath == currentWorktreePath,
		}

		latestWorktreeInfo = append(latestWorktreeInfo, worktree)
	}

	gwt.allWorktree = latestWorktreeInfo
}

// ------------------------------------
//
//	Resolve a submodule `.git/modules/...` git dir to its real working tree path
//	by reading the `core.worktree` value from that git dir's config
//
// ------------------------------------
func processSubmoduleActualDirPath(submoduleWorktreeGitPath string) (string, error) {
	var realPath string
	worktreeRelative, error := getSubmoduleWorktreePath(filepath.Join(submoduleWorktreeGitPath, "config"))
	if error != nil {
		return realPath, error
	}

	// core.worktree is stored relative to the git dir, join them to get the absolute checkout path
	realPath = filepath.Clean(filepath.Join(submoduleWorktreeGitPath, worktreeRelative))
	return realPath, nil
}

// ------------------------------------
//
//	Report whether the given path points inside a submodule or nested submodule git dir (`.git/modules/...`)
//
// ------------------------------------
func isSubmoduleGitDirPath(path string) bool {
	normalized := filepath.ToSlash(filepath.Clean(path))
	return strings.Contains(normalized, "/.git/modules/")
}

// ------------------------------------
//
//	Report whether the path is a real worktree on disk, an existing directory that contains a `.git` entry
//
// ------------------------------------
func isUsableWorktreePath(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}

	// .git is a dir for the main worktree and a file for linked worktrees, either is valid
	gitPath := filepath.Join(path, ".git")
	_, gitErr := os.Stat(gitPath)
	return gitErr == nil
}
