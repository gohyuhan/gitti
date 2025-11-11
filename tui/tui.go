package tui

import (
	"gitti/api"
	"gitti/api/git"
	"gitti/tui/constant"

	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
)

func NewGittiModel(repoPath string, gitState *api.GitState) *GittiModel {
	vp := viewport.New()
	vp.SoftWrap = false
	vp.MouseWheelEnabled = true
	vp.SetHorizontalStep(1)
	vp.MouseWheelDelta = 1
	gitti := &GittiModel{
		CurrentSelectedComponent:         constant.ModifiedFilesComponent,
		CurrentSelectedComponentIndex:    2,
		TotalComponentCount:              4,
		RepoPath:                         repoPath,
		Width:                            0,
		Height:                           0,
		CurrentRepoBranchesInfoList:      list.New([]list.Item{}, gitBranchItemDelegate{}, 0, 0),
		CurrentRepoModifiedFilesInfoList: list.New([]list.Item{}, gitModifiedFilesItemDelegate{}, 0, 0),
		CurrentRepoStashInfoList:         list.New([]list.Item{}, gitStashItemDelegate{}, 0, 0),
		DetailPanelParentComponent:       "",
		DetailPanelViewport:              vp,
		DetailPanelViewportOffset:        0,
		ListNavigationIndexPosition:      GittiComponentsCurrentListNavigationIndexPosition{LocalBranchComponent: 0, ModifiedFilesComponent: 0, StashComponent: 0},
		PopUpType:                        constant.NoPopUp,
		PopUpModel:                       struct{}{},
		GitState:                         gitState,
		GlobalKeyBindingKeyMapLargestLen: 0,
	}
	gitti.IsRenderInit.Store(false)
	gitti.ShowPopUp.Store(false)
	gitti.IsTyping.Store(false)

	return gitti
}

// -----------------------------------------------------------------------------
// Bubble Tea standard functions
// -----------------------------------------------------------------------------

func (m *GittiModel) Init() tea.Cmd {
	return nil
}

func (m *GittiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		// recompute layout instantly
		tuiWindowSizing(m)
		if m.IsRenderInit.CompareAndSwap(false, true) {
			initBranchList(m)
			initModifiedFilesList(m)
			initStashList(m)
		}
	case tea.KeyMsg:
		return gittiKeyInteraction(msg, m)
	case GitUpdateMsg:
		updateEvent := string(msg)
		switch updateEvent {
		case git.GIT_BRANCH_UPDATE:
			initBranchList(m)
		case git.GIT_FILES_STATUS_UPDATE:
			initModifiedFilesList(m)
			renderDetailComponentPanelViewPort(m)
		case git.GIT_STASH_UPDATE:
			initStashList(m)
		case git.GIT_COMMIT_OUTPUT_UPDATE:
			updatePopUpCommitOutputViewPort(m)
		case git.GIT_REMOTE_PUSH_OUTPUT_UPDATE:
			updateGitRemotePushOutputViewport(m)
		case git.GIT_PULL_OUTPUT_UPDATE:
			updateGitPullOutputViewport(m)
		}
		return m, nil
	case tea.MouseMsg:
		return GittiMouseInteraction(msg, m)
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
		}
	}

	var cmd tea.Cmd
	m.CurrentRepoBranchesInfoList, cmd = m.CurrentRepoBranchesInfoList.Update(msg)
	cmds = append(cmds, cmd)
	m.CurrentRepoModifiedFilesInfoList, cmd = m.CurrentRepoModifiedFilesInfoList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *GittiModel) View() tea.View {
	var v tea.View
	v.SetContent(gittiMainPageView(m))
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}
