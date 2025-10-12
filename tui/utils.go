package tui

import (
	"fmt"
	"gitti/types"
	"sort"

	"github.com/charmbracelet/bubbles/list"
)

// variables for indicating which panel/components/container or whatever the hell you wanna call it that the user is currently landed or selected, so that they can do precious action related to the part of whatever the hell you wanna call it
var (
	None = "0"

	localBranchComponent  = "B1"
	filesChangesComponent = "B2"
	fileDiffComponent     = "B3"
)

// this is for tab ( there will be 4 tab for now, initialzation tab(only accesible when user's repo was not git initialized yet), home tab, commit logs tab, about gitti tab )
var (
	initializationTab = "A"
	homeTab           = "B"
	commitLogsTab     = "C"
	aboutGittiTab     = "D"
)

func ProcessGitUpdate(m *GittiModel, gitInfo types.GitInfo) {
	m.CurrentCheckedOutBranch = gitInfo.CurrentCheckedOutBranch
	m.CurrentSelectedFiles = gitInfo.CurrentSelectedFile
	m.AllRepoBranches = gitInfo.AllBranches
	m.AllChangedFiles = gitInfo.AllChangedFiles

	InitBranchList(m)
	return
}

func InitBranchList(m *GittiModel) {
	keys := make([]string, 0, len(m.AllRepoBranches))
	for k := range m.AllRepoBranches {
		keys = append(keys, k)
	}
	sort.Strings(keys) // alphabetically sort keys

	items := []list.Item{
		item(fmt.Sprintf("* %s", m.CurrentCheckedOutBranch)),
	}

	for _, k := range keys {
		branch := m.AllRepoBranches[k]
		if !branch.CurrentCheckout {
			items = append(items, item(branch.Name))
		}
	}

	m.CurrentRepoBranchesInfo = list.New(items, itemDelegate{}, m.HomeTabLeftPanelWidth, m.HomeTabLocalBranchesPanelHeight)
	m.CurrentRepoBranchesInfo.Title = "[b]  Branches:"
	m.CurrentRepoBranchesInfo.SetShowStatusBar(false)
	m.CurrentRepoBranchesInfo.SetFilteringEnabled(false)
	m.CurrentRepoBranchesInfo.SetShowHelp(false)
	m.CurrentRepoBranchesInfo.Styles.Title = titleStyle
	m.CurrentRepoBranchesInfo.Styles.PaginationStyle = paginationStyle

	if m.NavigationIndexPosition.LocalBranchComponent > len(m.CurrentRepoBranchesInfo.Items())-1 {
		m.CurrentRepoBranchesInfo.Select(len(m.CurrentRepoBranchesInfo.Items()) - 1)
	} else {
		m.CurrentRepoBranchesInfo.Select(m.NavigationIndexPosition.LocalBranchComponent)
	}

	return
}
