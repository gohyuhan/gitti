package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'i' key interaction.
//	Responsibility: Opens the interactive rebase option popup when the commit log panel
//	is focused and at least one commit exists in the list.
//
// ------------------------------------
func handleNonTypingiKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.CurrentSelectedComponent == constant.CommitLogOrRefLogComponentPanel {
			// interactive rebase only possible when there is commit
			if len(m.CurrentRepoCommitLogInfoList.Items()) > 0 {
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(false)
				m.PopUpType = constant.InteractiveRebaseOptionPopUp
				interactiverebasePopUp.InitInteractiveRebaseOptionPopUpModel(m)
			}
		}
	}

	return m, nil
}
