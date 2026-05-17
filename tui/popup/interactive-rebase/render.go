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
//	Render the interactive rebase operation type selection popup, showing a
//	titled list of options (fixup/squash, reword, drop) for the user to choose.
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

// *************************************************************************************
//                        INTERACTIVE REBASE - FIXUP / SQUASH
// *************************************************************************************

// ------------------------------------
//
//	Render the split-pane fixup/squash commit selection popup. The left pane
//	shows the commit list and the right pane shows a preview of the selected
//	commits. The active pane is highlighted with the selected border style.
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
			style.NewStyle.Faint(true).Render(i18n.LANGUAGEMAPPING.InteractiveRebaseFixupSquashWarning),
		)

		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the fixup/squash commit message editing popup, showing a message text
//	input and a description textarea pre-filled with the combined commit message
//	from all selected commits.
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
//	Render the fixup/squash rebase output popup, showing a scrollable viewport
//	of rebase command output. Displays a spinner while the rebase is in progress,
//	and colors the viewport border red on error or green on success.
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

// *************************************************************************************
//
//	INTERACTIVE REBASE - REWORD
//
// *************************************************************************************

// ------------------------------------
//
//	Render the reword commit selection popup, showing a titled list of commits.
//	If a selection validation error is present, it is rendered below the list;
//	otherwise the reword warning hint is shown.
//
// ------------------------------------
func RenderInteractiveRebaseRewordSelectionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseRewordSelectionPopUpModel)
	if ok {
		popUpWidth := int(float64(m.Width) * 0.9)
		listWidth := popUpWidth - 2
		height := int(float64(m.Height)*0.8) - 2

		popUp.CommitList.SetWidth(listWidth)
		popUp.CommitList.SetHeight(height)

		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.InteractiveRebaseReword)

		var content string
		if popUp.SelectionError != nil {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				popUp.CommitList.View(),
				style.NewStyle.Foreground(style.ColorError).Render(popUp.SelectionError.Error()),
			)
		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				popUp.CommitList.View(),
				style.NewStyle.Faint(true).Render(i18n.LANGUAGEMAPPING.InteractiveRebaseRewordWarning),
			)
		}

		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the reword commit message editing popup, showing a message text input
//	and a description textarea pre-filled with the selected commit's existing message.
//
// ------------------------------------
func RenderInteractiveRebaseRewordCommitPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseRewordCommitPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxInteractiveRebaseRewordCommitPopUpWidth, int(float64(m.Width)*0.8))
		popUp.MessageTextInput.SetWidth(popUpWidth - 6)
		popUp.DescriptionTextAreaInput.SetWidth(popUpWidth - 6)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			style.TitleStyle.Render(i18n.LANGUAGEMAPPING.InteractiveRebaseRewordCommitMessageTitle),
			popUp.MessageTextInput.View(),
			style.TitleStyle.Render(i18n.LANGUAGEMAPPING.InteractiveRebaseRewordCommitDescriptionTitle),
			popUp.DescriptionTextAreaInput.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the reword rebase output popup, showing a scrollable viewport of
//	rebase command output. Displays a spinner while the rebase is in progress,
//	and colors the viewport border red on error or green on success.
//
// ------------------------------------
func RenderInteractiveRebaseRewordOutputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseRewordOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxInteractiveRebaseRewordOutputPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.InteractiveRebaseRewordOutputPopUpTitle)
		logViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpInteractiveRebaseRewordOutputviewportHeight + 2)
		if popUp.HasError.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.RewordOutputViewport.SetWidth(popUpWidth - 4)
		popUp.RewordOutputViewport.SetYOffset(popUp.RewordOutputViewport.YOffset())
		logViewPort := logViewPortStyle.Render(popUp.RewordOutputViewport.View())

		var content string
		// Show spinner above viewport when processing
		if popUp.IsProcessing.Load() {
			processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.InteractiveRebaseRewording)
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

// *************************************************************************************
//
//	INTERACTIVE REBASE - DROP
//
// *************************************************************************************

// ------------------------------------
//
//	Render the drop commit selection popup with a scrollable commit list,
//	selection state, and inline validation error or warning message
//
// ------------------------------------
func RenderInteractiveRebaseDropSelectionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseDropSelectionPopUpModel)
	if ok {
		popUpWidth := int(float64(m.Width) * 0.9)
		listWidth := popUpWidth - 2
		height := int(float64(m.Height)*0.8) - 2

		popUp.CommitList.SetWidth(listWidth)
		popUp.CommitList.SetHeight(height)

		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.InteractiveRebaseDrop)

		var content string
		if popUp.SelectionError != nil {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				popUp.CommitList.View(),
				style.NewStyle.Foreground(style.ColorError).Render(popUp.SelectionError.Error()),
			)
		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				popUp.CommitList.View(),
				style.NewStyle.Faint(true).Render(i18n.LANGUAGEMAPPING.InteractiveRebaseDropWarning),
			)
		}

		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the drop output popup with a log viewport, spinner during processing,
//	and border color reflecting success/error state
//
// ------------------------------------
func RenderInteractiveRebaseDropOutputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*InteractiveRebaseDropOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxInteractiveRebaseDropOutputPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.InteractiveRebaseDropOutputPopUpTitle)
		logViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpInteractiveRebaseDropOutputviewportHeight + 2)
		if popUp.HasError.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.DropOutputViewport.SetWidth(popUpWidth - 4)
		popUp.DropOutputViewport.SetYOffset(popUp.DropOutputViewport.YOffset())
		logViewPort := logViewPortStyle.Render(popUp.DropOutputViewport.View())

		var content string
		// Show spinner above viewport when processing
		if popUp.IsProcessing.Load() {
			processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.InteractiveRebaseDropping)
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
