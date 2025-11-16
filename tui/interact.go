package tui

import (
	"gitti/api"
	"gitti/api/git"
	"gitti/tui/constant"
	"gitti/utils"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/google/uuid"
)

const (
	AUTHOR_GITHUB   = "https://github.com/gohyuhan"
	AUTHOR_LINKEDIN = "https://my.linkedin.com/in/yu-han-goh-209480200"
)

// the function to handle bubbletea key interactions
func gittiKeyInteraction(msg tea.KeyMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	// global key binding
	switch msg.String() {
	case "ctrl+c":
		api.GITDAEMON.Stop()
		return m, tea.Quit
	case "ctrl+s":
		gitStageAllChangesService(m)
		return m, nil
	case "ctrl+u":
		gitUnstageAllChangesService(m)
		return m, nil
	case "ctrl+g":
		utils.OpenBrowser(AUTHOR_GITHUB)
		return m, nil
	case "ctrl+l":
		utils.OpenBrowser(AUTHOR_LINKEDIN)
		return m, nil
	}

	if m.IsTyping.Load() {
		return handleTypingKeyBindingInteraction(msg, m)
	} else {
		return handleNonTypingGlobalKeyBindingInteraction(msg, m)
	}
}

// typing is currently only on pop up model, so we can safely process it without checking if they were on pop up or not
func handleTypingKeyBindingInteraction(msg tea.KeyMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		switch m.PopUpType {
		case constant.CommitPopUp:
			gitCommitCancelService(m)
		case constant.AmendCommitPopUp:
			gitAmendCommitCancelService(m)
		case constant.AddRemotePromptPopUp:
			gitAddRemoteCancelService(m)
		case constant.CreateNewBranchPopUp:
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpType = constant.NoPopUp
			m.PopUpModel = nil
		case constant.GitStashMessagePopUp:
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpType = constant.NoPopUp
			m.PopUpModel = nil
		}
		return m, nil

		// in typing mode, tab is move to next input
	case "tab":
		switch m.PopUpType {
		case constant.CommitPopUp:
			popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
			if ok {
				popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
				switch popUp.CurrentActiveInputIndex {
				case 1:
					popUp.MessageTextInput.Focus()
					popUp.DescriptionTextAreaInput.Blur()
				case 2:
					popUp.MessageTextInput.Blur()
					popUp.DescriptionTextAreaInput.Focus()
				}
			}
		case constant.AmendCommitPopUp:
			popUp, ok := m.PopUpModel.(*GitAmendCommitPopUpModel)
			if ok {
				popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
				switch popUp.CurrentActiveInputIndex {
				case 1:
					popUp.MessageTextInput.Focus()
					popUp.DescriptionTextAreaInput.Blur()
				case 2:
					popUp.MessageTextInput.Blur()
					popUp.DescriptionTextAreaInput.Focus()
				}
			}
		case constant.AddRemotePromptPopUp:
			popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
			if ok {
				popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
				switch popUp.CurrentActiveInputIndex {
				case 1:
					popUp.RemoteNameTextInput.Focus()
					popUp.RemoteUrlTextInput.Blur()
				case 2:
					popUp.RemoteNameTextInput.Blur()
					popUp.RemoteUrlTextInput.Focus()
				}
			}
		}
		return m, nil

	// in typing mode, shift+tab is move to previous input
	case "shift+tab":
		switch m.PopUpType {
		case constant.CommitPopUp:
			popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
			if ok {
				popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
				switch popUp.CurrentActiveInputIndex {
				case 1:
					popUp.MessageTextInput.Focus()
					popUp.DescriptionTextAreaInput.Blur()
				case 2:
					popUp.MessageTextInput.Blur()
					popUp.DescriptionTextAreaInput.Focus()
				}
			}
		case constant.AmendCommitPopUp:
			popUp, ok := m.PopUpModel.(*GitAmendCommitPopUpModel)
			if ok {
				popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
				switch popUp.CurrentActiveInputIndex {
				case 1:
					popUp.MessageTextInput.Focus()
					popUp.DescriptionTextAreaInput.Blur()
				case 2:
					popUp.MessageTextInput.Blur()
					popUp.DescriptionTextAreaInput.Focus()
				}
			}
		case constant.AddRemotePromptPopUp:
			popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
			if ok {
				popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
				switch popUp.CurrentActiveInputIndex {
				case 1:
					popUp.RemoteNameTextInput.Focus()
					popUp.RemoteUrlTextInput.Blur()
				case 2:
					popUp.RemoteNameTextInput.Blur()
					popUp.RemoteUrlTextInput.Focus()
				}
			}

		}

	case "ctrl+enter":
		switch m.PopUpType {
		case constant.CommitPopUp:
			popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
			if ok {
				// once they start for commit process, reinit the input focus
				popUp.MessageTextInput.Focus()
				popUp.DescriptionTextAreaInput.Blur()
				popUp.CurrentActiveInputIndex = 1
				// start a seperate thread commit them and set the value of msg and desc to "" if committed successfully
				// also do not start any git operation is message is no provided
				if !popUp.IsProcessing.Load() && len(popUp.MessageTextInput.Value()) > 0 {
					gitCommitService(m, popUp.IsAmendCommit)
					// Start spinner ticking
					return m, popUp.Spinner.Tick
				}
			}
		case constant.AmendCommitPopUp:
			popUp, ok := m.PopUpModel.(*GitAmendCommitPopUpModel)
			if ok {
				// once they start for amend commit process, reinit the input focus
				popUp.MessageTextInput.Focus()
				popUp.DescriptionTextAreaInput.Blur()
				popUp.CurrentActiveInputIndex = 1
				if !popUp.IsProcessing.Load() && len(popUp.MessageTextInput.Value()) > 0 {
					gitAmendCommitService(m, popUp.IsAmendCommit)
					// Start spinner ticking
					return m, popUp.Spinner.Tick
				}
			}

		case constant.AddRemotePromptPopUp:
			popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
			if ok {
				// once they start for commit process, reinit the input focus
				popUp.RemoteNameTextInput.Focus()
				popUp.RemoteUrlTextInput.Blur()
				popUp.CurrentActiveInputIndex = 1
				// start a seperate thread that stage the current selected files and commit them and set the value of msg and desc to "" if committed successfully
				// also do not start any git operation is message is no provided
				if !popUp.IsProcessing.Load() && len(popUp.RemoteNameTextInput.Value()) > 0 && len(popUp.RemoteUrlTextInput.Value()) > 0 {
					gitAddRemoteService(m)
				}
			}

		case constant.CreateNewBranchPopUp:
			popUp, ok := m.PopUpModel.(*CreateNewBranchPopUpModel)
			if ok {
				// we direclty close the pop up and trigger the branch creation operation
				validBranchName, _ := api.IsBranchNameValid(popUp.NewBranchNameInput.Value())
				if len(validBranchName) > 0 {
					switch popUp.CreateType {
					case git.NEWBRANCH:
						gitCreateNewBranchService(m, validBranchName)
					case git.NEWBRANCHANDSWITCH:
						gitCreateNewBranchAndSwitchService(m, validBranchName)
					}
					m.ShowPopUp.Store(false)
					m.IsTyping.Store(false)
					m.PopUpType = constant.NoPopUp
					m.PopUpModel = nil
				}
			}

		case constant.GitStashMessagePopUp:
			popUp, ok := m.PopUpModel.(*GitStashMessagePopUpModel)
			if ok {
				msg := popUp.StashMessageInput.Value()
				switch popUp.StashType {
				case git.STASHALL:
					gitStashAllService(m, msg)
				case git.STASHINDIVIDUAL:
					gitStashIndividualFileService(m, popUp.FilePathName, msg)
				}
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		}
		return m, nil
	}

	// for input typing update
	switch m.PopUpType {
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				var cmd tea.Cmd
				popUp.MessageTextInput, cmd = popUp.MessageTextInput.Update(msg)
				return m, cmd

			case 2:
				var cmd tea.Cmd
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*GitAmendCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				var cmd tea.Cmd
				popUp.MessageTextInput, cmd = popUp.MessageTextInput.Update(msg)
				return m, cmd

			case 2:
				var cmd tea.Cmd
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				var cmd tea.Cmd
				popUp.RemoteNameTextInput, cmd = popUp.RemoteNameTextInput.Update(msg)
				return m, cmd

			case 2:
				var cmd tea.Cmd
				popUp.RemoteUrlTextInput, cmd = popUp.RemoteUrlTextInput.Update(msg)
				return m, cmd
			}
		}
	case constant.CreateNewBranchPopUp:
		popUp, ok := m.PopUpModel.(*CreateNewBranchPopUpModel)
		if ok {
			var cmd tea.Cmd
			popUp.NewBranchNameInput, cmd = popUp.NewBranchNameInput.Update(msg)
			return m, cmd
		}
	case constant.GitStashMessagePopUp:
		popUp, ok := m.PopUpModel.(*GitStashMessagePopUpModel)
		if ok {
			var cmd tea.Cmd
			popUp.StashMessageInput, cmd = popUp.StashMessageInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func handleNonTypingGlobalKeyBindingInteraction(msg tea.KeyMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "?":
		m.ShowPopUp.Store(true)
		m.IsTyping.Store(false)
		m.PopUpType = constant.GlobalKeyBindingPopUp
		initGlobalKeyBindingPopUpModel(m)
		return m, nil

	case "1":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedComponent != constant.LocalBranchComponent {
				m.CurrentSelectedComponent = constant.LocalBranchComponent
				m.CurrentSelectedComponentIndex = 1
				leftPanelDynamicResize(m)
				renderDetailComponentPanelViewPort(m)
			}
		}
		return m, nil

	case "2":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedComponent != constant.ModifiedFilesComponent {
				m.CurrentSelectedComponent = constant.ModifiedFilesComponent
				m.CurrentSelectedComponentIndex = 2
				leftPanelDynamicResize(m)
				renderDetailComponentPanelViewPort(m)
			}
		}
		return m, nil

	case "3":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedComponent != constant.StashComponent {
				m.CurrentSelectedComponent = constant.StashComponent
				m.CurrentSelectedComponentIndex = 3
				leftPanelDynamicResize(m)
				renderDetailComponentPanelViewPort(m)
			}
		}
		return m, nil

	case "A":
		if !m.ShowPopUp.Load() {
			m.ShowPopUp.Store(true)
			m.PopUpType = constant.AmendCommitPopUp
			m.GitState.GitCommit.ClearGitCommitOutput()

			initGitAmendCommitPopUpModel(m)

			m.IsTyping.Store(true)
		}
		return m, nil

	case "c":
		if !m.ShowPopUp.Load() {
			m.ShowPopUp.Store(true)
			m.PopUpType = constant.CommitPopUp
			m.GitState.GitCommit.ClearGitCommitOutput()

			// if the current pop up model is not commit pop up model, then init it
			if popUp, ok := m.PopUpModel.(*GitCommitPopUpModel); !ok {
				initGitCommitPopUpModel(m)
			} else {
				popUp.GitCommitOutputViewport.SetContent("")
				popUp.SessionID = uuid.New()
			}
			m.IsTyping.Store(true)
		}
		return m, nil

	case "d":
		if !m.ShowPopUp.Load() && m.CurrentSelectedComponent == constant.StashComponent {
			selectedStashId := m.CurrentRepoStashInfoList.SelectedItem()
			gitStashDropService(m, selectedStashId.(gitStashItem).Id)
		}
		return m, nil

	case "n":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedComponent == constant.LocalBranchComponent {
				m.PopUpType = constant.ChooseNewBranchTypePopUp
				m.IsTyping.Store(false)
				m.ShowPopUp.Store(true)
				if _, ok := m.PopUpModel.(*ChooseNewBranchTypeOptionPopUpModel); !ok {
					initChooseNewBranchTypePopUpModel(m)
				}
			}
		}
		return m, nil

	case "p":
		if !m.ShowPopUp.Load() {
			// first we need to check if there are any push/pull origin origin for this repo
			// if not we prompt the user to add a new remote origin
			if !m.GitState.GitRemote.CheckRemoteExist() {
				m.ShowPopUp.Store(true)
				m.PopUpType = constant.AddRemotePromptPopUp
				// if the current pop up model is not commit pop up model, then init it
				if popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel); !ok {
					initAddRemotePromptPopUpModel(m, true)
				} else {
					popUp.AddRemoteOutputViewport.SetContent("")
				}
				m.IsTyping.Store(true)
			} else {
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				remotes := m.GitState.GitRemote.Remote()
				if len(remotes) == 1 {
					m.PopUpType = constant.ChoosePushTypePopUp
					// if the current pop up model is not commit pop up model, then init it and start git push service
					initChoosePushTypePopUpModel(m, remotes[0].Name)
				} else if len(remotes) > 1 {
					// if remote is more than 1 let user choose which remote to push to first before pushing
					m.PopUpType = constant.ChooseRemotePopUp
					if _, ok := m.PopUpModel.(*ChooseRemotePopUpModel); !ok {
						initGitRemotePushChooseRemotePopUpModel(m, remotes)
					}
				}
			}
		}
		return m, nil

	case "P":
		if !m.ShowPopUp.Load() {
			// first we need to check if there are any push/pull origin for this repo
			// if not we prompt the user to add a new remote origin
			if !m.GitState.GitRemote.CheckRemoteExist() {
				m.ShowPopUp.Store(true)
				m.PopUpType = constant.AddRemotePromptPopUp
				// if the current pop up model is not commit pop up model, then init it
				if popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel); !ok {
					initAddRemotePromptPopUpModel(m, true)
				} else {
					popUp.AddRemoteOutputViewport.SetContent("")
				}
				m.IsTyping.Store(true)
			} else {
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				m.PopUpType = constant.ChooseGitPullTypePopUp
				initChooseGitPullTypePopUp(m)
			}
		}
		return m, nil

	case "s":
		if m.CurrentSelectedComponent == constant.ModifiedFilesComponent {
			currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
			var filePathName string
			if currentSelectedModifiedFile != nil {
				filePathName = currentSelectedModifiedFile.(gitModifiedFilesItem).FilePathname
				m.PopUpType = constant.GitStashMessagePopUp
				m.ShowPopUp.Store(true)
				initGitStashMessagePopUpModel(m, filePathName, git.STASHINDIVIDUAL)
				m.IsTyping.Store(true)
			}
		}
		return m, nil

	case "S":
		if m.CurrentSelectedComponent == constant.ModifiedFilesComponent {
			currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
			var filePathName string
			if currentSelectedModifiedFile != nil {
				filePathName = currentSelectedModifiedFile.(gitModifiedFilesItem).FilePathname
				m.PopUpType = constant.GitStashMessagePopUp
				m.ShowPopUp.Store(true)
				initGitStashMessagePopUpModel(m, filePathName, git.STASHALL)
				m.IsTyping.Store(true)
			}
		}
		return m, nil

	case "q", "Q":
		// only work when there is no pop up
		if !m.ShowPopUp.Load() {
			api.GITDAEMON.Stop()
			return m, tea.Quit
		}

	case "backspace":
		if !m.ShowPopUp.Load() && m.CurrentSelectedComponent == constant.StashComponent {
			selectedStashId := m.CurrentRepoStashInfoList.SelectedItem()
			gitStashPopService(m, selectedStashId.(gitStashItem).Id)
		}
		return m, nil

	case "enter":
		if !m.ShowPopUp.Load() {
			switch m.CurrentSelectedComponent {
			case constant.ModifiedFilesComponent:
				if len(m.CurrentRepoModifiedFilesInfoList.Items()) > 0 {
					m.CurrentSelectedComponent = constant.DetailComponent
					m.DetailPanelParentComponent = constant.ModifiedFilesComponent
				}
			case constant.LocalBranchComponent:
				currentSelectedLocalBranch := m.CurrentRepoBranchesInfoList.SelectedItem().(gitBranchItem)
				// only proceed if the local branch selected is not current checkedout branch
				// we can't switch from current checkout branch to current checkout branch, do we
				if !currentSelectedLocalBranch.IsCheckedOut {
					m.PopUpType = constant.ChooseSwitchBranchTypePopUp
					m.IsTyping.Store(false)
					m.ShowPopUp.Store(true)
					initChooseSwitchBranchTypePopUpModel(m, currentSelectedLocalBranch.BranchName)
				}
			}
		} else {
			switch m.PopUpType {
			case constant.ChooseRemotePopUp:
				popUp, ok := m.PopUpModel.(*ChooseRemotePopUpModel)
				if ok {
					remote := popUp.RemoteList.SelectedItem()
					m.PopUpType = constant.ChoosePushTypePopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					initChoosePushTypePopUpModel(m, remote.(gitRemoteItem).Name)
				}

			case constant.ChoosePushTypePopUp:
				popUp, ok := m.PopUpModel.(*ChoosePushTypePopUpModel)
				if ok {
					m.PopUpType = constant.GitRemotePushPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					selectedOption := popUp.PushOptionList.SelectedItem()
					return initGitRemotePushPopUpModelAndStartGitRemotePushService(m, popUp.RemoteName, selectedOption.(gitPushOptionItem).pushType)
				}

			case constant.ChooseNewBranchTypePopUp:
				popUp, ok := m.PopUpModel.(*ChooseNewBranchTypeOptionPopUpModel)
				if ok {
					m.PopUpType = constant.CreateNewBranchPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(true)
					selectedOption := popUp.NewBranchTypeOptionList.SelectedItem()
					initCreateNewBranchPopUpModel(m, selectedOption.(gitNewBranchTypeOptionItem).newBranchType)
				}

			case constant.ChooseSwitchBranchTypePopUp:
				popUp, ok := m.PopUpModel.(*ChooseSwitchBranchTypePopUpModel)
				if ok {
					m.PopUpType = constant.SwitchBranchOutputPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					selectedOption := popUp.SwitchTypeOptionList.SelectedItem().(gitSwitchBranchTypeOptionItem)
					branchName := popUp.BranchName
					initSwitchBranchOutputPopUpModel(m, branchName, selectedOption.switchBranchType)
					popUp, ok := m.PopUpModel.(*SwitchBranchOutputPopUpModel)
					if ok {
						popUp.IsProcessing.Store(true) // set it directly first
						gitSwitchBranchService(m, branchName, selectedOption.switchBranchType)
						return m, popUp.Spinner.Tick
					}
				}

			case constant.ChooseGitPullTypePopUp:
				popUp, ok := m.PopUpModel.(*ChooseGitPullTypePopUpModel)
				if ok {
					m.PopUpType = constant.GitPullOutputPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					selectedOption := popUp.PullTypeOptionList.SelectedItem().(gitPullTypeOptionItem)
					initGitPullOutputPopUpModel(m)
					popUp, ok := m.PopUpModel.(*GitPullOutputPopUpModel)
					if ok {
						popUp.IsProcessing.Store(true) // set it directly first
						// start the git pull service
						gitPullService(m, selectedOption.PullType)
						return m, popUp.Spinner.Tick
					}
				}
			}
		}
		return m, nil

	case "tab":
		// next component navigation
		nextNavigation := m.CurrentSelectedComponentIndex + 1
		if nextNavigation < len(constant.ComponentNavigationList) {
			m.CurrentSelectedComponentIndex = nextNavigation
			m.CurrentSelectedComponent = constant.ComponentNavigationList[nextNavigation]
			leftPanelDynamicResize(m)
			renderDetailComponentPanelViewPort(m)
		}
		return m, nil
	case "shift+tab":
		// previous component navigation
		previousNavigation := m.CurrentSelectedComponentIndex - 1
		if previousNavigation >= 0 {
			m.CurrentSelectedComponentIndex = previousNavigation
			m.CurrentSelectedComponent = constant.ComponentNavigationList[previousNavigation]
			leftPanelDynamicResize(m)
			renderDetailComponentPanelViewPort(m)
		}
		return m, nil

	case "space":
		if !m.ShowPopUp.Load() {
			switch m.CurrentSelectedComponent {
			case constant.ModifiedFilesComponent:
				currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
				var filePathName string
				if currentSelectedModifiedFile != nil {
					filePathName = currentSelectedModifiedFile.(gitModifiedFilesItem).FilePathname
					gitStageOrUnstageService(m, filePathName)
				}

			case constant.StashComponent:
				selectedStashId := m.CurrentRepoStashInfoList.SelectedItem()
				gitStashApplyService(m, selectedStashId.(gitStashItem).Id)
			}
		}
		return m, nil

	case "esc":
		if m.ShowPopUp.Load() {
			switch m.PopUpType {
			case constant.GlobalKeyBindingPopUp:
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			case constant.GitRemotePushPopUp:
				gitRemotePushCancelService(m)
			case constant.GitPullOutputPopUp:
				gitPullCancelService(m)
			case constant.ChooseRemotePopUp:
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			case constant.ChoosePushTypePopUp:
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			case constant.ChooseNewBranchTypePopUp:
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			case constant.ChooseSwitchBranchTypePopUp:
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			case constant.SwitchBranchOutputPopUp:
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			case constant.ChooseGitPullTypePopUp:
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
			return m, nil
		} else {
			switch m.CurrentSelectedComponent {
			case constant.DetailComponent:
				m.CurrentSelectedComponent = m.DetailPanelParentComponent
				m.DetailPanelParentComponent = ""
			}
		}
		return m, nil

	case "up", "k":
		if !m.ShowPopUp.Load() {
			switch m.CurrentSelectedComponent {
			case constant.LocalBranchComponent:
				// we don't use the list native Update() because we track the current selected index
				if m.CurrentRepoBranchesInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoBranchesInfoList.Index() - 1
					m.CurrentRepoBranchesInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.LocalBranchComponent = latestIndex
				}
			case constant.ModifiedFilesComponent:
				// we don't use the list native Update() because we need to also render the diff view as well as track the current selected index
				if m.CurrentRepoModifiedFilesInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoModifiedFilesInfoList.Index() - 1
					m.CurrentRepoModifiedFilesInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.ModifiedFilesComponent = latestIndex
					reinitAndRenderDetailComponentPanelViewPort(m)
				}
			case constant.StashComponent:
				// we don't use the list native Update() because we need to also render the diff view as well as track the current selected index
				if m.CurrentRepoStashInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoStashInfoList.Index() - 1
					m.CurrentRepoStashInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.StashComponent = latestIndex
				}
			case constant.DetailComponent:
				m.DetailPanelViewport, cmd = m.DetailPanelViewport.Update(msg)
				return m, cmd
			}
		} else {
			return upDownKeyMsgUpdateForPopUp(msg, m)
		}
		return m, nil

	case "down", "j":
		if !m.ShowPopUp.Load() {
			switch m.CurrentSelectedComponent {
			case constant.LocalBranchComponent:
				// we don't use the list native Update() because we track the current selected index
				if m.CurrentRepoBranchesInfoList.Index() < len(m.CurrentRepoBranchesInfoList.Items())-1 {
					latestIndex := m.CurrentRepoBranchesInfoList.Index() + 1
					m.CurrentRepoBranchesInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.LocalBranchComponent = latestIndex
				}
			case constant.ModifiedFilesComponent:
				// we don't use the list native Update() because we need to also render the diff view as well as track the current selected index
				if m.CurrentRepoModifiedFilesInfoList.Index() < len(m.CurrentRepoModifiedFilesInfoList.Items())-1 {
					latestIndex := m.CurrentRepoModifiedFilesInfoList.Index() + 1
					m.CurrentRepoModifiedFilesInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.ModifiedFilesComponent = latestIndex
					reinitAndRenderDetailComponentPanelViewPort(m)
				}
			case constant.StashComponent:
				// we don't use the list native Update() because we need to also render the diff view as well as track the current selected index
				if m.CurrentRepoStashInfoList.Index() < len(m.CurrentRepoStashInfoList.Items())-1 {
					latestIndex := m.CurrentRepoStashInfoList.Index() + 1
					m.CurrentRepoStashInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.StashComponent = latestIndex
				}
			case constant.DetailComponent:
				m.DetailPanelViewport, cmd = m.DetailPanelViewport.Update(msg)
				return m, cmd
			}
		} else {
			return upDownKeyMsgUpdateForPopUp(msg, m)
		}
		return m, nil

	case "left", "h":
		if !m.ShowPopUp.Load() {
			m.DetailPanelViewport.MoveLeft(1)
		} else {
			switch m.PopUpType {
			case constant.CommitPopUp:
				popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
				if ok {
					popUp.GitCommitOutputViewport, cmd = popUp.GitCommitOutputViewport.Update(msg)
					return m, cmd
				}
			}
		}

	case "right", "l":
		if !m.ShowPopUp.Load() {
			m.DetailPanelViewport.MoveRight(1)
		} else {
			switch m.PopUpType {
			case constant.CommitPopUp:
				popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
				if ok {
					popUp.GitCommitOutputViewport, cmd = popUp.GitCommitOutputViewport.Update(msg)
					return m, cmd
				}
			}
		}
	}
	return m, nil
}

