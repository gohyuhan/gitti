package commit

import (
	"context"
	"sync/atomic"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
)

// ---------------------------------
// # For git commit process pop up model
//
// ---------------------------------
type GitCommitPopUpModel struct {
	IsAmendCommit            bool            // to indicate is this is a normal commit or an amend commit operation
	MessageTextInput         textinput.Model // input index 1
	DescriptionTextAreaInput textarea.Model  // input index 2
	TotalInputCount          int             // to tell us how many input were there
	CurrentActiveInputIndex  int             // to tell us which input should be shown as highlighted/focus and be updated
	GitCommitOutputViewport  viewport.Model  // to log out the output from git operation
	Spinner                  spinner.Model   // spinner for showing processing state
	InitialCommitStarted     atomic.Bool     // indicated that this pop up session has start the first commit action
	IsProcessing             atomic.Bool     // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                 atomic.Bool     // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess           atomic.Bool     // has the process sucessfuly executed
	IsCancelled              atomic.Bool     // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git commit operation
	CancelFunc context.CancelFunc
}

// ---------------------------------
//
// # For git amend commit process pop up model
//
// ---------------------------------
type GitAmendCommitPopUpModel struct {
	IsAmendCommit                bool            // to indicate is this is a normal commit or an amend commit operation
	MessageTextInput             textinput.Model // input index 1
	DescriptionTextAreaInput     textarea.Model  // input index 2
	TotalInputCount              int             // to tell us how many input were there
	CurrentActiveInputIndex      int             // to tell us which input should be shown as highlighted/focus and be updated
	GitAmendCommitOutputViewport viewport.Model  // to log out the output from git operation
	Spinner                      spinner.Model   // spinner for showing processing state
	InitialCommitStarted         atomic.Bool     // indicated that this pop up session has start the first commit action
	IsProcessing                 atomic.Bool     // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                     atomic.Bool     // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess               atomic.Bool     // has the process sucessfuly executed
	IsCancelled                  atomic.Bool     // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git amend commit operation
	CancelFunc context.CancelFunc
}
