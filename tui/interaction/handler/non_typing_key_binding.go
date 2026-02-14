package handler

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/component/branch"
	"github.com/gohyuhan/gitti/tui/component/commitlog"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/component/stash"
	"github.com/gohyuhan/gitti/tui/component/tag"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/layout"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	discardPopUp "github.com/gohyuhan/gitti/tui/popup/discard"
	filesPopUp "github.com/gohyuhan/gitti/tui/popup/files"
	keybindingPopUp "github.com/gohyuhan/gitti/tui/popup/keybinding"
	logPopUp "github.com/gohyuhan/gitti/tui/popup/log"
	pullPopUp "github.com/gohyuhan/gitti/tui/popup/pull"
	pushPopUp "github.com/gohyuhan/gitti/tui/popup/push"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	resolvePopUp "github.com/gohyuhan/gitti/tui/popup/resolve"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

func handleNonTypingGlobalKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	m.ShowPopUp.Store(true)
	m.IsTyping.Store(false)
	m.PopUpType = constant.KeybindingAndFeatureInstructionsPopUp
	keybindingPopUp.InitKeybindingAndFeatureInstructionsPopUpModel(m)
	return m, nil
}

func handleNonTyping1KeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent != constant.LocalBranchOrTagComponentPanel {
			m.CurrentSelectedComponent = constant.LocalBranchOrTagComponentPanel
			m.CurrentSelectedComponentIndex = 1
			m.DetailPanelParentComponent = ""
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	}
	return m, nil
}

func handleNonTyping2KeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent != constant.ModifiedFilesComponentPanel {
			m.CurrentSelectedComponent = constant.ModifiedFilesComponentPanel
			m.CurrentSelectedComponentIndex = 2
			m.DetailPanelParentComponent = ""
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	}
	return m, nil
}

func handleNonTyping3KeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent != constant.CommitLogComponentPanel {
			m.CurrentSelectedComponent = constant.CommitLogComponentPanel
			m.CurrentSelectedComponentIndex = 3
			m.DetailPanelParentComponent = ""
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	}
	return m, nil
}

func handleNonTyping4KeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent != constant.StashComponentPanel {
			m.CurrentSelectedComponent = constant.StashComponentPanel
			m.CurrentSelectedComponentIndex = 4
			m.DetailPanelParentComponent = ""
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	}
	return m, nil
}

func handleNonTypingaKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.GitCherryPickPopUp:
			m.ShowPopUp.Store(true)
			m.PopUpType = constant.GitCherryPickApplyConfirmPopUp
			m.PopUpModel = nil // we don't need to initialize the pop up model, as we are just showing the pop up and we don't need to hold any state or info
			m.IsTyping.Store(false)
		case constant.GitEditCherryPickPopUp:
			m.ShowPopUp.Store(true)
			m.PopUpType = constant.GitCherryPickApplyConfirmPopUp
			m.PopUpModel = nil // we don't need to initialize the pop up model, as we are just showing the pop up and we don't need to hold any state or info
			m.IsTyping.Store(false)
		}
	}
	return m, nil
}

func handleNonTypingAKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
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

func handleNonTypingCKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		services.GitStateUniversalUtilsContinueService(m)
	}
	return m, nil
}

func handleNonTypingdKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagComponentPanel:
			switch m.CurrentLocalBranchOrTagComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
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
			case constant.SHOW_TAG:
				selectedTagItem := m.CurrentRepoTagInfoList.SelectedItem()
				if selectedTagItem != nil {
					tagItem := selectedTagItem.(tag.GitTagItem)
					tagPopUp.InitChooseDeleteTagOptionPopUpModel(m, tagItem.TagName)
					m.PopUpType = constant.ChooseDeleteTagOptionPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
				}
			}
		case constant.StashComponentPanel:
			selectedStashId := m.CurrentRepoStashInfoList.SelectedItem()
			if selectedStashId != nil {
				stashPopUp.InitGitStashConfirmPromptPopUpModel(m, git.DROPSTASH, "", selectedStashId.(stash.GitStashItem).Id, selectedStashId.(stash.GitStashItem).Message)
				m.PopUpType = constant.GitStashConfirmPromptPopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
			}
		case constant.ModifiedFilesComponentPanel:
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
		case constant.DetailComponentPanel, constant.DetailComponentPanelTwo:
			if m.IsLineEditingState.Load() {
				currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
				var filePathName string
				if currentSelectedModifiedFile != nil {
					promptForDiscardLineChangeConfirm := false
					switch m.CurrentSelectedComponent {
					case constant.DetailComponentPanel:
						switch m.LineEditingIndexPositionAndInfo.DetailPanelViewportStageType {
						case constant.STAGE:
							promptForDiscardLineChangeConfirm = false
						case constant.UNSTAGE:
							promptForDiscardLineChangeConfirm = true
						}
					case constant.DetailComponentPanelTwo:
						switch m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportStageType {
						case constant.STAGE:
							promptForDiscardLineChangeConfirm = false
						case constant.UNSTAGE:
							promptForDiscardLineChangeConfirm = true
						}
					}

					if promptForDiscardLineChangeConfirm {
						m.ShowPopUp.Store(true)
						m.IsTyping.Store(false)
						m.PopUpType = constant.GitDiscardFileLineChangeConfirmPopUp
						filesPopUp.InitGitDiscardFileLineChangeConfirmPopUp(m)
					} else {
						filePathName = currentSelectedModifiedFile.(files.GitModifiedFilesItem).FilePathname
						services.GitDiscardLineFileChangeService(m, filePathName)
					}
				}
			}
		}
	} else {
		switch m.PopUpType {
		case constant.GitEditCherryPickPopUp:
			popUp, ok := m.PopUpModel.(*logPopUp.GitEditCherryPickPopUpModel)
			if ok {
				selectedCherryPickedCommiyLogItem := popUp.CherryPickedCommitLog.SelectedItem()
				if selectedCherryPickedCommiyLogItem != nil {
					selectedCherryPickedCommiyLogHash := selectedCherryPickedCommiyLogItem.(logPopUp.GitEditCherryPickItem).Hash
					_, exist := m.CherryPickedCommitInfo.CherryPickedCommitMap[selectedCherryPickedCommiyLogHash]
					if exist {
						delete(m.CherryPickedCommitInfo.CherryPickedCommitMap, selectedCherryPickedCommiyLogHash)
					}
					if len(m.CherryPickedCommitInfo.CherryPickedCommitMap) < 1 {
						utils.ReinitCherryPickedCommitInfo(m)
						m.ShowPopUp.Store(false)
						m.IsTyping.Store(false)
						m.PopUpType = constant.NoPopUp
						m.PopUpModel = nil
					} else {
						// reinit the list after a removal
						logPopUp.InitGitEditCherryPickPopUp(m, popUp.CherryPickedCommitLog.Index())
					}
				}
			}
		}
	}
	return m, nil
}

func handleNonTypingeKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.GitCherryPickPopUp:
			m.PopUpType = constant.GitEditCherryPickPopUp
			logPopUp.InitGitEditCherryPickPopUp(m, 0)
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
		case constant.GitCherryPickApplyConfirmPopUp:
			m.PopUpType = constant.GitEditCherryPickPopUp
			logPopUp.InitGitEditCherryPickPopUp(m, 0)
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
		}
	} else {
		switch m.CurrentSelectedComponent {
		case constant.ModifiedFilesComponentPanel:
			currentSelectedFileItem := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
			if currentSelectedFileItem != nil {
				currentSelectedFile := currentSelectedFileItem.(files.GitModifiedFilesItem)
				cmd, isNonTerminalEditor := utils.ReturnEditorLaunchCommand(currentSelectedFile.FilePathname, m.UserSetEditor)
				if isNonTerminalEditor {
					cmd.Start()
					return m, nil
				} else {
					return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
						return types.EditorFinishedMsg{
							Err: err,
						}
					})
				}
			}
		case constant.LogComponentPanel:
			go func() {
				m.GittiLogger.ExportLogging()
			}()
		}
	}
	return m, nil
}

func handleNonTypingfKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if (m.CurrentSelectedComponent == constant.LocalBranchOrTagComponentPanel || m.DetailPanelParentComponent == constant.LocalBranchOrTagComponentPanel) &&
			m.CurrentLocalBranchOrTagComponentShowing == constant.SHOW_TAG {
			if !m.GitOperations.GitRemote.CheckRemoteExist() {
				// if no remote found, we add one
				m.PopUpType = constant.AddRemotePromptPopUp
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
				remotes := m.GitOperations.GitRemote.FetchRemote()
				if len(remotes) == 1 {
					m.PopUpType = constant.ChooseFetchTagOptionPopUp
					// only one remote found so, we will default to that remote
					tagPopUp.InitChooseFetchTagOptionPopUpModel(m, remotes[0].Name)
				} else if len(remotes) > 1 {
					// if remote is more than 1 let user choose which remote
					m.PopUpType = constant.ChooseRemotePopUp
					remotePopUp.InitChooseRemotePopUpModel(m, remotes, constant.TAGFETCHACTION)
				}
			}
		}
	}
	return m, nil
}

func handleNonTypingLKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	// enter line stage state
	services.EnterOrReinitLineEditingStateService(m)
	m.TuiUpdateChannel <- constant.DETAIL_COMPONENT_PANEL_UPDATED
	return m, nil
}

func handleNonTypingnKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagComponentPanel:
			switch m.CurrentLocalBranchOrTagComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				m.PopUpType = constant.ChooseNewBranchTypePopUp
				m.IsTyping.Store(false)
				m.ShowPopUp.Store(true)
				if _, ok := m.PopUpModel.(*branchPopUp.ChooseNewBranchTypeOptionPopUpModel); !ok {
					branchPopUp.InitChooseNewBranchTypePopUpModel(m)
				}
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
			remotes := m.GitOperations.GitRemote.PushRemote()
			if len(remotes) == 1 {
				m.PopUpType = constant.ChoosePushTypePopUp
				// if the current pop up model is not commit pop up model, then init it and start git push service
				pushPopUp.InitChoosePushTypePopUpModel(m, remotes[0].Name)
			} else if len(remotes) > 1 {
				// if remote is more than 1 let user choose which remote to push to first before pushing
				m.PopUpType = constant.ChooseRemotePopUp
				if _, ok := m.PopUpModel.(*remotePopUp.ChooseRemotePopUpModel); !ok {
					remotePopUp.InitChooseRemotePopUpModel(m, remotes, constant.PUSHACTION)
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
		case constant.ModifiedFilesComponentPanel:
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
		case constant.CommitLogComponentPanel:
			selectedCommit := m.CurrentRepoCommitLogInfoList.SelectedItem()
			if m.CurrentRepoCommitLogInfoList.Items() != nil {
				parsedCommit := selectedCommit.(commitlog.GitCommitLogItem)
				commitPopUp.InitGitResetToSelectedCommitTypeOptionPopUpModel(
					m,
					parsedCommit.Hash,
					parsedCommit.Message,
					parsedCommit.Author,
				)
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				m.PopUpType = constant.GitResetToSelectedCommitTypeOptionPopUp
			}
		}
	}
	return m, nil
}

func handleNonTypingRKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.CommitLogComponentPanel:
			if len(m.CurrentRepoCommitLogInfoList.Items()) > 1 {
				commitPopUp.InitGitResetLatestCommitTypeOptionPopUpModel(m)
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				m.PopUpType = constant.GitResetLatestCommitTypeOptionPopUp
			}
		}
	}
	return m, nil
}

func handleNonTypingsKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.CurrentSelectedComponent == constant.ModifiedFilesComponentPanel {
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
	if m.CurrentSelectedComponent == constant.ModifiedFilesComponentPanel {
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

func handleNonTypingtKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch m.CurrentSelectedComponent {
	case constant.CommitLogComponentPanel:
		currentSelectedCommit := m.CurrentRepoCommitLogInfoList.SelectedItem()
		if currentSelectedCommit != nil {
			commit := currentSelectedCommit.(commitlog.GitCommitLogItem)
			m.PopUpType = constant.CreateTagPopUp
			tagPopUp.InitCreateTagPopUpModel(m, commit.Hash, commit.Message)
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
	if !m.ShowPopUp.Load() && m.CurrentSelectedComponent == constant.StashComponentPanel {
		selectedStashId := m.CurrentRepoStashInfoList.SelectedItem()
		if selectedStashId != nil {
			stashPopUp.InitGitStashConfirmPromptPopUpModel(m, git.POPSTASH, "", selectedStashId.(stash.GitStashItem).Id, selectedStashId.(stash.GitStashItem).Message)
			m.PopUpType = constant.GitStashConfirmPromptPopUp
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
		}
	} else if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.GitEditCherryPickPopUp:
			_, ok := m.PopUpModel.(*logPopUp.GitEditCherryPickPopUpModel)
			if ok {
				utils.ReinitCherryPickedCommitInfo(m)
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		}
	}
	return m, nil
}

func handleNonTypingEnterKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.ModifiedFilesComponentPanel:
			if len(m.CurrentRepoModifiedFilesInfoList.Items()) > 0 {
				m.CurrentSelectedComponent = constant.DetailComponentPanel
				m.DetailPanelParentComponent = constant.ModifiedFilesComponentPanel
			}
		case constant.CommitLogComponentPanel:
			if len(m.CurrentRepoCommitLogInfoList.Items()) > 0 {
				m.CurrentSelectedComponent = constant.DetailComponentPanel
				m.DetailPanelParentComponent = constant.CommitLogComponentPanel
			}
		case constant.StashComponentPanel:
			if len(m.CurrentRepoStashInfoList.Items()) > 0 {
				m.CurrentSelectedComponent = constant.DetailComponentPanel
				m.DetailPanelParentComponent = constant.StashComponentPanel
			}
		case constant.LocalBranchOrTagComponentPanel:
			switch m.CurrentLocalBranchOrTagComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
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
		case constant.LogComponentPanel:
			m.CurrentSelectedComponent = constant.DetailComponentPanel
			m.DetailPanelParentComponent = constant.LogComponentPanel
		}
	} else {
		switch m.PopUpType {
		case constant.ChooseRemotePopUp:
			popUp, ok := m.PopUpModel.(*remotePopUp.ChooseRemotePopUpModel)
			if ok {
				remote := popUp.RemoteList.SelectedItem()
				remoteName := remote.(remotePopUp.GitRemoteItem).Name
				switch popUp.Action {
				case constant.PUSHACTION:
					m.PopUpType = constant.ChoosePushTypePopUp
					m.IsTyping.Store(false)
					m.ShowPopUp.Store(true)
					pushPopUp.InitChoosePushTypePopUpModel(m, remoteName)
				case constant.CREATEBRANCHBASEDONREMOTE:
					m.IsTyping.Store(true)
					m.ShowPopUp.Store(true)
					m.PopUpType = constant.CreateBranchBasedOnRemotePopUp
					// only one remote found so, we will default to that remote
					branchPopUp.InitCreateBranchBasedOnRemotePopUp(m, remoteName)
				case constant.TAGPUSHACTION:
					selectedTag := m.CurrentRepoTagInfoList.SelectedItem()
					if selectedTag != nil {
						m.IsTyping.Store(false)
						m.ShowPopUp.Store(true)
						m.PopUpType = constant.ChoosePushTagOptionPopUp
						tagPopUp.InitChoosePushTagOptionPopUpModel(m, remoteName, selectedTag.(tag.GitTagItem).TagName)
					} else {
						m.IsTyping.Store(false)
						m.ShowPopUp.Store(false)
						m.PopUpType = constant.NoPopUp
						m.PopUpModel = nil
					}
				case constant.TAGFETCHACTION:
					m.IsTyping.Store(false)
					m.ShowPopUp.Store(true)
					m.PopUpType = constant.ChooseFetchTagOptionPopUp
					tagPopUp.InitChooseFetchTagOptionPopUpModel(m, remoteName)
				}
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
				selectedOption := popUp.NewBranchTypeOptionList.SelectedItem()
				newBranchType := selectedOption.(branchPopUp.GitNewBranchTypeOptionItem).NewBranchType
				if newBranchType == git.NEWBRANCHBASEDONREMOTE {
					if !m.GitOperations.GitRemote.CheckRemoteExist() {
						// if no remote found, we add one
						m.PopUpType = constant.AddRemotePromptPopUp
						if popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel); !ok {
							remotePopUp.InitAddRemotePromptPopUpModel(m, true)
						} else {
							popUp.AddRemoteOutputViewport.SetContent("")
						}
						m.ShowPopUp.Store(true)
						m.IsTyping.Store(true)
					} else {
						m.ShowPopUp.Store(true)
						remotes := m.GitOperations.GitRemote.FetchRemote()
						if len(remotes) == 1 {
							m.IsTyping.Store(true)
							m.PopUpType = constant.CreateBranchBasedOnRemotePopUp
							// only one remote found so, we will default to that remote
							branchPopUp.InitCreateBranchBasedOnRemotePopUp(m, remotes[0].Name)
						} else if len(remotes) > 1 {
							m.IsTyping.Store(false)
							// if remote is more than 1 let user choose which remote
							m.PopUpType = constant.ChooseRemotePopUp
							if _, ok := m.PopUpModel.(*remotePopUp.ChooseRemotePopUpModel); !ok {
								remotePopUp.InitChooseRemotePopUpModel(m, remotes, constant.CREATEBRANCHBASEDONREMOTE)
							}
						}
					}
				} else {
					m.PopUpType = constant.CreateNewBranchPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(true)
					branchPopUp.InitCreateNewBranchPopUpModel(m, newBranchType)
				}
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
		case constant.GitResetLatestCommitTypeOptionPopUp:
			popUp, ok := m.PopUpModel.(*commitPopUp.GitResetLatestCommitTypeOptionPopUpModel)
			if ok {
				selectedResetLatestCommitType := popUp.ResetLatestCommitTypeOptionList.SelectedItem()
				if selectedResetLatestCommitType != nil {
					resetType := selectedResetLatestCommitType.(commitPopUp.GitResetLatestCommitTypeOptionItem).ResetType
					commitPopUp.InitGitResetLatestCommitConfirmPromptPopUpModel(m, resetType)
					_, ok = m.PopUpModel.(*commitPopUp.GitResetLatestCommitConfirmPromptPopUpModel)
					if ok {
						m.PopUpType = constant.GitResetLatestCommitConfirmPromptPopUp
						m.ShowPopUp.Store(true)
						m.IsTyping.Store(false)
					}
				}
			}
		case constant.GitResetLatestCommitConfirmPromptPopUp:
			popUp, ok := m.PopUpModel.(*commitPopUp.GitResetLatestCommitConfirmPromptPopUpModel)
			if ok {
				selectedResetLatestCommitType := popUp.GitResetLatestCommitType
				services.GitResetLatestCommitService(m, selectedResetLatestCommitType)
				m.PopUpType = constant.NoPopUp
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
			}
		case constant.GitResetToSelectedCommitTypeOptionPopUp:
			popUp, ok := m.PopUpModel.(*commitPopUp.GitResetToSelectedCommitTypeOptionPopUpModel)
			if ok {
				selectedResetToSelectedCommitType := popUp.ResetToSelectedCommitTypeOptionList.SelectedItem()
				if selectedResetToSelectedCommitType != nil {
					resetType := selectedResetToSelectedCommitType.(commitPopUp.GitResetToSelectedCommitTypeOptionItem).ResetType
					commitPopUp.InitGitResetToSelectedCommitConfirmPromptPopUpModel(m, resetType, popUp.SelectedCommitHash, popUp.CommitInfoMessage, popUp.CommitInfoAuthor)
					_, ok = m.PopUpModel.(*commitPopUp.GitResetToSelectedCommitConfirmPromptPopUpModel)
					if ok {
						m.PopUpType = constant.GitResetToSelectedCommitConfirmPromptPopUp
						m.ShowPopUp.Store(true)
						m.IsTyping.Store(false)
					}
				}
			}
		case constant.GitResetToSelectedCommitConfirmPromptPopUp:
			popUp, ok := m.PopUpModel.(*commitPopUp.GitResetToSelectedCommitConfirmPromptPopUpModel)
			if ok {
				selectedResetToSelectedCommitType := popUp.GitResetToSelectedCommitType
				services.GitResetToSelectedCommitService(m, selectedResetToSelectedCommitType, popUp.SelectedCommitHash)
				m.PopUpType = constant.NoPopUp
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
			}
		case constant.GitCherryPickOptionSelectionPopUp:
			popUp, ok := m.PopUpModel.(*logPopUp.GitCherryPickOptionSelectionPopUpModel)
			if ok {
				selectedCherryPickType := popUp.CherryPickedOpsOption.SelectedItem()
				if selectedCherryPickType != nil {
					cherryPickType := selectedCherryPickType.(logPopUp.CherryPickOpsOptionItem).CherryPickOpsType
					switch cherryPickType {
					case constant.CHERRYPICK:
						m.PopUpType = constant.GitCherryPickPopUp
						logPopUp.InitGitCherryPickPopUp(m, m.CheckOutBranch)
					case constant.EDITCHERRYPICK:
						m.PopUpType = constant.GitEditCherryPickPopUp
						logPopUp.InitGitEditCherryPickPopUp(m, 0)
					case constant.APPLYCHERRYPICK:
						m.ShowPopUp.Store(true)
						m.PopUpType = constant.GitCherryPickApplyConfirmPopUp
						m.PopUpModel = nil // we don't need to initialize the pop up model, as we are just showing the pop up and we don't need to hold any state or info
						m.IsTyping.Store(false)
					}
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
				}
			}
		case constant.GitCherryPickApplyConfirmPopUp:
			services.GitCherryPickService(m, m.CherryPickedCommitInfo.CherryPickedCommitMap)
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpModel = nil
			m.PopUpType = constant.NoPopUp
		case constant.GitDiscardFileLineChangeConfirmPopUp:
			if m.IsLineEditingState.Load() {
				currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
				var filePathName string
				if currentSelectedModifiedFile != nil {
					filePathName = currentSelectedModifiedFile.(files.GitModifiedFilesItem).FilePathname
					services.GitDiscardLineFileChangeService(m, filePathName)
					m.ShowPopUp.Store(false)
					m.IsTyping.Store(false)
					m.PopUpModel = nil
					m.PopUpType = constant.NoPopUp
				}
			}
		case constant.CreateTagConfirmationPopUp:
			popUp, ok := m.PopUpModel.(*tagPopUp.CreateTagConfirmationPopUpModel)
			if ok {
				services.CreateNewTagService(m, popUp.CommitHash, popUp.TagName, popUp.TagMessage)
				m.PopUpType = constant.NoPopUp
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
			}
		case constant.ChooseDeleteTagOptionPopUp:
			popUp, ok := m.PopUpModel.(*tagPopUp.ChooseDeleteTagOptionPopUpModel)
			if ok {
				selectedDeleteType := popUp.DeleteOptionList.SelectedItem()
				if selectedDeleteType != nil {
					deleteType := selectedDeleteType.(tagPopUp.DeleteTagOptionItem).DeleteTagType
					switch deleteType {
					case git.TAGDELETELOCAL:
						m.ShowPopUp.Store(true)
						m.IsTyping.Store(false)
						m.PopUpType = constant.DeleteTagOutputPopUp
						tagPopUp.InitDeleteTagOutputPopUpModel(m, popUp.TagName)
						services.DeleteTagService(m, "", popUp.TagName, git.TAGDELETELOCAL)

						// to return the tick for the spinner
						tagOutputPopUp, ok := m.PopUpModel.(*tagPopUp.DeleteTagOutputPopUpModel)
						if ok {
							return m, tagOutputPopUp.Spinner.Tick
						}
					case git.TAGDELETEREMOTE:
						// first we need to check if there are any origin for this repo
						// if not we prompt the user to add a new remote origin
						if !m.GitOperations.GitRemote.CheckRemoteExist() {
							m.PopUpType = constant.AddRemotePromptPopUp
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
							remotes := m.GitOperations.GitRemote.PushRemote()
							if len(remotes) == 1 {
								m.PopUpType = constant.DeleteTagOutputPopUp
								tagPopUp.InitDeleteTagOutputPopUpModel(m, popUp.TagName)
								services.DeleteTagService(m, remotes[0].Name, popUp.TagName, git.TAGDELETEREMOTE)

								// to return the tick for the spinner
								tagOutputPopUp, ok := m.PopUpModel.(*tagPopUp.DeleteTagOutputPopUpModel)
								if ok {
									return m, tagOutputPopUp.Spinner.Tick
								}
							} else if len(remotes) > 1 {
								m.PopUpType = constant.ChooseRemoteForDeleteRemoteTagPopUp
								tagPopUp.InitChooseRemoteForDeleteRemoteTagPopUpModel(m, remotes, popUp.TagName, deleteType)
							}
						}
					}
				}
			}
		case constant.ChooseRemoteForDeleteRemoteTagPopUp:
			popUp, ok := m.PopUpModel.(*tagPopUp.ChooseRemoteForDeleteRemoteTagPopUpModel)
			if ok {
				selectedRemote := popUp.RemoteList.SelectedItem()
				if selectedRemote != nil {
					remote := selectedRemote.(tagPopUp.GitRemoteForDeleteRemoteTagItem)
					m.PopUpType = constant.DeleteTagOutputPopUp
					tagPopUp.InitDeleteTagOutputPopUpModel(m, popUp.TagName)
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					services.DeleteTagService(m, remote.Name, popUp.TagName, git.TAGDELETEREMOTE)

					// to return the tick for the spinner
					tagOutputPopUp, ok := m.PopUpModel.(*tagPopUp.DeleteTagOutputPopUpModel)
					if ok {
						return m, tagOutputPopUp.Spinner.Tick
					}
				}
			}
		case constant.ChoosePushTagOptionPopUp:
			popUp, ok := m.PopUpModel.(*tagPopUp.ChoosePushTagOptionPopUpModel)
			if ok {
				selectedPushType := popUp.PushOptionList.SelectedItem()
				originName := popUp.RemoteName
				tagName := popUp.TagName
				if selectedPushType != nil {
					m.PopUpType = constant.PushTagOutputPopUp
					tagPopUp.InitPushTagOutputPopUpModel(m)
					popUp, ok := m.PopUpModel.(*tagPopUp.PushTagOutputPopUpModel)
					if ok {
						services.GitPushTagService(m, originName, tagName, selectedPushType.(tagPopUp.PushTagOptionItem).PushTagType)
						return m, popUp.Spinner.Tick
					}
				} else {
					m.ShowPopUp.Store(false)
					m.IsTyping.Store(false)
					m.PopUpModel = nil
					m.PopUpType = constant.NoPopUp
					return m, nil
				}
			}
		case constant.ChooseFetchTagOptionPopUp:
			popUp, ok := m.PopUpModel.(*tagPopUp.ChooseFetchTagOptionPopUpModel)
			if ok {
				selectedFetchType := popUp.FetchOptionList.SelectedItem()
				originName := popUp.RemoteName
				if selectedFetchType != nil {
					m.PopUpType = constant.FetchTagOutputPopUp
					tagPopUp.InitFetchTagOutputPopUpModel(m)
					popUp, ok := m.PopUpModel.(*tagPopUp.FetchTagOutputPopUpModel)
					if ok {
						services.GitFetchTagService(m, originName, selectedFetchType.(tagPopUp.FetchTagOptionItem).FetchTagType)
						return m, popUp.Spinner.Tick
					}
				} else {
					m.ShowPopUp.Store(false)
					m.IsTyping.Store(false)
					m.PopUpModel = nil
					m.PopUpType = constant.NoPopUp
					return m, nil
				}
			}
		}
	}
	return m, nil
}

func handleNonTypingTabKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	nextNavigation := m.CurrentSelectedComponentIndex + 1
	if nextNavigation < len(constant.ComponentPanelNavigationList) {
		m.CurrentSelectedComponentIndex = nextNavigation
		m.CurrentSelectedComponent = constant.ComponentPanelNavigationList[nextNavigation]
		m.DetailPanelParentComponent = ""
		layout.LeftPanelDynamicResize(m)
		services.FetchDetailComponentPanelInfoService(m, true)
	}
	return m, nil
}

func handleNonTypingShiftTabKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	previousNavigation := m.CurrentSelectedComponentIndex - 1
	if previousNavigation >= 0 {
		m.CurrentSelectedComponentIndex = previousNavigation
		m.CurrentSelectedComponent = constant.ComponentPanelNavigationList[previousNavigation]
		m.DetailPanelParentComponent = ""
		layout.LeftPanelDynamicResize(m)
		services.FetchDetailComponentPanelInfoService(m, true)
	}
	return m, nil
}

func handleNonTypingSpaceKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.ModifiedFilesComponentPanel:
			currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
			var filePathName string
			if currentSelectedModifiedFile != nil {
				filePathName = currentSelectedModifiedFile.(files.GitModifiedFilesItem).FilePathname
				services.GitStageOrUnstageService(m, filePathName)
			}

		case constant.StashComponentPanel:
			selectedStashId := m.CurrentRepoStashInfoList.SelectedItem()
			if selectedStashId != nil {
				stashPopUp.InitGitStashConfirmPromptPopUpModel(m, git.APPLYSTASH, "", selectedStashId.(stash.GitStashItem).Id, selectedStashId.(stash.GitStashItem).Message)
				m.PopUpType = constant.GitStashConfirmPromptPopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
			}
		case constant.DetailComponentPanel, constant.DetailComponentPanelTwo:
			if m.IsLineEditingState.Load() {
				currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
				var filePathName string
				if currentSelectedModifiedFile != nil {
					filePathName = currentSelectedModifiedFile.(files.GitModifiedFilesItem).FilePathname
					services.GitStageLineOrUnstageLineService(m, filePathName)
				}
			}
		}
	} else {
		switch m.PopUpType {
		case constant.GitCherryPickPopUp:
			popUp, ok := m.PopUpModel.(*logPopUp.GitCherryPickPopUpModel)
			if ok {
				selectedCommitLog := popUp.CurrentBranchCherryPickCommitLog.SelectedItem()
				if selectedCommitLog != nil {
					cherryPickedCommitLog := selectedCommitLog.(logPopUp.GitCherryPickItem)
					CherryPickedCommitLogItem, ok := m.CherryPickedCommitInfo.CherryPickedCommitMap[cherryPickedCommitLog.Hash]
					if ok {
						delete(m.CherryPickedCommitInfo.CherryPickedCommitMap, CherryPickedCommitLogItem.Hash)
						if len(m.CherryPickedCommitInfo.CherryPickedCommitMap) < 1 {
							utils.ReinitCherryPickedCommitInfo(m)
						}
					} else {
						m.CherryPickedCommitInfo.CherryPickedCommitMap[cherryPickedCommitLog.Hash] = git.CherryPickedCommitLog{
							Hash:                 cherryPickedCommitLog.Hash,
							Message:              cherryPickedCommitLog.Message,
							Author:               cherryPickedCommitLog.Author,
							FromBranch:           cherryPickedCommitLog.FromBranch,
							UserSelectedSequence: m.CherryPickedCommitInfo.LatestSequenceCounter,
						}
						m.CherryPickedCommitInfo.LatestSequenceCounter++
					}
				}
			}
		}
	}
	return m, nil
}

func handleNonTypingEscKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.GitRemotePushPopUp:
			services.GitRemotePushCancelService(m)
		case constant.GitPullOutputPopUp:
			services.GitPullCancelService(m)
		case constant.PushTagOutputPopUp:
			services.GitPushTagCancelService(m)
		case constant.FetchTagOutputPopUp:
			services.GitFetchTagCancelService(m)
		case constant.DeleteTagOutputPopUp:
			services.DeleteTagCancelService(m)
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
		case constant.GitDeleteBranchOutputPopUp:
			popUp, ok := m.PopUpModel.(*branchPopUp.GitDeleteBranchOutputPopUpModel)
			if ok && !popUp.IsProcessing.Load() {
				// only close when done processing
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		case constant.CreateBranchBasedOnRemoteOutputPopUp:
			// Block ESC during create new branch based on remote operation - operation must complete
			popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemoteOutputPopUpModel)
			if ok && !popUp.IsProcessing.Load() {
				// only close when done processing
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		case constant.KeybindingAndFeatureInstructionsPopUp,
			constant.ChooseRemotePopUp,
			constant.ChoosePushTypePopUp,
			constant.ChooseNewBranchTypePopUp,
			constant.ChooseSwitchBranchTypePopUp,
			constant.ChooseGitPullTypePopUp,
			constant.GitDiscardTypeOptionPopUp,
			constant.GitDiscardConfirmPromptPopUp,
			constant.GitStashConfirmPromptPopUp,
			constant.GitResolveConflictOptionPopUp,
			constant.GitDeleteBranchConfirmPromptPopUp,
			constant.GitResetLatestCommitTypeOptionPopUp,
			constant.GitResetLatestCommitConfirmPromptPopUp,
			constant.GitResetToSelectedCommitTypeOptionPopUp,
			constant.GitResetToSelectedCommitConfirmPromptPopUp,
			constant.GitCherryPickPopUp,
			constant.GitEditCherryPickPopUp,
			constant.GitCherryPickOptionSelectionPopUp,
			constant.GitCherryPickApplyConfirmPopUp,
			constant.GitDiscardFileLineChangeConfirmPopUp,
			constant.CreateTagConfirmationPopUp,
			constant.ChooseDeleteTagOptionPopUp,
			constant.ChooseRemoteForDeleteRemoteTagPopUp,
			constant.ChoosePushTagOptionPopUp,
			constant.ChooseFetchTagOptionPopUp:
			// simple closing of the pop up
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpType = constant.NoPopUp
			m.PopUpModel = nil
		}
		return m, nil
	} else {
		switch m.CurrentSelectedComponent {
		case constant.DetailComponentPanel:
			if m.IsLineEditingState.Load() {
				m.IsLineEditingState.Store(false)
				m.TuiUpdateChannel <- constant.DETAIL_COMPONENT_PANEL_UPDATED
			} else {
				m.CurrentSelectedComponent = m.DetailPanelParentComponent
				m.DetailPanelParentComponent = ""
			}
		case constant.DetailComponentPanelTwo:
			if m.IsLineEditingState.Load() {
				m.IsLineEditingState.Store(false)
				m.TuiUpdateChannel <- constant.DETAIL_COMPONENT_PANEL_UPDATED
			} else {
				m.CurrentSelectedComponent = m.DetailPanelParentComponent
				m.DetailPanelParentComponent = ""
			}
		}
	}
	return m, nil
}

