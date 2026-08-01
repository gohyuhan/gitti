package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/component/branch"
	"github.com/gohyuhan/gitti/tui/component/commitlog"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/component/reflog"
	"github.com/gohyuhan/gitti/tui/constant"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	remotePopUp "github.com/gohyuhan/gitti/tui/popup/remote"
	resolvePopUp "github.com/gohyuhan/gitti/tui/popup/resolve"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'r' key interaction.
//	Responsibility: Contextual "resolve" or "reset" operations depending on the view:
//	- In Local Branch View: Opens a popup to rebase the current branch to a specific commit. (the selector must be on the current checkout branch to trigger this)
//	- In Modified Files Panel: Opens the conflict resolution popup for a conflicted file.
//	- In Commit Log Panel: Opens a popup to reset the repository HEAD to a specifically
//	  selected older commit (soft, mixed, or hard reset).
//
// ------------------------------------
func handleNonTypingrKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				selectedBranch := m.CurrentRepoBranchesInfoList.SelectedItem()
				if selectedBranch != nil {
					parsedBranch := selectedBranch.(branch.GitBranchItem)
					if parsedBranch.IsCheckedOut {
						_ = m.GitOperations.GitRemote.CheckRemoteExist(false)
						m.IsTyping.Store(false)
						m.ShowPopUp.Store(true)
						remotes := m.GitOperations.GitRemote.FetchRemote()
						m.PopUpType = constant.ChooseRemotePopUp
						// add a use local branch option which the name and url is "" (empty string) and fetch/push is false
						remotes = append([]git.GitRemoteInfo{{Name: "", Url: "", Fetch: false, Push: false}}, remotes...)
						remotePopUp.InitChooseRemotePopUpModel(m, remotes, constant.REBASEACTION)
					} else {
						return m, nil
					}
				}
			}
		case constant.ModifiedFilesComponentPanel:
			currentSelectedFileItem := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
			if currentSelectedFileItem != nil {
				currentSelectedFile := currentSelectedFileItem.(files.GitModifiedFilesItem)
				// return early if the file has no conflict
				if !currentSelectedFile.HasConflict {
					return m, nil
				}
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				m.PopUpType = constant.GitResolveConflictOptionPopUp
				resolvePopUp.InitGitResolveConflictOptionPopUpModel(m, currentSelectedFile.FilePathname)
			}
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
				selectedCommit := m.CurrentRepoCommitLogInfoList.SelectedItem()
				if selectedCommit != nil {
					parsedCommit := selectedCommit.(commitlog.GitCommitLogItem)
					commitPopUp.InitGitResetToSelectedCommitTypeOptionPopUpModel(
						m,
						parsedCommit.Hash,
						parsedCommit.Message,
						parsedCommit.Author,
					)
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					m.PopUpType = constant.GitResetToSelectedCommitTypeOptionPopUp
				}
			case constant.SHOW_REFLOG:
				selectedRefLog := m.CurrentRepoRefLogInfoList.SelectedItem()
				if selectedRefLog != nil {
					parsedRefLog := selectedRefLog.(reflog.GitRefLogItem)
					commitPopUp.InitGitResetToSelectedCommitTypeOptionPopUpModel(
						m,
						parsedRefLog.Hash,
						"",
						"",
					)
					m.ShowPopUp.Store(true)
					m.IsTyping.Store(false)
					m.PopUpType = constant.GitResetToSelectedCommitTypeOptionPopUp
				}
			}
		}
	}
	return m, nil
}
