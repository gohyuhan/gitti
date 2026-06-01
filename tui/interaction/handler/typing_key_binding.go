package handler

import (
	"context"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/constant"
	blamePopUp "github.com/gohyuhan/gitti/tui/popup/blame"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	rebasePopUp "github.com/gohyuhan/gitti/tui/popup/rebase"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	worktreePopUp "github.com/gohyuhan/gitti/tui/popup/worktree"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
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
		constant.WorktreeAddNewWorktreePopUp:
		m.ShowPopUp.Store(false)
		m.IsTyping.Store(false)
		m.PopUpType = constant.NoPopUp
		m.PopUpModel = nil
	}
	return m, nil
}

// ------------------------------------
//
//	Handle Tab key interaction during typing state.
//	Responsibility: Facilitates forward navigation between multiple input fields
//	within complex popups (like commit creation, cherry-pick editing, or fixup/squash commit editing).
//
// ------------------------------------
func handleTypingTabKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch m.PopUpType {
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
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
			popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
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
			popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
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
			popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
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
			popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
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
			popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
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
			popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
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
			popUp.CurrentActiveInputIndex = min(popUp.CurrentActiveInputIndex+1, popUp.TotalInputCount)
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

// ------------------------------------
//
//	Handle Ctrl+e key interaction during typing state.
//	Responsibility: Acts as an explicit submission/confirmation trigger for multi-line
//	or complex inputs (e.g., executing a commit from the commit popup or fixup/squash from the fixup/squash commit popup).
//
// ------------------------------------
func handleTypingCtrleKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch m.PopUpType {
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			// once they start for commit process, reinit the input focus
			popUp.MessageTextInput.Focus()
			popUp.DescriptionTextAreaInput.Blur()
			popUp.CurrentActiveInputIndex = 1
			// start a seperate thread commit them and set the value of msg and desc to "" if committed successfully
			// also do not start any git operation is message is no provided
			if !popUp.IsProcessing.Load() && utf8.RuneCountInString(popUp.MessageTextInput.Value()) > 0 {
				if m.GitCommitRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
					gitArgs := m.GitOperations.GitCommit.GitCommitWithSigning(popUp.MessageTextInput.Value(), popUp.DescriptionTextAreaInput.Value(), false)
					return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.COMMIT_WITH_SIGNING_OPS)
				} else {
					services.GitCommitService(m, popUp.IsAmendCommit)
					popUp.InitialCommitStarted.Store(true)
					// Start spinner ticking
					return m, popUp.Spinner.Tick
				}
			}
		}
	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
		if ok {
			// once they start for amend commit process, reinit the input focus
			popUp.MessageTextInput.Focus()
			popUp.DescriptionTextAreaInput.Blur()
			popUp.CurrentActiveInputIndex = 1
			if !popUp.IsProcessing.Load() && utf8.RuneCountInString(popUp.MessageTextInput.Value()) > 0 {
				if m.GitCommitRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
					gitArgs := m.GitOperations.GitCommit.GitCommitWithSigning(popUp.MessageTextInput.Value(), popUp.DescriptionTextAreaInput.Value(), true)
					return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.AMEND_COMMIT_WITH_SIGNING_OPS)
				} else {
					services.GitAmendCommitService(m, popUp.IsAmendCommit)
					popUp.InitialCommitStarted.Store(true)
					// Start spinner ticking
					return m, popUp.Spinner.Tick
				}
			}
		}
	case constant.CreateTagPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.CreateTagPopUpModel)
		if ok {
			if utf8.RuneCountInString(popUp.TagNameInput.Value()) < 1 {
				return m, nil
			}
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
			tagPopUp.InitCreateTagConfirmationPopUpModel(m, popUp.TagNameInput.Value(), popUp.TagMessageTextAreaInput.Value(), popUp.CommitHash, popUp.CommitMessage)
			m.PopUpType = constant.CreateTagConfirmationPopUp
		}
	case constant.InteractiveRebaseFixupSquashCommitPopUp:
		popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashCommitPopUpModel)
		if ok {
			ogRetrievedCommitsList := popUp.OriginalRetrievedCommitList
			sortedSelectedCommits := popUp.SortedSelectedCommits
			fixupSquashCommitMessage := popUp.MessageTextInput.Value()
			fixupSquashCommitDescription := popUp.DescriptionTextAreaInput.Value()

			if utf8.RuneCountInString(fixupSquashCommitMessage) > 0 && len(ogRetrievedCommitsList) > 1 && len(sortedSelectedCommits) > 1 {
				// Switch to output popup before starting execution so errors/progress are visible immediately.
				interactiverebasePopUp.InitInteractiveRebaseFixupSquashOutputPopUpModel(m)
				popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashOutputPopUpModel)
				if !ok {
					return m, nil
				}
				m.PopUpType = constant.InteractiveRebaseFixupSquashOutputPopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				if m.GitCommitRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
					// Signing path returns prepared exec command; tea.ExecProcess handles terminal suspension.
					executor, cleanupCallbackFunc, fixupSquashErr := m.GitOperations.GitInteractiveRebase.GitInteractiveRebaseFixupSquashWithSigning(context.TODO(), ogRetrievedCommitsList, sortedSelectedCommits, fixupSquashCommitMessage, fixupSquashCommitDescription)
					if fixupSquashErr != nil {
						popUp.HasError.Store(true)
						popUp.FixupSquashOutputViewport.SetContent(fixupSquashErr.Error())
						return m, nil
					}
					return utils.SuspendGittiUIForGitOperationRequireSigningWithExecAndCleanUp(m, executor, cleanupCallbackFunc, logging.INTERACTIVE_REBASE_FIXUP_SQUASH)
				} else {
					if ok {
						popUp.IsProcessing.Store(true)
						// Non-signing path runs async service with cancellable context.
						services.InteractiveRebaseFixupSquashService(m, ogRetrievedCommitsList, sortedSelectedCommits, fixupSquashCommitMessage, fixupSquashCommitDescription)

						// Start spinner ticking
						return m, popUp.Spinner.Tick
					}
				}
			}
		}
	case constant.InteractiveRebaseRewordCommitPopUp:
		popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseRewordCommitPopUpModel)
		if ok {
			ogRetrievedCommitsList := popUp.OriginalRetrievedCommitList
			selectedCommit := popUp.SelectedCommit
			rewordCommitMessage := popUp.MessageTextInput.Value()
			rewordCommitDescription := popUp.DescriptionTextAreaInput.Value()

			if utf8.RuneCountInString(rewordCommitMessage) > 0 && len(ogRetrievedCommitsList) > 0 {
				// Switch to output popup before starting execution so errors/progress are visible immediately.
				interactiverebasePopUp.InitInteractiveRebaseRewordOutputPopUpModel(m)
				popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseRewordOutputPopUpModel)
				if !ok {
					return m, nil
				}
				m.PopUpType = constant.InteractiveRebaseRewordOutputPopUp
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				if m.GitCommitRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
					// Signing path returns prepared exec command; tea.ExecProcess handles terminal suspension.
					executor, cleanupCallbackFunc, rewordErr := m.GitOperations.GitInteractiveRebase.GitInteractiveRebaseRewordWithSigning(context.TODO(), ogRetrievedCommitsList, selectedCommit, rewordCommitMessage, rewordCommitDescription)
					if rewordErr != nil {
						popUp.HasError.Store(true)
						popUp.RewordOutputViewport.SetContent(rewordErr.Error())
						return m, nil
					}
					return utils.SuspendGittiUIForGitOperationRequireSigningWithExecAndCleanUp(m, executor, cleanupCallbackFunc, logging.INTERACTIVE_REBASE_REWORD)
				} else {
					if ok {
						popUp.IsProcessing.Store(true)
						// Non-signing path runs async service with cancellable context.
						services.InteractiveRebaseRewordService(m, ogRetrievedCommitsList, selectedCommit, rewordCommitMessage, rewordCommitDescription)

						// Start spinner ticking
						return m, popUp.Spinner.Tick
					}
				}
			}
		}
	}
	return m, nil
}

