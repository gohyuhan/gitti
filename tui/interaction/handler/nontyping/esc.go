package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle Esc key interaction.
//	Responsibility: Universal "cancel" or "back" operation.
//	- If operations are processing (pushing, pulling, etc.), blocking esc prevents interruption.
//	- Cancels out of active background git services if possible.
//	- Dismisses almost any active open popup, returning the user to the underlying view.
//	- If in Detail line-editing mode, exits that specific mode back to the parent panel.
//
// ------------------------------------
func handleNonTypingEscKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.GitRemotePushPopUp:
			services.GitRemotePushCancelService(m)
			m.PopUpModel = nil
		case constant.GitPullOutputPopUp:
			services.GitPullCancelService(m)
			m.PopUpModel = nil
		case constant.PushTagOutputPopUp:
			services.GitPushTagCancelService(m)
			m.PopUpModel = nil
		case constant.FetchTagOutputPopUp:
			services.GitFetchTagCancelService(m)
			m.PopUpModel = nil
		case constant.DeleteTagOutputPopUp:
			services.DeleteTagCancelService(m)
			m.PopUpModel = nil
		case constant.GitRebaseOutputPopUp:
			services.GitRebaseCancelService(m)
			m.PopUpModel = nil
		case constant.BranchMergeOutputPopUp:
			services.GitMergeCancelService(m)
			m.PopUpModel = nil
		case constant.InteractiveRebaseFixupSquashOutputPopUp:
			services.InteractiveRebaseFixupSquashCancelService(m)
			m.PopUpModel = nil
		case constant.InteractiveRebaseRewordOutputPopUp:
			services.InteractiveRebaseRewordCancelService(m)
			m.PopUpModel = nil
		case constant.InteractiveRebaseDropOutputPopUp:
			services.InteractiveRebaseDropCancelService(m)
			m.PopUpModel = nil
		case constant.WorktreeAddNewWorktreeOutputPopUp:
			services.AddNewWorktreeCancelService(m)
			m.PopUpModel = nil
		case constant.SwitchBranchOutputPopUp:
			// Block ESC during branch switching - operation must complete
			popUp, ok := m.PopUpModel.(*branchPopUp.SwitchBranchOutputPopUpModel)
			if ok && !popUp.IsProcessing.Load() {
				// only close when done processing
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		case constant.GitStashOperationOutputPopUp:
			// Block ESC during stash operation - operation must complete
			popUp, ok := m.PopUpModel.(*stashPopUp.GitStashOperationOutputPopUpModel)
			if ok && !popUp.IsProcessing.Load() {
				// only close when done processing
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		case constant.GitDeleteBranchOutputPopUp:
			popUp, ok := m.PopUpModel.(*branchPopUp.GitDeleteBranchOutputPopUpModel)
			if ok && !popUp.IsProcessing.Load() {
				// only close when done processing
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		case constant.CreateBranchBasedOnRemoteOutputPopUp:
			// Block ESC during create new branch based on remote operation - operation must complete
			popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemoteOutputPopUpModel)
			if ok && !popUp.IsProcessing.Load() {
				// only close when done processing
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		case constant.KeybindingAndFeatureInstructionsPopUp,
			constant.ChooseRemotePopUp,
			constant.ChoosePushTypePopUp,
			constant.ChooseNewBranchTypePopUp,
			constant.ChooseSwitchBranchTypePopUp,
			constant.ChooseGitPullTypePopUp,
			constant.GitDiscardTypeOptionPopUp,
			constant.GitDiscardConfirmPromptPopUp,
			constant.GitStashConfirmPromptPopUp,
			constant.GitResolveConflictOptionPopUp,
			constant.GitDeleteBranchConfirmPromptPopUp,
			constant.GitResetLatestCommitTypeOptionPopUp,
			constant.GitResetLatestCommitConfirmPromptPopUp,
			constant.GitResetToSelectedCommitTypeOptionPopUp,
			constant.GitResetToSelectedCommitConfirmPromptPopUp,
			constant.GitCherryPickPopUp,
			constant.GitEditCherryPickPopUp,
			constant.GitCherryPickOptionSelectionPopUp,
			constant.GitCherryPickApplyConfirmPopUp,
			constant.GitDiscardFileLineChangeConfirmPopUp,
			constant.CreateTagConfirmationPopUp,
			constant.ChooseDeleteTagOptionPopUp,
			constant.ChooseRemoteForDeleteRemoteTagPopUp,
			constant.ChoosePushTagOptionPopUp,
			constant.ChooseFetchTagOptionPopUp,
			constant.RemoveRemoteConfirmationPopUp,
			constant.RemoteAsTrackingUpstreamConfirmationPopUp,
			constant.GitRevertParentOptionSelectionPopUp,
			constant.GitRevertConfirmationPopUp,
			constant.GitCherryPickFromRefLogApplyConfirmationPopUp,
			constant.ChooseRemoteBranchOptionPopUp,
			constant.ChooseBranchOptionForMergePopUp,
			constant.InteractiveRebaseOptionPopUp,
			constant.InteractiveRebaseFixupSquashSelectionPopUp,
			constant.InteractiveRebaseRewordSelectionPopUp,
			constant.InteractiveRebaseDropSelectionPopUp,
			constant.WorktreeRemoveWorktreeConfirmationPopUp:
			// simple closing of the pop up
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpType = constant.NoPopUp
			m.PopUpModel = nil
		}
		return m, nil
	} else {
		switch m.CurrentSelectedComponent {
		case constant.DetailComponentPanel, constant.DetailComponentPanelTwo:
			if m.IsLineEditingState.Load() {
				m.IsLineEditingState.Store(false)
				m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
					Event: constant.DETAIL_COMPONENT_PANEL_LAYOUT_UPDATED_EVENT,
				}
			} else {
				m.CurrentSelectedComponent = m.DetailPanelParentComponent
				m.DetailPanelParentComponent = ""
			}
		}
	}
	return m, nil
}
