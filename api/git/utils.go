package git

import "strings"

func processGeneralGitOpsOutputIntoStringArray(dirtyGitOutput []byte) []string {
	var cleanedStringArray []string
	cleanedStringArray = strings.Split(strings.TrimSpace(string(dirtyGitOutput)), "\n")

	return cleanedStringArray
}
