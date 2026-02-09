package tag

import (
	"strings"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

func UpdateDeleteTagOutputViewPort(m *types.GittiModel, deleteTagOutput []string) {
	popUp, ok := m.PopUpModel.(*DeleteTagOutputPopUpModel)
	if ok {
		popUp.DeleteTagOutputViewport.SetWidth(min(constant.MaxDeleteTagOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		popUp.DeleteTagOutputViewport.SetYOffset(popUp.DeleteTagOutputViewport.YOffset())
		var gitOpsOutputLogString strings.Builder
		for _, line := range deleteTagOutput {
			logLine := style.NewStyle.Render(line)
			gitOpsOutputLogString.WriteString(logLine)
			gitOpsOutputLogString.WriteRune('\n')
		}
		popUp.DeleteTagOutputViewport.SetContent(gitOpsOutputLogString.String())
		popUp.DeleteTagOutputViewport.PageDown()
	}
}
