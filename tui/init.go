package tui

import (
	"fmt"

	"gitti/api/git"
	"gitti/i18n"
	"gitti/tui/constant"
	"gitti/tui/style"
	"gitti/tui/utils"

	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/spinner"
	"github.com/charmbracelet/bubbles/v2/textarea"
	"github.com/charmbracelet/bubbles/v2/textinput"
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/google/uuid"
) // this was for various components part init or reinit function due to update or newly create

func initBranchList(m *GittiModel) {
	items := []list.Item{
		gitBranchItem(m.GitState.GitBranch.CurrentCheckOut()),
	}

	for _, branch := range m.GitState.GitBranch.AllBranches() {
		items = append(items, gitBranchItem(branch))
	}

	m.CurrentRepoBranchesInfoList = list.New(items, gitBranchItemDelegate{}, m.HomeTabLeftPanelWidth, m.HomeTabLocalBranchesPanelHeight)
	m.CurrentRepoBranchesInfoList.SetShowStatusBar(false)
	m.CurrentRepoBranchesInfoList.SetFilteringEnabled(false)
	m.CurrentRepoBranchesInfoList.SetShowHelp(false)
	m.CurrentRepoBranchesInfoList.Title = utils.TruncateString(fmt.Sprintf("[0]  %s:", i18n.LANGUAGEMAPPING.Branches), m.HomeTabLeftPanelWidth-constant.ListItemOrTitleWidthPad-2)
	m.CurrentRepoBranchesInfoList.Styles.Title = style.TitleStyle
	m.CurrentRepoBranchesInfoList.Styles.PaginationStyle = style.PaginationStyle

	if m.NavigationIndexPosition.LocalBranchComponent > len(m.CurrentRepoBranchesInfoList.Items())-1 {
		m.CurrentRepoBranchesInfoList.Select(len(m.CurrentRepoBranchesInfoList.Items()) - 1)
	} else {
		m.CurrentRepoBranchesInfoList.Select(m.NavigationIndexPosition.LocalBranchComponent)
	}
}

func initModifiedFilesList(m *GittiModel) {
	latestModifiedFilesArray := m.GitState.GitFiles.FilesStatus()
	items := make([]list.Item, 0, len(latestModifiedFilesArray))
	for _, modifiedFile := range latestModifiedFilesArray {
		items = append(items, gitModifiedFilesItem(modifiedFile))
	}

	// get the previous selected file and see if it was within the new list if yes get the latest position of the previous selected file
	previousSelectedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
	selectedFilesPosition := -1

	for index, item := range items {
		if item == previousSelectedFile {
			selectedFilesPosition = index
			break
		}
	}

	m.CurrentRepoModifiedFilesInfoList = list.New(items, gitModifiedFilesItemDelegate{}, m.HomeTabLeftPanelWidth, m.HomeTabChangedFilesPanelHeight)
	m.CurrentRepoModifiedFilesInfoList.SetShowStatusBar(false)
	m.CurrentRepoModifiedFilesInfoList.SetFilteringEnabled(false)
	m.CurrentRepoModifiedFilesInfoList.SetShowHelp(false)
	m.CurrentRepoModifiedFilesInfoList.SetShowPagination(false)
	m.CurrentRepoModifiedFilesInfoList.Title = utils.TruncateString(fmt.Sprintf("[1] 📄%s:", i18n.LANGUAGEMAPPING.ModifiedFiles), m.HomeTabLeftPanelWidth-constant.ListItemOrTitleWidthPad-2)
	m.CurrentRepoModifiedFilesInfoList.Styles.Title = style.TitleStyle

	if len(items) < 1 {
		return
	}

	if selectedFilesPosition >= 0 {
		m.CurrentRepoModifiedFilesInfoList.Select(selectedFilesPosition)
		m.NavigationIndexPosition.ModifiedFilesComponent = selectedFilesPosition
	} else {
		if m.NavigationIndexPosition.ModifiedFilesComponent > len(m.CurrentRepoModifiedFilesInfoList.Items())-1 {
			m.CurrentRepoModifiedFilesInfoList.Select(len(m.CurrentRepoModifiedFilesInfoList.Items()) - 1)
		} else {
			m.CurrentRepoModifiedFilesInfoList.Select(m.NavigationIndexPosition.ModifiedFilesComponent)
		}
	}
}

