package tui

import (
	"gitti/api/git"
	"gitti/i18n"

	"github.com/charmbracelet/lipgloss/v2"
)

// -----------------------------------------------------------------------------
//
//	Pop Up Type
//
// -----------------------------------------------------------------------------
const (
	NoPopUp              = "NoPopUp"
	CommitPopUp          = "CommitPopUp"          // IsTyping will be true
	AddRemotePromptPopUp = "AddRemotePromptPopUp" // IsTyping will be true
	ChoosePushTypePopUp  = "ChoosePushTypePopUp"  // IsTyping will be false
	ChooseRemotePopUp    = "ChooseRemotePopUp"    // IsTyping will be false
	GitRemotePushPopUp   = "GitRemotePushPopUp"   // IsTyping will be false
)

const AUTOCLOSEINTERVAL = 500

// -----------------------------------------------------------------------------
//
//	Functions that related to the rendering of pop up
//
// -----------------------------------------------------------------------------
// render the PopUp and the content within it will be a determine dynamically
func renderPopUpComponent(m *GittiModel) string {
	var popUp string

	switch m.PopUpType {
	case CommitPopUp:
		popUp = renderGitCommitPopUp(m)
	case AddRemotePromptPopUp:
		popUp = renderAddRemotePromptPopUp(m)
	case GitRemotePushPopUp:
		popUp = renderGitRemotePushPopUp(m)
	case ChooseRemotePopUp:
		popUp = renderChooseRemotePopUp(m)
	case ChoosePushTypePopUp:
		popUp = renderChoosePushTypePopUp(m)
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
		popUpWidth := min(maxCommitPopUpWidth, int(float64(m.Width)*0.8))
		popUp.MessageTextInput.SetWidth(popUpWidth - 4)
		popUp.DescriptionTextAreaInput.SetWidth(popUpWidth - 4)

		// Rendered content
		title := titleStyle.Render(i18n.LANGUAGEMAPPING.CommitPopUpMessageTitle)
		inputView := popUp.MessageTextInput.View()
		descLabel := titleStyle.Render(i18n.LANGUAGEMAPPING.CommitPopUpDescriptionTitle)
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
			logViewPortStyle := panelBorderStyle.
				Width(popUpWidth - 2).
				Height(popUpGitCommitOutputViewPortHeight + 2)
			if popUp.HasError {
				logViewPortStyle = panelBorderStyle.
					BorderForeground(colorError)
			} else if popUp.ProcessSuccess {
				logViewPortStyle = panelBorderStyle.
					BorderForeground(colorAccent)
			}

			logViewPort := logViewPortStyle.Render(popUp.GitCommitOutputViewport.View())

			// Show spinner above viewport when processing
			if popUp.IsProcessing {
				processingText := spinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.CommitPopUpProcessing)
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
		return popUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// to update the commit output log for git commit
// this also take care of log by pre commit and post commit
func updatePopUpCommitOutputViewPort(m *GittiModel) {
	popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
	if ok {
		popUp.GitCommitOutputViewport.SetWidth(min(maxCommitPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		var gitCommitOutputLog string
		logs := git.GITCOMMIT.GitCommitOutput
		for _, line := range logs {
			logLine := newStyle.Render(line)
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
		popUpWidth := min(maxAddRemotePromptPopUpWidth, int(float64(m.Width)*0.8))
		popUp.RemoteNameTextInput.SetWidth(popUpWidth - 4)
		popUp.RemoteUrlTextInput.SetWidth(popUpWidth - 4)

		noInitialRemote := popUp.NoInitialRemote

		// Rendered content
		addRemotePrompt := promptTitleStyle.Render(i18n.LANGUAGEMAPPING.AddRemotePopUpPrompt)
		remoteNameTitle := titleStyle.Render(i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteNameTitle)
		remoteNameInputView := popUp.RemoteNameTextInput.View()
		remoteUrlTitle := titleStyle.Render(i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteUrlTitle)
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
			logViewPortStyle := panelBorderStyle.
				Width(popUpWidth - 2).
				Height(popUpAddRemoteOutputViewPortHeight + 2)
			if popUp.HasError {
				logViewPortStyle = panelBorderStyle.
					BorderForeground(colorError)
			} else if popUp.ProcessSuccess {
				logViewPortStyle = panelBorderStyle.
					BorderForeground(colorAccent)
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
		return popUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

func updateAddRemoteOutputViewport(m *GittiModel, outputLog []string) {
	popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
	if ok {
		popUp.AddRemoteOutputViewport.SetWidth(min(maxAddRemotePromptPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		var addRemoteLog string
		for _, line := range outputLog {
			logLine := newStyle.Render(line)
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
		popUpWidth := min(maxChooseRemotePopUpWidth, int(float64(m.Width)*0.8))
		title := titleStyle.Render(i18n.LANGUAGEMAPPING.ChooseRemoteTitle)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			popUp.RemoteList.View(),
		)
		return popUpBorderStyle.Width(popUpWidth).Render(content)
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
		popUpWidth := min(maxChoosePushTypePopUpWidth, int(float64(m.Width)*0.8))
		title := titleStyle.Width(popUpWidth).Render(i18n.LANGUAGEMAPPING.GitRemotePushOptionTitle)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			popUp.PushOptionList.View(),
		)
		return popUpBorderStyle.Width(popUpWidth).Render(content)
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
		popUpWidth := min(maxGitRemotePushPopUpWidth, int(float64(m.Width)*0.8))
		title := titleStyle.Render(i18n.LANGUAGEMAPPING.GitRemotePushPopUpTitle)
		logViewPortStyle := panelBorderStyle.
			Width(popUpWidth - 2).
			Height(popUpGitCommitOutputViewPortHeight + 2)
		if popUp.HasError {
			logViewPortStyle = panelBorderStyle.
				BorderForeground(colorError)
		} else if popUp.ProcessSuccess {
			logViewPortStyle = panelBorderStyle.
				BorderForeground(colorAccent)
		}

		logViewPort := logViewPortStyle.Render(popUp.GitRemotePushOutputViewport.View())

		var content string
		// Show spinner above viewport when processing
		if popUp.IsProcessing {
			processingText := spinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.GitRemotePushPopUpProcessing)
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
		return popUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

func updateGitRemotePushOutputViewport(m *GittiModel) {
	popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
	if ok {
		popUp.GitRemotePushOutputViewport.SetWidth(min(maxGitRemotePushPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		logs := git.GITCOMMIT.GitRemotePushOutput
		var GitPushLog string
		for _, line := range logs {
			logLine := newStyle.Render(line)
			GitPushLog += logLine + "\n"
		}
		popUp.GitRemotePushOutputViewport.SetContent(GitPushLog)
		popUp.GitRemotePushOutputViewport.ViewDown()
	}
}
