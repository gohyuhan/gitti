package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gohyuhan/gitti/utils"
)

// ------------------------------------------
//
//  LOGGING will be the only part that doesn't support i18n, and will always be EN
//
// ------------------------------------------

type LogItem struct {
	OpsTimeString    string
	OpsType          string
	OpsCommand       string // this is use to record either the underlying triggered command or description of the OPS (git ops will always be a command, other will depends)
	OpsSeverityLevel string // either INFO, WARN or ERROR
	OpsDescription   string
}

type GittiLogging struct {
	logs            []LogItem
	maxLogsCount    int
	updateChannel   chan string
	showLatestXLogs int // this is used to control how many latest logs to show in the log component
}

// ----------------------------------
//
//	initialize GittiLogging
//
// ----------------------------------
func InitGittiLogging(maxLogsCount int, updateChannel chan string, showLatestXLogs int) *GittiLogging {
	return &GittiLogging{
		logs:            make([]LogItem, 0, maxLogsCount),
		maxLogsCount:    maxLogsCount,
		updateChannel:   updateChannel,
		showLatestXLogs: showLatestXLogs,
	}
}

// ----------------------------------
//
//	Get logs will only return the latest 3 log item, this was used for the log component, to get the full list, used GetFullLogs
//
// ----------------------------------
func (gl *GittiLogging) GetLogs() []LogItem {
	if len(gl.logs) > gl.showLatestXLogs {
		return gl.logs[len(gl.logs)-gl.showLatestXLogs-1:] // we get only the latest 3 log items
	} else {
		return gl.logs
	}
}

// ----------------------------------
//
//	Return the complete log history
//
// ----------------------------------
func (gl *GittiLogging) GetFullLogs() []LogItem {
	return gl.logs
}

// ----------------------------------
//
//	Create and append a new log entry, evicting the oldest entry if at capacity
//
// ----------------------------------
func (gl *GittiLogging) RegisterNewLog(logOpsType string, logOpsCommand string, logOpsSeverityLevel string, logOpsDescription string, isGitOps bool) {
	if isGitOps {
		logOpsCommand = "git " + logOpsCommand
	}
	newLogItem := LogItem{
		OpsTimeString:    time.Now().Format("2006-01-02T15:04:05-0700"),
		OpsType:          logOpsType,
		OpsCommand:       logOpsCommand,
		OpsSeverityLevel: logOpsSeverityLevel,
		OpsDescription:   logOpsDescription,
	}
	if len(gl.logs) < gl.maxLogsCount {
		gl.logs = append(gl.logs, newLogItem)
	} else {
		gl.logs = append(gl.logs[1:], newLogItem)
	}

	gl.updateChannel <- NEW_LOG_UPDATE
}

// ----------------------------------
//
//	Export all logs to a timestamped file in the user's Downloads directory
//
// ----------------------------------
func (gl *GittiLogging) ExportLogging() {
	gl.RegisterNewLog(EXPORT_LOGGING_OPS, "Export Logging Requested", INFO, "", false)
	exportDir, err := utils.GetDownloadsDir()
	if err != nil {
		gl.RegisterNewLog(EXPORT_LOGGING_OPS, "", ERROR, fmt.Sprintf("[%s ERROR]:%s", EXPORT_LOGGING_OPS, err.Error()), false)
		return
	}

	var contentString strings.Builder

	for _, log := range gl.GetFullLogs() {
		opsTimeStringCharCount := utf8.RuneCountInString(log.OpsTimeString)
		switch log.OpsSeverityLevel {
		case INFO:
			contentString.WriteString(log.OpsTimeString)
			contentString.WriteString("  ")
			contentString.WriteString("[INFO]")
			contentString.WriteString(" ")
			contentString.WriteString(log.OpsType)
			contentString.WriteString("\n")
			contentString.WriteString(strings.Repeat(" ", opsTimeStringCharCount+2))
			contentString.WriteString(log.OpsCommand)
		case WARN:
			contentString.WriteString(log.OpsTimeString)
			contentString.WriteString("  ")
			contentString.WriteString("[WARN]")
			contentString.WriteString(" ")
			contentString.WriteString(log.OpsType)
			contentString.WriteString("\n")
			contentString.WriteString(strings.Repeat(" ", opsTimeStringCharCount+2))
			contentString.WriteString(log.OpsDescription)
		case ERROR:
			contentString.WriteString(log.OpsTimeString)
			contentString.WriteString("  ")
			contentString.WriteString("[ERROR]")
			contentString.WriteString(" ")
			contentString.WriteString(log.OpsType)
			contentString.WriteString("\n")
			contentString.WriteString(strings.Repeat(" ", opsTimeStringCharCount+2))
			contentString.WriteString(log.OpsDescription)
		}
		contentString.WriteString("\n")
	}

	now := time.Now()
	filename := fmt.Sprintf("gitti-log-%s.log", now.Format("2006-01-02T15:04:05-0700"))

	path := filepath.Join(exportDir, filename)

	// Write the content (create new file or overwrite if exists)
	err = os.WriteFile(path, []byte(contentString.String()), 0o644)
	if err != nil {
		gl.RegisterNewLog(EXPORT_LOGGING_OPS, "", ERROR, fmt.Sprintf("[%s ERROR]: %s", EXPORT_LOGGING_OPS, err.Error()), false)
	} else {
		gl.RegisterNewLog(EXPORT_LOGGING_OPS, fmt.Sprintf("Log exported successfully: %s", path), INFO, "", false)
	}
}
