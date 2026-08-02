package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/component/remote"
	"github.com/gohyuhan/gitti/tui/constant"
	commitLogPopUp "github.com/gohyuhan/gitti/tui/popup/commitlog"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'e' key interaction.
//	Responsibility: Contextual "edit" action.
//	- In popups: Switches a cherry-pick view into "edit cherry-pick" mode.
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
