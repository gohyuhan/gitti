package branch

import (
	"slices"

	"charm.land/bubbles/v2/list"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	Handle the async branch-switch result event. Clears IsProcessing, sets
//	ProcessSuccess or HasError accordingly, and loads the command output into
//	the switch-branch output viewport.
//
// ------------------------------------
func UpdateSwitchBranchResultEvent(m *types.GittiModel, updateData types.GitSwitchBranchResultEventDataStructure) {
	popUp, ok := m.PopUpModel.(*SwitchBranchOutputPopUpModel)
	if ok {
		popUp.IsProcessing.Store(false)
		if updateData.Success {
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(true)
		} else {
			popUp.HasError.Store(true)
			popUp.ProcessSuccess.Store(false)
		}
		popUp.SwitchBranchOutputViewport.SetContentLines(updateData.Result)
		popUp.SwitchBranchOutputViewport.PageDown()
	}
}

// ------------------------------------
//
//	Handle the async branch-deletion result event. Clears IsProcessing, sets
//	ProcessSuccess or HasError accordingly, and loads the command output into
//	the branch-delete output viewport.
//
// ------------------------------------
func UpdateDeleteBranchResultEvent(m *types.GittiModel, updateData types.GitDeleteBranchResultEventDataStructure) {
	popUp, ok := m.PopUpModel.(*GitDeleteBranchOutputPopUpModel)
	if ok {
		popUp.IsProcessing.Store(false)
		if updateData.Success {
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(true)
		} else {
			popUp.HasError.Store(true)
			popUp.ProcessSuccess.Store(false)
		}
		popUp.BranchDeleteOutputViewport.SetContentLines(updateData.Result)
		popUp.BranchDeleteOutputViewport.PageDown()
	}
}

// ------------------------------------
//
//	Handle the async create-branch-from-remote result event. Clears IsProcessing,
//	sets ProcessSuccess or HasError accordingly, and loads the command output into
//	the create-branch-based-on-remote output viewport.
//
// ------------------------------------
func UpdateCreateNewBranchBasedOnRemoteResultEvent(m *types.GittiModel, updateData types.GitCreateNewBranchBasedOnRemoteResultEventDataStructure) {
	popUp, ok := m.PopUpModel.(*CreateBranchBasedOnRemoteOutputPopUpModel)
	if ok {
		popUp.IsProcessing.Store(false)
		if updateData.Success {
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(true)
		} else {
			popUp.HasError.Store(true)
			popUp.ProcessSuccess.Store(false)
		}
		popUp.CreateBranchBasedOnRemoteOutputViewport.SetContentLines(updateData.Result)
		popUp.CreateBranchBasedOnRemoteOutputViewport.PageDown()
	}
}

// ------------------------------------
//
//	Handle the create-branch-from-remote invalid event by dismissing the popup:
//	clears the typing flag, hides the popup, resets the popup type and model.
//
// ------------------------------------
func UpdateCreateNewBranchBasedOnRemoteInvalidEvent(m *types.GittiModel, _ types.GitCreateNewBranchBasedOnRemoteInvalidEventDataStructure) {
	m.IsTyping.Store(false)
	m.ShowPopUp.Store(false)
	m.PopUpType = constant.NoPopUp
	m.PopUpModel = nil
}

// ------------------------------------
//
//	Rebuild both the available-branch and selected-branch lists after a select
//	or unselect action. Moving a branch to/from the selected list resets both
//	list.Model instances while preserving the nearest valid cursor positions.
//
// ------------------------------------
func UpdateChooseBranchOptionForMergePopUpModel(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*ChooseBranchOptionForMergePopUpModel)
	if ok {
		branchOptionItems := []list.Item{}
		selectedBranchOptionItems := []list.Item{}
		width := (min(constant.MaxChooseBranchOptionForMergePopUpWidth, int(float64(m.Width)*0.8)) - 4)
		branchOptionListCurrentIndex := popUp.BranchOptionList.Index()
		selectedBranchOptionListCurrentIndex := popUp.SelectedBranchList.Index()

		// a string of branch that have been selected
		var tmpSelectedBranchStringList []string

		if popUp.BranchOptionSectionSelected.Load() {
			// if user is currently selecting branch for merging, we include the selected one to the selected branch list
			selectedBranch := popUp.BranchOptionList.SelectedItem()
			if selectedBranch == nil {
				return
			}

			var newSelectedBranchStringList []string
			for _, item := range popUp.SelectedBranchList.Items() {
				branchItem := item.(GitMergeBranchOptionItem)
				newSelectedBranchStringList = append(newSelectedBranchStringList, branchItem.BranchName)
			}
			newSelectedBranchStringList = append(newSelectedBranchStringList, selectedBranch.(GitMergeBranchOptionItem).BranchName)
			tmpSelectedBranchStringList = newSelectedBranchStringList

		} else if popUp.SelectedBranchSectionSelected.Load() {
			// if user is currently unselecting branch for merging, we exclude the selected one from the selected branch list
			unselectedBranch := popUp.SelectedBranchList.SelectedItem()
			if unselectedBranch == nil {
				return
			}
			var newSelectedBranchStringList []string
			unselected := unselectedBranch.(GitMergeBranchOptionItem)
			for _, item := range popUp.SelectedBranchList.Items() {
				branchItem := item.(GitMergeBranchOptionItem)
				if branchItem.BranchName != unselected.BranchName {
					newSelectedBranchStringList = append(newSelectedBranchStringList, branchItem.BranchName)
				}
			}
			tmpSelectedBranchStringList = newSelectedBranchStringList
		} else {
			return
		}

		// reinit the current branch option list and selected branch list

		for _, branch := range m.GitOperations.GitBranch.AllBranches() {
			if branch.IsCheckedOut {
				continue
			}
			if slices.Contains(tmpSelectedBranchStringList, branch.BranchName) {
				selectedBranchOptionItems = append(selectedBranchOptionItems, GitMergeBranchOptionItem{BranchName: branch.BranchName})
			} else {
				branchOptionItems = append(branchOptionItems, GitMergeBranchOptionItem{BranchName: branch.BranchName})
			}
		}

		// for selecting branch for git merge
		cBOFMBOL := list.New(branchOptionItems, GitMergeBranchOptionItemDelegate{}, width, constant.PopUpChooseBranchOptionForMergeBranchOptionHeight)
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
		cBOFMSBOL := list.New(selectedBranchOptionItems, GitMergeBranchOptionItemDelegate{}, width, constant.PopUpChooseBranchOptionForMergeSelectedBranchOptionHeight)
		cBOFMSBOL.SetShowPagination(false)
		cBOFMSBOL.SetShowStatusBar(false)
		cBOFMSBOL.SetFilteringEnabled(false)
		cBOFMSBOL.SetShowTitle(false)

		// Custom Help Model for Count Display
		cBOFMSBOL.SetShowHelp(true)
		cBOFMSBOL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
		cBOFMSBOL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
		cBOFMSBOL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &cBOFMSBOL, constant.MaxChooseBranchOptionForMergePopUpWidth)

		// reinit the selection index
		maxBranchIdx := max(0, len(branchOptionItems)-1)
		cBOFMBOL.Select(min(branchOptionListCurrentIndex, maxBranchIdx))

		maxSelectedBranchIdx := max(0, len(selectedBranchOptionItems)-1)
		cBOFMSBOL.Select(min(selectedBranchOptionListCurrentIndex, maxSelectedBranchIdx))

		// assign the newly constructed lists back to the popUp model
		popUp.BranchOptionList = cBOFMBOL
		popUp.SelectedBranchList = cBOFMSBOL
	}
}

// ------------------------------------
//
//	Handle the async branch-merge result event. Clears IsProcessing and sets
//	ProcessSuccess or HasError based on success, then writes all output lines
//	to the merge output viewport.
//
// ------------------------------------
func UpdateMergeViewport(m *types.GittiModel, updateData types.MergeResultEventDataStructure) {
	success := updateData.Success
	outputResult := updateData.Result
	popUp, ok := m.PopUpModel.(*BranchMergeOutputPopUpModel)
	if ok && !popUp.IsCancelled.Load() {
		if success && !popUp.IsProcessing.Load() {
			popUp.ProcessSuccess.Store(true)
			popUp.HasError.Store(false)
		} else if !success && !popUp.IsProcessing.Load() {
			popUp.ProcessSuccess.Store(false)
			popUp.HasError.Store(true)
		}
		popUp.IsProcessing.Store(false) // update the processing status
		popUp.BranchMergeOutputViewport.SetContentLines(outputResult)
	}
}
