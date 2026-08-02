package initialize

import (
	branchComponent "github.com/gohyuhan/gitti/tui/component/branch"
	commitlogComponent "github.com/gohyuhan/gitti/tui/component/commitlog"
	filesComponent "github.com/gohyuhan/gitti/tui/component/files"
	reflogComponent "github.com/gohyuhan/gitti/tui/component/reflog"
	remoteComponent "github.com/gohyuhan/gitti/tui/component/remote"
	stashComponent "github.com/gohyuhan/gitti/tui/component/stash"
	tagComponent "github.com/gohyuhan/gitti/tui/component/tag"
	worktreeComponent "github.com/gohyuhan/gitti/tui/component/worktree"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"

	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/settings"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Build a fresh GittiModel for app startup: creates the detail/log/line-editing
//	viewports, initialises every list component, wires the channels, git
//	operations and logger, and checks whether git signing is required for
//	commit/tag/push. Returns the newly allocated model.
//
// ------------------------------------
func InitGittiModel(tuiUpdateChannel chan interface{}, repoPath string, repoName string, gitOperations *api.GitOperations, gittiLogger *logging.GittiLogging, daemonUpdateChannel chan string, gitUpdateChannel chan string) *types.GittiModel {
	vp := viewport.New()
	vp.SoftWrap = false
	vp.MouseWheelEnabled = true
	vp.SetHorizontalStep(1)
	vp.MouseWheelDelta = 1

	vpTwo := viewport.New()
	vpTwo.SoftWrap = false
	vpTwo.MouseWheelEnabled = true
	vpTwo.SetHorizontalStep(1)
	vpTwo.MouseWheelDelta = 1

	logVp := viewport.New()
	logVp.SoftWrap = false
	logVp.MouseWheelEnabled = false
	logVp.SetHorizontalStep(1)
	logVp.MouseWheelDelta = 1

	lineEditingIndexCursorVp := viewport.New()
	lineEditingIndexCursorVp.SoftWrap = false
	lineEditingIndexCursorVp.MouseWheelEnabled = false
	lineEditingIndexCursorVp.SetHorizontalStep(0)
	lineEditingIndexCursorVp.MouseWheelDelta = 0

	lineEditingIndexCursorVpTwo := viewport.New()
	lineEditingIndexCursorVpTwo.SoftWrap = false
	lineEditingIndexCursorVpTwo.MouseWheelEnabled = false
	lineEditingIndexCursorVpTwo.SetHorizontalStep(0)
	lineEditingIndexCursorVpTwo.MouseWheelDelta = 0

	gittiModel := &types.GittiModel{
		GittiLogger:                   gittiLogger,
		DaemonUpdateChannel:           daemonUpdateChannel,
		GitUpdateChannel:              gitUpdateChannel,
		TuiUpdateChannel:              tuiUpdateChannel,
		UserSetEditor:                 settings.GITTICONFIGSETTINGS.Editor,
		CurrentSelectedComponent:      constant.ModifiedFilesComponentPanel,
		CurrentSelectedComponentIndex: 2,
		CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing: constant.SHOW_LOCAL_BRANCH,
		CurrentCommitLogOrRefLogComponentShowing:                  constant.SHOW_COMMITLOG,
		TotalComponentCount:                                       4,
		RepoPath:                                                  repoPath,
		RepoName:                                                  repoName,
		CheckOutBranch:                                            "",
		RemoteSyncLocalState:                                      "",
		RemoteSyncRemoteState:                                     "",
		CurrentGitRepoStatus:                                      "",
		BranchUpStream:                                            "",
		TrackedUpstreamOrBranchIcon:                               "",
		Width:                                                     0,
		Height:                                                    0,
		WindowLeftPanelRatio:                                      settings.GITTICONFIGSETTINGS.LeftPanelWidthRatio,
		CurrentRepoBranchesInfoList:                               list.New([]list.Item{}, branchComponent.GitBranchItemDelegate{}, 0, 0),
		CurrentRepoTagInfoList:                                    list.New([]list.Item{}, tagComponent.GitTagItemDelegate{}, 0, 0),
		CurrentRepoModifiedFilesInfoList:                          list.New([]list.Item{}, filesComponent.GitModifiedFilesItemDelegate{}, 0, 0),
		CurrentRepoCommitLogInfoList:                              list.New([]list.Item{}, commitlogComponent.GitCommitLogItemDelegate{}, 0, 0),
		CurrentRepoRefLogInfoList:                                 list.New([]list.Item{}, reflogComponent.GitRefLogItemDelegate{}, 0, 0),
		CurrentRepoStashInfoList:                                  list.New([]list.Item{}, stashComponent.GitStashItemDelegate{}, 0, 0),
		CurrentRepoRemoteInfoList:                                 list.New([]list.Item{}, remoteComponent.GitRemoteItemDelegate{}, 0, 0),
		CurrentRepoWorktreeInfoList:                               list.New([]list.Item{}, worktreeComponent.GitWorktreeItemDelegate{}, 0, 0),
		DetailPanelParentComponent:                                "",
		DetailPanelViewport:                                       vp,
		DetailPanelViewportOffset:                                 0,
		DetailPanelTwoViewport:                                    vpTwo,
		DetailPanelTwoViewportOffset:                              0,
		DetailComponentPanelLayout:                                constant.HORIZONTAL,
		CurrentLogComponentViewport:                               logVp,
		ListNavigationIndexPosition:                               types.GittiComponentsCurrentListNavigationIndexPosition{LocalBranchComponent: 0, ModifiedFilesComponent: 0, StashComponent: 0, RefLogComponent: 0, TagComponent: 0, RemoteComponent: 0, WorktreeComponent: 0},
		PopUpType:                                                 constant.NoPopUp,
		PopUpModel:                                                struct{}{},
		GitOperations:                                             gitOperations,
		GlobalKeyBindingKeyMapLargestLen:                          0,
		LocalBranchComponentKeyBindingKeyMapLargestLen:            0,
		TagComponentKeyBindingKeyMapLargestLen:                    0,
		RemoteComponentKeyBindingKeyMapLargestLen:                 0,
		ModifiedFilesComponentKeyBindingKeyMapLargestLen:          0,
		CommitLogComponentKeyBindingKeyMapLargestLen:              0,
		RefLogComponentKeyBindingKeyMapLargestLen:                 0,
		StashComponentKeyBindingKeyMapLargestLen:                  0,
		LogComponentKeyBindingKeyMapLargestLen:                    0,
		DetailComponentKeyBindingKeyMapLargestLen:                 0,
		LineEditingIndexPositionAndInfo:                           types.GittiLineEditingIndexPositionAndInfo{},
		LineEditingIndexCursorViewport:                            lineEditingIndexCursorVp,
		LineEditingIndexCursorTwoViewport:                         lineEditingIndexCursorVpTwo,
		CherryPickedCommitInfo:                                    types.CherryPickedCommitInfo{LatestSequenceCounter: 0, CherryPickedCommitMap: make(map[string]git.CherryPickedCommitLog)},
		PanelFilterQuery:                                          make(map[string]string),
	}
	gittiModel.IsRenderInit.Store(false)
	gittiModel.ShowPopUp.Store(false)
	gittiModel.IsTyping.Store(false)
	gittiModel.IsDetailComponentPanelInfoFetchProcessing.Store(false)
	gittiModel.ShowDetailPanelTwo.Store(false)
	gittiModel.IsLineEditingState.Store(false)

	commitRequireSigning, tagRequireSigning, pushRequireSigning := api.CheckSigningRequiredOperation()
	gittiModel.GitCommitRequireSigning = commitRequireSigning
	gittiModel.GitTagRequireSigning = tagRequireSigning
	gittiModel.GitPushRequireSigning = pushRequireSigning

	return gittiModel
}

// ------------------------------------
//
//	Reset an existing GittiModel in place for a worktree switch. Mutates the model
//	through its pointer (the model holds atomic fields, so it cannot be replaced by
//	a struct copy) so the running bubbletea program keeps the same reference. Every
//	field is reset to its startup default to avoid cross-worktree state leaking,
//	EXCEPT Width/Height which are preserved since the terminal size is unchanged.
//
// ------------------------------------
func ReinitGittiModel(m *types.GittiModel, repoPath string, repoName string, gitOperations *api.GitOperations) {
	vp := viewport.New()
	vp.SoftWrap = false
	vp.MouseWheelEnabled = true
	vp.SetHorizontalStep(1)
	vp.MouseWheelDelta = 1

	vpTwo := viewport.New()
	vpTwo.SoftWrap = false
	vpTwo.MouseWheelEnabled = true
	vpTwo.SetHorizontalStep(1)
	vpTwo.MouseWheelDelta = 1

	lineEditingIndexCursorVp := viewport.New()
	lineEditingIndexCursorVp.SoftWrap = false
	lineEditingIndexCursorVp.MouseWheelEnabled = false
	lineEditingIndexCursorVp.SetHorizontalStep(0)
	lineEditingIndexCursorVp.MouseWheelDelta = 0

	lineEditingIndexCursorVpTwo := viewport.New()
	lineEditingIndexCursorVpTwo.SoftWrap = false
	lineEditingIndexCursorVpTwo.MouseWheelEnabled = false
	lineEditingIndexCursorVpTwo.SetHorizontalStep(0)
	lineEditingIndexCursorVpTwo.MouseWheelDelta = 0

	// reinit all state on worktree switch to avoid cross-worktree contamination.
	// Width/Height are intentionally preserved (terminal size is unchanged).
	m.UserSetEditor = settings.GITTICONFIGSETTINGS.Editor
	m.CurrentSelectedComponent = constant.ModifiedFilesComponentPanel
	m.CurrentSelectedComponentIndex = 2
	m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing = constant.SHOW_LOCAL_BRANCH
	m.CurrentCommitLogOrRefLogComponentShowing = constant.SHOW_COMMITLOG
	m.TotalComponentCount = 4
	m.RepoPath = repoPath
	m.RepoName = repoName
	m.CheckOutBranch = ""
	m.RemoteSyncLocalState = ""
	m.RemoteSyncRemoteState = ""
	m.CurrentGitRepoStatus = ""
	m.BranchUpStream = ""
	m.TrackedUpstreamOrBranchIcon = ""
	m.WindowLeftPanelRatio = settings.GITTICONFIGSETTINGS.LeftPanelWidthRatio
	m.CurrentRepoBranchesInfoList = list.New([]list.Item{}, branchComponent.GitBranchItemDelegate{}, 0, 0)
	m.CurrentRepoTagInfoList = list.New([]list.Item{}, tagComponent.GitTagItemDelegate{}, 0, 0)
	m.CurrentRepoModifiedFilesInfoList = list.New([]list.Item{}, filesComponent.GitModifiedFilesItemDelegate{}, 0, 0)
	m.CurrentRepoCommitLogInfoList = list.New([]list.Item{}, commitlogComponent.GitCommitLogItemDelegate{}, 0, 0)
	m.CurrentRepoRefLogInfoList = list.New([]list.Item{}, reflogComponent.GitRefLogItemDelegate{}, 0, 0)
	m.CurrentRepoStashInfoList = list.New([]list.Item{}, stashComponent.GitStashItemDelegate{}, 0, 0)
	m.CurrentRepoRemoteInfoList = list.New([]list.Item{}, remoteComponent.GitRemoteItemDelegate{}, 0, 0)
	m.CurrentRepoWorktreeInfoList = list.New([]list.Item{}, worktreeComponent.GitWorktreeItemDelegate{}, 0, 0)
	m.DetailPanelParentComponent = ""
	m.DetailPanelViewport = vp
	m.DetailPanelViewportOffset = 0
	m.DetailPanelTwoViewport = vpTwo
	m.DetailPanelTwoViewportOffset = 0
	m.DetailComponentPanelLayout = constant.HORIZONTAL
	m.ListNavigationIndexPosition = types.GittiComponentsCurrentListNavigationIndexPosition{LocalBranchComponent: 0, ModifiedFilesComponent: 0, StashComponent: 0, RefLogComponent: 0, TagComponent: 0, RemoteComponent: 0, WorktreeComponent: 0}
	m.PopUpType = constant.NoPopUp
	m.PopUpModel = struct{}{}
	m.GitOperations = gitOperations
	m.GlobalKeyBindingKeyMapLargestLen = 0
	m.LocalBranchComponentKeyBindingKeyMapLargestLen = 0
	m.TagComponentKeyBindingKeyMapLargestLen = 0
	m.RemoteComponentKeyBindingKeyMapLargestLen = 0
	m.ModifiedFilesComponentKeyBindingKeyMapLargestLen = 0
	m.CommitLogComponentKeyBindingKeyMapLargestLen = 0
	m.RefLogComponentKeyBindingKeyMapLargestLen = 0
	m.StashComponentKeyBindingKeyMapLargestLen = 0
	m.LogComponentKeyBindingKeyMapLargestLen = 0
	m.DetailComponentKeyBindingKeyMapLargestLen = 0
	m.LineEditingIndexPositionAndInfo = types.GittiLineEditingIndexPositionAndInfo{}
	m.LineEditingIndexCursorViewport = lineEditingIndexCursorVp
	m.LineEditingIndexCursorTwoViewport = lineEditingIndexCursorVpTwo
	m.CherryPickedCommitInfo = types.CherryPickedCommitInfo{LatestSequenceCounter: 0, CherryPickedCommitMap: make(map[string]git.CherryPickedCommitLog)}
	m.PanelFilterQuery = make(map[string]string)

	m.IsRenderInit.Store(false)
	m.IsPanelFiltering.Store(false)
	m.ShowPopUp.Store(false)
	m.IsTyping.Store(false)
	m.IsDetailComponentPanelInfoFetchProcessing.Store(false)
	m.ShowDetailPanelTwo.Store(false)
	m.IsLineEditingState.Store(false)

	commitRequireSigning, tagRequireSigning, pushRequireSigning := api.CheckSigningRequiredOperation()
	m.GitCommitRequireSigning = commitRequireSigning
	m.GitTagRequireSigning = tagRequireSigning
	m.GitPushRequireSigning = pushRequireSigning

	m.GittiLogger.ClearLogs()
	m.CurrentLogComponentViewport.SetContent("")
}
