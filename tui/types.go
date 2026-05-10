package tui

import (
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	GittiAppModel is the top-level Bubbletea model that wraps the internal
//	GittiModel pointer and implements the tea.Model interface.
//
// ------------------------------------
type GittiAppModel struct {
	model *types.GittiModel
}

// ------------------------------------
//
//	GitUpdateMsg carries a git daemon event string sent from the background
//	listener goroutine into the Bubbletea runtime via p.Send.
//
// ------------------------------------
type GitUpdateMsg string
