package branch

import (
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

// ---------------------------------
//
// create a new branch and remain on current branch
//
// ---------------------------------
type CreateNewBranchPopUpModel struct {
	NewBranchNameInput textinput.Model
	CreateType         string
}

// ---------------------------------
//
// choose on how to create the new branch, just create or create and move changes
//
// ---------------------------------
type ChooseNewBranchTypeOptionPopUpModel struct {
	NewBranchTypeOptionList list.Model
}

// ---------------------------------
//
// choose a switch type when switching branch
//
// ---------------------------------
type ChooseSwitchBranchTypePopUpModel struct {
	SwitchTypeOptionList list.Model
	BranchName           string
}

// ---------------------------------
//
// # A pop up to show branch switch result
//
// ---------------------------------
type SwitchBranchOutputPopUpModel struct {
	BranchName                 string // the branch name of the branch it was switching to
	SwitchType                 string
	SwitchBranchOutputViewport viewport.Model // to log out the output from git operation
	Spinner                    spinner.Model  // spinner for showing processing state
	IsProcessing               atomic.Bool    // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                   atomic.Bool    // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess             atomic.Bool    // has the process sucessfuly executed
}

// ---------------------------------
//
// for branch delete confirm prompt pop up
//
// ---------------------------------
type GitDeleteBranchConfirmPromptPopUpModel struct {
	BranchName string
}

// ---------------------------------
//
// for branch delete output result pop up
//
// ---------------------------------
type GitDeleteBranchOutputPopUpModel struct {
	BranchDeleteOutputViewport viewport.Model
	Spinner                    spinner.Model
	IsProcessing               atomic.Bool // indicator to prevent multiple thread spawning reacting to the key binding trigger
	HasError                   atomic.Bool // indicate if git commit exitcode is not 0 (meaning have error)
	ProcessSuccess             atomic.Bool // has the process sucessfuly executed
}

// ---------------------------------
//
// for new branch option selection option
//
// ---------------------------------
type (
	GitNewBranchTypeOptionDelegate struct{}
	GitNewBranchTypeOptionItem     struct {
		Name          string
		Info          string
		NewBranchType string
	}
)

func (i GitNewBranchTypeOptionItem) FilterValue() string {
	return i.Name
}

// for new branch creation type selection
func (d GitNewBranchTypeOptionDelegate) Height() int                             { return 1 }
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

// ---------------------------------
//
// for switch branch option selection option
//
// ---------------------------------
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

// for switch branch type selection
func (d GitSwitchBranchTypeOptionDelegate) Height() int                             { return 1 }
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
