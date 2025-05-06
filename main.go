package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initModel(), tea.WithInput(os.Stdin))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed with error:%v\n", err)
	}
}
