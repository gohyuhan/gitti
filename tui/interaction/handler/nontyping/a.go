package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'a' key interaction.
//	Responsibility: Used contextually within popups. Specifically handles 'apply' actions
//	for git cherry-pick popups, shifting the model state to confirm the application
//	of cherry-picked commits.
//
// ------------------------------------
func handleNonTypingaKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.GitCherryPickPopUp, constant.GitEditCherryPickPopUp:
			m.ShowPopUp.Store(true)
			m.PopUpType = constant.GitCherryPickApplyConfirmPopUp
			m.PopUpModel = nil // we don't need to initialize the pop up model, as we are just showing the pop up and we don't need to hold any state or info
			m.IsTyping.Store(false)
		}
	}
	return m, nil
}
