package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle ']' key interaction.
//	Responsibility: Detail panel forward navigation.
//	Switches the active focus from DetailComponentPanel to DetailComponentPanelTwo (the secondary detail panel), if it is currently visible.
//
// ------------------------------------
func handleNonTypingRightBracketKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		// handle detail component panel switching
		if m.CurrentSelectedComponent == constant.DetailComponentPanel && m.ShowDetailPanelTwo.Load() {
			m.CurrentSelectedComponent = constant.DetailComponentPanelTwo
		}
	}
	return m, nil
}
