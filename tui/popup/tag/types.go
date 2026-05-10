package tag

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync/atomic"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	CreateTagPopUpModel holds two inputs for the tag creation form: a single-line
//	text input for the tag name and a dynamic-height textarea for the optional
//	tag message. CommitHash and CommitMessage identify the target commit.
//	CurrentActiveInputIndex tracks which input is focused (1 = name, 2 = message).
//
// ------------------------------------
type CreateTagPopUpModel struct {
	TagNameInput            textinput.Model
	TagMessageTextAreaInput textarea.Model
	CommitHash              string
	CommitMessage           string
	CurrentActiveInputIndex int
	TotalInputCount         int
}

// ------------------------------------
//
//	CreateTagConfirmationPopUpModel holds the collected tag creation parameters
//	shown in the confirmation prompt before the tag is created: tag name, optional
//	message, target commit hash, and commit message.
//
// ------------------------------------
type CreateTagConfirmationPopUpModel struct {
	TagName       string
	TagMessage    string
	CommitHash    string
	CommitMessage string
}

// ------------------------------------
//
//	ChooseDeleteTagOptionPopUpModel holds the tag name and the option list for
//	the delete-tag scope selection popup, which offers local-only or remote
//	deletion.
//
// ------------------------------------
type ChooseDeleteTagOptionPopUpModel struct {
	TagName          string
	DeleteOptionList list.Model
}

// ------------------------------------
//
//	ChooseRemoteForDeleteRemoteTagPopUpModel holds the remote list and the
//	selected tag name and deletion option type for the remote-selection popup
//	shown when there is more than one configured remote.
//
// ------------------------------------
type ChooseRemoteForDeleteRemoteTagPopUpModel struct {
	RemoteList       list.Model
	TagName          string
	DeleteOptionType string
}

// ------------------------------------
//
//	DeleteTagOutputPopUpModel holds the output viewport and spinner for the tag
//	deletion operation. Atomic flags IsProcessing, IsCancelled, HasError, and
//	ProcessSuccess drive spinner visibility and border color. CancelFunc aborts
//	the in-flight git operation.
//
// ------------------------------------
type DeleteTagOutputPopUpModel struct {
	TagName                 string
	DeleteTagOutputViewport viewport.Model // to log out the output from git operation
	Spinner                 spinner.Model  // spinner for showing processing state
	IsProcessing            atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                atomic.Bool    // indicate if git delete tag exitcode is not 0 (meaning have error)
	ProcessSuccess          atomic.Bool    // has the process sucessfuly executed
	IsCancelled             atomic.Bool    // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git delete tag operation
	CancelFunc context.CancelFunc
}

// ------------------------------------
//
//	DeleteTagOptionDelegate renders each tag deletion scope row (name + info).
//	DeleteTagOptionItem carries the display name, info string, and the deletion
//	type constant (local or remote) for a single tag deletion choice.
//
// ------------------------------------
type (
	DeleteTagOptionDelegate struct{}
	DeleteTagOptionItem     struct {
		Name          string
		Info          string
		DeleteTagType string
	}
)

func (i DeleteTagOptionItem) FilterValue() string {
	return i.Name
}

