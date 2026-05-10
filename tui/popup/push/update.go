package push

import (
	"strings"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

func UpdatePopUpGitRemotePushOutputViewport(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
	if ok {
		popUp.GitRemotePushOutputViewport.SetWidth(min(constant.MaxGitRemotePushPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		popUp.GitRemotePushOutputViewport.SetYOffset(popUp.GitRemotePushOutputViewport.YOffset())
		logs := m.GitOperations.GitCommit.GitRemotePushOutput()
		var GitPushLogString strings.Builder
		for _, line := range logs {
			logLine := style.NewStyle.Render(line)
			GitPushLogString.WriteString(logLine)
			GitPushLogString.WriteRune('\n')
		}
		popUp.GitRemotePushOutputViewport.SetContent(GitPushLogString.String())
		popUp.GitRemotePushOutputViewport.PageDown()
	}
}

func UpdateGitPushResultEvent(m *types.GittiModel, data types.GitPushResultEventDataInterface) {
	popUp, ok := m.PopUpModel.(*GitRemotePushPopUpModel)
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
