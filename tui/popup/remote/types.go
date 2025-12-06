package remote

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync/atomic"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ---------------------------------
//
// # For add Remote prompt pop up
//
// ---------------------------------
type AddRemotePromptPopUpModel struct {
	RemoteNameTextInput     textinput.Model // input index 1
	RemoteUrlTextInput      textinput.Model // input index 2
	TotalInputCount         int             // to tell us how many input were there
	CurrentActiveInputIndex int             // to tell us which input should be shown as highlighted/focus and be updated
	AddRemoteOutputViewport viewport.Model  // to log out the output from git operation
	IsProcessing            atomic.Bool     // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                atomic.Bool     // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess          atomic.Bool     // has the process sucessfuly executed
	NoInitialRemote         bool            // indicate if this repo has no remote yet or user just wanted to add more remote
	IsCancelled             atomic.Bool     // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git remote add operation
	CancelFunc context.CancelFunc
}

// ---------------------------------
//
// choose a remote to push to
//
// ---------------------------------
type ChooseRemotePopUpModel struct {
	RemoteList list.Model
}

// ---------------------------------
//
// for list component of git remote
//
// ---------------------------------
type (
	GitRemoteItemDelegate struct{}
	GitRemoteItem         struct {
		Name string
		Url  string
	}
)

func (i GitRemoteItem) FilterValue() string {
	return i.Name
}

// for list component of git remote
func (d GitRemoteItemDelegate) Height() int                             { return 1 }
func (d GitRemoteItemDelegate) Spacing() int                            { return 0 }
func (d GitRemoteItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitRemoteItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitRemoteItem)
	if !ok {
		return
	}

	nameStr := fmt.Sprintf("   %s", i.Name)
	urlStr := fmt.Sprintf("    %s", i.Url)

	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad - 2

	nameStr = utils.TruncateString(nameStr, componentWidth)
	urlStr = utils.TruncateString(urlStr, componentWidth)

	nameRendered := style.ItemStyle.Render(nameStr)
	urlRendered := style.ItemStyle.Faint(true).Render(urlStr)
	fullStr := nameRendered + "\n" + "  " + urlRendered

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
