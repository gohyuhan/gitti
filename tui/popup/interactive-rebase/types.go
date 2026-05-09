package interactiverebase

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
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ---------------------------------
//
// choose a rebase option, fixup/squash, reword, drop
//
// ---------------------------------
type InteractiveRebaseOptionPopUpModel struct {
	InteractiveRebaseOptionList list.Model
}

// ---------------------------------
//
// for interactive rebase option selection option
//
// ---------------------------------
type (
	InteractiveRebaseOptionDelegate struct{}
	InteractiveRebaseOptionItem     struct {
		Name                  string
		Info                  string
		InteractiveRebaseType string
	}
)

// ------------------------------------
//
//	Returns filter text for interactive rebase option list item
//
// ------------------------------------
func (i InteractiveRebaseOptionItem) FilterValue() string {
	return i.Name
}

// ------------------------------------
//
//	Returns list delegate row height for interactive rebase option items
//
// ------------------------------------
func (d InteractiveRebaseOptionDelegate) Height() int { return 2 }

// ------------------------------------
//
//	Returns list delegate row spacing for interactive rebase option items
//
// ------------------------------------
func (d InteractiveRebaseOptionDelegate) Spacing() int { return 0 }

// ------------------------------------
//
//	No-op update hook for interactive rebase option list delegate
//
// ------------------------------------
func (d InteractiveRebaseOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// ------------------------------------
//
//	Renders one interactive rebase option row with truncation and selected styling
//
// ------------------------------------
func (d InteractiveRebaseOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(InteractiveRebaseOptionItem)
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

// ---------------------------------
//
// for fixup/squash commit selection
//
// ---------------------------------
type InteractiveRebaseFixupSquashSelectionPopUpModel struct {
	CommitList                          list.Model
	CommitFixupSquashViewport           viewport.Model
	OriginalRetrievedCommitList         []git.CommitInfo
	SelectedCommitHashMap               map[string]git.CommitInfo // key is commit hash
	IsCommitListSelected                bool                      // to track if user was in the list section or not, default is here
	IsCommitFixupSquashViewportSelected bool                      // to track if user was in the viewport section or not
	SortedSelectedCommits               []git.CommitInfo
	SelectionError                      error
}

// ---------------------------------
//
// bubble tea list for selecting commit for fixup squash
//
// ---------------------------------
type (
	InteractiveRebaseFixupSquashSelectionDelegate struct {
		SelectedCommitHashMap *map[string]git.CommitInfo // key is commit hash
	}
	InteractiveRebaseFixupSquashSelectionItem struct {
		Hash        string
		Message     string
		Author      string
		Description string
		Parent      []string
		CommitOrder int
	}
)

// ------------------------------------
//
//	Returns filter text for fixup/squash commit selection list item
//
// ------------------------------------
func (i InteractiveRebaseFixupSquashSelectionItem) FilterValue() string {
	return i.Hash
}

// ------------------------------------
//
//	Returns list delegate row height for fixup/squash commit selection items
//
// ------------------------------------
func (d InteractiveRebaseFixupSquashSelectionDelegate) Height() int { return 2 }

// ------------------------------------
//
//	Returns list delegate row spacing for fixup/squash commit selection items
//
// ------------------------------------
func (d InteractiveRebaseFixupSquashSelectionDelegate) Spacing() int { return 0 }

// ------------------------------------
//
//	No-op update hook for fixup/squash commit selection list delegate
//
// ------------------------------------
func (d InteractiveRebaseFixupSquashSelectionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

// ------------------------------------
//
//	Renders one fixup/squash commit selection row with selected checkbox state and selected styling
//
// ------------------------------------
func (d InteractiveRebaseFixupSquashSelectionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(InteractiveRebaseFixupSquashSelectionItem)
	if !ok {
		return
	}

	var firstStr string
	var secondStr string
	var str string

	var selected bool
	if d.SelectedCommitHashMap != nil {
		_, selected = (*d.SelectedCommitHashMap)[i.Hash]
	}
	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad

	if selected {
		firstStr = fmt.Sprintf("[X]  %s  |  %s", i.Hash[:7], i.Author)
	} else {
		firstStr = fmt.Sprintf("[ ]  %s  |  %s", i.Hash[:7], i.Author)
	}

	firstStr = utils.TruncateString(firstStr, componentWidth)
	secondStr = utils.TruncateString(fmt.Sprintf("         %s", i.Message), componentWidth)

	// except for the first string, all other string should be faded
	secondStr = style.ItemStyle.Faint(true).Render(secondStr)

	str = fmt.Sprintf("%s\n%s", firstStr, secondStr)

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

	fmt.Fprint(w, fn(str))
}

// ---------------------------------
//
// Pop up model for editing the final commit message and description before executing the fixup/squash rebase
//
// ---------------------------------
type InteractiveRebaseFixupSquashCommitPopUpModel struct {
	MessageTextInput            textinput.Model // input index 1
	DescriptionTextAreaInput    textarea.Model  // input index 2
	TotalInputCount             int             // to tell us how many input were there
	CurrentActiveInputIndex     int             // to tell us which input should be shown as highlighted/focus and be updated
	SortedSelectedCommits       []git.CommitInfo
	OriginalRetrievedCommitList []git.CommitInfo
}

// ---------------------------------
//
// Pop up model for displaying fixup/squash rebase output, with spinner and atomic state flags for processing/error/success/cancellation
//
// ---------------------------------
type InteractiveRebaseFixupSquashOutputPopUpModel struct {
	FixupSquashOutputViewport viewport.Model // to log out the output from git operation
	Spinner                   spinner.Model  // spinner for showing processing state
	IsProcessing              atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                  atomic.Bool    // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess            atomic.Bool    // has the process sucessfuly executed
	IsCancelled               atomic.Bool    // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git rebase operation
	CancelFunc context.CancelFunc
}
