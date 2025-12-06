package keybinding

import (
	"fmt"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	For Global Key binding pop up
//
// ------------------------------------
func RenderGlobalKeyBindingPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GlobalKeyBindingPopUpModel)
	if ok {
		keyBindingLine := "\n"

		// this will usually only be run once for the entire gitti session
		if m.GlobalKeyBindingKeyMapLargestLen < 1 {
			maxLen := 0
			for _, line := range i18n.LANGUAGEMAPPING.GlobalKeyBinding {
				if l := len(line.KeyBindingLine); l > maxLen {
					maxLen = l
				}
			}
			m.GlobalKeyBindingKeyMapLargestLen = maxLen
		}
		for _, line := range i18n.LANGUAGEMAPPING.GlobalKeyBinding {
			switch line.LineType {
			case i18n.TITLE:
				keyBindingLine += " " + fmt.Sprintf("%*s", m.GlobalKeyBindingKeyMapLargestLen, line.KeyBindingLine) +
					"  " +
					style.GlobalKeyBindingTitleLineStyle.Render(line.TitleOrInfoLine) +
					"\n"
			case i18n.INFO:
				keyBindingLine += " " + style.GlobalKeyBindingKeyMappingLineStyle.Render(fmt.Sprintf("%*s", m.GlobalKeyBindingKeyMapLargestLen, line.KeyBindingLine)) +
					"  " +
					line.TitleOrInfoLine +
					"\n"
			case i18n.WARN:
				keyBindingLine += " " + style.GlobalKeyBindingKeyMappingLineStyle.Render(fmt.Sprintf("%s", line.KeyBindingLine)) +
					line.TitleOrInfoLine +
					"\n"
			}
		}
		height := min(constant.PopUpGlobalKeyBindingViewPortHeight, int(float64(m.Height)*0.8))
		width := min(constant.MaxGlobalKeyBindingPopUpWidth, int(float64(m.Width)*0.8)-4)
		popUp.GlobalKeyBindingViewport.SetWidth(width)
		popUp.GlobalKeyBindingViewport.SetYOffset(popUp.GlobalKeyBindingViewport.YOffset())
		popUp.GlobalKeyBindingViewport.SetHeight(height)
		popUp.GlobalKeyBindingViewport.SetContent(keyBindingLine)
		return style.GlobalKeyBindingPopUpStyle.Render(popUp.GlobalKeyBindingViewport.View())
	}
	return ""
}
