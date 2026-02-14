package commitlog

import (
	"fmt"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
)

func ConstructCommitLogComponentTitle(widthLimit int) string {
	return style.TitleCurrentComponentStyle.Render(fmt.Sprintf("[3] \ue729 %s", i18n.LANGUAGEMAPPING.CommitLog))
}
