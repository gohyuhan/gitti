package tui

import (
	"gitti/api/git"
	"time"
	// "time"
)

// service.go was to bridge api and the needs of the terminal interface logic so that it can be compatible and feels smooth and not cluncy

func gitCommitService(m *GittiModel) {
	go func() {
		// set to is processing and remove the log content in UI and also in GITCOMMIT in memory
		m.PopUpModel.(*GitCommitPopUpModel).IsProcessing = true
		m.PopUpModel.(*GitCommitPopUpModel).GitCommitOutputViewport.SetContent("")
		git.GITCOMMIT.ClearGitCommitOutput()
		// retrieve the value of commit message and desc
		message := m.PopUpModel.(*GitCommitPopUpModel).MessageTextInput.Value()
		description := m.PopUpModel.(*GitCommitPopUpModel).DescriptionTextAreaInput.Value()
		if len(message) < 1 {
			return
		}
		// stage the changes based on what was chosen by user
		git.GITCOMMIT.GitStage()
		// and commit it
		exitStatusCode := git.GITCOMMIT.GitCommit(message, description)
		m.PopUpModel.(*GitCommitPopUpModel).IsProcessing = false // update the processing status

		// if sucessful exitcode will be 0
		if exitStatusCode == 0 {
			m.PopUpModel.(*GitCommitPopUpModel).MessageTextInput.Reset()
			m.PopUpModel.(*GitCommitPopUpModel).DescriptionTextAreaInput.Reset()
			git.GITCOMMIT.KillCommit()                                                 // kill the cmd process if exist
			time.Sleep(2 * time.Second)                                                // a 1 seconds delay before closing the pop up
			git.GITCOMMIT.ClearGitCommitOutput()                                       // clear the git commit output log
			m.ShowPopUp = false                                                        // close the pop up
			m.IsTyping = false                                                         // and reset typing mode
			m.PopUpModel.(*GitCommitPopUpModel).GitCommitOutputViewport.SetContent("") // set the git commit output viewport to nothing
		}
	}()

}

func gitCommitCancelService(m *GittiModel) {
	go func() {
		git.GITCOMMIT.KillCommit()                                                 // kill the cmd process if exist
		git.GITCOMMIT.ClearGitCommitOutput()                                       // clear the git commit output log
		m.PopUpModel.(*GitCommitPopUpModel).IsProcessing = false                   // mark the status to not processing
		m.ShowPopUp = false                                                        // close the pop up
		m.IsTyping = false                                                         // and reset typing mode
		m.PopUpModel.(*GitCommitPopUpModel).GitCommitOutputViewport.SetContent("") // set the git commit output viewport to nothing
	}()
}
