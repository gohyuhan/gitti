package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	minWidth  = 90
	minHeight = 28

	padding                             = 1
	mainPageKeyBindingLayoutPanelHeight = 1
)

var (
	// Base colors
	colorPrimary   = lipgloss.Color("#00BFFF")
	colorSecondary = lipgloss.Color("#AAAAAA")
	colorHighlight = lipgloss.Color("#FFD700")
	colorAccent    = lipgloss.Color("#32CD32")
	colorError     = lipgloss.Color("#FF5555")
	colorBasic     = lipgloss.Color("#FFFFFF")

	// list component style
	titleStyle = lipgloss.NewStyle().Foreground(colorHighlight).
			Underline(true).
			Bold(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(2)

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
				Foreground(lipgloss.Color("#FF6B6B"))

	diffNewLineStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#98FB98"))

	panelBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#555555"))
	selectedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(colorBasic)
)
