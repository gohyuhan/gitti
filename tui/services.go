package tui

import (
	"context"
	"fmt"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
)

// service.go was to bridge api and the needs of the terminal interface logic so that it can be compatible and feels smooth and not cluncy

// ------------------------------------
//
//	For Git Commit
//
// ------------------------------------
func gitCommitService(m *GittiModel, isAmendCommit bool) {
	ctx, cancel := context.WithCancel(context.Background())

	popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
	if ok {
		popUp.CancelFunc = cancel
	}

	go func(ctx context.Context) {
		defer cancel()

		// set to is processing and remove the log content in UI and also in GITCOMMIT in memory
		popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
		var message string
		var description string
		var exitStatusCode int
		if ok {
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
			popUp.IsProcessing.Store(false)
			return
		}
		// stage the changes based on what was chosen by user and commit it
		exitStatusCode = m.GitOperations.GitCommit.GitCommit(ctx, message, description, isAmendCommit)

		popUp, ok = m.PopUpModel.(*GitCommitPopUpModel)
		if ok && !popUp.IsCancelled.Load() {
			popUp.IsProcessing.Store(false) // update the processing status
			// if sucessful exitcode will be 0
			if exitStatusCode == 0 && !popUp.IsProcessing.Load() {
				popUp.ProcessSuccess.Store(true)
				popUp.MessageTextInput.Reset()
				popUp.DescriptionTextAreaInput.Reset()
				popUp.IsProcessing.Store(false)
				popUp.HasError.Store(false)
			} else if exitStatusCode != 0 && !popUp.IsProcessing.Load() {
				popUp.HasError.Store(true)
			}
		}
	}(ctx)
}

func gitCommitCancelService(m *GittiModel) {
	popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
	if ok {
		popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		if popUp.CancelFunc != nil {
			popUp.CancelFunc() // Cancel the context, which terminates the command and goroutine
		}
	}
	m.GitOperations.GitCommit.ClearGitCommitOutput() // clear the git commit output log

	m.ShowPopUp.Store(false) // close the pop up
	m.IsTyping.Store(false)  // reset typing mode
	m.PopUpType = constant.NoPopUp
	if ok {
		popUp.GitCommitOutputViewport.SetContent("") // set the git commit output viewport to nothing
		popUp.IsProcessing.Store(false)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
	}
}

// ------------------------------------
//
//	For Git Amend Commit
//
// ------------------------------------
func gitAmendCommitService(m *GittiModel, isAmendCommit bool) {
	ctx, cancel := context.WithCancel(context.Background())

	popUp, ok := m.PopUpModel.(*GitAmendCommitPopUpModel)
	if ok {
		popUp.CancelFunc = cancel
	}

	go func(ctx context.Context) {
		defer cancel()

		// set to is processing and remove the log content in UI and also in GITCOMMIT in memory
		popUp, ok := m.PopUpModel.(*GitAmendCommitPopUpModel)
		var message string
		var description string
		var exitStatusCode int
		if ok {
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
			popUp.IsProcessing.Store(false)
			return
		}
		// stage the changes based on what was chosen by user and commit it
		exitStatusCode = m.GitOperations.GitCommit.GitCommit(ctx, message, description, isAmendCommit)

		popUp, ok = m.PopUpModel.(*GitAmendCommitPopUpModel)
		if ok && !popUp.IsCancelled.Load() {
			popUp.IsProcessing.Store(false) // update the processing status
			// if sucessful exitcode will be 0
			if exitStatusCode == 0 && !popUp.IsProcessing.Load() {
				popUp.ProcessSuccess.Store(true)
				popUp.MessageTextInput.Reset()
				popUp.DescriptionTextAreaInput.Reset()
				popUp.IsProcessing.Store(false)
				popUp.HasError.Store(false)
			} else if exitStatusCode != 0 && !popUp.IsProcessing.Load() {
				popUp.HasError.Store(true)
			}
		}
	}(ctx)
}

func gitAmendCommitCancelService(m *GittiModel) {
	popUp, ok := m.PopUpModel.(*GitAmendCommitPopUpModel)
	if ok {
		popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		if popUp.CancelFunc != nil {
			popUp.CancelFunc() // Cancel the context, which terminates the command and goroutine
		}
	}
	m.GitOperations.GitCommit.ClearGitCommitOutput() // clear the git commit output log

	m.ShowPopUp.Store(false) // close the pop up
	m.IsTyping.Store(false)  // reset typing mode
	m.PopUpType = constant.NoPopUp
	if ok {
		popUp.GitAmendCommitOutputViewport.SetContent("") // set the git commit output viewport to nothing
		popUp.IsProcessing.Store(false)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
	}
}

// ------------------------------------
//
//	For Adding Git Remote
//
// ------------------------------------
func gitAddRemoteService(m *GittiModel) {
	ctx, cancel := context.WithCancel(context.Background())

	popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
	if ok {
		popUp.CancelFunc = cancel
	}

	go func(ctx context.Context) {
		defer cancel()

		// set to is processing and remove the log content in UI and also in GITCOMMIT in memory
		popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
		var remoteName string
		var remoteUrl string
		var exitStatusCode int
		if ok {
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
		gitAddRemoteResult, exitStatusCode := m.GitOperations.GitRemote.GitAddRemote(ctx, remoteName, remoteUrl)
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
				popUp.HasError.Store(false)
				popUp.ProcessSuccess.Store(true)
			} else if exitStatusCode != 0 && !popUp.IsProcessing.Load() {
				popUp.HasError.Store(true)
				updateAddRemoteOutputViewport(m, gitAddRemoteResult)
			}
		}
	}(ctx)
}

