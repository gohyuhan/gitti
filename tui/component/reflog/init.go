package reflog

import (
	"charm.land/bubbles/v2/list"
	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ----------------------------------
//
//	init the list component for Reflog Component
//	return bool was to tell if we need to reinit the detail component panel or not
//
// ----------------------------------
func InitGitRefLogList(m *types.GittiModel) bool {
	latestRefLogArray := m.GitOperations.GitRefLog.RefLog()
	items := make([]list.Item, 0, len(latestRefLogArray))

	// get the previous selected file and see if it was within the new list if yes get the latest position of the previous selected file
	previousSelectedRefLog := m.CurrentRepoRefLogInfoList.SelectedItem()
	selectedRefLogPosition := -1

	titleWidthLimit := m.WindowLeftPanelWidth - constant.ListItemOrTitleWidthPad - 2

	if previousSelectedRefLog != nil {
		for index, refLog := range latestRefLogArray {
			// we use reflog hash here to determine if it was the same file as the reflog hash is unique
			if refLog.Hash == previousSelectedRefLog.(GitRefLogItem).Hash {
				selectedRefLogPosition = index
			}
			items = append(items, GitRefLogItem(refLog))
		}
	} else {
		for _, refLog := range latestRefLogArray {
			items = append(items, GitRefLogItem(refLog))
		}
	}

	previousRefLogCount := len(m.CurrentRepoRefLogInfoList.Items())

	m.CurrentRepoRefLogInfoList = list.New(items, GitRefLogItemDelegate{}, m.WindowLeftPanelWidth, m.RefLogComponentPanelHeight)
	m.CurrentRepoRefLogInfoList.SetShowPagination(false)
	m.CurrentRepoRefLogInfoList.SetShowStatusBar(false)
	m.CurrentRepoRefLogInfoList.SetFilteringEnabled(false)
	m.CurrentRepoRefLogInfoList.SetShowFilter(false)
	m.CurrentRepoRefLogInfoList.Title = ansi.Truncate(ConstructRefLogComponentTitle(titleWidthLimit), titleWidthLimit, "...")
	m.CurrentRepoRefLogInfoList.Styles.Title = style.TitleStyle
	m.CurrentRepoRefLogInfoList.Styles.TitleBar = style.NewStyle
	m.CurrentRepoRefLogInfoList.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)

	// Custom Help Model for Count Display
	m.CurrentRepoRefLogInfoList.SetShowHelp(true)
	m.CurrentRepoRefLogInfoList.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	m.CurrentRepoRefLogInfoList.AdditionalShortHelpKeys = utils.ListCounterHelper(m, &m.CurrentRepoRefLogInfoList)

	if len(items) < 1 {
		return len(items) != previousRefLogCount
	}

	if selectedRefLogPosition >= 0 {
		m.CurrentRepoRefLogInfoList.Select(selectedRefLogPosition)
		m.ListNavigationIndexPosition.RefLogComponent = selectedRefLogPosition
	} else {
		if m.ListNavigationIndexPosition.RefLogComponent > len(m.CurrentRepoRefLogInfoList.Items())-1 {
			m.CurrentRepoRefLogInfoList.Select(len(m.CurrentRepoRefLogInfoList.Items()) - 1)
			m.ListNavigationIndexPosition.RefLogComponent = len(m.CurrentRepoRefLogInfoList.Items()) - 1
		} else {
			m.CurrentRepoRefLogInfoList.Select(m.ListNavigationIndexPosition.RefLogComponent)
		}
	}

	if previousSelectedRefLog == m.CurrentRepoRefLogInfoList.SelectedItem() {
		return false
	}
	return true
}
