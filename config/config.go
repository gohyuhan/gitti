package config

import (
	"fmt"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/settings"
)

// ----------------------------------
//
//	Set the application language code for i18n and persist to config
//
// ----------------------------------
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

// ----------------------------------
//
//	set the default git init branch name only for gitti
//
// ----------------------------------
func SetInitBranch(branchName string) {
	settings.UpdateDefaultBranch(branchName, false, "")
	fmt.Printf(i18n.LANGUAGEMAPPING.GittiDefaultBranchSet+"\n", branchName)
	os.Exit(0)
}

// ----------------------------------
//
//	set the default git init branch name only for gitti
//
// ----------------------------------
func SetAutoUpdate(autoUpdateString string) {
	if strings.ToLower(autoUpdateString) == "true" {
		settings.UpdateAutoUpdate(true)
		lipgloss.Println(i18n.LANGUAGEMAPPING.UpdaterAutoUpdaterEnable)
		os.Exit(0)
	} else if strings.ToLower(autoUpdateString) == "false" {
		settings.UpdateAutoUpdate(false)
		lipgloss.Println(i18n.LANGUAGEMAPPING.UpdaterAutoUpdaterDisable)
		os.Exit(0)
	} else {
		lipgloss.Println(i18n.LANGUAGEMAPPING.UpdaterAutoUpdaterSetError)
		os.Exit(1)
	}
}

// ----------------------------------
//
//	set the default git init branch name for both gitti and git
//
// ----------------------------------
func SetGlobalInitBranch(branchName string, cwd string) {
	settings.UpdateDefaultBranch(branchName, true, cwd)
	fmt.Printf(i18n.LANGUAGEMAPPING.GittiDefaultAndGitDefaultBranchSet+"\n", branchName)
	os.Exit(0)
}

// ----------------------------------
//
//	set the max commit log count to be retrieve and display
//
// ----------------------------------
func SetMaxCommitLogCount(maxCommitLogCount int) {
	if maxCommitLogCount < 10 {
		lipgloss.Println(i18n.LANGUAGEMAPPING.MaxCommitLogCountSetError)
		os.Exit(1)
	}
	settings.UpdateMaxCommitLogCount(maxCommitLogCount)
	fmt.Printf(i18n.LANGUAGEMAPPING.MaxCommitLogCountSet+"\n", maxCommitLogCount)
	os.Exit(0)
}

// ----------------------------------
//
//	set the max reflog count to be retrieve and display
//
// ----------------------------------
func SetMaxRefLogCount(maxRefLogCount int) {
	if maxRefLogCount < 10 {
		lipgloss.Println(i18n.LANGUAGEMAPPING.MaxRefLogCountSetError)
		os.Exit(1)
	}
	settings.UpdateMaxRefLogCount(maxRefLogCount)
	fmt.Printf(i18n.LANGUAGEMAPPING.MaxRefLogCountSet+"\n", maxRefLogCount)
	os.Exit(0)
}

// ----------------------------------
//
//	Set the override for signing UI suspend behavior and persist to config
//
// ----------------------------------
func SetOverrideSigningUISuspend(override string) {
	if strings.ToLower(override) == "true" {
		settings.UpdateOverrideSigningUISuspend(true)
		lipgloss.Println(i18n.LANGUAGEMAPPING.OverrideSigningUISuspendEnabled)
		os.Exit(0)
	} else if strings.ToLower(override) == "false" {
		settings.UpdateOverrideSigningUISuspend(false)
		lipgloss.Println(i18n.LANGUAGEMAPPING.OverrideSigningUISuspendDisabled)
		os.Exit(0)
	} else {
		lipgloss.Println(i18n.LANGUAGEMAPPING.OverrideSigningUISuspendSetError)
		os.Exit(1)
	}
}

// ----------------------------------
//
//	Set the ffMerge behaviour (ff or non ff for merge ops) and persist to config
//
// ----------------------------------
func SetFfMerge(ffMerge string) {
	if strings.ToLower(ffMerge) == "true" {
		settings.UpdateFfMerge(true)
		lipgloss.Println(i18n.LANGUAGEMAPPING.FfMergeEnabled)
		os.Exit(0)
	} else if strings.ToLower(ffMerge) == "false" {
		settings.UpdateFfMerge(false)
		lipgloss.Println(i18n.LANGUAGEMAPPING.FfMergeDisabled)
		os.Exit(0)
	} else {
		lipgloss.Println(i18n.LANGUAGEMAPPING.FfMergeSetError)
		os.Exit(1)
	}
}

