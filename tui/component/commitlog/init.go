package commitlog

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
//	Rebuild the commit log list widget from the latest git log data, preserve the previously
//	selected commit by hash, and return true if the selected commit changed (signals that the
//	detail panel needs to be reinitialized).
//
// ------------------------------------
func InitGitCommitLogList(m *types.GittiModel) bool {
	latestGitCommitLog := m.GitOperations.GitCommitLog.GitCommitLogOutput()
	latestGitCommitLogItemArray := make([]list.Item, 0, len(latestGitCommitLog))

	// get the previous selected commit log and see if it was within the new list if yes get the latest position of the previous selected file
	previousSelectedCommitLog := m.CurrentRepoCommitLogInfoList.SelectedItem()
	var prevHash string
	if previousSelectedCommitLog != nil {
		prevHash = previousSelectedCommitLog.(GitCommitLogItem).Hash
	}
	selectedCommitLogPosition := -1

	titleWidthLimit := m.WindowLeftPanelWidth - constant.ListItemOrTitleWidthPad - 2

	if previousSelectedCommitLog != nil {
		for index, commitLog := range latestGitCommitLog {
			//  we use hash here to determine if it was the same commit log as the hash is unique
			if commitLog.Hash == prevHash {
				selectedCommitLogPosition = index
			}
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
	} else {
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
	}

	previousCommitLogCount := len(m.CurrentRepoCommitLogInfoList.Items())

	m.CurrentRepoCommitLogInfoList = list.New(latestGitCommitLogItemArray, GitCommitLogItemDelegate{}, m.WindowLeftPanelWidth, m.CommitLogComponentPanelHeight)
	m.CurrentRepoCommitLogInfoList.SetShowPagination(false)
	m.CurrentRepoCommitLogInfoList.SetShowStatusBar(false)
	m.CurrentRepoCommitLogInfoList.SetFilteringEnabled(false)
	m.CurrentRepoCommitLogInfoList.SetShowFilter(false)
	m.CurrentRepoCommitLogInfoList.Title = ansi.Truncate(ConstructCommitLogComponentTitle(titleWidthLimit), titleWidthLimit, "...")
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
