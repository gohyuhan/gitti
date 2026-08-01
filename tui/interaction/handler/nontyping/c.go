package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'c' key interaction.
//	Responsibility: Initiates the "New Commit" flow. Opens the standard commit popup,
//	resets the commit output view, initializes the input model if necessary,
//	and places the user in typing mode to enter commit details.
//
// ------------------------------------
func handleNonTypingcKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		m.GitOperations.GitCommit.ClearGitCommitOutput()
		// if the current pop up model is not commit pop up model, then init it
		if popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel); !ok {
			commitPopUp.InitGitCommitPopUpModel(m)
		} else {
			popUp.InitialCommitStarted.Store(false)
			popUp.GitCommitOutputViewport.SetContent("")
		}
		m.PopUpType = constant.CommitPopUp
		m.ShowPopUp.Store(true)
		m.IsTyping.Store(true)
	}
	return m, nil
}
