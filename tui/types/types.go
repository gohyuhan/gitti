package types

import (
	"context"
	"sync/atomic"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	"github.com/gohyuhan/gitti/api"
)

type GittiModel struct {
	IsRenderInit                              atomic.Bool // to indicate if the render has been initialized, this will be check by function that run once only after the screen is rendered
	TuiUpdateChannel                          chan string
	CurrentSelectedComponent                  string
	CurrentSelectedComponentIndex             int
	TotalComponentCount                       int
	RepoPath                                  string
	RepoName                                  string
	CheckOutBranch                            string
	RemoteSyncLocalState                      string
	RemoteSyncRemoteState                     string
	BranchUpStream                            string
	TrackedUpstreamOrBranchIcon               string
	Width                                     int
	Height                                    int
	WindowLeftPanelWidth                      int // this is the left part of the window
	DetailComponentPanelWidth                 int // this is the right part of the window, will always be for detail component panel only
	WindowCoreContentHeight                   int // this is the height of the part where key binding panel is not included
	DetailComponentPanelHeight                int
	LocalBranchesComponentPanelHeight         int
	ModifiedFilesComponentPanelHeight         int
	CommitLogComponentPanelHeight             int
	StashComponentPanelHeight                 int
	CurrentRepoBranchesInfoList               list.Model
	CurrentRepoModifiedFilesInfoList          list.Model
	CurrentRepoCommitLogInfoList              list.Model
	CurrentRepoStashInfoList                  list.Model
	DetailPanelParentComponent                string // this is to store the parent component that cause a move into the detail panel component, so that we can return back to the correct one
	DetailPanelViewport                       viewport.Model
	DetailPanelViewportOffset                 int
	DetailPanelTwoViewport                    viewport.Model
	DetailPanelTwoViewportOffset              int
	ShowDetailPanelTwo                        atomic.Bool
	DetailComponentPanelLayout                string
	ListNavigationIndexPosition               GittiComponentsCurrentListNavigationIndexPosition
	ShowPopUp                                 atomic.Bool
	PopUpType                                 string
	PopUpModel                                interface{}
	IsTyping                                  atomic.Bool
	GitOperations                             *api.GitOperations
	GlobalKeyBindingKeyMapLargestLen          int                // this was use for global key binding pop up styling, we save it once so we don't have to recompute
	DetailComponentPanelInfoFetchCancelFunc   context.CancelFunc // this was to cancel the fetch detail oepration
	IsDetailComponentPanelInfoFetchProcessing atomic.Bool
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
}
