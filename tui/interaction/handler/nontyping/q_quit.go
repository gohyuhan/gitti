package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'q' or 'Q' key interaction.
//	Responsibility: Global "quit" or "exit" operation.
//	If no popup is currently active, this forcefully stops any running git daemon
//	and triggers Bubble Tea's quit command, closing the application.
//
// ------------------------------------
func handleNonTypingqQKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if api.GITDAEMON != nil {
			api.GITDAEMON.Stop()
		}
		return m, tea.Quit
	}
	return m, nil
}
