package tag

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	For Creating a tag
//
// ------------------------------------
func RenderCreateTagPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*CreateTagPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxCreateTagPopUpWidth, int(float64(m.Width)*0.8))
		popUp.TagNameInput.SetWidth(popUpWidth - 6)
		popUp.TagMessageTextAreaInput.SetWidth(popUpWidth - 6)

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			i18n.LANGUAGEMAPPING.CreateTagPopUpNameTitle,
			popUp.TagNameInput.View(),
			i18n.LANGUAGEMAPPING.CreateTagPopUpMessageTitle,
			popUp.TagMessageTextAreaInput.View(),
		)

		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	For Confirming Tag Creation
//
// ------------------------------------
func RenderCreateTagConfirmationPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*CreateTagConfirmationPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxCreateTagConfirmationPopUpWidth, int(float64(m.Width)*0.8))
		content := fmt.Sprintf(
			i18n.LANGUAGEMAPPING.CreateTagConfirmation,
			style.NewStyle.Foreground(style.ColorYellowWarm).Render(popUp.CommitHash),
			style.NewStyle.Foreground(style.ColorYellowWarm).Render(popUp.CommitMessage),
			style.NewStyle.Foreground(style.ColorPurpleVibrant).Render(popUp.TagName),
			style.NewStyle.Foreground(style.ColorBlueMuted).Render(popUp.TagMessage),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	For Choosing a Tag delete option type
//
// ------------------------------------
func RenderChooseDeleteTagOptionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChooseDeleteTagOptionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChooseDeleteTagOptionPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.ChooseDeleteTagOptionTitle)
		popUp.DeleteOptionList.SetWidth(popUpWidth - 4)
		popUp.DeleteOptionList.SetHeight(constant.PopUpChooseDeleteTagOptionHeight)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.DeleteOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	For Choosing a Remote for git delete remote tag if there is more than 1
//
// ------------------------------------
func RenderChooseRemoteForDeleteRemoteTagPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChooseRemoteForDeleteRemoteTagPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChooseRemoteForDeleteRemoteTagPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.ChooseRemoteTitle)
		popUp.RemoteList.SetWidth(popUpWidth - 4)
		popUp.RemoteList.SetHeight(constant.PopUpChooseRemoteForDeleteRemoteTagHeight)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.RemoteList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	For Delete Tag Output
//
// ------------------------------------
func RenderDeleteTagOutputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*DeleteTagOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxDeleteTagOutputPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.DeleteTagOutputPopUpTitle, popUp.TagName))
		logViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpDeleteTagOutputViewportHeight + 2)
		if popUp.HasError.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.DeleteTagOutputViewport.SetWidth(popUpWidth - 4)
		popUp.DeleteTagOutputViewport.SetYOffset(popUp.DeleteTagOutputViewport.YOffset())
		logViewPort := logViewPortStyle.Render(popUp.DeleteTagOutputViewport.View())

		var content string
		// Show spinner above viewport when processing
		if popUp.IsProcessing.Load() {
			processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.DeleteTagDeleting)
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
