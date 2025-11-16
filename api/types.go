package api

import "gitti/api/git"

type GitState struct {
	GitBranch *git.GitBranch
	GitCommit *git.GitCommit
	GitFiles  *git.GitFiles
	GitPull   *git.GitPull
	GitStash  *git.GitStash
}

type GitRepoPath struct {
	// having both these path is to support submodule
	AbsoluteGitRepoPath string // this is the most root level path where .git folder is located
	TopLevelRepoPath    string // this is the path where the top level .git file/folder is located at
	RepoName            string
}
