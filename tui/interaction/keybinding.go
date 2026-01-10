package interaction

import (
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/interaction/handler"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/utils"

	tea "charm.land/bubbletea/v2"
)

// the function to handle bubbletea key interactions
func GittiKeyInteraction(msg tea.KeyMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	// global key binding
	switch msg.String() {
	case "ctrl+c":
		if api.GITDAEMON != nil {
			api.GITDAEMON.Stop()
		}
		return m, tea.Quit
	case "ctrl+s":
		services.GitStageAllChangesService(m)
		return m, nil
	case "ctrl+u":
		services.GitUnstageAllChangesService(m)
		return m, nil
	case "ctrl+g":
		utils.OpenBrowser(constant.AUTHOR_GITHUB)
		return m, nil
	case "ctrl+l":
		utils.OpenBrowser(constant.AUTHOR_LINKEDIN)
		return m, nil
	}

	if m.IsTyping.Load() {
		return handler.HandleTypingKeyBindingInteraction(msg, m)
	} else {
		return handler.HandleNonTypingGlobalKeyBindingInteraction(msg, m)
	}
}
