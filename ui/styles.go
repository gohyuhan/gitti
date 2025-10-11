package ui

import "github.com/charmbracelet/lipgloss"

const (
	padding = 1

	mainPageLayoutLeftPanelWidthRatio  = 0.4
	mainPageLayoutRightPanelWidthRatio = 0.6

	mainPageLocalBranchesPanelHeightRatio  = 0.25
	mainPageChangedFilesHeightRatio = 0.75

	mainPageLayoutTitlePanelHeight      = 1
	mainPageKeyBindingLayoutPanelHeight = 1
)

var (
	minWidth  = 90
	minHeight = 30

	// Base colors
	colorPrimary   = lipgloss.Color("#00BFFF")
	colorSecondary = lipgloss.Color("#AAAAAA")
	colorHighlight = lipgloss.Color("#FFD700")
	colorAccent    = lipgloss.Color("#32CD32")
	colorError     = lipgloss.Color("#FF5555")

	// Styles
	topBarStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Background(lipgloss.Color("#1E1E1E")).
			Bold(true).
			PaddingLeft(1).
			Underline(true)

	bottomBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#303030")).
			PaddingLeft(1)

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(colorHighlight).
				Underline(true).
				Bold(true)

	listItemCheckedStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				SetString("[✓]")

	listItemUncheckedStyle = lipgloss.NewStyle().
				Foreground(colorError).
				SetString("[ ]")

	diffOldLineStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B"))

	diffNewLineStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#98FB98"))

	panelBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#555555")).
				Padding(0, 1)
)
