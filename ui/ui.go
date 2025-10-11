package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func NewGittiModel() GittiModel {
	return GittiModel{
		width: 90,
		height: 30,
	}
}

// -----------------------------------------------------------------------------
// Bubble Tea standard functions
// -----------------------------------------------------------------------------

func (m GittiModel) Init() tea.Cmd {
	return nil
}


func (m GittiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "Q", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m GittiModel) View() string {
	return GittiMainPageView(m)
}
