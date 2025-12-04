package tui

import (
	"fmt"

	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"

	"charm.land/lipgloss/v2"
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
	case constant.GlobalKeyBindingPopUp:
		popUp = renderGlobalKeyBindingPopUp(m)
	case constant.CommitPopUp:
		popUp = renderGitCommitPopUp(m)
	case constant.AmendCommitPopUp:
		popUp = renderGitAmendCommitPopUp(m)
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
	case constant.ChooseGitPullTypePopUp:
		popUp = renderChooseGitPullTypePopUp(m)
	case constant.GitPullOutputPopUp:
		popUp = renderGitPullOutputPopUp(m)
	case constant.GitStashMessagePopUp:
		popUp = renderGitStashMessagePopUp(m)
	case constant.GitDiscardTypeOptionPopUp:
		popUp = renderGitDiscardTypeOptionPopUp(m)
	case constant.GitDiscardConfirmPromptPopUp:
		popUp = renderGitDiscardConfirmPromptPopup(m)
	case constant.GitStashOperationOutputPopUp:
		popUp = renderGitStashOperationOutputPopUp(m)
	case constant.GitStashConfirmPromptPopUp:
		popUp = renderGitStashConfirmPromptPopUp(m)
	case constant.GitResolveConflictOptionPopUp:
		popUp = renderGitResolveConflictOptionPopUp(m)
	}
	return popUp
}

// ------------------------------------
//
//	For Global Key binding pop up
//
// ------------------------------------
func renderGlobalKeyBindingPopUp(m *GittiModel) string {
	popUp, ok := m.PopUpModel.(*GlobalKeyBindingPopUpModel)
	if ok {
		keyBindingLine := "\n"

		// this will usually only be run once for the entire gitti session
		if m.GlobalKeyBindingKeyMapLargestLen < 1 {
			maxLen := 0
			for _, line := range i18n.LANGUAGEMAPPING.GlobalKeyBinding {
				if l := len(line.KeyBindingLine); l > maxLen {
					maxLen = l
				}
			}
			m.GlobalKeyBindingKeyMapLargestLen = maxLen
		}
		for _, line := range i18n.LANGUAGEMAPPING.GlobalKeyBinding {
			switch line.LineType {
			case i18n.TITLE:
				keyBindingLine += " " + fmt.Sprintf("%*s", m.GlobalKeyBindingKeyMapLargestLen, line.KeyBindingLine) +
					"  " +
					style.GlobalKeyBindingTitleLineStyle.Render(line.TitleOrInfoLine) +
					"\n"
			case i18n.INFO:
				keyBindingLine += " " + style.GlobalKeyBindingKeyMappingLineStyle.Render(fmt.Sprintf("%*s", m.GlobalKeyBindingKeyMapLargestLen, line.KeyBindingLine)) +
					"  " +
					line.TitleOrInfoLine +
					"\n"
			case i18n.WARN:
				keyBindingLine += " " + style.GlobalKeyBindingKeyMappingLineStyle.Render(fmt.Sprintf("%s", line.KeyBindingLine)) +
					line.TitleOrInfoLine +
					"\n"
			}
		}
		height := min(constant.PopUpGlobalKeyBindingViewPortHeight, int(float64(m.Height)*0.8))
		width := min(constant.MaxGlobalKeyBindingPopUpWidth, int(float64(m.Width)*0.8)-4)
		popUp.GlobalKeyBindingViewport.SetWidth(width)
		popUp.GlobalKeyBindingViewport.SetHeight(height)
		popUp.GlobalKeyBindingViewport.SetContent(keyBindingLine)
		return style.GlobalKeyBindingPopUpStyle.Render(popUp.GlobalKeyBindingViewport.View())
	}
	return ""
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
					BorderForeground(style.ColorGreenSoft)
			}
			popUp.GitCommitOutputViewport.SetWidth(popUpWidth - 4)
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

