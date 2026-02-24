package tui

import (
	"fmt"

	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/settings"
	branchComponent "github.com/gohyuhan/gitti/tui/component/branch"
	commitlogComponent "github.com/gohyuhan/gitti/tui/component/commitlog"
	filesComponent "github.com/gohyuhan/gitti/tui/component/files"
	logComponent "github.com/gohyuhan/gitti/tui/component/log"
	reflogComponent "github.com/gohyuhan/gitti/tui/component/reflog"
	remoteComponent "github.com/gohyuhan/gitti/tui/component/remote"
	stashComponent "github.com/gohyuhan/gitti/tui/component/stash"
	tagComponent "github.com/gohyuhan/gitti/tui/component/tag"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/interaction"
	"github.com/gohyuhan/gitti/tui/layout"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	pullPopUp "github.com/gohyuhan/gitti/tui/popup/pull"
	pushPopUp "github.com/gohyuhan/gitti/tui/popup/push"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

// ----------------------------------
//
//	Initialize and return a new GittiAppModel with necessary dependencies and viewport setup
//
// ----------------------------------
func NewGittiAppModel(tuiUpdateChannel chan string, repoPath string, repoName string, gitOperations *api.GitOperations, gittiLogger *logging.GittiLogging, daemonUpdateChannel chan string) *GittiAppModel {
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
		TuiUpdateChannel:              tuiUpdateChannel,
		UserSetEditor:                 settings.GITTICONFIGSETTINGS.Editor,
		CurrentSelectedComponent:      constant.ModifiedFilesComponentPanel,
		CurrentSelectedComponentIndex: 2,
		CurrentLocalBranchOrTagOrRemoteComponentShowing:  constant.SHOW_LOCAL_BRANCH,
		CurrentCommitLogOrRefLogComponentShowing:         constant.SHOW_COMMITLOG,
		TotalComponentCount:                              4,
		RepoPath:                                         repoPath,
		RepoName:                                         repoName,
		CheckOutBranch:                                   "",
		RemoteSyncLocalState:                             "",
		RemoteSyncRemoteState:                            "",
		CurrentGitRepoStatus:                             "",
		BranchUpStream:                                   "",
		TrackedUpstreamOrBranchIcon:                      "",
		Width:                                            0,
		Height:                                           0,
		WindowLeftPanelRatio:                             settings.GITTICONFIGSETTINGS.LeftPanelWidthRatio,
		CurrentRepoBranchesInfoList:                      list.New([]list.Item{}, branchComponent.GitBranchItemDelegate{}, 0, 0),
		CurrentRepoTagInfoList:                           list.New([]list.Item{}, tagComponent.GitTagItemDelegate{}, 0, 0),
		CurrentRepoModifiedFilesInfoList:                 list.New([]list.Item{}, filesComponent.GitModifiedFilesItemDelegate{}, 0, 0),
		CurrentRepoCommitLogInfoList:                     list.New([]list.Item{}, commitlogComponent.GitCommitLogItemDelegate{}, 0, 0),
		CurrentRepoRefLogInfoList:                        list.New([]list.Item{}, reflogComponent.GitRefLogItemDelegate{}, 0, 0),
		CurrentRepoStashInfoList:                         list.New([]list.Item{}, stashComponent.GitStashItemDelegate{}, 0, 0),
		CurrentRepoRemoteInfoList:                        list.New([]list.Item{}, remoteComponent.GitRemoteItemDelegate{}, 0, 0),
		DetailPanelParentComponent:                       "",
		DetailPanelViewport:                              vp,
		DetailPanelViewportOffset:                        0,
		DetailPanelTwoViewport:                           vpTwo,
		DetailPanelTwoViewportOffset:                     0,
		DetailComponentPanelLayout:                       constant.HORIZONTAL,
		CurrentLogComponentViewport:                      logVp,
		ListNavigationIndexPosition:                      types.GittiComponentsCurrentListNavigationIndexPosition{LocalBranchComponent: 0, ModifiedFilesComponent: 0, StashComponent: 0, RefLogComponent: 0, TagComponent: 0, RemoteComponent: 0},
		PopUpType:                                        constant.NoPopUp,
		PopUpModel:                                       struct{}{},
		GitOperations:                                    gitOperations,
		GlobalKeyBindingKeyMapLargestLen:                 0,
		LocalBranchComponentKeyBindingKeyMapLargestLen:   0,
		TagComponentKeyBindingKeyMapLargestLen:           0,
		RemoteComponentKeyBindingKeyMapLargestLen:        0,
		ModifiedFilesComponentKeyBindingKeyMapLargestLen: 0,
		CommitLogComponentKeyBindingKeyMapLargestLen:     0,
		RefLogComponentKeyBindingKeyMapLargestLen:        0,
		StashComponentKeyBindingKeyMapLargestLen:         0,
		LogComponentKeyBindingKeyMapLargestLen:           0,
		DetailComponentKeyBindingKeyMapLargestLen:        0,
		LineEditingIndexPositionAndInfo:                  types.GittiLineEditingIndexPositionAndInfo{},
		LineEditingIndexCursorViewport:                   lineEditingIndexCursorVp,
		LineEditingIndexCursorTwoViewport:                lineEditingIndexCursorVpTwo,
		CherryPickedCommitInfo:                           types.CherryPickedCommitInfo{LatestSequenceCounter: 0, CherryPickedCommitMap: make(map[string]git.CherryPickedCommitLog)},
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

	return &GittiAppModel{model: gittiModel}
}

// -----------------------------------------------------------------------------
// Bubble Tea standard functions
// -----------------------------------------------------------------------------

// ----------------------------------
//
//	Bubble tea Init function, called once when the program starts
//
// ----------------------------------
func (gAM *GittiAppModel) Init() tea.Cmd {
	return nil
}

// ----------------------------------
//
//	Bubble tea Update function, handles all incoming messages and state changes
//
// ----------------------------------
func (gAM *GittiAppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	m := gAM.model
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		// recompute layout instantly
		layout.TuiWindowSizing(m)
		// Initialize list components once, immediately after the first window resize.
		// Valid dimensions are required to calculate item layouts (specifically text truncation);
		// initializing earlier would cause the UI layout to break.
		if m.IsRenderInit.CompareAndSwap(false, true) {
			branchComponent.InitBranchList(m)
			tagComponent.InitTagList(m)
			remoteComponent.InitRemoteList(m)
			filesComponent.InitModifiedFilesList(m)
			commitlogComponent.InitGitCommitLogList(m)
			reflogComponent.InitGitRefLogList(m)
			stashComponent.InitStashList(m)
		}
	case tea.KeyMsg:
		model, cmd := interaction.GittiKeyInteraction(msg, m)
		gAM.model = model
		return gAM, cmd
	case GitUpdateMsg:
		updateEvent := string(msg)
		switch updateEvent {
		case constant.DETAIL_COMPONENT_PANEL_UPDATED:
			layout.UpdateDetailComponentViewportLayout(m)
		case git.GIT_BRANCH_UPDATE:
			branchComponent.InitBranchList(m)
			if m.CurrentSelectedComponent == constant.LocalBranchOrTagOrRemoteComponentPanel && m.CurrentLocalBranchOrTagOrRemoteComponentShowing == constant.SHOW_LOCAL_BRANCH {
				services.FetchDetailComponentPanelInfoService(m, false)
			}
		case git.GIT_TAG_UPDATE:
			needReinit := tagComponent.InitTagList(m)
			if m.CurrentSelectedComponent == constant.LocalBranchOrTagOrRemoteComponentPanel && m.CurrentLocalBranchOrTagOrRemoteComponentShowing == constant.SHOW_TAG {
				services.FetchDetailComponentPanelInfoService(m, needReinit)
			}
		case git.GIT_REMOTE_UPDATE:
			needReinit := remoteComponent.InitRemoteList(m)
			if m.CurrentSelectedComponent == constant.LocalBranchOrTagOrRemoteComponentPanel && m.CurrentLocalBranchOrTagOrRemoteComponentShowing == constant.SHOW_REMOTE {
				services.FetchDetailComponentPanelInfoService(m, needReinit)
			}
		case git.GIT_FILES_STATUS_UPDATE:
			needReinit := filesComponent.InitModifiedFilesList(m)
			if m.CurrentSelectedComponent == constant.ModifiedFilesComponentPanel || m.DetailPanelParentComponent == constant.ModifiedFilesComponentPanel {
				services.FetchDetailComponentPanelInfoService(m, needReinit)
			}
		case git.GIT_STATE_UPDATE:
			gAM.updateGitRepoState()
		case git.GIT_COMMITLOG_UPDATE:
			needReinit := commitlogComponent.InitGitCommitLogList(m)
			if m.CurrentSelectedComponent == constant.CommitLogOrRefLogComponentPanel && m.CurrentCommitLogOrRefLogComponentShowing == constant.SHOW_COMMITLOG {
				services.FetchDetailComponentPanelInfoService(m, needReinit)
			}
		case git.GIT_REFLOG_UPDATE:
			needReinit := reflogComponent.InitGitRefLogList(m)
			if m.CurrentSelectedComponent == constant.CommitLogOrRefLogComponentPanel && m.CurrentCommitLogOrRefLogComponentShowing == constant.SHOW_REFLOG {
				services.FetchDetailComponentPanelInfoService(m, needReinit)
			}
		case git.GIT_STASH_UPDATE:
			needReinit := stashComponent.InitStashList(m)
			if m.CurrentSelectedComponent == constant.StashComponentPanel {
				services.FetchDetailComponentPanelInfoService(m, needReinit)
			}
		case git.GIT_COMMIT_OUTPUT_UPDATE:
			commitPopUp.UpdatePopUpCommitOutputViewPort(m)
		case git.GIT_AMEND_COMMIT_OUTPUT_UPDATE:
			commitPopUp.UpdatePopUpAmendCommitOutputViewPort(m)
		case git.GIT_REMOTE_PUSH_OUTPUT_UPDATE:
			pushPopUp.UpdatePopUpGitRemotePushOutputViewport(m)
		case git.GIT_PULL_OUTPUT_UPDATE:
			pullPopUp.UpdatePopUpGitPullOutputViewport(m)
		case git.GIT_REMOTE_SYNC_STATUS_AND_UPSTREAM_UPDATE:
			gAM.updateGitRemoteStatusSyncLineStringAndUpStream()
		case git.GIT_EDIT_LINE_DETAILS_AND_FILES_UPDATE:
			needReinit := filesComponent.InitModifiedFilesList(m)
			services.FetchDetailComponentPanelInfoService(m, needReinit)
		case logging.NEW_LOG_UPDATE:
			content := logComponent.InitGittiLogViewport(m, true, nil)
			m.CurrentLogComponentViewport.SetContent(content)
			m.CurrentLogComponentViewport.GotoBottom()
			if m.CurrentSelectedComponent == constant.LogComponentPanel || m.DetailPanelParentComponent == constant.LogComponentPanel {
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case git.GIT_TAG_PUSH_OUTPUT_UPDATE:
			tagPopUp.UpdatePushTagOutputViewPort(m)
		case git.GIT_TAG_FETCH_OUTPUT_UPDATE:
			tagPopUp.UpdateFetchTagOutputViewPort(m)
		}
		return gAM, nil
	case types.EditorFinishedMsg:
		return gAM, nil
	case types.GitOperationRequiredSigningFinishedMsg:
		if msg.Err != nil {
			m.GittiLogger.RegisterNewLog(msg.GitOperationOpsTypeForLogging, "", logging.ERROR, fmt.Sprintf("[%s ERROR] %s", msg.GitOperationOpsTypeForLogging, msg.Err.Error()), false)
		} else {
			// After a signed Git operation (executed directly in the terminal) completes successfully,
			// we perform a global state reset. This safely clears any active popups and resets temporary data
			// (like cherry-pick selections) for all supported operations, eliminating the need for operation-specific cleanup logic.

			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpModel = nil
			m.PopUpType = constant.NoPopUp
			utils.ReinitCherryPickedCommitInfo(m)
		}
		return gAM, nil
	case tea.MouseMsg:
		model, cmd := interaction.GittiMouseInteraction(msg, m)
		gAM.model = model
		return gAM, cmd
	}

	// Update spinners in popups when they are processing
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.CommitPopUp:
			if commitPopup, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel); ok && commitPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				commitPopup.Spinner, cmd = commitPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.AmendCommitPopUp:
			if amendCommitPopup, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel); ok && amendCommitPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				amendCommitPopup.Spinner, cmd = amendCommitPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.GitRemotePushPopUp:
			if pushPopup, ok := m.PopUpModel.(*pushPopUp.GitRemotePushPopUpModel); ok && pushPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				pushPopup.Spinner, cmd = pushPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.GitPullOutputPopUp:
			if pullPopup, ok := m.PopUpModel.(*pullPopUp.GitPullOutputPopUpModel); ok && pullPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				pullPopup.Spinner, cmd = pullPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.SwitchBranchOutputPopUp:
			if pullPopup, ok := m.PopUpModel.(*branchPopUp.SwitchBranchOutputPopUpModel); ok && pullPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				pullPopup.Spinner, cmd = pullPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.GitStashOperationOutputPopUp:
			if stashPopup, ok := m.PopUpModel.(*stashPopUp.GitStashOperationOutputPopUpModel); ok && stashPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				stashPopup.Spinner, cmd = stashPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.GitDeleteBranchOutputPopUp:
			if branchPopup, ok := m.PopUpModel.(*branchPopUp.GitDeleteBranchOutputPopUpModel); ok && branchPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				branchPopup.Spinner, cmd = branchPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.CreateBranchBasedOnRemoteOutputPopUp:
			if branchPopup, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemoteOutputPopUpModel); ok && branchPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				branchPopup.Spinner, cmd = branchPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.DeleteTagOutputPopUp:
			if tagPopup, ok := m.PopUpModel.(*tagPopUp.DeleteTagOutputPopUpModel); ok && tagPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				tagPopup.Spinner, cmd = tagPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.PushTagOutputPopUp:
			tagPopUp, ok := m.PopUpModel.(*tagPopUp.PushTagOutputPopUpModel)
			if ok {
				var cmd tea.Cmd
				tagPopUp.Spinner, cmd = tagPopUp.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.FetchTagOutputPopUp:
			if tagPopup, ok := m.PopUpModel.(*tagPopUp.FetchTagOutputPopUpModel); ok && tagPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				tagPopup.Spinner, cmd = tagPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}
	return gAM, tea.Batch(cmds...)
}

