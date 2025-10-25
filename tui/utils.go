package tui

import (
	"fmt"
	"gitti/api/git"
	"strings"

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

func ProcessGitUpdate(m *GittiModel) {
	InitBranchList(m)
	InitModifiedFilesList(m)
	RenderModifiedFilesDiffViewPort(m)
	return
}

func InitBranchList(m *GittiModel) {
	items := []list.Item{
		itemString(fmt.Sprintf("* %s", git.GITBRANCH.CurrentCheckOut)),
	}

	for _, branch := range git.GITBRANCH.AllBranches {
		items = append(items, itemString(fmt.Sprintf("  %s", branch)))
	}

	m.CurrentRepoBranchesInfo = list.New(items, itemStringDelegate{}, m.HomeTabLeftPanelWidth, m.HomeTabLocalBranchesPanelHeight)
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
	for _, modifiedFiles := range latestModifiedFilesArray {
		if modifiedFiles.SelectedForStage {
			items = append(items, itemString(fmt.Sprintf(" [X] %s", modifiedFiles.FileName)))
		} else {
			items = append(items, itemString(fmt.Sprintf(" [ ] %s", modifiedFiles.FileName)))
		}
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

	m.CurrentRepoModifiedFilesInfo = list.New(items, itemStringDelegate{}, m.HomeTabLeftPanelWidth, m.HomeTabChangedFilesPanelHeight)
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
	var fileName string
	if currentSelectedModifiedFile != nil {
		fileName = ReturnFileNameFromSelectedItemList(currentSelectedModifiedFile)
	} else {
		m.CurrentSelectedFileDiffViewport.SetContent("")
		return
	}

	vpLine := fmt.Sprintf("[ %s ]\n\n", fileName)
	previousDiffRowNum := 0
	modifiedDiffRowNum := 0
	fileDiff := git.GITFILES.GetFilesDiffInfo(fileName)
	if len(fileDiff) < 1 {
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
			style = diffOldLineStyle
			modifiedDiffRowNum += 1
			rowNum = fmt.Sprintf("|%*s|%*v|  ", diffDigitLength, "", diffDigitLength, modifiedDiffRowNum)
		case git.RemoveLine:
			style = diffNewLineStyle
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

func ReturnFileNameFromSelectedItemList(currentSelectedModifiedFile list.Item) string {
	fileName := ""
	trimmedString := strings.Trim(string(currentSelectedModifiedFile.(itemString)), " ")
	var splitResult []string
	splitResult = strings.Split(trimmedString, "[ ]")
	if len(splitResult) < 2 {
		splitResult = strings.Split(trimmedString, "[X]")
	}
	fileName = strings.TrimSpace(splitResult[1])

	return fileName
}
