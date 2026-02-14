package layout

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	gitticonst "github.com/gohyuhan/gitti/constant"
	"github.com/gohyuhan/gitti/i18n"
	branchComponent "github.com/gohyuhan/gitti/tui/component/branch"
	filesComponent "github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// -----------------------------------------------------------------------------
//
//	Functions that help construct the view
//
// -----------------------------------------------------------------------------
// render the Gitti Status Panel
func renderGitStatusComponentPanel(m *types.GittiModel) string {
	borderStyle := style.PanelBorderStyle
	if m.CurrentSelectedComponent == constant.GitStatusComponentPanel {
		borderStyle = style.SelectedBorderStyle
	}
	if m.CurrentGitRepoStatus == "" {
		var remoteSyncStateLineString string
		additionalWidth := 0

		if m.RemoteSyncLocalState == "" || m.RemoteSyncRemoteState == "" {
			remoteSyncStateLineString = style.ErrorStyle.Render("\uf00d")
			additionalWidth += 1
		} else {
			local := style.LocalStatusStyle.Render(fmt.Sprintf("%s↑", m.RemoteSyncLocalState))
			remote := style.RemoteStatusStyle.Render(fmt.Sprintf("%s↓", m.RemoteSyncRemoteState))

			remoteSyncStateLineString = local + " " + remote
			additionalWidth += 3 + lipgloss.Width(m.RemoteSyncLocalState) + lipgloss.Width(m.RemoteSyncRemoteState)
		}

		trackedUpStreamOrBranchName := m.CheckOutBranch
		if m.BranchUpStream != "" {
			trackedUpStreamOrBranchName = m.BranchUpStream
		}

		repoTrackBranchName := fmt.Sprintf(" %s -> %s %s", m.RepoName, m.TrackedUpstreamOrBranchIcon, trackedUpStreamOrBranchName)

		// the max width is the window width - padding - the length of RemoteSyncStateLineString
		repoTrackBranchName = utils.TruncateString(repoTrackBranchName, m.WindowLeftPanelWidth-constant.ListItemOrTitleWidthPad-additionalWidth)

		return borderStyle.
			Width(m.WindowLeftPanelWidth).
			Height(1).
			Render(fmt.Sprintf("%s%s", remoteSyncStateLineString, repoTrackBranchName))
	} else {
		var gitStateInProgress string

		gitStateInProgress = utils.TruncateString(fmt.Sprintf(i18n.LANGUAGEMAPPING.GitCertainStateStillInProgress, m.CurrentGitRepoStatus), m.WindowLeftPanelWidth-constant.ListItemOrTitleWidthPad-2)

		return borderStyle.
			Width(m.WindowLeftPanelWidth).
			Height(1).
			Render(fmt.Sprintf("%s %s", lipgloss.NewStyle().Foreground(style.ColorError).Render("!"), gitStateInProgress))
	}
}

