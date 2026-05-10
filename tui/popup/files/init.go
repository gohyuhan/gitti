package files

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Initialize the discard-line-change confirmation popup. Creates a non-soft-wrap
//	viewport pre-filled with the diff line at the cursor position from whichever
//	detail panel is currently selected (DetailComponentPanel or DetailComponentPanelTwo).
//
// ------------------------------------
func InitGitDiscardFileLineChangeConfirmPopUpModel(m *types.GittiModel) {
	vp := viewport.New()
	vp.SoftWrap = false
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpGitDiscardFileLineChangeViewportHeight)
	vp.SetWidth(min(constant.MaxGitDiscardConfirmPromptPopupWidth, int(float64(m.Width)*0.8)) - 4)

	var discardFileLineChange strings.Builder
	discardFileLineChange.WriteString("  ")

	switch m.CurrentSelectedComponent {
	case constant.DetailComponentPanel:
		ogArrayIndex := m.LineEditingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex - m.LineEditingIndexPositionAndInfo.DetailPanelViewportOverflowIndexCount
		discardFileLineChange.WriteString(m.DetailPanelViewportOGStringArray[ogArrayIndex])
	case constant.DetailComponentPanelTwo:
		ogArrayIndex := m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportActualCurrentIndex - m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportOverflowIndexCount
		discardFileLineChange.WriteString(m.DetailPanelTwoViewportOGStringArray[ogArrayIndex])
	}

	vp.SetContent(discardFileLineChange.String())

	popUpModel := &GitDiscardFileLineChangeConfirmPopUpModel{
		DiscardFileLineChangeViewport: vp,
	}
	m.PopUpModel = popUpModel
}
