package files

import (
	"fmt"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
)

func ConstructModifiedFilesComponentTitle(widthLimit int) string {
	return style.TitleCurrentComponentStyle.Render(fmt.Sprintf("[2] \ueae9 %s", i18n.LANGUAGEMAPPING.ModifiedFiles))
}
