package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/component/commitlog"
	"github.com/gohyuhan/gitti/tui/constant"
	commitLogPopUp "github.com/gohyuhan/gitti/tui/popup/commitlog"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle Ctrl+r key interaction.
//	Responsibility: Contextual "revert" mechanism.
//	Specifically triggered in the Commit Log Panel to initiate a `git revert`
//	of the currently selected commit, handling both standard commits and merge commits
//	(by optionally prompting for the parent parent-line to revert against).
//
// ------------------------------------
func handleNonTypingCtrlrKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
				selectedCommitLog := m.CurrentRepoCommitLogInfoList.SelectedItem()
				if selectedCommitLog != nil {
					parsedCommitLog := selectedCommitLog.(commitlog.GitCommitLogItem)
					commitHashParentInfos := services.GetCommitHashParentInfoService(m, parsedCommitLog.Hash)
					if commitHashParentInfos != nil {
						if len(commitHashParentInfos) > 1 {
							commitLogPopUp.InitGitRevertParentOptionSelectionPopUpModel(m, parsedCommitLog.Hash, commitHashParentInfos)
							m.ShowPopUp.Store(true)
							m.IsTyping.Store(false)
							m.PopUpType = constant.GitRevertParentOptionSelectionPopUp
						} else {
							commitLogPopUp.InitGitRevertConfirmationPopUpModel(m, parsedCommitLog.Hash, 0)
							m.ShowPopUp.Store(true)
							m.IsTyping.Store(false)
							m.PopUpType = constant.GitRevertConfirmationPopUp
						}
					} else {
						commitLogPopUp.InitGitRevertConfirmationPopUpModel(m, parsedCommitLog.Hash, 0)
						m.ShowPopUp.Store(true)
						m.IsTyping.Store(false)
						m.PopUpType = constant.GitRevertConfirmationPopUp
					}
				}
			}
		}
	}
	return m, nil
}
