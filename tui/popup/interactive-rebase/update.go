package interactiverebase

import (
	"fmt"
	"strings"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// *************************************************************************************
//                        INTERACTIVE REBASE - FIXUP / SQUASH
// *************************************************************************************

// ------------------------------------
//
//	Rebuild the fixup/squash preview viewport content from the current selection.
//	If a validation error is present, shows the error message; otherwise renders
//	an oldest-to-newest chain of selected commit hashes with arrows, ending with
//	the newest commit in the original list if it was not selected.
//
// ------------------------------------
func UpdateInteractiveRebaseFixupSquashViewport(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseFixupSquashSelectionPopUpModel)
	if ok {
		var content strings.Builder

		sortedSelectedCommitArray := popUp.SortedSelectedCommits

		if popUp.SelectionError != nil {
			errorMsg := style.NewStyle.Foreground(style.ColorError).Render(popUp.SelectionError.Error())
			content.WriteString(errorMsg)
			content.WriteRune('\n')
		} else {
			// Render oldest->newest selected chain; last selected rendered after loop.
			for _, commit := range sortedSelectedCommitArray[:len(sortedSelectedCommitArray)-1] {
				commitHashAuthorLine := style.NewStyle.Foreground(style.ColorPurpleSoft).Render(fmt.Sprintf(" %s | %s ", commit.Hash[:7], commit.Author))
				arrowLine := style.NewStyle.Foreground(style.ColorYellowSoft).Render("          ↑")
				content.WriteString(commitHashAuthorLine)
				content.WriteRune('\n')
				content.WriteString(arrowLine)
				content.WriteRune('\n')
			}

			lastSelectedCommit := sortedSelectedCommitArray[len(sortedSelectedCommitArray)-1]
			commitHashAuthorLine := style.NewStyle.Foreground(style.ColorPurpleSoft).Render(fmt.Sprintf(" %s | %s ", lastSelectedCommit.Hash[:7], lastSelectedCommit.Author))
			content.WriteString(commitHashAuthorLine)
			content.WriteRune('\n')

			// Original list order is newest->oldest.
			newestCommit := popUp.OriginalRetrievedCommitList[0]

			if lastSelectedCommit.Hash != newestCommit.Hash {
				content.WriteString("          .")
				content.WriteRune('\n')
				content.WriteString("          .")
				content.WriteRune('\n')
				content.WriteString("          .")
				content.WriteRune('\n')
				commitHashAuthorLine := style.NewStyle.Foreground(style.ColorPurpleSoft).Render(fmt.Sprintf(" %s | %s ", newestCommit.Hash[:7], newestCommit.Author))
				content.WriteString(commitHashAuthorLine)
				content.WriteRune('\n')
			}
		}

		popUp.CommitFixupSquashViewport.SetContent(content.String())

	}
}

// ------------------------------------
//
//	Handle the async commit-info fetch event. Stores the retrieved commit list
//	on the matching popup model and populates the commit list with the received
//	items. Currently handles the fixup/squash selection popup type.
//
// ------------------------------------
func UpdateInteractiveRebaseFixupSquashFetchedCommitInfoList(m *types.GittiModel, updateData types.InteractiveRebaseFetchCommitInfoListEventDataStructure) {
	switch updateData.PopUpModel {
	case constant.InteractiveRebaseFixupSquashSelectionPopUp:
		popUp, ok := m.PopUpModel.(*InteractiveRebaseFixupSquashSelectionPopUpModel)
		if ok {
			popUp.OriginalRetrievedCommitList = updateData.CommitInfos
			popUp.CommitList.SetItems(updateData.ListItems)
		}
	}
}

// ------------------------------------
//
//	Handle the async fixup/squash rebase result event. Clears IsProcessing and
//	sets ProcessSuccess on success or HasError on failure. Loads the output lines
//	into the viewport. No-ops if the popup was cancelled.
//
// ------------------------------------
func UpdateInteractiveRebaseFixupSquashResultEvent(m *types.GittiModel, updateData types.InteractiveRebaseFixupSquashResultEventDataStructure) {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseFixupSquashOutputPopUpModel)
	if ok && !popUp.IsCancelled.Load() {
		if updateData.Success && !popUp.IsProcessing.Load() {
			popUp.ProcessSuccess.Store(true)
			popUp.HasError.Store(false)
		} else if !updateData.Success && !popUp.IsProcessing.Load() {
			popUp.ProcessSuccess.Store(false)
			popUp.HasError.Store(true)
		}
		popUp.FixupSquashOutputViewport.SetContentLines(updateData.Result)
		popUp.IsProcessing.Store(false)
	}
}

// *************************************************************************************
//
//	INTERACTIVE REBASE - REWORD
//
// *************************************************************************************

// ------------------------------------
//
//	Handle the async commit-info fetch event for reword. Stores the retrieved
//	commit list on the matching popup model and populates the commit list with
//	the received items.
//
// ------------------------------------
func UpdateInteractiveRebaseRewordFetchedCommitInfoList(m *types.GittiModel, updateData types.InteractiveRebaseFetchCommitInfoListEventDataStructure) {
	switch updateData.PopUpModel {
	case constant.InteractiveRebaseRewordSelectionPopUp:
		popUp, ok := m.PopUpModel.(*InteractiveRebaseRewordSelectionPopUpModel)
		if ok {
			popUp.OriginalRetrievedCommitList = updateData.CommitInfos
			popUp.CommitList.SetItems(updateData.ListItems)
		}
	}
}

// ------------------------------------
//
//	Handle the async reword rebase result event. Clears IsProcessing and
//	sets ProcessSuccess on success or HasError on failure. Loads the output
//	lines into the viewport. No-ops if the popup was cancelled.
//
// ------------------------------------
func UpdateInteractiveRebaseRewordResultEvent(m *types.GittiModel, updateData types.InteractiveRebaseRewordResultEventDataStructure) {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseRewordOutputPopUpModel)
	if ok && !popUp.IsCancelled.Load() {
		if updateData.Success && !popUp.IsProcessing.Load() {
			popUp.ProcessSuccess.Store(true)
			popUp.HasError.Store(false)
		} else if !updateData.Success && !popUp.IsProcessing.Load() {
			popUp.ProcessSuccess.Store(false)
			popUp.HasError.Store(true)
		}
		popUp.RewordOutputViewport.SetContentLines(updateData.Result)
		popUp.IsProcessing.Store(false)
	}
}

// *************************************************************************************
//                           INTERACTIVE REBASE - DROP
// *************************************************************************************
