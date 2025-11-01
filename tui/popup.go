package tui

import (
	"gitti/api/git"
	"gitti/i18n"

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
	case CommitPopUp:
		popUp = renderGitCommitPopUp(m)
	}

	return popUp
}

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
			logViewPort := panelBorderStyle.
				Width(popUpWidth - 2).
				Height(popUpGitCommitOutputViewPortHeight + 2).
				Render(popUp.GitCommitOutputViewport.View())
			if popUp.HasError {
				logViewPort = panelBorderStyle.
					BorderForeground(colorError).
					Width(popUpWidth - 2).
					Height(popUpGitCommitOutputViewPortHeight + 2).
					Render(popUp.GitCommitOutputViewport.View())
			}
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
			logLine := lipgloss.NewStyle().Foreground(colorBasic).Render(line)
			gitCommitOutputLog += logLine + "\n"
		}
		popUp.GitCommitOutputViewport.SetContent(gitCommitOutputLog)
		popUp.GitCommitOutputViewport.ViewDown()
	}
}
