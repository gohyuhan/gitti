package services

import (
	"context"

	"github.com/gohyuhan/gitti/tui/constant"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/types"
)

func CreateNewTagService(m *types.GittiModel, commitHash string, tagName string, tagMessage string) {
	go func() {
		m.GitOperations.GitTag.CreateNewTag(commitHash, tagName, tagMessage)
	}()
}

// ------------------------------------
//
//	For Tag delete
//
// ------------------------------------
func DeleteTagService(m *types.GittiModel, remoteName string, tagName string, deleteType string) {
	ctx, cancel := context.WithCancel(context.Background())

	popUp, ok := m.PopUpModel.(*tagPopUp.DeleteTagOutputPopUpModel)
	if ok {
		popUp.CancelFunc = cancel
	}

	go func(ctx context.Context) {
		defer cancel()

		popUp, ok := m.PopUpModel.(*tagPopUp.DeleteTagOutputPopUpModel)
		if ok {
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
			popUp.IsProcessing.Store(true)
			popUp.IsCancelled.Store(false)
		} else {
			return
		}

		output, success := m.GitOperations.GitTag.GitDeleteTag(ctx, remoteName, tagName, deleteType)
		popUp, ok = m.PopUpModel.(*tagPopUp.DeleteTagOutputPopUpModel)
		if ok && !popUp.IsCancelled.Load() {
			popUp.IsProcessing.Store(false) // update the processing status
			// update the output viewport
			tagPopUp.UpdateDeleteTagOutputViewPort(m, output)
			// if successful exitcode will be 0
			if success && !popUp.IsProcessing.Load() {
				popUp.ProcessSuccess.Store(true)
				popUp.IsProcessing.Store(false)
				popUp.HasError.Store(false)
			} else if !success && !popUp.IsProcessing.Load() {
				popUp.ProcessSuccess.Store(false)
				popUp.IsProcessing.Store(false)
				popUp.HasError.Store(true)
			}
		}
	}(ctx)
}

// ------------------------------------
//
//	For Cancelling tag delete
//
// ------------------------------------
func DeleteTagCancelService(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*tagPopUp.DeleteTagOutputPopUpModel)
	if ok {
		popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		if popUp.CancelFunc != nil {
			popUp.CancelFunc() // Cancel the context, which terminates the command and goroutine
		}

		popUp.DeleteTagOutputViewport.SetContent("") // set the delete tag output viewport to nothing
		popUp.IsProcessing.Store(false)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
	}

	m.ShowPopUp.Store(false) // close the pop up
	m.IsTyping.Store(false)  // reset typing mode
	m.PopUpType = constant.NoPopUp
}

// ------------------------------------
//
//	For Git Tag Push
//
// ------------------------------------
func GitPushTagService(m *types.GittiModel, originName string, tagName string, pushType string) {
	// Initialize a cancellable context for the push operation
	ctx, cancel := context.WithCancel(context.Background())

	// Store the cancel function in the popup model so the task can be aborted from the UI
	popUp, ok := m.PopUpModel.(*tagPopUp.PushTagOutputPopUpModel)
	if ok {
		popUp.CancelFunc = cancel
	}

	// Execute the push operation in a separate goroutine to maintain UI responsiveness
	go func(ctx context.Context) {
		defer cancel()

		// Prepare the popup state: reset errors and set processing flags
		popUp, ok := m.PopUpModel.(*tagPopUp.PushTagOutputPopUpModel)
		var exitStatusCode int
		if ok {
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
			popUp.IsProcessing.Store(true)
			popUp.IsCancelled.Store(false)
		} else {
			return
		}

		// Perform the actual git tag push operation via the git service
		exitStatusCode = m.GitOperations.GitTag.GitPushTag(ctx, originName, tagName, pushType)

		// After the operation completes, update the UI state if the popup hasn't been cancelled
		popUp, ok = m.PopUpModel.(*tagPopUp.PushTagOutputPopUpModel)
		if ok && !popUp.IsCancelled.Load() {
			popUp.IsProcessing.Store(false) // Update the processing status

			// Check the exit status: 0 typically indicates success
			if exitStatusCode == 0 && !popUp.IsProcessing.Load() {
				popUp.ProcessSuccess.Store(true)
				popUp.IsProcessing.Store(false)
				popUp.HasError.Store(false)
			} else if exitStatusCode != 0 && !popUp.IsProcessing.Load() {
				// If the command failed, set the error flag
				popUp.HasError.Store(true)
			}
		}
	}(ctx)
}

