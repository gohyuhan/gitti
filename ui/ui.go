package ui

import (
	"fmt"
	"gitti/api/git"
	"gitti/types"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// implementation for list compoenent
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

func NewGittiModel() GittiModel {
	repoPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Failed to get current working directory: %v", err))
	}
	gitti := GittiModel{
		RepoPath:                          repoPath,
		Width:                             0,
		Height:                            0,
		AllRepoBranches:                   []types.BranchesInfo{},
		CurrentRepoBranchesInfo:           list.New([]list.Item{}, itemDelegate{}, 0, 0),
		CurrentCheckedOutBranch:           "",
		CurrentSelectedFiles:              "",
		CurrentSelectedFilesIndexPosition: 0,
		RemoteOrigin:                      "",
		UserName:                          "",
		UserEmail:                         "",
	}

	isSuccess, statusCode, getGitInfoErr, gitInfo := git.GetInitialGitInfo(
		gitti.RepoPath,
	)

	gitti.CurrentCheckedOutBranch = gitInfo.CurrentCheckedOutBranch
	gitti.AllRepoBranches = gitInfo.AllBranches
	gitti.CurrentSelectedFiles = gitInfo.CurrentSelectedFile
	gitti.AllChangedFiles = gitInfo.AllChangedFiles

	if !isSuccess && getGitInfoErr != nil {
		panic(fmt.Sprintf("[%v], %s", statusCode, getGitInfoErr.Error()))
	}

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
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "Q", "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.CurrentRepoBranchesInfo, cmd = m.CurrentRepoBranchesInfo.Update(msg)
	return m, cmd
	return m, nil
}

func (m GittiModel) View() string {
	return GittiMainPageView(m)
}
