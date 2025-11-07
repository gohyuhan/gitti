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

const (
	PUSH               = "PUSH"
	FORCEPUSHSAFE      = "FORCEPUSHSAFE"
	FORCEPUSHDANGEROUS = "FORCEPUSHDANGEROUS"
)

type GitCommit struct {
	errorLog                          []error
	gitStashAndCommitProcess          *exec.Cmd
	gitRemotePushProcess              *exec.Cmd
	gitAddRemoteProcess               *exec.Cmd
	gitCommitOutput                   []string
	gitRemotePushOutput               []string
	updateChannel                     chan string
	gitStashAndCommitProcessMutex     sync.Mutex
	gitRemotePushProcessMutex         sync.Mutex
	gitAddRemoteProcessMutex          sync.Mutex
	isGitStashAndCommitProcessRunning atomic.Bool
	isGitRemotePushProcessRunning     atomic.Bool
	isGitAddRemoteProcessRunning      atomic.Bool
	remote                            []GitRemote
}

type GitRemote struct {
	Name string
	Url  string
}

func InitGitCommit(updateChannel chan string) *GitCommit {
	gitCommit := GitCommit{
		gitStashAndCommitProcess: nil,
		gitRemotePushProcess:     nil,
		gitAddRemoteProcess:      nil,
		gitCommitOutput:          []string{},
		gitRemotePushOutput:      []string{},
		updateChannel:            updateChannel,
		remote:                   []GitRemote{},
	}

	gitCommit.isGitStashAndCommitProcessRunning.Store(false)
	gitCommit.isGitRemotePushProcessRunning.Store(false)
	gitCommit.isGitAddRemoteProcessRunning.Store(false)

	return &gitCommit
}

// ----------------------------------
//
//	Return git commit output
//
// ----------------------------------
func (gc *GitCommit) GitCommitOutput() []string {
	return gc.gitCommitOutput
}

// ----------------------------------
//
//	Return remote
//
// ----------------------------------
func (gc *GitCommit) Remote() []GitRemote {
	return gc.remote
}

// ----------------------------------
//
//	Return git remote push output
//
// ----------------------------------
func (gc *GitCommit) GitRemotePushOutput() []string {
	return gc.gitRemotePushOutput
}

func (gc *GitCommit) GitFetch() {
	gitArgs := []string{"fetch"}
	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	// Disable interactive prompts for credentials
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	_, err := cmd.Output()
	if err != nil {
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", err))
	}
}

// ----------------------------------
//
//	Related to Git Stage & Commit
//
// ----------------------------------
func (gc *GitCommit) GitStageAndCommit(message, description string, stagedFiles []string) int {
	if !gc.isGitStashAndCommitProcessRunning.CompareAndSwap(false, true) {
		return -1
	}
	gc.ClearGitCommitOutput()

	defer func() {
		// ensure cleanup even if Start or Wait fails
		gc.gitStashAndCommitProcessMutex.Lock()
		gc.gitStageAndCommitProcessReset()
		gc.gitReset()
		gc.gitStashAndCommitProcessMutex.Unlock()
	}()

	gc.gitStashAndCommitProcessMutex.Lock()

	stageGitArgs := []string{"add"}
	if len(stagedFiles) > 0 {
		stageGitArgs = append(stageGitArgs, stagedFiles...)
	}

	stageCmd := cmd.GittiCmd.RunGitCmd(stageGitArgs)
	gc.gitStashAndCommitProcess = stageCmd

	if len(stageGitArgs) == 1 {
		// No files selected, nothing to stage
		gc.gitStashAndCommitProcessMutex.Unlock()
		return -1
	}

	if _, err := stageCmd.Output(); err != nil {
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[GIT STAGE ERROR]: %w", err))
		gc.gitStashAndCommitProcessMutex.Unlock()
		return -1
	}

	gitArgs := []string{"commit", "-m", message}
	if len(description) > 0 {
		gitArgs = append(gitArgs, "-m", description)
	}

	commitCmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	gc.gitStashAndCommitProcess = commitCmd

	// Combine stderr into stdout
	stdout, err := commitCmd.StdoutPipe()
	if err != nil {
		gc.gitStashAndCommitProcessMutex.Unlock()
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[PIPE ERROR]: %w", err))
		return -1
	}
	commitCmd.Stderr = commitCmd.Stdout

	// Start the process while still holding the mutex
	if err := commitCmd.Start(); err != nil {
		gc.gitStashAndCommitProcessMutex.Unlock()
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[START ERROR]: %w", err))
		return -1
	}

	// Process is now started and can be killed safely
	gc.gitStashAndCommitProcessMutex.Unlock()

	// Stream combined output
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Text()
			gc.gitCommitOutput = append(gc.gitCommitOutput, line)
			select {
			case gc.updateChannel <- GIT_COMMIT_OUTPUT_UPDATE:
			default:
			}
		}
	}()

	waitErr := commitCmd.Wait()
	wg.Wait()

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gc.errorLog = append(gc.errorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", waitErr))
			return status
		}
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[UNEXPECTED ERROR]: %w", waitErr))
		return -1
	}

	return 0
}

func (gc *GitCommit) ClearGitCommitOutput() {
	gc.gitCommitOutput = []string{}
}

// This method will not be responsible to set the process state but will be the function that trigger the action will be responsible to reset the status with defer
func (gc *GitCommit) KillGitStageAndCommitCmd() {
	gc.gitStashAndCommitProcessMutex.Lock()
	defer gc.gitStashAndCommitProcessMutex.Unlock()

	if gc.gitStashAndCommitProcess != nil && gc.gitStashAndCommitProcess.Process != nil {
		_ = gc.gitStashAndCommitProcess.Process.Kill()
	}
}

