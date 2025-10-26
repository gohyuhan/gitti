package tui

import (
	"fmt"
	"gitti/api/git"
	"gitti/settings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// variables for indicating which panel/components/container or whatever the hell you wanna call it that the user is currently landed or selected, so that they can do precious action related to the part of whatever the hell you wanna call it
const (
	None = "0"

	LocalBranchComponent   = "C1"
	ModifiedFilesComponent = "C2"
	FileDiffComponent      = "C3"
)

func TuiWindowSizing(m *GittiModel) {
	// Compute panel widths
	m.HomeTabLeftPanelWidth = int(float64(m.Width) * settings.GITTICONFIGSETTINGS.LeftPanelWidthRatio)
	m.HomeTabFileDiffPanelWidth = m.Width - m.HomeTabLeftPanelWidth - 4 // adjust for borders/padding

	m.HomeTabCoreContentHeight = m.Height - mainPageKeyBindingLayoutPanelHeight - 2*padding
	m.HomeTabFileDiffPanelHeight = m.HomeTabCoreContentHeight
	m.HomeTabLocalBranchesPanelHeight = int(float64(m.HomeTabCoreContentHeight)*settings.GITTICONFIGSETTINGS.GitBranchComponentHeightRatio) - 2*padding
	m.HomeTabChangedFilesPanelHeight = m.HomeTabCoreContentHeight - m.HomeTabLocalBranchesPanelHeight - 2*padding

	// update all components Width and Height
	m.CurrentRepoBranchesInfo.SetWidth(m.HomeTabLeftPanelWidth)
	m.CurrentRepoBranchesInfo.SetHeight(m.HomeTabLocalBranchesPanelHeight)

	// update viewport
	m.CurrentSelectedFileDiffViewport.Height = m.HomeTabFileDiffPanelHeight - 1 //some margin
	m.CurrentSelectedFileDiffViewport.Width = m.HomeTabFileDiffPanelWidth
	m.CurrentSelectedFileDiffViewportOffset = max(0, int(m.CurrentSelectedFileDiffViewport.HorizontalScrollPercent()*float64(m.CurrentSelectedFileDiffViewportOffset))-1)
	m.CurrentSelectedFileDiffViewport.SetXOffset(m.CurrentSelectedFileDiffViewportOffset)
	m.CurrentSelectedFileDiffViewport.SetYOffset(m.CurrentSelectedFileDiffViewport.YOffset)
}

func ProcessGitUpdate(m *GittiModel) {
	InitBranchList(m)
	InitModifiedFilesList(m)
	RenderModifiedFilesDiffViewPort(m)
	return
}

func InitBranchList(m *GittiModel) {
	items := []list.Item{
		gitBranchItem(git.GITBRANCH.CurrentCheckOut),
	}

	for _, branch := range git.GITBRANCH.AllBranches {
		items = append(items, gitBranchItem(branch))
	}

	m.CurrentRepoBranchesInfo = list.New(items, gitBranchItemDelegate{}, m.HomeTabLeftPanelWidth, m.HomeTabLocalBranchesPanelHeight)
	m.CurrentRepoBranchesInfo.Title = "[b]  Branches:"
	m.CurrentRepoBranchesInfo.SetShowStatusBar(false)
	m.CurrentRepoBranchesInfo.SetFilteringEnabled(false)
	m.CurrentRepoBranchesInfo.SetShowHelp(false)
	m.CurrentRepoBranchesInfo.Styles.Title = titleStyle
	m.CurrentRepoBranchesInfo.Styles.PaginationStyle = paginationStyle

	if m.NavigationIndexPosition.LocalBranchComponent > len(m.CurrentRepoBranchesInfo.Items())-1 {
		m.CurrentRepoBranchesInfo.Select(len(m.CurrentRepoBranchesInfo.Items()) - 1)
	} else {
		m.CurrentRepoBranchesInfo.Select(m.NavigationIndexPosition.LocalBranchComponent)
	}

	return
}

func InitModifiedFilesList(m *GittiModel) {
	latestModifiedFilesArray := git.GITFILES.FilesStatus
	items := make([]list.Item, 0, len(latestModifiedFilesArray))
	for _, modifiedFile := range latestModifiedFilesArray {
		items = append(items, gitModifiedFilesItem(modifiedFile))
	}

	// get the previous selected file and see if it was within the new list if yes get the latest position of the previous selected file
	previousSelectedFile := m.CurrentRepoBranchesInfo.SelectedItem()
	selectedFilesPosition := -1

	for index, item := range items {
		if item == previousSelectedFile {
			selectedFilesPosition = index
			break
		}
	}

	m.CurrentRepoModifiedFilesInfo = list.New(items, gitModifiedFilesItemDelegate{}, m.HomeTabLeftPanelWidth, m.HomeTabChangedFilesPanelHeight)
	m.CurrentRepoModifiedFilesInfo.Title = "[f] 📄 Modified Files:"
	m.CurrentRepoModifiedFilesInfo.SetShowStatusBar(false)
	m.CurrentRepoModifiedFilesInfo.SetFilteringEnabled(false)
	m.CurrentRepoModifiedFilesInfo.SetShowHelp(false)
	m.CurrentRepoModifiedFilesInfo.SetShowPagination(false)
	m.CurrentRepoModifiedFilesInfo.Styles.Title = titleStyle

	if len(items) < 1 {
		return
	}

	if selectedFilesPosition >= 0 {
		m.CurrentRepoModifiedFilesInfo.Select(selectedFilesPosition)
		m.NavigationIndexPosition.ModifiedFilesComponent = selectedFilesPosition
	} else {
		if m.NavigationIndexPosition.ModifiedFilesComponent > len(m.CurrentRepoModifiedFilesInfo.Items())-1 {
			m.CurrentRepoModifiedFilesInfo.Select(len(m.CurrentRepoModifiedFilesInfo.Items()) - 1)
		} else {
			m.CurrentRepoModifiedFilesInfo.Select(m.NavigationIndexPosition.ModifiedFilesComponent)
		}
	}
	return
}

// reinit and render diff file viewport
func ReinitAndRenderModifiedFileDiffViewPort(m *GittiModel) {
	m.CurrentSelectedFileDiffViewportOffset = 0
	m.CurrentSelectedFileDiffViewport.SetXOffset(0)
	m.CurrentSelectedFileDiffViewport.SetYOffset(0)
	RenderModifiedFilesDiffViewPort(m)
}

// for the current selected modified file preview viewport
func RenderModifiedFilesDiffViewPort(m *GittiModel) {
	currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfo.SelectedItem()
	var fileStatus git.FileStatus
	if currentSelectedModifiedFile != nil {
		fileStatus = git.FileStatus(currentSelectedModifiedFile.(gitModifiedFilesItem))
	} else {
		m.CurrentSelectedFileDiffViewport.SetContent("")
		return
	}

	vpLine := fmt.Sprintf("[ %s ]\n\n", fileStatus.FileName)
	previousDiffRowNum := 0
	modifiedDiffRowNum := 0
	fileDiff := git.GITFILES.GetFilesDiffInfo(fileStatus)
	if fileDiff == nil {
		vpLine += "The current selected file type is not supported for preview"
		m.CurrentSelectedFileDiffViewport.SetContent(vpLine)
		return
	}
	diffDigitLength := len(fmt.Sprintf("%d", len(fileDiff))) + 1
	for _, Line := range fileDiff {
		var diffLine string
		var rowNum string
		style := lipgloss.NewStyle()
		switch Line.Type {
		case git.AddLine:
			style = diffNewLineStyle
			modifiedDiffRowNum += 1
			rowNum = fmt.Sprintf("|%*s|%*v|  ", diffDigitLength, "", diffDigitLength, modifiedDiffRowNum)
		case git.RemoveLine:
			style = diffOldLineStyle
			previousDiffRowNum += 1
			rowNum = fmt.Sprintf("|%*v|%*s|  ", diffDigitLength, previousDiffRowNum, diffDigitLength, "")
		default:
			previousDiffRowNum += 1
			modifiedDiffRowNum += 1
			rowNum = fmt.Sprintf("|%*v|%*v|  ", diffDigitLength, previousDiffRowNum, diffDigitLength, modifiedDiffRowNum)
		}

		diffLine = style.Render(Line.Line)
		vpLine += rowNum + diffLine + "\n"
	}
	m.CurrentSelectedFileDiffViewport.SetContent(vpLine)
}

// to prevent wrapping and overflow doesn't seems to be supported natively in bubbletea yet
func TruncateString(s string, width int) string {
	if utf8.RuneCountInString(s) <= width {
		return s
	}
	runes := []rune(s)
	if width > 1 {
		return string(runes[:width-1]) + "…"
	}
	return string(runes[:width])
}
