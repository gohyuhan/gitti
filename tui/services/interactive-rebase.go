package services

import (
	"context"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/constant"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	"github.com/gohyuhan/gitti/tui/types"
)

// *************************************************************************************
//
//	INTERACTIVE REBASE - FIXUP / SQUASH
//
// *************************************************************************************
// ------------------------------------
//
//	Runs fixup/squash interactive rebase asynchronously, streams command output to popup viewport,
//	and updates popup processing/success/error state
//
// ------------------------------------
func InteractiveRebaseFixupSquashService(m *types.GittiModel, originalRetrievedGitCommitInfo []git.CommitInfo, sortedSelectedCommits []git.CommitInfo, fixupSquashCommitMessage string, fixupSquashCommitDescription string) {

	popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashOutputPopUpModel)
	if ok {
		ctx, cancel := context.WithCancel(context.Background())
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsProcessing.Store(true)
		popUp.IsCancelled.Store(false)
		popUp.CancelFunc = cancel

		go func(ctx context.Context) {
			defer cancel()

			fixupSquashResult, fixupSquashErr := m.GitOperations.GitInteractiveRebase.GitInteractiveRebaseFixupSquash(ctx, originalRetrievedGitCommitInfo, sortedSelectedCommits, fixupSquashCommitMessage, fixupSquashCommitDescription)
			data := types.InteractiveRebaseFixupSquashResultEventDataStructure{
				Result:  fixupSquashResult,
				Success: fixupSquashErr == nil,
			}
			m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
				Event: constant.INTERACTIVE_REBASE_FIXUP_SQUASH_RESULT_EVENT,
				Data:  data,
			}
		}(ctx)
	} else {
		return
	}

}

// ------------------------------------
//
//	Cancels an in-progress fixup/squash rebase, tears down the output popup state, and resets all atomic flags
//
// ------------------------------------
func InteractiveRebaseFixupSquashCancelService(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashOutputPopUpModel)
	if ok {
		// Set cancellation marker first, then cancel context.
		popUp.IsCancelled.Store(true)
		if popUp.CancelFunc != nil {
			popUp.CancelFunc()
		}
	}
	m.ShowPopUp.Store(false)
	m.IsTyping.Store(false)
	m.PopUpType = constant.NoPopUp
	if ok {
		popUp.FixupSquashOutputViewport.SetContent("")
		popUp.IsProcessing.Store(false)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
	}
}

// *************************************************************************************
//
//	INTERACTIVE REBASE - REWORD
//
// *************************************************************************************
func InteractiveRebaseRewordService(m *types.GittiModel, originalRetrievedGitCommitInfo []git.CommitInfo, selectedCommit git.CommitInfo, rewordCommitMessage string, rewordCommitDescription string) {

	popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseRewordOutputPopUpModel)
	if ok {
		ctx, cancel := context.WithCancel(context.Background())
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsProcessing.Store(true)
		popUp.IsCancelled.Store(false)
		popUp.CancelFunc = cancel

		go func(ctx context.Context) {
			defer cancel()

			rewordResult, rewordErr := m.GitOperations.GitInteractiveRebase.GitInteractiveRebaseReword(ctx, originalRetrievedGitCommitInfo, selectedCommit, rewordCommitMessage, rewordCommitDescription)
			data := types.InteractiveRebaseRewordResultEventDataStructure{
				Result:  rewordResult,
				Success: rewordErr == nil,
			}
			m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
				Event: constant.INTERACTIVE_REBASE_REWORD_RESULT_EVENT,
				Data:  data,
			}
		}(ctx)
	} else {
		return
	}

}

func InteractiveRebaseRewordCancelService(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseRewordOutputPopUpModel)
	if ok {
		// Set cancellation marker first, then cancel context.
		popUp.IsCancelled.Store(true)
		if popUp.CancelFunc != nil {
			popUp.CancelFunc()
		}
	}
	m.ShowPopUp.Store(false)
	m.IsTyping.Store(false)
	m.PopUpType = constant.NoPopUp
	if ok {
		popUp.RewordOutputViewport.SetContent("")
		popUp.IsProcessing.Store(false)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
	}
}
