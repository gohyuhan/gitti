package services

import (
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/types"
)

// services was to bridge api and the needs of the terminal interface logic so that it can be compatible and feels smooth and not clunky
// ------------------------------------
//
//	For Git stash operation
//	* Stash operations are not cancellable in Gitti, because interrupting
//	  the process mid-operation could leave the repository in a partial or
//	  inconsistent state (stash applied only halfway).
//
// ------------------------------------
func GitStashOperationService(m *types.GittiModel, filePathName string, stashId string, stashMessage string) {
	go func() {
		popUp, ok := m.PopUpModel.(*stashPopUp.GitStashOperationOutputPopUpModel)
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

		popUp, ok = m.PopUpModel.(*stashPopUp.GitStashOperationOutputPopUpModel)
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
