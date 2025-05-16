package styles

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	BaseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.HiddenBorder()).BorderForeground(lipgloss.Color("240"))
)

type Theme struct {
	TableHeader lipgloss.AdaptiveColor
	FaintBorder lipgloss.AdaptiveColor
	TableRows   lipgloss.AdaptiveColor

	TableRowSelectedForeground lipgloss.AdaptiveColor
	TableRowSelectedBackground lipgloss.AdaptiveColor

	SidebarTitle    lipgloss.AdaptiveColor
	SidebarUrl      lipgloss.AdaptiveColor
	SidebarLocation lipgloss.AdaptiveColor
	SidebarDateTime lipgloss.AdaptiveColor
	// SidebarDescirption glamour.Render

	StatusBackground lipgloss.AdaptiveColor
	StatusForeground lipgloss.AdaptiveColor
}

var DefaultTheme = &Theme{

	TableHeader:                lipgloss.AdaptiveColor{Light: "", Dark: "2"},
	TableRowSelectedForeground: lipgloss.AdaptiveColor{Light: "", Dark: "229"},
	TableRowSelectedBackground: lipgloss.AdaptiveColor{Light: "", Dark: "57"},
	FaintBorder:                lipgloss.AdaptiveColor{Light: "", Dark: "240"},
	TableRows:                  lipgloss.AdaptiveColor{Light: "", Dark: "240"},

	SidebarTitle:    lipgloss.AdaptiveColor{Light: "", Dark: "2"},
	SidebarUrl:      lipgloss.AdaptiveColor{Light: "", Dark: "3"},
	SidebarLocation: lipgloss.AdaptiveColor{Light: "", Dark: "4"},
	SidebarDateTime: lipgloss.AdaptiveColor{Light: "", Dark: "5"},
	//SidebarDescirption:

	StatusBackground: lipgloss.AdaptiveColor{Light: "240", Dark: "#340"},
	StatusForeground: lipgloss.AdaptiveColor{Light: "#3C3C3c", Dark: "240"},
}

func GetTableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		Foreground(DefaultTheme.TableHeader).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderBottomForeground(DefaultTheme.FaintBorder)

	s.Selected = s.Selected.
		Foreground(DefaultTheme.TableRowSelectedForeground).
		Background(DefaultTheme.TableRowSelectedBackground).
		Bold(false)
	return s
}
