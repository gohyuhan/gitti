package stash

import (
	"fmt"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
)

func ConstructStashComponentTitle(widthLimit int) string {
	return style.TitleCurrentComponentStyle.Render(fmt.Sprintf("[4] \uf0c7 %s", i18n.LANGUAGEMAPPING.Stash))
}
