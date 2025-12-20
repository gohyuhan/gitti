package commitlog

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
)

// ---------------------------------
//
// for list component of git branch
//
// ---------------------------------
type Cell struct {
	Char    rune
	ColorID int
}

type (
	GitCommitLogItemDelegate struct{}
	GitCommitLogItem         struct {
		Hash         string
		Parents      []string
		Message      string
		Author       string
		LaneCharList []Cell
		ColorID      int
	}
)

func (i GitCommitLogItem) FilterValue() string {
	return i.Hash
}

// for list component of Git branch
func (d GitCommitLogItemDelegate) Height() int                             { return 1 }
func (d GitCommitLogItemDelegate) Spacing() int                            { return 0 }
func (d GitCommitLogItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitCommitLogItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitCommitLogItem)
	if !ok {
		return
	}

	var sb strings.Builder
	sb.Grow(3)
	for idx, part := range strings.Fields(i.Author) {
		if idx > 2 {
			break
		}
		for _, r := range part {
			sb.WriteRune(r)
			break
		}
	}
	nameShortForm := sb.String()

	var commitGraphLine strings.Builder

	for _, char := range i.LaneCharList {
		commitGraphLine.WriteString(
			style.NewStyle.Foreground(style.GetColor(char.ColorID)).Render(string(char.Char)),
		)
	}

	var lineBuilder strings.Builder
	lineBuilder.WriteString(style.NewStyle.Foreground(style.ColorYellowWarm).Render(i.Hash[:7]))
	lineBuilder.WriteString(" ")
	lineBuilder.WriteString(style.NewStyle.Foreground(style.GetColor(i.ColorID)).Render(fmt.Sprintf("%-*s", 3, nameShortForm)))
	lineBuilder.WriteString(" ")
	lineBuilder.WriteString(commitGraphLine.String())
	lineBuilder.WriteString(" ")
	lineBuilder.WriteString(style.NewStyle.Render(i.Message))

	strContent := lineBuilder.String()

	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad
	needTruncate := false

	if lipgloss.Width(strContent) > componentWidth {
		needTruncate = true
		componentWidth -= 3
	}

	var fn func(...string) string
	if index == m.Index() {
		fn = func(s ...string) string {
			return style.SelectedItemStyle.Render("❯ " + strings.Join(s, " "))
		}
	} else {
		fn = func(s ...string) string {
			return style.ItemStyle.Render("  " + strings.Join(s, " "))
		}
	}

	str := style.NewStyle.MaxWidth(componentWidth).Render(strContent)

	if needTruncate {
		str += "..."
	}

	fmt.Fprint(w, fn(str))
}
