package git

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/logging"
)

type TagInfo struct {
	TagName string
}

type GitTag struct {
	updateChannel    chan string
	allTag           []TagInfo
	logging          *logging.GittiLogging
	gitProcessLock   *GitProcessLock
	tagPushOutput    []string
	tagPushOutputMu  sync.RWMutex
	tagFetchOutput   []string
	tagFetchOutputMu sync.RWMutex
}

func InitGitTag(updateChannel chan string, gitProcessLock *GitProcessLock, logging *logging.GittiLogging) *GitTag {
	gitTag := GitTag{
		updateChannel:  updateChannel,
		gitProcessLock: gitProcessLock,
		logging:        logging,
	}
	return &gitTag
}

// ----------------------------------
//
//	Return all tag
//
// ----------------------------------
func (gt *GitTag) AllTag() []TagInfo {
	return gt.allTag
}

// ----------------------------------
//
//	Return git tag push output
//
// ----------------------------------
func (gt *GitTag) TagPushOutput() []string {
	gt.tagPushOutputMu.RLock()
	defer gt.tagPushOutputMu.RUnlock()

	copied := make([]string, len(gt.tagPushOutput))
	copy(copied, gt.tagPushOutput)
	return copied
}

// ----------------------------------
//
//	Fetch the latest tag available
//
// ----------------------------------
func (gt *GitTag) GetLatestGitTag() {
	var latestTags []TagInfo
	gitArgs := []string{"tag"}

	getLatestGitTagCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	getLatestGitTagCmdOutput, getLatestGitTagCmdErr := getLatestGitTagCmdExecutor.Output()

	if getLatestGitTagCmdErr != nil {
		gt.logging.RegisterNewLog(logging.GET_LATEST_TAG_INFO_OPS, "", logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.GET_LATEST_TAG_INFO_OPS, getLatestGitTagCmdErr.Error()), true)
		return
	}

	parsedOutput := processGeneralGitOpsOutputIntoStringArray(getLatestGitTagCmdOutput)
	if len(parsedOutput) < 1 {
		gt.allTag = []TagInfo{}
		return
	}
	for index := range parsedOutput {
		tag := TagInfo{
			TagName: parsedOutput[index],
		}

		latestTags = append(latestTags, tag)
	}

	gt.allTag = latestTags
}

// ----------------------------------
//
//	Create a new git tag, return the output and is success boolean
//
// ----------------------------------
func (gt *GitTag) CreateNewTag(commitHash string, newTagName string, tagMessage string) {
	if !gt.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer func() {
		gt.gitProcessLock.ReleaseGitOpsLock()
	}()

	gitArgs := []string{"tag", "-a", newTagName, commitHash, "-m", tagMessage}
	createNewTagCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, true)
	gt.logging.RegisterNewLog(logging.CREATE_NEW_TAG_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	createNewTagCmdErr := createNewTagCmdExecutor.Run()
	if createNewTagCmdErr != nil {
		gt.logging.RegisterNewLog(logging.CREATE_NEW_TAG_OPS, "", logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.CREATE_NEW_TAG_OPS, createNewTagCmdErr.Error()), true)
	}

	return
}

// ----------------------------------
//
//	Fetch the detail associate with the provided tag
//
// ----------------------------------
func (gt *GitTag) ShowGitTagDetail(ctx context.Context, tagName string) []string {
	gitArgs := []string{"show", tagName}

	tagDetailCmdExecutor := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, true)
	tagDetailCmdOutput, tagDetailCmdErr := tagDetailCmdExecutor.Output()

	if tagDetailCmdErr != nil {
		if ctx.Err() != nil {
			// This catches context.Canceled
			gt.logging.RegisterNewLog(logging.TAG_DETAIL_OPS, "", logging.WARN, fmt.Sprintf("[%s CANCELLED]", logging.TAG_DETAIL_OPS), true)
			return nil
		} else {
			gt.logging.RegisterNewLog(logging.TAG_DETAIL_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.TAG_DETAIL_OPS, tagDetailCmdErr.Error()), true)
			return nil
		}
	}

	tagDetail := processGeneralGitOpsOutputIntoStringArray(tagDetailCmdOutput)

	return tagDetail
}

