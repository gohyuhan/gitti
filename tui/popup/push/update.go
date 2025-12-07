package push

import (
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
		var GitPushLog string
		for _, line := range logs {
			logLine := style.NewStyle.Render(line)
			GitPushLog += logLine + "\n"
		}
		popUp.GitRemotePushOutputViewport.SetContent(GitPushLog)
		popUp.GitRemotePushOutputViewport.PageDown()
	}
}
