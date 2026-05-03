package commitlog

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
		// we are dividing by 2 because the list is a single string that have 1 \n that break into 2 line
		popUpHeight := max(constant.PopUpGitCherryPickPopUpHeight, int(float64(m.Height)*0.8))
		title := style.TitleStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.CherryPickTitle, style.NewStyle.Foreground(style.ColorYellowWarm).Render(popUp.CurrentBranchName)))
		popUp.CurrentBranchCherryPickCommitLog.SetWidth(popUpWidth - 4)
		popUp.CurrentBranchCherryPickCommitLog.SetHeight(popUpHeight - 2)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.CurrentBranchCherryPickCommitLog.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Height(popUpHeight).Render(content)
	}
	return ""
}

func RenderGitEditCherryPickPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitEditCherryPickPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitEditCherryPickPopUpWidth, int(float64(m.Width)*0.8))
		// we are dividing by 3 because the list is a single string that have 2 \n that break into 3 line
		popUpHeight := max(constant.PopUpGitEditCherryPickPopUpHeight, int(float64(m.Height)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.EditCherryPickTitle)
		popUp.CherryPickedCommitLog.SetWidth(popUpWidth - 4)
		popUp.CherryPickedCommitLog.SetHeight(popUpHeight - 2)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.CherryPickedCommitLog.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Height(popUpHeight).Render(content)
	}
	return ""
}

func RenderGitCherryPickApplyConfirmPopUp(m *types.GittiModel) string {
	popUpWidth := min(constant.MaxGitCherryPickApplyConfirmPopUpWidth, int(float64(m.Width)*0.8))
	confirmMessage := fmt.Sprintf(i18n.LANGUAGEMAPPING.CherryPickApplyConfirmTitle, style.NewStyle.Foreground(style.ColorYellowWarm).Render(m.CheckOutBranch))
	title := style.TitleStyle.Render(confirmMessage)
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
	)
	return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
}

func RenderGitRevertParentOptionSelectionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitRevertParentOptionSelectionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitRevertParentOptionSelectionPopUpWidth, int(float64(m.Width)*0.8))
		popUpHeight := max(constant.PopUpGitRevertParentOptionSelectionHeight, int(float64(m.Height)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitRevertParentOptionSelectionTitle)
		popUp.GitRevertParentOption.SetWidth(popUpWidth - 4)
		popUp.GitRevertParentOption.SetHeight(popUpHeight - 2)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.GitRevertParentOption.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

func RenderGitRevertConfirmationPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitRevertConfirmationPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitRevertConfirmationPopUpWidth, int(float64(m.Width)*0.8))
		return style.PopUpBorderStyle.Width(popUpWidth).Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitRevertConfirmationTitle, style.NewStyle.Foreground(style.ColorYellowWarm).Render(popUp.CommitHash)))
	}
	return ""
}
