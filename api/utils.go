package api

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gitti/api/git"
	"gitti/executor"
	"gitti/i18n"
	"gitti/settings"
)

func IsGitInstalled(repoPath string) {
	gitArgs := []string{"--version"}

	cmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	cmdExecutor.Dir = repoPath
	err := cmdExecutor.Run()
	if err != nil {
		_, notInSystem := err.(*exec.Error) // check if git is not installed wihitn the system, exec Error means it the executable was no within the system
		if notInSystem {
			fmt.Println(i18n.LANGUAGEMAPPING.GitNotInstalledError)
			os.Exit(1)
		}
	}
}

// IsRepoGitInitialized checks if the given path is a Git repository
func IsRepoGitInitialized(repoPath string) GitRepoPath {
	gitPathInfo, err := getGitPathInfo()
	if err != nil {
		// .git does not exist or some other error
		PromptUserForGitInitConfirmation(repoPath)
	}

	return gitPathInfo
}

func PromptUserForGitInitConfirmation(repoPath string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(i18n.LANGUAGEMAPPING.GitNotInitPrompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToUpper(input))

	switch input {
	case "Y":
		git.GitInit(repoPath, settings.GITTICONFIGSETTINGS.GitInitDefaultBranch)
	case "N":
		fmt.Println(i18n.LANGUAGEMAPPING.GitInitRefuse)
		os.Exit(0)
	default:
		fmt.Println(i18n.LANGUAGEMAPPING.GitInitPromptInvalidInput)
		os.Exit(1)
	}
}

func InitGitState(updateChannel chan string) *GitState {
	gitProcessLock := git.InitGitProcessLock()
	return &GitState{
		GitBranch: git.InitGitBranch(gitProcessLock),
		GitCommit: git.InitGitCommit(updateChannel, gitProcessLock),
		GitFiles:  git.InitGitFile(gitProcessLock),
		GitPull:   git.InitGitPull(updateChannel, gitProcessLock),
		GitStash:  git.InitGitStash(gitProcessLock),
	}
}

func IsBranchNameValid(branchName string) (string, bool) {
	var modifiedBranchName string

	modifiedBranchName = branchName
	modifiedBranchName = strings.TrimSpace(strings.ReplaceAll(modifiedBranchName, " ", "-"))

	if modifiedBranchName != branchName {
		return modifiedBranchName, false
	}

	return branchName, true
}

func getGitPathInfo() (GitRepoPath, error) {
	// get the most absolute git folder path
	absGitPathArgs := []string{"rev-parse", "--absolute-git-dir"}
	absGitPathCmd := executor.GittiCmdExecutor.RunGitCmd(absGitPathArgs, false)
	absGitPathOutput, absGitPathErr := absGitPathCmd.Output()

	if absGitPathErr != nil {
		return GitRepoPath{}, fmt.Errorf("not git initialized")
	}

	// get the top level git path
	topLevelGitPathArgs := []string{"rev-parse", "--show-toplevel"}
	topLevelGitPathCmd := executor.GittiCmdExecutor.RunGitCmd(topLevelGitPathArgs, false)
	topLevelGitPathOutput, topLevelGitPathErr := topLevelGitPathCmd.Output()
	if topLevelGitPathErr != nil {
		return GitRepoPath{}, fmt.Errorf("not git initialized")
	}

	gitRepoPath := GitRepoPath{
		AbsoluteGitRepoPath: strings.TrimSpace(string(absGitPathOutput)),
		TopLevelRepoPath:    strings.TrimSpace(string(topLevelGitPathOutput)),
	}

	return gitRepoPath, nil
}
