package log

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

func RenderGitCherryPickOptionSelectionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitCherryPickOptionSelectionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitCherryPickOptionSelectionPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.CherryPickOpsSelectionTitle)
		popUp.CherryPickedOpsOption.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.CherryPickedOpsOption.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

func RenderGitCherryPickPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitCherryPickPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitCherryPickPopUpWidth, int(float64(m.Width)*0.8))
		popUpHeight := max(constant.PopUpGitCherryPickPopUpHeight, int(float64(m.Height)*0.5)-3)
		title := style.TitleStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.CherryPickTitle, popUp.CurrentBranchName))
		popUp.CurrentBranchCherryPickCommitLog.SetWidth(popUpWidth - 4)
		popUp.CurrentBranchCherryPickCommitLog.SetHeight(popUpHeight)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.CurrentBranchCherryPickCommitLog.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

func RenderGitEditCherryPickPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitEditCherryPickPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitEditCherryPickPopUpWidth, int(float64(m.Width)*0.8))
		popUpHeight := max(constant.PopUpGitEditCherryPickPopUpHeight, int(float64(m.Height)*0.5)-3)
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.EditCherryPickTitle)
		popUp.CherryPickedCommitLog.SetWidth(popUpWidth - 4)
		popUp.CherryPickedCommitLog.SetHeight(popUpHeight)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.CherryPickedCommitLog.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
