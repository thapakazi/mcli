package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model.View.
func (m Model) View() string {
	if m.Loading {
		return "Loading...\n"
	}
	if m.Err != nil {
		return fmt.Sprintf("Error: %v\nPress q to quit.", m.Err)
	}

	// Render details view
	if m.ViewMode == "details" && m.SelectedEvent != nil {
		// Render event details view
		s := strings.Builder{}

		headerStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

		titleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF3333")).
			Bold(true)

		dateStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF"))

		metadataStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00"))

		urlStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF33FF")).
			Underline(true)

		bodyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			PaddingLeft(2)

		// Header
		s.WriteString(headerStyle.Render("Event Details (Press Enter or Esc to go back, q to quit)\n\n"))

		d := m.SelectedEvent
		details := []string{}

		details = append(details, titleStyle.Render(fmt.Sprintf("Title: %s", d.Title)))

		if d.GroupName != "" {
			details = append(details, metadataStyle.Render(fmt.Sprintf("Group: %s", d.GroupName)))
		}

		details = append(details,
			dateStyle.Render(fmt.Sprintf("Date: %s", d.DateTime)),
			metadataStyle.Render(fmt.Sprintf("Type: %s", d.EventType)),
			dateStyle.Render(fmt.Sprintf("Venue: %s, %s, %s", d.VenueName, d.City, d.State)),
		)

		if d.TicketCount > 0 {
			details = append(details, metadataStyle.Render(fmt.Sprintf("Tickets Remaining: %d/%d", d.TicketRemaining, d.TicketCount)))
		}

		if d.RsvpsCount > 0 {
			details = append(details, metadataStyle.Render(fmt.Sprintf("RSVPs: %d", d.RsvpsCount)))
		}

		if d.TicketPrice != nil {
			details = append(details, metadataStyle.Render(fmt.Sprintf("Price: %s", formatTicketPrice(d.TicketPrice))))
		}

		details = append(details, urlStyle.Render(fmt.Sprintf("URL: %s", d.URL)))

		if d.Description != "" {
			formattedDesc := formatDescription(d.Description)
			details = append(details, fmt.Sprintf("\nDescription:\n%s", formattedDesc))
		}

		s.WriteString(bodyStyle.Render(strings.Join(details, "\n")))

		return s.String()
	}

	// Render event list view
	if len(m.Events) == 0 {
		return "No events available.\nPress q to quit."
	}

	normalStyle := lipgloss.NewStyle().
		PaddingLeft(2)

	highlightStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		Background(lipgloss.Color("1")).
		Foreground(lipgloss.Color("15")).
		Width(80).
		Align(lipgloss.Left)

	// Build event list
	s := strings.Builder{}
	s.WriteString("Events List (Use ↑/↓ to navigate, / to filter, Enter for details, q to quit)\n\n")

	if m.Filtering {
		s.WriteString("Filter: " + m.FilterInput.View() + "\n\n")
	}

	if len(m.FilteredEvents) == 0 {
		s.WriteString("No events match the filter.\n")
		return s.String()
	}

	start := m.ViewportTop
	end := m.ViewportTop + m.ViewportHeight
	if end > len(m.FilteredEvents) {
		end = len(m.FilteredEvents)
	}

	for i := start; i < end; i++ {
		event := m.FilteredEvents[i]
		prefix := "  "
		if m.Cursor == i {
			prefix = "> "
		}

		line := fmt.Sprintf("%s%s", prefix, event.Title)
		if m.Cursor == i {
			line = highlightStyle.Render(line)
		} else {
			line = normalStyle.Render(line)
		}
		s.WriteString(line + "\n")
	}

	if m.ViewportTop > 0 {
		s.WriteString("↑ More events above...\n")
	}
	if end < len(m.FilteredEvents) {
		s.WriteString("↓ More events below...\n")
	}

	return s.String()
}
