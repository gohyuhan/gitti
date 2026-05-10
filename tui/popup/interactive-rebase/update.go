package interactiverebase

import (
	"fmt"
	"strings"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Updates the fixup/squash preview viewport from current selected commits and validation state
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

func UpdateInteractiveRebaseFetchedCommitInfoList(m *types.GittiModel, UpdateData types.InteractiveRebaseFetchCommitInfoListEventDataInterface) {
	switch UpdateData.PopUpModel {
	case constant.InteractiveRebaseFixupSquashSelectionPopUp:
		popUp, ok := m.PopUpModel.(*InteractiveRebaseFixupSquashSelectionPopUpModel)
		if ok {
			popUp.OriginalRetrievedCommitList = UpdateData.CommitInfos
			popUp.CommitList.SetItems(UpdateData.ListItems)
		}
	}
}

func UpdateInteractiveRebaseFixupSquashResultEvent(m *types.GittiModel, updateData types.InteractiveRebaseFixupSquashResultEventDataInterface) {
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
