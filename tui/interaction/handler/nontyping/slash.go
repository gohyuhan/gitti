package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle '/' key interaction.
//	Responsibility: Focuses the Log Component Panel.
//	Used as a quick jumping shortcut to view the internal console/error logs of Gitti.
//
// ------------------------------------
func handleNonTypingSlashKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent != constant.LogComponentPanel {
			m.CurrentSelectedComponent = constant.LogComponentPanel
			m.DetailPanelParentComponent = ""
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	}
	return m, nil
}
