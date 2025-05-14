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
	logger   *utils.Logger
	termSize termSize
	loading  bool
	err      error
}

func NewModel(debug bool) model {
	return model{
		loading: true,
		logger:  utils.NewLogger(debug),
	}
}

// Call fetchEvents to populate the table
func (m model) Init() tea.Cmd {
	m.logger.GetLogger().Debug("Init Called")
	return utils.FetchEventCmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {

	case utils.FetchErrorMsg:
		m.logger.GetLogger().Debug("update/tea.FetchErrorMsg")
		m.loading = false
		m.err = msg.Err
		return m, nil

	case utils.FetchSuccessMsg:
		m.logger.GetLogger().Debug("update/tea.FetchSuccessMsg")
		m.loading = false
		m.Events = msg.Events
		m.table = tui.NewTable(m.DisplayedEvents(""))
		m.sidebar = tui.NewSidebar(0)
		m.filter = tui.NewFilter()

	case tea.WindowSizeMsg:
		m.logger.GetLogger().Debug("update/tea.WindowSizeMsg", "type", msg)
		m.termSize.height = msg.Height
		m.termSize.width = msg.Width
		m.AdjustViewports()
		m.DebugLayout()
		return m, nil

	case tea.KeyMsg:
		m.logger.GetLogger().Info("update/key pressed", "key", msg.String())
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
				m.logger.GetLogger().Info("filtering list", "text", m.filter.Text)
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
			m.logger.GetLogger().Info("filteredEvents", "count", len(filteredEvents))
			if m.sidebar.IsVisible() {
				cursor := m.table.Cursor()
				event := filteredEvents[cursor]
				m.sidebar.UpdateSidebarConntent(event, m.termSize.height)
				m.logger.GetLogger().Info("Inspecting details on", "event", event.ID)
			}
			m.AdjustViewports()
			m.DebugLayout()
			return m, nil

		case "/":
			m.logger.GetLogger().Info("Filtering the entries")
			m.filter.ToggleFilterView()
			tableHeight := m.termSize.height
			if m.filter.IsFiltering() {
				tableHeight = m.termSize.height - 2 // input filed area
			}
			m.table.SetHeight(tableHeight)
			m.filter.Input.Focus()
			return m, textinput.Blink

		case "`":
			m.logger.ToggleDebugView()
			return m, nil
		case "up":
			if m.logger.IsDebugViewShown() {
				m.logger.ScrollUp()
				return m, nil
			}
		case "down":
			if m.logger.IsDebugViewShown() {
				m.logger.ScrollDown()
				return m, nil
			}
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

	tableView := m.table.View(m.termSize.width)

	// if !m.logger.IsDebugViewShown() {
	// 	return styles.BaseStyle.
	// 		Width(m.termSize.width - 2). // Adjust for border
	// 		Render(mainContent)
	// }
	if m.sidebar.IsVisible() {

		//m.logger.GetLogger().Debug("Descirption", "desc", sideBarContent)
		return styles.BaseStyle.
			Width(m.termSize.width - 2). // Adjust for border
			Render(
				lipgloss.JoinHorizontal(lipgloss.Left, tableView, m.sidebar.View()),
			)
	}

	if m.filter.IsFiltering() {
		return lipgloss.JoinVertical(lipgloss.Left, tableView, m.filter.View())
	}

	// debugPanel := m.logger.RenderDebugPanel(m.termSize.width - 2) // Adjust for border
	// return styles.BaseStyle.
	// 	Width(m.termSize.width - 2). // Adjust for border
	// 	Render(lipgloss.JoinVertical(lipgloss.Left, mainContent, debugPanel))
	return tableView
}

// DisplayedEvents returns either the filtered slice (if filter.Text != "")
// or the full list.
func (m model) DisplayedEvents(filter string) []types.Event {
	if filter == "" {
		return m.Events
	}
	return tui.FilterEvents(m.Events, filter)
}

func (m *model) SortByDateTime() {

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
		m.table.SetHeight(m.termSize.height - 1)
	}

}

func (m model) DebugLayout() {
	m.logger.GetLogger().Debug("table", "width", m.table.Width())
	m.logger.GetLogger().Debug("table", "height", m.table.Height())
	m.logger.GetLogger().Debug("sidebar", "width", m.sidebar.Width)
	m.logger.GetLogger().Debug("sidebar", "height", m.sidebar.GetHeight())
}
