package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// -----------------------------------------------------------------------------
// Gitti Main Page View
// -----------------------------------------------------------------------------
func renderMainPageTitleBar(width int) string {
	title := lipgloss.NewStyle().
		Foreground(colorHighlight).
		Bold(true).
		Render("GitDeskTUI")

	repoPath := lipgloss.NewStyle().
		Foreground(colorSecondary).
		Render("/user/project/repo")

	branch := lipgloss.NewStyle().
		Foreground(colorAccent).
		Underline(true).
		Render("main")

	status := lipgloss.NewStyle().
		Foreground(colorSecondary).
		Render("✔ Clean")

	line := fmt.Sprintf("%s | %s | Branch: %s | %s", title, repoPath, branch, status)
	return topBarStyle.Width(width).Height(mainPageLayoutTitlePanelHeight).Render(line)
}

// Render the Local Branches panel (top 25%)
func renderLocalBranchesPanel(width int, height int) string {
	content := sectionTitleStyle.Render("Local Branches:") + "\n" +
		"  * main (current)\n" +
		"  - feature/login-ui\n" +
		"  - bugfix/db-timeout\n"

	return panelBorderStyle.
		Width(width).
		Height(height).
		Render(content)
}

// Render the Changed Files panel (bottom 75%)
func renderChangedFilesPanel(width int, height int) string {
	content := sectionTitleStyle.Render("Changed Files:") + "\n" +
		fmt.Sprintf("  %s main.go\n", listItemCheckedStyle) +
		fmt.Sprintf("  %s ui/view.go\n", listItemUncheckedStyle) +
		fmt.Sprintf("  %s internal/git/commit.go\n", listItemCheckedStyle)

	return panelBorderStyle.
		Width(width).
		Height(height).
		Render(content)
}

func renderFileDiffPanel(width int, height int) string {
	diffTitle := sectionTitleStyle.Render("Diff Viewer:") + "\n\n"

	diffContent := diffOldLineStyle.Render("- func oldLine() {}\n") +
		diffNewLineStyle.Render("+ func newLine() {}\n") +
		lipgloss.NewStyle().Render("  fmt.Println(\"Hello World\")\n")

	return panelBorderStyle.Width(width).Height(height).Render(diffTitle + diffContent)
}

func renderBottomBar(width int) string {
	keys := "[q] Quit    [c] Commit    [p] Push    [f] Fetch"
	padding := max(0, (width-lipgloss.Width(keys))/2)
	return bottomBarStyle.Width(width).Height(mainPageKeyBindingLayoutPanelHeight).Render(strings.Repeat(" ", padding) + keys)
}

func GittiMainPageView(m GittiModel) string {
	if m.width < minWidth || m.height < minHeight {
		return "Terminal too small — resize to continue."
	}

	// Compute panel widths
	leftPanelWidth := int(float64(m.width) * mainPageLayoutLeftPanelWidthRatio)
	fileDiffPanelWidth := m.width - leftPanelWidth - 4 // adjust for borders/padding

	coreContentHeight := m.height - mainPageLayoutTitlePanelHeight - padding - mainPageKeyBindingLayoutPanelHeight - padding
	fileDiffPanelHeight := coreContentHeight
	localBranchesPanelHeight := int(float64(coreContentHeight)*mainPageLocalBranchesPanelHeightRatio) - padding
	changedFilesPanelHeight := int(float64(coreContentHeight)*mainPageChangedFilesHeightRatio)

	// --- Components ---
	topBar := renderMainPageTitleBar(m.width)
	localBranchesPanel := renderLocalBranchesPanel(leftPanelWidth, localBranchesPanelHeight)
	changedFilesPanel := renderChangedFilesPanel(leftPanelWidth, changedFilesPanelHeight)
	fileDiffPanel := renderFileDiffPanel(fileDiffPanelWidth, fileDiffPanelHeight)
	bottomBar := renderBottomBar(m.width)

	leftPanel:= lipgloss.JoinVertical(lipgloss.Left, localBranchesPanel, changedFilesPanel)

	// Combine panels horizontally with explicit top alignment
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, fileDiffPanel)

	// Stack vertically with explicit left alignment
	mainView := lipgloss.JoinVertical(lipgloss.Left, topBar, content, bottomBar)

	return mainView
}