// ------------------------------------
//
//	Handle Enter key interaction during typing state.
//	Responsibility: Context-dependent behavior for text input.
//	- Single-line inputs (e.g., prompts): Submits the input and triggers the associated action.
//	- Multi-line inputs (e.g., commit message body): Inserts a newline character into the text area.
//
// ------------------------------------
func handleTypingEnterKeyBindingInteraction(m *types.GittiModel, msg tea.KeyPressMsg) (*types.GittiModel, tea.Cmd) {
	switch m.PopUpType {
	case constant.AddRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.AddRemotePromptPopUpModel)
		if ok {
			// once they start for commit process, reinit the input focus
			popUp.RemoteNameTextInput.Focus()
			popUp.RemoteUrlTextInput.Blur()
			popUp.CurrentActiveInputIndex = 1
			// start a seperate thread that stage the current selected files and commit them and set the value of msg and desc to "" if committed successfully
			// also do not start any git operation is message is no provided
			if !popUp.IsProcessing.Load() && utf8.RuneCountInString(popUp.RemoteNameTextInput.Value()) > 0 && utf8.RuneCountInString(popUp.RemoteUrlTextInput.Value()) > 0 {
				services.GitAddRemoteService(m)
			}
		}
	case constant.EditRemotePromptPopUp:
		popUp, ok := m.PopUpModel.(*remotePopUp.EditRemotePromptPopUpModel)
		if ok {
			newRemoteName := popUp.NewRemoteNameTextInput.Value()
			newRemoteUrl := popUp.NewRemoteUrlTextInput.Value()

			services.GitEditRemoteNameAndUrlService(m, popUp.OldRemoteName, newRemoteName, popUp.OldRemoteUrl, newRemoteUrl)
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpType = constant.NoPopUp
			m.PopUpModel = nil
		}
	case constant.CreateNewBranchPopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateNewBranchPopUpModel)
		if ok {
			// we direclty close the pop up and trigger the branch creation operation
			validBranchName, _ := api.IsBranchNameValid(popUp.NewBranchNameInput.Value())
			if len(validBranchName) > 0 {
				switch popUp.CreateType {
				case git.NEWBRANCH:
					services.GitCreateNewBranchService(m, validBranchName)
				case git.NEWBRANCHANDSWITCH:
					services.GitCreateNewBranchAndSwitchService(m, validBranchName)
				case git.NEWBRANCHBASEDONCOMMITHASH:
					services.GitCreateNewBranchBasedOnCommitHashService(m, validBranchName, popUp.CommitHash)
				}
				m.ShowPopUp.Store(false)
				m.IsTyping.Store(false)
				m.PopUpType = constant.NoPopUp
				m.PopUpModel = nil
			}
		}

	case constant.GitStashMessagePopUp:
		popUp, ok := m.PopUpModel.(*stashPopUp.GitStashMessagePopUpModel)
		if ok {
			msg := popUp.StashMessageInput.Value()
			switch popUp.StashType {
			case git.STASHALL:
				stashPopUp.InitGitStashConfirmPromptPopUpModel(m, git.STASHALL, "", "", msg)
			case git.STASHFILE:
				stashPopUp.InitGitStashConfirmPromptPopUpModel(m, git.STASHFILE, popUp.FilePathName, "", msg)
			}
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
			m.PopUpType = constant.GitStashConfirmPromptPopUp
		}

	case constant.CreateBranchBasedOnRemotePopUp:
		popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemotePopUpModel)
		if ok {
			// we direclty close the pop up and trigger the branch creation operation
			validBranchName, _ := api.IsBranchNameValid(popUp.RemoteBranchNameInput.Value())
			remoteOrigin := popUp.RemoteOrigin
			if utf8.RuneCountInString(validBranchName) > 0 {
				branchPopUp.InitCreateBranchBasedOnRemoteOutputPopUpModel(m)
				popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemoteOutputPopUpModel)
				if ok {
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					m.PopUpType = constant.CreateBranchBasedOnRemoteOutputPopUp
					popUp.IsProcessing.Store(true)
					services.CreateNewBranchBasedOnRemoteService(m, remoteOrigin, validBranchName, git.NEWBRANCHBASEDONREMOTEUSERINPUT)
					return m, popUp.Spinner.Tick
				} else {
					m.ShowPopUp.Store(false)
					m.IsTyping.Store(false)
					m.PopUpType = constant.NoPopUp
				}
			}
		}

	case constant.GitRebaseBranchInputPopUp:
		popUp, ok := m.PopUpModel.(*rebasePopUp.GitRebaseBranchInputPopUpModel)
		if ok {
			// we direclty close the pop up and trigger the branch creation operation
			validBranchName, _ := api.IsBranchNameValid(popUp.BranchNameInput.Value())
			remoteOrigin := popUp.Remote
			if utf8.RuneCountInString(validBranchName) > 0 {
				rebasePopUp.InitGitRebaseOutputPopUpModel(m)
				popUp, ok := m.PopUpModel.(*rebasePopUp.GitRebaseOutputPopUpModel)
				if ok {
					if m.GitCommitRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
						gitArgs := m.GitOperations.GitRebase.GitRebaseWithSigning(remoteOrigin, validBranchName)
						return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.GIT_REBASE_WITH_SIGNING_OPS)
					} else {
						m.ShowPopUp.Store(true)
						m.IsTyping.Store(false)
						m.PopUpType = constant.GitRebaseOutputPopUp
						popUp.IsProcessing.Store(true)
						services.GitRebaseService(m, remoteOrigin, validBranchName)
						return m, popUp.Spinner.Tick
					}
				} else {
					m.ShowPopUp.Store(false)
					m.IsTyping.Store(false)
					m.PopUpType = constant.NoPopUp
				}
			}
		}
	case constant.WorktreeAddNewWorktreePopUp:
		popUp, ok := m.PopUpModel.(*worktreePopUp.WorktreeAddNewWorktreePopUpModel)
		if ok {
			// we directly close the pop up and trigger the new worktree creation operation
			newWorktreeName := popUp.WorktreeNameTextInput.Value()
			checkoutWorktreeBranchName := popUp.WorktreeBranchNameTextInput.Value()
			if utf8.RuneCountInString(newWorktreeName) > 0 {
				worktreePopUp.InitWorktreeAddNewWorktreeOutputPopUpModel(m)
				popUp, ok := m.PopUpModel.(*worktreePopUp.WorktreeAddNewWorktreeOutputPopUpModel)
				if ok {
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					m.PopUpType = constant.WorktreeAddNewWorktreeOutputPopUp
					popUp.IsProcessing.Store(true)
					services.AddNewWorktreeService(m, newWorktreeName, checkoutWorktreeBranchName)
					return m, popUp.Spinner.Tick
				} else {
					m.ShowPopUp.Store(false)
					m.IsTyping.Store(false)
					m.PopUpType = constant.NoPopUp
				}
			}
		}

	// the following is to handle the change line for textarea input
	case constant.CommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel)
		if ok {
			if popUp.CurrentActiveInputIndex == 2 {
				var cmd tea.Cmd
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}

	case constant.AmendCommitPopUp:
		popUp, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel)
		if ok {
			if popUp.CurrentActiveInputIndex == 2 {
				var cmd tea.Cmd
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}

	case constant.CreateTagPopUp:
		popUp, ok := m.PopUpModel.(*tagPopUp.CreateTagPopUpModel)
		if ok {
			if popUp.CurrentActiveInputIndex == 2 {
				var cmd tea.Cmd
				popUp.TagMessageTextAreaInput, cmd = popUp.TagMessageTextAreaInput.Update(msg)
				return m, cmd
			}
		}

	case constant.BlamePopUp:
		popUp, ok := m.PopUpModel.(*blamePopUp.BlamePopUpModel)
		if ok {
			selectedFilepath := popUp.CurrentGitTrackedFilesPathList.SelectedItem()
			if selectedFilepath == nil {
				return m, nil
			}
			parsedFilePath := selectedFilepath.(blamePopUp.CurrentGitTrackedFilesPathItem).FilePath
			popUp.ShowBlameInfoView(parsedFilePath)
			services.GetFileGitBlameInfoService(m, parsedFilePath)
			return m, nil
		}

	case constant.InteractiveRebaseFixupSquashCommitPopUp:
		popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashCommitPopUpModel)
		if ok {
			if popUp.CurrentActiveInputIndex == 2 {
				var cmd tea.Cmd
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}

	case constant.InteractiveRebaseRewordCommitPopUp:
		popUp, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseRewordCommitPopUpModel)
		if ok {
			if popUp.CurrentActiveInputIndex == 2 {
				var cmd tea.Cmd
				popUp.DescriptionTextAreaInput, cmd = popUp.DescriptionTextAreaInput.Update(msg)
				return m, cmd
			}
		}
	}
	return m, nil
}

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
	}
	return m, nil
}

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
	}

	err := clipboard.WriteAll(content)
	if err != nil {
		// TODO: log the error
	}

	return m, nil
}

