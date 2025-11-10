package tui

import (
	"fmt"
	"strings"

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
// Render the Local Branches panel
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

// Render the Changed Files panel
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
		Render(m.DetailPanelViewport.View())
}

func renderKeyBindingPanel(width int, m *GittiModel) string {
	keys:=[]string{""} // to prevent a misconfiguration on key binding will not crash the program
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
		case constant.ChooseNewBranchTypePopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChooseNewBranchTypePopUp
		case constant.CreateNewBranchPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForCreateNewBranchPopUp
		case constant.ChooseSwitchBranchTypePopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChooseSwitchBranchTypePopUp
		case constant.SwitchBranchOutputPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChooseSwitchBranchTypePopUp
			popUp, ok := m.PopUpModel.(*SwitchBranchOutputPopUpModel)
			if ok {
				if popUp.IsProcessing.Load() {
					keys = []string{"..."} // nothing can be done during switching, only force quit gitti is possible
				}
			}
		case constant.ChooseGitPullTypePopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChooseGitPullTypePopUp
		case constant.GitPullOutputPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitPullOutputPopUp
		}
	} else {
		switch m.CurrentSelectedContainer {
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


	var keyBindingLine string
	keyBindingLine = "  "+strings.Join(keys, "  |  ")
	keyBindingLine =  utils.TruncateString(keyBindingLine, width) 

	return style.BottomKeyBindingStyle.
		Width(width).
		Height(constant.MainPageKeyBindingLayoutPanelHeight).
		Render(keyBindingLine)
}

// for the current selected modified file preview viewport
func renderDetailPanelViewPort(m *GittiModel) {
	currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
	var fileStatus git.FileStatus
	if currentSelectedModifiedFile != nil {
		fileStatus = git.FileStatus(currentSelectedModifiedFile.(gitModifiedFilesItem))
	} else {
		m.DetailPanelViewport.SetContent("")
		return
	}

	vpLine := fmt.Sprintf("[ %s ]\n\n", fileStatus.FileName)
	previousDiffRowNum := 0
	modifiedDiffRowNum := 0
	fileDiff := m.GitState.GitFiles.GetFilesDiffInfo(fileStatus)
	if fileDiff == nil {
		vpLine += i18n.LANGUAGEMAPPING.FileTypeUnSupportedPreview
		m.DetailPanelViewport.SetContent(vpLine)
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
	m.DetailPanelViewport.SetContent(vpLine)
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
	m.DetailPanelViewport.SetHeight(m.HomeTabFileDiffPanelHeight) //some margin
	m.DetailPanelViewport.SetWidth(m.HomeTabFileDiffPanelWidth - 2)
	m.DetailPanelViewportOffset = max(0, int(m.DetailPanelViewport.HorizontalScrollPercent()*float64(m.DetailPanelViewportOffset))-1)
	m.DetailPanelViewport.SetXOffset(m.DetailPanelViewportOffset)
	m.DetailPanelViewport.SetYOffset(m.DetailPanelViewport.YOffset)

	// update list of viewport component width within pop up
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
		case constant.ChooseNewBranchTypePopUp:
			popUp, exist := m.PopUpModel.(*ChooseNewBranchTypeOptionPopUpModel)
			if exist {
				width := (min(constant.MaxChooseNewBranchTypePopUpWidth, int(float64(m.Width)*0.8)) - 4)
				popUp.NewBranchTypeOptionList.SetWidth(width)
			}
		case constant.ChooseSwitchBranchTypePopUp:
			popUp, exist := m.PopUpModel.(*ChooseSwitchBranchTypePopUpModel)
			if exist {
				width := (min(constant.MaxChooseSwitchBranchTypePopUpWidth, int(float64(m.Width)*0.8)) - 4)
				popUp.SwitchTypeOptionList.SetWidth(width)
			}
		case constant.ChooseGitPullTypePopUp:
			popUp, exist := m.PopUpModel.(*ChooseGitPullTypePopUpModel)
			if exist {
				width := (min(constant.MaxChooseGitPullTypePopUpWidth, int(float64(m.Width)*0.8)) - 4)
				popUp.PullTypeOptionList.SetWidth(width)
			}
		case constant.GitPullOutputPopUp:
			popUp, exist := m.PopUpModel.(*GitPullOutputPopUpModel)
			if exist {
				width := (min(constant.MaxGitPullOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)
				popUp.GitPullOutputViewport.SetWidth(width)
			}
		}
	}
}
