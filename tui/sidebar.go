package tui

import (
	"fmt"
	"mcli/types"
	"mcli/utils"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

type Sidebar struct {
	visible  bool
	viewport viewport.Model
	Width    int
	Height   int
}

func NewSidebar() Sidebar {

	utils.Logger.Info("NewSidebar called")
	width, height := 20, 20               // initial default width,height
	vp := viewport.New(width-2, height-4) // -2 for border and space to left, -4 for 2 space at top and bottom
	vp.Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))
	sidebar := Sidebar{
		visible:  false,
		viewport: vp,
		Width:    width,
		Height:   height,
	}
	utils.Logger.Info("Sidebar", "sidebar", sidebar)
	return sidebar
}
func (s *Sidebar) IsVisible() bool {
	return s.visible
}
func (s *Sidebar) ToggleSidebarView() {
	s.visible = !s.visible
}

func (s *Sidebar) UpdateSidebarContent(event types.Event, height int) {

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")).Render(event.Title)

	description, _ := glamour.Render(event.Description, "dark")
	if event.Description == "" {
		description, _ = glamour.Render("press R to fetch description", "dark")
	}

	url := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3")).Render(event.Url)

	parsedTime, _, _, _ := utils.ParseAndCompareDateTime(event.DateTime)
	date := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")).Render(parsedTime.String())
	styledDescription := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5")).Render("Description:\n------------")
	location := event.Location

	sidebarText := fmt.Sprintf(
		"%s\n\n🔗: %s\n\n📍: %s\n\n📅:%s\n\n%s\n%s",
		title, url, location, date, styledDescription, description,
	)
	// Split into lines
	lines := strings.Split(sidebarText, "\n")

	// Truncate each line to viewport width
	for i, line := range lines {
		lines[i] = truncate.StringWithTail(line, uint(90), "...")
	}
	truncatedContent := strings.Join(lines, "\n")

	// Calculate the actual content height (number of lines)
	// contentHeight := strings.Count(renderedContent, "\n") + 1

	//s.viewport.Height = height - 6
	//s.viewport.Width = s.Width - 4
	s.viewport.SetContent(truncatedContent)
}

func (s *Sidebar) Update(msg tea.Msg) (Sidebar, tea.Cmd) {
	var cmd tea.Cmd
	s.viewport, cmd = s.viewport.Update(msg)
	return *s, cmd
}

func (s Sidebar) View() string {
	if !s.visible {
		return ""
	}

	sidebarStyle := lipgloss.NewStyle().
		Width(s.Width).
		Border(lipgloss.NormalBorder(),
			false,
			false,
			false,
			true, // left border
		).BorderForeground(lipgloss.Color("63")). // Light blue border
		PaddingTop(3).
		PaddingLeft(2)

	return sidebarStyle.Render(s.viewport.View())
}

func (s Sidebar) GetHeight() int { return s.viewport.Height }
func (s Sidebar) GetWidth() int  { return s.Width }

func (s *Sidebar) SetSidebarViewportHeight(height int) {
	s.viewport.Height = height
}

func (s *Sidebar) SetSidebarViewportWidth(width int) {
	s.viewport.Width = width
}
