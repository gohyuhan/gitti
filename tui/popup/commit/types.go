package commit

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
// # For git commit process pop up model
//
// ---------------------------------
type GitCommitPopUpModel struct {
	IsAmendCommit            bool            // to indicate is this is a normal commit or an amend commit operation
	MessageTextInput         textinput.Model // input index 1
	DescriptionTextAreaInput textarea.Model  // input index 2
	TotalInputCount          int             // to tell us how many input were there
	CurrentActiveInputIndex  int             // to tell us which input should be shown as highlighted/focus and be updated
	GitCommitOutputViewport  viewport.Model  // to log out the output from git operation
	Spinner                  spinner.Model   // spinner for showing processing state
	InitialCommitStarted     atomic.Bool     // indicated that this pop up session has start the first commit action
	IsProcessing             atomic.Bool     // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                 atomic.Bool     // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess           atomic.Bool     // has the process sucessfuly executed
	IsCancelled              atomic.Bool     // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git commit operation
	CancelFunc context.CancelFunc
}

// ---------------------------------
//
// # For git amend commit process pop up model
//
// ---------------------------------
type GitAmendCommitPopUpModel struct {
	IsAmendCommit                bool            // to indicate is this is a normal commit or an amend commit operation
	MessageTextInput             textinput.Model // input index 1
	DescriptionTextAreaInput     textarea.Model  // input index 2
	TotalInputCount              int             // to tell us how many input were there
	CurrentActiveInputIndex      int             // to tell us which input should be shown as highlighted/focus and be updated
	GitAmendCommitOutputViewport viewport.Model  // to log out the output from git operation
	Spinner                      spinner.Model   // spinner for showing processing state
	InitialCommitStarted         atomic.Bool     // indicated that this pop up session has start the first commit action
	IsProcessing                 atomic.Bool     // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                     atomic.Bool     // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess               atomic.Bool     // has the process sucessfuly executed
	IsCancelled                  atomic.Bool     // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git amend commit operation
	CancelFunc context.CancelFunc
}

// ---------------------------------
//
// # A pop up to prompt for git reset type that will apply for the latest commit
//
// ---------------------------------
type GitResetLatestCommitTypeOptionPopUpModel struct {
	ResetLatestCommitTypeOptionList list.Model
}

// ---------------------------------
//
// # To prompt user for confirmation for reset type on the latest commit
//
// ---------------------------------
type GitResetLatestCommitConfirmPromptPopUpModel struct {
	GitResetLatestCommitType string
}

// ---------------------------------
//
// for reset latest commit option selection option
//
// ---------------------------------
type (
	GitResetLatestCommitTypeOptionDelegate struct{}
	GitResetLatestCommitTypeOptionItem     struct {
		Name      string
		Info      string
		ResetType string
	}
)

func (i GitResetLatestCommitTypeOptionItem) FilterValue() string {
	return i.Name
}

// for reset laytes commit type selection
func (d GitResetLatestCommitTypeOptionDelegate) Height() int                             { return 1 }
func (d GitResetLatestCommitTypeOptionDelegate) Spacing() int                            { return 0 }
func (d GitResetLatestCommitTypeOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitResetLatestCommitTypeOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitResetLatestCommitTypeOptionItem)
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
