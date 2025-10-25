package tui

import (
	"gitti/api"

	tea "github.com/charmbracelet/bubbletea"
)

// the function to handle bubbletea key interactions
func GittiKeyInteraction(msg tea.KeyMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "Q", "esc":
		api.GITDAEMON.Stop()
		return m, tea.Quit
	case "b", "B":
		if m.CurrentSelectedContainer != localBranchComponent {
			m.CurrentSelectedContainer = localBranchComponent
		} else {
			m.CurrentSelectedContainer = None
		}
		return m, nil
	case "up", "k":
		if m.CurrentSelectedContainer == localBranchComponent && m.CurrentRepoBranchesInfo.Index() > 0 {
			latestIndex := m.CurrentRepoBranchesInfo.Index() - 1
			m.CurrentRepoBranchesInfo.Select(latestIndex)
			m.NavigationIndexPosition.LocalBranchComponent = latestIndex
		}
		return m, nil
	case "down", "j":
		if m.CurrentSelectedContainer == localBranchComponent && m.CurrentRepoBranchesInfo.Index() < len(m.CurrentRepoBranchesInfo.Items())-1 {
			latestIndex := m.CurrentRepoBranchesInfo.Index() + 1
			m.CurrentRepoBranchesInfo.Select(latestIndex)
			m.NavigationIndexPosition.LocalBranchComponent = latestIndex
		}
		return m, nil
	}
	return m, nil
}
