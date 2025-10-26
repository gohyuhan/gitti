package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var KeyBindingHashmap = map[string]interface{}{
	NoneSelected: []string{"[b] branch component", "[f] files component", "[esc] quit"},
	LocalBranchComponent: map[string]interface{}{
		"IsCheckOut": []string{"[s] stash all file(s)", "[u] unstage all file(s)", "[esc] unselect component"},
		"Default":    []string{"[enter] switch branch", "[s] stash all file(s)", "[u] unstage all file(s)", "[esc] unselect component"},
		"None":       []string{"[esc] unselect component"},
	},
	ModifiedFilesComponent: map[string]interface{}{
		"IsStaged": []string{"[s] unstage this change", "[enter] view modified content", "[esc] unselect component"},
		"Default":  []string{"[s] stage this change", "[enter] view modified content", "[esc] unselect component"},
		"None":     []string{"[esc] unselect component"},
	},
	FileDiffComponent: []string{"[←/→] move left and right", "[↑/↓] move up and down", "[esc] back to file compoenent"},
}

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
		keys = KeyBindingHashmap[NoneSelected].([]string)
	case LocalBranchComponent:
		CurrentSelectedBranch := m.CurrentRepoBranchesInfoList.SelectedItem()
		if CurrentSelectedBranch == nil {
			keys = KeyBindingHashmap[LocalBranchComponent].(map[string]interface{})["None"].([]string)
		} else {
			isCurrentSelectedBranchCheckedOutBranch := CurrentSelectedBranch.(gitBranchItem).IsCheckedOut
			if isCurrentSelectedBranchCheckedOutBranch {
				keys = KeyBindingHashmap[LocalBranchComponent].(map[string]interface{})["IsCheckOut"].([]string)
			} else {
				keys = KeyBindingHashmap[LocalBranchComponent].(map[string]interface{})["Default"].([]string)
			}
		}
	case ModifiedFilesComponent:
		CurrentSelectedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
		if CurrentSelectedFile == nil {
			keys = KeyBindingHashmap[ModifiedFilesComponent].(map[string]interface{})["None"].([]string)
		} else {
			isFileStaged := CurrentSelectedFile.(gitModifiedFilesItem).SelectedForStage
			if isFileStaged {
				keys = KeyBindingHashmap[ModifiedFilesComponent].(map[string]interface{})["IsStaged"].([]string)
			} else {
				keys = KeyBindingHashmap[ModifiedFilesComponent].(map[string]interface{})["Default"].([]string)
			}
		}
	case FileDiffComponent:
		keys = KeyBindingHashmap[FileDiffComponent].([]string)
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
			Render("Terminal too small — resize to continue.")

		// Styles for the metric labels and values
		labelStyle := lipgloss.NewStyle()
		passStyle := lipgloss.NewStyle().Foreground(colorAccent)
		failStyle := lipgloss.NewStyle().Foreground(colorError)

		// Height
		heightStatus := passStyle.Render(fmt.Sprintf("Current height: %v", m.Height))
		if m.Height < minHeight {
			heightStatus = failStyle.Render(fmt.Sprintf("Current height: %v", m.Height))
		}

		// Width
		widthStatus := passStyle.Render(fmt.Sprintf("Current width: %v", m.Width))
		if m.Width < minWidth {
			widthStatus = failStyle.Render(fmt.Sprintf("Current width: %v", m.Width))
		}

		// Combine formatted text
		warningLine := lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			fmt.Sprintf(
				"\n%s %d\n%s %d\n%s\n%s",
				labelStyle.Render("Minimum required height:"), minHeight,
				labelStyle.Render("Minimum required width:"), minWidth,
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
