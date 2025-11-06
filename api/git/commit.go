package git

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"gitti/cmd"
	"gitti/i18n"
)

var GITCOMMIT *GitCommit

const (
	PUSH               = "PUSH"
	FORCEPUSHSAFE      = "FORCEPUSHSAFE"
	FORCEPUSHDANGEROUS = "FORCEPUSHDANGEROUS"
)

type GitCommit struct {
	ErrorLog                          []error
	GitStashAndCommitProcess          *exec.Cmd
	GitRemotePushProcess              *exec.Cmd
	GitAddRemoteProcess               *exec.Cmd
	GitCommitOutput                   []string
	GitRemotePushOutput               []string
	UpdateChannel                     chan string
	GitStashAndCommitProcessMutex     sync.Mutex
	GitRemotePushProcessMutex         sync.Mutex
	GitAddRemoteProcessMutex          sync.Mutex
	isGitStashAndCommitProcessRunning atomic.Bool
	isGitRemotePushProcessRunning     atomic.Bool
	isGitAddRemoteProcessRunning      atomic.Bool
	Remote                            []GitRemote
}

type GitRemote struct {
	Name string
	Url  string
}

func InitGitCommit(updateChannel chan string) {
	gitCommit := GitCommit{
		GitStashAndCommitProcess: nil,
		GitRemotePushProcess:     nil,
		GitAddRemoteProcess:      nil,
		GitCommitOutput:          []string{},
		GitRemotePushOutput:      []string{},
		UpdateChannel:            updateChannel,
		Remote:                   []GitRemote{},
	}

	gitCommit.isGitStashAndCommitProcessRunning.Store(false)
	gitCommit.isGitRemotePushProcessRunning.Store(false)
	gitCommit.isGitAddRemoteProcessRunning.Store(false)

	GITCOMMIT = &gitCommit
}

func (gc *GitCommit) GitFetch() {
	gitArgs := []string{"fetch"}
	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	// Disable interactive prompts for credentials
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	_, err := cmd.Output()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", err))
	}
}

// ----------------------------------
//
//	Related to Git Stash & Commit
//
// ----------------------------------
func (gc *GitCommit) GitStashAndCommit(message, description string) int {
	if !gc.isGitStashAndCommitProcessRunning.CompareAndSwap(false, true) {
		return -1
	}
	gc.ClearGitCommitOutput()
	gc.GitStashAndCommitProcessMutex.Lock()

	defer func() {
		// ensure cleanup even if Start or Wait fails
		gc.GitStashAndCommitProcessMutex.Lock()
		gc.gitStashAndCommitProcessReset()
		gc.gitReset()
		gc.GitStashAndCommitProcessMutex.Unlock()
	}()

	stashGitArgs := []string{"add"}

	// stage selected files.
	GITFILES.GitFilesMutex.Lock()
	for _, files := range GITFILES.FilesStatus {
		if files.SelectedForStage {
			stashGitArgs = append(stashGitArgs, files.FileName)
		}
	}
	GITFILES.GitFilesMutex.Unlock()
	stashCmd := cmd.GittiCmd.RunGitCmd(stashGitArgs)
	gc.GitStashAndCommitProcess = stashCmd

	if len(stashGitArgs) == 1 {
		// No files selected, nothing to stage
		gc.GitStashAndCommitProcessMutex.Unlock()
		return -1
	}

	if _, err := stashCmd.Output(); err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT STAGE ERROR]: %w", err))
		gc.GitStashAndCommitProcessMutex.Unlock()
		return -1
	}

	gitArgs := []string{"commit", "-m", message}
	if len(description) > 0 {
		gitArgs = append(gitArgs, "-m", description)
	}

	commitCmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	gc.GitStashAndCommitProcess = commitCmd
	gc.GitStashAndCommitProcessMutex.Unlock()

	// Combine stderr into stdout
	stdout, err := commitCmd.StdoutPipe()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[PIPE ERROR]: %w", err))
		return -1
	}
	commitCmd.Stderr = commitCmd.Stdout

	if err := commitCmd.Start(); err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[START ERROR]: %w", err))
		return -1
	}

	// Stream combined output
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Text()
			gc.GitCommitOutput = append(gc.GitCommitOutput, line)
			select {
			case gc.UpdateChannel <- GIT_COMMIT_OUTPUT_UPDATE:
			default:
			}
		}
	}()

	waitErr := commitCmd.Wait()
	wg.Wait()

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", waitErr))
			return status
		}
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[UNEXPECTED ERROR]: %w", waitErr))
		return -1
	}

	return 0
}

