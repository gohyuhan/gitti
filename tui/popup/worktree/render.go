package worktree

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Render the add-new-worktree input popup: a titled worktree-name input followed
//	by a titled optional branch-name input, both resized to the current popup width.
//
// ------------------------------------
func RenderWorktreeAddNewWorktreePopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*WorktreeAddNewWorktreePopUpModel)
	if ok {
		popUpWidth := min(constant.MaxWorktreeAddNewWorktreePopUpWidth, int(float64(m.Width)*0.8))
		popUp.WorktreeNameTextInput.SetWidth(popUpWidth - 6)
		popUp.WorktreeBranchNameTextInput.SetWidth(popUpWidth - 6)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			style.TitleStyle.Render(i18n.LANGUAGEMAPPING.AddNewWorktreeTitle),
			popUp.WorktreeNameTextInput.View(),
			style.TitleStyle.Render(i18n.LANGUAGEMAPPING.NewWorktreeBranchTitle),
			popUp.WorktreeBranchNameTextInput.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the add-new-worktree output popup: shows the command output viewport
//	with a spinner and processing label while running, and colors the viewport
//	border red on error or green on success once finished.
//
// ------------------------------------
func RenderWorktreeAddNewWorktreeOutputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*WorktreeAddNewWorktreeOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxWorktreeAddNewWorktreeOutputPopUpWidth, int(float64(m.Width)*0.8))

		outputViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpWorktreeAddNewWorktreeOutputViewportHeight + 2)
		if popUp.HasError.Load() {
			outputViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			outputViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.AddNewWorktreeOutputViewport.SetWidth(popUpWidth - 4)
		popUp.AddNewWorktreeOutputViewport.SetYOffset(popUp.AddNewWorktreeOutputViewport.YOffset())
		outputViewPort := outputViewPortStyle.Render(popUp.AddNewWorktreeOutputViewport.View())

		var content string
		if popUp.IsProcessing.Load() {
			processingText := popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.AddingNewWorktree
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				i18n.LANGUAGEMAPPING.NewWorktreeTitle,
				processingText,
				outputViewPort,
			)

		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				i18n.LANGUAGEMAPPING.NewWorktreeTitle,
				outputViewPort,
			)
		}
		return style.PopUpBorderStyle.Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the lock-reason input popup: a title showing the target worktree path
//	followed by the optional lock-reason input, resized to the current popup width.
//
// ------------------------------------
func RenderWorktreeLockReasonInputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*WorktreeLockReasonInputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxWorktreeLockReasonInputPopUpWidth, int(float64(m.Width)*0.8))
		popUp.WorktreeLockReasonTextInput.SetWidth(popUpWidth - 6)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			style.TitleStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.WorktreeLockReasonTitle, popUp.WorktreePath)),
			popUp.WorktreeLockReasonTextInput.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the remove-worktree confirmation popup: a confirmation prompt followed
//	by the target worktree path, resized to the current popup width.
//
// ------------------------------------
func RenderWorktreeRemoveWorktreeConfirmationPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*WorktreeRemoveWorktreeConfirmationPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxWorktreeRemoveWorktreeConfirmationPopUpWdith, int(float64(m.Width)*0.8))
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			style.TitleStyle.Render(i18n.LANGUAGEMAPPING.WorktreeRemoveConfirmation),
			style.NewStyle.Foreground(style.ColorPurpleVibrant).Render(fmt.Sprintf("[%s]", popUp.WorktreePath)),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
