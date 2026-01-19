package log

import (
	"sort"

	"charm.land/bubbles/v2/list"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/component/commitlog"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

func InitGitCherryPickOptionSelectionPopUp(m *types.GittiModel) {
	newCherryPickOpsOption := []CherryPickOpsOptionItem{
		{
			Name:              i18n.LANGUAGEMAPPING.CherryPickOpsTitle,
			Info:              i18n.LANGUAGEMAPPING.CherryPickOpsDescription,
			CherryPickOpsType: constant.CHERRYPICK,
		},
		{
			Name:              i18n.LANGUAGEMAPPING.EditCherryPickOpsTitle,
			Info:              i18n.LANGUAGEMAPPING.EditCherryPickOpsDescription,
			CherryPickOpsType: constant.EDITCHERRYPICK,
		},
		{
			Name:              i18n.LANGUAGEMAPPING.ApplyCherryPickOpsTitle,
			Info:              i18n.LANGUAGEMAPPING.ApplyCherryPickOpsDescription,
			CherryPickOpsType: constant.APPLYCHERRYPICK,
		},
	}

	items := make([]list.Item, 0, len(newCherryPickOpsOption))
	for _, newCherryPickOption := range newCherryPickOpsOption {
		items = append(items, CherryPickOpsOptionItem(newCherryPickOption))
	}
	width := (min(constant.MaxGitCherryPickOptionSelectionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	gCPOOL := list.New(items, CherryPickOpsOptionDelegate{}, width, constant.PopUpGitCherryPickOptionSelectionHeight)
	gCPOOL.SetShowPagination(false)
	gCPOOL.SetShowStatusBar(false)
	gCPOOL.SetFilteringEnabled(false)
	gCPOOL.SetShowTitle(false)

	// Custom Help Model for Count Display
	gCPOOL.SetShowHelp(true)
	gCPOOL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	gCPOOL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	gCPOOL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &gCPOOL, constant.MaxGitCherryPickOptionSelectionPopUpWidth)

	popUpModel := &GitCherryPickOptionSelectionPopUpModel{
		CherryPickedOpsOption: gCPOOL,
	}

	m.PopUpModel = popUpModel
}

func InitGitCherryPickPopUp(m *types.GittiModel, branchName string) {
	items := make([]list.Item, 0, len(m.CurrentRepoCommitLogInfoList.Items()))
	for _, item := range m.CurrentRepoCommitLogInfoList.Items() {
		if commitItem, ok := item.(commitlog.GitCommitLogItem); ok {
			items = append(items, GitCherryPickItem{
				Hash:       commitItem.Hash,
				Message:    commitItem.Message,
				Author:     commitItem.Author,
				FromBranch: branchName,
			})
		}
	}

	width := (min(constant.MaxGitCherryPickPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	// we are dividing by 2 because the list is a single string that have 1 \n that break into 2 line
	height := max(constant.PopUpGitCherryPickPopUpHeight, int((float64(m.Height)*0.8)/2)-3)
	gCPL := list.New(items, GitCherryPickDelegate{&m.CherryPickedCommitInfo.CherryPickedCommitMap}, width, height)
	gCPL.SetShowPagination(false)
	gCPL.SetShowStatusBar(false)
	gCPL.SetFilteringEnabled(false)
	gCPL.SetShowTitle(false)

	// Custom Help Model for Count Display
	gCPL.SetShowHelp(true)
	gCPL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	gCPL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	gCPL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &gCPL, constant.MaxGitCherryPickPopUpWidth)

	popUpModel := &GitCherryPickPopUpModel{
		CurrentBranchName:                branchName,
		CurrentBranchCherryPickCommitLog: gCPL,
	}

	m.PopUpModel = popUpModel
}

func InitGitEditCherryPickPopUp(m *types.GittiModel, selectionIndex int) {
	// selectionIndex is need here to retain the user selection index cursor when we reinit the list due to removal or any modification to the list
	// it will be 0 when user newly enter the cherry pick edit UI
	items := make([]list.Item, 0, len(m.CherryPickedCommitInfo.CherryPickedCommitMap))
	sortedCommits := make([]git.CherryPickedCommitLog, 0, len(m.CherryPickedCommitInfo.CherryPickedCommitMap))
	for _, commitItem := range m.CherryPickedCommitInfo.CherryPickedCommitMap {
		sortedCommits = append(sortedCommits, commitItem)
	}
	sort.Slice(sortedCommits, func(i, j int) bool {
		return sortedCommits[i].UserSelectedSequence < sortedCommits[j].UserSelectedSequence
	})
	for _, commitItem := range sortedCommits {
		items = append(items, GitEditCherryPickItem{
			Hash:       commitItem.Hash,
			Message:    commitItem.Message,
			Author:     commitItem.Author,
			FromBranch: commitItem.FromBranch,
		})
	}

	width := (min(constant.MaxGitEditCherryPickPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	// we are dividing by 3 because the list is a single string that have 2 \n that break into 3 line
	height := max(constant.PopUpGitEditCherryPickPopUpHeight, int((float64(m.Height)*0.8)/3)-3)
	gECPL := list.New(items, GitEditCherryPickDelegate{}, width, height)
	gECPL.SetShowPagination(false)
	gECPL.SetShowStatusBar(false)
	gECPL.SetFilteringEnabled(false)
	gECPL.SetShowTitle(false)

	// Custom Help Model for Count Display
	gECPL.SetShowHelp(true)
	gECPL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	gECPL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	gECPL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &gECPL, constant.MaxGitEditCherryPickPopUpWidth)

	// init the cursor index position
	// only trigger if the selectedIndex is > 0
	if selectionIndex > 0 {
		totalItemCount := len(items)
		if selectionIndex >= totalItemCount {
			gECPL.Select(totalItemCount - 1)
		} else {
			gECPL.Select(selectionIndex)
		}
	}

	popUpModel := &GitEditCherryPickPopUpModel{
		CherryPickedCommitLog: gECPL,
	}

	m.PopUpModel = popUpModel
}
