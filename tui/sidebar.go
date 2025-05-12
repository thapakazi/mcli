package tui

import (
	"fmt"
	"mcli/types"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Sidebar struct {
	visible  bool
	viewport viewport.Model
	Width    int
}

func NewSidebar(width int) Sidebar {

	vp := viewport.New(width, 20) // initial height
	vp.Style = lipgloss.NewStyle().
		Background(lipgloss.Color("#222222")).
		Foreground(lipgloss.Color("#FFFFFF"))
	return Sidebar{
		visible:  false,
		viewport: vp,
		Width:    width, // manually setitng 30 for now, TODO make it dynamic
	}
}
func (s *Sidebar) IsVisible() bool {
	return s.visible
}
func (s *Sidebar) ToggleSidebarView() {
	s.visible = !s.visible
}

func (s *Sidebar) UpdateSidebarConntent(event types.Event) {

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")).Render(event.Title)

	description, _ := glamour.Render(event.Description, "dark")
	url := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3")).Render(event.Url)
	date := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")).Render(event.DateTime)
	source := event.Source

	sidebarText := fmt.Sprintf(
		"%s\n\n🔗: %s\n\n＃: %s\n\n📅:%s\n\nDescription:\n-----------\n%s",
		title, url, source, date, description,
	)

	// Style the sidebar content
	contentStyle := lipgloss.NewStyle().
		Width(s.Width - 4). // Account for padding
		Foreground(lipgloss.Color("15"))

	// Wrap the text to fit the viewport width
	renderedContent := contentStyle.Render(sidebarText)

	// Calculate the actual content height (number of lines)
	contentHeight := strings.Count(renderedContent, "\n") + 1

	s.viewport.Height = contentHeight
	s.viewport.SetContent(renderedContent)
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
