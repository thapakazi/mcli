package model

import (
	"fmt"

	"github.com/charmbracelet/glamour"
)

// formatTicketPrice converts ticketPrice to a string for display.
func formatTicketPrice(price interface{}) string {
	if price == nil {
		return "N/A"
	}
	switch v := price.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("$%.2f", v)
	case int:
		return fmt.Sprintf("$%d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// formatDescription applies markdown rendering using Glamour to the description.
func formatDescription(desc string) string {
	// Use Glamour to render the Markdown description
	rendered, _ := glamour.Render(desc, "dark") // Enable color rendering for the terminal

	// Return the formatted markdown output
	return rendered
}