func gitAddRemoteCancelService(m *GittiModel) {
	popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
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
//	For Git Remote Push
//
// ------------------------------------
func gitRemotePushService(m *GittiModel, remoteName string, pushType string) {
	ctx, cancel := context.WithCancel(context.Background())

	popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
	if ok {
		popUp.CancelFunc = cancel
	}

	go func(ctx context.Context) {
		defer cancel()

		// set to is processing and remove the log content in UI and also in GITCOMMIT in memory
		popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
		var exitStatusCode int
		if ok {
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
			popUp.IsProcessing.Store(true)
			popUp.IsCancelled.Store(false)
		} else {
			return
		}
		exitStatusCode = m.GitOperations.GitCommit.GitPush(ctx, remoteName, pushType, m.CheckOutBranch)
		popUp, ok = m.PopUpModel.(*GitRemotePushPopUpModel)
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

func gitRemotePushCancelService(m *GittiModel) {
	popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
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
//	For Git Switching brnach ( only switch or switch while bringing changes )
//
// ------------------------------------
func gitSwitchBranchService(m *GittiModel, branchName string, switchType string) {
	go func() {
		popUp, ok := m.PopUpModel.(*SwitchBranchOutputPopUpModel)

		if ok {
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
			popUp.IsProcessing.Store(true)
		} else {
			return
		}

		var gitOpsOutput []string
		var success bool
		switch switchType {
		case git.SWITCHBRANCH:
			gitOpsOutput, success = m.GitOperations.GitBranch.GitSwitchBranch(branchName)
		case git.SWITCHBRANCHWITHCHANGES:
			gitOpsOutput, success = m.GitOperations.GitBranch.GitSwitchBranchWithChanges(branchName)
		}

		updateSwitchBranchOutputViewPort(m, gitOpsOutput)

		if success {
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(true)
			popUp.IsProcessing.Store(false)
		} else {
			popUp.HasError.Store(true)
			popUp.ProcessSuccess.Store(false)
			popUp.IsProcessing.Store(false)
		}
	}()
}

// ------------------------------------
//
//	For Git Pull
//
// ------------------------------------
func gitPullService(m *GittiModel, pullType string) {
	ctx, cancel := context.WithCancel(context.Background())

	popUp, ok := m.PopUpModel.(*GitPullOutputPopUpModel)
	if ok {
		popUp.CancelFunc = cancel
	}

	go func(ctx context.Context) {
		defer cancel()

		// set to is processing and remove the log content in UI and also in GITCOMMIT in memory
		popUp, ok := m.PopUpModel.(*GitPullOutputPopUpModel)
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
		popUp, ok = m.PopUpModel.(*GitPullOutputPopUpModel)
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

func gitPullCancelService(m *GittiModel) {
	popUp, ok := m.PopUpModel.(*GitPullOutputPopUpModel)
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

// ------------------------------------
//
//	For create new branch
//
// ------------------------------------
func gitCreateNewBranchService(m *GittiModel, validBranchName string) {
	go func() {
		if len(validBranchName) < 1 {
			return
		}
		m.GitOperations.GitBranch.GitCreateNewBranch(validBranchName)
	}()
}

// ------------------------------------
//
//	For create new branch and switch
//
// ------------------------------------
func gitCreateNewBranchAndSwitchService(m *GittiModel, validBranchName string) {
	go func() {
		m.GitOperations.GitBranch.GitCreateNewBranchAndSwitch(validBranchName)
	}()
}

// ------------------------------------
//
//	For Git Individual file stage or unstage
//
// ------------------------------------
func gitStageOrUnstageService(m *GittiModel, filePathName string) {
	go func() {
		m.GitOperations.GitFiles.StageOrUnstageFile(filePathName)
	}()
}

// ------------------------------------
//
//	For Git Stage All
//
// ------------------------------------
func gitStageAllChangesService(m *GittiModel) {
	go func() {
		m.GitOperations.GitFiles.StageAllChanges()
	}()
}

// ------------------------------------
//
//	For Git Unstage All
//
// ------------------------------------
func gitUnstageAllChangesService(m *GittiModel) {
	go func() {
		m.GitOperations.GitFiles.UnstageAllChanges()
	}()
}

// ------------------------------------
//
//	For Git stash all
//
// ------------------------------------
func gitStashAllService(m *GittiModel, msg string) {
	go func() {
		m.GitOperations.GitStash.GitStashAll(msg)
	}()
}

// ------------------------------------
//
//	For Git stash individual file
//
// ------------------------------------
func gitStashIndividualFileService(m *GittiModel, filePathName string, msg string) {
	go func() {
		m.GitOperations.GitStash.GitStashFile(filePathName, msg)
	}()
}

// ------------------------------------
//
//	For Git stash Apply
//
// ------------------------------------
func gitStashApplyService(m *GittiModel, filePathName string) {
	go func() {
		m.GitOperations.GitStash.GitStashApply(filePathName)
	}()
}

// ------------------------------------
//
//	For Git stash Pop
//
// ------------------------------------
func gitStashPopService(m *GittiModel, filePathName string) {
	go func() {
		m.GitOperations.GitStash.GitStashPop(filePathName)
	}()
}

// ------------------------------------
//
//	For Git stash drop
//
// ------------------------------------
func gitStashDropService(m *GittiModel, filePathName string) {
	go func() {
		m.GitOperations.GitStash.GitStashDrop(filePathName)
	}()
}

// ------------------------------------
//
//	For Git discard file changes
//
// ------------------------------------
func gitDiscardFileChangesService(m *GittiModel, filePathName string, discardType string) {
	go func() {
		m.GitOperations.GitFiles.DiscardFileChanges(filePathName, discardType)
	}()
}
