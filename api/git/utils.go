package git

const (
	GENERAL_GIT_UPDATE      = "GENERAL_GIT_UPDATE"
	GIT_COMMIT_UPDATE       = "GIT_COMMIT_UPDATE"
	GIT_FILES_STATUS_UPDATE = "GIT_FILES_STATUS_UPDATE"
)

func GetUpdatedGitInfo(updateChannel chan string) {
	GITFILES.GetGitFilesStatus()
	GITBRANCH.GetLatestBranchesinfo()

	// not included in v0.1.0
	// go func() {
	// 	GITCOMMIT.GetLatestGitCommitLogInfoAndDAG(updateChannel)
	// }()

	updateChannel <- GENERAL_GIT_UPDATE
}
