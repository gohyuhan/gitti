package interaction

import (
	"unicode/utf8"

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
//	(entered with 'F'). Printable runes append to the query, backspace deletes
//	the last rune, enter exits typing mode keeping the query applied, and esc
//	clears the query and exits. The focused list is rebuilt on every edit.
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
	case "backspace":
		query := []rune(m.PanelFilterQuery[filterKey])
		if len(query) > 0 {
			m.PanelFilterQuery[filterKey] = string(query[:len(query)-1])
			reinitFilteredList(m, filterKey)
		}
	case "space":
		m.PanelFilterQuery[filterKey] += " "
		reinitFilteredList(m, filterKey)
	default:
		key := msg.String()
		if utf8.RuneCountInString(key) == 1 {
			m.PanelFilterQuery[filterKey] += key
			reinitFilteredList(m, filterKey)
		}
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
