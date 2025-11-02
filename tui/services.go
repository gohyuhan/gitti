package tui

import (
	"fmt"
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
		if ok {
			popUp.HasError = false
			popUp.ProcessSuccess = false
			popUp.IsProcessing = true
			popUp.GitCommitOutputViewport.SetContent("")
			git.GITCOMMIT.ClearGitCommitOutput()
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
		popUp, ok = m.PopUpModel.(*GitCommitPopUpModel)
		if ok {
			popUp.IsProcessing = false // update the processing status
			// if sucessful exitcode will be 0
			if exitStatusCode == 0 && popUp.IsProcessing == false {
				popUp.ProcessSuccess = true
				popUp.MessageTextInput.Reset()
				popUp.DescriptionTextAreaInput.Reset()
			} else if exitStatusCode != 0 && popUp.IsProcessing == false {
				popUp.HasError = true
			}
		}
		return
	}()
}

func gitCommitCancelService(m *GittiModel) {
	go func() {
		git.GITCOMMIT.KillGitCommitCmd()     // kill the cmd process if exist
		git.GITCOMMIT.ClearGitCommitOutput() // clear the git commit output log
		m.ShowPopUp = false                  // close the pop up
		m.IsTyping = false                   // and reset typing mode
		popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
		if ok {
			popUp.GitCommitOutputViewport.SetContent("") // set the git commit output viewport to nothing
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
		if ok {
			popUp.HasError = false
			popUp.ProcessSuccess = false
			popUp.IsProcessing = true
			popUp.AddRemoteOutputViewport.SetContent("")
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
		if ok {
			popUp.IsProcessing = false // update the processing status
			// if sucessful exitcode will be 0
			if exitStatusCode == 0 && popUp.IsProcessing == false {
				popUp.ProcessSuccess = true
				popUp.RemoteNameTextInput.Reset()
				popUp.RemoteUrlTextInput.Reset()
				popUp.NoInitialRemote = false
				gitAddRemoteResult = append(gitAddRemoteResult, fmt.Sprintf(i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteAddSuccess, remoteName, remoteUrl))
			} else if exitStatusCode != 0 && popUp.IsProcessing == false {
				popUp.HasError = true
			}
			updateAddRemoteOutputViewport(m, gitAddRemoteResult)
		}
		return
	}()

}

func gitAddRemoteCancelService(m *GittiModel) {
	go func() {
		m.ShowPopUp = false // close the pop up
		m.IsTyping = false  // and reset typing mode
		popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
		if ok {
			popUp.AddRemoteOutputViewport.SetContent("") // set the git commit output viewport to nothing
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
func gitRemotePushService(m *GittiModel, originName string) {
	go func() {
		// git push
		popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)

		if ok {
			popUp.HasError = false
			popUp.ProcessSuccess = false
			popUp.IsProcessing = true
			popUp.GitRemotePushOutputViewport.SetContent("")
			git.GITCOMMIT.ClearGitRemotePushOutput()
		} else {
			return
		}
		var exitStatusCode int
		exitStatusCode = git.GITCOMMIT.GitPush(originName)
		popUp, ok = m.PopUpModel.(*GitRemotePushPopUpModel)
		if ok {
			popUp.IsProcessing = false // update the processing status
			// if sucessful exitcode will be 0
			if exitStatusCode == 0 && popUp.IsProcessing == false {
				popUp.ProcessSuccess = true
			} else if exitStatusCode != 0 && popUp.IsProcessing == false {
				popUp.HasError = true
			}
		}
		return
	}()
}

func gitRemotePushCancelService(m *GittiModel) {
	go func() {
		git.GITCOMMIT.KillGitRemotePushCmd()     // kill the cmd process if exist
		git.GITCOMMIT.ClearGitRemotePushOutput() // clear the git commit output log
		m.ShowPopUp = false                      // close the pop up
		m.IsTyping = false                       // and reset typing mode
		popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
		if ok {
			popUp.GitRemotePushOutputViewport.SetContent("") // set the git commit output viewport to nothing
			popUp.HasError = false
			popUp.ProcessSuccess = false
		}
	}()
}
