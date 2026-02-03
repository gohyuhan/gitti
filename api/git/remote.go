package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/logging"
)

type GitRemote struct {
	updateChannel                 chan string
	gitProcessLock                *GitProcessLock
	remote                        []GitRemoteInfo
	remoteSyncStatus              RemoteSyncStatus
	upStreamRemoteIcon            string
	currentBranchUpStream         string
	currentBranchUpStreamWithIcon string
	logging                       *logging.GittiLogging
}

type GitRemoteInfo struct {
	Name string
	Url  string
}

type RemoteSyncStatus struct {
	Local  string
	Remote string
}

func InitGitRemote(updateChannel chan string, gitProcessLock *GitProcessLock, logging *logging.GittiLogging) *GitRemote {
	gitRemote := GitRemote{
		updateChannel:                 updateChannel,
		gitProcessLock:                gitProcessLock,
		remote:                        []GitRemoteInfo{},
		remoteSyncStatus:              RemoteSyncStatus{},
		upStreamRemoteIcon:            "",
		currentBranchUpStream:         "",
		currentBranchUpStreamWithIcon: "",
		logging:                       logging,
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
func (gr *GitRemote) GitAddRemote(ctx context.Context, originName string, url string) ([]string, int) {
	if !gr.gitProcessLock.CanProceedWithGitOps() {
		return []string{gr.gitProcessLock.OtherProcessRunningWarning()}, -1
	}
	defer func() {
		gr.gitProcessLock.ReleaseGitOpsLock()
	}()

	if !isValidGitRemoteURL(url) {
		errMsg := "Invalid remote URL format"
		if i18n.LANGUAGEMAPPING != nil {
			errMsg = i18n.LANGUAGEMAPPING.AddRemotePopUpInvalidRemoteUrlFormat
		}
		return []string{errMsg}, -1
	}

	gitArgs := []string{"remote", "add", originName, url}
	cmd := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, false)

	// CombinedOutput starts and waits for the command
	gitOutput, err := cmd.CombinedOutput()
	gr.logging.RegisterNewLog(logging.ADD_REMOTE_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)

	gitAddRemoteOutput := processGeneralGitOpsOutputIntoStringArray(gitOutput)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gr.logging.RegisterNewLog(logging.ADD_REMOTE_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.ADD_REMOTE_OPS, err.Error()), true)
			return gitAddRemoteOutput, status
		}
		gr.logging.RegisterNewLog(logging.ADD_REMOTE_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.ADD_REMOTE_OPS, err.Error()), true)
		return gitAddRemoteOutput, -1

	}
	return gitAddRemoteOutput, 0
}

func (gr *GitRemote) CheckRemoteExist() bool {
	gitArgs := []string{"remote", "-v"}
	cmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	gitOutput, err := cmd.Output()
	gr.logging.RegisterNewLog(logging.CHECK_REMOTE_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	if err != nil {
		gr.logging.RegisterNewLog(logging.CHECK_REMOTE_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.CHECK_REMOTE_OPS, err.Error()), true)
		return false
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

		if strings.TrimSpace(remoteLinePart[2]) == "(push)" {
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
func (gr *GitRemote) GetLatestRemoteSyncStatusAndUpstream(needFetch bool, userTriggered bool) {
	upstreamIcon, upstream, _ := hasUpstreamWithIcon()
	gr.upStreamRemoteIcon = upstreamIcon
	gr.currentBranchUpStream = upstream

	if needFetch {
		gitFetch(gr.logging, userTriggered)
	}

	gitArgs := []string{"rev-list", "--left-right", "--count", "HEAD...@{upstream}"}

	remoteSyncStatusCmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	remoteSyncStatusOutput, remoteSyncStatusErr := remoteSyncStatusCmd.Output()
	if remoteSyncStatusErr != nil {
		gr.remoteSyncStatus = RemoteSyncStatus{}
		return
	}

	parsedOutput := strings.TrimSpace(string(remoteSyncStatusOutput))
	parts := strings.Fields(parsedOutput)

	if len(parts) < 2 {
		gr.logging.RegisterNewLog(logging.CHECK_REMOTE_SYNC_STATUS_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: Invalid output format", logging.CHECK_REMOTE_SYNC_STATUS_OPS), true)
		gr.remoteSyncStatus = RemoteSyncStatus{}
		return
	}

	gr.remoteSyncStatus = RemoteSyncStatus{
		Local:  parts[0],
		Remote: parts[1],
	}
}
