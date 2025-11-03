package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"

	"gitti/api"
	"gitti/i18n"
	"gitti/tui"
)

func main() {
	repoPath, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s: %v", i18n.LANGUAGEMAPPING.FailToGetCWD, err)
		os.Exit(1)
	}

	initGlobalSettingAndLanguage()
	langCode := flag.String("language", "", i18n.LANGUAGEMAPPING.FlagLangCode)
	defaultInitBranch := flag.String("init-dbranch", "", i18n.LANGUAGEMAPPING.FlagInitDefaultBranch)
	applyToSystemGit := flag.Bool("global", false, i18n.LANGUAGEMAPPING.FlagGlobal)

	flag.Parse()

	switch {
	case *langCode != "":
		setLanguage(*langCode)
	case *defaultInitBranch != "" && *applyToSystemGit:
		setGlobalInitBranch(*defaultInitBranch, repoPath)
	case *defaultInitBranch != "" && !*applyToSystemGit:
		setInitBranch(*defaultInitBranch)
	default:
		// create the channel that will be the bring to emit update event back to main thread
		updateChannel := make(chan string)

		// initialization
		initGitAndAPI(repoPath, updateChannel)

		// start the Git Daemon
		api.GITDAEMON.Start()

		gittiUiModel := tui.NewGittiModel(repoPath)
		gitti := tea.NewProgram(
			&gittiUiModel,
		)

		tui.StartGitUpdateListener(gitti, updateChannel)

		if _, err := gitti.Run(); err != nil {
			api.GITDAEMON.Stop()
			fmt.Printf("%s: %v", i18n.LANGUAGEMAPPING.TuiRunFail, err)
			os.Exit(1)
		}
	}
}
