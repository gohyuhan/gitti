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
//	Handle '2' key interaction.
//	Responsibility: Switches the active view focus to the second main layout tab
//	(Modified Files Component Panel), triggering an update to the detail panel
//	to reflect currently modified, added, or deleted repository files.
//
// ------------------------------------
func handleNonTyping2KeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent != constant.ModifiedFilesComponentPanel {
			m.CurrentSelectedComponent = constant.ModifiedFilesComponentPanel
			m.CurrentSelectedComponentIndex = 2
			m.DetailPanelParentComponent = ""
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	}
	return m, nil
}
