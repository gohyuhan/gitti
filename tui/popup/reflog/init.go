package reflog

import "github.com/gohyuhan/gitti/tui/types"

func InitGitCherryPickFromRefLogApplyConfirmationPopUpModel(m *types.GittiModel, hash string, head string, action string, actionInfo string) {
	popUpModel := &GitCherryPickFromRefLogApplyConfirmationPopUpModel{
		Hash:       hash,
		Head:       head,
		Action:     action,
		ActionInfo: actionInfo,
	}

	m.PopUpModel = popUpModel
}
