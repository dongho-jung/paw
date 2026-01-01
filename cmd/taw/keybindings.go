package main

import (
	"fmt"

	"github.com/donghojung/taw/internal/tmux"
)

// buildKeybindings creates tmux keybindings for TAW.
// New simplified hotkey scheme:
//   - Alt+Tab: Cycle panes
//   - Alt+Left/Right: Navigate windows
//   - Ctrl+R: Command palette (fzf-based fuzzy search)
//   - Ctrl+C/D twice: Exit session
func buildKeybindings(tawBin, sessionName string) []tmux.BindOpts {
	// Command palette command
	cmdPalette := fmt.Sprintf("run-shell '%s internal command-palette %s'", tawBin, sessionName)

	// Double quit command (Ctrl+C/D twice to exit)
	// The key is passed as an argument to double-quit, which will forward it to the pane
	// if it's not a double-quit. This avoids tmux command chaining issues with ';'.
	cmdDoubleQuitC := fmt.Sprintf("run-shell -b '%s internal double-quit %s C-c'", tawBin, sessionName)
	cmdDoubleQuitD := fmt.Sprintf("run-shell -b '%s internal double-quit %s C-d'", tawBin, sessionName)

	return []tmux.BindOpts{
		// Navigation (Alt-based)
		{Key: "M-Tab", Command: "select-pane -t :.+", NoPrefix: true},
		{Key: "M-Left", Command: "previous-window", NoPrefix: true},
		{Key: "M-Right", Command: "next-window", NoPrefix: true},

		// Command palette (Ctrl+R)
		{Key: "C-r", Command: cmdPalette, NoPrefix: true},

		// Double quit (Ctrl+C/D twice to exit)
		{Key: "C-c", Command: cmdDoubleQuitC, NoPrefix: true},
		{Key: "C-d", Command: cmdDoubleQuitD, NoPrefix: true},
	}
}
