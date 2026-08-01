package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/component/stash"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitLogPopUp "github.com/gohyuhan/gitti/tui/popup/commitlog"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	Handle Space key interaction.
//	Responsibility: Contextual "toggle" or "apply" action.
//	- Modified / Detail Panels: Stages or unstages the selected file or specific hunk/line.
//	- Stash Panel: Triggers applying a stash (without dropping it).
//	- Cherry Pick Popup: Selects/toggles the currently highlighted commit to be added to the cherry-pick queue.
//	- Merge Popup: Toggles the currently highlighted branch between the available and selected branch lists.
//	- Interactive Rebase Fixup/Squash Popup: Toggles the highlighted commit's inclusion in the fixup/squash target set.
//	- Interactive Rebase Drop Popup: Toggles the highlighted commit's inclusion in the drop target set.
//
// ------------------------------------
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

		case constant.InteractiveRebaseFixupSquashSelectionPopUp:
			popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashSelectionPopUpModel)
			if ok {
				if popUp.IsCommitListSelected {
					selectedItem := popUp.CommitList.SelectedItem()
					if selectedItem != nil {
						// check if the selected item is pick for selection or not
						// if not, include it for selection, else unselect it
						parsedSelectedItem := selectedItem.(interactiverebasePopUp.InteractiveRebaseFixupSquashSelectionItem)
						_, exist := popUp.SelectedCommitHashMap[parsedSelectedItem.Hash]
						if exist {
							delete(popUp.SelectedCommitHashMap, parsedSelectedItem.Hash)
						} else {
							popUp.SelectedCommitHashMap[parsedSelectedItem.Hash] = git.CommitInfo(parsedSelectedItem)
						}
						interactiverebasePopUp.InteractiveRebaseFixupSquashSelectionValidationAndSort(m)
						interactiverebasePopUp.UpdateInteractiveRebaseFixupSquashViewport(m)
					}
				}
			}
		case constant.InteractiveRebaseDropSelectionPopUp:
			popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseDropSelectionPopUpModel)
			if ok {
				selectedItem := popUp.CommitList.SelectedItem()
				if selectedItem != nil {
					// check if the selected item is pick for selection or not
					// if not, include it for selection, else unselect it
					parsedSelectedItem := selectedItem.(interactiverebasePopUp.InteractiveRebaseDropSelectionItem)
					_, exist := popUp.SelectedCommitHashMap[parsedSelectedItem.Hash]
					if exist {
						delete(popUp.SelectedCommitHashMap, parsedSelectedItem.Hash)
					} else {
						popUp.SelectedCommitHashMap[parsedSelectedItem.Hash] = git.CommitInfo(parsedSelectedItem)
					}
					interactiverebasePopUp.InteractiveRebaseDropSelectionValidationAndSort(m)
				}
			}
		}
	}
	return m, nil
}
