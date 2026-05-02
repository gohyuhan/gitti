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
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	init the popup model for creating a new branch
//
// ------------------------------------
func InitCreateNewBranchPopUpModel(m *types.GittiModel, createType string, commitHash string) {
	NewBranchNameInput := textinput.New()
	NewBranchNameInput.Placeholder = i18n.LANGUAGEMAPPING.CreateNewBranchPrompt
	NewBranchNameInput.Focus()
	NewBranchNameInput.SetVirtualCursor(true)

	NewBranchNameInput.SetWidth(min(constant.MaxCreateNewBranchPopUpWidth, int(float64(m.Width)*0.8)) - 6)
	m.PopUpModel = &CreateNewBranchPopUpModel{
		NewBranchNameInput: NewBranchNameInput,
		CreateType:         createType,
		CommitHash:         commitHash,
	}
}

// ------------------------------------
//
//	init the popup model for choosing new branch creation option
//
// ------------------------------------
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
		{
			Name:          i18n.LANGUAGEMAPPING.CreateNewBranchBasedOnRemoteUserInputTitle,
			Info:          i18n.LANGUAGEMAPPING.CreateNewBranchBasedOnRemoteUserInputDescription,
			NewBranchType: git.NEWBRANCHBASEDONREMOTEUSERINPUT,
		},
		{
			Name:          i18n.LANGUAGEMAPPING.CreateNewBranchBasedOnRemoteUserSelectionTitle,
			Info:          i18n.LANGUAGEMAPPING.CreateNewBranchBasedOnRemoteUserSelectionDescription,
			NewBranchType: git.NEWBRANCHBASEDONREMOTEUSERSELECT,
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
	nBTOL.SetShowTitle(false)

	// Custom Help Model for Count Display
	nBTOL.SetShowHelp(true)
	nBTOL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	nBTOL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	nBTOL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &nBTOL, constant.MaxChooseNewBranchTypePopUpWidth)

	m.PopUpModel = &ChooseNewBranchTypeOptionPopUpModel{
		NewBranchTypeOptionList: nBTOL,
	}
}

// ------------------------------------
//
//	init the popup model for switching branch
//
// ------------------------------------
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
	sBTOL.SetShowTitle(false)

	// Custom Help Model for Count Display
	sBTOL.SetShowHelp(true)
	sBTOL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	sBTOL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	sBTOL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &sBTOL, constant.MaxChooseSwitchBranchTypePopUpWidth)

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

func InitGitDeleteBranchConfirmPromptPopUpModel(m *types.GittiModel, branchName string) {
	popUpModel := &GitDeleteBranchConfirmPromptPopUpModel{
		BranchName: branchName,
	}
	m.PopUpModel = popUpModel
}

