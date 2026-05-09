package services

import (
	"context"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/constant"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Runs fixup/squash interactive rebase asynchronously, streams command output to popup viewport,
//	and updates popup processing/success/error state
//
// ------------------------------------
func InteractiveRebaseFixupSquashService(m *types.GittiModel, originalRetrievedGitCommitInfo []git.CommitInfo, sortedSelectedCommits []git.CommitInfo, fixupSquashCommitMessage string, fixupSquashCommitDescription string) {
	ctx, cancel := context.WithCancel(context.Background())

	popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashOutputPopUpModel)
	if ok {
		popUp.CancelFunc = cancel
	}

	go func(ctx context.Context) {
		defer cancel()

		// Reset state before starting new fixup/squash execution.
		popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashOutputPopUpModel)
		if ok {
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
			popUp.IsProcessing.Store(true)
			popUp.IsCancelled.Store(false)
		} else {
			return
		}
		fixupSquashResult, fixupSquashErr := m.GitOperations.GitInteractiveRebase.GitInteractiveRebaseFixupSquash(ctx, originalRetrievedGitCommitInfo, sortedSelectedCommits, fixupSquashCommitMessage, fixupSquashCommitDescription)
		popUp.FixupSquashOutputViewport.SetContentLines(fixupSquashResult)
		popUp, ok = m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashOutputPopUpModel)
		if ok && !popUp.IsCancelled.Load() {
			popUp.IsProcessing.Store(false)
			if fixupSquashErr == nil && !popUp.IsProcessing.Load() {
				popUp.ProcessSuccess.Store(true)
				popUp.IsProcessing.Store(false)
				popUp.HasError.Store(false)
			} else if fixupSquashErr != nil && !popUp.IsProcessing.Load() {
				popUp.HasError.Store(true)
			}
		}
	}(ctx)
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