// ----------------------------------
//
//	Push new tag to remote
//
// ----------------------------------
func (gt *GitTag) GitTagPush(ctx context.Context, originName string, tagName string, pushType string) int {
	if !gt.gitProcessLock.CanProceedWithGitOps() {
		return -1
	}
	defer func() {
		gt.gitProcessLock.ReleaseGitOpsLock()
	}()

	// Reset the output buffer
	gt.ClearGitTagPushOutput()

	var gitArgs []string
	switch pushType {
	case TAGPUSH:
		gitArgs = []string{"push", "--progress", originName, tagName}
	case TAGPUSHALL:
		gitArgs = []string{"push", "--progress", originName, "--tags"}
	case TAGPUSHFORCE:
		gitArgs = []string{"push", "--progress", originName, tagName, "--force"}
	case TAGPUSHALLFORCE:
		gitArgs = []string{"push", "--progress", originName, "--tags", "--force"}
	default:
		// default will be same as TAGPUSHALL
		gitArgs = []string{"push", "--progress", originName, "--tags"}
	}

	cmd := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, true)

	// Combine stderr into stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		gt.logging.RegisterNewLog(logging.TAG_PUSH_OPS, "", logging.ERROR, fmt.Sprintf("[PIPE ERROR]: %s", err.Error()), false)
		return -1
	}
	cmd.Stderr = cmd.Stdout

	// Start the process
	if err := cmd.Start(); err != nil {
		gt.logging.RegisterNewLog(logging.TAG_PUSH_OPS, "", logging.ERROR, fmt.Sprintf("[START ERROR]: %s", err.Error()), false)
		return -1
	} else {
		gt.logging.RegisterNewLog(logging.TAG_PUSH_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	}

	// Stream combined output
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		scanner.Split(splitOnCarriageReturnOrNewline)
		cursorIndex := 0
		lastSent := time.Time{}
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return // Stop immediately on cancel
			default:
				gt.tagPushOutputMu.Lock()
				updatedCursorIndex, updatedTagPushOutput := handleProgressOutputStream(cursorIndex, scanner, gt.tagPushOutput)
				gt.tagPushOutput = updatedTagPushOutput
				cursorIndex = updatedCursorIndex
				gt.tagPushOutputMu.Unlock()
				if time.Since(lastSent) >= STREAMUPDATETHROTTLEMS*time.Millisecond {
					select {
					case gt.updateChannel <- GIT_TAG_PUSH_OUTPUT_UPDATE:
						lastSent = time.Now()
					default:
					}
				}
			}
		}
		// trigger an update once it ends
		gt.updateChannel <- GIT_TAG_PUSH_OUTPUT_UPDATE
	}()

	waitErr := cmd.Wait()
	wg.Wait()

	if ctx.Err() != nil {
		gt.logging.RegisterNewLog(logging.TAG_PUSH_OPS, "", logging.WARN, fmt.Sprintf("[%s CANCELLED]", logging.TAG_PUSH_OPS), true)
		return -1
	}

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gt.logging.RegisterNewLog(logging.TAG_PUSH_OPS, "", logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.TAG_PUSH_OPS, exitErr.Error()), true)
			return status
		}
		gt.logging.RegisterNewLog(logging.TAG_PUSH_OPS, "", logging.ERROR, fmt.Sprintf("[UNEXPECTED ERROR]: %s", waitErr.Error()), false)
		return -1
	}

	return 0
}

func (gt *GitTag) ClearGitTagPushOutput() {
	gt.tagPushOutputMu.Lock()
	defer gt.tagPushOutputMu.Unlock()
	gt.tagPushOutput = []string{}
}

