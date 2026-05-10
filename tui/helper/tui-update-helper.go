package helper

import (
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/layout"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	"github.com/gohyuhan/gitti/tui/types"
)

func GittiTuiUpdateEventHelper(m *types.GittiModel, msg types.GittiTuiUpdateMsg) {
	updateMsg := types.GittiTuiUpdateMsg(msg)
	updateEvent := updateMsg.Event
	switch updateEvent {
	case constant.DETAIL_COMPONENT_PANEL_LAYOUT_STATE_UPDATED_EVENT:
		layout.UpdateDetailComponentViewportContentAndState(m, updateMsg.Data.(types.DetailPanelStateAndLayoutUpdateEventDataInterface))
		layout.UpdateDetailComponentViewportLayout(m)
	case constant.DETAIL_COMPONENT_PANEL_LAYOUT_UPDATED_EVENT:
		layout.UpdateDetailComponentViewportLayout(m)
	case constant.DETAIL_COMPONENT_PANEL_LAYOUT_STATE_REINIT_EVENT:
		layout.DetailComponentReinit(m)
	case constant.GIT_SWITCH_BRANCH_RESULT_EVENT:
		branchPopUp.UpdateSwitchBranchResultEvent(m, updateMsg.Data.(types.GitSwitchBranchResultEventDataInterface))
	case constant.GIT_DELETE_BRANCH_RESULT_EVENT:
		branchPopUp.UpdateDeleteBranchResultEvent(m, updateMsg.Data.(types.GitDeleteBranchResultEventDataInterface))
	case constant.GIT_CREATE_NEW_BRANCH_BASED_ON_REMOTE_RESULT_EVENT:
		branchPopUp.UpdateCreateNewBranchBasedOnRemoteResultEvent(m, updateMsg.Data.(types.GitCreateNewBranchBasedOnRemoteResultEventDataInterface))
	case constant.GIT_CREATE_NEW_BRANCH_BASED_ON_REMOTE_INVALID_EVENT:
		branchPopUp.UpdateCreateNewBranchBasedOnRemoteInvalidEvent(m, updateMsg.Data.(types.GitCreateNewBranchBasedOnRemoteInvalidEventDataInterface))
	case constant.GIT_MERGE_RESULT_EVENT:
		branchPopUp.UpdateMergeViewport(m, updateMsg.Data.(types.MergeResultEventDataInterface))
	case constant.INTERACTIVE_REBASE_FIXUP_SQUASH_RESULT_EVENT:
		interactiverebasePopUp.UpdateInteractiveRebaseFixupSquashResultEvent(m, updateMsg.Data.(types.InteractiveRebaseFixupSquashResultEventDataInterface))
	case constant.INTERACTIVE_REBASE_FETCH_COMMITS_INFO_EVENT:
		interactiverebasePopUp.UpdateInteractiveRebaseFetchedCommitInfoList(m, updateMsg.Data.(types.InteractiveRebaseFetchCommitInfoListEventDataInterface))
	}
}
