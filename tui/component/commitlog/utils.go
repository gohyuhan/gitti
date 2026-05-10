package commitlog

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
)

// those utf-8 icons for the component can be found at https://www.nerdfonts.com/cheat-sheet

// ------------------------------------
//
//	Build the styled tab-bar title string for the commit log panel, with the commit log tab
//	highlighted as active. Falls back to abbreviated icon/letter labels when the rendered
//	width exceeds widthLimit.
//
// ------------------------------------
func ConstructCommitLogComponentTitle(widthLimit int) string {
	title := fmt.Sprintf("%s %s %s %s %s %s",
		style.TitleCurrentComponentStyle.Render("[3]"),
		style.TitleCurrentComponentStyle.Render("\ue729"),
		style.TitleCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.CommitLog),
		style.TitleNonCurrentComponentStyle.Render("•"),
		style.TitleNonCurrentComponentStyle.Render("\uf4ed"),
		style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.RefLog),
	)
	if lipgloss.Width(title) > widthLimit {
		title = fmt.Sprintf("%s %s %s %s %s %s",
			style.TitleCurrentComponentStyle.Render("[3]"),
			style.TitleCurrentComponentStyle.Render("\ue729"),
			style.TitleCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.CommitLog),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\uf4ed"),
			style.TitleNonCurrentComponentStyle.Render("R"),
		)
	}
	if lipgloss.Width(title) > widthLimit {
		title = fmt.Sprintf("%s %s %s %s %s %s",
			style.TitleCurrentComponentStyle.Render("[3]"),
			style.TitleCurrentComponentStyle.Render("\ue729"),
			style.TitleCurrentComponentStyle.Render("C"),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\uf4ed"),
			style.TitleNonCurrentComponentStyle.Render("R"),
		)
	}
	return title
}
