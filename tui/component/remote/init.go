package remote

import (
	"charm.land/bubbles/v2/list"
	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// those utf-8 icons for the component can be found at https://www.nerdfonts.com/cheat-sheet

// init the list component for Branch Component
func InitRemoteList(m *types.GittiModel) bool {
	latestRemoteArray := []list.Item{}

	previousSelectedRemote := m.CurrentRepoRemoteInfoList.SelectedItem()
	selectedRemotePosition := -1

	titleWidthLimit := m.WindowLeftPanelWidth - constant.ListItemOrTitleWidthPad - 2

	if previousSelectedRemote != nil {
		previousSelectedRemoteInfo := previousSelectedRemote.(GitRemoteItem)
		for index, remote := range m.GitOperations.GitRemote.Remote() {
			// we use branch name here to determine if it was the same branch as the branch name is unique
			// we need to +1 to the index as there will always be 1 item already in the list which is the checkout branch
			if remote.Name == previousSelectedRemoteInfo.Name {
				selectedRemotePosition = index
			}
			latestRemoteArray = append(latestRemoteArray, GitRemoteItem(remote))
		}
	} else {
		for _, remote := range m.GitOperations.GitRemote.Remote() {
			latestRemoteArray = append(latestRemoteArray, GitRemoteItem(remote))
		}
	}

	previousRemoteCount := len(m.CurrentRepoRemoteInfoList.Items())

	m.CurrentRepoRemoteInfoList = list.New(latestRemoteArray, GitRemoteItemDelegate{}, m.WindowLeftPanelWidth, m.LocalBranchesComponentPanelHeight)
	m.CurrentRepoRemoteInfoList.SetShowPagination(false)
	m.CurrentRepoRemoteInfoList.SetShowStatusBar(false)
	m.CurrentRepoRemoteInfoList.SetFilteringEnabled(false)
	m.CurrentRepoRemoteInfoList.SetShowFilter(false)
	m.CurrentRepoRemoteInfoList.Title = ansi.Truncate(ConstructRemoteComponentTitle(titleWidthLimit), titleWidthLimit, "...")
	m.CurrentRepoRemoteInfoList.Styles.Title = style.TitleStyle
	m.CurrentRepoRemoteInfoList.Styles.PaginationStyle = style.PaginationStyle
	m.CurrentRepoRemoteInfoList.Styles.TitleBar = style.NewStyle
	m.CurrentRepoRemoteInfoList.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)

	// Custom Help Model for Count Display
	m.CurrentRepoRemoteInfoList.SetShowHelp(true)
	m.CurrentRepoRemoteInfoList.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	m.CurrentRepoRemoteInfoList.AdditionalShortHelpKeys = utils.ListCounterHelper(m, &m.CurrentRepoRemoteInfoList)

	if len(latestRemoteArray) < 1 {
		return len(latestRemoteArray) != previousRemoteCount
	}

	if selectedRemotePosition >= 0 {
		m.CurrentRepoRemoteInfoList.Select(selectedRemotePosition)
		m.ListNavigationIndexPosition.RemoteComponent = selectedRemotePosition
	} else {
		if m.ListNavigationIndexPosition.RemoteComponent > len(m.CurrentRepoRemoteInfoList.Items())-1 {
			m.CurrentRepoRemoteInfoList.Select(len(m.CurrentRepoRemoteInfoList.Items()) - 1)
			m.ListNavigationIndexPosition.RemoteComponent = len(m.CurrentRepoRemoteInfoList.Items()) - 1
		} else {
			m.CurrentRepoRemoteInfoList.Select(m.ListNavigationIndexPosition.RemoteComponent)
		}
	}

	if previousSelectedRemote == m.CurrentRepoRemoteInfoList.SelectedItem() {
		return false
	}

	return true
}
