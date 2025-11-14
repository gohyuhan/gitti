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
	"gitti/cmd"
	"gitti/config"
	"gitti/i18n"
	"gitti/tui"
)

func main() {
	repoPath, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s: %v", i18n.LANGUAGEMAPPING.FailToGetCWD, err)
		os.Exit(1)
	}

	config.InitGlobalSettingAndLanguage()
	langCode := flag.String("language", "", i18n.LANGUAGEMAPPING.FlagLangCode)
	defaultInitBranch := flag.String("init-dbranch", "", i18n.LANGUAGEMAPPING.FlagInitDefaultBranch)
	applyToSystemGit := flag.Bool("global", false, i18n.LANGUAGEMAPPING.FlagGlobal)

	flag.Parse()

	cmd.InitCmd(repoPath)

	switch {
	case *langCode != "":
		config.SetLanguage(*langCode)
	case *defaultInitBranch != "" && *applyToSystemGit:
		config.SetGlobalInitBranch(*defaultInitBranch, repoPath)
	case *defaultInitBranch != "" && !*applyToSystemGit:
		config.SetInitBranch(*defaultInitBranch)
	default:
		// create the channel that will be the bring to emit update event back to main thread
		updateChannel := make(chan string)

		// initialization
		gitState := config.InitGitAndAPI(repoPath, updateChannel)

		// start the Git Daemon
		api.GITDAEMON.Start()

		gittiUiModel := tui.NewGittiModel(repoPath, gitState)
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
