package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/component/reflog"
	"github.com/gohyuhan/gitti/tui/component/tag"
	"github.com/gohyuhan/gitti/tui/constant"
	commitLogPopUp "github.com/gohyuhan/gitti/tui/popup/commitlog"
	reflogPopUp "github.com/gohyuhan/gitti/tui/popup/reflog"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle Ctrl+p key interaction.
//	Responsibility: Contextual "cherry-pick mode" management.
//	- Active Cherry-pick Popup: Switches between simple cherry-pick view and apply-confirm view.
//	- Log Panel / Branch Panel: Triggers opening the initial cherry-pick interface or option
//	  selection menu if commits are already queued.
//
// ------------------------------------
func handleNonTypingCtrlpKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.GitEditCherryPickPopUp, constant.GitCherryPickApplyConfirmPopUp:
			m.PopUpType = constant.GitCherryPickPopUp
			commitLogPopUp.InitGitCherryPickPopUpModel(m, m.CheckOutBranch)
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
		}
	} else {
		if (m.CurrentSelectedComponent == constant.CommitLogOrRefLogComponentPanel || m.DetailPanelParentComponent == constant.CommitLogOrRefLogComponentPanel) &&
			len(m.CurrentRepoCommitLogInfoList.Items()) > 0 && m.CurrentCommitLogOrRefLogComponentShowing == constant.SHOW_COMMITLOG {
			if len(m.CherryPickedCommitInfo.CherryPickedCommitMap) < 1 {
				m.ShowPopUp.Store(true)
				m.PopUpType = constant.GitCherryPickPopUp
				m.IsTyping.Store(false)
				commitLogPopUp.InitGitCherryPickPopUpModel(m, m.CheckOutBranch)
			} else {
				m.ShowPopUp.Store(true)
				m.PopUpType = constant.GitCherryPickOptionSelectionPopUp
				m.IsTyping.Store(false)
				commitLogPopUp.InitGitCherryPickOptionSelectionPopUpModel(m)
			}
		} else if (m.CurrentSelectedComponent == constant.CommitLogOrRefLogComponentPanel || m.DetailPanelParentComponent == constant.CommitLogOrRefLogComponentPanel) &&
			len(m.CurrentRepoRefLogInfoList.Items()) > 0 && m.CurrentCommitLogOrRefLogComponentShowing == constant.SHOW_REFLOG {
			selectedReflog := m.CurrentRepoRefLogInfoList.SelectedItem()
			if selectedReflog != nil {
				parsedSelectedReflog := selectedReflog.(reflog.GitRefLogItem)
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				m.PopUpType = constant.GitCherryPickFromRefLogApplyConfirmationPopUp
				reflogPopUp.InitGitCherryPickFromRefLogApplyConfirmationPopUpModel(m, parsedSelectedReflog.Hash, parsedSelectedReflog.Head, parsedSelectedReflog.Action, parsedSelectedReflog.ActionInfo)
			}
		} else if (m.CurrentSelectedComponent == constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel || m.DetailPanelParentComponent == constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel) &&
			len(m.CurrentRepoTagInfoList.Items()) > 0 && m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing == constant.SHOW_TAG {
			selectedTag := m.CurrentRepoTagInfoList.SelectedItem()
			if selectedTag == nil {
				return m, nil
			}
			if !m.GitOperations.GitRemote.CheckRemoteExist(false) {
				// if no remote found, we add one
				showAddRemotePromptPopUp(m)
			} else {
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				remotes := m.GitOperations.GitRemote.PushRemote()
				if len(remotes) == 1 {
					m.PopUpType = constant.ChoosePushTagOptionPopUp
					// only one remote found so, we will default to that remote
					tagPopUp.InitChoosePushTagOptionPopUpModel(m, remotes[0].Name, selectedTag.(tag.GitTagItem).TagName)
				} else if len(remotes) > 1 {
					// if remote is more than 1 let user choose which remote
					m.PopUpType = constant.ChooseRemotePopUp
					remotePopUp.InitChooseRemotePopUpModel(m, remotes, constant.TAGPUSHACTION)
				}
			}
		} else if (m.CurrentSelectedComponent == constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel || m.DetailPanelParentComponent == constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel) &&
			len(m.CurrentRepoWorktreeInfoList.Items()) > 0 && m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing == constant.SHOW_WORKTREE {
			services.PruneWorktreesService(m)
		}
	}
	return m, nil
}
