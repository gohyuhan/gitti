package rebase

import (
	"context"
	"sync/atomic"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
)

// ---------------------------------
//
// # A pop up to show git rebase result
//
// ---------------------------------
type GitRebaseOutputPopUpModel struct {
	GitRebaseOutputViewport viewport.Model // to log out the output from git operation
	Spinner                 spinner.Model  // spinner for showing processing state
	IsProcessing            atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                atomic.Bool    // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess          atomic.Bool    // has the process sucessfuly executed
	IsCancelled             atomic.Bool    // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git pull operation
	CancelFunc context.CancelFunc
}

// ---------------------------------
//
// # A pop up for branch input
//
// ---------------------------------
type GitRebaseBranchInputPopUpModel struct {
	BranchNameInput textinput.Model
	Remote          string
}
