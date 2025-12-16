package layout

import (
	"github.com/gohyuhan/gitti/tui/constant"
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

func UpdateDetailComponentViewportLayout(m *types.GittiModel) {
	if m.ShowDetailPanelTwo.Load() {
		// vertical layout
		// Since terminal characters are usually about twice as tall as they are wide,
		// we weight the height by 2 to approximate visual "squareness".
		splitHeight := int(m.DetailComponentPanelHeight / 2)
		splitWidth := int(m.DetailComponentPanelWidth / 2)

		if m.DetailComponentPanelHeight*2 > m.DetailComponentPanelWidth {
			m.DetailComponentPanelLayout = constant.VERTICAL
			m.DetailPanelViewport.SetHeight(splitHeight - 1)
			m.DetailPanelViewport.SetWidth(m.DetailComponentPanelWidth - 2)
			m.DetailPanelTwoViewport.SetHeight(m.DetailComponentPanelHeight - splitHeight - 1)
			m.DetailPanelTwoViewport.SetWidth(m.DetailComponentPanelWidth - 2)
		} else {
			// horizontal layout
			m.DetailComponentPanelLayout = constant.HORIZONTAL
			m.DetailPanelViewport.SetHeight(m.DetailComponentPanelHeight)
			m.DetailPanelViewport.SetWidth(splitWidth - 2)
			m.DetailPanelTwoViewport.SetHeight(m.DetailComponentPanelHeight)
			m.DetailPanelTwoViewport.SetWidth(m.DetailComponentPanelWidth - splitWidth - 2)
		}
	} else {
		m.DetailPanelViewport.SetHeight(m.DetailComponentPanelHeight)
		m.DetailPanelViewport.SetWidth(m.DetailComponentPanelWidth - 2)
	}
}
