package rebase

import (
	"strings"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

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

func UpdateGitRebaseResultEvent(m *types.GittiModel, data types.GitRebaseResultEventDataInterface) {
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
