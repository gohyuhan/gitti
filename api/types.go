package api

import "github.com/gohyuhan/gitti/api/git"

type GitOperations struct {
	GitBranch              *git.GitBranch
	GitCommit              *git.GitCommit
	GitFiles               *git.GitFiles
	GitPull                *git.GitPull
	GitRebase              *git.GitRebase
	GitStash               *git.GitStash
	GitRemote              *git.GitRemote
	GitCommitLog           *git.GitCommitLog
	GitRefLog              *git.GitRefLog
	GitTag                 *git.GitTag
	GitStateUniversalUtils *git.GitStateUniversalUtils
}

type GitRepoPath struct {
	// having both these path is to support submodule
	AbsoluteGitRepoPath string // this is the most root level path where .git folder is located
	TopLevelRepoPath    string // this is the path where the top level .git file/folder is located at
	RepoName            string
}
