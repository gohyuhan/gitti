package commitfiles

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/utils"
)

type (
	GitCommitFileItemDelegate struct{}
	GitCommitFileItem         struct {
		CommitHash   string
		Status       string // "M", "A", "D", "R", etc.
		FilePathname string
		OldPathname  string
	}
)

func (i GitCommitFileItem) FilterValue() string {
	if i.OldPathname != "" {
		return fmt.Sprintf("%s -> %s", i.OldPathname, i.FilePathname)
	}
	return i.FilePathname
}

func (d GitCommitFileItemDelegate) Height() int                             { return 1 }
func (d GitCommitFileItemDelegate) Spacing() int                            { return 0 }
func (d GitCommitFileItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d GitCommitFileItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GitCommitFileItem)
	if !ok {
		return
	}

	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad - 6
	filePathName := utils.TruncateString(i.FilePathname, componentWidth)

	var statusStyled string
	statusPrefix := i.Status
	if len(statusPrefix) > 0 {
		switch statusPrefix[0] {
		case 'A':
			statusStyled = style.StagedFileStyle.Render("A")
		case 'M':
			statusStyled = style.UnstagedFileStyle.Render("M")
		case 'D':
			statusStyled = style.ErrorStyle.Render("D")
		case 'R':
			statusStyled = style.NewStyle.Foreground(style.ColorCyanSoft).Render("R")
		default:
			statusStyled = style.ItemStyle.Render(string(statusPrefix[0]))
		}
	}

	str := fmt.Sprintf(" %s %s", statusStyled, filePathName)

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
