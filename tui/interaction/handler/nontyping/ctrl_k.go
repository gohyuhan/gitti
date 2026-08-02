package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle Ctrl+k key interaction.
//	Responsibility: Global "Skip" operation.
//	Used primarily to skip the current commit during an interactive rebase
//	or similar multi-step git states where skipping is a valid option.
//
// ------------------------------------
func handleNonTypingCtrlkKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		services.GitStateUniversalUtilsSkipService(m)
	}
	return m, nil
}
