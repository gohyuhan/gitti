package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/layout"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'L' key interaction.
//	Responsibility: Enters "Line Editing Mode" (also known as hunk staging mode).
//	Allows the user to specifically stage, unstage or discard individual lines of a modified file
//	instead of the entire file, updating the detail panel UI to reflect this state.
//
// ------------------------------------
func handleNonTypingLKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	// enter line stage state
	layout.EnterOrReinitLineEditingState(m)
	m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
		Event: constant.DETAIL_COMPONENT_PANEL_LAYOUT_UPDATED_EVENT,
	}
	return m, nil
}
