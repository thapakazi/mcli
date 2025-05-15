package styles

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF"))
	sourceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	BaseStyle   = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
)

func GetTableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		Foreground(lipgloss.Color("2")).
		Bold(true)
		//Border(lipgloss.NormalBorder(), false, false, true, false) // skipping border bottom, it messing up with layout
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	return s
}
