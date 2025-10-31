package tui

import (
	"gitti/api/git"
	"time"
	// "time"
)

// service.go was to bridge api and the needs of the terminal interface logic so that it can be compatible and feels smooth and not cluncy

func gitCommitService(m *GittiModel) {
	go func() {
		m.PopUpModel.(*GitCommitPopUpModel).IsProcessing = true
		message := m.PopUpModel.(*GitCommitPopUpModel).MessageTextInput.Value()
		description := m.PopUpModel.(*GitCommitPopUpModel).DescriptionTextAreaInput.Value()
		if len(message) < 1 {
			return
		}
		git.GITCOMMIT.GitStage()
		exitStatusCode := git.GITCOMMIT.GitCommit(message, description)
		m.PopUpModel.(*GitCommitPopUpModel).IsProcessing = false // update the processing status
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
