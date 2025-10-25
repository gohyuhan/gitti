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

// for list component
type itemStringDelegate struct{}
type itemString string

type GitUpdateMsg string
