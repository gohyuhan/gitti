package git

import (
	"fmt"
	"strings"

	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/logging"
)

type GitStateUniversalUtils struct {
	currentGitState string
	gitProcessLock  *GitProcessLock
	GitAbsolutePath string
	logging         *logging.GittiLogging
}

// ------------------------------------
//
//	Initialize git state utility with repo path, process lock, and logger dependencies
//
// ------------------------------------
func InitGitStateUniversalUtils(gitAbsolutePath string, gitProcessLock *GitProcessLock, logging *logging.GittiLogging) *GitStateUniversalUtils {
	gitStateUniversalUtils := &GitStateUniversalUtils{
		currentGitState: "",
		gitProcessLock:  gitProcessLock,
		GitAbsolutePath: gitAbsolutePath,
		logging:         logging,
	}

	return gitStateUniversalUtils
}

// ------------------------------------
//
//	Return cached current git state string
//
// ------------------------------------
func (gSUU *GitStateUniversalUtils) GetCurrentGitState() string {
	return gSUU.currentGitState
}

// ------------------------------------
//
//	Detect active git state and execute corresponding continue command
//
// ------------------------------------
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
	err := continueCmdExecutor.Run()
	gSUU.logging.RegisterNewLog(logging.CONTINUE, strings.Join(gitArgs, " "), logging.INFO, "", true)
	if err != nil {
		gSUU.logging.RegisterNewLog(logging.CONTINUE, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.CONTINUE, err.Error()), true)
	}
}

// ------------------------------------
//
//	GitUniversalContinueWithSigning constructs a git continue command for terminal execution when signing is required.
//	When signing is enabled, gitti UI is suspended and the continue operation is executed directly in the terminal,
//	allowing the user to interact with the signing prompt (e.g., GPG passphrase).
//	This function detects the current git operation in progress and builds the appropriate continue command.
//
// ------------------------------------
func (gSUU *GitStateUniversalUtils) GitUniversalContinueWithSigning() []string {
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

	return gitArgs
}

// ------------------------------------
//
//	Detect active git state and execute corresponding abort command
//
// ------------------------------------
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
	err := abortCmdExecutor.Run()
	gSUU.logging.RegisterNewLog(logging.ABORT, strings.Join(gitArgs, " "), logging.INFO, "", true)
	if err != nil {
		gSUU.logging.RegisterNewLog(logging.ABORT, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.ABORT, err.Error()), true)
	}
}

// ------------------------------------
//
//	Detect active git state and execute corresponding skip command
//
// ------------------------------------
func (gSUU *GitStateUniversalUtils) GitUniversalSkip() {
	if !gSUU.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer func() {
		gSUU.gitProcessLock.ReleaseGitOpsLock()
	}()

	var gitArgs []string

	if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "rebase-apply/applying") {
		gitArgs = []string{"am", "--skip"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "rebase-merge") {
		gitArgs = []string{"rebase", "--skip"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "rebase-apply") {
		gitArgs = []string{"rebase", "--skip"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "CHERRY_PICK_HEAD") {
		gitArgs = []string{"cherry-pick", "--skip"}
	} else if checkIfFileExistWithinDotGitFolder(gSUU.GitAbsolutePath, "REVERT_HEAD") {
		gitArgs = []string{"revert", "--skip"}
	}

	if len(gitArgs) < 1 {
		return
	}

	skipCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	err := skipCmdExecutor.Run()
	gSUU.logging.RegisterNewLog(logging.SKIP, strings.Join(gitArgs, " "), logging.INFO, "", true)
	if err != nil {
		gSUU.logging.RegisterNewLog(logging.SKIP, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.SKIP, err.Error()), true)
	}
}

// ------------------------------------
//
//	Refresh cached git state by checking operation marker files under .git
//
// ------------------------------------
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
