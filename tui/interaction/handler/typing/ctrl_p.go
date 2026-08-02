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
//	Handle Ctrl+p key interaction during typing state.
//	Responsibility: Reads text from the system clipboard and inserts it at the cursor
//	position within the currently focused input field or text area.
//
// ------------------------------------
func handleTypingCtrlpKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	content, err := clipboard.ReadAll()
	if err != nil {
		return m, nil
	}
	msg := tea.PasteMsg{
		Content: content,
	}
	switch m.PopUpType {
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				var cmd tea.Cmd
				popUp.MessageTextInput, cmd = popUp.MessageTextInput.Update(msg)
				return m, cmd

			case 2:
				var cmd tea.Cmd
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				var cmd tea.Cmd
				popUp.MessageTextInput, cmd = popUp.MessageTextInput.Update(msg)
				return m, cmd

			case 2:
				var cmd tea.Cmd
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				var cmd tea.Cmd
				popUp.RemoteNameTextInput, cmd = popUp.RemoteNameTextInput.Update(msg)
				return m, cmd

			case 2:
				var cmd tea.Cmd
				popUp.RemoteUrlTextInput, cmd = popUp.RemoteUrlTextInput.Update(msg)
				return m, cmd
			}
		}
	case constant.CreateNewBranchPopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateNewBranchPopUpModel)
		if ok {
			var cmd tea.Cmd
			popUp.NewBranchNameInput, cmd = popUp.NewBranchNameInput.Update(msg)
			return m, cmd
		}
	case constant.GitStashMessagePopUp:
		popUp, ok := m.PopUpModel.(*stashPopUp.GitStashMessagePopUpModel)
		if ok {
			var cmd tea.Cmd
			popUp.StashMessageInput, cmd = popUp.StashMessageInput.Update(msg)
			return m, cmd
		}
	case constant.CreateBranchBasedOnRemotePopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemotePopUpModel)
		if ok {
			var cmd tea.Cmd
			popUp.RemoteBranchNameInput, cmd = popUp.RemoteBranchNameInput.Update(msg)
			return m, cmd
		}
	case constant.CreateTagPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.CreateTagPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				var cmd tea.Cmd
				popUp.TagNameInput, cmd = popUp.TagNameInput.Update(msg)
				return m, cmd

			case 2:
				var cmd tea.Cmd
				popUp.TagMessageTextAreaInput, cmd = popUp.TagMessageTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	case constant.GitRebaseBranchInputPopUp:
		popUp, ok := m.PopUpModel.(*rebasePopUp.GitRebaseBranchInputPopUpModel)
		if ok {
			var cmd tea.Cmd
			popUp.BranchNameInput, cmd = popUp.BranchNameInput.Update(msg)
			return m, cmd
		}
	case constant.InteractiveRebaseFixupSquashCommitPopUp:
		popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				var cmd tea.Cmd
				popUp.MessageTextInput, cmd = popUp.MessageTextInput.Update(msg)
				return m, cmd

			case 2:
				var cmd tea.Cmd
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	case constant.InteractiveRebaseRewordCommitPopUp:
		popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseRewordCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				var cmd tea.Cmd
				popUp.MessageTextInput, cmd = popUp.MessageTextInput.Update(msg)
				return m, cmd

			case 2:
				var cmd tea.Cmd
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	case constant.WorktreeAddNewWorktreePopUp:
		popUp, ok := m.PopUpModel.(*worktreePopUp.WorktreeAddNewWorktreePopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				var cmd tea.Cmd
				popUp.WorktreeNameTextInput, cmd = popUp.WorktreeNameTextInput.Update(msg)
				return m, cmd

			case 2:
				var cmd tea.Cmd
				popUp.WorktreeBranchNameTextInput, cmd = popUp.WorktreeBranchNameTextInput.Update(msg)
				return m, cmd
			}
		}
	case constant.WorktreeLockReasonInputPopUp:
		popUp, ok := m.PopUpModel.(*worktreePopUp.WorktreeLockReasonInputPopUpModel)
		if ok {
			var cmd tea.Cmd
			popUp.WorktreeLockReasonTextInput, cmd = popUp.WorktreeLockReasonTextInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}
