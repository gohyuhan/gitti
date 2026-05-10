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

type DetailPanelStateAndLayoutUpdateEventDataInterface struct {
	ContentLine              string
	ContentLine2             string
	OgLineDiff1              []string
	OgLineDiff2              []string
	SetForDetailComponentTwo bool
}

type MergeResultEventDataInterface struct {
	Result  []string
	Success bool
}

type GitSwitchBranchResultEventDataInterface struct {
	Result  []string
	Success bool
}

type GitDeleteBranchResultEventDataInterface struct {
	Result  []string
	Success bool
}

type GitCreateNewBranchBasedOnRemoteResultEventDataInterface struct {
	Result  []string
	Success bool
}

type GitCreateNewBranchBasedOnRemoteInvalidEventDataInterface struct {
	RemoteName string
	BranchName string
}

type GitDeleteTagResultEventDataInterface struct {
	Result  []string
	Success bool
}

type GitPushTagResultEventDataInterface struct {
	Success bool
}

type GitFetchTagResultEventDataInterface struct {
	Success bool
}

type GitStashOperationResultEventDataInterface struct {
	Result  []string
	Success bool
}

type GitAddRemoteResultEventDataInterface struct {
	Result  []string
	Success bool
}

type GitRebaseResultEventDataInterface struct {
	Success bool
}

type GitPushResultEventDataInterface struct {
	Success bool
}

type GitCommitResultEventDataInterface struct {
	Success bool
}

type GitAmendCommitResultEventDataInterface struct {
	Success bool
}

type InteractiveRebaseFixupSquashResultEventDataInterface struct {
	Result  []string
	Success bool
}

type InteractiveRebaseFetchCommitInfoListEventDataInterface struct {
	PopUpModel  string
	CommitInfos []git.CommitInfo
	ListItems   []list.Item
}
