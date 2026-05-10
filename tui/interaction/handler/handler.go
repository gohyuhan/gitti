package handler

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/layout"
	blamePopUp "github.com/gohyuhan/gitti/tui/popup/blame"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	rebasePopUp "github.com/gohyuhan/gitti/tui/popup/rebase"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle key presses in typing mode (always inside a popup). Routes esc, tab,
//	shift+tab, ctrl+e, enter, ctrl+p, ctrl+y, and arrow keys to dedicated
//	handlers, then forwards remaining input to the active popup's text
//	input or textarea field.
//
// ------------------------------------
func HandleTypingKeyBindingInteraction(msg tea.KeyPressMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd

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
		return handleTypingEnterKeyBindingInteraction(m, msg)

	// to paste clipboard content into the current input field
	case "ctrl+p":
		return handleTypingCtrlpKeyBindingInteraction(m)

	// to copy content from current input field to clipboard
	case "ctrl+y":
		return handleTypingCtrlyKeyBindingInteraction(m)

	case "up":
		m, cmd = handleTypingUpKeyBindingInteraction(msg, m)

	case "down":
		m, cmd = handleTypingDownKeyBindingInteraction(msg, m)

	case "left":
		m, cmd = handleTypingLeftKeyBindingInteraction(msg, m)

	case "right":
		m, cmd = handleTypingRightKeyBindingInteraction(msg, m)
	}

	// for input typing update
	switch m.PopUpType {
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.MessageTextInput, cmd = popUp.MessageTextInput.Update(msg)
				return m, cmd

			case 2:
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.MessageTextInput, cmd = popUp.MessageTextInput.Update(msg)
				return m, cmd

			case 2:
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.RemoteNameTextInput, cmd = popUp.RemoteNameTextInput.Update(msg)
				return m, cmd

			case 2:
				popUp.RemoteUrlTextInput, cmd = popUp.RemoteUrlTextInput.Update(msg)
				return m, cmd
			}
		}
	case constant.EditRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.EditRemotePromptPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.NewRemoteNameTextInput, cmd = popUp.NewRemoteNameTextInput.Update(msg)
				return m, cmd

			case 2:
				popUp.NewRemoteUrlTextInput, cmd = popUp.NewRemoteUrlTextInput.Update(msg)
				return m, cmd
			}
		}
	case constant.CreateNewBranchPopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateNewBranchPopUpModel)
		if ok {
			popUp.NewBranchNameInput, cmd = popUp.NewBranchNameInput.Update(msg)
			return m, cmd
		}
	case constant.GitStashMessagePopUp:
		popUp, ok := m.PopUpModel.(*stashPopUp.GitStashMessagePopUpModel)
		if ok {
			popUp.StashMessageInput, cmd = popUp.StashMessageInput.Update(msg)
			return m, cmd
		}
	case constant.CreateBranchBasedOnRemotePopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemotePopUpModel)
		if ok {
			popUp.RemoteBranchNameInput, cmd = popUp.RemoteBranchNameInput.Update(msg)
			return m, cmd
		}
	case constant.CreateTagPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.CreateTagPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.TagNameInput, cmd = popUp.TagNameInput.Update(msg)
				return m, cmd

			case 2:
				popUp.TagMessageTextAreaInput, cmd = popUp.TagMessageTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	case constant.GitRebaseBranchInputPopUp:
		popUp, ok := m.PopUpModel.(*rebasePopUp.GitRebaseBranchInputPopUpModel)
		if ok {
			popUp.BranchNameInput, cmd = popUp.BranchNameInput.Update(msg)
			return m, cmd
		}
	case constant.BlamePopUp:
		popUp, ok := m.PopUpModel.(*blamePopUp.BlamePopUpModel)
		if ok && !popUp.ShowingBlameInfo && msg.String() != "up" && msg.String() != "down" {
			popUp.FilterInput, cmd = popUp.FilterInput.Update(msg)
			popUp.FilterValue = popUp.FilterInput.Value()
			popUp.CurrentGitTrackedFilesPathList.SetFilterText(popUp.FilterValue)
			return m, cmd
		}
	case constant.InteractiveRebaseFixupSquashCommitPopUp:
		popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashCommitPopUpModel)
		if ok {
			switch popUp.CurrentActiveInputIndex {
			case 1:
				popUp.MessageTextInput, cmd = popUp.MessageTextInput.Update(msg)
				return m, cmd

			case 2:
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	}
	return m, cmd
}

// ------------------------------------
//
//	Handle key presses in non-typing mode. Dispatches each key (letters, arrows,
//	special keys, and panel-resize +/-) to its dedicated handler function.
//
// ------------------------------------
func HandleNonTypingGlobalKeyBindingInteraction(msg tea.KeyPressMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
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

	case "a":
		return handleNonTypingaKeyBindingInteraction(m)

	case "A":
		return handleNonTypingAKeyBindingInteraction(m)

	case "b":
		return handleNonTypingbKeyBindingInteraction(m)

	case "c":
		return handleNonTypingcKeyBindingInteraction(m)

	case "C":
		return handleNonTypingCKeyBindingInteraction(m)

	case "d":
		return handleNonTypingdKeyBindingInteraction(m)

	case "e":
		return handleNonTypingeKeyBindingInteraction(m)

	case "f":
		return handleNonTypingfKeyBindingInteraction(m)

	case "i":
		return handleNonTypingiKeyBindingInteraction(m)

	case "L":
		return handleNonTypingLKeyBindingInteraction(m)

	case "m":
		return handleNonTypingmKeyBindingInteraction(m)

	case "n":
		return handleNonTypingnKeyBindingInteraction(m)

	case "p":
		return handleNonTypingpKeyBindingInteraction(m)

	case "P":
		return handleNonTypingPKeyBindingInteraction(m)

	case "r":
		return handleNonTypingrKeyBindingInteraction(m)

	case "R":
		return handleNonTypingRKeyBindingInteraction(m)

	case "s":
		return handleNonTypingsKeyBindingInteraction(m)

	case "S":
		return handleNonTypingSKeyBindingInteraction(m)

	case "t":
		return handleNonTypingtKeyBindingInteraction(m)

	case "[":
		return handleNonTypingLeftBracketKeyBindingInteraction(m)

	case "]":
		return handleNonTypingRightBracketKeyBindingInteraction(m)

	case "/":
		return handleNonTypingSlashKeyBindingInteraction(m)

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
	case "-":
		if !m.ShowPopUp.Load() {
			m.WindowLeftPanelRatio = max(settings.MINLEFTPANELWIDTHRATIO, m.WindowLeftPanelRatio-0.01)
			layout.TuiWindowSizing(m)
		}
		return m, nil
	case "+":
		if !m.ShowPopUp.Load() {
			m.WindowLeftPanelRatio = min(settings.MAXLEFTPANELWIDTHRATIO, m.WindowLeftPanelRatio+0.01)
			layout.TuiWindowSizing(m)
		}
		return m, nil
	case "<":
		return handleNonTypingLeftAngleBracketKeyBindingInteraction(m)

	case ">":
		return handleNonTypingRightAngleBracketKeyBindingInteraction(m)

	case "ctrl+a":
		return handleNonTypingCtrlaKeyBindingInteraction(m)
	case "ctrl+p":
		return handleNonTypingCtrlpKeyBindingInteraction(m)
	case "ctrl+k":
		return handleNonTypingCtrlkKeyBindingInteraction(m)
	case "ctrl+r":
		return handleNonTypingCtrlrKeyBindingInteraction(m)
	}
	return m, nil
}
