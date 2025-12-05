package tui

import (
	"context"
	"fmt"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// service.go was to bridge api and the needs of the terminal interface logic so that it can be compatible and feels smooth and not clunky
// ------------------------------------
//
//	For Git Commit
//
// ------------------------------------
func gitCommitService(m *types.GittiModel, isAmendCommit bool) {
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

func gitCommitCancelService(m *types.GittiModel) {
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
func gitAmendCommitService(m *types.GittiModel, isAmendCommit bool) {
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

func gitAmendCommitCancelService(m *types.GittiModel) {
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
func gitAddRemoteService(m *types.GittiModel) {
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

func gitAddRemoteCancelService(m *types.GittiModel) {
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
func gitRemotePushService(m *types.GittiModel, remoteName string, pushType string) {
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

func gitRemotePushCancelService(m *types.GittiModel) {
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
func gitSwitchBranchService(m *types.GittiModel, branchName string, switchType string) {
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
func gitPullService(m *types.GittiModel, pullType string) {
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

func gitPullCancelService(m *types.GittiModel) {
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
func gitCreateNewBranchService(m *types.GittiModel, validBranchName string) {
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
func gitCreateNewBranchAndSwitchService(m *types.GittiModel, validBranchName string) {
	go func() {
		m.GitOperations.GitBranch.GitCreateNewBranchAndSwitch(validBranchName)
	}()
}

// ------------------------------------
//
//	For Git Individual file stage or unstage
//
// ------------------------------------
func gitStageOrUnstageService(m *types.GittiModel, filePathName string) {
	go func() {
		m.GitOperations.GitFiles.StageOrUnstageFile(filePathName)
	}()
}

// ------------------------------------
//
//	For Git Stage All
//
// ------------------------------------
func gitStageAllChangesService(m *types.GittiModel) {
	go func() {
		m.GitOperations.GitFiles.StageAllChanges()
	}()
}

// ------------------------------------
//
//	For Git Unstage All
//
// ------------------------------------
func gitUnstageAllChangesService(m *types.GittiModel) {
	go func() {
		m.GitOperations.GitFiles.UnstageAllChanges()
	}()
}

// ------------------------------------
//
//	For Git stash operation
//	* Stash operations are not cancellable in Gitti, because interrupting
//	  the process mid-operation could leave the repository in a partial or
//	  inconsistent state (stash applied only halfway).
//
// ------------------------------------
func gitStashOperationService(m *types.GittiModel, filePathName string, stashId string, stashMessage string) {
	go func() {
		popUp, ok := m.PopUpModel.(*GitStashOperationOutputPopUpModel)
		if ok {
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
			popUp.IsProcessing.Store(true)
		} else {
			return
		}

		var exitStatusCode int
		var resultOutput []string

		switch popUp.StashOperationType {
		case git.STASHALL:
			resultOutput, exitStatusCode = m.GitOperations.GitStash.GitStashAll(stashMessage)
		case git.STASHFILE:
			resultOutput, exitStatusCode = m.GitOperations.GitStash.GitStashFile(filePathName, stashMessage)
		case git.APPLYSTASH:
			resultOutput, exitStatusCode = m.GitOperations.GitStash.GitStashApply(stashId)
		case git.POPSTASH:
			resultOutput, exitStatusCode = m.GitOperations.GitStash.GitStashPop(stashId)
		case git.DROPSTASH:
			resultOutput, exitStatusCode = m.GitOperations.GitStash.GitStashDrop(stashId)
		}

		popUp, ok = m.PopUpModel.(*GitStashOperationOutputPopUpModel)
		if ok {
			popUp.IsProcessing.Store(false) // update the processing status
			// if sucessful exitcode will be 0
			if exitStatusCode == 0 && !popUp.IsProcessing.Load() {
				popUp.ProcessSuccess.Store(true)
				popUp.HasError.Store(false)
			} else if exitStatusCode != 0 && !popUp.IsProcessing.Load() {
				popUp.HasError.Store(true)
			}
			popUp.GitStashOperationOutputViewport.SetContentLines(resultOutput)
		}
	}()
}

// ------------------------------------
//
//	For Git discard file changes
//
// ------------------------------------
func gitDiscardFileChangesService(m *types.GittiModel, filePathName string, discardType string) {
	go func() {
		m.GitOperations.GitFiles.DiscardFileChanges(filePathName, discardType)
	}()
}

// ------------------------------------
//
//	For resolving git file conflict
//
// ------------------------------------
func gitResolveConflictService(m *types.GittiModel, filePathName string, resolveType string) {
	go func() {
		m.GitOperations.GitFiles.GitResolveConflict(filePathName, resolveType)
	}()
}

// ------------------------------------
//
//		For fetching detail component panel info
//	  * it can be for stash info, commit info etc
//
// ------------------------------------
func fetchDetailComponentPanelInfoService(m *types.GittiModel, reinit bool) {
	// For non-reinit calls (refreshing current view), abort if already processing.
	// This avoids looping a cancel and execution cycle which would end up blocking
	// a slightly longer processing process.
	//
	// If not processing, we proceed to fetch to ensure we capture any updates (e.g., file changes,
	// amends), as we lack specific context on whether the underlying data has changed.
	//
	// If `reinit` is true (context switch), we bypass this check to cancel the active fetch
	// and start the new one immediately.
	if !reinit && m.IsDetailComponentPanelInfoFetchProcessing.Load() {
		return
	}

	// Cancel any existing operation first
	if m.DetailComponentPanelInfoFetchCancelFunc != nil {
		m.DetailComponentPanelInfoFetchCancelFunc()
	}

	// Wait for the previous goroutine to finish (its defer will set processing to false),
	// then atomically set it to true before starting a new one.
	for !m.IsDetailComponentPanelInfoFetchProcessing.CompareAndSwap(false, true) {
		// The previous goroutine is still running, wait a tiny bit
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.DetailComponentPanelInfoFetchCancelFunc = cancel
	go func(ctx context.Context) {
		defer func() {
			m.IsDetailComponentPanelInfoFetchProcessing.Store(false)
			cancel()
		}()

		var contentLine string
		var theCurrentSelectedComponent string
		// reinit and render detail component panel viewport
		if reinit {
			m.DetailPanelViewport.SetContent(style.NewStyle.Render(i18n.LANGUAGEMAPPING.Loading))
			m.DetailPanelViewportOffset = 0
			m.DetailPanelViewport.SetXOffset(0)
			m.DetailPanelViewport.SetYOffset(0)
		}
		if m.CurrentSelectedComponent == constant.DetailComponent {
			// if the current selected one is the detail component itself, the current selected one will be its parent (the component that led into the detail component)
			theCurrentSelectedComponent = m.DetailPanelParentComponent
		} else {
			theCurrentSelectedComponent = m.CurrentSelectedComponent
		}
		switch theCurrentSelectedComponent {
		case constant.ModifiedFilesComponent:
			contentLine = generateModifiedFileDetailPanelContent(ctx, m)
		case constant.StashComponent:
			contentLine = generateStashDetailPanelContent(ctx, m)
		default:
			contentLine = generateAboutGittiContent()
		}

		select {
		case <-ctx.Done():
			return
		default:
			if contentLine == "" {
				// if the content will be empty, render about gitti for detail panel
				contentLine = generateAboutGittiContent()
			}

			m.DetailPanelViewport.SetContent(contentLine)
			m.TuiUpdateChannel <- constant.DETAIL_COMPONENT_PANEL_UPDATED
			return
		}
	}(ctx)
}

// for modified file detail panel view
func generateModifiedFileDetailPanelContent(ctx context.Context, m *types.GittiModel) string {
	currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
	var fileStatus git.FileStatus
	if currentSelectedModifiedFile != nil {
		fileStatus = git.FileStatus(currentSelectedModifiedFile.(gitModifiedFilesItem))
	} else {
		return ""
	}

	vpLine := fmt.Sprintf("[ %s ]\n\n", fileStatus.FilePathname)
	fileDiffLines := m.GitOperations.GitFiles.GetFilesDiffInfo(ctx, fileStatus)
	if fileDiffLines == nil {
		vpLine += i18n.LANGUAGEMAPPING.FileTypeUnSupportedPreview
		return vpLine
	}
	for _, line := range fileDiffLines {
		line = style.NewStyle.Render(line)
		vpLine += line + "\n"
	}
	return vpLine
}

// for stash detail panel view
func generateStashDetailPanelContent(ctx context.Context, m *types.GittiModel) string {
	currentSelectedStash := m.CurrentRepoStashInfoList.SelectedItem()
	var stash gitStashItem
	if currentSelectedStash != nil {
		stash = currentSelectedStash.(gitStashItem)
	} else {
		return ""
	}

	vpLine := fmt.Sprintf(
		"[%s]\n[%s]\n\n",
		style.StashIdStyle.Render(stash.Id),
		style.StashMessageStyle.Render(stash.Message),
	)

	stashDetail := m.GitOperations.GitStash.GitStashDetail(ctx, stash.Id)
	if len(stashDetail) < 1 {
		return ""
	}

	for _, Line := range stashDetail {
		line := style.NewStyle.Render(Line)
		vpLine += line + "\n"
	}
	return vpLine
}
