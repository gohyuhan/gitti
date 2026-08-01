package handler

import (
	tea "charm.land/bubbletea/v2"

	"github.com/gohyuhan/gitti/tui/interaction/handler/nontyping"
	"github.com/gohyuhan/gitti/tui/interaction/handler/typing"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Entry point for key presses while an input is focused.
//	The per-key handlers live in the typing package, one file per key.
//
// ------------------------------------
func HandleTypingKeyBindingInteraction(msg tea.KeyPressMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	return typing.Handle(msg, m)
}

// ------------------------------------
//
//	Entry point for key presses while no input is focused.
//	The per-key handlers live in the nontyping package, one file per key.
//
// ------------------------------------
func HandleNonTypingGlobalKeyBindingInteraction(msg tea.KeyPressMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	return nontyping.Handle(msg, m)
}
