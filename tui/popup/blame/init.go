package blame

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	Initialize the blame popup model with file list, filter input, and viewport
//
// ------------------------------------
func InitBlamePopUpModel(m *types.GittiModel) {
	width := int(float64(m.Width)*0.9) - 4
	height := int(float64(m.Height) * 0.9)
	listHeight := height - 2 // account for title and custom filter input; list's built-in filtering is disabled in favour of manual filter input

	var items []list.Item
	for _, filepath := range m.GitOperations.GitBlame.GetCurrentGitTrackedFiles() {
		items = append(items, CurrentGitTrackedFilesPathItem{FilePath: filepath})
	}
	cGTFPL := list.New(items, CurrentGitTrackedFilesPathDelegate{}, width, listHeight)
	cGTFPL.SetShowPagination(false)
	cGTFPL.SetShowStatusBar(false)
	cGTFPL.SetFilteringEnabled(true)
	cGTFPL.SetShowTitle(false)

	filterInput := textinput.New()
	filterInput.SetValue("")
	filterInput.Placeholder = i18n.LANGUAGEMAPPING.BlameFilePathFilterPlaceholder
	filterInput.Focus()
	filterInput.SetVirtualCursor(true)

	// Custom Help Model for Count Display
	cGTFPL.SetShowHelp(true)
	cGTFPL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	cGTFPL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	cGTFPL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &cGTFPL, width)

	vp := viewport.New()
	vp.SoftWrap = false
	vp.SetHorizontalStep(2)
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(height)
	vp.SetWidth(width)

	popUpModel := &BlamePopUpModel{
		CurrentGitTrackedFilesPathList: cGTFPL,
		FilterInput:                    filterInput,
		BlameViewport:                  vp,
		ShowingBlameInfo:               false,
		HasFilePathChosen:              false,
		FilterValue:                    "",
		SelectedFilePath:               "",
	}

	m.PopUpModel = popUpModel
}
