package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle '[' key interaction.
//	Responsibility: Detail panel backward navigation.
//	Switches the active focus from DetailComponentPanelTwo back to DetailComponentPanel (the primary detail panel).
//
// ------------------------------------
func handleNonTypingLeftBracketKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		// handle detail component panel switching
		if m.CurrentSelectedComponent == constant.DetailComponentPanelTwo {
			m.CurrentSelectedComponent = constant.DetailComponentPanel
		}
	}
	return m, nil
}
