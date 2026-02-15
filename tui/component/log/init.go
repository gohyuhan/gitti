package log

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

func InitGittiLogViewport(m *types.GittiModel, ForLogComponent bool, ctx context.Context) string {
	var logsArray []logging.LogItem

	if ForLogComponent {
		logsArray = m.GittiLogger.GetLogs()
	} else {
		logsArray = m.GittiLogger.GetFullLogs()
	}

	var contentString strings.Builder

	for _, log := range logsArray {
		// Check for context cancellation during iteration (only if ctx is not nil)
		if ctx != nil {
			select {
			case <-ctx.Done():
				return contentString.String()
			default:
			}
		}
		opsTimeStringCharCount := utf8.RuneCountInString(log.OpsTimeString)
		switch log.OpsSeverityLevel {
		case logging.INFO:
			contentString.WriteString(log.OpsTimeString)
			contentString.WriteString("  ")
			contentString.WriteString("[INFO]")
			contentString.WriteString(" ")
			contentString.WriteString(log.OpsType)
			contentString.WriteRune('\n')
			contentString.WriteString(strings.Repeat(" ", opsTimeStringCharCount+2))
			contentString.WriteString(log.OpsCommand)
		case logging.WARN:
			contentString.WriteString(style.NewStyle.Foreground(style.ColorYellowWarm).Render(log.OpsTimeString))
			contentString.WriteString("  ")
			contentString.WriteString(style.NewStyle.Foreground(style.ColorYellowWarm).Render("[WARN]"))
			contentString.WriteString(" ")
			contentString.WriteString(style.NewStyle.Foreground(style.ColorYellowWarm).Render(log.OpsType))
			contentString.WriteRune('\n')
			logDescriptionArray := strings.Split(log.OpsDescription, "\n")
			for _, logDescription := range logDescriptionArray {
				contentString.WriteString(strings.Repeat(" ", opsTimeStringCharCount+2))
				contentString.WriteString(style.NewStyle.Foreground(style.ColorYellowWarm).Render(logDescription))
				contentString.WriteRune('\n')
			}
		case logging.ERROR:
			contentString.WriteString(style.NewStyle.Foreground(style.ColorError).Render(log.OpsTimeString))
			contentString.WriteString("  ")
			contentString.WriteString(style.NewStyle.Foreground(style.ColorError).Render("[ERROR]"))
			contentString.WriteString(" ")
			contentString.WriteString(style.NewStyle.Foreground(style.ColorError).Render(log.OpsType))
			contentString.WriteRune('\n')
			logDescriptionArray := strings.Split(log.OpsDescription, "\n")
			for _, logDescription := range logDescriptionArray {
				contentString.WriteString(strings.Repeat(" ", opsTimeStringCharCount+2))
				contentString.WriteString(style.NewStyle.Foreground(style.ColorError).Render(logDescription))
				contentString.WriteRune('\n')
			}
		}

		contentString.WriteRune('\n')
	}
	return contentString.String()
}
