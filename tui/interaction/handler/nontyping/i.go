package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/constant"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 'i' in the files and commit log panels.
//
// ------------------------------------
func handleNonTypingiKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		switch m.CurrentSelectedComponent {
		case constant.ModifiedFilesComponentPanel:
			item, ok := m.CurrentRepoModifiedFilesInfoList.SelectedItem().(files.GitModifiedFilesItem)
			if !ok {
				return m, nil
			}
			if m.IgnoreInProgress {
				m.GittiLogger.RegisterNewLog(logging.IGNORE_FILE_OPS, "", logging.WARN, i18n.LANGUAGEMAPPING.IgnoreInProgress, false)
				return m, nil
			}
			if item.IndexState != "?" || item.WorkTree != "?" {
				m.GittiLogger.RegisterNewLog(logging.IGNORE_FILE_OPS, "", logging.WARN, i18n.LANGUAGEMAPPING.IgnoreTracked, false)
				return m, nil
			}
			gitFiles := m.GitOperations.GitFiles
			filePath := item.NewFilePathname
			m.IgnoreInProgress = true
			return m, func() tea.Msg {
				result, err := gitFiles.IgnoreFile(filePath)
				refreshErr := gitFiles.RefreshFilesStatus()
				return types.IgnoreFinishedMsg{GitFiles: gitFiles, FilePath: filePath, Result: result, Err: err, RefreshErr: refreshErr}
			}
		case constant.CommitLogOrRefLogComponentPanel:
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
