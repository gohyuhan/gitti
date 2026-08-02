package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
	"golang.org/x/text/width"
)

// ------------------------------------
//
//	TruncateString trims string s to fit within given display width,
//	accounting for wide CJK characters, and appends "…" if truncated.
//
// ------------------------------------
func TruncateString(s string, maxWidth int) string {
	displayWidth := 0
	runes := []rune(s)
	var result []rune

	for _, r := range runes {
		prop := width.LookupRune(r)
		k := 1
		if prop.Kind() == width.EastAsianWide || prop.Kind() == width.EastAsianFullwidth {
			k = 2
		}

		if displayWidth+k > maxWidth {
			break
		}

		displayWidth += k
		result = append(result, r)
	}

	if len(result) < len(runes) {
		// Add ellipsis if we have room
		if len(result) >= 2 {
			result = append(result[:len(result)-1], '…')
		} else if len(result) == 1 {
			result = []rune{'…'}
		}
	}

	return string(result)
}

// ------------------------------------
//
//	CurrentPanelFilterKey returns the panel filter query map key for the list that is
//	currently focused (resolving the shared branch/tag/remote/worktree and
//	commitlog/reflog panel slots to the component actually showing). Returns ""
//	when the focused component has no filterable list.
//
// ------------------------------------
func CurrentPanelFilterKey(m *types.GittiModel) string {
	switch m.CurrentSelectedComponent {
	case constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel:
		return m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing
	case constant.ModifiedFilesComponentPanel:
		return constant.ModifiedFilesComponentPanel
	case constant.CommitLogOrRefLogComponentPanel:
		return m.CurrentCommitLogOrRefLogComponentShowing
	case constant.StashComponentPanel:
		return constant.StashComponentPanel
	}
	return ""
}

// ------------------------------------
//
//	FilterListItems filters list items by a case-insensitive substring match of the
//	query against each item's FilterValue(). It also recomputes the position of the
//	previously selected item within the filtered result (matched by FilterValue, as
//	some item types are not comparable). An empty query returns the input unchanged.
//
// ------------------------------------
func FilterListItems(items []list.Item, query string, previousSelected list.Item, previousPosition int) ([]list.Item, int) {
	if query == "" {
		return items, previousPosition
	}
	loweredQuery := strings.ToLower(query)
	filtered := make([]list.Item, 0, len(items))
	position := -1
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.FilterValue()), loweredQuery) {
			if previousSelected != nil && item.FilterValue() == previousSelected.FilterValue() {
				position = len(filtered)
			}
			filtered = append(filtered, item)
		}
	}
	return filtered, position
}

// ------------------------------------
//
//	ListCounterHelper returns a function that generates a counter display e.g. ("3/10")
//	showing the current item position in the list for the main left panel. When a panel
//	filter query is set for filterKey, it is appended as "/query" (with a trailing block
//	cursor while the query is being typed).
//
// ------------------------------------
func ListCounterHelper(m *types.GittiModel, list *list.Model, filterKey string) func() []key.Binding {
	return func() []key.Binding {
		currentIndex := list.Index() + 1
		totalCount := len(list.Items())
		countStr := fmt.Sprintf("%d/%d", currentIndex, totalCount)
		if totalCount == 0 {
			countStr = "0/0"
		}
		if m.IsPanelFiltering.Load() && CurrentPanelFilterKey(m) == filterKey {
			countStr = fmt.Sprintf("%s  /%s█", countStr, m.PanelFilterQuery[filterKey])
		} else if query := m.PanelFilterQuery[filterKey]; query != "" {
			countStr = fmt.Sprintf("%s  /%s", countStr, query)
		}
		countStr = TruncateString(countStr, m.WindowLeftPanelWidth-constant.ListItemOrTitleWidthPad-2)
		return []key.Binding{
			key.NewBinding(
				key.WithKeys(countStr),
				key.WithHelp(countStr, ""),
			),
		}
	}
}

// ------------------------------------
//
//	PopUpListCounterHelper returns a function that generates a counter display e.g. ("3/10")
//	showing the current item position in the list for pop-up dialogs.
//
// ------------------------------------
func PopUpListCounterHelper(m *types.GittiModel, list *list.Model, maxWidth int) func() []key.Binding {
	return func() []key.Binding {
		currentIndex := list.Index() + 1
		totalCount := len(list.Items())
		countStr := fmt.Sprintf("%d/%d", currentIndex, totalCount)
		width := (min(maxWidth, int(float64(m.Width)*0.8)) - 4)
		countStr = TruncateString(countStr, width-constant.ListItemOrTitleWidthPad-2)
		if totalCount == 0 {
			countStr = "0/0"
		}
		return []key.Binding{
			key.NewBinding(
				key.WithKeys(countStr),
				key.WithHelp(countStr, ""),
			),
		}
	}
}

// ------------------------------------
//
//	ReturnEditorLaunchCommand creates a command to launch the specified editor with the given file.
//	Returns the exec.Cmd to run and a boolean indicating if it's a non-terminal editor (true for GUI editors like VS Code).
//
// ------------------------------------
func ReturnEditorLaunchCommand(fileName string, userSetEditor string) (*exec.Cmd, bool) {
	filepath := "."
	if fileName != "" {
		filepath = fileName
	}
	var isNonTerminalEditor bool
	var editorCommand string

	switch strings.ToLower(userSetEditor) {
	case "nano":
		editorCommand = "nano"
		isNonTerminalEditor = false
	case "vim":
		editorCommand = "vim"
		isNonTerminalEditor = false
	case "neovim":
		editorCommand = "nvim"
		isNonTerminalEditor = false
	case "vscode":
		editorCommand = "code"
		isNonTerminalEditor = true
	case "zed":
		editorCommand = "zed"
		isNonTerminalEditor = true
	case "cursor":
		editorCommand = "cursor"
		isNonTerminalEditor = true
	case "windsurf":
		editorCommand = "windsurf"
		isNonTerminalEditor = true
	case "antigravity":
		editorCommand = "antigravity"
		isNonTerminalEditor = true
	default:
		editorCommand = "vi"
		isNonTerminalEditor = false
	}

	cmd := exec.Command(editorCommand, []string{filepath}...)
	return cmd, isNonTerminalEditor
}

