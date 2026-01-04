// Package tui provides terminal user interface components for PAW.
package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// HelpViewer provides an interactive help viewer with vim-like navigation.
type HelpViewer struct {
	lines         []string
	scrollPos     int
	horizontalPos int
	width         int
	height        int
}

// NewHelpViewer creates a new help viewer with the given content.
func NewHelpViewer(content string) *HelpViewer {
	lines := strings.Split(content, "\n")
	// Remove last empty line if present
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return &HelpViewer{
		lines: lines,
	}
}

// Init initializes the help viewer.
func (m *HelpViewer) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model.
func (m *HelpViewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// handleKey handles keyboard input.
func (m *HelpViewer) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	// Close on q, Esc, or Ctrl+/ (which is ctrl+_)
	case "q", "esc", "ctrl+_", "ctrl+/", "ctrl+shift+/":
		return m, tea.Quit

	case "down", "j":
		m.scrollDown(1)

	case "up", "k":
		m.scrollUp(1)

	case "left", "h":
		if m.horizontalPos > 0 {
			m.horizontalPos -= 10
			if m.horizontalPos < 0 {
				m.horizontalPos = 0
			}
		}

	case "right", "l":
		m.horizontalPos += 10

	case "g", "home":
		m.scrollPos = 0
		m.horizontalPos = 0

	case "G", "end":
		m.scrollToEnd()
		m.horizontalPos = 0

	case "pgup", "ctrl+b":
		m.scrollUp(m.contentHeight())

	case "pgdown", "ctrl+f", " ":
		m.scrollDown(m.contentHeight())

	case "ctrl+u":
		m.scrollUp(m.contentHeight() / 2)

	case "ctrl+d":
		m.scrollDown(m.contentHeight() / 2)
	}

	return m, nil
}

// scrollUp scrolls up by n lines.
func (m *HelpViewer) scrollUp(n int) {
	m.scrollPos -= n
	if m.scrollPos < 0 {
		m.scrollPos = 0
	}
}

// scrollDown scrolls down by n lines.
func (m *HelpViewer) scrollDown(n int) {
	max := len(m.lines) - m.contentHeight()
	if max < 0 {
		max = 0
	}
	m.scrollPos += n
	if m.scrollPos > max {
		m.scrollPos = max
	}
}

// View renders the help viewer.
func (m *HelpViewer) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Loading...")
	}

	var sb strings.Builder

	// Calculate visible lines
	contentHeight := m.contentHeight()
	endPos := m.scrollPos + contentHeight
	if endPos > len(m.lines) {
		endPos = len(m.lines)
	}

	// Render visible lines
	for i := m.scrollPos; i < endPos; i++ {
		line := m.lines[i]

		// Apply horizontal scroll
		if m.horizontalPos < len(line) {
			line = line[m.horizontalPos:]
		} else {
			line = ""
		}

		// Truncate to screen width
		if len(line) > m.width {
			line = line[:m.width]
		}

		// Pad to full width
		line = fmt.Sprintf("%-*s", m.width, line)
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	// Pad remaining lines
	for i := endPos - m.scrollPos; i < contentHeight; i++ {
		sb.WriteString(strings.Repeat(" ", m.width))
		sb.WriteString("\n")
	}

	// Status bar
	statusStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("240")).
		Foreground(lipgloss.Color("252"))

	var status string
	if len(m.lines) > 0 {
		status = fmt.Sprintf(" Lines %d-%d of %d ", m.scrollPos+1, endPos, len(m.lines))
	} else {
		status = " (empty) "
	}

	// Keybindings hint
	hint := "↑↓j/k:scroll g/G:top/end ⌃/:close"
	padding := m.width - len(status) - len(hint)
	if padding < 0 {
		hint = "⌃/:close"
		padding = m.width - len(status) - len(hint)
		if padding < 0 {
			padding = 0
		}
	}

	statusLine := statusStyle.Render(
		status + strings.Repeat(" ", padding) + hint,
	)

	sb.WriteString(statusLine)

	v := tea.NewView(sb.String())
	v.AltScreen = true
	return v
}

// contentHeight returns the height available for content.
func (m *HelpViewer) contentHeight() int {
	// Reserve 1 line for status bar
	h := m.height - 1
	if h < 1 {
		h = 1
	}
	return h
}

// scrollToEnd scrolls to the end of the content.
func (m *HelpViewer) scrollToEnd() {
	max := len(m.lines) - m.contentHeight()
	if max < 0 {
		max = 0
	}
	m.scrollPos = max
}

// RunHelpViewer runs the help viewer with the given content.
func RunHelpViewer(content string) error {
	m := NewHelpViewer(content)
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}
