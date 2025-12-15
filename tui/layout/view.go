package layout

import (
	"fmt"

	"charm.land/lipgloss/v2"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/popup"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// -----------------------------------------------------------------------------
// Gitti Main Page View
// -----------------------------------------------------------------------------
func GittiMainPageView(m *types.GittiModel) string {
	if m.Width < constant.MinWidth || m.Height < constant.MinHeight {
		title := style.NewStyle.
			Bold(true).
			Render(i18n.LANGUAGEMAPPING.TerminalSizeWarning)

		// Styles for the metric labels and values
		labelStyle := style.NewStyle
		passStyle := style.NewStyle.Foreground(style.ColorGreenSoft)
		failStyle := style.NewStyle.Foreground(style.ColorError)

		// Height
		heightStatus := passStyle.Render(fmt.Sprintf("%s: %v", i18n.LANGUAGEMAPPING.CurrentTerminalHeight, m.Height))
		if m.Height < constant.MinHeight {
			heightStatus = failStyle.Render(fmt.Sprintf("%s: %v", i18n.LANGUAGEMAPPING.CurrentTerminalWidth, m.Height))
		}

		// Width
		widthStatus := passStyle.Render(fmt.Sprintf("%s: %v", i18n.LANGUAGEMAPPING.CurrentTerminalHeight, m.Width))
		if m.Width < constant.MinWidth {
			widthStatus = failStyle.Render(fmt.Sprintf("%s: %v", i18n.LANGUAGEMAPPING.CurrentTerminalWidth, m.Width))
		}

		// Combine formatted text
		warningLine := lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			fmt.Sprintf(
				"\n%s %d\n%s %d\n%s\n%s",
				labelStyle.Render(fmt.Sprintf("%s: ", i18n.LANGUAGEMAPPING.MinimumTerminalHeight)), constant.MinHeight,
				labelStyle.Render(fmt.Sprintf("%s: ", i18n.LANGUAGEMAPPING.MinimumTerminalWidth)), constant.MinWidth,
				heightStatus,
				widthStatus,
			),
		)

		centered := style.NewStyle.
			Width(m.Width).
			Height(m.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Render(warningLine)

		return centered
	}

	// --- Components ---
	GitStatusPanel := renderGitStatusComponentPanel(m)
	localBranchesPanel := renderLocalBranchesComponentPanel(m.WindowLeftPanelWidth, m.LocalBranchesComponentPanelHeight, m)
	modifiedFilesPanel := renderModifiedFilesComponentPanel(m.WindowLeftPanelWidth, m.ModifiedFilesComponentPanelHeight, m)
	commitLogPanel := renderCommitLogComponentPanel(m.WindowLeftPanelWidth, m.CommitLogComponentPanelHeight, m)
	stashFilesPanel := renderStashComponentPanel(m.WindowLeftPanelWidth, m.StashComponentPanelHeight, m)
	detailPanel := renderDetailComponentPanel(m.DetailComponentPanelWidth, m.DetailComponentPanelHeight, m)
	bottomBar := renderKeyBindingComponentPanel(m.Width, m)

	leftPanel := lipgloss.JoinVertical(lipgloss.Left, GitStatusPanel, localBranchesPanel, modifiedFilesPanel, commitLogPanel, stashFilesPanel)

	// Combine panels horizontally with explicit top alignment
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, detailPanel)

	// Stack vertically with explicit left alignment
	mainView := lipgloss.JoinVertical(lipgloss.Left, content, bottomBar)

	if m.ShowPopUp.Load() {
		// --- SETUP THE CANVAS ---
		// Create a new canvas that will hold all our layers.
		canvas := lipgloss.NewCanvas()

		// Create the base layer from our main UI string.
		// It has no offset (X:0, Y:0) and is the bottom-most layer (Z:0).
		baseLayer := lipgloss.NewLayer(mainView)
		canvas.AddLayers(baseLayer)
		// Render the popup view into a string.
		popUpComponent := popup.RenderPopUpComponent(m)

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
