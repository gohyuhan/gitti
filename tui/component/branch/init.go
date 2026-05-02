package branch

import (
	"charm.land/bubbles/v2/list"
	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	init the list component for Branch Component
//
// ------------------------------------
func InitBranchList(m *types.GittiModel) {
	currentCheckOut := m.GitOperations.GitBranch.CurrentCheckOut()
	latestBranchArray := []list.Item{
		GitBranchItem(currentCheckOut),
	}

	m.CheckOutBranch = currentCheckOut.BranchName

	previousSelectedBranch := m.CurrentRepoBranchesInfoList.SelectedItem()
	selectedBranchPosition := -1

	titleWidthLimit := m.WindowLeftPanelWidth - constant.ListItemOrTitleWidthPad - 2

	if previousSelectedBranch != nil {
		previousSelectedBranchInfo := previousSelectedBranch.(GitBranchItem)
		for index, branch := range m.GitOperations.GitBranch.AllBranches() {
			// we use branch name here to determine if it was the same branch as the branch name is unique
			// we need to +1 to the index as there will always be 1 item already in the list which is the checkout branch
			if branch.BranchName == previousSelectedBranchInfo.BranchName {
				selectedBranchPosition = index + 1
			}
			latestBranchArray = append(latestBranchArray, GitBranchItem(branch))
		}

		// if it previous selected branch name is the current checkout one, the position will be 0, as the checkout branch will always be the first in the list
		if currentCheckOut.BranchName == previousSelectedBranchInfo.BranchName {
			selectedBranchPosition = 0
		}
	} else {
		for _, branch := range m.GitOperations.GitBranch.AllBranches() {
			latestBranchArray = append(latestBranchArray, GitBranchItem(branch))
		}
	}

	m.CurrentRepoBranchesInfoList = list.New(latestBranchArray, GitBranchItemDelegate{}, m.WindowLeftPanelWidth, m.LocalBranchesComponentPanelHeight)
	m.CurrentRepoBranchesInfoList.SetShowPagination(false)
	m.CurrentRepoBranchesInfoList.SetShowStatusBar(false)
	m.CurrentRepoBranchesInfoList.SetFilteringEnabled(false)
	m.CurrentRepoBranchesInfoList.SetShowFilter(false)

	m.CurrentRepoBranchesInfoList.Title = ansi.Truncate(ConstructLocalBranchComponentTitle(titleWidthLimit), titleWidthLimit, "...")
	m.CurrentRepoBranchesInfoList.Styles.Title = style.TitleStyle
	m.CurrentRepoBranchesInfoList.Styles.PaginationStyle = style.PaginationStyle
	m.CurrentRepoBranchesInfoList.Styles.TitleBar = style.NewStyle
	m.CurrentRepoBranchesInfoList.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)

	// Custom Help Model for Count Display
	m.CurrentRepoBranchesInfoList.SetShowHelp(true)
	m.CurrentRepoBranchesInfoList.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	m.CurrentRepoBranchesInfoList.AdditionalShortHelpKeys = utils.ListCounterHelper(m, &m.CurrentRepoBranchesInfoList)

	if selectedBranchPosition >= 0 {
		m.CurrentRepoBranchesInfoList.Select(selectedBranchPosition)
		m.ListNavigationIndexPosition.LocalBranchComponent = selectedBranchPosition
	} else {
		if m.ListNavigationIndexPosition.LocalBranchComponent > len(m.CurrentRepoBranchesInfoList.Items())-1 {
			m.CurrentRepoBranchesInfoList.Select(len(m.CurrentRepoBranchesInfoList.Items()) - 1)
			m.ListNavigationIndexPosition.LocalBranchComponent = len(m.CurrentRepoBranchesInfoList.Items()) - 1
		} else {
			m.CurrentRepoBranchesInfoList.Select(m.ListNavigationIndexPosition.LocalBranchComponent)
		}
	}
}
