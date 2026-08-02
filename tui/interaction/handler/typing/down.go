package typing

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	blamePopUp "github.com/gohyuhan/gitti/tui/popup/blame"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle down arrow in typing mode. In the blame popup, navigates the file
//	list down when the file selector is active, or scrolls the blame viewport
//	down when blame info is shown.
//
// ------------------------------------
func handleTypingDownKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.BlamePopUp:
			popUp, ok := m.PopUpModel.(*blamePopUp.BlamePopUpModel)
			if ok {
				if !popUp.ShowingBlameInfo {
					popUp.CurrentGitTrackedFilesPathList.CursorDown()
				} else {
					popUp.BlameViewport.ScrollDown(1)
				}
			}
		}
	}

	return m, nil
}
