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
type CommitCherryPickPopUpModel struct {
	CurrentBranchName      string
	CurrentBranchCommitLog list.Model
}

// ---------------------------------
//
// pop up to edit (mainly removal for now) cherry picked commit
//
// ---------------------------------
type CommitCherryPickEditPopUpModel struct {
	CherryPickedCommitLog list.Model
}

type (
	GitCherryPickDelegate struct {
		CherryPickedCommit *[]git.CherryPickedCommitLog
		CherryPickedMap    map[string]struct{}
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

	var str string

	var selected bool
	if d.CherryPickedMap != nil {
		_, selected = d.CherryPickedMap[i.Hash]
	}

	if selected {
		str = fmt.Sprintf("[X]  %s", i.Hash[:7])
	} else {
		str = fmt.Sprintf("[ ]  %s", i.Hash[:7])
	}

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
