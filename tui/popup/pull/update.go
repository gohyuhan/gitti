package pull

import (
	"strings"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Sync the pull output viewport content and scroll position from the latest
//	git pull output lines. Called whenever the underlying output buffer changes
//	during an in-progress or completed pull operation.
//
// ------------------------------------
func UpdatePopUpGitPullOutputViewport(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*GitPullOutputPopUpModel)
	if ok {
		popUp.GitPullOutputViewport.SetWidth(min(constant.MaxGitPullOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		popUp.GitPullOutputViewport.SetYOffset(popUp.GitPullOutputViewport.YOffset())
		logs := m.GitOperations.GitPull.GetGitPullOutput()
		var GitPullLogString strings.Builder
		for _, line := range logs {
			logLine := style.NewStyle.Render(line)
			GitPullLogString.WriteString(logLine)
			GitPullLogString.WriteRune('\n')
		}
		popUp.GitPullOutputViewport.SetContent(GitPullLogString.String())
		popUp.GitPullOutputViewport.PageDown()
	}
}

// ------------------------------------
//
//	Handle the async git pull result event. Clears the IsProcessing flag and
//	sets ProcessSuccess on success, or sets HasError on failure. No-ops if
//	the popup is not the active pull output popup or the operation was cancelled.
//
// ------------------------------------
func UpdateGitPullResultEvent(m *types.GittiModel, data types.GitPullResultEventDataStructure) {
	popUp, ok := m.PopUpModel.(*GitPullOutputPopUpModel)
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
