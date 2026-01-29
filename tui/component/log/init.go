package log

import (
	"context"
	"strings"

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

		switch log.OpsSeverityLevel {
		case logging.INFO:
			contentString.WriteString("[INFO]")
			contentString.WriteString("   ")
			contentString.WriteString(log.OpsType)
			contentString.WriteString("\n")
			contentString.WriteString("         ")
			contentString.WriteString(log.OpsCommand)
		case logging.WARN:
			contentString.WriteString(style.NewStyle.Foreground(style.ColorYellowWarm).Render("[WARN]"))
			contentString.WriteString("   ")
			contentString.WriteString(style.NewStyle.Foreground(style.ColorYellowWarm).Render(log.OpsType))
			contentString.WriteString("\n")
			contentString.WriteString("         ")
			contentString.WriteString(style.NewStyle.Foreground(style.ColorYellowWarm).Render(log.OpsDescription))
		case logging.ERROR:
			contentString.WriteString(style.NewStyle.Foreground(style.ColorError).Render("[ERROR]"))
			contentString.WriteString("  ")
			contentString.WriteString(style.NewStyle.Foreground(style.ColorError).Render(log.OpsType))
			contentString.WriteString("\n")
			contentString.WriteString("         ")
			contentString.WriteString(style.NewStyle.Foreground(style.ColorError).Render(log.OpsDescription))
		}

		contentString.WriteString("\n")
	}
	return contentString.String()
}
