package main

//                              ,----,       ,----,
//                            ,/   .`|     ,/   .`|
//    ,----..      ,---,    ,`   .'  :   ,`   .'  :   ,---,
//   /   /   \  ,`--.' |  ;    ;     / ;    ;     /,`--.' |
//  |   :     : |   :  :.'___,/    ,'.'___,/    ,' |   :  :
//  .   |  ;. / :   |  '|    :     | |    :     |  :   |  '
//  .   ; /--`  |   :  |;    |.';  ; ;    |.';  ;  |   :  |
//  ;   | ;  __ '   '  ;`----'  |  | `----'  |  |  '   '  ;
//  |   : |.' .'|   |  |    '   :  ;     '   :  ;  |   |  |
//  .   | '_.' :'   :  ;    |   |  '     |   |  '  '   :  ;
//  '   ; : \  ||   |  '    '   :  |     '   :  |  |   |  '
//  '   | '/  .''   :  |    ;   |.'      ;   |.'   '   :  |
//  |   :    /  ;   |.'     '---'        '---'     ;   |.'
//   \   \ .'   '---'                              '---'
//    `---`

// By Yu Han Goh <software engineer>
//               <site:     https://yh.boredui.com>
//               <github:   https://github.com/gohyuhan>
//               <linkedin: https://my.linkedin.com/in/yu-han-goh-209480200>

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"

	"gitti/api"
	"gitti/config"
	"gitti/executor"
	"gitti/i18n"
	"gitti/settings"
	"gitti/tui"
	"gitti/updater"
)

func main() {
	repoPath, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s: %v", i18n.LANGUAGEMAPPING.FailToGetCWD, err)
		os.Exit(1)
	}

	// setting and config need to be the first thing to be initialized
	config.InitGlobalSettingAndLanguage()
	langCode := flag.String("language", "", i18n.LANGUAGEMAPPING.FlagLangCode)
	defaultInitBranch := flag.String("init-dbranch", "", i18n.LANGUAGEMAPPING.FlagInitDefaultBranch)
	autoUpdate := flag.String("auto-update", "", i18n.LANGUAGEMAPPING.FlagAutoUpdate)
	updatePrompt := flag.Bool("update", false, i18n.LANGUAGEMAPPING.FlagUpdate)
	applyToSystemGit := flag.Bool("global", false, i18n.LANGUAGEMAPPING.FlagGlobal)

	flag.Parse()

	// the Cmd Shoule be initialized right after gitti setting
	executor.InitCmdExecutor(repoPath)

	switch {
	case *langCode != "":
		config.SetLanguage(*langCode)
	case *defaultInitBranch != "" && *applyToSystemGit:
		config.SetGlobalInitBranch(*defaultInitBranch, repoPath)
	case *defaultInitBranch != "" && !*applyToSystemGit:
		config.SetInitBranch(*defaultInitBranch)
	case *autoUpdate != "":
		config.SetAutoUpdate(*autoUpdate)
	case *updatePrompt:
		updater.Update()
	default:
		// create the channel that will be the bring to emit update event back to main thread
		updateChannel := make(chan string)

		// initialization
		gitOperations, gitRepoPathInfo := config.InitGitAndAPI(repoPath, updateChannel)

		// check for update if user allows it
		if settings.GITTICONFIGSETTINGS.AutoUpdate {
			updater.AutoUpdater()
		}

		// start the Git Daemon
		api.GITDAEMON.Start()

		gittiUiModel := tui.NewGittiModel(repoPath, gitRepoPathInfo.RepoName, gitOperations)
		gitti := tea.NewProgram(
			gittiUiModel,
		)

		tui.StartGitUpdateListener(gitti, updateChannel)

		if _, err := gitti.Run(); err != nil {
			api.GITDAEMON.Stop()
			fmt.Printf("%s: %v", i18n.LANGUAGEMAPPING.TuiRunFail, err)
			os.Exit(1)
		}
	}
}
