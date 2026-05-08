package git

import (
	"cmp"
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/logging"
)

// represent the info of commit that got cherry picked
type CommitInfo struct {
	Hash        string
	Message     string
	Author      string
	Description string
	Parent      []string
	CommitOrder int // latest commit have the smallest integer value
}

type GitInteractiveRebase struct {
	gitCommitInfo     []CommitInfo // this is on demand, everytime a ops need it, we reassign a new one to prevent race conditions
	updateChannel     chan string
	maxCommitLogCount string
	gitProcessLock    *GitProcessLock
	logging           *logging.GittiLogging
}

// ------------------------------------
//
//	Init Git Interactive Rebase
//
// ------------------------------------
func InitGitInteractiveRebase(updateChannel chan string, gitProcessLock *GitProcessLock, maxCommitLogCountInt int, logging *logging.GittiLogging) *GitInteractiveRebase {
	maxCommitLogCount := strconv.Itoa(maxCommitLogCountInt)
	gitInteractiveRebase := GitInteractiveRebase{
		gitCommitInfo:     make([]CommitInfo, 0),
		gitProcessLock:    gitProcessLock,
		updateChannel:     updateChannel,
		maxCommitLogCount: maxCommitLogCount,
		logging:           logging,
	}
	return &gitInteractiveRebase
}

// ------------------------------------
//
//	Return commit log output
//
// ------------------------------------
func (gIR *GitInteractiveRebase) GitCommitInfo() []CommitInfo {
	copied := make([]CommitInfo, len(gIR.gitCommitInfo))
	copy(copied, gIR.gitCommitInfo)
	return copied
}

// ------------------------------------
//
//	Get the Commit Info for Interactive Rebase
//
// ------------------------------------
func (gIR *GitInteractiveRebase) GetCommitInfos() {
	gIR.clearGitCommitInfos()
	// 1. Prepare git command
	gitArgs := []string{
		"log",
		"--topo-order",
		"--no-decorate",
		"--no-notes",
		"--pretty=format:%H%x00%s%x00%an%x00%b%x00%P%x01",
		"-n", gIR.maxCommitLogCount,
		"--",
	}

	var commitInfos []CommitInfo

	commitInfoCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	commitInfoOutput, commitInfoErr := commitInfoCmdExecutor.Output()
	if commitInfoErr != nil {
		gIR.gitCommitInfo = commitInfos
		return
	}

	parsedCommitInfoArray := strings.Split(string(commitInfoOutput), "\x01")
	for _, commitInfo := range parsedCommitInfoArray {
		parsedCommitInfo := strings.SplitN(commitInfo, SEPARATOR, 5)
		if len(parsedCommitInfo) < 5 {
			continue
		}
		var commitParent []string
		if utf8.RuneCountInString(parsedCommitInfo[4]) > 0 {
			commitParent = strings.Fields(parsedCommitInfo[4])
		}
		commitInfos = append(commitInfos, CommitInfo{
			Hash:        strings.TrimSpace(parsedCommitInfo[0]),
			Message:     strings.TrimSpace(parsedCommitInfo[1]),
			Author:      strings.TrimSpace(parsedCommitInfo[2]),
			Description: strings.TrimSpace(parsedCommitInfo[3]),
			Parent:      commitParent,
			CommitOrder: len(commitInfos),
		})
	}

	gIR.gitCommitInfo = commitInfos
}

// ------------------------------------
//
//	Interactive Rebase - Fixup
//	Squashes all selected commits into one with a new message. No editor is opened.
//
//	* sortedSelectedCommitInfos must be sorted oldest → latest (largest CommitOrder first)
//	* at least 2 commits must be selected
//	* the oldest (base) selected commit must not be a merge commit
//	* merge commits inside the selected range are silently skipped and dropped
//
//	Flow:
//	* builds affected range: HEAD down to oldest selected commit
//	* todo: pick oldest-selected → fixup rest-of-selected → exec amend with new message → pick non-selected
//	* GIT_SEQUENCE_EDITOR is a temp shell script that replaces the todo file non-interactively
//	* message applied via `git commit --amend -F <tempfile>` — no editor flag
//
// ------------------------------------
func (gIR *GitInteractiveRebase) GitInteractiveRebaseFixup(ctx context.Context, sortedSelectedCommitInfos []CommitInfo, newCommitMessage string, newCommitDesceription string) ([]string, error) {
	if !gIR.gitProcessLock.CanProceedWithGitOps() {
		return []string{}, fmt.Errorf("%s", gIR.gitProcessLock.OtherProcessRunningWarning())
	}
	defer func() {
		gIR.gitProcessLock.ReleaseGitOpsLock()
	}()

	fixupCmd, fixupCleanup, fixupErr := gIR.interactiveRebaseFixup(ctx, sortedSelectedCommitInfos, newCommitMessage, newCommitDesceription, false)
	if fixupErr != nil {
		return []string{}, fixupErr
	}
	if fixupCleanup != nil {
		defer fixupCleanup()
	}

	fixupOutput, runErr := fixupCmd.CombinedOutput()
	parsedFixupOutput := processGeneralGitOpsOutputIntoStringArray(fixupOutput)

	if runErr != nil {
		if ctx.Err() != nil {
			return parsedFixupOutput, ctx.Err()
		}
		return parsedFixupOutput, runErr
	}

	return parsedFixupOutput, nil
}

