package branch

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync/atomic"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	CreateNewBranchPopUpModel holds a focused text input for the new branch name,
//	the creation type constant, and an optional commit hash to branch from.
//
// ------------------------------------
type CreateNewBranchPopUpModel struct {
	NewBranchNameInput textinput.Model
	CreateType         string
	CommitHash         string
}

// ------------------------------------
//
//	ChooseNewBranchTypeOptionPopUpModel holds the option list for the new branch
//	type selection popup, offering four modes: create-only, create-and-switch,
//	create-from-remote-input, and create-from-remote-selection.
//
// ------------------------------------
type ChooseNewBranchTypeOptionPopUpModel struct {
	NewBranchTypeOptionList list.Model
}

// ------------------------------------
//
//	ChooseSwitchBranchTypePopUpModel holds the option list and the target branch
//	name for the switch type selection popup, offering clean-switch and
//	switch-with-changes modes.
//
// ------------------------------------
type ChooseSwitchBranchTypePopUpModel struct {
	SwitchTypeOptionList list.Model
	BranchName           string
}

// ------------------------------------
//
//	SwitchBranchOutputPopUpModel holds the scrollable output viewport and dot
//	spinner for the branch-switch progress popup, the target branch name, the
//	switch type, and atomic flags (IsProcessing, HasError, ProcessSuccess).
//
// ------------------------------------
type SwitchBranchOutputPopUpModel struct {
	BranchName                 string // the branch name of the branch it was switching to
	SwitchType                 string
	SwitchBranchOutputViewport viewport.Model // to log out the output from git operation
	Spinner                    spinner.Model  // spinner for showing processing state
	IsProcessing               atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                   atomic.Bool    // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess             atomic.Bool    // has the process sucessfuly executed
}

// ------------------------------------
//
//	GitDeleteBranchConfirmPromptPopUpModel holds the target branch name for the
//	deletion confirmation prompt popup.
//
// ------------------------------------
type GitDeleteBranchConfirmPromptPopUpModel struct {
	BranchName string
}

// ------------------------------------
//
//	GitDeleteBranchOutputPopUpModel holds the scrollable output viewport and dot
//	spinner for the branch-deletion progress popup, plus atomic flags
//	(IsProcessing, HasError, ProcessSuccess).
//
// ------------------------------------
type GitDeleteBranchOutputPopUpModel struct {
	BranchDeleteOutputViewport viewport.Model
	Spinner                    spinner.Model
	IsProcessing               atomic.Bool // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                   atomic.Bool // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess             atomic.Bool // has the process sucessfuly executed
}

// ------------------------------------
//
//	GitNewBranchTypeOptionDelegate renders each new-branch type row (name +
//	description) in the choose new branch type option list.
//	GitNewBranchTypeOptionItem carries the display name, description, and the
//	new branch type constant for a single creation mode choice.
//
// ------------------------------------
type (
	GitNewBranchTypeOptionDelegate struct{}
	GitNewBranchTypeOptionItem     struct {
		Name          string
		Info          string
		NewBranchType string
	}
)

// ------------------------------------
//
//	CreateBranchBasedOnRemotePopUpModel holds the remote origin name and a
//	focused text input for the user to type the remote branch name to track.
//
// ------------------------------------
type CreateBranchBasedOnRemotePopUpModel struct {
	RemoteOrigin          string // remote origin name
	RemoteBranchNameInput textinput.Model
}

// ------------------------------------
//
//	CreateBranchBasedOnRemoteOutputPopUpModel holds the scrollable output viewport
//	and dot spinner for the create-branch-from-remote progress popup, plus atomic
//	flags (IsProcessing, HasError, ProcessSuccess).
//
// ------------------------------------
type CreateBranchBasedOnRemoteOutputPopUpModel struct {
	CreateBranchBasedOnRemoteOutputViewport viewport.Model // to log out the output from git operation
	Spinner                                 spinner.Model  // spinner for showing processing state
	IsProcessing                            atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                                atomic.Bool    // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess                          atomic.Bool    // has the process sucessfuly executed
}

func (i GitNewBranchTypeOptionItem) FilterValue() string {
	return i.Name
}

