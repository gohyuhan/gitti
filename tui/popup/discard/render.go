package discard

import (
	"fmt"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"

	"charm.land/lipgloss/v2"
)

// ------------------------------------
//
//	For Discard file changes type list selection
//
// ------------------------------------
func RenderGitDiscardTypeOptionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitDiscardTypeOptionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitDiscardTypeOptionPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitDiscardTypeOptionTitle)
		popUp.DiscardTypeOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.DiscardTypeOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	For Discard file changes confirmation prompt
//
// ------------------------------------
func RenderGitDiscardConfirmPromptPopup(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitDiscardConfirmPromptPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitDiscardConfirmPromptPopupWidth, int(float64(m.Width)*0.8))
		var content string
		switch popUp.DiscardType {
		case git.DISCARDWHOLE:
			content = style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardWholeConfirmation, popUp.FilePathName))
		case git.DISCARDUNSTAGE:
			content = style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardUnstageConfirmation, popUp.FilePathName))
		case git.DISCARDUNTRACKED:
			content = style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardUntrackedConfirmation, popUp.FilePathName))
		case git.DISCARDNEWLYADDEDORCOPIED:
			content = style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardNewlyAddedorCopyConfirmation, popUp.FilePathName))
		case git.DISCARDANDREVERTRENAME:
			content = style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardAndRevertRenameConfirmation, popUp.FilePathName))
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
