package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oug-t/difi/internal/tree"
)

type TreeDelegate struct {
	Focused bool
}

func (d TreeDelegate) Height() int                               { return 1 }
func (d TreeDelegate) Spacing() int                              { return 0 }
func (d TreeDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d TreeDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(tree.TreeItem)
	if !ok {
		return
	}

	title := i.Title()

	if index == m.Index() {
		style := lipgloss.NewStyle().
			Background(lipgloss.Color("237")). // Dark gray background
			Foreground(lipgloss.Color("255")). // White text
			Bold(true)

		if !d.Focused {
			style = style.Foreground(lipgloss.Color("245"))
		}

		fmt.Fprint(w, style.Render(title))
	} else {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
		fmt.Fprint(w, style.Render(title))
	}
}
