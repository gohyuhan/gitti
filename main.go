package main

import (
	"fmt"
	"gitti/ui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	gitti := tea.NewProgram(
		ui.NewGittiModel(),
		tea.WithAltScreen(), // ← enables full-screen TUI mode
		tea.WithMouseCellMotion(),
	)
	if _, err := gitti.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
