package stash

import "github.com/gohyuhan/gitti/tui/types"

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
