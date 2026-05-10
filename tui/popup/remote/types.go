package remote

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync/atomic"
	"unicode/utf8"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	AddRemotePromptPopUpModel holds the two text inputs (remote name and URL),
//	an output viewport for the add-remote result, atomic state flags, a cancel
//	function for aborting the operation, and a flag indicating whether this repo
//	has no remote configured yet.
//
// ------------------------------------
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

// ------------------------------------
//
//	ChooseRemotePopUpModel holds a list of configured remotes and the action
//	string (e.g. "push") so the user can select which remote to target.
//
// ------------------------------------
type ChooseRemotePopUpModel struct {
	RemoteList list.Model
	Action     string
}

// ------------------------------------
//
//	GitRemoteItem holds the name, URL, and fetch/push flags for one configured
//	git remote. GitRemoteItemDelegate renders each row as a two-line entry with
//	the remote name on the first line and the URL (faint) on the second line.
//	Rows with empty name/URL render the "use local branch" rebase option instead.
//
// ------------------------------------
type (
	GitRemoteItemDelegate struct{}
	GitRemoteItem         struct {
		Name  string
		Url   string
		Fetch bool
		Push  bool
	}
)

func (i GitRemoteItem) FilterValue() string {
	return i.Name
}

func (d GitRemoteItemDelegate) Height() int                             { return 2 }
func (d GitRemoteItemDelegate) Spacing() int                            { return 0 }
func (d GitRemoteItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitRemoteItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitRemoteItem)
	if !ok {
		return
	}

	var nameStr string
	var urlStr string
	if utf8.RuneCountInString(i.Name) < 1 && utf8.RuneCountInString(i.Url) < 1 {
		nameStr = fmt.Sprintf("   %s", i18n.LANGUAGEMAPPING.GitRebaseUseLocalBranch)
		urlStr = fmt.Sprintf("    %s", i18n.LANGUAGEMAPPING.GitRebaseUseLocalBranchDesc)
	} else {
		nameStr = fmt.Sprintf("   %s", i.Name)
		urlStr = fmt.Sprintf("    %s", i.Url)
	}

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
//	RemoveRemoteConfirmationPopUpModel holds the name, URL, and fetch/push flags
//	of the remote selected for deletion, displayed in the confirmation popup.
//
// ------------------------------------
type RemoveRemoteConfirmationPopUpModel struct {
	RemoteName string
	RemoteUrl  string
	Fetch      bool
	Push       bool
}

// ------------------------------------
//
//	RemoteAsTrackingUpstreamConfirmationPopUpModel holds the remote name and URL
//	that will be set as the tracking upstream for the currently checked-out branch.
//
// ------------------------------------
type RemoteAsTrackingUpstreamConfirmationPopUpModel struct {
	RemoteName string
	RemoteUrl  string
}

// ------------------------------------
//
//	EditRemotePromptPopUpModel holds the original remote name and URL alongside
//	two text inputs pre-filled with those values, allowing the user to rename
//	the remote or update its URL before submitting the edit.
//
// ------------------------------------
type EditRemotePromptPopUpModel struct {
	OldRemoteName           string          // the old remote name
	OldRemoteUrl            string          // the old remote url
	NewRemoteNameTextInput  textinput.Model // input index 1
	NewRemoteUrlTextInput   textinput.Model // input index 2
	TotalInputCount         int             // to tell us how many input were there
	CurrentActiveInputIndex int             // to tell us which input should be shown as highlighted/focus and be updated
}
