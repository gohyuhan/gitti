package git

import (
	"regexp"
	"strings"
)

func processGeneralGitOpsOutputIntoStringArray(dirtyGitOutput []byte) []string {
	var cleanedStringArray []string
	cleanedStringArray = strings.Split(strings.TrimSpace(string(dirtyGitOutput)), "\n")

	return cleanedStringArray
}

func cleanGitOutput(s string) string {
	// remove carriage returns
	s = strings.ReplaceAll(s, "\r", "")

	// remove ANSI escape sequences
	re := regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)
	s = re.ReplaceAllString(s, "")

	return s
}
