package tui

import (
	"mcli/tui/styles"
	"mcli/types"
	"mcli/utils"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Table struct {
	model table.Model
}

func getTableColumns(width int, isSidebarVisible bool) []table.Column {

	locationWidth := int(float64(width) * 0.2)
	dateWidth := int(float64(width) * 0.2)
	remaining := width - locationWidth - dateWidth
	eventWidth := int(float64(remaining) * 1)

	utils.Logger.Info("getTableColumns", "width", width)
	utils.Logger.Info("getTableColumns", "locationWidth", locationWidth)
	utils.Logger.Info("getTableColumns", "dateWidth", dateWidth)
	utils.Logger.Info("getTableColumns", "eventWidth", eventWidth)

	// if showTitleOnly is set, i,e when sidebar is visible
	if isSidebarVisible {
		locationWidth = 0
		dateWidth = 10
	}
	return []table.Column{
		{Title: "🚀", Width: 2},
		{Title: "Event", Width: eventWidth},
		{Title: "Location", Width: locationWidth},
		{Title: "Date", Width: dateWidth},
	}
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
	width := 0 // initial width size of table, will be adjusted dynamically
	showTitleOnly := false
	t := table.New(
		table.WithColumns(getTableColumns(width, showTitleOnly)),
		table.WithRows(createTableRows(events)),
		table.WithFocused(true),
	)
	t.SetStyles(styles.GetTableStyles())
	return Table{
		model: t,
	}
}

// dynamicaly adjust the column width
func (t *Table) UpdateColumnWidth(termWidth int, sidebarVisible bool, events []types.Event) {

	showTitleOnly := false
	if sidebarVisible || termWidth < 40 {
		showTitleOnly = true
	}

	columns := getTableColumns(t.model.Width(), showTitleOnly)
	t.model.SetColumns(columns)
	t.model.SetRows(createTableRows(events))
}

// set table width and update column widths
func (t *Table) SetWidth(width int, termWidth int, sidebarVisible bool, events []types.Event) {
	t.model.SetWidth(width)
	t.UpdateColumnWidth(termWidth, sidebarVisible, events)
}

// get table width
func (t *Table) Width() int {
	return t.model.Width()
}

// set table height
func (t *Table) SetHeight(height int) {
	t.model.SetHeight(height)
}

// get table height
func (t *Table) Height() int {
	return t.model.Height()
}

// pass down updates to the table
func (t *Table) Update(msg tea.Msg) (Table, tea.Cmd) {
	var cmd tea.Cmd
	t.model, cmd = t.model.Update(msg)
	return *t, cmd
}

// render table with appropirate styling
func (t *Table) View(termWidth int) string {
	tableStyle := lipgloss.NewStyle().Width(t.model.Width()).MaxWidth(termWidth)
	return tableStyle.Render(t.model.View())
}

// SelectedRow returns the currently selected row.
func (t Table) Cursor() int {
	// Ensure we don't return an empty row if there are no rows
	return t.model.Cursor()
}
