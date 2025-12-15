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
			m.IsDetailComponentPanelInfoFetchProcessing.Store(false)
			cancel()
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
		if m.CurrentSelectedComponent == constant.DetailComponent {
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

	vpLine1 := fmt.Sprintf("[ %s ]\n\n", fileStatus.FilePathname)
	var vpLine2 string
	var fileDiffLines1 []string
	var fileDiffLines2 []string
	getDiffTypeForVpLine1 := git.GETCOMBINEDDIFF

	// indicating that the file is not in conflict state and have both staged and unstaged changes
	if !fileStatus.HasConflict && fileStatus.IndexState != " " && fileStatus.WorkTree != " " {
		shouldRenderDetailComponentPanelTwo = true
		vpLine2 = fmt.Sprintf("%s\n\n[ %s ]\n\n", i18n.LANGUAGEMAPPING.UnstagedTitle, fileStatus.FilePathname)
		fileDiffLines2 = m.GitOperations.GitFiles.GetFilesDiffInfo(ctx, fileStatus, git.GETUNSTAGEDDIFF)

		if fileDiffLines2 == nil {
			vpLine2 += i18n.LANGUAGEMAPPING.FileTypeUnSupportedPreview
		} else {
			for _, line := range fileDiffLines2 {
				line = style.NewStyle.Render(line)
				vpLine2 += line + "\n"
			}
		}

		getDiffTypeForVpLine1 = git.GETSTAGEDDIFF
		vpLine1 = fmt.Sprintf("%s\n\n[ %s ]\n\n", i18n.LANGUAGEMAPPING.StagedTitle, fileStatus.FilePathname)
	}

	fileDiffLines1 = m.GitOperations.GitFiles.GetFilesDiffInfo(ctx, fileStatus, getDiffTypeForVpLine1)
	if fileDiffLines1 == nil {
		vpLine1 += i18n.LANGUAGEMAPPING.FileTypeUnSupportedPreview
	} else {
		for _, line := range fileDiffLines1 {
			line = style.NewStyle.Render(line)
			vpLine1 += line + "\n"
		}
	}

	return vpLine1, vpLine2, shouldRenderDetailComponentPanelTwo
}

// for commit log detail panel view
func generateCommitLogDetailPanelContent(ctx context.Context, m *types.GittiModel) string {
	currentSelectedCommitLog := m.CurrentRepoCommitLogInfoList.SelectedItem()
	var commitLogItem commitlog.GitCommitLogItem
	var vpLine string
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
		vpLine += line + "\n"
	}
	return vpLine
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

	vpLine := fmt.Sprintf(
		"[%s]\n[%s]\n\n",
		style.StashIdStyle.Render(stashItem.Id),
		style.StashMessageStyle.Render(stashItem.Message),
	)

	stashDetail := m.GitOperations.GitStash.GitStashDetail(ctx, stashItem.Id)
	if len(stashDetail) < 1 {
		return ""
	}

	for _, Line := range stashDetail {
		line := style.NewStyle.Render(Line)
		vpLine += line + "\n"
	}
	return vpLine
}

// for about gitti content
func generateAboutGittiContent() string {
	var vpLine string

	logoLineArray := style.GradientLines(constant.GittiAsciiArtLogo)
	aboutLines := i18n.LANGUAGEMAPPING.AboutGitti

	vpLine += strings.Join(logoLineArray, "\n") + "\n"
	vpLine += strings.Join(aboutLines, "\n")

	return vpLine
}
