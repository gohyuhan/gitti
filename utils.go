package main

import (
	"fmt"
	"os"

	"gitti/api"
	"gitti/api/git"
	"gitti/i18n"
	"gitti/settings"
)

func setLanguage(langCode string) {
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
func setInitBranch(branchName string) {
	settings.UpdatedDefaultBranch(branchName, false, "")
	fmt.Printf(i18n.LANGUAGEMAPPING.GittiDefaultBranchSet+"\n", branchName)
	os.Exit(0)
}

// set the default git init branch name for both gitti and git
func setGlobalInitBranch(branchName string, cwd string) {
	settings.UpdatedDefaultBranch(branchName, true, cwd)
	fmt.Printf(i18n.LANGUAGEMAPPING.GittiDefaultAndGitDefaultBranchSet+"\n", branchName)
	os.Exit(0)
}

func initGitAndAPI(repoPath string, updateChannel chan string) {
	// check if git is installed in system if not, exit(1)
	api.IsGitInstalled(repoPath)
	// check if the user repo is git inited, is not prompt user to init it
	api.IsRepoGitInitialized(repoPath)
	// various initialization
	git.InitGitBranch(repoPath)
	git.InitGitFile(repoPath, updateChannel)
	git.InitGitCommit(repoPath, updateChannel)
	// git.InitGitCommitLog(repoPath, false) // not included in v0.1.0
	api.InitGitDaemon(repoPath, updateChannel)
}

func initGlobalSettingAndLanguage() {
	settings.InitOrReadConfig()
	i18n.InitGittiLanguageMapping(settings.GITTICONFIGSETTINGS.LanguageCode)
}
