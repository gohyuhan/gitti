package interactiverebase

import (
	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	choose interactive rebase option
//
// ------------------------------------
func RenderInteractiveRebaseOptionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseOptionPopUp)
	if ok {
		popUpWidth := min(constant.MaxInteractiveRebaseOptionPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.ChooseInteractiveRebaseOption)
		popUp.InteractiveRebaseOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.InteractiveRebaseOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
