package tui

import (
	"fmt"
	"io"
	"strings"

	"gitti/api"
	"gitti/api/git"

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
		CurrentSelectedContainer:              ModifiedFilesComponent,
		RepoPath:                              repoPath,
		Width:                                 0,
		Height:                                0,
		CurrentRepoBranchesInfoList:           list.New([]list.Item{}, gitBranchItemDelegate{}, 0, 0),
		CurrentRepoModifiedFilesInfoList:      list.New([]list.Item{}, gitModifiedFilesItemDelegate{}, 0, 0),
		CurrentSelectedFileDiffViewport:       vp,
		CurrentSelectedFileDiffViewportOffset: 0,
		NavigationIndexPosition:               GittiComponentsCurrentNavigationIndexPosition{LocalBranchComponent: 0, ModifiedFilesComponent: 0},
		PopUpType:                             NoPopUp,
		PopUpModel:                            struct{}{},
		GitState:                              gitState,
	}
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
	case tea.KeyMsg:
		return gittiKeyInteraction(msg, m)
	case GitUpdateMsg:
		updateEvent := string(msg)
		switch updateEvent {
		case git.GIT_BRANCH_UPDATE:
			initBranchList(m)
		case git.GIT_FILES_STATUS_UPDATE:
			initModifiedFilesList(m)
			renderModifiedFilesDiffViewPort(m)
		case git.GIT_COMMIT_OUTPUT_UPDATE:
			updatePopUpCommitOutputViewPort(m)
		case git.GIT_REMOTE_PUSH_OUTPUT_UPDATE:
			updateGitRemotePushOutputViewport(m)
		}
		return m, nil
	case tea.MouseMsg:
		return GittiMouseInteraction(msg, m)
	}

	// Update spinners in popups when they are processing
	if m.ShowPopUp.Load() {
		if commitPopup, ok := m.PopUpModel.(*GitCommitPopUpModel); ok && commitPopup.IsProcessing.Load() {
			var cmd tea.Cmd
			commitPopup.Spinner, cmd = commitPopup.Spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
		if pushPopup, ok := m.PopUpModel.(*GitRemotePushPopUpModel); ok && pushPopup.IsProcessing.Load() {
			var cmd tea.Cmd
			pushPopup.Spinner, cmd = pushPopup.Spinner.Update(msg)
			cmds = append(cmds, cmd)
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

// -----------------------------------------------------------------------------
// implementation for list compoenent
// -----------------------------------------------------------------------------
// to record the current navigation index position
func (d gitModifiedFilesItemDelegate) Height() int                             { return 1 }
func (d gitModifiedFilesItemDelegate) Spacing() int                            { return 0 }
func (d gitModifiedFilesItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gitModifiedFilesItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(gitModifiedFilesItem)
	if !ok {
		return
	}

	str := fmt.Sprintf(" [ ] %s", i.FileName)
	if i.SelectedForStage {
		str = fmt.Sprintf(" [X] %s", i.FileName)
	}

	componentWidth := m.Width() - listItemOrTitleWidthPad

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Width(m.Width() - 2).Reverse(true).Render("> " + strings.Join(s, " "))
		}
	}
	str = truncateString(str, componentWidth)

	fmt.Fprint(w, fn(str))
}

// for list component of git branch
func (d gitBranchItemDelegate) Height() int                             { return 1 }
func (d gitBranchItemDelegate) Spacing() int                            { return 0 }
func (d gitBranchItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gitBranchItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(gitBranchItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("   %s", i.BranchName)
	if i.IsCheckedOut {
		str = fmt.Sprintf(" * %s", i.BranchName)
	}

	componentWidth := m.Width() - listItemOrTitleWidthPad

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Width(m.Width() - 2).Reverse(true).Render("> " + strings.Join(s, " "))
		}
	}
	str = truncateString(str, componentWidth)

	fmt.Fprint(w, fn(str))
}

// for list component of git remote
func (d gitRemoteItemDelegate) Height() int                             { return 1 }
func (d gitRemoteItemDelegate) Spacing() int                            { return 0 }
func (d gitRemoteItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gitRemoteItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(gitRemoteItem)
	if !ok {
		return
	}

	nameStr := fmt.Sprintf("   %s", i.Name)
	urlStr := fmt.Sprintf("    %s", i.Url)

	componentWidth := m.Width() - listItemOrTitleWidthPad

	nameStr = truncateString(nameStr, componentWidth)
	urlStr = truncateString(urlStr, componentWidth)

	nameRendered := itemStyle.Render(nameStr)
	urlRendered := itemStyle.Faint(true).Render(urlStr)
	fullStr := nameRendered + "\n" + urlRendered

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Width(m.Width() - 2).Reverse(true).Render(strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(fullStr))
}

// for push selection option
func (d gitPushOptionDelegate) Height() int                             { return 1 }
func (d gitPushOptionDelegate) Spacing() int                            { return 0 }
func (d gitPushOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gitPushOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(gitPushOptionItem)
	if !ok {
		return
	}

	nameStr := fmt.Sprintf("   %s", i.Name)
	urlStr := fmt.Sprintf("    %s", i.Info)

	componentWidth := m.Width() - listItemOrTitleWidthPad

	nameStr = truncateString(nameStr, componentWidth)
	urlStr = truncateString(urlStr, componentWidth)

	nameRendered := itemStyle.Render(nameStr)
	urlRendered := itemStyle.Faint(true).Render(urlStr)
	fullStr := nameRendered + "\n" + urlRendered

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Width(m.Width() - 2).Reverse(true).Render(strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(fullStr))
}
