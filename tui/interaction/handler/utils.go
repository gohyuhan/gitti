package handler

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	"github.com/gohyuhan/gitti/tui/popup/commit"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	commitLogPopUp "github.com/gohyuhan/gitti/tui/popup/commitlog"
	discardPopUp "github.com/gohyuhan/gitti/tui/popup/discard"
	keybindingPopUp "github.com/gohyuhan/gitti/tui/popup/keybinding"
	pullPopUp "github.com/gohyuhan/gitti/tui/popup/pull"
	pushPopUp "github.com/gohyuhan/gitti/tui/popup/push"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	resolvePopUp "github.com/gohyuhan/gitti/tui/popup/resolve"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ----------------------------------
//
//	Handle up/down key messages for pop-up navigation and scrolling.
//	Responsibility: Routes vertical navigation events (Up/k, Down/j) to the currently
//	active popup's internal lists or text viewports, ensuring the correct popup state is updated.
//
// ----------------------------------
func UpDownKeyPressMsgUpdateForPopUp(msg tea.KeyPressMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	// for within pop up component
	switch m.PopUpType {
	// following is for list component
	case constant.ChooseRemotePopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.ChooseRemotePopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.RemoteList.Index() > 0 {
					latestIndex := popUp.RemoteList.Index() - 1
					popUp.RemoteList.Select(latestIndex)
				}
			case "down", "j":
				if popUp.RemoteList.Index() < len(popUp.RemoteList.Items())-1 {
					latestIndex := popUp.RemoteList.Index() + 1
					popUp.RemoteList.Select(latestIndex)
				}
			}
			popUp.RemoteList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.RemoteList, constant.MaxChooseRemotePopUpWidth)
			return m, nil
		}
	case constant.ChoosePushTypePopUp:
		popUp, ok := m.PopUpModel.(*pushPopUp.ChoosePushTypePopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.PushOptionList.Index() > 0 {
					latestIndex := popUp.PushOptionList.Index() - 1
					popUp.PushOptionList.Select(latestIndex)
				}
			case "down", "j":
				if popUp.PushOptionList.Index() < len(popUp.PushOptionList.Items())-1 {
					latestIndex := popUp.PushOptionList.Index() + 1
					popUp.PushOptionList.Select(latestIndex)
				}
			}
			popUp.PushOptionList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.PushOptionList, constant.MaxChoosePushTypePopUpWidth)
			return m, nil
		}
	case constant.ChooseNewBranchTypePopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.ChooseNewBranchTypeOptionPopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.NewBranchTypeOptionList.Index() > 0 {
					latestIndex := popUp.NewBranchTypeOptionList.Index() - 1
					popUp.NewBranchTypeOptionList.Select(latestIndex)
				}
			case "down", "j":
				if popUp.NewBranchTypeOptionList.Index() < len(popUp.NewBranchTypeOptionList.Items())-1 {
					latestIndex := popUp.NewBranchTypeOptionList.Index() + 1
					popUp.NewBranchTypeOptionList.Select(latestIndex)
				}
			}
			popUp.NewBranchTypeOptionList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.NewBranchTypeOptionList, constant.MaxChooseNewBranchTypePopUpWidth)
			return m, nil
		}
	case constant.ChooseSwitchBranchTypePopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.ChooseSwitchBranchTypePopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.SwitchTypeOptionList.Index() > 0 {
					latestIndex := popUp.SwitchTypeOptionList.Index() - 1
					popUp.SwitchTypeOptionList.Select(latestIndex)
				}
			case "down", "j":
				if popUp.SwitchTypeOptionList.Index() < len(popUp.SwitchTypeOptionList.Items())-1 {
					latestIndex := popUp.SwitchTypeOptionList.Index() + 1
					popUp.SwitchTypeOptionList.Select(latestIndex)
				}
			}
			popUp.SwitchTypeOptionList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.SwitchTypeOptionList, constant.MaxChooseSwitchBranchTypePopUpWidth)
			return m, nil
		}
	case constant.ChooseGitPullTypePopUp:
		popUp, ok := m.PopUpModel.(*pullPopUp.ChooseGitPullTypePopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.PullTypeOptionList.Index() > 0 {
					latestIndex := popUp.PullTypeOptionList.Index() - 1
					popUp.PullTypeOptionList.Select(latestIndex)
				}
			case "down", "j":
				if popUp.PullTypeOptionList.Index() < len(popUp.PullTypeOptionList.Items())-1 {
					latestIndex := popUp.PullTypeOptionList.Index() + 1
					popUp.PullTypeOptionList.Select(latestIndex)
				}
			}
			popUp.PullTypeOptionList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.PullTypeOptionList, constant.MaxChooseGitPullTypePopUpWidth)
			return m, nil
		}
	case constant.GitDiscardTypeOptionPopUp:
		popUp, ok := m.PopUpModel.(*discardPopUp.GitDiscardTypeOptionPopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.DiscardTypeOptionList.Index() > 0 {
					latestIndex := popUp.DiscardTypeOptionList.Index() - 1
					popUp.DiscardTypeOptionList.Select(latestIndex)
				}
			case "down", "j":
				if popUp.DiscardTypeOptionList.Index() < len(popUp.DiscardTypeOptionList.Items())-1 {
					latestIndex := popUp.DiscardTypeOptionList.Index() + 1
					popUp.DiscardTypeOptionList.Select(latestIndex)
				}
			}
			popUp.DiscardTypeOptionList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.DiscardTypeOptionList, constant.MaxGitDiscardTypeOptionPopUpWidth)
			return m, nil
		}
	case constant.GitResolveConflictOptionPopUp:
		popUp, ok := m.PopUpModel.(*resolvePopUp.GitResolveConflictOptionPopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.ResolveConflictOptionList.Index() > 0 {
					latestIndex := popUp.ResolveConflictOptionList.Index() - 1
					popUp.ResolveConflictOptionList.Select(latestIndex)
				}
			case "down", "j":
				if popUp.ResolveConflictOptionList.Index() < len(popUp.ResolveConflictOptionList.Items())-1 {
					latestIndex := popUp.ResolveConflictOptionList.Index() + 1
					popUp.ResolveConflictOptionList.Select(latestIndex)
				}
			}
			popUp.ResolveConflictOptionList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.ResolveConflictOptionList, constant.MaxGitResolveConflictOptionPopUpWidth)
			return m, nil
		}
	case constant.GitResetLatestCommitTypeOptionPopUp:
		popUp, ok := m.PopUpModel.(*commit.GitResetLatestCommitTypeOptionPopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.ResetLatestCommitTypeOptionList.Index() > 0 {
					latestIndex := popUp.ResetLatestCommitTypeOptionList.Index() - 1
					popUp.ResetLatestCommitTypeOptionList.Select(latestIndex)
				}
			case "down", "j":
				if popUp.ResetLatestCommitTypeOptionList.Index() < len(popUp.ResetLatestCommitTypeOptionList.Items())-1 {
					latestIndex := popUp.ResetLatestCommitTypeOptionList.Index() + 1
					popUp.ResetLatestCommitTypeOptionList.Select(latestIndex)
				}
			}
			popUp.ResetLatestCommitTypeOptionList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.ResetLatestCommitTypeOptionList, constant.MaxGitResetLatestCommitTypeOptionPopUpWidth)
			return m, nil
		}
	case constant.GitResetToSelectedCommitTypeOptionPopUp:
		popUp, ok := m.PopUpModel.(*commit.GitResetToSelectedCommitTypeOptionPopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.ResetToSelectedCommitTypeOptionList.Index() > 0 {
					latestIndex := popUp.ResetToSelectedCommitTypeOptionList.Index() - 1
					popUp.ResetToSelectedCommitTypeOptionList.Select(latestIndex)
				}
			case "down", "j":
				if popUp.ResetToSelectedCommitTypeOptionList.Index() < len(popUp.ResetToSelectedCommitTypeOptionList.Items())-1 {
					latestIndex := popUp.ResetToSelectedCommitTypeOptionList.Index() + 1
					popUp.ResetToSelectedCommitTypeOptionList.Select(latestIndex)
				}
			}
			popUp.ResetToSelectedCommitTypeOptionList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.ResetToSelectedCommitTypeOptionList, constant.MaxGitResetToSelectedCommitTypeOptionPopUpWidth)
			return m, nil
		}

	case constant.GitCherryPickOptionSelectionPopUp:
		popUp, ok := m.PopUpModel.(*commitLogPopUp.GitCherryPickOptionSelectionPopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.CherryPickedOpsOption.Index() > 0 {
					latestIndex := popUp.CherryPickedOpsOption.Index() - 1
					popUp.CherryPickedOpsOption.Select(latestIndex)
				}
			case "down", "j":
				if popUp.CherryPickedOpsOption.Index() < len(popUp.CherryPickedOpsOption.Items())-1 {
					latestIndex := popUp.CherryPickedOpsOption.Index() + 1
					popUp.CherryPickedOpsOption.Select(latestIndex)
				}
			}
			popUp.CherryPickedOpsOption.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.CherryPickedOpsOption, constant.MaxGitCherryPickOptionSelectionPopUpWidth)
			return m, nil
		}
	case constant.GitCherryPickPopUp:
		popUp, ok := m.PopUpModel.(*commitLogPopUp.GitCherryPickPopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.CurrentBranchCherryPickCommitLog.Index() > 0 {
					latestIndex := popUp.CurrentBranchCherryPickCommitLog.Index() - 1
					popUp.CurrentBranchCherryPickCommitLog.Select(latestIndex)
				}
			case "down", "j":
				if popUp.CurrentBranchCherryPickCommitLog.Index() < len(popUp.CurrentBranchCherryPickCommitLog.Items())-1 {
					latestIndex := popUp.CurrentBranchCherryPickCommitLog.Index() + 1
					popUp.CurrentBranchCherryPickCommitLog.Select(latestIndex)
				}
			}
			popUp.CurrentBranchCherryPickCommitLog.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.CurrentBranchCherryPickCommitLog, constant.MaxGitCherryPickPopUpWidth)
			return m, nil
		}
	case constant.GitEditCherryPickPopUp:
		popUp, ok := m.PopUpModel.(*commitLogPopUp.GitEditCherryPickPopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.CherryPickedCommitLog.Index() > 0 {
					latestIndex := popUp.CherryPickedCommitLog.Index() - 1
					popUp.CherryPickedCommitLog.Select(latestIndex)
				}
			case "down", "j":
				if popUp.CherryPickedCommitLog.Index() < len(popUp.CherryPickedCommitLog.Items())-1 {
					latestIndex := popUp.CherryPickedCommitLog.Index() + 1
					popUp.CherryPickedCommitLog.Select(latestIndex)
				}
			}
			popUp.CherryPickedCommitLog.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.CherryPickedCommitLog, constant.MaxGitEditCherryPickPopUpWidth)
			return m, nil
		}
	case constant.ChooseDeleteTagOptionPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.ChooseDeleteTagOptionPopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.DeleteOptionList.Index() > 0 {
					latestIndex := popUp.DeleteOptionList.Index() - 1
					popUp.DeleteOptionList.Select(latestIndex)
				}
			case "down", "j":
				if popUp.DeleteOptionList.Index() < len(popUp.DeleteOptionList.Items())-1 {
					latestIndex := popUp.DeleteOptionList.Index() + 1
					popUp.DeleteOptionList.Select(latestIndex)
				}
			}
			popUp.DeleteOptionList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.DeleteOptionList, constant.MaxChooseDeleteTagOptionPopUpWidth)
			return m, nil
		}
	case constant.ChooseRemoteForDeleteRemoteTagPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.ChooseRemoteForDeleteRemoteTagPopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.RemoteList.Index() > 0 {
					latestIndex := popUp.RemoteList.Index() - 1
					popUp.RemoteList.Select(latestIndex)
				}
			case "down", "j":
				if popUp.RemoteList.Index() < len(popUp.RemoteList.Items())-1 {
					latestIndex := popUp.RemoteList.Index() + 1
					popUp.RemoteList.Select(latestIndex)
				}
			}
			popUp.RemoteList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.RemoteList, constant.MaxChooseRemoteForDeleteRemoteTagPopUpWidth)
			return m, nil
		}

	case constant.ChoosePushTagOptionPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.ChoosePushTagOptionPopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.PushOptionList.Index() > 0 {
					latestIndex := popUp.PushOptionList.Index() - 1
					popUp.PushOptionList.Select(latestIndex)
				}
			case "down", "j":
				if popUp.PushOptionList.Index() < len(popUp.PushOptionList.Items())-1 {
					latestIndex := popUp.PushOptionList.Index() + 1
					popUp.PushOptionList.Select(latestIndex)
				}
			}
			popUp.PushOptionList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.PushOptionList, constant.MaxChoosePushTagOptionPopUpWidth)
			return m, nil
		}
	case constant.ChooseFetchTagOptionPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.ChooseFetchTagOptionPopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.FetchOptionList.Index() > 0 {
					latestIndex := popUp.FetchOptionList.Index() - 1
					popUp.FetchOptionList.Select(latestIndex)
				}
			case "down", "j":
				if popUp.FetchOptionList.Index() < len(popUp.FetchOptionList.Items())-1 {
					latestIndex := popUp.FetchOptionList.Index() + 1
					popUp.FetchOptionList.Select(latestIndex)
				}
			}
			popUp.FetchOptionList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.FetchOptionList, constant.MaxChooseFetchTagOptionPopUpWidth)
			return m, nil
		}

	case constant.GitRevertParentOptionSelectionPopUp:
		popUp, ok := m.PopUpModel.(*commitLogPopUp.GitRevertParentOptionSelectionPopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.GitRevertParentOption.Index() > 0 {
					latestIndex := popUp.GitRevertParentOption.Index() - 1
					popUp.GitRevertParentOption.Select(latestIndex)
				}
			case "down", "j":
				if popUp.GitRevertParentOption.Index() < len(popUp.GitRevertParentOption.Items())-1 {
					latestIndex := popUp.GitRevertParentOption.Index() + 1
					popUp.GitRevertParentOption.Select(latestIndex)
				}
			}
			popUp.GitRevertParentOption.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.GitRevertParentOption, constant.MaxGitRevertParentOptionSelectionPopUpWidth)
			return m, nil
		}
	case constant.ChooseRemoteBranchOptionPopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.ChooseRemoteBranchOptionPopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.RemoteBranchOptionList.Index() > 0 {
					latestIndex := popUp.RemoteBranchOptionList.Index() - 1
					popUp.RemoteBranchOptionList.Select(latestIndex)
				}
			case "down", "j":
				if popUp.RemoteBranchOptionList.Index() < len(popUp.RemoteBranchOptionList.Items())-1 {
					latestIndex := popUp.RemoteBranchOptionList.Index() + 1
					popUp.RemoteBranchOptionList.Select(latestIndex)
				}
			}
			popUp.RemoteBranchOptionList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.RemoteBranchOptionList, constant.MaxChooseRemoteBranchOptionPopUpWidth)
			return m, nil
		}

	case constant.ChooseBranchOptionForMergePopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.ChooseBranchOptionForMergePopUpModel)
		if ok {
			switch msg.String() {
			case "up", "k":
				if popUp.BranchOptionSectionSelected.Load() {
					if popUp.BranchOptionList.Index() > 0 {
						latestIndex := popUp.BranchOptionList.Index() - 1
						popUp.BranchOptionList.Select(latestIndex)
					}
				} else if popUp.SelectedBranchSectionSelected.Load() {
					if popUp.SelectedBranchList.Index() > 0 {
						latestIndex := popUp.SelectedBranchList.Index() - 1
						popUp.SelectedBranchList.Select(latestIndex)
					}
				}
			case "down", "j":
				if popUp.BranchOptionSectionSelected.Load() {
					if popUp.BranchOptionList.Index() < len(popUp.BranchOptionList.Items())-1 {
						latestIndex := popUp.BranchOptionList.Index() + 1
						popUp.BranchOptionList.Select(latestIndex)
					}
				} else if popUp.SelectedBranchSectionSelected.Load() {
					if popUp.SelectedBranchList.Index() < len(popUp.SelectedBranchList.Items())-1 {
						latestIndex := popUp.SelectedBranchList.Index() + 1
						popUp.SelectedBranchList.Select(latestIndex)
					}
				}
			}
			popUp.BranchOptionList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.BranchOptionList, constant.MaxChooseBranchOptionForMergePopUpWidth)
			popUp.SelectedBranchList.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.SelectedBranchList, constant.MaxChooseBranchOptionForMergePopUpWidth)
			return m, nil
		}

	// following is for viewport
	case constant.KeybindingAndFeatureInstructionsPopUp:
		popUp, ok := m.PopUpModel.(*keybindingPopUp.KeybindingAndFeatureInstructionsPopUpModel)
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
	case constant.BranchMergeOutputPopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.BranchMergeOutputPopUpModel)
		if ok {
			popUp.BranchMergeOutputViewport, cmd = popUp.BranchMergeOutputViewport.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle mouse wheel up/down messages for pop-up scrolling.
//	Responsibility: Routes vertical mouse scroll events to the appropriate internal list
//	or viewport of the currently active popup, matching keyboard navigation behavior.
//
// ----------------------------------
func UpDownMouseMsgUpdateForPopUp(msg tea.MouseMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	// for pop up that have viewport
	switch m.PopUpType {
	case constant.KeybindingAndFeatureInstructionsPopUp:
		popUp, ok := m.PopUpModel.(*keybindingPopUp.KeybindingAndFeatureInstructionsPopUpModel)
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
	case constant.BranchMergeOutputPopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.BranchMergeOutputPopUpModel)
		if ok {
			popUp.BranchMergeOutputViewport, cmd = popUp.BranchMergeOutputViewport.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}
