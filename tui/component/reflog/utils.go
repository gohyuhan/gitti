package reflog

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
)

// those utf-8 icons for the component can be found at https://www.nerdfonts.com/cheat-sheet

func ConstructRefLogComponentTitle(widthLimit int) string {
	title := fmt.Sprintf("%s %s %s %s %s %s",
		style.TitleCurrentComponentStyle.Render("[3]"),
		style.TitleNonCurrentComponentStyle.Render("\ue729"),
		style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.CommitLog),
		style.TitleNonCurrentComponentStyle.Render("•"),
		style.TitleCurrentComponentStyle.Render("\uf4ed"),
		style.TitleCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.RefLog),
	)
	if lipgloss.Width(title) > widthLimit {
		title = fmt.Sprintf("%s %s %s %s %s %s",
			style.TitleCurrentComponentStyle.Render("[3]"),
			style.TitleNonCurrentComponentStyle.Render("\ue729"),
			style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.CommitLog),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleCurrentComponentStyle.Render("\uf4ed"),
			style.TitleCurrentComponentStyle.Render("R"),
		)
	}
	if lipgloss.Width(title) > widthLimit {
		title = fmt.Sprintf("%s %s %s %s %s %s",
			style.TitleCurrentComponentStyle.Render("[3]"),
			style.TitleNonCurrentComponentStyle.Render("\ue729"),
			style.TitleNonCurrentComponentStyle.Render("C"),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleCurrentComponentStyle.Render("\uf4ed"),
			style.TitleCurrentComponentStyle.Render("R"),
		)
	}
	return title
}
