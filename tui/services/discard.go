package services

import (
	"github.com/gohyuhan/gitti/tui/types"
)

// services was to bridge api and the needs of the terminal interface logic so that it can be compatible and feels smooth and not clunky
// ------------------------------------
//
//	For Git discard file changes
//
// ------------------------------------
func GitDiscardFileChangesService(m *types.GittiModel, filePathName string, discardType string) {
	go func() {
		m.GitOperations.GitFiles.DiscardFileChanges(filePathName, discardType)
	}()
}
