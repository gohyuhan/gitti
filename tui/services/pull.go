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
	popUp, ok := m.PopUpModel.(*pullPopUp.GitPullOutputPopUpModel)
	if ok {
		ctx, cancel := context.WithCancel(context.Background())
		popUp.CancelFunc = cancel
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsProcessing.Store(true)
		popUp.IsCancelled.Store(false)

		go func(ctx context.Context) {
			defer cancel()

			exitStatusCode := m.GitOperations.GitPull.GitPull(ctx, pullType)
			data := types.GitPullResultEventDataStructure{
				Success: exitStatusCode == 0,
			}
			m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
				Event: constant.GIT_PULL_RESULT_EVENT,
				Data:  data,
			}
		}(ctx)
	}
}

// ------------------------------------
//
//	Cancel the current git pull operation and clean up pop-up state
//
// ------------------------------------
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
