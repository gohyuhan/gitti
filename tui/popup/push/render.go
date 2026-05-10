package push

import (
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"

	"charm.land/lipgloss/v2"
)

// ------------------------------------
//
//	Render the push type selection popup. Shows a titled bordered list where the
//	user chooses between a normal push, safe force-push (--force-with-lease), or
//	dangerous force-push (--force).
//
// ------------------------------------
func RenderChoosePushTypePopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChoosePushTypePopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChoosePushTypePopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Width(popUpWidth).Render(i18n.LANGUAGEMAPPING.GitRemotePushOptionTitle)
		popUp.PushOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.PushOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the git push progress popup. Shows a scrollable output viewport with
//	a border that turns red on error or green on success. Displays a spinner
//	above the viewport while IsProcessing is true.
//
// ------------------------------------
func RenderGitRemotePushPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitRemotePushPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitRemotePushPopUpTitle)
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
		popUp.GitRemotePushOutputViewport.SetWidth(popUpWidth - 4)
		popUp.GitRemotePushOutputViewport.SetYOffset(popUp.GitRemotePushOutputViewport.YOffset())
		logViewPort := logViewPortStyle.Render(popUp.GitRemotePushOutputViewport.View())

		var content string
		// Show spinner above viewport when processing
		if popUp.IsProcessing.Load() {
			processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.GitRemotePushPopUpProcessing)
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
