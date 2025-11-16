package git

import (
	"bufio"
	"bytes"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"gitti/executor"
)

func processGeneralGitOpsOutputIntoStringArray(dirtyGitOutput []byte) []string {
	var cleanedStringArray []string
	cleanedStringArray = strings.Split(strings.TrimSpace(string(dirtyGitOutput)), "\n")

	return cleanedStringArray
}

func splitOnCarriageReturnOrNewline(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Find the first delimiter.
	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		// Check which delimiter we found.
		if data[i] == '\r' {
			// It's a carriage return. Return the token *including* the \r.
			return i + 1, data[0 : i+1], nil
		}
		// It's a newline. Return the token *excluding* the \n.
		return i + 1, data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func handleProgressOutputStream(cursorIndex int, scanner *bufio.Scanner, outputArray []string) (int, []string) {
	// line counter was to determine when a line end with \r,
	// we should replace the latest line in the array or append because this is a new line
	line := scanner.Text()
	isReplaceLine := strings.HasSuffix(line, "\r")
	lineContent := strings.TrimRight(line, "\r")

	// if the cursor is larger or equal then the output array, append into the output outputArray
	// if not replace the latest itme in the array
	if cursorIndex >= len(outputArray) {
		outputArray = append(outputArray, lineContent)
	} else {
		// Otherwise, update the line the cursor is pointing to.
		outputArray[cursorIndex] = lineContent
	}

	// if it was not a replace string line, increment the cursorindex
	if !isReplaceLine {
		return cursorIndex + 1, outputArray
	}

	return cursorIndex, outputArray
}

// check if the format for git remote is correct and valid
func isValidGitRemoteURL(remote string) bool {
	// Check HTTPS style
	if strings.HasPrefix(remote, "https://") || strings.HasPrefix(remote, "http://") {
		_, err := url.ParseRequestURI(remote)
		return err == nil
	}

	// Check SSH style (e.g. git@github.com:user/repo.git)
	sshPattern := `^[\w.-]+@[\w.-]+:[\w./-]+(\.git)?$`
	matched, _ := regexp.MatchString(sshPattern, remote)
	return matched
}

// ----------------------------------
//
//	Related to Git Init
//
// ----------------------------------
func GitInit(repoPath string, initBranchName string) {
	initGitArgs := []string{"init"}

	initCmd := executor.GittiCmdExecutor.RunGitCmd(initGitArgs, false)
	_, initErr := initCmd.Output()
	if initErr != nil {
		fmt.Printf("[GIT INIT ERROR]: %v", initErr)
		os.Exit(1)
	}

	// set the branch
	checkoutBranchGitArgs := []string{"checkout", "-b", initBranchName}

	checkoutBranchCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(checkoutBranchGitArgs, false)
	_, checkoutBranchErr := checkoutBranchCmdExecutor.Output()
	if checkoutBranchErr != nil {
		fmt.Printf("[GIT INIT ERROR]: %v", checkoutBranchErr)
		os.Exit(1)
	}
}

// ----------------------------------
//
//	Related to Git check upstream existence
//
// ----------------------------------
func hasUpStream() (string, bool) {
	gitArgs := []string{"rev-parse", "--abbrev-ref", "@{u}"}

	checkUpStreamCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	checkUpStreamOutput, checkUpStreamErr := checkUpStreamCmdExecutor.Output()
	if checkUpStreamErr != nil {
		return "", false
	}

	return strings.TrimSpace(string(checkUpStreamOutput)), true
}
