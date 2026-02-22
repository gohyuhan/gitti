package services

import (
	"regexp"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
)

// services was to bridge api and the needs of the terminal interface logic so that it can be compatible and feels smooth and not clunky
// ------------------------------------
//
//	For Git Individual file stage or unstage
//
// ------------------------------------
func GitStageOrUnstageService(m *types.GittiModel, filePathName string) {
	go func() {
		m.GitOperations.GitFiles.StageOrUnstageFile(filePathName)
	}()
}

// ------------------------------------
//
//	For Git Stage All
//
// ------------------------------------
func GitStageAllChangesService(m *types.GittiModel) {
	go func() {
		m.GitOperations.GitFiles.StageAllChanges()
	}()
}

// ------------------------------------
//
//	For Git Unstage All
//
// ------------------------------------
func GitUnstageAllChangesService(m *types.GittiModel) {
	go func() {
		m.GitOperations.GitFiles.UnstageAllChanges()
	}()
}

// ------------------------------------
//
//	Strip ANSI escape codes from a string array for clean text processing
//
// ------------------------------------
func stripAnsi(strArray []string) []string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntgry=><~]))"
	re := regexp.MustCompile(ansi)
	for i, str := range strArray {
		strArray[i] = re.ReplaceAllString(str, "")
	}
	return strArray
}

// ------------------------------------
//
//	Stage or unstage a single diff line for the selected file
//
// ------------------------------------
func GitStageLineOrUnstageLineService(m *types.GittiModel, filePathName string) {
	switch m.CurrentSelectedComponent {
	case constant.DetailComponentPanel:
		contentArray := stripAnsi(m.DetailPanelViewportOGStringArray)
		overflowIndexCount := m.LineEditingIndexPositionAndInfo.DetailPanelViewportOverflowIndexCount
		actualCurrentIndex := m.LineEditingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex

		switch m.LineEditingIndexPositionAndInfo.DetailPanelViewportStageType {
		case constant.STAGE:
			go func() {
				m.GitOperations.GitFiles.UnstageLine(filePathName, contentArray, overflowIndexCount, actualCurrentIndex)
			}()
		case constant.UNSTAGE:
			go func() {
				m.GitOperations.GitFiles.StageLine(filePathName, contentArray, overflowIndexCount, actualCurrentIndex)
			}()
		default:
			return
		}
	case constant.DetailComponentPanelTwo:
		contentArray := stripAnsi(m.DetailPanelTwoViewportOGStringArray)
		overflowIndexCount := m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportOverflowIndexCount
		actualCurrentIndex := m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportActualCurrentIndex

		switch m.LineEditingIndexPositionAndInfo.DetailPanelTwoViewportStageType {
		case constant.STAGE:
			go func() {
				m.GitOperations.GitFiles.UnstageLine(filePathName, contentArray, overflowIndexCount, actualCurrentIndex)
			}()
		case constant.UNSTAGE:
			go func() {
				m.GitOperations.GitFiles.StageLine(filePathName, contentArray, overflowIndexCount, actualCurrentIndex)
			}()
		default:
			return
		}
	}
}
