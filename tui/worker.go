package tui

import (
	tea "charm.land/bubbletea/v2"
)

func StartGitUpdateListener(p *tea.Program, updateReceiverChannel chan string) {
	go func() {
		for updateEvent := range updateReceiverChannel {
			// Push message into the Bubble Tea runtime
			p.Send(GitUpdateMsg(updateEvent))
		}
	}()
}

func StartTuiUpdateListener(p *tea.Program, updateReceiverChannel chan string) {
	go func() {
		for updateEvent := range updateReceiverChannel {
			// Push message into the Bubble Tea runtime
			p.Send(GitUpdateMsg(updateEvent))
		}
	}()
}

func StartLoggingUpdateListener(p *tea.Program, updateReceiverChannel chan string) {
	go func() {
		for updateEvent := range updateReceiverChannel {
			// Push message into the Bubble Tea runtime
			p.Send(GitUpdateMsg(updateEvent))
		}
	}()
}
