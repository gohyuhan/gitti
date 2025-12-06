package remote

import (
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

func UpdateAddRemoteOutputViewport(m *types.GittiModel, outputLog []string) {
	popUp, ok := m.PopUpModel.(*AddRemotePromptPopUpModel)
	if ok {
		popUp.AddRemoteOutputViewport.SetWidth(min(constant.MaxAddRemotePromptPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		popUp.AddRemoteOutputViewport.SetYOffset(popUp.AddRemoteOutputViewport.YOffset())
		var addRemoteLog string
		for _, line := range outputLog {
			logLine := style.NewStyle.Render(line)
			addRemoteLog += logLine + "\n"
		}
		popUp.AddRemoteOutputViewport.SetContent(addRemoteLog)
		popUp.AddRemoteOutputViewport.PageDown()
	}
}
