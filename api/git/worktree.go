package git

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/logging"
)

type WorktreeInfo struct {
	WorktreePath        string
	IsMain              bool
	IsInCurrentWorktree bool
	IsLocked            bool
	LockReason          string
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
				worktree.LockReason = strings.TrimSpace(strings.TrimPrefix(parsedWorktreeInfo, "locked"))
			}
		}

		if skipRecord {
			continue
		}

		latestWorktreeInfo = append(latestWorktreeInfo, worktree)
	}

	gwt.allWorktree = latestWorktreeInfo
}

// ------------------------------------
//
//	Add a new worktree as a sibling of the current worktree, named newWorktreeName.
//	If checkoutWorktreeBranchName is non-empty, it is passed to `git worktree add`
//	to check out an existing local or remote branch; a name that is neither a local
//	nor remote branch will error. If empty, git derives a branch from the path.
//	Acquires the git ops lock first; returns the lock-busy warning if another git
//	process is running. Returns the command output and whether the add succeeded.
//
// ------------------------------------
func (gwt *GitWorktree) AddNewWorktree(ctx context.Context, newWorktreeName string, checkoutWorktreeBranchName string) (string, bool) {
	if !gwt.gitProcessLock.CanProceedWithGitOps() {
		gwt.logging.RegisterNewLog(logging.ADD_NEW_WORKTREE_OPS, "", logging.WARN, fmt.Sprintf("[WARN]: %s", gwt.gitProcessLock.OtherProcessRunningWarning()), false)
		return gwt.gitProcessLock.OtherProcessRunningWarning(), false
	}
	defer func() {
		gwt.gitProcessLock.ReleaseGitOpsLock()
	}()

	success := false
	newWorktreePath := filepath.Clean(filepath.Join(gwt.currentWorktreePath, fmt.Sprintf("../%s", newWorktreeName)))
	gitArgs := []string{"worktree", "add", newWorktreePath}
	if utf8.RuneCountInString(checkoutWorktreeBranchName) > 0 {
		gitArgs = append(gitArgs, checkoutWorktreeBranchName)
	}

	newWorktreeCmdExecutor := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, false)
	gwt.logging.RegisterNewLog(logging.ADD_NEW_WORKTREE_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	newWorktreeOutput, newWorktreeError := newWorktreeCmdExecutor.CombinedOutput()

	if newWorktreeError != nil {
		gwt.logging.RegisterNewLog(logging.ADD_NEW_WORKTREE_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.ADD_NEW_WORKTREE_OPS, newWorktreeError.Error()), true)
	} else {
		success = true
	}

	return string(newWorktreeOutput), success
}

// ------------------------------------
//
//	Remove the worktree at worktreePath via `git worktree remove`. No --force is
//	passed, so git refuses to remove a worktree with uncommitted/untracked changes
//	or a locked worktree, returning that error as the output.
//	Acquires the git ops lock first; returns the lock-busy warning if another git
//	process is running. Returns the command output and whether the remove succeeded.
//
// ------------------------------------
func (gwt *GitWorktree) RemoveWorktree(ctx context.Context, worktreePath string) (string, bool) {
	if !gwt.gitProcessLock.CanProceedWithGitOps() {
		gwt.logging.RegisterNewLog(logging.REMOVE_WORKTREE_OPS, "", logging.WARN, fmt.Sprintf("[WARN]: %s", gwt.gitProcessLock.OtherProcessRunningWarning()), false)
		return gwt.gitProcessLock.OtherProcessRunningWarning(), false
	}
	defer func() {
		gwt.gitProcessLock.ReleaseGitOpsLock()
	}()

	success := false
	gitArgs := []string{"worktree", "remove", worktreePath}
	removeWorktreeCmdExecutor := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, false)
	gwt.logging.RegisterNewLog(logging.REMOVE_WORKTREE_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	removeWorktreeOutput, removeWorktreeError := removeWorktreeCmdExecutor.CombinedOutput()

	if removeWorktreeError != nil {
		gwt.logging.RegisterNewLog(logging.REMOVE_WORKTREE_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.REMOVE_WORKTREE_OPS, removeWorktreeError.Error()), true)
	} else {
		success = true
	}

	return string(removeWorktreeOutput), success
}

// ------------------------------------
//
//	Prune all stale worktree admin entries in one pass via `git worktree prune`,
//	removing the bookkeeping under `.git/worktrees/*` for worktrees whose
//	working dir is gone or otherwise unusable. Valid worktrees are untouched.
//	Acquires the git ops lock first; no-ops if another git process is running.
//
// ------------------------------------
func (gwt *GitWorktree) PruneWorktrees() {
	if !gwt.gitProcessLock.CanProceedWithGitOps() {
		gwt.logging.RegisterNewLog(logging.PRUNE_WORKTREES_OPS, "", logging.WARN, fmt.Sprintf("[WARN]: %s", gwt.gitProcessLock.OtherProcessRunningWarning()), false)
		return
	}
	defer func() {
		gwt.gitProcessLock.ReleaseGitOpsLock()
	}()

	gitArgs := []string{"worktree", "prune"}
	pruneWorktreesCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	pruneWorktreesCmdExecutor.Run()
}

// ------------------------------------
//
//	Lock the worktree at worktreePath via `git worktree lock`, marking it so it is
//	never pruned or removed (useful for worktrees on removable or network media).
//	If lockReason is non-empty it is recorded with `--reason` for later inspection.
//	Acquires the git ops lock first; no-ops if another git process is running.
//
// ------------------------------------
func (gwt *GitWorktree) LockWorktree(lockReason string, worktreePath string) {
	if !gwt.gitProcessLock.CanProceedWithGitOps() {
		gwt.logging.RegisterNewLog(logging.LOCK_WORKTREE_OPS, "", logging.WARN, fmt.Sprintf("[WARN]: %s", gwt.gitProcessLock.OtherProcessRunningWarning()), false)
		return
	}
	defer func() {
		gwt.gitProcessLock.ReleaseGitOpsLock()
	}()
	gitArgs := []string{"worktree", "lock"}
	if utf8.RuneCountInString(lockReason) > 0 {
		gitArgs = append(gitArgs, "--reason", lockReason)
	}
	gitArgs = append(gitArgs, worktreePath)
	lockWorktreesCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	lockWorktreesCmdExecutor.Run()
}

// ------------------------------------
//
//	Unlock the worktree at worktreePath via `git worktree unlock`, reversing a prior
//	lock so it can be pruned or removed again.
//	Acquires the git ops lock first; no-ops if another git process is running.
//
// ------------------------------------
func (gwt *GitWorktree) UnlockWorktree(worktreePath string) {
	if !gwt.gitProcessLock.CanProceedWithGitOps() {
		gwt.logging.RegisterNewLog(logging.UNLOCK_WORKTREE_OPS, "", logging.WARN, fmt.Sprintf("[WARN]: %s", gwt.gitProcessLock.OtherProcessRunningWarning()), false)
		return
	}
	defer func() {
		gwt.gitProcessLock.ReleaseGitOpsLock()
	}()
	gitArgs := []string{"worktree", "unlock", worktreePath}
	unlockWorktreesCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	unlockWorktreesCmdExecutor.Run()
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
