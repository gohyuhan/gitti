package tui

import (
	"context"
	"sync/atomic"

	"github.com/gohyuhan/gitti/api"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
)

// ---------------------------------
//
// # Main Model for TUI
//
// ---------------------------------
type GittiModel struct {
	IsRenderInit                              atomic.Bool // to indicate if the render has been initialized, this will be check by function that run once only after the screen is rendered
	TuiUpdateChannel                          chan string
	CurrentSelectedComponent                  string
	CurrentSelectedComponentIndex             int
	TotalComponentCount                       int
	RepoPath                                  string
	RepoName                                  string
	CheckOutBranch                            string
	RemoteSyncStateLineString                 string
	BranchUpStream                            string
	TrackedUpstreamOrBranchIcon               string
	Width                                     int
	Height                                    int
	WindowLeftPanelWidth                      int // this is the left part of the window
	DetailComponentPanelWidth                 int // this is the right part of the window, will always be for detail component panel only
	WindowCoreContentHeight                   int // this is the height of the part where key binding panel is not included
	DetailComponentPanelHeight                int
	LocalBranchesComponentPanelHeight         int
	ModifiedFilesComponentPanelHeight         int
	StashComponentPanelHeight                 int
	CurrentRepoBranchesInfoList               list.Model
	CurrentRepoModifiedFilesInfoList          list.Model
	CurrentRepoStashInfoList                  list.Model
	DetailPanelParentComponent                string // this is to store the parent component that cause a move into the detail panel component, so that we can return back to the correct one
	DetailPanelViewport                       viewport.Model
	DetailPanelViewportOffset                 int
	ListNavigationIndexPosition               GittiComponentsCurrentListNavigationIndexPosition
	ShowPopUp                                 atomic.Bool
	PopUpType                                 string
	PopUpModel                                interface{}
	IsTyping                                  atomic.Bool
	GitOperations                             *api.GitOperations
	GlobalKeyBindingKeyMapLargestLen          int                // this was use for global key binding pop up styling, we save it once so we don't have to recompute
	DetailComponentPanelInfoFetchCancelFunc   context.CancelFunc // this was to cancel the fetch detail oepration
	IsDetailComponentPanelInfoFetchProcessing atomic.Bool
}

// ---------------------------------
//
// # A pop up helper for global keybinding
//
// ---------------------------------
type GlobalKeyBindingPopUpModel struct {
	GlobalKeyBindingViewport viewport.Model
}

// ---------------------------------
//
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

// ---------------------------------
//
// # For add Remote prompt pop up
//
// ---------------------------------
type AddRemotePromptPopUpModel struct {
	RemoteNameTextInput     textinput.Model // input index 1
	RemoteUrlTextInput      textinput.Model // input index 2
	TotalInputCount         int             // to tell us how many input were there
	CurrentActiveInputIndex int             // to tell us which input should be shown as highlighted/focus and be updated
	AddRemoteOutputViewport viewport.Model  // to log out the output from git operation
	IsProcessing            atomic.Bool     // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                atomic.Bool     // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess          atomic.Bool     // has the process sucessfuly executed
	NoInitialRemote         bool            // indicate if this repo has no remote yet or user just wanted to add more remote
	IsCancelled             atomic.Bool     // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git remote add operation
	CancelFunc context.CancelFunc
}

// ---------------------------------
//
// # For Remote push process pop up
//
// ---------------------------------
type GitRemotePushPopUpModel struct {
	GitRemotePushOutputViewport viewport.Model // to log out the output from git operation
	Spinner                     spinner.Model  // spinner for showing processing state
	IsProcessing                atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                    atomic.Bool    // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess              atomic.Bool    // has the process sucessfuly executed
	IsCancelled                 atomic.Bool    // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git push operation
	CancelFunc context.CancelFunc
}

// ---------------------------------
//
// user choose how do they want to push the commit, push /  push --force / push --force-with-lease
//
// ---------------------------------
type ChoosePushTypePopUpModel struct {
	PushOptionList list.Model
	RemoteName     string
}

// ---------------------------------
//
// choose a remote to push to
//
// ---------------------------------
type ChooseRemotePopUpModel struct {
	RemoteList list.Model
}

// ---------------------------------
//
// create a new branch and remain on current branch
//
// ---------------------------------
type CreateNewBranchPopUpModel struct {
	NewBranchNameInput textinput.Model
	CreateType         string
}

// ---------------------------------
//
// choose on how to create the new branch, just create or create and move changes
//
// ---------------------------------
type ChooseNewBranchTypeOptionPopUpModel struct {
	NewBranchTypeOptionList list.Model
}

// ---------------------------------
//
// choose a switch type when switching branch
//
// ---------------------------------
type ChooseSwitchBranchTypePopUpModel struct {
	SwitchTypeOptionList list.Model
	BranchName           string
}

// ---------------------------------
//
// # A pop up to show branch switch result
//
// ---------------------------------
type SwitchBranchOutputPopUpModel struct {
	BranchName                 string // the branch name of the branch it was switching to
	SwitchType                 string
	SwitchBranchOutputViewport viewport.Model // to log out the output from git operation
	Spinner                    spinner.Model  // spinner for showing processing state
	IsProcessing               atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                   atomic.Bool    // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess             atomic.Bool    // has the process sucessfuly executed
}

