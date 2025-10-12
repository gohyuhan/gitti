package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// -----------------------------------------------------------------------------
// Gitti Main Page View
// -----------------------------------------------------------------------------
func renderTitleBar(width int, currentTab string) string {
	distributedWidth := int(float64(width) * 0.33)

	home := "[1] Home"
	commitLogs := "[2] Commit Logs"
	about := "[3] About Gitti"

	// Define styles for active and inactive tabs
	activeTabStyle := topBarHighLightStyle
	inactiveTabStyle := topBarStyle
	// Apply styles based on currentTab
	var homeStyled, commitLogsStyled, aboutStyled string
	if currentTab == homeTab {
		homeStyled = activeTabStyle.Render(home)
		commitLogsStyled = inactiveTabStyle.Render(commitLogs)
		aboutStyled = inactiveTabStyle.Render(about)
	} else if currentTab == commitLogsTab {
		homeStyled = inactiveTabStyle.Render(home)
		commitLogsStyled = activeTabStyle.Render(commitLogs)
		aboutStyled = inactiveTabStyle.Render(about)
	} else if currentTab == aboutGittiTab {
		homeStyled = inactiveTabStyle.Render(home)
		commitLogsStyled = inactiveTabStyle.Render(commitLogs)
		aboutStyled = activeTabStyle.Render(about)
	} else {
		// Fallback: no highlight if currentTab is invalid
		homeStyled = inactiveTabStyle.Render(home)
		commitLogsStyled = inactiveTabStyle.Render(commitLogs)
		aboutStyled = inactiveTabStyle.Render(about)
	}

	// Calculate spacing
	homeWidth := max(0, (distributedWidth-lipgloss.Width(homeStyled))/2)
	commitLogsWidth := max(0, (distributedWidth-lipgloss.Width(commitLogsStyled))/2)
	aboutWidth := max(0, (distributedWidth-lipgloss.Width(aboutStyled))/2)

	// Combine styled tabs with spacing
	titleLine := strings.Repeat(" ", homeWidth) + homeStyled +
		strings.Repeat(" ", homeWidth) + strings.Repeat(" ", commitLogsWidth) + commitLogsStyled +
		strings.Repeat(" ", commitLogsWidth) + strings.Repeat(" ", aboutWidth) + aboutStyled +
		strings.Repeat(" ", aboutWidth)

	return topBarStyle.Width(width).Height(mainPageLayoutTitlePanelHeight).Render(titleLine)
}

// Render the Local Branches panel (top 25%)
func renderLocalBranchesPanel(width int, height int, m GittiModel) string {
	return panelBorderStyle.
		Width(width).
		Height(height).
		Render(m.CurrentRepoBranchesInfo.View())
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

func renderFileDiffPanel(width int, height int, m GittiModel) string {
	diffTitle := sectionTitleStyle.Render("Diff Viewer:") + "\n\n"

	diffContent := diffOldLineStyle.Render("- func oldLine() {}\n") +
		diffNewLineStyle.Render("+ func newLine() {}\n") +
		lipgloss.NewStyle().Render(fmt.Sprintf("all: %v, current: %v", m.AllRepoBranches, m.CurrentCheckedOutBranch))

	return panelBorderStyle.Width(width).Height(height).Render(diffTitle + diffContent)
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

func GittiMainPageView(m GittiModel) string {
	if m.Width < minWidth || m.Height < minHeight {
		return "Terminal too small — resize to continue."
	}

	keys := []string{"[c] Commit", "[p] Push", "[f] Fetch", "[q] Quit"}

	// --- Components ---
	topBar := renderTitleBar(m.Width, m.CurrentTab)
	localBranchesPanel := renderLocalBranchesPanel(m.HomeTabLeftPanelWidth, m.HomeTabLocalBranchesPanelHeight, m)
	changedFilesPanel := renderChangedFilesPanel(m.HomeTabLeftPanelWidth, m.HomeTabChangedFilesPanelHeight)
	fileDiffPanel := renderFileDiffPanel(m.HomeTabFileDiffPanelWidth, m.HomeTabFileDiffPanelHeight, m)
	bottomBar := renderKeyBindingPanel(keys, m.Width)

	leftPanel := lipgloss.JoinVertical(lipgloss.Left, localBranchesPanel, changedFilesPanel)

	// Combine panels horizontally with explicit top alignment
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, fileDiffPanel)

	// Stack vertically with explicit left alignment
	mainView := lipgloss.JoinVertical(lipgloss.Left, topBar, content, bottomBar)

	return mainView
}
