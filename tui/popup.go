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
	popUpWidth := min(maxCommitPopUpWidth, int(float64(m.Width)*0.8))
	m.PopUpModel.(*GitCommitPopUpModel).MessageTextInput.SetWidth(popUpWidth - 4)
	m.PopUpModel.(*GitCommitPopUpModel).DescriptionTextAreaInput.SetWidth(popUpWidth - 4)

	// Rendered content
	title := titleStyle.Render(i18n.LANGUAGEMAPPING.CommitPopUpMessageTitle)
	inputView := m.PopUpModel.(*GitCommitPopUpModel).MessageTextInput.View()
	descLabel := titleStyle.Render(i18n.LANGUAGEMAPPING.CommitPopUpDescriptionTitle)
	descView := m.PopUpModel.(*GitCommitPopUpModel).DescriptionTextAreaInput.View()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		inputView,
		"", // 1-line padding
		descLabel,
		descView,
	)
	gitCommitOutputViewport := m.PopUpModel.(*GitCommitPopUpModel).GitCommitOutputViewport
	if gitCommitOutputViewport.GetContent() != "" {
		logViewPort := panelBorderStyle.Width(popUpWidth - 2).Height(popUpGitCommitOutputViewPortHeight + 2).Render(gitCommitOutputViewport.View())
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

// to update the commit output log for git commit
// this also take care of log by pre commit and post commit
func updatePopUpCommitOutputViewPort(m *GittiModel) {
	gitCommitOutputViewport := &m.PopUpModel.(*GitCommitPopUpModel).GitCommitOutputViewport
	gitCommitOutputViewport.SetWidth(min(maxCommitPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	var gitCommitOutputLog string
	for _, line := range git.GITCOMMIT.GitCommitOutput {
		logLine := lipgloss.NewStyle().Foreground(colorBasic).Render(line)
		gitCommitOutputLog += logLine + "\n"
	}
	gitCommitOutputViewport.SetContent(gitCommitOutputLog)
}
