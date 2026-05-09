package interactiverebase

import (
	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	choose interactive rebase option
//
// ------------------------------------
func RenderInteractiveRebaseOptionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseOptionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxInteractiveRebaseOptionPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.ChooseInteractiveRebaseOption)
		popUp.InteractiveRebaseOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.InteractiveRebaseOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the split-pane fixup/squash commit selection popup;
//	the active pane (list or detail viewport) is highlighted with the selected border style
//
// ------------------------------------
func RenderInteractiveRebaseFixupSquashSelectionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseFixupSquashSelectionPopUpModel)
	if ok {
		popUpWidth := int(float64(m.Width) * 0.9)
		innerWidth := popUpWidth - 2
		listWidth := int(float64(innerWidth) * 0.65)
		vpWidth := innerWidth - listWidth
		height := int(float64(m.Height)*0.8) - 2

		popUp.CommitList.SetWidth(listWidth - 2)
		popUp.CommitFixupSquashViewport.SetWidth(vpWidth - 2)
		popUp.CommitList.SetHeight(height)
		popUp.CommitFixupSquashViewport.SetHeight(height)

		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.InteractiveRebaseFixupSquash)

		var commitSelectionView string
		var commitViewportView string

		if popUp.IsCommitListSelected.Load() && !popUp.IsCommitFixupSquashViewportSelected.Load() {
			commitSelectionView = style.SelectedBorderStyle.Width(listWidth).Render(popUp.CommitList.View())
			commitViewportView = style.PanelBorderStyle.Width(vpWidth).Render(popUp.CommitFixupSquashViewport.View())
		} else if !popUp.IsCommitListSelected.Load() && popUp.IsCommitFixupSquashViewportSelected.Load() {
			commitSelectionView = style.PanelBorderStyle.Width(listWidth).Render(popUp.CommitList.View())
			commitViewportView = style.SelectedBorderStyle.Width(vpWidth).Render(popUp.CommitFixupSquashViewport.View())
		}

		innerContent := lipgloss.JoinHorizontal(
			lipgloss.Top,
			commitSelectionView,
			commitViewportView,
		)

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			innerContent,
		)

		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
