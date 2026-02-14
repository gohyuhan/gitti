package handler

import (
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
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
	case constant.CreateNewBranchPopUp,
		constant.GitStashMessagePopUp,
		constant.CreateBranchBasedOnRemotePopUp,
		constant.CreateTagPopUp,
		constant.EditRemotePromptPopUp:
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
	case constant.EditRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.EditRemotePromptPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.NewRemoteNameTextInput.Focus()
				popUp.NewRemoteUrlTextInput.Blur()
			case 2:
				popUp.NewRemoteNameTextInput.Blur()
				popUp.NewRemoteUrlTextInput.Focus()
			}
		}
	case constant.CreateTagPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.CreateTagPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.TagNameInput.Focus()
				popUp.TagMessageTextAreaInput.Blur()
			case 2:
				popUp.TagNameInput.Blur()
				popUp.TagMessageTextAreaInput.Focus()
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
	case constant.EditRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.EditRemotePromptPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.NewRemoteNameTextInput.Focus()
				popUp.NewRemoteUrlTextInput.Blur()
			case 2:
				popUp.NewRemoteNameTextInput.Blur()
				popUp.NewRemoteUrlTextInput.Focus()
			}
		}
	case constant.CreateTagPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.CreateTagPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.TagNameInput.Focus()
				popUp.TagMessageTextAreaInput.Blur()
			case 2:
				popUp.TagNameInput.Blur()
				popUp.TagMessageTextAreaInput.Focus()
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
	case constant.CreateTagPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.CreateTagPopUpModel)
		if ok {
			if utf8.RuneCountInString(popUp.TagNameInput.Value()) < 1 {
				return m, nil
			}
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
			tagPopUp.InitCreateTagConfirmationPopUpModel(m, popUp.TagNameInput.Value(), popUp.TagMessageTextAreaInput.Value(), popUp.CommitHash, popUp.CommitMessage)
			m.PopUpType = constant.CreateTagConfirmationPopUp
		}
	}
	return m, nil
}

func handleTypingEnterKeyBindingInteraction(m *types.GittiModel, msg tea.KeyMsg) (*types.GittiModel, tea.Cmd) {
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
			if !popUp.IsProcessing.Load() && utf8.RuneCountInString(popUp.RemoteNameTextInput.Value()) > 0 && utf8.RuneCountInString(popUp.RemoteUrlTextInput.Value()) > 0 {
				services.GitAddRemoteService(m)
			}
		}
	case constant.EditRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.EditRemotePromptPopUpModel)
		if ok {
			newRemoteName := popUp.NewRemoteNameTextInput.Value()
			newRemoteUrl := popUp.NewRemoteUrlTextInput.Value()

			services.GitEditRemoteNameAndUrlService(m, popUp.OldRemoteName, newRemoteName, popUp.OldRemoteUrl, newRemoteUrl)
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpType = constant.NoPopUp
			m.PopUpModel = nil
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

	// the following is to handle the change line for textarea input
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			if popUp.CurrentActiveInputIndex == 2 {
				var cmd tea.Cmd
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}

	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
		if ok {
			if popUp.CurrentActiveInputIndex == 2 {
				var cmd tea.Cmd
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}

	case constant.CreateTagPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.CreateTagPopUpModel)
		if ok {
			if popUp.CurrentActiveInputIndex == 2 {
				var cmd tea.Cmd
				popUp.TagMessageTextAreaInput, cmd = popUp.TagMessageTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	}
	return m, nil
}

func handleTypingCtrlpKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	content, err := clipboard.ReadAll()
	if err != nil {
		return m, nil
	}
	msg := tea.PasteMsg{
		Content: content,
	}
	switch m.PopUpType {
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				var cmd tea.Cmd
				popUp.MessageTextInput, cmd = popUp.MessageTextInput.Update(msg)
				return m, cmd

			case 2:
				var cmd tea.Cmd
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				var cmd tea.Cmd
				popUp.MessageTextInput, cmd = popUp.MessageTextInput.Update(msg)
				return m, cmd

			case 2:
				var cmd tea.Cmd
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				var cmd tea.Cmd
				popUp.RemoteNameTextInput, cmd = popUp.RemoteNameTextInput.Update(msg)
				return m, cmd

			case 2:
				var cmd tea.Cmd
				popUp.RemoteUrlTextInput, cmd = popUp.RemoteUrlTextInput.Update(msg)
				return m, cmd
			}
		}
	case constant.CreateNewBranchPopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateNewBranchPopUpModel)
		if ok {
			var cmd tea.Cmd
			popUp.NewBranchNameInput, cmd = popUp.NewBranchNameInput.Update(msg)
			return m, cmd
		}
	case constant.GitStashMessagePopUp:
		popUp, ok := m.PopUpModel.(*stashPopUp.GitStashMessagePopUpModel)
		if ok {
			var cmd tea.Cmd
			popUp.StashMessageInput, cmd = popUp.StashMessageInput.Update(msg)
			return m, cmd
		}
	case constant.CreateBranchBasedOnRemotePopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemotePopUpModel)
		if ok {
			var cmd tea.Cmd
			popUp.RemoteBranchNameInput, cmd = popUp.RemoteBranchNameInput.Update(msg)
			return m, cmd
		}
	case constant.CreateTagPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.CreateTagPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				var cmd tea.Cmd
				popUp.TagNameInput, cmd = popUp.TagNameInput.Update(msg)
				return m, cmd

			case 2:
				var cmd tea.Cmd
				popUp.TagMessageTextAreaInput, cmd = popUp.TagMessageTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	}
	return m, nil
}

func handleTypingCtrlyKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var content string
	switch m.PopUpType {
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				content = popUp.MessageTextInput.Value()
			case 2:
				content = popUp.DescriptionTextAreaInput.Value()
			}
		}
	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				content = popUp.MessageTextInput.Value()
			case 2:
				content = popUp.DescriptionTextAreaInput.Value()
			}
		}
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				content = popUp.RemoteNameTextInput.Value()
			case 2:
				content = popUp.RemoteUrlTextInput.Value()
			}
		}
	case constant.CreateNewBranchPopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateNewBranchPopUpModel)
		if ok {
			content = popUp.NewBranchNameInput.Value()
		}
	case constant.GitStashMessagePopUp:
		popUp, ok := m.PopUpModel.(*stashPopUp.GitStashMessagePopUpModel)
		if ok {
			content = popUp.StashMessageInput.Value()
		}
	case constant.CreateBranchBasedOnRemotePopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemotePopUpModel)
		if ok {
			content = popUp.RemoteBranchNameInput.Value()
		}
	case constant.CreateTagPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.CreateTagPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				content = popUp.TagNameInput.Value()
			case 2:
				content = popUp.TagMessageTextAreaInput.Value()
			}
		}
	}

	err := clipboard.WriteAll(content)
	if err != nil {
		// TODO: log the error
	}

	return m, nil
}