func handleNonTypingUpkKeyBindingInteraction(msg tea.KeyMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagComponentPanel:
			// we don't use the list native Update() because we track the current selected index
			switch m.CurrentLocalBranchOrTagComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				if m.CurrentRepoBranchesInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoBranchesInfoList.Index() - 1
					m.CurrentRepoBranchesInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.LocalBranchComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			case constant.SHOW_TAG:
				if m.CurrentRepoTagInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoTagInfoList.Index() - 1
					m.CurrentRepoTagInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.TagComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			}
		case constant.ModifiedFilesComponentPanel:
			// we don't use the list native Update() because we need to also track the current selected index
			if m.CurrentRepoModifiedFilesInfoList.Index() > 0 {
				latestIndex := m.CurrentRepoModifiedFilesInfoList.Index() - 1
				m.CurrentRepoModifiedFilesInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.ModifiedFilesComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.CommitLogComponentPanel:
			// we don't use the list native Update() because we need to also track the current selected index
			if m.CurrentRepoCommitLogInfoList.Index() > 0 {
				latestIndex := m.CurrentRepoCommitLogInfoList.Index() - 1
				m.CurrentRepoCommitLogInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.CommitLogComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.StashComponentPanel:
			// we don't use the list native Update() because we need to also track the current selected index
			if m.CurrentRepoStashInfoList.Index() > 0 {
				latestIndex := m.CurrentRepoStashInfoList.Index() - 1
				m.CurrentRepoStashInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.StashComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.DetailComponentPanel:
			if m.IsLineEditingState.Load() {
				m.LineEditingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex = max(0, m.LineEditingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex-1)
				if m.LineEditingIndexPositionAndInfo.DetailPanelViewportIndexPosition < 1 {
					m.DetailPanelViewport.ScrollUp(1)
				} else {
					m.LineEditingIndexPositionAndInfo.DetailPanelViewportIndexPosition -= 1
				}
				services.SetLineEditingCursorViewportContent(m, m.DetailPanelViewport.VisibleLineCount(), m.DetailPanelTwoViewport.VisibleLineCount())
			} else {
				m.DetailPanelViewport, cmd = m.DetailPanelViewport.Update(msg)
				return m, cmd
			}
		case constant.DetailComponentPanelTwo:
			if m.IsLineEditingState.Load() {
				m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportActualCurrentIndex = max(0, m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportActualCurrentIndex-1)
				if m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportIndexPosition < 1 {
					m.DetailPanelTwoViewport.ScrollUp(1)
				} else {
					m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportIndexPosition -= 1
				}
				services.SetLineEditingCursorViewportContent(m, m.DetailPanelViewport.VisibleLineCount(), m.DetailPanelTwoViewport.VisibleLineCount())
			} else {
				m.DetailPanelTwoViewport, cmd = m.DetailPanelTwoViewport.Update(msg)
				return m, cmd
			}
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
		case constant.LocalBranchOrTagComponentPanel:
			// we don't use the list native Update() because we track the current selected index
			switch m.CurrentLocalBranchOrTagComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				if m.CurrentRepoBranchesInfoList.Index() < len(m.CurrentRepoBranchesInfoList.Items())-1 {
					latestIndex := m.CurrentRepoBranchesInfoList.Index() + 1
					m.CurrentRepoBranchesInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.LocalBranchComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			case constant.SHOW_TAG:
				if m.CurrentRepoTagInfoList.Index() < len(m.CurrentRepoTagInfoList.Items())-1 {
					latestIndex := m.CurrentRepoTagInfoList.Index() + 1
					m.CurrentRepoTagInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.TagComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			}
		case constant.ModifiedFilesComponentPanel:
			// we don't use the list native Update() because we need to also track the current selected index
			if m.CurrentRepoModifiedFilesInfoList.Index() < len(m.CurrentRepoModifiedFilesInfoList.Items())-1 {
				latestIndex := m.CurrentRepoModifiedFilesInfoList.Index() + 1
				m.CurrentRepoModifiedFilesInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.ModifiedFilesComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.CommitLogComponentPanel:
			// we don't use the list native Update() because we need to also track the current selected index
			if m.CurrentRepoCommitLogInfoList.Index() < len(m.CurrentRepoCommitLogInfoList.Items())-1 {
				latestIndex := m.CurrentRepoCommitLogInfoList.Index() + 1
				m.CurrentRepoCommitLogInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.CommitLogComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.StashComponentPanel:
			// we don't use the list native Update() because we need to also track the current selected index
			if m.CurrentRepoStashInfoList.Index() < len(m.CurrentRepoStashInfoList.Items())-1 {
				latestIndex := m.CurrentRepoStashInfoList.Index() + 1
				m.CurrentRepoStashInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.StashComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.DetailComponentPanel:
			if m.IsLineEditingState.Load() {
				m.LineEditingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex = min(m.DetailPanelViewport.TotalLineCount()-1, m.LineEditingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex+1)
				if m.LineEditingIndexPositionAndInfo.DetailPanelViewportIndexPosition >= m.DetailPanelViewport.VisibleLineCount()-1 {
					m.DetailPanelViewport.ScrollDown(1)
				} else {
					m.LineEditingIndexPositionAndInfo.DetailPanelViewportIndexPosition += 1
				}
				services.SetLineEditingCursorViewportContent(m, m.DetailPanelViewport.VisibleLineCount(), m.DetailPanelTwoViewport.VisibleLineCount())
			} else {
				m.DetailPanelViewport, cmd = m.DetailPanelViewport.Update(msg)
				return m, cmd
			}
		case constant.DetailComponentPanelTwo:
			if m.IsLineEditingState.Load() {
				m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportActualCurrentIndex = min(m.DetailPanelTwoViewport.TotalLineCount()-1, m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportActualCurrentIndex+1)
				if m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportIndexPosition >= m.DetailPanelTwoViewport.VisibleLineCount()-1 {
					m.DetailPanelTwoViewport.ScrollDown(1)
				} else {
					m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportIndexPosition += 1
				}
				services.SetLineEditingCursorViewportContent(m, m.DetailPanelViewport.VisibleLineCount(), m.DetailPanelTwoViewport.VisibleLineCount())
			} else {
				m.DetailPanelTwoViewport, cmd = m.DetailPanelTwoViewport.Update(msg)
				return m, cmd
			}
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
		case constant.DetailComponentPanel:
			m.DetailPanelViewport.ScrollLeft(1)
		case constant.DetailComponentPanelTwo:
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
		case constant.GitDiscardFileLineChangeConfirmPopUp:
			popUp, ok := m.PopUpModel.(*filesPopUp.GitDiscardFileLineChangeConfirmPopUpModel)
			if ok {
				popUp.DiscardFileLineChangeViewport.ScrollLeft(1)
			}
		case constant.KeybindingAndFeatureInstructionsPopUp:
			popUp, ok := m.PopUpModel.(*keybindingPopUp.KeybindingAndFeatureInstructionsPopUpModel)
			if ok {
				scrollSpeed := 1
				if strings.ToUpper(settings.GITTICONFIGSETTINGS.LanguageCode) != "EN" {
					// other than en, all other i18n we support are both zh and jp which each rune takes up twice the width of en character
					scrollSpeed = 2
				}

				popUp.GlobalKeyBindingViewport.ScrollLeft(scrollSpeed)
			}
		}
	}
	return m, nil
}

func handleNonTypingRightlKeyBindingInteraction(msg tea.KeyMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.DetailComponentPanel:
			m.DetailPanelViewport.ScrollRight(1)
		case constant.DetailComponentPanelTwo:
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
		case constant.GitDiscardFileLineChangeConfirmPopUp:
			popUp, ok := m.PopUpModel.(*filesPopUp.GitDiscardFileLineChangeConfirmPopUpModel)
			if ok {
				popUp.DiscardFileLineChangeViewport.ScrollRight(1)
			}
		case constant.KeybindingAndFeatureInstructionsPopUp:
			popUp, ok := m.PopUpModel.(*keybindingPopUp.KeybindingAndFeatureInstructionsPopUpModel)
			if ok {
				scrollSpeed := 1
				if strings.ToUpper(settings.GITTICONFIGSETTINGS.LanguageCode) != "EN" {
					// other than en, all other i18n we support are both zh and jp which each rune takes up twice the width of en character
					scrollSpeed = 2
				}
				popUp.GlobalKeyBindingViewport.ScrollRight(scrollSpeed)
			}
		}
	}
	return m, nil
}

// handleNonTypingLeftBracketKeyBindingInteraction handles the '[' key not only for navigation but contextually to switch to the previous detail component panel
func handleNonTypingLeftBracketKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		// handle detail component panel switching
		if m.CurrentSelectedComponent == constant.DetailComponentPanelTwo {
			m.CurrentSelectedComponent = constant.DetailComponentPanel
		}
	}
	return m, nil
}

