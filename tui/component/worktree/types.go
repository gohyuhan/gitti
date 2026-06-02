package worktree

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	GitWorktreeItem holds the info of a single git worktree.
//	GitWorktreeItemDelegate renders each row as the truncated worktree path.
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

	var fn func(str string, shouldFaint bool) string
	if index == m.Index() {
		fn = func(s string, shouldFaint bool) string {
			return style.SelectedItemStyle.Faint(shouldFaint).Render("❯ " + s)
		}
	} else {
		fn = func(s string, shouldFaint bool) string {
			return style.ItemStyle.Faint(shouldFaint).Render("  " + s)
		}
	}
	str := utils.TruncateString(i.WorktreePath, componentWidth)

	fmt.Fprint(w, fn(str, i.IsPrunable))
}
