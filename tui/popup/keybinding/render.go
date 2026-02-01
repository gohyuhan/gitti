package keybinding

import (
	"strings"
	"unicode/utf8"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	For Key binding and Instructions pop up
//
// ------------------------------------
func RenderKeyBindingPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GlobalKeyBindingPopUpModel)
	if ok {
		var keyBindingLine strings.Builder
		keyBindingLine.WriteString("\n")

		// context aware, only show the key binding of the current user selected component panel related keybinding
		renderSelectedComponentKeyBindingPart(m, &keyBindingLine)

		// Global key binding and its instructions
		renderGlobalKeyBindingPart(m, &keyBindingLine)

		height := max(constant.PopUpGlobalKeyBindingViewPortHeight, int(float64(m.Height)*0.8))
		width := min(constant.MaxGlobalKeyBindingPopUpWidth, int(float64(m.Width)*0.8)-4)
		popUp.GlobalKeyBindingViewport.SetWidth(width)
		popUp.GlobalKeyBindingViewport.SetYOffset(popUp.GlobalKeyBindingViewport.YOffset())
		popUp.GlobalKeyBindingViewport.SetHeight(height)
		popUp.GlobalKeyBindingViewport.SetContent(keyBindingLine.String())
		return style.KeyBindingPopUpStyle.Render(popUp.GlobalKeyBindingViewport.View())
	}
	return ""
}

func renderGlobalKeyBindingPart(m *types.GittiModel, keyBindingLine *strings.Builder) {
	// this will usually only be run once for the entire gitti session
	// to determine the largest keybinding line so that everything is render nicely
	if m.GlobalKeyBindingKeyMapLargestLen < 1 {
		maxLen := 0
		for _, line := range i18n.LANGUAGEMAPPING.GlobalKeyBinding {
			if l := utf8.RuneCountInString(line.KeyBindingLine); l > maxLen {
				maxLen = l
			}
		}
		m.GlobalKeyBindingKeyMapLargestLen = maxLen
	}

	for _, line := range i18n.LANGUAGEMAPPING.GlobalKeyBinding {
		keyBindingLine.WriteString(" ")
		switch line.LineType {
		case i18n.TITLE:
			keyBindingLine.WriteString(line.KeyBindingLine)
			pad := m.GlobalKeyBindingKeyMapLargestLen - utf8.RuneCountInString(line.KeyBindingLine)
			if pad > 0 {
				for range pad {
					keyBindingLine.WriteString(" ")
				}
			}
			keyBindingLine.WriteString("  ")
			keyBindingLine.WriteString(style.KeyBindingTitleLineStyle.Render(line.TitleOrInfoLine))
		case i18n.INFO:
			keyBindingLine.WriteString(style.KeyBindingKeyMappingLineStyle.Render(line.KeyBindingLine))
			pad := m.GlobalKeyBindingKeyMapLargestLen - utf8.RuneCountInString(line.KeyBindingLine)
			if pad > 0 {
				for range pad {
					keyBindingLine.WriteString(" ")
				}
			}
			keyBindingLine.WriteString("  ")
			keyBindingLine.WriteString(line.TitleOrInfoLine)
		case i18n.WARN:
			keyBindingLine.WriteString(style.KeyBindingKeyMappingLineStyle.Render(line.KeyBindingLine))
			keyBindingLine.WriteString(style.KeyBindingKeyMappingWarnStyle.Render(line.TitleOrInfoLine))
		}

		keyBindingLine.WriteString(" \n")
	}
}

func renderSelectedComponentKeyBindingPart(m *types.GittiModel, keyBindingLine *strings.Builder) {
	// this will usually only be run once for the entire gitti session
	// to determine the largest keybinding line so that everything is render nicely

	var selectedComponenti18nKeybinding []i18n.KeyBindingMappingFormat
	var selectedComponentKeyBindingKeyMapLargestLen *int
	switch m.CurrentSelectedComponent {
	case constant.LocalBranchComponent:
		selectedComponentKeyBindingKeyMapLargestLen = &m.LocalBranchComponentKeyBindingKeyMapLargestLen
		selectedComponenti18nKeybinding = i18n.LANGUAGEMAPPING.LocalBranchComponentKeyBinding
	case constant.ModifiedFilesComponent:
		selectedComponentKeyBindingKeyMapLargestLen = &m.ModifiedFilesComponentKeyBindingKeyMapLargestLen
		selectedComponenti18nKeybinding = i18n.LANGUAGEMAPPING.ModifiedFilesComponentKeyBinding
	case constant.CommitLogComponent:
		selectedComponentKeyBindingKeyMapLargestLen = &m.CommitLogComponentKeyBindingKeyMapLargestLen
		selectedComponenti18nKeybinding = i18n.LANGUAGEMAPPING.CommitLogComponentKeyBinding
	case constant.StashComponent:
		selectedComponentKeyBindingKeyMapLargestLen = &m.StashComponentKeyBindingKeyMapLargestLen
		selectedComponenti18nKeybinding = i18n.LANGUAGEMAPPING.StashComponentKeyBinding
	case constant.DetailComponent, constant.DetailComponentTwo:
		selectedComponentKeyBindingKeyMapLargestLen = &m.DetailComponentKeyBindingKeyMapLargestLen
		selectedComponenti18nKeybinding = i18n.LANGUAGEMAPPING.DetailComponentKeyBinding
	case constant.LogComponent:
		selectedComponentKeyBindingKeyMapLargestLen = &m.LogComponentKeyBindingKeyMapLargestLen
		selectedComponenti18nKeybinding = i18n.LANGUAGEMAPPING.LogComponentKeyBinding
	case constant.GitStatusComponent:
		return
	}

	if *selectedComponentKeyBindingKeyMapLargestLen < 1 {
		maxLen := 0
		for _, line := range selectedComponenti18nKeybinding {
			if l := utf8.RuneCountInString(line.KeyBindingLine); l > maxLen {
				maxLen = l
			}
		}
		*selectedComponentKeyBindingKeyMapLargestLen = maxLen
	}

	for _, line := range selectedComponenti18nKeybinding {
		keyBindingLine.WriteString(" ")
		switch line.LineType {
		case i18n.TITLE:
			keyBindingLine.WriteString(line.KeyBindingLine)
			pad := *selectedComponentKeyBindingKeyMapLargestLen - utf8.RuneCountInString(line.KeyBindingLine)
			if pad > 0 {
				for range pad {
					keyBindingLine.WriteString(" ")
				}
			}
			keyBindingLine.WriteString("  ")
			keyBindingLine.WriteString(style.KeyBindingTitleLineStyle.Render(line.TitleOrInfoLine))
		case i18n.INFO:
			keyBindingLine.WriteString(style.KeyBindingKeyMappingLineStyle.Render(line.KeyBindingLine))
			pad := *selectedComponentKeyBindingKeyMapLargestLen - utf8.RuneCountInString(line.KeyBindingLine)
			if pad > 0 {
				for range pad {
					keyBindingLine.WriteString(" ")
				}
			}
			keyBindingLine.WriteString("  ")
			keyBindingLine.WriteString(line.TitleOrInfoLine)
		case i18n.WARN:
			keyBindingLine.WriteString(style.KeyBindingKeyMappingLineStyle.Render(line.KeyBindingLine))
			keyBindingLine.WriteString(style.KeyBindingKeyMappingWarnStyle.Render(line.TitleOrInfoLine))
		}

		keyBindingLine.WriteString(" \n")
	}
	keyBindingLine.WriteString("\n")
	keyBindingLine.WriteString("\n")
	keyBindingLine.WriteString("\n")
}
