package interaction

import (
	tea "charm.land/bubbletea/v2"
	branchComponent "github.com/gohyuhan/gitti/tui/component/branch"
	commitlogComponent "github.com/gohyuhan/gitti/tui/component/commitlog"
	filesComponent "github.com/gohyuhan/gitti/tui/component/files"
	reflogComponent "github.com/gohyuhan/gitti/tui/component/reflog"
	remoteComponent "github.com/gohyuhan/gitti/tui/component/remote"
	stashComponent "github.com/gohyuhan/gitti/tui/component/stash"
	tagComponent "github.com/gohyuhan/gitti/tui/component/tag"
	worktreeComponent "github.com/gohyuhan/gitti/tui/component/worktree"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	Handle key presses while the user is typing a panel list filter query
//	(entered with 'F'). Other keys edit the query through the panel filter text
//	input, enter exits typing mode keeping the query applied, and esc clears the
//	query and exits. The focused list is rebuilt whenever the query changes.
//
// ------------------------------------
func handlePanelFilterKeyInput(msg tea.KeyPressMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	filterKey := utils.CurrentPanelFilterKey(m)
	if filterKey == "" {
		m.IsPanelFiltering.Store(false)
		return m, nil
	}

	switch msg.String() {
	case "enter":
		m.IsPanelFiltering.Store(false)
	case "esc":
		if m.PanelFilterQuery[filterKey] != "" {
			delete(m.PanelFilterQuery, filterKey)
			reinitFilteredList(m, filterKey)
		}
		m.IsPanelFiltering.Store(false)
	default:
		var cmd tea.Cmd
		m.PanelFilterInput, cmd = m.PanelFilterInput.Update(msg)
		if query := m.PanelFilterInput.Value(); query != m.PanelFilterQuery[filterKey] {
			m.PanelFilterQuery[filterKey] = query
			reinitFilteredList(m, filterKey)
		}
		return m, cmd
	}
	return m, nil
}

// ------------------------------------
//
//	Rebuild the list identified by filterKey so the current filter query is
//	applied, then refresh the detail panel to follow the new selection. Mirrors
//	the per-event reinit calls in tui.go's GitUpdateMsg handling.
//
// ------------------------------------
func reinitFilteredList(m *types.GittiModel, filterKey string) {
	switch filterKey {
	case constant.SHOW_LOCAL_BRANCH:
		branchComponent.InitBranchList(m)
		services.FetchDetailComponentPanelInfoService(m, false)
	case constant.SHOW_TAG:
		needReinit := tagComponent.InitTagList(m)
		services.FetchDetailComponentPanelInfoService(m, needReinit)
	case constant.SHOW_REMOTE:
		needReinit := remoteComponent.InitRemoteList(m)
		services.FetchDetailComponentPanelInfoService(m, needReinit)
	case constant.SHOW_WORKTREE:
		needReinit := worktreeComponent.InitWorktreeList(m)
		services.FetchDetailComponentPanelInfoService(m, needReinit)
	case constant.ModifiedFilesComponentPanel:
		needReinit := filesComponent.InitModifiedFilesList(m)
		services.FetchDetailComponentPanelInfoService(m, needReinit)
	case constant.SHOW_COMMITLOG:
		needReinit := commitlogComponent.InitGitCommitLogList(m)
		services.FetchDetailComponentPanelInfoService(m, needReinit)
	case constant.SHOW_REFLOG:
		needReinit := reflogComponent.InitGitRefLogList(m)
		services.FetchDetailComponentPanelInfoService(m, needReinit)
	case constant.StashComponentPanel:
		needReinit := stashComponent.InitStashList(m)
		services.FetchDetailComponentPanelInfoService(m, needReinit)
	}
}