func GittiMouseInteraction(msg tea.MouseMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "wheelleft":
		if !m.ShowPopUp.Load() {
			m.DetailPanelViewport.MoveLeft(1)
		}

	case "wheelright":
		if !m.ShowPopUp.Load() {
			m.DetailPanelViewport.MoveRight(1)
		}

	case "wheelup":
		if !m.ShowPopUp.Load() {
			m.DetailPanelViewport, cmd = m.DetailPanelViewport.Update(msg)
			return m, cmd
		} else {
			return upDownMouseMsgUpdateForPopUp(msg, m)
		}

	case "wheeldown":
		if !m.ShowPopUp.Load() {
			m.DetailPanelViewport, cmd = m.DetailPanelViewport.Update(msg)
			return m, cmd
		} else {
			return upDownMouseMsgUpdateForPopUp(msg, m)
		}
	}
	return m, nil
}

func upDownKeyMsgUpdateForPopUp(msg tea.KeyMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	// for within pop up component
	switch m.PopUpType {
	// following is for list component
	case constant.GlobalKeyBindingPopUp:
		popUp, ok := m.PopUpModel.(*GlobalKeyBindingPopUpModel)
		if ok {
			popUp.GlobalKeyBindingViewport, cmd = popUp.GlobalKeyBindingViewport.Update(msg)
			return m, cmd
		}
	case constant.ChooseRemotePopUp:
		popUp, ok := m.PopUpModel.(*ChooseRemotePopUpModel)
		if ok {
			popUp.RemoteList, cmd = popUp.RemoteList.Update(msg)
			return m, cmd
		}
	case constant.ChoosePushTypePopUp:
		popUp, ok := m.PopUpModel.(*ChoosePushTypePopUpModel)
		if ok {
			popUp.PushOptionList, cmd = popUp.PushOptionList.Update(msg)
			return m, cmd
		}
	case constant.ChooseNewBranchTypePopUp:
		popUp, ok := m.PopUpModel.(*ChooseNewBranchTypeOptionPopUpModel)
		if ok {
			popUp.NewBranchTypeOptionList, cmd = popUp.NewBranchTypeOptionList.Update(msg)
			return m, cmd
		}
	case constant.ChooseSwitchBranchTypePopUp:
		popUp, ok := m.PopUpModel.(*ChooseSwitchBranchTypePopUpModel)
		if ok {
			popUp.SwitchTypeOptionList, cmd = popUp.SwitchTypeOptionList.Update(msg)
			return m, cmd
		}
	case constant.ChooseGitPullTypePopUp:
		popUp, ok := m.PopUpModel.(*ChooseGitPullTypePopUpModel)
		if ok {
			popUp.PullTypeOptionList, cmd = popUp.PullTypeOptionList.Update(msg)
			return m, cmd
		}
	// following is for viewport
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
		if ok {
			popUp.GitCommitOutputViewport, cmd = popUp.GitCommitOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.GitRemotePushPopUp:
		popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
		if ok {
			popUp.GitRemotePushOutputViewport, cmd = popUp.GitRemotePushOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.GitPullOutputPopUp:
		popUp, ok := m.PopUpModel.(*GitPullOutputPopUpModel)
		if ok {
			popUp.GitPullOutputViewport, cmd = popUp.GitPullOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
		if ok {
			popUp.AddRemoteOutputViewport, cmd = popUp.AddRemoteOutputViewport.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func upDownMouseMsgUpdateForPopUp(msg tea.MouseMsg, m *GittiModel) (*GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	// for pop up that have viewport
	switch m.PopUpType {
	case constant.GlobalKeyBindingPopUp:
		popUp, ok := m.PopUpModel.(*GlobalKeyBindingPopUpModel)
		if ok {
			popUp.GlobalKeyBindingViewport, cmd = popUp.GlobalKeyBindingViewport.Update(msg)
			return m, cmd
		}
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
		if ok {
			popUp.GitCommitOutputViewport, cmd = popUp.GitCommitOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.GitRemotePushPopUp:
		popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
		if ok {
			popUp.GitRemotePushOutputViewport, cmd = popUp.GitRemotePushOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.GitPullOutputPopUp:
		popUp, ok := m.PopUpModel.(*GitPullOutputPopUpModel)
		if ok {
			popUp.GitPullOutputViewport, cmd = popUp.GitPullOutputViewport.Update(msg)
			return m, cmd
		}
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
		if ok {
			popUp.AddRemoteOutputViewport, cmd = popUp.AddRemoteOutputViewport.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}