// ---------------------------------
//
// choose a pull type, git pull, git pull rebase or git pull merge
//
// ---------------------------------
type ChooseGitPullTypePopUpModel struct {
	PullTypeOptionList list.Model
}

// ---------------------------------
//
// # A pop up to show git pull result
//
// ---------------------------------
type GitPullOutputPopUpModel struct {
	PullType              string
	GitPullOutputViewport viewport.Model // to log out the output from git operation
	Spinner               spinner.Model  // spinner for showing processing state
	IsProcessing          atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError              atomic.Bool    // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess        atomic.Bool    // has the process sucessfuly executed
	IsCancelled           atomic.Bool    // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git pull operation
	CancelFunc context.CancelFunc
}

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
// # A pop up to prompt for git discard type for file,
//
//	this will only be available when there are both changes in stage and unstage (index and worktree)
//
// ---------------------------------
type GitDiscardTypeOptionPopUpModel struct {
	DiscardTypeOptionList list.Model
	FilePathName          string
}

// ---------------------------------
//
// # To prompt user for confirmation
//
// ---------------------------------
type GitDiscardConfirmPromptPopUpModel struct {
	DiscardType  string
	FilePathName string
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

// ---------------------------------
//
// for resolve conflict option pop up
//
// ---------------------------------
type GitResolveConflictOptionPopUpModel struct {
	FilePathName              string
	ResolveConflictOptionList list.Model
}

// ---------------------------------
//
// to record the current navigation index position
//
// ---------------------------------
type GittiComponentsCurrentListNavigationIndexPosition struct {
	LocalBranchComponent   int
	ModifiedFilesComponent int
	StashComponent         int
}

// ---------------------------------
//
// for list component of git branch
//
// ---------------------------------
type (
	gitBranchItemDelegate struct{}
	gitBranchItem         struct {
		BranchName   string
		IsCheckedOut bool
	}
)

func (i gitBranchItem) FilterValue() string {
	return i.BranchName
}

// ---------------------------------
//
// for list component of git modified files
//
// ---------------------------------
type (
	gitModifiedFilesItemDelegate struct{}
	gitModifiedFilesItem         struct {
		FilePathname string
		IndexState   string
		WorkTree     string
		HasConflict  bool
	}
)

func (i gitModifiedFilesItem) FilterValue() string {
	return i.FilePathname
}

// ---------------------------------
//
// for list component of git stashed files
//
// ---------------------------------
type (
	gitStashItemDelegate struct{}
	gitStashItem         struct {
		Id      string
		Message string
	}
)

func (i gitStashItem) FilterValue() string {
	return i.Message
}

// ---------------------------------
//
// for list component of git remote
//
// ---------------------------------
type (
	gitRemoteItemDelegate struct{}
	gitRemoteItem         struct {
		Name string
		Url  string
	}
)

func (i gitRemoteItem) FilterValue() string {
	return i.Name
}

// ---------------------------------
//
// for push selection option
//
// ---------------------------------
type (
	gitPushOptionDelegate struct{}
	gitPushOptionItem     struct {
		Name     string
		Info     string
		pushType string
	}
)

func (i gitPushOptionItem) FilterValue() string {
	return i.Name
}

// ---------------------------------
//
// for new branch option selection option
//
// ---------------------------------
type (
	gitNewBranchTypeOptionDelegate struct{}
	gitNewBranchTypeOptionItem     struct {
		Name          string
		Info          string
		newBranchType string
	}
)

func (i gitNewBranchTypeOptionItem) FilterValue() string {
	return i.Name
}

// ---------------------------------
//
// for switch branch option selection option
//
// ---------------------------------
type (
	gitSwitchBranchTypeOptionDelegate struct{}
	gitSwitchBranchTypeOptionItem     struct {
		Name             string
		Info             string
		switchBranchType string
	}
)

func (i gitSwitchBranchTypeOptionItem) FilterValue() string {
	return i.Name
}

// ---------------------------------
//
// for pull option selection option
//
// ---------------------------------
type (
	gitPullTypeOptionDelegate struct{}
	gitPullTypeOptionItem     struct {
		Name     string
		Info     string
		PullType string
	}
)

func (i gitPullTypeOptionItem) FilterValue() string {
	return i.Name
}

// ---------------------------------
//
// for discard option selection option
//
// ---------------------------------
type (
	gitDiscardTypeOptionDelegate struct{}
	gitDiscardTypeOptionItem     struct {
		Name        string
		Info        string
		DiscardType string
	}
)

func (i gitDiscardTypeOptionItem) FilterValue() string {
	return i.Name
}

// ---------------------------------
//
// for resolve conflict option selection option
//
// ---------------------------------
type (
	gitResolveConflictOptionDelegate struct{}
	gitResolveConflictOptionItem     struct {
		Name        string
		Info        string
		ResolveType string
	}
)

func (i gitResolveConflictOptionItem) FilterValue() string {
	return i.Name
}

// ---------------------------------
//
// tea msg
//
// ---------------------------------
type GitUpdateMsg string
