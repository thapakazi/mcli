package main

import (
	"mcli/tui"
	"mcli/tui/styles"
	"mcli/types"
	"mcli/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type termSize struct {
	height int
	width  int
}

type model struct {
	Events   types.Events
	table    tui.Table
	sidebar  tui.Sidebar
	filter   tui.Filter
	termSize termSize
	loading  bool
	err      error
}

func NewModel() model {
	return model{
		loading: true,
		table:   tui.NewTable(types.Events{}),
		sidebar: tui.NewSidebar(0),
		filter:  tui.NewFilter(),
	}
}

// Call fetchEvents to populate the table
func (m model) Init() tea.Cmd {
	utils.Logger.Debug("Init Called")
	//return tea.Batch(utils.FetchEventCmd, tea.WindowSize())
	// m.table = tui.NewTable(m.DisplayedEvents(""))
	// m.sidebar = tui.NewSidebar(0)
	// m.filter = tui.NewFilter()
	return utils.FetchEventCmd

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {

	case utils.FetchErrorMsg:
		utils.Logger.Debug("update/tea.FetchErrorMsg")
		m.loading = false
		m.err = msg.Err
		return m, nil

	case utils.FetchSuccessMsg:
		utils.Logger.Debug("update/tea.FetchSuccessMsg")
		m.loading = false
		m.Events = msg.Events
		return m, tea.WindowSize()

	case tea.WindowSizeMsg:
		utils.Logger.Debug("update/tea.WindowSizeMsg", "type", msg)
		m.termSize.height = msg.Height
		m.termSize.width = msg.Width
		m.AdjustViewports()
		m.DebugLayout()
		return m, nil

	case tea.KeyMsg:
		utils.Logger.Info("update/key pressed", "key", msg.String())
		if m.filter.IsFiltering() {
			switch msg.String() {
			case "esc":
				m.filter.ToggleFilterView()
				m.table = tui.NewTable(m.DisplayedEvents(""))
			case "enter":
				m.filter.ToggleFilterView()
				m.filter.Text = m.filter.Input.Value()
				m.table = tui.NewTable(m.DisplayedEvents(m.filter.Text))
				m.filter.Input.Blur()
			default:
				var cmd tea.Cmd
				m.filter, cmd = m.filter.Update(msg)
				m.filter.Text = m.filter.Input.Value()
				m.table = tui.NewTable(m.DisplayedEvents(m.filter.Text))
				utils.Logger.Info("filtering list", "text", m.filter.Text)
				return m, cmd
			}
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "y":
			m.sidebar.ToggleSidebarView()

			filteredEvents := m.DisplayedEvents(m.filter.Text)
			utils.Logger.Info("filteredEvents", "count", len(filteredEvents))
			if m.sidebar.IsVisible() {
				cursor := m.table.Cursor()
				event := filteredEvents[cursor]
				m.sidebar.UpdateSidebarConntent(event, m.termSize.height)
				utils.Logger.Info("Inspecting details on", "event", event.ID)
			}
			m.AdjustViewports()
			m.DebugLayout()
			return m, nil

		case "/":
			utils.Logger.Info("Filtering the entries")
			m.filter.ToggleFilterView()
			tableHeight := m.termSize.height
			if m.filter.IsFiltering() {
				tableHeight = m.termSize.height - 2 // input filed area
			}
			m.table.SetHeight(tableHeight)
			m.filter.Input.Focus()
			return m, textinput.Blink

		}

		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	// Ignore table/sidebar updates while loading or if there's an error
	if m.loading || m.err != nil {
		return m, nil
	}
	if !m.sidebar.IsVisible() {
		m.table, cmd = m.table.Update(msg)
	} else if m.filter.IsFiltering() {
		m.filter, cmd = m.filter.Update(msg)
	} else {
		m.sidebar, cmd = m.sidebar.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.loading {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("6")). // Cyan
			Render("Loading...")
	}

	// Show error if fetch failed
	if m.err != nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Render("Error: " + m.err.Error())
	}
	if len(m.Events) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("3")).
			Render("No events found\n")
	}

	tableView := m.table.View(m.termSize.width - 2)
	if m.sidebar.IsVisible() {

		return styles.BaseStyle.
			Render(lipgloss.JoinHorizontal(
				lipgloss.Left,
				tableView,
				m.sidebar.View()))
	}

	if m.filter.IsFiltering() {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			tableView,
			m.filter.View())
	}
	return styles.BaseStyle.Render(tableView)

}

// DisplayedEvents returns either the filtered slice (if filter.Text != "")
// or the full list.
func (m model) DisplayedEvents(filter string) []types.Event {
	if filter == "" {
		return m.Events
	}
	return tui.FilterEvents(m.Events, filter)
}

func (m *model) AdjustViewports() {

	// set sidebarWidth be ceratin % of total available width
	sidebarWidth := int(float64(m.termSize.width) * 0.65)
	m.sidebar.Width = max(sidebarWidth, 10) // Ensure minimum width

	// adjust table dimension if not loading
	if !m.loading && m.err == nil {

		// table spans over the available terminal width
		tableWidth := m.termSize.width
		sidebarIsVisible := m.sidebar.IsVisible()
		if sidebarIsVisible {
			tableWidth = m.termSize.width - sidebarWidth
		}
		m.table.SetWidth(tableWidth, m.termSize.width, sidebarIsVisible, m.DisplayedEvents(m.filter.Text))
		m.table.SetHeight(m.termSize.height - 3)
	}
}

func (m *model) DebugLayout() {
	utils.Logger.Debug("table", "width", m.table.Width())
	utils.Logger.Debug("table", "height", m.table.Height())
	utils.Logger.Debug("sidebar", "width", m.sidebar.Width)
	utils.Logger.Debug("sidebar", "height", m.sidebar.GetHeight())
}
