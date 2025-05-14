package tui

import (
	"mcli/types"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Filter struct {
	Input   textinput.Model
	Text    string
	visible bool
}

func NewFilter() Filter {

	filterInput := textinput.New()
	filterInput.Placeholder = "Filter by title, location, description..."
	filterInput.CharLimit = 50
	filterInput.Width = 50

	return Filter{
		Input:   filterInput,
		Text:    "",
		visible: false,
	}
}

func (f *Filter) IsFiltering() bool {
	return f.visible
}

func (f *Filter) ToggleFilterView() {
	f.visible = !f.visible
}

func (f *Filter) View() string {
	if !f.visible {
		return ""
	}
	// TODO: stylize the input field
	return f.Input.View()
}
func (f *Filter) Update(msg tea.Msg) (Filter, tea.Cmd) {
	var cmd tea.Cmd
	f.Input, cmd = f.Input.Update(msg)
	return *f, cmd
}

// FilterEvents returns only those events whose title, location or description
// contains the query (case‑insensitive).
func FilterEvents(events []types.Event, query string) []types.Event {
	q := strings.ToLower(query)
	var filtered []types.Event
	for _, e := range events {
		if strings.Contains(strings.ToLower(e.Title), q) ||
			strings.Contains(strings.ToLower(e.Location), q) ||
			strings.Contains(strings.ToLower(e.Description), q) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
