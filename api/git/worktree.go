package git

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/logging"
)

type WorktreeInfo struct {
	WorktreePath        string
	IsMain              bool
	IsInCurrentWorktree bool
	IsLocked            bool
	IsPrunable          bool
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
	for outputIndex, worktreeInlineInfo := range parsedOutput {
		worktree := WorktreeInfo{IsMain: outputIndex == 0}

		// Split a single record into its attribute lines (worktree/HEAD/branch/...)
		parsedWorktreeInfoArray := processOutputIntoStringArrayWithCustomSeperator([]byte(worktreeInlineInfo), WORKTREE_FIELD_SEPARATOR)
		if len(parsedWorktreeInfoArray) < 1 {
			continue
		}

		// The `worktree <path>` attribute is always the first line of a record, a record
		// without a resolvable path is malformed and skipped entirely
		skipRecord := false
		for index, parsedWorktreeInfo := range parsedWorktreeInfoArray {
			if index == 0 {
				if !strings.HasPrefix(parsedWorktreeInfo, "worktree ") {
					gwt.logging.RegisterNewLog(logging.GET_LATEST_WORKTREE_INFO, "", logging.WARN, "skip malformed worktree record", false)
					skipRecord = true
					break
				}

				actualWorktreePath := strings.TrimSpace(strings.TrimPrefix(parsedWorktreeInfo, "worktree "))
				if actualWorktreePath == "" {
					skipRecord = true
					break
				}

				// A submodule worktree may be reported as its `.git/modules/...` git dir, resolve it back to the real checkout
				if isSubmoduleGitDirPath(actualWorktreePath) {
					processedPath, processError := processSubmoduleActualDirPath(actualWorktreePath)
					if processError != nil {
						gwt.logging.RegisterNewLog(logging.GET_LATEST_WORKTREE_INFO, "", logging.WARN, fmt.Sprintf("skip unresolvable submodule worktree path: %s", processError.Error()), false)
						skipRecord = true
						break
					}
					actualWorktreePath = processedPath
				}

				normalizedWorktreePath := filepath.Clean(actualWorktreePath)
				worktree.WorktreePath = normalizedWorktreePath
				worktree.IsInCurrentWorktree = normalizedWorktreePath == currentWorktreePath
			}

			if strings.HasPrefix(parsedWorktreeInfo, "prunable") {
				worktree.IsPrunable = true
			}

			if strings.HasPrefix(parsedWorktreeInfo, "locked") {
				worktree.IsLocked = true
			}
		}

		if skipRecord {
			continue
		}

		latestWorktreeInfo = append(latestWorktreeInfo, worktree)
	}

	gwt.allWorktree = latestWorktreeInfo
}

func (gwt *GitWorktree) AddNewWorktree(newWorktreeName string) {
	if !gwt.gitProcessLock.CanProceedWithGitOps() {
		gwt.logging.RegisterNewLog(logging.ADD_NEW_WORKTREE_OPS, "", logging.WARN, fmt.Sprintf("[WARN]: %s", gwt.gitProcessLock.OtherProcessRunningWarning()), false)
		return
	}
	defer func() {
		gwt.gitProcessLock.ReleaseGitOpsLock()
	}()
	newWorktreePath := filepath.Clean(filepath.Join(gwt.currentWorktreePath, fmt.Sprintf("../%s", newWorktreeName)))
	gitArgs := []string{"worktree", "add", newWorktreePath}

	newWorktreeCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	gwt.logging.RegisterNewLog(logging.ADD_NEW_WORKTREE_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	newWorktreeError := newWorktreeCmdExecutor.Run()

	if newWorktreeError != nil {
		gwt.logging.RegisterNewLog(logging.ADD_NEW_WORKTREE_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.ADD_NEW_WORKTREE_OPS, newWorktreeError.Error()), true)
	}
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
