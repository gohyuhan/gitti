package services

import (
	"context"

	"github.com/gohyuhan/gitti/tui/constant"
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
