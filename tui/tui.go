package tui

import (
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/settings"
	branchComponent "github.com/gohyuhan/gitti/tui/component/branch"
	commitlogComponent "github.com/gohyuhan/gitti/tui/component/commitlog"
	filesComponent "github.com/gohyuhan/gitti/tui/component/files"
	stashComponent "github.com/gohyuhan/gitti/tui/component/stash"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/interaction"
	"github.com/gohyuhan/gitti/tui/layout"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	pullPopUp "github.com/gohyuhan/gitti/tui/popup/pull"
	pushPopUp "github.com/gohyuhan/gitti/tui/popup/push"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	"github.com/gohyuhan/gitti/tui/services"
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

	vpTwo := viewport.New()
	vpTwo.SoftWrap = false
	vpTwo.MouseWheelEnabled = true
	vpTwo.SetHorizontalStep(1)
	vpTwo.MouseWheelDelta = 1

	gittiModel := &types.GittiModel{
		TuiUpdateChannel:                 tuiUpdateChannel,
		UserSetEditor:                    settings.GITTICONFIGSETTINGS.Editor,
		CurrentSelectedComponent:         constant.ModifiedFilesComponent,
		CurrentSelectedComponentIndex:    2,
		TotalComponentCount:              4,
		RepoPath:                         repoPath,
		RepoName:                         repoName,
		CheckOutBranch:                   "",
		RemoteSyncLocalState:             "",
		RemoteSyncRemoteState:            "",
		BranchUpStream:                   "",
		TrackedUpstreamOrBranchIcon:      "",
		Width:                            0,
		Height:                           0,
		WindowLeftPanelRatio:             settings.GITTICONFIGSETTINGS.LeftPanelWidthRatio,
		CurrentRepoBranchesInfoList:      list.New([]list.Item{}, branchComponent.GitBranchItemDelegate{}, 0, 0),
		CurrentRepoModifiedFilesInfoList: list.New([]list.Item{}, filesComponent.GitModifiedFilesItemDelegate{}, 0, 0),
		CurrentRepoCommitLogInfoList:     list.New([]list.Item{}, commitlogComponent.GitCommitLogItemDelegate{}, 0, 0),
		CurrentRepoStashInfoList:         list.New([]list.Item{}, stashComponent.GitStashItemDelegate{}, 0, 0),
		DetailPanelParentComponent:       "",
		DetailPanelViewport:              vp,
		DetailPanelViewportOffset:        0,
		DetailPanelTwoViewport:           vpTwo,
		DetailPanelTwoViewportOffset:     0,
		DetailComponentPanelLayout:       constant.HORIZONTAL,
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
	gittiModel.ShowDetailPanelTwo.Store(false)

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
		layout.TuiWindowSizing(m)
		// Initialize list components once, immediately after the first window resize.
		// Valid dimensions are required to calculate item layouts (specifically text truncation);
		// initializing earlier would cause the UI layout to break.
		if m.IsRenderInit.CompareAndSwap(false, true) {
			branchComponent.InitBranchList(m)
			filesComponent.InitModifiedFilesList(m)
			commitlogComponent.InitGitCommitLogList(m)
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
			layout.UpdateDetailComponentViewportLayout(gAM.model)
			return gAM, nil
		case git.GIT_BRANCH_UPDATE:
			branchComponent.InitBranchList(m)
			if m.CurrentSelectedComponent == constant.LocalBranchComponent {
				services.FetchDetailComponentPanelInfoService(m, false)
			}
		case git.GIT_FILES_STATUS_UPDATE:
			needReinit := filesComponent.InitModifiedFilesList(m)
			if m.CurrentSelectedComponent == constant.ModifiedFilesComponent {
				services.FetchDetailComponentPanelInfoService(m, needReinit)
			}
		case git.GIT_LOG_UPDATE:
			needReinit := commitlogComponent.InitGitCommitLogList(m)
			if m.CurrentSelectedComponent == constant.CommitLogComponent {
				services.FetchDetailComponentPanelInfoService(m, needReinit)
			}
		case git.GIT_STASH_UPDATE:
			needReinit := stashComponent.InitStashList(m)
			if m.CurrentSelectedComponent == constant.StashComponent {
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
		}
		return gAM, nil
	case types.EditorFinishedMsg:
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
		}
	}
	return gAM, tea.Batch(cmds...)
}

func (gAM *GittiAppModel) View() tea.View {
	var v tea.View
	v.SetContent(layout.GittiMainPageView(gAM.model))
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
	m.RemoteSyncLocalState = remoteSynsStatusInfo.Local
	m.RemoteSyncRemoteState = remoteSynsStatusInfo.Remote
}
