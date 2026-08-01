package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle '<' key interaction.
//	Responsibility: Component sub-state cycling (Backward).
//	- No popup: In the 'Local Branch Or Tag Or Remote' panel, rotates the visible list to the
//	  left/previous entity (e.g., from Remotes -> Tags -> Branches). In the Commit Log panel,
//	  switches from Reflog back to Commit Log view.
//	- Interactive Rebase Fixup/Squash Popup: Shifts focus to the commit list pane.
//
// ------------------------------------
func handleNonTypingLeftAngleBracketKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
			// do nothing, as local branch will be the most left option in the local branch or tag or remote or worktree component panel
			case constant.SHOW_TAG:
				m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing = constant.SHOW_LOCAL_BRANCH
				services.FetchDetailComponentPanelInfoService(m, true)
			case constant.SHOW_REMOTE:
				m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing = constant.SHOW_TAG
				services.FetchDetailComponentPanelInfoService(m, true)
			case constant.SHOW_WORKTREE:
				m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing = constant.SHOW_REMOTE
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
			// do nothing, as commit log will be the most left option in the commit log or reflog component panel
			case constant.SHOW_REFLOG:
				m.CurrentCommitLogOrRefLogComponentShowing = constant.SHOW_COMMITLOG
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		}
	} else {
		switch m.PopUpType {
		case constant.InteractiveRebaseFixupSquashSelectionPopUp:
			popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashSelectionPopUpModel)
			if ok {
				popUp.IsCommitListSelected = true
				popUp.IsCommitFixupSquashViewportSelected = false
			}
		}
	}
	return m, nil
}
