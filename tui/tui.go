package tui

import (
	"fmt"
	"gitti/api"
	"gitti/api/git"
	"gitti/types"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func NewGittiModel(repoPath string, gitWorkerDaemon *api.GittiDaemonWorker) GittiModel {
	gitti := GittiModel{
		CurrentTab:                        homeTab,
		CurrentSelectedContainer:          None,
		RepoPath:                          repoPath,
		Width:                             0,
		Height:                            0,
		AllRepoBranches:                   map[string]types.BranchesInfo{},
		CurrentRepoBranchesInfo:           list.New([]list.Item{}, itemDelegate{}, 0, 0),
		CurrentCheckedOutBranch:           "",
		CurrentSelectedFiles:              "",
		CurrentSelectedFilesIndexPosition: 0,
		RemoteOrigin:                      "",
		UserName:                          "",
		UserEmail:                         "",
		GitWorkerDaemon:                   *gitWorkerDaemon,
		NavigationIndexPosition:           GittiComponentsCurrentNavigationIndexPosition{LocalBranchComponent: 0, FilesChangesComponent: 0},
	}

	isSuccess, statusCode, getGitInfoErr, gitInfo := git.GetGitInfo(
		gitti.RepoPath,
	)

	gitti.CurrentCheckedOutBranch = gitInfo.CurrentCheckedOutBranch
	gitti.AllRepoBranches = gitInfo.AllBranches
	gitti.CurrentSelectedFiles = gitInfo.CurrentSelectedFile
	gitti.AllChangedFiles = gitInfo.AllChangedFiles

	if !isSuccess && getGitInfoErr != nil {
		panic(fmt.Sprintf("[%v], %s", statusCode, getGitInfoErr.Error()))
	}

	InitBranchList(&gitti)

	gitti.GitWorkerDaemon.Start()

	return gitti
}

// -----------------------------------------------------------------------------
// Bubble Tea standard functions
// -----------------------------------------------------------------------------

func (m GittiModel) Init() tea.Cmd {
	return nil
}

func (m GittiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		// Compute panel widths
		m.HomeTabLeftPanelWidth = int(float64(m.Width) * mainPageLayoutLeftPanelWidthRatio)
		m.HomeTabFileDiffPanelWidth = m.Width - m.HomeTabLeftPanelWidth - 4 // adjust for borders/padding

		m.HomeTabCoreContentHeight = m.Height - mainPageLayoutTitlePanelHeight - padding - mainPageKeyBindingLayoutPanelHeight - padding
		m.HomeTabFileDiffPanelHeight = m.HomeTabCoreContentHeight
		m.HomeTabLocalBranchesPanelHeight = int(float64(m.HomeTabCoreContentHeight)*mainPageLocalBranchesPanelHeightRatio) - padding
		m.HomeTabChangedFilesPanelHeight = int(float64(m.HomeTabCoreContentHeight) * mainPageChangedFilesHeightRatio)

		// update all components Width and Height
		m.CurrentRepoBranchesInfo.SetWidth(m.HomeTabLeftPanelWidth)
		m.CurrentRepoBranchesInfo.SetHeight(m.HomeTabLocalBranchesPanelHeight)
	case tea.KeyMsg:
		return GittiKeyInteraction(msg, m)
	case GitUpdateMsg:
		ProcessGitUpdate(&m, types.GitInfo(msg))
		return m, nil
	}

	var cmd tea.Cmd
	m.CurrentRepoBranchesInfo, cmd = m.CurrentRepoBranchesInfo.Update(msg)
	return m, cmd
}

func (m GittiModel) View() string {
	return GittiMainPageView(m)
}

// -----------------------------------------------------------------------------
// implementation for list compoenent
// -----------------------------------------------------------------------------

func (i item) FilterValue() string                             { return "" }
func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
