package services

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	"github.com/gohyuhan/gitti/tui/types"
)

// services was to bridge api and the needs of the terminal interface logic so that it can be compatible and feels smooth and not clunky
// ------------------------------------
//
//	For Git Switching branch ( only switch or switch while bringing changes )
//
// ------------------------------------
func GitSwitchBranchService(m *types.GittiModel, branchName string, switchType string) {
	go func() {
		popUp, ok := m.PopUpModel.(*branchPopUp.SwitchBranchOutputPopUpModel)

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

		branchPopUp.UpdateSwitchBranchOutputViewPort(m, gitOpsOutput)

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
//	For create new branch
//
// ------------------------------------
func GitCreateNewBranchService(m *types.GittiModel, validBranchName string) {
	go func() {
		if utf8.RuneCountInString(validBranchName) < 1 {
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
func GitCreateNewBranchAndSwitchService(m *types.GittiModel, validBranchName string) {
	go func() {
		m.GitOperations.GitBranch.GitCreateNewBranchAndSwitch(validBranchName)
	}()
}

// ------------------------------------
//
//	For branch delete
//
// ------------------------------------
func GitDeleteBranchService(m *types.GittiModel, branchName string) {
	go func() {
		result, success := m.GitOperations.GitBranch.DeleteLocalBranch(branchName)
		popUp, ok := m.PopUpModel.(*branchPopUp.GitDeleteBranchOutputPopUpModel)
		if ok {
			if success {
				popUp.HasError.Store(false)
				popUp.ProcessSuccess.Store(true)
			} else {
				popUp.HasError.Store(true)
				popUp.ProcessSuccess.Store(false)
			}
			popUp.IsProcessing.Store(false)
			popUp.BranchDeleteOutputViewport.SetContentLines(result)
			popUp.BranchDeleteOutputViewport.PageDown()
		}
	}()
}

// ------------------------------------
//
//	For create new branch based on remote branch
//
// ------------------------------------
func CreateNewBranchBasedOnRemoteService(m *types.GittiModel, remoteName string, branchName string, newBranchCreateType string) {
	go func() {
		if newBranchCreateType == git.NEWBRANCHBASEDONREMOTEUSERSELECT {
			parts := strings.SplitN(branchName, "/", 2)
			if len(parts) < 2 {
				m.GittiLogger.RegisterNewLog(logging.RETRIEVE_LATEST_REMOTE_BRANCH_INFO, "", logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.RETRIEVE_LATEST_REMOTE_BRANCH_INFO, "Invalid Remote Branch Naming"), false)
				m.IsTyping.Store(false)
				m.ShowPopUp.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil

				return
			}
			remoteName = parts[0]
			branchName = parts[1]
		}
		result, success := m.GitOperations.GitBranch.GitCreateNewBranchBasedOnRemote(remoteName, branchName)
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemoteOutputPopUpModel)
		if ok {
			if success {
				popUp.HasError.Store(false)
				popUp.ProcessSuccess.Store(true)
			} else {
				popUp.HasError.Store(true)
				popUp.ProcessSuccess.Store(false)
			}
			popUp.IsProcessing.Store(false)
			popUp.CreateBranchBasedOnRemoteOutputViewport.SetContentLines(result)
			popUp.CreateBranchBasedOnRemoteOutputViewport.PageDown()
		}
	}()
}

// ------------------------------------
//
//	For create new branch based on commit hash (trigger from reflog component panel)
//
// ------------------------------------
func GitCreateNewBranchBasedOnCommitHashService(m *types.GittiModel, validBranchName string, commitHash string) {
	go func() {
		if utf8.RuneCountInString(validBranchName) < 1 || utf8.RuneCountInString(commitHash) < 1 {
			return
		}
		m.GitOperations.GitBranch.GitCreateNewBranchBasedOnCommitHash(validBranchName, commitHash)
	}()
}

// ------------------------------------
//
//	For Git Merge
//
// ------------------------------------
func GitMergeService(m *types.GittiModel, branchesName []string) {
	ctx, cancel := context.WithCancel(context.Background())

	popUp, ok := m.PopUpModel.(*branchPopUp.BranchMergeOutputPopUpModel)
	if ok {
		popUp.CancelFunc = cancel
	}

	go func(ctx context.Context) {
		defer cancel()

		// set to is processing and remove the log content in UI and also in GITCOMMIT in memory
		popUp, ok := m.PopUpModel.(*branchPopUp.BranchMergeOutputPopUpModel)
		if ok {
			popUp.HasError.Store(false)
			popUp.ProcessSuccess.Store(false)
			popUp.IsProcessing.Store(true)
			popUp.IsCancelled.Store(false)
		} else {
			return
		}
		outputResult, success := m.GitOperations.GitBranch.GitMerge(ctx, branchesName)
		popUp, ok = m.PopUpModel.(*branchPopUp.BranchMergeOutputPopUpModel)
		if ok && !popUp.IsCancelled.Load() {
			popUp.IsProcessing.Store(false) // update the processing status
			if success && !popUp.IsProcessing.Load() {
				popUp.ProcessSuccess.Store(true)
				popUp.HasError.Store(false)
			} else if !success && !popUp.IsProcessing.Load() {
				popUp.ProcessSuccess.Store(false)
				popUp.HasError.Store(true)
			}

			var content strings.Builder
			for index := range outputResult {
				content.WriteString(outputResult[index])
				content.WriteRune('\n')
			}

			popUp.BranchMergeOutputViewport.SetContent(content.String())
		}
	}(ctx)
}

// ------------------------------------
//
//	Cancel the current git merge operation and clean up pop-up state
//
// ------------------------------------
func GitMergeCancelService(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*branchPopUp.BranchMergeOutputPopUpModel)
	if ok {
		popUp.IsCancelled.Store(true) // set cancellation flag first to prevent race condition
		if popUp.CancelFunc != nil {
			popUp.CancelFunc() // Cancel the context, which terminates the command and goroutine
		}
	}
	m.ShowPopUp.Store(false) // close the pop up
	m.IsTyping.Store(false)  // and reset typing mode
	m.PopUpType = constant.NoPopUp
	if ok {
		popUp.BranchMergeOutputViewport.SetContent("") // set the git rebase output viewport to nothing
		popUp.IsProcessing.Store(false)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
	}
}
