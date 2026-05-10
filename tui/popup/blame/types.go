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

// ------------------------------------
//
//	BlamePopUpModel holds the file-selection list, filter text input, and blame
//	output viewport for the blame popup. ShowingBlameInfo/HasFilePathChosen
//	control which view is active; SelectedFilePath tracks the current file.
//
// ------------------------------------
type BlamePopUpModel struct {
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
func (bPM *BlamePopUpModel) ResetSelectedBlameFile() {
	bPM.ShowingBlameInfo = false
	bPM.HasFilePathChosen = false
	bPM.SelectedFilePath = ""
	bPM.BlameViewport.SetContent("")
}

// ------------------------------------
//
//	Transition the blame popup to blame-view state for the given file path,
//	setting ShowingBlameInfo and HasFilePathChosen to true and storing the path.
//
// ------------------------------------
func (bPM *BlamePopUpModel) ShowBlameInfoView(filePath string) {
	bPM.ShowingBlameInfo = true
	bPM.HasFilePathChosen = true
	bPM.SelectedFilePath = filePath
}

// ------------------------------------
//
//	CurrentGitTrackedFilesPathDelegate renders each file path row as a single
//	line in the git-tracked file selection list inside the blame popup.
//	CurrentGitTrackedFilesPathItem wraps a single file path string.
//
// ------------------------------------
type (
	CurrentGitTrackedFilesPathDelegate struct{}
	CurrentGitTrackedFilesPathItem     struct {
		FilePath string
	}
)

func (i CurrentGitTrackedFilesPathItem) FilterValue() string {
	return i.FilePath
}

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
