package rebase

import (
	"fmt"
	"unicode/utf8"

	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

func RenderGitRebaseBranchInputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitRebaseBranchInputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitRebaseBranchInputPopUpWidth, int(float64(m.Width)*0.8))
		popUp.BranchNameInput.SetWidth(popUpWidth - 6)
		var title string
		if utf8.RuneCountInString(popUp.Remote) > 0 {
			title = style.TitleStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitRebaseBranchInputPopUpTitleForRemoteBranch, popUp.Remote))
		} else {
			title = style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitRebaseBranchInputPopUpTitleForLocalBranch)
		}
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.BranchNameInput.View(),
		)
		modifiedBranchName, isValid := api.IsBranchNameValid(popUp.BranchNameInput.Value())
		if !isValid {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				popUp.BranchNameInput.View(),
				style.BranchInvalidWarningStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.NewBranchInvalidWarning, modifiedBranchName)),
			)
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	for git rebase output
//
// ------------------------------------
func RenderGitRebaseOutputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitRebaseOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitRebaseOutputPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitRebaseTitle)
		logViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpGitRebaseOutputViewportHeight + 2)
		if popUp.HasError.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.GitRebaseOutputViewport.SetWidth(popUpWidth - 4)
		popUp.GitRebaseOutputViewport.SetYOffset(popUp.GitRebaseOutputViewport.YOffset())
		logViewPort := logViewPortStyle.Render(popUp.GitRebaseOutputViewport.View())

		var content string
		// Show spinner above viewport when processing
		if popUp.IsProcessing.Load() {
			processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.GitRebaseProcessing)
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				"",
				processingText,
				logViewPort,
			)
		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				logViewPort,
			)
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
