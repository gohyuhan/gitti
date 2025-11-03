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

	"gitti/i18n"
)

var GITCOMMIT *GitCommit

type GitCommit struct {
	RepoPath                  string
	ErrorLog                  []error
	GitCommitProcess          *exec.Cmd
	GitRemotePushProcess      *exec.Cmd
	GitAddRemoteProcess       *exec.Cmd
	GitStageProcess           *exec.Cmd
	GitCommitOutput           []string
	GitRemotePushOutput       []string
	UpdateChannel             chan string
	GitCommitProcessMutex     sync.Mutex
	GitRemotePushProcessMutex sync.Mutex
	GitAddRemoteProcessMutex  sync.Mutex
	GitStageProcessMutex      sync.Mutex
	Remote                    []GitRemote
}

type GitRemote struct {
	Name string
	Url  string
}

func InitGitCommit(repoPath string, updateChannel chan string) {
	gitCommit := GitCommit{
		RepoPath:             repoPath,
		GitCommitProcess:     nil,
		GitRemotePushProcess: nil,
		GitAddRemoteProcess:  nil,
		GitStageProcess:      nil,
		GitCommitOutput:      []string{},
		GitRemotePushOutput:  []string{},
		UpdateChannel:        updateChannel,
		Remote:               []GitRemote{},
	}
	GITCOMMIT = &gitCommit
}

func (gc *GitCommit) GitFetch() {
	gitArgs := []string{"fetch"}
	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gc.RepoPath
	// Disable interactive prompts for credentials
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	_, err := cmd.Output()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", err))
	}
}

// ----------------------------------
//
//	Related to Git Stage
//
// ----------------------------------
func (gc *GitCommit) GitStage() {
	// First, reset the staging area to avoid conflicts with previously staged files.
	// Also prevent to stage file that was stage else where like on git cli or other git related program
	// but decided to not include that file for commit when using gitti
	resetCmd := exec.Command("git", "reset")
	resetCmd.Dir = gc.RepoPath
	if _, err := resetCmd.Output(); err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT RESET ERROR]: %w", err))
	}

	// Now, stage selected files.
	gitArgs := []string{"add"}
	for _, files := range GITFILES.FilesStatus {
		if files.SelectedForStage {
			gitArgs = append(gitArgs, files.FileName)
		}
	}

	if len(gitArgs) == 1 {
		// No files selected, nothing to stage
		return
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gc.RepoPath

	gc.GitStageProcessMutex.Lock()
	gc.GitStageProcess = cmd
	gc.GitStageProcessMutex.Unlock()
	defer func() {
		gc.GitStageProcessMutex.Lock()
		gc.GitStageProcess = nil
		gc.GitStageProcessMutex.Unlock()
	}()

	if _, err := cmd.Output(); err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT STAGE ERROR]: %w", err))
	}
}

// KillGitStageCmd forcefully terminates any running git stage process.
// It is safe to call this method even if no process is running.
// This method is thread-safe and can be called from multiple goroutines.
func (gc *GitCommit) KillGitStageCmd() {
	gc.GitStageProcessMutex.Lock()
	defer gc.GitStageProcessMutex.Unlock()

	if gc.GitStageProcess != nil && gc.GitStageProcess.Process != nil {
		_ = gc.GitStageProcess.Process.Kill()
		gc.GitStageProcess = nil
	}
}

// ----------------------------------
//
//	Related to Git Commit
//
// ----------------------------------
func (gc *GitCommit) GitCommit(message, description string) int {
	gc.ClearGitCommitOutput()
	gitArgs := []string{"commit", "-m", message}
	if len(description) > 0 {
		gitArgs = append(gitArgs, "-m", description)
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gc.RepoPath

	// Combine stderr into stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[PIPE ERROR]: %w", err))
		return -1
	}
	cmd.Stderr = cmd.Stdout

	gc.GitCommitProcessMutex.Lock()
	gc.GitCommitProcess = cmd
	gc.GitCommitProcessMutex.Unlock()
	defer func() {
		// ensure cleanup even if Start or Wait fails
		gc.GitCommitProcessMutex.Lock()
		gc.GitCommitProcess = nil
		gc.GitCommitProcessMutex.Unlock()
	}()

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
			gc.GitCommitOutput = append(gc.GitCommitOutput, line)
			select {
			case gc.UpdateChannel <- GIT_COMMIT_OUTPUT_UPDATE:
			default:
			}
		}
	}()

	waitErr := cmd.Wait()
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

func (gc *GitCommit) KillGitCommitCmd() {
	gc.GitCommitProcessMutex.Lock()
	defer gc.GitCommitProcessMutex.Unlock()

	if gc.GitCommitProcess != nil && gc.GitCommitProcess.Process != nil {
		_ = gc.GitCommitProcess.Process.Kill()
		gc.GitCommitProcess = nil
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
func (gc *GitCommit) GitPush(originName string) int {
	gc.ClearGitRemotePushOutput()
	gitArgs := []string{"push", "-u", originName, GITBRANCH.CurrentCheckOut.BranchName}
	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gc.RepoPath
	// Disable interactive prompts for credentials
	cmd.Env = append(os.Environ(), "GIT_ASKPASS=true", "GIT_TERMINAL_PROMPT=0")
	// Combine stderr into stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[PIPE ERROR]: %w", err))
		return -1
	}
	cmd.Stderr = cmd.Stdout

	gc.GitRemotePushProcessMutex.Lock()
	gc.GitRemotePushProcess = cmd
	gc.GitRemotePushProcessMutex.Unlock()
	defer func() {
		// ensure cleanup even if Start or Wait fails
		gc.GitRemotePushProcessMutex.Lock()
		gc.GitRemotePushProcess = nil
		gc.GitRemotePushProcessMutex.Unlock()
	}()

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

func (gc *GitCommit) KillGitRemotePushCmd() {
	gc.GitRemotePushProcessMutex.Lock()
	defer gc.GitRemotePushProcessMutex.Unlock()

	if gc.GitRemotePushProcess != nil && gc.GitRemotePushProcess.Process != nil {
		_ = gc.GitRemotePushProcess.Process.Kill()
		gc.GitRemotePushProcess = nil
	}
}

// ----------------------------------
//
//	Related to Git Remote
//
// ----------------------------------
func (gc *GitCommit) GitAddRemote(originName string, url string) ([]string, int) {
	if !isValidGitRemoteURL(url) {
		return []string{i18n.LANGUAGEMAPPING.AddRemotePopUpInvalidRemoteUrlFormat}, -1
	}
	gitArgs := []string{"remote", "add", originName, url}
	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gc.RepoPath

	gc.GitAddRemoteProcessMutex.Lock()
	gc.GitAddRemoteProcess = cmd
	gc.GitAddRemoteProcessMutex.Unlock()
	defer func() {
		gc.GitAddRemoteProcessMutex.Lock()
		gc.GitAddRemoteProcess = nil
		gc.GitAddRemoteProcessMutex.Unlock()
	}()

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
func (gc *GitCommit) KillGitAddRemoteCmd() {
	gc.GitAddRemoteProcessMutex.Lock()
	defer gc.GitAddRemoteProcessMutex.Unlock()

	if gc.GitAddRemoteProcess != nil && gc.GitAddRemoteProcess.Process != nil {
		_ = gc.GitAddRemoteProcess.Process.Kill()
		gc.GitAddRemoteProcess = nil
	}
}

func (gc *GitCommit) CheckRemoteExist() bool {
	gitArgs := []string{"remote", "-v"}
	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gc.RepoPath
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

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = repoPath
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
