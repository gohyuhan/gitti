package api

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/settings"
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

		// reassign again if user choose to init the repo after prompt
		gitPathInfo, err = getGitPathInfo()
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

func InitGitOperations(updateChannel chan string) *GitOperations {
	gitProcessLock := git.InitGitProcessLock()
	return &GitOperations{
		GitBranch:    git.InitGitBranch(gitProcessLock),
		GitCommit:    git.InitGitCommit(updateChannel, gitProcessLock),
		GitFiles:     git.InitGitFile(updateChannel, gitProcessLock),
		GitPull:      git.InitGitPull(updateChannel, gitProcessLock),
		GitStash:     git.InitGitStash(gitProcessLock),
		GitRemote:    git.InitGitRemote(updateChannel, gitProcessLock),
		GitCommitLog: git.InitGitCommitLog(updateChannel, gitProcessLock),
	}
}

func IsBranchNameValid(branchName string) (string, bool) {
	// Git-invalid characters anywhere (except space which we replace with "-")
	// These characters must be removed entirely.
	invalidChars := regexp.MustCompile(`[~^:?*\[\\]`) // characters removed fully
	controlChars := regexp.MustCompile(`[\x00-\x1F\x7F]`)

	modified := strings.TrimSpace(branchName)
	afterModified := ""

	for modified != afterModified {
		afterModified = modified
		modified = strings.ReplaceAll(modified, " ", "-") // space → dash
		modified = invalidChars.ReplaceAllString(modified, "")
		modified = controlChars.ReplaceAllString(modified, "")

		// Remove special disallowed sequences
		modified = strings.ReplaceAll(modified, "..", "")
		modified = strings.ReplaceAll(modified, "/./", "/")
		modified = strings.ReplaceAll(modified, "@{", "")
		modified = strings.ReplaceAll(modified, "//", "/")
	}

	// loop check till the prefix and suffix is clean and valid
	prefixClean := false
	suffixClean := false
	for !prefixClean || !suffixClean {
		// mark the prefix and suffix as clean first
		prefixClean = true
		suffixClean = true

		// prefix
		if strings.HasPrefix(modified, "/") {
			modified = strings.TrimLeft(modified, "/")
			prefixClean = false
		}
		if strings.HasPrefix(modified, ".") {
			modified = strings.TrimLeft(modified, ".")
			prefixClean = false
		}
		if strings.HasPrefix(modified, "refs/") {
			modified = strings.TrimPrefix(modified, "refs/")
			prefixClean = false
		}
		if strings.HasPrefix(modified, "-") {
			modified = strings.TrimLeft(modified, "-")
			prefixClean = false
		}

		// suffix
		if strings.HasSuffix(modified, "/") {
			modified = strings.TrimRight(modified, "/")
			suffixClean = false
		}
		if strings.HasSuffix(modified, ".") {
			modified = strings.TrimRight(modified, ".")
			suffixClean = false
		}
		if strings.HasSuffix(modified, ".lock") {
			modified = strings.TrimSuffix(modified, ".lock")
			suffixClean = false
		}
	}

	if modified == "@" {
		modified = ""
	}

	// Determine if original was already valid
	isValid := (modified == branchName)

	return modified, isValid
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

	repoName := filepath.Base(strings.TrimSpace(string(topLevelGitPathOutput)))

	gitRepoPath := GitRepoPath{
		AbsoluteGitRepoPath: strings.TrimSpace(string(absGitPathOutput)),
		TopLevelRepoPath:    strings.TrimSpace(string(topLevelGitPathOutput)),
		RepoName:            repoName,
	}

	return gitRepoPath, nil
}
