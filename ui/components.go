package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// -----------------------------------------------------------------------------
// Gitti Main Page View
// -----------------------------------------------------------------------------
func renderTitleBar(width int) string {

	distributedWidth := int(float64(width) * 0.33)

	home := "[1] Home"
	commitLogs := "[2] Commit Logs"
	about := "[3] About Gitti"

	homeWidth := max(0, (distributedWidth-lipgloss.Width(home))/2)
	commitLogsWidth := max(0, (distributedWidth-lipgloss.Width(commitLogs))/2)
	aboutWidth := max(0, (distributedWidth-lipgloss.Width(about))/2)

	titleLine := strings.Repeat(" ", homeWidth) + home + strings.Repeat(" ", homeWidth) + strings.Repeat(" ", commitLogsWidth) + commitLogs + strings.Repeat(" ", commitLogsWidth) + strings.Repeat(" ", aboutWidth) + about + strings.Repeat(" ", aboutWidth)
	return topBarStyle.Width(width).Height(mainPageLayoutTitlePanelHeight).Render(titleLine)
}

// Render the Local Branches panel (top 25%)
func renderLocalBranchesPanel(width int, height int, m GittiModel) string {
	items := []list.Item{
		item(fmt.Sprintf("* %s", m.CurrentCheckedOutBranch)),
	}

	for _, branch := range m.AllRepoBranches {
		if !branch.CurrentCheckout {
			items = append(items, item(branch.Name))
		}
	}

	m.CurrentRepoBranchesInfo = list.New(items, itemDelegate{}, width, height)
	m.CurrentRepoBranchesInfo.Title = "[b]  Branches:"
	m.CurrentRepoBranchesInfo.SetShowStatusBar(false)
	m.CurrentRepoBranchesInfo.SetFilteringEnabled(false)
	m.CurrentRepoBranchesInfo.SetShowHelp(false)
	m.CurrentRepoBranchesInfo.Styles.Title = titleStyle
	m.CurrentRepoBranchesInfo.Styles.PaginationStyle = paginationStyle

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

	// Compute panel widths
	leftPanelWidth := int(float64(m.Width) * mainPageLayoutLeftPanelWidthRatio)
	fileDiffPanelWidth := m.Width - leftPanelWidth - 4 // adjust for borders/padding

	coreContentHeight := m.Height - mainPageLayoutTitlePanelHeight - padding - mainPageKeyBindingLayoutPanelHeight - padding
	fileDiffPanelHeight := coreContentHeight
	localBranchesPanelHeight := int(float64(coreContentHeight)*mainPageLocalBranchesPanelHeightRatio) - padding
	changedFilesPanelHeight := int(float64(coreContentHeight) * mainPageChangedFilesHeightRatio)

	// --- Components ---
	topBar := renderTitleBar(m.Width)
	localBranchesPanel := renderLocalBranchesPanel(leftPanelWidth, localBranchesPanelHeight, m)
	changedFilesPanel := renderChangedFilesPanel(leftPanelWidth, changedFilesPanelHeight)
	fileDiffPanel := renderFileDiffPanel(fileDiffPanelWidth, fileDiffPanelHeight, m)
	bottomBar := renderKeyBindingPanel(keys, m.Width)

	leftPanel := lipgloss.JoinVertical(lipgloss.Left, localBranchesPanel, changedFilesPanel)

	// Combine panels horizontally with explicit top alignment
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, fileDiffPanel)

	// Stack vertically with explicit left alignment
	mainView := lipgloss.JoinVertical(lipgloss.Left, topBar, content, bottomBar)

	return mainView
}
