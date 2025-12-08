package files

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// those utf-8 icons for the component can be found at https://www.nerdfonts.com/cheat-sheet

// init the list component for Modified Files Component
// return bool was to tell if we need to reinit the detail component panel or not
func InitModifiedFilesList(m *types.GittiModel) bool {
	latestModifiedFilesArray := m.GitOperations.GitFiles.FilesStatus()
	items := make([]list.Item, 0, len(latestModifiedFilesArray))
	for _, modifiedFile := range latestModifiedFilesArray {
		items = append(items, GitModifiedFilesItem(modifiedFile))
	}

	// get the previous selected file and see if it was within the new list if yes get the latest position of the previous selected file
	previousSelectedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
	selectedFilesPosition := -1

	for index, item := range items {
		if item == previousSelectedFile {
			selectedFilesPosition = index
			break
		}
	}

	previousModifiedFilesCount := len(m.CurrentRepoModifiedFilesInfoList.Items())

	m.CurrentRepoModifiedFilesInfoList = list.New(items, GitModifiedFilesItemDelegate{}, m.WindowLeftPanelWidth, m.ModifiedFilesComponentPanelHeight)
	m.CurrentRepoModifiedFilesInfoList.SetShowPagination(false)
	m.CurrentRepoModifiedFilesInfoList.SetShowStatusBar(false)
	m.CurrentRepoModifiedFilesInfoList.SetFilteringEnabled(false)
	m.CurrentRepoModifiedFilesInfoList.SetShowFilter(false)
	m.CurrentRepoModifiedFilesInfoList.Title = utils.TruncateString(fmt.Sprintf("[2] \ueae9 %s:", i18n.LANGUAGEMAPPING.ModifiedFiles), m.WindowLeftPanelWidth-constant.ListItemOrTitleWidthPad-2)
	m.CurrentRepoModifiedFilesInfoList.Styles.Title = style.TitleStyle
	m.CurrentRepoModifiedFilesInfoList.Styles.TitleBar = style.NewStyle

	// Custom Help Model for Count Display
	m.CurrentRepoModifiedFilesInfoList.SetShowHelp(true)
	m.CurrentRepoModifiedFilesInfoList.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	m.CurrentRepoModifiedFilesInfoList.AdditionalShortHelpKeys = utils.ListCounterHelper(m, &m.CurrentRepoModifiedFilesInfoList)

	if len(items) < 1 {
		return len(items) != previousModifiedFilesCount
	}

	if selectedFilesPosition >= 0 {
		m.CurrentRepoModifiedFilesInfoList.Select(selectedFilesPosition)
		m.ListNavigationIndexPosition.ModifiedFilesComponent = selectedFilesPosition
	} else {
		if m.ListNavigationIndexPosition.ModifiedFilesComponent > len(m.CurrentRepoModifiedFilesInfoList.Items())-1 {
			m.CurrentRepoModifiedFilesInfoList.Select(len(m.CurrentRepoModifiedFilesInfoList.Items()) - 1)
			m.ListNavigationIndexPosition.ModifiedFilesComponent = len(m.CurrentRepoModifiedFilesInfoList.Items()) - 1
		} else {
			m.CurrentRepoModifiedFilesInfoList.Select(m.ListNavigationIndexPosition.ModifiedFilesComponent)
		}
	}

	if previousSelectedFile == m.CurrentRepoModifiedFilesInfoList.SelectedItem() {
		return false
	}
	return true
}
