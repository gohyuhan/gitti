package tui

import (
	"fmt"
	"gitti/api/git"
	"gitti/i18n"
	"gitti/settings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/textarea"
	"github.com/charmbracelet/bubbles/v2/textinput"
	"github.com/charmbracelet/bubbles/v2/viewport"
	"github.com/charmbracelet/lipgloss/v2"
)

// variables for indicating which panel/components/container or whatever the hell you wanna call it that the user is currently landed or selected, so that they can do precious action related to the part of whatever the hell you wanna call it
const (
	NoneSelected = "0"

	LocalBranchComponent   = "C1"
	ModifiedFilesComponent = "C2"
	FileDiffComponent      = "C3"
)

const (
	None        = "None"
	CommitPopUp = "CommitPopUp"
)

func TuiWindowSizing(m *GittiModel) {
	// Compute panel widths
	m.HomeTabLeftPanelWidth = min(int(float64(m.Width)*settings.GITTICONFIGSETTINGS.LeftPanelWidthRatio), maxLeftPanelWidth)
	m.HomeTabFileDiffPanelWidth = m.Width - m.HomeTabLeftPanelWidth

	m.HomeTabCoreContentHeight = m.Height - mainPageKeyBindingLayoutPanelHeight - 2*padding
	m.HomeTabFileDiffPanelHeight = m.HomeTabCoreContentHeight

	leftPanelRemainingHeight := m.HomeTabCoreContentHeight - 3 // this is after reserving the height for the gitti version panel
	m.HomeTabLocalBranchesPanelHeight = int(float64(leftPanelRemainingHeight)*settings.GITTICONFIGSETTINGS.GitBranchComponentHeightRatio) - 2*padding
	m.HomeTabChangedFilesPanelHeight = leftPanelRemainingHeight - m.HomeTabLocalBranchesPanelHeight - 2*padding

	// update all components Width and Height
	m.CurrentRepoBranchesInfoList.SetWidth(m.HomeTabLeftPanelWidth - 2)
	m.CurrentRepoBranchesInfoList.SetHeight(m.HomeTabLocalBranchesPanelHeight)

	m.CurrentRepoModifiedFilesInfoList.SetWidth(m.HomeTabLeftPanelWidth - 2)
	m.CurrentRepoModifiedFilesInfoList.SetHeight(m.HomeTabChangedFilesPanelHeight)

	// update viewport
	m.CurrentSelectedFileDiffViewport.SetHeight(m.HomeTabFileDiffPanelHeight) //some margin
	m.CurrentSelectedFileDiffViewport.SetWidth(m.HomeTabFileDiffPanelWidth - 2)
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

	m.CurrentRepoBranchesInfoList = list.New(items, gitBranchItemDelegate{}, m.HomeTabLeftPanelWidth, m.HomeTabLocalBranchesPanelHeight)
	m.CurrentRepoBranchesInfoList.Title = fmt.Sprintf("[b]  %s:", i18n.LANGUAGEMAPPING.Branches)
	m.CurrentRepoBranchesInfoList.SetShowStatusBar(false)
	m.CurrentRepoBranchesInfoList.SetFilteringEnabled(false)
	m.CurrentRepoBranchesInfoList.SetShowHelp(false)
	m.CurrentRepoBranchesInfoList.Styles.Title = titleStyle
	m.CurrentRepoBranchesInfoList.Styles.PaginationStyle = paginationStyle

	if m.NavigationIndexPosition.LocalBranchComponent > len(m.CurrentRepoBranchesInfoList.Items())-1 {
		m.CurrentRepoBranchesInfoList.Select(len(m.CurrentRepoBranchesInfoList.Items()) - 1)
	} else {
		m.CurrentRepoBranchesInfoList.Select(m.NavigationIndexPosition.LocalBranchComponent)
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
	previousSelectedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
	selectedFilesPosition := -1

	for index, item := range items {
		if item == previousSelectedFile {
			selectedFilesPosition = index
			break
		}
	}

	m.CurrentRepoModifiedFilesInfoList = list.New(items, gitModifiedFilesItemDelegate{}, m.HomeTabLeftPanelWidth, m.HomeTabChangedFilesPanelHeight)
	m.CurrentRepoModifiedFilesInfoList.Title = fmt.Sprintf("[f] 📄%s:", i18n.LANGUAGEMAPPING.ModifiedFiles)
	m.CurrentRepoModifiedFilesInfoList.SetShowStatusBar(false)
	m.CurrentRepoModifiedFilesInfoList.SetFilteringEnabled(false)
	m.CurrentRepoModifiedFilesInfoList.SetShowHelp(false)
	m.CurrentRepoModifiedFilesInfoList.SetShowPagination(false)
	m.CurrentRepoModifiedFilesInfoList.Styles.Title = titleStyle

	if len(items) < 1 {
		return
	}

	if selectedFilesPosition >= 0 {
		m.CurrentRepoModifiedFilesInfoList.Select(selectedFilesPosition)
		m.NavigationIndexPosition.ModifiedFilesComponent = selectedFilesPosition
	} else {
		if m.NavigationIndexPosition.ModifiedFilesComponent > len(m.CurrentRepoModifiedFilesInfoList.Items())-1 {
			m.CurrentRepoModifiedFilesInfoList.Select(len(m.CurrentRepoModifiedFilesInfoList.Items()) - 1)
		} else {
			m.CurrentRepoModifiedFilesInfoList.Select(m.NavigationIndexPosition.ModifiedFilesComponent)
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
	currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
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
		vpLine += i18n.LANGUAGEMAPPING.FileTypeUnSupportedPreview
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
		return string(runes[:width-3]) + "…  "
	}
	return string(runes[:width])
}

func InitGitCommitPopUpModel(m *GittiModel) {
	CommitMessageTextInput := textinput.New()
	CommitMessageTextInput.Placeholder = i18n.LANGUAGEMAPPING.CommitPopUpMessageInputPlaceHolder
	CommitMessageTextInput.Focus()

	CommitDescriptionTextAreaInput := textarea.New()
	CommitDescriptionTextAreaInput.Placeholder = i18n.LANGUAGEMAPPING.CommitPopUpCommitDescriptionInputPlaceHolder
	CommitDescriptionTextAreaInput.ShowLineNumbers = false
	CommitDescriptionTextAreaInput.SetHeight(5)
	CommitDescriptionTextAreaInput.Blur()

	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1

	m.PopUpModel = &GitCommitPopUpModel{
		MessageTextInput:         CommitMessageTextInput,
		DescriptionTextAreaInput: CommitDescriptionTextAreaInput,
		TotalInputCount:          2,
		CurrentActiveInputIndex:  1,
		GitCommitOutputViewport:  vp,
	}
}
