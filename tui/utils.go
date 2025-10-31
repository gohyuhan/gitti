package tui

import (
	"gitti/settings"
	"unicode/utf8"
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

func tuiWindowSizing(m *GittiModel) {
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

func processGitUpdate(m *GittiModel) {
	initBranchList(m)
	initModifiedFilesList(m)
	renderModifiedFilesDiffViewPort(m)
	return
}

// to prevent wrapping and overflow doesn't seems to be supported natively in bubbletea yet
func truncateString(s string, width int) string {
	if utf8.RuneCountInString(s) <= width {
		return s
	}
	runes := []rune(s)
	if width > 1 {
		return string(runes[:width-3]) + "…  "
	}
	return string(runes[:width])
}
