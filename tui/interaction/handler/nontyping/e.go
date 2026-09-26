package nontyping

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/tui/component/branch"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/component/remote"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitLogPopUp "github.com/gohyuhan/gitti/tui/popup/commitlog"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'e' key interaction.
//	Responsibility: Contextual "edit" action.
//	- In popups: Switches a cherry-pick view into "edit cherry-pick" mode.
//	- In Local Branch Panel: Opens the prompt to rename the selected local branch.
//	- In Remote Panel: Opens the prompt to edit an existing remote's URL/Name.
//	- In Modified Files Panel: Launches the user's defined system editor (handling terminal/GUI diffs).
//	- In Log Panel: Triggers exporting the internal application logs.
//
// ------------------------------------
func handleNonTypingeKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.GitCherryPickPopUp, constant.GitCherryPickApplyConfirmPopUp:
			m.PopUpType = constant.GitEditCherryPickPopUp
			commitLogPopUp.InitGitEditCherryPickPopUpModel(m, 0)
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(false)
		}
	} else {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				selectedBranch := m.CurrentRepoBranchesInfoList.SelectedItem()
				if selectedBranch != nil {
					// renaming in the middle of a merge/rebase/cherry-pick/revert could break that operation
					if m.CurrentGitRepoStatus != "" {
						m.GittiLogger.RegisterNewLog(logging.RENAME_LOCAL_BRANCH_OPS, "", logging.WARN, fmt.Sprintf("Cannot rename branch while %s is in progress", m.CurrentGitRepoStatus), false)
						return m, nil
					}
					branchItem := selectedBranch.(branch.GitBranchItem)
					branchPopUp.InitRenameBranchPopUpModel(m, branchItem.BranchName)
					m.PopUpType = constant.CreateNewBranchPopUp
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(true)
				}
			case constant.SHOW_REMOTE:
				selectedRemote := m.CurrentRepoRemoteInfoList.SelectedItem()
				if selectedRemote != nil {
					remoteItem := selectedRemote.(remote.GitRemoteItem)
					m.PopUpType = constant.EditRemotePromptPopUp
					remotePopUp.InitEditRemotePromptPopUpModel(m, remoteItem.Name, remoteItem.Url)
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(true)
				}
			}
		case constant.ModifiedFilesComponentPanel:
			currentSelectedFileItem := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
			if currentSelectedFileItem != nil {
				currentSelectedFile := currentSelectedFileItem.(files.GitModifiedFilesItem)
				return launchEditor(m, currentSelectedFile.FilePathname)
			}
		case constant.LogComponentPanel:
			go func() {
				m.GittiLogger.ExportLogging()
			}()
		}
	}
	return m, nil
}
