package tui

import (
	"fmt"
	"strings"

	gitticonst "github.com/gohyuhan/gitti/constant"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/component/branch"
	"github.com/gohyuhan/gitti/tui/constant"
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

	trackedUpStreamOrBranchName := m.CheckOutBranch
	if m.BranchUpStream != "" {
		trackedUpStreamOrBranchName = m.BranchUpStream
	}

	repoTrackBranchName := fmt.Sprintf(" %s -> %s %s", m.RepoName, m.TrackedUpstreamOrBranchIcon, trackedUpStreamOrBranchName)

	// the max width is the window width - padding - the length of RemoteSyncStateLineString (max of 5)
	repoTrackBranchName = utils.TruncateString(repoTrackBranchName, m.WindowLeftPanelWidth-constant.ListItemOrTitleWidthPad-5)

	return borderStyle.
		Width(m.WindowLeftPanelWidth).
		Height(1).
		Render(fmt.Sprintf("%s%s", m.RemoteSyncStateLineString, repoTrackBranchName))
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
	borderStyle := style.PanelBorderStyle
	if m.CurrentSelectedComponent == constant.DetailComponent {
		borderStyle = style.SelectedBorderStyle
	}
	return borderStyle.
		Width(width).
		Height(height).
		Render(m.DetailPanelViewport.View())
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
			popUp, ok := m.PopUpModel.(*SwitchBranchOutputPopUpModel)
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
			popUp, ok := m.PopUpModel.(*GitStashOperationOutputPopUpModel)
			if ok {
				if popUp.IsProcessing.Load() {
					keys = []string{"..."} // nothing can be done during stash operation, only force quit gitti is possible
				}
			}
		case constant.GitStashConfirmPromptPopUp:
			keys = i18n.LANGUAGEMAPPING.KeyBindingForGitStashConfirmPromptPopUp
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
				isCurrentSelectedBranchCheckedOutBranch := CurrentSelectedBranch.(branch.GitBranchItem).IsCheckedOut
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
				file := CurrentSelectedFile.(gitModifiedFilesItem)
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
	processedWidth := width - len(gitticonst.APPVERSION) - 3
	keyBindingLine = utils.TruncateString(keyBindingLine, processedWidth)
	versionLine := style.NewStyle.Foreground(style.ColorYellowWarm).Render(gitticonst.APPVERSION)
	content := fmt.Sprintf("%-*s  %s", processedWidth, keyBindingLine, versionLine)

	return style.BottomKeyBindingStyle.
		Width(width).
		Height(constant.MainPageKeyBindingLayoutPanelHeight).
		Render(content)
}

// to update the width and height of all components
func tuiWindowSizing(m *types.GittiModel) {
	// Compute panel widths
	m.WindowLeftPanelWidth = min(int(float64(m.Width)*settings.GITTICONFIGSETTINGS.LeftPanelWidthRatio), constant.MaxLeftPanelWidth)
	m.DetailComponentPanelWidth = m.Width - m.WindowLeftPanelWidth

	m.WindowCoreContentHeight = m.Height - constant.MainPageKeyBindingLayoutPanelHeight - 2*constant.Padding
	m.DetailComponentPanelHeight = m.WindowCoreContentHeight

	// update the dynamic size of the left panel
	leftPanelDynamicResize(m)

	// update viewport
	m.DetailPanelViewport.SetHeight(m.DetailComponentPanelHeight) // some margin
	m.DetailPanelViewport.SetWidth(m.DetailComponentPanelWidth - 2)
	m.DetailPanelViewportOffset = max(0, int(m.DetailPanelViewport.HorizontalScrollPercent()*float64(m.DetailPanelViewportOffset))-1)
	m.DetailPanelViewport.SetXOffset(m.DetailPanelViewportOffset)
	m.DetailPanelViewport.SetYOffset(m.DetailPanelViewport.YOffset())
}

func leftPanelDynamicResize(m *types.GittiModel) {
	leftPanelRemainingHeight := m.WindowCoreContentHeight - 7 // this is after reserving the height for the gitti version panel and also Padding

	// we minus 2 if GitStatusComponent is not the one chosen is because GitStatusComponent
	// and the one that got selected will not be account in to the dynamic height calculation
	// ( gitti status component's height is fix at 3, while the selected one will always get 40% )
	componentWithDynamicHeight := (len(constant.ComponentNavigationList) - 2)
	unSelectedComponentPanelHeightPerComponent := (int(float64(leftPanelRemainingHeight)*(1.0-constant.SelectedLeftPanelComponentHeightRatio)) / componentWithDynamicHeight)
	selectedComponentPanelHeight := leftPanelRemainingHeight - (unSelectedComponentPanelHeightPerComponent * componentWithDynamicHeight)
	m.LocalBranchesComponentPanelHeight = unSelectedComponentPanelHeightPerComponent
	m.ModifiedFilesComponentPanelHeight = unSelectedComponentPanelHeightPerComponent
	m.StashComponentPanelHeight = unSelectedComponentPanelHeightPerComponent

	switch m.CurrentSelectedComponent {
	case constant.LocalBranchComponent:
		m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
	case constant.ModifiedFilesComponent:
		m.ModifiedFilesComponentPanelHeight = selectedComponentPanelHeight
	case constant.StashComponent:
		m.StashComponentPanelHeight = selectedComponentPanelHeight
	case constant.GitStatusComponent:
		// if it was the Gitti status component panel that got selected (because its height is fix),
		// the next panel will get the selected height which is the branch component panel
		m.LocalBranchesComponentPanelHeight = selectedComponentPanelHeight
	}
	// update all components Width and Height
	m.CurrentRepoBranchesInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoBranchesInfoList.SetHeight(m.LocalBranchesComponentPanelHeight)

	m.CurrentRepoModifiedFilesInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoModifiedFilesInfoList.SetHeight(m.ModifiedFilesComponentPanelHeight)

	m.CurrentRepoStashInfoList.SetWidth(m.WindowLeftPanelWidth - 2)
	m.CurrentRepoStashInfoList.SetHeight(m.StashComponentPanelHeight)
}

// for about gitti content
func generateAboutGittiContent() string {
	var vpLine string

	logoLineArray := style.GradientLines(gittiAsciiArtLogo)
	aboutLines := i18n.LANGUAGEMAPPING.AboutGitti

	vpLine += strings.Join(logoLineArray, "\n") + "\n"
	vpLine += strings.Join(aboutLines, "\n")

	return vpLine
}
