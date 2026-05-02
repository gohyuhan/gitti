package git

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/logging"
)

const (
	MAX_AUTHOR_CHAR_LENGTH = 16
)

type GitBlame struct {
	updateChannel  chan string
	gitProcessLock *GitProcessLock
	logging        *logging.GittiLogging
}

// LineBlameInfo holds per-line data parsed from `git blame --line-porcelain` output.
type LineBlameInfo struct {
	CommitHash            string
	Author                string
	AuthorMail            string
	AuthorTime            string
	AuthorTimeZone        string
	Committer             string
	CommitterMail         string
	CommitterTime         string
	CommitterTimeZone     string
	CommitSummary         string
	FileName              string
	CommittedLine         string
	ConsolidatedBlameInfo string // this is the blame info for display (not including the code line)
}

// ------------------------------------
//
//	Initialize the git blame handler with shared dependencies
//
// ------------------------------------
func InitGitBlame(updateChannel chan string, gitProcessLock *GitProcessLock, logging *logging.GittiLogging) *GitBlame {
	gitBlame := GitBlame{
		gitProcessLock: gitProcessLock,
		updateChannel:  updateChannel,
		logging:        logging,
	}
	return &gitBlame
}

// ------------------------------------
//
//	Return all files currently tracked by git in the repo via `git ls-files`
//
// ------------------------------------
func (gb *GitBlame) GetCurrentGitTrackedFiles() []string {
	gitArgs := []string{"ls-files", "--cached", "--others", "--exclude-standard"}
	getCurrentGitTrackedFilesCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	gb.logging.RegisterNewLog(logging.GET_CURRENT_GIT_TRACKED_FILES, strings.Join(gitArgs, " "), logging.INFO, "", true)
	getCurrentGitTrackedFilesCmdExecutorOutput, getCurrentGitTrackedFilesCmdExecutorErr := getCurrentGitTrackedFilesCmdExecutor.Output()

	if getCurrentGitTrackedFilesCmdExecutorErr != nil {
		gb.logging.RegisterNewLog(logging.GET_CURRENT_GIT_TRACKED_FILES, "", logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.GET_CURRENT_GIT_TRACKED_FILES, getCurrentGitTrackedFilesCmdExecutorErr.Error()), true)
		return []string{}
	}

	parsedOutput := processGeneralGitOpsOutputIntoStringArray(getCurrentGitTrackedFilesCmdExecutorOutput)
	return parsedOutput
}

