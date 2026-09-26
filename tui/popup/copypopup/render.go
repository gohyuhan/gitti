package copypopup

import (
	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

func Render(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*Model)
	if !ok {
		return ""
	}
	width := min(constant.MaxChooseCopyValuePopUpWidth, int(float64(m.Width)*0.8))
	popUp.Options.SetWidth(width - 4)
	title := style.TitleStyle.Width(width).Render(i18n.LANGUAGEMAPPING.CopyPopupTitle)
	content := lipgloss.JoinVertical(lipgloss.Left, title, popUp.Options.View())
	return style.PopUpBorderStyle.Width(width).Render(content)
}
