package tui

import (
	tea "charm.land/bubbletea/v2"
)

func StartGitUpdateListener(p *tea.Program, updateChannel chan string) {
	go func() {
		for updateEvent := range updateChannel {
			// Push message into the Bubble Tea runtime
			p.Send(GitUpdateMsg(updateEvent))
		}
	}()
}

func StartTuiUpdateListener(p *tea.Program, updateChannel chan string) {
	go func() {
		for updateEvent := range updateChannel {
			// Push message into the Bubble Tea runtime
			p.Send(GitUpdateMsg(updateEvent))
		}
	}()
}
