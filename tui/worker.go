package tui

import (
	tea "charm.land/bubbletea/v2"
)

// ----------------------------------
//
//	Listen for git update events and forward them to the TUI program
//
// ----------------------------------
func StartGitUpdateListener(p *tea.Program, updateReceiverChannel chan string) {
	go func() {
		for updateEvent := range updateReceiverChannel {
			// Push message into the Bubble Tea runtime
			p.Send(GitUpdateMsg(updateEvent))
		}
	}()
}

// ----------------------------------
//
//	Listen for TUI update events and forward them to the TUI program
//
// ----------------------------------
func StartTuiUpdateListener(p *tea.Program, updateReceiverChannel chan string) {
	go func() {
		for updateEvent := range updateReceiverChannel {
			// Push message into the Bubble Tea runtime
			p.Send(GitUpdateMsg(updateEvent))
		}
	}()
}

// ----------------------------------
//
//	Listen for logging update events and forward them to the TUI program
//
// ----------------------------------
func StartLoggingUpdateListener(p *tea.Program, updateReceiverChannel chan string) {
	go func() {
		for updateEvent := range updateReceiverChannel {
			// Push message into the Bubble Tea runtime
			p.Send(GitUpdateMsg(updateEvent))
		}
	}()
}
