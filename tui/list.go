package tui

import (
	"fmt"
	"io"
	"strings"

	"gitti/tui/constant"
	"gitti/tui/style"
	"gitti/tui/utils"

	"github.com/charmbracelet/bubbles/v2/list"
	tea "github.com/charmbracelet/bubbletea/v2"
)

// -----------------------------------------------------------------------------
// implementation for list compoenent
// -----------------------------------------------------------------------------

// for list component of git branch
func (d gitBranchItemDelegate) Height() int                             { return 1 }
func (d gitBranchItemDelegate) Spacing() int                            { return 0 }
func (d gitBranchItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gitBranchItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(gitBranchItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("   %s", i.BranchName)
	if i.IsCheckedOut {
		str = fmt.Sprintf(" * %s", i.BranchName)
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

// for list component of modified files
func (d gitModifiedFilesItemDelegate) Height() int                             { return 1 }
func (d gitModifiedFilesItemDelegate) Spacing() int                            { return 0 }
func (d gitModifiedFilesItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gitModifiedFilesItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(gitModifiedFilesItem)
	if !ok {
		return
	}

	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad - 5
	filePathName := utils.TruncateString(i.FilePathname, componentWidth)

	indexState := style.StagedFileStyle.Render(i.IndexState)
	workTree := style.UnstagedFileStyle.Render(i.WorkTree)
	if i.IndexState == "?" {
		indexState = style.UnstagedFileStyle.Render(i.IndexState)
	}

	str := fmt.Sprintf(" %s%s %s", indexState, workTree, filePathName)

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

// for list component of stash
func (d gitStashItemDelegate) Height() int                             { return 1 }
func (d gitStashItemDelegate) Spacing() int                            { return 0 }
func (d gitStashItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gitStashItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(gitStashItem)
	if !ok {
		return
	}

	str := fmt.Sprintf(" %s", i.Message)

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

// for list component of git remote
func (d gitRemoteItemDelegate) Height() int                             { return 1 }
func (d gitRemoteItemDelegate) Spacing() int                            { return 0 }
func (d gitRemoteItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gitRemoteItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(gitRemoteItem)
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

// for push selection option
func (d gitPushOptionDelegate) Height() int                             { return 1 }
func (d gitPushOptionDelegate) Spacing() int                            { return 0 }
func (d gitPushOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gitPushOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(gitPushOptionItem)
	if !ok {
		return
	}

	nameStr := fmt.Sprintf("   %s", i.Name)
	urlStr := fmt.Sprintf("    %s", i.Info)

	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad - 2

	nameStr = utils.TruncateString(nameStr, componentWidth)
	urlStr = utils.TruncateString(urlStr, componentWidth)

	nameRendered := style.ItemStyle.Render(nameStr)
	infoRendered := style.ItemStyle.Faint(true).Render(urlStr)
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

// for new branch creation type selection
func (d gitNewBranchTypeOptionDelegate) Height() int                             { return 1 }
func (d gitNewBranchTypeOptionDelegate) Spacing() int                            { return 0 }
func (d gitNewBranchTypeOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gitNewBranchTypeOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(gitNewBranchTypeOptionItem)
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

// for switch branch type selection
func (d gitSwitchBranchTypeOptionDelegate) Height() int                             { return 1 }
func (d gitSwitchBranchTypeOptionDelegate) Spacing() int                            { return 0 }
func (d gitSwitchBranchTypeOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gitSwitchBranchTypeOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(gitSwitchBranchTypeOptionItem)
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

// for pull type selection
func (d gitPullTypeOptionDelegate) Height() int                             { return 1 }
func (d gitPullTypeOptionDelegate) Spacing() int                            { return 0 }
func (d gitPullTypeOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gitPullTypeOptionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(gitPullTypeOptionItem)
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
