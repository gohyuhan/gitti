package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'R' key interaction.
//	Responsibility: Contextual "reset latest" operation.
//	Specifically in the Commit Log Panel, this opens a popup allowing the user
//	to reset the repository HEAD exactly one commit backwards (undoing the latest commit),
//	offering soft, mixed, or hard reset options.
//
// ------------------------------------
func handleNonTypingRKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
				if len(m.CurrentRepoCommitLogInfoList.Items()) > 1 {
					commitPopUp.InitGitResetLatestCommitTypeOptionPopUpModel(m)
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					m.PopUpType = constant.GitResetLatestCommitTypeOptionPopUp
				}
			}
		}
	}
	return m, nil
}
