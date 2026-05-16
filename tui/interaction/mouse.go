package interaction

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/interaction/handler"
	blamePopUp "github.com/gohyuhan/gitti/tui/popup/blame"
	filesPopUp "github.com/gohyuhan/gitti/tui/popup/files"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	keybindingPopUp "github.com/gohyuhan/gitti/tui/popup/keybinding"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle mouse interactions.
//	Responsibility: Translates mouse wheel events (up, down, left, right) into
//	scrolling actions. Depending on the current view state:
//	- Active Popup: Scrolls within popups (instructions, discard confirmation, interactive rebase fixup/squash detail viewport, etc.).
//	- Detail Panels: Horizontally or vertically scrolls the main viewports.
//
// ------------------------------------
func GittiMouseInteraction(msg tea.MouseMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch msg.Mouse().Button {
	case tea.MouseWheelLeft:
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedComponent == constant.DetailComponentPanelTwo {
				m.DetailPanelTwoViewport.ScrollLeft(1)
			} else {
				m.DetailPanelViewport.ScrollLeft(1)
			}
		} else {
			switch m.PopUpType {
			case constant.GitDiscardFileLineChangeConfirmPopUp:
				popUp, ok := m.PopUpModel.(*filesPopUp.GitDiscardFileLineChangeConfirmPopUpModel)
				if ok {
					popUp.DiscardFileLineChangeViewport.ScrollLeft(1)
				}
			case constant.KeybindingAndFeatureInstructionsPopUp:
				popUp, ok := m.PopUpModel.(*keybindingPopUp.KeybindingAndFeatureInstructionsPopUpModel)
				if ok {
					scrollSpeed := 1
					if strings.ToUpper(settings.GITTICONFIGSETTINGS.LanguageCode) != "EN" {
						// other than en, all other i18n we support are both zh and jp which each rune takes up twice the width of en character
						scrollSpeed = 2
					}
					popUp.GlobalKeyBindingViewport.ScrollLeft(scrollSpeed)
				}
			case constant.BlamePopUp:
				popUp, ok := m.PopUpModel.(*blamePopUp.BlamePopUpModel)
				if ok && popUp.ShowingBlameInfo {
					popUp.BlameViewport.ScrollLeft(1)
				}
			case constant.InteractiveRebaseFixupSquashSelectionPopUp:
				popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashSelectionPopUpModel)
				if ok && popUp.IsCommitFixupSquashViewportSelected {
					popUp.CommitFixupSquashViewport.ScrollLeft(1)
				}
			case constant.InteractiveRebaseFixupSquashOutputPopUp:
				popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashOutputPopUpModel)
				if ok {
					popUp.FixupSquashOutputViewport.ScrollLeft(1)
				}
			case constant.InteractiveRebaseRewordOutputPopUp:
				popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseRewordOutputPopUpModel)
				if ok {
					popUp.RewordOutputViewport.ScrollLeft(1)
				}
			}
		}

	case tea.MouseWheelRight:
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedComponent == constant.DetailComponentPanelTwo {
				m.DetailPanelTwoViewport.ScrollRight(1)
			} else {
				m.DetailPanelViewport.ScrollRight(1)
			}
		} else {
			switch m.PopUpType {
			case constant.GitDiscardFileLineChangeConfirmPopUp:
				popUp, ok := m.PopUpModel.(*filesPopUp.GitDiscardFileLineChangeConfirmPopUpModel)
				if ok {
					popUp.DiscardFileLineChangeViewport.ScrollRight(1)
				}
			case constant.KeybindingAndFeatureInstructionsPopUp:
				popUp, ok := m.PopUpModel.(*keybindingPopUp.KeybindingAndFeatureInstructionsPopUpModel)
				if ok {
					scrollSpeed := 1
					if strings.ToUpper(settings.GITTICONFIGSETTINGS.LanguageCode) != "EN" {
						// other than en, all other i18n we support are both zh and jp which each rune takes up twice the width of en character
						scrollSpeed = 2
					}
					popUp.GlobalKeyBindingViewport.ScrollRight(scrollSpeed)
				}
			case constant.BlamePopUp:
				popUp, ok := m.PopUpModel.(*blamePopUp.BlamePopUpModel)
				if ok && popUp.ShowingBlameInfo {
					popUp.BlameViewport.ScrollRight(1)
				}
			case constant.InteractiveRebaseFixupSquashSelectionPopUp:
				popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashSelectionPopUpModel)
				if ok && popUp.IsCommitFixupSquashViewportSelected {
					popUp.CommitFixupSquashViewport.ScrollRight(1)
				}
			case constant.InteractiveRebaseFixupSquashOutputPopUp:
				popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashOutputPopUpModel)
				if ok {
					popUp.FixupSquashOutputViewport.ScrollRight(1)
				}
			case constant.InteractiveRebaseRewordOutputPopUp:
				popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseRewordOutputPopUpModel)
				if ok {
					popUp.RewordOutputViewport.ScrollRight(1)
				}
			}
		}

	case tea.MouseWheelUp, tea.MouseWheelDown:
		if !m.ShowPopUp.Load() {
			if !m.IsLineEditingState.Load() {
				var vpToBeScrolled *viewport.Model
				if m.CurrentSelectedComponent == constant.DetailComponentPanelTwo {
					vpToBeScrolled = &m.DetailPanelTwoViewport
				} else {
					vpToBeScrolled = &m.DetailPanelViewport
				}

				switch msg.Mouse().Button {
				case tea.MouseWheelUp:
					vpToBeScrolled.ScrollUp(1)
				case tea.MouseWheelDown:
					vpToBeScrolled.ScrollDown(1)
				}
			}
			return m, nil
		} else {
			return handler.UpDownMouseMsgUpdateForPopUp(msg, m)
		}
	}
	return m, nil
}
