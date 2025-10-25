package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"gitti/api/git"
)

// variables for indicating which panel/components/container or whatever the hell you wanna call it that the user is currently landed or selected, so that they can do precious action related to the part of whatever the hell you wanna call it
var (
	None = "0"

	localBranchComponent  = "B1"
	filesChangesComponent = "B2"
	fileDiffComponent     = "B3"
)

func ProcessGitUpdate(m *GittiModel) {
	InitBranchList(m)
	return
}

func InitBranchList(m *GittiModel) {
	items := []list.Item{
		item(fmt.Sprintf("* %s", git.GITBRANCH.CurrentCheckOut)),
	}

	for _, branch := range git.GITBRANCH.AllBranches {
		items = append(items, item(fmt.Sprintf("  %s", branch)))
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
