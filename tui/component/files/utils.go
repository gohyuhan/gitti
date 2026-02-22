package files

import (
	"fmt"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
)

// those utf-8 icons for the component can be found at https://www.nerdfonts.com/cheat-sheet

func ConstructModifiedFilesComponentTitle(widthLimit int) string {
	return style.TitleCurrentComponentStyle.Render(fmt.Sprintf("[2] \ueae9 %s", i18n.LANGUAGEMAPPING.ModifiedFiles))
}
