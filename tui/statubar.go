package tui

import (
	"mcli/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// StatusBar represents a one-line status bar component.
type StatusBar struct {
	helpText     string // Text to display on the left (e.g., help menu)
	FilteredText string // Text to display on the right (e.g., current filter)
	Width        int    // Width of the status bar, typically the terminal width
}

// NewStatusBar creates a new StatusBar instance.
func NewStatusBar(helpText, filteredText string, width int) StatusBar {
	return StatusBar{
		helpText:     helpText,
		FilteredText: filteredText,
		Width:        width,
	}
}

// View renders the status bar as a single line with help text on the left and filter text on the right.
func (s StatusBar) View() string {
	// Prepare left and right content
	left := s.helpText
	right := s.FilteredText

	// Truncate text if it exceeds half the width to prevent overlap
	if lipgloss.Width(left) > s.Width/2 {
		left = lipgloss.NewStyle().MaxWidth(s.Width / 2).Render(left)
	}
	if lipgloss.Width(right) > s.Width/2 {
		right = lipgloss.NewStyle().MaxWidth(s.Width / 2).Render(right)
	}

	// Join the left and right parts with a spacer in between
	statusBar := lipgloss.JoinHorizontal(
		lipgloss.Left,
		left,
		lipgloss.NewStyle().Width(s.Width-lipgloss.Width(left)-lipgloss.Width(right)).Render(""),
		right,
	)

	// Apply styling to the entire status bar
	style := lipgloss.NewStyle().
		Background(styles.DefaultTheme.StatusBackground).
		Foreground(styles.DefaultTheme.StatusForeground).
		Width(s.Width)

	return style.Render(statusBar)
}
