package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	blamePopUp "github.com/gohyuhan/gitti/tui/popup/blame"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'b' key interaction.
//	Responsibility: Opens the git blame popup for the currently focused file;
//	sets IsTyping to true so subsequent key events are routed to the popup.
//
// ------------------------------------
func handleNonTypingbKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		m.ShowPopUp.Store(true)
		m.PopUpType = constant.BlamePopUp

		blamePopUp.InitBlamePopUpModel(m)

		m.IsTyping.Store(true)
	}
	return m, nil
}
