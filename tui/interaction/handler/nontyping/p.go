package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	pushPopUp "github.com/gohyuhan/gitti/tui/popup/push"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'p' key interaction.
//	Responsibility: Initiates the "push" operation.
//	Checks for existing remotes; if none exist, prompts to add one.
//	If remotes exist, it figures out the appropriate remote to push to,
//	often opening a popup for the user to select the push argument (e.g., force push)
//	or the specific remote if multiple exist.
//
// ------------------------------------
func handleNonTypingpKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		// first we need to check if there are any push/pull origin origin for this repo
		// if not we prompt the user to add a new remote origin
		if !m.GitOperations.GitRemote.CheckRemoteExist(false) {
			showAddRemotePromptPopUp(m)
		} else {
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
			remotes := m.GitOperations.GitRemote.PushRemote()
			if len(remotes) == 1 {
				m.PopUpType = constant.ChoosePushTypePopUp
				// if the current pop up model is not commit pop up model, then init it and start git push service
				pushPopUp.InitChoosePushTypePopUpModel(m, remotes[0].Name)
			} else if len(remotes) > 1 {
				// if remote is more than 1 let user choose which remote to push to first before pushing
				m.PopUpType = constant.ChooseRemotePopUp
				if _, ok := m.PopUpModel.(*remotePopUp.ChooseRemotePopUpModel); !ok {
					remotePopUp.InitChooseRemotePopUpModel(m, remotes, constant.PUSHACTION)
				}
			}
		}
	}
	return m, nil
}
