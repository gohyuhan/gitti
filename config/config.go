package config

import (
	"fmt"
	"os"

	"gitti/api"
	"gitti/executor"
	"gitti/i18n"
	"gitti/settings"
)

func SetLanguage(langCode string) {
	if i18n.IsLanguageCodeSupported(langCode) {
		settings.UpdateLanguageCode(langCode)
		fmt.Printf(i18n.LANGUAGEMAPPING.LanguageSet+"\n", langCode)
		os.Exit(0)
	} else {
		fmt.Printf(i18n.LANGUAGEMAPPING.LanguageNotSupportedPanic+"\n", langCode, i18n.SUPPORTED_LANGUAGE_CODE)
		os.Exit(1)
	}
}

// set the default git init branch name only for gitti
func SetInitBranch(branchName string) {
	settings.UpdateDefaultBranch(branchName, false, "")
	fmt.Printf(i18n.LANGUAGEMAPPING.GittiDefaultBranchSet+"\n", branchName)
	os.Exit(0)
}

// set the default git init branch name for both gitti and git
func SetGlobalInitBranch(branchName string, cwd string) {
	settings.UpdateDefaultBranch(branchName, true, cwd)
	fmt.Printf(i18n.LANGUAGEMAPPING.GittiDefaultAndGitDefaultBranchSet+"\n", branchName)
	os.Exit(0)
}

func InitGitAndAPI(repoPath string, updateChannel chan string) *api.GitState {
	// check if git is installed in system if not, exit(1)
	api.IsGitInstalled(repoPath)
	// check if the user repo is git inited, is not prompt user to init it
	gitRepoPathInfo := api.IsRepoGitInitialized(repoPath)

	// after we successfully get the gitRepoPathInfo back we need to update the current cmd executor dir
	executor.GittiCmdExecutor.UpdateRepoPath(gitRepoPathInfo.TopLevelRepoPath)
	// various initialization
	gitState := api.InitGitState(updateChannel)
	// git.InitGitCommitLog(false) // not included in v0.1.0
	api.InitGitDaemon(gitRepoPathInfo.AbsoluteGitRepoPath, updateChannel, gitState)

	return gitState

}

func InitGlobalSettingAndLanguage() {
	settings.InitOrReadConfig()
	i18n.InitGittiLanguageMapping(settings.GITTICONFIGSETTINGS.LanguageCode)
}
