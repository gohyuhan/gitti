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

func (c *Cmd) RunGitCmd(gitArgs []string, colorized bool) *exec.Cmd {
	if colorized {
		gitArgs = append([]string{"-c", "color.ui=always"}, gitArgs...)
	}
	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = c.repoPath

	return cmd
}
