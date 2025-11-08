package tui

import (
	"gitti/api"
	"gitti/api/git"
	"gitti/tui/constant"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/google/uuid"
)

// the function to handle bubbletea key interactions
func gittiKeyInteraction(msg tea.KeyMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	if m.IsTyping.Load() {
		return handleTypingKeyBindingInteraction(msg, m)
	} else {
		return handleNonTypingGlobalKeyBindingInteraction(msg, m)
	}
}

// typing is currently only on pop up model, so we can safely process it without checking if they were on pop up or not
func handleTypingKeyBindingInteraction(msg tea.KeyMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		api.GITDAEMON.Stop()
		return m, tea.Quit
	case "esc":
		switch m.PopUpType {
		case constant.CommitPopUp:
			gitCommitCancelService(m)
		case constant.AddRemotePromptPopUp:
			gitAddRemoteCancelService(m)
		case constant.CreateNewBranchPopUp:
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpModel = nil
		}
		return m, nil
	// in typing mode, tab is move to next input
	case "tab":
		switch m.PopUpType {
		case constant.CommitPopUp:
			popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
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
			popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
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

	// in typing mode, shift+tab is move to previous input
	case "shift+tab":
		switch m.PopUpType {
		case constant.CommitPopUp:
			popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
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
			popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
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

	case "ctrl+enter":
		switch m.PopUpType {
		case constant.CommitPopUp:
			popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
			if ok {
				// once they start for commit process, reinit the input focus
				popUp.MessageTextInput.Focus()
				popUp.DescriptionTextAreaInput.Blur()
				popUp.CurrentActiveInputIndex = 1
				// start a seperate thread that stage the current selected files and commit them and set the value of msg and desc to "" if committed successfully
				// also do not start any git operation is message is no provided
				if !popUp.IsProcessing.Load() {
					gitCommitService(m)
					// Start spinner ticking
					return m, popUp.Spinner.Tick
				}
			}
		case constant.AddRemotePromptPopUp:
			popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
			if ok {
				// once they start for commit process, reinit the input focus
				popUp.RemoteNameTextInput.Focus()
				popUp.RemoteUrlTextInput.Blur()
				popUp.CurrentActiveInputIndex = 1
				// start a seperate thread that stage the current selected files and commit them and set the value of msg and desc to "" if committed successfully
				// also do not start any git operation is message is no provided
				if !popUp.IsProcessing.Load() {
					gitAddRemoteService(m)
				}
			}
		case constant.CreateNewBranchPopUp:
			popUp, ok := m.PopUpModel.(*CreateNewBranchPopUpModel)
			if ok {
				// we direclty trigger the branch creation operation and close the pop up, we will assume this always result in success
				validBranchName, _ := api.IsBranchNameValid(popUp.NewBranchNameInput.Value())
				switch popUp.CreateType {
				case git.NEWBRANCH:
					m.GitState.GitBranch.GitCreateNewBranch(validBranchName)
				case git.NEWBRANCHANDSWITCH:
					m.GitState.GitBranch.GitCreateNewBranchAndSwitch(validBranchName)
				}
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpModel = nil
			}
		}
		return m, nil
	}
	switch m.PopUpType {
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
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
		popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
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
		popUp, ok := m.PopUpModel.(*CreateNewBranchPopUpModel)
		if ok {
			var cmd tea.Cmd
			popUp.NewBranchNameInput, cmd = popUp.NewBranchNameInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func handleNonTypingGlobalKeyBindingInteraction(msg tea.KeyMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+c":
		api.GITDAEMON.Stop()
		return m, tea.Quit
	case "n", "N":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedContainer == constant.LocalBranchComponent {
				m.PopUpType = constant.ChooseNewBranchTypePopUp
				m.IsTyping.Store(false)
				m.ShowPopUp.Store(true)
				if _, ok := m.PopUpModel.(*ChooseNewBranchTypeOptionPopUpModel); !ok {
					initChooseNewBranchTypePopUpModel(m)
				}
			}
		}
		return m, nil
	case "b", "B":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedContainer != constant.LocalBranchComponent {
				m.CurrentSelectedContainer = constant.LocalBranchComponent
			} else {
				m.CurrentSelectedContainer = constant.NoneSelected
			}
		}
		return m, nil
	case "f", "F":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedContainer != constant.ModifiedFilesComponent {
				m.CurrentSelectedContainer = constant.ModifiedFilesComponent
			} else {
				m.CurrentSelectedContainer = constant.NoneSelected
			}
		}
		return m, nil
	case "c", "C":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedContainer == constant.ModifiedFilesComponent {
				m.ShowPopUp.Store(true)
				m.PopUpType = constant.CommitPopUp
				m.GitState.GitCommit.ClearGitCommitOutput()

				// if the current pop up model is not commit pop up model, then init it
				if popUp, ok := m.PopUpModel.(*GitCommitPopUpModel); !ok {
					initGitCommitPopUpModel(m)
				} else {
					popUp.GitCommitOutputViewport.SetContent("")
					popUp.SessionID = uuid.New()
				}
				m.IsTyping.Store(true)
			}
		}
		return m, nil
	case "p", "P":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedContainer == constant.ModifiedFilesComponent || m.CurrentSelectedContainer == constant.NoneSelected || m.CurrentSelectedContainer == constant.LocalBranchComponent {
				// first we need to check if there are any push origin for this repo
				// if not we prompt the user to add a new remote origin
				if !m.GitState.GitCommit.CheckRemoteExist() {
					m.ShowPopUp.Store(true)
					m.PopUpType = constant.AddRemotePromptPopUp
					// if the current pop up model is not commit pop up model, then init it
					if popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel); !ok {
						initAddRemotePromptPopUpModel(m, true)
					} else {
						popUp.AddRemoteOutputViewport.SetContent("")
					}
					m.IsTyping.Store(true)
				} else {
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					remotes := m.GitState.GitCommit.Remote()
					if len(remotes) == 1 {
						m.PopUpType = constant.ChoosePushTypePopUp
						// if the current pop up model is not commit pop up model, then init it and start git push service
						initChoosePushTypePopUpModel(m, remotes[0].Name)
					} else if len(remotes) > 1 {
						// if remote is more than 1 let user choose which remote to push to first before pushing
						m.PopUpType = constant.ChooseRemotePopUp
						if _, ok := m.PopUpModel.(*ChooseRemotePopUpModel); !ok {
							initGitRemotePushChooseRemotePopUpModel(m, remotes)
						}
					}
				}
			}
		}
		return m, nil

	case "enter":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedContainer == constant.ModifiedFilesComponent && len(m.CurrentRepoModifiedFilesInfoList.Items()) > 0 {
				m.CurrentSelectedContainer = constant.FileDiffComponent
			}
		} else {
			switch m.PopUpType {
			case constant.ChooseRemotePopUp:
				popUp, ok := m.PopUpModel.(*ChooseRemotePopUpModel)
				if ok {
					remote := popUp.RemoteList.SelectedItem()
					m.PopUpType = constant.ChoosePushTypePopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					initChoosePushTypePopUpModel(m, remote.(gitRemoteItem).Name)
				}
			case constant.ChoosePushTypePopUp:
				popUp, ok := m.PopUpModel.(*ChoosePushTypePopUpModel)
				if ok {
					m.PopUpType = constant.GitRemotePushPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					selectedOption := popUp.PushOptionList.SelectedItem()
					return initGitRemotePushPopUpModelAndStartGitRemotePushService(m, popUp.RemoteName, selectedOption.(gitPushOptionItem).pushType)
				}
			case constant.ChooseNewBranchTypePopUp:
				popUp, ok := m.PopUpModel.(*ChooseNewBranchTypeOptionPopUpModel)
				if ok {
					m.PopUpType = constant.CreateNewBranchPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(true)
					selectedOption := popUp.NewBranchTypeOptionList.SelectedItem()
					initCreateNewBranchPopUpModel(m, selectedOption.(gitNewBranchTypeOptionItem).newBranchType)
				}
			}
		}
		return m, nil
	case "esc":
		if m.ShowPopUp.Load() {
			switch m.PopUpType {
			case constant.GitRemotePushPopUp:
				gitRemotePushCancelService(m)
			case constant.ChooseRemotePopUp:
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
			case constant.ChoosePushTypePopUp:
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpModel = nil
			case constant.ChooseNewBranchTypePopUp:
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpModel = nil
			}
			return m, nil
		} else {
			switch m.CurrentSelectedContainer {
			case constant.NoneSelected:
				api.GITDAEMON.Stop()
				return m, tea.Quit
			case constant.FileDiffComponent:
				m.CurrentSelectedContainer = constant.ModifiedFilesComponent
			case constant.LocalBranchComponent:
				m.CurrentSelectedContainer = constant.NoneSelected
			case constant.ModifiedFilesComponent:
				m.CurrentSelectedContainer = constant.NoneSelected
			}
		}
		return m, nil
	case "s", "S":
		if m.CurrentSelectedContainer == constant.ModifiedFilesComponent {
			currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
			var fileName string
			if currentSelectedModifiedFile != nil {
				fileName = currentSelectedModifiedFile.(gitModifiedFilesItem).FileName
				m.GitState.GitFiles.ToggleFilesStageStatus(fileName)
			}
		}
		return m, nil
	case "up", "k":
		if !m.ShowPopUp.Load() {
			switch m.CurrentSelectedContainer {
			case constant.LocalBranchComponent:
				// we don't use the list native Update() because we track the current selected index
				if m.CurrentRepoBranchesInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoBranchesInfoList.Index() - 1
					m.CurrentRepoBranchesInfoList.Select(latestIndex)
					m.NavigationIndexPosition.LocalBranchComponent = latestIndex
				}
			case constant.ModifiedFilesComponent:
				// we don't use the list native Update() because we need to also render the diff view as well as track the current selected index
				if m.CurrentRepoModifiedFilesInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoModifiedFilesInfoList.Index() - 1
					m.CurrentRepoModifiedFilesInfoList.Select(latestIndex)
					m.NavigationIndexPosition.ModifiedFilesComponent = latestIndex
					reinitAndRenderModifiedFileDiffViewPort(m)
				}
			case constant.FileDiffComponent:
				m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
				return m, cmd
			}
		} else {
			// for within pop up component
			switch m.PopUpType {
			case constant.ChooseRemotePopUp:
				popUp, ok := m.PopUpModel.(*ChooseRemotePopUpModel)
				if ok {
					popUp.RemoteList, cmd = popUp.RemoteList.Update(msg)
					return m, cmd
				}
			case constant.ChoosePushTypePopUp:
				popUp, ok := m.PopUpModel.(*ChoosePushTypePopUpModel)
				if ok {
					popUp.PushOptionList, cmd = popUp.PushOptionList.Update(msg)
					return m, cmd
				}
			case constant.ChooseNewBranchTypePopUp:
				popUp, ok := m.PopUpModel.(*ChooseNewBranchTypeOptionPopUpModel)
				if ok {
					popUp.NewBranchTypeOptionList, cmd = popUp.NewBranchTypeOptionList.Update(msg)
					return m, cmd
				}
			}
		}
		return m, nil
	case "down", "j":
		if !m.ShowPopUp.Load() {
			switch m.CurrentSelectedContainer {
			case constant.LocalBranchComponent:
				// we don't use the list native Update() because we track the current selected index
				if m.CurrentRepoBranchesInfoList.Index() < len(m.CurrentRepoBranchesInfoList.Items())-1 {
					latestIndex := m.CurrentRepoBranchesInfoList.Index() + 1
					m.CurrentRepoBranchesInfoList.Select(latestIndex)
					m.NavigationIndexPosition.LocalBranchComponent = latestIndex
				}
			case constant.ModifiedFilesComponent:
				// we don't use the list native Update() because we need to also render the diff view as well as track the current selected index
				if m.CurrentRepoModifiedFilesInfoList.Index() < len(m.CurrentRepoModifiedFilesInfoList.Items())-1 {
					latestIndex := m.CurrentRepoModifiedFilesInfoList.Index() + 1
					m.CurrentRepoModifiedFilesInfoList.Select(latestIndex)
					m.NavigationIndexPosition.ModifiedFilesComponent = latestIndex
					reinitAndRenderModifiedFileDiffViewPort(m)
				}
			case constant.FileDiffComponent:
				m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
				return m, cmd
			}
		} else {
			// for within pop up component
			switch m.PopUpType {
			case constant.ChooseRemotePopUp:
				popUp, ok := m.PopUpModel.(*ChooseRemotePopUpModel)
				if ok {
					popUp.RemoteList, cmd = popUp.RemoteList.Update(msg)
					return m, cmd
				}
			case constant.ChoosePushTypePopUp:
				popUp, ok := m.PopUpModel.(*ChoosePushTypePopUpModel)
				if ok {
					popUp.PushOptionList, cmd = popUp.PushOptionList.Update(msg)
					return m, cmd
				}
			case constant.ChooseNewBranchTypePopUp:
				popUp, ok := m.PopUpModel.(*ChooseNewBranchTypeOptionPopUpModel)
				if ok {
					popUp.NewBranchTypeOptionList, cmd = popUp.NewBranchTypeOptionList.Update(msg)
					return m, cmd
				}
			}
		}
		return m, nil
	case "left", "h":
		if !m.ShowPopUp.Load() {
			m.CurrentSelectedFileDiffViewport.MoveLeft(1)
		} else {
			switch m.PopUpType {
			case constant.CommitPopUp:
				popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
				if ok {
					popUp.GitCommitOutputViewport, cmd = popUp.GitCommitOutputViewport.Update(msg)
					return m, cmd
				}
			}
		}
	case "right", "l":
		if !m.ShowPopUp.Load() {
			m.CurrentSelectedFileDiffViewport.MoveRight(1)
		} else {
			switch m.PopUpType {
			case constant.CommitPopUp:
				popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
				if ok {
					popUp.GitCommitOutputViewport, cmd = popUp.GitCommitOutputViewport.Update(msg)
					return m, cmd
				}
			}
		}
	}
	return m, nil
}

func GittiMouseInteraction(msg tea.MouseMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "wheelleft":
		if !m.ShowPopUp.Load() {
			m.CurrentSelectedFileDiffViewport.MoveLeft(1)
		}
	case "wheelright":
		if !m.ShowPopUp.Load() {
			m.CurrentSelectedFileDiffViewport.MoveRight(1)
		}
	case "wheelup":
		if !m.ShowPopUp.Load() {
			m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
			return m, cmd
		} else {
			switch m.PopUpType {
			case constant.CommitPopUp:
				popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
				if ok {
					popUp.GitCommitOutputViewport, cmd = popUp.GitCommitOutputViewport.Update(msg)
					return m, cmd
				}
			}
		}
	case "wheeldown":
		if !m.ShowPopUp.Load() {
			m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
			return m, cmd
		} else {
			switch m.PopUpType {
			case constant.CommitPopUp:
				popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
				if ok {
					popUp.GitCommitOutputViewport, cmd = popUp.GitCommitOutputViewport.Update(msg)
					return m, cmd
				}
			}
		}

	}
	return m, nil
}
