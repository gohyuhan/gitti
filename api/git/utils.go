package git

const (
	GENERAL_GIT_UPDATE = "GENERAL_GIT_UPDATE"
	GIT_COMMIT_UPDATE  = "GIT_COMMIT_UPDATE"
)

func GetUpdatedGitInfo(updateChannel chan string) {
	GITFILES.GetGitFilesStatus()
	GITBRANCH.GetLatestBranchesinfo()
	go func() {
		GITCOMMIT.GetLatestGitCommitInfoAndDAG(updateChannel)
	}()

	updateChannel <- GENERAL_GIT_UPDATE
}

// Contains is a generic helper function to check for the existence of an item in a slice.
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
