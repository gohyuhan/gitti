package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/component/branch"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle Ctrl+o key interaction.
//	Responsibility: In the Local Branch Panel, opens the new pull request (GitLab: merge request)
//	page of the selected branch in the browser.
//
// ------------------------------------
func handleNonTypingCtrloKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() &&
		m.CurrentSelectedComponent == constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel &&
		m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing == constant.SHOW_LOCAL_BRANCH {
		selectedBranch := m.CurrentRepoBranchesInfoList.SelectedItem()
		if selectedBranch != nil {
			services.OpenBranchWebPageService(m, selectedBranch.(branch.GitBranchItem).BranchName, git.WEBPAGEPULLREQUEST)
		}
	}
	return m, nil
}
