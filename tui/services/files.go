package services

import (
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
)

func GitDiscardLineFileChangeService(m *types.GittiModel, filePathName string) {
	switch m.CurrentSelectedComponent {
	case constant.DetailComponentPanel:
		contentArray := stripAnsi(m.DetailPanelViewportOGStringArray)
		overflowIndexCount := m.LineEditingIndexPositionAndInfo.DetailPanelViewportOverflowIndexCount
		actualCurrentIndex := m.LineEditingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex

		switch m.LineEditingIndexPositionAndInfo.DetailPanelViewportStageType {
		case constant.STAGE:
			go func() {
				m.GitOperations.GitFiles.GitDiscardFileLineChange(filePathName, contentArray, overflowIndexCount, actualCurrentIndex, git.STAGE)
			}()
		case constant.UNSTAGE:
			go func() {
				m.GitOperations.GitFiles.GitDiscardFileLineChange(filePathName, contentArray, overflowIndexCount, actualCurrentIndex, git.UNSTAGE)
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
				m.GitOperations.GitFiles.GitDiscardFileLineChange(filePathName, contentArray, overflowIndexCount, actualCurrentIndex, git.STAGE)
			}()
		case constant.UNSTAGE:
			go func() {
				m.GitOperations.GitFiles.GitDiscardFileLineChange(filePathName, contentArray, overflowIndexCount, actualCurrentIndex, git.UNSTAGE)
			}()
		default:
			return
		}
	}
}
