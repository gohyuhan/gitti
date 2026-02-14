package branch

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
)

func ConstructLocalBranchComponentTitle(titleWidthLimit int) string {
	title := fmt.Sprintf("%s %s %s %s %s %s %s %s %s",
		style.TitleCurrentComponentStyle.Render("[1]"),
		style.TitleCurrentComponentStyle.Render("\uf418"),
		style.TitleCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Branches),
		style.TitleNonCurrentComponentStyle.Render("•"),
		style.TitleNonCurrentComponentStyle.Render("\uf412"),
		style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Tag),
		style.TitleNonCurrentComponentStyle.Render("•"),
		style.TitleNonCurrentComponentStyle.Render("\ueb39"),
		style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Remote),
	)
	if lipgloss.Width(title) > titleWidthLimit {
		title = fmt.Sprintf("%s %s %s %s %s %s %s %s %s",
			style.TitleCurrentComponentStyle.Render("[1]"),
			style.TitleCurrentComponentStyle.Render("\uf418"),
			style.TitleCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Branches),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\uf412"),
			style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Tag),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\ueb39"),
			style.TitleNonCurrentComponentStyle.Render("R"),
		)
	}
	if lipgloss.Width(title) > titleWidthLimit {
		title = fmt.Sprintf("%s %s %s %s %s %s %s %s %s",
			style.TitleCurrentComponentStyle.Render("[1]"),
			style.TitleCurrentComponentStyle.Render("\uf418"),
			style.TitleCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Branches),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\uf412"),
			style.TitleNonCurrentComponentStyle.Render("T"),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\ueb39"),
			style.TitleNonCurrentComponentStyle.Render("R"),
		)
	}
	if lipgloss.Width(title) > titleWidthLimit {
		title = fmt.Sprintf("%s %s %s %s %s %s %s %s %s",
			style.TitleCurrentComponentStyle.Render("[1]"),
			style.TitleCurrentComponentStyle.Render("\uf418"),
			style.TitleCurrentComponentStyle.Render("B"),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\uf412"),
			style.TitleNonCurrentComponentStyle.Render("T"),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleNonCurrentComponentStyle.Render("\ueb39"),
			style.TitleNonCurrentComponentStyle.Render("R"),
		)
	}
	return title
}
