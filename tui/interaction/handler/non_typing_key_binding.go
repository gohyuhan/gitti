package handler

import (
	"slices"
	"strings"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/component/branch"
	"github.com/gohyuhan/gitti/tui/component/commitlog"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/component/reflog"
	"github.com/gohyuhan/gitti/tui/component/remote"
	"github.com/gohyuhan/gitti/tui/component/stash"
	"github.com/gohyuhan/gitti/tui/component/tag"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/layout"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	commitLogPopUp "github.com/gohyuhan/gitti/tui/popup/commitlog"
	discardPopUp "github.com/gohyuhan/gitti/tui/popup/discard"
	filesPopUp "github.com/gohyuhan/gitti/tui/popup/files"
	keybindingPopUp "github.com/gohyuhan/gitti/tui/popup/keybinding"
	pullPopUp "github.com/gohyuhan/gitti/tui/popup/pull"
	pushPopUp "github.com/gohyuhan/gitti/tui/popup/push"
	rebasePopUp "github.com/gohyuhan/gitti/tui/popup/rebase"
	reflogPopUp "github.com/gohyuhan/gitti/tui/popup/reflog"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	resolvePopUp "github.com/gohyuhan/gitti/tui/popup/resolve"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ----------------------------------
//
//	Handle global non-typing key interactions (e.g., '?').
//	Responsibility: Opens the global Keybinding and Feature Instructions pop-up,
//	suspending the typing state to display application-wide shortcuts and help info.
//
// ----------------------------------
func handleNonTypingGlobalKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	m.ShowPopUp.Store(true)
	m.IsTyping.Store(false)
	m.PopUpType = constant.KeybindingAndFeatureInstructionsPopUp
	keybindingPopUp.InitKeybindingAndFeatureInstructionsPopUpModel(m)
	return m, nil
}

