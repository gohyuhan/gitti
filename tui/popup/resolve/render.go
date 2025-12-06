package resolve

import (
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"

	"charm.land/lipgloss/v2"
)

// ------------------------------------
//
//	For resolve conflict option list
//
// ------------------------------------
func RenderGitResolveConflictOptionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitResolveConflictOptionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitResolveConflictOptionPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitResolveConflictOptionTitle)

		popUp.ResolveConflictOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.ResolveConflictOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
