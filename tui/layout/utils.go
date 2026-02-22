package layout

import (
	"github.com/charmbracelet/x/ansi"
	branchComponent "github.com/gohyuhan/gitti/tui/component/branch"
	commitlogComponent "github.com/gohyuhan/gitti/tui/component/commitlog"
	filesComponent "github.com/gohyuhan/gitti/tui/component/files"
	remoteComponent "github.com/gohyuhan/gitti/tui/component/remote"
	stashComponent "github.com/gohyuhan/gitti/tui/component/stash"
	tagComponent "github.com/gohyuhan/gitti/tui/component/tag"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ----------------------------------
//
//	To update the width and height of all components
//
// ----------------------------------
func TuiWindowSizing(m *types.GittiModel) {
	// Compute panel widths
	m.WindowLeftPanelWidth = int(float64(m.Width) * m.WindowLeftPanelRatio)
	m.DetailComponentPanelWidth = m.Width - m.WindowLeftPanelWidth

	m.WindowCoreContentHeight = m.Height - constant.MainPageKeyBindingLayoutPanelHeight - 2*constant.Padding

	// calculate the log component height based on the ratio of the window core content height
	logComponentHeight := int(float64(m.WindowCoreContentHeight) * constant.LogComponentHeightRatio)
	// if the calculated log component height is less than the minimum log component height,
	// set the log component height to the minimum log component height
	logComponentHeight = max(logComponentHeight, constant.MinLogComponentHeight)
	m.LogComponentPanelHeight = logComponentHeight
	m.DetailComponentPanelHeight = m.WindowCoreContentHeight - 2*constant.Padding - logComponentHeight

	// update the dynamic size of the left panel
	LeftPanelDynamicResize(m)

	// reconstruct the component title
	titleWidthLimit := m.WindowLeftPanelWidth - constant.ListItemOrTitleWidthPad - 2
	m.CurrentRepoBranchesInfoList.Title = ansi.Truncate(branchComponent.ConstructLocalBranchComponentTitle(titleWidthLimit), titleWidthLimit, "...")
	m.CurrentRepoTagInfoList.Title = ansi.Truncate(tagComponent.ConstructTagComponentTitle(titleWidthLimit), titleWidthLimit, "...")
	m.CurrentRepoRemoteInfoList.Title = ansi.Truncate(remoteComponent.ConstructRemoteComponentTitle(titleWidthLimit), titleWidthLimit, "...")
	m.CurrentRepoModifiedFilesInfoList.Title = ansi.Truncate(filesComponent.ConstructModifiedFilesComponentTitle(titleWidthLimit), titleWidthLimit, "...")
	m.CurrentRepoCommitLogInfoList.Title = ansi.Truncate(commitlogComponent.ConstructCommitLogComponentTitle(titleWidthLimit), titleWidthLimit, "...")
	m.CurrentRepoStashInfoList.Title = ansi.Truncate(stashComponent.ConstructStashComponentTitle(titleWidthLimit), titleWidthLimit, "...")

	// update viewport of detail panel
	UpdateDetailComponentViewportLayout(m)
	m.DetailPanelViewportOffset = max(0, int(m.DetailPanelViewport.HorizontalScrollPercent()*float64(m.DetailPanelViewportOffset))-1)
	m.DetailPanelTwoViewportOffset = max(0, int(m.DetailPanelTwoViewport.HorizontalScrollPercent()*float64(m.DetailPanelTwoViewportOffset))-1)
	m.DetailPanelViewport.SetXOffset(m.DetailPanelViewportOffset)
	m.DetailPanelViewport.SetYOffset(m.DetailPanelViewport.YOffset())
	m.DetailPanelTwoViewport.SetXOffset(m.DetailPanelTwoViewportOffset)
	m.DetailPanelTwoViewport.SetYOffset(m.DetailPanelTwoViewport.YOffset())

	// to recalculate the viewport of detail panel if it was in line editing mode so that
	// it matches exactly the position of the selected line in the viewport
	if m.IsLineEditingState.Load() {
		services.EnterOrReinitLineEditingStateService(m)
	}

	// log panel (the height is fixed)
	m.CurrentLogComponentViewport.SetWidth(m.DetailComponentPanelWidth - 2)
	m.CurrentLogComponentViewport.SetHeight(m.LogComponentPanelHeight)
}

// ----------------------------------
//
//	Dynamically resize the left panel components to fill the available height
//
// ----------------------------------
func LeftPanelDynamicResize(m *types.GittiModel) {
	var unSelectedComponentPanelHeightPerComponent int
	var selectedComponentPanelHeight int
	remainingHeight := 0
	// this is after reserving the height for the gitti status panel and also Padding
	leftPanelRemainingHeight := m.WindowCoreContentHeight - 1 - ((len(constant.ComponentPanelNavigationList) - 1) * 2)

	// we minus 2 if GitStatusComponent is not the one chosen is because GitStatusComponent
	// and the one that got selected will not be account in to the dynamic height calculation
	// ( gitti status component's height is fix at 3, while the selected one will always get 40% )
	componentWithDynamicHeight := (len(constant.ComponentPanelNavigationList) - 2)

	// because log component is not a member of left panel, so if this panel was selected, we need top adjust the component with dynamic height as no component in left panel is selected in that case
	if m.CurrentSelectedComponent == constant.LogComponentPanel {
		componentWithDynamicHeight = (len(constant.ComponentPanelNavigationList) - 1)
		unSelectedComponentPanelHeightPerComponent = int(leftPanelRemainingHeight / componentWithDynamicHeight)
		// there will be possibility of some height remaining after divide and turn to int, so we use the original - (divided height*the count of component)
		// the remainingHeight will be added to the height for modified file component
		remainingHeight = leftPanelRemainingHeight - (unSelectedComponentPanelHeightPerComponent * componentWithDynamicHeight)
	} else {
		unSelectedComponentPanelHeightPerComponent = int(int(float64(leftPanelRemainingHeight)*(1.0-constant.SelectedLeftPanelComponentHeightRatio)) / componentWithDynamicHeight)
		selectedComponentPanelHeight = leftPanelRemainingHeight - (unSelectedComponentPanelHeightPerComponent * componentWithDynamicHeight)
	}

	if unSelectedComponentPanelHeightPerComponent < constant.MinUnSelectedComponentPanelHeight {
		unSelectedComponentPanelHeightPerComponent = constant.MinUnSelectedComponentPanelHeight
		selectedComponentPanelHeight = leftPanelRemainingHeight - (unSelectedComponentPanelHeightPerComponent * componentWithDynamicHeight)
	}

	m.LocalBranchesComponentPanelHeight = unSelectedComponentPanelHeightPerComponent
	m.TagsComponentPanelHeight = unSelectedComponentPanelHeightPerComponent
	m.RemoteComponentPanelHeight = unSelectedComponentPanelHeightPerComponent
	m.ModifiedFilesComponentPanelHeight = unSelectedComponentPanelHeightPerComponent + remainingHeight
	m.CommitLogComponentPanelHeight = unSelectedComponentPanelHeightPerComponent
	m.StashComponentPanelHeight = unSelectedComponentPanelHeightPerComponent

	switch m.CurrentSelectedComponent {
	case constant.LocalBranchOrTagOrRemoteComponentPanel:
		m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
		m.TagsComponentPanelHeight = selectedComponentPanelHeight
		m.RemoteComponentPanelHeight = selectedComponentPanelHeight
	case constant.ModifiedFilesComponentPanel:
		m.ModifiedFilesComponentPanelHeight = selectedComponentPanelHeight
	case constant.CommitLogComponentPanel:
		m.CommitLogComponentPanelHeight = selectedComponentPanelHeight
	case constant.StashComponentPanel:
		m.StashComponentPanelHeight = selectedComponentPanelHeight
	case constant.GitStatusComponentPanel:
		// if it was the Gitti status component panel that got selected (because its height is fix),
		// the next panel will get the selected height which is the branch component panel
		m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
	case constant.DetailComponentPanelTwo:
		switch m.DetailPanelParentComponent {
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
			m.TagsComponentPanelHeight = selectedComponentPanelHeight
			m.RemoteComponentPanelHeight = selectedComponentPanelHeight
		case constant.ModifiedFilesComponentPanel:
			m.ModifiedFilesComponentPanelHeight = selectedComponentPanelHeight
		case constant.CommitLogComponentPanel:
			m.CommitLogComponentPanelHeight = selectedComponentPanelHeight
		case constant.StashComponentPanel:
			m.StashComponentPanelHeight = selectedComponentPanelHeight
		}
	case constant.DetailComponentPanel:
		switch m.DetailPanelParentComponent {
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
			m.TagsComponentPanelHeight = selectedComponentPanelHeight
			m.RemoteComponentPanelHeight = selectedComponentPanelHeight
		case constant.ModifiedFilesComponentPanel:
			m.ModifiedFilesComponentPanelHeight = selectedComponentPanelHeight
		case constant.CommitLogComponentPanel:
			m.CommitLogComponentPanelHeight = selectedComponentPanelHeight
		case constant.StashComponentPanel:
			m.StashComponentPanelHeight = selectedComponentPanelHeight
		}
	}

	// update all components Width and Height
	m.CurrentRepoBranchesInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoBranchesInfoList.SetHeight(m.LocalBranchesComponentPanelHeight)

	m.CurrentRepoTagInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoTagInfoList.SetHeight(m.TagsComponentPanelHeight)

	m.CurrentRepoRemoteInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoRemoteInfoList.SetHeight(m.RemoteComponentPanelHeight)

	m.CurrentRepoModifiedFilesInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoModifiedFilesInfoList.SetHeight(m.ModifiedFilesComponentPanelHeight)

	m.CurrentRepoCommitLogInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoCommitLogInfoList.SetHeight(m.CommitLogComponentPanelHeight)

	m.CurrentRepoStashInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoStashInfoList.SetHeight(m.StashComponentPanelHeight)
}

// ------------------------------------
//
//			For Update Detail Component Viewport Layout
//		  * this was to update the layout for detail component viewport
//	   * it will handle the layout for both line editing mode and normal mode
//	   * it will also handle the split view for line editing mode (one for staged, one for unstaged)
//
// ------------------------------------
func UpdateDetailComponentViewportLayout(m *types.GittiModel) {
	availableHeight := m.DetailComponentPanelHeight
	// set the cursor viewport width
	m.LineEditingIndexCursorViewport.SetWidth(3)
	m.LineEditingIndexCursorTwoViewport.SetWidth(3)

	// if it is in line editing mode, we need to minus 3 for the title
	if m.IsLineEditingState.Load() {
		availableHeight -= 3
	}

	if m.ShowDetailPanelTwo.Load() {
		// vertical layout
		// Since terminal characters are usually about twice as tall as they are wide,
		// we weight the height by 2 to approximate visual "squareness".
		splitHeight := int(availableHeight / 2)
		splitWidth := int(m.DetailComponentPanelWidth / 2)

		if availableHeight*2 > m.DetailComponentPanelWidth {
			// Vertical split (Top/Bottom)
			m.DetailComponentPanelLayout = constant.VERTICAL
			m.DetailPanelViewport.SetHeight(splitHeight - 1)
			m.DetailPanelTwoViewport.SetHeight(availableHeight - splitHeight - 1)
			m.LineEditingIndexCursorViewport.SetHeight(splitHeight - 1)
			m.LineEditingIndexCursorTwoViewport.SetHeight(availableHeight - splitHeight - 1)

			// Adjust width based on mode
			if m.IsLineEditingState.Load() {
				m.DetailPanelViewport.SetWidth(m.DetailComponentPanelWidth - 2 - 3)
				m.DetailPanelTwoViewport.SetWidth(m.DetailComponentPanelWidth - 2 - 3)
			} else {
				m.DetailPanelViewport.SetWidth(m.DetailComponentPanelWidth - 2)
				m.DetailPanelTwoViewport.SetWidth(m.DetailComponentPanelWidth - 2)
			}
		} else {
			// Horizontal split (Left/Right)
			m.DetailComponentPanelLayout = constant.HORIZONTAL
			m.DetailPanelViewport.SetHeight(availableHeight)
			m.DetailPanelTwoViewport.SetHeight(availableHeight)
			m.LineEditingIndexCursorViewport.SetHeight(availableHeight)
			m.LineEditingIndexCursorTwoViewport.SetHeight(availableHeight)

			// Adjust width based on mode
			if m.IsLineEditingState.Load() {
				m.DetailPanelViewport.SetWidth(splitWidth - 2 - 3)
				m.DetailPanelTwoViewport.SetWidth(m.DetailComponentPanelWidth - splitWidth - 2 - 3)
			} else {
				m.DetailPanelViewport.SetWidth(splitWidth - 2)
				m.DetailPanelTwoViewport.SetWidth(m.DetailComponentPanelWidth - splitWidth - 2)
			}
		}
	} else {
		// Single Panel View
		m.DetailPanelViewport.SetHeight(availableHeight)
		m.LineEditingIndexCursorViewport.SetHeight(availableHeight)

		// Adjust width based on mode
		if m.IsLineEditingState.Load() {
			m.DetailPanelViewport.SetWidth(m.DetailComponentPanelWidth - 2 - 3)
		} else {
			m.DetailPanelViewport.SetWidth(m.DetailComponentPanelWidth - 2)
		}
	}
}
