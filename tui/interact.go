package tui

import (
	"gitti/api"

	tea "github.com/charmbracelet/bubbletea/v2"
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
		case CommitPopUp:
			gitCommitCancelService(m)
		case AddRemotePromptPopUp:
			gitAddRemoteCancelService(m)
		}
		return m, nil
	// in typing mode, tab is move to next input
	case "tab":
		switch m.PopUpType {
		case CommitPopUp:
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
		case AddRemotePromptPopUp:
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
		case CommitPopUp:
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
		case AddRemotePromptPopUp:
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
		case CommitPopUp:
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
		case AddRemotePromptPopUp:
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

		}
		return m, nil
	}
	switch m.PopUpType {
	case CommitPopUp:
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
	case AddRemotePromptPopUp:
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

	}
	return m, nil
}

func handleNonTypingGlobalKeyBindingInteraction(msg tea.KeyMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+c":
		api.GITDAEMON.Stop()
		return m, tea.Quit
	case "b", "B":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedContainer != LocalBranchComponent {
				m.CurrentSelectedContainer = LocalBranchComponent
			} else {
				m.CurrentSelectedContainer = NoneSelected
			}
		}
		return m, nil
	case "f", "F":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedContainer != ModifiedFilesComponent {
				m.CurrentSelectedContainer = ModifiedFilesComponent
			} else {
				m.CurrentSelectedContainer = NoneSelected
			}
		}
		return m, nil
	case "c", "C":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedContainer == ModifiedFilesComponent {
				m.ShowPopUp.Store(true)
				m.PopUpType = CommitPopUp
				m.GitState.GitCommit.ClearGitCommitOutput()

				// if the current pop up model is not commit pop up model, then init it
				if popUp, ok := m.PopUpModel.(*GitCommitPopUpModel); !ok {
					initGitCommitPopUpModel(m)
				} else {
					popUp.GitCommitOutputViewport.SetContent("")
				}
				m.IsTyping.Store(true)
			}
		}
		return m, nil
	case "p", "P":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedContainer == ModifiedFilesComponent || m.CurrentSelectedContainer == NoneSelected || m.CurrentSelectedContainer == LocalBranchComponent {
				// first we need to check if there are any push origin for this repo
				// if not we prompt the user to add a new remote origin
				if !m.GitState.GitCommit.CheckRemoteExist() {
					m.ShowPopUp.Store(true)
					m.PopUpType = AddRemotePromptPopUp
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
						m.PopUpType = ChoosePushTypePopUp
						// if the current pop up model is not commit pop up model, then init it and start git push service
						initChoosePushTypePopUpModel(m, remotes[0].Name)
					} else if len(remotes) > 1 {
						// if remote is more than 1 let user choose which remote to push to first before pushing
						m.PopUpType = ChooseRemotePopUp
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
			if m.CurrentSelectedContainer == ModifiedFilesComponent && len(m.CurrentRepoModifiedFilesInfoList.Items()) > 0 {
				m.CurrentSelectedContainer = FileDiffComponent
			}
		} else {
			switch m.PopUpType {
			case ChooseRemotePopUp:
				popUp, ok := m.PopUpModel.(*ChooseRemotePopUpModel)
				if ok {
					remote := popUp.RemoteList.SelectedItem()
					m.PopUpType = ChoosePushTypePopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					initChoosePushTypePopUpModel(m, remote.(gitRemoteItem).Name)
				}
			case ChoosePushTypePopUp:
				popUp, ok := m.PopUpModel.(*ChoosePushTypePopUpModel)
				if ok {
					m.PopUpType = GitRemotePushPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					selectedOption := popUp.PushOptionList.SelectedItem()
					return initGitRemotePushPopUpModelAndStartGitRemotePushService(m, popUp.RemoteName, selectedOption.(gitPushOptionItem).pushType)
				}
			}
		}
		return m, nil
	case "esc":
		if m.ShowPopUp.Load() {
			switch m.PopUpType {
			case GitRemotePushPopUp:
				gitRemotePushCancelService(m)
			case ChooseRemotePopUp:
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)

			case ChoosePushTypePopUp:
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpModel = nil
			}
			return m, nil
		} else {
			switch m.CurrentSelectedContainer {
			case NoneSelected:
				api.GITDAEMON.Stop()
				return m, tea.Quit
			case FileDiffComponent:
				m.CurrentSelectedContainer = ModifiedFilesComponent
			case LocalBranchComponent:
				m.CurrentSelectedContainer = NoneSelected
			case ModifiedFilesComponent:
				m.CurrentSelectedContainer = NoneSelected
			}
		}
		return m, nil
	case "s", "S":
		if m.CurrentSelectedContainer == ModifiedFilesComponent {
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
			case LocalBranchComponent:
				// we don't use the list native Update() because we track the current selected index
				if m.CurrentRepoBranchesInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoBranchesInfoList.Index() - 1
					m.CurrentRepoBranchesInfoList.Select(latestIndex)
					m.NavigationIndexPosition.LocalBranchComponent = latestIndex
				}
			case ModifiedFilesComponent:
				// we don't use the list native Update() because we need to also render the diff view as well as track the current selected index
				if m.CurrentRepoModifiedFilesInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoModifiedFilesInfoList.Index() - 1
					m.CurrentRepoModifiedFilesInfoList.Select(latestIndex)
					m.NavigationIndexPosition.ModifiedFilesComponent = latestIndex
					reinitAndRenderModifiedFileDiffViewPort(m)
				}
			case FileDiffComponent:
				m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
				return m, cmd
			}
		} else {
			// for within pop up component
			switch m.PopUpType {
			case ChooseRemotePopUp:
				popUp, ok := m.PopUpModel.(*ChooseRemotePopUpModel)
				if ok {
					popUp.RemoteList, cmd = popUp.RemoteList.Update(msg)
					return m, cmd
				}
			case ChoosePushTypePopUp:
				popUp, ok := m.PopUpModel.(*ChoosePushTypePopUpModel)
				if ok {
					popUp.PushOptionList, cmd = popUp.PushOptionList.Update(msg)
					return m, cmd
				}
			}
		}
		return m, nil
	case "down", "j":
		if !m.ShowPopUp.Load() {
			switch m.CurrentSelectedContainer {
			case LocalBranchComponent:
				// we don't use the list native Update() because we track the current selected index
				if m.CurrentRepoBranchesInfoList.Index() < len(m.CurrentRepoBranchesInfoList.Items())-1 {
					latestIndex := m.CurrentRepoBranchesInfoList.Index() + 1
					m.CurrentRepoBranchesInfoList.Select(latestIndex)
					m.NavigationIndexPosition.LocalBranchComponent = latestIndex
				}
			case ModifiedFilesComponent:
				// we don't use the list native Update() because we need to also render the diff view as well as track the current selected index
				if m.CurrentRepoModifiedFilesInfoList.Index() < len(m.CurrentRepoModifiedFilesInfoList.Items())-1 {
					latestIndex := m.CurrentRepoModifiedFilesInfoList.Index() + 1
					m.CurrentRepoModifiedFilesInfoList.Select(latestIndex)
					m.NavigationIndexPosition.ModifiedFilesComponent = latestIndex
					reinitAndRenderModifiedFileDiffViewPort(m)
				}
			case FileDiffComponent:
				m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
				return m, cmd
			}
		} else {
			// for within pop up component
			switch m.PopUpType {
			case ChooseRemotePopUp:
				popUp, ok := m.PopUpModel.(*ChooseRemotePopUpModel)
				if ok {
					popUp.RemoteList, cmd = popUp.RemoteList.Update(msg)
					return m, cmd
				}
			case ChoosePushTypePopUp:
				popUp, ok := m.PopUpModel.(*ChoosePushTypePopUpModel)
				if ok {
					popUp.PushOptionList, cmd = popUp.PushOptionList.Update(msg)
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
			case CommitPopUp:
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
			case CommitPopUp:
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
			case CommitPopUp:
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
			case CommitPopUp:
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