// ------------------------------------
//
//	Parse `git blame --line-porcelain` output for the given file and return per-line blame info
//
// ------------------------------------
func (gb *GitBlame) GetFileGitBlameInfo(filePath string) (int, int, []LineBlameInfo) {
	var lineBlameInfo []LineBlameInfo
	largestConsolidatedBlameInfoLineLength := 0 // this is the length for the consolidated blame info (commit hash [first 7 chars] + time ago + author)
	largestLineLength := 0                      // this is the length for the consolidated blame info + code line, which determines the viewport largest width

	gitArgs := []string{"blame", "--line-porcelain", filePath}
	getFileGitBlameInfoCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	gb.logging.RegisterNewLog(logging.GET_FILE_BLAME_INFO, strings.Join(gitArgs, " "), logging.INFO, "", true)
	getFileGitBlameInfoCmdExecutorOutput, getFileGitBlameInfoCmdExecutorErr := getFileGitBlameInfoCmdExecutor.Output()

	if getFileGitBlameInfoCmdExecutorErr != nil {
		gb.logging.RegisterNewLog(logging.GET_FILE_BLAME_INFO, "", logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.GET_FILE_BLAME_INFO, getFileGitBlameInfoCmdExecutorErr.Error()), true)
		return largestConsolidatedBlameInfoLineLength, largestLineLength, lineBlameInfo
	}

	parsedOutput := processGeneralGitOpsOutputIntoStringArray(getFileGitBlameInfoCmdExecutorOutput)

	var currentBlameInfo LineBlameInfo
	for _, parsedLine := range parsedOutput {
		switch {
		case strings.HasPrefix(parsedLine, "boundary"):
			continue
		case strings.HasPrefix(parsedLine, "author "):
			currentBlameInfo.Author = sanitizeBlameLine(strings.SplitN(parsedLine, " ", 2)[1])
		case strings.HasPrefix(parsedLine, "author-mail "):
			currentBlameInfo.AuthorMail = sanitizeBlameLine(strings.SplitN(parsedLine, " ", 2)[1])
		case strings.HasPrefix(parsedLine, "author-time "):
			currentBlameInfo.AuthorTime = sanitizeBlameLine(strings.SplitN(parsedLine, " ", 2)[1])
		case strings.HasPrefix(parsedLine, "author-tz "):
			currentBlameInfo.AuthorTimeZone = sanitizeBlameLine(strings.SplitN(parsedLine, " ", 2)[1])
		case strings.HasPrefix(parsedLine, "committer "):
			currentBlameInfo.Committer = sanitizeBlameLine(strings.SplitN(parsedLine, " ", 2)[1])
		case strings.HasPrefix(parsedLine, "committer-mail "):
			currentBlameInfo.CommitterMail = sanitizeBlameLine(strings.SplitN(parsedLine, " ", 2)[1])
		case strings.HasPrefix(parsedLine, "committer-time "):
			currentBlameInfo.CommitterTime = sanitizeBlameLine(strings.SplitN(parsedLine, " ", 2)[1])
		case strings.HasPrefix(parsedLine, "committer-tz "):
			currentBlameInfo.CommitterTimeZone = sanitizeBlameLine(strings.SplitN(parsedLine, " ", 2)[1])
		case strings.HasPrefix(parsedLine, "summary "):
			currentBlameInfo.CommitSummary = sanitizeBlameLine(strings.SplitN(parsedLine, " ", 2)[1])
		case strings.HasPrefix(parsedLine, "previous "):
			continue
		case strings.HasPrefix(parsedLine, "filename "):
			currentBlameInfo.FileName = sanitizeBlameLine(strings.SplitN(parsedLine, " ", 2)[1])
		case strings.HasPrefix(parsedLine, "\t"):
			shortHash := "!@#$%^&"
			if utf8.RuneCountInString(currentBlameInfo.CommitHash) >= 7 {
				shortHash = currentBlameInfo.CommitHash[:7]
			}
			currentBlameInfo.CommittedLine = normalizeBlameCodeLine(parsedLine[1:]) // strip the \t
			consolidateBlameInfo := sanitizeBlameLine(shortHash + " " +
				timeAgo(currentBlameInfo.AuthorTime) + " " +
				ansi.Truncate(currentBlameInfo.Author, MAX_AUTHOR_CHAR_LENGTH, "...")) // build the consolidate blame

			currentConsolidateBlameInfoLen := ansi.StringWidth(consolidateBlameInfo)
			currentLineLength := ansi.StringWidth(currentBlameInfo.CommittedLine) + currentConsolidateBlameInfoLen
			if currentConsolidateBlameInfoLen > largestConsolidatedBlameInfoLineLength {
				largestConsolidatedBlameInfoLineLength = currentConsolidateBlameInfoLen
			}
			if currentLineLength > largestLineLength {
				largestLineLength = currentLineLength
			}
			currentBlameInfo.ConsolidatedBlameInfo = sanitizeBlameLine(consolidateBlameInfo)
			lineBlameInfo = append(lineBlameInfo, currentBlameInfo)
			currentBlameInfo = LineBlameInfo{} // reinit after append
		default:
			// first line of each porcelain block: "<hash> <orig-lineno> <final-lineno>" — no keyword prefix
			currentBlameInfo.CommitHash = sanitizeBlameLine(strings.SplitN(parsedLine, " ", 2)[0])
		}
	}

	return largestConsolidatedBlameInfoLineLength, largestLineLength, lineBlameInfo
}

// ------------------------------------
//
//	Sanitize blame line
//
// ------------------------------------
func sanitizeBlameLine(s string) string {
	return ansi.Strip(s)
}

// ------------------------------------
//
//	Normalize blame code line
//
// ------------------------------------
func normalizeBlameCodeLine(s string) string {
	return strings.ReplaceAll(sanitizeBlameLine(s), "\t", "    ")
}
