package resolve

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// for resolve conflict option list popup
func InitGitResolveConflictOptionPopUpModel(m *types.GittiModel, filePathName string) {
	resolveConflictOption := []GitResolveConflictOptionItem{
		{
			Name:        i18n.LANGUAGEMAPPING.GitResolveConflictReset,
			Info:        fmt.Sprintf(i18n.LANGUAGEMAPPING.GitResolveConflictResetInfo, filePathName),
			ResolveType: git.RESETCONFLICT,
		},
		{
			Name:        i18n.LANGUAGEMAPPING.GitResolveConflictAcceptLocalChanges,
			Info:        fmt.Sprintf(i18n.LANGUAGEMAPPING.GitResolveConflictAcceptLocalChangesInfo, filePathName),
			ResolveType: git.CONFLICTACCEPTLOCALCHANGES,
		},
		{
			Name:        i18n.LANGUAGEMAPPING.GitResolveConflictAcceptIncomingChanges,
			Info:        fmt.Sprintf(i18n.LANGUAGEMAPPING.GitResolveConflictAcceptIncomingChangesInfo, filePathName),
			ResolveType: git.CONFLICTACCEPTINCOMINGCHANGES,
		},
	}

	items := make([]list.Item, 0, len(resolveConflictOption))
	for _, resolveConflictOption := range resolveConflictOption {
		items = append(items, GitResolveConflictOptionItem(resolveConflictOption))
	}

	width := (min(constant.MaxGitResolveConflictOptionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	rROL := list.New(items, GitResolveConflictOptionDelegate{}, width, constant.PopUpGitResolveConflictOptionPopUpHeight)
	rROL.SetShowPagination(false)
	rROL.SetShowStatusBar(false)
	rROL.SetFilteringEnabled(false)
	rROL.SetShowTitle(false)

	// Custom Help Model for Count Display
	rROL.SetShowHelp(true)
	rROL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	rROL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	rROL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &rROL, constant.MaxGitResolveConflictOptionPopUpWidth)

	popUpModel := &GitResolveConflictOptionPopUpModel{
		ResolveConflictOptionList: rROL,
		FilePathName:              filePathName,
	}

	m.PopUpModel = popUpModel
}
