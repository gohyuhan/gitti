package reflog

import "github.com/gohyuhan/gitti/tui/types"

// ------------------------------------
//
//	Initialize the cherry-pick confirmation popup model with the commit hash,
//	HEAD reference, reflog action, and action detail string for the entry the
//	user selected in the reflog panel.
//
// ------------------------------------
func InitGitCherryPickFromRefLogApplyConfirmationPopUpModel(m *types.GittiModel, hash string, head string, action string, actionInfo string) {
	popUpModel := &GitCherryPickFromRefLogApplyConfirmationPopUpModel{
		Hash:       hash,
		Head:       head,
		Action:     action,
		ActionInfo: actionInfo,
	}

	m.PopUpModel = popUpModel
}
