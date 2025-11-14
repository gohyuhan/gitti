package style

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss/v2"
)

var (
	// Base colors
	ColorPrimary    = lipgloss.Color("#009ACD")
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
	GlobalKeyBindingPopUpStyle = NewStyle.
					Border(lipgloss.ThickBorder()).
					Padding(0).
					Margin(0).
					BorderForeground(ColorPrompt)
	GlobalKeyBindingTitleLineStyle = NewStyle.
					Foreground(ColorPrompt)
	GlobalKeyBindingKeyMappingLineStyle = NewStyle.
						Foreground(ColorKeyBinding)
	StagedFileStyle = NewStyle.
			Foreground(ColorAccent)
	UnstagedFileStyle = NewStyle.
				Foreground(ColorError)
)

func GradientLines(lines []string) []string {
	colored := make([]string, len(lines))

	// Tunable values
	startHue := 200.0 // degrees
	hueStep := 12.0   // per line
	sat := 0.70       // 0–1
	light := 0.65     // 0–1

	// inline HSL→RGB→HEX conversion
	hslToHex := func(h, s, l float64) string {
		c := (1 - math.Abs(2*l-1)) * s
		x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
		m := l - c/2

		var r, g, b float64
		switch {
		case h < 60:
			r, g, b = c, x, 0
		case h < 120:
			r, g, b = x, c, 0
		case h < 180:
			r, g, b = 0, c, x
		case h < 240:
			r, g, b = 0, x, c
		case h < 300:
			r, g, b = x, 0, c
		default:
			r, g, b = c, 0, x
		}

		R := int((r + m) * 255)
		G := int((g + m) * 255)
		B := int((b + m) * 255)

		return fmt.Sprintf("#%02x%02x%02x", R, G, B)
	}

	for i, line := range lines {
		h := math.Mod(startHue+float64(i)*hueStep, 360)
		hex := hslToHex(h, sat, light)

		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color(hex))

		colored[i] = style.Render(line)
	}
	return colored
}