// ----------------------------------
//
//	Handle '1' key interaction.
//	Responsibility: Switches the active view focus to the first main layout tab
//	(Local Branch, Tag, or Remote Component Panel), triggering an update to the
//	detail panel dynamically.
//
// ----------------------------------
func handleNonTyping1KeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent != constant.LocalBranchOrTagOrRemoteComponentPanel {
			m.CurrentSelectedComponent = constant.LocalBranchOrTagOrRemoteComponentPanel
			m.CurrentSelectedComponentIndex = 1
			m.DetailPanelParentComponent = ""
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle '2' key interaction.
//	Responsibility: Switches the active view focus to the second main layout tab
//	(Modified Files Component Panel), triggering an update to the detail panel
//	to reflect currently modified, added, or deleted repository files.
//
// ----------------------------------
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

// ----------------------------------
//
//	Handle '3' key interaction.
//	Responsibility: Switches the active view focus to the third main layout tab
//	(Commit Log Component Panel), refreshing the detail view to display git log history.
//
// ----------------------------------
func handleNonTyping3KeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent != constant.CommitLogOrRefLogComponentPanel {
			m.CurrentSelectedComponent = constant.CommitLogOrRefLogComponentPanel
			m.CurrentSelectedComponentIndex = 3
			m.DetailPanelParentComponent = ""
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle '4' key interaction.
//	Responsibility: Switches the active view focus to the fourth main layout tab
//	(Stash Component Panel), updating the detail panel to show current git stashes.
//
// ----------------------------------
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

// ----------------------------------
//
//	Handle 'a' key interaction.
//	Responsibility: Used contextually within popups. Specifically handles 'apply' actions
//	for git cherry-pick popups, shifting the model state to confirm the application
//	of cherry-picked commits.
//
// ----------------------------------
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

// ----------------------------------
//
//	Handle 'A' key interaction.
//	Responsibility: Initiates the "Amend Commit" flow. Opens the amend commit popup,
//	clears previous commit output, initializes the input model for the message,
//	and switches the application to typing state.
//
// ----------------------------------
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

// ----------------------------------
//
//	Handle 'c' key interaction.
//	Responsibility: Initiates the "New Commit" flow. Opens the standard commit popup,
//	resets the commit output view, initializes the input model if necessary,
//	and places the user in typing mode to enter commit details.
//
// ----------------------------------
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

// ----------------------------------
//
//	Handle 'C' key interaction.
//	Responsibility: Handles continuation operations. Primarily triggers "git commit --continue",
//	"git rebase --continue", or "git cherry-pick --continue". Also handles the specific UI
//	suspension needed for GPG signing processes during continuations.
//
// ----------------------------------
func handleNonTypingCKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.GitCommitRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
			gitArgs := m.GitOperations.GitStateUniversalUtils.GitUniversalContinueWithSigning()
			if len(gitArgs) < 1 {
				return m, nil
			}
			return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.COMMIT_WITH_SIGNING_OPS)
		} else {
			services.GitStateUniversalUtilsContinueService(m)
			return m, nil
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle 'd' key interaction.
//	Responsibility: Multi-contextual "discard" or "delete" operation. Depending on the focused panel, it handles:
//	- Deleting local branches, remote references, or tags (prompting confirmations).
//	- Dropping git stashes.
//	- Discarding changes in modified files (with complex checking for staged, unstaged, untracked, and renamed states).
//	- Discarding specific hunk/line changes when in detail line-editing mode.
//	Also handles popup dismissals for specific cherry-pick list items.
//
// ----------------------------------
func handleNonTypingdKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteComponentShowing {
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
			case constant.SHOW_REMOTE:
				selectedRemoteItem := m.CurrentRepoRemoteInfoList.SelectedItem()
				if selectedRemoteItem != nil {
					remoteItem := selectedRemoteItem.(remote.GitRemoteItem)
					remotePopUp.InitRemoveRemoteConfirmationPopUpModel(m, remoteItem.Name, remoteItem.Url, remoteItem.Fetch, remoteItem.Push)
					m.PopUpType = constant.RemoveRemoteConfirmationPopUp
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
			popUp, ok := m.PopUpModel.(*commitLogPopUp.GitEditCherryPickPopUpModel)
			if ok {
				selectedCherryPickedCommiyLogItem := popUp.CherryPickedCommitLog.SelectedItem()
				if selectedCherryPickedCommiyLogItem != nil {
					selectedCherryPickedCommiyLogHash := selectedCherryPickedCommiyLogItem.(commitLogPopUp.GitEditCherryPickItem).Hash
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
						commitLogPopUp.InitGitEditCherryPickPopUp(m, popUp.CherryPickedCommitLog.Index())
					}
				}
			}
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle 'e' key interaction.
//	Responsibility: Contextual "edit" action.
//	- In popups: Switches a cherry-pick view into "edit cherry-pick" mode.
//	- In Remote Panel: Opens the prompt to edit an existing remote's URL/Name.
//	- In Modified Files Panel: Launches the user's defined system editor (handling terminal/GUI diffs).
//	- In Log Panel: Triggers exporting the internal application logs.
//
// ----------------------------------
func handleNonTypingeKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.GitCherryPickPopUp:
			m.PopUpType = constant.GitEditCherryPickPopUp
			commitLogPopUp.InitGitEditCherryPickPopUp(m, 0)
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
		case constant.GitCherryPickApplyConfirmPopUp:
			m.PopUpType = constant.GitEditCherryPickPopUp
			commitLogPopUp.InitGitEditCherryPickPopUp(m, 0)
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
		}
	} else {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteComponentShowing {
			case constant.SHOW_REMOTE:
				selectedRemote := m.CurrentRepoRemoteInfoList.SelectedItem()
				if selectedRemote != nil {
					remoteItem := selectedRemote.(remote.GitRemoteItem)
					m.PopUpType = constant.EditRemotePromptPopUp
					remotePopUp.InitEditRemotePromptPopUpModel(m, remoteItem.Name, remoteItem.Url)
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(true)
				}
			}
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

// ----------------------------------
//
//	Handle 'f' key interaction.
//	Responsibility: Initiates the "fetch" operation. Specifically checks if the user is in the
//	Tag component view. If no remotes exist, prompts to add one. If remotes exist,
//	either directly fetches tags (if 1 remote) or opens a popup to select which remote to fetch tags from.
//
// ----------------------------------
func handleNonTypingfKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if (m.CurrentSelectedComponent == constant.LocalBranchOrTagOrRemoteComponentPanel || m.DetailPanelParentComponent == constant.LocalBranchOrTagOrRemoteComponentPanel) &&
			m.CurrentLocalBranchOrTagOrRemoteComponentShowing == constant.SHOW_TAG {
			if !m.GitOperations.GitRemote.CheckRemoteExist(false) {
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

// ----------------------------------
//
//	Handle 'L' key interaction.
//	Responsibility: Enters "Line Editing Mode" (also known as hunk staging mode).
//	Allows the user to specifically stage, unstage or discard individual lines of a modified file
//	instead of the entire file, updating the detail panel UI to reflect this state.
//
// ----------------------------------
func handleNonTypingLKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	// enter line stage state
	services.EnterOrReinitLineEditingStateService(m)
	m.TuiUpdateChannel <- constant.DETAIL_COMPONENT_PANEL_UPDATED
	return m, nil
}

// ----------------------------------
//
//	Handle 'm' key interaction.
//	Responsibility: Initiates the "Merge" workflow.
//	In the Local Branch panel, opens the branch selection popup allowing the user to
//	choose one or more branches to merge into the currently checked-out branch.
//
// ----------------------------------
func handleNonTypingmKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				m.PopUpType = constant.ChooseBranchOptionForMergePopUp
				m.IsTyping.Store(false)
				m.ShowPopUp.Store(true)
				branchPopUp.InitChooseBranchOptionForMergePopUpModel(m)
			}
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle 'n' key interaction.
//	Responsibility: Contextual "new" operation. Depending on the focused view:
//	- In Local Branch View: Opens popup to create a new branch (optionally based on a remote).
//	- In Remote View: Opens a prompt to add a new remote connection to the repository.
//
// ----------------------------------
func handleNonTypingnKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				m.PopUpType = constant.ChooseNewBranchTypePopUp
				m.IsTyping.Store(false)
				m.ShowPopUp.Store(true)
				if _, ok := m.PopUpModel.(*branchPopUp.ChooseNewBranchTypeOptionPopUpModel); !ok {
					branchPopUp.InitChooseNewBranchTypePopUpModel(m)
				}
			case constant.SHOW_REMOTE:
				m.PopUpType = constant.AddRemotePromptPopUp
				m.IsTyping.Store(true)
				m.ShowPopUp.Store(true)
				if m.GitOperations.GitRemote.CheckRemoteExist(false) {
					remotePopUp.InitAddRemotePromptPopUpModel(m, false)
				} else {
					remotePopUp.InitAddRemotePromptPopUpModel(m, true)
				}
			}
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_REFLOG:
				selectedReflog := m.CurrentRepoRefLogInfoList.SelectedItem()
				if selectedReflog != nil {
					parsedReflog := selectedReflog.(reflog.GitRefLogItem)
					m.PopUpType = constant.CreateNewBranchPopUp
					m.IsTyping.Store(true)
					m.ShowPopUp.Store(true)
					branchPopUp.InitCreateNewBranchPopUpModel(m, git.NEWBRANCHBASEDONCOMMITHASH, parsedReflog.Hash)
				}
			}
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle 'p' key interaction.
//	Responsibility: Initiates the "push" operation.
//	Checks for existing remotes; if none exist, prompts to add one.
//	If remotes exist, it figures out the appropriate remote to push to,
//	often opening a popup for the user to select the push argument (e.g., force push)
//	or the specific remote if multiple exist.
//
// ----------------------------------
func handleNonTypingpKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		// first we need to check if there are any push/pull origin origin for this repo
		// if not we prompt the user to add a new remote origin
		if !m.GitOperations.GitRemote.CheckRemoteExist(false) {
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

// ----------------------------------
//
//	Handle 'P' key interaction.
//	Responsibility: Initiates the "pull" operation.
//	Checks for existing remotes; if none exist, prompts to add one.
//	If remotes exist, opens a popup allowing the user to select the type
//	of pull operation (e.g., normal pull, pull rebase).
//
// ----------------------------------
func handleNonTypingPKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		// first we need to check if there are any push/pull origin for this repo
		// if not we prompt the user to add a new remote origin
		if !m.GitOperations.GitRemote.CheckRemoteExist(false) {
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

// ----------------------------------
//
//	Handle 'r' key interaction.
//	Responsibility: Contextual "resolve" or "reset" operations depending on the view:
//	- In Local Branch View: Opens a popup to rebase the current branch to a specific commit. (the selector must be on the current checkout branch to trigger this)
//	- In Modified Files Panel: Opens the conflict resolution popup for a conflicted file.
//	- In Commit Log Panel: Opens a popup to reset the repository HEAD to a specifically
//	  selected older commit (soft, mixed, or hard reset).
//
// ----------------------------------
func handleNonTypingrKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				selectedBranch := m.CurrentRepoBranchesInfoList.SelectedItem()
				if selectedBranch != nil {
					parsedBranch := selectedBranch.(branch.GitBranchItem)
					if parsedBranch.IsCheckedOut {
						_ = m.GitOperations.GitRemote.CheckRemoteExist(false)
						m.IsTyping.Store(false)
						m.ShowPopUp.Store(true)
						remotes := m.GitOperations.GitRemote.FetchRemote()
						m.PopUpType = constant.ChooseRemotePopUp
						// add a use local branch option which the name and url is "" (empty string) and fetch/push is false
						remotes = append([]git.GitRemoteInfo{{Name: "", Url: "", Fetch: false, Push: false}}, remotes...)
						remotePopUp.InitChooseRemotePopUpModel(m, remotes, constant.REBASEACTION)
					} else {
						return m, nil
					}
				}
			}
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
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
				selectedCommit := m.CurrentRepoCommitLogInfoList.SelectedItem()
				if selectedCommit != nil {
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
			case constant.SHOW_REFLOG:
				selectedRefLog := m.CurrentRepoRefLogInfoList.SelectedItem()
				if selectedRefLog != nil {
					parsedRefLog := selectedRefLog.(reflog.GitRefLogItem)
					commitPopUp.InitGitResetToSelectedCommitTypeOptionPopUpModel(
						m,
						parsedRefLog.Hash,
						"",
						"",
					)
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					m.PopUpType = constant.GitResetToSelectedCommitTypeOptionPopUp
				}
			}
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle 'R' key interaction.
//	Responsibility: Contextual "reset latest" operation.
//	Specifically in the Commit Log Panel, this opens a popup allowing the user
//	to reset the repository HEAD exactly one commit backwards (undoing the latest commit),
//	offering soft, mixed, or hard reset options.
//
// ----------------------------------
func handleNonTypingRKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
				if len(m.CurrentRepoCommitLogInfoList.Items()) > 1 {
					commitPopUp.InitGitResetLatestCommitTypeOptionPopUpModel(m)
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					m.PopUpType = constant.GitResetLatestCommitTypeOptionPopUp
				}
			}
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle 's' key interaction.
//	Responsibility: Contextual "stash file" operation.
//	In the Modified Files Panel, if a valid (non-conflicted) uncommitted file is selected,
//	this triggers a popup to stash *only* that specific file, prompting for a stash message.
//
// ----------------------------------
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

