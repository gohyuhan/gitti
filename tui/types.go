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
	AllChangedFiles                   map[string]string
	RemoteOrigin                      string
	UserName                          string
	UserEmail                         string
	GitWorkerDaemon                   api.GittiDaemonWorker
}

// for list component
type itemDelegate struct{}
type item string

type GitUpdateMsg types.GitInfo
