package layout

import (
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
)

// to update the width and height of all components
func TuiWindowSizing(m *types.GittiModel) {
	// Compute panel widths
	m.WindowLeftPanelWidth = min(int(float64(m.Width)*settings.GITTICONFIGSETTINGS.LeftPanelWidthRatio), constant.MaxLeftPanelWidth)
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
	leftPanelRemainingHeight := m.WindowCoreContentHeight - 7 // this is after reserving the height for the gitti version panel and also Padding

	// we minus 2 if GitStatusComponent is not the one chosen is because GitStatusComponent
	// and the one that got selected will not be account in to the dynamic height calculation
	// ( gitti status component's height is fix at 3, while the selected one will always get 40% )
	componentWithDynamicHeight := (len(constant.ComponentNavigationList) - 2)
	unSelectedComponentPanelHeightPerComponent := (int(float64(leftPanelRemainingHeight)*(1.0-constant.SelectedLeftPanelComponentHeightRatio)) / componentWithDynamicHeight)
	selectedComponentPanelHeight := leftPanelRemainingHeight - (unSelectedComponentPanelHeightPerComponent * componentWithDynamicHeight)
	m.LocalBranchesComponentPanelHeight = unSelectedComponentPanelHeightPerComponent
	m.ModifiedFilesComponentPanelHeight = unSelectedComponentPanelHeightPerComponent
	m.StashComponentPanelHeight = unSelectedComponentPanelHeightPerComponent

	switch m.CurrentSelectedComponent {
	case constant.LocalBranchComponent:
		m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
	case constant.ModifiedFilesComponent:
		m.ModifiedFilesComponentPanelHeight = selectedComponentPanelHeight
	case constant.StashComponent:
		m.StashComponentPanelHeight = selectedComponentPanelHeight
	case constant.GitStatusComponent:
		// if it was the Gitti status component panel that got selected (because its height is fix),
		// the next panel will get the selected height which is the branch component panel
		m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
	}
	// update all components Width and Height
	m.CurrentRepoBranchesInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoBranchesInfoList.SetHeight(m.LocalBranchesComponentPanelHeight)

	m.CurrentRepoModifiedFilesInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoModifiedFilesInfoList.SetHeight(m.ModifiedFilesComponentPanelHeight)

	m.CurrentRepoStashInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoStashInfoList.SetHeight(m.StashComponentPanelHeight)
}

func UpdateDetailComponentViewportLayout(m *types.GittiModel) {
	if m.ShowDetailPanelTwo.Load() {
		// vertical layout
		// Since terminal characters are usually about twice as tall as they are wide,
		// we weight the height by 2 to approximate visual "squareness".
		if m.DetailComponentPanelHeight*2 > m.DetailComponentPanelWidth {
			m.DetailComponentPanelLayout = constant.VERTICAL
			m.DetailPanelViewport.SetHeight(m.DetailComponentPanelHeight/2 - 1)
			m.DetailPanelViewport.SetWidth(m.DetailComponentPanelWidth - 2)
			m.DetailPanelTwoViewport.SetHeight(m.DetailComponentPanelHeight/2 - 1)
			m.DetailPanelTwoViewport.SetWidth(m.DetailComponentPanelWidth - 2)
		} else {
			// horizontal layout
			m.DetailComponentPanelLayout = constant.HORIZONTAL
			m.DetailPanelViewport.SetHeight(m.DetailComponentPanelHeight)
			m.DetailPanelViewport.SetWidth(m.DetailComponentPanelWidth/2 - 2)
			m.DetailPanelTwoViewport.SetHeight(m.DetailComponentPanelHeight)
			m.DetailPanelTwoViewport.SetWidth(m.DetailComponentPanelWidth/2 - 2)
		}
	} else {
		m.DetailPanelViewport.SetHeight(m.DetailComponentPanelHeight)
		m.DetailPanelViewport.SetWidth(m.DetailComponentPanelWidth - 2)
	}
}
