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

// ------------------------------------
//
//	GitDiscardTypeOptionPopUpModel holds the discard type option list and the
//	target file path for the discard type selection popup. Only shown when a
//	file has both staged and unstaged changes (index and worktree).
//
// ------------------------------------
type GitDiscardTypeOptionPopUpModel struct {
	DiscardTypeOptionList list.Model
	FileDisplayName       string // the files list label, for display only
	FilePathName          string // the exact path, for the discard operation
}

// ------------------------------------
//
//	GitDiscardConfirmPromptPopUpModel holds the discard type and the target file
//	path for the discard confirmation popup, shown before the discard executes.
//
// ------------------------------------
type GitDiscardConfirmPromptPopUpModel struct {
	DiscardType     string
	FileDisplayName string // the files list label, for display only
	FilePathName    string // the exact path, for the discard operation
	CommitHash      string // for commit file discard
}

// ------------------------------------
//
//	GitDiscardTypeOptionDelegate renders each discard option row (name + info)
//	in the discard type selection list.
//	GitDiscardTypeOptionItem carries the display name, info string, and the
//	discard type constant for a single discard action choice.
//
// ------------------------------------
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

func (d GitDiscardTypeOptionDelegate) Height() int                             { return 2 }
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
