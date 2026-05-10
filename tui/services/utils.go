package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/component/commitlog"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/component/log"
	"github.com/gohyuhan/gitti/tui/component/reflog"
	"github.com/gohyuhan/gitti/tui/component/remote"
	"github.com/gohyuhan/gitti/tui/component/stash"
	"github.com/gohyuhan/gitti/tui/component/tag"
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

	ctx, cancel := context.WithCancel(context.Background())
	m.DetailComponentPanelInfoFetchCancelFunc = cancel
	go func(ctx context.Context) {
		defer cancel()

		// Wait for the previous goroutine to finish (its defer will set processing to false),
		// then atomically set it to true before starting a new one.

		// this is the only place for something in a bubbletea model gets modified directly in a goroutine.
		// As this is a flagging, we will allowed this, all other update must go thru channel event
		for !m.IsDetailComponentPanelInfoFetchProcessing.CompareAndSwap(false, true) {
			select {
			case <-ctx.Done():
				return
			default:
				// The previous goroutine is still running, wait a bit
				time.Sleep(10 * time.Millisecond)
			}
		}
		defer m.IsDetailComponentPanelInfoFetchProcessing.Store(false)

		var contentLine string
		var contentLine2 string // fro detail panel 2nd (only used for files changes to show staged and unstaged diff in seperated panel)
		var ogDiffLine1 []string
		var ogDiffLine2 []string
		setForDetailComponentTwo := false
		var theCurrentSelectedComponent string
		// reinit and render detail component panel viewport
		if reinit {
			m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
				Event: constant.DETAIL_COMPONENT_PANEL_LAYOUT_STATE_REINIT_EVENT,
			}
		}
		if m.CurrentSelectedComponent == constant.DetailComponentPanel || m.CurrentSelectedComponent == constant.DetailComponentPanelTwo {
			// if the current selected one is the detail component itself, the current selected one will be its parent (the component that led into the detail component)
			theCurrentSelectedComponent = m.DetailPanelParentComponent
		} else {
			theCurrentSelectedComponent = m.CurrentSelectedComponent
		}
		switch theCurrentSelectedComponent {
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			switch m.CurrentLocalBranchOrTagOrRemoteComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				contentLine = generateAboutGittiContent()
			case constant.SHOW_TAG:
				contentLine = generateTagDetailPanelContent(ctx, m)
			case constant.SHOW_REMOTE:
				contentLine = generateRemoteDetailPanelContent(m)
			}
		case constant.ModifiedFilesComponentPanel:
			contentLine, contentLine2, ogDiffLine1, ogDiffLine2, setForDetailComponentTwo = generateBothModifiedFileDetailPanelContent(ctx, m)
		case constant.CommitLogOrRefLogComponentPanel:
			contentLine = generateCommitLogOrRefLogDetailPanelContent(ctx, m)
		case constant.StashComponentPanel:
			contentLine = generateStashDetailPanelContent(ctx, m)
		case constant.LogComponentPanel:
			contentLine = generateLogDetailPanelContent(ctx, m)
		default:
			contentLine = generateAboutGittiContent()
		}

		select {
		case <-ctx.Done():
			return
		default:
			if contentLine == "" {
				contentLine = generateAboutGittiContent()
			}
			dataInterface := types.DetailPanelStateAndLayoutUpdateInterface{
				ContentLine:              contentLine,
				ContentLine2:             contentLine2,
				OgLineDiff1:              ogDiffLine1,
				OgLineDiff2:              ogDiffLine2,
				SetForDetailComponentTwo: setForDetailComponentTwo,
			}
			m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
				Event: constant.DETAIL_COMPONENT_PANEL_LAYOUT_STATE_UPDATED_EVENT,
				Data:  dataInterface,
			}
			return
		}
	}(ctx)
}

// ------------------------------------
//
//	for tag detail panel view
//
// ------------------------------------
func generateTagDetailPanelContent(ctx context.Context, m *types.GittiModel) string {
	currentSelectedTag := m.CurrentRepoTagInfoList.SelectedItem()
	var tagItem tag.GitTagItem
	if currentSelectedTag != nil {
		tagItem = currentSelectedTag.(tag.GitTagItem)
	} else {
		return ""
	}

	var vpLine strings.Builder

	tagDetail := m.GitOperations.GitTag.ShowGitTagDetail(ctx, tagItem.TagName)
	if len(tagDetail) < 1 {
		return ""
	}

	for _, Line := range tagDetail {
		line := style.NewStyle.Render(Line)
		vpLine.WriteString(line)

		vpLine.WriteRune('\n')
	}
	return vpLine.String()
}

