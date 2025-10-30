package tui

import (
	"github.com/charmbracelet/lipgloss/v2"
)

const (
	minWidth  = 90
	minHeight = 28

	padding                             = 1
	mainPageKeyBindingLayoutPanelHeight = 1

	maxLeftPanelWidth   = 80
	maxCommitPopUpWidth = 100
)

var (
	// Base colors
	colorPrimary   = lipgloss.Color("#00BFFF")
	colorSecondary = lipgloss.Color("#AAAAAA")
	colorHighlight = lipgloss.Color("#FFD700")
	colorAccent    = lipgloss.Color("#98FB98")
	colorError     = lipgloss.Color("#FF6B6B")
	colorBasic     = lipgloss.Color("#FFFFFF")

	// list component style
	titleStyle = lipgloss.NewStyle().Foreground(colorHighlight).
			Underline(true).
			Bold(true)
	itemStyle         = lipgloss.NewStyle()
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	paginationStyle   = lipgloss.NewStyle()

	// Styles
	topBarStyle = lipgloss.NewStyle().
			Foreground(colorBasic).
			Background(lipgloss.Color("#1E1E1E")).
			Bold(false)
	topBarHighLightStyle = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Background(lipgloss.Color("#1E1E1E")).
				Bold(true)

	bottomBarStyle = lipgloss.NewStyle().
			Foreground(colorBasic)

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(colorHighlight).
				Underline(true).
				Bold(true)

	diffOldLineStyle = lipgloss.NewStyle().
				Foreground(colorError)

	diffNewLineStyle = lipgloss.NewStyle().
				Foreground(colorAccent)

	panelBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(0).
				Margin(0).
				BorderForeground(lipgloss.Color("#555555"))
	selectedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				Padding(0).
				Margin(0).
				BorderForeground(colorBasic)
	popUpBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				Padding(0).
				Margin(0).
				BorderForeground(colorBasic)
)
