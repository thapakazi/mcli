package tui

import (
	"fmt"
	"mcli/tui/styles"
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
	vp.Style = lipgloss.NewStyle()
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

	title := lipgloss.NewStyle().Bold(true).Foreground(styles.DefaultTheme.TableHeader).Render(event.Title)

	description, _ := glamour.Render(event.Description, "dark")
	if event.Description == "" {
		description, _ = glamour.Render("press R to fetch description", "dark")
	}

	url := lipgloss.NewStyle().Bold(true).Foreground(styles.DefaultTheme.SidebarUrl).Render(event.Url)

	parsedTime, _, _, _ := utils.ParseAndCompareDateTime(event.DateTime)
	date := lipgloss.NewStyle().Bold(true).Foreground(styles.DefaultTheme.SidebarDateTime).Render(parsedTime.String())
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

	s.viewport.SetContent(truncatedContent)

}

func (s *Sidebar) Update(msg tea.Msg) (Sidebar, tea.Cmd) {
	var cmd tea.Cmd
	s.viewport, cmd = s.viewport.Update(msg)

	utils.Logger.Debug("Sidebar / Update", "msg", msg)
	return *s, cmd
}

func (s Sidebar) View() string {

	utils.Logger.Debug("Inside Sidebar View", "visible", s.visible)
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

		// s.viewport.Width = 30
		// s.viewport.Height = 30

	rendered := sidebarStyle.Render(s.viewport.View())
	//utils.Logger.Debug("Sidebar rendered", "content", rendered)
	return rendered
}

func (s Sidebar) GetHeight() int        { return s.Height }
func (s *Sidebar) SetHeight(height int) { s.Height = height }
func (s *Sidebar) SetViewportHeight(height int) {
	s.viewport.Height = height
}

func (s Sidebar) GetWidth() int       { return s.Width }
func (s *Sidebar) SetWidth(width int) { s.Width = width }
func (s *Sidebar) SetViewportWidth(width int) {
	s.viewport.Width = width
}
