package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	Launches the user's defined system editor on the given path (handling terminal/GUI diffs).
//	A non-terminal (GUI) editor is started in the background, a terminal editor takes over the
//	terminal and reports back with types.EditorFinishedMsg.
//
// ------------------------------------
func launchEditor(m *types.GittiModel, path string) (*types.GittiModel, tea.Cmd) {
	cmd, isNonTerminalEditor := utils.ReturnEditorLaunchCommand(path, m.UserSetEditor)
	if isNonTerminalEditor {
		cmd.Start()
		return m, nil
	} else {
		return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
			return types.EditorFinishedMsg{
				Err: err,
			}
		})
	}
}
