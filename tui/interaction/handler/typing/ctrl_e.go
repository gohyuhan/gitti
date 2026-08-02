package typing

import (
	"context"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/constant"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

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
