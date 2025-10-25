package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"gitti/api/git"
	"io"
	"strings"
)

func NewGittiModel(repoPath string) GittiModel {
	gitti := GittiModel{
		CurrentSelectedContainer:              ModifiedFilesComponent,
		RepoPath:                              repoPath,
		Width:                                 0,
		Height:                                0,
		CurrentRepoBranchesInfo:               list.New([]list.Item{}, itemStringDelegate{}, 0, 0),
		CurrentRepoModifiedFilesInfo:          list.New([]list.Item{}, itemStringDelegate{}, 0, 0),
		CurrentSelectedFileDiffViewport:       viewport.New(0, 0),
		CurrentSelectedFileDiffViewportOffset: 0,
		NavigationIndexPosition:               GittiComponentsCurrentNavigationIndexPosition{LocalBranchComponent: 0, ModifiedFilesComponent: 0},
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

		// Compute panel widths
		m.HomeTabLeftPanelWidth = int(float64(m.Width) * mainPageLayoutLeftPanelWidthRatio)
		m.HomeTabFileDiffPanelWidth = m.Width - m.HomeTabLeftPanelWidth - 4 // adjust for borders/padding

		m.HomeTabCoreContentHeight = m.Height - mainPageKeyBindingLayoutPanelHeight - 2*padding
		m.HomeTabFileDiffPanelHeight = m.HomeTabCoreContentHeight
		m.HomeTabLocalBranchesPanelHeight = int(float64(m.HomeTabCoreContentHeight)*mainPageLocalBranchesPanelHeightRatio) - 2*padding
		m.HomeTabChangedFilesPanelHeight = m.HomeTabCoreContentHeight - m.HomeTabLocalBranchesPanelHeight - 2*padding

		// update all components Width and Height
		m.CurrentRepoBranchesInfo.SetWidth(m.HomeTabLeftPanelWidth)
		m.CurrentRepoBranchesInfo.SetHeight(m.HomeTabLocalBranchesPanelHeight)

		// update viewport
		m.CurrentSelectedFileDiffViewport.Height = m.HomeTabFileDiffPanelHeight - 1 //some margin
		m.CurrentSelectedFileDiffViewport.Width = m.HomeTabFileDiffPanelWidth
		m.CurrentSelectedFileDiffViewportOffset = max(0, int(m.CurrentSelectedFileDiffViewport.HorizontalScrollPercent()*float64(m.CurrentSelectedFileDiffViewportOffset))-1)
		m.CurrentSelectedFileDiffViewport.SetXOffset(m.CurrentSelectedFileDiffViewportOffset)

	case tea.KeyMsg:
		return GittiKeyInteraction(msg, m)
	case GitUpdateMsg:
		updateEvent := string(msg)
		switch updateEvent {
		case git.GIT_FILES_STATUS_UPDATE:
			InitModifiedFilesList(m)
		default:
			ProcessGitUpdate(m)
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.CurrentRepoBranchesInfo, cmd = m.CurrentRepoBranchesInfo.Update(msg)
	m.CurrentRepoModifiedFilesInfo, cmd = m.CurrentRepoModifiedFilesInfo.Update(msg)
	if m.CurrentSelectedContainer == FileDiffComponent {
		m.CurrentSelectedFileDiffViewport, cmd = m.CurrentSelectedFileDiffViewport.Update(msg)
	}
	return m, cmd
}

func (m *GittiModel) View() string {
	return GittiMainPageView(m)
}

// -----------------------------------------------------------------------------
// implementation for list compoenent
// -----------------------------------------------------------------------------

func (i itemString) FilterValue() string                             { return "" }
func (d itemStringDelegate) Height() int                             { return 1 }
func (d itemStringDelegate) Spacing() int                            { return 0 }
func (d itemStringDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemStringDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(itemString)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