// ----------------------------------
//
//	Enable or disable commit graph write and persist to config
//
// ----------------------------------
func SetAllowCommitGraphWrite(allow string) {
	if strings.ToLower(allow) == "true" {
		settings.UpdateAllowCommitGraphWrite(true)
		lipgloss.Println(i18n.LANGUAGEMAPPING.AllowCommitGraphWriteEnabled)
		os.Exit(0)
	} else if strings.ToLower(allow) == "false" {
		settings.UpdateAllowCommitGraphWrite(false)
		lipgloss.Println(i18n.LANGUAGEMAPPING.AllowCommitGraphWriteDisabled)
		os.Exit(0)
	} else {
		lipgloss.Println(i18n.LANGUAGEMAPPING.AllowCommitGraphWriteSetError)
		os.Exit(1)
	}
}

// ----------------------------------
//
//	Set the maximum number of log entries to retain and persist to config
//
// ----------------------------------
func SetMaxLogCount(maxLog int) {
	if maxLog < settings.GITTICONFIGSETTINGS.ShowXLog {
		fmt.Printf(i18n.LANGUAGEMAPPING.MaxLogCountSetError+"\n", maxLog, settings.GITTICONFIGSETTINGS.ShowXLog)
		os.Exit(1)
	} else {
		settings.UpdateMaxLogCount(maxLog)
		fmt.Printf(i18n.LANGUAGEMAPPING.MaxLogCountSet+"\n", maxLog)
		os.Exit(0)
	}
}

// ----------------------------------
//
//	Set the number of latest log entries to display in the log component and persist to config
//
// ----------------------------------
func SetShowXLog(x int) {
	if x > settings.GITTICONFIGSETTINGS.MaxLogCount {
		fmt.Printf(i18n.LANGUAGEMAPPING.ShowXLogSetError+"\n", x, settings.GITTICONFIGSETTINGS.MaxLogCount)
		os.Exit(1)
	} else {
		settings.UpdateShowXLog(x)
		fmt.Printf(i18n.LANGUAGEMAPPING.ShowXLogSet+"\n", x)
		os.Exit(0)
	}
}

// ----------------------------------
//
//	Initialize git verification, repo path resolution, and all API operations
//
// ----------------------------------
func InitGitAndAPI(repoPath string, updateChannel chan string, gittiLogging *logging.GittiLogging, daemonUpdateChannel chan string) (*api.GitOperations, api.GitRepoPath) {
	// check if git is installed in system if not, exit(1)
	api.IsGitInstalled(repoPath)
	// check if the user repo is git inited, is not prompt user to init it
	gitRepoPathInfo := api.IsRepoGitInitialized(repoPath)
	// after we successfully get the gitRepoPathInfo back we need to update the current cmd executor dir
	executor.GittiCmdExecutor.UpdateRepoPath(gitRepoPathInfo.TopLevelRepoPath)
	// various initialization
	gitOperations := api.InitGitOperations(gitRepoPathInfo.AbsoluteGitRepoPath, updateChannel, gittiLogging)
	api.InitGitDaemon(gitRepoPathInfo.AbsoluteGitRepoPath, updateChannel, gitOperations, settings.GITTICONFIGSETTINGS.AllowCommitGraphWrite, daemonUpdateChannel, gittiLogging)

	return gitOperations, gitRepoPathInfo
}

// ----------------------------------
//
//	Load settings from config file and initialize language mapping
//
// ----------------------------------
func InitGlobalSettingAndLanguage() {
	settings.InitOrReadConfig()
	i18n.InitGittiLanguageMapping(settings.GITTICONFIGSETTINGS.LanguageCode)
}

// ----------------------------------
//
//	Interactive CLI menu to select and persist the default editor
//
// ----------------------------------
func ChooseAndSetEditor() {
	editors := []string{
		"Vim",
		"Neovim",
		"Nano",
		"VSCode",
		"Zed",
		"Cursor",
		"Windsurf",
		"Antigravity",
	}

	// 1. Print the menu
	lipgloss.Println(i18n.LANGUAGEMAPPING.EditorTitle)
	lipgloss.Println(i18n.LANGUAGEMAPPING.EditorDescription)
	for i, name := range editors {
		fmt.Printf("[%d] %s\n", i+1, name)
	}
	fmt.Print(i18n.LANGUAGEMAPPING.EditorInstruction)

	// 2. Read input
	var selection int
	_, err := fmt.Scan(&selection)

	// 3. Validate input
	if err != nil || selection < 1 || selection > len(editors) {
		// clear the invalid input from buffer if needed, or just loop
		lipgloss.Println(i18n.LANGUAGEMAPPING.EditorSetError)
		os.Exit(1)
	}

	// 4. Map to internal value
	choice := editors[selection-1]
	fmt.Printf("Selected: %s\n", choice)

	// Save the result to your config
	settings.UpdateEditor(strings.ToLower(choice))
	fmt.Printf(i18n.LANGUAGEMAPPING.EditorSetSuccess, choice)
	os.Exit(0)
}
