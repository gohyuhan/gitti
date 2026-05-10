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
//	Render the new-branch type selection popup, showing a list of creation
//	options: create only, create-and-switch, create-based-on-remote-input, and
//	create-based-on-remote-selection.
//
// ------------------------------------
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

// ------------------------------------
//
//	Render the new-branch name input popup with a focused text input. Appends an
//	inline warning if the current input value is not a valid branch name.
//
// ------------------------------------
func RenderCreateNewBranchPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*CreateNewBranchPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxCreateNewBranchPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.CreateNewBranchTitle)
		popUp.NewBranchNameInput.SetWidth(popUpWidth - 6)
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
//	Render the switch-branch type selection popup, showing a titled list of
//	switch options (switch clean, switch-with-changes) for the target branch.
//
// ------------------------------------
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

// ------------------------------------
//
//	Render the switch-branch output popup, showing a scrollable viewport of
//	command output. When switching-with-changes, conflicts can occur, so output
//	is always displayed. Shows a spinner while processing, and colors the viewport
//	border red on error or green on success.
//
// ------------------------------------
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
//	Render the branch-deletion confirmation popup, showing the localized prompt
//	with the target branch name highlighted in yellow.
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
//	Render the branch-deletion output popup, showing command output in a
//	viewport with a spinner while processing. Colors the border red on error
//	or green on success.
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

// ------------------------------------
//
//	Render the create-branch-from-remote input popup, showing the remote origin
//	and a focused branch-name text input. Appends an inline warning if the
//	current input value is not a valid branch name.
//
// ------------------------------------
func RenderCreateBranchBasedOnRemotePopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*CreateBranchBasedOnRemotePopUpModel)
	if ok {
		popUpWidth := min(constant.MaxCreateBranchBasedOnRemotePopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.CreateNewBranchBasedOnRemoteUserInputTitle)
		popUp.RemoteBranchNameInput.SetWidth(popUpWidth - 6)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			i18n.LANGUAGEMAPPING.RemoteOriginTitle,
			fmt.Sprintf(" %s", popUp.RemoteOrigin),
			i18n.LANGUAGEMAPPING.EnterRemoteBranchTitle,
			popUp.RemoteBranchNameInput.View(),
		)
		modifiedBranchName, isValid := api.IsBranchNameValid(popUp.RemoteBranchNameInput.Value())
		if !isValid {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				"",
				i18n.LANGUAGEMAPPING.RemoteOriginTitle,
				fmt.Sprintf(" %s", popUp.RemoteOrigin),
				i18n.LANGUAGEMAPPING.EnterRemoteBranchTitle,
				popUp.RemoteBranchNameInput.View(),
				style.BranchInvalidWarningStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.NewBranchInvalidWarning, modifiedBranchName)),
			)
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the create-branch-from-remote output popup, showing command output
//	in a viewport with a spinner while processing. Colors the border red on
//	error or green on success.
//
// ------------------------------------
func RenderCreateBranchBasedOnRemoteOutputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*CreateBranchBasedOnRemoteOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxCreateBranchBasedOnRemoteOutputPopUpWidth, int(float64(m.Width)*0.8))

		outputViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpCreateBranchBasedOnRemoteOutputViewportHeight + 2)
		if popUp.HasError.Load() {
			outputViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			outputViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.CreateBranchBasedOnRemoteOutputViewport.SetWidth(popUpWidth - 4)
		popUp.CreateBranchBasedOnRemoteOutputViewport.SetYOffset(popUp.CreateBranchBasedOnRemoteOutputViewport.YOffset())
		outputViewPort := outputViewPortStyle.Render(popUp.CreateBranchBasedOnRemoteOutputViewport.View())

		var content string
		if popUp.IsProcessing.Load() {
			processingText := popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.CreatingNewBranchBasedOnRemoteProcessing
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				i18n.LANGUAGEMAPPING.CreatingNewBranchBasedOnRemoteTitle,
				processingText,
				outputViewPort,
			)

		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				i18n.LANGUAGEMAPPING.CreatingNewBranchBasedOnRemoteTitle,
				outputViewPort,
			)
		}
		return style.PopUpBorderStyle.Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the remote-branch selection popup, showing a titled list of all
//	remote branches the user can choose to track locally.
//
// ------------------------------------
func RenderChooseRemoteBranchOptionPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChooseRemoteBranchOptionPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChooseRemoteBranchOptionPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.ChooseRemoteBranchOptionTitle)
		popUp.RemoteBranchOptionList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.RemoteBranchOptionList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the branch-selection popup for git merge. Shows two bordered panels:
//	the top panel lists available (unselected) branches, the bottom panel lists
//	already-selected branches. The active panel is highlighted with a selected
//	border; the inactive panel uses the default panel border.
//
// ------------------------------------
func RenderChooseBranchOptionForMergePopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChooseBranchOptionForMergePopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChooseBranchOptionForMergePopUpWidth, int(float64(m.Width)*0.8))
		branchOptionTitle := style.TitleStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.ChooseBranchOptionForGitMergeTitle, style.NewStyle.Foreground(style.ColorYellowWarm).Render(m.CheckOutBranch)))
		selectedBranchOptionTitle := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.SelectedBranchOptionForGitMergeTitle)
		popUp.BranchOptionList.SetWidth(popUpWidth - 4)
		popUp.SelectedBranchList.SetWidth(popUpWidth - 4)

		branchOptionListView := popUp.BranchOptionList.View()
		selectedBranchList := popUp.SelectedBranchList.View()

		borderWidth := popUpWidth - 2

		if popUp.BranchOptionSectionSelected.Load() {
			branchOptionListView = style.SelectedBorderStyle.Width(borderWidth).Render(popUp.BranchOptionList.View())
			selectedBranchList = style.PanelBorderStyle.Width(borderWidth).Render(popUp.SelectedBranchList.View())
		} else if popUp.SelectedBranchSectionSelected.Load() {
			branchOptionListView = style.PanelBorderStyle.Width(borderWidth).Render(popUp.BranchOptionList.View())
			selectedBranchList = style.SelectedBorderStyle.Width(borderWidth).Render(popUp.SelectedBranchList.View())
		}

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			branchOptionTitle,
			branchOptionListView,
			selectedBranchOptionTitle,
			selectedBranchList,
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	Render the git merge output popup, showing a scrollable viewport of merge
//	command output. Displays a spinner while the merge is in progress, and colors
//	the viewport border red on error or green on success.
//
// ------------------------------------
func RenderBranchMergeOutputPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*BranchMergeOutputPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxBranchMergeOutputPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.GitMergeOutputTitle)
		logViewPortStyle := style.PanelBorderStyle.
			Width(popUpWidth - 2).
			Height(constant.PopUpBranchMergeOutputViewportHeight + 2)
		if popUp.HasError.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorError)
		} else if popUp.ProcessSuccess.Load() {
			logViewPortStyle = style.PanelBorderStyle.
				BorderForeground(style.ColorGreenSoft)
		}
		popUp.BranchMergeOutputViewport.SetWidth(popUpWidth - 4)
		popUp.BranchMergeOutputViewport.SetYOffset(popUp.BranchMergeOutputViewport.YOffset())
		logViewPort := logViewPortStyle.Render(popUp.BranchMergeOutputViewport.View())

		var content string
		// Show spinner above viewport when processing
		if popUp.IsProcessing.Load() {
			processingText := style.SpinnerStyle.Render(popUp.Spinner.View() + " " + i18n.LANGUAGEMAPPING.BranchMerging)
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
