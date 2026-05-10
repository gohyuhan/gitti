package stash

import (
	"sync/atomic"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
)

// ------------------------------------
//
//	GitStashMessagePopUpModel holds a focused text input for the user to type an
//	optional stash message. FilePathName is set when stashing a single file;
//	StashType distinguishes stash-file from stash-all operations.
//
// ------------------------------------
type GitStashMessagePopUpModel struct {
	StashMessageInput textinput.Model
	FilePathName      string
	StashType         string
}

// ------------------------------------
//
//	GitStashOperationOutputPopUpModel holds the output viewport and spinner for a
//	stash operation (stash-all, stash-file, apply, drop, or pop). Atomic flags
//	IsProcessing, HasError, and ProcessSuccess drive the border color and spinner
//	visibility in the render function.
//
// ------------------------------------
type GitStashOperationOutputPopUpModel struct {
	StashOperationType              string
	GitStashOperationOutputViewport viewport.Model // to log out the output from git operation
	Spinner                         spinner.Model  // spinner for showing processing state
	IsProcessing                    atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                        atomic.Bool    // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess                  atomic.Bool    // has the process sucessfuly executed
}

// ------------------------------------
//
//	GitStashConfirmPromptPopUpModel holds the details for the stash confirmation
//	prompt shown before any stash operation executes. StashOperationType selects
//	the localized confirmation message; FilePathName is used for single-file
//	stash; StashMessage and StashId are used for apply/drop/pop confirmations.
//
// ------------------------------------
type GitStashConfirmPromptPopUpModel struct {
	StashOperationType string
	FilePathName       string
	StashMessage       string
	StashId            string
}
