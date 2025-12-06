package stash

import (
	"sync/atomic"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
)

// ---------------------------------
//
// # A pop up to prompt for git stash message
//
// ---------------------------------
type GitStashMessagePopUpModel struct {
	StashMessageInput textinput.Model
	FilePathName      string
	StashType         string
}

// ---------------------------------
//
// for stash operation pop up with viewport (for showing the result of stash operation)
//
// ---------------------------------
type GitStashOperationOutputPopUpModel struct {
	StashOperationType              string
	GitStashOperationOutputViewport viewport.Model // to log out the output from git operation
	Spinner                         spinner.Model  // spinner for showing processing state
	IsProcessing                    atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                        atomic.Bool    // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess                  atomic.Bool    // has the process sucessfuly executed
}

// ---------------------------------
//
// for stash operation confirm prompt pop up (for prompting user for confirmation)
//   - for stash, stash all, drop, apply, discard
//
// ---------------------------------
type GitStashConfirmPromptPopUpModel struct {
	StashOperationType string
	FilePathName       string
	StashMessage       string
	StashId            string
}
