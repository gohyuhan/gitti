package commit

import (
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"

	"charm.land/lipgloss/v2"
)

// ------------------------------------
//
//	For Git Commit
//
// ------------------------------------
func RenderGitCommitPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxCommitPopUpWidth, int(float64(m.Width)*0.8))
		popUp.MessageTextInput.SetWidth(popUpWidth - 4)
		popUp.DescriptionTextAreaInput.SetWidth(popUpWidth - 4)

		// Rendered content
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.CommitPopUpMessageTitle)
		inputView := popUp.MessageTextInput.View()
		descLabel := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.CommitPopUpDescriptionTitle)
		descView := popUp.DescriptionTextAreaInput.View()

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			inputView,
			descLabel,
			descView,
		)
		if popUp.InitialCommitStarted.Load() {
			logViewPortStyle := style.PanelBorderStyle.
				Width(popUpWidth - 2).
				Height(constant.PopUpGitCommitOutputViewPortHeight + 2)
			if popUp.HasError.Load() {
				logViewPortStyle = style.PanelBorderStyle.
					BorderForeground(style.ColorError)
			} else if popUp.ProcessSuccess.Load() {
				logViewPortStyle = style.PanelBorderStyle.
					BorderForeground(style.ColorGreenSoft)
			}
			popUp.GitCommitOutputViewport.SetWidth(popUpWidth - 4)
			popUp.GitCommitOutputViewport.SetYOffset(popUp.GitCommitOutputViewport.YOffset())
			logViewPort := logViewPortStyle.Render(popUp.GitCommitOutputViewport.View())

			// Show spinner above viewport when processing
			if popUp.IsProcessing.Load() {
				processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.CommitPopUpProcessing)
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					title,
					inputView,
					descLabel,
					descView,
					"",
					processingText,
					logViewPort,
				)
			} else {
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					title,
					inputView,
					descLabel,
					descView,
					logViewPort,
				)
			}
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	For Git Commit (Amend)
//
// ------------------------------------
func RenderGitAmendCommitPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitAmendCommitPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxAmendCommitPopUpWidth, int(float64(m.Width)*0.8))
		popUp.MessageTextInput.SetWidth(popUpWidth - 4)
		popUp.DescriptionTextAreaInput.SetWidth(popUpWidth - 4)

		// Rendered content
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.CommitPopUpMessageTitleAmendVersion)
		inputView := popUp.MessageTextInput.View()
		descLabel := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.CommitPopUpDescriptionTitleAmendVersion)
		descView := popUp.DescriptionTextAreaInput.View()

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			inputView,
			descLabel,
			descView,
		)
		if popUp.InitialCommitStarted.Load() {
			logViewPortStyle := style.PanelBorderStyle.
				Width(popUpWidth - 2).
				Height(constant.PopUpGitCommitOutputViewPortHeight + 2)
			if popUp.HasError.Load() {
				logViewPortStyle = style.PanelBorderStyle.
					BorderForeground(style.ColorError)
			} else if popUp.ProcessSuccess.Load() {
				logViewPortStyle = style.PanelBorderStyle.
					BorderForeground(style.ColorGreenSoft)
			}
			popUp.GitAmendCommitOutputViewport.SetWidth(popUpWidth - 4)
			popUp.GitAmendCommitOutputViewport.SetYOffset(popUp.GitAmendCommitOutputViewport.YOffset())
			logViewPort := logViewPortStyle.Render(popUp.GitAmendCommitOutputViewport.View())

			// Show spinner above viewport when processing
			if popUp.IsProcessing.Load() {
				processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.CommitPopUpProcessing)
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					title,
					inputView,
					descLabel,
					descView,
					"",
					processingText,
					logViewPort,
				)
			} else {
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					title,
					inputView,
					descLabel,
					descView,
					logViewPort,
				)
			}
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
