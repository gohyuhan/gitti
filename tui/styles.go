package tui

import (
	"github.com/charmbracelet/lipgloss/v2"
)

const (
	minWidth  = 90
	minHeight = 28

	padding                             = 1
	mainPageKeyBindingLayoutPanelHeight = 1

	listItemOrTitleWidthPad = 5

	maxLeftPanelWidth            = 80
	maxCommitPopUpWidth          = 100
	maxAddRemotePromptPopUpWidth = 100
	maxGitRemotePushPopUpWidth   = 100

	popUpGitCommitOutputViewPortHeight     = 10
	popUpAddRemoteOutputViewPortHeight     = 2
	popUpGitRemotePushOutputViewportHeight = 10
)

var (
	// Base colors
	colorPrimary    = lipgloss.Color("#00BFFF")
	colorSecondary  = lipgloss.Color("#AAAAAA")
	colorHighlight  = lipgloss.Color("#FFD700")
	colorAccent     = lipgloss.Color("#98FB98")
	colorError      = lipgloss.Color("#FF6B6B")
	colorBasic      = lipgloss.Color("#FFFFFF")
	colorFade       = lipgloss.Color("#555555")
	colorPrompt     = lipgloss.Color("#DB74ED")
	colorTitle      = lipgloss.Color("#FF4500")
	colorKeyBinding = lipgloss.Color("#AAF0F0")

	// lipgloss empty new style
	newStyle = lipgloss.NewStyle()

	// list component style
	itemStyle         = newStyle
	selectedItemStyle = newStyle.Foreground(colorPrimary)
	paginationStyle   = newStyle

	// Styles
	titleStyle = newStyle.Foreground(colorTitle).
			Bold(true)
	promptTitleStyle = newStyle.Foreground(colorPrompt).
				Bold(true)
	bottomKeyBindingStyle = newStyle.
				Foreground(colorKeyBinding)
	diffOldLineStyle = newStyle.
				Foreground(colorError)
	diffNewLineStyle = newStyle.
				Foreground(colorAccent)
	panelBorderStyle = newStyle.
				Border(lipgloss.RoundedBorder()).
				Padding(0).
				Margin(0).
				BorderForeground(colorFade)
	selectedBorderStyle = newStyle.
				Border(lipgloss.DoubleBorder()).
				Padding(0).
				Margin(0).
				BorderForeground(colorBasic)
	popUpBorderStyle = newStyle.
				Border(lipgloss.ThickBorder()).
				Padding(0).
				Margin(0).
				BorderForeground(colorBasic)
	spinnerStyle = newStyle.
			Foreground(colorPrimary)
)
