package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	Events []Event
	Err    error
}

func initModel() model {
	return model{
		Events: []Event{},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			log.Println("enter pressed")
			return m, nil
		case "j", "k":
			log.Println("navigation...")
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	return "Loading..."
}
