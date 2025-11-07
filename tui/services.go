package tui

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"gitti/i18n"
	"gitti/tui/constant"
)

// service.go was to bridge api and the needs of the terminal interface logic so that it can be compatible and feels smooth and not cluncy

// ------------------------------------
//
//	For Git Commit
//
// ------------------------------------
func gitCommitService(m *GittiModel) {
	go func() {
		// set to is processing and remove the log content in UI and also in GITCOMMIT in memory
		popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
		var message string
		var description string
		var exitStatusCode int
		var sessionID uuid.UUID
		if ok {
			sessionID = popUp.SessionID // Capture the session ID at start
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
			popUp.IsProcessing.Store(true)
			popUp.IsCancelled.Store(false)
			// retrieve the value of commit message and desc
			message = popUp.MessageTextInput.Value()
			description = popUp.DescriptionTextAreaInput.Value()
		} else {
			return
		}
		if len(message) < 1 {
			return
		}
		// stage the changes based on what was chosen by user and commit it
		exitStatusCode = m.GitState.GitCommit.GitStageAndCommit(message, description, m.GitState.GitFiles.GetSelectedForStageFiles())

		popUp, ok = m.PopUpModel.(*GitCommitPopUpModel)
		if ok && !popUp.IsCancelled.Load() {
			popUp.IsProcessing.Store(false) // update the processing status
			// if sucessful exitcode will be 0
			if exitStatusCode == 0 && !popUp.IsProcessing.Load() {
				popUp.ProcessSuccess.Store(true)
				popUp.MessageTextInput.Reset()
				popUp.DescriptionTextAreaInput.Reset()
				time.Sleep(constant.AUTOCLOSEINTERVAL * time.Millisecond)
				// Check if user cancelled during sleep and verify this is still the same popup session
				popUp, ok = m.PopUpModel.(*GitCommitPopUpModel)
				if ok && !popUp.IsCancelled.Load() && popUp.SessionID == sessionID {
					m.GitState.GitCommit.ClearGitCommitOutput() // clear the git commit output log
					m.ShowPopUp.Store(false)                    // close the pop up
					m.IsTyping.Store(false)
					popUp.GitCommitOutputViewport.SetContent("") // set the git commit output viewport to nothing
					popUp.IsProcessing.Store(false)
					popUp.HasError.Store(false)
					popUp.ProcessSuccess.Store(false)
				}
			} else if exitStatusCode != 0 && !popUp.IsProcessing.Load() {
				popUp.HasError.Store(true)
			}
		}
	}()
}

func gitCommitCancelService(m *GittiModel) {
	go func() {
		popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
		if ok {
			popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		}
		// Clean up git processes and state
		m.GitState.GitCommit.KillGitStageAndCommitCmd() // kill the git stash and commit cmd process if exist
		m.GitState.GitCommit.ClearGitCommitOutput()     // clear the git commit output log

		m.ShowPopUp.Store(false) // close the pop up
		m.IsTyping.Store(false)  // reset typing mode
		if ok {
			popUp.GitCommitOutputViewport.SetContent("") // set the git commit output viewport to nothing
			popUp.IsProcessing.Store(false)
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
		}
	}()
}

