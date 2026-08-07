package reflog

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
)

// ------------------------------------
//
//	GitRefLogItem holds the full info string, description, HEAD reference, action,
//	action detail, and short hash for one reflog entry. GitRefLogItemDelegate renders
//	each row as a yellow 7-char hash followed by the truncated description.
//
// ------------------------------------
type (
	GitRefLogItemDelegate struct{}
	GitRefLogItem         struct {
		FullInfo   string
		InfoDesc   string
		Head       string
		Action     string
		ActionInfo string
		Hash       string
	}
)

func (i GitRefLogItem) FilterValue() string {
	return i.Hash + " " + i.InfoDesc
}

func (d GitRefLogItemDelegate) Height() int                             { return 1 }
func (d GitRefLogItemDelegate) Spacing() int                            { return 0 }
func (d GitRefLogItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitRefLogItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitRefLogItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s %s", style.NewStyle.Foreground(style.ColorYellowWarm).Render(i.Hash[:7]), i.InfoDesc)

	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad

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
	str = ansi.Truncate(str, componentWidth, "...")

	fmt.Fprint(w, fn(str))
}
