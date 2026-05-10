package reflog

// ------------------------------------
//
//	GitCherryPickFromRefLogApplyConfirmationPopUpModel holds the commit hash,
//	HEAD reference, reflog action, and action detail for the cherry-pick
//	confirmation popup triggered from the reflog panel.
//
// ------------------------------------
type GitCherryPickFromRefLogApplyConfirmationPopUpModel struct {
	Hash       string
	Head       string
	Action     string
	ActionInfo string
}
