package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	Handle 'F' key interaction.
//	Responsibility: Enters panel list filter typing mode for the currently focused
//	list panel (branch/tag/remote/worktree, modified files, commit log/reflog, stash).
//	While active, key presses edit the panel's filter query instead of triggering
//	their normal actions (enter keeps the query, esc clears it).
//
// ------------------------------------
func handleNonTypingFKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() && !m.IsLineEditingState.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel,
			constant.ModifiedFilesComponentPanel,
			constant.CommitLogOrRefLogComponentPanel,
			constant.StashComponentPanel:
			if m.PanelFilterQuery == nil {
				m.PanelFilterQuery = make(map[string]string)
			}
			m.PanelFilterInput.SetValue(m.PanelFilterQuery[utils.CurrentPanelFilterKey(m)])
			m.PanelFilterInput.CursorEnd()
			m.IsPanelFiltering.Store(true)
		}
	}
	return m, nil
}
