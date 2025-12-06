package discard

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
// # A pop up to prompt for git discard type for file,
//
//	this will only be available when there are both changes in stage and unstage (index and worktree)
//
// ---------------------------------
type GitDiscardTypeOptionPopUpModel struct {
	DiscardTypeOptionList list.Model
	FilePathName          string
}

// ---------------------------------
//
// # To prompt user for confirmation
//
// ---------------------------------
type GitDiscardConfirmPromptPopUpModel struct {
	DiscardType  string
	FilePathName string
}

// ---------------------------------
//
// for discard option selection option
//
// ---------------------------------
type (
	GitDiscardTypeOptionDelegate struct{}
	GitDiscardTypeOptionItem     struct {
		Name        string
		Info        string
		DiscardType string
	}
)

func (i GitDiscardTypeOptionItem) FilterValue() string {
	return i.Name
}

// for discard type selection
func (d GitDiscardTypeOptionDelegate) Height() int                             { return 1 }
func (d GitDiscardTypeOptionDelegate) Spacing() int                            { return 0 }
func (d GitDiscardTypeOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitDiscardTypeOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitDiscardTypeOptionItem)
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
