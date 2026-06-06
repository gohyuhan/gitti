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
	GitBlame               *git.GitBlame
	GitInteractiveRebase   *git.GitInteractiveRebase
	GitWorktree            *git.GitWorktree
}

type GitRepoPath struct {
	// having both these path is to support submodule
	AbsoluteGitRepoPath  string // this is the most root level path where .git folder is located
	RepoMainGitDirPath   string // the common (main) git dir; for linked worktrees this is the shared .git, used as the file-watch root so worktree add/remove is observed
	TopLevelRepoPath     string // this is the path where the top level .git file/folder is located at
	AbsoluteWorktreePath string // the absolute path for the current repo worktree
	RepoName             string
}
