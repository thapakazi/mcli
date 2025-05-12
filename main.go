package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug logging to debug.log")
	flag.Parse()
	p := tea.NewProgram(
		NewModel((*debug)),
		tea.WithInput(os.Stdin),
		tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed with error:%v\n", err)
	}
}
