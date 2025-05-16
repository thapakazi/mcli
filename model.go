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

// termSize holds the terminal dimensions
type termSize struct {
	height int
	width  int
}

// model represents the application state
type model struct {
	Events    types.Events
	table     tui.Table
	sidebar   tui.Sidebar
	statusbar tui.StatusBar
	filter    tui.Filter
	termSize  termSize
	loading   bool
	err       error
}

// NewModel initializes the application model
func NewModel() model {
	return model{
		loading:   true,
		table:     tui.NewTable(types.Events{}),
		sidebar:   tui.NewSidebar(),
		filter:    tui.NewFilter(),
		statusbar: tui.NewStatusBar("Press 'q' to quit, '/' to filter, 'y' for deatils", "", 80),
	}
}

// Init starts the application by fetching events
func (m model) Init() tea.Cmd {
	utils.Logger.Debug("Init Called")
	return utils.FetchEventCmd
}

// Update handles incoming messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.AdjustViewports()
		return m, nil

	case tea.WindowSizeMsg:
		utils.Logger.Debug("update/tea.WindowSizeMsg", "type", msg)
		m.termSize.height = msg.Height
		m.termSize.width = msg.Width
		m.statusbar.Width = msg.Width - 2
		m.AdjustViewports()
		m.DebugLayout()
		return m, nil

	case tea.KeyMsg:
		utils.Logger.Info("update/key pressed", "key", msg.String())
		if m.filter.IsFiltering() {
			switch msg.String() {
			case "esc":
				m.filter.ToggleFilterView()
				m.filter.Text = ""
				m.table.SetRows(tui.CreateTableRows(m.DisplayedEvents("")))
				m.statusbar.FilteredText = "" // Clear filter text
				m.AdjustViewports()
			case "enter":
				m.filter.ToggleFilterView()
				m.filter.Text = m.filter.Input.Value()
				m.table.SetRows(tui.CreateTableRows(m.DisplayedEvents(m.filter.Text)))
				m.statusbar.FilteredText = "/" + m.filter.Text // Update filter text
				m.AdjustViewports()
			default:
				var cmd tea.Cmd
				m.filter, cmd = m.filter.Update(msg)
				m.filter.Text = m.filter.Input.Value()
				m.table.SetRows(tui.CreateTableRows(m.DisplayedEvents(m.filter.Text)))
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
			if m.sidebar.IsVisible() && len(filteredEvents) > 0 {
				cursor := m.table.Cursor()
				//if cursor < len(filteredEvents) {
				event := filteredEvents[cursor]
				m.sidebar.UpdateSidebarContent(event, m.termSize.height)
				utils.Logger.Info("Inspecting details on", "event", event.ID)
				//	}
			}
			m.AdjustViewports()
			m.DebugLayout()
			var cmd tea.Cmd
			m.sidebar, cmd = m.sidebar.Update(msg)
			return m, cmd
		case "/":
			utils.Logger.Info("Filtering the entries")
			m.filter.ToggleFilterView()
			m.AdjustViewports()
			if m.filter.IsFiltering() {
				m.filter.Input.Focus()
				return m, textinput.Blink
			}
			return m, nil
		}

		var cmd tea.Cmd
		m.table.Model, cmd = m.table.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the current state of the application
func (m model) View() string {
	if m.loading {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Loading...")
	}
	if m.err != nil {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("Error: " + m.err.Error())
	}
	if len(m.Events) == 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("No events found\n")
	}

	// start from table rendering
	renderedView := m.table.View()

	// if user is filtering the text
	if m.filter.IsFiltering() {
		renderedView = lipgloss.JoinVertical(lipgloss.Top, renderedView, m.filter.View())
	}

	// if sidebar is visible
	if m.sidebar.IsVisible() {
		renderedView = lipgloss.JoinHorizontal(lipgloss.Left, renderedView, m.sidebar.View())
	}

	// Add status bar
	statusBarView := m.statusbar.View()
	renderedView = lipgloss.JoinVertical(lipgloss.Left, renderedView, statusBarView)

	return styles.BaseStyle.Render(renderedView)
}

// DisplayedEvents returns the current list of events based on the filter
func (m model) DisplayedEvents(filter string) []types.Event {
	if filter == "" {
		return m.Events
	}
	return tui.FilterEvents(m.Events, filter)
}

func (m *model) AdjustViewports() {

	// Calculate statubar & filter height
	statusbarHeight := 1
	filterHeight := 0
	if m.filter.IsFiltering() {
		filterHeight = 1
	}

	// Calculate table height
	tableHeight := m.termSize.height - statusbarHeight - filterHeight - 2 // 2 for border(head/tail)
	m.table.SetHeight(tableHeight)
	m.table.SetRows(tui.CreateTableRows(m.DisplayedEvents(m.filter.Text)))

	// Calculate table width
	m.statusbar.Width = m.termSize.width - 2
	tableWidth := m.termSize.width
	if m.sidebar.IsVisible() {
		m.sidebar.Width = int(float64(m.termSize.width) * 0.6) // Sidebar takes X%, adjustable
		tableWidth = m.termSize.width - m.sidebar.Width

		m.sidebar.SetViewportWidth(m.sidebar.Width - 3) // Border (1) + PaddingLeft (2)
		m.sidebar.SetViewportHeight(tableHeight - 3)    // PaddingTop (3)

	}
	m.table.SetWidth(tableWidth)
	m.table.AdjustColumns(tableWidth, m.sidebar.IsVisible())

}

// DebugLayout logs the current layout dimensions for debugging
func (m *model) DebugLayout() {
	utils.Logger.Debug("table", "width", m.table.Width())
	utils.Logger.Debug("table", "height", m.table.Height())
	utils.Logger.Debug("sidebar", "width", m.sidebar.Width)
	utils.Logger.Debug("sidebar", "height", m.sidebar.GetHeight())
}
