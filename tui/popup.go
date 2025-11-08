package tui

import (
	"fmt"

	"gitti/api"
	"gitti/api/git"
	"gitti/i18n"
	"gitti/tui/constant"
	"gitti/tui/style"

	"github.com/charmbracelet/lipgloss/v2"
)

// -----------------------------------------------------------------------------
//
//	Functions that related to the rendering of pop up
//
// -----------------------------------------------------------------------------
// render the PopUp and the content within it will be a determine dynamically
func renderPopUpComponent(m *GittiModel) string {
	var popUp string

	switch m.PopUpType {
	case constant.CommitPopUp:
		popUp = renderGitCommitPopUp(m)
	case constant.AddRemotePromptPopUp:
		popUp = renderAddRemotePromptPopUp(m)
	case constant.GitRemotePushPopUp:
		popUp = renderGitRemotePushPopUp(m)
	case constant.ChooseRemotePopUp:
		popUp = renderChooseRemotePopUp(m)
	case constant.ChoosePushTypePopUp:
		popUp = renderChoosePushTypePopUp(m)
	case constant.ChooseNewBranchTypePopUp:
		popUp = renderChooseNewBranchTypePopUp(m)
	case constant.CreateNewBranchPopUp:
		popUp = renderCreateNewBranchPopUp(m)
	case constant.ChooseSwitchBranchTypePopUp:
		popUp = renderChooseSwitchBranchTypePopUp(m)
	case constant.SwitchBranchOutputPopUp:
		popUp = renderSwitchBranchOutputPopUp(m)
	}
	return popUp
}

// ------------------------------------
//
//	For Git Commit
//
// ------------------------------------
func renderGitCommitPopUp(m *GittiModel) string {
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
			"", // 1-line padding
			descLabel,
			descView,
		)
		if popUp.GitCommitOutputViewport.GetContent() != "" {
			logViewPortStyle := style.PanelBorderStyle.
				Width(popUpWidth - 2).
				Height(constant.PopUpGitCommitOutputViewPortHeight + 2)
			if popUp.HasError.Load() {
				logViewPortStyle = style.PanelBorderStyle.
					BorderForeground(style.ColorError)
			} else if popUp.ProcessSuccess.Load() {
				logViewPortStyle = style.PanelBorderStyle.
					BorderForeground(style.ColorAccent)
			}

			logViewPort := logViewPortStyle.Render(popUp.GitCommitOutputViewport.View())

			// Show spinner above viewport when processing
			if popUp.IsProcessing.Load() {
				processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.CommitPopUpProcessing)
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					title,
					inputView,
					"", // 1-line padding
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
					"", // 1-line padding
					descLabel,
					descView,
					"",
					logViewPort,
				)
			}
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// to update the commit output log for git commit
// this also take care of log by pre commit and post commit
func updatePopUpCommitOutputViewPort(m *GittiModel) {
	popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
	if ok {
		popUp.GitCommitOutputViewport.SetWidth(min(constant.MaxCommitPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		var gitCommitOutputLog string
		logs := m.GitState.GitCommit.GitCommitOutput()
		for _, line := range logs {
			logLine := style.NewStyle.Render(line)
			gitCommitOutputLog += logLine + "\n"
		}
		popUp.GitCommitOutputViewport.SetContent(gitCommitOutputLog)
		popUp.GitCommitOutputViewport.ViewDown()
	}
}

// ------------------------------------
//
//	For Adding Git Remote
//
// ------------------------------------
func renderAddRemotePromptPopUp(m *GittiModel) string {
	popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxAddRemotePromptPopUpWidth, int(float64(m.Width)*0.8))
		popUp.RemoteNameTextInput.SetWidth(popUpWidth - 4)
		popUp.RemoteUrlTextInput.SetWidth(popUpWidth - 4)

		noInitialRemote := popUp.NoInitialRemote

		// Rendered content
		addRemotePrompt := style.PromptTitleStyle.Render(i18n.LANGUAGEMAPPING.AddRemotePopUpPrompt)
		remoteNameTitle := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteNameTitle)
		remoteNameInputView := popUp.RemoteNameTextInput.View()
		remoteUrlTitle := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteUrlTitle)
		remoteUrlTitleInputView := popUp.RemoteUrlTextInput.View()

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			remoteNameTitle,
			remoteNameInputView,
			"", // 1-line padding
			remoteUrlTitle,
			remoteUrlTitleInputView,
		)
		if noInitialRemote {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				addRemotePrompt,
				"",
				remoteNameTitle,
				remoteNameInputView,
				"", // 1-line padding
				remoteUrlTitle,
				remoteUrlTitleInputView,
			)
		}
		if popUp.AddRemoteOutputViewport.GetContent() != "" {
			logViewPortStyle := style.PanelBorderStyle.
				Width(popUpWidth - 2).
				Height(constant.PopUpAddRemoteOutputViewPortHeight + 2)
			if popUp.HasError.Load() {
				logViewPortStyle = style.PanelBorderStyle.
					BorderForeground(style.ColorError)
			} else if popUp.ProcessSuccess.Load() {
				logViewPortStyle = style.PanelBorderStyle.
					BorderForeground(style.ColorAccent)
			}

			logViewPort := logViewPortStyle.Render(popUp.AddRemoteOutputViewport.View())
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				remoteNameTitle,
				remoteNameInputView,
				"", // 1-line padding
				remoteUrlTitle,
				remoteUrlTitleInputView,
				"",
				logViewPort,
			)
			if noInitialRemote {
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					addRemotePrompt,
					"",
					remoteNameTitle,
					remoteNameInputView,
					"", // 1-line padding
					remoteUrlTitle,
					remoteUrlTitleInputView,
					"",
					logViewPort,
				)
			}
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

