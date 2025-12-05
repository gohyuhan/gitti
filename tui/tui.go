package tui

import (
	"fmt"

	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/component/branch"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

func NewGittiAppModel(tuiUpdateChannel chan string, repoPath string, repoName string, gitOperations *api.GitOperations) *GittiAppModel {
	vp := viewport.New()
	vp.SoftWrap = false
	vp.MouseWheelEnabled = true
	vp.SetHorizontalStep(1)
	vp.MouseWheelDelta = 1
	gittiModel := &types.GittiModel{
		TuiUpdateChannel:                 tuiUpdateChannel,
		CurrentSelectedComponent:         constant.ModifiedFilesComponent,
		CurrentSelectedComponentIndex:    2,
		TotalComponentCount:              4,
		RepoPath:                         repoPath,
		RepoName:                         repoName,
		CheckOutBranch:                   "",
		RemoteSyncStateLineString:        "",
		BranchUpStream:                   "",
		TrackedUpstreamOrBranchIcon:      "",
		Width:                            0,
		Height:                           0,
		CurrentRepoBranchesInfoList:      list.New([]list.Item{}, branch.GitBranchItemDelegate{}, 0, 0),
		CurrentRepoModifiedFilesInfoList: list.New([]list.Item{}, gitModifiedFilesItemDelegate{}, 0, 0),
		CurrentRepoStashInfoList:         list.New([]list.Item{}, gitStashItemDelegate{}, 0, 0),
		DetailPanelParentComponent:       "",
		DetailPanelViewport:              vp,
		DetailPanelViewportOffset:        0,
		ListNavigationIndexPosition:      types.GittiComponentsCurrentListNavigationIndexPosition{LocalBranchComponent: 0, ModifiedFilesComponent: 0, StashComponent: 0},
		PopUpType:                        constant.NoPopUp,
		PopUpModel:                       struct{}{},
		GitOperations:                    gitOperations,
		GlobalKeyBindingKeyMapLargestLen: 0,
	}
	gittiModel.IsRenderInit.Store(false)
	gittiModel.ShowPopUp.Store(false)
	gittiModel.IsTyping.Store(false)
	gittiModel.IsDetailComponentPanelInfoFetchProcessing.Store(false)

	return &GittiAppModel{model: gittiModel}
}

// -----------------------------------------------------------------------------
// Bubble Tea standard functions
// -----------------------------------------------------------------------------

func (gAM *GittiAppModel) Init() tea.Cmd {
	return nil
}

func (gAM *GittiAppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	m := gAM.model
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		// recompute layout instantly
		tuiWindowSizing(m)
		// Initialize list components once, immediately after the first window resize.
		// Valid dimensions are required to calculate item layouts (specifically text truncation);
		// initializing earlier would cause the UI layout to break.
		if m.IsRenderInit.CompareAndSwap(false, true) {
			branch.InitBranchList(m)
			initModifiedFilesList(m)
			initStashList(m)
		}
	case tea.KeyMsg:
		model, cmd := gittiKeyInteraction(msg, m)
		gAM.model = model
		return gAM, cmd
	case GitUpdateMsg:
		updateEvent := string(msg)
		switch updateEvent {
		case constant.DETAIL_COMPONENT_PANEL_UPDATED:
			return gAM, nil
		case git.GIT_BRANCH_UPDATE:
			branch.InitBranchList(m)
			if m.CurrentSelectedComponent == constant.LocalBranchComponent {
				fetchDetailComponentPanelInfoService(m, false)
			}
		case git.GIT_FILES_STATUS_UPDATE:
			needReinit := initModifiedFilesList(m)
			if m.CurrentSelectedComponent == constant.ModifiedFilesComponent {
				fetchDetailComponentPanelInfoService(m, needReinit)
			}
		case git.GIT_STASH_UPDATE:
			needReinit := initStashList(m)
			if m.CurrentSelectedComponent == constant.StashComponent {
				fetchDetailComponentPanelInfoService(m, needReinit)
			}
		case git.GIT_COMMIT_OUTPUT_UPDATE:
			updatePopUpCommitOutputViewPort(m)
		case git.GIT_AMEND_COMMIT_OUTPUT_UPDATE:
			updatePopUpAmendCommitOutputViewPort(m)
		case git.GIT_REMOTE_PUSH_OUTPUT_UPDATE:
			updatePopUpGitRemotePushOutputViewport(m)
		case git.GIT_PULL_OUTPUT_UPDATE:
			updatePopUpGitPullOutputViewport(m)
		case git.GIT_REMOTE_SYNC_STATUS_AND_UPSTREAM_UPDATE:
			gAM.updateGitRemoteStatusSyncLineStringAndUpStream()
		}
		return gAM, nil
	case tea.MouseMsg:
		model, cmd := GittiMouseInteraction(msg, m)
		gAM.model = model
		return gAM, cmd
	}

	// Update spinners in popups when they are processing
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.CommitPopUp:
			if commitPopup, ok := m.PopUpModel.(*GitCommitPopUpModel); ok && commitPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				commitPopup.Spinner, cmd = commitPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.AmendCommitPopUp:
			if amendCommitPopup, ok := m.PopUpModel.(*GitAmendCommitPopUpModel); ok && amendCommitPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				amendCommitPopup.Spinner, cmd = amendCommitPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.GitRemotePushPopUp:
			if pushPopup, ok := m.PopUpModel.(*GitRemotePushPopUpModel); ok && pushPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				pushPopup.Spinner, cmd = pushPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.GitPullOutputPopUp:
			if pullPopup, ok := m.PopUpModel.(*GitPullOutputPopUpModel); ok && pullPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				pullPopup.Spinner, cmd = pullPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.SwitchBranchOutputPopUp:
			if pullPopup, ok := m.PopUpModel.(*SwitchBranchOutputPopUpModel); ok && pullPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				pullPopup.Spinner, cmd = pullPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.GitStashOperationOutputPopUp:
			if stashPopup, ok := m.PopUpModel.(*GitStashOperationOutputPopUpModel); ok && stashPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				stashPopup.Spinner, cmd = stashPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}
	return gAM, tea.Batch(cmds...)
}

func (gAM *GittiAppModel) View() tea.View {
	var v tea.View
	v.SetContent(gittiMainPageView(gAM.model))
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func (gAM *GittiAppModel) updateGitRemoteStatusSyncLineStringAndUpStream() {
	m := gAM.model
	// set branch upstream
	m.TrackedUpstreamOrBranchIcon = m.GitOperations.GitRemote.UpStreamRemoteIcon()
	m.BranchUpStream = m.GitOperations.GitRemote.CurrentBranchUpStream()

	// set remote sync status
	remoteSynsStatusInfo := m.GitOperations.GitRemote.RemoteSyncStatus()
	if remoteSynsStatusInfo.Local == "" || remoteSynsStatusInfo.Remote == "" {
		m.RemoteSyncStateLineString = style.ErrorStyle.Render("\uf00d")
		return
	}

	local := style.LocalStatusStyle.Render(fmt.Sprintf("%s↑", remoteSynsStatusInfo.Local))
	remote := style.RemoteStatusStyle.Render(fmt.Sprintf("%s↓", remoteSynsStatusInfo.Remote))

	m.RemoteSyncStateLineString = local + " " + remote
}
