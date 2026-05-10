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

// ------------------------------------
//
//	Entry point for all Bubbletea key-press events. Handles global shortcuts
//	(ctrl+c quit, ctrl+s/u stage/unstage all, ctrl+g/l open browser links,
//	ctrl+f fetch), then dispatches to typing or non-typing key handlers based
//	on m.IsTyping.
//
// ------------------------------------
func GittiKeyInteraction(msg tea.KeyPressMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
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
	case "ctrl+f":
		services.GitFetchService(m)
		return m, nil
	}

	// for typing mode, it will always and must be a pop up
	if m.IsTyping.Load() {
		return handler.HandleTypingKeyBindingInteraction(msg, m)
	} else {
		return handler.HandleNonTypingGlobalKeyBindingInteraction(msg, m)
	}
}
