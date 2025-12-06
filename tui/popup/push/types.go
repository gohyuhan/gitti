package push

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync/atomic"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ---------------------------------
//
// user choose how do they want to push the commit, push /  push --force / push --force-with-lease
//
// ---------------------------------
type ChoosePushTypePopUpModel struct {
	PushOptionList list.Model
	RemoteName     string
}

// ---------------------------------
//
// for push selection option
//
// ---------------------------------
type (
	GitPushOptionDelegate struct{}
	GitPushOptionItem     struct {
		Name     string
		Info     string
		PushType string
	}
)

func (i GitPushOptionItem) FilterValue() string {
	return i.Name
}

// for push selection option
func (d GitPushOptionDelegate) Height() int                             { return 1 }
func (d GitPushOptionDelegate) Spacing() int                            { return 0 }
func (d GitPushOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitPushOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitPushOptionItem)
	if !ok {
		return
	}

	nameStr := fmt.Sprintf("   %s", i.Name)
	urlStr := fmt.Sprintf("    %s", i.Info)

	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad - 2

	nameStr = utils.TruncateString(nameStr, componentWidth)
	urlStr = utils.TruncateString(urlStr, componentWidth)

	nameRendered := style.ItemStyle.Render(nameStr)
	infoRendered := style.ItemStyle.Faint(true).Render(urlStr)
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

// ---------------------------------
//
// # For Remote push process pop up
//
// ---------------------------------
type GitRemotePushPopUpModel struct {
	GitRemotePushOutputViewport viewport.Model // to log out the output from git operation
	Spinner                     spinner.Model  // spinner for showing processing state
	IsProcessing                atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                    atomic.Bool    // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess              atomic.Bool    // has the process sucessfuly executed
	IsCancelled                 atomic.Bool    // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git push operation
	CancelFunc context.CancelFunc
}
