package commitlog

import (
	"fmt"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
)

// those utf-8 icons for the component can be found at https://www.nerdfonts.com/cheat-sheet

func ConstructCommitLogComponentTitle(widthLimit int) string {
	return style.TitleCurrentComponentStyle.Render(fmt.Sprintf("[3] \ue729 %s", i18n.LANGUAGEMAPPING.CommitLog))
}
