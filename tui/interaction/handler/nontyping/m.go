package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'm' key interaction.
//	Responsibility: Initiates the "Merge" workflow.
//	In the Local Branch panel, opens the branch selection popup allowing the user to
//	choose one or more branches to merge into the currently checked-out branch.
//
// ------------------------------------
func handleNonTypingmKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				m.PopUpType = constant.ChooseBranchOptionForMergePopUp
				m.IsTyping.Store(false)
				m.ShowPopUp.Store(true)
				branchPopUp.InitChooseBranchOptionForMergePopUpModel(m)
			}
		}
	}
	return m, nil
}
