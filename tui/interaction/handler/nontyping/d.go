package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/component/branch"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/component/remote"
	"github.com/gohyuhan/gitti/tui/component/stash"
	"github.com/gohyuhan/gitti/tui/component/tag"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitLogPopUp "github.com/gohyuhan/gitti/tui/popup/commitlog"
	discardPopUp "github.com/gohyuhan/gitti/tui/popup/discard"
	filesPopUp "github.com/gohyuhan/gitti/tui/popup/files"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	Handle 'd' key interaction.
//	Responsibility: Multi-contextual "discard" or "delete" operation. Depending on the focused panel, it handles:
//	- Deleting local branches, remote references, or tags (prompting confirmations).
//	- Dropping git stashes.
//	- Discarding changes in modified files (with complex checking for staged, unstaged, untracked, and renamed states).
//	- Discarding specific hunk/line changes when in detail line-editing mode.
//	Also handles popup dismissals for specific cherry-pick list items.
//
// ------------------------------------
func handleNonTypingdKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
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
					discardPopUp.InitGitDiscardTypeOptionPopUpModel(m, currentSelectedFile.FilePathname, true, false)
				} else if currentSelectedFile.IndexState == "R" && currentSelectedFile.WorkTree != " " {
					// a staged rename with unstaged modification
					m.PopUpType = constant.GitDiscardTypeOptionPopUp
					discardPopUp.InitGitDiscardTypeOptionPopUpModel(m, currentSelectedFile.FilePathname, false, true)
				} else if currentSelectedFile.IndexState == "?" && currentSelectedFile.WorkTree == "?" {
					// newly added untracked file
					m.PopUpType = constant.GitDiscardConfirmPromptPopUp
					discardPopUp.InitGitDiscardConfirmPromptPopUpModel(m, currentSelectedFile.FilePathname, git.DISCARDUNTRACKED)
				} else if currentSelectedFile.IndexState != "A" && currentSelectedFile.IndexState != "C" && currentSelectedFile.IndexState != "R" && currentSelectedFile.IndexState != "?" && currentSelectedFile.IndexState != " " && currentSelectedFile.WorkTree != " " {
					// tracked file with both staged and unstaged modification (beside A, C and  )
					m.PopUpType = constant.GitDiscardTypeOptionPopUp
					discardPopUp.InitGitDiscardTypeOptionPopUpModel(m, currentSelectedFile.FilePathname, false, false)
				} else if (currentSelectedFile.IndexState == "A" && currentSelectedFile.WorkTree == " ") || (currentSelectedFile.IndexState == "C" && currentSelectedFile.WorkTree == " ") {
					// newly added tracked / copied file
					m.PopUpType = constant.GitDiscardConfirmPromptPopUp
					discardPopUp.InitGitDiscardConfirmPromptPopUpModel(m, currentSelectedFile.FilePathname, git.DISCARDNEWLYADDEDORCOPIED)
				} else if currentSelectedFile.IndexState == "R" && currentSelectedFile.WorkTree == " " {
					// a staged rename
					m.PopUpType = constant.GitDiscardConfirmPromptPopUp
					discardPopUp.InitGitDiscardConfirmPromptPopUpModel(m, currentSelectedFile.FilePathname, git.DISCARDANDREVERTRENAME)
				} else {
					// tracked file with only unstaged modification
					m.PopUpType = constant.GitDiscardConfirmPromptPopUp
					discardPopUp.InitGitDiscardConfirmPromptPopUpModel(m, currentSelectedFile.FilePathname, git.DISCARDWHOLE)
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
						filesPopUp.InitGitDiscardFileLineChangeConfirmPopUpModel(m)
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
						commitLogPopUp.InitGitEditCherryPickPopUpModel(m, popUp.CherryPickedCommitLog.Index())
					}
				}
			}
		}
	}
	return m, nil
}