// ------------------------------------
//
//	For Adding Git Remote
//
// ------------------------------------
func gitAddRemoteService(m *GittiModel) {
	go func() {
		// set to is processing and remove the log content in UI and also in GITCOMMIT in memory
		popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
		var remoteName string
		var remoteUrl string
		var exitStatusCode int
		var sessionID uuid.UUID
		if ok {
			sessionID = popUp.SessionID // Capture the session ID at start
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
			popUp.IsProcessing.Store(true)
			popUp.IsCancelled.Store(false)

			// retrieve the value of remote name and remote url
			remoteName = popUp.RemoteNameTextInput.Value()
			remoteUrl = popUp.RemoteUrlTextInput.Value()
		} else {
			return
		}
		if len(remoteName) < 1 || len(remoteUrl) < 1 {
			return
		}
		gitAddRemoteResult, exitStatusCode := m.GitState.GitCommit.GitAddRemote(remoteName, remoteUrl)
		popUp, ok = m.PopUpModel.(*AddRemotePromptPopUpModel)
		if ok && !popUp.IsCancelled.Load() {
			popUp.IsProcessing.Store(false) // update the processing status
			// if sucessful exitcode will be 0
			if exitStatusCode == 0 && !popUp.IsProcessing.Load() {
				popUp.ProcessSuccess.Store(true)
				popUp.RemoteNameTextInput.Reset()
				popUp.RemoteUrlTextInput.Reset()
				popUp.NoInitialRemote = false
				gitAddRemoteResult = append(gitAddRemoteResult, fmt.Sprintf(i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteAddSuccess, remoteName, remoteUrl))
				updateAddRemoteOutputViewport(m, gitAddRemoteResult)
				time.Sleep(constant.AUTOCLOSEINTERVAL * time.Millisecond)
				// Check if user cancelled during sleep and verify this is still the same popup session
				popUp, ok = m.PopUpModel.(*AddRemotePromptPopUpModel)
				if ok && !popUp.IsCancelled.Load() && popUp.SessionID == sessionID {
					m.ShowPopUp.Store(false) // close the pop up
					m.IsTyping.Store(false)
					popUp.AddRemoteOutputViewport.SetContent("") // set the git commit output viewport to nothing
					popUp.IsProcessing.Store(false)
					popUp.HasError.Store(false)
					popUp.ProcessSuccess.Store(false)
				}
			} else if exitStatusCode != 0 && !popUp.IsProcessing.Load() {
				popUp.HasError.Store(true)
				updateAddRemoteOutputViewport(m, gitAddRemoteResult)
			}
		}
	}()
}

func gitAddRemoteCancelService(m *GittiModel) {
	go func() {
		popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
		if ok {
			popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		}
		// Clean up git remote add process
		m.GitState.GitCommit.KillGitAddRemoteCmd() // kill the cmd process if exist

		m.ShowPopUp.Store(false) // close the pop up
		m.IsTyping.Store(false)  // reset typing mode
		if ok {
			popUp.AddRemoteOutputViewport.SetContent("") // set the git commit output viewport to nothing
			popUp.IsProcessing.Store(false)
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
		}
	}()
}

// ------------------------------------
//
//	For Git Push
//
// ------------------------------------
func gitRemotePushService(m *GittiModel, originName string, pushType string) {
	go func() {
		// git push
		popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
		var sessionID uuid.UUID

		if ok {
			sessionID = popUp.SessionID // Capture the session ID at start
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
			popUp.IsProcessing.Store(true)
			popUp.IsCancelled.Store(false)
		} else {
			return
		}

		exitStatusCode := m.GitState.GitCommit.GitPush(m.GitState.GitBranch.CurrentCheckOut().BranchName, originName, originName)
		popUp, ok = m.PopUpModel.(*GitRemotePushPopUpModel)
		if ok && !popUp.IsCancelled.Load() {
			popUp.IsProcessing.Store(false) // update the processing status
			// if sucessful exitcode will be 0
			if exitStatusCode == 0 && !popUp.IsProcessing.Load() {
				popUp.ProcessSuccess.Store(true)
				time.Sleep(constant.AUTOCLOSEINTERVAL * time.Millisecond)
				// Check if user cancelled during sleep and verify this is still the same popup session
				popUp, ok = m.PopUpModel.(*GitRemotePushPopUpModel)
				if ok && !popUp.IsCancelled.Load() && popUp.SessionID == sessionID {
					m.GitState.GitCommit.ClearGitRemotePushOutput() // clear the git commit output log
					m.ShowPopUp.Store(false)                        // close the pop up
					m.IsTyping.Store(false)
					popUp.GitRemotePushOutputViewport.SetContent("") // set the git commit output viewport to nothing
					popUp.IsProcessing.Store(false)
					popUp.HasError.Store(false)
					popUp.ProcessSuccess.Store(false)
				}
			} else if exitStatusCode != 0 && !popUp.IsProcessing.Load() {
				popUp.HasError.Store(true)
			}
		}
	}()
}

func gitRemotePushCancelService(m *GittiModel) {
	go func() {
		popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
		if ok {
			popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		}
		m.GitState.GitCommit.KillGitRemotePushCmd()     // kill the cmd process if exist
		m.GitState.GitCommit.ClearGitRemotePushOutput() // clear the git commit output log
		m.ShowPopUp.Store(false)                        // close the pop up
		m.IsTyping.Store(false)                         // and reset typing mode
		if ok {
			popUp.GitRemotePushOutputViewport.SetContent("") // set the git commit output viewport to nothing
			popUp.IsProcessing.Store(false)
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
		}
	}()
}