func (gc *GitCommit) ClearGitCommitOutput() {
	gc.GitCommitOutput = []string{}
}

// This method will not be responsible to set the process state but will be the function that trigger the action will be responsible to reset the status with defer
func (gc *GitCommit) KillGitStashAndCommitCmd() {
	gc.GitStashAndCommitProcessMutex.Lock()
	defer gc.GitStashAndCommitProcessMutex.Unlock()

	if gc.GitStashAndCommitProcess != nil && gc.GitStashAndCommitProcess.Process != nil {
		_ = gc.GitStashAndCommitProcess.Process.Kill()
	}
}

func (gc *GitCommit) gitStashAndCommitProcessReset() {
	gc.GitStashAndCommitProcess = nil
	gc.isGitStashAndCommitProcessRunning.Store(false)
}

// ----------------------------------
//
//	Related to Git Reset
//
// ----------------------------------
func (gc *GitCommit) gitReset() {
	resetCmd := cmd.GittiCmd.RunGitCmd([]string{"reset"})
	if _, err := resetCmd.Output(); err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT RESET ERROR]: %w", err))
	}
}

// func (gc *GitCommit) GitPull() {
//
// }

// ----------------------------------
//
//	Related to Git Push
//
// ----------------------------------
func (gc *GitCommit) GitPush(originName string, pushType string) int {
	if !gc.isGitRemotePushProcessRunning.CompareAndSwap(false, true) {
		return -1
	}
	gc.ClearGitRemotePushOutput()
	gc.GitRemotePushProcessMutex.Lock()
	gitArgs := []string{"push", "-u", originName, GITBRANCH.CurrentCheckOut.BranchName}
	switch pushType {
	case FORCEPUSHSAFE:
		gitArgs = []string{"push", "--force-with-lease", "-u", originName, GITBRANCH.CurrentCheckOut.BranchName}
	case FORCEPUSHDANGEROUS:
		gitArgs = []string{"push", "--force", "-u", originName, GITBRANCH.CurrentCheckOut.BranchName}
	default:
		gitArgs = []string{"push", "-u", originName, GITBRANCH.CurrentCheckOut.BranchName}
	}
	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	// Disable interactive prompts for credentials
	cmd.Env = append(os.Environ(), "GIT_ASKPASS=true", "GIT_TERMINAL_PROMPT=0")

	gc.GitRemotePushProcess = cmd
	gc.GitRemotePushProcessMutex.Unlock()
	defer func() {
		// ensure cleanup even if Start or Wait fails
		gc.GitRemotePushProcessMutex.Lock()
		gc.resetGitRemotePushProcesstatus()
		gc.GitRemotePushProcessMutex.Unlock()
	}()

	// Combine stderr into stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[PIPE ERROR]: %w", err))
		return -1
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[START ERROR]: %w", err))
		return -1
	}

	// Stream combined output
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Text()
			gc.GitRemotePushOutput = append(gc.GitRemotePushOutput, line)
			select {
			case gc.UpdateChannel <- GIT_REMOTE_PUSH_OUTPUT_UPDATE:
			default:
			}
		}
	}()

	waitErr := cmd.Wait()
	wg.Wait()

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT PUSH ERROR]: %w", waitErr))
			return status
		}
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[UNEXPECTED ERROR]: %w", waitErr))
		return -1
	}
	return 0
}

func (gc *GitCommit) ClearGitRemotePushOutput() {
	gc.GitRemotePushOutput = []string{}
}

// This method will not be responsible to set the process state but will be the function that trigger the action will be responsible to reset the status with defer
func (gc *GitCommit) KillGitRemotePushCmd() {
	gc.GitRemotePushProcessMutex.Lock()
	defer gc.GitRemotePushProcessMutex.Unlock()

	if gc.GitRemotePushProcess != nil && gc.GitRemotePushProcess.Process != nil {
		_ = gc.GitRemotePushProcess.Process.Kill()
	}
}

