package main

import (
	"fmt"

	"github.com/dongho-jung/paw/internal/tmux"
)

// buildKeybindings creates tmux keybindings for PAW.
// Keyboard shortcuts:
//   - Ctrl+N: New task
//   - Ctrl+K: Send Ctrl+C (double-press to cancel task)
//   - Ctrl+F: Finish task
//   - Ctrl+]: Branch menu (↑ merge to main, ↓ sync from main)
//   - Ctrl+Q: Quit paw
//   - Ctrl+T: Toggle tasks
//   - Ctrl+O: Toggle logs
//   - Ctrl+G: Toggle git status
//   - Ctrl+B: Toggle bottom (shell)
//   - Ctrl+/: Toggle help
//   - Alt+Left/Right: Move window
//   - Alt+Tab: Cycle pane (in task windows) / Cycle options (in new task window)
func buildKeybindings(pawBin, sessionName string) []tmux.BindOpts {
	// Command shortcuts
	cmdNewTask := fmt.Sprintf("run-shell '%s internal toggle-new %s'", pawBin, sessionName)
	cmdCtrlC := fmt.Sprintf("run-shell '%s internal ctrl-c %s'", pawBin, sessionName)
	cmdDoneTask := fmt.Sprintf("run-shell '%s internal done-task %s'", pawBin, sessionName)
	cmdBranchMenu := fmt.Sprintf("run-shell '%s internal branch-menu %s'", pawBin, sessionName)
	cmdQuit := "detach-client"
	cmdToggleTasks := fmt.Sprintf("run-shell '%s internal toggle-task-list %s'", pawBin, sessionName)
	cmdToggleLogs := fmt.Sprintf("run-shell '%s internal toggle-log %s'", pawBin, sessionName)
	cmdToggleGitStatus := fmt.Sprintf("run-shell '%s internal toggle-git-status %s'", pawBin, sessionName)
	cmdToggleBottom := fmt.Sprintf("run-shell '%s internal popup-shell %s'", pawBin, sessionName)
	cmdToggleHelp := fmt.Sprintf("run-shell '%s internal toggle-help %s'", pawBin, sessionName)

	// Alt+Tab: context-aware - pass through to TUI in new task window, cycle panes otherwise
	// #{m:pattern,string} checks if string matches pattern (⭐️* = starts with ⭐️)
	// Use "Escape Tab" instead of "M-Tab" because send-keys M-Tab may not produce
	// the correct escape sequence (\x1b\x09) that bubbletea expects for "alt+tab"
	// -F flag is required so tmux evaluates the format as a boolean, not as a shell command
	cmdAltTab := `if -F "#{m:⭐️*,#{window_name}}" "send-keys Escape Tab" "select-pane -t :.+"`

	return []tmux.BindOpts{
		// Navigation (Alt-based)
		{Key: "M-Tab", Command: cmdAltTab, NoPrefix: true},
		{Key: "M-Left", Command: "previous-window", NoPrefix: true},
		{Key: "M-Right", Command: "next-window", NoPrefix: true},

		// Task commands (Ctrl-based)
		{Key: "C-n", Command: cmdNewTask, NoPrefix: true},
		{Key: "C-k", Command: cmdCtrlC, NoPrefix: true},
		{Key: "C-f", Command: cmdDoneTask, NoPrefix: true},
		{Key: "C-]", Command: cmdBranchMenu, NoPrefix: true},
		{Key: "C-q", Command: cmdQuit, NoPrefix: true},

		// Toggle commands (Ctrl-based)
		{Key: "C-t", Command: cmdToggleTasks, NoPrefix: true},
		{Key: "C-o", Command: cmdToggleLogs, NoPrefix: true},
		{Key: "C-g", Command: cmdToggleGitStatus, NoPrefix: true},
		{Key: "C-b", Command: cmdToggleBottom, NoPrefix: true},
		{Key: "C-_", Command: cmdToggleHelp, NoPrefix: true}, // Ctrl+/ sends C-_

	}
}
