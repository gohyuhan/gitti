package style

import (
	"fmt"
	"math"

	"charm.land/lipgloss/v2"
)

var (
	// Base colors - Stunning gradient-inspired theme
	ColorBlueSoft      = lipgloss.Color("#82AAFF") // Soft periwinkle blue (from your reference)
	ColorBlueMuted     = lipgloss.Color("#A5B7E8") // Muted lavender-blue for secondary text
	ColorYellowWarm    = lipgloss.Color("#F0D278") // Warm golden highlight (from your reference)
	ColorYellowSoft    = lipgloss.Color("#F5E6A3") // Lighter golden yellow for emphasis
	ColorGreenSoft     = lipgloss.Color("#98FB98") // Keep original - DO NOT MODIFY
	ColorError         = lipgloss.Color("#FF6B6B") // Keep original - DO NOT MODIFY
	ColorBlueVeryLight = lipgloss.Color("#E8F0FF") // Soft white with blue tint for readability
	ColorBlueGrayMuted = lipgloss.Color("#6B7A9E") // Muted blue-gray for subtle elements
	ColorPurpleSoft    = lipgloss.Color("#B496FF") // Beautiful lavender purple (from your reference)
	ColorPurpleVibrant = lipgloss.Color("#9F7AEA") // Rich purple for titles
	ColorCyanSoft      = lipgloss.Color("#7DD3FC") // Sky blue for key bindings

	// lipgloss empty new style
	NewStyle = lipgloss.NewStyle()

	// list component style
	ItemStyle         = NewStyle
	SelectedItemStyle = NewStyle.Foreground(ColorBlueSoft)
	PaginationStyle   = NewStyle

	// Styles
	TitleStyle = NewStyle.Foreground(ColorPurpleVibrant).
			Bold(true)
	PromptTitleStyle = NewStyle.Foreground(ColorPurpleSoft).
				Bold(true)
	BottomKeyBindingStyle = NewStyle.
				Foreground(ColorCyanSoft)
	PanelBorderStyle = NewStyle.
				Border(lipgloss.RoundedBorder()).
				Padding(0).
				Margin(0).
				BorderForeground(ColorBlueGrayMuted)
	SelectedBorderStyle = NewStyle.
				Border(lipgloss.DoubleBorder()).
				Padding(0).
				Margin(0).
				BorderForeground(ColorBlueVeryLight)
	PopUpBorderStyle = NewStyle.
				Border(lipgloss.ThickBorder()).
				Padding(0).
				Margin(0).
				BorderForeground(ColorBlueVeryLight)
	SpinnerStyle = NewStyle.
			Foreground(ColorBlueSoft)
	BranchInvalidWarningStyle = NewStyle.
					Foreground(ColorBlueMuted).
					Faint(true)

	GlobalKeyBindingPopUpStyle = NewStyle.
					Border(lipgloss.ThickBorder()).
					Padding(0).
					Margin(0).
					BorderForeground(ColorPurpleSoft)
	GlobalKeyBindingTitleLineStyle = NewStyle.
					Foreground(ColorPurpleSoft)
	GlobalKeyBindingKeyMappingLineStyle = NewStyle.
						Foreground(ColorCyanSoft)

	DiffOldLineStyle = NewStyle.
				Foreground(ColorError)
	DiffNewLineStyle = NewStyle.
				Foreground(ColorGreenSoft)

	StagedFileStyle = NewStyle.
			Foreground(ColorGreenSoft)
	UnstagedFileStyle = NewStyle.
				Foreground(ColorError)

	LocalStatusStyle = NewStyle.
				Foreground(ColorGreenSoft)
	RemoteStatusStyle = NewStyle.
				Foreground(ColorError)

	StashIdStyle       = NewStyle.Foreground(ColorYellowWarm)
	StashMessageStyle  = NewStyle.Foreground(ColorYellowSoft)
	StashFilePathStyle = NewStyle.Foreground(ColorCyanSoft)

	ErrorStyle = NewStyle.
			Foreground(ColorError)
)

func GradientLines(lines []string) []string {
	colored := make([]string, len(lines))

	// // Tunable values
	// startHue := 200.0 // degrees
	// hueStep := 12.0   // per line
	// sat := 0.70       // 0–1
	// light := 0.65     // 0–1

	// Enhanced gradient values for stunning visual effect
	startHue := 220.0 // degrees - start at beautiful blue-purple
	hueStep := 8.0    // per line - smoother transitions
	sat := 0.75       // 0–1 - increased saturation for vibrancy
	light := 0.68     // 0–1 - optimized brightness for both light/dark terminals

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
