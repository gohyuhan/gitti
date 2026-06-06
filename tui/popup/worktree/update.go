package worktree

import "github.com/gohyuhan/gitti/tui/types"

// ------------------------------------
//
//	Write the add-new-worktree command result into the output popup viewport and
//	flip the processing/success/error state flags accordingly. No-op if the popup
//	was cancelled.
//
// ------------------------------------
func UpdateWorktreeAddNewWorktreeOutputViewPort(m *types.GittiModel, updateData types.WorktreeNewWorktreeResultEventDataStructure) {
	popUp, ok := m.PopUpModel.(*WorktreeAddNewWorktreeOutputPopUpModel)
	if ok && !popUp.IsCancelled.Load() {
		popUp.AddNewWorktreeOutputViewport.SetContentLines(updateData.Result)
		popUp.IsProcessing.Store(false)
		if updateData.Success {
			popUp.ProcessSuccess.Store(true)
			popUp.HasError.Store(false)
		} else {
			popUp.ProcessSuccess.Store(false)
			popUp.HasError.Store(true)
		}
	}
}
