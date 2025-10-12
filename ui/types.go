package ui

import (
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
	AllRepoBranches                   []types.BranchesInfo
	CurrentRepoBranchesInfo           list.Model
	CurrentCheckedOutBranch           string
	CurrentSelectedFiles              string
	CurrentSelectedFilesIndexPosition int
	AllChangedFiles                   []string
	RemoteOrigin                      string
	UserName                          string
	UserEmail                         string
}

type itemDelegate struct{}
type item string
