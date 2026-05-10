package services

import (
	"context"
	"fmt"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	For Adding Git Remote
//
// ------------------------------------
func GitAddRemoteService(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel)
	if ok {
		ctx, cancel := context.WithCancel(context.Background())
		popUp.CancelFunc = cancel
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsProcessing.Store(true)
		popUp.IsCancelled.Store(false)

		remoteName := popUp.RemoteNameTextInput.Value()
		remoteUrl := popUp.RemoteUrlTextInput.Value()
		if len(remoteName) < 1 || len(remoteUrl) < 1 {
			popUp.IsProcessing.Store(false)
			cancel()
			return
		}

		go func(ctx context.Context, remoteName string, remoteUrl string) {
			defer cancel()

			gitAddRemoteResult, exitStatusCode := m.GitOperations.GitRemote.GitAddRemote(ctx, remoteName, remoteUrl)
			if exitStatusCode == 0 {
				gitAddRemoteResult = append(gitAddRemoteResult, fmt.Sprintf(i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteAddSuccess, remoteName, remoteUrl))
			}

			data := types.GitAddRemoteResultEventDataInterface{
				Result:  gitAddRemoteResult,
				Success: exitStatusCode == 0,
			}
			m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
				Event: constant.GIT_ADD_REMOTE_RESULT_EVENT,
				Data:  data,
			}
		}(ctx, remoteName, remoteUrl)
	} else {
		return
	}
}

// ------------------------------------
//
//	Cancel the current add remote operation and clean up pop-up state
//
// ------------------------------------
func GitAddRemoteCancelService(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel)
	if ok {
		popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		if popUp.CancelFunc != nil {
			popUp.CancelFunc() // Cancel the context, which terminates the command and goroutine
		}
	}

	m.ShowPopUp.Store(false) // close the pop up
	m.IsTyping.Store(false)  // reset typing mode
	m.PopUpType = constant.NoPopUp
	if ok {
		popUp.AddRemoteOutputViewport.SetContent("") // set the git commit output viewport to nothing
		popUp.IsProcessing.Store(false)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
	}
}

// ------------------------------------
//
//	For Removal of Remote
//
// ------------------------------------
func GitRemoveRemoteService(m *types.GittiModel, remoteName string) {
	go func() {
		m.GitOperations.GitRemote.GitRemoveRemote(remoteName)
	}()
}

// ------------------------------------
//
//	For setting remote as tracking upstream
//
// ------------------------------------
func GitSetRemoteAsTrackingUpstreamService(m *types.GittiModel, remoteName string) {
	go func() {
		m.GitOperations.GitRemote.GitSetRemoteAsTrackingUpstream(remoteName, m.CheckOutBranch)
	}()
}

// ------------------------------------
//
//	For editing remote name and url
//
// ------------------------------------
func GitEditRemoteNameAndUrlService(m *types.GittiModel, oldRemoteName string, newRemoteName string, oldRemoteUrl string, newRemoteUrl string) {
	go func() {
		if oldRemoteName != newRemoteName {
			m.GitOperations.GitRemote.GitChangeRemoteName(oldRemoteName, newRemoteName)
		}

		if oldRemoteUrl != newRemoteUrl {
			m.GitOperations.GitRemote.GitChangeRemoteUrl(newRemoteName, newRemoteUrl)
		}
	}()
}