// GitInteractiveRebaseFixupWithSigning is the signing variant of GitInteractiveRebaseFixup.
// ------------------------------------
//
//	Interactive Rebase - Fixup (Signing variant)
//	Returns the rebase cmd without running it — caller (bubbletea) suspends the TUI,
//	hands the cmd to the OS, and resumes after git finishes.
//
//	* gitti lock is released before the cmd runs — intentional
//	* git's own repo lock (.git/rebase-merge/) prevents concurrent ops during execution
//
// ------------------------------------
func (gIR *GitInteractiveRebase) GitInteractiveRebaseFixupWithSigning(ctx context.Context, sortedSelectedCommitInfos []CommitInfo, newCommitMessage string, newCommitDesceription string) (*exec.Cmd, func(), error) {
	if !gIR.gitProcessLock.CanProceedWithGitOps() {
		return nil, nil, fmt.Errorf("%s", gIR.gitProcessLock.OtherProcessRunningWarning())
	}
	defer func() {
		gIR.gitProcessLock.ReleaseGitOpsLock()
	}()

	return gIR.interactiveRebaseFixup(ctx, sortedSelectedCommitInfos, newCommitMessage, newCommitDesceription, true)
}

func (gIR *GitInteractiveRebase) interactiveRebaseFixup(ctx context.Context, sortedSelectedCommitInfos []CommitInfo, newCommitMessage string, newCommitDesceription string, signing bool) (*exec.Cmd, func(), error) {
	if len(sortedSelectedCommitInfos) < 2 {
		return nil, nil, fmt.Errorf("%s", i18n.LANGUAGEMAPPING.InteractiveRebaseFixupMustHaveAtLeastTwoSelectedError)
	}

	// base is merge commit, does not allow that
	if len(sortedSelectedCommitInfos[0].Parent) > 1 {
		return nil, nil, fmt.Errorf("%s", i18n.LANGUAGEMAPPING.InteractiveRebaseFixupBaseSelectionMustNotBeMergeCommitError)
	}

	// oldest selected commit's CommitOrder equals its index in gitCommitInfo (which is latest-first)
	// so gitCommitInfo[:targetCommitPosition+1] gives us HEAD..oldest-selected
	targetCommitPosition := sortedSelectedCommitInfos[0].CommitOrder
	if targetCommitPosition < 0 || targetCommitPosition >= len(gIR.gitCommitInfo) {
		return nil, nil, fmt.Errorf("%s", i18n.LANGUAGEMAPPING.InteractiveRebaseFixupPositionMismatchError)
	}

	// copy so we don't mutate gitCommitInfo, then sort oldest → latest for todo construction
	sortedAffectedCommitInfos := make([]CommitInfo, targetCommitPosition+1)
	copy(sortedAffectedCommitInfos, gIR.gitCommitInfo[:targetCommitPosition+1])
	slices.SortFunc(sortedAffectedCommitInfos, func(a, b CommitInfo) int {
		return cmp.Compare(b.CommitOrder, a.CommitOrder) // largest CommitOrder first = oldest to latest
	})

	// build the `exec git commit --amend -F <tempfile>` line that applies the new message
	execAmendCommitCommandString, commitMsgTempPath, execErr := gIR.buildRebaseAmendExec(newCommitMessage, newCommitDesceription)
	if execErr != nil {
		return nil, nil, execErr
	}
	fixupTodoString := gIR.constructFixupTodo(sortedAffectedCommitInfos, sortedSelectedCommitInfos, execAmendCommitCommandString)

	// write todo to a temp file; the sequence editor script will cp it into git's todo path
	sequenceEditorScriptPath, cleanupFn, buildCmdErr := gIR.buildNonInteractiveFixupRebaseCmd(fixupTodoString)
	if buildCmdErr != nil {
		os.Remove(commitMsgTempPath)
		return nil, cleanupFn, buildCmdErr
	}
	// chain commitMsgTempPath removal into the existing cleanup so all temp files are swept together
	if cleanupFn != nil {
		oldCleanup := cleanupFn
		cleanupFn = func() {
			oldCleanup()
			os.Remove(commitMsgTempPath)
		}
	} else {
		cleanupFn = func() {
			os.Remove(commitMsgTempPath)
		}
	}

	// rebase from just before the oldest affected commit; use --root only if it has no parent
	gitArgs := []string{
		"rebase",
		"-i",
		"--root",
	}
	if len(sortedAffectedCommitInfos[0].Parent) > 0 {
		gitArgs = []string{
			"rebase",
			"-i",
			sortedAffectedCommitInfos[0].Parent[0],
		}
	}
	// signing path uses no context — bubbletea owns execution, gitti must not cancel it
	rebaseCmd := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, false)
	if signing {
		rebaseCmd = executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	}
	// override git's sequence editor with our script so no interactive editor is ever opened
	rebaseCmd.Env = append(rebaseCmd.Env, fmt.Sprintf("GIT_SEQUENCE_EDITOR=%s", sequenceEditorScriptPath))

	return rebaseCmd, cleanupFn, nil
}

