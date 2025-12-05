package branch

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// init the list component for Branch Component
func InitBranchList(m *types.GittiModel) {
	currentCheckOut := m.GitOperations.GitBranch.CurrentCheckOut()
	latestBranchArray := []list.Item{
		GitBranchItem(currentCheckOut),
	}

	m.CheckOutBranch = currentCheckOut.BranchName

	for _, branch := range m.GitOperations.GitBranch.AllBranches() {
		latestBranchArray = append(latestBranchArray, GitBranchItem(branch))
	}

	m.CurrentRepoBranchesInfoList = list.New(latestBranchArray, GitBranchItemDelegate{}, m.WindowLeftPanelWidth, m.LocalBranchesComponentPanelHeight)
	m.CurrentRepoBranchesInfoList.SetShowPagination(false)
	m.CurrentRepoBranchesInfoList.SetShowStatusBar(false)
	m.CurrentRepoBranchesInfoList.SetFilteringEnabled(false)
	m.CurrentRepoBranchesInfoList.SetShowFilter(false)
	m.CurrentRepoBranchesInfoList.Title = utils.TruncateString(fmt.Sprintf("[1] \uf418 %s:", i18n.LANGUAGEMAPPING.Branches), m.WindowLeftPanelWidth-constant.ListItemOrTitleWidthPad-2)
	m.CurrentRepoBranchesInfoList.Styles.Title = style.TitleStyle
	m.CurrentRepoBranchesInfoList.Styles.PaginationStyle = style.PaginationStyle
	m.CurrentRepoBranchesInfoList.Styles.TitleBar = style.NewStyle

	// Custom Help Model for Count Display
	m.CurrentRepoBranchesInfoList.SetShowHelp(true)
	m.CurrentRepoBranchesInfoList.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	m.CurrentRepoBranchesInfoList.AdditionalShortHelpKeys = func() []key.Binding {
		currentIndex := m.CurrentRepoBranchesInfoList.Index() + 1
		totalCount := len(m.CurrentRepoBranchesInfoList.Items())
		countStr := fmt.Sprintf("%d/%d", currentIndex, totalCount)
		countStr = utils.TruncateString(countStr, m.WindowLeftPanelWidth-constant.ListItemOrTitleWidthPad-2)
		if totalCount == 0 {
			countStr = "0/0"
		}
		return []key.Binding{
			key.NewBinding(
				key.WithKeys(countStr),
				key.WithHelp(countStr, ""),
			),
		}
	}

	if m.ListNavigationIndexPosition.LocalBranchComponent > len(m.CurrentRepoBranchesInfoList.Items())-1 {
		m.CurrentRepoBranchesInfoList.Select(len(m.CurrentRepoBranchesInfoList.Items()) - 1)
		m.ListNavigationIndexPosition.LocalBranchComponent = len(m.CurrentRepoBranchesInfoList.Items()) - 1
	} else {
		m.CurrentRepoBranchesInfoList.Select(m.ListNavigationIndexPosition.LocalBranchComponent)
	}
}
