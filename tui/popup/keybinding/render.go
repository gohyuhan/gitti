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
//	For Key binding/Feature Instructions pop up
//
// ------------------------------------
func RenderKeyBindingAndFeatureInstructionsPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*KeybindingAndFeatureInstructionsPopUpModel)
	if ok {
		var contentLine strings.Builder
		contentLine.WriteString("\n")

		// context aware, only show the key binding of the current user selected component panel related keybinding
		renderSelectedComponentKeyBindingPart(m, &contentLine)

		// Global key binding and its instructions
		renderGlobalKeyBindingPart(m, &contentLine)

		// available features and its instuctions and steps
		renderFeaturesInstuctionsAndStepsPart(&contentLine)

		height := max(constant.PopUpGlobalKeyBindingViewPortHeight, int(float64(m.Height)*0.8))
		width := min(constant.MaxKeybindingAndFeatureInstructionsPopUpWidth, int(float64(m.Width)*0.8)-4)
		popUp.GlobalKeyBindingViewport.SetWidth(width)
		popUp.GlobalKeyBindingViewport.SetYOffset(popUp.GlobalKeyBindingViewport.YOffset())
		popUp.GlobalKeyBindingViewport.SetHeight(height)
		popUp.GlobalKeyBindingViewport.SetContent(contentLine.String())
		return style.KeyBindingPopUpStyle.Render(popUp.GlobalKeyBindingViewport.View())
	}
	return ""
}

func renderGlobalKeyBindingPart(m *types.GittiModel, contentLine *strings.Builder) {
	// this will usually only be run once for the entire gitti session
	// to determine the largest keybinding line so that everything is render nicely
	if m.GlobalKeyBindingKeyMapLargestLen < 1 {
		maxLen := 0
		for index := range i18n.LANGUAGEMAPPING.GlobalKeyBinding {
			if l := utf8.RuneCountInString(i18n.LANGUAGEMAPPING.GlobalKeyBinding[index].KeyBindingLine); l > maxLen {
				maxLen = l
			}
		}
		m.GlobalKeyBindingKeyMapLargestLen = maxLen
	}

	for index := range i18n.LANGUAGEMAPPING.GlobalKeyBinding {
		globalKeybinding := i18n.LANGUAGEMAPPING.GlobalKeyBinding[index]
		contentLine.WriteString(" ")
		switch globalKeybinding.LineType {
		case i18n.TITLE:
			contentLine.WriteString(globalKeybinding.KeyBindingLine)
			pad := m.GlobalKeyBindingKeyMapLargestLen - utf8.RuneCountInString(globalKeybinding.KeyBindingLine)
			if pad > 0 {
				for range pad {
					contentLine.WriteString(" ")
				}
			}
			contentLine.WriteString("  ")
			contentLine.WriteString(style.KeyBindingAndFeatureInstructionsTitleLineStyle.Render(globalKeybinding.TitleOrInfoLine))
		case i18n.INFO:
			contentLine.WriteString(style.KeyBindingKeyMappingInfoStyle.Render(globalKeybinding.KeyBindingLine))
			pad := m.GlobalKeyBindingKeyMapLargestLen - utf8.RuneCountInString(globalKeybinding.KeyBindingLine)
			if pad > 0 {
				for range pad {
					contentLine.WriteString(" ")
				}
			}
			contentLine.WriteString("  ")
			contentLine.WriteString(globalKeybinding.TitleOrInfoLine)
		case i18n.WARN:
			contentLine.WriteString(style.KeyBindingKeyMappingInfoStyle.Render(globalKeybinding.KeyBindingLine))
			contentLine.WriteString(style.KeyBindingAndFeatureInstructionsWarnStyle.Render(globalKeybinding.TitleOrInfoLine))
		}

		contentLine.WriteString(" \n")
	}
	contentLine.WriteString("\n")
	contentLine.WriteString("\n")
	contentLine.WriteString("\n")
}

func renderSelectedComponentKeyBindingPart(m *types.GittiModel, contentLine *strings.Builder) {
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
		for index := range selectedComponenti18nKeybinding {
			if l := utf8.RuneCountInString(selectedComponenti18nKeybinding[index].KeyBindingLine); l > maxLen {
				maxLen = l
			}
		}
		*selectedComponentKeyBindingKeyMapLargestLen = maxLen
	}

	for index := range selectedComponenti18nKeybinding {
		selectedComponentKeyBinding := selectedComponenti18nKeybinding[index]
		contentLine.WriteString(" ")
		switch selectedComponentKeyBinding.LineType {
		case i18n.TITLE:
			contentLine.WriteString(selectedComponentKeyBinding.KeyBindingLine)
			pad := *selectedComponentKeyBindingKeyMapLargestLen - utf8.RuneCountInString(selectedComponentKeyBinding.KeyBindingLine)
			if pad > 0 {
				for range pad {
					contentLine.WriteString(" ")
				}
			}
			contentLine.WriteString("  ")
			contentLine.WriteString(style.KeyBindingAndFeatureInstructionsTitleLineStyle.Render(selectedComponentKeyBinding.TitleOrInfoLine))
		case i18n.INFO:
			contentLine.WriteString(style.KeyBindingKeyMappingInfoStyle.Render(selectedComponentKeyBinding.KeyBindingLine))
			pad := *selectedComponentKeyBindingKeyMapLargestLen - utf8.RuneCountInString(selectedComponentKeyBinding.KeyBindingLine)
			if pad > 0 {
				for range pad {
					contentLine.WriteString(" ")
				}
			}
			contentLine.WriteString("  ")
			contentLine.WriteString(selectedComponentKeyBinding.TitleOrInfoLine)
		case i18n.WARN:
			contentLine.WriteString(style.KeyBindingKeyMappingInfoStyle.Render(selectedComponentKeyBinding.KeyBindingLine))
			contentLine.WriteString(style.KeyBindingAndFeatureInstructionsWarnStyle.Render(selectedComponentKeyBinding.TitleOrInfoLine))
		}

		contentLine.WriteString(" \n")
	}
	contentLine.WriteString("\n")
	contentLine.WriteString("\n")
	contentLine.WriteString("\n")
}

func renderFeaturesInstuctionsAndStepsPart(contentLine *strings.Builder) {
	for index := range i18n.LANGUAGEMAPPING.FeatureInstructions {
		featureInstruction := i18n.LANGUAGEMAPPING.FeatureInstructions[index]
		contentLine.WriteString(" ")
		switch featureInstruction.LineType {
		case i18n.TITLE:
			contentLine.WriteString(featureInstruction.Feature)
			contentLine.WriteString("       ")
			contentLine.WriteString(style.KeyBindingAndFeatureInstructionsTitleLineStyle.Render(featureInstruction.InstructionLines...))
			contentLine.WriteString(" \n")
		case i18n.INFO:
			contentLine.WriteString(style.FeatureInfoLineStyle.Bold(true).Render("// " + featureInstruction.Feature + " //"))
			for index := range featureInstruction.InstructionLines {
				instructionLine := featureInstruction.InstructionLines[index]
				contentLine.WriteString(" \n")
				contentLine.WriteString("   ")
				contentLine.WriteString(instructionLine)
			}
			contentLine.WriteString(" \n")
		case i18n.WARN:
			contentLine.WriteString(style.FeatureInfoLineStyle.Render(featureInstruction.Feature))
			contentLine.WriteString(style.KeyBindingAndFeatureInstructionsWarnStyle.Render(featureInstruction.InstructionLines...))
		}

		contentLine.WriteString(" \n")
	}
	contentLine.WriteString("\n")
	contentLine.WriteString("\n")
	contentLine.WriteString("\n")
}
