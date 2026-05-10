package files

import (
	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Render the discard-line-change confirmation popup, showing the localized
//	title and a viewport with the specific diff hunk the user is about to discard.
//
// ------------------------------------
func RenderGitDiscardFileLineChangeConfirmPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitDiscardFileLineChangeConfirmPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitDiscardFileLineChangeConfirmPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitDiscardFileLineChangeConfirmTitle)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.DiscardFileLineChangeViewport.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
