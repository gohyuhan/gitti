package stash

import "github.com/gohyuhan/gitti/tui/types"

// ------------------------------------
//
//	Handle the async stash operation result event. Clears IsProcessing, then sets
//	ProcessSuccess on success or HasError on failure, and loads the output lines
//	into the viewport.
//
// ------------------------------------
func UpdateGitStashOperationResultEvent(m *types.GittiModel, updateData types.GitStashOperationResultEventDataInterface) {
	popUp, ok := m.PopUpModel.(*GitStashOperationOutputPopUpModel)
	if ok {
		popUp.IsProcessing.Store(false)
		if updateData.Success && !popUp.IsProcessing.Load() {
			popUp.ProcessSuccess.Store(true)
			popUp.HasError.Store(false)
		} else if !updateData.Success && !popUp.IsProcessing.Load() {
			popUp.HasError.Store(true)
		}
		popUp.GitStashOperationOutputViewport.SetContentLines(updateData.Result)
	}
}
