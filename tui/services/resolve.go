package services

import "github.com/gohyuhan/gitti/tui/types"

// ------------------------------------
//
//	For resolving git file conflict
//
// ------------------------------------
func GitResolveConflictService(m *types.GittiModel, filePathName string, resolveType string) {
	go func() {
		m.GitOperations.GitFiles.GitResolveConflict(filePathName, resolveType)
	}()
}