func InitGitDeleteBranchOutputPopUpModel(m *types.GittiModel) {
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpGitDeleteBranchOutputViewportHeight)
	vp.SetWidth(min(constant.MaxGitDeleteBranchOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &GitDeleteBranchOutputPopUpModel{
		BranchDeleteOutputViewport: vp,
		Spinner:                    s,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)

	m.PopUpModel = popUpModel
}

func InitCreateBranchBasedOnRemotePopUp(m *types.GittiModel, remoteOrigin string) {
	remoteBranchNameInput := textinput.New()
	remoteBranchNameInput.Placeholder = i18n.LANGUAGEMAPPING.EnterRemoteBranchPrompt
	remoteBranchNameInput.Focus()
	remoteBranchNameInput.SetVirtualCursor(true)

	remoteBranchNameInput.SetWidth(min(constant.MaxCreateNewBranchPopUpWidth, int(float64(m.Width)*0.8)) - 6)

	popUpModel := &CreateBranchBasedOnRemotePopUpModel{
		RemoteOrigin:          remoteOrigin,
		RemoteBranchNameInput: remoteBranchNameInput,
	}

	m.PopUpModel = popUpModel
}

func InitCreateBranchBasedOnRemoteOutputPopUp(m *types.GittiModel) {
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpCreateBranchBasedOnRemoteOutputViewportHeight)
	vp.SetWidth(min(constant.MaxCreateBranchBasedOnRemoteOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &CreateBranchBasedOnRemoteOutputPopUpModel{
		CreateBranchBasedOnRemoteOutputViewport: vp,
		Spinner:                                 s,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)

	m.PopUpModel = popUpModel
}

func InitChooseRemoteBranchOptionPopUpModel(m *types.GittiModel) {
	remoteBranches := m.GitOperations.GitBranch.RemoteBranches()
	items := make([]list.Item, 0, len(remoteBranches))
	for _, remoteBranch := range remoteBranches {
		items = append(items, RemoteBranchItem(remoteBranch))
	}
	width := (min(constant.MaxChooseRemoteBranchOptionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	cRBOL := list.New(items, RemoteBranchItemDelegate{}, width, constant.PopUpChooseRemoteBranchOptionHeight)
	cRBOL.SetShowPagination(false)
	cRBOL.SetShowStatusBar(false)
	cRBOL.SetFilteringEnabled(false)
	cRBOL.SetShowTitle(false)

	// Custom Help Model for Count Display
	cRBOL.SetShowHelp(true)
	cRBOL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	cRBOL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	cRBOL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &cRBOL, constant.MaxChooseRemoteBranchOptionPopUpWidth)

	m.PopUpModel = &ChooseRemoteBranchOptionPopUpModel{
		RemoteBranchOptionList: cRBOL,
	}
}

func InitChooseBranchOptionForMergePopUpModel(m *types.GittiModel) {
	branches := m.GitOperations.GitBranch.AllBranches()
	items := make([]list.Item, 0, len(branches))
	for _, branch := range branches {
		if branch.IsCheckedOut {
			continue
		}
		items = append(items, GitMergeBranchOptionItem{BranchName: branch.BranchName})
	}
	width := (min(constant.MaxChooseBranchOptionForMergePopUpWidth, int(float64(m.Width)*0.8)) - 4)

	// for selecting branch for git merge
	cBOFMBOL := list.New(items, GitMergeBranchOptionItemDelegate{}, width, constant.PopUpChooseBranchOptionForMergeBranchOptionHeight)
	cBOFMBOL.SetShowPagination(false)
	cBOFMBOL.SetShowStatusBar(false)
	cBOFMBOL.SetFilteringEnabled(false)
	cBOFMBOL.SetShowTitle(false)

	// Custom Help Model for Count Display
	cBOFMBOL.SetShowHelp(true)
	cBOFMBOL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	cBOFMBOL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	cBOFMBOL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &cBOFMBOL, constant.MaxChooseBranchOptionForMergePopUpWidth)

	// for ALREADY selected branch for git merge
	cBOFMSBOL := list.New([]list.Item{}, GitMergeBranchOptionItemDelegate{}, width, constant.PopUpChooseBranchOptionForMergeSelectedBranchOptionHeight)
	cBOFMSBOL.SetShowPagination(false)
	cBOFMSBOL.SetShowStatusBar(false)
	cBOFMSBOL.SetFilteringEnabled(false)
	cBOFMSBOL.SetShowTitle(false)

	// Custom Help Model for Count Display
	cBOFMSBOL.SetShowHelp(true)
	cBOFMSBOL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	cBOFMSBOL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	cBOFMSBOL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &cBOFMSBOL, constant.MaxChooseBranchOptionForMergePopUpWidth)

	popUpModel := &ChooseBranchOptionForMergePopUpModel{
		BranchOptionList:   cBOFMBOL,
		SelectedBranchList: cBOFMSBOL,
	}
	popUpModel.BranchOptionSectionSelected.Store(true)
	popUpModel.SelectedBranchSectionSelected.Store(false)
	m.PopUpModel = popUpModel
}

func InitBranchMergeOutputPopUpModel(m *types.GittiModel) {
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpBranchMergeOutputViewportHeight)
	vp.SetWidth(min(constant.MaxBranchMergeOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &BranchMergeOutputPopUpModel{
		BranchMergeOutputViewport: vp,
		Spinner:                   s,
	}

	popUpModel.IsProcessing.Store(false)
	popUpModel.IsCancelled.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)

	m.PopUpModel = popUpModel
}
