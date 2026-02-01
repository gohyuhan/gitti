package interaction

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/interaction/handler"
	filesPopUp "github.com/gohyuhan/gitti/tui/popup/files"
	keybindingPopUp "github.com/gohyuhan/gitti/tui/popup/keybinding"
	"github.com/gohyuhan/gitti/tui/types"
)

func GittiMouseInteraction(msg tea.MouseMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "wheelleft":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedComponent == constant.DetailComponentTwo {
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
			case constant.GlobalKeyBindingPopUp:
				popUp, ok := m.PopUpModel.(*keybindingPopUp.GlobalKeyBindingPopUpModel)
				if ok {
					scrollSpeed := 1
					if strings.ToUpper(settings.GITTICONFIGSETTINGS.LanguageCode) != "EN" {
						// other than en, all other i18n we support are both zh and jp which each rune takes up twice the width of en character
						scrollSpeed = 2
					}
					popUp.GlobalKeyBindingViewport.ScrollLeft(scrollSpeed)
				}
			}
		}

	case "wheelright":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedComponent == constant.DetailComponentTwo {
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
			case constant.GlobalKeyBindingPopUp:
				popUp, ok := m.PopUpModel.(*keybindingPopUp.GlobalKeyBindingPopUpModel)
				if ok {
					scrollSpeed := 1
					if strings.ToUpper(settings.GITTICONFIGSETTINGS.LanguageCode) != "EN" {
						// other than en, all other i18n we support are both zh and jp which each rune takes up twice the width of en character
						scrollSpeed = 2
					}
					popUp.GlobalKeyBindingViewport.ScrollRight(scrollSpeed)
				}
			}
		}

	case "wheelup":
		if !m.ShowPopUp.Load() {
			if !m.IsLineEditingState.Load() {
				if m.CurrentSelectedComponent == constant.DetailComponentTwo {
					m.DetailPanelTwoViewport, cmd = m.DetailPanelTwoViewport.Update(msg)
				} else {
					m.DetailPanelViewport, cmd = m.DetailPanelViewport.Update(msg)
				}
				return m, cmd
			}
			return m, nil
		} else {
			return handler.UpDownMouseMsgUpdateForPopUp(msg, m)
		}

	case "wheeldown":
		if !m.ShowPopUp.Load() {
			if !m.IsLineEditingState.Load() {
				if m.CurrentSelectedComponent == constant.DetailComponentTwo {
					m.DetailPanelTwoViewport, cmd = m.DetailPanelTwoViewport.Update(msg)
				} else {
					m.DetailPanelViewport, cmd = m.DetailPanelViewport.Update(msg)
				}
				return m, cmd
			}
			return m, nil
		} else {
			return handler.UpDownMouseMsgUpdateForPopUp(msg, m)
		}
	}
	return m, nil
}
