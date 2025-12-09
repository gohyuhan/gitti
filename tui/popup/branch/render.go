package branch

import (
	"fmt"

	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"

	"charm.land/lipgloss/v2"
)

// ------------------------------------
//
//	For Creating New Git branch
//
// ------------------------------------
// pop up that confirm the option for creating a new branch, just create or create and move everything to the new branch
func RenderChooseNewBranchTypePopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChooseNewBranchTypeOptionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChooseNewBranchTypePopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.ChooseNewBranchTypeTitle)
		popUp.NewBranchTypeOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.NewBranchTypeOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// to prompt user for new branch name and then proceed to trigger the creation of branch and optionally move changes
func RenderCreateNewBranchPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*CreateNewBranchPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxCreateNewBranchPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.CreateNewBranchTitle)
		popUp.NewBranchNameInput.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.NewBranchNameInput.View(),
		)
		modifiedBranchName, isValid := api.IsBranchNameValid(popUp.NewBranchNameInput.Value())
		if !isValid {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				popUp.NewBranchNameInput.View(),
				style.BranchInvalidWarningStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.NewBranchInvalidWarning, modifiedBranchName)),
			)
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	For Switching Git branch
//
// ------------------------------------
// pop up that confirm the option for switching a branch, just switch or switch to the branch while bringing all the changes
func RenderChooseSwitchBranchTypePopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChooseSwitchBranchTypePopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChooseSwitchBranchTypePopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.ChooseSwitchBranchTypeTitle, popUp.BranchName))
		popUp.SwitchTypeOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.SwitchTypeOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// pop up to render the output of the switch branch operation
// because we allow switching with bring changes over, there is conflict possiblities there fore we need to show the output
// so that the user is aware of it
func RenderSwitchBranchOutputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*SwitchBranchOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxSwitchBranchOutputPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.SwitchBranchSwitchingToPopUpTitle, popUp.BranchName))
		logViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpSwitchBranchOutputViewPortHeight + 2)
		if popUp.HasError.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.SwitchBranchOutputViewport.SetWidth(popUpWidth - 4)
		popUp.SwitchBranchOutputViewport.SetYOffset(popUp.SwitchBranchOutputViewport.YOffset())
		logViewPort := logViewPortStyle.Render(popUp.SwitchBranchOutputViewport.View())

		var content string
		// Show spinner above viewport when processing
		if popUp.IsProcessing.Load() {
			processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.SwitchBranchPopUpSwitchProcessing)
			if popUp.SwitchType == git.SWITCHBRANCHWITHCHANGES {
				processingText = style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.SwitchBranchPopUpSwitchWithChangesProcessing)
			}
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

// ------------------------------------
//
//	For Git delete branch confirmation prompt
//
// ------------------------------------
func RenderGitDeleteBranchConfirmPromptPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitDeleteBranchConfirmPromptPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitDeleteBranchConfirmPromptPopUpWidth, int(float64(m.Width)*0.8))
		deleteConfirmationPrompt := fmt.Sprintf(i18n.LANGUAGEMAPPING.GitDeleteBranchComfirmPrompt, style.NewStyle.Foreground(style.ColorYellowWarm).Render(popUp.BranchName))

		return style.PopUpBorderStyle.Width(popUpWidth).Render(deleteConfirmationPrompt)
	}

	return ""
}

// ------------------------------------
//
//	For Git delete branch output result
//
// ------------------------------------
func RenderGitDeleteBranchOutputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*GitDeleteBranchOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxGitDeleteBranchOutputPopUpWidth, int(float64(m.Width)*0.8))

		outputViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpGitDeleteBranchOutputViewportHeight + 2)
		if popUp.HasError.Load() {
			outputViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			outputViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.BranchDeleteOutputViewport.SetWidth(popUpWidth - 4)
		popUp.BranchDeleteOutputViewport.SetYOffset(popUp.BranchDeleteOutputViewport.YOffset())
		outputViewPort := outputViewPortStyle.Render(popUp.BranchDeleteOutputViewport.View())

		var content string
		if popUp.IsProcessing.Load() {
			processingText := popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.DeletingBranch
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				i18n.LANGUAGEMAPPING.GitDeleteBranchTitle,
				processingText,
				outputViewPort,
			)

		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				i18n.LANGUAGEMAPPING.GitDeleteBranchTitle,
				outputViewPort,
			)
		}
		return style.PopUpBorderStyle.Render(content)
	}
	return ""
}
