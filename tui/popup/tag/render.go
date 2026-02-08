package tag

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

func RenderCreateTagPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*CreateTagPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxCreateTagPopUpWidth, int(float64(m.Width)*0.8))
		popUp.TagNameInput.SetWidth(popUpWidth - 6)
		popUp.TagMessageTextAreaInput.SetWidth(popUpWidth - 6)

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			i18n.LANGUAGEMAPPING.CreateTagPopUpNameTitle,
			popUp.TagNameInput.View(),
			i18n.LANGUAGEMAPPING.CreateTagPopUpMessageTitle,
			popUp.TagMessageTextAreaInput.View(),
		)

		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

func RenderCreateTagConfirmationPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*CreateTagConfirmationPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxCreateTagConfirmationPopUp, int(float64(m.Width)*0.8))
		content := fmt.Sprintf(
			i18n.LANGUAGEMAPPING.CreateTagConfirmation,
			style.NewStyle.Foreground(style.ColorYellowWarm).Render(popUp.CommitHash),
			style.NewStyle.Foreground(style.ColorYellowWarm).Render(popUp.CommitMessage),
			style.NewStyle.Foreground(style.ColorPurpleVibrant).Render(popUp.TagName),
			style.NewStyle.Foreground(style.ColorBlueMuted).Render(popUp.TagMessage),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
