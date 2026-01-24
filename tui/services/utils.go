package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/component/commitlog"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/component/stash"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// services was to bridge api and the needs of the terminal interface logic so that it can be compatible and feels smooth and not clunky
// ------------------------------------
//
//		For fetching detail component panel info
//	  * it can be for stash info, commit info etc
//
// ------------------------------------
func FetchDetailComponentPanelInfoService(m *types.GittiModel, reinit bool) {
	// For non-reinit calls (refreshing current view), abort if already processing.
	// This avoids looping a cancel and execution cycle which would end up blocking
	// a slightly longer processing process.
	//
	// If not processing, we proceed to fetch to ensure we capture any updates (e.g., file changes,
	// amends), as we lack specific context on whether the underlying data has changed.
	//
	// If `reinit` is true (context switch), we bypass this check to cancel the active fetch
	// and start the new one immediately.
	if !reinit && m.IsDetailComponentPanelInfoFetchProcessing.Load() {
		return
	}

	// Cancel any existing operation first
	if m.DetailComponentPanelInfoFetchCancelFunc != nil {
		m.DetailComponentPanelInfoFetchCancelFunc()
	}

	// Wait for the previous goroutine to finish (its defer will set processing to false),
	// then atomically set it to true before starting a new one.
	for !m.IsDetailComponentPanelInfoFetchProcessing.CompareAndSwap(false, true) {
		// The previous goroutine is still running, wait a tiny bit
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.DetailComponentPanelInfoFetchCancelFunc = cancel
	go func(ctx context.Context) {
		defer func() {
			cancel()
			m.IsDetailComponentPanelInfoFetchProcessing.Store(false)
		}()

		var contentLine string
		var contentLine2 string // fro detail panel 2nd (only used for files changes to show staged and unstaged diff in seperated panel)
		setForDetailComponentTwo := false
		var theCurrentSelectedComponent string
		// reinit and render detail component panel viewport
		if reinit {
			m.DetailPanelViewport.SetContent(style.NewStyle.Render(i18n.LANGUAGEMAPPING.Loading))
			m.ShowDetailPanelTwo.Store(false)
			m.DetailPanelViewportOffset = 0
			m.DetailPanelViewport.SetXOffset(0)
			m.DetailPanelViewport.SetYOffset(0)
			m.DetailPanelTwoViewportOffset = 0
			m.DetailPanelTwoViewport.SetXOffset(0)
			m.DetailPanelTwoViewport.SetYOffset(0)
		}
		if m.CurrentSelectedComponent == constant.DetailComponent || m.CurrentSelectedComponent == constant.DetailComponentTwo {
			// if the current selected one is the detail component itself, the current selected one will be its parent (the component that led into the detail component)
			theCurrentSelectedComponent = m.DetailPanelParentComponent
		} else {
			theCurrentSelectedComponent = m.CurrentSelectedComponent
		}
		switch theCurrentSelectedComponent {
		case constant.ModifiedFilesComponent:
			contentLine, contentLine2, setForDetailComponentTwo = generateBothModifiedFileDetailPanelContent(ctx, m)
		case constant.CommitLogComponent:
			contentLine = generateCommitLogDetailPanelContent(ctx, m)
		case constant.StashComponent:
			contentLine = generateStashDetailPanelContent(ctx, m)
		default:
			contentLine = generateAboutGittiContent()
		}

		select {
		case <-ctx.Done():
			return
		default:
			if contentLine == "" {
				// if the content will be empty, render about gitti for detail panel
				contentLine = generateAboutGittiContent()
			}
			m.DetailPanelViewport.SetContent(contentLine)

			if setForDetailComponentTwo {
				m.DetailPanelTwoViewport.SetContent(contentLine2)
				m.ShowDetailPanelTwo.Store(true)
			} else {
				// if the detail component two is selected, switch to detail component as it is not set for detail component two
				// as it will hide the detail component two viewport
				if m.CurrentSelectedComponent == constant.DetailComponentTwo {
					m.CurrentSelectedComponent = constant.DetailComponent
				}
			}
			if m.IsLineStagingState.Load() {
				EnterOrReinitLineStagingStateService(m)
			}
			m.TuiUpdateChannel <- constant.DETAIL_COMPONENT_PANEL_UPDATED
			return
		}
	}(ctx)
}

// for modified file detail panel view
func generateBothModifiedFileDetailPanelContent(ctx context.Context, m *types.GittiModel) (string, string, bool) {
	shouldRenderDetailComponentPanelTwo := false
	currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
	var fileStatus git.FileStatus
	if currentSelectedModifiedFile != nil {
		fileStatus = git.FileStatus(currentSelectedModifiedFile.(files.GitModifiedFilesItem))
	} else {
		return "", "", shouldRenderDetailComponentPanelTwo
	}

	var vpLine1 strings.Builder
	vpLine1.WriteString(fmt.Sprintf("[ %s ]\n\n", fileStatus.FilePathname))

	var vpLine2 strings.Builder
	var fileDiffLines1 []string
	var fileDiffLines2 []string
	getDiffTypeForVpLine1 := git.GETCOMBINEDDIFF

	// indicating that the file is not in conflict state and have both staged and unstaged changes and they are not in "?" for both index and worktree
	if !fileStatus.HasConflict && fileStatus.IndexState != " " && fileStatus.WorkTree != " " && fileStatus.IndexState != "?" && fileStatus.WorkTree != "?" {
		shouldRenderDetailComponentPanelTwo = true
		vpLine2.WriteString(fmt.Sprintf("%s\n[ %s ]\n\n", i18n.LANGUAGEMAPPING.UnstagedTitle, fileStatus.FilePathname))
		fileDiffLines2 = m.GitOperations.GitFiles.GetFilesDiffInfo(ctx, fileStatus, git.GETUNSTAGEDDIFF)
		m.DetailPanelTwoViewportOGStringArray = fileDiffLines2

		if fileDiffLines2 == nil {
			vpLine2.WriteString(i18n.LANGUAGEMAPPING.FileTypeUnSupportedPreview)
		} else {
			for _, line := range fileDiffLines2 {
				line = style.NewStyle.Render(line)
				vpLine2.WriteString(line + "\n")
			}
		}

		getDiffTypeForVpLine1 = git.GETSTAGEDDIFF
		vpLine1.Reset()
		vpLine1.WriteString(fmt.Sprintf("%s\n[ %s ]\n\n", i18n.LANGUAGEMAPPING.StagedTitle, fileStatus.FilePathname))
	}

	fileDiffLines1 = m.GitOperations.GitFiles.GetFilesDiffInfo(ctx, fileStatus, getDiffTypeForVpLine1)
	m.DetailPanelViewportOGStringArray = fileDiffLines1
	if fileDiffLines1 == nil {
		vpLine1.WriteString(i18n.LANGUAGEMAPPING.FileTypeUnSupportedPreview)
	} else {
		for _, line := range fileDiffLines1 {
			line = style.NewStyle.Render(line)
			vpLine1.WriteString(line + "\n")
		}
	}

	return vpLine1.String(), vpLine2.String(), shouldRenderDetailComponentPanelTwo
}

// for commit log detail panel view
func generateCommitLogDetailPanelContent(ctx context.Context, m *types.GittiModel) string {
	currentSelectedCommitLog := m.CurrentRepoCommitLogInfoList.SelectedItem()
	var commitLogItem commitlog.GitCommitLogItem
	var vpLine strings.Builder
	if currentSelectedCommitLog != nil {
		commitLogItem = currentSelectedCommitLog.(commitlog.GitCommitLogItem)
	} else {
		return ""
	}

	commitLogDetail := m.GitOperations.GitCommitLog.GitCommitLogDetail(ctx, commitLogItem.Hash)
	if len(commitLogDetail) < 1 {
		return ""
	}

	for _, Line := range commitLogDetail {
		line := style.NewStyle.Render(Line)
		vpLine.WriteString(line + "\n")
	}
	return vpLine.String()
}

// for stash detail panel view
func generateStashDetailPanelContent(ctx context.Context, m *types.GittiModel) string {
	currentSelectedStash := m.CurrentRepoStashInfoList.SelectedItem()
	var stashItem stash.GitStashItem
	if currentSelectedStash != nil {
		stashItem = currentSelectedStash.(stash.GitStashItem)
	} else {
		return ""
	}

	var vpLine strings.Builder
	vpLine.WriteString(fmt.Sprintf(
		"[%s]\n[%s]\n\n",
		style.StashIdStyle.Render(stashItem.Id),
		style.StashMessageStyle.Render(stashItem.Message),
	))

	stashDetail := m.GitOperations.GitStash.GitStashDetail(ctx, stashItem.Id)
	if len(stashDetail) < 1 {
		return ""
	}

	for _, Line := range stashDetail {
		line := style.NewStyle.Render(Line)
		vpLine.WriteString(line + "\n")
	}
	return vpLine.String()
}

// for about gitti content
func generateAboutGittiContent() string {
	var vpLine strings.Builder

	logoLineArray := style.GradientLines(constant.GittiAsciiArtLogo)
	aboutLines := i18n.LANGUAGEMAPPING.AboutGitti

	vpLine.WriteString(strings.Join(logoLineArray, "\n") + "\n")
	vpLine.WriteString(strings.Join(aboutLines, "\n"))

	return vpLine.String()
}

// ------------------------------------
//
//			For Enter Line Staging State Service
//		  * this was to enter or reinit the line staging state
//		  * it will calculate the index position for both visible and actual content
//	   * it will also determine if the current file selected has staged/unstaged changes or both
//	   * lastly it will set the cursor viewport content
//
// ------------------------------------
func EnterOrReinitLineStagingStateService(m *types.GittiModel) {
	if !((m.CurrentSelectedComponent == constant.DetailComponent || m.CurrentSelectedComponent == constant.DetailComponentTwo) && m.DetailPanelParentComponent == constant.ModifiedFilesComponent && m.CurrentRepoModifiedFilesInfoList.SelectedItem() != nil) {
		// we are eligible to be in line staging mode, so we need to reset the state
		m.IsLineStagingState.Store(false)
		// reinit the index position
		m.LineStagingIndexPositionAndInfo = types.GittiLineStagingIndexPositionAndInfo{}
		return
	}

	detailPanelViewportVisibleIndex := m.DetailPanelViewport.VisibleLineCount()
	detailPanelTwoViewportVisibleIndex := m.DetailPanelTwoViewport.VisibleLineCount()
	detailPanelViewportTotalIndex := m.DetailPanelViewport.TotalLineCount()
	detailPanelTwoViewportTotalIndex := m.DetailPanelTwoViewport.TotalLineCount()

	if !m.IsLineStagingState.Load() {
		// we first confirm that we are not in line staging mode yet, but we are supposed to be
		// so we need to init the line staging state
		currentSelectedFileItem := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
		currentSelectedFile := currentSelectedFileItem.(files.GitModifiedFilesItem)

		var detailPanelViewportStageType string
		var detailPanelTwoViewportStageType string
		detailPanelViewportOverflowIndexCount := 0
		detailPanelTwoViewportOverflowIndexCount := 0

		if currentSelectedFile.HasConflict {
			// if the file is in conflict state, we cannot do line staging
			m.IsLineStagingState.Store(false)
			detailPanelViewportStageType = ""
			detailPanelTwoViewportStageType = ""
		} else {
			m.IsLineStagingState.Store(true)
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
		m.LineStagingIndexPositionAndInfo = types.GittiLineStagingIndexPositionAndInfo{
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
		SetLineStagingCursorViewportContent(m, detailPanelViewportVisibleIndex, detailPanelTwoViewportVisibleIndex)
	} else if m.IsLineStagingState.Load() {
		// we are already in line staging mode, so we need to update the state
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
			// if the file is in conflict state, we cannot do line staging
			m.IsLineStagingState.Store(false)
			detailPanelViewportStageType = ""
			detailPanelTwoViewportStageType = ""
		} else {
			m.IsLineStagingState.Store(true)
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
		detailPanelViewportActualCurrentIndex = m.LineStagingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex
		detailPanelTwoViewportActualCurrentIndex = m.LineStagingIndexPositionAndInfo.DetailPanelTwoViewportActualCurrentIndex
		if detailPanelViewportOverflowIndexCount < m.LineStagingIndexPositionAndInfo.DetailPanelViewportOverflowIndexCount &&
			detailPanelTwoViewportOverflowIndexCount < m.LineStagingIndexPositionAndInfo.DetailPanelTwoViewportIndexPosition {
			detailPanelViewportActualCurrentIndex -= 1
			detailPanelTwoViewportActualCurrentIndex -= 1
		} else if detailPanelViewportOverflowIndexCount > m.LineStagingIndexPositionAndInfo.DetailPanelViewportOverflowIndexCount &&
			detailPanelTwoViewportOverflowIndexCount > m.LineStagingIndexPositionAndInfo.DetailPanelTwoViewportIndexPosition {
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
		m.LineStagingIndexPositionAndInfo = types.GittiLineStagingIndexPositionAndInfo{
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
		SetLineStagingCursorViewportContent(m, detailPanelViewportVisibleIndex, detailPanelTwoViewportVisibleIndex)
	}
}

// ------------------------------------
//
//			For Set Line Staging Cursor Viewport Content
//		  * this was to set the cursor viewport content
//	   * needed because we are using a dual viewport setup for line staging mode (one for cursor, one for content)
//
// ------------------------------------
func SetLineStagingCursorViewportContent(m *types.GittiModel, detailPanelViewportVisibleIndex int, detailPanelTwoViewportVisibleIndex int) {
	// set the cursor viewport
	var cursorVpLine strings.Builder
	var cursorVpTwoLine strings.Builder
	for index := range detailPanelViewportVisibleIndex {
		if index == m.LineStagingIndexPositionAndInfo.DetailPanelViewportIndexPosition {
			cursorVpLine.WriteString(style.SelectedItemStyle.Render("❯  ") + "\n")
		} else {
			cursorVpLine.WriteString(style.NewStyle.Render("   ") + "\n")
		}
	}

	for index := range detailPanelTwoViewportVisibleIndex {
		if index == m.LineStagingIndexPositionAndInfo.DetailPanelTwoViewportIndexPosition {
			cursorVpTwoLine.WriteString(style.SelectedItemStyle.Render("❯  ") + "\n")
		} else {
			cursorVpTwoLine.WriteString(style.NewStyle.Render("   ") + "\n")
		}
	}
	m.LineStagingIndexCursorViewport.SetContent(cursorVpLine.String())
	m.LineStagingIndexCursorTwoViewport.SetContent(cursorVpTwoLine.String())
}
