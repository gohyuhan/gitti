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
//	Handle '>' key interaction.
//	Responsibility: Component sub-state cycling (Forward).
//	- No popup: In the 'Local Branch Or Tag Or Remote' panel, rotates the visible list to the
//	  right/next entity (e.g., from Branches -> Tags -> Remotes). In the Commit Log panel,
//	  switches from Commit Log to Reflog view.
//	- Interactive Rebase Fixup/Squash Popup: Shifts focus to the detail viewport pane.
//
// ------------------------------------
func handleNonTypingRightAngleBracketKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing = constant.SHOW_TAG
				services.FetchDetailComponentPanelInfoService(m, true)
			case constant.SHOW_TAG:
				m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing = constant.SHOW_REMOTE
				services.FetchDetailComponentPanelInfoService(m, true)
			case constant.SHOW_REMOTE:
				m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing = constant.SHOW_WORKTREE
				services.FetchDetailComponentPanelInfoService(m, true)
			case constant.SHOW_WORKTREE:
				// do nothing, as worktree is currently the most right option in the local branch or tag or remote or worktree component panel
			}
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
				m.CurrentCommitLogOrRefLogComponentShowing = constant.SHOW_REFLOG
				services.FetchDetailComponentPanelInfoService(m, true)
			case constant.SHOW_REFLOG:
				// do nothing, as reflog is currently the most right option in the commit log or reflog component panel
			}
		}
	} else {
		switch m.PopUpType {
		case constant.InteractiveRebaseFixupSquashSelectionPopUp:
			popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashSelectionPopUpModel)
			if ok {
				popUp.IsCommitListSelected = false
				popUp.IsCommitFixupSquashViewportSelected = true
			}
		}
	}
	return m, nil
}
