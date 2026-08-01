package typing

import (
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
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
	case constant.WorktreeLockReasonInputPopUp:
		popUp, ok := m.PopUpModel.(*worktreePopUp.WorktreeLockReasonInputPopUpModel)
		if ok {
			lockReason := popUp.WorktreeLockReasonTextInput.Value()
			m.ShowPopUp.Store(false)
			m.IsTyping.Store(false)
			m.PopUpType = constant.NoPopUp
			services.LockWorktreeService(m, popUp.WorktreePath, lockReason)
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
