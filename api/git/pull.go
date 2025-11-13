package git

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
	"sync/atomic"

	"gitti/cmd"
)

type GitPull struct {
	errorLog                []error
	gitPullOutput           []string
	gitPullProcessMutex     sync.Mutex
	gitPullProcessCmd       *exec.Cmd
	isGitPullProcessRunning atomic.Bool
	updateChannel           chan string
}

func InitGitPull(updateChannel chan string) *GitPull {
	gitPull := &GitPull{
		errorLog:      []error{},
		gitPullOutput: []string{},
		updateChannel: updateChannel,
	}
	gitPull.gitPullProcessCmd = nil
	gitPull.isGitPullProcessRunning.Store(false)

	return gitPull
}

// --------------------------------
//
// return the git pull output
//
// --------------------------------
func (gp *GitPull) GetGitPullOutput() []string {
	return gp.gitPullOutput
}

// --------------------------------
//
// # Git Pull and will operate differently based on the user selection type
//
// --------------------------------
func (gp *GitPull) GitPull(pullType string) int {
	if !gp.isGitPullProcessRunning.CompareAndSwap(false, true) {
		return -1
	}
	defer func() {
		gp.gitPullProcessMutex.Lock()
		gp.gitPullProcessReset()
		gp.gitPullProcessMutex.Unlock()
	}()

	gp.ClearGitPullOutput()
	gp.gitPullProcessMutex.Lock()
	var gitPullArgs []string
	switch pullType {
	case GITPULL:
		gitPullArgs = []string{"pull", "--no-edit"}
	case GITPULLREBASE:
		gitPullArgs = []string{"pull", "--rebase", "--autostash", "--no-edit"}
	case GITPULLMERGE:
		gitPullArgs = []string{"pull", "--no-rebase", "--no-edit"}
	}

	cmd := cmd.GittiCmd.RunGitCmd(gitPullArgs, true)

	gp.gitPullProcessCmd = cmd

	// Combine stderr into stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		gp.gitPullProcessMutex.Unlock()
		gp.errorLog = append(gp.errorLog, fmt.Errorf("[PIPE ERROR]: %w", err))
		return -1
	}
	cmd.Stderr = cmd.Stdout

	// Start the process while still holding the mutex
	if err := cmd.Start(); err != nil {
		gp.gitPullProcessMutex.Unlock()
		gp.errorLog = append(gp.errorLog, fmt.Errorf("[START ERROR]: %w", err))
		return -1
	}

	// Process is now started and can be killed safely
	gp.gitPullProcessMutex.Unlock()

	// Stream combined output
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Text()
			gp.gitPullOutput = append(gp.gitPullOutput, line)
			select {
			case gp.updateChannel <- GIT_PULL_OUTPUT_UPDATE:
			default:
			}
		}
	}()

	waitErr := cmd.Wait()
	wg.Wait()

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gp.errorLog = append(gp.errorLog, fmt.Errorf("[GIT PULL ERROR]: %w", waitErr))
			return status
		}
		gp.errorLog = append(gp.errorLog, fmt.Errorf("[UNEXPECTED ERROR]: %w", waitErr))
		return -1
	}
	return 0
}

// --------------------------------
//
// # Clear the Git Process Output
//
// --------------------------------
func (gp *GitPull) ClearGitPullOutput() {
	gp.gitPullOutput = []string{}
}

// --------------------------------
//
// # Kill the git pull process
//
// --------------------------------
func (gp *GitPull) KillGitPullCmd() {
	gp.gitPullProcessMutex.Lock()
	defer gp.gitPullProcessMutex.Unlock()

	if gp.gitPullProcessCmd != nil && gp.gitPullProcessCmd.Process != nil {
		_ = gp.gitPullProcessCmd.Process.Kill()
	}
}

// --------------------------------
//
// # Reset the git pull process after done
//
// --------------------------------
func (gp *GitPull) gitPullProcessReset() {
	gp.gitPullProcessCmd = nil
	gp.isGitPullProcessRunning.Store(false)
}
