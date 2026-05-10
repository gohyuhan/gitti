package remote

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
//	GitRemoteItem holds the name, URL, and fetch/push flags for one configured
//	git remote. GitRemoteItemDelegate renders each row as the truncated remote name.
//
// ------------------------------------
type (
	GitRemoteItemDelegate struct{}
	GitRemoteItem         struct {
		Name  string
		Url   string
		Fetch bool
		Push  bool
	}
)

func (i GitRemoteItem) FilterValue() string {
	return i.Name
}

func (d GitRemoteItemDelegate) Height() int                             { return 1 }
func (d GitRemoteItemDelegate) Spacing() int                            { return 0 }
func (d GitRemoteItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitRemoteItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitRemoteItem)
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
	str := utils.TruncateString(i.Name, componentWidth)

	fmt.Fprint(w, fn(str))
}
