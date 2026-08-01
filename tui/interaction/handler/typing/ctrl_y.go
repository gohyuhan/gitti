package typing

import (
	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	rebasePopUp "github.com/gohyuhan/gitti/tui/popup/rebase"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	worktreePopUp "github.com/gohyuhan/gitti/tui/popup/worktree"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle Ctrl+y key interaction during typing state.
//	Responsibility: Copies the entire content of the currently focused input field
//	or text area into the system clipboard.
//
// ------------------------------------
func handleTypingCtrlyKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var content string
	switch m.PopUpType {
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				content = popUp.MessageTextInput.Value()
			case 2:
				content = popUp.DescriptionTextAreaInput.Value()
			}
		}
	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				content = popUp.MessageTextInput.Value()
			case 2:
				content = popUp.DescriptionTextAreaInput.Value()
			}
		}
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				content = popUp.RemoteNameTextInput.Value()
			case 2:
				content = popUp.RemoteUrlTextInput.Value()
			}
		}
	case constant.CreateNewBranchPopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateNewBranchPopUpModel)
		if ok {
			content = popUp.NewBranchNameInput.Value()
		}
	case constant.GitStashMessagePopUp:
		popUp, ok := m.PopUpModel.(*stashPopUp.GitStashMessagePopUpModel)
		if ok {
			content = popUp.StashMessageInput.Value()
		}
	case constant.CreateBranchBasedOnRemotePopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemotePopUpModel)
		if ok {
			content = popUp.RemoteBranchNameInput.Value()
		}
	case constant.CreateTagPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.CreateTagPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				content = popUp.TagNameInput.Value()
			case 2:
				content = popUp.TagMessageTextAreaInput.Value()
			}
		}
	case constant.GitRebaseBranchInputPopUp:
		popUp, ok := m.PopUpModel.(*rebasePopUp.GitRebaseBranchInputPopUpModel)
		if ok {
			content = popUp.BranchNameInput.Value()
		}
	case constant.InteractiveRebaseFixupSquashCommitPopUp:
		popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				content = popUp.MessageTextInput.Value()
			case 2:
				content = popUp.DescriptionTextAreaInput.Value()
			}
		}
	case constant.InteractiveRebaseRewordCommitPopUp:
		popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseRewordCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				content = popUp.MessageTextInput.Value()
			case 2:
				content = popUp.DescriptionTextAreaInput.Value()
			}
		}
	case constant.WorktreeAddNewWorktreePopUp:
		popUp, ok := m.PopUpModel.(*worktreePopUp.WorktreeAddNewWorktreePopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				content = popUp.WorktreeNameTextInput.Value()
			case 2:
				content = popUp.WorktreeBranchNameTextInput.Value()
			}
		}
	case constant.WorktreeLockReasonInputPopUp:
		popUp, ok := m.PopUpModel.(*worktreePopUp.WorktreeLockReasonInputPopUpModel)
		if ok {
			content = popUp.WorktreeLockReasonTextInput.Value()
		}
	}

	err := clipboard.WriteAll(content)
	if err != nil {
		// TODO: log the error
	}

	return m, nil
}