// reinit and render diff file viewport
func reinitAndRenderModifiedFileDiffViewPort(m *GittiModel) {
	m.DetailPanelViewportOffset = 0
	m.DetailPanelViewport.SetXOffset(0)
	m.DetailPanelViewport.SetYOffset(0)
	renderDetailPanelViewPort(m)
}

// init the popup model for git commit
func initGitCommitPopUpModel(m *GittiModel) {
	CommitMessageTextInput := textinput.New()
	CommitMessageTextInput.Placeholder = i18n.LANGUAGEMAPPING.CommitPopUpMessageInputPlaceHolder
	CommitMessageTextInput.Focus()
	CommitMessageTextInput.VirtualCursor = true

	CommitDescriptionTextAreaInput := textarea.New()
	CommitDescriptionTextAreaInput.ShowLineNumbers = false
	CommitDescriptionTextAreaInput.Placeholder = i18n.LANGUAGEMAPPING.CommitPopUpCommitDescriptionInputPlaceHolder
	CommitDescriptionTextAreaInput.SetHeight(5)
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

	// Generate a unique UUID for this popup session
	newSessionID := uuid.New()

	popUpModel := &GitCommitPopUpModel{
		MessageTextInput:         CommitMessageTextInput,
		DescriptionTextAreaInput: CommitDescriptionTextAreaInput,
		TotalInputCount:          2,
		CurrentActiveInputIndex:  1,
		GitCommitOutputViewport:  vp,
		Spinner:                  s,
		SessionID:                newSessionID,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}

// init the popup model for prompting user to add remote origin
func initAddRemotePromptPopUpModel(m *GittiModel, noInitialRemote bool) {
	RemoteNameTextInput := textinput.New()
	RemoteNameTextInput.Placeholder = i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteNamePlaceHolder
	RemoteNameTextInput.Focus()
	RemoteNameTextInput.VirtualCursor = true

	RemoteUrlTextInput := textinput.New()
	RemoteUrlTextInput.Placeholder = i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteUrlPlaceHolder
	RemoteUrlTextInput.Blur()
	RemoteUrlTextInput.VirtualCursor = true

	// for git add remote output viewport, we will not have any interaction for it as usually it will be a one line for error log or also for our custom success message
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpAddRemoteOutputViewPortHeight)
	vp.SetWidth(min(constant.MaxAddRemotePromptPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	// Generate a unique UUID for this popup session
	newSessionID := uuid.New()

	popUpModel := &AddRemotePromptPopUpModel{
		RemoteNameTextInput:     RemoteNameTextInput,
		RemoteUrlTextInput:      RemoteUrlTextInput,
		TotalInputCount:         2,
		CurrentActiveInputIndex: 1,
		AddRemoteOutputViewport: vp,
		NoInitialRemote:         noInitialRemote,
		SessionID:               newSessionID,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}

// init the popup model for push output log
func initGitRemotePushPopUpModel(m *GittiModel) {
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

	// Generate a unique UUID for this popup session
	newSessionID := uuid.New()

	popUpModel := &GitRemotePushPopUpModel{
		GitRemotePushOutputViewport: vp,
		Spinner:                     s,
		SessionID:                   newSessionID,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}

func initGitRemotePushPopUpModelAndStartGitRemotePushService(m *GittiModel, remoteName string, pushType string) (*GittiModel, tea.Cmd) {
	m.GitState.GitCommit.ClearGitRemotePushOutput()
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
func initGitRemotePushChooseRemotePopUpModel(m *GittiModel, remoteList []git.GitRemote) {
	items := make([]list.Item, 0, len(remoteList))
	for _, remote := range remoteList {
		items = append(items, gitRemoteItem(remote))
	}
	width := (min(constant.MaxChooseRemotePopUpWidth, int(float64(m.Width)*0.8)) - 4)
	rL := list.New(items, gitRemoteItemDelegate{}, width, constant.PopUpChooseRemoteHeight)
	rL.SetShowStatusBar(false)
	rL.SetFilteringEnabled(false)
	rL.SetShowHelp(false)
	rL.SetShowTitle(false)

	m.PopUpModel = &ChooseRemotePopUpModel{
		RemoteList: rL,
	}
}

// init the popup model for choosing push type
func initChoosePushTypePopUpModel(m *GittiModel, remoteName string) {
	pushTypeOption := []gitPushOptionItem{
		{
			Name:     i18n.LANGUAGEMAPPING.NormalPush,
			Info:     "git push",
			pushType: git.PUSH,
		},
		{
			Name:     i18n.LANGUAGEMAPPING.ForcePushSafe,
			Info:     "git push --force",
			pushType: git.FORCEPUSHSAFE,
		},
		{
			Name:     i18n.LANGUAGEMAPPING.ForcePushDangerous,
			Info:     "git push --force-with-lease",
			pushType: git.FORCEPUSHDANGEROUS,
		},
	}

	items := make([]list.Item, 0, len(pushTypeOption))
	for _, pushOption := range pushTypeOption {
		items = append(items, gitPushOptionItem(pushOption))
	}
	width := (min(constant.MaxChoosePushTypePopUpWidth, int(float64(m.Width)*0.8)) - 4)
	pOL := list.New(items, gitPushOptionDelegate{}, width, constant.PopUpChoosePushTypeHeight)
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
func initCreateNewBranchPopUpModel(m *GittiModel, createType string) {
	NewBranchNameInput := textinput.New()
	NewBranchNameInput.Placeholder = i18n.LANGUAGEMAPPING.CreateNewBranchPrompt
	NewBranchNameInput.Focus()
	NewBranchNameInput.VirtualCursor = true

	m.PopUpModel = &CreateNewBranchPopUpModel{
		NewBranchNameInput: NewBranchNameInput,
		CreateType:         createType,
	}
}

// init the popup model for choosing new branch creation option
func initChooseNewBranchTypePopUpModel(m *GittiModel) {
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
	nBTOL.SetShowStatusBar(false)
	nBTOL.SetFilteringEnabled(false)
	nBTOL.SetShowHelp(false)
	nBTOL.SetShowTitle(false)

	m.PopUpModel = &ChooseNewBranchTypeOptionPopUpModel{
		NewBranchTypeOptionList: nBTOL,
	}
}

// init the popup model for switching branch
func initChooseSwitchBranchTypePopUpModel(m *GittiModel, branchName string) {
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
	sBTOL.SetShowStatusBar(false)
	sBTOL.SetFilteringEnabled(false)
	sBTOL.SetShowHelp(false)
	sBTOL.SetShowTitle(false)

	m.PopUpModel = &ChooseSwitchBranchTypePopUpModel{
		SwitchTypeOptionList: sBTOL,
		BranchName:           branchName,
	}
}

func initSwitchBranchOutputPopUpModel(m *GittiModel, branchName string, switchType string) {
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

func initChooseGitPullTypePopUp(m *GittiModel) {
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
	cGPOL.SetShowStatusBar(false)
	cGPOL.SetFilteringEnabled(false)
	cGPOL.SetShowHelp(false)
	cGPOL.SetShowTitle(false)

	popUpModel := &ChooseGitPullTypePopUpModel{
		PullTypeOptionList: cGPOL,
	}

	m.PopUpModel = popUpModel
}

func initGitPullOutputPopUpModel(m *GittiModel) {
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
		Spinner:                    s,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}
