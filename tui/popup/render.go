package popup

import (
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/popup/branch"
	"github.com/gohyuhan/gitti/tui/popup/commit"
	"github.com/gohyuhan/gitti/tui/popup/discard"
	"github.com/gohyuhan/gitti/tui/popup/keybinding"
	"github.com/gohyuhan/gitti/tui/popup/pull"
	"github.com/gohyuhan/gitti/tui/popup/push"
	"github.com/gohyuhan/gitti/tui/popup/remote"
	"github.com/gohyuhan/gitti/tui/popup/resolve"
	"github.com/gohyuhan/gitti/tui/popup/stash"
	"github.com/gohyuhan/gitti/tui/types"
)

// -----------------------------------------------------------------------------
//
//	Functions that related to the rendering of pop up
//
// -----------------------------------------------------------------------------
// render the PopUp and the content within it will be a determine dynamically
func RenderPopUpComponent(m *types.GittiModel) string {
	var popUp string

	switch m.PopUpType {
	case constant.GlobalKeyBindingPopUp:
		popUp = keybinding.RenderGlobalKeyBindingPopUp(m)
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
	}
	return popUp
}
