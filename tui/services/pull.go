package services

import (
	"context"

	"github.com/gohyuhan/gitti/tui/constant"
	pullPopUp "github.com/gohyuhan/gitti/tui/popup/pull"
	"github.com/gohyuhan/gitti/tui/types"
)

// services was to bridge api and the needs of the terminal interface logic so that it can be compatible and feels smooth and not clunky
// ------------------------------------
//
//	For Git Pull
//
// ------------------------------------
func GitPullService(m *types.GittiModel, pullType string) {
	ctx, cancel := context.WithCancel(context.Background())

	popUp, ok := m.PopUpModel.(*pullPopUp.GitPullOutputPopUpModel)
	if ok {
		popUp.CancelFunc = cancel
	}

	go func(ctx context.Context) {
		defer cancel()

		// set to is processing and remove the log content in UI and also in GITCOMMIT in memory
		popUp, ok := m.PopUpModel.(*pullPopUp.GitPullOutputPopUpModel)
		var exitStatusCode int
		if ok {
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
			popUp.IsProcessing.Store(true)
			popUp.IsCancelled.Store(false)
		} else {
			return
		}
		exitStatusCode = m.GitOperations.GitPull.GitPull(ctx, pullType)
		popUp, ok = m.PopUpModel.(*pullPopUp.GitPullOutputPopUpModel)
		if ok && !popUp.IsCancelled.Load() {
			popUp.IsProcessing.Store(false) // update the processing status
			// if successful exitcode will be 0
			if exitStatusCode == 0 && !popUp.IsProcessing.Load() {
				popUp.ProcessSuccess.Store(true)
				popUp.IsProcessing.Store(false)
				popUp.HasError.Store(false)
			} else if exitStatusCode != 0 && !popUp.IsProcessing.Load() {
				popUp.HasError.Store(true)
			}
		}
	}(ctx)
}

func GitPullCancelService(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*pullPopUp.GitPullOutputPopUpModel)
	if ok {
		popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		if popUp.CancelFunc != nil {
			popUp.CancelFunc() // Cancel the context, which terminates the command and goroutine
		}
	}
	m.GitOperations.GitPull.ClearGitPullOutput() // clear the git commit output log
	m.ShowPopUp.Store(false)                     // close the pop up
	m.IsTyping.Store(false)                      // and reset typing mode
	m.PopUpType = constant.NoPopUp
	if ok {
		popUp.GitPullOutputViewport.SetContent("") // set the git commit output viewport to nothing
		popUp.IsProcessing.Store(false)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
	}
}