func (gc *GitCommit) resetGitRemotePushProcesstatus() {
	gc.GitRemotePushProcess = nil
	gc.isGitRemotePushProcessRunning.Store(false)
}

// ----------------------------------
//
//	Related to Git Remote
//
// ----------------------------------
func (gc *GitCommit) GitAddRemote(originName string, url string) ([]string, int) {
	if !gc.isGitAddRemoteProcessRunning.CompareAndSwap(false, true) {
		return []string{}, -1
	}
	gc.GitAddRemoteProcessMutex.Lock()
	gitArgs := []string{"remote", "add", originName, url}
	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)

	gc.GitAddRemoteProcess = cmd
	gc.GitAddRemoteProcessMutex.Unlock()
	defer func() {
		gc.GitAddRemoteProcessMutex.Lock()
		gc.gitAddRemoteProcessReset()
		gc.GitAddRemoteProcessMutex.Unlock()
	}()
	if !isValidGitRemoteURL(url) {
		return []string{i18n.LANGUAGEMAPPING.AddRemotePopUpInvalidRemoteUrlFormat}, -1
	}

	gitOutput, err := cmd.CombinedOutput()

	gitAddRemoteOutput := strings.Split(strings.TrimSpace(string(gitOutput)), "\n")
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT ADD REMOTE ERROR]: %w", err))
			return gitAddRemoteOutput, status
		}
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[UNEXPECTED ERROR]: %w", err))
		return gitAddRemoteOutput, -1

	}
	return gitAddRemoteOutput, 0
}

// KillGitAddRemoteCmd forcefully terminates any running git remote add process.
// It is safe to call this method even if no process is running.
// This method is thread-safe and can be called from multiple goroutines.
// This method will not be responsible to set the process state but will be the function that trigger the action will be responsible to reset the status with defer
func (gc *GitCommit) KillGitAddRemoteCmd() {
	gc.GitAddRemoteProcessMutex.Lock()
	defer gc.GitAddRemoteProcessMutex.Unlock()

	if gc.GitAddRemoteProcess != nil && gc.GitAddRemoteProcess.Process != nil {
		_ = gc.GitAddRemoteProcess.Process.Kill()
	}
}

func (gc *GitCommit) gitAddRemoteProcessReset() {
	gc.GitAddRemoteProcess = nil
	gc.isGitAddRemoteProcessRunning.Store(false)
}

func (gc *GitCommit) CheckRemoteExist() bool {
	gitArgs := []string{"remote", "-v"}
	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	gitOutput, err := cmd.Output()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", err))
	}
	remotes := strings.SplitSeq(strings.TrimSpace(string(gitOutput)), "\n")
	var remoteStruct []GitRemote
	for remote := range remotes {
		if !strings.HasSuffix(remote, "(push)") {
			continue
		}
		remoteLinePart := strings.Fields(remote)
		if len(remoteLinePart) < 2 {
			continue
		}
		remoteStruct = append(remoteStruct, GitRemote{
			Name: remoteLinePart[0],
			Url:  remoteLinePart[1],
		})
	}
	gc.Remote = remoteStruct
	return len(gc.Remote) > 0
}

// ----------------------------------
//
//	Related to Git Init
//
// ----------------------------------
func GitInit(repoPath string) {
	gitArgs := []string{"init"}

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("[GIT INIT ERROR]: %v", err)
		os.Exit(1)
	}
}

// check if the format for git remote is correct and valid
func isValidGitRemoteURL(remote string) bool {
	// Check HTTPS style
	if strings.HasPrefix(remote, "https://") || strings.HasPrefix(remote, "http://") {
		_, err := url.ParseRequestURI(remote)
		return err == nil
	}

	// Check SSH style (e.g. git@github.com:user/repo.git)
	sshPattern := `^[\w.-]+@[\w.-]+:[\w./-]+(\.git)?$`
	matched, _ := regexp.MatchString(sshPattern, remote)
	return matched
}
