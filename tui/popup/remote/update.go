package remote

import (
	"github.com/gohyuhan/gitti/tui/types"
)

func UpdateAddRemoteResultEvent(m *types.GittiModel, updateData types.GitAddRemoteResultEventDataInterface) {
	popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
	if !ok || popUp.IsCancelled.Load() {
		return
	}

	popUp.IsProcessing.Store(false)
	popUp.AddRemoteOutputViewport.SetContentLines(updateData.Result)
	popUp.AddRemoteOutputViewport.PageDown()

	if updateData.Success {
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(true)
		popUp.NoInitialRemote = false
		popUp.RemoteNameTextInput.Reset()
		popUp.RemoteUrlTextInput.Reset()
		return
	}

	popUp.HasError.Store(true)
	popUp.ProcessSuccess.Store(false)
}
