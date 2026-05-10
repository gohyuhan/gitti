package rebase

import (
	"strings"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Sync the rebase output viewport content and scroll position from the latest
//	git rebase output lines. Called whenever the underlying output buffer changes
//	during an in-progress or completed rebase operation.
//
// ------------------------------------
func UpdatePopUpGitRebaseOutputViewport(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*GitRebaseOutputPopUpModel)
	if ok {
		popUp.GitRebaseOutputViewport.SetWidth(min(constant.MaxGitRebaseOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		popUp.GitRebaseOutputViewport.SetYOffset(popUp.GitRebaseOutputViewport.YOffset())
		logs := m.GitOperations.GitRebase.GetGitRebaseOutput()
		var GitRebaseLogString strings.Builder
		for _, line := range logs {
			logLine := style.NewStyle.Render(line)
			GitRebaseLogString.WriteString(logLine)
			GitRebaseLogString.WriteRune('\n')
		}
		popUp.GitRebaseOutputViewport.SetContent(GitRebaseLogString.String())
		popUp.GitRebaseOutputViewport.PageDown()
	}
}

// ------------------------------------
//
//	Handle the async git rebase result event. Clears the IsProcessing flag and
//	sets ProcessSuccess on success, or sets HasError on failure. No-ops if
//	the popup is not the active rebase output popup or the operation was cancelled.
//
// ------------------------------------
func UpdateGitRebaseResultEvent(m *types.GittiModel, data types.GitRebaseResultEventDataStructure) {
	popUp, ok := m.PopUpModel.(*GitRebaseOutputPopUpModel)
	if !ok || popUp.IsCancelled.Load() {
		return
	}

	popUp.IsProcessing.Store(false)
	if data.Success {
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(true)
		return
	}

	popUp.HasError.Store(true)
	popUp.ProcessSuccess.Store(false)
}
