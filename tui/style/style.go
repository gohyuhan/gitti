package style

import (
	"github.com/charmbracelet/lipgloss/v2"
)

var (
	// Base colors
	ColorPrimary    = lipgloss.Color("#00BFFF")
	ColorSecondary  = lipgloss.Color("#AAAAAA")
	ColorHighlight  = lipgloss.Color("#FFD700")
	ColorAccent     = lipgloss.Color("#98FB98")
	ColorError      = lipgloss.Color("#FF6B6B")
	ColorBasic      = lipgloss.Color("#FFFFFF")
	ColorFade       = lipgloss.Color("#555555")
	ColorPrompt     = lipgloss.Color("#DB74ED")
	ColorTitle      = lipgloss.Color("#FF4500")
	ColorKeyBinding = lipgloss.Color("#AAF0F0")

	// lipgloss empty new style
	NewStyle = lipgloss.NewStyle()

	// list component style
	ItemStyle         = NewStyle
	SelectedItemStyle = NewStyle.Foreground(ColorPrimary)
	PaginationStyle   = NewStyle

	// Styles
	TitleStyle = NewStyle.Foreground(ColorTitle).
			Bold(true)
	PromptTitleStyle = NewStyle.Foreground(ColorPrompt).
				Bold(true)
	BottomKeyBindingStyle = NewStyle.
				Foreground(ColorKeyBinding)
	DiffOldLineStyle = NewStyle.
				Foreground(ColorError)
	DiffNewLineStyle = NewStyle.
				Foreground(ColorAccent)
	PanelBorderStyle = NewStyle.
				Border(lipgloss.RoundedBorder()).
				Padding(0).
				Margin(0).
				BorderForeground(ColorFade)
	SelectedBorderStyle = NewStyle.
				Border(lipgloss.DoubleBorder()).
				Padding(0).
				Margin(0).
				BorderForeground(ColorBasic)
	PopUpBorderStyle = NewStyle.
				Border(lipgloss.ThickBorder()).
				Padding(0).
				Margin(0).
				BorderForeground(ColorBasic)
	SpinnerStyle = NewStyle.
			Foreground(ColorPrimary)
	BranchInvalidWarningStyle = NewStyle.
					Foreground(ColorSecondary).
					Faint(true)
)
