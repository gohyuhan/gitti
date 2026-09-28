package commitfiles

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/list"
	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	Initialize or update the commit files list for the currently drilled-down commit.
//
// ------------------------------------
func InitCommitFilesList(ctx context.Context, m *types.GittiModel, commitHash string, commitSubject string) bool {
	touchedFiles := m.GitOperations.GitCommitLog.GetCommitTouchedFiles(ctx, commitHash)
	items := make([]list.Item, 0, len(touchedFiles))

	for _, file := range touchedFiles {
		items = append(items, GitCommitFileItem{
			CommitHash:   commitHash,
			Status:       file.Status,
			FilePathname: file.FilePathname,
			OldPathname:  file.OldPathname,
		})
	}

	m.CurrentDrillDownCommitHash = commitHash
	m.CurrentDrillDownCommitSubject = commitSubject

	previousPosition := m.ListNavigationIndexPosition.CommitFilesComponent
	var previousSelected list.Item
	if len(m.CurrentRepoCommitFilesList.Items()) > 0 {
		previousSelected = m.CurrentRepoCommitFilesList.SelectedItem()
	}

	// Apply panel filtering if active
	items, selectedPosition := utils.FilterListItems(items, m.PanelFilterQuery[constant.CommitFilesComponentPanel], previousSelected, previousPosition)

	shortHash := commitHash
	if len(shortHash) > 7 {
		shortHash = shortHash[:7]
	}

	titleWidthLimit := m.WindowLeftPanelWidth - constant.ListItemOrTitleWidthPad - 2
	titleText := fmt.Sprintf("Commit Files: %s (%d files)", shortHash, len(items))

	m.CurrentRepoCommitFilesList = list.New(items, GitCommitFileItemDelegate{}, m.WindowLeftPanelWidth, m.CommitLogComponentPanelHeight)
	m.CurrentRepoCommitFilesList.SetShowPagination(false)
	m.CurrentRepoCommitFilesList.SetShowStatusBar(false)
	m.CurrentRepoCommitFilesList.SetFilteringEnabled(false)
	m.CurrentRepoCommitFilesList.SetShowFilter(false)
	m.CurrentRepoCommitFilesList.Title = ansi.Truncate(titleText, titleWidthLimit, "...")
	m.CurrentRepoCommitFilesList.Styles.Title = style.TitleStyle
	m.CurrentRepoCommitFilesList.Styles.TitleBar = style.NewStyle
	m.CurrentRepoCommitFilesList.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)

	m.CurrentRepoCommitFilesList.SetShowHelp(true)
	m.CurrentRepoCommitFilesList.KeyMap = list.KeyMap{}
	m.CurrentRepoCommitFilesList.AdditionalShortHelpKeys = utils.ListCounterHelper(m, &m.CurrentRepoCommitFilesList, constant.CommitFilesComponentPanel)

	if len(items) > 0 {
		targetIndex := 0
		if selectedPosition >= 0 && selectedPosition < len(items) {
			targetIndex = selectedPosition
		}
		m.CurrentRepoCommitFilesList.Select(targetIndex)
		m.ListNavigationIndexPosition.CommitFilesComponent = targetIndex
		return true
	}

	return false
}
