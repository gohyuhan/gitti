package services

import (
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/constant"
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
	popUp, ok := m.PopUpModel.(*stashPopUp.GitStashOperationOutputPopUpModel)
	if ok {
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(false)
		popUp.IsProcessing.Store(true)
	} else {
		return
	}

	go func() {
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

		data := types.GitStashOperationResultEventDataInterface{
			Result:  resultOutput,
			Success: exitStatusCode == 0,
		}
		m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
			Event: constant.GIT_STASH_OPERATION_RESULT_EVENT,
			Data:  data,
		}
	}()
}
