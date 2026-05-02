package tag

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
func InitTagList(m *types.GittiModel) bool {
	latestTagArray := []list.Item{}

	previousSelectedTag := m.CurrentRepoTagInfoList.SelectedItem()
	selectedTagPosition := -1

	titleWidthLimit := m.WindowLeftPanelWidth - constant.ListItemOrTitleWidthPad - 2

	if previousSelectedTag != nil {
		previousSelectedTagInfo := previousSelectedTag.(GitTagItem)
		for index, tag := range m.GitOperations.GitTag.AllTag() {
			// we use branch name here to determine if it was the same branch as the branch name is unique
			// we need to +1 to the index as there will always be 1 item already in the list which is the checkout branch
			if tag.TagName == previousSelectedTagInfo.TagName {
				selectedTagPosition = index
			}
			latestTagArray = append(latestTagArray, GitTagItem(tag))
		}
	} else {
		for _, tag := range m.GitOperations.GitTag.AllTag() {
			latestTagArray = append(latestTagArray, GitTagItem(tag))
		}
	}

	previousTagsCount := len(m.CurrentRepoTagInfoList.Items())

	m.CurrentRepoTagInfoList = list.New(latestTagArray, GitTagItemDelegate{}, m.WindowLeftPanelWidth, m.LocalBranchesComponentPanelHeight)
	m.CurrentRepoTagInfoList.SetShowPagination(false)
	m.CurrentRepoTagInfoList.SetShowStatusBar(false)
	m.CurrentRepoTagInfoList.SetFilteringEnabled(false)
	m.CurrentRepoTagInfoList.SetShowFilter(false)
	m.CurrentRepoTagInfoList.Title = ansi.Truncate(ConstructTagComponentTitle(titleWidthLimit), titleWidthLimit, "...")
	m.CurrentRepoTagInfoList.Styles.Title = style.TitleStyle
	m.CurrentRepoTagInfoList.Styles.PaginationStyle = style.PaginationStyle
	m.CurrentRepoTagInfoList.Styles.TitleBar = style.NewStyle
	m.CurrentRepoTagInfoList.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)

	// Custom Help Model for Count Display
	m.CurrentRepoTagInfoList.SetShowHelp(true)
	m.CurrentRepoTagInfoList.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	m.CurrentRepoTagInfoList.AdditionalShortHelpKeys = utils.ListCounterHelper(m, &m.CurrentRepoTagInfoList)

	if len(latestTagArray) < 1 {
		return len(latestTagArray) != previousTagsCount
	}

	if selectedTagPosition >= 0 {
		m.CurrentRepoTagInfoList.Select(selectedTagPosition)
		m.ListNavigationIndexPosition.TagComponent = selectedTagPosition
	} else {
		if m.ListNavigationIndexPosition.TagComponent > len(m.CurrentRepoTagInfoList.Items())-1 {
			m.CurrentRepoTagInfoList.Select(len(m.CurrentRepoTagInfoList.Items()) - 1)
			m.ListNavigationIndexPosition.TagComponent = len(m.CurrentRepoTagInfoList.Items()) - 1
		} else {
			m.CurrentRepoTagInfoList.Select(m.ListNavigationIndexPosition.TagComponent)
		}
	}

	if previousSelectedTag == m.CurrentRepoTagInfoList.SelectedItem() {
		return false
	}

	return true
}
