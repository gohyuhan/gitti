package layout

import (
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// to update the width and height of all components
func TuiWindowSizing(m *types.GittiModel) {
	// Compute panel widths
	m.WindowLeftPanelWidth = int(float64(m.Width) * m.WindowLeftPanelRatio)
	m.DetailComponentPanelWidth = m.Width - m.WindowLeftPanelWidth

	m.WindowCoreContentHeight = m.Height - constant.MainPageKeyBindingLayoutPanelHeight - 2*constant.Padding
	m.DetailComponentPanelHeight = m.WindowCoreContentHeight

	// update the dynamic size of the left panel
	LeftPanelDynamicResize(m)

	// update viewport
	UpdateDetailComponentViewportLayout(m)
	m.DetailPanelViewportOffset = max(0, int(m.DetailPanelViewport.HorizontalScrollPercent()*float64(m.DetailPanelViewportOffset))-1)
	m.DetailPanelTwoViewportOffset = max(0, int(m.DetailPanelTwoViewport.HorizontalScrollPercent()*float64(m.DetailPanelTwoViewportOffset))-1)
	m.DetailPanelViewport.SetXOffset(m.DetailPanelViewportOffset)
	m.DetailPanelViewport.SetYOffset(m.DetailPanelViewport.YOffset())
	m.DetailPanelTwoViewport.SetXOffset(m.DetailPanelTwoViewportOffset)
	m.DetailPanelTwoViewport.SetYOffset(m.DetailPanelTwoViewport.YOffset())

	if m.IsLineEditingState.Load() {
		services.EnterOrReinitLineEditingStateService(m)
	}
}

func LeftPanelDynamicResize(m *types.GittiModel) {
	// this is after reserving the height for the gitti status panel and also Padding
	leftPanelRemainingHeight := m.WindowCoreContentHeight - 1 - ((len(constant.ComponentNavigationList) - 1) * 2)

	// we minus 2 if GitStatusComponent is not the one chosen is because GitStatusComponent
	// and the one that got selected will not be account in to the dynamic height calculation
	// ( gitti status component's height is fix at 3, while the selected one will always get 40% )
	componentWithDynamicHeight := (len(constant.ComponentNavigationList) - 2)
	unSelectedComponentPanelHeightPerComponent := int(int(float64(leftPanelRemainingHeight)*(1.0-constant.SelectedLeftPanelComponentHeightRatio)) / componentWithDynamicHeight)
	selectedComponentPanelHeight := leftPanelRemainingHeight - (unSelectedComponentPanelHeightPerComponent * componentWithDynamicHeight)
	m.LocalBranchesComponentPanelHeight = unSelectedComponentPanelHeightPerComponent
	m.ModifiedFilesComponentPanelHeight = unSelectedComponentPanelHeightPerComponent
	m.CommitLogComponentPanelHeight = unSelectedComponentPanelHeightPerComponent
	m.StashComponentPanelHeight = unSelectedComponentPanelHeightPerComponent

	switch m.CurrentSelectedComponent {
	case constant.LocalBranchComponent:
		m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
	case constant.ModifiedFilesComponent:
		m.ModifiedFilesComponentPanelHeight = selectedComponentPanelHeight
	case constant.CommitLogComponent:
		m.CommitLogComponentPanelHeight = selectedComponentPanelHeight
	case constant.StashComponent:
		m.StashComponentPanelHeight = selectedComponentPanelHeight
	case constant.GitStatusComponent:
		// if it was the Gitti status component panel that got selected (because its height is fix),
		// the next panel will get the selected height which is the branch component panel
		m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
	case constant.DetailComponentTwo:
		switch m.DetailPanelParentComponent {
		case constant.LocalBranchComponent:
			m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
		case constant.ModifiedFilesComponent:
			m.ModifiedFilesComponentPanelHeight = selectedComponentPanelHeight
		case constant.CommitLogComponent:
			m.CommitLogComponentPanelHeight = selectedComponentPanelHeight
		case constant.StashComponent:
			m.StashComponentPanelHeight = selectedComponentPanelHeight
		}
	case constant.DetailComponent:
		switch m.DetailPanelParentComponent {
		case constant.LocalBranchComponent:
			m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
		case constant.ModifiedFilesComponent:
			m.ModifiedFilesComponentPanelHeight = selectedComponentPanelHeight
		case constant.CommitLogComponent:
			m.CommitLogComponentPanelHeight = selectedComponentPanelHeight
		case constant.StashComponent:
			m.StashComponentPanelHeight = selectedComponentPanelHeight
		}
	}

	// update all components Width and Height
	m.CurrentRepoBranchesInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoBranchesInfoList.SetHeight(m.LocalBranchesComponentPanelHeight)

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
