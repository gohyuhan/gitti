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
//	Render the interactive rebase option selection popup (fixup/squash, reword, drop)
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

		if popUp.IsCommitListSelected && !popUp.IsCommitFixupSquashViewportSelected {
			commitSelectionView = style.SelectedBorderStyle.Width(listWidth).Render(popUp.CommitList.View())
			commitViewportView = style.PanelBorderStyle.Width(vpWidth).Render(popUp.CommitFixupSquashViewport.View())
		} else if !popUp.IsCommitListSelected && popUp.IsCommitFixupSquashViewportSelected {
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
			style.NewStyle.Faint(true).Render(i18n.LANGUAGEMAPPING.InteractiveRebaseFixupWarning),
		)

		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the commit message and description input popup for the fixup/squash rebase operation
//
// ------------------------------------
func RenderInteractiveRebaseFixupSquashCommitPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseFixupSquashCommitPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxInteractiveRebaseFixupSquashCommitPopUpWidth, int(float64(m.Width)*0.8))
		popUp.MessageTextInput.SetWidth(popUpWidth - 6)
		popUp.DescriptionTextAreaInput.SetWidth(popUpWidth - 6)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			style.TitleStyle.Render(i18n.LANGUAGEMAPPING.InteractiveRebaseFixupSquashCommitMessageTitle),
			popUp.MessageTextInput.View(),
			style.TitleStyle.Render(i18n.LANGUAGEMAPPING.InteractiveRebaseFixupSquashCommitDescriptionTitle),
			popUp.DescriptionTextAreaInput.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the fixup/squash rebase output popup; shows a spinner while processing,
//	and changes the viewport border color to reflect error or success state
//
// ------------------------------------
func RenderInteractiveRebaseFixupSquashOutputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseFixupSquashOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxInteractiveRebaseFixupSquashOutputPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.InteractiveRebaseFixupSquashOutputPopUpTitle)
		logViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpInteractiveRebaseFixupSquashOutputviewportHeight + 2)
		if popUp.HasError.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.FixupSquashOutputViewport.SetWidth(popUpWidth - 4)
		popUp.FixupSquashOutputViewport.SetYOffset(popUp.FixupSquashOutputViewport.YOffset())
		logViewPort := logViewPortStyle.Render(popUp.FixupSquashOutputViewport.View())

		var content string
		// Show spinner above viewport when processing
		if popUp.IsProcessing.Load() {
			processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.InteractiveRebaseFixupSquashing)
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				"",
				processingText,
				logViewPort,
			)
		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				logViewPort,
			)
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
