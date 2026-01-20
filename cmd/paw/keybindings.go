package main

import (
	"fmt"

	"github.com/dongho-jung/paw/internal/tmux"
)

// KeybindingsContext contains the context needed for building keybindings.
type KeybindingsContext struct {
	PawBin      string
	SessionName string
	PawDir      string
	ProjectDir  string
	DisplayName string
}

// shellPassthrough wraps a command to pass through the key in shell pane.
// When in shell pane (Ctrl+B popup), the key is sent to shell instead of executing PAW command.
func shellPassthrough(key, pawCmd string) string {
	// #{==:a,b} checks equality, #{@option} gets user option value
	// If current pane is the shell pane, send the key; otherwise run PAW command
	return fmt.Sprintf(`if -F "#{==:#{pane_id},#{@paw_shell_pane_id}}" "send-keys %s" "%s"`, key, pawCmd)
}

// buildKeybindings creates tmux keybindings for PAW.
// Keyboard shortcuts:
//   - Ctrl+N: New task
//   - Ctrl+F: Finish task (shows action picker)
//   - Ctrl+P: Command palette
//   - Ctrl+Q: Quit paw
//   - Ctrl+O: Toggle logs
//   - Ctrl+G: Toggle git viewer (status/log/graph modes)
//   - Ctrl+B: Toggle bottom (shell)
//   - Ctrl+/: Toggle help
//   - Ctrl+R: Toggle history search (in new task window only)
//   - Ctrl+T: Toggle template picker (in new task window only)
//   - Ctrl+J: Toggle project picker (switch between PAW sessions)
//   - Ctrl+Y: Edit prompts (open prompt picker)
//   - Alt+Left/Right: Move window
//   - Alt+Tab: Cycle pane forward (in task windows) / Cycle options (in new task window)
//   - Alt+Shift+Tab: Cycle pane backward (in task windows) / Cycle options backward (in new task window)
//
// Note: Most Ctrl keybindings pass through to shell when in the shell pane (Ctrl+B popup),
// except Ctrl+B (toggle shell), Ctrl+Q (quit), and Ctrl+F (finish task).
func buildKeybindings(ctx KeybindingsContext) []tmux.BindOpts {
	// Environment variables for proper context resolution in subdirectory sessions
	// These are embedded directly in the keybindings (not tmux format variables)
	envPrefix := fmt.Sprintf(`PAW_DIR="%s" PROJECT_DIR="%s" DISPLAY_NAME="%s" `,
		ctx.PawDir, ctx.ProjectDir, ctx.DisplayName)

	// Command shortcuts - all commands include env vars for proper context resolution
	cmdNewTask := fmt.Sprintf("run-shell '%s%s internal toggle-new %s'", envPrefix, ctx.PawBin, ctx.SessionName)
	cmdDoneTask := fmt.Sprintf("run-shell '%s%s internal done-task %s'", envPrefix, ctx.PawBin, ctx.SessionName)
	cmdQuit := "detach-client"
	cmdToggleLogs := fmt.Sprintf("run-shell '%s%s internal toggle-log %s'", envPrefix, ctx.PawBin, ctx.SessionName)
	cmdToggleGitStatus := fmt.Sprintf("run-shell '%s%s internal toggle-git-status %s'", envPrefix, ctx.PawBin, ctx.SessionName)
	cmdToggleBottom := fmt.Sprintf("run-shell '%s%s internal popup-shell %s'", envPrefix, ctx.PawBin, ctx.SessionName)
	cmdToggleHelp := fmt.Sprintf("run-shell '%s%s internal toggle-help %s'", envPrefix, ctx.PawBin, ctx.SessionName)
	cmdToggleCmdPalette := fmt.Sprintf("run-shell '%s%s internal toggle-cmd-palette %s'", envPrefix, ctx.PawBin, ctx.SessionName)
	cmdToggleProjectPicker := fmt.Sprintf("run-shell '%s%s internal toggle-project-picker %s'", envPrefix, ctx.PawBin, ctx.SessionName)
	cmdTogglePromptPicker := fmt.Sprintf("run-shell '%s%s internal toggle-prompt-picker %s'", envPrefix, ctx.PawBin, ctx.SessionName)

	// Alt+Tab: context-aware - pass through to TUI in new task window, cycle panes otherwise
	// #{m:pattern,string} checks if string matches pattern (⭐️* = starts with ⭐️)
	// Use "Escape Tab" instead of "M-Tab" because send-keys M-Tab may not produce
	// the correct escape sequence (\x1b\x09) that bubbletea expects for "alt+tab"
	// -F flag is required so tmux evaluates the format as a boolean, not as a shell command
	cmdAltTab := `if -F "#{m:⭐️*,#{window_name}}" "send-keys Escape Tab" "select-pane -t :.+"`
	cmdAltShiftTab := `if -F "#{m:⭐️*,#{window_name}}" "send-keys Escape BTab" "select-pane -t :.-"`

	// Ctrl+R: context-aware with shell passthrough
	// Priority: shell pane > new task window > pass through
	cmdCtrlRBase := fmt.Sprintf(`if -F "#{m:⭐️*,#{window_name}}" "run-shell '%s%s internal toggle-history %s'" "send-keys C-r"`, envPrefix, ctx.PawBin, ctx.SessionName)
	cmdCtrlR := shellPassthrough("C-r", cmdCtrlRBase)

	// Ctrl+T: context-aware with shell passthrough
	cmdCtrlTBase := fmt.Sprintf(`if -F "#{m:⭐️*,#{window_name}}" "run-shell '%s%s internal toggle-template %s'" "send-keys C-t"`, envPrefix, ctx.PawBin, ctx.SessionName)
	cmdCtrlT := shellPassthrough("C-t", cmdCtrlTBase)

	return []tmux.BindOpts{
		// Navigation (Alt-based)
		{Key: "M-Tab", Command: cmdAltTab, NoPrefix: true},
		{Key: "M-BTab", Command: cmdAltShiftTab, NoPrefix: true},
		{Key: "M-Left", Command: "previous-window", NoPrefix: true},
		{Key: "M-Right", Command: "next-window", NoPrefix: true},

		// Task commands (Ctrl-based)
		// These pass through to shell in shell pane, except Ctrl+F and Ctrl+Q
		{Key: "C-n", Command: shellPassthrough("C-n", cmdNewTask), NoPrefix: true},
		{Key: "C-f", Command: cmdDoneTask, NoPrefix: true}, // Always works (finish task)
		{Key: "C-p", Command: shellPassthrough("C-p", cmdToggleCmdPalette), NoPrefix: true},
		{Key: "C-q", Command: cmdQuit, NoPrefix: true}, // Always works (quit)

		// Toggle commands (Ctrl-based)
		// These pass through to shell in shell pane, except Ctrl+B
		{Key: "C-o", Command: shellPassthrough("C-o", cmdToggleLogs), NoPrefix: true},
		{Key: "C-g", Command: shellPassthrough("C-g", cmdToggleGitStatus), NoPrefix: true},
		{Key: "C-b", Command: cmdToggleBottom, NoPrefix: true},                             // Always works (toggle shell)
		{Key: "C-_", Command: shellPassthrough("C-_", cmdToggleHelp), NoPrefix: true},      // Ctrl+/ sends C-_
		{Key: "C-r", Command: cmdCtrlR, NoPrefix: true},                                    // History search in new task window
		{Key: "C-t", Command: cmdCtrlT, NoPrefix: true},                                    // Template picker in new task window
		{Key: "C-j", Command: shellPassthrough("C-j", cmdToggleProjectPicker), NoPrefix: true}, // Project picker
		{Key: "C-y", Command: shellPassthrough("C-y", cmdTogglePromptPicker), NoPrefix: true},  // Prompt editor
	}
}
