package cmd

import "os/exec"

type Cmd struct {
	repoPath string
}

var GittiCmd *Cmd

func InitCmd(repoPath string) {
	GittiCmd = &Cmd{
		repoPath: repoPath,
	}
}

func (c *Cmd) RunGitCmd(gitArgs []string) *exec.Cmd {
	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = c.repoPath

	return cmd
}
