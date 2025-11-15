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

	"gitti/executor"
	"gitti/i18n"
	"gitti/utils"
)

type GitCommit struct {
	errorLog                  []error
	gitCommitProcess          *exec.Cmd
	gitRemotePushProcess      *exec.Cmd
	gitAddRemoteProcess       *exec.Cmd
	gitCommitOutput           []string
	gitRemotePushOutput       []string
	updateChannel             chan string
	gitCommitProcessMutex     sync.Mutex
	gitRemotePushProcessMutex sync.Mutex
	gitAddRemoteProcessMutex  sync.Mutex
	gitProcessLock            *GitProcessLock
	remote                    []GitRemote
}

type GitRemote struct {
	Name string
	Url  string
}

type LatestCommitMsgAndDesc struct {
	Message     string
	Description string
}

func InitGitCommit(updateChannel chan string, gitProcessLock *GitProcessLock) *GitCommit {
	gitCommit := GitCommit{
		gitCommitProcess:     nil,
		gitRemotePushProcess: nil,
		gitAddRemoteProcess:  nil,
		gitCommitOutput:      []string{},
		gitRemotePushOutput:  []string{},
		updateChannel:        updateChannel,
		gitProcessLock:       gitProcessLock,
		remote:               []GitRemote{},
	}

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
	cmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	// Disable interactive prompts for credentials
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	_, err := cmd.Output()
	if err != nil {
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", err))
	}
}

