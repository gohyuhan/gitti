package reflog

import (
	"fmt"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Render the cherry-pick confirmation popup, displaying the commit hash and
//	HEAD reference in color alongside the reflog action and action detail, so
//	the user can confirm before applying the cherry-pick.
//
// ------------------------------------
func RenderGitCherryPickFromRefLogApplyConfirmationPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitCherryPickFromRefLogApplyConfirmationPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitCherryPickFromRefLogApplyConfirmationPopUpWidth, int(float64(m.Width)*0.8))
		HashString := style.NewStyle.Foreground(style.ColorYellowWarm).Render(popUp.Hash)
		HeadString := style.NewStyle.Foreground(style.ColorPurpleVibrant).Render(popUp.Head)
		content := fmt.Sprintf(i18n.LANGUAGEMAPPING.GitCherryPickFromRefLogApplyConfirmationTitle, HashString, HeadString, popUp.Action, popUp.ActionInfo)

		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
