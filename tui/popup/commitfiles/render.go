package commitfiles

import (
	"fmt"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

func RenderGitCheckoutCommitFileConfirmPromptPopup(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitCheckoutCommitFileConfirmPromptPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitDiscardConfirmPromptPopupWidth, int(float64(m.Width)*0.8))
		shortHash := popUp.CommitHash
		if len(shortHash) > 7 {
			shortHash = shortHash[:7]
		}
		content := style.NewStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitCheckoutCommitFileConfirmation, popUp.FileDisplayName, shortHash))
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
