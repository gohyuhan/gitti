package commit

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
//	Render the git commit popup with message and description text inputs. Once a
//	commit has started, appends a scrollable output viewport below the inputs;
//	adds a spinner line while processing, and colors the viewport border red on
//	error or green on success.
//
// ------------------------------------
func RenderGitCommitPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxCommitPopUpWidth, int(float64(m.Width)*0.8))
		popUp.MessageTextInput.SetWidth(popUpWidth - 6)
		popUp.DescriptionTextAreaInput.SetWidth(popUpWidth - 6)

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
//	Render the git amend-commit popup with pre-filled message and description
//	inputs. Once an amend has started, appends a scrollable output viewport;
//	adds a spinner line while processing, and colors the viewport border red on
//	error or green on success.
//
// ------------------------------------
func RenderGitAmendCommitPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitAmendCommitPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxAmendCommitPopUpWidth, int(float64(m.Width)*0.8))
		popUp.MessageTextInput.SetWidth(popUpWidth - 6)
		popUp.DescriptionTextAreaInput.SetWidth(popUpWidth - 6)

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

// ------------------------------------
//
//	Render the reset-latest-commit type selection popup, showing a titled list
//	of soft/hard/mixed reset options for the user to choose from.
//
// ------------------------------------
func RenderGitResetLatestCommitTypeOptionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitResetLatestCommitTypeOptionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitResetLatestCommitTypeOptionPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitResetLatestCommitTypeOptionTitle)
		popUp.ResetLatestCommitTypeOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.ResetLatestCommitTypeOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the reset-latest-commit confirmation popup, displaying a localized
//	confirmation message that matches the chosen reset type (soft/hard/mixed).
//
// ------------------------------------
func RenderGitResetLatestCommitConfirmPromptPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitResetLatestCommitConfirmPromptPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitResetLatestCommitConfirmPromptPopUpWidth, int(float64(m.Width)*0.8))
		var content string
		switch popUp.GitResetLatestCommitType {
		case git.RESETSOFT:
			content = style.NewStyle.Render(i18n.LANGUAGEMAPPING.GitResetLatestCommitSoftConfirmation)
		case git.RESETHARD:
			content = style.NewStyle.Render(i18n.LANGUAGEMAPPING.GitResetLatestCommitHardConfirmation)
		case git.RESETMIXED:
			content = style.NewStyle.Render(i18n.LANGUAGEMAPPING.GitResetLatestCommitMixedConfirmation)
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the reset-to-selected-commit type selection popup, showing the target
//	commit hash in the title alongside a soft/hard/mixed option list.
//
// ------------------------------------
func RenderGitResetToSelectedCommitTypeOptionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitResetToSelectedCommitTypeOptionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitResetToSelectedCommitTypeOptionPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitResetToSelectedCommitTypeOptionTitle, popUp.SelectedCommitHash))
		popUp.ResetToSelectedCommitTypeOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.ResetToSelectedCommitTypeOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the reset-to-selected-commit confirmation popup, showing the
//	localized message for the chosen reset type with the target commit hash,
//	author, and message displayed below when available.
//
// ------------------------------------
func RenderGitResetToSelectedCommitConfirmPromptPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitResetToSelectedCommitConfirmPromptPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitResetToSelectedCommitConfirmPromptPopUpWidth, int(float64(m.Width)*0.8))
		var content string
		var title string
		switch popUp.GitResetToSelectedCommitType {
		case git.RESETSOFT:
			title = style.NewStyle.Render(i18n.LANGUAGEMAPPING.GitResetToSelectedCommitSoftConfirmation)
		case git.RESETHARD:
			title = style.NewStyle.Render(i18n.LANGUAGEMAPPING.GitResetToSelectedCommitHardConfirmation)
		case git.RESETMIXED:
			title = style.NewStyle.Render(i18n.LANGUAGEMAPPING.GitResetToSelectedCommitMixedConfirmation)
		}

		if popUp.CommitInfoAuthor == "" || popUp.CommitInfoMessage == "" {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				"",
				style.NewStyle.Foreground(style.ColorPurpleVibrant).Render(popUp.SelectedCommitHash),
			)
		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				"",
				style.NewStyle.Foreground(style.ColorPurpleVibrant).Render(popUp.SelectedCommitHash),
				style.NewStyle.Foreground(style.ColorYellowSoft).Render(popUp.CommitInfoAuthor),
				style.NewStyle.Foreground(style.ColorYellowWarm).Render(popUp.CommitInfoMessage),
			)
		}

		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
