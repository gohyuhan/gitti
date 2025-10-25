package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
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
	CurrentRepoBranchesInfo               list.Model
	CurrentRepoModifiedFilesInfo          list.Model
	CurrentSelectedFileDiffViewport       viewport.Model
	CurrentSelectedFileDiffViewportOffset int
	NavigationIndexPosition               GittiComponentsCurrentNavigationIndexPosition
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
