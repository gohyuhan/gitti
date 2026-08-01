package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/component/stash"
	"github.com/gohyuhan/gitti/tui/component/worktree"
	"github.com/gohyuhan/gitti/tui/constant"
	commitLogPopUp "github.com/gohyuhan/gitti/tui/popup/commitlog"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	worktreePopUp "github.com/gohyuhan/gitti/tui/popup/worktree"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	Handle Backspace key interaction.
//	Responsibility: Contextual "pop" or "revert" action.
//	- In Stash Panel: Triggers a popup to 'pop' (apply and drop) the selected stash.
//	- In Edit Cherry Pick Popup: Clears the cherry-picked commit list and dismisses the popup.
//
// ------------------------------------
func handleNonTypingBackspaceKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.StashComponentPanel:
			selectedStashId := m.CurrentRepoStashInfoList.SelectedItem()
			if selectedStashId != nil {
				stashPopUp.InitGitStashConfirmPromptPopUpModel(m, git.POPSTASH, "", selectedStashId.(stash.GitStashItem).Id, selectedStashId.(stash.GitStashItem).Message)
				m.PopUpType = constant.GitStashConfirmPromptPopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
			}
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
			case constant.SHOW_WORKTREE:
				selectedWorktreeItem := m.CurrentRepoWorktreeInfoList.SelectedItem()
				if selectedWorktreeItem != nil {
					worktreeItem := selectedWorktreeItem.(worktree.GitWorktreeItem)
					m.PopUpType = constant.WorktreeRemoveWorktreeConfirmationPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					worktreePopUp.InitWorktreeRemoveWorktreeConfirmationPopUpModel(m, worktreeItem.WorktreePath)
				}
			}
		}
	} else if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.GitEditCherryPickPopUp:
			_, ok := m.PopUpModel.(*commitLogPopUp.GitEditCherryPickPopUpModel)
			if ok {
				utils.ReinitCherryPickedCommitInfo(m)
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		}
	}
	return m, nil
}
