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
//               <site:     https://yuhangoh.com>
//               <github:   https://github.com/gohyuhan>
//               <linkedin: https://my.linkedin.com/in/yu-han-goh-209480200>

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/config"
	"github.com/gohyuhan/gitti/constant"
	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui"
	"github.com/gohyuhan/gitti/updater"
)

func main() {
	repoPath, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s: %v", i18n.LANGUAGEMAPPING.FailToGetCWD, err)
		os.Exit(1)
	}

	// setting and config need to be the first thing to be initialized
	config.InitGlobalSettingAndLanguage()
	showVersion := flag.Bool("version", false, i18n.LANGUAGEMAPPING.FlagVersion)
	langCode := flag.String("language", "", i18n.LANGUAGEMAPPING.FlagLangCode)
	defaultInitBranch := flag.String("init-dbranch", "", i18n.LANGUAGEMAPPING.FlagInitDefaultBranch)
	autoUpdate := flag.String("auto-update", "", i18n.LANGUAGEMAPPING.FlagAutoUpdate)
	updatePrompt := flag.Bool("update", false, i18n.LANGUAGEMAPPING.FlagUpdate)
	applyToSystemGit := flag.Bool("global", false, i18n.LANGUAGEMAPPING.FlagGlobal)
	setEditor := flag.Bool("editor", false, i18n.LANGUAGEMAPPING.FlagEditor)
	setMaxCommitLogCount := flag.Int("max-commit-log-count", 0, i18n.LANGUAGEMAPPING.FlagMaxCommitLogCount)
	allowCommitGraphWrite := flag.String("allow-commit-graph-write", "", i18n.LANGUAGEMAPPING.FlagAllowCommitGraphWrite)

	flag.Parse()

	// the Cmd Shoule be initialized right after gitti setting
	executor.InitCmdExecutor(repoPath)

	switch {
	case *showVersion:
		fmt.Println(constant.APPVERSION)
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
	case *setEditor:
		config.ChooseAndSetEditor()
	case *setMaxCommitLogCount > 0:
		config.SetMaxCommitLogCount(*setMaxCommitLogCount)
	case *allowCommitGraphWrite != "":
		config.SetAllowCommitGraphWrite(*allowCommitGraphWrite)
	default:
		// create the channel that will be the bring to emit update event back to main thread
		gitUpdateChannel := make(chan string)
		tuiUpdateChannel := make(chan string)
		loggingUpdateChannel := make(chan string)
		gittiLogging := logging.InitGittiLogging(300, loggingUpdateChannel, 3)

		// initialization
		gitOperations, gitRepoPathInfo := config.InitGitAndAPI(repoPath, gitUpdateChannel, gittiLogging)

		// check for update if user allows it
		if settings.GITTICONFIGSETTINGS.AutoUpdate {
			updater.AutoUpdater()
		}

		lastMouseSignal := time.Now()
		mouseThrottleFrequency := 8 * time.Millisecond

		// throttle the mouse signal to ~120 fps
		mouseThrottle := func(m tea.Model, msg tea.Msg) tea.Msg {
			if _, ok := msg.(tea.MouseMsg); ok {
				if time.Since(lastMouseSignal) < mouseThrottleFrequency {
					return nil // Drop the message entirely
				}
				lastMouseSignal = time.Now()
			}
			return msg
		}

		gittiAppModel := tui.NewGittiAppModel(tuiUpdateChannel, repoPath, gitRepoPathInfo.RepoName, gitOperations, gittiLogging)
		gitti := tea.NewProgram(
			gittiAppModel,
			tea.WithFilter(mouseThrottle),
		)

		tui.StartGitUpdateListener(gitti, gitUpdateChannel)
		tui.StartTuiUpdateListener(gitti, tuiUpdateChannel)
		tui.StartLoggingUpdateListener(gitti, loggingUpdateChannel)

		// start the Git Daemon
		if api.GITDAEMON != nil {
			api.GITDAEMON.Start()
		}

		if _, err := gitti.Run(); err != nil {
			if api.GITDAEMON != nil {
				api.GITDAEMON.Stop()
			}
			fmt.Printf("%s: %v", i18n.LANGUAGEMAPPING.TuiRunFail, err)
			os.Exit(1)
		}
	}
}
