package commit

import (
	"strings"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Sync the commit output viewport content and scroll position from the latest
//	git commit output lines. Called whenever the underlying output buffer changes
//	during an in-progress or completed commit operation.
//
// ------------------------------------
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

// ------------------------------------
//
//	Sync the amend-commit output viewport content and scroll position from the
//	latest git commit output lines, including output from pre-commit and
//	post-commit hooks. Called whenever the output buffer changes during an
//	in-progress or completed amend operation.
//
// ------------------------------------
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

// ------------------------------------
//
//	Handle the async git commit result event. Clears the IsProcessing flag and
//	resets both text inputs on success; sets HasError on failure. No-ops if the
//	popup is not the active commit popup or the operation was cancelled.
//
// ------------------------------------
func UpdateGitCommitResultEvent(m *types.GittiModel, data types.GitCommitResultEventDataStructure) {
	popUp, ok := m.PopUpModel.(*GitCommitPopUpModel)
	if !ok || popUp.IsCancelled.Load() {
		return
	}

	popUp.IsProcessing.Store(false)
	if data.Success {
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(true)
		popUp.MessageTextInput.Reset()
		popUp.DescriptionTextAreaInput.Reset()
		return
	}

	popUp.HasError.Store(true)
	popUp.ProcessSuccess.Store(false)
}

// ------------------------------------
//
//	Handle the async git amend-commit result event. Clears the IsProcessing flag
//	and resets both text inputs on success; sets HasError on failure. No-ops if
//	the popup is not the active amend-commit popup or the operation was cancelled.
//
// ------------------------------------
func UpdateGitAmendCommitResultEvent(m *types.GittiModel, data types.GitAmendCommitResultEventDataStructure) {
	popUp, ok := m.PopUpModel.(*GitAmendCommitPopUpModel)
	if !ok || popUp.IsCancelled.Load() {
		return
	}

	popUp.IsProcessing.Store(false)
	if data.Success {
		popUp.HasError.Store(false)
		popUp.ProcessSuccess.Store(true)
		popUp.MessageTextInput.Reset()
		popUp.DescriptionTextAreaInput.Reset()
		return
	}

	popUp.HasError.Store(true)
	popUp.ProcessSuccess.Store(false)
}