// ------------------------------------
//
//	Construct Fixup Todo
//	Builds the rebase todo file content for the fixup operation.
//
//	Todo structure (oldest → latest order):
//	  pick  <oldest-selected>          ← anchor: becomes the squashed commit
//	  fixup <each-newer-selected>      ← absorbed into pick; message discarded (no editor)
//	  exec git commit --amend -F <f>   ← replace squashed message with user's new message
//	  pick  <non-selected-1>           ← replayed on top unchanged
//	  pick  <non-selected-N>
//
//	* merge commits in selected range are skipped (dropped from rebase)
//	* merge commits in non-selected range are also skipped (cannot be replayed by pick)
//
// ------------------------------------
func (gIR *GitInteractiveRebase) constructFixupTodo(sortedAffectedCommitInfos []CommitInfo, sortedSelectedCommitInfos []CommitInfo, execAmendCommitCommandString string) string {
	var fixupTodoString strings.Builder

	// oldest selected commit is the rebase anchor — `pick` preserves it as the surviving commit
	fixupTodoString.WriteString("pick ")
	fixupTodoString.WriteString(sortedSelectedCommitInfos[0].Hash)
	fixupTodoString.WriteRune('\n')

	// remaining selected commits are absorbed via `fixup` — their messages are silently discarded
	for _, commitInfo := range sortedSelectedCommitInfos[1:] {
		// skip any merge commit
		if len(commitInfo.Parent) > 1 {
			continue
		}
		fixupTodoString.WriteString("fixup ")
		fixupTodoString.WriteString(commitInfo.Hash)
		fixupTodoString.WriteRune('\n')
	}

	// after fixup the surviving commit has the oldest selected commit's message;
	// `exec` replaces it with the user-supplied message via --amend -F (no editor opened)
	fixupTodoString.WriteString(execAmendCommitCommandString)
	fixupTodoString.WriteRune('\n')

	// build set of selected hashes for O(1) exclusion below
	selectedSet := make(map[string]struct{}, len(sortedSelectedCommitInfos))
	for _, c := range sortedSelectedCommitInfos {
		selectedSet[c.Hash] = struct{}{}
	}

	// non-selected commits in the affected range are replayed unchanged on top of the squashed commit
	filtered := make([]CommitInfo, 0, len(sortedAffectedCommitInfos))
	for _, c := range sortedAffectedCommitInfos {
		if _, found := selectedSet[c.Hash]; !found {
			filtered = append(filtered, c)
		}
	}

	for _, commitInfo := range filtered {
		// skip any merge commit
		if len(commitInfo.Parent) > 1 {
			continue
		}
		fixupTodoString.WriteString("pick ")
		fixupTodoString.WriteString(commitInfo.Hash)
		fixupTodoString.WriteRune('\n')
	}

	return fixupTodoString.String()
}

