package main

import (
	"fmt"

	"mcli/types"
	"mcli/utils"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type windowSize struct {
	height int
	width  int
}

type model struct {
	Events     []types.Event
	table      table.Model
	windowSize windowSize
	Err        error
	logger     *utils.Logger
}

// updateTableMsg is a custom message to trigger a table update
type updateTableMsg struct{}

// updateTableCmd returns a command that dispatches an updateTableMsg
func updateTableCmd() tea.Cmd {
	return func() tea.Msg {
		return updateTableMsg{}
	}
}

func initModel(debug bool) model {
	defaultSize := windowSize{height: 10, width: 80}
	logger := utils.NewLogger(debug)
	return model{
		Events:     []types.Event{},
		table:      utils.CreateTable([]types.Event{}, defaultSize.height, defaultSize.width),
		windowSize: defaultSize,
		logger:     logger,
	}
}

func (m model) Init() tea.Cmd {
	m.logger.GetLogger().Debug("Init Called")
	return func() tea.Msg {
		events, err := utils.FetchEvents()
		return types.EventsMsg{
			Events: events,
			Err:    err,
		}
	}
}

// updateTable updates the table with the current events and window dimensions,
// adjusting the height based on whether the debug view is shown.
func (m *model) updateTable() {
	adjustedHeight := m.windowSize.height
	if m.logger.IsDebugViewShown() {
		adjustedHeight = m.windowSize.height - m.logger.GetDebugPanelHeight()
	}
	m.table = utils.CreateTable(m.Events, adjustedHeight, m.windowSize.width)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.logger.GetLogger().Debug("update/tea.WindowSizeMsg", "type", msg)
		m.windowSize = windowSize{
			height: msg.Height,
			width:  msg.Width,
		}
		return m, updateTableCmd()

	case types.EventsMsg:
		m.logger.GetLogger().Debug("update/tea.EventMsg")
		m.Events = msg.Events
		m.Err = msg.Err
		return m, updateTableCmd()

	case updateTableMsg:
		m.updateTable()
		return m, nil

	case tea.KeyMsg:
		m.logger.GetLogger().Info("update/key pressed", "key", msg.String())
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			m.logger.GetLogger().Info("enter pressed")
			return m, nil
		case "`":
			m.logger.ToggleDebugView()
			return m, updateTableCmd()
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
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.Err != nil {
		return fmt.Sprintf("Error fetching events. %s\n", m.Err.Error())
	}
	if len(m.Events) == 0 {
		return "No events found\n"
	}

	mainContent := m.table.View()
	if !m.logger.IsDebugViewShown() {
		return utils.BaseStyle.
			Width(m.windowSize.width - 2). // Adjust for border
			Render(mainContent)
	}

	debugPanel := m.logger.RenderDebugPanel(m.windowSize.width - 2) // Adjust for border
	return utils.BaseStyle.
		Width(m.windowSize.width - 2). // Adjust for border
		Render(lipgloss.JoinVertical(lipgloss.Left, mainContent, debugPanel))
}
