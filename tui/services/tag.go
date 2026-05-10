package services

import (
	"context"

	"github.com/gohyuhan/gitti/tui/constant"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Create a new tag at the specified commit with a name and message
//
// ------------------------------------
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

	popUp, ok := m.PopUpModel.(*tagPopUp.DeleteTagOutputPopUpModel)
	if ok {
		ctx, cancel := context.WithCancel(context.Background())
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsProcessing.Store(true)
		popUp.IsCancelled.Store(false)
		popUp.CancelFunc = cancel

		go func(ctx context.Context) {
			defer cancel()

			output, success := m.GitOperations.GitTag.GitDeleteTag(ctx, remoteName, tagName, deleteType)
			data := types.GitDeleteTagResultEventDataInterface{
				Result:  output,
				Success: success,
			}
			m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
				Event: constant.GIT_DELETE_TAG_RESULT_EVENT,
				Data:  data,
			}
		}(ctx)
	} else {
		return
	}
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

	// Store the cancel function in the popup model so the task can be aborted from the UI
	popUp, ok := m.PopUpModel.(*tagPopUp.PushTagOutputPopUpModel)
	if ok {
		// Initialize a cancellable context for the push operation
		ctx, cancel := context.WithCancel(context.Background())
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsProcessing.Store(true)
		popUp.IsCancelled.Store(false)
		popUp.CancelFunc = cancel
		// Execute the push operation in a separate goroutine to maintain UI responsiveness
		go func(ctx context.Context) {
			defer cancel()

			var exitStatusCode int
			// Perform the actual git tag push operation via the git service
			exitStatusCode = m.GitOperations.GitTag.GitPushTag(ctx, originName, tagName, pushType)

			data := types.GitPushTagResultEventDataInterface{
				Success: exitStatusCode == 0,
			}
			m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
				Event: constant.GIT_PUSH_TAG_RESULT_EVENT,
				Data:  data,
			}
		}(ctx)
	} else {
		return
	}

}

// ------------------------------------
//
//	Cancel the current tag push operation and clean up pop-up state
//
// ------------------------------------
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
	// Store the cancel function in the popup model so the task can be aborted from the UI
	popUp, ok := m.PopUpModel.(*tagPopUp.FetchTagOutputPopUpModel)
	if ok {
		// Initialize a cancellable context for the fetch operation
		ctx, cancel := context.WithCancel(context.Background())
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsProcessing.Store(true)
		popUp.IsCancelled.Store(false)
		popUp.CancelFunc = cancel

		// Execute the push operation in a separate goroutine to maintain UI responsiveness
		go func(ctx context.Context) {
			defer cancel()

			var exitStatusCode int
			// Perform the actual git tag push operation via the git service
			exitStatusCode = m.GitOperations.GitTag.GitFetchTag(ctx, originName, fetchType)

			data := types.GitFetchTagResultEventDataInterface{
				Success: exitStatusCode == 0,
			}
			m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
				Event: constant.GIT_FETCH_TAG_RESULT_EVENT,
				Data:  data,
			}
		}(ctx)
	} else {
		return
	}

}

// ------------------------------------
//
//	Cancel the current tag fetch operation and clean up pop-up state
//
// ------------------------------------
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
