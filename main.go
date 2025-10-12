package main

import (
	"fmt"
	"gitti/api"
	"gitti/tui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	repoPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Failed to get current working directory: %v", err))
	}
	gitWorkerDaemon := api.NewGitWorkerDaemon(repoPath)

	gittiUiModel := tui.NewGittiModel(repoPath, gitWorkerDaemon)
	gitti := tea.NewProgram(
		gittiUiModel,
		tea.WithAltScreen(), // ← enables full-screen TUI mode
		tea.WithMouseCellMotion(),
	)

	tui.StartGitUpdateListener(gitti, gitWorkerDaemon)

	if _, err := gitti.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		gitWorkerDaemon.Stop()
		os.Exit(1)
	}

}
