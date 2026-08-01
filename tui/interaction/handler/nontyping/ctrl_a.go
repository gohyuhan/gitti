package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle Ctrl+a key interaction.
//	Responsibility: Global "Abort" operation.
//	Attempts to actively abort ongoing multi-step git state operations
//	like an incomplete merge, rebase, or cherry-pick.
//
// ------------------------------------
func handleNonTypingCtrlaKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		services.GitStateUniversalUtilsAbortService(m)
	}
	return m, nil
}