// ------------------------------------
//
//	Generate remote detail panel content
//
// ------------------------------------
func generateRemoteDetailPanelContent(m *types.GittiModel) string {
	currentSelectedRemote := m.CurrentRepoRemoteInfoList.SelectedItem()
	var remoteItem remote.GitRemoteItem
	if currentSelectedRemote != nil {
		remoteItem = currentSelectedRemote.(remote.GitRemoteItem)
	} else {
		return ""
	}

	var vpLine strings.Builder

	vpLine.WriteString("[")
	vpLine.WriteString(style.NewStyle.Foreground(style.ColorYellowWarm).Render(remoteItem.Name))
	vpLine.WriteString("]")
	vpLine.WriteRune('\n')
	vpLine.WriteRune('\n')
	// Calculate the length of all labels to align URL, Fetch, and Push values
	urlLabel := "URL:"
	fetchLabel := i18n.LANGUAGEMAPPING.Fetch
	pushLabel := i18n.LANGUAGEMAPPING.Push
	urlLen := len([]rune(urlLabel))
	fetchLen := len([]rune(fetchLabel))
	pushLen := len([]rune(pushLabel))
	maxLen := max(urlLen, max(fetchLen, pushLen)) + 1 // plus 1 for spacing

	// Render URL with padding
	vpLine.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render(urlLabel))
	for i := 0; i < maxLen-urlLen; i++ {
		vpLine.WriteString(" ")
	}
	vpLine.WriteString(style.NewStyle.Foreground(style.ColorPurpleSoft).Render(remoteItem.Url))
	vpLine.WriteRune('\n')

	// Render Fetch with padding
	vpLine.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render(fetchLabel))
	for i := 0; i < maxLen-fetchLen; i++ {
		vpLine.WriteString(" ")
	}
	if remoteItem.Fetch {
		vpLine.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("["))
		vpLine.WriteString(style.NewStyle.Foreground(style.ColorPurpleSoft).Render("X"))
		vpLine.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("]"))
	} else {
		vpLine.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("["))
		vpLine.WriteString(" ")
		vpLine.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("]"))
	}
	vpLine.WriteRune('\n')

	// Render Push with padding
	vpLine.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render(pushLabel))
	for i := 0; i < maxLen-pushLen; i++ {
		vpLine.WriteString(" ")
	}
	if remoteItem.Push {
		vpLine.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("["))
		vpLine.WriteString(style.NewStyle.Foreground(style.ColorPurpleSoft).Render("X"))
		vpLine.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("]"))
	} else {
		vpLine.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("["))
		vpLine.WriteString(" ")
		vpLine.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("]"))
	}
	vpLine.WriteRune('\n')

	return vpLine.String()
}

// ------------------------------------
//
//	for modified file detail panel view
//
// ------------------------------------
func generateBothModifiedFileDetailPanelContent(ctx context.Context, m *types.GittiModel) (string, string, []string, []string, bool) {
	shouldRenderDetailComponentPanelTwo := false
	var fileDiffLines1 []string
	var fileDiffLines2 []string
	currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
	var fileStatus git.FileStatus
	if currentSelectedModifiedFile != nil {
		fileStatus = git.FileStatus(currentSelectedModifiedFile.(files.GitModifiedFilesItem))
	} else {
		return "", "", fileDiffLines1, fileDiffLines2, shouldRenderDetailComponentPanelTwo
	}

	var vpLine1 strings.Builder
	vpLine1.WriteString(fmt.Sprintf("[ %s ]\n\n", fileStatus.FilePathname))

	var vpLine2 strings.Builder
	getDiffTypeForVpLine1 := git.GETCOMBINEDDIFF

	// indicating that the file is not in conflict state and have both staged and unstaged changes and they are not in "?" for both index and worktree
	if !fileStatus.HasConflict && fileStatus.IndexState != " " && fileStatus.WorkTree != " " && fileStatus.IndexState != "?" && fileStatus.WorkTree != "?" {
		shouldRenderDetailComponentPanelTwo = true
		vpLine2.WriteString(fmt.Sprintf("%s\n[ %s ]\n\n", i18n.LANGUAGEMAPPING.UnstagedTitle, fileStatus.FilePathname))
		fileDiffLines2 = m.GitOperations.GitFiles.GetFilesDiffInfo(ctx, fileStatus, git.GETUNSTAGEDDIFF)

		if fileDiffLines2 == nil {
			vpLine2.WriteString(i18n.LANGUAGEMAPPING.FileTypeUnSupportedPreview)
		} else {
			for _, line := range fileDiffLines2 {
				line = style.NewStyle.Render(line)
				vpLine2.WriteString(line)
				vpLine2.WriteRune('\n')
			}
		}

		getDiffTypeForVpLine1 = git.GETSTAGEDDIFF
		vpLine1.Reset()
		vpLine1.WriteString(fmt.Sprintf("%s\n[ %s ]\n\n", i18n.LANGUAGEMAPPING.StagedTitle, fileStatus.FilePathname))
	} else if !fileStatus.HasConflict && fileStatus.IndexState != " " && fileStatus.WorkTree == " " && fileStatus.IndexState != "?" && fileStatus.WorkTree != "?" {
		getDiffTypeForVpLine1 = git.GETSTAGEDDIFF
	} else if !fileStatus.HasConflict && fileStatus.IndexState == " " && fileStatus.WorkTree != " " && fileStatus.IndexState != "?" && fileStatus.WorkTree != "?" {
		getDiffTypeForVpLine1 = git.GETUNSTAGEDDIFF
	}

	fileDiffLines1 = m.GitOperations.GitFiles.GetFilesDiffInfo(ctx, fileStatus, getDiffTypeForVpLine1)
	if fileDiffLines1 == nil {
		vpLine1.WriteString(i18n.LANGUAGEMAPPING.FileTypeUnSupportedPreview)
	} else {
		for _, line := range fileDiffLines1 {
			line = style.NewStyle.Render(line)
			vpLine1.WriteString(line)
			vpLine1.WriteRune('\n')
		}
	}

	return vpLine1.String(), vpLine2.String(), fileDiffLines1, fileDiffLines2, shouldRenderDetailComponentPanelTwo
}

