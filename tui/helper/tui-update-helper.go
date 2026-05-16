package helper

import (
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/layout"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	pullPopUp "github.com/gohyuhan/gitti/tui/popup/pull"
	pushPopUp "github.com/gohyuhan/gitti/tui/popup/push"
	rebasePopUp "github.com/gohyuhan/gitti/tui/popup/rebase"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	Dispatch a GittiTuiUpdateMsg to the appropriate popup or layout update
//	function based on the event type constant in msg.Event. Covers detail panel
//	state/layout changes, branch/tag/stash/remote/rebase/commit/pull results,
//	and interactive-rebase events.
//
// ------------------------------------
func GittiTuiUpdateEventHelper(m *types.GittiModel, msg types.GittiTuiUpdateMsg) {
	updateMsg := msg
	updateEvent := updateMsg.Event
	switch updateEvent {
	case constant.DETAIL_COMPONENT_PANEL_LAYOUT_STATE_UPDATED_EVENT:
		layout.UpdateDetailComponentViewportContentAndState(m, updateMsg.Data.(types.DetailPanelStateAndLayoutUpdateEventDataStructure))
		layout.UpdateDetailComponentViewportLayout(m)
	case constant.DETAIL_COMPONENT_PANEL_LAYOUT_UPDATED_EVENT:
		layout.UpdateDetailComponentViewportLayout(m)
	case constant.DETAIL_COMPONENT_PANEL_LAYOUT_STATE_REINIT_EVENT:
		layout.DetailComponentReinit(m)
	case constant.GIT_SWITCH_BRANCH_RESULT_EVENT:
		branchPopUp.UpdateSwitchBranchResultEvent(m, updateMsg.Data.(types.GitSwitchBranchResultEventDataStructure))
	case constant.GIT_DELETE_BRANCH_RESULT_EVENT:
		branchPopUp.UpdateDeleteBranchResultEvent(m, updateMsg.Data.(types.GitDeleteBranchResultEventDataStructure))
	case constant.GIT_CREATE_NEW_BRANCH_BASED_ON_REMOTE_RESULT_EVENT:
		branchPopUp.UpdateCreateNewBranchBasedOnRemoteResultEvent(m, updateMsg.Data.(types.GitCreateNewBranchBasedOnRemoteResultEventDataStructure))
	case constant.GIT_CREATE_NEW_BRANCH_BASED_ON_REMOTE_INVALID_EVENT:
		branchPopUp.UpdateCreateNewBranchBasedOnRemoteInvalidEvent(m, updateMsg.Data.(types.GitCreateNewBranchBasedOnRemoteInvalidEventDataStructure))
	case constant.GIT_MERGE_RESULT_EVENT:
		branchPopUp.UpdateMergeViewport(m, updateMsg.Data.(types.MergeResultEventDataStructure))
	case constant.GIT_DELETE_TAG_RESULT_EVENT:
		tagPopUp.UpdateDeleteTagResultEvent(m, updateMsg.Data.(types.GitDeleteTagResultEventDataStructure))
	case constant.GIT_PUSH_TAG_RESULT_EVENT:
		tagPopUp.UpdatePushTagResultEvent(m, updateMsg.Data.(types.GitPushTagResultEventDataStructure))
	case constant.GIT_FETCH_TAG_RESULT_EVENT:
		tagPopUp.UpdateFetchTagResultEvent(m, updateMsg.Data.(types.GitFetchTagResultEventDataStructure))
	case constant.GIT_STASH_OPERATION_RESULT_EVENT:
		stashPopUp.UpdateGitStashOperationResultEvent(m, updateMsg.Data.(types.GitStashOperationResultEventDataStructure))
	case constant.GIT_ADD_REMOTE_RESULT_EVENT:
		remotePopUp.UpdateAddRemoteResultEvent(m, updateMsg.Data.(types.GitAddRemoteResultEventDataStructure))
	case constant.GIT_REBASE_RESULT_EVENT:
		rebasePopUp.UpdateGitRebaseResultEvent(m, updateMsg.Data.(types.GitRebaseResultEventDataStructure))
	case constant.GIT_PUSH_RESULT_EVENT:
		pushPopUp.UpdateGitPushResultEvent(m, updateMsg.Data.(types.GitPushResultEventDataStructure))
	case constant.GIT_COMMIT_RESULT_EVENT:
		commitPopUp.UpdateGitCommitResultEvent(m, updateMsg.Data.(types.GitCommitResultEventDataStructure))
	case constant.GIT_AMEND_COMMIT_RESULT_EVENT:
		commitPopUp.UpdateGitAmendCommitResultEvent(m, updateMsg.Data.(types.GitAmendCommitResultEventDataStructure))
	case constant.GIT_PULL_RESULT_EVENT:
		pullPopUp.UpdateGitPullResultEvent(m, updateMsg.Data.(types.GitPullResultEventDataStructure))
	case constant.INTERACTIVE_REBASE_FIXUP_SQUASH_RESULT_EVENT:
		interactiverebasePopUp.UpdateInteractiveRebaseFixupSquashResultEvent(m, updateMsg.Data.(types.InteractiveRebaseFixupSquashResultEventDataStructure))
	case constant.INTERACTIVE_REBASE_FIXUP_SQUASH_FETCH_COMMITS_INFO_EVENT:
		interactiverebasePopUp.UpdateInteractiveRebaseFixupSquashFetchedCommitInfoList(m, updateMsg.Data.(types.InteractiveRebaseFetchCommitInfoListEventDataStructure))
	case constant.INTERACTIVE_REBASE_REWORD_FETCH_COMMITS_INFO_EVENT:
		interactiverebasePopUp.UpdateInteractiveRebaseRewordFetchedCommitInfoList(m, updateMsg.Data.(types.InteractiveRebaseFetchCommitInfoListEventDataStructure))
	case constant.REINIT_CHERRY_PICKED_COMMIT_INFO_EVENT:
		utils.ReinitCherryPickedCommitInfo(m)
	}
}
