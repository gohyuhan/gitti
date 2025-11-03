package tui

import (
	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/spinner"
	"github.com/charmbracelet/bubbles/v2/textarea"
	"github.com/charmbracelet/bubbles/v2/textinput"
	"github.com/charmbracelet/bubbles/v2/viewport"
	"github.com/google/uuid"
)

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
	ShowPopUp                             bool
	PopUpType                             string
	PopUpModel                            interface{}
	IsTyping                              bool
}

type GitCommitPopUpModel struct {
	MessageTextInput         textinput.Model // input index 1
	DescriptionTextAreaInput textarea.Model  // input index 2
	TotalInputCount          int             // to tell us how many input were there
	CurrentActiveInputIndex  int             // to tell us which input should be shown as highlighted/focus and be updated
	GitCommitOutputViewport  viewport.Model  // to log out the output from git operation
	Spinner                  spinner.Model   // spinner for showing processing state
	IsProcessing             bool            // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                 bool            // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess           bool            // has the process sucessfuly executed
	IsCancelled              bool            // flag to indicate if the operation was cancelled by user
	// SessionID is a unique UUID for each popup instance to prevent
	// stale goroutines from affecting new popups
	SessionID uuid.UUID
}

type AddRemotePromptPopUpModel struct {
	RemoteNameTextInput     textinput.Model // input index 1
	RemoteUrlTextInput      textinput.Model // input index 2
	TotalInputCount         int             // to tell us how many input were there
	CurrentActiveInputIndex int             // to tell us which input should be shown as highlighted/focus and be updated
	AddRemoteOutputViewport viewport.Model  // to log out the output from git operation
	IsProcessing            bool            // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                bool            // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess          bool            // has the process sucessfuly executed
	NoInitialRemote         bool            // indicate if this repo has no remote yet or user just wanted to add more remote
	IsCancelled             bool            // flag to indicate if the operation was cancelled by user
	// SessionID is a unique UUID for each popup instance to prevent
	// stale goroutines from affecting new popups
	SessionID uuid.UUID
}

type GitRemotePushPopUpModel struct {
	GitRemotePushOutputViewport viewport.Model // to log out the output from git operation
	Spinner                     spinner.Model  // spinner for showing processing state
	IsProcessing                bool           // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                    bool           // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess              bool           // has the process sucessfuly executed
	IsCancelled                 bool           // flag to indicate if the operation was cancelled by user
	// SessionID is a unique UUID for each popup instance to prevent
	// stale goroutines from affecting new popups
	SessionID uuid.UUID
}

type ChooseRemotePopUpModel struct {
	RemoteList list.Model
}

// to record the current navigation index position
type GittiComponentsCurrentNavigationIndexPosition struct {
	LocalBranchComponent   int
	ModifiedFilesComponent int
}

// for list component of git branch
type gitBranchItemDelegate struct{}
type gitBranchItem struct {
	BranchName   string
	IsCheckedOut bool
}

func (i gitBranchItem) FilterValue() string {
	return i.BranchName
}

// for list component of git modified files
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

// for list component of git remote
type gitRemoteItemDelegate struct{}
type gitRemoteItem struct {
	Name string
	Url  string
}

func (i gitRemoteItem) FilterValue() string {
	return i.Name
}

// tea msg
type GitUpdateMsg string
