package stash

import (
	"fmt"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"

	"charm.land/lipgloss/v2"
)

// ------------------------------------
//
//	Render the stash message prompt popup. Shows a titled border box containing
//	the text input where the user enters an optional stash message before the
//	stash operation executes.
//
// ------------------------------------
func RenderGitStashMessagePopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitStashMessagePopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitStashMessagePopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitStashMessageTitle)
		popUp.StashMessageInput.SetWidth(popUpWidth - 6)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.StashMessageInput.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the stash operation output popup. Colors the viewport border red on
//	error or green on success. When IsProcessing is true, shows a spinner line
//	above the viewport. The localized title and processing text are selected by
//	StashOperationType (stash-all, stash-file, apply, drop, or pop).
//
// ------------------------------------
func RenderGitStashOperationOutputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitStashOperationOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitStashOperationOutputPopUpWidth, int(float64(m.Width)*0.8))
		logViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpGitStashOperationOutputViewPortHeight + 2)
		if popUp.HasError.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.GitStashOperationOutputViewport.SetWidth(popUpWidth - 4)
		popUp.GitStashOperationOutputViewport.SetYOffset(popUp.GitStashOperationOutputViewport.YOffset())
		logViewPort := logViewPortStyle.Render(popUp.GitStashOperationOutputViewport.View())

		var title string
		var processingText string

		switch popUp.StashOperationType {
		case git.STASHALL:
			title = style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitStashAllTitle)
			processingText = style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.GitStashAllProcessing)
		case git.STASHFILE:
			title = style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitStashFileTitle)
			processingText = style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.GitStashFileProcessing)
		case git.APPLYSTASH:
			title = style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitStashApplyTitle)
			processingText = style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.GitStashApplyProcessing)
		case git.DROPSTASH:
			title = style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitStashDropTitle)
			processingText = style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.GitStashDropProcessing)
		case git.POPSTASH:
			title = style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitStashPopTitle)
			processingText = style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.GitStashPopProcessing)
		}

		var content string
		// Show spinner above viewport when processing
		if popUp.IsProcessing.Load() {
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
//	Render the stash confirmation prompt popup. Selects the localized confirmation
//	message by StashOperationType (stash-all, stash-file, apply, drop, or pop),
//	interpolating the styled file path, stash message, or stash ID as needed.
//
// ------------------------------------
func RenderGitStashConfirmPromptPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitStashConfirmPromptPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitStashConfirmPromptPopUpWidth, int(float64(m.Width)*0.8))
		var content string
		fpn := style.StashFilePathStyle.Render(popUp.FilePathName)
		msg := style.StashMessageStyle.Render(popUp.StashMessage)
		id := style.StashIdStyle.Render(popUp.StashId)
		switch popUp.StashOperationType {
		case git.STASHALL:
			content = style.NewStyle.Render(i18n.LANGUAGEMAPPING.GitStashAllConfirmation)
		case git.STASHFILE:
			content = style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitStashFileConfirmation, fpn))
		case git.APPLYSTASH:
			content = style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitApplyStashConfirmation, msg, id))
		case git.DROPSTASH:
			content = style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDropStashConfirmation, msg, id))
		case git.POPSTASH:
			content = style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitPopStashConfirmation, msg, id))
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
