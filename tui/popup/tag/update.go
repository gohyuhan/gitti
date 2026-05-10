package tag

import (
	"strings"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Refresh the delete tag output viewport. Resizes to the current terminal
//	width, renders each output line with the default style, and scrolls to the
//	bottom.
//
// ------------------------------------
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

// ------------------------------------
//
//	Handle the async delete tag result event. No-ops if the popup was cancelled.
//	Clears IsProcessing, refreshes the output viewport, then sets ProcessSuccess
//	on success or HasError on failure.
//
// ------------------------------------
func UpdateDeleteTagResultEvent(m *types.GittiModel, updateData types.GitDeleteTagResultEventDataStructure) {
	popUp, ok := m.PopUpModel.(*DeleteTagOutputPopUpModel)
	if ok && !popUp.IsCancelled.Load() {
		popUp.IsProcessing.Store(false)
		UpdateDeleteTagOutputViewPort(m, updateData.Result)
		if updateData.Success && !popUp.IsProcessing.Load() {
			popUp.ProcessSuccess.Store(true)
			popUp.HasError.Store(false)
		} else if !updateData.Success && !popUp.IsProcessing.Load() {
			popUp.ProcessSuccess.Store(false)
			popUp.HasError.Store(true)
		}
	}
}

// ------------------------------------
//
//	Refresh the push tag output viewport. Resizes to the current terminal width,
//	reads the latest output lines from GitTag.PushTagOutput(), renders each line,
//	and scrolls to the bottom.
//
// ------------------------------------
func UpdatePushTagOutputViewPort(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*PushTagOutputPopUpModel)
	if ok {
		popUp.PushTagOutputViewport.SetWidth(min(constant.MaxPushTagOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		popUp.PushTagOutputViewport.SetYOffset(popUp.PushTagOutputViewport.YOffset())
		var gitOpsOutputLogString strings.Builder
		pushTagOutput := m.GitOperations.GitTag.PushTagOutput()
		for _, line := range pushTagOutput {
			logLine := style.NewStyle.Render(line)
			gitOpsOutputLogString.WriteString(logLine)
			gitOpsOutputLogString.WriteRune('\n')
		}
		popUp.PushTagOutputViewport.SetContent(gitOpsOutputLogString.String())
		popUp.PushTagOutputViewport.PageDown()
	}
}

// ------------------------------------
//
//	Handle the async push tag result event. No-ops if the popup was cancelled.
//	Clears IsProcessing, refreshes the output viewport, then sets ProcessSuccess
//	on success or HasError on failure.
//
// ------------------------------------
func UpdatePushTagResultEvent(m *types.GittiModel, updateData types.GitPushTagResultEventDataStructure) {
	popUp, ok := m.PopUpModel.(*PushTagOutputPopUpModel)
	if ok && !popUp.IsCancelled.Load() {
		popUp.IsProcessing.Store(false)
		UpdatePushTagOutputViewPort(m)
		if updateData.Success && !popUp.IsProcessing.Load() {
			popUp.ProcessSuccess.Store(true)
			popUp.HasError.Store(false)
		} else if !updateData.Success && !popUp.IsProcessing.Load() {
			popUp.HasError.Store(true)
		}
	}
}

// ------------------------------------
//
//	Refresh the fetch tag output viewport. Resizes to the current terminal width,
//	reads the latest output lines from GitTag.FetchTagOutput(), renders each
//	line, and scrolls to the bottom.
//
// ------------------------------------
func UpdateFetchTagOutputViewPort(m *types.GittiModel) {
	popUp, ok := m.PopUpModel.(*FetchTagOutputPopUpModel)
	if ok {
		popUp.FetchTagOutputViewport.SetWidth(min(constant.MaxFetchTagOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)
		popUp.FetchTagOutputViewport.SetYOffset(popUp.FetchTagOutputViewport.YOffset())
		var gitOpsOutputLogString strings.Builder
		fetchTagOutput := m.GitOperations.GitTag.FetchTagOutput()
		for _, line := range fetchTagOutput {
			logLine := style.NewStyle.Render(line)
			gitOpsOutputLogString.WriteString(logLine)
			gitOpsOutputLogString.WriteRune('\n')
		}
		popUp.FetchTagOutputViewport.SetContent(gitOpsOutputLogString.String())
		popUp.FetchTagOutputViewport.PageDown()
	}
}

// ------------------------------------
//
//	Handle the async fetch tag result event. No-ops if the popup was cancelled.
//	Clears IsProcessing, refreshes the output viewport, then sets ProcessSuccess
//	on success or HasError on failure.
//
// ------------------------------------
func UpdateFetchTagResultEvent(m *types.GittiModel, updateData types.GitFetchTagResultEventDataStructure) {
	popUp, ok := m.PopUpModel.(*FetchTagOutputPopUpModel)
	if ok && !popUp.IsCancelled.Load() {
		popUp.IsProcessing.Store(false)
		UpdateFetchTagOutputViewPort(m)
		if updateData.Success && !popUp.IsProcessing.Load() {
			popUp.ProcessSuccess.Store(true)
			popUp.HasError.Store(false)
		} else if !updateData.Success && !popUp.IsProcessing.Load() {
			popUp.HasError.Store(true)
		}
	}
}