// ----------------------------------
//
//	Related to Git Commit
//
// ----------------------------------
func (gc *GitCommit) GitCommit(message, description string, isAmendCommit bool) int {
	if !gc.gitProcessLock.CanProceedWithGitOps() {
		return -1
	}

	defer func() {
		// ensure cleanup even if Start or Wait fails
		gc.gitCommitProcessMutex.Lock()
		gc.gitCommitProcessReset()
		gc.gitCommitProcessMutex.Unlock()
	}()

	gc.gitCommitProcessMutex.Lock()

	gc.ClearGitCommitOutput()
	gitArgs := []string{"commit", "-m", message}
	if isAmendCommit {
		gitArgs = []string{"commit", "--amend", "-m", message}
	}
	if len(description) > 0 {
		gitArgs = append(gitArgs, "-m", description)
	}

	commitCmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, true)
	gc.gitCommitProcess = commitCmd

	// Combine stderr into stdout
	stdout, err := commitCmd.StdoutPipe()
	if err != nil {
		gc.gitCommitProcessMutex.Unlock()
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[PIPE ERROR]: %w", err))
		return -1
	}
	commitCmd.Stderr = commitCmd.Stdout

	// Start the process while still holding the mutex
	if err := commitCmd.Start(); err != nil {
		gc.gitCommitProcessMutex.Unlock()
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[START ERROR]: %w", err))
		return -1
	}

	// Process is now started and can be killed safely
	gc.gitCommitProcessMutex.Unlock()

	// Stream combined output
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		scanner.Split(splitOnCarriageReturnOrNewline)
		cursorIndex := 0
		for scanner.Scan() {
			updatedCursorIndex, updatedGitCommitOutput := handleProgressOutputStream(cursorIndex, scanner, gc.gitCommitOutput)
			gc.gitCommitOutput = updatedGitCommitOutput
			cursorIndex = updatedCursorIndex
			if isAmendCommit {
				gc.updateChannel <- GIT_AMEND_COMMIT_OUTPUT_UPDATE
			} else {
				gc.updateChannel <- GIT_COMMIT_OUTPUT_UPDATE
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
func (gc *GitCommit) KillGitCommitCmd() {
	gc.gitCommitProcessMutex.Lock()
	defer gc.gitCommitProcessMutex.Unlock()

	if gc.gitCommitProcess != nil && gc.gitCommitProcess.Process != nil {
		_ = gc.gitCommitProcess.Process.Kill()
	}
}

func (gc *GitCommit) gitCommitProcessReset() {
	gc.gitCommitProcess = nil
	gc.gitProcessLock.ReleaseGitOpsLock()
}

// ----------------------------------
//
//	Related to Git Commit (Amend)
//
// ----------------------------------
func (gc *GitCommit) GetLatestCommitMsgAndDesc() LatestCommitMsgAndDesc {
	gitArgs := []string{"log", "-1", "--pretty=format:%s%n%b", "HEAD"}
	latestCommitCmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	commitMsgAndDesc, cmdErr := latestCommitCmd.Output()
	if cmdErr != nil {
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[GET LATEST COMMIT INFO ERROR]: %w", cmdErr))
		return LatestCommitMsgAndDesc{}
	}

	parsed := strings.SplitN(string(commitMsgAndDesc), "\n", 2)
	title := parsed[0]
	description := ""
	if len(parsed) > 1 {
		description = parsed[1]
	}

	return LatestCommitMsgAndDesc{
		Message:     title,
		Description: description,
	}
}

// ----------------------------------
//
//	Related to Git Push
//
// ----------------------------------
func (gc *GitCommit) GitPush(currentCheckOutBranch string, originName string, pushType string) int {
	if !gc.gitProcessLock.CanProceedWithGitOps() {
		return -1
	}
	defer func() {
		// ensure cleanup even if Start or Wait fails
		gc.gitRemotePushProcessMutex.Lock()
		gc.resetGitRemotePushProcesstatus()
		gc.gitRemotePushProcessMutex.Unlock()
	}()

	gc.gitRemotePushProcessMutex.Lock()
	gc.ClearGitRemotePushOutput()
	var gitArgs []string
	switch pushType {
	case FORCEPUSHSAFE:
		gitArgs = []string{"push", "--progress", "--force-with-lease", "-u", originName, currentCheckOutBranch}
	case FORCEPUSHDANGEROUS:
		gitArgs = []string{"push", "--progress", "--force", "-u", originName, currentCheckOutBranch}
	default:
		gitArgs = []string{"push", "--progress", "-u", originName, currentCheckOutBranch}
	}
	cmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, true)
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
		scanner.Split(splitOnCarriageReturnOrNewline)
		cursorIndex := 0
		for scanner.Scan() {
			updatedCursorIndex, updatedGitRemotePushOutput := handleProgressOutputStream(cursorIndex, scanner, gc.gitRemotePushOutput)
			gc.gitRemotePushOutput = updatedGitRemotePushOutput
			cursorIndex = updatedCursorIndex
			gc.updateChannel <- GIT_REMOTE_PUSH_OUTPUT_UPDATE
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
	gc.gitProcessLock.ReleaseGitOpsLock()
}

// ----------------------------------
//
//	Related to Git Remote
//
// ----------------------------------
func (gc *GitCommit) GitAddRemote(originName string, url string) ([]string, int) {
	if !gc.gitProcessLock.CanProceedWithGitOps() {
		return []string{gc.gitProcessLock.OtherProcessRunningWarning()}, -1
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
	cmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	gc.gitAddRemoteProcess = cmd

	// CombinedOutput starts and waits for the command
	gitOutput, err := cmd.CombinedOutput()
	gc.gitAddRemoteProcessMutex.Unlock()

	gitAddRemoteOutput := processGeneralGitOpsOutputIntoStringArray(gitOutput)
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
	gc.gitProcessLock.ReleaseGitOpsLock()
}

func (gc *GitCommit) CheckRemoteExist() bool {
	gitArgs := []string{"remote", "-v"}
	cmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	gitOutput, err := cmd.Output()
	if err != nil {
		gc.errorLog = append(gc.errorLog, fmt.Errorf("[GIT COMMIT ERROR]: %w", err))
	}
	remotes := strings.SplitSeq(strings.TrimSpace(string(gitOutput)), "\n")
	var remoteStruct []GitRemote
	for remote := range remotes {
		remoteLinePart := strings.Fields(remote)
		if len(remoteLinePart) < 2 {
			continue
		}

		r := GitRemote{
			Name: remoteLinePart[0],
			Url:  remoteLinePart[1],
		}

		if !utils.Contains(remoteStruct, r) {
			remoteStruct = append(remoteStruct, GitRemote{
				Name: remoteLinePart[0],
				Url:  remoteLinePart[1],
			})

		}
	}
	gc.remote = remoteStruct
	return len(gc.remote) > 0
}

// ----------------------------------
//
//	Related to Git Init
//
// ----------------------------------
func GitInit(repoPath string, initBranchName string) {
	initGitArgs := []string{"init"}

	initCmd := executor.GittiCmdExecutor.RunGitCmd(initGitArgs, false)
	_, initErr := initCmd.Output()
	if initErr != nil {
		fmt.Printf("[GIT INIT ERROR]: %v", initErr)
		os.Exit(1)
	}

	// set the branch
	checkoutBranchGitArgs := []string{"checkout", "-b", initBranchName}

	checkoutBranchCmd := executor.GittiCmdExecutor.RunGitCmd(checkoutBranchGitArgs, false)
	_, checkoutBranchErr := checkoutBranchCmd.Output()
	if checkoutBranchErr != nil {
		fmt.Printf("[GIT INIT ERROR]: %v", checkoutBranchErr)
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
