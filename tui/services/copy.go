package services

import (
	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/tui/types"
)

// StartCopy allows one clipboard write at a time so an earlier write cannot
// finish after a later one and replace the user's most recent selection.
func StartCopy(m *types.GittiModel, value string) (tea.Cmd, bool) {
	if m.CopyInProgress {
		m.GittiLogger.RegisterNewLog(logging.COPY_VALUE_OPS, "", logging.WARN, i18n.LANGUAGEMAPPING.CopyInProgress, false)
		return nil, false
	}
	m.CopyInProgress = true
	return func() tea.Msg {
		return types.CopyFinishedMsg{Value: value, Err: clipboard.WriteAll(value)}
	}, true
}
