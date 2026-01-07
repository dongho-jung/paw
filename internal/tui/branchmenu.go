package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// BranchAction represents the selected action.
type BranchAction int

const (
	BranchActionCancel BranchAction = iota
	BranchActionMerge               // ↑ Merge to main (default ← task)
	BranchActionSync                // ↓ Sync from main (default → task)
)

// BranchMenu is a simple menu for branch operations.
type BranchMenu struct {
	action BranchAction
}

// NewBranchMenu creates a new branch menu.
func NewBranchMenu() *BranchMenu {
	return &BranchMenu{}
}

// Init initializes the menu.
func (m *BranchMenu) Init() tea.Cmd {
	return nil
}

// Update handles key events.
func (m *BranchMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.action = BranchActionMerge
			return m, tea.Quit
		case "down", "j":
			m.action = BranchActionSync
			return m, tea.Quit
		default:
			// Any other key cancels
			m.action = BranchActionCancel
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the menu.
func (m *BranchMenu) View() tea.View {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	content := fmt.Sprintf(
		"%s\n\n  %s  %s\n  %s  %s\n\n%s",
		titleStyle.Render("Branch Actions"),
		keyStyle.Render("↑"),
		itemStyle.Render("Merge to main (default ← task)"),
		keyStyle.Render("↓"),
		itemStyle.Render("Sync from main (default → task)"),
		dimStyle.Render("Press any other key to cancel"),
	)

	return tea.NewView(content)
}

// Action returns the selected action.
func (m *BranchMenu) Action() BranchAction {
	return m.action
}

// RunBranchMenu runs the branch menu and returns the selected action.
func RunBranchMenu() (BranchAction, error) {
	m := NewBranchMenu()
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return BranchActionCancel, err
	}

	menu := finalModel.(*BranchMenu)
	return menu.Action(), nil
}
