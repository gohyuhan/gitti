package handler

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/component/branch"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/component/stash"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/layout"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	discardPopUp "github.com/gohyuhan/gitti/tui/popup/discard"
	keybindingPopUp "github.com/gohyuhan/gitti/tui/popup/keybinding"
	pullPopUp "github.com/gohyuhan/gitti/tui/popup/pull"
	pushPopUp "github.com/gohyuhan/gitti/tui/popup/push"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	resolvePopUp "github.com/gohyuhan/gitti/tui/popup/resolve"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

func handleNonTypingGlobalKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	m.ShowPopUp.Store(true)
	m.IsTyping.Store(false)
	m.PopUpType = constant.GlobalKeyBindingPopUp
	keybindingPopUp.InitGlobalKeyBindingPopUpModel(m)
	return m, nil
}

func handleNonTyping1KeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent != constant.LocalBranchComponent {
			m.CurrentSelectedComponent = constant.LocalBranchComponent
			m.CurrentSelectedComponentIndex = 1
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	}
	return m, nil
}

func handleNonTyping2KeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent != constant.ModifiedFilesComponent {
			m.CurrentSelectedComponent = constant.ModifiedFilesComponent
			m.CurrentSelectedComponentIndex = 2
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	}
	return m, nil
}

func handleNonTyping3KeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent != constant.StashComponent {
			m.CurrentSelectedComponent = constant.StashComponent
			m.CurrentSelectedComponentIndex = 3
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	}
	return m, nil
}

func handleNonTypingaKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		m.ShowPopUp.Store(true)
		m.PopUpType = constant.AmendCommitPopUp
		m.GitOperations.GitCommit.ClearGitCommitOutput()

		commitPopUp.InitGitAmendCommitPopUpModel(m)

		m.IsTyping.Store(true)
	}
	return m, nil
}

func handleNonTypingcKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		m.GitOperations.GitCommit.ClearGitCommitOutput()
		// if the current pop up model is not commit pop up model, then init it
		if popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel); !ok {
			commitPopUp.InitGitCommitPopUpModel(m)
		} else {
			popUp.InitialCommitStarted.Store(false)
			popUp.GitCommitOutputViewport.SetContent("")
		}
		m.PopUpType = constant.CommitPopUp
		m.ShowPopUp.Store(true)
		m.IsTyping.Store(true)
	}
	return m, nil
}

func handleNonTypingdKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchComponent:
			selectedBranchItem := m.CurrentRepoBranchesInfoList.SelectedItem()
			if selectedBranchItem != nil {
				branchItem := selectedBranchItem.(branch.GitBranchItem)
				if branchItem.IsCheckedOut {
					return m, nil
				} else {
					branchPopUp.InitGitDeleteBranchConfirmPromptPopUpModel(m, branchItem.BranchName)
					m.PopUpType = constant.GitDeleteBranchConfirmPromptPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
				}
			}
		case constant.StashComponent:
			selectedStashId := m.CurrentRepoStashInfoList.SelectedItem()
			if selectedStashId != nil {
				stashPopUp.InitGitStashConfirmPromptPopUpModel(m, git.DROPSTASH, "", selectedStashId.(stash.GitStashItem).Id, selectedStashId.(stash.GitStashItem).Message)
				m.PopUpType = constant.GitStashConfirmPromptPopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
			}
		case constant.ModifiedFilesComponent:
			currentSelectedFileItem := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
			if currentSelectedFileItem != nil {
				currentSelectedFile := currentSelectedFileItem.(files.GitModifiedFilesItem)
				// return early if the file has conflict (we should not allow discard on conflict files but resolve option instead)
				if currentSelectedFile.HasConflict {
					return m, nil
				}
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)

				// determine the pop up state
				if (currentSelectedFile.IndexState == "A" && currentSelectedFile.WorkTree != " ") || (currentSelectedFile.IndexState == "C" && currentSelectedFile.WorkTree != " ") {
					// indicating the files is a newly added tracked / copied file with unstage modification (modified or delete )
					m.PopUpType = constant.GitDiscardTypeOptionPopUp
					discardPopUp.InitGitDiscardTypeOptionPopUp(m, currentSelectedFile.FilePathname, true, false)
				} else if currentSelectedFile.IndexState == "R" && currentSelectedFile.WorkTree != " " {
					// a staged rename with unstaged modification
					m.PopUpType = constant.GitDiscardTypeOptionPopUp
					discardPopUp.InitGitDiscardTypeOptionPopUp(m, currentSelectedFile.FilePathname, false, true)
				} else if currentSelectedFile.IndexState == "?" && currentSelectedFile.WorkTree == "?" {
					// newly added untracked file
					m.PopUpType = constant.GitDiscardConfirmPromptPopUp
					discardPopUp.InitGitDiscardConfirmPromptPopupModel(m, currentSelectedFile.FilePathname, git.DISCARDUNTRACKED)
				} else if currentSelectedFile.IndexState != "A" && currentSelectedFile.IndexState != "C" && currentSelectedFile.IndexState != "R" && currentSelectedFile.IndexState != "?" && currentSelectedFile.IndexState != " " && currentSelectedFile.WorkTree != " " {
					// tracked file with both staged and unstaged modification (beside A, C and  )
					m.PopUpType = constant.GitDiscardTypeOptionPopUp
					discardPopUp.InitGitDiscardTypeOptionPopUp(m, currentSelectedFile.FilePathname, false, false)
				} else if (currentSelectedFile.IndexState == "A" && currentSelectedFile.WorkTree == " ") || (currentSelectedFile.IndexState == "C" && currentSelectedFile.WorkTree == " ") {
					// newly added tracked / copied file
					m.PopUpType = constant.GitDiscardConfirmPromptPopUp
					discardPopUp.InitGitDiscardConfirmPromptPopupModel(m, currentSelectedFile.FilePathname, git.DISCARDNEWLYADDEDORCOPIED)
				} else if currentSelectedFile.IndexState == "R" && currentSelectedFile.WorkTree == " " {
					// a staged rename
					m.PopUpType = constant.GitDiscardConfirmPromptPopUp
					discardPopUp.InitGitDiscardConfirmPromptPopupModel(m, currentSelectedFile.FilePathname, git.DISCARDANDREVERTRENAME)
				} else {
					// tracked file with only unstaged modification
					m.PopUpType = constant.GitDiscardConfirmPromptPopUp
					discardPopUp.InitGitDiscardConfirmPromptPopupModel(m, currentSelectedFile.FilePathname, git.DISCARDWHOLE)
				}
			}
		}
	}
	return m, nil
}

func handleNonTypingnKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent == constant.LocalBranchComponent {
			m.PopUpType = constant.ChooseNewBranchTypePopUp
			m.IsTyping.Store(false)
			m.ShowPopUp.Store(true)
			if _, ok := m.PopUpModel.(*branchPopUp.ChooseNewBranchTypeOptionPopUpModel); !ok {
				branchPopUp.InitChooseNewBranchTypePopUpModel(m)
			}
		}
	}
	return m, nil
}

func handleNonTypingpKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		// first we need to check if there are any push/pull origin origin for this repo
		// if not we prompt the user to add a new remote origin
		if !m.GitOperations.GitRemote.CheckRemoteExist() {
			m.PopUpType = constant.AddRemotePromptPopUp
			// if the current pop up model is not commit pop up model, then init it
			if popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel); !ok {
				remotePopUp.InitAddRemotePromptPopUpModel(m, true)
			} else {
				popUp.AddRemoteOutputViewport.SetContent("")
			}
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(true)
		} else {
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
			remotes := m.GitOperations.GitRemote.Remote()
			if len(remotes) == 1 {
				m.PopUpType = constant.ChoosePushTypePopUp
				// if the current pop up model is not commit pop up model, then init it and start git push service
				pushPopUp.InitChoosePushTypePopUpModel(m, remotes[0].Name)
			} else if len(remotes) > 1 {
				// if remote is more than 1 let user choose which remote to push to first before pushing
				m.PopUpType = constant.ChooseRemotePopUp
				if _, ok := m.PopUpModel.(*remotePopUp.ChooseRemotePopUpModel); !ok {
					remotePopUp.InitGitRemotePushChooseRemotePopUpModel(m, remotes)
				}
			}
		}
	}
	return m, nil
}

func handleNonTypingPKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		// first we need to check if there are any push/pull origin for this repo
		// if not we prompt the user to add a new remote origin
		if !m.GitOperations.GitRemote.CheckRemoteExist() {
			m.PopUpType = constant.AddRemotePromptPopUp
			// if the current pop up model is not commit pop up model, then init it
			if popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel); !ok {
				remotePopUp.InitAddRemotePromptPopUpModel(m, true)
			} else {
				popUp.AddRemoteOutputViewport.SetContent("")
			}
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(true)
		} else {
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
			m.PopUpType = constant.ChooseGitPullTypePopUp
			pullPopUp.InitChooseGitPullTypePopUp(m)
		}
	}
	return m, nil
}

func handleNonTypingrKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.ModifiedFilesComponent:
			currentSelectedFileItem := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
			if currentSelectedFileItem != nil {
				currentSelectedFile := currentSelectedFileItem.(files.GitModifiedFilesItem)
				// return early if the file has no conflict
				if !currentSelectedFile.HasConflict {
					return m, nil
				}
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				m.PopUpType = constant.GitResolveConflictOptionPopUp
				resolvePopUp.InitGitResolveConflictOptionPopUpModel(m, currentSelectedFile.FilePathname)
			}
		}
	}
	return m, nil
}

func handleNonTypingsKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.CurrentSelectedComponent == constant.ModifiedFilesComponent {
		currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
		var filePathName string
		if currentSelectedModifiedFile != nil {
			selectedFile := currentSelectedModifiedFile.(files.GitModifiedFilesItem)
			// return early if the file is in a conflict status
			if selectedFile.HasConflict {
				return m, nil
			}
			filePathName = selectedFile.FilePathname
			m.PopUpType = constant.GitStashMessagePopUp
			stashPopUp.InitGitStashMessagePopUpModel(m, filePathName, git.STASHFILE)
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(true)
		}
	}
	return m, nil
}

func handleNonTypingSKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.CurrentSelectedComponent == constant.ModifiedFilesComponent {
		currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
		var filePathName string
		if currentSelectedModifiedFile != nil {
			filePathName = currentSelectedModifiedFile.(files.GitModifiedFilesItem).FilePathname
			m.PopUpType = constant.GitStashMessagePopUp
			stashPopUp.InitGitStashMessagePopUpModel(m, filePathName, git.STASHALL)
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(true)
		}
	}
	return m, nil
}

func handleNonTypingqQKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if api.GITDAEMON != nil {
			api.GITDAEMON.Stop()
		}
		return m, tea.Quit
	}
	return m, nil
}

func handleNonTypingBackspaceKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() && m.CurrentSelectedComponent == constant.StashComponent {
		selectedStashId := m.CurrentRepoStashInfoList.SelectedItem()
		if selectedStashId != nil {
			stashPopUp.InitGitStashConfirmPromptPopUpModel(m, git.POPSTASH, "", selectedStashId.(stash.GitStashItem).Id, selectedStashId.(stash.GitStashItem).Message)
			m.PopUpType = constant.GitStashConfirmPromptPopUp
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
		}
	}
	return m, nil
}

func handleNonTypingEnterKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.ModifiedFilesComponent:
			if len(m.CurrentRepoModifiedFilesInfoList.Items()) > 0 {
				m.CurrentSelectedComponent = constant.DetailComponent
				m.DetailPanelParentComponent = constant.ModifiedFilesComponent
			}
		case constant.StashComponent:
			if len(m.CurrentRepoStashInfoList.Items()) > 0 {
				m.CurrentSelectedComponent = constant.DetailComponent
				m.DetailPanelParentComponent = constant.StashComponent
			}
		case constant.LocalBranchComponent:
			currentSelectedLocalBranch := m.CurrentRepoBranchesInfoList.SelectedItem().(branch.GitBranchItem)
			// only proceed if the local branch selected is not current checkedout branch
			// we can't switch from current checkout branch to current checkout branch, do we
			if !currentSelectedLocalBranch.IsCheckedOut {
				m.PopUpType = constant.ChooseSwitchBranchTypePopUp
				m.IsTyping.Store(false)
				m.ShowPopUp.Store(true)
				branchPopUp.InitChooseSwitchBranchTypePopUpModel(m, currentSelectedLocalBranch.BranchName)
			}
		}
	} else {
		switch m.PopUpType {
		case constant.ChooseRemotePopUp:
			popUp, ok := m.PopUpModel.(*remotePopUp.ChooseRemotePopUpModel)
			if ok {
				remote := popUp.RemoteList.SelectedItem()
				m.PopUpType = constant.ChoosePushTypePopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				pushPopUp.InitChoosePushTypePopUpModel(m, remote.(remotePopUp.GitRemoteItem).Name)
			}

		case constant.ChoosePushTypePopUp:
			popUp, ok := m.PopUpModel.(*pushPopUp.ChoosePushTypePopUpModel)
			if ok {
				m.PopUpType = constant.GitRemotePushPopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				selectedOption := popUp.PushOptionList.SelectedItem()
				return services.InitGitRemotePushPopUpModelAndStartGitRemotePushService(m, popUp.RemoteName, selectedOption.(pushPopUp.GitPushOptionItem).PushType)
			}

		case constant.ChooseNewBranchTypePopUp:
			popUp, ok := m.PopUpModel.(*branchPopUp.ChooseNewBranchTypeOptionPopUpModel)
			if ok {
				m.PopUpType = constant.CreateNewBranchPopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(true)
				selectedOption := popUp.NewBranchTypeOptionList.SelectedItem()
				branchPopUp.InitCreateNewBranchPopUpModel(m, selectedOption.(branchPopUp.GitNewBranchTypeOptionItem).NewBranchType)
			}

		case constant.ChooseSwitchBranchTypePopUp:
			popUp, ok := m.PopUpModel.(*branchPopUp.ChooseSwitchBranchTypePopUpModel)
			if ok {
				m.PopUpType = constant.SwitchBranchOutputPopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				selectedOption := popUp.SwitchTypeOptionList.SelectedItem().(branchPopUp.GitSwitchBranchTypeOptionItem)
				branchName := popUp.BranchName
				branchPopUp.InitSwitchBranchOutputPopUpModel(m, branchName, selectedOption.SwitchBranchType)
				popUp, ok := m.PopUpModel.(*branchPopUp.SwitchBranchOutputPopUpModel)
				if ok {
					popUp.IsProcessing.Store(true) // set it directly first
					services.GitSwitchBranchService(m, branchName, selectedOption.SwitchBranchType)
					return m, popUp.Spinner.Tick
				}
			}

		case constant.ChooseGitPullTypePopUp:
			popUp, ok := m.PopUpModel.(*pullPopUp.ChooseGitPullTypePopUpModel)
			if ok {
				m.PopUpType = constant.GitPullOutputPopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				selectedOption := popUp.PullTypeOptionList.SelectedItem().(pullPopUp.GitPullTypeOptionItem)
				pullPopUp.InitGitPullOutputPopUpModel(m)
				popUp, ok := m.PopUpModel.(*pullPopUp.GitPullOutputPopUpModel)
				if ok {
					popUp.IsProcessing.Store(true) // set it directly first
					// start the git pull service
					services.GitPullService(m, selectedOption.PullType)
					return m, popUp.Spinner.Tick
				}
			}
		case constant.GitDiscardTypeOptionPopUp:
			popUp, ok := m.PopUpModel.(*discardPopUp.GitDiscardTypeOptionPopUpModel)
			if ok {
				m.PopUpType = constant.GitDiscardConfirmPromptPopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				selectedOption := popUp.DiscardTypeOptionList.SelectedItem().(discardPopUp.GitDiscardTypeOptionItem)
				discardPopUp.InitGitDiscardConfirmPromptPopupModel(m, popUp.FilePathName, selectedOption.DiscardType)
			}
		case constant.GitDiscardConfirmPromptPopUp:
			popUp, ok := m.PopUpModel.(*discardPopUp.GitDiscardConfirmPromptPopUpModel)
			if ok {
				services.GitDiscardFileChangesService(m, popUp.FilePathName, popUp.DiscardType)
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		case constant.GitStashConfirmPromptPopUp:
			popUp, ok := m.PopUpModel.(*stashPopUp.GitStashConfirmPromptPopUpModel)
			if ok {
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				m.PopUpType = constant.GitStashOperationOutputPopUp
				stashPopUp.InitGitStashOperationOutputPopUpModel(m, popUp.StashOperationType)
				outputPopUp, ok := m.PopUpModel.(*stashPopUp.GitStashOperationOutputPopUpModel)
				if ok {
					services.GitStashOperationService(m, popUp.FilePathName, popUp.StashId, popUp.StashMessage)
					return m, outputPopUp.Spinner.Tick
				}
			}
		case constant.GitResolveConflictOptionPopUp:
			popUp, ok := m.PopUpModel.(*resolvePopUp.GitResolveConflictOptionPopUpModel)
			if ok {
				selectedResolveType := popUp.ResolveConflictOptionList.SelectedItem().(resolvePopUp.GitResolveConflictOptionItem)
				services.GitResolveConflictService(m, popUp.FilePathName, selectedResolveType.ResolveType)
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
			}
		case constant.GitDeleteBranchConfirmPromptPopUp:
			popUp, ok := m.PopUpModel.(*branchPopUp.GitDeleteBranchConfirmPromptPopUpModel)
			branchName := popUp.BranchName
			if ok {
				branchPopUp.InitGitDeleteBranchOutputPopUpModel(m)
				popUp, ok := m.PopUpModel.(*branchPopUp.GitDeleteBranchOutputPopUpModel)
				if ok {
					popUp.IsProcessing.Store(true)
					m.PopUpType = constant.GitDeleteBranchOutputPopUp
					services.GitDeleteBranchService(m, branchName)
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					return m, popUp.Spinner.Tick
				}
			}
		}
	}
	return m, nil
}

func handleNonTypingTabKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	nextNavigation := m.CurrentSelectedComponentIndex + 1
	if nextNavigation < len(constant.ComponentNavigationList) {
		m.CurrentSelectedComponentIndex = nextNavigation
		m.CurrentSelectedComponent = constant.ComponentNavigationList[nextNavigation]
		layout.LeftPanelDynamicResize(m)
		services.FetchDetailComponentPanelInfoService(m, true)
	}
	return m, nil
}

func handleNonTypingShiftTabKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	previousNavigation := m.CurrentSelectedComponentIndex - 1
	if previousNavigation >= 0 {
		m.CurrentSelectedComponentIndex = previousNavigation
		m.CurrentSelectedComponent = constant.ComponentNavigationList[previousNavigation]
		layout.LeftPanelDynamicResize(m)
		services.FetchDetailComponentPanelInfoService(m, true)
	}
	return m, nil
}

func handleNonTypingSpaceKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.ModifiedFilesComponent:
			currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
			var filePathName string
			if currentSelectedModifiedFile != nil {
				filePathName = currentSelectedModifiedFile.(files.GitModifiedFilesItem).FilePathname
				services.GitStageOrUnstageService(m, filePathName)
			}

		case constant.StashComponent:
			selectedStashId := m.CurrentRepoStashInfoList.SelectedItem()
			if selectedStashId != nil {
				stashPopUp.InitGitStashConfirmPromptPopUpModel(m, git.APPLYSTASH, "", selectedStashId.(stash.GitStashItem).Id, selectedStashId.(stash.GitStashItem).Message)
				m.PopUpType = constant.GitStashConfirmPromptPopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
			}
		}
	}
	return m, nil
}

func handleNonTypingEscKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.GlobalKeyBindingPopUp:
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpType = constant.NoPopUp
			m.PopUpModel = nil
		case constant.GitRemotePushPopUp:
			services.GitRemotePushCancelService(m)
		case constant.GitPullOutputPopUp:
			services.GitPullCancelService(m)
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
			// Block ESC during branch switching - operation must complete
			popUp, ok := m.PopUpModel.(*branchPopUp.SwitchBranchOutputPopUpModel)
			if ok && !popUp.IsProcessing.Load() {
				// only close when done processing
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		case constant.ChooseGitPullTypePopUp:
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpType = constant.NoPopUp
			m.PopUpModel = nil
		case constant.GitDiscardTypeOptionPopUp:
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpType = constant.NoPopUp
			m.PopUpModel = nil
		case constant.GitDiscardConfirmPromptPopUp:
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpType = constant.NoPopUp
			m.PopUpModel = nil
		case constant.GitStashOperationOutputPopUp:
			// Block ESC during stash operation - operation must complete
			popUp, ok := m.PopUpModel.(*stashPopUp.GitStashOperationOutputPopUpModel)
			if ok && !popUp.IsProcessing.Load() {
				// only close when done processing
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		case constant.GitStashConfirmPromptPopUp:
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpType = constant.NoPopUp
			m.PopUpModel = nil

		case constant.GitResolveConflictOptionPopUp:
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpType = constant.NoPopUp
			m.PopUpModel = nil

		case constant.GitDeleteBranchConfirmPromptPopUp:
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpType = constant.NoPopUp
			m.PopUpModel = nil

		case constant.GitDeleteBranchOutputPopUp:
			popUp, ok := m.PopUpModel.(*branchPopUp.GitDeleteBranchOutputPopUpModel)
			if ok && !popUp.IsProcessing.Load() {
				// only close when done processing
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		}

		return m, nil
	} else {
		switch m.CurrentSelectedComponent {
		case constant.DetailComponent:
			m.CurrentSelectedComponent = m.DetailPanelParentComponent
			m.DetailPanelParentComponent = ""
		case constant.DetailComponentTwo:
			m.CurrentSelectedComponent = m.DetailPanelParentComponent
			m.DetailPanelParentComponent = ""
		}
	}
	return m, nil
}

