package typing

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	worktreePopUp "github.com/gohyuhan/gitti/tui/popup/worktree"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle Shift+Tab key interaction during typing state.
//	Responsibility: Facilitates backward navigation between multiple input fields
//	within complex popups (like commit creation, cherry-pick editing, or fixup/squash commit editing).
//
// ------------------------------------
func handleTypingShiftTabKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch m.PopUpType {
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.MessageTextInput.Focus()
				popUp.DescriptionTextAreaInput.Blur()
			case 2:
				popUp.MessageTextInput.Blur()
				popUp.DescriptionTextAreaInput.Focus()
			}
		}
	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.MessageTextInput.Focus()
				popUp.DescriptionTextAreaInput.Blur()
			case 2:
				popUp.MessageTextInput.Blur()
				popUp.DescriptionTextAreaInput.Focus()
			}
		}
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.RemoteNameTextInput.Focus()
				popUp.RemoteUrlTextInput.Blur()
			case 2:
				popUp.RemoteNameTextInput.Blur()
				popUp.RemoteUrlTextInput.Focus()
			}
		}
	case constant.EditRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.EditRemotePromptPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.NewRemoteNameTextInput.Focus()
				popUp.NewRemoteUrlTextInput.Blur()
			case 2:
				popUp.NewRemoteNameTextInput.Blur()
				popUp.NewRemoteUrlTextInput.Focus()
			}
		}
	case constant.CreateTagPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.CreateTagPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.TagNameInput.Focus()
				popUp.TagMessageTextAreaInput.Blur()
			case 2:
				popUp.TagNameInput.Blur()
				popUp.TagMessageTextAreaInput.Focus()
			}
		}
	case constant.InteractiveRebaseFixupSquashCommitPopUp:
		popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashCommitPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.MessageTextInput.Focus()
				popUp.DescriptionTextAreaInput.Blur()
			case 2:
				popUp.MessageTextInput.Blur()
				popUp.DescriptionTextAreaInput.Focus()
			}
		}
	case constant.InteractiveRebaseRewordCommitPopUp:
		popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseRewordCommitPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.MessageTextInput.Focus()
				popUp.DescriptionTextAreaInput.Blur()
			case 2:
				popUp.MessageTextInput.Blur()
				popUp.DescriptionTextAreaInput.Focus()
			}
		}
	case constant.WorktreeAddNewWorktreePopUp:
		popUp, ok := m.PopUpModel.(*worktreePopUp.WorktreeAddNewWorktreePopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = max(popUp.CurrentActiveInputIndex-1, 1)
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.WorktreeNameTextInput.Focus()
				popUp.WorktreeBranchNameTextInput.Blur()
			case 2:
				popUp.WorktreeNameTextInput.Blur()
				popUp.WorktreeBranchNameTextInput.Focus()
			}
		}
	}
	return m, nil
}