// to update the commit output log for git commit
// this also take care of log by pre commit and post commit
func updatePopUpCommitOutputViewPort(m *GittiModel) {
	popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
	if ok {
		popUp.GitCommitOutputViewport.SetWidth(min(constant.MaxCommitPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		var gitCommitOutputLog string
		logs := m.GitOperations.GitCommit.GitCommitOutput()
		for _, line := range logs {
			logLine := style.NewStyle.Render(line)
			gitCommitOutputLog += logLine + "\n"
		}
		popUp.GitCommitOutputViewport.SetContent(gitCommitOutputLog)
		popUp.GitCommitOutputViewport.PageDown()
	}
}

// ------------------------------------
//
//	For Git Commit (Amend)
//
// ------------------------------------
func renderGitAmendCommitPopUp(m *GittiModel) string {
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
		if popUp.GitAmendCommitOutputViewport.GetContent() != "" {
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

// to update the amend commit output log for git amend commit
// this also take care of log by pre commit and post commit
func updatePopUpAmendCommitOutputViewPort(m *GittiModel) {
	popUp, ok := m.PopUpModel.(*GitAmendCommitPopUpModel)
	if ok {
		popUp.GitAmendCommitOutputViewport.SetWidth(min(constant.MaxAmendCommitPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		var gitCommitOutputLog string
		logs := m.GitOperations.GitCommit.GitCommitOutput()
		for _, line := range logs {
			logLine := style.NewStyle.Render(line)
			gitCommitOutputLog += logLine + "\n"
		}
		popUp.GitAmendCommitOutputViewport.SetContent(gitCommitOutputLog)
		popUp.GitAmendCommitOutputViewport.PageDown()
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
			remoteUrlTitle,
			remoteUrlTitleInputView,
		)
		if noInitialRemote {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				addRemotePrompt,
				remoteNameTitle,
				remoteNameInputView,
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
					BorderForeground(style.ColorGreenSoft)
			}
			popUp.AddRemoteOutputViewport.SetWidth(popUpWidth - 4)
			logViewPort := logViewPortStyle.Render(popUp.AddRemoteOutputViewport.View())
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				remoteNameTitle,
				remoteNameInputView,
				remoteUrlTitle,
				remoteUrlTitleInputView,
				logViewPort,
			)
			if noInitialRemote {
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					addRemotePrompt,
					remoteNameTitle,
					remoteNameInputView,
					remoteUrlTitle,
					remoteUrlTitleInputView,
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
		popUp.AddRemoteOutputViewport.PageDown()
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
		popUp.RemoteList.SetWidth(popUpWidth - 4)
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
//	For Choosing a push option
//
// ------------------------------------
func renderChoosePushTypePopUp(m *GittiModel) string {
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
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.GitRemotePushOutputViewport.SetWidth(popUpWidth - 4)
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

func updatePopUpGitRemotePushOutputViewport(m *GittiModel) {
	popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
	if ok {
		popUp.GitRemotePushOutputViewport.SetWidth(min(constant.MaxGitRemotePushPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		logs := m.GitOperations.GitCommit.GitRemotePushOutput()
		var GitPushLog string
		for _, line := range logs {
			logLine := style.NewStyle.Render(line)
			GitPushLog += logLine + "\n"
		}
		popUp.GitRemotePushOutputViewport.SetContent(GitPushLog)
		popUp.GitRemotePushOutputViewport.PageDown()
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
		popUp.NewBranchTypeOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
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
			popUp.NewBranchNameInput.View(),
		)
		modifiedBranchName, isValid := api.IsBranchNameValid(popUp.NewBranchNameInput.Value())
		if !isValid {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				popUp.NewBranchNameInput.View(),
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
		popUp.SwitchTypeOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.SwitchTypeOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// pop up to render the output of the switch branch operation
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
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.SwitchBranchOutputViewport.SetWidth(popUpWidth - 4)
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
			popUp.SwitchBranchOutputViewport.PageDown()
		}
	}
}

// ------------------------------------
//
//	For Git Pull
//
// ------------------------------------
// choose git pull option
func renderChooseGitPullTypePopUp(m *GittiModel) string {
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

// for git pull output
func renderGitPullOutputPopUp(m *GittiModel) string {
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

func updatePopUpGitPullOutputViewport(m *GittiModel) {
	popUp, ok := m.PopUpModel.(*GitPullOutputPopUpModel)
	if ok {
		popUp.GitPullOutputViewport.SetWidth(min(constant.MaxGitPullOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		logs := m.GitOperations.GitPull.GetGitPullOutput()
		var GitPullLog string
		for _, line := range logs {
			logLine := style.NewStyle.Render(line)
			GitPullLog += logLine + "\n"
		}
		popUp.GitPullOutputViewport.SetContent(GitPullLog)
		popUp.GitPullOutputViewport.PageDown()
	}
}

// ------------------------------------
//
//	For Git Stash to prompt for stash message
//
// ------------------------------------
func renderGitStashMessagePopUp(m *GittiModel) string {
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
//	For Discard file changes type list selection
//
// ------------------------------------
func renderGitDiscardTypeOptionPopUp(m *GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitDiscardTypeOptionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitDiscardTypeOptionPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitDiscardTypeOptionTitle)
		popUp.DiscardTypeOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.DiscardTypeOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	For Discard file changes confirmation prompt
//
// ------------------------------------
func renderGitDiscardConfirmPromptPopup(m *GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitDiscardConfirmPromptPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitDiscardConfirmPromptPopupWidth, int(float64(m.Width)*0.8))
		var content string
		switch popUp.DiscardType {
		case git.DISCARDWHOLE:
			content = style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardWholeConfirmation, popUp.FilePathName))
		case git.DISCARDUNSTAGE:
			content = style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardUnstageConfirmation, popUp.FilePathName))
		case git.DISCARDUNTRACKED:
			content = style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardUntrackedConfirmation, popUp.FilePathName))
		case git.DISCARDNEWLYADDEDORCOPIED:
			content = style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardNewlyAddedorCopyConfirmation, popUp.FilePathName))
		case git.DISCARDANDREVERTRENAME:
			content = style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardAndRevertRenameConfirmation, popUp.FilePathName))
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	For stash operation output
//
// ------------------------------------
func renderGitStashOperationOutputPopUp(m *GittiModel) string {
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
func renderGitStashConfirmPromptPopUp(m *GittiModel) string {
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

// ------------------------------------
//
//	For resolve conflict option list
//
// ------------------------------------
func renderGitResolveConflictOptionPopUp(m *GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitResolveConflictOptionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitResolveConflictOptionPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitResolveConflictOptionTitle)

		popUp.ResolveConflictOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.ResolveConflictOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
