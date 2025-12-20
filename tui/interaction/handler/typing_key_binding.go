package handler

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

func handleTypingESCKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch m.PopUpType {
	case constant.CommitPopUp:
		services.GitCommitCancelService(m)
	case constant.AmendCommitPopUp:
		services.GitAmendCommitCancelService(m)
	case constant.AddRemotePromptPopUp:
		services.GitAddRemoteCancelService(m)
	case constant.CreateNewBranchPopUp:
		m.ShowPopUp.Store(false)
		m.IsTyping.Store(false)
		m.PopUpType = constant.NoPopUp
		m.PopUpModel = nil
	case constant.GitStashMessagePopUp:
		m.ShowPopUp.Store(false)
		m.IsTyping.Store(false)
		m.PopUpType = constant.NoPopUp
		m.PopUpModel = nil
	case constant.CreateBranchBasedOnRemotePopUp:
		m.ShowPopUp.Store(false)
		m.IsTyping.Store(false)
		m.PopUpType = constant.NoPopUp
		m.PopUpModel = nil
	}
	return m, nil
}

func handleTypingTabKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch m.PopUpType {
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.MessageTextInput.Focus()
				popUp.DescriptionTextAreaInput.Blur()
			case 2:
				popUp.MessageTextInput.Blur()
				popUp.DescriptionTextAreaInput.Focus()
			}
		}
	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.MessageTextInput.Focus()
				popUp.DescriptionTextAreaInput.Blur()
			case 2:
				popUp.MessageTextInput.Blur()
				popUp.DescriptionTextAreaInput.Focus()
			}
		}
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.RemoteNameTextInput.Focus()
				popUp.RemoteUrlTextInput.Blur()
			case 2:
				popUp.RemoteNameTextInput.Blur()
				popUp.RemoteUrlTextInput.Focus()
			}
		}
	}
	return m, nil
}

func handleTypingShiftTabKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch m.PopUpType {
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.MessageTextInput.Focus()
				popUp.DescriptionTextAreaInput.Blur()
			case 2:
				popUp.MessageTextInput.Blur()
				popUp.DescriptionTextAreaInput.Focus()
			}
		}
	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.MessageTextInput.Focus()
				popUp.DescriptionTextAreaInput.Blur()
			case 2:
				popUp.MessageTextInput.Blur()
				popUp.DescriptionTextAreaInput.Focus()
			}
		}
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.RemoteNameTextInput.Focus()
				popUp.RemoteUrlTextInput.Blur()
			case 2:
				popUp.RemoteNameTextInput.Blur()
				popUp.RemoteUrlTextInput.Focus()
			}
		}
	}
	return m, nil
}

func handleTypingCtrleKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch m.PopUpType {
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			// once they start for commit process, reinit the input focus
			popUp.MessageTextInput.Focus()
			popUp.DescriptionTextAreaInput.Blur()
			popUp.CurrentActiveInputIndex = 1
			// start a seperate thread commit them and set the value of msg and desc to "" if committed successfully
			// also do not start any git operation is message is no provided
			if !popUp.IsProcessing.Load() && len(popUp.MessageTextInput.Value()) > 0 {
				services.GitCommitService(m, popUp.IsAmendCommit)
				popUp.InitialCommitStarted.Store(true)
				// Start spinner ticking
				return m, popUp.Spinner.Tick
			}
		}
	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
		if ok {
			// once they start for amend commit process, reinit the input focus
			popUp.MessageTextInput.Focus()
			popUp.DescriptionTextAreaInput.Blur()
			popUp.CurrentActiveInputIndex = 1
			if !popUp.IsProcessing.Load() && len(popUp.MessageTextInput.Value()) > 0 {
				services.GitAmendCommitService(m, popUp.IsAmendCommit)
				popUp.InitialCommitStarted.Store(true)
				// Start spinner ticking
				return m, popUp.Spinner.Tick
			}
		}
	}
	return m, nil
}

func handleTypingEnterKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch m.PopUpType {
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel)
		if ok {
			// once they start for commit process, reinit the input focus
			popUp.RemoteNameTextInput.Focus()
			popUp.RemoteUrlTextInput.Blur()
			popUp.CurrentActiveInputIndex = 1
			// start a seperate thread that stage the current selected files and commit them and set the value of msg and desc to "" if committed successfully
			// also do not start any git operation is message is no provided
			if !popUp.IsProcessing.Load() && len(popUp.RemoteNameTextInput.Value()) > 0 && len(popUp.RemoteUrlTextInput.Value()) > 0 {
				services.GitAddRemoteService(m)
			}
		}

	case constant.CreateNewBranchPopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateNewBranchPopUpModel)
		if ok {
			// we direclty close the pop up and trigger the branch creation operation
			validBranchName, _ := api.IsBranchNameValid(popUp.NewBranchNameInput.Value())
			if len(validBranchName) > 0 {
				switch popUp.CreateType {
				case git.NEWBRANCH:
					services.GitCreateNewBranchService(m, validBranchName)
				case git.NEWBRANCHANDSWITCH:
					services.GitCreateNewBranchAndSwitchService(m, validBranchName)
				}
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		}

	case constant.GitStashMessagePopUp:
		popUp, ok := m.PopUpModel.(*stashPopUp.GitStashMessagePopUpModel)
		if ok {
			msg := popUp.StashMessageInput.Value()
			switch popUp.StashType {
			case git.STASHALL:
				stashPopUp.InitGitStashConfirmPromptPopUpModel(m, git.STASHALL, "", "", msg)
			case git.STASHFILE:
				stashPopUp.InitGitStashConfirmPromptPopUpModel(m, git.STASHFILE, popUp.FilePathName, "", msg)
			}
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
			m.PopUpType = constant.GitStashConfirmPromptPopUp
		}

	case constant.CreateBranchBasedOnRemotePopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemotePopUpModel)
		if ok {
			// we direclty close the pop up and trigger the branch creation operation
			validBranchName, _ := api.IsBranchNameValid(popUp.RemoteBranchNameInput.Value())
			remoteOrigin := popUp.RemoteOrigin
			if len(validBranchName) > 0 {
				branchPopUp.InitCreateBranchBasedOnRemoteOutputPopUp(m)
				popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemoteOutputPopUpModel)
				if ok {
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					m.PopUpType = constant.CreateBranchBasedOnRemoteOutputPopUp
					popUp.IsProcessing.Store(true)
					services.CreateNewBranchBasedOnRemoteService(m, remoteOrigin, validBranchName)
					return m, popUp.Spinner.Tick
				} else {
					m.ShowPopUp.Store(false)
					m.IsTyping.Store(false)
					m.PopUpType = constant.NoPopUp
				}
			}
		}

	}
	return m, nil
}
