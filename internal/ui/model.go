package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/oug-t/difi/internal/git"
	"github.com/oug-t/difi/internal/tree"
)

const TargetBranch = "main"

type Focus int

const (
	FocusTree Focus = iota
	FocusDiff
)

type Model struct {
	fileTree     list.Model
	diffViewport viewport.Model

	// Data
	selectedPath  string
	currentBranch string
	repoName      string

	// Diff State
	diffContent string
	diffLines   []string
	diffCursor  int

	// UI State
	focus    Focus
	showHelp bool

	width, height int
}

func NewModel() Model {
	files, _ := git.ListChangedFiles(TargetBranch)
	items := tree.Build(files)

	l := list.New(items, listDelegate{}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	m := Model{
		fileTree:      l,
		diffViewport:  viewport.New(0, 0),
		focus:         FocusTree,
		currentBranch: git.GetCurrentBranch(),
		repoName:      git.GetRepoName(),
		showHelp:      false,
	}

	if len(items) > 0 {
		if first, ok := items[0].(tree.TreeItem); ok {
			m.selectedPath = first.FullPath
		}
	}
	return m
}

func (m Model) Init() tea.Cmd {
	if m.selectedPath != "" {
		return git.DiffCmd(TargetBranch, m.selectedPath)
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()

	case tea.KeyMsg:
		// Toggle Help
		if msg.String() == "?" {
			m.showHelp = !m.showHelp
			m.updateSizes()
			return m, nil
		}

		// Quit
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Navigation
		switch msg.String() {
		case "tab":
			if m.focus == FocusTree {
				m.focus = FocusDiff
			} else {
				m.focus = FocusTree
			}

		case "l", "]", "ctrl+l", "right":
			m.focus = FocusDiff

		case "h", "[", "ctrl+h", "left":
			m.focus = FocusTree

		// Editing
		case "e", "enter":
			if m.selectedPath != "" {
				line := 0
				if m.focus == FocusDiff {
					line = git.CalculateFileLine(m.diffContent, m.diffCursor)
				} else {
					line = git.CalculateFileLine(m.diffContent, 0)
				}
				return m, git.OpenEditorCmd(m.selectedPath, line)
			}

		// Diff Cursor
		case "j", "down":
			if m.focus == FocusDiff {
				if m.diffCursor < len(m.diffLines)-1 {
					m.diffCursor++
					if m.diffCursor >= m.diffViewport.YOffset+m.diffViewport.Height {
						m.diffViewport.LineDown(1)
					}
				}
			}
		case "k", "up":
			if m.focus == FocusDiff {
				if m.diffCursor > 0 {
					m.diffCursor--
					if m.diffCursor < m.diffViewport.YOffset {
						m.diffViewport.LineUp(1)
					}
				}
			}
		}
	}

	// Update Components
	if m.focus == FocusTree {
		m.fileTree, cmd = m.fileTree.Update(msg)
		cmds = append(cmds, cmd)

		if item, ok := m.fileTree.SelectedItem().(tree.TreeItem); ok && !item.IsDir {
			if item.FullPath != m.selectedPath {
				m.selectedPath = item.FullPath
				m.diffCursor = 0
				m.diffViewport.GotoTop()
				cmds = append(cmds, git.DiffCmd(TargetBranch, m.selectedPath))
			}
		}
	}

	switch msg := msg.(type) {
	case git.DiffMsg:
		m.diffContent = msg.Content
		m.diffLines = strings.Split(msg.Content, "\n")
		m.diffViewport.SetContent(msg.Content)

	case git.EditorFinishedMsg:
		return m, git.DiffCmd(TargetBranch, m.selectedPath)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateSizes() {
	reservedHeight := 1
	if m.showHelp {
		reservedHeight += 6
	}

	contentHeight := m.height - reservedHeight
	if contentHeight < 1 {
		contentHeight = 1
	}

	treeWidth := int(float64(m.width) * 0.20)
	if treeWidth < 20 {
		treeWidth = 20
	}

	m.fileTree.SetSize(treeWidth, contentHeight)
	m.diffViewport.Width = m.width - treeWidth - 2
	m.diffViewport.Height = contentHeight
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	treeStyle := PaneStyle
	if m.focus == FocusTree {
		treeStyle = FocusedPaneStyle
	} else {
		treeStyle = PaneStyle
	}

	treeView := treeStyle.Copy().
		Width(m.fileTree.Width()).
		Height(m.fileTree.Height()).
		Render(m.fileTree.View())

	var renderedDiff strings.Builder
	start := m.diffViewport.YOffset
	end := start + m.diffViewport.Height
	if end > len(m.diffLines) {
		end = len(m.diffLines)
	}

	for i := start; i < end; i++ {
		line := m.diffLines[i]
		if m.focus == FocusDiff && i == m.diffCursor {
			line = SelectedItemStyle.Render(line)
		} else {
			line = "  " + line
		}
		renderedDiff.WriteString(line + "\n")
	}

	diffView := DiffStyle.Copy().
		Width(m.diffViewport.Width).
		Height(m.diffViewport.Height).
		Render(renderedDiff.String())

	mainPanes := lipgloss.JoinHorizontal(lipgloss.Top, treeView, diffView)

	// Status Bar
	repoSection := StatusKeyStyle.Render(" " + m.repoName)
	divider := StatusDividerStyle.Render("│")
	branchSection := StatusBarStyle.Render(fmt.Sprintf(" %s ↔ %s", m.currentBranch, TargetBranch))

	leftStatus := lipgloss.JoinHorizontal(lipgloss.Center, repoSection, divider, branchSection)
	rightStatus := StatusBarStyle.Render("? Help")

	statusBar := StatusBarStyle.Copy().
		Width(m.width).
		Render(lipgloss.JoinHorizontal(lipgloss.Top,
			leftStatus,
			lipgloss.PlaceHorizontal(m.width-lipgloss.Width(leftStatus)-lipgloss.Width(rightStatus), lipgloss.Right, rightStatus),
		))

	// Help Drawer
	var finalView string
	if m.showHelp {
		col1 := lipgloss.JoinVertical(lipgloss.Left,
			HelpTextStyle.Render("↑/k   Move Up"),
			HelpTextStyle.Render("↓/j   Move Down"),
		)
		col2 := lipgloss.JoinVertical(lipgloss.Left,
			HelpTextStyle.Render("←/h   Left Panel"),
			HelpTextStyle.Render("→/l   Right Panel"),
		)
		col3 := lipgloss.JoinVertical(lipgloss.Left,
			HelpTextStyle.Render("Tab   Switch Panel"),
			HelpTextStyle.Render("Ent/e Edit File"),
		)
		col4 := lipgloss.JoinVertical(lipgloss.Left,
			HelpTextStyle.Render("q     Quit"),
			HelpTextStyle.Render("?     Close Help"),
		)

		helpDrawer := HelpDrawerStyle.Copy().
			Width(m.width).
			Render(lipgloss.JoinHorizontal(lipgloss.Top,
				col1,
				lipgloss.NewStyle().Width(4).Render(""),
				col2,
				lipgloss.NewStyle().Width(4).Render(""),
				col3,
				lipgloss.NewStyle().Width(4).Render(""),
				col4,
			))

		finalView = lipgloss.JoinVertical(lipgloss.Top, mainPanes, helpDrawer, statusBar)
	} else {
		finalView = lipgloss.JoinVertical(lipgloss.Top, mainPanes, statusBar)
	}

	return finalView
}

type listDelegate struct{}

func (d listDelegate) Height() int                               { return 1 }
func (d listDelegate) Spacing() int                              { return 0 }
func (d listDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d listDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(tree.TreeItem)
	if !ok {
		return
	}
	str := i.Title()
	if index == m.Index() {
		fmt.Fprint(w, SelectedItemStyle.Render(str))
	} else {
		fmt.Fprint(w, ItemStyle.Render(str))
	}
}
