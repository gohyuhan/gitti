package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/component/worktree"
	"github.com/gohyuhan/gitti/tui/constant"
	worktreePopUp "github.com/gohyuhan/gitti/tui/popup/worktree"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

func handleNonTypingoKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
			case constant.SHOW_WORKTREE:
				selectedWorktreeItem := m.CurrentRepoWorktreeInfoList.SelectedItem()
				if selectedWorktreeItem != nil {
					parsedSelectedWorktree := selectedWorktreeItem.(worktree.GitWorktreeItem)
					// main worktree cannot be lock and unlock, silent return
					if parsedSelectedWorktree.IsMain {
						return m, nil
					}
					if parsedSelectedWorktree.IsLocked {
						services.UnlockWorktreeService(m, parsedSelectedWorktree.WorktreePath)
					} else {
						m.PopUpType = constant.WorktreeLockReasonInputPopUp
						m.IsTyping.Store(true)
						m.ShowPopUp.Store(true)
						worktreePopUp.InitWorktreeLockReasonInputPopUpModel(m, parsedSelectedWorktree.WorktreePath)
					}
				}
			}
		}
	}
	return m, nil
}
