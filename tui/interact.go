package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"gitti/api"
	"gitti/api/git"
)

// the function to handle bubbletea key interactions
func GittiKeyInteraction(msg tea.KeyMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
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

	case "enter":
		if m.CurrentSelectedContainer == ModifiedFilesComponent {
			m.CurrentSelectedContainer = FileDiffComponent
		}
		return m, nil

	case "esc":
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
			m.CurrentSelectedFileDiffViewportOffset = max(0, m.CurrentSelectedFileDiffViewportOffset-1)
			m.CurrentSelectedFileDiffViewport.SetXOffset(m.CurrentSelectedFileDiffViewportOffset)
		}
		return m, nil

	case "right", "l":
		switch m.CurrentSelectedContainer {
		case FileDiffComponent:
			if m.CurrentSelectedFileDiffViewport.HorizontalScrollPercent() < 1 {
				m.CurrentSelectedFileDiffViewportOffset = m.CurrentSelectedFileDiffViewportOffset + 1
			}
			m.CurrentSelectedFileDiffViewport.SetXOffset(m.CurrentSelectedFileDiffViewportOffset)
		}
		return m, nil
	}

	return m, nil
}

func GittiMouseInteraction(msg tea.MouseMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "wheel left":
		switch m.CurrentSelectedContainer {
		case FileDiffComponent:
			m.CurrentSelectedFileDiffViewportOffset = max(0, m.CurrentSelectedFileDiffViewportOffset-1)
			m.CurrentSelectedFileDiffViewport.SetXOffset(m.CurrentSelectedFileDiffViewportOffset)
		}
	case "wheel right":
		switch m.CurrentSelectedContainer {
		case FileDiffComponent:
			if m.CurrentSelectedFileDiffViewport.HorizontalScrollPercent() < 1 {
				m.CurrentSelectedFileDiffViewportOffset = m.CurrentSelectedFileDiffViewportOffset + 1
			}
			m.CurrentSelectedFileDiffViewport.SetXOffset(m.CurrentSelectedFileDiffViewportOffset)
		}
	case "wheel up":
		switch m.CurrentSelectedContainer {
		case FileDiffComponent:
			m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
			return m, cmd
		}
	case "wheel down":
		switch m.CurrentSelectedContainer {
		case FileDiffComponent:
			m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}
