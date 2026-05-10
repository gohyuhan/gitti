package layout

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/gitti/i18n"
	branchComponent "github.com/gohyuhan/gitti/tui/component/branch"
	commitlogComponent "github.com/gohyuhan/gitti/tui/component/commitlog"
	"github.com/gohyuhan/gitti/tui/component/files"
	filesComponent "github.com/gohyuhan/gitti/tui/component/files"
	reflogComponent "github.com/gohyuhan/gitti/tui/component/reflog"
	remoteComponent "github.com/gohyuhan/gitti/tui/component/remote"
	stashComponent "github.com/gohyuhan/gitti/tui/component/stash"
	tagComponent "github.com/gohyuhan/gitti/tui/component/tag"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	To update the width and height of all components
//
// ------------------------------------
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
	m.CurrentRepoRefLogInfoList.Title = ansi.Truncate(reflogComponent.ConstructRefLogComponentTitle(titleWidthLimit), titleWidthLimit, "...")
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
		EnterOrReinitLineEditingState(m)
	}

	// log panel (the height is fixed)
	m.CurrentLogComponentViewport.SetWidth(m.DetailComponentPanelWidth - 2)
	m.CurrentLogComponentViewport.SetHeight(m.LogComponentPanelHeight)
}

// ------------------------------------
//
//	Dynamically resize the left panel components to fill the available height
//
// ------------------------------------
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
	m.RefLogComponentPanelHeight = unSelectedComponentPanelHeightPerComponent
	m.StashComponentPanelHeight = unSelectedComponentPanelHeightPerComponent

	switch m.CurrentSelectedComponent {
	case constant.LocalBranchOrTagOrRemoteComponentPanel:
		m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
		m.TagsComponentPanelHeight = selectedComponentPanelHeight
		m.RemoteComponentPanelHeight = selectedComponentPanelHeight
	case constant.ModifiedFilesComponentPanel:
		m.ModifiedFilesComponentPanelHeight = selectedComponentPanelHeight
	case constant.CommitLogOrRefLogComponentPanel:
		m.CommitLogComponentPanelHeight = selectedComponentPanelHeight
		m.RefLogComponentPanelHeight = selectedComponentPanelHeight
	case constant.StashComponentPanel:
		m.StashComponentPanelHeight = selectedComponentPanelHeight
	case constant.GitStatusComponentPanel:
		// if it was the Gitti status component panel that got selected (because its height is fix),
		// the next panel will get the selected height which is the branch component panel
		m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
		m.TagsComponentPanelHeight = selectedComponentPanelHeight
		m.RemoteComponentPanelHeight = selectedComponentPanelHeight
	case constant.DetailComponentPanelTwo:
		switch m.DetailPanelParentComponent {
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
			m.TagsComponentPanelHeight = selectedComponentPanelHeight
			m.RemoteComponentPanelHeight = selectedComponentPanelHeight
		case constant.ModifiedFilesComponentPanel:
			m.ModifiedFilesComponentPanelHeight = selectedComponentPanelHeight
		case constant.CommitLogOrRefLogComponentPanel:
			m.CommitLogComponentPanelHeight = selectedComponentPanelHeight
			m.RefLogComponentPanelHeight = selectedComponentPanelHeight
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
		case constant.CommitLogOrRefLogComponentPanel:
			m.CommitLogComponentPanelHeight = selectedComponentPanelHeight
			m.RefLogComponentPanelHeight = selectedComponentPanelHeight
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

	m.CurrentRepoRefLogInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoRefLogInfoList.SetHeight(m.RefLogComponentPanelHeight)

	m.CurrentRepoStashInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoStashInfoList.SetHeight(m.StashComponentPanelHeight)
}

func UpdateDetailComponentViewportContentAndState(m *types.GittiModel, UpdateData types.DetailPanelStateAndLayoutUpdateInterface) {
	needToScrollToBottom := m.CurrentSelectedComponent == constant.LogComponentPanel
	m.DetailPanelViewport.SetContent(UpdateData.ContentLine)
	m.DetailPanelViewportOGStringArray = UpdateData.OgLineDiff1
	m.DetailPanelTwoViewportOGStringArray = UpdateData.OgLineDiff2

	if needToScrollToBottom {
		m.DetailPanelViewport.GotoBottom()
	}

	if UpdateData.SetForDetailComponentTwo {
		m.DetailPanelTwoViewport.SetContent(UpdateData.ContentLine2)
		m.ShowDetailPanelTwo.Store(true)
	} else {
		// if the detail component two is selected, switch to detail component as it is not set for detail component two
		// as it will hide the detail component two viewport
		if m.CurrentSelectedComponent == constant.DetailComponentPanelTwo {
			m.CurrentSelectedComponent = constant.DetailComponentPanel
		}
	}
	if m.IsLineEditingState.Load() {
		EnterOrReinitLineEditingState(m)
	}
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

// ------------------------------------
//
//			For Enter Line Editing State
//		  * this was to enter or reinit the line editing state
//		  * it will calculate the index position for both visible and actual content
//	   * it will also determine if the current file selected has staged/unstaged changes or both
//	   * lastly it will set the cursor viewport content
//
// ------------------------------------
func EnterOrReinitLineEditingState(m *types.GittiModel) {
	if !((m.CurrentSelectedComponent == constant.DetailComponentPanel || m.CurrentSelectedComponent == constant.DetailComponentPanelTwo) && m.DetailPanelParentComponent == constant.ModifiedFilesComponentPanel && m.CurrentRepoModifiedFilesInfoList.SelectedItem() != nil) {
		// we are eligible to be in line editing mode, so we need to reset the state
		m.IsLineEditingState.Store(false)
		// reinit the index position
		m.LineEditingIndexPositionAndInfo = types.GittiLineEditingIndexPositionAndInfo{}
		return
	}

	detailPanelViewportVisibleIndex := m.DetailPanelViewport.VisibleLineCount()
	detailPanelTwoViewportVisibleIndex := m.DetailPanelTwoViewport.VisibleLineCount()
	detailPanelViewportTotalIndex := m.DetailPanelViewport.TotalLineCount()
	detailPanelTwoViewportTotalIndex := m.DetailPanelTwoViewport.TotalLineCount()

	if !m.IsLineEditingState.Load() {
		// we first confirm that we are not in line editing mode yet, but we are supposed to be
		// so we need to init the line editing state
		currentSelectedFileItem := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
		currentSelectedFile := currentSelectedFileItem.(files.GitModifiedFilesItem)

		var detailPanelViewportStageType string
		var detailPanelTwoViewportStageType string
		detailPanelViewportOverflowIndexCount := 0
		detailPanelTwoViewportOverflowIndexCount := 0

		if currentSelectedFile.HasConflict {
			// if the file is in conflict state, we cannot do line editing
			m.IsLineEditingState.Store(false)
			detailPanelViewportStageType = ""
			detailPanelTwoViewportStageType = ""
		} else {
			m.IsLineEditingState.Store(true)
			// determine the pop up state
			if currentSelectedFile.IndexState == "?" && currentSelectedFile.WorkTree == "?" {
				// newly added untracked file
				detailPanelViewportStageType = constant.UNSTAGE
				detailPanelTwoViewportStageType = constant.NOSTAGESTATUS
				detailPanelViewportOverflowIndexCount = 2
				detailPanelTwoViewportOverflowIndexCount = 2
			} else if currentSelectedFile.IndexState != " " && currentSelectedFile.WorkTree != " " {
				// tracked file with both staged and unstaged modification
				detailPanelViewportStageType = constant.STAGE
				detailPanelTwoViewportStageType = constant.UNSTAGE
				detailPanelViewportOverflowIndexCount = 3
				detailPanelTwoViewportOverflowIndexCount = 3
			} else if currentSelectedFile.IndexState != " " && currentSelectedFile.WorkTree == " " {
				// tracked file with only staged modification
				detailPanelViewportStageType = constant.STAGE
				detailPanelTwoViewportStageType = constant.NOSTAGESTATUS
				detailPanelViewportOverflowIndexCount = 2
				detailPanelTwoViewportOverflowIndexCount = 2
			} else {
				// tracked file with only unstaged modification
				detailPanelViewportStageType = constant.UNSTAGE
				detailPanelTwoViewportStageType = constant.NOSTAGESTATUS
				detailPanelViewportOverflowIndexCount = 2
				detailPanelTwoViewportOverflowIndexCount = 2
			}
		}

		// reinit the index position
		m.LineEditingIndexPositionAndInfo = types.GittiLineEditingIndexPositionAndInfo{
			DetailPanelViewportIndexPosition:         0,
			DetailPanelTwoViewportIndexPosition:      0,
			DetailPanelViewportStageType:             detailPanelViewportStageType,
			DetailPanelTwoViewportStageType:          detailPanelTwoViewportStageType,
			DetailPanelViewportActualCurrentIndex:    0,
			DetailPanelTwoViewportActualCurrentIndex: 0,
			DetailPanelViewportOverflowIndexCount:    detailPanelViewportOverflowIndexCount,
			DetailPanelTwoViewportOverflowIndexCount: detailPanelTwoViewportOverflowIndexCount,
		}

		// set the cursor viewport
		SetLineEditingCursorViewportContent(m, detailPanelViewportVisibleIndex, detailPanelTwoViewportVisibleIndex)
	} else if m.IsLineEditingState.Load() {
		// we are already in line editing mode, so we need to update the state
		// this usually happens when we resize the window or similar event
		currentSelectedFileItem := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
		currentSelectedFile := currentSelectedFileItem.(files.GitModifiedFilesItem)

		var detailPanelViewportStageType string
		var detailPanelTwoViewportStageType string

		var detailPanelViewportIndexPosition int
		var detailPanelTwoViewportIndexPosition int
		var detailPanelViewportActualCurrentIndex int
		var detailPanelTwoViewportActualCurrentIndex int
		var detailPanelViewportOverflowIndexCount int
		var detailPanelTwoViewportOverflowIndexCount int

		if currentSelectedFile.HasConflict {
			// if the file is in conflict state, we cannot do line editing
			m.IsLineEditingState.Store(false)
			detailPanelViewportStageType = ""
			detailPanelTwoViewportStageType = ""
		} else {
			m.IsLineEditingState.Store(true)
			// determine the pop up state
			if currentSelectedFile.IndexState == "?" && currentSelectedFile.WorkTree == "?" {
				// newly added untracked file
				detailPanelViewportStageType = constant.UNSTAGE
				detailPanelTwoViewportStageType = constant.NOSTAGESTATUS
				detailPanelViewportOverflowIndexCount = 2
				detailPanelTwoViewportOverflowIndexCount = 2
			} else if currentSelectedFile.IndexState != " " && currentSelectedFile.WorkTree != " " {
				// tracked file with both staged and unstaged modification
				detailPanelViewportStageType = constant.STAGE
				detailPanelTwoViewportStageType = constant.UNSTAGE
				detailPanelViewportOverflowIndexCount = 3
				detailPanelTwoViewportOverflowIndexCount = 3
			} else if currentSelectedFile.IndexState != " " && currentSelectedFile.WorkTree == " " {
				// tracked file with only staged modification
				detailPanelViewportStageType = constant.STAGE
				detailPanelTwoViewportStageType = constant.NOSTAGESTATUS
				detailPanelViewportOverflowIndexCount = 2
				detailPanelTwoViewportOverflowIndexCount = 2
			} else {
				// tracked file with only unstaged modification
				detailPanelViewportStageType = constant.UNSTAGE
				detailPanelTwoViewportStageType = constant.NOSTAGESTATUS
				detailPanelViewportOverflowIndexCount = 2
				detailPanelTwoViewportOverflowIndexCount = 2
			}
		}

		// recalculate the actual position
		detailPanelViewportActualCurrentIndex = m.LineEditingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex
		detailPanelTwoViewportActualCurrentIndex = m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportActualCurrentIndex
		if detailPanelViewportOverflowIndexCount < m.LineEditingIndexPositionAndInfo.DetailPanelViewportOverflowIndexCount &&
			detailPanelTwoViewportOverflowIndexCount < m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportIndexPosition {
			detailPanelViewportActualCurrentIndex -= 1
			detailPanelTwoViewportActualCurrentIndex -= 1
		} else if detailPanelViewportOverflowIndexCount > m.LineEditingIndexPositionAndInfo.DetailPanelViewportOverflowIndexCount &&
			detailPanelTwoViewportOverflowIndexCount > m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportIndexPosition {
			detailPanelViewportActualCurrentIndex += 1
			detailPanelTwoViewportActualCurrentIndex += 1
		}

		// make sure the actual current index is not out of bound
		detailPanelViewportActualCurrentIndex = min(detailPanelViewportActualCurrentIndex, detailPanelViewportTotalIndex-1)
		detailPanelTwoViewportActualCurrentIndex = min(detailPanelTwoViewportActualCurrentIndex, detailPanelTwoViewportTotalIndex-1)

		// recalculate the visible index position
		// The relationship is: Actual Index = Viewport Offset + Visible Index
		// So: Visible Index = Actual Index - Viewport Offset

		// For Panel 1
		currentOffset := m.DetailPanelViewport.YOffset()
		calculatedVisibleIndex := detailPanelViewportActualCurrentIndex - currentOffset

		if calculatedVisibleIndex < 0 {
			// Cursor is above the current view, scroll up to make it the first line
			m.DetailPanelViewport.SetYOffset(detailPanelViewportActualCurrentIndex)
			detailPanelViewportIndexPosition = 0
		} else if calculatedVisibleIndex >= detailPanelViewportVisibleIndex {
			// Cursor is below the current view, scroll down to make it the last visible line
			newOffset := detailPanelViewportActualCurrentIndex - (detailPanelViewportVisibleIndex - 1)
			m.DetailPanelViewport.SetYOffset(newOffset)
			detailPanelViewportIndexPosition = detailPanelViewportVisibleIndex - 1
		} else {
			// Cursor is within view
			detailPanelViewportIndexPosition = calculatedVisibleIndex
		}

		// For Panel 2
		currentOffsetTwo := m.DetailPanelTwoViewport.YOffset()
		calculatedVisibleIndexTwo := detailPanelTwoViewportActualCurrentIndex - currentOffsetTwo

		if calculatedVisibleIndexTwo < 0 {
			// Cursor is above, scroll up
			m.DetailPanelTwoViewport.SetYOffset(detailPanelTwoViewportActualCurrentIndex)
			detailPanelTwoViewportIndexPosition = 0
		} else if calculatedVisibleIndexTwo >= detailPanelTwoViewportVisibleIndex {
			// Cursor is below, scroll down
			newOffsetTwo := detailPanelTwoViewportActualCurrentIndex - (detailPanelTwoViewportVisibleIndex - 1)
			m.DetailPanelTwoViewport.SetYOffset(newOffsetTwo)
			detailPanelTwoViewportIndexPosition = detailPanelTwoViewportVisibleIndex - 1
		} else {
			// Cursor is within view
			detailPanelTwoViewportIndexPosition = calculatedVisibleIndexTwo
		}

		// reinit the index position
		m.LineEditingIndexPositionAndInfo = types.GittiLineEditingIndexPositionAndInfo{
			DetailPanelViewportIndexPosition:         detailPanelViewportIndexPosition,
			DetailPanelTwoViewportIndexPosition:      detailPanelTwoViewportIndexPosition,
			DetailPanelViewportStageType:             detailPanelViewportStageType,
			DetailPanelTwoViewportStageType:          detailPanelTwoViewportStageType,
			DetailPanelViewportActualCurrentIndex:    detailPanelViewportActualCurrentIndex,
			DetailPanelTwoViewportActualCurrentIndex: detailPanelTwoViewportActualCurrentIndex,
			DetailPanelViewportOverflowIndexCount:    detailPanelViewportOverflowIndexCount,
			DetailPanelTwoViewportOverflowIndexCount: detailPanelTwoViewportOverflowIndexCount,
		}

		// set the cursor viewport
		SetLineEditingCursorViewportContent(m, detailPanelViewportVisibleIndex, detailPanelTwoViewportVisibleIndex)
	}
}

// ------------------------------------
//
//			For Set Line Editing Cursor Viewport Content
//		  * this was to set the cursor viewport content
//	   * needed because we are using a dual viewport setup for line editing mode (one for cursor, one for content)
//
// ------------------------------------
func SetLineEditingCursorViewportContent(m *types.GittiModel, detailPanelViewportVisibleIndex int, detailPanelTwoViewportVisibleIndex int) {
	// set the cursor viewport
	var cursorVpLine strings.Builder
	var cursorVpTwoLine strings.Builder

	for index := range detailPanelViewportVisibleIndex {
		if index == m.LineEditingIndexPositionAndInfo.DetailPanelViewportIndexPosition {
			cursorVpLine.WriteString(style.SelectedItemStyle.Render("❯  "))
			cursorVpLine.WriteRune('\n')
		} else {
			cursorVpLine.WriteString(style.NewStyle.Render("   "))
			cursorVpLine.WriteRune('\n')
		}
	}

	for index := range detailPanelTwoViewportVisibleIndex {
		if index == m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportIndexPosition {
			cursorVpTwoLine.WriteString(style.SelectedItemStyle.Render("❯  "))
			cursorVpTwoLine.WriteRune('\n')
		} else {
			cursorVpTwoLine.WriteString(style.NewStyle.Render("   "))
			cursorVpTwoLine.WriteRune('\n')
		}
	}
	m.LineEditingIndexCursorViewport.SetContent(cursorVpLine.String())
	m.LineEditingIndexCursorTwoViewport.SetContent(cursorVpTwoLine.String())
}

func DetailComponentReinit(m *types.GittiModel) {
	m.DetailPanelViewport.SetContent(style.NewStyle.Render(i18n.LANGUAGEMAPPING.Loading))
	m.ShowDetailPanelTwo.Store(false)
	m.DetailPanelViewportOffset = 0
	m.DetailPanelViewport.SetXOffset(0)
	m.DetailPanelViewport.SetYOffset(0)
	m.DetailPanelTwoViewportOffset = 0
	m.DetailPanelTwoViewport.SetXOffset(0)
	m.DetailPanelTwoViewport.SetYOffset(0)
}
