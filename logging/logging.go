package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gohyuhan/gitti/utils"
)

// ------------------------------------
//
//   Logging module - EN only, no i18n support
//
// ------------------------------------

type LogItem struct {
	OpsTimeString    string
	OpsType          string
	OpsCommand       string // this is use to record either the underlying triggered command or description of the OPS (git ops will always be a command, other will depends)
	OpsSeverityLevel string // either INFO, WARN or ERROR
	OpsDescription   string
}

type GittiLogging struct {
	logsMu          sync.RWMutex
	logs            []LogItem
	maxLogsCount    int
	updateChannel   chan string
	showLatestXLogs int // this is used to control how many latest logs to show in the log component
}

// ------------------------------------
//
//	Initialize GittiLogging with max log count and update channel
//
// ------------------------------------
func InitGittiLogging(maxLogsCount int, updateChannel chan string, showLatestXLogs int) *GittiLogging {
	return &GittiLogging{
		logs:            make([]LogItem, 0, maxLogsCount),
		maxLogsCount:    maxLogsCount,
		updateChannel:   updateChannel,
		showLatestXLogs: showLatestXLogs,
	}
}

// ------------------------------------
//
//	Get latest logs for display in log component
//
// ------------------------------------
func (gl *GittiLogging) GetLogs() []LogItem {
	gl.logsMu.RLock()
	defer gl.logsMu.RUnlock()
	visible := gl.logs
	if len(gl.logs) > gl.showLatestXLogs {
		visible = gl.logs[len(gl.logs)-gl.showLatestXLogs-1:] // we get only the latest 3 log items
	}
	copied := make([]LogItem, len(visible))
	copy(copied, visible)
	return copied
}

// ------------------------------------
//
//	Return the complete log history
//
// ------------------------------------
func (gl *GittiLogging) GetFullLogs() []LogItem {
	gl.logsMu.RLock()
	defer gl.logsMu.RUnlock()
	copied := make([]LogItem, len(gl.logs))
	copy(copied, gl.logs)
	return copied
}

func (gl *GittiLogging) ClearLogs() {
	gl.logsMu.Lock()
	defer gl.logsMu.Unlock()
	gl.logs = make([]LogItem, 0, gl.maxLogsCount)
}

// ------------------------------------
//
//	Create and append a new log entry, evicting the oldest entry if at capacity
//
// ------------------------------------
func (gl *GittiLogging) RegisterNewLog(logOpsType string, logOpsCommand string, logOpsSeverityLevel string, logOpsDescription string, isGitOps bool) {
	if isGitOps {
		logOpsCommand = "git " + logOpsCommand
	}
	// commands carry exact file paths and messages; escape control characters so a name with a
	// newline or a terminal escape sequence stays on one line in the log panel and exported log
	logOpsCommand = utils.EscapeControlCharacters(logOpsCommand)
	newLogItem := LogItem{
		OpsTimeString:    time.Now().Format("2006-01-02T15:04:05-0700"),
		OpsType:          logOpsType,
		OpsCommand:       logOpsCommand,
		OpsSeverityLevel: logOpsSeverityLevel,
		OpsDescription:   logOpsDescription,
	}
	gl.logsMu.Lock()
	if len(gl.logs) < gl.maxLogsCount {
		gl.logs = append(gl.logs, newLogItem)
	} else {
		gl.logs = append(gl.logs[1:], newLogItem)
	}
	gl.logsMu.Unlock()

	// The UI rebuilds from the full log snapshot. A queued update covers newer logs too.
	select {
	case gl.updateChannel <- NEW_LOG_UPDATE:
	default:
	}
}

// ------------------------------------
//
//	Export all logs to a timestamped file in the user's Downloads directory
//
// ------------------------------------
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
