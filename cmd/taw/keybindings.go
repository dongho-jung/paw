package main

import (
	"fmt"

	"github.com/donghojung/taw/internal/tmux"
)

// keyMapping defines a hotkey with its English and Korean equivalents
type keyMapping struct {
	english string // English key (e.g., "n")
	korean  string // Korean 2-벌식 key (e.g., "ㅜ")
}

// Korean 2-벌식 keyboard layout mapping
var koreanKeyMap = map[string]string{
	"n": "ㅜ",
	"t": "ㅅ",
	"e": "ㄷ",
	"m": "ㅡ",
	"p": "ㅔ",
	"u": "ㅕ",
	"l": "ㅣ",
	"q": "ㅂ",
}

// hotkeyDef defines a hotkey action
type hotkeyDef struct {
	key     string // English key
	command string
}

// buildKeybindings creates tmux keybindings for both English and Korean layouts
func buildKeybindings(tawBin, sessionName string) []tmux.BindOpts {
	// Command templates
	cmdToggleNew := fmt.Sprintf("run-shell '%s internal toggle-new %s'", tawBin, sessionName)
	cmdToggleTaskList := fmt.Sprintf("run-shell '%s internal toggle-task-list %s'", tawBin, sessionName)
	cmdEndTaskUI := fmt.Sprintf("run-shell '%s internal end-task-ui %s #{window_id}'", tawBin, sessionName)
	cmdMergeCompleted := fmt.Sprintf("run-shell '%s internal merge-completed %s'", tawBin, sessionName)
	cmdPopupShell := fmt.Sprintf("run-shell '%s internal popup-shell %s'", tawBin, sessionName)
	cmdQuickTask := fmt.Sprintf("run-shell '%s internal quick-task %s'", tawBin, sessionName)
	cmdToggleLog := fmt.Sprintf("run-shell '%s internal toggle-log %s'", tawBin, sessionName)
	cmdToggleHelp := fmt.Sprintf("run-shell '%s internal toggle-help %s'", tawBin, sessionName)

	// Hotkey definitions (English key -> command)
	hotkeys := []hotkeyDef{
		{"n", cmdToggleNew},
		{"t", cmdToggleTaskList},
		{"e", cmdEndTaskUI},
		{"m", cmdMergeCompleted},
		{"p", cmdPopupShell},
		{"u", cmdQuickTask},
		{"l", cmdToggleLog},
		{"/", cmdToggleHelp},
		{"q", "detach"},
	}

	// Build bindings
	var bindings []tmux.BindOpts

	// Navigation keys (language-independent)
	bindings = append(bindings,
		tmux.BindOpts{Key: "M-Tab", Command: "select-pane -t :.+", NoPrefix: true},
		tmux.BindOpts{Key: "M-Left", Command: "previous-window", NoPrefix: true},
		tmux.BindOpts{Key: "M-Right", Command: "next-window", NoPrefix: true},
	)

	// Add English and Korean bindings for each hotkey
	for _, hk := range hotkeys {
		// English binding
		bindings = append(bindings, tmux.BindOpts{
			Key:      "M-" + hk.key,
			Command:  hk.command,
			NoPrefix: true,
		})

		// Korean binding (if mapping exists)
		if korean, ok := koreanKeyMap[hk.key]; ok {
			bindings = append(bindings, tmux.BindOpts{
				Key:      "M-" + korean,
				Command:  hk.command,
				NoPrefix: true,
			})
		}
	}

	return bindings
}
