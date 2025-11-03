package tui

import (
	"gitti/api"
	"gitti/api/git"

	tea "github.com/charmbracelet/bubbletea/v2"
)

// the function to handle bubbletea key interactions
func gittiKeyInteraction(msg tea.KeyMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	if m.IsTyping {
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
				if !popUp.IsProcessing {
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
				if !popUp.IsProcessing {
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
		if m.CurrentSelectedContainer != LocalBranchComponent {
			m.CurrentSelectedContainer = LocalBranchComponent
		} else {
			m.CurrentSelectedContainer = NoneSelected
		}
		return m, nil
	case "f", "F":
		if m.CurrentSelectedContainer != ModifiedFilesComponent {
			m.CurrentSelectedContainer = ModifiedFilesComponent
		} else {
			m.CurrentSelectedContainer = NoneSelected
		}
		return m, nil
	case "c", "C":
		if m.CurrentSelectedContainer == ModifiedFilesComponent {
			m.ShowPopUp = true
			m.PopUpType = CommitPopUp
			// if the current pop up model is not commit pop up model, then init it
			if _, ok := m.PopUpModel.(*GitCommitPopUpModel); !ok {
				initGitCommitPopUpModel(m)
			}
			m.IsTyping = true
		}
		return m, nil
	case "p", "P":
		if m.CurrentSelectedContainer == ModifiedFilesComponent || m.CurrentSelectedContainer == NoneSelected || m.CurrentSelectedContainer == LocalBranchComponent {
			// first we need to check if there are any push origin for this repo
			// if not we prompt the user to add a new remote origin
			if !git.GITCOMMIT.CheckRemoteExist() {
				m.ShowPopUp = true
				m.PopUpType = AddRemotePromptPopUp
				// if the current pop up model is not commit pop up model, then init it
				if _, ok := m.PopUpModel.(*AddRemotePromptPopUpModel); !ok {
					initAddRemotePromptPopUpModel(m, true)
				}
				m.IsTyping = true
			} else {
				m.ShowPopUp = true
				if len(git.GITCOMMIT.Remote) == 1 {
					m.PopUpType = GitRemotePushPopUp
					// if the current pop up model is not commit pop up model, then init it and start git push service
					return initGitRemotePushPopUpModelAndStartGitRemotePushService(m, git.GITCOMMIT.Remote[0].Name)
				} else if len(git.GITCOMMIT.Remote) > 1 {
					// if remote is more than 1 let user choose which remote to push to first before pushing
					m.PopUpType = ChooseRemotePopUp
					if _, ok := m.PopUpModel.(*ChooseRemotePopUpModel); !ok {
						initGitRemotePushChooseRemotePopUpModel(m, git.GITCOMMIT.Remote)
					}
				}
			}
		}
		return m, nil

	case "enter":
		if !m.ShowPopUp {
			if m.CurrentSelectedContainer == ModifiedFilesComponent && len(m.CurrentRepoModifiedFilesInfoList.Items()) > 0 {
				m.CurrentSelectedContainer = FileDiffComponent
			}
		} else {
			switch m.PopUpType {
			case ChooseRemotePopUp:
				popUp, ok := m.PopUpModel.(*ChooseRemotePopUpModel)
				if ok {
					remote := popUp.RemoteList.SelectedItem()
					m.PopUpType = GitRemotePushPopUp
					m.ShowPopUp = true
					return initGitRemotePushPopUpModelAndStartGitRemotePushService(m, remote.(gitRemoteItem).Name)
				}
			}
		}
		return m, nil
	case "esc":
		if m.ShowPopUp {
			switch m.PopUpType {
			case GitRemotePushPopUp:
				gitRemotePushCancelService(m)
			case ChooseRemotePopUp:
				m.ShowPopUp = false
				m.IsTyping = false
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
				git.GITFILES.ToggleFilesStageStatus(fileName)
			}
		}
		return m, nil
	case "up", "k":
		if !m.ShowPopUp {
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
			}
		}
		return m, nil
	case "down", "j":
		if !m.ShowPopUp {
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
			}
		}
		return m, nil
	case "left", "h":
		if !m.ShowPopUp {
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
		if !m.ShowPopUp {
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
		if !m.ShowPopUp {
			m.CurrentSelectedFileDiffViewport.MoveLeft(1)
		}
	case "wheelright":
		if !m.ShowPopUp {
			m.CurrentSelectedFileDiffViewport.MoveRight(1)
		}
	case "wheelup":
		if !m.ShowPopUp {
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
		if !m.ShowPopUp {
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
