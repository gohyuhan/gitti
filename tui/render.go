package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/v2"

	"gitti/api/git"
	"gitti/i18n"
	"gitti/settings"
	"gitti/tui/constant"
	"gitti/tui/style"
	"gitti/tui/utils"
)

// -----------------------------------------------------------------------------
//
//	Functions that help construct the view
//
// -----------------------------------------------------------------------------
// Render the Local Branches panel (top 25%)
func renderLocalBranchesPanel(width int, height int, m *GittiModel) string {
	borderStyle := style.PanelBorderStyle
	if m.CurrentSelectedContainer == constant.LocalBranchComponent {
		borderStyle = style.SelectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(m.CurrentRepoBranchesInfoList.View())
}

// Render the Changed Files panel (bottom 75%)
func renderChangedFilesPanel(width int, height int, m *GittiModel) string {
	borderStyle := style.PanelBorderStyle
	if m.CurrentSelectedContainer == constant.ModifiedFilesComponent {
		borderStyle = style.SelectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(m.CurrentRepoModifiedFilesInfoList.View())
}

func renderFileDiffPanel(width int, height int, m *GittiModel) string {
	borderStyle := style.PanelBorderStyle
	if m.CurrentSelectedContainer == constant.FileDiffComponent {
		borderStyle = style.SelectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(m.CurrentSelectedFileDiffViewport.View())
}

func renderKeyBindingPanel(width int, m *GittiModel) string {
	var keys []string
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.CommitPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForCommitPopUp
		case constant.AddRemotePromptPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForAddRemotePromptPopUp
		case constant.GitRemotePushPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitRemotePushPopUp
		case constant.ChooseRemotePopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChooseRemotePopUp
		case constant.ChoosePushTypePopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChoosePushTypePopUp
		}
	} else {
		switch m.CurrentSelectedContainer {
		case constant.NoneSelected:
			keys = i18n.LANGUAGEMAPPING.KeyBindingNoneSelected
		case constant.LocalBranchComponent:
			CurrentSelectedBranch := m.CurrentRepoBranchesInfoList.SelectedItem()
			if CurrentSelectedBranch == nil {
				keys = i18n.LANGUAGEMAPPING.KeyBindingLocalBranchComponentNone
			} else {
				isCurrentSelectedBranchCheckedOutBranch := CurrentSelectedBranch.(gitBranchItem).IsCheckedOut
				if isCurrentSelectedBranchCheckedOutBranch {
					keys = i18n.LANGUAGEMAPPING.KeyBindingLocalBranchComponentIsCheckOut
				} else {
					keys = i18n.LANGUAGEMAPPING.KeyBindingLocalBranchComponentDefault
				}
			}
		case constant.ModifiedFilesComponent:
			CurrentSelectedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
			if CurrentSelectedFile == nil {
				keys = i18n.LANGUAGEMAPPING.KeyBindingModifiedFilesComponentNone
			} else {
				isFileStaged := CurrentSelectedFile.(gitModifiedFilesItem).SelectedForStage
				if isFileStaged {
					keys = i18n.LANGUAGEMAPPING.KeyBindingModifiedFilesComponentIsStaged
				} else {
					keys = i18n.LANGUAGEMAPPING.KeyBindingModifiedFilesComponentDefault
				}
			}
		case constant.FileDiffComponent:
			keys = i18n.LANGUAGEMAPPING.KeyBindingFileDiffComponent
		}
	}

	distributedWidth := (width / len(keys))

	var keyBindingLine string

	for _, key := range keys {
		truncated := utils.TruncateString(key, distributedWidth) // truncate manually
		cell := style.NewStyle.
			Width(distributedWidth).    // fixed box width
			MaxWidth(distributedWidth). // disallow overflow expansion
			Align(lipgloss.Center).
			Render(truncated)
		keyBindingLine += cell
	}

	return style.BottomKeyBindingStyle.
		Width(width).
		Height(constant.MainPageKeyBindingLayoutPanelHeight).
		Render(keyBindingLine)
}

// for the current selected modified file preview viewport
func renderModifiedFilesDiffViewPort(m *GittiModel) {
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
	fileDiff := m.GitState.GitFiles.GetFilesDiffInfo(fileStatus)
	if fileDiff == nil {
		vpLine += i18n.LANGUAGEMAPPING.FileTypeUnSupportedPreview
		m.CurrentSelectedFileDiffViewport.SetContent(vpLine)
		return
	}
	diffDigitLength := len(fmt.Sprintf("%d", len(fileDiff))) + 1
	for _, Line := range fileDiff {
		var diffLine string
		var rowNum string
		lineStyle := style.NewStyle
		switch Line.Type {
		case git.AddLine:
			lineStyle = style.DiffNewLineStyle
			modifiedDiffRowNum += 1
			rowNum = fmt.Sprintf("|%*s|%*v|  ", diffDigitLength, "", diffDigitLength, modifiedDiffRowNum)
		case git.RemoveLine:
			lineStyle = style.DiffOldLineStyle
			previousDiffRowNum += 1
			rowNum = fmt.Sprintf("|%*v|%*s|  ", diffDigitLength, previousDiffRowNum, diffDigitLength, "")
		default:
			previousDiffRowNum += 1
			modifiedDiffRowNum += 1
			rowNum = fmt.Sprintf("|%*v|%*v|  ", diffDigitLength, previousDiffRowNum, diffDigitLength, modifiedDiffRowNum)
		}

		diffLine = lineStyle.Render(Line.Line)
		vpLine += rowNum + diffLine + "\n"
	}
	m.CurrentSelectedFileDiffViewport.SetContent(vpLine)
}

// to update the width and height of all components
func tuiWindowSizing(m *GittiModel) {
	// Compute panel widths
	m.HomeTabLeftPanelWidth = min(int(float64(m.Width)*settings.GITTICONFIGSETTINGS.LeftPanelWidthRatio), constant.MaxLeftPanelWidth)
	m.HomeTabFileDiffPanelWidth = m.Width - m.HomeTabLeftPanelWidth

	m.HomeTabCoreContentHeight = m.Height - constant.MainPageKeyBindingLayoutPanelHeight - 2*constant.Padding
	m.HomeTabFileDiffPanelHeight = m.HomeTabCoreContentHeight

	leftPanelRemainingHeight := m.HomeTabCoreContentHeight - 3 // this is after reserving the height for the gitti version panel
	m.HomeTabLocalBranchesPanelHeight = int(float64(leftPanelRemainingHeight)*settings.GITTICONFIGSETTINGS.GitBranchComponentHeightRatio) - 2*constant.Padding
	m.HomeTabChangedFilesPanelHeight = leftPanelRemainingHeight - m.HomeTabLocalBranchesPanelHeight - 2*constant.Padding

	// update all components Width and Height
	m.CurrentRepoBranchesInfoList.SetWidth(m.HomeTabLeftPanelWidth - 2)
	m.CurrentRepoBranchesInfoList.SetHeight(m.HomeTabLocalBranchesPanelHeight)
	// m.CurrentRepoBranchesInfoList.Title = truncateString(fmt.Sprintf("[b]  %s:", i18n.LANGUAGEMAPPING.Branches), m.HomeTabLeftPanelWidth - listItemOrTitleWidthPad -2 )

	m.CurrentRepoModifiedFilesInfoList.SetWidth(m.HomeTabLeftPanelWidth - 2)
	m.CurrentRepoModifiedFilesInfoList.SetHeight(m.HomeTabChangedFilesPanelHeight)
	// m.CurrentRepoModifiedFilesInfoList.Title = truncateString(fmt.Sprintf("[f] 📄%s:", i18n.LANGUAGEMAPPING.ModifiedFiles), m.HomeTabLeftPanelWidth - listItemOrTitleWidthPad - 2)

	// update viewport
	m.CurrentSelectedFileDiffViewport.SetHeight(m.HomeTabFileDiffPanelHeight) //some margin
	m.CurrentSelectedFileDiffViewport.SetWidth(m.HomeTabFileDiffPanelWidth - 2)
	m.CurrentSelectedFileDiffViewportOffset = max(0, int(m.CurrentSelectedFileDiffViewport.HorizontalScrollPercent()*float64(m.CurrentSelectedFileDiffViewportOffset))-1)
	m.CurrentSelectedFileDiffViewport.SetXOffset(m.CurrentSelectedFileDiffViewportOffset)
	m.CurrentSelectedFileDiffViewport.SetYOffset(m.CurrentSelectedFileDiffViewport.YOffset)

	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.CommitPopUp:
			popUp, exist := m.PopUpModel.(*GitCommitPopUpModel)
			if exist {
				width := (min(constant.MaxCommitPopUpWidth, int(float64(m.Width)*0.8)) - 4)
				popUp.GitCommitOutputViewport.SetWidth(width)
			}
		case constant.GitRemotePushPopUp:
			popUp, exist := m.PopUpModel.(*GitRemotePushPopUpModel)
			if exist {
				width := (min(constant.MaxGitRemotePushPopUpWidth, int(float64(m.Width)*0.8)) - 4)
				popUp.GitRemotePushOutputViewport.SetWidth(width)
			}
		case constant.ChoosePushTypePopUp:
			popUp, exist := m.PopUpModel.(*ChoosePushTypePopUpModel)
			if exist {
				width := (min(constant.MaxChoosePushTypePopUpWidth, int(float64(m.Width)*0.8)) - 4)
				popUp.PushOptionList.SetWidth(width)
			}
		case constant.ChooseRemotePopUp:
			popUp, exist := m.PopUpModel.(*ChooseRemotePopUpModel)
			if exist {
				width := (min(constant.MaxChooseRemotePopUpWidth, int(float64(m.Width)*0.8)) - 4)
				popUp.RemoteList.SetWidth(width)
			}
		}
	}

}
