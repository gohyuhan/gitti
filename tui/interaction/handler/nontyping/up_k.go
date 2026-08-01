package nontyping

import (
	"github.com/gohyuhan/gitti/tui/interaction/handler/keyutil"

	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/layout"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle Up/k key interaction.
//	Responsibility: Vertical upward navigation.
//	- Native lists: Moves selection highlight up one item.
//	- Detail Viewports (text viewers & line-editing): Scrolls text up by one line.
//	- Popups: Navigates upward in popup selection lists natively.
//
// ------------------------------------
func handleNonTypingUpkKeyBindingInteraction(msg tea.KeyPressMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
			// we don't use the list native Update() because we track the current selected index
			switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				if m.CurrentRepoBranchesInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoBranchesInfoList.Index() - 1
					m.CurrentRepoBranchesInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.LocalBranchComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			case constant.SHOW_TAG:
				if m.CurrentRepoTagInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoTagInfoList.Index() - 1
					m.CurrentRepoTagInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.TagComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			case constant.SHOW_REMOTE:
				if m.CurrentRepoRemoteInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoRemoteInfoList.Index() - 1
					m.CurrentRepoRemoteInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.RemoteComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			case constant.SHOW_WORKTREE:
				if m.CurrentRepoWorktreeInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoWorktreeInfoList.Index() - 1
					m.CurrentRepoWorktreeInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.WorktreeComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			}
		case constant.ModifiedFilesComponentPanel:
			// we don't use the list native Update() because we need to also track the current selected index
			if m.CurrentRepoModifiedFilesInfoList.Index() > 0 {
				latestIndex := m.CurrentRepoModifiedFilesInfoList.Index() - 1
				m.CurrentRepoModifiedFilesInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.ModifiedFilesComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
				// we don't use the list native Update() because we need to also track the current selected index
				if m.CurrentRepoCommitLogInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoCommitLogInfoList.Index() - 1
					m.CurrentRepoCommitLogInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.CommitLogComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			case constant.SHOW_REFLOG:
				// we don't use the list native Update() because we need to also track the current selected index
				if m.CurrentRepoRefLogInfoList.Index() > 0 {
					latestIndex := m.CurrentRepoRefLogInfoList.Index() - 1
					m.CurrentRepoRefLogInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.RefLogComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			}
		case constant.StashComponentPanel:
			// we don't use the list native Update() because we need to also track the current selected index
			if m.CurrentRepoStashInfoList.Index() > 0 {
				latestIndex := m.CurrentRepoStashInfoList.Index() - 1
				m.CurrentRepoStashInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.StashComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.DetailComponentPanel:
			if m.IsLineEditingState.Load() {
				m.LineEditingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex = max(0, m.LineEditingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex-1)
				if m.LineEditingIndexPositionAndInfo.DetailPanelViewportIndexPosition < 1 {
					m.DetailPanelViewport.ScrollUp(1)
				} else {
					m.LineEditingIndexPositionAndInfo.DetailPanelViewportIndexPosition -= 1
				}
				layout.SetLineEditingCursorViewportContent(m, m.DetailPanelViewport.VisibleLineCount(), m.DetailPanelTwoViewport.VisibleLineCount())
			} else {
				m.DetailPanelViewport.ScrollUp(1)
				return m, nil
			}
		case constant.DetailComponentPanelTwo:
			if m.IsLineEditingState.Load() {
				m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportActualCurrentIndex = max(0, m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportActualCurrentIndex-1)
				if m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportIndexPosition < 1 {
					m.DetailPanelTwoViewport.ScrollUp(1)
				} else {
					m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportIndexPosition -= 1
				}
				layout.SetLineEditingCursorViewportContent(m, m.DetailPanelViewport.VisibleLineCount(), m.DetailPanelTwoViewport.VisibleLineCount())
			} else {
				m.DetailPanelTwoViewport.ScrollUp(1)
				return m, nil
			}
		}
	} else {
		return keyutil.UpDownKeyPressMsgUpdateForPopUp(msg, m)
	}
	return m, nil
}
