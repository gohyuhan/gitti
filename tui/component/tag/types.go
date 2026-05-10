package tag

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
	GitTagItemDelegate struct{}
	GitTagItem         struct {
		TagName string
	}
)

func (i GitTagItem) FilterValue() string {
	return i.TagName
}

func (d GitTagItemDelegate) Height() int                             { return 1 }
func (d GitTagItemDelegate) Spacing() int                            { return 0 }
func (d GitTagItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitTagItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitTagItem)
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
	str := utils.TruncateString(i.TagName, componentWidth)

	fmt.Fprint(w, fn(str))
}
