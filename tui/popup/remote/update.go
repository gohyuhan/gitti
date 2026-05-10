package remote

import (
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle the async add-remote result event. Writes the result lines to the
//	output viewport, clears IsProcessing, and either resets both inputs with
//	ProcessSuccess on success, or sets HasError on failure. No-ops if the popup
//	is not the active add-remote popup or the operation was cancelled.
//
// ------------------------------------
func UpdateAddRemoteResultEvent(m *types.GittiModel, updateData types.GitAddRemoteResultEventDataStructure) {
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
