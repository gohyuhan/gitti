package interaction

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/layout"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle a left mouse click on the main page. Maps the click coordinates to a
//	panel (left column panels stack vertically, right column is detail above
//	log), focuses that panel like its keybinding would, and for the left list
//	panels also selects the clicked list row. Ignored while a popup, line
//	editing, or panel filter typing is active.
//
// ------------------------------------
func handleLeftMouseClick(msg tea.MouseClickMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.ShowPopUp.Load() || m.IsLineEditingState.Load() || m.IsPanelFiltering.Load() {
		return m, nil
	}

	x := msg.Mouse().X
	y := msg.Mouse().Y

	// left column: each panel renders with a 1-cell border on every side,
	// so a panel's total footprint is its content height + 2
	if x <= m.WindowLeftPanelWidth+1 {
		leftPanels := []struct {
			component string
			index     int
			height    int
		}{
			{constant.GitStatusComponentPanel, 0, 1 + 2},
			{constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel, 1, m.LocalBranchesComponentPanelHeight + 2},
			{constant.ModifiedFilesComponentPanel, 2, m.ModifiedFilesComponentPanelHeight + 2},
			{constant.CommitLogOrRefLogComponentPanel, 3, m.CommitLogComponentPanelHeight + 2},
			{constant.StashComponentPanel, 4, m.StashComponentPanelHeight + 2},
		}

		top := 0
		for _, panel := range leftPanels {
			if y < top+panel.height {
				// row within the list body: skip the top border and the list title row.
				// selection must be resolved before focusing, as focusing resizes the panels
				selectionChanged := selectListItemFromClick(m, panel.component, y-top-2)
				if m.CurrentSelectedComponent != panel.component {
					m.CurrentSelectedComponent = panel.component
					m.CurrentSelectedComponentIndex = panel.index
					m.DetailPanelParentComponent = ""
					layout.LeftPanelDynamicResize(m)
					services.FetchDetailComponentPanelInfoService(m, true)
				} else if selectionChanged {
					services.FetchDetailComponentPanelInfoService(m, true)
				}
				return m, nil
			}
			top += panel.height
		}
		return m, nil
	}

	// right column: detail panel on top, log panel below
	if y >= m.DetailComponentPanelHeight {
		if m.CurrentSelectedComponent != constant.LogComponentPanel {
			m.CurrentSelectedComponent = constant.LogComponentPanel
			m.DetailPanelParentComponent = ""
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	} else {
		// entering the detail panel keeps the parent so esc returns to it,
		// mirroring the enter keybinding's drill-in behavior
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel,
			constant.ModifiedFilesComponentPanel,
			constant.CommitLogOrRefLogComponentPanel,
			constant.StashComponentPanel,
			constant.LogComponentPanel:
			m.DetailPanelParentComponent = m.CurrentSelectedComponent
			m.CurrentSelectedComponent = constant.DetailComponentPanel
		}
	}
	return m, nil
}

// ------------------------------------
//
//	Select the list row of the clicked left panel. itemRow is the row within the
//	list body (0 = first visible item); the absolute index is resolved through
//	the list's paginator page. Returns true when the selection changed.
//
// ------------------------------------
func selectListItemFromClick(m *types.GittiModel, component string, itemRow int) bool {
	if itemRow < 0 {
		return false
	}

	var clickedList *list.Model
	var navigationIndex *int
	switch component {
	case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
		switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
		case constant.SHOW_LOCAL_BRANCH:
			clickedList = &m.CurrentRepoBranchesInfoList
			navigationIndex = &m.ListNavigationIndexPosition.LocalBranchComponent
		case constant.SHOW_TAG:
			clickedList = &m.CurrentRepoTagInfoList
			navigationIndex = &m.ListNavigationIndexPosition.TagComponent
		case constant.SHOW_REMOTE:
			clickedList = &m.CurrentRepoRemoteInfoList
			navigationIndex = &m.ListNavigationIndexPosition.RemoteComponent
		case constant.SHOW_WORKTREE:
			clickedList = &m.CurrentRepoWorktreeInfoList
			navigationIndex = &m.ListNavigationIndexPosition.WorktreeComponent
		}
	case constant.ModifiedFilesComponentPanel:
		clickedList = &m.CurrentRepoModifiedFilesInfoList
		navigationIndex = &m.ListNavigationIndexPosition.ModifiedFilesComponent
	case constant.CommitLogOrRefLogComponentPanel:
		switch m.CurrentCommitLogOrRefLogComponentShowing {
		case constant.SHOW_COMMITLOG:
			clickedList = &m.CurrentRepoCommitLogInfoList
			navigationIndex = &m.ListNavigationIndexPosition.CommitLogComponent
		case constant.SHOW_REFLOG:
			clickedList = &m.CurrentRepoRefLogInfoList
			navigationIndex = &m.ListNavigationIndexPosition.RefLogComponent
		}
	case constant.StashComponentPanel:
		clickedList = &m.CurrentRepoStashInfoList
		navigationIndex = &m.ListNavigationIndexPosition.StashComponent
	}
	if clickedList == nil {
		return false
	}

	totalCount := len(clickedList.Items())
	page := clickedList.Paginator.Page
	perPage := clickedList.Paginator.PerPage
	itemsOnPage := min(perPage, totalCount-page*perPage)
	if itemRow >= itemsOnPage {
		return false
	}

	index := page*perPage + itemRow
	if index < 0 || index >= totalCount || index == clickedList.Index() {
		return false
	}

	clickedList.Select(index)
	*navigationIndex = index
	return true
}
