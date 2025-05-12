package main

import (
	"mcli/tui"
	"mcli/tui/styles"
	"mcli/types"
	"mcli/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type termSize struct {
	height int
	width  int
}

type model struct {
	Events   []types.Event
	table    tui.Table
	sidebar  tui.Sidebar
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
		m.table = tui.NewTable(m.Events)
		m.sidebar = tui.NewSidebar(0)

	case tea.WindowSizeMsg:
		m.logger.GetLogger().Debug("update/tea.WindowSizeMsg", "type", msg)
		m.termSize.height = msg.Height
		m.termSize.width = msg.Width

		// set sidebarWidth be ceratin % of total available width
		sidebarWidth := int(float64(m.termSize.width) * 0.65)
		if sidebarWidth < 10 {
			sidebarWidth = 10 // Ensure minimum width
		}
		m.sidebar.Width = sidebarWidth

		// adjust table dimension if not loading
		if !m.loading && m.err == nil {

			tableWidth := m.termSize.width
			sidebarIsVisible := m.sidebar.IsVisible()
			if sidebarIsVisible {
				tableWidth = m.termSize.width - sidebarWidth
			}
			m.table.SetWidth(tableWidth, m.termSize.width, sidebarIsVisible, m.Events)
			m.table.SetHeight(m.termSize.height - 1)
		}

		return m, nil

	case tea.KeyMsg:
		m.logger.GetLogger().Info("update/key pressed", "key", msg.String())
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "y":
			m.sidebar.ToggleSidebarView()

			tableWidth := m.termSize.width
			if m.sidebar.IsVisible() {
				tableWidth = m.termSize.width - m.sidebar.Width
				cursor := m.table.Cursor()
				event := m.Events[cursor]
				m.sidebar.UpdateSidebarConntent(event)
				m.logger.GetLogger().Info("Inspecting details on", "event", event.ID)
			}
			m.table.SetWidth(tableWidth, m.termSize.width, m.sidebar.IsVisible(), m.Events)

			return m, nil
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

	// debugPanel := m.logger.RenderDebugPanel(m.termSize.width - 2) // Adjust for border
	// return styles.BaseStyle.
	// 	Width(m.termSize.width - 2). // Adjust for border
	// 	Render(lipgloss.JoinVertical(lipgloss.Left, mainContent, debugPanel))
	return tableView
}
