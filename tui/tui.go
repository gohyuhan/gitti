package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"gitti/api/git"
	"io"
	"strings"
)

func NewGittiModel(repoPath string) GittiModel {
	vp := viewport.New()
	vp.SoftWrap = false
	vp.MouseWheelEnabled = true
	vp.SetHorizontalStep(1)
	vp.MouseWheelDelta = 1
	gitti := GittiModel{
		CurrentSelectedContainer:              ModifiedFilesComponent,
		RepoPath:                              repoPath,
		Width:                                 0,
		Height:                                0,
		CurrentRepoBranchesInfoList:           list.New([]list.Item{}, gitBranchItemDelegate{}, 0, 0),
		CurrentRepoModifiedFilesInfoList:      list.New([]list.Item{}, gitModifiedFilesItemDelegate{}, 0, 0),
		CurrentSelectedFileDiffViewport:       vp,
		CurrentSelectedFileDiffViewportOffset: 0,
		NavigationIndexPosition:               GittiComponentsCurrentNavigationIndexPosition{LocalBranchComponent: 0, ModifiedFilesComponent: 0},
		ShowPopUp:                             false,
		PopUpType:                             None,
		PopUpModel:                            struct{}{},
		IsTyping:                              false,
	}

	return gitti
}

// -----------------------------------------------------------------------------
// Bubble Tea standard functions
// -----------------------------------------------------------------------------

func (m *GittiModel) Init() tea.Cmd {
	return nil
}

func (m *GittiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case git.GIT_FILES_STATUS_UPDATE:
			initModifiedFilesList(m)
		case git.GIT_COMMIT_OUTPUT_UPDATE:
			updatePopUpCommitOutputViewPort(m)
		default:
			processGitUpdate(m)
		}
		return m, nil
	case tea.MouseMsg:
		return GittiMouseInteraction(msg, m)
	}

	var cmd tea.Cmd
	m.CurrentRepoBranchesInfoList, cmd = m.CurrentRepoBranchesInfoList.Update(msg)
	m.CurrentRepoModifiedFilesInfoList, cmd = m.CurrentRepoModifiedFilesInfoList.Update(msg)

	return m, cmd
}

func (m *GittiModel) View() tea.View {
	var v tea.View
	v.SetContent(themeStyle.Render(gittiMainPageView(m)))
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

// -----------------------------------------------------------------------------
// implementation for list compoenent
// -----------------------------------------------------------------------------
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

	componentWidth := m.Width() - 5
	str = truncateString(str, componentWidth)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

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

	componentWidth := m.Width() - 5
	str = truncateString(str, componentWidth)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
