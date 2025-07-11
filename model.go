package main

import (
	"fmt"
	"mcli/cmdprompt"
	"mcli/tui"
	"mcli/tui/styles"
	"mcli/types"
	"mcli/utils"
	"net/url"
	"strings"

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
	cmdPrompt *cmdprompt.CommandPrompt
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
		cmdPrompt: cmdprompt.New(":", handleCommand),
		statusbar: tui.NewStatusBar("Press 'q' to quit, '/' to filter, 'y' for deatils", "", 80),
	}
}

// Init starts the application by fetching events
func (m model) Init() tea.Cmd {
	utils.Logger.Debug("Init Called")
	m.cmdPrompt.Init()
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
				filterText := m.filter.Text
				if filterText != "" {
					filterText = "/" + filterText
				}
				m.statusbar.FilteredText = filterText // Update filter text
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

		if m.sidebar.IsVisible() {
			switch msg.String() {
			case "q", "esc":
				m.sidebar.ToggleSidebarView()
				m.sidebar.Viewport.GotoTop()
				m.AdjustViewports()
				return m, nil
			}
		}

		//handle command prompt
		consumed, updatedPrompt, _cmd := m.cmdPrompt.Update(msg, handleCommand)
		m.cmdPrompt = updatedPrompt
		if consumed {
			// If CommandPrompt handled the message, return early
			if m.cmdPrompt.GetOutput() == "Quitting..." {
				return m, tea.Quit
			}
			return m, _cmd
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "y":
			m.sidebar.ToggleSidebarView()
			// always render viewport from top
			m.sidebarMovement(msg)
		case "/":
			utils.Logger.Info("Filtering the entries")
			m.filter.ToggleFilterView()
			m.AdjustViewports()
			if m.filter.IsFiltering() {
				m.filter.Input.Focus()
				return m, textinput.Blink
			}
			return m, nil
		case "r":
			return m, utils.FetchEventCmd

		case "o":
			// open link in browser
			events := m.DisplayedEvents(m.filter.Text)
			utils.OpenURL(events[m.table.Cursor()].Url)
			return m, nil

		case "h":
			m.table.MoveDown(1)
			if m.sidebar.IsVisible() {
				m.sidebarMovement(msg)
			}

		case "l":
			m.table.MoveUp(1)
			if m.sidebar.IsVisible() {
				m.sidebarMovement(msg)
			}
		}

		var cmd tea.Cmd
		if m.sidebar.IsVisible() {
			m.sidebar, cmd = m.sidebar.Update(msg)
			return m, cmd
		}

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

	// Include CommandPrompt's view in the main app's view
	cmdBarView := m.cmdPrompt.View()
	renderedView = lipgloss.JoinVertical(lipgloss.Left, renderedView, cmdBarView)

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

func (m *model) sidebarMovement(msg tea.Msg) (tea.Model, tea.Cmd) {

	m.sidebar.Viewport.GotoTop()
	filteredEvents := m.DisplayedEvents(m.filter.Text)
	if m.sidebar.IsVisible() && len(filteredEvents) > 0 {
		cursor := m.table.Cursor()
		event := filteredEvents[cursor]
		m.sidebar.UpdateSidebarContent(event, m.termSize.height)
		utils.Logger.Info("Inspecting details on", "event", event.ID)
	}
	m.AdjustViewports()
	m.DebugLayout()
	var cmd tea.Cmd
	m.sidebar, cmd = m.sidebar.Update(msg)
	return m, cmd
}

func handleCommand(command string) (string, tea.Cmd) {

	var availableOpts = []string{"refresh", "fetch", "quit", "help"}
	_cmd := strings.Split(command, " ")
	switch strings.ToLower(_cmd[0]) {
	case "":
		// no command entered
		return "", nil
	case "help":
		return fmt.Sprintf("Check available opts: %s", strings.Join(availableOpts, ",")), nil
	case "quit":
		return "Quitting...", nil
	case "refresh":
		return "Refreshing list", nil
	case "fetch":
		args := strings.Join(_cmd[1:], " ")

		return fmt.Sprintf("Fetching events for %s", args), FetchEventByLocation(url.PathEscape(args))
	default:
		return fmt.Sprintf("Unknown command: %s", command), nil
	}
}

func FetchEventByLocation(loc string) tea.Cmd {
	return func() tea.Msg {
		return utils.FetchEventByLocationCmd(loc)
	}
}
