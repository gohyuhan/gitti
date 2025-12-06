package tui

import (
	"github.com/gohyuhan/gitti/tui/types"
)

// ---------------------------------
//
// # Main Model for TUI
//
// ---------------------------------
type GittiAppModel struct {
	model *types.GittiModel
}

// ---------------------------------
//
// tea msg
//
// ---------------------------------
type GitUpdateMsg string
