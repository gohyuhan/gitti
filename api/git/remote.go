package git

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"gitti/executor"
	"gitti/i18n"
	"gitti/utils"
)

type GitRemote struct {
	errorLog                      []error
	gitAddRemoteProcess           *exec.Cmd
	updateChannel                 chan string
	gitAddRemoteProcessMutex      sync.Mutex
	gitProcessLock                *GitProcessLock
	remote                        []GitRemoteInfo
	remoteSyncStatus              RemoteSyncStatus
	upStreamRemoteIcon            string
	currentBranchUpStream         string
	currentBranchUpStreamWithIcon string
}

type GitRemoteInfo struct {
	Name string
	Url  string
}

type RemoteSyncStatus struct {
	Local  string
	Remote string
}

func InitGitRemote(updateChannel chan string, gitProcessLock *GitProcessLock) *GitRemote {
	gitRemote := GitRemote{
		gitAddRemoteProcess:           nil,
		updateChannel:                 updateChannel,
		gitProcessLock:                gitProcessLock,
		remote:                        []GitRemoteInfo{},
		remoteSyncStatus:              RemoteSyncStatus{},
		upStreamRemoteIcon:            "",
		currentBranchUpStream:         "",
		currentBranchUpStreamWithIcon: "",
	}

	return &gitRemote
}

// ----------------------------------
//
//	Return remote
//
// ----------------------------------
func (gr *GitRemote) Remote() []GitRemoteInfo {
	return gr.remote
}

// ----------------------------------
//
//	Return remote sync status
//
// ----------------------------------
func (gr *GitRemote) RemoteSyncStatus() RemoteSyncStatus {
	return gr.remoteSyncStatus
}

// ----------------------------------
//
//	Return current upstream icon
//
// ----------------------------------
func (gr *GitRemote) UpStreamRemoteIcon() string {
	return gr.upStreamRemoteIcon
}

// ----------------------------------
//
//	Return current branch upstream
//
// ----------------------------------
func (gr *GitRemote) CurrentBranchUpStream() string {
	return gr.currentBranchUpStream
}

// ----------------------------------
//
//	Related to Git Remote
//
// ----------------------------------
func (gr *GitRemote) GitAddRemote(originName string, url string) ([]string, int) {
	if !gr.gitProcessLock.CanProceedWithGitOps() {
		return []string{gr.gitProcessLock.OtherProcessRunningWarning()}, -1
	}
	defer func() {
		gr.gitAddRemoteProcessMutex.Lock()
		gr.gitAddRemoteProcessReset()
		gr.gitAddRemoteProcessMutex.Unlock()
	}()

	if !isValidGitRemoteURL(url) {
		errMsg := "Invalid remote URL format"
		if i18n.LANGUAGEMAPPING != nil {
			errMsg = i18n.LANGUAGEMAPPING.AddRemotePopUpInvalidRemoteUrlFormat
		}
		return []string{errMsg}, -1
	}

	gr.gitAddRemoteProcessMutex.Lock()
	gitArgs := []string{"remote", "add", originName, url}
	cmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	gr.gitAddRemoteProcess = cmd

	// CombinedOutput starts and waits for the command
	gitOutput, err := cmd.CombinedOutput()
	gr.gitAddRemoteProcessMutex.Unlock()

	gitAddRemoteOutput := processGeneralGitOpsOutputIntoStringArray(gitOutput)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gr.errorLog = append(gr.errorLog, fmt.Errorf("[GIT ADD REMOTE ERROR]: %w", err))
			return gitAddRemoteOutput, status
		}
		gr.errorLog = append(gr.errorLog, fmt.Errorf("[UNEXPECTED ERROR]: %w", err))
		return gitAddRemoteOutput, -1

	}
	return gitAddRemoteOutput, 0
}

// KillGitAddRemoteCmd forcefully terminates any running git remote add process.
// It is safe to call this method even if no process is running.
// This method is thread-safe and can be called from multiple goroutines.
// This method will not be responsible to set the process state but will be the function that trigger the action will be responsible to reset the status with defer
func (gr *GitRemote) KillGitAddRemoteCmd() {
	gr.gitAddRemoteProcessMutex.Lock()
	defer gr.gitAddRemoteProcessMutex.Unlock()

	if gr.gitAddRemoteProcess != nil && gr.gitAddRemoteProcess.Process != nil {
		_ = gr.gitAddRemoteProcess.Process.Kill()
	}
}

func (gr *GitRemote) gitAddRemoteProcessReset() {
	gr.gitAddRemoteProcess = nil
	gr.gitProcessLock.ReleaseGitOpsLock()
}

func (gr *GitRemote) CheckRemoteExist() bool {
	gitArgs := []string{"remote", "-v"}
	cmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	gitOutput, err := cmd.Output()
	if err != nil {
		gr.errorLog = append(gr.errorLog, fmt.Errorf("[GIT REMOTE ERROR]: %w", err))
	}
	remotes := strings.SplitSeq(strings.TrimSpace(string(gitOutput)), "\n")
	var remoteStruct []GitRemoteInfo
	for remote := range remotes {
		remoteLinePart := strings.Fields(remote)
		if len(remoteLinePart) < 2 {
			continue
		}

		r := GitRemoteInfo{
			Name: remoteLinePart[0],
			Url:  remoteLinePart[1],
		}

		if !utils.Contains(remoteStruct, r) {
			remoteStruct = append(remoteStruct, r)

		}
	}
	gr.remote = remoteStruct
	return len(gr.remote) > 0
}

// ----------------------------------
//
//	Related to Git Remote sync status and upstream, will be call by system
//
// ----------------------------------
func (gr *GitRemote) GetLatestRemoteSyncStatusAndUpstream() {
	upstreamIcon, upstream, _ := hasUpstreamWithIcon()
	gr.upStreamRemoteIcon = upstreamIcon
	gr.currentBranchUpStream = upstream

	gitArgs := []string{"rev-list", "--left-right", "--count", "HEAD...@{upstream}"}

	remoteSyncStatusCmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	remoteSyncStatusOutput, remoteSyncStatusErr := remoteSyncStatusCmd.Output()
	if remoteSyncStatusErr != nil {
		gr.errorLog = append(gr.errorLog, fmt.Errorf("[GIT REMOTE SYNC STATUS ERROR]: %w", remoteSyncStatusErr))
		gr.remoteSyncStatus = RemoteSyncStatus{}
		return
	}

	parsedOutput := strings.TrimSpace(string(remoteSyncStatusOutput))
	parts := strings.Fields(parsedOutput)

	if len(parts) < 2 {
		gr.errorLog = append(gr.errorLog, fmt.Errorf("[GIT REMOTE SYNC STATUS ERROR]: %w", remoteSyncStatusErr))
		gr.remoteSyncStatus = RemoteSyncStatus{}
		return
	}

	gr.remoteSyncStatus = RemoteSyncStatus{
		Local:  parts[0],
		Remote: parts[1],
	}
}
