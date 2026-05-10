package branch

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
//	GitBranchItem holds the name and checkout status of a branch for display
//	in the branch list. GitBranchItemDelegate implements list.ItemDelegate,
//	rendering each row with a "* " prefix for the checked-out branch and
//	"❯" for the selected row.
//
// ------------------------------------
type (
	GitBranchItemDelegate struct{}
	GitBranchItem         struct {
		BranchName   string
		IsCheckedOut bool
	}
)

func (i GitBranchItem) FilterValue() string {
	return i.BranchName
}

func (d GitBranchItemDelegate) Height() int                             { return 1 }
func (d GitBranchItemDelegate) Spacing() int                            { return 0 }
func (d GitBranchItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitBranchItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitBranchItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("   %s", i.BranchName)
	if i.IsCheckedOut {
		str = fmt.Sprintf(" * %s", i.BranchName)
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
	str = utils.TruncateString(str, componentWidth)

	fmt.Fprint(w, fn(str))
}
