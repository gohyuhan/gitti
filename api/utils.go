package api

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gitti/api/git"
	"gitti/i18n"
)

func GetUpdatedGitInfo(updateChannel chan string) {
	git.GITFILES.GetGitFilesStatus()
	git.GITBRANCH.GetLatestBranchesinfo()

	// not included in v0.1.0
	// go func() {
	// 	GITCOMMIT.GetLatestGitCommitLogInfoAndDAG(updateChannel)
	// }()

	updateChannel <- git.GENERAL_GIT_UPDATE
}

func IsGitInstalled(repoPath string) {
	gitArgs := []string{"--version"}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = repoPath
	err := cmd.Run()
	if err != nil {
		_, notInSystem := err.(*exec.Error) // check if git is not installed wihitn the system, exec Error means it the executable was no within the system
		if notInSystem {
			fmt.Println(i18n.LANGUAGEMAPPING.GitNotInstalledError)
			os.Exit(1)
		}
	}
}

// IsRepoGitInitialized checks if the given path is a Git repository
func IsRepoGitInitialized(repoPath string) {
	gitPath := filepath.Join(repoPath, ".git")

	info, err := os.Stat(gitPath)
	if err != nil || !info.IsDir() {
		// .git does not exist or some other error
		PromptUserForGitInitConfirmation(repoPath)
	}
}

func PromptUserForGitInitConfirmation(repoPath string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(i18n.LANGUAGEMAPPING.GitNotInitPrompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToUpper(input))

	switch input {
	case "Y":
		git.GitInit(repoPath)
	case "N":
		fmt.Println(i18n.LANGUAGEMAPPING.GitInitRefuse)
		os.Exit(0)
	default:
		fmt.Println(i18n.LANGUAGEMAPPING.GitInitPromptInvalidInput)
		os.Exit(1)
	}
}
