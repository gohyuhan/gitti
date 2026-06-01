package types

import (
	"charm.land/bubbles/v2/list"
	"github.com/gohyuhan/gitti/api/git"
)

// ---------------------------------
//
// tea msg (interface, include custom data structure)
//
// ---------------------------------
type GittiTuiUpdateMsg struct {
	Event string
	Data  interface{}
}

type DetailPanelStateAndLayoutUpdateEventDataStructure struct {
	ContentLine              string
	ContentLine2             string
	OgLineDiff1              []string
	OgLineDiff2              []string
	SetForDetailComponentTwo bool
}

type MergeResultEventDataStructure struct {
	Result  []string
	Success bool
}

type GitSwitchBranchResultEventDataStructure struct {
	Result  []string
	Success bool
}

type GitDeleteBranchResultEventDataStructure struct {
	Result  []string
	Success bool
}

type GitCreateNewBranchBasedOnRemoteResultEventDataStructure struct {
	Result  []string
	Success bool
}

type GitCreateNewBranchBasedOnRemoteInvalidEventDataStructure struct {
	RemoteName string
	BranchName string
}

type GitDeleteTagResultEventDataStructure struct {
	Result  []string
	Success bool
}

type GitPushTagResultEventDataStructure struct {
	Success bool
}

type GitFetchTagResultEventDataStructure struct {
	Success bool
}

type GitStashOperationResultEventDataStructure struct {
	Result  []string
	Success bool
}

type GitAddRemoteResultEventDataStructure struct {
	Result  []string
	Success bool
}

type GitRebaseResultEventDataStructure struct {
	Success bool
}

type GitPushResultEventDataStructure struct {
	Success bool
}

type GitCommitResultEventDataStructure struct {
	Success bool
}

type GitAmendCommitResultEventDataStructure struct {
	Success bool
}

type GitPullResultEventDataStructure struct {
	Success bool
}

type InteractiveRebaseFixupSquashResultEventDataStructure struct {
	Result  []string
	Success bool
}

type InteractiveRebaseRewordResultEventDataStructure struct {
	Result  []string
	Success bool
}

type InteractiveRebaseDropResultEventDataStructure struct {
	Result  []string
	Success bool
}

type InteractiveRebaseFetchCommitInfoListEventDataStructure struct {
	PopUpModel  string
	CommitInfos []git.CommitInfo
	ListItems   []list.Item
}

type WorktreeNewWorktreeResultEventDataStructure struct {
	Result  []string
	Success bool
}