// ----------------------------------
//
//	Fetch tags from remote
//
// ----------------------------------
func (gt *GitTag) GitTagFetch(ctx context.Context, originName string, fetchType string) int {
	if !gt.gitProcessLock.CanProceedWithGitOps() {
		return -1
	}
	defer func() {
		gt.gitProcessLock.ReleaseGitOpsLock()
	}()

	gt.ClearGitTagFetchOutput()
	var gitArgs []string
	switch fetchType {
	case TAGFETCH:
		gitArgs = []string{"fetch", "--progress", originName, "--tags"}
	case TAGFETCHOVERWRITE:
		gitArgs = []string{"fetch", "--progress", originName, "--tags", "--force"}
	case TAGFETCHPRUNE:
		gitArgs = []string{"fetch", "--progress", originName, "--tags", "--prune-tags"}
	case TAGFETCHMIRROR:
		gitArgs = []string{"fetch", "--progress", originName, "--tags", "--force", "--prune-tags"}
	default:
		// default will be same as TAGFETCH
		gitArgs = []string{"fetch", "--progress", originName, "--tags"}
	}

	cmdExecutor := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, true)

	// Combine stderr into stdout
	stdout, err := cmdExecutor.StdoutPipe()
	if err != nil {
		gt.logging.RegisterNewLog(logging.TAG_FETCH_OPS, "", logging.ERROR, fmt.Sprintf("[PIPE ERROR]: %s", err.Error()), false)
		return -1
	}
	cmdExecutor.Stderr = cmdExecutor.Stdout

	// Start the process
	if err := cmdExecutor.Start(); err != nil {
		gt.logging.RegisterNewLog(logging.TAG_FETCH_OPS, "", logging.ERROR, fmt.Sprintf("[START ERROR]: %s", err.Error()), false)
		return -1
	} else {
		gt.logging.RegisterNewLog(logging.TAG_FETCH_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	}

	// Stream combined output
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		scanner.Split(splitOnCarriageReturnOrNewline)
		cursorIndex := 0
		lastSent := time.Time{}
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return // Stop immediately on cancel
			default:
				gt.tagFetchOutputMu.Lock()
				updatedCursorIndex, updatedTagFetchOutput := handleProgressOutputStream(cursorIndex, scanner, gt.tagFetchOutput)
				gt.tagFetchOutput = updatedTagFetchOutput
				cursorIndex = updatedCursorIndex
				gt.tagFetchOutputMu.Unlock()
				if time.Since(lastSent) >= STREAMUPDATETHROTTLEMS*time.Millisecond {
					select {
					case gt.updateChannel <- GIT_TAG_FETCH_OUTPUT_UPDATE:
						lastSent = time.Now()
					default:
					}
				}
			}
		}
		// trigger an update once it ends
		gt.updateChannel <- GIT_TAG_FETCH_OUTPUT_UPDATE
	}()

	waitErr := cmdExecutor.Wait()
	wg.Wait()

	if ctx.Err() != nil {
		gt.logging.RegisterNewLog(logging.TAG_FETCH_OPS, "", logging.WARN, fmt.Sprintf("[%s CANCELLED]: %s", logging.TAG_FETCH_OPS, ctx.Err().Error()), true)
		return -1
	}

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			status := exitErr.ExitCode()
			gt.logging.RegisterNewLog(logging.TAG_FETCH_OPS, "", logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.TAG_FETCH_OPS, waitErr.Error()), true)
			return status
		}
		gt.logging.RegisterNewLog(logging.TAG_FETCH_OPS, "", logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.TAG_FETCH_OPS, waitErr.Error()), true)
		return -1
	}
	return 0
}

func (gt *GitTag) ClearGitTagFetchOutput() {
	gt.tagFetchOutputMu.Lock()
	defer gt.tagFetchOutputMu.Unlock()
	gt.tagFetchOutput = []string{}
}

// ----------------------------------
//
//	Delete tags from local/remote
//
// ----------------------------------
func (gt *GitTag) GitDeleteTag(ctx context.Context, originName string, tagName string, deleteType string) ([]string, bool) {
	if !gt.gitProcessLock.CanProceedWithGitOps() {
		return []string{gt.gitProcessLock.OtherProcessRunningWarning()}, false
	}
	defer func() {
		gt.gitProcessLock.ReleaseGitOpsLock()
	}()

	var gitArgs []string
	switch deleteType {
	case TAGDELETELOCAL:
		gitArgs = []string{"tag", "-d", tagName}
	case TAGDELETEREMOTE:
		gitArgs = []string{"push", originName, "--delete", tagName}
	default:
		gitArgs = []string{"tag", "-d", tagName}
	}

	deleteTagCmdExecutor := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, true)
	gt.logging.RegisterNewLog(logging.TAG_DELETE_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	deleteTagCmdOutput, deleteTagCmdErr := deleteTagCmdExecutor.CombinedOutput()

	parsedDeleteTagOutput := processGeneralGitOpsOutputIntoStringArray(deleteTagCmdOutput)

	if deleteTagCmdErr != nil {
		if ctx.Err() != nil {
			// This catches context.Canceled
			gt.logging.RegisterNewLog(logging.TAG_DELETE_OPS, "", logging.WARN, fmt.Sprintf("[%s CANCELLED]", logging.TAG_DELETE_OPS), true)
			return parsedDeleteTagOutput, false
		} else {
			gt.logging.RegisterNewLog(logging.TAG_DELETE_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.TAG_DELETE_OPS, deleteTagCmdErr.Error()), true)
			return parsedDeleteTagOutput, false
		}
	}

	return parsedDeleteTagOutput, true
}
