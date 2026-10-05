package commitfiles

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/list"
	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

var (
	cachedCommitFiles []list.Item
	cachedCommitHash  string
)

// ------------------------------------
//
//	Load touched files for a commit into memory
//
// ------------------------------------
func LoadCommitFiles(ctx context.Context, m *types.GittiModel, commitHash string) {
	cachedCommitHash = commitHash
	touchedFiles := m.GitOperations.GitCommitLog.GetCommitTouchedFiles(ctx, commitHash)
	cachedCommitFiles = make([]list.Item, 0, len(touchedFiles))

	for _, file := range touchedFiles {
		cachedCommitFiles = append(cachedCommitFiles, GitCommitFileItem{
			CommitHash:   commitHash,
			Status:       file.Status,
			FilePathname: file.FilePathname,
			OldPathname:  file.OldPathname,
		})
	}
}

// ------------------------------------
//
//	Clear cached commit files when exiting the panel
//
// ------------------------------------
func ClearCachedCommitFiles() {
	cachedCommitFiles = nil
	cachedCommitHash = ""
}

// ------------------------------------
//
//	Initialize or update the commit files list from in-memory cache.
//
// ------------------------------------
func InitCommitFilesList(m *types.GittiModel) bool {
	items := make([]list.Item, len(cachedCommitFiles))
	copy(items, cachedCommitFiles)

	previousPosition := m.ListNavigationIndexPosition.CommitFilesComponent
	var previousSelected list.Item
	if len(m.CurrentRepoCommitFilesList.Items()) > 0 {
		previousSelected = m.CurrentRepoCommitFilesList.SelectedItem()
	}

	// Apply panel filtering if active
	items, selectedPosition := utils.FilterListItems(items, m.PanelFilterQuery[constant.CommitFilesComponentPanel], previousSelected, previousPosition)

	shortHash := cachedCommitHash
	if len(shortHash) > 7 {
		shortHash = shortHash[:7]
	}

	titleWidthLimit := m.WindowLeftPanelWidth - constant.ListItemOrTitleWidthPad - 2
	titleText := fmt.Sprintf(i18n.LANGUAGEMAPPING.CommitFilesTitleFormat, shortHash, len(items))

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
