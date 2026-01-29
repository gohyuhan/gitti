package logging

// ------------------------------------------
//
//  LOGGING will be the only part that doesn't support i18n, and will always be EN
//
// ------------------------------------------

type LogItem struct {
	opsType          string
	opsCommand       string
	opsSeverityLevel string // either INFO, WARN or ERROR
	opsDescription   string
}

type GittiLogging struct {
	logs          []LogItem
	maxLogsCount  int
	updateChannel chan string
}

// initialize GittiLogging
func InitGittiLogging(maxLogsCount int, updateChannel chan string) *GittiLogging {
	return &GittiLogging{
		logs:          make([]LogItem, 0, maxLogsCount),
		maxLogsCount:  maxLogsCount,
		updateChannel: updateChannel,
	}
}

// Get logs will only return the latest 3 log item, this was used for the log component, to get the full list, used GetFullLogs
func (gl *GittiLogging) GetLogs() []LogItem {
	if len(gl.logs) > 3 {
		return gl.logs[len(gl.logs)-4:] // we get only the latest 3 log items
	} else {
		return gl.logs
	}
}

func (gl *GittiLogging) GetFullLogs() []LogItem {
	return gl.logs
}

func (gl *GittiLogging) RegisterNewLog(logOpsType string, logOpsCommand string, logOpsSeverityLevel string, logOpsDescription string, isGitOps bool) {
	if isGitOps {
		logOpsCommand = "git " + logOpsCommand
	}
	newLogItem := LogItem{
		opsType:          logOpsType,
		opsCommand:       logOpsCommand,
		opsSeverityLevel: logOpsSeverityLevel,
		opsDescription:   logOpsDescription,
	}
	if len(gl.logs) < gl.maxLogsCount {
		gl.logs = append(gl.logs, newLogItem)
	} else {
		gl.logs = append(gl.logs[1:], newLogItem)
	}

	gl.updateChannel <- NEW_LOG_UPDATE
}
