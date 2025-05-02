package model

import (
	"fmt"
	"strings"

	"github.com/thapakazi/mcli/api"
	"github.com/thapakazi/mcli/types"
	"github.com/thapakazi/mcli/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// Import glamour for Markdown rendering
)

// Model holds the application state.
type Model struct {
	Events         []types.Event
	FilteredEvents []types.Event
	Cursor         int
	ViewportTop    int
	ViewportHeight int
	Err            error
	Loading        bool
	TermHeight     int
	Filtering      bool
	FilterInput    textinput.Model
	FilterText     string
	ViewMode       string             // "list" or "details"
	SelectedEvent  *types.EventDetail // Details of the selected event
}

// NewModel initializes the application state.
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Filter by title..."
	ti.CharLimit = 50
	ti.Width = 50

	return Model{
		Events:         []types.Event{},
		FilteredEvents: []types.Event{},
		Cursor:         0,
		ViewportTop:    0,
		ViewportHeight: 10,
		Loading:        true,
		TermHeight:     0,
		Filtering:      false,
		FilterInput:    ti,
		FilterText:     "",
		ViewMode:       "list",
		SelectedEvent:  nil,
	}
}

// Init implements tea.Model.Init.
func (m Model) Init() tea.Cmd {
	return tea.Batch(func() tea.Msg {
		events, err := api.FetchEvents()
		if err != nil {
			return types.ErrMsg{Err: err}
		}
		return types.EventsMsg(events)
	}, tea.WindowSize())
}

// Update implements tea.Model.Update.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.TermHeight = msg.Height
		m.ViewportHeight = msg.Height - 4
		if m.Filtering {
			m.ViewportHeight -= 2
		}
		if m.ViewportHeight < 1 {
			m.ViewportHeight = 1
		}
		m.AdjustViewport()
		return m, nil

	case types.EventsMsg:
		m.Events = msg
		m.FilteredEvents = m.Events
		m.Loading = false
		if len(m.Events) == 0 {
			m.Err = fmt.Errorf("no events found")
		}
		m.AdjustViewport()
	case types.EventDetailMsg:
		m.Loading = false
		if msg.Err != nil {
			m.Err = msg.Err
			m.ViewMode = "list"
			m.SelectedEvent = nil
			return m, nil
		}
		m.SelectedEvent = msg.Detail
	case types.ErrMsg:
		m.Err = msg.Err
		m.Loading = false
		return m, tea.Quit
	case tea.KeyMsg:
		if m.Err != nil {
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			return m, nil
		}

		if m.Filtering {
			switch msg.String() {
			case "esc":
				m.Filtering = false
				m.FilterText = ""
				m.FilterInput.SetValue("")
				m.FilteredEvents = m.Events
				m.Cursor = 0
				m.ViewportTop = 0
				m.ViewportHeight = m.TermHeight - 4
				m.AdjustViewport()
			case "enter":
				m.Filtering = false
				m.FilterText = m.FilterInput.Value()
				m.ViewportHeight = m.TermHeight - 4
				m.AdjustViewport()
			default:
				var cmd tea.Cmd
				m.FilterInput, cmd = m.FilterInput.Update(msg)
				m.FilterText = m.FilterInput.Value()
				m.filterEvents()
				m.Cursor = 0
				m.ViewportTop = 0
				m.AdjustViewport()
				return m, cmd
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "/":
			if m.ViewMode == "list" {
				m.Filtering = true
				m.FilterInput.Focus()
				m.ViewportHeight = m.TermHeight - 6
				m.AdjustViewport()
				return m, textinput.Blink
			}
		case "up", "k":
			if m.ViewMode == "list" {
				if m.Cursor > 0 {
					m.Cursor--
					m.AdjustViewport()
				}
			}
		case "down", "j":
			if m.ViewMode == "list" {
				if m.Cursor < len(m.FilteredEvents)-1 {
					m.Cursor++
					m.AdjustViewport()
				}
			}
		case "enter":
			if m.ViewMode == "list" && len(m.FilteredEvents) > 0 {
				m.ViewMode = "details"
				m.Loading = true
				selectedEvent := m.FilteredEvents[m.Cursor]
				return m, func() tea.Msg {
					detail, err := api.FetchEventDetail(selectedEvent)
					return types.EventDetailMsg{Detail: detail, Err: err}
				}
			} else if m.ViewMode == "details" {
				m.ViewMode = "list"
				m.SelectedEvent = nil
				m.Loading = false
				m.AdjustViewport()
			}
		case "esc":
			if m.ViewMode == "details" {
				m.ViewMode = "list"
				m.SelectedEvent = nil
				m.Loading = false
				m.AdjustViewport()
			}
		case "o": // Open the URL in the browser
			if m.ViewMode == "details" && m.SelectedEvent != nil {
				utils.OpenURL(m.SelectedEvent.URL)
			}
		}
	}
	return m, nil
}


