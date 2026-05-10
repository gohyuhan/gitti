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
	popUp, ok := m.PopUpModel.(*branchPopUp.SwitchBranchOutputPopUpModel)

	if ok {
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsProcessing.Store(true)
	} else {
		return
	}

	go func() {
		var gitOpsOutput []string
		var success bool
		switch switchType {
		case git.SWITCHBRANCH:
			gitOpsOutput, success = m.GitOperations.GitBranch.GitSwitchBranch(branchName)
		case git.SWITCHBRANCHWITHCHANGES:
			gitOpsOutput, success = m.GitOperations.GitBranch.GitSwitchBranchWithChanges(branchName)
		}

		data := types.GitSwitchBranchResultEventDataInterface{
			Result:  gitOpsOutput,
			Success: success,
		}
		m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
			Event: constant.GIT_SWITCH_BRANCH_RESULT_EVENT,
			Data:  data,
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
	popUp, ok := m.PopUpModel.(*branchPopUp.GitDeleteBranchOutputPopUpModel)
	if ok {
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsProcessing.Store(true)
	} else {
		return
	}

	go func() {
		result, success := m.GitOperations.GitBranch.DeleteLocalBranch(branchName)
		data := types.GitDeleteBranchResultEventDataInterface{
			Result:  result,
			Success: success,
		}
		m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
			Event: constant.GIT_DELETE_BRANCH_RESULT_EVENT,
			Data:  data,
		}
	}()
}

// ------------------------------------
//
//	For create new branch based on remote branch
//
// ------------------------------------
func CreateNewBranchBasedOnRemoteService(m *types.GittiModel, remoteName string, branchName string, newBranchCreateType string) {
	popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemoteOutputPopUpModel)
	if ok {
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsProcessing.Store(true)
	} else {
		return
	}

	go func() {
		if newBranchCreateType == git.NEWBRANCHBASEDONREMOTEUSERSELECT {
			parts := strings.SplitN(branchName, "/", 2)
			if len(parts) < 2 {
				m.GittiLogger.RegisterNewLog(logging.RETRIEVE_LATEST_REMOTE_BRANCH_INFO, "", logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.RETRIEVE_LATEST_REMOTE_BRANCH_INFO, "Invalid Remote Branch Naming"), false)
				data := types.GitCreateNewBranchBasedOnRemoteInvalidEventDataInterface{
					RemoteName: remoteName,
					BranchName: branchName,
				}
				m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
					Event: constant.GIT_CREATE_NEW_BRANCH_BASED_ON_REMOTE_INVALID_EVENT,
					Data:  data,
				}
				return
			}
			remoteName = parts[0]
			branchName = parts[1]
		}
		result, success := m.GitOperations.GitBranch.GitCreateNewBranchBasedOnRemote(remoteName, branchName)
		data := types.GitCreateNewBranchBasedOnRemoteResultEventDataInterface{
			Result:  result,
			Success: success,
		}
		m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
			Event: constant.GIT_CREATE_NEW_BRANCH_BASED_ON_REMOTE_RESULT_EVENT,
			Data:  data,
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
		popUp.IsProcessing.Store(true)
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsCancelled.Store(false)
		popUp.CancelFunc = cancel
	}

	go func(ctx context.Context) {
		defer cancel()

		outputResult, success := m.GitOperations.GitBranch.GitMerge(ctx, branchesName)
		data := types.MergeResultEventDataInterface{
			Result:  outputResult,
			Success: success,
		}
		m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
			Event: constant.GIT_MERGE_RESULT_EVENT,
			Data:  data,
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
