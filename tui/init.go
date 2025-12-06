package tui

import (
	"fmt"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
) // this was for various components part init or reinit function due to update or newly create

// those utf-8 icons for the component can be found at https://www.nerdfonts.com/cheat-sheet
// for bubbletea list component, we can't get rid of the "No items." natively for now as we couldn't access into it or modify it
// see https://github.com/charmbracelet/bubbles/blob/master/list/list.go#L1222

// init the list component for Stash info Component
// return bool was to tell if we need to reinit the detail component panel or not
func initStashList(m *types.GittiModel) bool {
	latestStashArray := m.GitOperations.GitStash.AllStash()
	items := make([]list.Item, 0, len(latestStashArray))
	for _, stashInfo := range latestStashArray {
		items = append(items, gitStashItem(stashInfo))
	}

	// get the previous selected file and see if it was within the new list if yes get the latest position of the previous selected file
	previousSelectedStash := m.CurrentRepoStashInfoList.SelectedItem()
	selectedFilesPosition := -1

	for index, item := range items {
		if item == previousSelectedStash {
			selectedFilesPosition = index
			break
		}
	}

	m.CurrentRepoStashInfoList = list.New(items, gitStashItemDelegate{}, m.WindowLeftPanelWidth, m.StashComponentPanelHeight)
	m.CurrentRepoStashInfoList.SetShowPagination(false)
	m.CurrentRepoStashInfoList.SetShowStatusBar(false)
	m.CurrentRepoStashInfoList.SetFilteringEnabled(false)
	m.CurrentRepoStashInfoList.SetShowFilter(false)
	m.CurrentRepoStashInfoList.Title = utils.TruncateString(fmt.Sprintf("[3] \ueaf7 %s:", i18n.LANGUAGEMAPPING.Stash), m.WindowLeftPanelWidth-constant.ListItemOrTitleWidthPad-2)
	m.CurrentRepoStashInfoList.Styles.Title = style.TitleStyle
	m.CurrentRepoStashInfoList.Styles.TitleBar = style.NewStyle

	// Custom Help Model for Count Display
	m.CurrentRepoStashInfoList.SetShowHelp(true)
	m.CurrentRepoStashInfoList.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	m.CurrentRepoStashInfoList.AdditionalShortHelpKeys = func() []key.Binding {
		currentIndex := m.CurrentRepoStashInfoList.Index() + 1
		totalCount := len(m.CurrentRepoStashInfoList.Items())
		countStr := fmt.Sprintf("%d/%d", currentIndex, totalCount)
		countStr = utils.TruncateString(countStr, m.WindowLeftPanelWidth-constant.ListItemOrTitleWidthPad-2)
		if totalCount == 0 {
			countStr = "0/0"
		}
		return []key.Binding{
			key.NewBinding(
				key.WithKeys(countStr),
				key.WithHelp(countStr, ""),
			),
		}
	}

	if len(items) < 1 {
		return true
	}

	if selectedFilesPosition >= 0 {
		m.CurrentRepoStashInfoList.Select(selectedFilesPosition)
		m.ListNavigationIndexPosition.StashComponent = selectedFilesPosition
	} else {
		if m.ListNavigationIndexPosition.StashComponent > len(m.CurrentRepoStashInfoList.Items())-1 {
			m.CurrentRepoStashInfoList.Select(len(m.CurrentRepoStashInfoList.Items()) - 1)
			m.ListNavigationIndexPosition.StashComponent = len(m.CurrentRepoStashInfoList.Items()) - 1
		} else {
			m.CurrentRepoStashInfoList.Select(m.ListNavigationIndexPosition.StashComponent)
		}
	}

	if previousSelectedStash == m.CurrentRepoStashInfoList.SelectedItem() {
		return false
	}
	return true
}

