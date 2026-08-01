package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/layout"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle '1' key interaction.
//	Responsibility: Switches the active view focus to the first main layout tab
//	(Local Branch, Tag, or Remote Component Panel), triggering an update to the
//	detail panel dynamically.
//
// ------------------------------------
func handleNonTyping1KeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent != constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel {
			m.CurrentSelectedComponent = constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel
			m.CurrentSelectedComponentIndex = 1
			m.DetailPanelParentComponent = ""
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	}
	return m, nil
}
