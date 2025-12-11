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

// for bubbletea list component, we can't get rid of the "No items." natively for now as there was no exposed api
// see https://github.com/charmbracelet/bubbles/blob/master/list/list.go#L1222
// we are using ReplaceAll as a hack for now to replace "No items." with ""

// -----------------------------------------------------------------------------
//
//	Functions that help construct the view
//
// -----------------------------------------------------------------------------
// render the Gitti Status Panel
func renderGitStatusComponentPanel(m *types.GittiModel) string {
	borderStyle := style.PanelBorderStyle
	if m.CurrentSelectedComponent == constant.GitStatusComponent {
		borderStyle = style.SelectedBorderStyle
	}

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
}

// Render the Local Branches panel
func renderLocalBranchesComponentPanel(width int, height int, m *types.GittiModel) string {
	borderStyle := style.PanelBorderStyle
	if m.CurrentSelectedComponent == constant.LocalBranchComponent {
		borderStyle = style.SelectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(strings.ReplaceAll(m.CurrentRepoBranchesInfoList.View(), "No items.", ""))
}

// Render the Changed Files panel
func renderModifiedFilesComponentPanel(width int, height int, m *types.GittiModel) string {
	borderStyle := style.PanelBorderStyle
	if m.CurrentSelectedComponent == constant.ModifiedFilesComponent {
		borderStyle = style.SelectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(strings.ReplaceAll(m.CurrentRepoModifiedFilesInfoList.View(), "No items.", ""))
}

// Render the detail component part at the right of the window,
// however the content within it will be dynamic based on the current selected component
func renderDetailComponentPanel(width int, height int, m *types.GittiModel) string {
	detailComponentBorderStyle := style.PanelBorderStyle
	detailComponentTwoBorderStyle := style.PanelBorderStyle

	if m.CurrentSelectedComponent == constant.DetailComponent {
		detailComponentBorderStyle = style.SelectedBorderStyle
	} else if m.CurrentSelectedComponent == constant.DetailComponentTwo {
		detailComponentTwoBorderStyle = style.SelectedBorderStyle
	}

	var content string
	detailPanelWidth := width
	detailPanelHeight := height

	if m.ShowDetailPanelTwo.Load() {
		detailPanelHeight = (height / 2)
		detailPanelWidth = (width / 2)
		if m.DetailComponentPanelLayout == constant.HORIZONTAL {
			content = lipgloss.JoinHorizontal(
				lipgloss.Top,
				detailComponentBorderStyle.Width(detailPanelWidth).Height(height).Render(m.DetailPanelViewport.View()),
				detailComponentTwoBorderStyle.Width(detailPanelWidth).Height(height).Render(m.DetailPanelTwoViewport.View()),
			)
		} else {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				detailComponentBorderStyle.Width(width).Height(detailPanelHeight).Render(m.DetailPanelViewport.View()),
				detailComponentTwoBorderStyle.Width(width).Height(detailPanelHeight).Render(m.DetailPanelTwoViewport.View()),
			)
		}
	} else {
		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			detailComponentBorderStyle.Width(detailPanelWidth).Height(detailPanelHeight).Render(m.DetailPanelViewport.View()),
		)
	}

	return style.NewStyle.
		Width(width).
		Height(height).
		Render(content)
}

func renderStashComponentPanel(width int, height int, m *types.GittiModel) string {
	borderStyle := style.PanelBorderStyle
	if m.CurrentSelectedComponent == constant.StashComponent {
		borderStyle = style.SelectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(strings.ReplaceAll(m.CurrentRepoStashInfoList.View(), "No items.", ""))
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
		case constant.GlobalKeyBindingPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGlobalKeyBindingPopUp
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
		}
	} else {
		//-----------------------------
		//
		// for non-popup keybinding render
		//
		//-----------------------------
		switch m.CurrentSelectedComponent {
		case constant.GitStatusComponent:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitStatusComponent
		case constant.LocalBranchComponent:
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
		case constant.ModifiedFilesComponent:
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
		case constant.DetailComponent:
			keys = i18n.LANGUAGEMAPPING.KeyBindingKeyDetailComponent
		case constant.DetailComponentTwo:
			keys = i18n.LANGUAGEMAPPING.KeyBindingKeyDetailComponent
		case constant.StashComponent:
			if len(m.CurrentRepoStashInfoList.Items()) > 0 {
				keys = i18n.LANGUAGEMAPPING.KeyBindingKeyStashComponent
			} else {
				keys = i18n.LANGUAGEMAPPING.KeyBindingKeyStashComponentNone
			}
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
