package services

import "github.com/gohyuhan/gitti/tui/types"

func GitStateUniversalUtilsContinueService(m *types.GittiModel) {
	go func() {
		m.GitOperations.GitStateUniversalUtils.GitUniversalContinue()
	}()
}

func GitStateUniversalUtilsAbortService(m *types.GittiModel) {
	go func() {
		m.GitOperations.GitStateUniversalUtils.GitUniversalAbort()
	}()
}

func GitStateUniversalUtilsSkipService(m *types.GittiModel) {
	go func() {
		m.GitOperations.GitStateUniversalUtils.GitUniversalSkip()
	}()
}
