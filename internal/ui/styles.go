package ui

import "github.com/charmbracelet/lipgloss"

var (
	ColorBorder   = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	ColorFocus    = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#E5E5E5"}
	ColorText     = lipgloss.AdaptiveColor{Light: "#1F1F1F", Dark: "#F8F8F2"}
	ColorSubtle   = lipgloss.AdaptiveColor{Light: "#A8A8A8", Dark: "#626262"}
	ColorCursorBg = lipgloss.AdaptiveColor{Light: "#E5E5E5", Dark: "#3E3E3E"}

	ColorBarBg = lipgloss.AdaptiveColor{Light: "#F2F2F2", Dark: "#1F1F1F"}
	ColorBarFg = lipgloss.AdaptiveColor{Light: "#6E6E6E", Dark: "#9E9E9E"}

	PaneStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(ColorBorder)

	FocusedPaneStyle = PaneStyle.Copy().
				BorderForeground(ColorFocus)

	DiffStyle = lipgloss.NewStyle().Padding(0, 0)
	ItemStyle = lipgloss.NewStyle().PaddingLeft(2)

	SelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Background(ColorCursorBg).
				Foreground(ColorText).
				Bold(true).
				Width(1000)

	LineNumberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#707070")). // Solid gray, easy to read
			PaddingRight(1).
			Width(4)

	StatusBarStyle     = lipgloss.NewStyle().Foreground(ColorBarFg).Background(ColorBarBg).Padding(0, 1)
	StatusKeyStyle     = lipgloss.NewStyle().Foreground(ColorText).Background(ColorBarBg).Bold(true).Padding(0, 1)
	StatusDividerStyle = lipgloss.NewStyle().Foreground(ColorSubtle).Background(ColorBarBg).Padding(0, 0)
	HelpTextStyle      = lipgloss.NewStyle().Foreground(ColorSubtle).Padding(0, 1)
	HelpDrawerStyle    = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false, false, false).BorderForeground(ColorBorder).Padding(1, 2)
)
