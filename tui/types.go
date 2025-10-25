package tui

import (
	"github.com/charmbracelet/bubbles/list"
)

type GittiModel struct {
	CurrentSelectedContainer        string
	RepoPath                        string
	Width                           int
	Height                          int
	HomeTabLeftPanelWidth           int
	HomeTabFileDiffPanelWidth       int
	HomeTabCoreContentHeight        int
	HomeTabFileDiffPanelHeight      int
	HomeTabLocalBranchesPanelHeight int
	HomeTabChangedFilesPanelHeight  int
	CurrentRepoBranchesInfo         list.Model
	NavigationIndexPosition         GittiComponentsCurrentNavigationIndexPosition
}

// to record the current navigation index position
type GittiComponentsCurrentNavigationIndexPosition struct {
	LocalBranchComponent  int
	FilesChangesComponent int
}

// for list component
type itemDelegate struct{}
type item string

type GitUpdateMsg string
