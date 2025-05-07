package main

import (
	"fmt"
	"log"
	"mcli/types"
	"mcli/utils"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	Events       []types.Event
	table        table.Model
	windowHeight int
	Err          error
}

func initModel() model {
	return model{
		Events:       []types.Event{},
		table:        utils.CreateTable([]types.Event{}, 10),
		windowHeight: 10,
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
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.table = utils.CreateTable(m.Events, m.windowHeight)
		return m, nil
	case types.EventsMsg:
		m.Events = msg.Events
		m.Err = msg.Err
		m.table = utils.CreateTable(m.Events, m.windowHeight)
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
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.Err != nil {
		return fmt.Sprintf("Error fetching events. %s\n", m.Err.Error())
	}
	if len(m.Events) == 0 {
		return "No events found\n"
	}
	return utils.BaseStyle.Render(m.table.View())
}
