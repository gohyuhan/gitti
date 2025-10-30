package tui

import (
	"gitti/api"
	"gitti/api/git"

	tea "github.com/charmbracelet/bubbletea/v2"
)

// the function to handle bubbletea key interactions
func GittiKeyInteraction(msg tea.KeyMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	if m.IsTyping {
		return handleTypingKeyBindingInteraction(msg, m)
	} else {
		return handleNonTypingGlobalKeyBindingInteraction(msg, m)
	}
}

func handleTypingKeyBindingInteraction(msg tea.KeyMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		git.GITCOMMIT.KillCommit()
		m.ShowPopUp = false
		m.IsTyping = false
		return m, nil
	// in typing mode, tab is move to next input
	case "tab":
		switch m.PopUpType {
		case CommitPopUp:
			m.PopUpModel.(*GitCommitPopUpModel).CurrentActiveInputIndex = min(m.PopUpModel.(*GitCommitPopUpModel).CurrentActiveInputIndex+1, m.PopUpModel.(*GitCommitPopUpModel).TotalInputCount)
			switch m.PopUpModel.(*GitCommitPopUpModel).CurrentActiveInputIndex {
			case 1:
				m.PopUpModel.(*GitCommitPopUpModel).MessageTextInput.Focus()
				m.PopUpModel.(*GitCommitPopUpModel).DescriptionTextAreaInput.Blur()
			case 2:
				m.PopUpModel.(*GitCommitPopUpModel).MessageTextInput.Blur()
				m.PopUpModel.(*GitCommitPopUpModel).DescriptionTextAreaInput.Focus()
			}
		}
		return m, nil

	// in typing mode, shift+tab is move to previous input
	case "shift+tab":
		switch m.PopUpType {
		case CommitPopUp:
			m.PopUpModel.(*GitCommitPopUpModel).CurrentActiveInputIndex = max(m.PopUpModel.(*GitCommitPopUpModel).CurrentActiveInputIndex-1, 1)
			switch m.PopUpModel.(*GitCommitPopUpModel).CurrentActiveInputIndex {
			case 1:
				m.PopUpModel.(*GitCommitPopUpModel).MessageTextInput.Focus()
				m.PopUpModel.(*GitCommitPopUpModel).DescriptionTextAreaInput.Blur()
			case 2:
				m.PopUpModel.(*GitCommitPopUpModel).MessageTextInput.Blur()
				m.PopUpModel.(*GitCommitPopUpModel).DescriptionTextAreaInput.Focus()
			}
		}

	case "ctrl+enter":
		switch m.PopUpType {
		case CommitPopUp:
			// once they start for commit process, reinit the input focus
			m.PopUpModel.(*GitCommitPopUpModel).MessageTextInput.Focus()
			m.PopUpModel.(*GitCommitPopUpModel).CurrentActiveInputIndex = 1
			// start a seperate thread that stage the current selected files and commit them and set the value of msg and desc to "" if committed successfully
			// also do not start any git operation is message is no provided
			if !m.PopUpModel.(*GitCommitPopUpModel).IsProcessing {
				go func() {
					m.PopUpModel.(*GitCommitPopUpModel).IsProcessing = true
					message := m.PopUpModel.(*GitCommitPopUpModel).MessageTextInput.Value()
					description := m.PopUpModel.(*GitCommitPopUpModel).DescriptionTextAreaInput.Value()
					if len(message) < 1 {
						return
					}
					git.GITCOMMIT.GitStage()
					exitStatusCode := git.GITCOMMIT.GitCommit(message, description)
					if exitStatusCode == 0 {
						m.PopUpModel.(*GitCommitPopUpModel).MessageTextInput.SetValue("")
						m.PopUpModel.(*GitCommitPopUpModel).DescriptionTextAreaInput.SetValue("")
					}
				}()
			}
		}
		return m, nil
	}
	switch m.PopUpType {
	case CommitPopUp:
		commitPopUpModel := m.PopUpModel.(*GitCommitPopUpModel)

		switch commitPopUpModel.CurrentActiveInputIndex {
		case 1:
			var cmd tea.Cmd
			commitPopUpModel.MessageTextInput, cmd = commitPopUpModel.MessageTextInput.Update(msg)
			return m, cmd

		case 2:
			var cmd tea.Cmd
			commitPopUpModel.DescriptionTextAreaInput, cmd = commitPopUpModel.DescriptionTextAreaInput.Update(msg)
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
				InitGitCommitPopUpModel(m)
			}
			m.IsTyping = true
		}
		return m, nil
	case "enter":
		if m.CurrentSelectedContainer == ModifiedFilesComponent {
			m.CurrentSelectedContainer = FileDiffComponent
		}
		return m, nil
	case "esc":
		if m.ShowPopUp {
			m.ShowPopUp = false
			m.PopUpType = None
			m.PopUpModel = struct{}{}
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
		switch m.CurrentSelectedContainer {
		case LocalBranchComponent:
			if m.CurrentRepoBranchesInfoList.Index() > 0 {
				latestIndex := m.CurrentRepoBranchesInfoList.Index() - 1
				m.CurrentRepoBranchesInfoList.Select(latestIndex)
				m.NavigationIndexPosition.LocalBranchComponent = latestIndex
			}
		case ModifiedFilesComponent:
			if m.CurrentRepoModifiedFilesInfoList.Index() > 0 {
				latestIndex := m.CurrentRepoModifiedFilesInfoList.Index() - 1
				m.CurrentRepoModifiedFilesInfoList.Select(latestIndex)
				m.NavigationIndexPosition.ModifiedFilesComponent = latestIndex
				ReinitAndRenderModifiedFileDiffViewPort(m)
			}
		case FileDiffComponent:
			m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
			return m, cmd
		}
		return m, nil
	case "down", "j":
		switch m.CurrentSelectedContainer {
		case LocalBranchComponent:
			if m.CurrentRepoBranchesInfoList.Index() < len(m.CurrentRepoBranchesInfoList.Items())-1 {
				latestIndex := m.CurrentRepoBranchesInfoList.Index() + 1
				m.CurrentRepoBranchesInfoList.Select(latestIndex)
				m.NavigationIndexPosition.LocalBranchComponent = latestIndex

			}
		case ModifiedFilesComponent:
			if m.CurrentRepoModifiedFilesInfoList.Index() < len(m.CurrentRepoModifiedFilesInfoList.Items())-1 {
				latestIndex := m.CurrentRepoModifiedFilesInfoList.Index() + 1
				m.CurrentRepoModifiedFilesInfoList.Select(latestIndex)
				m.NavigationIndexPosition.ModifiedFilesComponent = latestIndex
				ReinitAndRenderModifiedFileDiffViewPort(m)
			}
		case FileDiffComponent:
			m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
			return m, cmd
		}
		return m, nil

	case "left", "h":
		switch m.CurrentSelectedContainer {
		case FileDiffComponent:
			m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
			return m, cmd
		}
		return m, nil

	case "right", "l":
		switch m.CurrentSelectedContainer {
		case FileDiffComponent:
			m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
			return m, cmd
		}
		return m, nil
	}
	return m, nil
}

func GittiMouseInteraction(msg tea.MouseMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "wheelleft":
		switch m.CurrentSelectedContainer {
		case FileDiffComponent:
			// mouse left right scroll doesn't seem to be supported yet
			m.CurrentSelectedFileDiffViewportOffset = max(0, m.CurrentSelectedFileDiffViewportOffset-1)
			m.CurrentSelectedFileDiffViewport.SetXOffset(m.CurrentSelectedFileDiffViewportOffset)
		}
	case "wheelright":
		switch m.CurrentSelectedContainer {
		case FileDiffComponent:
			// mouse left right scroll doesn't seem to be supported yet
			if m.CurrentSelectedFileDiffViewport.HorizontalScrollPercent() < 1 {
				m.CurrentSelectedFileDiffViewportOffset = m.CurrentSelectedFileDiffViewportOffset + 1
			}
			m.CurrentSelectedFileDiffViewport.SetXOffset(m.CurrentSelectedFileDiffViewportOffset)
		}
	case "wheelup":
		switch m.CurrentSelectedContainer {
		case FileDiffComponent:
			m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
			return m, cmd
		}
	case "wheeldown":
		switch m.CurrentSelectedContainer {
		case FileDiffComponent:
			m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}
