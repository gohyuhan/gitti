package nontyping

import (
	"context"
	"slices"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/component/branch"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/component/remote"
	"github.com/gohyuhan/gitti/tui/component/tag"
	"github.com/gohyuhan/gitti/tui/component/worktree"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/layout"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	commitLogPopUp "github.com/gohyuhan/gitti/tui/popup/commitlog"
	discardPopUp "github.com/gohyuhan/gitti/tui/popup/discard"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	pullPopUp "github.com/gohyuhan/gitti/tui/popup/pull"
	pushPopUp "github.com/gohyuhan/gitti/tui/popup/push"
	rebasePopUp "github.com/gohyuhan/gitti/tui/popup/rebase"
	reflogPopUp "github.com/gohyuhan/gitti/tui/popup/reflog"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	resolvePopUp "github.com/gohyuhan/gitti/tui/popup/resolve"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	worktreePopUp "github.com/gohyuhan/gitti/tui/popup/worktree"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	Handle Enter key interaction.
//	Responsibility: Core confirmation and drill-down action. Depends heavily on context:
//	- Active Popup (e.g., choosing remote, selecting branch type, confirming discard): Executes the chosen workflow or triggers git commands (push, pull, reset, discard, tag operations, etc.).
//	- Interactive Rebase Option Popup: Confirms the selected operation type and navigates to the next popup (fixup/squash, reword, or drop selection).
//	- Interactive Rebase Selection Popups: Validates selection and transitions to the output popup; forks into signing or async service path.
//	- Component Panels (when no popup active): Typically drills down into a "detail view" for the selected item (e.g., viewing diffs for a file, showing commit details), or triggers context menus like switching branches.
//
// ------------------------------------
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
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				selectedLocalBranchItem := m.CurrentRepoBranchesInfoList.SelectedItem()
				if selectedLocalBranchItem != nil {
					currentSelectedLocalBranch := selectedLocalBranchItem.(branch.GitBranchItem)
					// only proceed if the local branch selected is not current checkedout branch
					// we can't switch from current checkout branch to current checkout branch, do we
					if !currentSelectedLocalBranch.IsCheckedOut {
						m.PopUpType = constant.ChooseSwitchBranchTypePopUp
						m.IsTyping.Store(false)
						m.ShowPopUp.Store(true)
						branchPopUp.InitChooseSwitchBranchTypePopUpModel(m, currentSelectedLocalBranch.BranchName)
					}
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
			case constant.SHOW_WORKTREE:
				currentSelectedWorktree := m.CurrentRepoWorktreeInfoList.SelectedItem()
				if currentSelectedWorktree != nil {
					selectedWorktree := currentSelectedWorktree.(worktree.GitWorktreeItem)
					// can't switch into the worktree we're already in, and a prunable
					// (stale/missing) worktree is not a valid switch target
					if selectedWorktree.IsInCurrentWorktree || selectedWorktree.IsPrunable {
						return m, nil
					}
					services.SwitchWorktreeService(m, selectedWorktree.WorktreePath)
					layout.LeftPanelDynamicResize(m)
					services.FetchDetailComponentPanelInfoService(m, true)
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
				if remote != nil {
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
						branchPopUp.InitCreateBranchBasedOnRemotePopUpModel(m, remoteName)
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
			}

		case constant.ChoosePushTypePopUp:
			popUp, ok := m.PopUpModel.(*pushPopUp.ChoosePushTypePopUpModel)
			if ok {
				selectedOption := popUp.PushOptionList.SelectedItem()
				if selectedOption != nil {
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
			}

		case constant.ChooseNewBranchTypePopUp:
			popUp, ok := m.PopUpModel.(*branchPopUp.ChooseNewBranchTypeOptionPopUpModel)
			if ok {
				selectedOption := popUp.NewBranchTypeOptionList.SelectedItem()
				if selectedOption == nil {
					return m, nil
				}
				newBranchType := selectedOption.(branchPopUp.GitNewBranchTypeOptionItem).NewBranchType
				switch newBranchType {
				case git.NEWBRANCHBASEDONREMOTEUSERINPUT:
					if !m.GitOperations.GitRemote.CheckRemoteExist(false) {
						// if no remote found, we add one
						showAddRemotePromptPopUp(m)
					} else {
						m.ShowPopUp.Store(true)
						remotes := m.GitOperations.GitRemote.FetchRemote()
						if len(remotes) == 1 {
							m.IsTyping.Store(true)
							m.PopUpType = constant.CreateBranchBasedOnRemotePopUp
							// only one remote found so, we will default to that remote
							branchPopUp.InitCreateBranchBasedOnRemotePopUpModel(m, remotes[0].Name)
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
				selectedSwitchItem := popUp.SwitchTypeOptionList.SelectedItem()
				if selectedSwitchItem == nil {
					return m, nil
				}
				m.PopUpType = constant.SwitchBranchOutputPopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				selectedOption := selectedSwitchItem.(branchPopUp.GitSwitchBranchTypeOptionItem)
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
				selectedPullItem := popUp.PullTypeOptionList.SelectedItem()
				if selectedPullItem == nil {
					return m, nil
				}
				selectedOption := selectedPullItem.(pullPopUp.GitPullTypeOptionItem)

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
				selectedDiscardItem := popUp.DiscardTypeOptionList.SelectedItem()
				if selectedDiscardItem == nil {
					return m, nil
				}
				m.PopUpType = constant.GitDiscardConfirmPromptPopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				selectedOption := selectedDiscardItem.(discardPopUp.GitDiscardTypeOptionItem)
				discardPopUp.InitGitDiscardConfirmPromptPopUpModel(m, popUp.FilePathName, selectedOption.DiscardType)
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
				selectedResolveItem := popUp.ResolveConflictOptionList.SelectedItem()
				if selectedResolveItem == nil {
					return m, nil
				}
				selectedResolveType := selectedResolveItem.(resolvePopUp.GitResolveConflictOptionItem)
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
						commitLogPopUp.InitGitCherryPickPopUpModel(m, m.CheckOutBranch)
					case constant.EDITCHERRYPICK:
						m.PopUpType = constant.GitEditCherryPickPopUp
						commitLogPopUp.InitGitEditCherryPickPopUpModel(m, 0)
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
				copiedCherryPickedHashmap := make(map[string]git.CherryPickedCommitLog)
				for k, v := range m.CherryPickedCommitInfo.CherryPickedCommitMap {
					copiedCherryPickedHashmap[k] = v
				}
				services.GitCherryPickService(m, copiedCherryPickedHashmap)
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
							showAddRemotePromptPopUp(m)
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
					commitLogPopUp.InitGitRevertConfirmationPopUpModel(m, popUp.CommitHash, parsedSelectedParent.ParentOrder)
				} else {
					commitLogPopUp.InitGitRevertConfirmationPopUpModel(m, popUp.CommitHash, 1)
				}
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				m.PopUpType = constant.GitRevertConfirmationPopUp
				return m, nil
			}
		case constant.GitRevertConfirmationPopUp:
			popUp, ok := m.PopUpModel.(*commitLogPopUp.GitRevertConfirmationPopUpModel)
			if ok {
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
						branchPopUp.InitCreateBranchBasedOnRemoteOutputPopUpModel(m)
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
				if m.GitCommitRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
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
		case constant.InteractiveRebaseOptionPopUp:
			popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseOptionPopUpModel)
			if ok {
				selectedInteractiveRebaseOption := popUp.InteractiveRebaseOptionList.SelectedItem()
				if selectedInteractiveRebaseOption != nil {
					parsedSelected := selectedInteractiveRebaseOption.(interactiverebasePopUp.InteractiveRebaseOptionItem)
					switch parsedSelected.InteractiveRebaseType {
					case git.FIXUPSQUASH:
						m.ShowPopUp.Store(true)
						m.IsTyping.Store(false)
						m.PopUpType = constant.InteractiveRebaseFixupSquashSelectionPopUp
						interactiverebasePopUp.InitInteractiveRebaseFixupSquashSelectionPopUpModel(m)
					case git.REWORD:
						m.ShowPopUp.Store(true)
						m.IsTyping.Store(false)
						m.PopUpType = constant.InteractiveRebaseRewordSelectionPopUp
						interactiverebasePopUp.InitInteractiveRebaseRewordSelectionPopUpModel(m)
					case git.DROP:
						m.ShowPopUp.Store(true)
						m.IsTyping.Store(false)
						m.PopUpType = constant.InteractiveRebaseDropSelectionPopUp
						interactiverebasePopUp.InitInteractiveRebaseDropSelectionPopUpModel(m)
					}
				}
			}
		case constant.InteractiveRebaseFixupSquashSelectionPopUp:
			popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashSelectionPopUpModel)
			if ok {
				if popUp.SelectionError == nil {
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(true)
					m.PopUpType = constant.InteractiveRebaseFixupSquashCommitPopUp
					interactiverebasePopUp.InitInteractiveRebaseFixupSquashCommitPopUp(m, popUp.OriginalRetrievedCommitList, popUp.SortedSelectedCommits)
				}
			}

		case constant.InteractiveRebaseRewordSelectionPopUp:
			popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseRewordSelectionPopUpModel)
			if ok {
				selectedItem := popUp.CommitList.SelectedItem()
				// return early is none
				if selectedItem == nil {
					return m, nil
				}
				parsedSelectedItem := selectedItem.(interactiverebasePopUp.InteractiveRebaseRewordSelectionItem)
				selectedCommit := git.CommitInfo(parsedSelectedItem)
				interactiverebasePopUp.InteractiveRebaseRewordSelectionValidation(m, selectedCommit)

				if popUp.SelectionError == nil {
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(true)
					m.PopUpType = constant.InteractiveRebaseRewordCommitPopUp
					interactiverebasePopUp.InitInteractiveRebaseRewordCommitPopUp(m, popUp.OriginalRetrievedCommitList, selectedCommit)
				}
			}
		case constant.InteractiveRebaseDropSelectionPopUp:
			popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseDropSelectionPopUpModel)
			if ok {
				ogRetrievedCommitsList := popUp.OriginalRetrievedCommitList
				sortedSelectedCommits := popUp.SortedSelectedCommits
				if len(ogRetrievedCommitsList) > 1 && len(sortedSelectedCommits) > 0 && popUp.SelectionError == nil {
					// Switch to output popup before starting execution so errors/progress are visible immediately.
					interactiverebasePopUp.InitInteractiveRebaseDropOutputPopUpModel(m)
					popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseDropOutputPopUpModel)
					if !ok {
						return m, nil
					}
					m.PopUpType = constant.InteractiveRebaseDropOutputPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					if m.GitCommitRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
						// Signing path returns prepared exec command; tea.ExecProcess handles terminal suspension.
						executor, cleanupCallbackFunc, dropErr := m.GitOperations.GitInteractiveRebase.GitInteractiveRebaseDropWithSigning(context.TODO(), ogRetrievedCommitsList, sortedSelectedCommits)
						if dropErr != nil {
							popUp.HasError.Store(true)
							popUp.DropOutputViewport.SetContent(dropErr.Error())
							return m, nil
						}
						return utils.SuspendGittiUIForGitOperationRequireSigningWithExecAndCleanUp(m, executor, cleanupCallbackFunc, logging.INTERACTIVE_REBASE_DROP)
					} else {
						popUp.IsProcessing.Store(true)
						// Non-signing path runs async service with cancellable context.
						services.InteractiveRebaseDropService(m, ogRetrievedCommitsList, sortedSelectedCommits)

						// Start spinner ticking
						return m, popUp.Spinner.Tick
					}
				}
			}
		case constant.WorktreeRemoveWorktreeConfirmationPopUp:
			popUp, ok := m.PopUpModel.(*worktreePopUp.WorktreeRemoveWorktreeConfirmationPopUpModel)
			if ok {
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				services.RemoveWorktreeService(m, popUp.WorktreePath)
			}
		}
	}
	return m, nil
}
