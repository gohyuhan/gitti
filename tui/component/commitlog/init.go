package commitlog

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

// init the list component for commit log Component
func InitGitCommitLogList(m *types.GittiModel) bool {
	latestGitCommitLog := m.GitOperations.GitCommitLog.GitCommitLogOutput()
	var latestGitCommitLogItemArray []list.Item

	for _, commitLog := range latestGitCommitLog {
		laneCharList := make([]Cell, len(commitLog.LaneCharInfo))
		for i, c := range commitLog.LaneCharInfo {
			laneCharList[i] = Cell{
				Char:    c.Char,
				ColorID: c.ColorID,
			}
		}

		latestGitCommitLogItemArray = append(latestGitCommitLogItemArray, GitCommitLogItem{
			Hash:         commitLog.Hash,
			Parents:      commitLog.Parents,
			Message:      commitLog.Message,
			Author:       commitLog.Author,
			LaneCharList: laneCharList,
			ColorID:      commitLog.ColorID,
		})
	}

	// get the previous selected commit log and see if it was within the new list if yes get the latest position of the previous selected file
	previousSelectedCommitLog := m.CurrentRepoCommitLogInfoList.SelectedItem()
	var prevHash string
	if previousSelectedCommitLog != nil {
		prevHash = previousSelectedCommitLog.(GitCommitLogItem).Hash
	}
	selectedCommitLogPosition := -1

	if previousSelectedCommitLog != nil {
		for index, item := range latestGitCommitLogItemArray {
			if item.(GitCommitLogItem).Hash == prevHash {
				selectedCommitLogPosition = index
				break
			}
		}
	}

	previousCommitLogCount := len(m.CurrentRepoCommitLogInfoList.Items())

	m.CurrentRepoCommitLogInfoList = list.New(latestGitCommitLogItemArray, GitCommitLogItemDelegate{}, m.WindowLeftPanelWidth, m.CommitLogComponentPanelHeight)
	m.CurrentRepoCommitLogInfoList.SetShowPagination(false)
	m.CurrentRepoCommitLogInfoList.SetShowStatusBar(false)
	m.CurrentRepoCommitLogInfoList.SetFilteringEnabled(false)
	m.CurrentRepoCommitLogInfoList.SetShowFilter(false)
	m.CurrentRepoCommitLogInfoList.Title = utils.TruncateString(fmt.Sprintf("[3] \ue729 %s:", i18n.LANGUAGEMAPPING.CommitLog), m.WindowLeftPanelWidth-constant.ListItemOrTitleWidthPad-2)
	m.CurrentRepoCommitLogInfoList.Styles.Title = style.TitleStyle
	m.CurrentRepoCommitLogInfoList.Styles.PaginationStyle = style.PaginationStyle
	m.CurrentRepoCommitLogInfoList.Styles.TitleBar = style.NewStyle
	m.CurrentRepoCommitLogInfoList.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)

	// Custom Help Model for Count Display
	m.CurrentRepoCommitLogInfoList.SetShowHelp(true)
	m.CurrentRepoCommitLogInfoList.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	m.CurrentRepoCommitLogInfoList.AdditionalShortHelpKeys = utils.ListCounterHelper(m, &m.CurrentRepoCommitLogInfoList)

	if len(latestGitCommitLog) < 1 {
		return len(latestGitCommitLog) != previousCommitLogCount
	}

	if selectedCommitLogPosition >= 0 {
		m.CurrentRepoCommitLogInfoList.Select(selectedCommitLogPosition)
		m.ListNavigationIndexPosition.CommitLogComponent = selectedCommitLogPosition
	} else {
		if m.ListNavigationIndexPosition.CommitLogComponent > len(m.CurrentRepoCommitLogInfoList.Items())-1 {
			m.CurrentRepoCommitLogInfoList.Select(len(m.CurrentRepoCommitLogInfoList.Items()) - 1)
			m.ListNavigationIndexPosition.CommitLogComponent = len(m.CurrentRepoCommitLogInfoList.Items()) - 1
		} else {
			m.CurrentRepoCommitLogInfoList.Select(m.ListNavigationIndexPosition.CommitLogComponent)
		}
	}

	if previousSelectedCommitLog != nil {
		curr := m.CurrentRepoCommitLogInfoList.SelectedItem()
		if curr != nil && curr.(GitCommitLogItem).Hash == prevHash {
			return false
		}
	}
	return true
}
