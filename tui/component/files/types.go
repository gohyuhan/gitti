package files

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
// for list component of git modified files
//
// ---------------------------------
type (
	GitModifiedFilesItemDelegate struct{}
	GitModifiedFilesItem         struct {
		FilePathname string
		IndexState   string
		WorkTree     string
		HasConflict  bool
	}
)

func (i GitModifiedFilesItem) FilterValue() string {
	return i.FilePathname
}

// for list component of modified files
func (d GitModifiedFilesItemDelegate) Height() int                             { return 1 }
func (d GitModifiedFilesItemDelegate) Spacing() int                            { return 0 }
func (d GitModifiedFilesItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitModifiedFilesItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitModifiedFilesItem)
	if !ok {
		return
	}

	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad - 5
	filePathName := utils.TruncateString(i.FilePathname, componentWidth)

	indexState := style.StagedFileStyle.Render(i.IndexState)
	workTree := style.UnstagedFileStyle.Render(i.WorkTree)
	if i.IndexState == "?" {
		indexState = style.UnstagedFileStyle.Render(i.IndexState)
	}

	str := fmt.Sprintf(" %s%s %s", indexState, workTree, filePathName)

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

	fmt.Fprint(w, fn(str))
}
