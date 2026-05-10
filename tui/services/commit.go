package services

import (
	"context"

	"github.com/gohyuhan/gitti/tui/constant"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	"github.com/gohyuhan/gitti/tui/types"
)

// services was to bridge api and the needs of the terminal interface logic so that it can be compatible and feels smooth and not clunky
// ------------------------------------
//
//	For Git Commit
//
// ------------------------------------
func GitCommitService(m *types.GittiModel, isAmendCommit bool) {
	popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
	if ok {
		ctx, cancel := context.WithCancel(context.Background())
		popUp.CancelFunc = cancel
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsProcessing.Store(true)
		popUp.IsCancelled.Store(false)
		message := popUp.MessageTextInput.Value()
		description := popUp.DescriptionTextAreaInput.Value()
		if len(message) < 1 {
			popUp.IsProcessing.Store(false)
			cancel()
			return
		}

		go func(ctx context.Context, message string, description string) {
			defer cancel()

			exitStatusCode := m.GitOperations.GitCommit.GitCommit(ctx, message, description, isAmendCommit)
			data := types.GitCommitResultEventDataStructure{
				Success: exitStatusCode == 0,
			}
			m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
				Event: constant.GIT_COMMIT_RESULT_EVENT,
				Data:  data,
			}
		}(ctx, message, description)
	}
}

// ------------------------------------
//
//	Cancel the current git commit operation and clean up pop-up state
//
// ------------------------------------
func GitCommitCancelService(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
	if ok {
		popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		if popUp.CancelFunc != nil {
			popUp.CancelFunc() // Cancel the context, which terminates the command and goroutine
		}
	}
	m.GitOperations.GitCommit.ClearGitCommitOutput() // clear the git commit output log

	m.ShowPopUp.Store(false) // close the pop up
	m.IsTyping.Store(false)  // reset typing mode
	m.PopUpType = constant.NoPopUp
	if ok {
		popUp.GitCommitOutputViewport.SetContent("") // set the git commit output viewport to nothing
		popUp.IsProcessing.Store(false)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
	}
}

// ------------------------------------
//
//	For Git Amend Commit
//
// ------------------------------------
func GitAmendCommitService(m *types.GittiModel, isAmendCommit bool) {
	popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
	if ok {
		ctx, cancel := context.WithCancel(context.Background())
		popUp.CancelFunc = cancel
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsProcessing.Store(true)
		popUp.IsCancelled.Store(false)
		message := popUp.MessageTextInput.Value()
		description := popUp.DescriptionTextAreaInput.Value()
		if len(message) < 1 {
			popUp.IsProcessing.Store(false)
			cancel()
			return
		}

		go func(ctx context.Context, message string, description string) {
			defer cancel()

			exitStatusCode := m.GitOperations.GitCommit.GitCommit(ctx, message, description, isAmendCommit)
			data := types.GitAmendCommitResultEventDataStructure{
				Success: exitStatusCode == 0,
			}
			m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
				Event: constant.GIT_AMEND_COMMIT_RESULT_EVENT,
				Data:  data,
			}
		}(ctx, message, description)
	}
}

// ------------------------------------
//
//	Cancel the current git amend commit operation and clean up pop-up state
//
// ------------------------------------
func GitAmendCommitCancelService(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
	if ok {
		popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		if popUp.CancelFunc != nil {
			popUp.CancelFunc() // Cancel the context, which terminates the command and goroutine
		}
	}
	m.GitOperations.GitCommit.ClearGitCommitOutput() // clear the git commit output log

	m.ShowPopUp.Store(false) // close the pop up
	m.IsTyping.Store(false)  // reset typing mode
	m.PopUpType = constant.NoPopUp
	if ok {
		popUp.GitAmendCommitOutputViewport.SetContent("") // set the git commit output viewport to nothing
		popUp.IsProcessing.Store(false)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
	}
}

// ------------------------------------
//
//	For Git Reset Commit Service (apply to latest commit only)
//
// ------------------------------------
func GitResetLatestCommitService(m *types.GittiModel, resetType string) {
	go func() {
		m.GitOperations.GitCommit.GitResetLatestCommit(resetType)
	}()
}

// ------------------------------------
//
//	For Git Reset Commit Service (apply to selected commit [using commit hash])
//
// ------------------------------------
func GitResetToSelectedCommitService(m *types.GittiModel, resetType string, commitHash string) {
	go func() {
		m.GitOperations.GitCommit.GitResetToSelectedCommit(resetType, commitHash)
	}()
}
