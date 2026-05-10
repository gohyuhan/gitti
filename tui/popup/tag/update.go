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

func UpdateDeleteTagResultEvent(m *types.GittiModel, updateData types.GitDeleteTagResultEventDataInterface) {
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

func UpdatePushTagResultEvent(m *types.GittiModel, updateData types.GitPushTagResultEventDataInterface) {
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

func UpdateFetchTagResultEvent(m *types.GittiModel, updateData types.GitFetchTagResultEventDataInterface) {
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