func (d DeleteTagOptionDelegate) Height() int                             { return 2 }
func (d DeleteTagOptionDelegate) Spacing() int                            { return 0 }
func (d DeleteTagOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d DeleteTagOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(DeleteTagOptionItem)
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

// ------------------------------------
//
//	GitRemoteForDeleteRemoteTagItemDelegate renders each remote row (name + URL)
//	in the remote selection list for remote tag deletion.
//	GitRemoteForDeleteRemoteTagItem carries the remote name, URL, and fetch/push
//	flags for a single remote entry.
//
// ------------------------------------
type (
	GitRemoteForDeleteRemoteTagItemDelegate struct{}
	GitRemoteForDeleteRemoteTagItem         struct {
		Name  string
		Url   string
		Fetch bool
		Push  bool
	}
)

func (i GitRemoteForDeleteRemoteTagItem) FilterValue() string {
	return i.Name
}

func (d GitRemoteForDeleteRemoteTagItemDelegate) Height() int                             { return 2 }
func (d GitRemoteForDeleteRemoteTagItemDelegate) Spacing() int                            { return 0 }
func (d GitRemoteForDeleteRemoteTagItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitRemoteForDeleteRemoteTagItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitRemoteForDeleteRemoteTagItem)
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

// ------------------------------------
//
//	ChoosePushTagOptionPopUpModel holds the tag name, target remote name, and the
//	option list for the push tag scope selection popup (push tag, push all,
//	force-push tag, force-push all).
//
// ------------------------------------
type ChoosePushTagOptionPopUpModel struct {
	TagName        string
	RemoteName     string
	PushOptionList list.Model
}

// ------------------------------------
//
//	PushTagOptionDelegate renders each push option row (name + info) in the push
//	tag option list.
//	PushTagOptionItem carries the display name, info string, and the push type
//	constant for a single push action choice.
//
// ------------------------------------
type (
	PushTagOptionDelegate struct{}
	PushTagOptionItem     struct {
		Name        string
		Info        string
		PushTagType string
	}
)

func (i PushTagOptionItem) FilterValue() string {
	return i.Name
}

func (d PushTagOptionDelegate) Height() int                             { return 2 }
func (d PushTagOptionDelegate) Spacing() int                            { return 0 }
func (d PushTagOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d PushTagOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(PushTagOptionItem)
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

// ------------------------------------
//
//	PushTagOutputPopUpModel holds the output viewport and spinner for the tag
//	push operation. Atomic flags IsProcessing, IsCancelled, HasError, and
//	ProcessSuccess drive spinner visibility and border color. CancelFunc aborts
//	the in-flight git push.
//
// ------------------------------------
type PushTagOutputPopUpModel struct {
	PushTagOutputViewport viewport.Model // to log out the output from git operation
	Spinner               spinner.Model  // spinner for showing processing state
	IsProcessing          atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError              atomic.Bool    // indicate if git delete tag exitcode is not 0 (meaning have error)
	ProcessSuccess        atomic.Bool    // has the process sucessfuly executed
	IsCancelled           atomic.Bool    // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git delete tag operation
	CancelFunc context.CancelFunc
}

// ------------------------------------
//
//	ChooseFetchTagOptionPopUpModel holds the target remote name and the option
//	list for the fetch tag popup (fetch, fetch-overwrite, fetch-prune, mirror).
//
// ------------------------------------
type ChooseFetchTagOptionPopUpModel struct {
	RemoteName      string
	FetchOptionList list.Model
}

// ------------------------------------
//
//	FetchTagOptionDelegate renders each fetch option row (name + info) in the
//	fetch tag option list.
//	FetchTagOptionItem carries the display name, info string, and the fetch type
//	constant for a single fetch action choice.
//
// ------------------------------------
type (
	FetchTagOptionDelegate struct{}
	FetchTagOptionItem     struct {
		Name         string
		Info         string
		FetchTagType string
	}
)

func (i FetchTagOptionItem) FilterValue() string {
	return i.Name
}

func (d FetchTagOptionDelegate) Height() int                             { return 2 }
func (d FetchTagOptionDelegate) Spacing() int                            { return 0 }
func (d FetchTagOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d FetchTagOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(FetchTagOptionItem)
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

// ------------------------------------
//
//	FetchTagOutputPopUpModel holds the output viewport and spinner for the tag
//	fetch operation. Atomic flags IsProcessing, IsCancelled, HasError, and
//	ProcessSuccess drive spinner visibility and border color. CancelFunc aborts
//	the in-flight git fetch.
//
// ------------------------------------
type FetchTagOutputPopUpModel struct {
	FetchTagOutputViewport viewport.Model // to log out the output from git operation
	Spinner                spinner.Model  // spinner for showing processing state
	IsProcessing           atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError               atomic.Bool    // indicate if git delete tag exitcode is not 0 (meaning have error)
	ProcessSuccess         atomic.Bool    // has the process sucessfuly executed
	IsCancelled            atomic.Bool    // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git delete tag operation
	CancelFunc context.CancelFunc
}
