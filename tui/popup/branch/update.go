package branch

import (
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

func UpdateSwitchBranchOutputViewPort(m *types.GittiModel, gitOpsOutput []string) {
	popUp, ok := m.PopUpModel.(*SwitchBranchOutputPopUpModel)
	if ok {
		popUp.SwitchBranchOutputViewport.SetWidth(min(constant.MaxSwitchBranchOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		popUp.SwitchBranchOutputViewport.SetYOffset(popUp.SwitchBranchOutputViewport.YOffset())
		var gitOpsOutputLogs string
		for _, line := range gitOpsOutput {
			logLine := style.NewStyle.Render(line)
			gitOpsOutputLogs += logLine + "\n"
			popUp.SwitchBranchOutputViewport.SetContent(gitOpsOutputLogs)
			popUp.SwitchBranchOutputViewport.PageDown()
		}
	}
}
