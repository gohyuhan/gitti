package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/constant"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'S' key interaction.
//	Responsibility: Contextual "stash all" operation.
//	In the Modified Files Panel, this triggers a popup to stash *all* currently modified
//	tracked files in the repository, prompting the user for a stash message.
//
// ------------------------------------
func handleNonTypingSKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.CurrentSelectedComponent == constant.ModifiedFilesComponentPanel {
		currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
		var filePathName string
		if currentSelectedModifiedFile != nil {
			filePathName = currentSelectedModifiedFile.(files.GitModifiedFilesItem).FilePathname
			m.PopUpType = constant.GitStashMessagePopUp
			stashPopUp.InitGitStashMessagePopUpModel(m, filePathName, git.STASHALL)
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(true)
		}
	}
	return m, nil
}
