package nontyping

import (
	"github.com/gohyuhan/gitti/tui/constant"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Open the Add Remote prompt pop up.
//	If the current pop up model is not the add remote prompt pop up model, then init it.
//	If it is, reuse it and clear its output viewport.
//
// ------------------------------------
func showAddRemotePromptPopUp(m *types.GittiModel) {
	m.PopUpType = constant.AddRemotePromptPopUp
	if popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel); !ok {
		remotePopUp.InitAddRemotePromptPopUpModel(m, true)
	} else {
		popUp.AddRemoteOutputViewport.SetContent("")
	}
	m.ShowPopUp.Store(true)
	m.IsTyping.Store(true)
}
