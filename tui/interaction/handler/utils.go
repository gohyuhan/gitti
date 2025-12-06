package handler

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	discardPopUp "github.com/gohyuhan/gitti/tui/popup/discard"
	keybindingPopUp "github.com/gohyuhan/gitti/tui/popup/keybinding"
	pullPopUp "github.com/gohyuhan/gitti/tui/popup/pull"
	pushPopUp "github.com/gohyuhan/gitti/tui/popup/push"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	resolvePopUp "github.com/gohyuhan/gitti/tui/popup/resolve"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	"github.com/gohyuhan/gitti/tui/types"
)

func UpDownKeyMsgUpdateForPopUp(msg tea.KeyMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	// for within pop up component
	switch m.PopUpType {
	// following is for list component
	case constant.ChooseRemotePopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.ChooseRemotePopUpModel)
		if ok {
			popUp.RemoteList, cmd = popUp.RemoteList.Update(msg)
			return m, cmd
		}
	case constant.ChoosePushTypePopUp:
		popUp, ok := m.PopUpModel.(*pushPopUp.ChoosePushTypePopUpModel)
		if ok {
			popUp.PushOptionList, cmd = popUp.PushOptionList.Update(msg)
			return m, cmd
		}
	case constant.ChooseNewBranchTypePopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.ChooseNewBranchTypeOptionPopUpModel)
		if ok {
			popUp.NewBranchTypeOptionList, cmd = popUp.NewBranchTypeOptionList.Update(msg)
			return m, cmd
		}
	case constant.ChooseSwitchBranchTypePopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.ChooseSwitchBranchTypePopUpModel)
		if ok {
			popUp.SwitchTypeOptionList, cmd = popUp.SwitchTypeOptionList.Update(msg)
			return m, cmd
		}
	case constant.ChooseGitPullTypePopUp:
		popUp, ok := m.PopUpModel.(*pullPopUp.ChooseGitPullTypePopUpModel)
		if ok {
			popUp.PullTypeOptionList, cmd = popUp.PullTypeOptionList.Update(msg)
			return m, cmd
		}
	case constant.GitDiscardTypeOptionPopUp:
		popUp, ok := m.PopUpModel.(*discardPopUp.GitDiscardTypeOptionPopUpModel)
		if ok {
			popUp.DiscardTypeOptionList, cmd = popUp.DiscardTypeOptionList.Update(msg)
			return m, cmd
		}
	case constant.GitResolveConflictOptionPopUp:
		popUp, ok := m.PopUpModel.(*resolvePopUp.GitResolveConflictOptionPopUpModel)
		if ok {
			popUp.ResolveConflictOptionList, cmd = popUp.ResolveConflictOptionList.Update(msg)
			return m, cmd
		}

	// following is for viewport
	case constant.GlobalKeyBindingPopUp:
		popUp, ok := m.PopUpModel.(*keybindingPopUp.GlobalKeyBindingPopUpModel)
		if ok {
			popUp.GlobalKeyBindingViewport, cmd = popUp.GlobalKeyBindingViewport.Update(msg)
			return m, cmd
		}
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			popUp.GitCommitOutputViewport, cmd = popUp.GitCommitOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
		if ok {
			popUp.GitAmendCommitOutputViewport, cmd = popUp.GitAmendCommitOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.GitRemotePushPopUp:
		popUp, ok := m.PopUpModel.(*pushPopUp.GitRemotePushPopUpModel)
		if ok {
			popUp.GitRemotePushOutputViewport, cmd = popUp.GitRemotePushOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.GitPullOutputPopUp:
		popUp, ok := m.PopUpModel.(*pullPopUp.GitPullOutputPopUpModel)
		if ok {
			popUp.GitPullOutputViewport, cmd = popUp.GitPullOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.SwitchBranchOutputPopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.SwitchBranchOutputPopUpModel)
		if ok {
			popUp.SwitchBranchOutputViewport, cmd = popUp.SwitchBranchOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel)
		if ok {
			popUp.AddRemoteOutputViewport, cmd = popUp.AddRemoteOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.GitStashOperationOutputPopUp:
		popUp, ok := m.PopUpModel.(*stashPopUp.GitStashOperationOutputPopUpModel)
		if ok {
			popUp.GitStashOperationOutputViewport, cmd = popUp.GitStashOperationOutputViewport.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func UpDownMouseMsgUpdateForPopUp(msg tea.MouseMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	// for pop up that have viewport
	switch m.PopUpType {
	case constant.GlobalKeyBindingPopUp:
		popUp, ok := m.PopUpModel.(*keybindingPopUp.GlobalKeyBindingPopUpModel)
		if ok {
			popUp.GlobalKeyBindingViewport, cmd = popUp.GlobalKeyBindingViewport.Update(msg)
			return m, cmd
		}

	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			popUp.GitCommitOutputViewport, cmd = popUp.GitCommitOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
		if ok {
			popUp.GitAmendCommitOutputViewport, cmd = popUp.GitAmendCommitOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.GitRemotePushPopUp:
		popUp, ok := m.PopUpModel.(*pushPopUp.GitRemotePushPopUpModel)
		if ok {
			popUp.GitRemotePushOutputViewport, cmd = popUp.GitRemotePushOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.GitPullOutputPopUp:
		popUp, ok := m.PopUpModel.(*pullPopUp.GitPullOutputPopUpModel)
		if ok {
			popUp.GitPullOutputViewport, cmd = popUp.GitPullOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.SwitchBranchOutputPopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.SwitchBranchOutputPopUpModel)
		if ok {
			popUp.SwitchBranchOutputViewport, cmd = popUp.SwitchBranchOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel)
		if ok {
			popUp.AddRemoteOutputViewport, cmd = popUp.AddRemoteOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.GitStashOperationOutputPopUp:
		popUp, ok := m.PopUpModel.(*stashPopUp.GitStashOperationOutputPopUpModel)
		if ok {
			popUp.GitStashOperationOutputViewport, cmd = popUp.GitStashOperationOutputViewport.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}
