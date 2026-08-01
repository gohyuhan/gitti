package typing

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	blamePopUp "github.com/gohyuhan/gitti/tui/popup/blame"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle right arrow in typing mode. In the blame popup, scrolls the blame
//	viewport right when blame info is shown.
//
// ------------------------------------
func handleTypingRightKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.BlamePopUp:
			popUp, ok := m.PopUpModel.(*blamePopUp.BlamePopUpModel)
			if ok {
				if popUp.ShowingBlameInfo {
					popUp.BlameViewport.ScrollRight(1)
				}
			}
		}
	}

	return m, nil
}
