package files

import "charm.land/bubbles/v2/viewport"

// ------------------------------------
//
//	GitDiscardFileLineChangeConfirmPopUpModel holds a viewport that displays the
//	single diff line the user is about to discard, used in the confirmation popup.
//
// ------------------------------------
type GitDiscardFileLineChangeConfirmPopUpModel struct {
	DiscardFileLineChangeViewport viewport.Model
}
