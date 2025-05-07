package main

import (
	"fmt"
	"log"
	"mcli/types"
	"mcli/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	Events []types.Event
	Err    error
}

func initModel() model {
	return model{
		Events: []types.Event{},
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		events, err := utils.FetchEvents()
		return types.EventsMsg{
			Events: events,
			Err:    err,
		}

	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case types.EventsMsg:
		m.Events = msg.Events
		m.Err = msg.Err
		return m, nil
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
	if m.Err != nil {
		return fmt.Sprintf("Error fetching events. %s\n", m.Err.Error())
	}
	if len(m.Events) == 0 {
		return "No events found\n"
	}
	output := "Fetched events:\n"
	for _, event := range m.Events {
		output += fmt.Sprintf("%s: %+v\n", event.DateTime, event.Title)
	}
	return output
}
