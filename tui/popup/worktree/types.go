package worktree

import (
	"context"
	"sync/atomic"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
)

// ------------------------------------
//
//	Holds the input state for the add-new-worktree pop-up
//
// ------------------------------------
type WorktreeAddNewWorktreePopUpModel struct {
	WorktreeNameTextInput       textinput.Model // input index 1
	WorktreeBranchNameTextInput textinput.Model // input index 2
	TotalInputCount             int             // to tell us how many input were there
	CurrentActiveInputIndex     int             // to tell us which input should be shown as highlighted/focus and be updated
}

// ------------------------------------
//
//	Holds the output/processing state for the add-new-worktree git operation pop-up
//
// ------------------------------------
type WorktreeAddNewWorktreeOutputPopUpModel struct {
	AddNewWorktreeOutputViewport viewport.Model // to log out the output from git operation
	Spinner                      spinner.Model  // spinner for showing processing state
	IsProcessing                 atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                     atomic.Bool    // indicate if the add new worktree exitcode is not 0 (meaning have error)
	ProcessSuccess               atomic.Bool    // has the process successfully executed
	IsCancelled                  atomic.Bool    // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the add new worktree operation
	CancelFunc context.CancelFunc
}

// ------------------------------------
//
//	Holds the input state for the lock-reason pop-up of a worktree
//
// ------------------------------------
type WorktreeLockReasonInputPopUpModel struct {
	WorktreePath                string          // path of the worktree to be locked
	WorktreeLockReasonTextInput textinput.Model // input for the lock reason
}
