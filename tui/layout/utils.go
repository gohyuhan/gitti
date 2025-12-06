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
	m.DetailPanelViewport.SetHeight(m.DetailComponentPanelHeight) // some margin
	m.DetailPanelViewport.SetWidth(m.DetailComponentPanelWidth - 2)
	m.DetailPanelViewportOffset = max(0, int(m.DetailPanelViewport.HorizontalScrollPercent()*float64(m.DetailPanelViewportOffset))-1)
	m.DetailPanelViewport.SetXOffset(m.DetailPanelViewportOffset)
	m.DetailPanelViewport.SetYOffset(m.DetailPanelViewport.YOffset())
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
