package keybinding

import (
	"charm.land/bubbles/v2/viewport"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
)

// init the viewport pop up for showing info of global key binding
func InitGlobalKeyBindingPopUpModel(m *types.GittiModel) {
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(min(constant.PopUpGlobalKeyBindingViewPortHeight, int(float64(m.Height)*0.8)))
	vp.SetWidth(min(constant.MaxGlobalKeyBindingPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	m.PopUpModel = &GlobalKeyBindingPopUpModel{
		GlobalKeyBindingViewport: vp,
	}
}
