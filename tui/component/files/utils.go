package files

import (
	"fmt"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
)

// those utf-8 icons for the component can be found at https://www.nerdfonts.com/cheat-sheet

// ------------------------------------
//
//	Build the styled title string for the modified files panel, rendering the panel number,
//	icon, and localized label in the active/highlighted style.
//
// ------------------------------------
func ConstructModifiedFilesComponentTitle(widthLimit int) string {
	return style.TitleCurrentComponentStyle.Render(fmt.Sprintf("[2] \ueae9 %s", i18n.LANGUAGEMAPPING.ModifiedFiles))
}
