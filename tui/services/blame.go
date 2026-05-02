package services

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
	blamePopUp "github.com/gohyuhan/gitti/tui/popup/blame"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

const InfoLineSeparator = "  |  "

// ------------------------------------
//
//	Populate blame info for a file and display in blame popup viewport
//
// ------------------------------------
func GetFileGitBlameInfoService(m *types.GittiModel, filePath string) {
	largestBlameInfoLineLength, largestLineLength, lineBlameInfoArray := m.GitOperations.GitBlame.GetFileGitBlameInfo(filePath)

	rowSeperator := strings.Repeat("-", largestLineLength+ansi.StringWidth(InfoLineSeparator))
	var viewportContent strings.Builder
	var previousCommitHash string
	popUp, ok := m.PopUpModel.(*blamePopUp.BlamePoUpModel)
	if !ok {
		return
	}
	for _, lineBlameInfo := range lineBlameInfoArray {
		blameInfoColor := style.ColorPurpleVibrant
		var blameInfo string
		if previousCommitHash != lineBlameInfo.CommitHash {
			previousCommitHash = lineBlameInfo.CommitHash
			isUncommitedChange := strings.ReplaceAll(lineBlameInfo.CommitHash, "0", "") == ""
			if isUncommitedChange {
				blameInfoColor = style.ColorYellowSoft
			}
			viewportContent.WriteString(style.NewStyle.Foreground(style.ColorPurpleSoft).Faint(true).Render(rowSeperator))
			viewportContent.WriteRune('\n')
			blameInfo = lineBlameInfo.ConsolidatedBlameInfo
		} else {
			blameInfo = ""
		}
		blameInfoLength := ansi.StringWidth(blameInfo)
		if blameInfoLength < largestBlameInfoLineLength {
			blameInfo += strings.Repeat(" ", largestBlameInfoLineLength-blameInfoLength)
		}
		viewportContent.WriteString(style.NewStyle.Foreground(blameInfoColor).Render(blameInfo))
		viewportContent.WriteString(style.NewStyle.Foreground(style.ColorPurpleSoft).Faint(true).Render(InfoLineSeparator))
		viewportContent.WriteString(lineBlameInfo.ComittedLine)
		viewportContent.WriteRune('\n')
	}

	popUp.BlameViewport.SetContent(viewportContent.String())
}
