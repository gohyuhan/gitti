package pull

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
//	ChooseGitPullTypePopUpModel holds the option list for the pull strategy
//	selection popup, which offers three choices: git pull, git pull --rebase,
//	and git pull --no-rebase (merge).
//
// ------------------------------------
type ChooseGitPullTypePopUpModel struct {
	PullTypeOptionList list.Model
}

// ------------------------------------
//
//	GitPullOutputPopUpModel holds the scrollable output viewport and dot spinner
//	for the git pull progress popup, plus atomic flags (IsProcessing, HasError,
//	ProcessSuccess, IsCancelled) and a CancelFunc to abort the pull operation.
//
// ------------------------------------
type GitPullOutputPopUpModel struct {
	PullType              string
	GitPullOutputViewport viewport.Model // to log out the output from git operation
	Spinner               spinner.Model  // spinner for showing processing state
	IsProcessing          atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError              atomic.Bool    // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess        atomic.Bool    // has the process sucessfuly executed
	IsCancelled           atomic.Bool    // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git pull operation
	CancelFunc context.CancelFunc
}

// ------------------------------------
//
//	GitPullTypeOptionDelegate renders each pull strategy row (name + command info)
//	in the choose pull type option list.
//	GitPullTypeOptionItem carries the display name, command info string, and the
//	pull type constant for a single pull strategy choice.
//
// ------------------------------------
type (
	GitPullTypeOptionDelegate struct{}
	GitPullTypeOptionItem     struct {
		Name     string
		Info     string
		PullType string
	}
)

func (i GitPullTypeOptionItem) FilterValue() string {
	return i.Name
}

func (d GitPullTypeOptionDelegate) Height() int                             { return 2 }
func (d GitPullTypeOptionDelegate) Spacing() int                            { return 0 }
func (d GitPullTypeOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitPullTypeOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitPullTypeOptionItem)
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