// handleNonTypingRightBracketKeyBindingInteraction handles the ']' key not only for navigation but contextually to switch to the next detail component panel
func handleNonTypingRightBracketKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		// handle detail component panel switching
		if m.CurrentSelectedComponent == constant.DetailComponentPanel && m.ShowDetailPanelTwo.Load() {
			m.CurrentSelectedComponent = constant.DetailComponentPanelTwo
		}
	}
	return m, nil
}

func handleNonTypingSlashKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent != constant.LogComponentPanel {
			m.CurrentSelectedComponent = constant.LogComponentPanel
			m.DetailPanelParentComponent = ""
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	}
	return m, nil
}

func handleNonTypingLeftAngleBracketKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagComponentPanel:
			switch m.CurrentLocalBranchOrTagComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
			// do nothing, as local branch will be the most left option in the local branch or tagcomponent panel
			case constant.SHOW_TAG:
				m.CurrentLocalBranchOrTagComponentShowing = constant.SHOW_LOCAL_BRANCH
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		}
	}
	return m, nil
}

func handleNonTypingRightAngleBracketKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagComponentPanel:
			switch m.CurrentLocalBranchOrTagComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				m.CurrentLocalBranchOrTagComponentShowing = constant.SHOW_TAG
				services.FetchDetailComponentPanelInfoService(m, true)
			case constant.SHOW_TAG:
				// do nothing, as tag is currently the most right option in the local branch or tag component panel
			}
		}
	}
	return m, nil
}

func handleNonTypingCtrlaKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		services.GitStateUniversalUtilsAbortService(m)
	}
	return m, nil
}

func handleNonTypingCtrlkKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		services.GitStateUniversalUtilsSkipService(m)
	}
	return m, nil
}

func handleNonTypingCtrlpKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.GitEditCherryPickPopUp:
			m.PopUpType = constant.GitCherryPickPopUp
			logPopUp.InitGitCherryPickPopUp(m, m.CheckOutBranch)
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
		case constant.GitCherryPickApplyConfirmPopUp:
			m.PopUpType = constant.GitCherryPickPopUp
			logPopUp.InitGitCherryPickPopUp(m, m.CheckOutBranch)
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
		}
	} else {
		if (m.CurrentSelectedComponent == constant.CommitLogComponentPanel || m.DetailPanelParentComponent == constant.CommitLogComponentPanel) &&
			len(m.CurrentRepoCommitLogInfoList.Items()) > 0 {
			if len(m.CherryPickedCommitInfo.CherryPickedCommitMap) < 1 {
				m.ShowPopUp.Store(true)
				m.PopUpType = constant.GitCherryPickPopUp
				m.IsTyping.Store(false)
				logPopUp.InitGitCherryPickPopUp(m, m.CheckOutBranch)
			} else {
				m.ShowPopUp.Store(true)
				m.PopUpType = constant.GitCherryPickOptionSelectionPopUp
				m.IsTyping.Store(false)
				logPopUp.InitGitCherryPickOptionSelectionPopUp(m)
			}
		} else if (m.CurrentSelectedComponent == constant.LocalBranchOrTagComponentPanel || m.DetailPanelParentComponent == constant.LocalBranchOrTagComponentPanel) &&
			len(m.CurrentRepoTagInfoList.Items()) > 0 && m.CurrentLocalBranchOrTagComponentShowing == constant.SHOW_TAG {
			selectedTag := m.CurrentRepoTagInfoList.SelectedItem()
			if selectedTag == nil {
				return m, nil
			}
			if !m.GitOperations.GitRemote.CheckRemoteExist() {
				// if no remote found, we add one
				m.PopUpType = constant.AddRemotePromptPopUp
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
				remotes := m.GitOperations.GitRemote.PushRemote()
				if len(remotes) == 1 {
					m.PopUpType = constant.ChoosePushTagOptionPopUp
					// only one remote found so, we will default to that remote
					tagPopUp.InitChoosePushTagOptionPopUpModel(m, remotes[0].Name, selectedTag.(tag.GitTagItem).TagName)
				} else if len(remotes) > 1 {
					// if remote is more than 1 let user choose which remote
					m.PopUpType = constant.ChooseRemotePopUp
					remotePopUp.InitChooseRemotePopUpModel(m, remotes, constant.TAGPUSHACTION)
				}
			}
		}
	}
	return m, nil
}
