package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/component/reflog"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	worktreePopUp "github.com/gohyuhan/gitti/tui/popup/worktree"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'n' key interaction.
//	Responsibility: Contextual "new" operation. Depending on the focused view:
//	- In Local Branch View: Opens popup to create a new branch (optionally based on a remote).
//	- In Remote View: Opens a prompt to add a new remote connection to the repository.
//
// ------------------------------------
func handleNonTypingnKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				m.PopUpType = constant.ChooseNewBranchTypePopUp
				m.IsTyping.Store(false)
				m.ShowPopUp.Store(true)
				if _, ok := m.PopUpModel.(*branchPopUp.ChooseNewBranchTypeOptionPopUpModel); !ok {
					branchPopUp.InitChooseNewBranchTypePopUpModel(m)
				}
			case constant.SHOW_REMOTE:
				m.PopUpType = constant.AddRemotePromptPopUp
				m.IsTyping.Store(true)
				m.ShowPopUp.Store(true)
				if m.GitOperations.GitRemote.CheckRemoteExist(false) {
					remotePopUp.InitAddRemotePromptPopUpModel(m, false)
				} else {
					remotePopUp.InitAddRemotePromptPopUpModel(m, true)
				}
			case constant.SHOW_WORKTREE:
				m.PopUpType = constant.WorktreeAddNewWorktreePopUp
				m.IsTyping.Store(true)
				m.ShowPopUp.Store(true)
				worktreePopUp.InitWorktreeAddNewWorktreePopUpModel(m)
			}
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_REFLOG:
				selectedReflog := m.CurrentRepoRefLogInfoList.SelectedItem()
				if selectedReflog != nil {
					parsedReflog := selectedReflog.(reflog.GitRefLogItem)
					m.PopUpType = constant.CreateNewBranchPopUp
					m.IsTyping.Store(true)
					m.ShowPopUp.Store(true)
					branchPopUp.InitCreateNewBranchPopUpModel(m, git.NEWBRANCHBASEDONCOMMITHASH, parsedReflog.Hash)
				}
			}
		}
	}
	return m, nil
}
