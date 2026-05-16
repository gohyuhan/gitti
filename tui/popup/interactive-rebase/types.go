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

// ------------------------------------
//
//	InteractiveRebaseOptionPopUpModel holds the option list for the interactive
//	rebase type selection popup, where the user chooses between fixup/squash,
//	reword, and drop.
//
// ------------------------------------
type InteractiveRebaseOptionPopUpModel struct {
	InteractiveRebaseOptionList list.Model
}

// ------------------------------------
//
//	InteractiveRebaseOptionDelegate renders each interactive rebase option row
//	(name + info) in the operation type selection list.
//	InteractiveRebaseOptionItem carries the display name, info string, and the
//	interactive rebase type constant for a single operation choice.
//
// ------------------------------------
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

// *************************************************************************************
//                        INTERACTIVE REBASE - FIXUP / SQUASH
// *************************************************************************************

// ------------------------------------
//
//	InteractiveRebaseFixupSquashSelectionPopUpModel holds the commit-selection
//	list, the preview viewport, retrieved commit infos, the selected-commit map,
//	pane focus flags, the sorted selection array, and any current validation
//	error for the fixup/squash commit selection popup.
//
// ------------------------------------
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

// ------------------------------------
//
//	InteractiveRebaseFixupSquashSelectionDelegate renders each commit row
//	(checkbox, short hash, author, message) in the fixup/squash selection list,
//	marking already-selected commits with a checked box.
//	InteractiveRebaseFixupSquashSelectionItem carries the hash, message, author,
//	description, parent hashes, and commit order for a single candidate commit.
//
// ------------------------------------
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

// ------------------------------------
//
//	InteractiveRebaseFixupSquashCommitPopUpModel holds the message text input,
//	description textarea, input navigation state, the sorted selected commits,
//	and the original retrieved commit list for the fixup/squash commit message
//	editing popup shown before executing the rebase.
//
// ------------------------------------
type InteractiveRebaseFixupSquashCommitPopUpModel struct {
	MessageTextInput            textinput.Model // input index 1
	DescriptionTextAreaInput    textarea.Model  // input index 2
	TotalInputCount             int             // to tell us how many input were there
	CurrentActiveInputIndex     int             // to tell us which input should be shown as highlighted/focus and be updated
	SortedSelectedCommits       []git.CommitInfo
	OriginalRetrievedCommitList []git.CommitInfo
}

// ------------------------------------
//
//	InteractiveRebaseFixupSquashOutputPopUpModel holds the output viewport,
//	spinner, and atomic state flags (IsProcessing, HasError, ProcessSuccess,
//	IsCancelled) for the fixup/squash rebase output popup. CancelFunc allows
//	the in-progress git rebase operation to be cancelled by the user.
//
// ------------------------------------
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

// *************************************************************************************
//
//	INTERACTIVE REBASE - REWORD
//
// *************************************************************************************
type InteractiveRebaseRewordSelectionPopUpModel struct {
	CommitList                  list.Model
	OriginalRetrievedCommitList []git.CommitInfo
	SelectionError              error
}

type (
	InteractiveRebaseRewordSelectionDelegate struct{}
	InteractiveRebaseRewordSelectionItem     struct {
		Hash        string
		Message     string
		Author      string
		Description string
		Parent      []string
		CommitOrder int
	}
)

func (i InteractiveRebaseRewordSelectionItem) FilterValue() string {
	return i.Hash
}
func (d InteractiveRebaseRewordSelectionDelegate) Height() int  { return 2 }
func (d InteractiveRebaseRewordSelectionDelegate) Spacing() int { return 0 }
func (d InteractiveRebaseRewordSelectionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d InteractiveRebaseRewordSelectionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(InteractiveRebaseRewordSelectionItem)
	if !ok {
		return
	}

	var firstStr string
	var secondStr string
	var str string

	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad

	firstStr = fmt.Sprintf("  %s  |  %s", i.Hash[:7], i.Author)
	firstStr = utils.TruncateString(firstStr, componentWidth)
	secondStr = utils.TruncateString(fmt.Sprintf("         %s", i.Message), componentWidth)

	// if it was a merge commit (parent is more than 1), fade the selection
	if len(i.Parent) > 1 {
		firstStr = style.ItemStyle.Faint(true).Render(firstStr)
	}
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

type InteractiveRebaseRewordCommitPopUpModel struct {
	MessageTextInput            textinput.Model // input index 1
	DescriptionTextAreaInput    textarea.Model  // input index 2
	TotalInputCount             int             // to tell us how many input were there
	CurrentActiveInputIndex     int             // to tell us which input should be shown as highlighted/focus and be updated
	SelectedCommit              git.CommitInfo
	OriginalRetrievedCommitList []git.CommitInfo
}

type InteractiveRebaseRewordOutputPopUpModel struct {
	RewordOutputViewport viewport.Model // to log out the output from git operation
	Spinner              spinner.Model  // spinner for showing processing state
	IsProcessing         atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError             atomic.Bool    // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess       atomic.Bool    // has the process sucessfuly executed
	IsCancelled          atomic.Bool    // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git rebase operation
	CancelFunc context.CancelFunc
}

// *************************************************************************************
//                           INTERACTIVE REBASE - DROP
// *************************************************************************************
