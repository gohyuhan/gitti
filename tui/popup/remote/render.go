package remote

import (
	"fmt"
	"strings"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"

	"charm.land/lipgloss/v2"
)

// ------------------------------------
//
//	For Adding Git Remote
//
// ------------------------------------
func RenderAddRemotePromptPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxAddRemotePromptPopUpWidth, int(float64(m.Width)*0.8))
		popUp.RemoteNameTextInput.SetWidth(popUpWidth - 6)
		popUp.RemoteUrlTextInput.SetWidth(popUpWidth - 6)

		noInitialRemote := popUp.NoInitialRemote

		// Rendered content
		addRemotePrompt := style.PromptTitleStyle.Render(i18n.LANGUAGEMAPPING.AddRemotePopUpPrompt)
		remoteNameTitle := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteNameTitle)
		remoteNameInputView := popUp.RemoteNameTextInput.View()
		remoteUrlTitle := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteUrlTitle)
		remoteUrlTitleInputView := popUp.RemoteUrlTextInput.View()

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			remoteNameTitle,
			remoteNameInputView,
			remoteUrlTitle,
			remoteUrlTitleInputView,
		)
		if noInitialRemote {
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				addRemotePrompt,
				remoteNameTitle,
				remoteNameInputView,
				remoteUrlTitle,
				remoteUrlTitleInputView,
			)
		}
		if popUp.AddRemoteOutputViewport.GetContent() != "" {
			logViewPortStyle := style.PanelBorderStyle.
				Width(popUpWidth - 2).
				Height(constant.PopUpAddRemoteOutputViewPortHeight + 2)
			if popUp.HasError.Load() {
				logViewPortStyle = style.PanelBorderStyle.
					BorderForeground(style.ColorError)
			} else if popUp.ProcessSuccess.Load() {
				logViewPortStyle = style.PanelBorderStyle.
					BorderForeground(style.ColorGreenSoft)
			}
			popUp.AddRemoteOutputViewport.SetWidth(popUpWidth - 4)
			popUp.AddRemoteOutputViewport.SetYOffset(popUp.AddRemoteOutputViewport.YOffset())
			logViewPort := logViewPortStyle.Render(popUp.AddRemoteOutputViewport.View())
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				remoteNameTitle,
				remoteNameInputView,
				remoteUrlTitle,
				remoteUrlTitleInputView,
				logViewPort,
			)
			if noInitialRemote {
				content = lipgloss.JoinVertical(
					lipgloss.Left,
					addRemotePrompt,
					remoteNameTitle,
					remoteNameInputView,
					remoteUrlTitle,
					remoteUrlTitleInputView,
					logViewPort,
				)
			}
		}
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

// ------------------------------------
//
//	For Choosing a Remote for git push if there is more than 1
//
// ------------------------------------
func RenderChooseRemotePopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*ChooseRemotePopUpModel)
	if ok {
		popUpWidth := min(constant.MaxChooseRemotePopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(i18n.LANGUAGEMAPPING.ChooseRemoteTitle)
		popUp.RemoteList.SetWidth(popUpWidth - 4)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			popUp.RemoteList.View(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}

func RenderRemoveRemoteConfirmationPopUp(m *types.GittiModel) string {
	popUp, ok := m.PopUpModel.(*RemoveRemoteConfirmationPopUpModel)
	if ok {
		popUpWidth := min(constant.MaxRemoveRemoteConfirmationPopUpWidth, int(float64(m.Width)*0.8))
		title := style.TitleStyle.Render(fmt.Sprintf(i18n.LANGUAGEMAPPING.RemoveRemoteTitle, style.NewStyle.Foreground(style.ColorYellowWarm).Render(popUp.RemoteName)))

		var remoteInfo strings.Builder
		urlLabel := "URL:"
		fetchLabel := i18n.LANGUAGEMAPPING.Fetch
		pushLabel := i18n.LANGUAGEMAPPING.Push
		urlLen := len([]rune(urlLabel))
		fetchLen := len([]rune(fetchLabel))
		pushLen := len([]rune(pushLabel))
		maxLen := max(urlLen, max(fetchLen, pushLen)) + 1 // plus 1 for spacing

		// Render URL with padding
		remoteInfo.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render(urlLabel))
		for i := 0; i < maxLen-urlLen; i++ {
			remoteInfo.WriteString(" ")
		}
		remoteInfo.WriteString(style.NewStyle.Foreground(style.ColorPurpleSoft).Render(popUp.RemoteUrl))
		remoteInfo.WriteRune('\n')

		// Render Fetch with padding
		remoteInfo.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render(fetchLabel))
		for i := 0; i < maxLen-fetchLen; i++ {
			remoteInfo.WriteString(" ")
		}
		if popUp.Fetch {
			remoteInfo.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("["))
			remoteInfo.WriteString(style.NewStyle.Foreground(style.ColorPurpleSoft).Render("X"))
			remoteInfo.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("]"))
		} else {
			remoteInfo.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("["))
			remoteInfo.WriteString(" ")
			remoteInfo.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("]"))
		}
		remoteInfo.WriteRune('\n')

		// Render Push with padding
		remoteInfo.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render(pushLabel))
		for i := 0; i < maxLen-pushLen; i++ {
			remoteInfo.WriteString(" ")
		}
		if popUp.Push {
			remoteInfo.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("["))
			remoteInfo.WriteString(style.NewStyle.Foreground(style.ColorPurpleSoft).Render("X"))
			remoteInfo.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("]"))
		} else {
			remoteInfo.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("["))
			remoteInfo.WriteString(" ")
			remoteInfo.WriteString(style.NewStyle.Foreground(style.ColorPurpleVibrant).Render("]"))
		}

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			remoteInfo.String(),
		)
		return style.PopUpBorderStyle.Width(popUpWidth).Render(content)
	}
	return ""
}
