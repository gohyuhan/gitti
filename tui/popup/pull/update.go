package pull

import (
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

func UpdatePopUpGitPullOutputViewport(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*GitPullOutputPopUpModel)
	if ok {
		popUp.GitPullOutputViewport.SetWidth(min(constant.MaxGitPullOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		popUp.GitPullOutputViewport.SetYOffset(popUp.GitPullOutputViewport.YOffset())
		logs := m.GitOperations.GitPull.GetGitPullOutput()
		var GitPullLog string
		for _, line := range logs {
			logLine := style.NewStyle.Render(line)
			GitPullLog += logLine + "\n"
		}
		popUp.GitPullOutputViewport.SetContent(GitPullLog)
		popUp.GitPullOutputViewport.PageDown()
	}
}
