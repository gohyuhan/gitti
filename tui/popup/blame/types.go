package blame

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/utils"
)

type BlamePoUpModel struct {
	CurrentGitTrackedFilesPathList list.Model
	BlameViewport                  viewport.Model
	FilterInput                    textinput.Model
	ShowingBlameInfo               bool
	HasFilePathChosen              bool
	FilterValue                    string
	SelectedFilePath               string
}

// ------------------------------------
//
//	Reset the popup back to file-selection state, clearing any chosen file and blame view
//
// ------------------------------------
func (bPM *BlamePoUpModel) ResetSelectedBlameFile() {
	bPM.ShowingBlameInfo = false
	bPM.HasFilePathChosen = false
	bPM.SelectedFilePath = ""
	bPM.BlameViewport.SetContent("")
}

// ------------------------------------
//
//	Display blame information view for selected file path
//
// ------------------------------------
func (bPM *BlamePoUpModel) ShowBlameInfoView(filePath string) {
	bPM.ShowingBlameInfo = true
	bPM.HasFilePathChosen = true
	bPM.SelectedFilePath = filePath
}

// ---------------------------------
//
// bubble tea list for selecting a git-tracked file to view blame info
//
// ---------------------------------
type (
	CurrentGitTrackedFilesPathDelegate struct{}
	CurrentGitTrackedFilesPathItem     struct {
		FilePath string
	}
)

func (i CurrentGitTrackedFilesPathItem) FilterValue() string {
	return i.FilePath
}

// list delegate interface implementation for CurrentGitTrackedFilesPathDelegate
func (d CurrentGitTrackedFilesPathDelegate) Height() int                             { return 1 }
func (d CurrentGitTrackedFilesPathDelegate) Spacing() int                            { return 0 }
func (d CurrentGitTrackedFilesPathDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d CurrentGitTrackedFilesPathDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(CurrentGitTrackedFilesPathItem)
	if !ok {
		return
	}
	componentWidth := m.Width() - constant.ListItemOrTitleWidthPad

	filePathStr := utils.TruncateString(i.FilePath, componentWidth)

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

	fmt.Fprint(w, fn(filePathStr))
}
