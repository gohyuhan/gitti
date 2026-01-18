package log

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ---------------------------------
//
// pop up for cherry picking commits from current checkout branch
//
// ---------------------------------
type GitCherryPickPopUpModel struct {
	CurrentBranchName                string
	CurrentBranchCherryPickCommitLog list.Model
}

// ---------------------------------
//
// pop up to edit (mainly removal for now) cherry picked commit
//
// ---------------------------------
type GitEditCherryPickPopUpModel struct {
	CherryPickedCommitLog list.Model
}

// ---------------------------------
//
// pop up to choose to either cherry pick or edit cherry pick
//
// ---------------------------------
type GitCherryPickOptionSelectionPopUpModel struct {
	CherryPickedOpsOption list.Model
}

// ---------------------------------
//
// bubble tea list for selecting commit log for cherry pick
//
// ---------------------------------
type (
	GitCherryPickDelegate struct {
		CherryPickedCommit *[]git.CherryPickedCommitLog
		CherryPickedMap    *map[string]string
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

// for list component of Git branch
func (d GitCherryPickDelegate) Height() int                             { return 1 }
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
	secondStr = utils.TruncateString(fmt.Sprintf("       %s", i.Message), componentWidth)

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
// bubble tea list for selecting operation choice for cherry pick
//
// ---------------------------------
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

// for list component of Git branch
func (d CherryPickOpsOptionDelegate) Height() int                             { return 1 }
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
