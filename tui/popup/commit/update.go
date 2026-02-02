package commit

import (
	"strings"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

func UpdatePopUpCommitOutputViewPort(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
	if ok {
		popUp.GitCommitOutputViewport.SetWidth(min(constant.MaxCommitPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		popUp.GitCommitOutputViewport.SetYOffset(popUp.GitCommitOutputViewport.YOffset())
		var gitCommitOutputLogString strings.Builder
		logs := m.GitOperations.GitCommit.GitCommitOutput()
		for _, line := range logs {
			logLine := style.NewStyle.Render(line)
			gitCommitOutputLogString.WriteString(logLine)
			gitCommitOutputLogString.WriteRune('\n')
		}
		popUp.GitCommitOutputViewport.SetContent(gitCommitOutputLogString.String())
		popUp.GitCommitOutputViewport.PageDown()
	}
}

// to update the amend commit output log for git amend commit
// this also take care of log by pre commit and post commit
func UpdatePopUpAmendCommitOutputViewPort(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*GitAmendCommitPopUpModel)
	if ok {
		popUp.GitAmendCommitOutputViewport.SetWidth(min(constant.MaxAmendCommitPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		popUp.GitAmendCommitOutputViewport.SetYOffset(popUp.GitAmendCommitOutputViewport.YOffset())
		var gitCommitOutputLogString strings.Builder
		logs := m.GitOperations.GitCommit.GitCommitOutput()
		for _, line := range logs {
			logLine := style.NewStyle.Render(line)
			gitCommitOutputLogString.WriteString(logLine)
			gitCommitOutputLogString.WriteRune('\n')
		}
		popUp.GitAmendCommitOutputViewport.SetContent(gitCommitOutputLogString.String())
		popUp.GitAmendCommitOutputViewport.PageDown()
	}
}
