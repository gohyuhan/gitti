package tui

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"gitti/api/git"
	"gitti/i18n"
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
			popUp.HasError = false
			popUp.ProcessSuccess = false
			popUp.IsProcessing = true
			popUp.IsCancelled = false
			// retrieve the value of commit message and desc
			message = popUp.MessageTextInput.Value()
			description = popUp.DescriptionTextAreaInput.Value()
		} else {
			return
		}
		if len(message) < 1 {
			return
		}
		// stage the changes based on what was chosen by user
		git.GITCOMMIT.GitStage()
		// and commit it
		exitStatusCode = git.GITCOMMIT.GitCommit(message, description)

		// after the commit, we set back the Is Selected for Stage state.
		defer git.GITFILES.UpdateFilesStageStatusAfterCommit()
		popUp, ok = m.PopUpModel.(*GitCommitPopUpModel)
		if ok && !popUp.IsCancelled {
			popUp.IsProcessing = false // update the processing status
			// if sucessful exitcode will be 0
			if exitStatusCode == 0 && !popUp.IsProcessing {
				popUp.ProcessSuccess = true
				popUp.MessageTextInput.Reset()
				popUp.DescriptionTextAreaInput.Reset()
				time.Sleep(AUTOCLOSEINTERVAL * time.Millisecond)
				// Check if user cancelled during sleep and verify this is still the same popup session
				popUp, ok = m.PopUpModel.(*GitCommitPopUpModel)
				if ok && !popUp.IsCancelled && popUp.SessionID == sessionID {
					git.GITCOMMIT.ClearGitCommitOutput() // clear the git commit output log
					m.ShowPopUp = false                  // close the pop up
					m.IsTyping = false
					popUp.GitCommitOutputViewport.SetContent("") // set the git commit output viewport to nothing
					popUp.IsProcessing = false
					popUp.HasError = false
					popUp.ProcessSuccess = false
				}
			} else if exitStatusCode != 0 && !popUp.IsProcessing {
				popUp.HasError = true
			}
		}
	}()
}

func gitCommitCancelService(m *GittiModel) {
	go func() {
		popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
		if ok {
			popUp.IsCancelled = true // set cancellation flag first to prevent race condition
		}
		// Clean up git processes and state
		git.GITCOMMIT.KillGitCommitCmd()     // kill the git commit cmd process if exist
		git.GITCOMMIT.KillGitStageCmd()      // kill the git stage cmd process if exist
		git.GITCOMMIT.ClearGitCommitOutput() // clear the git commit output log

		m.ShowPopUp = false // close the pop up
		m.IsTyping = false  // reset typing mode
		if ok {
			popUp.GitCommitOutputViewport.SetContent("") // set the git commit output viewport to nothing
			popUp.IsProcessing = false
			popUp.HasError = false
			popUp.ProcessSuccess = false
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
			popUp.HasError = false
			popUp.ProcessSuccess = false
			popUp.IsProcessing = true
			popUp.IsCancelled = false

			// retrieve the value of remote name and remote url
			remoteName = popUp.RemoteNameTextInput.Value()
			remoteUrl = popUp.RemoteUrlTextInput.Value()
		} else {
			return
		}
		if len(remoteName) < 1 || len(remoteUrl) < 1 {
			return
		}
		gitAddRemoteResult, exitStatusCode := git.GITCOMMIT.GitAddRemote(remoteName, remoteUrl)
		popUp, ok = m.PopUpModel.(*AddRemotePromptPopUpModel)
		if ok && !popUp.IsCancelled {
			popUp.IsProcessing = false // update the processing status
			// if sucessful exitcode will be 0
			if exitStatusCode == 0 && !popUp.IsProcessing {
				popUp.ProcessSuccess = true
				popUp.RemoteNameTextInput.Reset()
				popUp.RemoteUrlTextInput.Reset()
				popUp.NoInitialRemote = false
				gitAddRemoteResult = append(gitAddRemoteResult, fmt.Sprintf(i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteAddSuccess, remoteName, remoteUrl))
				updateAddRemoteOutputViewport(m, gitAddRemoteResult)
				time.Sleep(AUTOCLOSEINTERVAL * time.Millisecond)
				// Check if user cancelled during sleep and verify this is still the same popup session
				popUp, ok = m.PopUpModel.(*AddRemotePromptPopUpModel)
				if ok && !popUp.IsCancelled && popUp.SessionID == sessionID {
					m.ShowPopUp = false // close the pop up
					m.IsTyping = false
					popUp.AddRemoteOutputViewport.SetContent("") // set the git commit output viewport to nothing
					popUp.IsProcessing = false
					popUp.HasError = false
					popUp.ProcessSuccess = false
				}
			} else if exitStatusCode != 0 && !popUp.IsProcessing {
				popUp.HasError = true
				updateAddRemoteOutputViewport(m, gitAddRemoteResult)
			}
		}
	}()
}

func gitAddRemoteCancelService(m *GittiModel) {
	go func() {
		popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
		if ok {
			popUp.IsCancelled = true // set cancellation flag first to prevent race condition
		}
		// Clean up git remote add process
		git.GITCOMMIT.KillGitAddRemoteCmd() // kill the cmd process if exist

		m.ShowPopUp = false // close the pop up
		m.IsTyping = false  // reset typing mode
		if ok {
			popUp.AddRemoteOutputViewport.SetContent("") // set the git commit output viewport to nothing
			popUp.IsProcessing = false
			popUp.HasError = false
			popUp.ProcessSuccess = false
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
			popUp.HasError = false
			popUp.ProcessSuccess = false
			popUp.IsProcessing = true
			popUp.IsCancelled = false
		} else {
			return
		}

		exitStatusCode := git.GITCOMMIT.GitPush(originName, pushType)
		popUp, ok = m.PopUpModel.(*GitRemotePushPopUpModel)
		if ok && !popUp.IsCancelled {
			popUp.IsProcessing = false // update the processing status
			// if sucessful exitcode will be 0
			if exitStatusCode == 0 && !popUp.IsProcessing {
				popUp.ProcessSuccess = true
				time.Sleep(AUTOCLOSEINTERVAL * time.Millisecond)
				// Check if user cancelled during sleep and verify this is still the same popup session
				popUp, ok = m.PopUpModel.(*GitRemotePushPopUpModel)
				if ok && !popUp.IsCancelled && popUp.SessionID == sessionID {
					git.GITCOMMIT.ClearGitRemotePushOutput() // clear the git commit output log
					m.ShowPopUp = false                      // close the pop up
					m.IsTyping = false
					popUp.GitRemotePushOutputViewport.SetContent("") // set the git commit output viewport to nothing
					popUp.IsProcessing = false
					popUp.HasError = false
					popUp.ProcessSuccess = false
				}
			} else if exitStatusCode != 0 && !popUp.IsProcessing {
				popUp.HasError = true
			}
		}
	}()
}

func gitRemotePushCancelService(m *GittiModel) {
	go func() {
		popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
		if ok {
			popUp.IsCancelled = true // set cancellation flag first to prevent race condition
		}
		git.GITCOMMIT.KillGitRemotePushCmd()     // kill the cmd process if exist
		git.GITCOMMIT.ClearGitRemotePushOutput() // clear the git commit output log
		m.ShowPopUp = false                      // close the pop up
		m.IsTyping = false                       // and reset typing mode
		if ok {
			popUp.GitRemotePushOutputViewport.SetContent("") // set the git commit output viewport to nothing
			popUp.IsProcessing = false
			popUp.HasError = false
			popUp.ProcessSuccess = false
		}
	}()
}
