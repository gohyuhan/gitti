package tui

import (
	"gitti/api"
	"gitti/types"

	"github.com/charmbracelet/bubbles/list"
)

type GittiModel struct {
	CurrentTab                        string
	CurrentSelectedContainer          string
	RepoPath                          string
	Width                             int
	Height                            int
	HomeTabLeftPanelWidth             int
	HomeTabFileDiffPanelWidth         int
	HomeTabCoreContentHeight          int
	HomeTabFileDiffPanelHeight        int
	HomeTabLocalBranchesPanelHeight   int
	HomeTabChangedFilesPanelHeight    int
	AllRepoBranches                   map[string]types.BranchesInfo
	CurrentRepoBranchesInfo           list.Model
	CurrentCheckedOutBranch           string
	CurrentSelectedFiles              string
	CurrentSelectedFilesIndexPosition int
	AllChangedFiles                   map[string]types.FilesInfo
	RemoteOrigin                      string
	UserName                          string
	UserEmail                         string
	GitWorkerDaemon                   api.GittiDaemonWorker
	NavigationIndexPosition           GittiComponentsCurrentNavigationIndexPosition
}

// to record the current navigation index position
type GittiComponentsCurrentNavigationIndexPosition struct {
	LocalBranchComponent  int
	FilesChangesComponent int
}

// for list component
type itemDelegate struct{}
type item string

type GitUpdateMsg types.GitInfo
