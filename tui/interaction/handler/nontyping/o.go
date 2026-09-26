package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/component/branch"
	"github.com/gohyuhan/gitti/tui/component/commitlog"
	"github.com/gohyuhan/gitti/tui/component/reflog"
	"github.com/gohyuhan/gitti/tui/component/worktree"
	"github.com/gohyuhan/gitti/tui/constant"
	worktreePopUp "github.com/gohyuhan/gitti/tui/popup/worktree"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'o' key interaction.
//	Responsibility: Contextual "open" action.
//	- In Local Branch Panel: Opens the selected branch's web page in the browser.
//	- In Worktree Panel: Locks (prompts for a reason) or unlocks the selected worktree.
//	- In Commit Log / RefLog Panel: Opens the selected commit's web page in the browser.
//
// ------------------------------------
func handleNonTypingoKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				selectedBranch := m.CurrentRepoBranchesInfoList.SelectedItem()
				if selectedBranch != nil {
					services.OpenBranchWebPageService(m, selectedBranch.(branch.GitBranchItem).BranchName, git.WEBPAGEBRANCH)
				}
			case constant.SHOW_WORKTREE:
				selectedWorktreeItem := m.CurrentRepoWorktreeInfoList.SelectedItem()
				if selectedWorktreeItem != nil {
					parsedSelectedWorktree := selectedWorktreeItem.(worktree.GitWorktreeItem)
					// main worktree cannot be lock and unlock, silent return
					if parsedSelectedWorktree.IsMain {
						return m, nil
					}
					if parsedSelectedWorktree.IsLocked {
						services.UnlockWorktreeService(m, parsedSelectedWorktree.WorktreePath)
					} else {
						m.PopUpType = constant.WorktreeLockReasonInputPopUp
						m.IsTyping.Store(true)
						m.ShowPopUp.Store(true)
						worktreePopUp.InitWorktreeLockReasonInputPopUpModel(m, parsedSelectedWorktree.WorktreePath)
					}
				}
			}
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
				selectedCommitLog := m.CurrentRepoCommitLogInfoList.SelectedItem()
				if selectedCommitLog != nil {
					services.OpenCommitWebPageService(m, selectedCommitLog.(commitlog.GitCommitLogItem).Hash)
				}
			case constant.SHOW_REFLOG:
				selectedRefLog := m.CurrentRepoRefLogInfoList.SelectedItem()
				if selectedRefLog != nil {
					services.OpenCommitWebPageService(m, selectedRefLog.(reflog.GitRefLogItem).Hash)
				}
			}
		}
	}
	return m, nil
}
