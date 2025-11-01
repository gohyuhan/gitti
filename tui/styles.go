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

	popUpGitCommitOutputViewPortHeight = 8
)

var (
	// Base colors
	colorPrimary   = lipgloss.Color("#00BFFF")
	colorSecondary = lipgloss.Color("#AAAAAA")
	colorHighlight = lipgloss.Color("#FFD700")
	colorAccent    = lipgloss.Color("#98FB98")
	colorError     = lipgloss.Color("#FF6B6B")
	colorBasic     = lipgloss.Color("#FFFFFF")
	colorFade      = lipgloss.Color("#555555")

	// list component style
	titleStyle = lipgloss.NewStyle().Foreground(colorHighlight).
			Underline(true).
			Bold(true)
	itemStyle         = lipgloss.NewStyle()
	selectedItemStyle = lipgloss.NewStyle().Foreground(colorPrimary)
	paginationStyle   = lipgloss.NewStyle()

	// Styles
	bottomBarStyle = lipgloss.NewStyle().
			Foreground(colorBasic).Faint(true)
	diffOldLineStyle = lipgloss.NewStyle().
				Foreground(colorError)
	diffNewLineStyle = lipgloss.NewStyle().
				Foreground(colorAccent)
	panelBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(0).
				Margin(0).
				BorderForeground(colorFade)
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
