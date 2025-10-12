package ui

import (
	"gitti/types"

	"github.com/charmbracelet/bubbles/list"
)

type GittiModel struct {
	RepoPath                          string
	Width                             int
	Height                            int
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
