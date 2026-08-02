package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	keybindingPopUp "github.com/gohyuhan/gitti/tui/popup/keybinding"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle global non-typing key interactions (e.g., '?').
//	Responsibility: Opens the global Keybinding and Feature Instructions pop-up,
//	suspending the typing state to display application-wide shortcuts and help info.
//
// ------------------------------------
func handleNonTypingGlobalKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	m.ShowPopUp.Store(true)
	m.IsTyping.Store(false)
	m.PopUpType = constant.KeybindingAndFeatureInstructionsPopUp
	keybindingPopUp.InitKeybindingAndFeatureInstructionsPopUpModel(m)
	return m, nil
}
