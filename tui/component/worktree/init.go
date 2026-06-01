package worktree

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
//	Rebuild the worktree list widget from the latest git worktree data, preserve the
//	previously selected worktree by path, and return true if the selection changed
//	(signals that the detail panel needs to be reinitialized).
//
// ------------------------------------
func InitWorktreeList(m *types.GittiModel) bool {
	latestWorktreeArray := []list.Item{}

	previousSelectedWorktree := m.CurrentRepoWorktreeInfoList.SelectedItem()
	selectedWorktreePosition := -1

	titleWidthLimit := m.WindowLeftPanelWidth - constant.ListItemOrTitleWidthPad - 2

	if previousSelectedWorktree != nil {
		previousSelectedWorktreeInfo := previousSelectedWorktree.(GitWorktreeItem)
		for index, worktree := range m.GitOperations.GitWorktree.AllWorktree() {
			if worktree.WorktreePath == previousSelectedWorktreeInfo.WorktreePath {
				selectedWorktreePosition = index
			}
			latestWorktreeArray = append(latestWorktreeArray, GitWorktreeItem(worktree))
		}
	} else {
		for _, worktree := range m.GitOperations.GitWorktree.AllWorktree() {
			latestWorktreeArray = append(latestWorktreeArray, GitWorktreeItem(worktree))
		}
	}

	previousWorktreeCount := len(m.CurrentRepoWorktreeInfoList.Items())

	m.CurrentRepoWorktreeInfoList = list.New(latestWorktreeArray, GitWorktreeItemDelegate{}, m.WindowLeftPanelWidth, m.LocalBranchesComponentPanelHeight)
	m.CurrentRepoWorktreeInfoList.SetShowPagination(false)
	m.CurrentRepoWorktreeInfoList.SetShowStatusBar(false)
	m.CurrentRepoWorktreeInfoList.SetFilteringEnabled(false)
	m.CurrentRepoWorktreeInfoList.SetShowFilter(false)
	m.CurrentRepoWorktreeInfoList.Title = ansi.Truncate(ConstructWorktreeComponentTitle(titleWidthLimit), titleWidthLimit, "...")
	m.CurrentRepoWorktreeInfoList.Styles.Title = style.TitleStyle
	m.CurrentRepoWorktreeInfoList.Styles.PaginationStyle = style.PaginationStyle
	m.CurrentRepoWorktreeInfoList.Styles.TitleBar = style.NewStyle
	m.CurrentRepoWorktreeInfoList.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)

	// Custom Help Model for Count Display
	m.CurrentRepoWorktreeInfoList.SetShowHelp(true)
	m.CurrentRepoWorktreeInfoList.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	m.CurrentRepoWorktreeInfoList.AdditionalShortHelpKeys = utils.ListCounterHelper(m, &m.CurrentRepoWorktreeInfoList)

	if len(latestWorktreeArray) < 1 {
		return len(latestWorktreeArray) != previousWorktreeCount
	}

	if selectedWorktreePosition >= 0 {
		m.CurrentRepoWorktreeInfoList.Select(selectedWorktreePosition)
		m.ListNavigationIndexPosition.WorktreeComponent = selectedWorktreePosition
	} else {
		if m.ListNavigationIndexPosition.WorktreeComponent > len(m.CurrentRepoWorktreeInfoList.Items())-1 {
			m.CurrentRepoWorktreeInfoList.Select(len(m.CurrentRepoWorktreeInfoList.Items()) - 1)
			m.ListNavigationIndexPosition.WorktreeComponent = len(m.CurrentRepoWorktreeInfoList.Items()) - 1
		} else {
			m.CurrentRepoWorktreeInfoList.Select(m.ListNavigationIndexPosition.WorktreeComponent)
		}
	}

	if previousSelectedWorktree == m.CurrentRepoWorktreeInfoList.SelectedItem() {
		return false
	}

	return true
}
