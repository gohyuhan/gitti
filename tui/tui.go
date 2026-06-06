package tui

import (
	"fmt"

	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/logging"
	branchComponent "github.com/gohyuhan/gitti/tui/component/branch"
	commitlogComponent "github.com/gohyuhan/gitti/tui/component/commitlog"
	filesComponent "github.com/gohyuhan/gitti/tui/component/files"
	logComponent "github.com/gohyuhan/gitti/tui/component/log"
	reflogComponent "github.com/gohyuhan/gitti/tui/component/reflog"
	remoteComponent "github.com/gohyuhan/gitti/tui/component/remote"
	stashComponent "github.com/gohyuhan/gitti/tui/component/stash"
	tagComponent "github.com/gohyuhan/gitti/tui/component/tag"
	worktreeComponent "github.com/gohyuhan/gitti/tui/component/worktree"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/helper"
	"github.com/gohyuhan/gitti/tui/initialize"
	"github.com/gohyuhan/gitti/tui/interaction"
	"github.com/gohyuhan/gitti/tui/layout"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	pullPopUp "github.com/gohyuhan/gitti/tui/popup/pull"
	pushPopUp "github.com/gohyuhan/gitti/tui/popup/push"
	rebasePopUp "github.com/gohyuhan/gitti/tui/popup/rebase"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"

	tea "charm.land/bubbletea/v2"
)

// ------------------------------------
//
//	Construct and return a new GittiAppModel. Creates the three detail viewports
//	(primary, secondary, log), two cursor viewports for line-editing mode,
//	initialises all list components, and checks whether git signing is required
//	for commit/tag/push operations.
//
// ------------------------------------
func NewGittiAppModel(tuiUpdateChannel chan interface{}, repoPath string, repoName string, gitOperations *api.GitOperations, gittiLogger *logging.GittiLogging, daemonUpdateChannel chan string, gitUpdateChannel chan string) *GittiAppModel {
	gittiModel := initialize.InitGittiModel(tuiUpdateChannel, repoPath, repoName, gitOperations, gittiLogger, daemonUpdateChannel, gitUpdateChannel)
	return &GittiAppModel{model: gittiModel}
}

// ------------------------------------
//
//	Bubble tea Init function, called once when the program starts
//
// ------------------------------------
func (gAM *GittiAppModel) Init() tea.Cmd {
	return nil
}

// ------------------------------------
//
//	Bubble tea Update function, handles all incoming messages and state changes
//
// ------------------------------------
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
	case tea.KeyPressMsg:
		model, cmd := interaction.GittiKeyInteraction(msg, m)
		gAM.model = model
		return gAM, cmd
	case GitUpdateMsg:
		updateEvent := string(msg)
		switch updateEvent {
		case git.GIT_BRANCH_UPDATE:
			branchComponent.InitBranchList(m)
			if m.CurrentSelectedComponent == constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel && m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing == constant.SHOW_LOCAL_BRANCH {
				services.FetchDetailComponentPanelInfoService(m, false)
			}
		case git.GIT_TAG_UPDATE:
			needReinit := tagComponent.InitTagList(m)
			if m.CurrentSelectedComponent == constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel && m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing == constant.SHOW_TAG {
				services.FetchDetailComponentPanelInfoService(m, needReinit)
			}
		case git.GIT_REMOTE_UPDATE:
			needReinit := remoteComponent.InitRemoteList(m)
			if m.CurrentSelectedComponent == constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel && m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing == constant.SHOW_REMOTE {
				services.FetchDetailComponentPanelInfoService(m, needReinit)
			}
		case git.GIT_WORKTREE_UPDATE:
			needReinit := worktreeComponent.InitWorktreeList(m)
			if m.CurrentSelectedComponent == constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel && m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing == constant.SHOW_WORKTREE {
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
		case git.GIT_REBASE_OUTPUT_UPDATE:
			rebasePopUp.UpdatePopUpGitRebaseOutputViewport(m)
		}
		return gAM, nil
	case types.EditorFinishedMsg:
		return gAM, nil
	case types.GitOperationRequiredSigningFinishedMsg:
		if msg.Err != nil {
			m.GittiLogger.RegisterNewLog(msg.GitOperationOpsTypeForLogging, "", logging.ERROR, fmt.Sprintf("[%s ERROR] %s", msg.GitOperationOpsTypeForLogging, msg.Err.Error()), false)
			if msg.GitOperationOpsTypeForLogging != logging.COMMIT_WITH_SIGNING_OPS && msg.GitOperationOpsTypeForLogging != logging.AMEND_COMMIT_WITH_SIGNING_OPS {
				utils.ResetPopUpModelStateForGitSigningOps(m)
				utils.ReinitCherryPickedCommitInfo(m)
			}
		} else {
			utils.ResetPopUpModelStateForGitSigningOps(m)
			utils.ReinitCherryPickedCommitInfo(m)
		}

		// ensure resumed signing flow does not preserve horizontal log scroll offset
		m.CurrentLogComponentViewport.SetXOffset(0)
		return gAM, nil
	case tea.MouseMsg:
		model, cmd := interaction.GittiMouseInteraction(msg, m)
		gAM.model = model
		return gAM, cmd

	case types.GittiTuiUpdateMsg:
		helper.GittiTuiUpdateEventHelper(m, msg)
	}

	// Update spinners in popups when they are processing
	cmds = helper.UpdateSpinner(m, msg, cmds)

	return gAM, tea.Batch(cmds...)
}

// ------------------------------------
//
//	Bubble tea View function, renders the UI based on the current model state
//
// ------------------------------------
func (gAM *GittiAppModel) View() tea.View {
	var v tea.View
	v.SetContent(layout.GittiMainPageView(gAM.model))
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

// ------------------------------------
//
//	Update the Git remote status, upstream branch, and sync line string in the model
//
// ------------------------------------
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

// ------------------------------------
//
//	Update the current repository state (e.g., REBASE, MERGE, etc.) in the model
//
// ------------------------------------
func (gAM *GittiAppModel) updateGitRepoState() {
	m := gAM.model
	m.CurrentGitRepoStatus = m.GitOperations.GitStateUniversalUtils.GetCurrentGitState()
}
