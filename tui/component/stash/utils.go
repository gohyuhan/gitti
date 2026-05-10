package stash

import (
	"fmt"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
)

// those utf-8 icons for the component can be found at https://www.nerdfonts.com/cheat-sheet

// ------------------------------------
//
//	Build the styled title string for the stash panel, rendering the panel number,
//	icon, and localized label in the active/highlighted style.
//
// ------------------------------------
func ConstructStashComponentTitle(widthLimit int) string {
	return style.TitleCurrentComponentStyle.Render(fmt.Sprintf("[4] \uf0c7 %s", i18n.LANGUAGEMAPPING.Stash))
}
