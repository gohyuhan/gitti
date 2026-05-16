package interactiverebase

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Validate and sort the current fixup/squash commit selection. Requires at
//	least two commits selected and the oldest (base) commit must not be a merge
//	commit. Sorts selected commits oldest-first by CommitOrder and stores the
//	result in SortedSelectedCommits; sets SelectionError on any validation failure.
//
// ------------------------------------
func InteractiveRebaseFixupSquashSelectionValidationAndSort(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseFixupSquashSelectionPopUpModel)
	if ok {
		popUp.SelectionError = nil
		var sortedSelectedCommitArray []git.CommitInfo

		for _, commit := range popUp.SelectedCommitHashMap {
			sortedSelectedCommitArray = append(sortedSelectedCommitArray, commit)
		}

		slices.SortFunc(sortedSelectedCommitArray, func(a, b git.CommitInfo) int {
			return cmp.Compare(b.CommitOrder, a.CommitOrder) // largest CommitOrder first = oldest to latest
		})

		popUp.SortedSelectedCommits = sortedSelectedCommitArray

		if len(sortedSelectedCommitArray) < 2 {
			popUp.SelectionError = fmt.Errorf("%s", i18n.LANGUAGEMAPPING.InteractiveRebaseFixupMustHaveAtLeastTwoSelectedError)
			return
		}

		if len(sortedSelectedCommitArray[0].Parent) > 1 {
			popUp.SelectionError = fmt.Errorf("%s", i18n.LANGUAGEMAPPING.InteractiveRebaseFixupBaseCommitCannotBeAMergeCommit)
		}
	}
}

// ------------------------------------
//
//	Validate the reword commit selection. Sets SelectionError if the selected
//	commit is a merge commit, which cannot be rewarded. Clears the error otherwise.
//
// ------------------------------------
func InteractiveRebaseRewordSelectionValidation(m *types.GittiModel, selectedCommit git.CommitInfo) {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseRewordSelectionPopUpModel)
	if ok {
		popUp.SelectionError = nil

		if len(selectedCommit.Parent) > 1 {
			popUp.SelectionError = fmt.Errorf("%s", i18n.LANGUAGEMAPPING.InteractiveRebaseRewordCommitCannotBeAMergeCommit)
		}
	}
}