// ----------------------------------
//
//	Handle 'S' key interaction.
//	Responsibility: Contextual "stash all" operation.
//	In the Modified Files Panel, this triggers a popup to stash *all* currently modified
//	tracked files in the repository, prompting the user for a stash message.
//
// ----------------------------------
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

// ----------------------------------
//
//	Handle 't' key interaction.
//	Responsibility: Initiates the "Create Tag" operation.
//	In the Commit Log Panel, selecting a specific commit and pressing 't'
//	opens a popup to create a new git tag referencing that exact commit hash.
//
// ----------------------------------
func handleNonTypingtKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch m.CurrentSelectedComponent {
	case constant.CommitLogOrRefLogComponentPanel:
		switch m.CurrentCommitLogOrRefLogComponentShowing {
		case constant.SHOW_COMMITLOG:
			currentSelectedCommit := m.CurrentRepoCommitLogInfoList.SelectedItem()
			if currentSelectedCommit != nil {
				commit := currentSelectedCommit.(commitlog.GitCommitLogItem)
				m.PopUpType = constant.CreateTagPopUp
				tagPopUp.InitCreateTagPopUpModel(m, commit.Hash, commit.Message)
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(true)
			}
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle 'q' or 'Q' key interaction.
//	Responsibility: Global "quit" or "exit" operation.
//	If no popup is currently active, this forcefully stops any running git daemon
//	and triggers Bubble Tea's quit command, closing the application.
//
// ----------------------------------
func handleNonTypingqQKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if api.GITDAEMON != nil {
			api.GITDAEMON.Stop()
		}
		return m, tea.Quit
	}
	return m, nil
}

// ----------------------------------
//
//	Handle Backspace key interaction.
//	Responsibility: Contextual "pop" or "revert" action.
//	- In Stash Panel: Triggers a popup to 'pop' (apply and drop) the selected stash.
//	- In Edit Cherry Pick Popup: Clears the cherry-picked commit list and dismisses the popup.
//
// ----------------------------------
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
			_, ok := m.PopUpModel.(*commitLogPopUp.GitEditCherryPickPopUpModel)
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

