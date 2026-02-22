package types

import (
	"context"
	"sync/atomic"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/logging"
)

type GittiModel struct {
	GittiLogger                                      *logging.GittiLogging
	DaemonUpdateChannel                              chan string // this is use to signal some trigger that will be handle by the daemon but we want it to be trigger manually also by user intention
	IsRenderInit                                     atomic.Bool // to indicate if the render has been initialized, this will be check by function that run once only after the screen is rendered
	UserSetEditor                                    string
	TuiUpdateChannel                                 chan string
	CurrentSelectedComponent                         string
	CurrentSelectedComponentIndex                    int
	CurrentLocalBranchOrTagOrRemoteComponentShowing  string
	CurrentCommitLogOrRefLogComponentShowing         string
	TotalComponentCount                              int
	RepoPath                                         string
	RepoName                                         string
	CheckOutBranch                                   string
	RemoteSyncLocalState                             string
	RemoteSyncRemoteState                            string
	CurrentGitRepoStatus                             string
	BranchUpStream                                   string
	TrackedUpstreamOrBranchIcon                      string
	Width                                            int
	Height                                           int
	WindowLeftPanelWidth                             int // this is the left part of the window
	WindowLeftPanelRatio                             float64
	DetailComponentPanelWidth                        int // this is the right part of the window, will always be for detail component panel only
	WindowCoreContentHeight                          int // this is the height of the part where key binding panel is not included
	DetailComponentPanelHeight                       int
	LocalBranchesComponentPanelHeight                int
	TagsComponentPanelHeight                         int
	RemoteComponentPanelHeight                       int
	ModifiedFilesComponentPanelHeight                int
	CommitLogComponentPanelHeight                    int
	StashComponentPanelHeight                        int
	LogComponentPanelHeight                          int
	CurrentRepoBranchesInfoList                      list.Model
	CurrentRepoTagInfoList                           list.Model
	CurrentRepoModifiedFilesInfoList                 list.Model
	CurrentRepoCommitLogInfoList                     list.Model
	CurrentRepoStashInfoList                         list.Model
	CurrentRepoRemoteInfoList                        list.Model
	DetailPanelParentComponent                       string // this is to store the parent component that cause a move into the detail panel component, so that we can return back to the correct one
	DetailPanelViewport                              viewport.Model
	DetailPanelViewportOffset                        int
	DetailPanelTwoViewport                           viewport.Model
	DetailPanelTwoViewportOffset                     int
	ShowDetailPanelTwo                               atomic.Bool
	DetailComponentPanelLayout                       string
	CurrentLogComponentViewport                      viewport.Model
	ListNavigationIndexPosition                      GittiComponentsCurrentListNavigationIndexPosition
	ShowPopUp                                        atomic.Bool
	PopUpType                                        string
	PopUpModel                                       interface{}
	IsTyping                                         atomic.Bool
	GitOperations                                    *api.GitOperations
	GlobalKeyBindingKeyMapLargestLen                 int                // this was use for global key binding pop up styling, we save it once so we don't have to recompute
	LocalBranchComponentKeyBindingKeyMapLargestLen   int                // this was use for local branch component key binding pop up styling, we save it once so we don't have to recompute
	TagComponentKeyBindingKeyMapLargestLen           int                // this was use for tag component key binding pop up styling, we save it once so we don't have to recompute
	RemoteComponentKeyBindingKeyMapLargestLen        int                // this was use for remote component key binding pop up styling, we save it once so we don't have to recompute
	ModifiedFilesComponentKeyBindingKeyMapLargestLen int                // this was use for modified files component key binding pop up styling, we save it once so we don't have to recompute
	CommitLogComponentKeyBindingKeyMapLargestLen     int                // this was use for commit log component key binding pop up styling, we save it once so we don't have to recompute
	StashComponentKeyBindingKeyMapLargestLen         int                // this was use for stash component key binding pop up styling, we save it once so we don't have to recompute
	LogComponentKeyBindingKeyMapLargestLen           int                // this was use for log component key binding pop up styling, we save it once so we don't have to recompute
	DetailComponentKeyBindingKeyMapLargestLen        int                // this was use for detail component key binding pop up styling, we save it once so we don't have to recompute
	DetailComponentPanelInfoFetchCancelFunc          context.CancelFunc // this was to cancel the fetch detail oepration
	IsDetailComponentPanelInfoFetchProcessing        atomic.Bool
	IsLineEditingState                               atomic.Bool
	LineEditingIndexPositionAndInfo                  GittiLineEditingIndexPositionAndInfo
	LineEditingIndexCursorViewport                   viewport.Model
	LineEditingIndexCursorTwoViewport                viewport.Model
	DetailPanelViewportOGStringArray                 []string
	DetailPanelTwoViewportOGStringArray              []string
	CherryPickedCommitInfo                           CherryPickedCommitInfo
	GitCommitRequireSigning                          bool
	GitTagRequireSigning                             bool
	GitPushRequireSigning                            bool
}

// ---------------------------------
//
// to record the cherry picked commit info like the sequence counter and the map of cherry picked commit
//
// ---------------------------------
type CherryPickedCommitInfo struct {
	LatestSequenceCounter int // this counter will only increase or reinit to 0, this help us to sort the Map when showing in the UI (doesn't mean will apply to cherry pick in the order)
	CherryPickedCommitMap map[string]git.CherryPickedCommitLog
}

// ---------------------------------
//
// to record the current navigation index position
//
// ---------------------------------
type GittiComponentsCurrentListNavigationIndexPosition struct {
	LocalBranchComponent   int
	ModifiedFilesComponent int
	CommitLogComponent     int
	StashComponent         int
	TagComponent           int
	RemoteComponent        int
}

// ---------------------------------
//
// to record the current navigation on detail component for line editing purpose
//
// ---------------------------------
type GittiLineEditingIndexPositionAndInfo struct {
	DetailPanelViewportIndexPosition         int    // for staged or unstaged panel (depends if there are both staged and unstaged changes)
	DetailPanelTwoViewportIndexPosition      int    // for unstaged panel
	DetailPanelViewportStageType             string // to tell it was currently showing staged or unstaged changes
	DetailPanelTwoViewportStageType          string // to tell it was currently showing staged or unstaged changes
	DetailPanelViewportActualCurrentIndex    int    // this is the actual current index of the viewport (including overflow)
	DetailPanelTwoViewportActualCurrentIndex int    // this is the actual current index of the viewport (including overflow)
	DetailPanelViewportOverflowIndexCount    int    // we need to record this as well as to calculate the actual current index due to we have added extra info at the top of the viewport that we added ourselves which is not from git
	DetailPanelTwoViewportOverflowIndexCount int    // we need to record this as well as to calculate the actual current index due to we have added extra info at the top of the viewport that we added ourselves which is not from git
}

// ---------------------------------
//
// # A bubbletea message to indicate that the editor has quit or close (apply only for terminal editor, external GUI like vscode/zed/cursor etc will not need this)
//
// ---------------------------------
type EditorFinishedMsg struct {
	Err error
}

type GitOperationRequiredSigningFinishedMsg struct {
	GitOperationOpsTypeForLogging string // for logging purpose
	Err                           error
}