func handleNonTypingUpkKeyBindingInteraction(msg tea.KeyMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchComponent:
			// we don't use the list native Update() because we track the current selected index
			if m.CurrentRepoBranchesInfoList.Index() > 0 {
				latestIndex := m.CurrentRepoBranchesInfoList.Index() - 1
				m.CurrentRepoBranchesInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.LocalBranchComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.ModifiedFilesComponent:
			// we don't use the list native Update() because we need to also track the current selected index
			if m.CurrentRepoModifiedFilesInfoList.Index() > 0 {
				latestIndex := m.CurrentRepoModifiedFilesInfoList.Index() - 1
				m.CurrentRepoModifiedFilesInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.ModifiedFilesComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.StashComponent:
			// we don't use the list native Update() because we need to also track the current selected index
			if m.CurrentRepoStashInfoList.Index() > 0 {
				latestIndex := m.CurrentRepoStashInfoList.Index() - 1
				m.CurrentRepoStashInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.StashComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.DetailComponent:
			m.DetailPanelViewport, cmd = m.DetailPanelViewport.Update(msg)
			return m, cmd
		case constant.DetailComponentTwo:
			m.DetailPanelTwoViewport, cmd = m.DetailPanelTwoViewport.Update(msg)
			return m, cmd
		}
	} else {
		return UpDownKeyMsgUpdateForPopUp(msg, m)
	}
	return m, nil
}

func handleNonTypingDownjKeyBindingInteraction(msg tea.KeyMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchComponent:
			// we don't use the list native Update() because we track the current selected index
			if m.CurrentRepoBranchesInfoList.Index() < len(m.CurrentRepoBranchesInfoList.Items())-1 {
				latestIndex := m.CurrentRepoBranchesInfoList.Index() + 1
				m.CurrentRepoBranchesInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.LocalBranchComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.ModifiedFilesComponent:
			// we don't use the list native Update() because we need to also track the current selected index
			if m.CurrentRepoModifiedFilesInfoList.Index() < len(m.CurrentRepoModifiedFilesInfoList.Items())-1 {
				latestIndex := m.CurrentRepoModifiedFilesInfoList.Index() + 1
				m.CurrentRepoModifiedFilesInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.ModifiedFilesComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.StashComponent:
			// we don't use the list native Update() because we need to also track the current selected index
			if m.CurrentRepoStashInfoList.Index() < len(m.CurrentRepoStashInfoList.Items())-1 {
				latestIndex := m.CurrentRepoStashInfoList.Index() + 1
				m.CurrentRepoStashInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.StashComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.DetailComponent:
			m.DetailPanelViewport, cmd = m.DetailPanelViewport.Update(msg)
			return m, cmd
		case constant.DetailComponentTwo:
			m.DetailPanelTwoViewport, cmd = m.DetailPanelTwoViewport.Update(msg)
			return m, cmd
		}
	} else {
		return UpDownKeyMsgUpdateForPopUp(msg, m)
	}
	return m, nil
}

func handleNonTypingLefthKeyBindingInteraction(msg tea.KeyMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.DetailComponent:
			m.DetailPanelViewport.ScrollLeft(1)
		case constant.DetailComponentTwo:
			m.DetailPanelTwoViewport.ScrollLeft(1)
		default:
			m.DetailPanelViewport.ScrollLeft(1)
		}
	} else {
		switch m.PopUpType {
		case constant.CommitPopUp:
			popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
			if ok {
				popUp.GitCommitOutputViewport, cmd = popUp.GitCommitOutputViewport.Update(msg)
				return m, cmd
			}
		}
	}
	return m, nil
}

func handleNonTypingRightlKeyBindingInteraction(msg tea.KeyMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.DetailComponent:
			m.DetailPanelViewport.ScrollRight(1)
		case constant.DetailComponentTwo:
			m.DetailPanelTwoViewport.ScrollRight(1)
		default:
			m.DetailPanelViewport.ScrollRight(1)
		}
	} else {
		switch m.PopUpType {
		case constant.CommitPopUp:
			popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
			if ok {
				popUp.GitCommitOutputViewport, cmd = popUp.GitCommitOutputViewport.Update(msg)
				return m, cmd
			}
		}
	}
	return m, nil
}

// handleNonTypingLeftBracketKeyBindingInteraction handles the '[' key not only for navigation but contextually to switch to the previous detail component panel
func handleNonTypingLeftBracketKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		// handle detail component panel switching
		if m.CurrentSelectedComponent == constant.DetailComponentTwo {
			m.CurrentSelectedComponent = constant.DetailComponent
		}
	}
	return m, nil
}

// handleNonTypingRightBracketKeyBindingInteraction handles the ']' key not only for navigation but contextually to switch to the next detail component panel
func handleNonTypingRightBracketKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		// handle detail component panel switching
		if m.CurrentSelectedComponent == constant.DetailComponent && m.ShowDetailPanelTwo.Load() {
			m.CurrentSelectedComponent = constant.DetailComponentTwo
		}
	}
	return m, nil
}
