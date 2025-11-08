package tui

import (
	"sync/atomic"

	"gitti/api"

	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/spinner"
	"github.com/charmbracelet/bubbles/v2/textarea"
	"github.com/charmbracelet/bubbles/v2/textinput"
	"github.com/charmbracelet/bubbles/v2/viewport"
	"github.com/google/uuid"
)

// ---------------------------------
//
// # Main Model for TUI
//
// ---------------------------------
type GittiModel struct {
	CurrentSelectedContainer              string
	RepoPath                              string
	Width                                 int
	Height                                int
	HomeTabLeftPanelWidth                 int
	HomeTabFileDiffPanelWidth             int
	HomeTabCoreContentHeight              int
	HomeTabFileDiffPanelHeight            int
	HomeTabLocalBranchesPanelHeight       int
	HomeTabChangedFilesPanelHeight        int
	CurrentRepoBranchesInfoList           list.Model
	CurrentRepoModifiedFilesInfoList      list.Model
	CurrentSelectedFileDiffViewport       viewport.Model
	CurrentSelectedFileDiffViewportOffset int
	NavigationIndexPosition               GittiComponentsCurrentNavigationIndexPosition
	ShowPopUp                             atomic.Bool
	PopUpType                             string
	PopUpModel                            interface{}
	IsTyping                              atomic.Bool
	GitState                              *api.GitState
}

// ---------------------------------
//
// # For git commit process pop up mdel
//
// ---------------------------------
type GitCommitPopUpModel struct {
	MessageTextInput         textinput.Model // input index 1
	DescriptionTextAreaInput textarea.Model  // input index 2
	TotalInputCount          int             // to tell us how many input were there
	CurrentActiveInputIndex  int             // to tell us which input should be shown as highlighted/focus and be updated
	GitCommitOutputViewport  viewport.Model  // to log out the output from git operation
	Spinner                  spinner.Model   // spinner for showing processing state
	IsProcessing             atomic.Bool     // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                 atomic.Bool     // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess           atomic.Bool     // has the process sucessfuly executed
	IsCancelled              atomic.Bool     // flag to indicate if the operation was cancelled by user
	// SessionID is a unique UUID for each popup instance to prevent
	// stale goroutines from affecting new popups
	SessionID uuid.UUID
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
	// SessionID is a unique UUID for each popup instance to prevent
	// stale goroutines from affecting new popups
	SessionID uuid.UUID
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
	// SessionID is a unique UUID for each popup instance to prevent
	// stale goroutines from affecting new popups
	SessionID uuid.UUID
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
// to record the current navigation index position
//
// ---------------------------------
type GittiComponentsCurrentNavigationIndexPosition struct {
	LocalBranchComponent   int
	ModifiedFilesComponent int
}

// ---------------------------------
//
// for list component of git branch
//
// ---------------------------------
type gitBranchItemDelegate struct{}
type gitBranchItem struct {
	BranchName   string
	IsCheckedOut bool
}

func (i gitBranchItem) FilterValue() string {
	return i.BranchName
}

// ---------------------------------
//
// for list component of git modified files
//
// ---------------------------------
type gitModifiedFilesItemDelegate struct{}
type gitModifiedFilesItem struct {
	FileName         string
	IndexState       string
	WorkTree         string
	SelectedForStage bool
}

func (i gitModifiedFilesItem) FilterValue() string {
	return i.FileName
}

// ---------------------------------
//
// for list component of git remote
//
// ---------------------------------
type gitRemoteItemDelegate struct{}
type gitRemoteItem struct {
	Name string
	Url  string
}

func (i gitRemoteItem) FilterValue() string {
	return i.Name
}

// ---------------------------------
//
// for push selection option
//
// ---------------------------------
type gitPushOptionDelegate struct{}
type gitPushOptionItem struct {
	Name     string
	Info     string
	pushType string
}

func (i gitPushOptionItem) FilterValue() string {
	return i.Name
}

// ---------------------------------
//
// for new branch option selection option
//
// ---------------------------------
type gitNewBranchTypeOptionDelegate struct{}
type gitNewBranchTypeOptionItem struct {
	Name          string
	Info          string
	newBranchType string
}

func (i gitNewBranchTypeOptionItem) FilterValue() string {
	return i.Name
}

// ---------------------------------
//
// for new branch option selection option
//
// ---------------------------------
type gitSwitchBranchTypeOptionDelegate struct{}
type gitSwitchBranchTypeOptionItem struct {
	Name             string
	Info             string
	switchBranchType string
}

func (i gitSwitchBranchTypeOptionItem) FilterValue() string {
	return i.Name
}

// ---------------------------------
//
// tea msg
//
// ---------------------------------
type GitUpdateMsg string
