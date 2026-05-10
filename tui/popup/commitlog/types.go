package commitlog

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	GitCherryPickPopUpModel holds the branch name and the commit-log list
//	for the cherry-pick selection popup, where the user picks commits from the
//	current branch to cherry-pick.
//
// ------------------------------------
type GitCherryPickPopUpModel struct {
	CurrentBranchName                string
	CurrentBranchCherryPickCommitLog list.Model
}

// ------------------------------------
//
//	GitEditCherryPickPopUpModel holds the ordered list of already-selected
//	cherry-pick commits. Used in the edit-cherry-pick popup, where the user
//	can remove or reorder commits before applying them.
//
// ------------------------------------
type GitEditCherryPickPopUpModel struct {
	CherryPickedCommitLog list.Model
}

// ------------------------------------
//
//	GitCherryPickOptionSelectionPopUpModel holds the operation-type list for
//	the cherry-pick action selection popup, where the user chooses between
//	pick, edit-and-pick, and apply.
//
// ------------------------------------
type GitCherryPickOptionSelectionPopUpModel struct {
	CherryPickedOpsOption list.Model
}

// ------------------------------------
//
//	GitCherryPickDelegate renders each commit-log row in the cherry-pick
//	selection list, showing a checkbox, short hash, author, and message.
//	GitCherryPickItem carries the hash, message, author, and source branch
//	for a single candidate commit.
//
// ------------------------------------
type (
	GitCherryPickDelegate struct {
		CherryPickedMap *map[string]git.CherryPickedCommitLog
	}
	GitCherryPickItem struct {
		Hash       string
		Message    string
		Author     string
		FromBranch string
	}
)

func (i GitCherryPickItem) FilterValue() string {
	return i.Hash
}

func (d GitCherryPickDelegate) Height() int                             { return 2 }
func (d GitCherryPickDelegate) Spacing() int                            { return 0 }
func (d GitCherryPickDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitCherryPickDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitCherryPickItem)
	if !ok {
		return
	}

	var firstStr string
	var secondStr string
	var str string

	var selected bool
	if d.CherryPickedMap != nil {
		_, selected = (*d.CherryPickedMap)[i.Hash]
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
//	GitEditCherryPickDelegate renders each row in the edit-cherry-pick list,
//	showing hash, message, author, and source branch across three lines.
//	GitEditCherryPickItem carries the same fields as GitCherryPickItem and is
//	used when displaying the already-selected commits for reordering/removal.
//
// ------------------------------------
type (
	GitEditCherryPickDelegate struct{}
	GitEditCherryPickItem     struct {
		Hash       string
		Message    string
		Author     string
		FromBranch string
	}
)

func (i GitEditCherryPickItem) FilterValue() string {
	return i.Hash
}

func (d GitEditCherryPickDelegate) Height() int                             { return 3 }
func (d GitEditCherryPickDelegate) Spacing() int                            { return 0 }
func (d GitEditCherryPickDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitEditCherryPickDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitEditCherryPickItem)
	if !ok {
		return
	}

	var firstStr string
	var secondStr string
	var thirdStr string
	var str string

	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad

	firstStr = utils.TruncateString(fmt.Sprintf("%s  |  %s", i.Hash[:7], i.Author), componentWidth)
	secondStr = utils.TruncateString(fmt.Sprintf("    %s", i.Message), componentWidth)
	thirdStr = utils.TruncateString(fmt.Sprintf("    %s  %s", i18n.LANGUAGEMAPPING.CherryPickedFromBranch, i.FromBranch), componentWidth)

	// except for the first string, all other string should be faded
	secondStr = style.ItemStyle.Faint(true).Render(secondStr)
	thirdStr = style.ItemStyle.Faint(true).Render(thirdStr)

	str = fmt.Sprintf("%s\n%s\n%s", firstStr, secondStr, thirdStr)

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
//	CherryPickOpsOptionDelegate renders each operation-type row (name + info)
//	in the cherry-pick action selection list.
//	CherryPickOpsOptionItem carries the display name, info string, and the
//	cherry-pick operation type constant for a single action choice.
//
// ------------------------------------
type (
	CherryPickOpsOptionDelegate struct{}
	CherryPickOpsOptionItem     struct {
		Name              string
		Info              string
		CherryPickOpsType string
	}
)

func (i CherryPickOpsOptionItem) FilterValue() string {
	return i.Name
}

func (d CherryPickOpsOptionDelegate) Height() int                             { return 2 }
func (d CherryPickOpsOptionDelegate) Spacing() int                            { return 0 }
func (d CherryPickOpsOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d CherryPickOpsOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(CherryPickOpsOptionItem)
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
//	GitRevertParentOptionSelectionPopUpModel holds the parent-option list and
//	the target commit hash for the revert-parent selection popup. This popup
//	appears when reverting a merge commit that has more than one parent, so the
//	user can choose which parent to revert to.
//
// ------------------------------------
type GitRevertParentOptionSelectionPopUpModel struct {
	GitRevertParentOption list.Model
	CommitHash            string
}

// ------------------------------------
//
//	GitRevertParentOptionDelegate renders each parent-commit row (message +
//	hash) in the revert-parent selection list.
//	GitRevertParentOptionItem carries the parent commit hash, message, and
//	parent order for a single revert-parent choice.
//
// ------------------------------------
type (
	GitRevertParentOptionDelegate struct{}
	GitRevertParentOptionItem     struct {
		CommitHash    string
		CommitMessage string
		ParentOrder   int
	}
)

func (i GitRevertParentOptionItem) FilterValue() string {
	return i.CommitHash
}

func (d GitRevertParentOptionDelegate) Height() int                             { return 2 }
func (d GitRevertParentOptionDelegate) Spacing() int                            { return 0 }
func (d GitRevertParentOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitRevertParentOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitRevertParentOptionItem)
	if !ok {
		return
	}

	commitMsgStr := fmt.Sprintf("   %s", i.CommitMessage)
	commitHashStr := fmt.Sprintf("    %s", i.CommitHash)

	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad - 2

	commitMsgStr = utils.TruncateString(commitMsgStr, componentWidth)
	commitHashStr = utils.TruncateString(commitHashStr, componentWidth)

	nameRendered := style.ItemStyle.Render(commitMsgStr)
	infoRendered := style.ItemStyle.Faint(true).Render(commitHashStr)
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
//	GitRevertConfirmationPopUpModel holds the target commit hash and parent
//	order for the git revert confirmation popup, shown before the revert is
//	executed.
//
// ------------------------------------
type GitRevertConfirmationPopUpModel struct {
	CommitHash  string
	ParentOrder int
}
