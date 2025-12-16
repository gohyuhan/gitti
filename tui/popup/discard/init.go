package discard

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

// for discard option list popup
func InitGitDiscardTypeOptionPopUp(m *types.GittiModel, filePathName string, newlyAddedOrCopiedFile bool, renameFile bool) {
	discardTypeOption := []GitDiscardTypeOptionItem{
		{
			Name:        i18n.LANGUAGEMAPPING.GitDiscardWhole,
			Info:        fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardWholeInfo, filePathName),
			DiscardType: git.DISCARDWHOLE,
		},
		{
			Name:        i18n.LANGUAGEMAPPING.GitDiscardUnstage,
			Info:        fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardUnstageInfo, filePathName),
			DiscardType: git.DISCARDUNSTAGE,
		},
	}

	if newlyAddedOrCopiedFile {
		discardTypeOption = []GitDiscardTypeOptionItem{
			{
				Name:        i18n.LANGUAGEMAPPING.GitDiscardWhole,
				Info:        fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardWholeInfo, filePathName),
				DiscardType: git.DISCARDNEWLYADDEDORCOPIED,
			},
			{
				Name:        i18n.LANGUAGEMAPPING.GitDiscardUnstage,
				Info:        fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardUnstageInfo, filePathName),
				DiscardType: git.DISCARDUNSTAGE,
			},
		}
	}

	if renameFile {
		discardTypeOption = []GitDiscardTypeOptionItem{
			{
				Name:        i18n.LANGUAGEMAPPING.GitDiscardAndRevertRename,
				Info:        fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardAndRevertRenameInfo, filePathName),
				DiscardType: git.DISCARDANDREVERTRENAME,
			},
			{
				Name:        i18n.LANGUAGEMAPPING.GitDiscardUnstage,
				Info:        fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDiscardUnstageInfo, filePathName),
				DiscardType: git.DISCARDUNSTAGE,
			},
		}
	}

	items := make([]list.Item, 0, len(discardTypeOption))
	for _, discardOption := range discardTypeOption {
		items = append(items, GitDiscardTypeOptionItem(discardOption))
	}

	width := (min(constant.MaxGitDiscardTypeOptionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	gDTOL := list.New(items, GitDiscardTypeOptionDelegate{}, width, constant.PopUpGitDiscardTypeOptionHeight)
	gDTOL.SetShowPagination(false)
	gDTOL.SetShowStatusBar(false)
	gDTOL.SetFilteringEnabled(false)
	gDTOL.SetShowTitle(false)

	// Custom Help Model for Count Display
	gDTOL.SetShowHelp(true)
	gDTOL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	gDTOL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	gDTOL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &gDTOL, constant.MaxGitDiscardTypeOptionPopUpWidth)

	popUpModel := &GitDiscardTypeOptionPopUpModel{
		DiscardTypeOptionList: gDTOL,
		FilePathName:          filePathName,
	}

	m.PopUpModel = popUpModel
}

// for discard confirm prompt
func InitGitDiscardConfirmPromptPopupModel(m *types.GittiModel, filePathName string, discardType string) {
	popUpModel := &GitDiscardConfirmPromptPopUpModel{
		FilePathName: filePathName,
		DiscardType:  discardType,
	}
	m.PopUpModel = popUpModel
}
