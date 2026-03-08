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
	"charm.land/lipgloss/v2"

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
	setMaxRefLogCount := flag.Int("max-reflog-count", 0, i18n.LANGUAGEMAPPING.FlagMaxRefLogCount)
	allowCommitGraphWrite := flag.String("allow-commit-graph-write", "", i18n.LANGUAGEMAPPING.FlagAllowCommitGraphWrite)
	setMaxLogCount := flag.Int("max-log-count", 0, i18n.LANGUAGEMAPPING.FlagMaxLogCount)
	setShowXLog := flag.Int("show-x-log", 0, i18n.LANGUAGEMAPPING.FlagShowXLog)
	overrideSigningUISuspend := flag.String("override-signing-ui-suspend", "", i18n.LANGUAGEMAPPING.FlagOverrideSigningUISuspend)

	flag.Parse()

	// the Cmd Shoule be initialized right after gitti setting
	executor.InitCmdExecutor(repoPath)

	switch {
	case *showVersion:
		lipgloss.Println(constant.APPVERSION)
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
	case *setMaxRefLogCount > 0:
		config.SetMaxRefLogCount(*setMaxRefLogCount)
	case *allowCommitGraphWrite != "":
		config.SetAllowCommitGraphWrite(*allowCommitGraphWrite)
	case *setMaxLogCount > 0:
		config.SetMaxLogCount(*setMaxLogCount)
	case *setShowXLog > 0:
		config.SetShowXLog(*setShowXLog)
	case *overrideSigningUISuspend != "":
		config.SetOverrideSigningUISuspend(*overrideSigningUISuspend)
	default:
		// create the channel that will be the bring to emit update event back to main thread
		gitUpdateChannel := make(chan string, 32)
		tuiUpdateChannel := make(chan string, 32)
		loggingUpdateChannel := make(chan string, 64)
		daemonUpdateChannel := make(chan string, 16)
		gittiLogging := logging.InitGittiLogging(settings.GITTICONFIGSETTINGS.MaxLogCount, loggingUpdateChannel, settings.GITTICONFIGSETTINGS.ShowXLog)

		// initialization
		gitOperations, gitRepoPathInfo := config.InitGitAndAPI(repoPath, gitUpdateChannel, gittiLogging, daemonUpdateChannel)

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

		gittiAppModel := tui.NewGittiAppModel(tuiUpdateChannel, gitRepoPathInfo.TopLevelRepoPath, gitRepoPathInfo.RepoName, gitOperations, gittiLogging, daemonUpdateChannel)
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
