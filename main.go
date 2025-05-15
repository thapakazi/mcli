package main

import (
	"fmt"
	"os"

	"mcli/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	// Initialize the global logger
	utils.InitLogger()

	// Log a startup message
	utils.Logger.Info("Program started")

	p := tea.NewProgram(
		NewModel(),
		tea.WithInput(os.Stdin),
		tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed with error:%v\n", err)
	}
}
