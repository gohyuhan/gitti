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
type (
	GitCommitLogItemDelegate struct{}
	GitCommitLogItem         struct {
		Hash       string
		Parents    []string
		Message    string
		Author     string
		LaneString string
		Color      string
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

	strContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		style.NewStyle.Foreground(style.ColorYellowWarm).Render(i.Hash[:7]),
		" ",
		style.NewStyle.Render(fmt.Sprintf("%s%-*s", i.Color, 3, nameShortForm)),
		" ",
		style.NewStyle.Render(i.LaneString),
		"  ",
		style.NewStyle.Render(i.Message),
	)

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
