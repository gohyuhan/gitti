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
//	Render the create tag popup. Shows a bordered box containing the localized
//	name title, the tag name text input, the localized message title, and the
//	message textarea, stacked vertically.
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
//	Render the tag creation confirmation popup. Interpolates the styled commit
//	hash, commit message, tag name, and tag message into the localized
//	confirmation template and wraps it in a bordered box.
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
//	Render the delete tag scope selection popup. Shows a titled bordered list
//	where the user chooses between deleting the tag locally or remotely.
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
//	Render the remote selection popup for remote tag deletion. Shows a titled
//	bordered list of available remotes (name + URL) when there is more than one
//	configured remote.
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
//	Render the delete tag output popup. Colors the viewport border red on error
//	or green on success. When IsProcessing is true, shows a spinner line above
//	the viewport. The title includes the tag name.
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

// ------------------------------------
//
//	Render the push tag option popup. Shows a titled bordered list where the user
//	chooses how to push the tag: push tag, push all tags, force-push tag, or
//	force-push all tags.
//
// ------------------------------------
func RenderPushTagOptionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChoosePushTagOptionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChoosePushTagOptionPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.ChoosePushTagOptionTitle)
		popUp.PushOptionList.SetWidth(popUpWidth - 4)
		popUp.PushOptionList.SetHeight(constant.PopUpChoosePushTagOptionHeight)
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
//	Render the push tag output popup. Colors the viewport border red on error or
//	green on success. When IsProcessing is true, shows a spinner line above the
//	viewport.
//
// ------------------------------------
func RenderPushTagOutputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*PushTagOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxPushTagOutputPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.PushTagOutputPopUpTitle)
		logViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpPushTagOutputViewportHeight + 2)
		if popUp.HasError.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.PushTagOutputViewport.SetWidth(popUpWidth - 4)
		popUp.PushTagOutputViewport.SetYOffset(popUp.PushTagOutputViewport.YOffset())
		logViewPort := logViewPortStyle.Render(popUp.PushTagOutputViewport.View())

		var content string
		// Show spinner above viewport when processing
		if popUp.IsProcessing.Load() {
			processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.PushTagPushing)
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

// ------------------------------------
//
//	Render the fetch tag option popup. Shows a titled bordered list where the
//	user chooses how to fetch tags: fetch, fetch-overwrite, fetch-prune, or
//	mirror.
//
// ------------------------------------
func RenderFetchTagOptionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChooseFetchTagOptionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChooseFetchTagOptionPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.ChooseFetchTagOptionTitle)
		popUp.FetchOptionList.SetWidth(popUpWidth - 4)
		popUp.FetchOptionList.SetHeight(constant.PopUpChooseFetchTagOptionHeight)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.FetchOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the fetch tag output popup. Colors the viewport border red on error or
//	green on success. When IsProcessing is true, shows a spinner line above the
//	viewport.
//
// ------------------------------------
func RenderFetchTagOutputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*FetchTagOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxFetchTagOutputPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.FetchTagOutputPopUpTitle)
		logViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpFetchTagOutputViewportHeight + 2)
		if popUp.HasError.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.FetchTagOutputViewport.SetWidth(popUpWidth - 4)
		popUp.FetchTagOutputViewport.SetYOffset(popUp.FetchTagOutputViewport.YOffset())
		logViewPort := logViewPortStyle.Render(popUp.FetchTagOutputViewport.View())

		var content string
		// Show spinner above viewport when processing
		if popUp.IsProcessing.Load() {
			processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.FetchTagFetching)
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
