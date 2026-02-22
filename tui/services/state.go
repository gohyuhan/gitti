package services

import "github.com/gohyuhan/gitti/tui/types"

// ------------------------------------
//
//	Continue the current in-progress git operation (rebase, merge, cherry-pick, etc.)
//
// ------------------------------------
func GitStateUniversalUtilsContinueService(m *types.GittiModel) {
	go func() {
		m.GitOperations.GitStateUniversalUtils.GitUniversalContinue()
	}()
}

// ------------------------------------
//
//	Abort the current in-progress git operation
//
// ------------------------------------
func GitStateUniversalUtilsAbortService(m *types.GittiModel) {
	go func() {
		m.GitOperations.GitStateUniversalUtils.GitUniversalAbort()
	}()
}

// ------------------------------------
//
//	Skip the current step of the in-progress git operation
//
// ------------------------------------
func GitStateUniversalUtilsSkipService(m *types.GittiModel) {
	go func() {
		m.GitOperations.GitStateUniversalUtils.GitUniversalSkip()
	}()
}
