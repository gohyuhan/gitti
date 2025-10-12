package tui

import (
	"gitti/api"

	tea "github.com/charmbracelet/bubbletea"
)

func StartGitUpdateListener(p *tea.Program, w *api.GittiDaemonWorker) {
	go func() {
		for gitInfo := range w.ListenToUpdateChannel() {
			// Push message into the Bubble Tea runtime
			p.Send(GitUpdateMsg(gitInfo))
		}
	}()
}