func (d GitNewBranchTypeOptionDelegate) Height() int                             { return 2 }
func (d GitNewBranchTypeOptionDelegate) Spacing() int                            { return 0 }
func (d GitNewBranchTypeOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitNewBranchTypeOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitNewBranchTypeOptionItem)
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
//	GitSwitchBranchTypeOptionDelegate renders each switch type row (name +
//	description) in the choose switch branch type option list.
//	GitSwitchBranchTypeOptionItem carries the display name, description, and the
//	switch branch type constant for a single switch mode choice.
//
// ------------------------------------
type (
	GitSwitchBranchTypeOptionDelegate struct{}
	GitSwitchBranchTypeOptionItem     struct {
		Name             string
		Info             string
		SwitchBranchType string
	}
)

func (i GitSwitchBranchTypeOptionItem) FilterValue() string {
	return i.Name
}

func (d GitSwitchBranchTypeOptionDelegate) Height() int                             { return 2 }
func (d GitSwitchBranchTypeOptionDelegate) Spacing() int                            { return 0 }
func (d GitSwitchBranchTypeOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitSwitchBranchTypeOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitSwitchBranchTypeOptionItem)
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
//	ChooseRemoteBranchOptionPopUpModel holds the list of remote branches for the
//	user to select one to track as a new local branch.
//
// ------------------------------------
type ChooseRemoteBranchOptionPopUpModel struct {
	RemoteBranchOptionList list.Model
}

// ------------------------------------
//
//	RemoteBranchItemDelegate renders each remote branch as a single row in the
//	remote branch selection list.
//	RemoteBranchItem wraps a single remote branch name and its checked-out state.
//
// ------------------------------------
type (
	RemoteBranchItemDelegate struct{}
	RemoteBranchItem         struct {
		BranchName   string
		IsCheckedOut bool
	}
)

func (i RemoteBranchItem) FilterValue() string {
	return i.BranchName
}

func (d RemoteBranchItemDelegate) Height() int                             { return 1 }
func (d RemoteBranchItemDelegate) Spacing() int                            { return 0 }
func (d RemoteBranchItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d RemoteBranchItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(RemoteBranchItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("  %s", i.BranchName)

	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad

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
	str = utils.TruncateString(str, componentWidth)

	fmt.Fprint(w, fn(str))
}

// ------------------------------------
//
//	ChooseBranchOptionForMergePopUpModel holds two lists for the git merge branch
//	selection popup: one for available (unselected) branches and one for already-
//	selected branches, plus atomic flags tracking which panel is focused.
//
// ------------------------------------
type ChooseBranchOptionForMergePopUpModel struct {
	BranchOptionList              list.Model
	BranchOptionSectionSelected   atomic.Bool
	SelectedBranchList            list.Model
	SelectedBranchSectionSelected atomic.Bool
}

// ------------------------------------
//
//	GitMergeBranchOptionItemDelegate renders each branch name as a single row in
//	the merge branch selection lists (available and selected panels).
//	GitMergeBranchOptionItem wraps a single branch name for use in those lists.
//
// ------------------------------------
type (
	GitMergeBranchOptionItemDelegate struct{}
	GitMergeBranchOptionItem         struct {
		BranchName string
	}
)

func (i GitMergeBranchOptionItem) FilterValue() string {
	return i.BranchName
}

func (d GitMergeBranchOptionItemDelegate) Height() int                             { return 1 }
func (d GitMergeBranchOptionItemDelegate) Spacing() int                            { return 0 }
func (d GitMergeBranchOptionItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitMergeBranchOptionItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitMergeBranchOptionItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("  %s", i.BranchName)

	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad

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
	str = utils.TruncateString(str, componentWidth)

	fmt.Fprint(w, fn(str))
}

// ------------------------------------
//
//	BranchMergeOutputPopUpModel holds the scrollable output viewport and dot
//	spinner for the git merge progress popup, plus atomic flags (IsProcessing,
//	HasError, ProcessSuccess, IsCancelled) and a CancelFunc to abort the merge.
//
// ------------------------------------
type BranchMergeOutputPopUpModel struct {
	BranchMergeOutputViewport viewport.Model // to log out the output from git operation
	Spinner                   spinner.Model  // spinner for showing processing state
	IsProcessing              atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                  atomic.Bool    // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess            atomic.Bool    // has the process sucessfuly executed
	IsCancelled               atomic.Bool    // flag to indicate if the operation was cancelled by user
	// CancelFunc is used to cancel the git merge operation
	CancelFunc context.CancelFunc
}
