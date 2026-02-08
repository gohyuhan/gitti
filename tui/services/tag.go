package services

import "github.com/gohyuhan/gitti/tui/types"

func CreateNewTagService(m *types.GittiModel, commitHash string, tagName string, tagMessage string) {
	go func() {
		m.GitOperations.GitTag.CreateNewTag(commitHash, tagName, tagMessage)
	}()
}
