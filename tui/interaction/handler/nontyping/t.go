package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/component/commitlog"
	"github.com/gohyuhan/gitti/tui/constant"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle 't' key interaction.
//	Responsibility: Initiates the "Create Tag" operation.
//	In the Commit Log Panel, selecting a specific commit and pressing 't'
//	opens a popup to create a new git tag referencing that exact commit hash.
//
// ------------------------------------
func handleNonTypingtKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	switch m.CurrentSelectedComponent {
	case constant.CommitLogOrRefLogComponentPanel:
		switch m.CurrentCommitLogOrRefLogComponentShowing {
		case constant.SHOW_COMMITLOG:
			currentSelectedCommit := m.CurrentRepoCommitLogInfoList.SelectedItem()
			if currentSelectedCommit != nil {
				commit := currentSelectedCommit.(commitlog.GitCommitLogItem)
				m.PopUpType = constant.CreateTagPopUp
				tagPopUp.InitCreateTagPopUpModel(m, commit.Hash, commit.Message)
				m.ShowPopUp.Store(true)
				m.IsTyping.Store(true)
			}
		}
	}
	return m, nil
}
