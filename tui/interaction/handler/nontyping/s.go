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
//	Handle 's' key interaction.
//	Responsibility: Contextual "stash file" operation.
//	In the Modified Files Panel, if a valid (non-conflicted) uncommitted file is selected,
//	this triggers a popup to stash *only* that specific file, prompting for a stash message.
//
// ------------------------------------
func handleNonTypingsKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if m.CurrentSelectedComponent == constant.ModifiedFilesComponentPanel {
		currentSelectedModifiedFile := m.CurrentRepoModifiedFilesInfoList.SelectedItem()
		var filePathName string
		if currentSelectedModifiedFile != nil {
			selectedFile := currentSelectedModifiedFile.(files.GitModifiedFilesItem)
			// return early if the file is in a conflict status
			if selectedFile.HasConflict {
				return m, nil
			}
			filePathName = selectedFile.FilePathname
			m.PopUpType = constant.GitStashMessagePopUp
			stashPopUp.InitGitStashMessagePopUpModel(m, filePathName, git.STASHFILE)
			m.ShowPopUp.Store(true)
			m.IsTyping.Store(true)
		}
	}
	return m, nil
}
