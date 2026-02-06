package tag

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
)

func ConstructTagComponentTitle(titleWidthLimit int) string {
	title := fmt.Sprintf("%s %s %s %s %s %s",
		style.TitleCurrentComponentStyle.Render("[1]"),
		style.TitleNonCurrentComponentStyle.Render("\uf418"),
		style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Branches),
		style.TitleNonCurrentComponentStyle.Render("•"),
		style.TitleCurrentComponentStyle.Render("\uf412"),
		style.TitleCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Tag),
	)
	if lipgloss.Width(title) > titleWidthLimit {
		title = fmt.Sprintf("%s %s %s %s %s %s",
			style.TitleCurrentComponentStyle.Render("[1]"),
			style.TitleNonCurrentComponentStyle.Render("\uf418"),
			style.TitleNonCurrentComponentStyle.Render(i18n.LANGUAGEMAPPING.Branches),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleCurrentComponentStyle.Render("\uf412"),
			style.TitleCurrentComponentStyle.Render("T"),
		)
	}
	if lipgloss.Width(title) > titleWidthLimit {
		title = fmt.Sprintf("%s %s %s %s %s %s",
			style.TitleCurrentComponentStyle.Render("[1]"),
			style.TitleNonCurrentComponentStyle.Render("\uf418"),
			style.TitleNonCurrentComponentStyle.Render("B"),
			style.TitleNonCurrentComponentStyle.Render("•"),
			style.TitleCurrentComponentStyle.Render("\uf412"),
			style.TitleCurrentComponentStyle.Render("T"),
		)
	}
	return title
}
