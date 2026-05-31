package worktree

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	GitTagItem holds the name of a single git tag.
//	GitTagItemDelegate renders each row as the truncated tag name.
//
// ------------------------------------
type (
	GitWorktreeItemDelegate struct{}
	GitWorktreeItem         struct {
		WorktreePath        string
		WorktreeHead        string
		WorktreeBranch      string
		IsMain              bool
		IsInCurrentWorktree bool
		IsLocked            bool
		LockReason          string
		IsPrunable          bool
	}
)

func (i GitWorktreeItem) FilterValue() string {
	return i.WorktreePath
}

func (d GitWorktreeItemDelegate) Height() int                             { return 1 }
func (d GitWorktreeItemDelegate) Spacing() int                            { return 0 }
func (d GitWorktreeItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitWorktreeItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitWorktreeItem)
	if !ok {
		return
	}

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
	str := utils.TruncateString(i.WorktreePath, componentWidth)

	fmt.Fprint(w, fn(str))
}