// ------------------------------------
//
//	Build Rebase Amend Exec Line
//	Writes the new commit message to a temp file and returns the `exec` line for the todo.
//
//	* description appended after a blank line (git paragraph format) if non-empty
//	* path is single-quote-escaped so spaces/special chars survive the shell exec line
//	* caller owns the temp file lifetime — returned path must be removed after rebase finishes
//
// ------------------------------------
func (gIR *GitInteractiveRebase) buildRebaseAmendExec(message, description string) (string, string, error) {
	// combine message + description
	fullMsg := message
	if strings.TrimSpace(description) != "" {
		fullMsg += "\n\n" + description
	}
	fullMsg += "\n"

	// create temp file
	f, err := os.CreateTemp("", "gitti-commit-msg-*")
	if err != nil {
		return "", "", err
	}

	// write content
	if _, err := f.WriteString(fullMsg); err != nil {
		f.Close()
		return "", "", err
	}

	if err := f.Close(); err != nil {
		return "", "", err
	}

	// shell escape path (safe for rebase exec)
	escapedPath := "'" + strings.ReplaceAll(f.Name(), "'", `'\''`) + "'"

	// construct exec string
	execStr := fmt.Sprintf(
		"exec git commit --amend -F %s",
		escapedPath,
	)

	return execStr, f.Name(), nil
}

// ------------------------------------
//
//	Build Non-Interactive Fixup Rebase Cmd
//	Writes the todo string to a temp file, then creates the sequence editor script.
//
//	* sequence editor script: `cp <todo-temp> "$1"` — replaces git's generated todo with ours
//	* GIT_SEQUENCE_EDITOR points to this script so git never opens an interactive editor
//	* returned cleanup removes both the todo temp file and the script temp file
//
// ------------------------------------
func (gIR *GitInteractiveRebase) buildNonInteractiveFixupRebaseCmd(fixupTodoString string) (string, func(), error) {
	todoFile, todoFileErr := os.CreateTemp("", "gitti-rebase-todo-*")
	if todoFileErr != nil {
		return "", nil, todoFileErr
	}

	if _, writeTodoFileErr := todoFile.WriteString(fixupTodoString); writeTodoFileErr != nil {
		todoFile.Close()
		os.Remove(todoFile.Name())
		return "", nil, writeTodoFileErr
	}
	if closeTodoFileErr := todoFile.Close(); closeTodoFileErr != nil {
		os.Remove(todoFile.Name())
		return "", nil, closeTodoFileErr
	}

	sequenceEditorScriptPath, sequenceEditorScriptErr := gIR.createSequenceEditorScript(todoFile.Name())
	if sequenceEditorScriptErr != nil {
		os.Remove(todoFile.Name())
		return "", nil, sequenceEditorScriptErr
	}

	cleanupFn := func() {
		os.Remove(sequenceEditorScriptPath)
		os.Remove(todoFile.Name())
	}

	return sequenceEditorScriptPath, cleanupFn, nil
}

// ------------------------------------
//
//	Create Sequence Editor Script
//	Emits a #!/bin/sh script that copies our pre-built todo over git's generated todo path.
//
//	* .sh extension used on all platforms — git-for-windows routes GIT_SEQUENCE_EDITOR
//	  through MSYS2's sh.exe, so .bat files cannot be used here
//	* chmod 0o700 required; git exec-calls the script directly by path
//
// ------------------------------------
func (gIR *GitInteractiveRebase) createSequenceEditorScript(todoFilePath string) (string, error) {
	tempScriptFile, createErr := os.CreateTemp("", "gitti-sequence-editor-*.sh")
	if createErr != nil {
		return "", createErr
	}

	scriptPath := tempScriptFile.Name()
	escapedTodoPath := "'" + strings.ReplaceAll(todoFilePath, "'", `'\''`) + "'"
	scriptContent := "#!/bin/sh\ncp " + escapedTodoPath + " \"$1\"\n"

	if _, writeErr := tempScriptFile.WriteString(scriptContent); writeErr != nil {
		tempScriptFile.Close()
		return "", writeErr
	}
	if closeErr := tempScriptFile.Close(); closeErr != nil {
		return "", closeErr
	}

	if chmodErr := os.Chmod(scriptPath, 0o700); chmodErr != nil {
		return "", chmodErr
	}

	return scriptPath, nil
}

// ------------------------------------
//
//	gIR internal clear git commit info func
//
// ------------------------------------
func (gIR *GitInteractiveRebase) clearGitCommitInfos() {
	gIR.gitCommitInfo = []CommitInfo{}
}
