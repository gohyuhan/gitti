package services

import (
	"slices"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/types"
)

func GitCherryPickService(m *types.GittiModel, cherryPickedCommitLogs map[string]git.CherryPickedCommitLog) {
	go func() {
		var sortedCherryPickedCommitLogs []git.CherryPickedCommitLog
		// turn the hashmap into array first
		for _, commitLogItem := range cherryPickedCommitLogs {
			sortedCherryPickedCommitLogs = append(sortedCherryPickedCommitLogs, commitLogItem)
		}

		// sort the array based on user selection sequence
		slices.SortFunc(sortedCherryPickedCommitLogs, func(a, b git.CherryPickedCommitLog) int {
			return a.UserSelectedSequence - b.UserSelectedSequence
		})

		// harvest the commit hash
		var cherryPickedCommitHashes []string
		for _, commitLog := range sortedCherryPickedCommitLogs {
			cherryPickedCommitHashes = append(cherryPickedCommitHashes, commitLog.Hash)
		}

		m.GitOperations.GitCommitLog.GitCherryPick(cherryPickedCommitHashes)
	}()
}
