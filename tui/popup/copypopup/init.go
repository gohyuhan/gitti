package copypopup

import (
	"charm.land/bubbles/v2/list"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

func Init(m *types.GittiModel, options []Option) uint64 {
	items := make([]list.Item, len(options))
	for i, option := range options {
		items[i] = option
	}
	width := min(constant.MaxChooseCopyValuePopUpWidth, int(float64(m.Width)*0.8)) - 4
	optionsList := list.New(items, optionDelegate{}, width, constant.PopUpChooseCopyValueHeight)
	optionsList.SetShowPagination(false)
	optionsList.SetShowStatusBar(false)
	optionsList.SetFilteringEnabled(false)
	optionsList.SetShowTitle(false)
	optionsList.SetShowHelp(true)
	optionsList.KeyMap = list.KeyMap{}
	optionsList.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	m.CopyPopupSequence++
	popUp := &Model{ID: m.CopyPopupSequence, Options: optionsList}
	popUp.Options.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &popUp.Options, constant.MaxChooseCopyValuePopUpWidth)
	m.PopUpModel = popUp
	m.PopUpType = constant.ChooseCopyValuePopUp
	m.ShowPopUp.Store(true)
	m.IsTyping.Store(false)
	return m.CopyPopupSequence
}
