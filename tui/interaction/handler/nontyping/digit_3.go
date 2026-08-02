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
//	Handle '3' key interaction.
//	Responsibility: Switches the active view focus to the third main layout tab
//	(Commit Log Component Panel), refreshing the detail view to display git log history.
//
// ------------------------------------
func handleNonTyping3KeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent != constant.CommitLogOrRefLogComponentPanel {
			m.CurrentSelectedComponent = constant.CommitLogOrRefLogComponentPanel
			m.CurrentSelectedComponentIndex = 3
			m.DetailPanelParentComponent = ""
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	}
	return m, nil
}
