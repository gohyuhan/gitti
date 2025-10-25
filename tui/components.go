package tui

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
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
		Render(m.CurrentRepoBranchesInfo.View())
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
		Render(m.CurrentRepoModifiedFilesInfo.View())
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

func renderKeyBindingPanel(keys []string, width int) string {
	distributedWidth := int(width / len(keys))
	keyBindingLine := ""

	for _, key := range keys {
		keyWidth := max(0, (distributedWidth-lipgloss.Width(key))/2)
		keyLine := strings.Repeat(" ", keyWidth) + key + strings.Repeat(" ", keyWidth)
		keyBindingLine += keyLine
	}

	return bottomBarStyle.Width(width).Height(mainPageKeyBindingLayoutPanelHeight).Render(keyBindingLine)
}

func GittiMainPageView(m *GittiModel) string {
	if m.Width < minWidth || m.Height < minHeight {
		return "Terminal too small — resize to continue."
	}

	keys := []string{"[c] Commit", "[p] Push", "[f] Fetch", "[q] Quit"}

	// --- Components ---
	localBranchesPanel := renderLocalBranchesPanel(m.HomeTabLeftPanelWidth, m.HomeTabLocalBranchesPanelHeight, m)
	changedFilesPanel := renderChangedFilesPanel(m.HomeTabLeftPanelWidth, m.HomeTabChangedFilesPanelHeight, m)
	fileDiffPanel := renderFileDiffPanel(m.HomeTabFileDiffPanelWidth, m.HomeTabFileDiffPanelHeight, m)
	bottomBar := renderKeyBindingPanel(keys, m.Width)

	leftPanel := lipgloss.JoinVertical(lipgloss.Left, localBranchesPanel, changedFilesPanel)

	// Combine panels horizontally with explicit top alignment
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, fileDiffPanel)

	// Stack vertically with explicit left alignment
	mainView := lipgloss.JoinVertical(lipgloss.Left, content, bottomBar)

	return mainView
}
