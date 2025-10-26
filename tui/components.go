package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"gitti/i18n"
)

// -----------------------------------------------------------------------------
// Gitti Main Page View
// -----------------------------------------------------------------------------
// Render the Local Branches panel (top 25%)
func renderLocalBranchesPanel(width int, height int, m *GittiModel) string {
	borderStyle := panelBorderStyle
	if m.CurrentSelectedContainer == LocalBranchComponent {
		borderStyle = selectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(m.CurrentRepoBranchesInfoList.View())
}

// Render the Changed Files panel (bottom 75%)
func renderChangedFilesPanel(width int, height int, m *GittiModel) string {
	borderStyle := panelBorderStyle
	if m.CurrentSelectedContainer == ModifiedFilesComponent {
		borderStyle = selectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(m.CurrentRepoModifiedFilesInfoList.View())
}

func renderFileDiffPanel(width int, height int, m *GittiModel) string {
	borderStyle := panelBorderStyle
	if m.CurrentSelectedContainer == FileDiffComponent {
		borderStyle = selectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		UnsetMaxWidth().
		UnsetMaxHeight().
		Render(m.CurrentSelectedFileDiffViewport.View())
}

func renderKeyBindingPanel(width int, m *GittiModel) string {
	var keys []string
	switch m.CurrentSelectedContainer {
	case NoneSelected:
		keys = i18n.LANGUAGEMAPPING.KeyBindingNoneSelected
	case LocalBranchComponent:
		CurrentSelectedBranch := m.CurrentRepoBranchesInfoList.SelectedItem()
		if CurrentSelectedBranch == nil {
			keys = i18n.LANGUAGEMAPPING.KeyBindingLocalBranchComponentNone
		} else {
			isCurrentSelectedBranchCheckedOutBranch := CurrentSelectedBranch.(gitBranchItem).IsCheckedOut
			if isCurrentSelectedBranchCheckedOutBranch {
				keys = i18n.LANGUAGEMAPPING.KeyBindingLocalBranchComponentIsCheckOut
			} else {
				keys = i18n.LANGUAGEMAPPING.KeyBindingLocalBranchComponentDefault
			}
		}
	case ModifiedFilesComponent:
		CurrentSelectedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
		if CurrentSelectedFile == nil {
			keys = i18n.LANGUAGEMAPPING.KeyBindingModifiedFilesComponentNone
		} else {
			isFileStaged := CurrentSelectedFile.(gitModifiedFilesItem).SelectedForStage
			if isFileStaged {
				keys = i18n.LANGUAGEMAPPING.KeyBindingModifiedFilesComponentIsStaged
			} else {
				keys = i18n.LANGUAGEMAPPING.KeyBindingModifiedFilesComponentDefault
			}
		}
	case FileDiffComponent:
		keys = i18n.LANGUAGEMAPPING.KeyBindingFileDiffComponent
	}

	distributedWidth := width / len(keys)
	var keyBindingLine string

	for _, key := range keys {
		truncatedKey := TruncateString(key, distributedWidth)
		cell := lipgloss.NewStyle().
			Width(distributedWidth).
			Align(lipgloss.Center).
			Render(truncatedKey)
		keyBindingLine += cell
	}

	keyBindingPanelWidth := lipgloss.Width(keyBindingLine)
	keyBindingLine = TruncateString(keyBindingLine, keyBindingPanelWidth)

	return bottomBarStyle.Width(width).Height(mainPageKeyBindingLayoutPanelHeight).Render(keyBindingLine)
}

func GittiMainPageView(m *GittiModel) string {
	if m.Width < minWidth || m.Height < minHeight {
		title := lipgloss.NewStyle().
			Bold(true).
			Render(i18n.LANGUAGEMAPPING.TerminalSizeWarning)

		// Styles for the metric labels and values
		labelStyle := lipgloss.NewStyle()
		passStyle := lipgloss.NewStyle().Foreground(colorAccent)
		failStyle := lipgloss.NewStyle().Foreground(colorError)

		// Height
		heightStatus := passStyle.Render(fmt.Sprintf("%s: %v", i18n.LANGUAGEMAPPING.CurrentTerminalHeight, m.Height))
		if m.Height < minHeight {
			heightStatus = failStyle.Render(fmt.Sprintf("%s: %v", i18n.LANGUAGEMAPPING.CurrentTerminalWidth, m.Height))
		}

		// Width
		widthStatus := passStyle.Render(fmt.Sprintf("%s: %v", i18n.LANGUAGEMAPPING.CurrentTerminalWidth, m.Width))
		if m.Width < minWidth {
			widthStatus = failStyle.Render(fmt.Sprintf("%s: %v", i18n.LANGUAGEMAPPING.CurrentTerminalWidth, m.Width))
		}

		// Combine formatted text
		warningLine := lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			fmt.Sprintf(
				"\n%s %d\n%s %d\n%s\n%s",
				labelStyle.Render(fmt.Sprintf("%s: ", i18n.LANGUAGEMAPPING.MinimumTerminalHeight)), minHeight,
				labelStyle.Render(fmt.Sprintf("%s: ", i18n.LANGUAGEMAPPING.MinimumTerminalWidth)), minWidth,
				heightStatus,
				widthStatus,
			),
		)

		centered := lipgloss.NewStyle().
			Width(m.Width).
			Height(m.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Render(warningLine)

		return centered
	}

	// --- Components ---
	localBranchesPanel := renderLocalBranchesPanel(m.HomeTabLeftPanelWidth, m.HomeTabLocalBranchesPanelHeight, m)
	changedFilesPanel := renderChangedFilesPanel(m.HomeTabLeftPanelWidth, m.HomeTabChangedFilesPanelHeight, m)
	fileDiffPanel := renderFileDiffPanel(m.HomeTabFileDiffPanelWidth, m.HomeTabFileDiffPanelHeight, m)
	bottomBar := renderKeyBindingPanel(m.Width, m)

	leftPanel := lipgloss.JoinVertical(lipgloss.Left, localBranchesPanel, changedFilesPanel)

	// Combine panels horizontally with explicit top alignment
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, fileDiffPanel)

	// Stack vertically with explicit left alignment
	mainView := lipgloss.JoinVertical(lipgloss.Left, content, bottomBar)

	return mainView
}
