package tui

import (
	"gitti/api/git"
)

// service.go was to bridge api and the needs of the terminal interface logic so that it can be compatible and feels smooth and not cluncy
func gitCommitService(m *GittiModel) {
	go func() {
		// set to is processing and remove the log content in UI and also in GITCOMMIT in memory
		popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
		var message string
		var description string
		var exitStatusCode int
		if ok {
			popUp.HasError = false
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
			if exitStatusCode == 0 && m.PopUpModel.(*GitCommitPopUpModel).IsProcessing == false {
				popUp.MessageTextInput.Reset()
				popUp.DescriptionTextAreaInput.Reset()
			} else if exitStatusCode != 0 && m.PopUpModel.(*GitCommitPopUpModel).IsProcessing == false {
				popUp.HasError = true
			}
		}
		return
	}()
}

func gitCommitCancelService(m *GittiModel) {
	go func() {
		git.GITCOMMIT.KillCommit()           // kill the cmd process if exist
		git.GITCOMMIT.ClearGitCommitOutput() // clear the git commit output log
		m.ShowPopUp = false                  // close the pop up
		m.IsTyping = false                   // and reset typing mode
		popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
		if ok {
			popUp.GitCommitOutputViewport.SetContent("") // set the git commit output viewport to nothing
			popUp.HasError = false
		}
	}()
}
