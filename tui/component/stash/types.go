package stash

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

// ---------------------------------
//
// for list component of git stashed files
//
// ---------------------------------
type (
	GitStashItemDelegate struct{}
	GitStashItem         struct {
		Id      string
		Message string
	}
)

func (i GitStashItem) FilterValue() string {
	return i.Message
}

// for list component of stash
func (d GitStashItemDelegate) Height() int                             { return 1 }
func (d GitStashItemDelegate) Spacing() int                            { return 0 }
func (d GitStashItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitStashItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitStashItem)
	if !ok {
		return
	}

	str := fmt.Sprintf(" %s", i.Message)

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
