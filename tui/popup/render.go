package popup

import (
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/popup/blame"
	"github.com/gohyuhan/gitti/tui/popup/branch"
	"github.com/gohyuhan/gitti/tui/popup/commit"
	"github.com/gohyuhan/gitti/tui/popup/commitlog"
	"github.com/gohyuhan/gitti/tui/popup/discard"
	"github.com/gohyuhan/gitti/tui/popup/files"
	interactiverebase "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	"github.com/gohyuhan/gitti/tui/popup/keybinding"
	"github.com/gohyuhan/gitti/tui/popup/pull"
	"github.com/gohyuhan/gitti/tui/popup/push"
	"github.com/gohyuhan/gitti/tui/popup/rebase"
	"github.com/gohyuhan/gitti/tui/popup/reflog"
	"github.com/gohyuhan/gitti/tui/popup/remote"
	"github.com/gohyuhan/gitti/tui/popup/resolve"
	"github.com/gohyuhan/gitti/tui/popup/stash"
	"github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Dispatch popup rendering to the appropriate package function based on
//	m.PopUpType. Returns an empty string when no popup type matches.
//
// ------------------------------------
func RenderPopUpComponent(m *types.GittiModel) string {
	var popUp string

	switch m.PopUpType {
	case constant.KeybindingAndFeatureInstructionsPopUp:
		popUp = keybinding.RenderKeyBindingAndFeatureInstructionsPopUp(m)
	case constant.CommitPopUp:
		popUp = commit.RenderGitCommitPopUp(m)
	case constant.AmendCommitPopUp:
		popUp = commit.RenderGitAmendCommitPopUp(m)
	case constant.AddRemotePromptPopUp:
		popUp = remote.RenderAddRemotePromptPopUp(m)
	case constant.GitRemotePushPopUp:
		popUp = push.RenderGitRemotePushPopUp(m)
	case constant.ChooseRemotePopUp:
		popUp = remote.RenderChooseRemotePopUp(m)
	case constant.ChoosePushTypePopUp:
		popUp = push.RenderChoosePushTypePopUp(m)
	case constant.ChooseNewBranchTypePopUp:
		popUp = branch.RenderChooseNewBranchTypePopUp(m)
	case constant.CreateNewBranchPopUp:
		popUp = branch.RenderCreateNewBranchPopUp(m)
	case constant.ChooseSwitchBranchTypePopUp:
		popUp = branch.RenderChooseSwitchBranchTypePopUp(m)
	case constant.SwitchBranchOutputPopUp:
		popUp = branch.RenderSwitchBranchOutputPopUp(m)
	case constant.ChooseGitPullTypePopUp:
		popUp = pull.RenderChooseGitPullTypePopUp(m)
	case constant.GitPullOutputPopUp:
		popUp = pull.RenderGitPullOutputPopUp(m)
	case constant.GitStashMessagePopUp:
		popUp = stash.RenderGitStashMessagePopUp(m)
	case constant.GitDiscardTypeOptionPopUp:
		popUp = discard.RenderGitDiscardTypeOptionPopUp(m)
	case constant.GitDiscardConfirmPromptPopUp:
		popUp = discard.RenderGitDiscardConfirmPromptPopup(m)
	case constant.GitStashOperationOutputPopUp:
		popUp = stash.RenderGitStashOperationOutputPopUp(m)
	case constant.GitStashConfirmPromptPopUp:
		popUp = stash.RenderGitStashConfirmPromptPopUp(m)
	case constant.GitResolveConflictOptionPopUp:
		popUp = resolve.RenderGitResolveConflictOptionPopUp(m)
	case constant.GitDeleteBranchConfirmPromptPopUp:
		popUp = branch.RenderGitDeleteBranchConfirmPromptPopUp(m)
	case constant.GitDeleteBranchOutputPopUp:
		popUp = branch.RenderGitDeleteBranchOutputPopUp(m)
	case constant.CreateBranchBasedOnRemotePopUp:
		popUp = branch.RenderCreateBranchBasedOnRemotePopUp(m)
	case constant.CreateBranchBasedOnRemoteOutputPopUp:
		popUp = branch.RenderCreateBranchBasedOnRemoteOutputPopUp(m)
	case constant.GitResetLatestCommitTypeOptionPopUp:
		popUp = commit.RenderGitResetLatestCommitTypeOptionPopUp(m)
	case constant.GitResetLatestCommitConfirmPromptPopUp:
		popUp = commit.RenderGitResetLatestCommitConfirmPromptPopUp(m)
	case constant.GitResetToSelectedCommitTypeOptionPopUp:
		popUp = commit.RenderGitResetToSelectedCommitTypeOptionPopUp(m)
	case constant.GitResetToSelectedCommitConfirmPromptPopUp:
		popUp = commit.RenderGitResetToSelectedCommitConfirmPromptPopUp(m)
	case constant.GitCherryPickOptionSelectionPopUp:
		popUp = commitlog.RenderGitCherryPickOptionSelectionPopUp(m)
	case constant.GitCherryPickPopUp:
		popUp = commitlog.RenderGitCherryPickPopUp(m)
	case constant.GitEditCherryPickPopUp:
		popUp = commitlog.RenderGitEditCherryPickPopUp(m)
	case constant.GitCherryPickApplyConfirmPopUp:
		popUp = commitlog.RenderGitCherryPickApplyConfirmPopUp(m)
	case constant.GitDiscardFileLineChangeConfirmPopUp:
		popUp = files.RenderGitDiscardFileLineChangeConfirmPopUp(m)
	case constant.CreateTagPopUp:
		popUp = tag.RenderCreateTagPopUp(m)
	case constant.CreateTagConfirmationPopUp:
		popUp = tag.RenderCreateTagConfirmationPopUp(m)
	case constant.ChooseDeleteTagOptionPopUp:
		popUp = tag.RenderChooseDeleteTagOptionPopUp(m)
	case constant.ChooseRemoteForDeleteRemoteTagPopUp:
		popUp = tag.RenderChooseRemoteForDeleteRemoteTagPopUp(m)
	case constant.DeleteTagOutputPopUp:
		popUp = tag.RenderDeleteTagOutputPopUp(m)
	case constant.ChoosePushTagOptionPopUp:
		popUp = tag.RenderPushTagOptionPopUp(m)
	case constant.PushTagOutputPopUp:
		popUp = tag.RenderPushTagOutputPopUp(m)
	case constant.ChooseFetchTagOptionPopUp:
		popUp = tag.RenderFetchTagOptionPopUp(m)
	case constant.FetchTagOutputPopUp:
		popUp = tag.RenderFetchTagOutputPopUp(m)
	case constant.RemoveRemoteConfirmationPopUp:
		popUp = remote.RenderRemoveRemoteConfirmationPopUp(m)
	case constant.RemoteAsTrackingUpstreamConfirmationPopUp:
		popUp = remote.RenderRemoteAsTrackingUpstreamConfirmationPopUp(m)
	case constant.EditRemotePromptPopUp:
		popUp = remote.RenderEditRemotePromptPopUp(m)
	case constant.GitRevertParentOptionSelectionPopUp:
		popUp = commitlog.RenderGitRevertParentOptionSelectionPopUp(m)
	case constant.GitRevertConfirmationPopUp:
		popUp = commitlog.RenderGitRevertConfirmationPopUp(m)
	case constant.GitCherryPickFromRefLogApplyConfirmationPopUp:
		popUp = reflog.RenderGitCherryPickFromRefLogApplyConfirmationPopUp(m)
	case constant.GitRebaseBranchInputPopUp:
		popUp = rebase.RenderGitRebaseBranchInputPopUp(m)
	case constant.GitRebaseOutputPopUp:
		popUp = rebase.RenderGitRebaseOutputPopUp(m)
	case constant.ChooseRemoteBranchOptionPopUp:
		popUp = branch.RenderChooseRemoteBranchOptionPopUp(m)
	case constant.ChooseBranchOptionForMergePopUp:
		popUp = branch.RenderChooseBranchOptionForMergePopUp(m)
	case constant.BranchMergeOutputPopUp:
		popUp = branch.RenderBranchMergeOutputPopUp(m)
	case constant.BlamePopUp:
		popUp = blame.RenderBlamePopUp(m)
	case constant.InteractiveRebaseOptionPopUp:
		popUp = interactiverebase.RenderInteractiveRebaseOptionPopUp(m)
	case constant.InteractiveRebaseFixupSquashSelectionPopUp:
		popUp = interactiverebase.RenderInteractiveRebaseFixupSquashSelectionPopUp(m)
	case constant.InteractiveRebaseFixupSquashCommitPopUp:
		popUp = interactiverebase.RenderInteractiveRebaseFixupSquashCommitPopUp(m)
	case constant.InteractiveRebaseFixupSquashOutputPopUp:
		popUp = interactiverebase.RenderInteractiveRebaseFixupSquashOutputPopUp(m)
	case constant.InteractiveRebaseRewordSelectionPopUp:
		popUp = interactiverebase.RenderInteractiveRebaseRewordSelectionPopUp(m)
	case constant.InteractiveRebaseRewordCommitPopUp:
		popUp = interactiverebase.RenderInteractiveRebaseRewordCommitPopUp(m)
	case constant.InteractiveRebaseRewordOutputPopUp:
		popUp = interactiverebase.RenderInteractiveRebaseRewordOutputPopUp(m)
	case constant.InteractiveRebaseDropSelectionPopUp:
		popUp = interactiverebase.RenderInteractiveRebaseDropSelectionPopUp(m)
	case constant.InteractiveRebaseDropOutputPopUp:
		popUp = interactiverebase.RenderInteractiveRebaseDropOutputPopUp(m)
	}
	return popUp
}
