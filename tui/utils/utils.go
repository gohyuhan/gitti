package utils

import (
	"golang.org/x/text/width"
)

// TruncateString trims string s to fit within given display width,
// accounting for wide CJK characters, and appends "…" if truncated.
func TruncateString(s string, maxWidth int) string {
	displayWidth := 0
	runes := []rune(s)
	var result []rune

	for _, r := range runes {
		prop := width.LookupRune(r)
		k := 1
		if prop.Kind() == width.EastAsianWide || prop.Kind() == width.EastAsianFullwidth {
			k = 2
		}

		if displayWidth+k > maxWidth {
			break
		}

		displayWidth += k
		result = append(result, r)
	}

	if len(result) < len(runes) {
		// Add ellipsis if we have room
		if len(result) >= 2 {
			result = append(result[:len(result)-1], '…')
		} else if len(result) == 1 {
			result = []rune{'…'}
		}
	}

	return string(result)
}