// ----------------------------------
//
//	Bubble tea View function, renders the UI based on the current model state
//
// ----------------------------------
func (gAM *GittiAppModel) View() tea.View {
	var v tea.View
	v.SetContent(layout.GittiMainPageView(gAM.model))
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

// ----------------------------------
//
//	Update the Git remote status, upstream branch, and sync line string in the model
//
// ----------------------------------
func (gAM *GittiAppModel) updateGitRemoteStatusSyncLineStringAndUpStream() {
	m := gAM.model
	// set branch upstream
	m.TrackedUpstreamOrBranchIcon = m.GitOperations.GitRemote.UpStreamRemoteIcon()
	m.BranchUpStream = m.GitOperations.GitRemote.CurrentBranchUpStream()

	// set remote sync status
	remoteSynsStatusInfo := m.GitOperations.GitRemote.RemoteSyncStatus()
	m.RemoteSyncLocalState = remoteSynsStatusInfo.Local
	m.RemoteSyncRemoteState = remoteSynsStatusInfo.Remote
}

// ----------------------------------
//
//	Update the current repository state (e.g., REBASE, MERGE, etc.) in the model
//
// ----------------------------------
func (gAM *GittiAppModel) updateGitRepoState() {
	m := gAM.model
	m.CurrentGitRepoStatus = m.GitOperations.GitStateUniversalUtils.GetCurrentGitState()
}