// ----------------------------------
//
//	Handle Enter key interaction.
//	Responsibility: Core confirmation and drill-down action. Depends heavily on context:
//	- Active Popup (e.g., choosing remote, selecting branch type, confirming discard): Executes the chosen workflow or triggers git commands (push, pull, reset, discard, tag operations, etc.).
//	- Component Panels (when no popup active): Typically drills down into a "detail view" for the selected item (e.g., viewing diffs for a file, showing commit details), or triggers context menus like switching branches.
//
// ----------------------------------
func handleNonTypingEnterKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.ModifiedFilesComponentPanel:
			if len(m.CurrentRepoModifiedFilesInfoList.Items()) > 0 {
				m.CurrentSelectedComponent = constant.DetailComponentPanel
				m.DetailPanelParentComponent = constant.ModifiedFilesComponentPanel
			}
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
				if len(m.CurrentRepoCommitLogInfoList.Items()) > 0 {
					m.CurrentSelectedComponent = constant.DetailComponentPanel
					m.DetailPanelParentComponent = constant.CommitLogOrRefLogComponentPanel
				}
			case constant.SHOW_REFLOG:
				if len(m.CurrentRepoRefLogInfoList.Items()) > 0 {
					m.CurrentSelectedComponent = constant.DetailComponentPanel
					m.DetailPanelParentComponent = constant.CommitLogOrRefLogComponentPanel
				}
			}
		case constant.StashComponentPanel:
			if len(m.CurrentRepoStashInfoList.Items()) > 0 {
				m.CurrentSelectedComponent = constant.DetailComponentPanel
				m.DetailPanelParentComponent = constant.StashComponentPanel
			}
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteComponentShowing {
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
			case constant.SHOW_REMOTE:
				currentSelectedRemote := m.CurrentRepoRemoteInfoList.SelectedItem()
				if currentSelectedRemote != nil {
					selectedRemote := currentSelectedRemote.(remote.GitRemoteItem)
					m.PopUpType = constant.RemoteAsTrackingUpstreamConfirmationPopUp
					m.IsTyping.Store(false)
					m.ShowPopUp.Store(true)
					remotePopUp.InitRemoteAsTrackingUpstreamConfirmationPopUpModel(m, selectedRemote.Name, selectedRemote.Url)
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
				case constant.REBASEACTION:
					m.IsTyping.Store(true)
					m.ShowPopUp.Store(true)
					m.PopUpType = constant.GitRebaseBranchInputPopUp
					rebasePopUp.InitGitRebaseBranchInputPopUpModel(m, remoteName)
				}
			}

		case constant.ChoosePushTypePopUp:
			popUp, ok := m.PopUpModel.(*pushPopUp.ChoosePushTypePopUpModel)
			if ok {
				selectedOption := popUp.PushOptionList.SelectedItem()
				if m.GitPushRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
					gitArgs := m.GitOperations.GitCommit.GitPushWithSigning(popUp.RemoteName, selectedOption.(pushPopUp.GitPushOptionItem).PushType, m.CheckOutBranch)
					return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.GIT_PUSH_WITH_SIGNING_OPS)
				} else {
					m.PopUpType = constant.GitRemotePushPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					return services.InitGitRemotePushPopUpModelAndStartGitRemotePushService(m, popUp.RemoteName, selectedOption.(pushPopUp.GitPushOptionItem).PushType)
				}
			}

		case constant.ChooseNewBranchTypePopUp:
			popUp, ok := m.PopUpModel.(*branchPopUp.ChooseNewBranchTypeOptionPopUpModel)
			if ok {
				selectedOption := popUp.NewBranchTypeOptionList.SelectedItem()
				newBranchType := selectedOption.(branchPopUp.GitNewBranchTypeOptionItem).NewBranchType
				switch newBranchType {
				case git.NEWBRANCHBASEDONREMOTEUSERINPUT:
					if !m.GitOperations.GitRemote.CheckRemoteExist(false) {
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
				case git.NEWBRANCHBASEDONREMOTEUSERSELECT:
					m.PopUpType = constant.ChooseRemoteBranchOptionPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					branchPopUp.InitChooseRemoteBranchOptionPopUpModel(m)
				default:
					m.PopUpType = constant.CreateNewBranchPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(true)
					branchPopUp.InitCreateNewBranchPopUpModel(m, newBranchType, "")
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
				selectedOption := popUp.PullTypeOptionList.SelectedItem().(pullPopUp.GitPullTypeOptionItem)

				if m.GitCommitRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
					gitArgs := m.GitOperations.GitPull.GitPullWithSigning(selectedOption.PullType)
					return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.GIT_PULL_WITH_SIGNING_OPS)
				} else {
					m.PopUpType = constant.GitPullOutputPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					pullPopUp.InitGitPullOutputPopUpModel(m)
					popUp, ok := m.PopUpModel.(*pullPopUp.GitPullOutputPopUpModel)
					if ok {
						popUp.IsProcessing.Store(true) // set it directly first
						// start the git pull service
						services.GitPullService(m, selectedOption.PullType)
						return m, popUp.Spinner.Tick
					}
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
			popUp, ok := m.PopUpModel.(*commitLogPopUp.GitCherryPickOptionSelectionPopUpModel)
			if ok {
				selectedCherryPickType := popUp.CherryPickedOpsOption.SelectedItem()
				if selectedCherryPickType != nil {
					cherryPickType := selectedCherryPickType.(commitLogPopUp.CherryPickOpsOptionItem).CherryPickOpsType
					switch cherryPickType {
					case constant.CHERRYPICK:
						m.PopUpType = constant.GitCherryPickPopUp
						commitLogPopUp.InitGitCherryPickPopUp(m, m.CheckOutBranch)
					case constant.EDITCHERRYPICK:
						m.PopUpType = constant.GitEditCherryPickPopUp
						commitLogPopUp.InitGitEditCherryPickPopUp(m, 0)
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
			if m.GitCommitRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
				var sortedCherryPickedCommitLogs []git.CherryPickedCommitLog
				// turn the hashmap into array first
				for _, commitLogItem := range m.CherryPickedCommitInfo.CherryPickedCommitMap {
					sortedCherryPickedCommitLogs = append(sortedCherryPickedCommitLogs, commitLogItem)
				}

				// sort the array based on user selection sequence
				slices.SortFunc(sortedCherryPickedCommitLogs, func(a, b git.CherryPickedCommitLog) int {
					return a.UserSelectedSequence - b.UserSelectedSequence
				})

				// harvest the commit hash
				var cherryPickedCommitHashes []string
				for _, commitLog := range sortedCherryPickedCommitLogs {
					cherryPickedCommitHashes = append(cherryPickedCommitHashes, commitLog.Hash)
				}

				gitArgs := m.GitOperations.GitCommitLog.GitCherryPickWithSigning(cherryPickedCommitHashes)
				return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.CHERRY_PICK_WITH_SIGNING_OPS)
			} else {
				services.GitCherryPickService(m, m.CherryPickedCommitInfo.CherryPickedCommitMap)
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpModel = nil
				m.PopUpType = constant.NoPopUp
				return m, nil
			}
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
				if m.GitTagRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
					gitArgs := m.GitOperations.GitTag.CreateNewTagWithSigning(popUp.CommitHash, popUp.TagName, popUp.TagMessage)
					return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.CREATE_NEW_TAG_WITH_SIGNING_OPS)
				} else {
					services.CreateNewTagService(m, popUp.CommitHash, popUp.TagName, popUp.TagMessage)
					m.PopUpType = constant.NoPopUp
					m.ShowPopUp.Store(false)
					m.IsTyping.Store(false)
					return m, nil
				}
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
						if !m.GitOperations.GitRemote.CheckRemoteExist(false) {
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
								if m.GitPushRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
									gitArgs := m.GitOperations.GitTag.GitDeleteRemoteTagWithSigning(remotes[0].Name, popUp.TagName)
									return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.TAG_REMOTE_DELETE_WITH_SIGNING_OPS)
								} else {
									m.PopUpType = constant.DeleteTagOutputPopUp
									tagPopUp.InitDeleteTagOutputPopUpModel(m, popUp.TagName)
									services.DeleteTagService(m, remotes[0].Name, popUp.TagName, git.TAGDELETEREMOTE)

									// to return the tick for the spinner
									tagOutputPopUp, ok := m.PopUpModel.(*tagPopUp.DeleteTagOutputPopUpModel)
									if ok {
										return m, tagOutputPopUp.Spinner.Tick
									}
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
					if m.GitPushRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
						gitArgs := m.GitOperations.GitTag.GitDeleteRemoteTagWithSigning(remote.Name, popUp.TagName)
						return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.TAG_REMOTE_DELETE_WITH_SIGNING_OPS)
					} else {
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
			}
		case constant.ChoosePushTagOptionPopUp:
			popUp, ok := m.PopUpModel.(*tagPopUp.ChoosePushTagOptionPopUpModel)
			if ok {
				selectedPushType := popUp.PushOptionList.SelectedItem()
				originName := popUp.RemoteName
				tagName := popUp.TagName
				if selectedPushType != nil {
					m.PopUpType = constant.PushTagOutputPopUp
					if m.GitPushRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
						gitArgs := m.GitOperations.GitTag.GitPushTagWithSigning(originName, tagName, selectedPushType.(tagPopUp.PushTagOptionItem).PushTagType)
						return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.TAG_PUSH_WITH_SIGNING_OPS)
					} else {
						tagPopUp.InitPushTagOutputPopUpModel(m)
						popUp, ok := m.PopUpModel.(*tagPopUp.PushTagOutputPopUpModel)
						if ok {
							services.GitPushTagService(m, originName, tagName, selectedPushType.(tagPopUp.PushTagOptionItem).PushTagType)
							return m, popUp.Spinner.Tick
						}
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
		case constant.RemoveRemoteConfirmationPopUp:
			popUp, ok := m.PopUpModel.(*remotePopUp.RemoveRemoteConfirmationPopUpModel)
			if ok {
				services.GitRemoveRemoteService(m, popUp.RemoteName)
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpModel = nil
				m.PopUpType = constant.NoPopUp
				return m, nil
			}
		case constant.RemoteAsTrackingUpstreamConfirmationPopUp:
			popUp, ok := m.PopUpModel.(*remotePopUp.RemoteAsTrackingUpstreamConfirmationPopUpModel)
			if ok {
				services.GitSetRemoteAsTrackingUpstreamService(m, popUp.RemoteName)
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpModel = nil
				m.PopUpType = constant.NoPopUp
				return m, nil
			}
		case constant.GitRevertParentOptionSelectionPopUp:
			popUp, ok := m.PopUpModel.(*commitLogPopUp.GitRevertParentOptionSelectionPopUpModel)
			if ok {
				selectedParent := popUp.GitRevertParentOption.SelectedItem()
				if selectedParent != nil {
					parsedSelectedParent := selectedParent.(commitLogPopUp.GitRevertParentOptionItem)
					commitLogPopUp.InitGitRevertConfirmationPopUp(m, popUp.CommitHash, parsedSelectedParent.ParentOrder)
				} else {
					commitLogPopUp.InitGitRevertConfirmationPopUp(m, popUp.CommitHash, 1)
				}
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				m.PopUpType = constant.GitRevertConfirmationPopUp
				return m, nil
			}
		case constant.GitRevertConfirmationPopUp:
			popUp, ok := m.PopUpModel.(*commitLogPopUp.GitRevertConfirmationPopUpModel)
			if ok {
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
				if m.GitCommitRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
					gitArgs := m.GitOperations.GitCommitLog.GitRevertCommitWithSigning(popUp.CommitHash, popUp.ParentOrder)
					return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.REVERT_COMMIT_WITH_SIGNING_OPS)
				} else {
					services.GitRevertCommitService(m, popUp.CommitHash, popUp.ParentOrder)
					return m, nil
				}
			}
		case constant.GitCherryPickFromRefLogApplyConfirmationPopUp:
			popUp, ok := m.PopUpModel.(*reflogPopUp.GitCherryPickFromRefLogApplyConfirmationPopUpModel)
			if ok {
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
				if m.GitCommitRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
					gitArgs := m.GitOperations.GitCommitLog.GitCherryPickWithSigning([]string{popUp.Hash})
					return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.CHERRY_PICK_WITH_SIGNING_OPS)
				} else {
					services.GitCherryPickReflogHashService(m, popUp.Hash)
					return m, nil
				}
			}
		case constant.ChooseRemoteBranchOptionPopUp:
			popUp, ok := m.PopUpModel.(*branchPopUp.ChooseRemoteBranchOptionPopUpModel)
			if ok {
				selectedRemoteBranch := popUp.RemoteBranchOptionList.SelectedItem()
				if selectedRemoteBranch != nil {
					branchName := selectedRemoteBranch.(branchPopUp.RemoteBranchItem).BranchName
					if utf8.RuneCountInString(branchName) > 0 {
						branchPopUp.InitCreateBranchBasedOnRemoteOutputPopUp(m)
						popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemoteOutputPopUpModel)
						if ok {
							m.ShowPopUp.Store(true)
							m.IsTyping.Store(false)
							m.PopUpType = constant.CreateBranchBasedOnRemoteOutputPopUp
							popUp.IsProcessing.Store(true)
							services.CreateNewBranchBasedOnRemoteService(m, "", branchName, git.NEWBRANCHBASEDONREMOTEUSERSELECT)
							return m, popUp.Spinner.Tick
						} else {
							m.ShowPopUp.Store(false)
							m.IsTyping.Store(false)
							m.PopUpType = constant.NoPopUp
						}
					}
				}
			}
		case constant.ChooseBranchOptionForMergePopUp:
			popUp, ok := m.PopUpModel.(*branchPopUp.ChooseBranchOptionForMergePopUpModel)
			if ok {
				var branchesNames []string
				for _, branch := range popUp.SelectedBranchList.Items() {
					branchesNames = append(branchesNames, branch.(branchPopUp.GitMergeBranchOptionItem).BranchName)
				}
				if m.GitPushRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
					gitArgs := m.GitOperations.GitBranch.GitMergeWithSigning(branchesNames)
					return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.GIT_MERGE_WITH_SIGNING_OPS)
				} else {
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					m.PopUpType = constant.BranchMergeOutputPopUp
					branchPopUp.InitBranchMergeOutputPopUpModel(m)
					popUp, ok := m.PopUpModel.(*branchPopUp.BranchMergeOutputPopUpModel)
					if ok {
						services.GitMergeService(m, branchesNames)
						return m, popUp.Spinner.Tick
					}
				}
			}
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle Tab key interaction.
//	Responsibility: Forward navigation.
//	- No popup: Cycles the primary active focus to the *next* main component panel in the UI layout
//	  (e.g., from Branches to Modified Files, then to Commit Log, etc.).
//	- Merge Popup: Shifts section focus from the branch option list to the selected branch list.
//
// ----------------------------------
func handleNonTypingTabKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		nextNavigation := m.CurrentSelectedComponentIndex + 1
		if nextNavigation < len(constant.ComponentPanelNavigationList) {
			m.CurrentSelectedComponentIndex = nextNavigation
			m.CurrentSelectedComponent = constant.ComponentPanelNavigationList[nextNavigation]
			m.DetailPanelParentComponent = ""
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	} else {
		switch m.PopUpType {
		case constant.ChooseBranchOptionForMergePopUp:
			popUp, ok := m.PopUpModel.(*branchPopUp.ChooseBranchOptionForMergePopUpModel)
			if ok {
				if popUp.BranchOptionSectionSelected.Load() {
					popUp.BranchOptionSectionSelected.Store(false)
					popUp.SelectedBranchSectionSelected.Store(true)
				}
			}
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle Shift+Tab key interaction.
//	Responsibility: Backward navigation.
//	- No popup: Cycles the primary active focus to the *previous* main component panel in the UI layout.
//	- Merge Popup: Shifts section focus from the selected branch list back to the branch option list.
//
// ----------------------------------
func handleNonTypingShiftTabKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		previousNavigation := m.CurrentSelectedComponentIndex - 1
		if previousNavigation >= 0 {
			m.CurrentSelectedComponentIndex = previousNavigation
			m.CurrentSelectedComponent = constant.ComponentPanelNavigationList[previousNavigation]
			m.DetailPanelParentComponent = ""
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	} else {
		switch m.PopUpType {
		case constant.ChooseBranchOptionForMergePopUp:
			popUp, ok := m.PopUpModel.(*branchPopUp.ChooseBranchOptionForMergePopUpModel)
			if ok {
				if popUp.SelectedBranchSectionSelected.Load() {
					popUp.BranchOptionSectionSelected.Store(true)
					popUp.SelectedBranchSectionSelected.Store(false)
				}
			}
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle Space key interaction.
//	Responsibility: Contextual "toggle" or "apply" action.
//	- Modified / Detail Panels: Stages or unstages the selected file or specific hunk/line.
//	- Stash Panel: Triggers applying a stash (without dropping it).
//	- Cherry Pick Popup: Selects/toggles the currently highlighted commit to be added to the cherry-pick queue.
//	- Merge Popup: Toggles the currently highlighted branch between the available and selected branch lists.
//
// ----------------------------------
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
			popUp, ok := m.PopUpModel.(*commitLogPopUp.GitCherryPickPopUpModel)
			if ok {
				selectedCommitLog := popUp.CurrentBranchCherryPickCommitLog.SelectedItem()
				if selectedCommitLog != nil {
					cherryPickedCommitLog := selectedCommitLog.(commitLogPopUp.GitCherryPickItem)
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
		case constant.ChooseBranchOptionForMergePopUp:
			_, ok := m.PopUpModel.(*branchPopUp.ChooseBranchOptionForMergePopUpModel)
			if ok {
				branchPopUp.UpdateChooseBranchOptionForMergePopUpModel(m)
			}
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle Esc key interaction.
//	Responsibility: Universal "cancel" or "back" operation.
//	- If operations are processing (pushing, pulling, etc.), blocking esc prevents interruption.
//	- Cancels out of active background git services if possible.
//	- Dismisses almost any active open popup, returning the user to the underlying view.
//	- If in Detail line-editing mode, exits that specific mode back to the parent panel.
//
// ----------------------------------
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
		case constant.GitRebaseOutputPopUp:
			services.GitRebaseCancelService(m)
		case constant.BranchMergeOutputPopUp:
			services.GitMergeCancelService(m)
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
			constant.ChooseFetchTagOptionPopUp,
			constant.RemoveRemoteConfirmationPopUp,
			constant.RemoteAsTrackingUpstreamConfirmationPopUp,
			constant.GitRevertParentOptionSelectionPopUp,
			constant.GitRevertConfirmationPopUp,
			constant.GitCherryPickFromRefLogApplyConfirmationPopUp,
			constant.ChooseRemoteBranchOptionPopUp,
			constant.ChooseBranchOptionForMergePopUp:
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

// ----------------------------------
//
//	Handle Up/k key interaction.
//	Responsibility: Vertical upward navigation.
//	- Native lists: Moves selection highlight up one item.
//	- Detail Viewports (text viewers & line-editing): Scrolls text up by one line.
//	- Popups: Navigates upward in popup selection lists natively.
//
// ----------------------------------
func handleNonTypingUpkKeyBindingInteraction(msg tea.KeyPressMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			// we don't use the list native Update() because we track the current selected index
			switch m.CurrentLocalBranchOrTagOrRemoteComponentShowing {
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
			case constant.SHOW_REMOTE:
				if m.CurrentRepoRemoteInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoRemoteInfoList.Index() - 1
					m.CurrentRepoRemoteInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.RemoteComponent = latestIndex
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
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
				// we don't use the list native Update() because we need to also track the current selected index
				if m.CurrentRepoCommitLogInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoCommitLogInfoList.Index() - 1
					m.CurrentRepoCommitLogInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.CommitLogComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			case constant.SHOW_REFLOG:
				// we don't use the list native Update() because we need to also track the current selected index
				if m.CurrentRepoRefLogInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoRefLogInfoList.Index() - 1
					m.CurrentRepoRefLogInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.RefLogComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
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
		return UpDownKeyPressMsgUpdateForPopUp(msg, m)
	}
	return m, nil
}

// ----------------------------------
//
//	Handle Down/j key interaction.
//	Responsibility: Vertical downward navigation.
//	- Native lists: Moves selection highlight down one item.
//	- Detail Viewports (text viewers & line-editing): Scrolls text down by one line.
//	- Popups: Navigates downward in popup selection lists natively.
//
// ----------------------------------
func handleNonTypingDownjKeyBindingInteraction(msg tea.KeyPressMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			// we don't use the list native Update() because we track the current selected index
			switch m.CurrentLocalBranchOrTagOrRemoteComponentShowing {
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
			case constant.SHOW_REMOTE:
				if m.CurrentRepoRemoteInfoList.Index() < len(m.CurrentRepoRemoteInfoList.Items())-1 {
					latestIndex := m.CurrentRepoRemoteInfoList.Index() + 1
					m.CurrentRepoRemoteInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.RemoteComponent = latestIndex
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
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
				// we don't use the list native Update() because we need to also track the current selected index
				if m.CurrentRepoCommitLogInfoList.Index() < len(m.CurrentRepoCommitLogInfoList.Items())-1 {
					latestIndex := m.CurrentRepoCommitLogInfoList.Index() + 1
					m.CurrentRepoCommitLogInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.CommitLogComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			case constant.SHOW_REFLOG:
				// we don't use the list native Update() because we need to also track the current selected index
				if m.CurrentRepoRefLogInfoList.Index() < len(m.CurrentRepoRefLogInfoList.Items())-1 {
					latestIndex := m.CurrentRepoRefLogInfoList.Index() + 1
					m.CurrentRepoRefLogInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.RefLogComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
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
		return UpDownKeyPressMsgUpdateForPopUp(msg, m)
	}
	return m, nil
}

// ----------------------------------
//
//	Handle Left/h key interaction.
//	Responsibility: Horizontal left navigation/scrolling.
//	Scrolls text viewers, detail panels, and wider popup views horizontally to the left.
//
// ----------------------------------
func handleNonTypingLefthKeyBindingInteraction(msg tea.KeyPressMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
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

// ----------------------------------
//
//	Handle Right/l key interaction.
//	Responsibility: Horizontal right navigation/scrolling.
//	Scrolls text viewers, detail panels, and wider popup views horizontally to the right.
//
// ----------------------------------
func handleNonTypingRightlKeyBindingInteraction(msg tea.KeyPressMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
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

// ----------------------------------
//
//	Handle '[' key interaction.
//	Responsibility: Detail panel backward navigation.
//	Switches the active focus from DetailComponentPanelTwo back to DetailComponentPanel (the primary detail panel).
//
// ----------------------------------
func handleNonTypingLeftBracketKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		// handle detail component panel switching
		if m.CurrentSelectedComponent == constant.DetailComponentPanelTwo {
			m.CurrentSelectedComponent = constant.DetailComponentPanel
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle ']' key interaction.
//	Responsibility: Detail panel forward navigation.
//	Switches the active focus from DetailComponentPanel to DetailComponentPanelTwo (the secondary detail panel), if it is currently visible.
//
// ----------------------------------
func handleNonTypingRightBracketKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		// handle detail component panel switching
		if m.CurrentSelectedComponent == constant.DetailComponentPanel && m.ShowDetailPanelTwo.Load() {
			m.CurrentSelectedComponent = constant.DetailComponentPanelTwo
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle '/' key interaction.
//	Responsibility: Focuses the Log Component Panel.
//	Used as a quick jumping shortcut to view the internal console/error logs of Gitti.
//
// ----------------------------------
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

// ----------------------------------
//
//	Handle '<' key interaction.
//	Responsibility: Component sub-state cycling (Backward).
//	In the unified 'Local Branch Or Tag Or Remote' panel, this specifically rotates the visible
//	list to the left/previous entity (e.g., from Remotes -> Tags -> Branches).
//
// ----------------------------------
func handleNonTypingLeftAngleBracketKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
			// do nothing, as local branch will be the most left option in the local branch or tag or remote component panel
			case constant.SHOW_TAG:
				m.CurrentLocalBranchOrTagOrRemoteComponentShowing = constant.SHOW_LOCAL_BRANCH
				services.FetchDetailComponentPanelInfoService(m, true)
			case constant.SHOW_REMOTE:
				m.CurrentLocalBranchOrTagOrRemoteComponentShowing = constant.SHOW_TAG
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
			// do nothing, as commit log will be the most left option in the commit log or reflog component panel
			case constant.SHOW_REFLOG:
				m.CurrentCommitLogOrRefLogComponentShowing = constant.SHOW_COMMITLOG
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle '>' key interaction.
//	Responsibility: Component sub-state cycling (Forward).
//	In the unified 'Local Branch Or Tag Or Remote' panel, this specifically rotates the visible
//	list to the right/next entity (e.g., from Branches -> Tags -> Remotes).
//
// ----------------------------------
func handleNonTypingRightAngleBracketKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				m.CurrentLocalBranchOrTagOrRemoteComponentShowing = constant.SHOW_TAG
				services.FetchDetailComponentPanelInfoService(m, true)
			case constant.SHOW_TAG:
				m.CurrentLocalBranchOrTagOrRemoteComponentShowing = constant.SHOW_REMOTE
				services.FetchDetailComponentPanelInfoService(m, true)
			case constant.SHOW_REMOTE:
				// do nothing, as remote is currently the most right option in the local branch or tag or remote component panel
			}
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
				m.CurrentCommitLogOrRefLogComponentShowing = constant.SHOW_REFLOG
				services.FetchDetailComponentPanelInfoService(m, true)
			case constant.SHOW_REFLOG:
				// do nothing, as reflog is currently the most right option in the commit log or reflog component panel
			}
		}
	}
	return m, nil
}

// ----------------------------------
//
//	Handle Ctrl+a key interaction.
//	Responsibility: Global "Abort" operation.
//	Attempts to actively abort ongoing multi-step git state operations
//	like an incomplete merge, rebase, or cherry-pick.
//
// ----------------------------------
func handleNonTypingCtrlaKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		services.GitStateUniversalUtilsAbortService(m)
	}
	return m, nil
}

// ----------------------------------
//
//	Handle Ctrl+k key interaction.
//	Responsibility: Global "Skip" operation.
//	Used primarily to skip the current commit during an interactive rebase
//	or similar multi-step git states where skipping is a valid option.
//
// ----------------------------------
func handleNonTypingCtrlkKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		services.GitStateUniversalUtilsSkipService(m)
	}
	return m, nil
}

// ----------------------------------
//
//	Handle Ctrl+p key interaction.
//	Responsibility: Contextual "cherry-pick mode" management.
//	- Active Cherry-pick Popup: Switches between simple cherry-pick view and apply-confirm view.
//	- Log Panel / Branch Panel: Triggers opening the initial cherry-pick interface or option
//	  selection menu if commits are already queued.
//
// ----------------------------------
func handleNonTypingCtrlpKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.GitEditCherryPickPopUp:
			m.PopUpType = constant.GitCherryPickPopUp
			commitLogPopUp.InitGitCherryPickPopUp(m, m.CheckOutBranch)
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
		case constant.GitCherryPickApplyConfirmPopUp:
			m.PopUpType = constant.GitCherryPickPopUp
			commitLogPopUp.InitGitCherryPickPopUp(m, m.CheckOutBranch)
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
		}
	} else {
		if (m.CurrentSelectedComponent == constant.CommitLogOrRefLogComponentPanel || m.DetailPanelParentComponent == constant.CommitLogOrRefLogComponentPanel) &&
			len(m.CurrentRepoCommitLogInfoList.Items()) > 0 && m.CurrentCommitLogOrRefLogComponentShowing == constant.SHOW_COMMITLOG {
			if len(m.CherryPickedCommitInfo.CherryPickedCommitMap) < 1 {
				m.ShowPopUp.Store(true)
				m.PopUpType = constant.GitCherryPickPopUp
				m.IsTyping.Store(false)
				commitLogPopUp.InitGitCherryPickPopUp(m, m.CheckOutBranch)
			} else {
				m.ShowPopUp.Store(true)
				m.PopUpType = constant.GitCherryPickOptionSelectionPopUp
				m.IsTyping.Store(false)
				commitLogPopUp.InitGitCherryPickOptionSelectionPopUp(m)
			}
		} else if (m.CurrentSelectedComponent == constant.CommitLogOrRefLogComponentPanel || m.DetailPanelParentComponent == constant.CommitLogOrRefLogComponentPanel) &&
			len(m.CurrentRepoRefLogInfoList.Items()) > 0 && m.CurrentCommitLogOrRefLogComponentShowing == constant.SHOW_REFLOG {
			selectedReflog := m.CurrentRepoRefLogInfoList.SelectedItem()
			if selectedReflog != nil {
				parsedSelectedReflog := selectedReflog.(reflog.GitRefLogItem)
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				m.PopUpType = constant.GitCherryPickFromRefLogApplyConfirmationPopUp
				reflogPopUp.InitGitCherryPickFromRefLogApplyConfirmationPopUpModel(m, parsedSelectedReflog.Hash, parsedSelectedReflog.Head, parsedSelectedReflog.Action, parsedSelectedReflog.ActionInfo)
			}
		} else if (m.CurrentSelectedComponent == constant.LocalBranchOrTagOrRemoteComponentPanel || m.DetailPanelParentComponent == constant.LocalBranchOrTagOrRemoteComponentPanel) &&
			len(m.CurrentRepoTagInfoList.Items()) > 0 && m.CurrentLocalBranchOrTagOrRemoteComponentShowing == constant.SHOW_TAG {
			selectedTag := m.CurrentRepoTagInfoList.SelectedItem()
			if selectedTag == nil {
				return m, nil
			}
			if !m.GitOperations.GitRemote.CheckRemoteExist(false) {
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

// ----------------------------------
//
//	Handle Ctrl+r key interaction.
//	Responsibility: Contextual "revert" mechanism.
//	Specifically triggered in the Commit Log Panel to initiate a `git revert`
//	of the currently selected commit, handling both standard commits and merge commits
//	(by optionally prompting for the parent parent-line to revert against).
//
// ----------------------------------
func handleNonTypingCtrlrKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
				selectedCommitLog := m.CurrentRepoCommitLogInfoList.SelectedItem()
				if selectedCommitLog != nil {
					parsedCommitLog := selectedCommitLog.(commitlog.GitCommitLogItem)
					commitHashParentInfos := services.GetCommitHashParentInfoService(m, parsedCommitLog.Hash)
					if commitHashParentInfos != nil {
						if len(commitHashParentInfos) > 1 {
							commitLogPopUp.InitGitRevertParentOptionSelectionPopUp(m, parsedCommitLog.Hash, commitHashParentInfos)
							m.ShowPopUp.Store(true)
							m.IsTyping.Store(false)
							m.PopUpType = constant.GitRevertParentOptionSelectionPopUp
						} else {
							commitLogPopUp.InitGitRevertConfirmationPopUp(m, parsedCommitLog.Hash, 0)
							m.ShowPopUp.Store(true)
							m.IsTyping.Store(false)
							m.PopUpType = constant.GitRevertConfirmationPopUp
						}
					} else {
						commitLogPopUp.InitGitRevertConfirmationPopUp(m, parsedCommitLog.Hash, 0)
						m.ShowPopUp.Store(true)
						m.IsTyping.Store(false)
						m.PopUpType = constant.GitRevertConfirmationPopUp
					}
				}
			}
		}
	}
	return m, nil
}