// Render the Local Branches panel
func renderLocalBranchesOrTagOrRemoteComponentPanel(width int, height int, m *types.GittiModel) string {
	borderStyle := style.PanelBorderStyle
	if m.CurrentSelectedComponent == constant.LocalBranchOrTagOrRemoteComponentPanel {
		borderStyle = style.SelectedBorderStyle
	}
	var content string
	switch m.CurrentLocalBranchOrTagComponentShowing {
	case constant.SHOW_LOCAL_BRANCH:
		content = m.CurrentRepoBranchesInfoList.View()
	case constant.SHOW_TAG:
		content = m.CurrentRepoTagInfoList.View()
	case constant.SHOW_REMOTE:
		content = m.CurrentRepoRemoteInfoList.View()
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(content)
}

// Render the Changed Files panel
func renderModifiedFilesComponentPanel(width int, height int, m *types.GittiModel) string {
	borderStyle := style.PanelBorderStyle
	if m.CurrentSelectedComponent == constant.ModifiedFilesComponentPanel {
		borderStyle = style.SelectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(m.CurrentRepoModifiedFilesInfoList.View())
}

// Render the Changed Files panel
func renderCommitLogComponentPanel(width int, height int, m *types.GittiModel) string {
	borderStyle := style.PanelBorderStyle
	if m.CurrentSelectedComponent == constant.CommitLogComponentPanel {
		borderStyle = style.SelectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(m.CurrentRepoCommitLogInfoList.View())
}

// Render the detail component part at the right of the window,
// however the content within it will be dynamic based on the current selected component
// ------------------------------------
//
//			For Render Detail Component Panel
//		  * this will render the detail component panel
//	   * it will handle the layout for both line editing mode and normal mode
//	   * it will also handle the split view for line editing mode (one for staged, one for unstaged)
//
// ------------------------------------
func renderDetailComponentPanel(width int, height int, m *types.GittiModel) string {
	detailComponentBorderStyle := style.PanelBorderStyle
	detailComponentTwoBorderStyle := style.PanelBorderStyle

	// determine the border style based on the current selected component
	switch m.CurrentSelectedComponent {
	case constant.DetailComponentPanel:
		detailComponentBorderStyle = style.SelectedBorderStyle
	case constant.DetailComponentPanelTwo:
		detailComponentTwoBorderStyle = style.SelectedBorderStyle
	}

	var content string
	ogHeight := height
	ogWidth := width

	// if it is in line editing mode, we need to minus 3 for the title
	if m.IsLineEditingState.Load() {
		ogHeight = height - 3
	}

	if m.IsLineEditingState.Load() {
		// we first define the border height
		borderHeight := 2
		if m.ShowDetailPanelTwo.Load() {
			detailPanelHeight := int(ogHeight / 2)
			detailPanelWidth := int(ogWidth / 2)
			if m.DetailComponentPanelLayout == constant.HORIZONTAL {
				// Horizontal Split Layout for Line Editing
				// Join Cursor Viewport + Content Viewport for Panel 1
				detailPanelViewportContent := lipgloss.JoinHorizontal(
					lipgloss.Top,
					style.NewStyle.Width(3).Height(ogHeight-borderHeight).Render(m.LineEditingIndexCursorViewport.View()),
					style.NewStyle.Width(detailPanelWidth-3).Height(ogHeight-borderHeight).Render(m.DetailPanelViewport.View()),
				)
				// Join Cursor Viewport + Content Viewport for Panel 2
				detailPanelTwoViewportContent := lipgloss.JoinHorizontal(
					lipgloss.Top,
					style.NewStyle.Width(3).Height(ogHeight-borderHeight).Render(m.LineEditingIndexCursorTwoViewport.View()),
					style.NewStyle.Width(detailPanelWidth-3).Height(ogHeight-borderHeight).Render(m.DetailPanelTwoViewport.View()),
				)

				// Combine both panels horizontally
				content = lipgloss.JoinHorizontal(
					lipgloss.Top,
					detailComponentBorderStyle.Width(detailPanelWidth).Height(ogHeight).Render(detailPanelViewportContent),
					detailComponentTwoBorderStyle.Width(width-detailPanelWidth).Height(ogHeight).Render(detailPanelTwoViewportContent),
				)
			} else {
				// Vertical Split Layout for Line Editing
				// Join Cursor + Content for Panel 1 (Top)
				detailPanelViewportContent := lipgloss.JoinHorizontal(
					lipgloss.Top,
					style.NewStyle.Width(3).Height(detailPanelHeight-borderHeight).Render(m.LineEditingIndexCursorViewport.View()),
					style.NewStyle.Width(ogWidth-3).Height(detailPanelHeight-borderHeight).Render(m.DetailPanelViewport.View()),
				)
				// Join Cursor + Content for Panel 2 (Bottom)
				detailPanelTwoViewportContent := lipgloss.JoinHorizontal(
					lipgloss.Top,
					style.NewStyle.Width(3).Height(ogHeight-detailPanelHeight-borderHeight).Render(m.LineEditingIndexCursorTwoViewport.View()),
					style.NewStyle.Width(ogWidth-3).Height(ogHeight-detailPanelHeight-borderHeight).Render(m.DetailPanelTwoViewport.View()),
				)

				// Combine both panels vertically
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					detailComponentBorderStyle.Width(ogWidth).Height(detailPanelHeight).Render(detailPanelViewportContent),
					detailComponentTwoBorderStyle.Width(ogWidth).Height(ogHeight-detailPanelHeight).Render(detailPanelTwoViewportContent),
				)
			}
		} else {
			// Single Panel Layout for Line Editing
			detailPanelViewportContent := lipgloss.JoinHorizontal(
				lipgloss.Top,
				style.NewStyle.Width(3).Height(ogHeight-borderHeight).Render(m.LineEditingIndexCursorViewport.View()),
				style.NewStyle.Width(ogWidth-3).Height(ogHeight-borderHeight).Render(m.DetailPanelViewport.View()),
			)
			content = lipgloss.JoinHorizontal(
				lipgloss.Top,
				detailComponentBorderStyle.Width(ogWidth).Height(ogHeight).Render(detailPanelViewportContent),
			)
		}
	} else {
		// Standard Rendering (Not Line Editing Mode)
		if m.ShowDetailPanelTwo.Load() {
			detailPanelHeight := int(ogHeight / 2)
			detailPanelWidth := int(ogWidth / 2)
			if m.DetailComponentPanelLayout == constant.HORIZONTAL {
				content = lipgloss.JoinHorizontal(
					lipgloss.Top,
					detailComponentBorderStyle.Width(detailPanelWidth).Height(ogHeight).Render(m.DetailPanelViewport.View()),
					detailComponentTwoBorderStyle.Width(width-detailPanelWidth).Height(ogHeight).Render(m.DetailPanelTwoViewport.View()),
				)
			} else {
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					detailComponentBorderStyle.Width(ogWidth).Height(detailPanelHeight).Render(m.DetailPanelViewport.View()),
					detailComponentTwoBorderStyle.Width(ogWidth).Height(ogHeight-detailPanelHeight).Render(m.DetailPanelTwoViewport.View()),
				)
			}
		} else {
			content = lipgloss.JoinHorizontal(
				lipgloss.Top,
				detailComponentBorderStyle.Width(ogWidth).Height(ogHeight).Render(m.DetailPanelViewport.View()),
			)
		}
	}

	// Add Title Block for Line Editing Mode
	if m.IsLineEditingState.Load() {
		inLineEditingModeNotifyBlock := style.PanelBorderStyle.Width(ogWidth).Render(utils.TruncateString(i18n.LANGUAGEMAPPING.LineEditingModeTitle, ogWidth-4))
		content = lipgloss.JoinVertical(
			lipgloss.Top,
			inLineEditingModeNotifyBlock,
			content,
		)
	}

	return style.NewStyle.
		Width(width).
		Height(height).
		Render(content)
}

func renderStashComponentPanel(width int, height int, m *types.GittiModel) string {
	borderStyle := style.PanelBorderStyle
	if m.CurrentSelectedComponent == constant.StashComponentPanel {
		borderStyle = style.SelectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(m.CurrentRepoStashInfoList.View())
}

func renderLogComponentPanel(width int, height int, m *types.GittiModel) string {
	borderStyle := style.PanelBorderStyle
	if m.CurrentSelectedComponent == constant.LogComponentPanel {
		borderStyle = style.SelectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(m.CurrentLogComponentViewport.View())
}

func renderKeyBindingComponentPanel(width int, m *types.GittiModel) string {
	keys := []string{""} // to prevent a misconfiguration on key binding will not crash the program

	if m.ShowPopUp.Load() {
		//-----------------------------
		//
		// for popup keybinding render
		//
		//-----------------------------
		switch m.PopUpType {
		case constant.CommitPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForCommitPopUp
		case constant.AmendCommitPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForAmendCommitPopUp
		case constant.AddRemotePromptPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForAddRemotePromptPopUp
		case constant.GitRemotePushPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitRemotePushPopUp
		case constant.ChooseRemotePopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChooseRemotePopUp
		case constant.ChoosePushTypePopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChoosePushTypePopUp
		case constant.ChooseNewBranchTypePopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChooseNewBranchTypePopUp
		case constant.CreateNewBranchPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForCreateNewBranchPopUp
		case constant.ChooseSwitchBranchTypePopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChooseSwitchBranchTypePopUp
		case constant.SwitchBranchOutputPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChooseSwitchBranchTypePopUp
			popUp, ok := m.PopUpModel.(*branchPopUp.SwitchBranchOutputPopUpModel)
			if ok {
				if popUp.IsProcessing.Load() {
					keys = []string{"..."} // nothing can be done during switching, only force quit gitti is possible
				}
			}
		case constant.ChooseGitPullTypePopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChooseGitPullTypePopUp
		case constant.GitPullOutputPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitPullOutputPopUp
		case constant.GitStashMessagePopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitStashMessagePopUp
		case constant.KeybindingAndFeatureInstructionsPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForKeybindingAndFeatureInstructionsPopUp
		case constant.GitDiscardTypeOptionPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitDiscardTypeOptionPopUp
		case constant.GitDiscardConfirmPromptPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitDiscardConfirmPromptPopUp
		case constant.GitStashOperationOutputPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitStashOperationOutputPopUp
			popUp, ok := m.PopUpModel.(*stashPopUp.GitStashOperationOutputPopUpModel)
			if ok {
				if popUp.IsProcessing.Load() {
					keys = []string{"..."} // nothing can be done during stash operation, only force quit gitti is possible
				}
			}
		case constant.GitStashConfirmPromptPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitStashConfirmPromptPopUp
		case constant.GitDeleteBranchConfirmPromptPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitDeleteBranchConfirmPromptPopUp
		case constant.GitDeleteBranchOutputPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitDeleteBranchOutputPopUp
			popUp, ok := m.PopUpModel.(*branchPopUp.GitDeleteBranchOutputPopUpModel)
			if ok {
				if popUp.IsProcessing.Load() {
					keys = []string{"..."} // nothing can be done during stash operation, only force quit gitti is possible
				}
			}
		case constant.CreateBranchBasedOnRemotePopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForCreateBranchBasedOnRemotePopUp
		case constant.CreateBranchBasedOnRemoteOutputPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForCreateBranchBasedOnRemoteOutputPopUp
			popUp, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemoteOutputPopUpModel)
			if ok {
				if popUp.IsProcessing.Load() {
					keys = []string{"..."} // nothing can be done during stash operation, only force quit gitti is possible
				}
			}
		case constant.GitResetLatestCommitTypeOptionPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitResetLatestCommitTypeOptionPopUp
		case constant.GitResetLatestCommitConfirmPromptPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitResetLatestCommitConfirmPromptPopUp
		case constant.GitResetToSelectedCommitTypeOptionPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitResetToSelectedCommitTypeOptionPopUp
		case constant.GitResetToSelectedCommitConfirmPromptPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitResetToSelectedCommitConfirmPromptPopUp
		case constant.GitCherryPickOptionSelectionPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitCherryPickOptionSelectionPopUp
		case constant.GitCherryPickPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitCherryPickPopUp
		case constant.GitEditCherryPickPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitEditCherryPickPopUp
		case constant.GitCherryPickApplyConfirmPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitCherryPickApplyConfirmPopUp
		case constant.GitDiscardFileLineChangeConfirmPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitDiscardFileLineChangeConfirmPopUp
		case constant.CreateTagPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForCreateTagPopUp
		case constant.CreateTagConfirmationPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForCreateTagConfirmationPopUp
		case constant.ChooseDeleteTagOptionPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChooseDeleteTagOptionPopUp
		case constant.ChooseRemoteForDeleteRemoteTagPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChooseRemoteForDeleteRemoteTagPopUp
		case constant.DeleteTagOutputPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForDeleteTagOutputPopUp
		case constant.ChoosePushTagOptionPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChoosePushTagOptionPopUp
		case constant.PushTagOutputPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForPushTagOutputPopUp
		case constant.ChooseFetchTagOptionPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForChooseFetchTagOptionPopUp
		case constant.FetchTagOutputPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForFetchTagOutputPopUp
		case constant.RemoveRemoteConfirmationPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForRemoveRemoteConfirmationPopUp
		case constant.RemoteAsTrackingUpstreamConfirmationPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForRemoteAsTrackingUpstreamConfirmationPopUp
		case constant.EditRemotePromptPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForEditRemotePromptPopUp
		}
	} else {
		//-----------------------------
		//
		// for non-popup keybinding render
		//
		//-----------------------------
		switch m.CurrentSelectedComponent {
		case constant.GitStatusComponentPanel:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitStatusComponent
		case constant.LocalBranchOrTagOrRemoteComponentPanel:
			switch m.CurrentLocalBranchOrTagComponentShowing {
			case constant.SHOW_LOCAL_BRANCH:
				CurrentSelectedBranch := m.CurrentRepoBranchesInfoList.SelectedItem()
				if CurrentSelectedBranch == nil {
					keys = i18n.LANGUAGEMAPPING.KeyBindingLocalBranchComponentNone
				} else {
					isCurrentSelectedBranchCheckedOutBranch := CurrentSelectedBranch.(branchComponent.GitBranchItem).IsCheckedOut
					if isCurrentSelectedBranchCheckedOutBranch {
						keys = i18n.LANGUAGEMAPPING.KeyBindingLocalBranchComponentIsCheckOut
					} else {
						keys = i18n.LANGUAGEMAPPING.KeyBindingLocalBranchComponentDefault
					}
				}
			case constant.SHOW_TAG:
				CurrentSelectedTag := m.CurrentRepoTagInfoList.SelectedItem()
				if CurrentSelectedTag == nil {
					keys = i18n.LANGUAGEMAPPING.KeyBindingTagComponentNone
				} else {
					keys = i18n.LANGUAGEMAPPING.KeyBindingTagComponentDefault
				}
			case constant.SHOW_REMOTE:
				CurrentSelectedRemote := m.CurrentRepoRemoteInfoList.SelectedItem()
				if CurrentSelectedRemote == nil {
					keys = i18n.LANGUAGEMAPPING.KeyBindingRemoteComponentNone
				} else {
					keys = i18n.LANGUAGEMAPPING.KeyBindingRemoteComponentDefault
				}
			}
		case constant.ModifiedFilesComponentPanel:
			CurrentSelectedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
			if CurrentSelectedFile == nil {
				keys = i18n.LANGUAGEMAPPING.KeyBindingModifiedFilesComponentNone
			} else {
				file := CurrentSelectedFile.(filesComponent.GitModifiedFilesItem)
				if file.HasConflict {
					keys = i18n.LANGUAGEMAPPING.KeyBindingModifiedFilesComponentConflict
				} else {
					if file.IndexState == "?" && file.WorkTree == "?" {
						// not tracked
						keys = i18n.LANGUAGEMAPPING.KeyBindingModifiedFilesComponentDefault
					} else if file.IndexState != " " && file.WorkTree != " " {
						// staged but have modification later
						keys = i18n.LANGUAGEMAPPING.KeyBindingModifiedFilesComponentDefault
					} else if file.IndexState != " " && file.WorkTree == " " {
						// staged and no latest modification
						keys = i18n.LANGUAGEMAPPING.KeyBindingModifiedFilesComponentIsStaged
					} else if file.IndexState == " " && file.WorkTree != " " {
						// tracked but not staged
						keys = i18n.LANGUAGEMAPPING.KeyBindingModifiedFilesComponentDefault
					}
				}
			}
		case constant.CommitLogComponentPanel:
			keys = i18n.LANGUAGEMAPPING.KeyBindingCommitLogComponent
		case constant.DetailComponentPanel:
			if m.IsLineEditingState.Load() {
				keys = i18n.LANGUAGEMAPPING.KeyBindingKeyDetailComponentLineEditing
			} else if m.DetailPanelParentComponent == constant.ModifiedFilesComponentPanel {
				keys = i18n.LANGUAGEMAPPING.KeyBindingKeyDetailComponentLineEditingEligible
			} else {
				keys = i18n.LANGUAGEMAPPING.KeyBindingKeyDetailComponent
			}

		case constant.DetailComponentPanelTwo:
			if m.IsLineEditingState.Load() {
				keys = i18n.LANGUAGEMAPPING.KeyBindingKeyDetailComponentLineEditing
			} else if m.DetailPanelParentComponent == constant.ModifiedFilesComponentPanel {
				keys = i18n.LANGUAGEMAPPING.KeyBindingKeyDetailComponentLineEditingEligible
			} else {
				keys = i18n.LANGUAGEMAPPING.KeyBindingKeyDetailComponent
			}
		case constant.StashComponentPanel:
			if len(m.CurrentRepoStashInfoList.Items()) > 0 {
				keys = i18n.LANGUAGEMAPPING.KeyBindingKeyStashComponent
			} else {
				keys = i18n.LANGUAGEMAPPING.KeyBindingKeyStashComponentNone
			}
		case constant.LogComponentPanel:
			keys = i18n.LANGUAGEMAPPING.KeyBindingLogComponent
		}
	}

	var keyBindingLine string
	keyBindingLine = strings.Join(keys, "  |  ")
	processedWidth := width - lipgloss.Width(gitticonst.APPVERSION) - 3
	keyBindingLine = utils.TruncateString(keyBindingLine, processedWidth)
	versionLine := style.NewStyle.Foreground(style.ColorYellowWarm).Render(gitticonst.APPVERSION)
	parsedKeyBindingLine := style.NewStyle.Width(processedWidth).Align(lipgloss.Left).Render(keyBindingLine)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		parsedKeyBindingLine,
		" ",
		versionLine,
	)

	return style.BottomKeyBindingStyle.
		Width(width).
		Height(constant.MainPageKeyBindingLayoutPanelHeight).
		Render(content)
}
