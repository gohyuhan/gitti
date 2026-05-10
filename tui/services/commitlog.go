package services

import (
	"slices"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Cherry-pick selected commits in the user-specified sequence order (cherry picked from commit logs)
//
// ------------------------------------
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
		m.TuiUpdateChannel <- types.GittiTuiUpdateMsg{
			Event: constant.REINIT_CHERRY_PICKED_COMMIT_INFO_EVENT,
		}
	}()

}

// ------------------------------------
//
//	Cherry-pick selected reflog hash (cherry picked from reflogs)
//
// ------------------------------------
func GitCherryPickReflogHashService(m *types.GittiModel, cherryPickedRefLogHash string) {
	go func() {
		m.GitOperations.GitCommitLog.GitCherryPick([]string{cherryPickedRefLogHash})
	}()
}

// ------------------------------------
//
//	Revert a specific commit by its hash
//
// ------------------------------------
func GitRevertCommitService(m *types.GittiModel, commitHash string, parentOrder int) {
	go func() {
		m.GitOperations.GitCommitLog.GitRevertCommit(commitHash, parentOrder)
	}()
}

// ------------------------------------
//
//	Return the parent info for a specific commit hash
//
// ------------------------------------
func GetCommitHashParentInfoService(m *types.GittiModel, commitHash string) []git.CommitHashParentInfo {
	return m.GitOperations.GitCommitLog.GetCommitHashParentInfo(commitHash)
}