func (gc *GitCommit) gitStageAndCommitProcessReset() {
	gc.gitStashAndCommitProcess = nil
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
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[GIT RESET ERROR]: %w", err))
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
func (gc *GitCommit) GitPush(currentCheckOutBranch string, originName string, pushType string) int {
	if !gc.isGitRemotePushProcessRunning.CompareAndSwap(false, true) {
		return -1
	}
	gc.ClearGitRemotePushOutput()

	defer func() {
		// ensure cleanup even if Start or Wait fails
		gc.gitRemotePushProcessMutex.Lock()
		gc.resetGitRemotePushProcesstatus()
		gc.gitRemotePushProcessMutex.Unlock()
	}()

	gc.gitRemotePushProcessMutex.Lock()
	gitArgs := []string{"push", "-u", originName, currentCheckOutBranch}
	switch pushType {
	case FORCEPUSHSAFE:
		gitArgs = []string{"push", "--force-with-lease", "-u", originName, currentCheckOutBranch}
	case FORCEPUSHDANGEROUS:
		gitArgs = []string{"push", "--force", "-u", originName, currentCheckOutBranch}
	default:
		gitArgs = []string{"push", "-u", originName, currentCheckOutBranch}
	}
	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	// Disable interactive prompts for credentials
	cmd.Env = append(os.Environ(), "GIT_ASKPASS=true", "GIT_TERMINAL_PROMPT=0")

	gc.gitRemotePushProcess = cmd

	// Combine stderr into stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		gc.gitRemotePushProcessMutex.Unlock()
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[PIPE ERROR]: %w", err))
		return -1
	}
	cmd.Stderr = cmd.Stdout

	// Start the process while still holding the mutex
	if err := cmd.Start(); err != nil {
		gc.gitRemotePushProcessMutex.Unlock()
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[START ERROR]: %w", err))
		return -1
	}

	// Process is now started and can be killed safely
	gc.gitRemotePushProcessMutex.Unlock()

	// Stream combined output
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Text()
			gc.gitRemotePushOutput = append(gc.gitRemotePushOutput, line)
			select {
			case gc.updateChannel <- GIT_REMOTE_PUSH_OUTPUT_UPDATE:
			default:
			}
		}
	}()

	waitErr := cmd.Wait()
	wg.Wait()

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gc.errorLog = append(gc.errorLog, fmt.Errorf("[GIT PUSH ERROR]: %w", waitErr))
			return status
		}
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[UNEXPECTED ERROR]: %w", waitErr))
		return -1
	}
	return 0
}

func (gc *GitCommit) ClearGitRemotePushOutput() {
	gc.gitRemotePushOutput = []string{}
}

// This method will not be responsible to set the process state but will be the function that trigger the action will be responsible to reset the status with defer
func (gc *GitCommit) KillGitRemotePushCmd() {
	gc.gitRemotePushProcessMutex.Lock()
	defer gc.gitRemotePushProcessMutex.Unlock()

	if gc.gitRemotePushProcess != nil && gc.gitRemotePushProcess.Process != nil {
		_ = gc.gitRemotePushProcess.Process.Kill()
	}
}

func (gc *GitCommit) resetGitRemotePushProcesstatus() {
	gc.gitRemotePushProcess = nil
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
	defer func() {
		gc.gitAddRemoteProcessMutex.Lock()
		gc.gitAddRemoteProcessReset()
		gc.gitAddRemoteProcessMutex.Unlock()
	}()

	if !isValidGitRemoteURL(url) {
		errMsg := "Invalid remote URL format"
		if i18n.LANGUAGEMAPPING != nil {
			errMsg = i18n.LANGUAGEMAPPING.AddRemotePopUpInvalidRemoteUrlFormat
		}
		return []string{errMsg}, -1
	}

	gc.gitAddRemoteProcessMutex.Lock()
	gitArgs := []string{"remote", "add", originName, url}
	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	gc.gitAddRemoteProcess = cmd

	// CombinedOutput starts and waits for the command
	gitOutput, err := cmd.CombinedOutput()
	gc.gitAddRemoteProcessMutex.Unlock()

	gitAddRemoteOutput := strings.Split(strings.TrimSpace(string(gitOutput)), "\n")
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gc.errorLog = append(gc.errorLog, fmt.Errorf("[GIT ADD REMOTE ERROR]: %w", err))
			return gitAddRemoteOutput, status
		}
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[UNEXPECTED ERROR]: %w", err))
		return gitAddRemoteOutput, -1

	}
	return gitAddRemoteOutput, 0
}

// KillGitAddRemoteCmd forcefully terminates any running git remote add process.
// It is safe to call this method even if no process is running.
// This method is thread-safe and can be called from multiple goroutines.
// This method will not be responsible to set the process state but will be the function that trigger the action will be responsible to reset the status with defer
func (gc *GitCommit) KillGitAddRemoteCmd() {
	gc.gitAddRemoteProcessMutex.Lock()
	defer gc.gitAddRemoteProcessMutex.Unlock()

	if gc.gitAddRemoteProcess != nil && gc.gitAddRemoteProcess.Process != nil {
		_ = gc.gitAddRemoteProcess.Process.Kill()
	}
}

func (gc *GitCommit) gitAddRemoteProcessReset() {
	gc.gitAddRemoteProcess = nil
	gc.isGitAddRemoteProcessRunning.Store(false)
}

func (gc *GitCommit) CheckRemoteExist() bool {
	gitArgs := []string{"remote", "-v"}
	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	gitOutput, err := cmd.Output()
	if err != nil {
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", err))
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
	gc.remote = remoteStruct
	return len(gc.remote) > 0
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
