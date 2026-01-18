package log

import (
	"charm.land/bubbles/v2/list"
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
			CherryPickOpsType: constant.EDITCHERRYPICk,
		},
	}

	items := make([]list.Item, 0, len(newCherryPickOpsOption))
	for _, newCherryPickOption := range newCherryPickOpsOption {
		items = append(items, CherryPickOpsOptionItem(newCherryPickOption))
	}
	width := (min(constant.MaxGitCherryPickOptionSelectionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	gCPL := list.New(items, CherryPickOpsOptionDelegate{}, width, constant.PopUpGitCherryPickOptionSelectionHeight)
	gCPL.SetShowPagination(false)
	gCPL.SetShowStatusBar(false)
	gCPL.SetFilteringEnabled(false)
	gCPL.SetShowTitle(false)

	// Custom Help Model for Count Display
	gCPL.SetShowHelp(true)
	gCPL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	gCPL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	gCPL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &gCPL, constant.MaxGitCherryPickOptionSelectionPopUpWidth)

	popUpModel := &GitCherryPickOptionSelectionPopUpModel{
		CherryPickedOpsOption: gCPL,
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
	height := max(constant.PopUpGitCherryPickPopUpHeight, int(float64(m.Height)*0.5)-3)
	gCPL := list.New(items, GitCherryPickDelegate{&m.CherryPickedCommitLogList, &m.CherryPickedCommitMap}, width, height)
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

func InitGitEditCherryPickPopUp(m *types.GittiModel) {
	items := make([]list.Item, 0, len(m.CherryPickedCommitLogList))
	for _, commitItem := range m.CherryPickedCommitLogList {
		items = append(items, GitCherryPickItem{
			Hash:       commitItem.Hash,
			Message:    commitItem.Message,
			Author:     commitItem.Author,
			FromBranch: commitItem.FromBranch,
		})
	}

	width := (min(constant.MaxGitEditCherryPickPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	height := max(constant.PopUpGitEditCherryPickPopUpHeight, int(float64(m.Height)*0.5)-3)
	cPCL := list.New(items, GitCherryPickDelegate{&m.CherryPickedCommitLogList, &m.CherryPickedCommitMap}, width, height)
	cPCL.SetShowPagination(false)
	cPCL.SetShowStatusBar(false)
	cPCL.SetFilteringEnabled(false)
	cPCL.SetShowTitle(false)

	// Custom Help Model for Count Display
	cPCL.SetShowHelp(true)
	cPCL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	cPCL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	cPCL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &cPCL, constant.MaxGitEditCherryPickPopUpWidth)

	popUpModel := &GitEditCherryPickPopUpModel{
		CherryPickedCommitLog: cPCL,
	}

	m.PopUpModel = popUpModel
}
