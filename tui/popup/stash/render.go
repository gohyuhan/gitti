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
//	For Git Stash to prompt for stash message
//
// ------------------------------------
func RenderGitStashMessagePopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitStashMessagePopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitStashMessagePopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitStashMessageTitle)
		popUp.StashMessageInput.SetWidth(popUpWidth - 4)
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
//	For stash operation output
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
//	For stash operation confirmation prompt
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