func GitPushTagCancelService(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*tagPopUp.PushTagOutputPopUpModel)
	if ok {
		popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		if popUp.CancelFunc != nil {
			popUp.CancelFunc() // Cancel the context, which terminates the command and goroutine
		}
	}
	m.GitOperations.GitTag.ClearGitPushTagOutput() // clear the git tag output log

	m.ShowPopUp.Store(false) // close the pop up
	m.IsTyping.Store(false)  // reset typing mode
	m.PopUpType = constant.NoPopUp
	if ok {
		popUp.PushTagOutputViewport.SetContent("") // set the git commit output viewport to nothing
		popUp.IsProcessing.Store(false)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
	}
}

// ------------------------------------
//
//	For Git Tag Fetch
//
// ------------------------------------
func GitFetchTagService(m *types.GittiModel, originName string, fetchType string) {
	// Initialize a cancellable context for the push operation
	ctx, cancel := context.WithCancel(context.Background())

	// Store the cancel function in the popup model so the task can be aborted from the UI
	popUp, ok := m.PopUpModel.(*tagPopUp.FetchTagOutputPopUpModel)
	if ok {
		popUp.CancelFunc = cancel
	}

	// Execute the push operation in a separate goroutine to maintain UI responsiveness
	go func(ctx context.Context) {
		defer cancel()

		// Prepare the popup state: reset errors and set processing flags
		popUp, ok := m.PopUpModel.(*tagPopUp.FetchTagOutputPopUpModel)
		var exitStatusCode int
		if ok {
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
			popUp.IsProcessing.Store(true)
			popUp.IsCancelled.Store(false)
		} else {
			return
		}

		// Perform the actual git tag push operation via the git service
		exitStatusCode = m.GitOperations.GitTag.GitFetchTag(ctx, originName, fetchType)

		// After the operation completes, update the UI state if the popup hasn't been cancelled
		popUp, ok = m.PopUpModel.(*tagPopUp.FetchTagOutputPopUpModel)
		if ok && !popUp.IsCancelled.Load() {
			popUp.IsProcessing.Store(false) // Update the processing status

			// Check the exit status: 0 typically indicates success
			if exitStatusCode == 0 && !popUp.IsProcessing.Load() {
				popUp.ProcessSuccess.Store(true)
				popUp.IsProcessing.Store(false)
				popUp.HasError.Store(false)
			} else if exitStatusCode != 0 && !popUp.IsProcessing.Load() {
				// If the command failed, set the error flag
				popUp.HasError.Store(true)
			}
		}
	}(ctx)
}

func GitFetchTagCancelService(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*tagPopUp.FetchTagOutputPopUpModel)
	if ok {
		popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		if popUp.CancelFunc != nil {
			popUp.CancelFunc() // Cancel the context, which terminates the command and goroutine
		}
	}
	m.GitOperations.GitTag.ClearGitFetchTagOutput() // clear the git tag output log

	m.ShowPopUp.Store(false) // close the pop up
	m.IsTyping.Store(false)  // reset typing mode
	m.PopUpType = constant.NoPopUp
	if ok {
		popUp.FetchTagOutputViewport.SetContent("") // set the git commit output viewport to nothing
		popUp.IsProcessing.Store(false)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
	}
}
