package nontyping

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/constant"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	filesPopUp "github.com/gohyuhan/gitti/tui/popup/files"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	keybindingPopUp "github.com/gohyuhan/gitti/tui/popup/keybinding"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle Left/h key interaction.
//	Responsibility: Horizontal left navigation/scrolling.
//	Scrolls text viewers, detail panels, and wider popup views horizontally to the left.
//
// ------------------------------------
func handleNonTypingLefthKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.DetailComponentPanel:
			m.DetailPanelViewport.ScrollLeft(1)
		case constant.DetailComponentPanelTwo:
			m.DetailPanelTwoViewport.ScrollLeft(1)
		default:
			m.DetailPanelViewport.ScrollLeft(1)
		}
	} else {
		switch m.PopUpType {
		case constant.CommitPopUp:
			popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
			if ok {
				popUp.GitCommitOutputViewport.ScrollLeft(1)
			}
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
		case constant.InteractiveRebaseFixupSquashSelectionPopUp:
			popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashSelectionPopUpModel)
			if ok && popUp.IsCommitFixupSquashViewportSelected {
				popUp.CommitFixupSquashViewport.ScrollLeft(1)
			}

		case constant.InteractiveRebaseFixupSquashOutputPopUp:
			popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashOutputPopUpModel)
			if ok && !popUp.IsProcessing.Load() {
				popUp.FixupSquashOutputViewport.ScrollLeft(1)
			}
		}
	}
	return m, nil
}
