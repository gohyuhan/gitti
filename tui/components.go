package tui

import (
	"fmt"
	"gitti/i18n"
	"gitti/settings"

	"github.com/charmbracelet/lipgloss/v2"
)

// -----------------------------------------------------------------------------
// Gitti Main Page View
// -----------------------------------------------------------------------------
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
	gittiVersionPanel := panelBorderStyle.Width(m.HomeTabLeftPanelWidth).Height(1).Render(fmt.Sprintf("%s %s", settings.AppName, settings.AppVersion))
	localBranchesPanel := renderLocalBranchesPanel(m.HomeTabLeftPanelWidth, m.HomeTabLocalBranchesPanelHeight, m)
	changedFilesPanel := renderChangedFilesPanel(m.HomeTabLeftPanelWidth, m.HomeTabChangedFilesPanelHeight, m)
	fileDiffPanel := renderFileDiffPanel(m.HomeTabFileDiffPanelWidth, m.HomeTabFileDiffPanelHeight, m)
	bottomBar := renderKeyBindingPanel(m.Width, m)

	leftPanel := lipgloss.JoinVertical(lipgloss.Left, gittiVersionPanel, localBranchesPanel, changedFilesPanel)

	// Combine panels horizontally with explicit top alignment
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, fileDiffPanel)

	// Stack vertically with explicit left alignment
	mainView := lipgloss.JoinVertical(lipgloss.Left, content, bottomBar)

	if m.ShowPopUp {
		// --- SETUP THE CANVAS ---
		// Create a new canvas that will hold all our layers.
		canvas := lipgloss.NewCanvas()

		// Create the base layer from our main UI string.
		// It has no offset (X:0, Y:0) and is the bottom-most layer (Z:0).
		baseLayer := lipgloss.NewLayer(mainView)
		canvas.AddLayers(baseLayer)
		// Render the popup view into a string.
		popUpComponent := renderPopUpComponent(m)

		// Calculate the X and Y coordinates to center the popup.
		// We need the popup's dimensions for this.
		popUpWidth := lipgloss.Width(popUpComponent)
		popUpHeight := lipgloss.Height(popUpComponent)
		x := (m.Width - popUpWidth) / 2
		y := (m.Height - popUpHeight) / 2

		// Create a new layer for the popup.
		// Position it using the calculated X and Y.
		// Give it a higher Z-index to ensure it's drawn on top.
		popUpLayer := lipgloss.NewLayer(popUpComponent).X(x).Y(y).Z(1)

		// Add the popup layer to the canvas.
		canvas.AddLayers(popUpLayer)
		// Render the entire canvas with all its layers into the final string.
		return canvas.Render()

	}

	return mainView
}

// -----------------------------------------------------------------------------
//
//	Functions that help construct the components part
//
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
		Render(m.CurrentSelectedFileDiffViewport.View())
}

func renderKeyBindingPanel(width int, m *GittiModel) string {
	var keys []string
	if m.ShowPopUp {
		switch m.PopUpType {
		case CommitPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForCommitPopUp
		}
	} else {
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

// render the PopUp and the content within it will be a determine dynamically
func renderPopUpComponent(m *GittiModel) string {
	var popUp string

	switch m.PopUpType {
	case CommitPopUp:
		popUp = createCommitPopUp(m)
	}

	return popUp
}

func createCommitPopUp(m *GittiModel) string {
	popUpWidth := min(maxCommitPopUpWidth, int(float64(m.Width)*0.8))
	m.PopUpModel.(*CommitPopUpModel).MessageTextInput.SetWidth(popUpWidth - 4)
	m.PopUpModel.(*CommitPopUpModel).DescriptionTextAreaInput.SetWidth(popUpWidth - 4)

	// Rendered content
	title := titleStyle.Render(i18n.LANGUAGEMAPPING.CommitPopUpMessageTitle)
	inputView := m.PopUpModel.(*CommitPopUpModel).MessageTextInput.View()
	descLabel := titleStyle.Render(i18n.LANGUAGEMAPPING.CommitPopUpDescriptionTitle)
	descView := m.PopUpModel.(*CommitPopUpModel).DescriptionTextAreaInput.View()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		inputView,
		"", // 1-line padding
		descLabel,
		descView,
		"",
	)
	return popUpBorderStyle.Width(popUpWidth).Render(content)
}
