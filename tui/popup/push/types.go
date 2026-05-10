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

// ------------------------------------
//
//	ChoosePushTypePopUpModel holds the option list and the target remote name for
//	the push type selection popup. The list offers three choices: normal push,
//	safe force-push (--force-with-lease), and dangerous force-push (--force).
//
// ------------------------------------
type ChoosePushTypePopUpModel struct {
	PushOptionList list.Model
	RemoteName     string
}

// ------------------------------------
//
//	GitPushOptionDelegate renders each push type row (name + command info) in
//	the choose push type option list.
//	GitPushOptionItem carries the display name, command info string, and the
//	push type constant for a single push mode choice.
//
// ------------------------------------
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

func (d GitPushOptionDelegate) Height() int                             { return 2 }
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

// ------------------------------------
//
//	GitRemotePushPopUpModel holds the scrollable output viewport and dot spinner
//	for the git push progress popup, plus atomic flags (IsProcessing, HasError,
//	ProcessSuccess, IsCancelled) and a CancelFunc to abort the push operation.
//
// ------------------------------------
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
