package stash

import (
	"fmt"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"

	"charm.land/bubbles/v2/list"
)

// those utf-8 icons for the component can be found at https://www.nerdfonts.com/cheat-sheet

// init the list component for Stash info Component
// return bool was to tell if we need to reinit the detail component panel or not
func InitStashList(m *types.GittiModel) bool {
	latestStashArray := m.GitOperations.GitStash.AllStash()
	items := make([]list.Item, 0, len(latestStashArray))
	for _, stashInfo := range latestStashArray {
		items = append(items, GitStashItem(stashInfo))
	}

	// get the previous selected file and see if it was within the new list if yes get the latest position of the previous selected file
	previousSelectedStash := m.CurrentRepoStashInfoList.SelectedItem()
	selectedFilesPosition := -1

	for index, item := range items {
		if item == previousSelectedStash {
			selectedFilesPosition = index
			break
		}
	}
	previousStashCount := len(m.CurrentRepoStashInfoList.Items())

	m.CurrentRepoStashInfoList = list.New(items, GitStashItemDelegate{}, m.WindowLeftPanelWidth, m.StashComponentPanelHeight)
	m.CurrentRepoStashInfoList.SetShowPagination(false)
	m.CurrentRepoStashInfoList.SetShowStatusBar(false)
	m.CurrentRepoStashInfoList.SetFilteringEnabled(false)
	m.CurrentRepoStashInfoList.SetShowFilter(false)
	m.CurrentRepoStashInfoList.Title = utils.TruncateString(fmt.Sprintf("[4] \ueaf7 %s:", i18n.LANGUAGEMAPPING.Stash), m.WindowLeftPanelWidth-constant.ListItemOrTitleWidthPad-2)
	m.CurrentRepoStashInfoList.Styles.Title = style.TitleStyle
	m.CurrentRepoStashInfoList.Styles.TitleBar = style.NewStyle
	m.CurrentRepoStashInfoList.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)

	// Custom Help Model for Count Display
	m.CurrentRepoStashInfoList.SetShowHelp(true)
	m.CurrentRepoStashInfoList.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	m.CurrentRepoStashInfoList.AdditionalShortHelpKeys = utils.ListCounterHelper(m, &m.CurrentRepoStashInfoList)

	if len(items) < 1 {
		return len(items) != previousStashCount
	}

	if selectedFilesPosition >= 0 {
		m.CurrentRepoStashInfoList.Select(selectedFilesPosition)
		m.ListNavigationIndexPosition.StashComponent = selectedFilesPosition
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
