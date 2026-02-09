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

// ---------------------------------
//
// # For Create Tag Pop Up Model
//
// ---------------------------------
type CreateTagPopUpModel struct {
	TagNameInput            textinput.Model
	TagMessageTextAreaInput textarea.Model
	CommitHash              string
	CommitMessage           string
	CurrentActiveInputIndex int
	TotalInputCount         int
}

// ---------------------------------
//
// # For Create Tag Confirmation Pop Up Model
//
// ---------------------------------
type CreateTagConfirmationPopUpModel struct {
	TagName       string
	TagMessage    string
	CommitHash    string
	CommitMessage string
}

// ---------------------------------
//
// # For Choose Delete Tag Option Pop Up Model
//
// ---------------------------------
type ChooseDeleteTagOptionPopUpModel struct {
	TagName          string
	DeleteOptionList list.Model
}

// ---------------------------------
//
// # For Choose Remote For Delete Remote Tag Pop Up Model
//
// ---------------------------------
type ChooseRemoteForDeleteRemoteTagPopUpModel struct {
	RemoteList       list.Model
	TagName          string
	DeleteOptionType string
}

// ------------------------------------
//
//	For git delete tag process pop up model
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
//	For tag deletion option selection option
//
// ------------------------------------
type (
	TagDeleteOptionDelegate struct{}
	TagDeleteOptionItem     struct {
		Name          string
		Info          string
		DeleteTagType string
	}
)

func (i TagDeleteOptionItem) FilterValue() string {
	return i.Name
}

// for tag deletion option selection
func (d TagDeleteOptionDelegate) Height() int                             { return 1 }
func (d TagDeleteOptionDelegate) Spacing() int                            { return 0 }
func (d TagDeleteOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d TagDeleteOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(TagDeleteOptionItem)
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
//	For list component of git remote for deleting remote tag
//
// ------------------------------------
type (
	GitRemoteForDeleteRemoteTagItemDelegate struct{}
	GitRemoteForDeleteRemoteTagItem         struct {
		Name string
		Url  string
	}
)

func (i GitRemoteForDeleteRemoteTagItem) FilterValue() string {
	return i.Name
}

// for list component of git remote for deleting remote tag
func (d GitRemoteForDeleteRemoteTagItemDelegate) Height() int                             { return 1 }
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
