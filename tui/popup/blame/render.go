package blame

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ----------------------------------
//
//	Render the blame popup: file-selection list with filter input, or blame viewport for chosen file
//
// ----------------------------------
func RenderBlamePopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*BlamePoUpModel)
	if ok {
		width := int(float64(m.Width) * 0.9)
		height := int(float64(m.Height) * 0.9)
		popUpComponentWidth := width - 4
		popUpComponentHeight := height - 2
		popUp.CurrentGitTrackedFilesPathList.SetWidth(popUpComponentWidth)
		popUp.CurrentGitTrackedFilesPathList.SetHeight(popUpComponentHeight)
		popUp.BlameViewport.SetWidth(popUpComponentWidth)
		popUp.BlameViewport.SetHeight(popUpComponentHeight)

		var content string

		if !popUp.HasFilePathChosen && !popUp.ShowingBlameInfo {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				i18n.LANGUAGEMAPPING.GitTrackedFileTitle,
				popUp.CurrentGitTrackedFilesPathList.View(),
				popUp.FilterInput.View(),
			)
		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				fmt.Sprintf(i18n.LANGUAGEMAPPING.BlameViewportTitle, popUp.SelectedFilePath),
				popUp.BlameViewport.View(),
			)
		}

		return style.PopUpBorderStyle.Width(width).Height(height).Render(content)
	}
	return ""
}