// ------------------------------------
//
//	Handle up arrow in typing mode. In the blame popup, navigates the file list
//	up when the file selector is active, or scrolls the blame viewport up when
//	blame info is shown.
//
// ------------------------------------
func handleTypingUpKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.BlamePopUp:
			popUp, ok := m.PopUpModel.(*blamePopUp.BlamePopUpModel)
			if ok {
				if !popUp.ShowingBlameInfo {
					popUp.CurrentGitTrackedFilesPathList.CursorUp()
				} else {
					popUp.BlameViewport.ScrollUp(1)
				}
			}
		}
	}

	return m, nil
}

// ------------------------------------
//
//	Handle down arrow in typing mode. In the blame popup, navigates the file
//	list down when the file selector is active, or scrolls the blame viewport
//	down when blame info is shown.
//
// ------------------------------------
func handleTypingDownKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.BlamePopUp:
			popUp, ok := m.PopUpModel.(*blamePopUp.BlamePopUpModel)
			if ok {
				if !popUp.ShowingBlameInfo {
					popUp.CurrentGitTrackedFilesPathList.CursorDown()
				} else {
					popUp.BlameViewport.ScrollDown(1)
				}
			}
		}
	}

	return m, nil
}

// ------------------------------------
//
//	Handle left arrow in typing mode. In the blame popup, scrolls the blame
//	viewport left when blame info is shown.
//
// ------------------------------------
func handleTypingLeftKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.BlamePopUp:
			popUp, ok := m.PopUpModel.(*blamePopUp.BlamePopUpModel)
			if ok {
				if popUp.ShowingBlameInfo {
					popUp.BlameViewport.ScrollLeft(1)
				}
			}
		}
	}

	return m, nil
}

// ------------------------------------
//
//	Handle right arrow in typing mode. In the blame popup, scrolls the blame
//	viewport right when blame info is shown.
//
// ------------------------------------
func handleTypingRightKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.BlamePopUp:
			popUp, ok := m.PopUpModel.(*blamePopUp.BlamePopUpModel)
			if ok {
				if popUp.ShowingBlameInfo {
					popUp.BlameViewport.ScrollRight(1)
				}
			}
		}
	}

	return m, nil
}
