package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'A' key interaction.
//	Responsibility: Initiates the "Amend Commit" flow. Opens the amend commit popup,
//	clears previous commit output, initializes the input model for the message,
//	and switches the application to typing state.
//
// ------------------------------------
func handleNonTypingAKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		m.ShowPopUp.Store(true)
		m.PopUpType = constant.AmendCommitPopUp
		m.GitOperations.GitCommit.ClearGitCommitOutput()

		commitPopUp.InitGitAmendCommitPopUpModel(m)

		m.IsTyping.Store(true)
	}
	return m, nil
}
