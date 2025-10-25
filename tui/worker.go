package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func StartGitUpdateListener(p *tea.Program, updateChannel chan string) {
	go func() {
		for updateEvent := range updateChannel {
			// Push message into the Bubble Tea runtime
			p.Send(GitUpdateMsg(updateEvent))
		}
	}()
}
