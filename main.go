package main

import (
	"fmt"
	"gitti/api"
	"gitti/api/git"
	"gitti/i18n"
	"gitti/settings"
	"gitti/tui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	repoPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Failed to get current working directory: %v", err))
	}

	// create the channel that will be the bring to emit update event back to main thread
	updateChannel := make(chan string)

	// various initialization
	settings.InitOrReadConfig("gitti")
	i18n.InitGittiLanguageMapping(settings.GITTICONFIGSETTINGS.LanguageCode)
	git.InitGitBranch(repoPath)
	git.InitGitFile(repoPath, updateChannel)
	// git.GitCommitInit(repoPath, false) // not included in v0.1.0
	api.InitGitDaemon(repoPath, updateChannel)

	// start the Git Daemon
	api.GITDAEMON.Start()

	gittiUiModel := tui.NewGittiModel(repoPath)
	gitti := tea.NewProgram(
		&gittiUiModel,
		tea.WithAltScreen(), // ← enables full-screen TUI mode
		tea.WithMouseCellMotion(),
	)

	tui.StartGitUpdateListener(gitti, updateChannel)

	if _, err := gitti.Run(); err != nil {
		api.GITDAEMON.Stop()
		panic(fmt.Sprintf("Alas, there's been an error: %v", err))
	}
}
