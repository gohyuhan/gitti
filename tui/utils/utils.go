package utils

import (
	"fmt"
	"os/exec"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
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

func ListCounterHelper(m *types.GittiModel, list *list.Model) func() []key.Binding {
	return func() []key.Binding {
		currentIndex := list.Index() + 1
		totalCount := len(list.Items())
		countStr := fmt.Sprintf("%d/%d", currentIndex, totalCount)
		countStr = TruncateString(countStr, m.WindowLeftPanelWidth-constant.ListItemOrTitleWidthPad-2)
		if totalCount == 0 {
			countStr = "0/0"
		}
		return []key.Binding{
			key.NewBinding(
				key.WithKeys(countStr),
				key.WithHelp(countStr, ""),
			),
		}
	}
}

func PopUpListCounterHelper(m *types.GittiModel, list *list.Model, maxWidth int) func() []key.Binding {
	return func() []key.Binding {
		currentIndex := list.Index() + 1
		totalCount := len(list.Items())
		countStr := fmt.Sprintf("%d/%d", currentIndex, totalCount)
		width := (min(maxWidth, int(float64(m.Width)*0.8)) - 4)
		countStr = TruncateString(countStr, width-constant.ListItemOrTitleWidthPad-2)
		if totalCount == 0 {
			countStr = "0/0"
		}
		return []key.Binding{
			key.NewBinding(
				key.WithKeys(countStr),
				key.WithHelp(countStr, ""),
			),
		}
	}
}

func ReturnEditorLaunchCommand(fileName string, userSetEditor string) (*exec.Cmd, bool) {
	filepath := "."
	if fileName != "" {
		filepath = fileName
	}
	var isNonTerminalEditor bool
	var editorCommand string

	switch strings.ToLower(userSetEditor) {
	case "nano":
		editorCommand = "nano"
		isNonTerminalEditor = false
	case "vim":
		editorCommand = "vim"
		isNonTerminalEditor = false
	case "neovim":
		editorCommand = "nvim"
		isNonTerminalEditor = false
	case "vscode":
		editorCommand = "code"
		isNonTerminalEditor = true
	case "zed":
		editorCommand = "zed"
		isNonTerminalEditor = true
	case "cursor":
		editorCommand = "cursor"
		isNonTerminalEditor = true
	case "windsurf":
		editorCommand = "windsurf"
		isNonTerminalEditor = true
	case "antigravity":
		editorCommand = "antigravity"
		isNonTerminalEditor = true
	default:
		editorCommand = "vi"
		isNonTerminalEditor = false
	}

	cmd := exec.Command(editorCommand, []string{filepath}...)
	return cmd, isNonTerminalEditor
}
