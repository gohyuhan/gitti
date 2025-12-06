package services

import (
	"github.com/gohyuhan/gitti/api/git"
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
func GitCreateNewBranchAndSwitchService(m *types.GittiModel, validBranchName string) {
	go func() {
		m.GitOperations.GitBranch.GitCreateNewBranchAndSwitch(validBranchName)
	}()
}
