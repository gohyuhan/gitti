package resolve

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
// for resolve conflict option pop up
//
// ---------------------------------
type GitResolveConflictOptionPopUpModel struct {
	FilePathName              string
	ResolveConflictOptionList list.Model
}

// ---------------------------------
//
// for resolve conflict option selection option
//
// ---------------------------------
type (
	GitResolveConflictOptionDelegate struct{}
	GitResolveConflictOptionItem     struct {
		Name        string
		Info        string
		ResolveType string
	}
)

func (i GitResolveConflictOptionItem) FilterValue() string {
	return i.Name
}

// for resolve conflict option selection
func (d GitResolveConflictOptionDelegate) Height() int                             { return 1 }
func (d GitResolveConflictOptionDelegate) Spacing() int                            { return 0 }
func (d GitResolveConflictOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitResolveConflictOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitResolveConflictOptionItem)
	if !ok {
		return
	}

	nameStr := fmt.Sprintf("   %s", i.Name)
	infoStr := fmt.Sprintf("    %s", i.Info)

	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad - 2

	nameStr = utils.TruncateString(nameStr, componentWidth)
	infoStr = utils.TruncateString(infoStr, componentWidth)

	nameRendered := style.ItemStyle.Render(nameStr)
	infoRendered := style.ItemStyle.Faint(true).Render(infoStr)
	fullStr := nameRendered + "\n" + "  " + infoRendered

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

	fmt.Fprint(w, fn(fullStr))
}
