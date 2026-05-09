package interactiverebase

import (
	"charm.land/bubbles/v2/list"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

func InitInteractiveRebaseOptionPopUp(m *types.GittiModel) {
	interactiveRebaseOption := []InteractiveRebaseOptionItem{
		{
			Name:                  i18n.LANGUAGEMAPPING.InteractiveRebaseFixupSquash,
			Info:                  i18n.LANGUAGEMAPPING.InteractiveRebaseFixupSquashDescription,
			InteractiveRebaseType: git.FIXUPSQUASH,
		},
		{
			Name:                  i18n.LANGUAGEMAPPING.InteractiveRebaseReword,
			Info:                  i18n.LANGUAGEMAPPING.InteractiveRebaseFeatureComingSoon,
			InteractiveRebaseType: git.REWORD,
		},
		{
			Name:                  i18n.LANGUAGEMAPPING.InteractiveRebaseDrop,
			Info:                  i18n.LANGUAGEMAPPING.InteractiveRebaseFeatureComingSoon,
			InteractiveRebaseType: git.DROP,
		},
	}

	items := make([]list.Item, 0, len(interactiveRebaseOption))
	for _, interactiveRebaseOption := range interactiveRebaseOption {
		items = append(items, InteractiveRebaseOptionItem(interactiveRebaseOption))
	}

	width := (min(constant.MaxInteractiveRebaseOptionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	iROL := list.New(items, InteractiveRebaseOptionDelegate{}, width, constant.PopUpInteractiveRebaseOptionHeight)
	iROL.SetShowPagination(false)
	iROL.SetShowStatusBar(false)
	iROL.SetFilteringEnabled(false)
	iROL.SetShowTitle(false)

	// Custom Help Model for Count Display
	iROL.SetShowHelp(true)
	iROL.KeyMap = list.KeyMap{}
	iROL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	iROL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &iROL, constant.MaxInteractiveRebaseOptionPopUpWidth)

	popUpModel := &InteractiveRebaseOptionPopUp{
		InteractiveRebaseOptionList: iROL,
	}

	m.PopUpModel = popUpModel
}
