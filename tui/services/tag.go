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
