package services

import (
	"context"
	"os"

	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/initialize"
	worktreePopUp "github.com/gohyuhan/gitti/tui/popup/worktree"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	For Adding New Worktree
//
// ------------------------------------
func AddNewWorktreeService(m *types.GittiModel, newWorktreeName string, worktreeBranchName string) {
	ctx, cancel := context.WithCancel(context.Background())

	popUp, ok := m.PopUpModel.(*worktreePopUp.WorktreeAddNewWorktreeOutputPopUpModel)
	if ok {
		popUp.IsProcessing.Store(true)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsCancelled.Store(false)
		popUp.CancelFunc = cancel
	}

	go func(ctx context.Context) {
		defer cancel()

		outputResult, success := m.GitOperations.GitWorktree.AddNewWorktree(ctx, newWorktreeName, worktreeBranchName)
		data := types.WorktreeNewWorktreeResultEventDataStructure{
			Result:  outputResult,
			Success: success,
		}
		m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
			Event: constant.WORKTREE_NEW_WORKTREE_RESULT_EVENT,
			Data:  data,
		}
	}(ctx)
}

// ------------------------------------
//
//	Cancel the current add new worktree operation and clean up pop-up state
//
// ------------------------------------
func AddNewWorktreeCancelService(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*worktreePopUp.WorktreeAddNewWorktreeOutputPopUpModel)
	if ok {
		popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		if popUp.CancelFunc != nil {
			popUp.CancelFunc() // Cancel the context, which terminates the command and goroutine
		}
	}
	m.ShowPopUp.Store(false) // close the pop up
	m.IsTyping.Store(false)  // and reset typing mode
	m.PopUpType = constant.NoPopUp
	if ok {
		popUp.AddNewWorktreeOutputViewport.SetContent("") // set the add new worktree output viewport to nothing
		popUp.IsProcessing.Store(false)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
	}
}

// ------------------------------------
//
//	Prune stale worktree administrative files in the background
//
// ------------------------------------
func PruneWorktreesService(m *types.GittiModel) {
	go func() {
		m.GitOperations.GitWorktree.PruneWorktrees()
	}()
}

// ------------------------------------
//
//	Remove the worktree at worktreePath in the background
//
// ------------------------------------
func RemoveWorktreeService(m *types.GittiModel, worktreePath string) {
	go func() {
		m.GitOperations.GitWorktree.RemoveWorktree(worktreePath)
	}()
}

// ------------------------------------
//
//	Lock the worktree at worktreePath in the background, with an optional lock reason
//
// ------------------------------------
func LockWorktreeService(m *types.GittiModel, worktreePath string, lockReason string) {
	go func() {
		m.GitOperations.GitWorktree.LockWorktree(worktreePath, lockReason)
	}()
}

// ------------------------------------
//
//	Unlock the worktree at worktreePath in the background
//
// ------------------------------------
func UnlockWorktreeService(m *types.GittiModel, worktreePath string) {
	go func() {
		m.GitOperations.GitWorktree.UnlockWorktree(worktreePath)
	}()
}

// ------------------------------------
//
//	Switch the running app to the worktree at worktreePath. Because all git state
//	is path-relative, every path-derived dependency is re-pointed: the cmd
//	executor dir, the resolved repo path info, a freshly built GitOperations, the
//	GittiModel state (via ReinitGittiModel), and the daemon's GitOperations. A full
//	info fetch is then triggered so the new worktree's state repopulates at once.
//	On failure to resolve the target worktree, the cwd and executor dir are rolled
//	back to the original repo and the switch is aborted.
//
// ------------------------------------
func SwitchWorktreeService(m *types.GittiModel, worktreePath string) {
	// point the cwd and cmd executor at the target worktree before resolving its
	// path info (the executor pins cmd.Dir, so it must be updated, not just cwd)
	currentRepoPathBeforeSwitch, _ := os.Getwd()
	os.Chdir(worktreePath)
	executor.GittiCmdExecutor.UpdateRepoPath(worktreePath)

	gitRepoPathInfo, gitRepoInfoErr := api.GetGitPathInfo()
	if gitRepoInfoErr != nil {
		// resolve failed, roll back cwd + executor dir to the original repo and abort
		os.Chdir(currentRepoPathBeforeSwitch)
		executor.GittiCmdExecutor.UpdateRepoPath(m.RepoPath)
		m.GittiLogger.RegisterNewLog(logging.SWITCH_WORKTREE_OPS, "", logging.ERROR, "[ERROR] fail to retrieve worktree info for switching", false)
		return
	}

	// re-point the executor at the resolved top level and rebuild GitOperations
	// against the new worktree's git/worktree paths
	executor.GittiCmdExecutor.UpdateRepoPath(gitRepoPathInfo.TopLevelRepoPath)

	gitOperations := api.InitGitOperations(gitRepoPathInfo.AbsoluteGitRepoPath, gitRepoPathInfo.AbsoluteWorktreePath, m.GitUpdateChannel, m.GittiLogger)

	// reset all model state in place (preserving terminal width/height) for the new worktree
	initialize.ReinitGittiModel(m, gitRepoPathInfo.TopLevelRepoPath, gitRepoPathInfo.RepoName, gitOperations)

	// rewire the daemon to the freshly rebuilt GitOperations, then trigger one full
	// fetch so the new worktree's git state repopulates immediately
	if api.GITDAEMON != nil {
		api.GITDAEMON.UpdateGitOperations(gitOperations)
		api.GITDAEMON.TriggerFullInfoFetch()
	}
}
