package stash

import (
	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"

	"charm.land/bubbles/v2/list"
)

// ------------------------------------
//
//	Rebuild the stash list widget from the latest git stash data, preserve the previously
//	selected stash entry by ID, and return true if the selection changed (signals that the
//	detail panel needs to be reinitialized).
//
// ------------------------------------
func InitStashList(m *types.GittiModel) bool {
	latestStashArray := m.GitOperations.GitStash.AllStash()
	items := make([]list.Item, 0, len(latestStashArray))

	// get the previous selected stash and see if it was within the new list if yes get the latest position of the previous selected stash
	previousSelectedStash := m.CurrentRepoStashInfoList.SelectedItem()
	selectedStashPosition := -1

	titleWidthLimit := m.WindowLeftPanelWidth - constant.ListItemOrTitleWidthPad - 2

	if previousSelectedStash != nil {
		for index, stashInfo := range latestStashArray {
			if stashInfo.Id == previousSelectedStash.(GitStashItem).Id {
				selectedStashPosition = index
			}
			items = append(items, GitStashItem(stashInfo))
		}
	} else {
		for _, stashInfo := range latestStashArray {
			items = append(items, GitStashItem(stashInfo))
		}
	}
	items, selectedStashPosition = utils.FilterListItems(items, m.PanelFilterQuery[constant.StashComponentPanel], previousSelectedStash, selectedStashPosition)

	previousStashCount := len(m.CurrentRepoStashInfoList.Items())

	m.CurrentRepoStashInfoList = list.New(items, GitStashItemDelegate{}, m.WindowLeftPanelWidth, m.StashComponentPanelHeight)
	m.CurrentRepoStashInfoList.SetShowPagination(false)
	m.CurrentRepoStashInfoList.SetShowStatusBar(false)
	m.CurrentRepoStashInfoList.SetFilteringEnabled(false)
	m.CurrentRepoStashInfoList.SetShowFilter(false)
	m.CurrentRepoStashInfoList.Title = ansi.Truncate(ConstructStashComponentTitle(titleWidthLimit), titleWidthLimit, "...")
	m.CurrentRepoStashInfoList.Styles.Title = style.TitleStyle
	m.CurrentRepoStashInfoList.Styles.TitleBar = style.NewStyle
	m.CurrentRepoStashInfoList.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)

	// Custom Help Model for Count Display
	m.CurrentRepoStashInfoList.SetShowHelp(true)
	m.CurrentRepoStashInfoList.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	m.CurrentRepoStashInfoList.AdditionalShortHelpKeys = utils.ListCounterHelper(m, &m.CurrentRepoStashInfoList, constant.StashComponentPanel)

	if len(items) < 1 {
		return len(items) != previousStashCount
	}

	if selectedStashPosition >= 0 {
		m.CurrentRepoStashInfoList.Select(selectedStashPosition)
		m.ListNavigationIndexPosition.StashComponent = selectedStashPosition
	} else {
		if m.ListNavigationIndexPosition.StashComponent > len(m.CurrentRepoStashInfoList.Items())-1 {
			m.CurrentRepoStashInfoList.Select(len(m.CurrentRepoStashInfoList.Items()) - 1)
			m.ListNavigationIndexPosition.StashComponent = len(m.CurrentRepoStashInfoList.Items()) - 1
		} else {
			m.CurrentRepoStashInfoList.Select(m.ListNavigationIndexPosition.StashComponent)
		}
	}

	if previousSelectedStash == m.CurrentRepoStashInfoList.SelectedItem() {
		return false
	}
	return true
}
