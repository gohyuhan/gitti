package main

import (
	"fmt"
	"gitti/api"
	"gitti/api/git"
	"gitti/tui"
	"os"
	"time"

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
	debounceMS := 500 * time.Millisecond
	git.InitGitBranch(repoPath)
	git.InitGitFile(repoPath, updateChannel)
	git.GitCommitInit(repoPath, false)
	api.InitGitDaemon(repoPath, debounceMS, updateChannel)

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
		fmt.Printf("Alas, there's been an error: %v", err)
		api.GITDAEMON.Stop()
		os.Exit(1)
	}

}
