package tui

import (
	"mcli/tui/styles"
	"mcli/types"
	"mcli/utils"

	"github.com/charmbracelet/bubbles/table"
)

type Table struct {
	table.Model
}

func getTableColumns(width int, isSidebarVisible bool) []table.Column {

	locationWidth := int(float64(width) * 0.2)
	dateWidth := int(float64(width) * 0.2)
	remaining := width - locationWidth - dateWidth
	eventWidth := int(float64(remaining) * 1)

	// if showTitleOnly is set, i,e when sidebar is visible
	if isSidebarVisible {
		locationWidth = 0
		dateWidth = 10
	}
	// utils.Logger.Info("getTableColumns", "width", width)
	// utils.Logger.Info("getTableColumns", "locationWidth", locationWidth)
	// utils.Logger.Info("getTableColumns", "dateWidth", dateWidth)
	// utils.Logger.Info("getTableColumns", "eventWidth", eventWidth)

	return []table.Column{
		{Title: "🚀", Width: 2},
		{Title: "Event", Width: eventWidth},
		{Title: "Location", Width: locationWidth},
		{Title: "Date", Width: dateWidth},
	}
}

func CreateTableRows(events []types.Event) []table.Row {
	var rows []table.Row
	for _, event := range events {
		sourceIcon := "?"
		switch event.Source {
		case "luma":
			sourceIcon = "✦︎"
		default:
			sourceIcon = "☘️"
		}
		title := event.Title
		_, _, dateTime, _ := utils.ParseAndCompareDateTime(event.DateTime)
		//if isFutureOrCurrent {
		rows = append(rows, table.Row{
			sourceIcon,
			title,
			event.Location,
			dateTime,
		})

		//}
	}
	return rows
}

// Initialize new table
func NewTable(events []types.Event) Table {
	width := 20 // initial width size of table, will be adjusted dynamically
	showTitleOnly := false
	t := table.New(
		table.WithColumns(getTableColumns(width, showTitleOnly)),
		table.WithRows(CreateTableRows(events)),
		table.WithFocused(true),
	)
	t.SetStyles(styles.GetTableStyles())
	return Table{t}
}

// dynamicaly adjust the column width
func (t *Table) AdjustColumns(termWidth int, sidebarVisible bool) {

	showTitleOnly := false
	if sidebarVisible || termWidth < 40 {
		showTitleOnly = true
	}

	columns := getTableColumns(t.Width(), showTitleOnly)
	t.SetColumns(columns)
}
