package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	pullPopUp "github.com/gohyuhan/gitti/tui/popup/pull"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'P' key interaction.
//	Responsibility: Initiates the "pull" operation.
//	Checks for existing remotes; if none exist, prompts to add one.
//	If remotes exist, opens a popup allowing the user to select the type
//	of pull operation (e.g., normal pull, pull rebase).
//
// ------------------------------------
func handleNonTypingPKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		// first we need to check if there are any push/pull origin for this repo
		// if not we prompt the user to add a new remote origin
		if !m.GitOperations.GitRemote.CheckRemoteExist(false) {
			showAddRemotePromptPopUp(m)
		} else {
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
			m.PopUpType = constant.ChooseGitPullTypePopUp
			pullPopUp.InitChooseGitPullTypePopUpModel(m)
		}
	}
	return m, nil
}
