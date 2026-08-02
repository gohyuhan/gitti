package typing

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	blamePopUp "github.com/gohyuhan/gitti/tui/popup/blame"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle ESC key interaction during typing state.
//	Responsibility: Acts as a primary cancellation mechanism when the user is in an input field.
//	It clears input models, resets state flags, and dismisses active popups without saving.
//
// ------------------------------------
func handleTypingESCKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch m.PopUpType {
	case constant.CommitPopUp:
		services.GitCommitCancelService(m)
	case constant.AmendCommitPopUp:
		services.GitAmendCommitCancelService(m)
	case constant.AddRemotePromptPopUp:
		services.GitAddRemoteCancelService(m)
	case constant.BlamePopUp:
		popUp, ok := m.PopUpModel.(*blamePopUp.BlamePopUpModel)
		if ok {
			if popUp.ShowingBlameInfo {
				popUp.ResetSelectedBlameFile()
			} else {
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		}
	case constant.CreateNewBranchPopUp,
		constant.GitStashMessagePopUp,
		constant.CreateBranchBasedOnRemotePopUp,
		constant.CreateTagPopUp,
		constant.EditRemotePromptPopUp,
		constant.GitRebaseBranchInputPopUp,
		constant.InteractiveRebaseFixupSquashCommitPopUp,
		constant.InteractiveRebaseRewordCommitPopUp,
		constant.WorktreeAddNewWorktreePopUp,
		constant.WorktreeLockReasonInputPopUp:
		m.ShowPopUp.Store(false)
		m.IsTyping.Store(false)
		m.PopUpType = constant.NoPopUp
		m.PopUpModel = nil
	}
	return m, nil
}
