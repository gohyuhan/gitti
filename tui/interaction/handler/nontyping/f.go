package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'f' key interaction.
//	Responsibility: Initiates the "fetch" operation. Specifically checks if the user is in the
//	Tag component view. If no remotes exist, prompts to add one. If remotes exist,
//	either directly fetches tags (if 1 remote) or opens a popup to select which remote to fetch tags from.
//
// ------------------------------------
func handleNonTypingfKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if (m.CurrentSelectedComponent == constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel || m.DetailPanelParentComponent == constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel) &&
			m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing == constant.SHOW_TAG {
			if !m.GitOperations.GitRemote.CheckRemoteExist(false) {
				// if no remote found, we add one
				showAddRemotePromptPopUp(m)
			} else {
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				remotes := m.GitOperations.GitRemote.FetchRemote()
				if len(remotes) == 1 {
					m.PopUpType = constant.ChooseFetchTagOptionPopUp
					// only one remote found so, we will default to that remote
					tagPopUp.InitChooseFetchTagOptionPopUpModel(m, remotes[0].Name)
				} else if len(remotes) > 1 {
					// if remote is more than 1 let user choose which remote
					m.PopUpType = constant.ChooseRemotePopUp
					remotePopUp.InitChooseRemotePopUpModel(m, remotes, constant.TAGFETCHACTION)
				}
			}
		}
	}
	return m, nil
}
