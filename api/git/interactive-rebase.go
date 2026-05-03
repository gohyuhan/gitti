package git

import (
	"cmp"
	"fmt"
	"os"
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
//			Interactive Rebase - fixup (techically the same as squash but we control the editing layer so fixup is more fit in this case)
//		  * selected commit info need to be sorted by the commit order to be oldest -> latest
//	   * to fixup/squash (which result in the same git result), at least 2 commit must be selected
//
// ------------------------------------
func (gIR *GitInteractiveRebase) interactiveRebaseFixup(sortedSelectedCommitInfos []CommitInfo, newCommitMessage string, newCommitDesceription string) error {
	if len(sortedSelectedCommitInfos) < 2 {
		return fmt.Errorf("%s", i18n.LANGUAGEMAPPING.InteractiveRebaseFixupMustHaveAtLeastTwoSelectedError)
	}

	// base is merge commit, does not allow that
	if len(sortedSelectedCommitInfos[0].Parent) > 1 {
		return fmt.Errorf("%s", i18n.LANGUAGEMAPPING.InteractiveRebaseFixupBaseSelectionMustNotBeMergeCommitError)
	}
	targetCommitPosition := sortedSelectedCommitInfos[0].CommitOrder

	// retrieve the affected range of commits and also sort it
	sortedAffectedCommitInfos := gIR.gitCommitInfo[:targetCommitPosition+1]
	slices.SortFunc(sortedAffectedCommitInfos, func(a, b CommitInfo) int {
		return cmp.Compare(b.CommitOrder, a.CommitOrder) // largest first = oldest to latest
	})

	execAmendCommitCommandString, execErr := gIR.buildRebaseAmendExec(newCommitMessage, newCommitDesceription)
	if execErr != nil {
		return execErr
	}
	// fixupTodoString := gIR.constructFixupTodo(sortedAffectedCommitInfos, sortedSelectedCommitInfos, execAmendCommitCommandString)
	gIR.constructFixupTodo(sortedAffectedCommitInfos, sortedSelectedCommitInfos, execAmendCommitCommandString)

	// TODO: construct the platform compatible script, .bash for mac/linux, .bat for window

	return nil
}

// internal fix up todo string construct function
func (gIR *GitInteractiveRebase) constructFixupTodo(sortedAffectedCommitInfos []CommitInfo, sortedSelectedCommitInfos []CommitInfo, execAmendCommitCommandString string) string {
	var fixupTodoString strings.Builder

	fixupTodoString.WriteString("pick ")
	fixupTodoString.WriteString(sortedSelectedCommitInfos[0].Hash)
	fixupTodoString.WriteRune('\n')

	for _, commitInfo := range sortedSelectedCommitInfos[1:] {
		// skip any merge commit
		if len(commitInfo.Parent) > 1 {
			continue
		}
		fixupTodoString.WriteString("fixup ")
		fixupTodoString.WriteString(commitInfo.Hash)
		fixupTodoString.WriteRune('\n')
	}

	fixupTodoString.WriteString(execAmendCommitCommandString)
	fixupTodoString.WriteRune('\n')

	// Create a set of selected commit hashes for O(1) lookup
	selectedSet := make(map[string]struct{}, len(sortedSelectedCommitInfos))
	for _, c := range sortedSelectedCommitInfos {
		selectedSet[c.Hash] = struct{}{}
	}

	// Filter affectedCommitInfos, preserving order
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

func (gIR *GitInteractiveRebase) buildRebaseAmendExec(message, description string) (string, error) {
	// combine message + description
	fullMsg := message
	if strings.TrimSpace(description) != "" {
		fullMsg += "\n\n" + description
	}
	fullMsg += "\n"

	// create temp file
	f, err := os.CreateTemp("", "gitti-commit-msg-*")
	if err != nil {
		return "", err
	}

	// write content
	if _, err := f.WriteString(fullMsg); err != nil {
		f.Close()
		return "", err
	}

	if err := f.Close(); err != nil {
		return "", err
	}

	// shell escape path (safe for rebase exec)
	escapedPath := "'" + strings.ReplaceAll(f.Name(), "'", `'\''`) + "'"

	// construct exec string
	execStr := fmt.Sprintf(
		"exec git commit --amend -F %s",
		escapedPath,
	)

	return execStr, nil
}

// ------------------------------------
//
//	gIR internal clear git commit info func
//
// ------------------------------------
func (gIR *GitInteractiveRebase) clearGitCommitInfos() {
	gIR.gitCommitInfo = []CommitInfo{}
}
