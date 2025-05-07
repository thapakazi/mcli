package utils

import (
	"mcli/types"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF"))
	sourceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	BaseStyle   = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
)

func getTableColumns() []table.Column {
	return []table.Column{
		{Title: "⚔️", Width: 2},
		{Title: "Event", Width: 90},
		{Title: "Location", Width: 50},
		{Title: "Date", Width: 12},
	}
}

func getTableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Foreground(lipgloss.Color("5")).
		Bold(true).
		BorderBottom(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	return s
}

func createTableRows(events []types.Event) []table.Row {
	var rows []table.Row
	for _, event := range events {

		sourceIcon := "?"
		switch event.Source {
		case "luma":
			sourceIcon = "✦︎"
		default:
			sourceIcon = "☘️"
		}
		title := titleStyle.Render(event.Title)
		rows = append(rows, table.Row{
			sourceIcon,
			title,
			event.Location,
			event.DateTime,
		})
	}
	return rows
}

func CreateTable(events []types.Event, height int) table.Model {

	availableHeight := max(height-5, 1)
	t := table.New(
		table.WithColumns(getTableColumns()),
		table.WithRows(createTableRows(events)),
		table.WithFocused(true),
		table.WithHeight(availableHeight),
	)
	t.SetStyles(getTableStyles())
	return t
}