// ------------------------------------
//
//	ReinitCherryPickedCommitInfo resets the cherry-pick tracking information in the model,
//	clearing all previously stored cherry-pick data.
//
// ------------------------------------
func ReinitCherryPickedCommitInfo(m *types.GittiModel) {
	m.CherryPickedCommitInfo.LatestSequenceCounter = 0
	m.CherryPickedCommitInfo.CherryPickedCommitMap = make(map[string]git.CherryPickedCommitLog)
}

// ------------------------------------
//
//	Suspend the Gitti UI and hand control to a git command that requires GPG
//	signing (e.g. signed commits/tags). Runs the command via tea.ExecProcess,
//	captures stderr silently (no direct terminal passthrough), sanitizes it, and
//	returns a GitOperationRequiredSigningFinishedMsg with the cleaned error output.
//
// ------------------------------------
func SuspendGittiUIForGitOperationRequireSigning(m *types.GittiModel, gitCommand []string, GitOperationOpsTypeForLogging string) (*types.GittiModel, tea.Cmd) {
	cmd := exec.Command("git", gitCommand...)
	var stderr bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	m.GittiLogger.RegisterNewLog(GitOperationOpsTypeForLogging, strings.Join(gitCommand, " "), logging.INFO, "", true)
	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		return buildSigningFinishedMsg(GitOperationOpsTypeForLogging, err, sanitizeGitSigningStderr(strings.TrimSpace(stderr.String())))
	})
}

// ------------------------------------
//
//	Suspend the Gitti UI and hand control to a pre-built exec.Cmd that requires
//	GPG signing. Runs cleanUpFunc after the process exits regardless of outcome,
//	captures stderr silently (no direct terminal passthrough), sanitizes it, then
//	returns a GitOperationRequiredSigningFinishedMsg with the cleaned error output.
//
// ------------------------------------
func SuspendGittiUIForGitOperationRequireSigningWithExecAndCleanUp(m *types.GittiModel, executor *exec.Cmd, cleanUpFunc func(), GitOperationOpsTypeForLogging string) (*types.GittiModel, tea.Cmd) {
	var stderr bytes.Buffer
	executor.Stderr = io.MultiWriter(os.Stderr, &stderr)

	m.GittiLogger.RegisterNewLog(GitOperationOpsTypeForLogging, strings.Join(executor.Args, " "), logging.INFO, "", true)
	return m, tea.ExecProcess(executor, func(err error) tea.Msg {
		cleanUpFunc()
		return buildSigningFinishedMsg(GitOperationOpsTypeForLogging, err, sanitizeGitSigningStderr(strings.TrimSpace(stderr.String())))
	})
}

var (
	ansiEscapePattern     *regexp.Regexp
	ansiEscapePatternOnce sync.Once
)

// ------------------------------------
//
//	Lazily initialize and return the compiled regex for ANSI escape sequence stripping.
//
// ------------------------------------
func getAnsiEscapePattern() *regexp.Regexp {
	ansiEscapePatternOnce.Do(func() {
		ansiEscapePattern = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)
	})
	return ansiEscapePattern
}

// ------------------------------------
//
//	Normalize git signing stderr for safe TUI rendering. Strips ANSI escapes,
//	converts CR/CRLF to LF, removes null bytes, and compresses repeated blank lines.
//
// ------------------------------------
func sanitizeGitSigningStderr(raw string) string {
	if raw == "" {
		return ""
	}
	raw = getAnsiEscapePattern().ReplaceAllString(raw, "")
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	raw = strings.ReplaceAll(raw, "\x00", "")

	lines := strings.Split(raw, "\n")
	cleaned := make([]string, 0, len(lines))
	prevEmpty := false
	for _, line := range lines {
		line = strings.TrimRight(line, " \t")
		isEmpty := strings.TrimSpace(line) == ""
		if isEmpty {
			if prevEmpty {
				continue
			}
			prevEmpty = true
			cleaned = append(cleaned, "")
			continue
		}
		prevEmpty = false
		cleaned = append(cleaned, line)
	}

	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}

// ------------------------------------
//
//	Build a GitOperationRequiredSigningFinishedMsg, appending stderr to the error
//	if present. Used by both signing functions to consolidate error-building logic.
//
// ------------------------------------
func buildSigningFinishedMsg(opsType string, err error, stderrStr string) types.GitOperationRequiredSigningFinishedMsg {
	if stderrStr != "" && err != nil {
		err = fmt.Errorf("%w\n\n%s", err, stderrStr)
	}
	return types.GitOperationRequiredSigningFinishedMsg{
		GitOperationOpsTypeForLogging: opsType,
		Err:                           err,
	}
}

// ------------------------------------
//
//	Clear the active popup after a GPG-signing operation completes. Hides the
//	popup, exits typing mode, and resets PopUpModel and PopUpType to their
//	no-popup defaults.
//
// ------------------------------------
func ResetPopUpModelStateForGitSigningOps(m *types.GittiModel) {
	m.ShowPopUp.Store(false)
	m.IsTyping.Store(false)
	m.PopUpModel = nil
	m.PopUpType = constant.NoPopUp
}
