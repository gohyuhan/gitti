package interactiverebase

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
// choose a rebase option, fixup/squash, reword, drop
//
// ---------------------------------
type InteractiveRebaseOptionPopUp struct {
	InteractiveRebaseOptionList list.Model
}

// ---------------------------------
//
// for interactive rebase option selection option
//
// ---------------------------------
type (
	InteractiveRebaseOptionDelegate struct{}
	InteractiveRebaseOptionItem     struct {
		Name                  string
		Info                  string
		InteractiveRebaseType string
	}
)

func (i InteractiveRebaseOptionItem) FilterValue() string {
	return i.Name
}

// for interactive rebase selection
func (d InteractiveRebaseOptionDelegate) Height() int                             { return 2 }
func (d InteractiveRebaseOptionDelegate) Spacing() int                            { return 0 }
func (d InteractiveRebaseOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d InteractiveRebaseOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(InteractiveRebaseOptionItem)
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
