package pull

import (
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"

	"charm.land/lipgloss/v2"
)

// ------------------------------------
//
//	Render the pull-type selection popup, showing a list of pull strategy options
//	(merge, rebase, fast-forward only) for the user to choose before pulling.
//
// ------------------------------------
func RenderChooseGitPullTypePopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChooseGitPullTypePopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChooseGitPullTypePopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.ChoosePullOptionPrompt)
		popUp.PullTypeOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.PullTypeOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the git pull output popup, showing a scrollable viewport of pull
//	command output. Displays a spinner while the pull is in progress, and
//	colors the viewport border red on error or green on success.
//
// ------------------------------------
func RenderGitPullOutputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitPullOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitPullOutputPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitPullTitle)
		logViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpGitPullOutputViewportHeight + 2)
		if popUp.HasError.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.GitPullOutputViewport.SetWidth(popUpWidth - 4)
		popUp.GitPullOutputViewport.SetYOffset(popUp.GitPullOutputViewport.YOffset())
		logViewPort := logViewPortStyle.Render(popUp.GitPullOutputViewport.View())

		var content string
		// Show spinner above viewport when processing
		if popUp.IsProcessing.Load() {
			processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.GitPullProcessing)
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
