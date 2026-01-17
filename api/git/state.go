package git

import (
	"github.com/gohyuhan/gitti/executor"
)

type GitStateUniversalUtils struct {
	currentGitState string
	errorLog        []error
	gitProcessLock  *GitProcessLock
	GitAbsolutePath string
}

// --------------------------------
//
// init the Git GitStateUniversalUtils
//
// --------------------------------
func InitGitStateUniversalUtils(gitAbsolutePath string, gitProcessLock *GitProcessLock) *GitStateUniversalUtils {
	gitStateUniversalUtils := &GitStateUniversalUtils{
		currentGitState: "",
		errorLog:        []error{},
		gitProcessLock:  gitProcessLock,
		GitAbsolutePath: gitAbsolutePath,
	}

	return gitStateUniversalUtils
}

// --------------------------------
//
// GetCurrentGitState returns the current Git state as a string.
//
// --------------------------------
func (gSUU *GitStateUniversalUtils) GetCurrentGitState() string {
	return gSUU.currentGitState
}

// --------------------------------
//
// GitUniversalContinue attempts to continue the current Git operation (rebase, merge, am, etc.)
// by detecting the state files within the .git folder.
//
// --------------------------------
func (gSUU *GitStateUniversalUtils) GitUniversalContinue() {
	if !gSUU.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer func() {
		gSUU.gitProcessLock.ReleaseGitOpsLock()
	}()

	var gitArgs []string

	if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "rebase-apply/applying") {
		gitArgs = []string{"am", "--continue"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "rebase-merge") {
		gitArgs = []string{"rebase", "--continue"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "rebase-apply") {
		gitArgs = []string{"rebase", "--continue"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "MERGE_HEAD") {
		gitArgs = []string{"merge", "--continue"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "CHERRY_PICK_HEAD") {
		gitArgs = []string{"cherry-pick", "--continue"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "REVERT_HEAD") {
		gitArgs = []string{"revert", "--continue"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "NOTES_MERGE_WORKTREE") {
		gitArgs = []string{"notes", "merge", "--commit"}
	}

	if len(gitArgs) < 1 {
		return
	}

	continueCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	continueCmdExecutor.Run()

	return
}

// --------------------------------
//
// GitUniversalAbort aborts the current Git operation (rebase, merge, am, etc.)
// by detecting the state files within the .git folder.
//
// --------------------------------
func (gSUU *GitStateUniversalUtils) GitUniversalAbort() {
	if !gSUU.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer func() {
		gSUU.gitProcessLock.ReleaseGitOpsLock()
	}()

	var gitArgs []string

	if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "rebase-apply/applying") {
		gitArgs = []string{"am", "--abort"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "rebase-merge") {
		gitArgs = []string{"rebase", "--abort"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "rebase-apply") {
		gitArgs = []string{"rebase", "--abort"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "MERGE_HEAD") {
		gitArgs = []string{"merge", "--abort"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "CHERRY_PICK_HEAD") {
		gitArgs = []string{"cherry-pick", "--abort"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "REVERT_HEAD") {
		gitArgs = []string{"revert", "--abort"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "NOTES_MERGE_WORKTREE") {
		gitArgs = []string{"notes", "merge", "--abort"}
	}

	if len(gitArgs) < 1 {
		return
	}

	abortCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	abortCmdExecutor.Run()

	return
}

// --------------------------------
//
// CheckCurrentGitState updates the currentGitState field by checking the existence of
// specific state files inside the .git directory.
//
// --------------------------------
func (gSUU *GitStateUniversalUtils) CheckCurrentGitState() {
	gSUU.currentGitState = ""
	if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "rebase-apply/applying") {
		gSUU.currentGitState = "AM"
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "rebase-merge") {
		gSUU.currentGitState = "REBASE"
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "rebase-apply") {
		gSUU.currentGitState = "REBASE"
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "MERGE_HEAD") {
		gSUU.currentGitState = "MERGE"
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "CHERRY_PICK_HEAD") {
		gSUU.currentGitState = "CHERRY-PICK"
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "REVERT_HEAD") {
		gSUU.currentGitState = "REVERT"
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "NOTES_MERGE_WORKTREE") {
		gSUU.currentGitState = "NOTES-MERGE"
	}
}
