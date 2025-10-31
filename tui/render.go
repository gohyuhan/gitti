package tui

import (
	"fmt"
	"gitti/api/git"
	"gitti/i18n"

	"github.com/charmbracelet/lipgloss/v2"
)

// -----------------------------------------------------------------------------
//
//	Functions that help construct the view
//
// -----------------------------------------------------------------------------
// Render the Local Branches panel (top 25%)
func renderLocalBranchesPanel(width int, height int, m *GittiModel) string {
	borderStyle := panelBorderStyle
	if m.CurrentSelectedContainer == LocalBranchComponent {
		borderStyle = selectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(m.CurrentRepoBranchesInfoList.View())
}

// Render the Changed Files panel (bottom 75%)
func renderChangedFilesPanel(width int, height int, m *GittiModel) string {
	borderStyle := panelBorderStyle
	if m.CurrentSelectedContainer == ModifiedFilesComponent {
		borderStyle = selectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(m.CurrentRepoModifiedFilesInfoList.View())
}

func renderFileDiffPanel(width int, height int, m *GittiModel) string {
	borderStyle := panelBorderStyle
	if m.CurrentSelectedContainer == FileDiffComponent {
		borderStyle = selectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(m.CurrentSelectedFileDiffViewport.View())
}

func renderKeyBindingPanel(width int, m *GittiModel) string {
	var keys []string
	if m.ShowPopUp {
		switch m.PopUpType {
		case CommitPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForCommitPopUp
		}
	} else {
		switch m.CurrentSelectedContainer {
		case NoneSelected:
			keys = i18n.LANGUAGEMAPPING.KeyBindingNoneSelected
		case LocalBranchComponent:
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
		case ModifiedFilesComponent:
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
		case FileDiffComponent:
			keys = i18n.LANGUAGEMAPPING.KeyBindingFileDiffComponent
		}
	}

	distributedWidth := width / len(keys)
	var keyBindingLine string

	for _, key := range keys {
		truncatedKey := truncateString(key, distributedWidth)
		cell := lipgloss.NewStyle().
			Width(distributedWidth).
			Align(lipgloss.Center).
			Render(truncatedKey)
		keyBindingLine += cell
	}

	keyBindingPanelWidth := lipgloss.Width(keyBindingLine)
	keyBindingLine = truncateString(keyBindingLine, keyBindingPanelWidth)

	return bottomBarStyle.Width(width).Height(mainPageKeyBindingLayoutPanelHeight).Render(keyBindingLine)
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
	m.CurrentSelectedFileDiffViewport.Update(nil)
}
