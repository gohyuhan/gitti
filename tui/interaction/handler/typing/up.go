package typing

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	blamePopUp "github.com/gohyuhan/gitti/tui/popup/blame"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle up arrow in typing mode. In the blame popup, navigates the file list
//	up when the file selector is active, or scrolls the blame viewport up when
//	blame info is shown.
//
// ------------------------------------
func handleTypingUpKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.BlamePopUp:
			popUp, ok := m.PopUpModel.(*blamePopUp.BlamePopUpModel)
			if ok {
				if !popUp.ShowingBlameInfo {
					popUp.CurrentGitTrackedFilesPathList.CursorUp()
				} else {
					popUp.BlameViewport.ScrollUp(1)
				}
			}
		}
	}

	return m, nil
}
