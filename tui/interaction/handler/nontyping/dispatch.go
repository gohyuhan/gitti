package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/layout"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle key presses in non-typing mode. Dispatches each key (letters, arrows,
//	special keys, and panel-resize +/-) to its dedicated handler function.
//
// ------------------------------------
func Handle(msg tea.KeyPressMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch msg.String() {
	case "?":
		return handleNonTypingGlobalKeyBindingInteraction(m)

	case "1":
		return handleNonTyping1KeyBindingInteraction(m)

	case "2":
		return handleNonTyping2KeyBindingInteraction(m)

	case "3":
		return handleNonTyping3KeyBindingInteraction(m)

	case "4":
		return handleNonTyping4KeyBindingInteraction(m)

	case "a":
		return handleNonTypingaKeyBindingInteraction(m)

	case "A":
		return handleNonTypingAKeyBindingInteraction(m)

	case "b":
		return handleNonTypingbKeyBindingInteraction(m)

	case "c":
		return handleNonTypingcKeyBindingInteraction(m)

	case "C":
		return handleNonTypingCKeyBindingInteraction(m)

	case "d":
		return handleNonTypingdKeyBindingInteraction(m)

	case "e":
		return handleNonTypingeKeyBindingInteraction(m)

	case "f":
		return handleNonTypingfKeyBindingInteraction(m)

	case "F":
		return handleNonTypingFKeyBindingInteraction(m)

	case "i":
		return handleNonTypingiKeyBindingInteraction(m)

	case "L":
		return handleNonTypingLKeyBindingInteraction(m)

	case "m":
		return handleNonTypingmKeyBindingInteraction(m)

	case "n":
		return handleNonTypingnKeyBindingInteraction(m)

	case "o":
		return handleNonTypingoKeyBindingInteraction(m)

	case "O":
		return handleNonTypingOKeyBindingInteraction(m)

	case "p":
		return handleNonTypingpKeyBindingInteraction(m)

	case "P":
		return handleNonTypingPKeyBindingInteraction(m)

	case "r":
		return handleNonTypingrKeyBindingInteraction(m)

	case "R":
		return handleNonTypingRKeyBindingInteraction(m)

	case "s":
		return handleNonTypingsKeyBindingInteraction(m)

	case "S":
		return handleNonTypingSKeyBindingInteraction(m)

	case "t":
		return handleNonTypingtKeyBindingInteraction(m)

	case "[":
		return handleNonTypingLeftBracketKeyBindingInteraction(m)

	case "]":
		return handleNonTypingRightBracketKeyBindingInteraction(m)

	case "/":
		return handleNonTypingSlashKeyBindingInteraction(m)

	case "q", "Q":
		// only work when there is no pop up
		return handleNonTypingqQKeyBindingInteraction(m)

	case "backspace":
		return handleNonTypingBackspaceKeyBindingInteraction(m)

	case "enter":
		return handleNonTypingEnterKeyBindingInteraction(m)

	case "tab":
		// next component navigation
		return handleNonTypingTabKeyBindingInteraction(m)

	case "shift+tab":
		// previous component navigation
		return handleNonTypingShiftTabKeyBindingInteraction(m)

	case "space":
		return handleNonTypingSpaceKeyBindingInteraction(m)

	case "esc":
		return handleNonTypingEscKeyBindingInteraction(m)

	case "up", "k":
		return handleNonTypingUpkKeyBindingInteraction(msg, m)

	case "down", "j":
		return handleNonTypingDownjKeyBindingInteraction(msg, m)

	case "left", "h":
		return handleNonTypingLefthKeyBindingInteraction(m)

	case "right", "l":
		return handleNonTypingRightlKeyBindingInteraction(m)
	case "-":
		if !m.ShowPopUp.Load() {
			m.WindowLeftPanelRatio = max(settings.MINLEFTPANELWIDTHRATIO, m.WindowLeftPanelRatio-0.01)
			layout.TuiWindowSizing(m)
		}
		return m, nil
	case "+":
		if !m.ShowPopUp.Load() {
			m.WindowLeftPanelRatio = min(settings.MAXLEFTPANELWIDTHRATIO, m.WindowLeftPanelRatio+0.01)
			layout.TuiWindowSizing(m)
		}
		return m, nil
	case "<":
		return handleNonTypingLeftAngleBracketKeyBindingInteraction(m)

	case ">":
		return handleNonTypingRightAngleBracketKeyBindingInteraction(m)

	case "ctrl+a":
		return handleNonTypingCtrlaKeyBindingInteraction(m)
	case "ctrl+p":
		return handleNonTypingCtrlpKeyBindingInteraction(m)
	case "ctrl+k":
		return handleNonTypingCtrlkKeyBindingInteraction(m)
	case "ctrl+r":
		return handleNonTypingCtrlrKeyBindingInteraction(m)
	}
	return m, nil
}
