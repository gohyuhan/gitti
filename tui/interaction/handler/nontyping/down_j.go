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
//	Handle Down/j key interaction.
//	Responsibility: Vertical downward navigation.
//	- Native lists: Moves selection highlight down one item.
//	- Detail Viewports (text viewers & line-editing): Scrolls text down by one line.
//	- Popups: Navigates downward in popup selection lists natively.
//
// ------------------------------------
func handleNonTypingDownjKeyBindingInteraction(msg tea.KeyPressMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
			// we don't use the list native Update() because we track the current selected index
			switch m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				if m.CurrentRepoBranchesInfoList.Index() < len(m.CurrentRepoBranchesInfoList.Items())-1 {
					latestIndex := m.CurrentRepoBranchesInfoList.Index() + 1
					m.CurrentRepoBranchesInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.LocalBranchComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			case constant.SHOW_TAG:
				if m.CurrentRepoTagInfoList.Index() < len(m.CurrentRepoTagInfoList.Items())-1 {
					latestIndex := m.CurrentRepoTagInfoList.Index() + 1
					m.CurrentRepoTagInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.TagComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			case constant.SHOW_REMOTE:
				if m.CurrentRepoRemoteInfoList.Index() < len(m.CurrentRepoRemoteInfoList.Items())-1 {
					latestIndex := m.CurrentRepoRemoteInfoList.Index() + 1
					m.CurrentRepoRemoteInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.RemoteComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			case constant.SHOW_WORKTREE:
				if m.CurrentRepoWorktreeInfoList.Index() < len(m.CurrentRepoWorktreeInfoList.Items())-1 {
					latestIndex := m.CurrentRepoWorktreeInfoList.Index() + 1
					m.CurrentRepoWorktreeInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.WorktreeComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			}
		case constant.ModifiedFilesComponentPanel:
			// we don't use the list native Update() because we need to also track the current selected index
			if m.CurrentRepoModifiedFilesInfoList.Index() < len(m.CurrentRepoModifiedFilesInfoList.Items())-1 {
				latestIndex := m.CurrentRepoModifiedFilesInfoList.Index() + 1
				m.CurrentRepoModifiedFilesInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.ModifiedFilesComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.CommitLogOrRefLogComponentPanel:
			switch m.CurrentCommitLogOrRefLogComponentShowing {
			case constant.SHOW_COMMITLOG:
				// we don't use the list native Update() because we need to also track the current selected index
				if m.CurrentRepoCommitLogInfoList.Index() < len(m.CurrentRepoCommitLogInfoList.Items())-1 {
					latestIndex := m.CurrentRepoCommitLogInfoList.Index() + 1
					m.CurrentRepoCommitLogInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.CommitLogComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			case constant.SHOW_REFLOG:
				// we don't use the list native Update() because we need to also track the current selected index
				if m.CurrentRepoRefLogInfoList.Index() < len(m.CurrentRepoRefLogInfoList.Items())-1 {
					latestIndex := m.CurrentRepoRefLogInfoList.Index() + 1
					m.CurrentRepoRefLogInfoList.Select(latestIndex)
					m.ListNavigationIndexPosition.RefLogComponent = latestIndex
					services.FetchDetailComponentPanelInfoService(m, true)
				}
			}
		case constant.StashComponentPanel:
			// we don't use the list native Update() because we need to also track the current selected index
			if m.CurrentRepoStashInfoList.Index() < len(m.CurrentRepoStashInfoList.Items())-1 {
				latestIndex := m.CurrentRepoStashInfoList.Index() + 1
				m.CurrentRepoStashInfoList.Select(latestIndex)
				m.ListNavigationIndexPosition.StashComponent = latestIndex
				services.FetchDetailComponentPanelInfoService(m, true)
			}
		case constant.DetailComponentPanel:
			if m.IsLineEditingState.Load() {
				m.LineEditingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex = min(m.DetailPanelViewport.TotalLineCount()-1, m.LineEditingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex+1)
				if m.LineEditingIndexPositionAndInfo.DetailPanelViewportIndexPosition >= m.DetailPanelViewport.VisibleLineCount()-1 {
					m.DetailPanelViewport.ScrollDown(1)
				} else {
					m.LineEditingIndexPositionAndInfo.DetailPanelViewportIndexPosition += 1
				}
				layout.SetLineEditingCursorViewportContent(m, m.DetailPanelViewport.VisibleLineCount(), m.DetailPanelTwoViewport.VisibleLineCount())
			} else {
				m.DetailPanelViewport.ScrollDown(1)
				return m, nil
			}
		case constant.DetailComponentPanelTwo:
			if m.IsLineEditingState.Load() {
				m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportActualCurrentIndex = min(m.DetailPanelTwoViewport.TotalLineCount()-1, m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportActualCurrentIndex+1)
				if m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportIndexPosition >= m.DetailPanelTwoViewport.VisibleLineCount()-1 {
					m.DetailPanelTwoViewport.ScrollDown(1)
				} else {
					m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportIndexPosition += 1
				}
				layout.SetLineEditingCursorViewportContent(m, m.DetailPanelViewport.VisibleLineCount(), m.DetailPanelTwoViewport.VisibleLineCount())
			} else {
				m.DetailPanelTwoViewport.ScrollDown(1)
				return m, nil
			}
		}
	} else {
		return keyutil.UpDownKeyPressMsgUpdateForPopUp(msg, m)
	}
	return m, nil
}