// init the viewport pop up for showing info of global key binding
func initGlobalKeyBindingPopUpModel(m *types.GittiModel) {
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(min(constant.PopUpGlobalKeyBindingViewPortHeight, int(float64(m.Height)*0.8)))
	vp.SetWidth(min(constant.MaxGlobalKeyBindingPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	m.PopUpModel = &GlobalKeyBindingPopUpModel{
		GlobalKeyBindingViewport: vp,
	}
}

// init the popup model for git commit
func initGitCommitPopUpModel(m *types.GittiModel) {
	commitMsg := ""
	commitDesc := ""
	commitMsgPlaceholder := i18n.LANGUAGEMAPPING.CommitPopUpMessageInputPlaceHolder
	commitDescPlaceholder := i18n.LANGUAGEMAPPING.CommitPopUpCommitDescriptionInputPlaceHolder

	CommitMessageTextInput := textinput.New()
	CommitMessageTextInput.SetValue(commitMsg)
	CommitMessageTextInput.Placeholder = commitMsgPlaceholder
	CommitMessageTextInput.Focus()
	CommitMessageTextInput.SetVirtualCursor(true)

	CommitDescriptionTextAreaInput := textarea.New()
	CommitDescriptionTextAreaInput.SetValue(commitDesc)
	CommitDescriptionTextAreaInput.ShowLineNumbers = false
	CommitDescriptionTextAreaInput.Placeholder = commitDescPlaceholder
	CommitDescriptionTextAreaInput.SetHeight(4)
	CommitDescriptionTextAreaInput.Blur()

	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpGitCommitOutputViewPortHeight)
	vp.SetWidth(min(constant.MaxCommitPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &GitCommitPopUpModel{
		IsAmendCommit:            false,
		MessageTextInput:         CommitMessageTextInput,
		DescriptionTextAreaInput: CommitDescriptionTextAreaInput,
		TotalInputCount:          2,
		CurrentActiveInputIndex:  1,
		GitCommitOutputViewport:  vp,
		Spinner:                  s,
	}
	popUpModel.InitialCommitStarted.Store(false)
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}

// init the popup model for git amend commit
func initGitAmendCommitPopUpModel(m *types.GittiModel) {
	commitMsgAndDesc := m.GitOperations.GitCommit.GetLatestCommitMsgAndDesc()
	commitMsg := commitMsgAndDesc.Message
	commitDesc := commitMsgAndDesc.Description
	commitMsgPlaceholder := i18n.LANGUAGEMAPPING.CommitPopUpMessageInputPlaceHolderAmendVersion
	commitDescPlaceholder := i18n.LANGUAGEMAPPING.CommitPopUpCommitDescriptionInputPlaceHolderAmendVersion

	CommitMessageTextInput := textinput.New()
	CommitMessageTextInput.SetValue(commitMsg)
	CommitMessageTextInput.Placeholder = commitMsgPlaceholder
	CommitMessageTextInput.Focus()
	CommitMessageTextInput.SetVirtualCursor(true)

	CommitDescriptionTextAreaInput := textarea.New()
	CommitDescriptionTextAreaInput.SetValue(commitDesc)
	CommitDescriptionTextAreaInput.ShowLineNumbers = false
	CommitDescriptionTextAreaInput.Placeholder = commitDescPlaceholder
	CommitDescriptionTextAreaInput.SetHeight(4)
	CommitDescriptionTextAreaInput.Blur()

	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpGitAmendCommitOutputViewPortHeight)
	vp.SetWidth(min(constant.MaxAmendCommitPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &GitAmendCommitPopUpModel{
		IsAmendCommit:                true,
		MessageTextInput:             CommitMessageTextInput,
		DescriptionTextAreaInput:     CommitDescriptionTextAreaInput,
		TotalInputCount:              2,
		CurrentActiveInputIndex:      1,
		GitAmendCommitOutputViewport: vp,
		Spinner:                      s,
	}
	popUpModel.InitialCommitStarted.Store(false)
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}

// init the popup model for prompting user to add remote origin
func initAddRemotePromptPopUpModel(m *types.GittiModel, noInitialRemote bool) {
	RemoteNameTextInput := textinput.New()
	if noInitialRemote {
		RemoteNameTextInput.SetValue("origin")
	}
	RemoteNameTextInput.Placeholder = i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteNamePlaceHolder
	RemoteNameTextInput.Focus()
	RemoteNameTextInput.SetVirtualCursor(true)

	RemoteUrlTextInput := textinput.New()
	RemoteUrlTextInput.Placeholder = i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteUrlPlaceHolder
	RemoteUrlTextInput.Blur()
	RemoteUrlTextInput.SetVirtualCursor(true)

	// for git add remote output viewport, we will not have any interaction for it as usually it will be a one line for error log or also for our custom success message
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpAddRemoteOutputViewPortHeight)
	vp.SetWidth(min(constant.MaxAddRemotePromptPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	popUpModel := &AddRemotePromptPopUpModel{
		RemoteNameTextInput:     RemoteNameTextInput,
		RemoteUrlTextInput:      RemoteUrlTextInput,
		TotalInputCount:         2,
		CurrentActiveInputIndex: 1,
		AddRemoteOutputViewport: vp,
		NoInitialRemote:         noInitialRemote,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}

// init the popup model for push output log
func initGitRemotePushPopUpModel(m *types.GittiModel) {
	// for git push output viewport,
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpGitRemotePushOutputViewportHeight)
	vp.SetWidth(min(constant.MaxGitRemotePushPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &GitRemotePushPopUpModel{
		GitRemotePushOutputViewport: vp,
		Spinner:                     s,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}

func initGitRemotePushPopUpModelAndStartGitRemotePushService(m *types.GittiModel, remoteName string, pushType string) (*types.GittiModel, tea.Cmd) {
	m.GitOperations.GitCommit.ClearGitRemotePushOutput()
	if popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel); !ok {
		initGitRemotePushPopUpModel(m)
	} else {
		popUp.GitRemotePushOutputViewport.SetContent("")
	}
	// then push it after init the git remote push pop up model
	gitRemotePushService(m, remoteName, pushType)
	// Start spinner ticking
	if pushPopup, ok := m.PopUpModel.(*GitRemotePushPopUpModel); ok {
		return m, pushPopup.Spinner.Tick
	}
	return m, nil
}

// init the popup model to choose remote to push to
func initGitRemotePushChooseRemotePopUpModel(m *types.GittiModel, remoteList []git.GitRemoteInfo) {
	items := make([]list.Item, 0, len(remoteList))
	for _, remote := range remoteList {
		items = append(items, gitRemoteItem(remote))
	}
	width := (min(constant.MaxChooseRemotePopUpWidth, int(float64(m.Width)*0.8)) - 4)
	rL := list.New(items, gitRemoteItemDelegate{}, width, constant.PopUpChooseRemoteHeight)
	rL.SetShowPagination(false)
	rL.SetShowStatusBar(false)
	rL.SetFilteringEnabled(false)
	rL.SetShowHelp(false)
	rL.SetShowTitle(false)

	m.PopUpModel = &ChooseRemotePopUpModel{
		RemoteList: rL,
	}
}

// init the popup model for choosing push type
func initChoosePushTypePopUpModel(m *types.GittiModel, remoteName string) {
	pushTypeOption := []gitPushOptionItem{
		{
			Name:     i18n.LANGUAGEMAPPING.NormalPush,
			Info:     "git push",
			pushType: git.PUSH,
		},
		{
			Name:     i18n.LANGUAGEMAPPING.ForcePushSafe,
			Info:     "git push --force-with-lease",
			pushType: git.FORCEPUSHSAFE,
		},
		{
			Name:     i18n.LANGUAGEMAPPING.ForcePushDangerous,
			Info:     "git push --force",
			pushType: git.FORCEPUSHDANGEROUS,
		},
	}

	items := make([]list.Item, 0, len(pushTypeOption))
	for _, pushOption := range pushTypeOption {
		items = append(items, gitPushOptionItem(pushOption))
	}
	width := (min(constant.MaxChoosePushTypePopUpWidth, int(float64(m.Width)*0.8)) - 4)
	pOL := list.New(items, gitPushOptionDelegate{}, width, constant.PopUpChoosePushTypeHeight)
	pOL.SetShowPagination(false)
	pOL.SetShowStatusBar(false)
	pOL.SetFilteringEnabled(false)
	pOL.SetShowHelp(false)
	pOL.SetShowTitle(false)

	m.PopUpModel = &ChoosePushTypePopUpModel{
		PushOptionList: pOL,
		RemoteName:     remoteName,
	}
}

// init the popup model for creating a new branch
func initCreateNewBranchPopUpModel(m *types.GittiModel, createType string) {
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
func initChooseNewBranchTypePopUpModel(m *types.GittiModel) {
	newBranchTypeOption := []gitNewBranchTypeOptionItem{
		{
			Name:          i18n.LANGUAGEMAPPING.CreateNewBranchTitle,
			Info:          i18n.LANGUAGEMAPPING.CreateNewBranchDescription,
			newBranchType: git.NEWBRANCH,
		},
		{
			Name:          i18n.LANGUAGEMAPPING.CreateNewBranchAndSwitchTitle,
			Info:          i18n.LANGUAGEMAPPING.CreateNewBranchAndSwitchDescription,
			newBranchType: git.NEWBRANCHANDSWITCH,
		},
	}

	items := make([]list.Item, 0, len(newBranchTypeOption))
	for _, newBranchOption := range newBranchTypeOption {
		items = append(items, gitNewBranchTypeOptionItem(newBranchOption))
	}
	width := (min(constant.MaxChooseNewBranchTypePopUpWidth, int(float64(m.Width)*0.8)) - 4)
	nBTOL := list.New(items, gitNewBranchTypeOptionDelegate{}, width, constant.PopUpChooseNewBranchTypeHeight)
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
func initChooseSwitchBranchTypePopUpModel(m *types.GittiModel, branchName string) {
	switchBranchTypeOption := []gitSwitchBranchTypeOptionItem{
		{
			Name:             i18n.LANGUAGEMAPPING.SwitchBranchTitle,
			Info:             i18n.LANGUAGEMAPPING.SwitchBranchDescription,
			switchBranchType: git.SWITCHBRANCH,
		},
		{
			Name:             i18n.LANGUAGEMAPPING.SwitchBranchWithChangesTitle,
			Info:             i18n.LANGUAGEMAPPING.SwitchBranchWithChangesDescription,
			switchBranchType: git.SWITCHBRANCHWITHCHANGES,
		},
	}

	items := make([]list.Item, 0, len(switchBranchTypeOption))
	for _, switchBranchOption := range switchBranchTypeOption {
		items = append(items, gitSwitchBranchTypeOptionItem(switchBranchOption))
	}

	width := (min(constant.MaxChooseSwitchBranchTypePopUpWidth, int(float64(m.Width)*0.8)) - 4)
	sBTOL := list.New(items, gitSwitchBranchTypeOptionDelegate{}, width, constant.PopUpChooseSwitchBranchTypeHeight)
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

func initSwitchBranchOutputPopUpModel(m *types.GittiModel, branchName string, switchType string) {
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

func initChooseGitPullTypePopUp(m *types.GittiModel) {
	pullTypeOption := []gitPullTypeOptionItem{
		{
			Name:     i18n.LANGUAGEMAPPING.GitPullOption,
			Info:     "git pull --no-edit",
			PullType: git.GITPULL,
		},
		{
			Name:     i18n.LANGUAGEMAPPING.GitPullRebaseOption,
			Info:     "git pull --rebase --autostash --no-edit",
			PullType: git.GITPULLREBASE,
		},
		{
			Name:     i18n.LANGUAGEMAPPING.GitPullMergeOption,
			Info:     "git pull --no-rebase --no-edit",
			PullType: git.GITPULLMERGE,
		},
	}

	items := make([]list.Item, 0, len(pullTypeOption))
	for _, pullOption := range pullTypeOption {
		items = append(items, gitPullTypeOptionItem(pullOption))
	}

	width := (min(constant.MaxChooseGitPullTypePopUpWidth, int(float64(m.Width)*0.8)) - 4)
	cGPOL := list.New(items, gitPullTypeOptionDelegate{}, width, constant.PopUpChooseGitPullTypeHeight)
	cGPOL.SetShowPagination(false)
	cGPOL.SetShowStatusBar(false)
	cGPOL.SetFilteringEnabled(false)
	cGPOL.SetShowHelp(false)
	cGPOL.SetShowTitle(false)

	popUpModel := &ChooseGitPullTypePopUpModel{
		PullTypeOptionList: cGPOL,
	}

	m.PopUpModel = popUpModel
}

func initGitPullOutputPopUpModel(m *types.GittiModel) {
	// for git pull output viewport
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpSwitchBranchOutputViewPortHeight)
	vp.SetWidth(min(constant.MaxSwitchBranchOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &GitPullOutputPopUpModel{
		GitPullOutputViewport: vp,
		Spinner:               s,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}

func initGitStashMessagePopUpModel(m *types.GittiModel, filePathName string, stashType string) {
	stashMessageTextInput := textinput.New()
	stashMessageTextInput.Placeholder = i18n.LANGUAGEMAPPING.GitStashMessagePlaceholder
	stashMessageTextInput.Focus()
	stashMessageTextInput.SetVirtualCursor(true)
	stashMessageTextInput.SetWidth(min(constant.MaxGitStashMessagePopUpWidth, int(float64(m.Width)*0.8)) - 4)

	popUpModel := &GitStashMessagePopUpModel{
		StashMessageInput: stashMessageTextInput,
		FilePathName:      filePathName,
		StashType:         stashType,
	}

	m.PopUpModel = popUpModel
}

// for discard option list popup
func initGitDiscardTypeOptionPopUp(m *types.GittiModel, filePathName string, newlyAddedOrCopiedFile bool, renameFile bool) {
	discardTypeOption := []gitDiscardTypeOptionItem{
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
		discardTypeOption = []gitDiscardTypeOptionItem{
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
		discardTypeOption = []gitDiscardTypeOptionItem{
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
		items = append(items, gitDiscardTypeOptionItem(discardOption))
	}

	width := (min(constant.MaxGitDiscardTypeOptionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	gDTOL := list.New(items, gitDiscardTypeOptionDelegate{}, width, constant.PopUpGitDiscardTypeOptionHeight)
	gDTOL.SetShowPagination(false)
	gDTOL.SetShowStatusBar(false)
	gDTOL.SetFilteringEnabled(false)
	gDTOL.SetShowHelp(false)
	gDTOL.SetShowTitle(false)

	popUpModel := &GitDiscardTypeOptionPopUpModel{
		DiscardTypeOptionList: gDTOL,
		FilePathName:          filePathName,
	}

	m.PopUpModel = popUpModel
}

// for discard confirm prompt
func initGitDiscardConfirmPromptPopup(m *types.GittiModel, filePathName string, discardType string) {
	popUpModel := &GitDiscardConfirmPromptPopUpModel{
		FilePathName: filePathName,
		DiscardType:  discardType,
	}
	m.PopUpModel = popUpModel
}

// for git stash output popup
func initGitStashOperationOutputPopUpModel(m *types.GittiModel, stashOperationType string) {
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpGitStashOperationOutputViewPortHeight)
	vp.SetWidth(min(constant.MaxGitStashOperationOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &GitStashOperationOutputPopUpModel{
		StashOperationType:              stashOperationType,
		GitStashOperationOutputViewport: vp,
		Spinner:                         s,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	m.PopUpModel = popUpModel
}

// for git stash operation confirm prompt
func initGitStashConfirmPromptPopUpModel(m *types.GittiModel, stashOperationType string, filePathName string, stashId string, stashMessage string) {
	popUpModel := &GitStashConfirmPromptPopUpModel{
		StashOperationType: stashOperationType,
		FilePathName:       filePathName,
		StashId:            stashId,
		StashMessage:       stashMessage,
	}
	m.PopUpModel = popUpModel
}

// for resolve conflict option list popup
func initGitResolveConflictOptionPopUpModel(m *types.GittiModel, filePathName string) {
	resolveConflictOption := []gitResolveConflictOptionItem{
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
		items = append(items, gitResolveConflictOptionItem(resolveConflictOption))
	}

	width := (min(constant.MaxGitResolveConflictOptionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	rROL := list.New(items, gitResolveConflictOptionDelegate{}, width, constant.PopUpGitResolveConflictOptionPopUpHeight)
	rROL.SetShowPagination(false)
	rROL.SetShowStatusBar(false)
	rROL.SetFilteringEnabled(false)
	rROL.SetShowHelp(false)
	rROL.SetShowTitle(false)

	popUpModel := &GitResolveConflictOptionPopUpModel{
		ResolveConflictOptionList: rROL,
		FilePathName:              filePathName,
	}

	m.PopUpModel = popUpModel
}
