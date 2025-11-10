package api

import "gitti/api/git"

type GitState struct {
	GitBranch *git.GitBranch
	GitCommit *git.GitCommit
	GitFiles  *git.GitFiles
	GitPull   *git.GitPull
}
