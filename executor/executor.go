package executor

import (
	"context"
	"os/exec"
)

type CmdExecutor struct {
	repoPath string
}

var GittiCmdExecutor *CmdExecutor

func InitCmdExecutor(repoPath string) {
	GittiCmdExecutor = &CmdExecutor{
		repoPath: repoPath,
	}
}

func (c *CmdExecutor) RunGitCmd(gitArgs []string, colorized bool) *exec.Cmd {
	if colorized {
		gitArgs = append([]string{"-c", "color.ui=always"}, gitArgs...)
	}
	gitArgs = append([]string{"--no-optional-locks"}, gitArgs...)
	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = c.repoPath

	return cmd
}

// RunGitCmdWithContext creates a git command that respects context cancellation.
// When the context is cancelled, the command will be terminated automatically.
func (c *CmdExecutor) RunGitCmdWithContext(ctx context.Context, gitArgs []string, colorized bool) *exec.Cmd {
	if colorized {
		gitArgs = append([]string{"-c", "color.ui=always"}, gitArgs...)
	}
	gitArgs = append([]string{"--no-optional-locks"}, gitArgs...)
	cmd := exec.CommandContext(ctx, "git", gitArgs...)
	cmd.Dir = c.repoPath

	return cmd
}

func (c *CmdExecutor) UpdateRepoPath(updatedRepoPath string) {
	c.repoPath = updatedRepoPath
}
