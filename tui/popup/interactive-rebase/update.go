package interactiverebase

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ----------------------------------------------------------------------------------------------------------
//
// UpdateInteractiveRebaseFixupSquashViewport updates fixup/squash review viewport from selected commits.
//
// ----------------------------------------------------------------------------------------------------------
func UpdateInteractiveRebaseFixupSquashViewport(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseFixupSquashSelectionPopUpModel)
	if ok {
		var sortedSelectedCommitArray []git.CommitInfo

		for _, commit := range popUp.SelectedCommitHashMap {
			sortedSelectedCommitArray = append(sortedSelectedCommitArray, commit)
		}

		if len(sortedSelectedCommitArray) < 2 {
			popUp.CommitFixupSquashViewport.SetContent(i18n.LANGUAGEMAPPING.InteractiveRebaseFixupMustHaveAtLeastTwoSelectedError)
			return
		}

		slices.SortFunc(sortedSelectedCommitArray, func(a, b git.CommitInfo) int {
			return cmp.Compare(b.CommitOrder, a.CommitOrder) // largest CommitOrder first = oldest to latest
		})

		var content strings.Builder

		// loop but exclude the last item
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

		newestCommit := popUp.OriginalRetrievedCommitList[0] // in original retrieved commit list the order was from newest commit to oldest

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

		popUp.CommitFixupSquashViewport.SetContent(content.String())

	}
}
