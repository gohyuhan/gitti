package git

import (
	"fmt"
	"strings"

	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/logging"
)

type GitRefLog struct {
	updateChannel  chan string
	gitProcessLock *GitProcessLock
	reflog         []GitRefLogInfo // all remote info
	logging        *logging.GittiLogging
	maxRefLogCount int
}

type GitRefLogInfo struct {
	FullInfo   string
	InfoDesc   string // InfoDesc is the description of the reflog entry which is a split of Head from FullInfo
	Head       string // Head is the head of the reflog entry
	Action     string // Action is the action of the reflog entry
	ActionInfo string // ActionInfo is the action info of the reflog entry
	Hash       string // Hash is the hash associate with the reflog entry
}

// ------------------------------------
//
//	Initialize the git reflog handler with shared dependencies
//
// ------------------------------------
func InitGitRefLog(updateChannel chan string, gitProcessLock *GitProcessLock, maxRefLogCount int, logging *logging.GittiLogging) *GitRefLog {
	gitRefLog := GitRefLog{
		updateChannel:  updateChannel,
		gitProcessLock: gitProcessLock,
		reflog:         []GitRefLogInfo{},
		logging:        logging,
		maxRefLogCount: maxRefLogCount,
	}

	return &gitRefLog
}

// ------------------------------------
//
//	Return reflog
//
// ------------------------------------
func (grl *GitRefLog) RefLog() []GitRefLogInfo {
	return grl.reflog
}

// ------------------------------------
//
//		Get latest RefLog entry
//	 * passive, called by daemon only
//
// ------------------------------------
func (grl *GitRefLog) GetLatestRefLog() {
	var latestRefLog []GitRefLogInfo

	gitArgs := []string{"reflog", "--no-abbrev", "-n", fmt.Sprintf("%d", grl.maxRefLogCount)}
	latestRefLogCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	reflogOutput, reflogErr := latestRefLogCmdExecutor.Output()

	if reflogErr != nil {
		return
	}

	parsedRefLogOutputArray := processGeneralGitOpsOutputIntoStringArray(reflogOutput)
	for _, entry := range parsedRefLogOutputArray {
		splitedRefLogEntry := strings.SplitN(entry, " ", 2)
		if len(splitedRefLogEntry) < 2 {
			continue
		}
		splitedHeadAndInfoDesc := strings.SplitN(splitedRefLogEntry[1], ":", 2)
		if len(splitedHeadAndInfoDesc) < 2 {
			continue
		}
		splitedActionAndActionInfoDesc := strings.SplitN(splitedHeadAndInfoDesc[1], ":", 2)
		if len(splitedActionAndActionInfoDesc) < 2 {
			continue
		}

		reflogInfo := GitRefLogInfo{
			FullInfo:   strings.TrimSpace(splitedRefLogEntry[1]),
			InfoDesc:   strings.TrimSpace(splitedHeadAndInfoDesc[1]),
			Head:       strings.TrimSpace(splitedHeadAndInfoDesc[0]),
			Action:     strings.TrimSpace(splitedActionAndActionInfoDesc[0]),
			ActionInfo: strings.TrimSpace(splitedActionAndActionInfoDesc[1]),
			Hash:       strings.TrimSpace(splitedRefLogEntry[0]),
		}

		latestRefLog = append(latestRefLog, reflogInfo)
	}

	grl.reflog = latestRefLog
}
