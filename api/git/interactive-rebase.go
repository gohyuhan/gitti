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

// ------------------------------------
//
//	Represents one commit entry used by the interactive rebase fixup/squash flow
//
// ------------------------------------
type CommitInfo struct {
	Hash        string
	Message     string
	Author      string
	Description string
	Parent      []string
	CommitOrder int // latest commit have the smallest integer value
}

type GitInteractiveRebase struct {
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
		gitProcessLock:    gitProcessLock,
		updateChannel:     updateChannel,
		maxCommitLogCount: maxCommitLogCount,
		logging:           logging,
	}
	return &gitInteractiveRebase
}

// ------------------------------------
//
//	Get the Commit Info for Interactive Rebase
//
// ------------------------------------
func (gIR *GitInteractiveRebase) GetCommitInfos() []CommitInfo {
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
		return commitInfos
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

	return commitInfos
}

// *************************************************************************************
//                        INTERACTIVE REBASE - FIXUP / SQUASH
// *************************************************************************************

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
func (gIR *GitInteractiveRebase) GitInteractiveRebaseFixupSquash(ctx context.Context, gitCommitInfo []CommitInfo, sortedSelectedCommitInfos []CommitInfo, newCommitMessage string, newCommitDescription string) ([]string, error) {
	if !gIR.gitProcessLock.CanProceedWithGitOps() {
		return []string{}, fmt.Errorf("%s", gIR.gitProcessLock.OtherProcessRunningWarning())
	}
	defer func() {
		gIR.gitProcessLock.ReleaseGitOpsLock()
	}()

	fixupCmd, fixupCleanup, fixupErr := gIR.interactiveRebaseFixupSquash(ctx, gitCommitInfo, sortedSelectedCommitInfos, newCommitMessage, newCommitDescription, false)
	if fixupErr != nil {
		return []string{}, fixupErr
	}
	if fixupCleanup != nil {
		defer fixupCleanup()
	}
	gIR.logging.RegisterNewLog(logging.INTERACTIVE_REBASE_FIXUP_SQUASH, strings.Join(fixupCmd.Args, " "), logging.INFO, "", true)
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
func (gIR *GitInteractiveRebase) GitInteractiveRebaseFixupSquashWithSigning(ctx context.Context, gitCommitInfo []CommitInfo, sortedSelectedCommitInfos []CommitInfo, newCommitMessage string, newCommitDescription string) (*exec.Cmd, func(), error) {
	if !gIR.gitProcessLock.CanProceedWithGitOps() {
		return nil, nil, fmt.Errorf("%s", gIR.gitProcessLock.OtherProcessRunningWarning())
	}
	defer func() {
		gIR.gitProcessLock.ReleaseGitOpsLock()
	}()

	return gIR.interactiveRebaseFixupSquash(ctx, gitCommitInfo, sortedSelectedCommitInfos, newCommitMessage, newCommitDescription, true)
}

// ------------------------------------
//
//	Builds the rebase command and cleanup callback for fixup/squash; validates selected commits,
//	generates todo and message temp files, and configures a non-interactive sequence editor
//
// ------------------------------------
func (gIR *GitInteractiveRebase) interactiveRebaseFixupSquash(ctx context.Context, gitCommitInfo []CommitInfo, sortedSelectedCommitInfos []CommitInfo, newCommitMessage string, newCommitDescription string, signing bool) (*exec.Cmd, func(), error) {
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
	if targetCommitPosition < 0 || targetCommitPosition >= len(gitCommitInfo) {
		return nil, nil, fmt.Errorf("%s", i18n.LANGUAGEMAPPING.InteractiveRebaseFixupPositionMismatchError)
	}

	// copy so we don't mutate gitCommitInfo, then sort oldest → latest for todo construction
	sortedAffectedCommitInfos := make([]CommitInfo, targetCommitPosition+1)
	copy(sortedAffectedCommitInfos, gitCommitInfo[:targetCommitPosition+1])
	slices.SortFunc(sortedAffectedCommitInfos, func(a, b CommitInfo) int {
		return cmp.Compare(b.CommitOrder, a.CommitOrder) // largest CommitOrder first = oldest to latest
	})

	// build the `exec git commit --amend -F <tempfile>` line that applies the new message
	execAmendCommitCommandString, commitMsgTempPath, execErr := gIR.buildRebaseAmendExec(newCommitMessage, newCommitDescription)
	if execErr != nil {
		return nil, nil, execErr
	}
	fixupTodoString := gIR.constructFixupTodo(sortedAffectedCommitInfos, sortedSelectedCommitInfos, execAmendCommitCommandString)

	// write todo to a temp file; the sequence editor script will cp it into git's todo path
	sequenceEditorScriptPath, cleanupFn, buildCmdErr := gIR.buildNonInteractiveTodoCmd(fixupTodoString)
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

// *************************************************************************************
//                           INTERACTIVE REBASE - REWORD
// *************************************************************************************

// ------------------------------------
//
//	Interactive Rebase - Reword
//	Renames a single commit's message in-place. No editor is opened.
//
//	* the selected commit must not be a merge commit
//	* all commits from HEAD down to the selected commit are replayed
//	* merge commits in the affected range are silently skipped and dropped
//
//	Flow:
//	* builds affected range: HEAD down to selected commit
//	* todo: pick selected → exec amend with new message → pick newer commits
//	* GIT_SEQUENCE_EDITOR is a temp shell script that replaces the todo file non-interactively
//	* message applied via `git commit --amend -F <tempfile>` — no editor flag
//
// ------------------------------------
func (gIR *GitInteractiveRebase) GitInteractiveRebaseReword(ctx context.Context, gitCommitInfo []CommitInfo, selectedCommitInfo CommitInfo, newCommitMessage string, newCommitDescription string) ([]string, error) {
	if !gIR.gitProcessLock.CanProceedWithGitOps() {
		return []string{}, fmt.Errorf("%s", gIR.gitProcessLock.OtherProcessRunningWarning())
	}
	defer func() {
		gIR.gitProcessLock.ReleaseGitOpsLock()
	}()

	rewordCmd, rewordCleanup, rewordErr := gIR.interactiveRebaseReword(ctx, gitCommitInfo, selectedCommitInfo, newCommitMessage, newCommitDescription, false)
	if rewordErr != nil {
		return []string{}, rewordErr
	}
	if rewordCleanup != nil {
		defer rewordCleanup()
	}
	gIR.logging.RegisterNewLog(logging.INTERACTIVE_REBASE_REWORD, strings.Join(rewordCmd.Args, " "), logging.INFO, "", true)
	rewordOutput, runErr := rewordCmd.CombinedOutput()
	parsedRewordOutput := processGeneralGitOpsOutputIntoStringArray(rewordOutput)

	if runErr != nil {
		if ctx.Err() != nil {
			return parsedRewordOutput, ctx.Err()
		}
		return parsedRewordOutput, runErr
	}

	return parsedRewordOutput, nil
}

// ------------------------------------
//
//	Interactive Rebase - Reword (Signing variant)
//	Returns the rebase cmd without running it — caller (bubbletea) suspends the TUI,
//	hands the cmd to the OS, and resumes after git finishes.
//
//	* gitti lock is released before the cmd runs — intentional
//	* git's own repo lock (.git/rebase-merge/) prevents concurrent ops during execution
//
// ------------------------------------
func (gIR *GitInteractiveRebase) GitInteractiveRebaseRewordWithSigning(ctx context.Context, gitCommitInfo []CommitInfo, selectedCommitInfo CommitInfo, newCommitMessage string, newCommitDescription string) (*exec.Cmd, func(), error) {
	if !gIR.gitProcessLock.CanProceedWithGitOps() {
		return nil, nil, fmt.Errorf("%s", gIR.gitProcessLock.OtherProcessRunningWarning())
	}
	defer func() {
		gIR.gitProcessLock.ReleaseGitOpsLock()
	}()

	return gIR.interactiveRebaseReword(ctx, gitCommitInfo, selectedCommitInfo, newCommitMessage, newCommitDescription, true)
}

// ------------------------------------
//
//	Builds the rebase command and cleanup callback for reword; validates the selected commit,
//	generates todo and message temp files, and configures a non-interactive sequence editor
//
// ------------------------------------
func (gIR *GitInteractiveRebase) interactiveRebaseReword(ctx context.Context, gitCommitInfo []CommitInfo, selectedCommitInfo CommitInfo, newCommitMessage string, newCommitDescription string, signing bool) (*exec.Cmd, func(), error) {
	// base is merge commit, does not allow that
	if len(selectedCommitInfo.Parent) > 1 {
		return nil, nil, fmt.Errorf("%s", i18n.LANGUAGEMAPPING.InteractiveRebaseRewordCommitCannotBeAMergeCommit)
	}

	// oldest selected commit's CommitOrder equals its index in gitCommitInfo (which is latest-first)
	// so gitCommitInfo[:targetCommitPosition+1] gives us HEAD..oldest-selected
	targetCommitPosition := selectedCommitInfo.CommitOrder
	if targetCommitPosition < 0 || targetCommitPosition >= len(gitCommitInfo) {
		return nil, nil, fmt.Errorf("%s", i18n.LANGUAGEMAPPING.InteractiveRebaseFixupPositionMismatchError)
	}

	// copy so we don't mutate gitCommitInfo, then sort oldest → latest for todo construction
	sortedAffectedCommitInfos := make([]CommitInfo, targetCommitPosition+1)
	copy(sortedAffectedCommitInfos, gitCommitInfo[:targetCommitPosition+1])
	slices.SortFunc(sortedAffectedCommitInfos, func(a, b CommitInfo) int {
		return cmp.Compare(b.CommitOrder, a.CommitOrder) // largest CommitOrder first = oldest to latest
	})

	// build the `exec git commit --amend -F <tempfile>` line that applies the new message
	execAmendCommitCommandString, commitMsgTempPath, execErr := gIR.buildRebaseAmendExec(newCommitMessage, newCommitDescription)
	if execErr != nil {
		return nil, nil, execErr
	}
	rewordTodoString := gIR.constructRewordTodo(sortedAffectedCommitInfos, selectedCommitInfo, execAmendCommitCommandString)

	// write todo to a temp file; the sequence editor script will cp it into git's todo path
	sequenceEditorScriptPath, cleanupFn, buildCmdErr := gIR.buildNonInteractiveTodoCmd(rewordTodoString)
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
//	Construct Reword Todo
//	Builds the rebase todo file content for the reword operation.
//
//	Todo structure (oldest → latest order):
//	  pick <selected>                  ← anchor: the commit whose message is being reworded
//	  exec git commit --amend -F <f>   ← replace message with user's new message (no editor)
//	  pick <newer-1>                   ← replayed on top unchanged
//	  pick <newer-N>
//
//	* merge commits in the affected range are skipped (dropped from rebase)
//
// ------------------------------------
func (gIR *GitInteractiveRebase) constructRewordTodo(sortedAffectedCommitInfos []CommitInfo, selectedCommitInfo CommitInfo, execAmendCommitCommandString string) string {
	var rewordTodoString strings.Builder

	// oldest selected commit is the rebase anchor — `pick` preserves it as the surviving commit
	rewordTodoString.WriteString("pick ")
	rewordTodoString.WriteString(selectedCommitInfo.Hash)
	rewordTodoString.WriteRune('\n')

	// after reword the surviving commit has the oldest selected commit's message;
	// `exec` replaces it with the user-supplied message via --amend -F (no editor opened)
	rewordTodoString.WriteString(execAmendCommitCommandString)
	rewordTodoString.WriteRune('\n')

	if len(sortedAffectedCommitInfos) > 1 {
		for _, commitInfo := range sortedAffectedCommitInfos[1:] {
			// skip any merge commit
			if len(commitInfo.Parent) > 1 {
				continue
			}
			rewordTodoString.WriteString("pick ")
			rewordTodoString.WriteString(commitInfo.Hash)
			rewordTodoString.WriteRune('\n')
		}
	}

	return rewordTodoString.String()
}

func (gIR *GitInteractiveRebase) buildNonInteractiveTodoCmd(todoString string) (string, func(), error) {
	todoFile, todoFileErr := os.CreateTemp("", "gitti-rebase-todo-*")
	if todoFileErr != nil {
		return "", nil, todoFileErr
	}

	if _, writeTodoFileErr := todoFile.WriteString(todoString); writeTodoFileErr != nil {
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

// *************************************************************************************
//                                 GENERAL UTILITIES
// *************************************************************************************

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
	fullMsg := strings.TrimSpace(message)
	if utf8.RuneCountInString(fullMsg) < 1 {
		return "", "", fmt.Errorf("%s", i18n.LANGUAGEMAPPING.CommitMessageMustBeProvided)
	}
	if strings.TrimSpace(description) != "" {
		fullMsg += "\n\n" + strings.TrimSpace(description)
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
