package tui

import (
	// "fmt"
	// "gitti/i18n"
	"gitti/settings"

	"golang.org/x/text/width"
)

// variables for indicating which panel/components/container or whatever the hell you wanna call it that the user is currently landed or selected, so that they can do precious action related to the part of whatever the hell you wanna call it
const (
	NoneSelected = "0"

	LocalBranchComponent   = "C1"
	ModifiedFilesComponent = "C2"
	FileDiffComponent      = "C3"
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

}

// truncateString trims string s to fit within given display width,
// accounting for wide CJK characters, and appends "…" if truncated.
func truncateString(s string, maxWidth int) string {
	displayWidth := 0
	runes := []rune(s)
	var result []rune

	for _, r := range runes {
		prop := width.LookupRune(r)
		k := 1
		if prop.Kind() == width.EastAsianWide || prop.Kind() == width.EastAsianFullwidth {
			k = 2
		}

		if displayWidth+k > maxWidth {
			break
		}

		displayWidth += k
		result = append(result, r)
	}

	if len(result) < len(runes) {
		// Add ellipsis and a space padding if possible
		if maxWidth-displayWidth >= 0 {
			result = append(result[:len(result)-2], '…')
		}
	}

	return string(result)
}
