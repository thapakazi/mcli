package model

import (
	"strings"

	"github.com/thapakazi/mcli/types"
)

// AdjustViewport ensures the cursor is visible within the viewport.
func (m *Model) AdjustViewport() {
	viewportBottom := m.ViewportTop + m.ViewportHeight - 1

	if m.Cursor < m.ViewportTop {
		m.ViewportTop = m.Cursor
	}

	if m.Cursor > viewportBottom {
		m.ViewportTop = m.Cursor - m.ViewportHeight + 1
	}

	if m.ViewportTop < 0 {
		m.ViewportTop = 0
	}

	if m.ViewportTop > len(m.FilteredEvents)-m.ViewportHeight && len(m.FilteredEvents) >= m.ViewportHeight {
		m.ViewportTop = len(m.FilteredEvents) - m.ViewportHeight
	}
}

// filterEvents filters events based on a search term.
func (m *Model) filterEvents() {
	if m.FilterText == "" {
		m.FilteredEvents = m.Events
		return
	}
	term := strings.ToLower(m.FilterText)
	var filtered []types.Event
	for _, event := range m.Events {
		if strings.Contains(strings.ToLower(event.Title), term) {
			filtered = append(filtered, event)
		}
	}
	m.FilteredEvents = filtered
}