// ------------------------------------
//
//	for commit log detail panel view
//
// ------------------------------------
func generateCommitLogOrRefLogDetailPanelContent(ctx context.Context, m *types.GittiModel) string {
	var hash string

	switch m.CurrentCommitLogOrRefLogComponentShowing {
	case constant.SHOW_COMMITLOG:
		currentSelectedCommitLog := m.CurrentRepoCommitLogInfoList.SelectedItem()
		if currentSelectedCommitLog != nil {
			hash = currentSelectedCommitLog.(commitlog.GitCommitLogItem).Hash
		} else {
			return ""
		}
	case constant.SHOW_REFLOG:
		currentSelectedRefLog := m.CurrentRepoRefLogInfoList.SelectedItem()
		if currentSelectedRefLog != nil {
			item := currentSelectedRefLog.(reflog.GitRefLogItem)
			hash = item.Hash
		} else {
			return ""
		}
	}

	var vpLine strings.Builder
	commitLogDetail := m.GitOperations.GitCommitLog.GitCommitLogDetail(ctx, hash)
	if len(commitLogDetail) < 1 {
		return ""
	}

	for _, Line := range commitLogDetail {
		line := style.NewStyle.Render(Line)
		vpLine.WriteString(line)
		vpLine.WriteRune('\n')
	}
	return vpLine.String()
}

// ------------------------------------
//
//	for stash detail panel view
//
// ------------------------------------
func generateStashDetailPanelContent(ctx context.Context, m *types.GittiModel) string {
	currentSelectedStash := m.CurrentRepoStashInfoList.SelectedItem()
	var stashItem stash.GitStashItem
	if currentSelectedStash != nil {
		stashItem = currentSelectedStash.(stash.GitStashItem)
	} else {
		return ""
	}

	var vpLine strings.Builder

	stashDetail := m.GitOperations.GitStash.GitStashDetail(ctx, stashItem.Id)
	if len(stashDetail) < 1 {
		return ""
	}

	vpLine.WriteString(fmt.Sprintf(
		"[%s]\n[%s]\n\n",
		style.StashIdStyle.Render(stashItem.Id),
		style.StashMessageStyle.Render(stashItem.Message),
	))

	for _, Line := range stashDetail {
		line := style.NewStyle.Render(Line)
		vpLine.WriteString(line)
		vpLine.WriteRune('\n')
	}
	return vpLine.String()
}

// ------------------------------------
//
//	Generate log detail panel content
//
// ------------------------------------
func generateLogDetailPanelContent(ctx context.Context, m *types.GittiModel) string {
	vpLine := log.InitGittiLogViewport(m, false, ctx)
	return vpLine
}

// ------------------------------------
//
//	for about gitti content
//
// ------------------------------------
func generateAboutGittiContent() string {
	var vpLine strings.Builder

	logoLineArray := style.GradientLines(constant.GittiAsciiArtLogo)
	aboutLines := i18n.LANGUAGEMAPPING.AboutGitti

	vpLine.WriteString(strings.Join(logoLineArray, "\n"))
	vpLine.WriteRune('\n')
	vpLine.WriteString(strings.Join(aboutLines, "\n"))

	return vpLine.String()
}

// ------------------------------------
//
//	Trigger async git fetch operation via daemon update channel
//
// ------------------------------------
func GitFetchService(m *types.GittiModel) {
	go func() {
		m.DaemonUpdateChannel <- git.GIT_FETCH
	}()
}
