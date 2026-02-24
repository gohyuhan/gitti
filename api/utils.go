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
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/settings"
)

// ----------------------------------
//
//	check if git was installed
//
// ----------------------------------
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

// ----------------------------------
//
//	IsRepoGitInitialized checks if the given path is a Git repository
//
// ----------------------------------
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

// ----------------------------------
//
//	to prompt user if they want to git init current dir if .git is not detected
//
// ----------------------------------
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

// ----------------------------------
//
//	Initialize all git operation handlers with the given path and shared dependencies
//
// ----------------------------------
func InitGitOperations(absolutePath string, updateChannel chan string, gittiLogging *logging.GittiLogging) *GitOperations {
	gitProcessLock := git.InitGitProcessLock(gittiLogging)
	return &GitOperations{
		GitBranch:              git.InitGitBranch(gitProcessLock, gittiLogging),
		GitCommit:              git.InitGitCommit(updateChannel, gitProcessLock, gittiLogging),
		GitFiles:               git.InitGitFile(updateChannel, gitProcessLock, gittiLogging),
		GitPull:                git.InitGitPull(updateChannel, gitProcessLock, gittiLogging),
		GitStash:               git.InitGitStash(gitProcessLock, gittiLogging),
		GitRemote:              git.InitGitRemote(updateChannel, gitProcessLock, gittiLogging),
		GitCommitLog:           git.InitGitCommitLog(updateChannel, gitProcessLock, settings.GITTICONFIGSETTINGS.MaxCommitLogCount, gittiLogging),
		GitRefLog:              git.InitGitRefLog(updateChannel, gitProcessLock, settings.GITTICONFIGSETTINGS.MaxRefLogCount, gittiLogging),
		GitTag:                 git.InitGitTag(updateChannel, gitProcessLock, gittiLogging),
		GitStateUniversalUtils: git.InitGitStateUniversalUtils(absolutePath, gitProcessLock, gittiLogging),
	}
}

// ----------------------------------
//
//	to check and validate if the input branch name is valid and also to generate the valid branch name
//
// ----------------------------------
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

// ----------------------------------
//
//	get the top-level and absolute git path
//	top-level eg: /Users/<username>/Projects/<example-project>
//	absolute-path eg: /Users/<username>/Projects/<example-project>/.git
//
// ----------------------------------
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

// ----------------------------------
//
//	check git operation if require signing
//
// ----------------------------------
func CheckSigningRequiredOperation() (bool, bool, bool) {
	var commitRequireSigning bool
	var tagRequireSigning bool
	var pushRequireSigning bool

	commitSigningGitArgs := []string{"config", "--type=bool", "commit.gpgsign"}
	tagSigningGitArgs := []string{"config", "--type=bool", "tag.gpgsign"}
	pushSigningGitArgs := []string{"config", "get", "push.gpgsign"}

	commitSigningCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(commitSigningGitArgs, false)
	tagSigningCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(tagSigningGitArgs, false)
	pushSigningCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(pushSigningGitArgs, false)

	commitRequireSigningOutput, commitRequireSigningErr := commitSigningCmdExecutor.Output()
	tagRequireSigningOutput, tagRequireSigningErr := tagSigningCmdExecutor.Output()
	pushRequireSigningOutput, pushRequireSigningErr := pushSigningCmdExecutor.Output()

	if commitRequireSigningErr != nil {
		commitRequireSigning = false
	} else {
		commitRequireSigning = strings.ToLower(strings.TrimSpace(string(commitRequireSigningOutput))) == "true"
	}

	if tagRequireSigningErr != nil {
		tagRequireSigning = false
	} else {
		tagRequireSigning = strings.ToLower(strings.TrimSpace(string(tagRequireSigningOutput))) == "true"
	}

	if pushRequireSigningErr != nil {
		pushRequireSigning = false
	} else {
		output := strings.ToLower(strings.TrimSpace(string(pushRequireSigningOutput)))
		pushRequireSigning = output == "true" || output == "if-asked"
	}

	return commitRequireSigning, tagRequireSigning, pushRequireSigning
}
