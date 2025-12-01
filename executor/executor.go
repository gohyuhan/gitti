package executor

import "os/exec"

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

func (c *CmdExecutor) UpdateRepoPath(updatedRepoPath string) {
	c.repoPath = updatedRepoPath
}