func updateAddRemoteOutputViewport(m *GittiModel, outputLog []string) {
	popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
	if ok {
		popUp.AddRemoteOutputViewport.SetWidth(min(constant.MaxAddRemotePromptPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		var addRemoteLog string
		for _, line := range outputLog {
			logLine := style.NewStyle.Render(line)
			addRemoteLog += logLine + "\n"
		}
		popUp.AddRemoteOutputViewport.SetContent(addRemoteLog)
		popUp.AddRemoteOutputViewport.ViewDown()
	}
}

// ------------------------------------
//
//	For Choosing a Remote for git push if there is more than 1
//
// ------------------------------------
func renderChooseRemotePopUp(m *GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChooseRemotePopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChooseRemotePopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.ChooseRemoteTitle)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			popUp.RemoteList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	For Choosing a push option
//
// ------------------------------------
func renderChoosePushTypePopUp(m *GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChoosePushTypePopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChoosePushTypePopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Width(popUpWidth).Render(i18n.LANGUAGEMAPPING.GitRemotePushOptionTitle)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			popUp.PushOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	For Git Push
//
// ------------------------------------
func renderGitRemotePushPopUp(m *GittiModel) string {
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
				BorderForeground(style.ColorAccent)
		}

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
				"",
				logViewPort,
			)
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

func updateGitRemotePushOutputViewport(m *GittiModel) {
	popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
	if ok {
		popUp.GitRemotePushOutputViewport.SetWidth(min(constant.MaxGitRemotePushPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		logs := m.GitState.GitCommit.GitRemotePushOutput()
		var GitPushLog string
		for _, line := range logs {
			logLine := style.NewStyle.Render(line)
			GitPushLog += logLine + "\n"
		}
		popUp.GitRemotePushOutputViewport.SetContent(GitPushLog)
		popUp.GitRemotePushOutputViewport.ViewDown()
	}
}

// ------------------------------------
//
//	For Creating New Git branch
//
// ------------------------------------
// pop up that confirm the option for creating a new branch, just create or create and move everything to the new branch
func renderChooseNewBranchTypePopUp(m *GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChooseNewBranchTypeOptionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChooseNewBranchTypePopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.ChooseNewBranchTypeTitle)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			popUp.NewBranchTypeOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// to prompt user for new branch name and then proceed to trigger the creation of branch and optionally move changes
func renderCreateNewBranchPopUp(m *GittiModel) string {
	popUp, ok := m.PopUpModel.(*CreateNewBranchPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxCreateNewBranchPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.CreateNewBranchTitle)
		popUp.NewBranchNameInput.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			popUp.NewBranchNameInput.View(),
		)
		modifiedBranchName, isValid := api.IsBranchNameValid(popUp.NewBranchNameInput.Value())
		if !isValid {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				"",
				popUp.NewBranchNameInput.View(),
				"",
				style.BranchInvalidWarningStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.NewBranchInvalidWarning, modifiedBranchName)),
			)

		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	For Switching Git branch
//
// ------------------------------------
// pop up that confirm the option for switching a branch, just switch or switch to the branch while bringing all the changes
func renderChooseSwitchBranchTypePopUp(m *GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChooseSwitchBranchTypePopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChooseSwitchBranchTypePopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.ChooseSwitchBranchTypeTitle, popUp.BranchName))
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			popUp.SwitchTypeOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// pop up to redner the output of the switch branch operation
// because we allow switching with bring changes over, there is conflict possiblities there fore we need to show the output
// so that the user is aware of it
func renderSwitchBranchOutputPopUp(m *GittiModel) string {
	popUp, ok := m.PopUpModel.(*SwitchBranchOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxSwitchBranchOutputPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.SwitchBranchSwitchingToPopUpTitle, popUp.BranchName))
		logViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpSwitchBranchOutputViewPortHeight + 2)
		if popUp.HasError.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorAccent)
		}

		logViewPort := logViewPortStyle.Render(popUp.SwitchBranchOutputViewport.View())

		var content string
		// Show spinner above viewport when processing
		if popUp.IsProcessing.Load() {
			processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.SwitchBranchPopUpSwitchProcessing)
			if popUp.SwitchType == git.SWITCHBRANCHWITHCHANGES {
				processingText = style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.SwitchBranchPopUpSwitchWithChangesProcessing)
			}
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
				"",
				logViewPort,
			)
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

func updateSwitchBranchOutputViewPort(m *GittiModel, gitOpsOutput []string) {
	popUp, ok := m.PopUpModel.(*SwitchBranchOutputPopUpModel)
	if ok {
		popUp.SwitchBranchOutputViewport.SetWidth(min(constant.MaxSwitchBranchOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		var gitOpsOutputLogs string
		for _, line := range gitOpsOutput {
			logLine := style.NewStyle.Render(line)
			gitOpsOutputLogs += logLine + "\n"
			popUp.SwitchBranchOutputViewport.SetContent(gitOpsOutputLogs)
			popUp.SwitchBranchOutputViewport.ViewDown()
		}
	}
}
