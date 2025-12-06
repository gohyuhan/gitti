package branch

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// init the popup model for creating a new branch
func InitCreateNewBranchPopUpModel(m *types.GittiModel, createType string) {
	NewBranchNameInput := textinput.New()
	NewBranchNameInput.Placeholder = i18n.LANGUAGEMAPPING.CreateNewBranchPrompt
	NewBranchNameInput.Focus()
	NewBranchNameInput.SetVirtualCursor(true)

	NewBranchNameInput.SetWidth(min(constant.MaxCreateNewBranchPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	m.PopUpModel = &CreateNewBranchPopUpModel{
		NewBranchNameInput: NewBranchNameInput,
		CreateType:         createType,
	}
}

// init the popup model for choosing new branch creation option
func InitChooseNewBranchTypePopUpModel(m *types.GittiModel) {
	newBranchTypeOption := []GitNewBranchTypeOptionItem{
		{
			Name:          i18n.LANGUAGEMAPPING.CreateNewBranchTitle,
			Info:          i18n.LANGUAGEMAPPING.CreateNewBranchDescription,
			NewBranchType: git.NEWBRANCH,
		},
		{
			Name:          i18n.LANGUAGEMAPPING.CreateNewBranchAndSwitchTitle,
			Info:          i18n.LANGUAGEMAPPING.CreateNewBranchAndSwitchDescription,
			NewBranchType: git.NEWBRANCHANDSWITCH,
		},
	}

	items := make([]list.Item, 0, len(newBranchTypeOption))
	for _, newBranchOption := range newBranchTypeOption {
		items = append(items, GitNewBranchTypeOptionItem(newBranchOption))
	}
	width := (min(constant.MaxChooseNewBranchTypePopUpWidth, int(float64(m.Width)*0.8)) - 4)
	nBTOL := list.New(items, GitNewBranchTypeOptionDelegate{}, width, constant.PopUpChooseNewBranchTypeHeight)
	nBTOL.SetShowPagination(false)
	nBTOL.SetShowStatusBar(false)
	nBTOL.SetFilteringEnabled(false)
	nBTOL.SetShowHelp(false)
	nBTOL.SetShowTitle(false)

	m.PopUpModel = &ChooseNewBranchTypeOptionPopUpModel{
		NewBranchTypeOptionList: nBTOL,
	}
}

// init the popup model for switching branch
func InitChooseSwitchBranchTypePopUpModel(m *types.GittiModel, branchName string) {
	switchBranchTypeOption := []GitSwitchBranchTypeOptionItem{
		{
			Name:             i18n.LANGUAGEMAPPING.SwitchBranchTitle,
			Info:             i18n.LANGUAGEMAPPING.SwitchBranchDescription,
			SwitchBranchType: git.SWITCHBRANCH,
		},
		{
			Name:             i18n.LANGUAGEMAPPING.SwitchBranchWithChangesTitle,
			Info:             i18n.LANGUAGEMAPPING.SwitchBranchWithChangesDescription,
			SwitchBranchType: git.SWITCHBRANCHWITHCHANGES,
		},
	}

	items := make([]list.Item, 0, len(switchBranchTypeOption))
	for _, switchBranchOption := range switchBranchTypeOption {
		items = append(items, GitSwitchBranchTypeOptionItem(switchBranchOption))
	}

	width := (min(constant.MaxChooseSwitchBranchTypePopUpWidth, int(float64(m.Width)*0.8)) - 4)
	sBTOL := list.New(items, GitSwitchBranchTypeOptionDelegate{}, width, constant.PopUpChooseSwitchBranchTypeHeight)
	sBTOL.SetShowPagination(false)
	sBTOL.SetShowStatusBar(false)
	sBTOL.SetFilteringEnabled(false)
	sBTOL.SetShowHelp(false)
	sBTOL.SetShowTitle(false)

	m.PopUpModel = &ChooseSwitchBranchTypePopUpModel{
		SwitchTypeOptionList: sBTOL,
		BranchName:           branchName,
	}
}

func InitSwitchBranchOutputPopUpModel(m *types.GittiModel, branchName string, switchType string) {
	// for git push output viewport,
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpSwitchBranchOutputViewPortHeight)
	vp.SetWidth(min(constant.MaxSwitchBranchOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &SwitchBranchOutputPopUpModel{
		BranchName:                 branchName,
		SwitchType:                 switchType,
		SwitchBranchOutputViewport: vp,
		Spinner:                    s,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	m.PopUpModel = popUpModel
}
