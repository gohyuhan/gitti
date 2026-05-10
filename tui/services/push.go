package services

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	pushPopUp "github.com/gohyuhan/gitti/tui/popup/push"
	"github.com/gohyuhan/gitti/tui/types"
)

// services was to bridge api and the needs of the terminal interface logic so that it can be compatible and feels smooth and not clunky
// ------------------------------------
//
//	For Git Remote Push
//
// ------------------------------------
func GitRemotePushService(m *types.GittiModel, remoteName string, pushType string) {
	popUp, ok := m.PopUpModel.(*pushPopUp.GitRemotePushPopUpModel)
	checkoutBranch := m.CheckOutBranch
	if ok {
		ctx, cancel := context.WithCancel(context.Background())
		popUp.CancelFunc = cancel
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsProcessing.Store(true)
		popUp.IsCancelled.Store(false)

		go func(ctx context.Context) {
			defer cancel()

			exitStatusCode := m.GitOperations.GitCommit.GitPush(ctx, remoteName, pushType, checkoutBranch)
			data := types.GitPushResultEventDataInterface{
				Success: exitStatusCode == 0,
			}
			m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
				Event: constant.GIT_PUSH_RESULT_EVENT,
				Data:  data,
			}
		}(ctx)
	}
}

// ------------------------------------
//
//	Cancel the current git push operation and clean up pop-up state
//
// ------------------------------------
func GitRemotePushCancelService(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*pushPopUp.GitRemotePushPopUpModel)
	if ok {
		popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		if popUp.CancelFunc != nil {
			popUp.CancelFunc() // Cancel the context, which terminates the command and goroutine
		}
	}
	m.GitOperations.GitCommit.ClearGitRemotePushOutput() // clear the git commit output log
	m.ShowPopUp.Store(false)                             // close the pop up
	m.IsTyping.Store(false)                              // and reset typing mode
	m.PopUpType = constant.NoPopUp
	if ok {
		popUp.GitRemotePushOutputViewport.SetContent("") // set the git commit output viewport to nothing
		popUp.IsProcessing.Store(false)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
	}
}

// ------------------------------------
//
//	Initialize the push pop-up model and start the push operation
//
// ------------------------------------
func InitGitRemotePushPopUpModelAndStartGitRemotePushService(m *types.GittiModel, remoteName string, pushType string) (*types.GittiModel, tea.Cmd) {
	m.GitOperations.GitCommit.ClearGitRemotePushOutput()
	if popUp, ok := m.PopUpModel.(*pushPopUp.GitRemotePushPopUpModel); !ok {
		pushPopUp.InitGitRemotePushPopUpModel(m)
	} else {
		popUp.GitRemotePushOutputViewport.SetContent("")
	}
	// then push it after init the git remote push pop up model
	GitRemotePushService(m, remoteName, pushType)
	// Start spinner ticking
	if pushPopup, ok := m.PopUpModel.(*pushPopUp.GitRemotePushPopUpModel); ok {
		return m, pushPopup.Spinner.Tick
	}
	return m, nil
}
