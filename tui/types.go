package tui

import (
	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/textarea"
	"github.com/charmbracelet/bubbles/v2/textinput"
	"github.com/charmbracelet/bubbles/v2/viewport"
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

type CommitPopUpModel struct {
	MessageTextInput         textinput.Model // input index 1
	DescriptionTextAreaInput textarea.Model  // input index 2
	TotalInputCount          int             // to tell us how many input were there
	CurrentActiveInputIndex  int             // to tell us which input should be shown as highlighted/focus and be updated
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

// tea msg
type GitUpdateMsg string
