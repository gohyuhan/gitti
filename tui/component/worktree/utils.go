package worktree

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
)

// those utf-8 icons for the component can be found at https://www.nerdfonts.com/cheat-sheet

// ------------------------------------
//
//	Build the styled tab-bar title string for the tag panel, with the tag tab
//	highlighted as active. Falls back to abbreviated icon/letter labels through up to
//	four levels when the rendered width exceeds titleWidthLimit.
//
// ------------------------------------
func ConstructWorktreeComponentTitle(titleWidthLimit int) string {
	title := fmt.Sprintf(
		"%s %s %s %s %s %s %s %s %s %s %s %s",
		style.TitleCurrentComponentStyle.Render("[1]"),
		style.TitleNonCurrentComponentStyle.Render("\uf418"),
		style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Branches),
		style.TitleNonCurrentComponentStyle.Render("•"),
		style.TitleNonCurrentComponentStyle.Render("\uf412"),
		style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Tag),
		style.TitleNonCurrentComponentStyle.Render("•"),
		style.TitleNonCurrentComponentStyle.Render("\ueb39"),
		style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Remote),
		style.TitleNonCurrentComponentStyle.Render("•"),
		style.TitleCurrentComponentStyle.Render("\ueb86"),
		style.TitleCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Worktree),
	)
	if lipgloss.Width(title) > titleWidthLimit {
		title = fmt.Sprintf(
			"%s %s %s %s %s %s %s %s %s %s %s %s",
			style.TitleCurrentComponentStyle.Render("[1]"),
			style.TitleNonCurrentComponentStyle.Render("\uf418"),
			style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Branches),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\uf412"),
			style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Tag),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\ueb39"),
			style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Remote),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleCurrentComponentStyle.Render("\ueb86"),
			style.TitleCurrentComponentStyle.Render("W"),
		)
	}
	if lipgloss.Width(title) > titleWidthLimit {
		title = fmt.Sprintf(
			"%s %s %s %s %s %s %s %s %s %s %s %s",
			style.TitleCurrentComponentStyle.Render("[1]"),
			style.TitleNonCurrentComponentStyle.Render("\uf418"),
			style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Branches),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\uf412"),
			style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Tag),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\ueb39"),
			style.TitleNonCurrentComponentStyle.Render("R"),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleCurrentComponentStyle.Render("\ueb86"),
			style.TitleCurrentComponentStyle.Render("W"),
		)
	}
	if lipgloss.Width(title) > titleWidthLimit {
		title = fmt.Sprintf(
			"%s %s %s %s %s %s %s %s %s %s %s %s",
			style.TitleCurrentComponentStyle.Render("[1]"),
			style.TitleNonCurrentComponentStyle.Render("\uf418"),
			style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Branches),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\uf412"),
			style.TitleNonCurrentComponentStyle.Render("T"),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\ueb39"),
			style.TitleNonCurrentComponentStyle.Render("R"),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleCurrentComponentStyle.Render("\ueb86"),
			style.TitleCurrentComponentStyle.Render("W"),
		)
	}
	if lipgloss.Width(title) > titleWidthLimit {
		title = fmt.Sprintf(
			"%s %s %s %s %s %s %s %s %s %s %s %s",
			style.TitleCurrentComponentStyle.Render("[1]"),
			style.TitleNonCurrentComponentStyle.Render("\uf418"),
			style.TitleNonCurrentComponentStyle.Render("B"),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\uf412"),
			style.TitleNonCurrentComponentStyle.Render("T"),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\ueb39"),
			style.TitleNonCurrentComponentStyle.Render("R"),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleCurrentComponentStyle.Render("\ueb86"),
			style.TitleCurrentComponentStyle.Render("W"),
		)
	}
	return title
}
