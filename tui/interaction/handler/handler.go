package handler

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	"github.com/gohyuhan/gitti/tui/types"
)

// typing is currently only on pop up model, so we can safely process it without checking if they were on pop up or not
func HandleTypingKeyBindingInteraction(msg tea.KeyMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return handleTypingESCKeyBindingInteraction(m)

		// in typing mode, tab is move to next input
	case "tab":
		return handleTypingTabKeyBindingInteraction(m)

	// in typing mode, shift+tab is move to previous input
	case "shift+tab":
		return handleTypingShiftTabKeyBindingInteraction(m)

	// because textare will use `enter` for line change and it will not be safe to use `enter` for submitting,
	// so `ctrl+e` will be used for submitting
	case "ctrl+e":
		return handleTypingCtrleKeyBindingInteraction(m)

	// because input mostly will no involve `enter` for change line, so `enter` can be safely used for submitting
	case "enter":
		return handleTypingEnterKeyBindingInteraction(m)
	}

	// for input typing update
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
	}
	return m, nil
}

func HandleNonTypingGlobalKeyBindingInteraction(msg tea.KeyMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch msg.String() {
	case "?":
		return handleNonTypingGlobalKeyBindingInteraction(m)

	case "1":
		return handleNonTyping1KeyBindingInteraction(m)

	case "2":
		return handleNonTyping2KeyBindingInteraction(m)

	case "3":
		return handleNonTyping3KeyBindingInteraction(m)

	case "4":
		return handleNonTyping4KeyBindingInteraction(m)

	case "A":
		return handleNonTypingaKeyBindingInteraction(m)

	case "c":
		return handleNonTypingcKeyBindingInteraction(m)

	case "d":
		return handleNonTypingdKeyBindingInteraction(m)

	case "n":
		return handleNonTypingnKeyBindingInteraction(m)

	case "p":
		return handleNonTypingpKeyBindingInteraction(m)

	case "P":
		return handleNonTypingPKeyBindingInteraction(m)

	case "r":
		return handleNonTypingrKeyBindingInteraction(m)

	case "s":
		return handleNonTypingsKeyBindingInteraction(m)

	case "S":
		return handleNonTypingSKeyBindingInteraction(m)

	case "[":
		return handleNonTypingLeftBracketKeyBindingInteraction(m)

	case "]":
		return handleNonTypingRightBracketKeyBindingInteraction(m)

	case "q", "Q":
		// only work when there is no pop up
		return handleNonTypingqQKeyBindingInteraction(m)

	case "backspace":
		return handleNonTypingBackspaceKeyBindingInteraction(m)

	case "enter":
		return handleNonTypingEnterKeyBindingInteraction(m)

	case "tab":
		// next component navigation
		return handleNonTypingTabKeyBindingInteraction(m)

	case "shift+tab":
		// previous component navigation
		return handleNonTypingShiftTabKeyBindingInteraction(m)

	case "space":
		return handleNonTypingSpaceKeyBindingInteraction(m)

	case "esc":
		return handleNonTypingEscKeyBindingInteraction(m)

	case "up", "k":
		return handleNonTypingUpkKeyBindingInteraction(msg, m)

	case "down", "j":
		return handleNonTypingDownjKeyBindingInteraction(msg, m)

	case "left", "h":
		return handleNonTypingLefthKeyBindingInteraction(msg, m)

	case "right", "l":
		return handleNonTypingRightlKeyBindingInteraction(msg, m)
	}
	return m, nil
}
