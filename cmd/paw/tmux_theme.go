package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/dongho-jung/paw/internal/logging"
	"github.com/dongho-jung/paw/internal/tmux"
)

// ThemePreset represents a tmux color theme preset name.
type ThemePreset string

const (
	// Auto-detection
	ThemeAuto ThemePreset = "auto"

	// Dark themes
	ThemeDark       ThemePreset = "dark"
	ThemeDarkBlue   ThemePreset = "dark-blue"
	ThemeDarkGreen  ThemePreset = "dark-green"
	ThemeDarkPurple ThemePreset = "dark-purple"
	ThemeDarkWarm   ThemePreset = "dark-warm"
	ThemeDarkMono   ThemePreset = "dark-mono"

	// Light themes
	ThemeLight       ThemePreset = "light"
	ThemeLightBlue   ThemePreset = "light-blue"
	ThemeLightGreen  ThemePreset = "light-green"
	ThemeLightPurple ThemePreset = "light-purple"
	ThemeLightWarm   ThemePreset = "light-warm"
	ThemeLightMono   ThemePreset = "light-mono"
)

// ValidThemePresets returns all valid theme preset names.
func ValidThemePresets() []ThemePreset {
	return []ThemePreset{
		ThemeAuto,
		ThemeDark, ThemeDarkBlue, ThemeDarkGreen, ThemeDarkPurple, ThemeDarkWarm, ThemeDarkMono,
		ThemeLight, ThemeLightBlue, ThemeLightGreen, ThemeLightPurple, ThemeLightWarm, ThemeLightMono,
	}
}

// IsValidThemePreset checks if a theme preset name is valid.
func IsValidThemePreset(name string) bool {
	for _, p := range ValidThemePresets() {
		if string(p) == name {
			return true
		}
	}
	return false
}

// tmuxThemeColors holds theme-specific colors for tmux UI elements.
type tmuxThemeColors struct {
	// Status bar
	statusFg string
	statusBg string

	// Window status
	windowFg        string // inactive window text color
	windowBg        string // inactive window background (usually "default")
	windowCurrentFg string // current window text color
	windowCurrentBg string // current window background

	// Pane borders
	paneBorderFg       string // inactive pane border
	paneActiveBorderFg string // active pane border

	// Popup
	popupBorderFg string // popup border color
}

// getThemeColors returns colors for the specified theme preset.
func getThemeColors(preset ThemePreset) tmuxThemeColors {
	switch preset {
	// ==================== DARK THEMES ====================
	case ThemeDark:
		return tmuxThemeColors{
			statusFg:           "colour252", // light gray
			statusBg:           "colour236", // dark gray
			windowFg:           "colour252", // light gray
			windowBg:           "default",
			windowCurrentFg:    "colour231", // white
			windowCurrentBg:    "colour24",  // blue
			paneBorderFg:       "colour238", // dim gray
			paneActiveBorderFg: "colour39",  // bright cyan
			popupBorderFg:      "colour244", // medium gray
		}

	case ThemeDarkBlue:
		return tmuxThemeColors{
			statusFg:           "colour153", // light blue
			statusBg:           "colour17",  // dark navy
			windowFg:           "colour153", // light blue
			windowBg:           "default",
			windowCurrentFg:    "colour231", // white
			windowCurrentBg:    "colour27",  // bright blue
			paneBorderFg:       "colour24",  // dark blue
			paneActiveBorderFg: "colour39",  // bright cyan-blue
			popupBorderFg:      "colour68",  // steel blue
		}

	case ThemeDarkGreen:
		return tmuxThemeColors{
			statusFg:           "colour157", // light green
			statusBg:           "colour22",  // dark green
			windowFg:           "colour157", // light green
			windowBg:           "default",
			windowCurrentFg:    "colour231", // white
			windowCurrentBg:    "colour28",  // green
			paneBorderFg:       "colour22",  // dark green
			paneActiveBorderFg: "colour46",  // bright green
			popupBorderFg:      "colour71",  // sea green
		}

	case ThemeDarkPurple:
		return tmuxThemeColors{
			statusFg:           "colour183", // light purple
			statusBg:           "colour53",  // dark purple
			windowFg:           "colour183", // light purple
			windowBg:           "default",
			windowCurrentFg:    "colour231", // white
			windowCurrentBg:    "colour93",  // purple
			paneBorderFg:       "colour53",  // dark purple
			paneActiveBorderFg: "colour135", // medium purple
			popupBorderFg:      "colour97",  // medium purple
		}

	case ThemeDarkWarm:
		return tmuxThemeColors{
			statusFg:           "colour223", // light orange/cream
			statusBg:           "colour52",  // dark red/brown
			windowFg:           "colour223", // light orange
			windowBg:           "default",
			windowCurrentFg:    "colour231", // white
			windowCurrentBg:    "colour166", // orange
			paneBorderFg:       "colour94",  // brown
			paneActiveBorderFg: "colour208", // bright orange
			popupBorderFg:      "colour137", // tan
		}

	case ThemeDarkMono:
		return tmuxThemeColors{
			statusFg:           "colour250", // light gray
			statusBg:           "colour235", // dark gray
			windowFg:           "colour250", // light gray
			windowBg:           "default",
			windowCurrentFg:    "colour255", // white
			windowCurrentBg:    "colour240", // medium gray
			paneBorderFg:       "colour238", // dim gray
			paneActiveBorderFg: "colour252", // light gray
			popupBorderFg:      "colour244", // medium gray
		}

	// ==================== LIGHT THEMES ====================
	case ThemeLight:
		return tmuxThemeColors{
			statusFg:           "colour236", // dark gray
			statusBg:           "colour253", // light gray
			windowFg:           "colour236", // dark gray
			windowBg:           "default",
			windowCurrentFg:    "colour231", // white
			windowCurrentBg:    "colour25",  // dark blue
			paneBorderFg:       "colour250", // light gray
			paneActiveBorderFg: "colour25",  // dark blue
			popupBorderFg:      "colour245", // medium gray
		}

	case ThemeLightBlue:
		return tmuxThemeColors{
			statusFg:           "colour17",  // dark navy
			statusBg:           "colour153", // light blue
			windowFg:           "colour17",  // dark navy
			windowBg:           "default",
			windowCurrentFg:    "colour231", // white
			windowCurrentBg:    "colour27",  // blue
			paneBorderFg:       "colour117", // light blue
			paneActiveBorderFg: "colour27",  // blue
			popupBorderFg:      "colour68",  // steel blue
		}

	case ThemeLightGreen:
		return tmuxThemeColors{
			statusFg:           "colour22",  // dark green
			statusBg:           "colour157", // light green
			windowFg:           "colour22",  // dark green
			windowBg:           "default",
			windowCurrentFg:    "colour231", // white
			windowCurrentBg:    "colour28",  // green
			paneBorderFg:       "colour158", // pale green
			paneActiveBorderFg: "colour28",  // green
			popupBorderFg:      "colour71",  // sea green
		}

	case ThemeLightPurple:
		return tmuxThemeColors{
			statusFg:           "colour53",  // dark purple
			statusBg:           "colour183", // light purple
			windowFg:           "colour53",  // dark purple
			windowBg:           "default",
			windowCurrentFg:    "colour231", // white
			windowCurrentBg:    "colour93",  // purple
			paneBorderFg:       "colour225", // pale pink
			paneActiveBorderFg: "colour93",  // purple
			popupBorderFg:      "colour97",  // medium purple
		}

	case ThemeLightWarm:
		return tmuxThemeColors{
			statusFg:           "colour94",  // brown
			statusBg:           "colour223", // cream/peach
			windowFg:           "colour94",  // brown
			windowBg:           "default",
			windowCurrentFg:    "colour231", // white
			windowCurrentBg:    "colour166", // orange
			paneBorderFg:       "colour223", // cream
			paneActiveBorderFg: "colour166", // orange
			popupBorderFg:      "colour137", // tan
		}

	case ThemeLightMono:
		return tmuxThemeColors{
			statusFg:           "colour235", // dark gray
			statusBg:           "colour254", // near white
			windowFg:           "colour235", // dark gray
			windowBg:           "default",
			windowCurrentFg:    "colour255", // white
			windowCurrentBg:    "colour240", // medium gray
			paneBorderFg:       "colour252", // light gray
			paneActiveBorderFg: "colour238", // dim gray
			popupBorderFg:      "colour245", // medium gray
		}

	default:
		// Fall back to dark theme
		return getThemeColors(ThemeDark)
	}
}

// applyTmuxTheme applies the specified theme preset to tmux.
func applyTmuxTheme(tm tmux.Client, preset ThemePreset) {
	colors := getThemeColors(preset)

	// Status bar style
	statusStyle := "fg=" + colors.statusFg + ",bg=" + colors.statusBg
	_ = tm.SetOption("status-style", statusStyle, true)

	// Window status format with theme colors
	windowFormat := "#[fg=" + colors.windowFg + "] #W "
	windowCurrentFormat := "#[fg=" + colors.windowCurrentFg + ",bg=" + colors.windowCurrentBg + ",bold] #W "
	_ = tm.SetOption("window-status-format", windowFormat, true)
	_ = tm.SetOption("window-status-current-format", windowCurrentFormat, true)

	// Pane borders
	_ = tm.SetOption("pane-border-style", "fg="+colors.paneBorderFg, true)
	_ = tm.SetOption("pane-active-border-style", "fg="+colors.paneActiveBorderFg+",bold", true)

	// Popup styling
	_ = tm.SetOption("popup-border-style", "fg="+colors.popupBorderFg, true)

	logging.Debug("Applied tmux theme: %s", preset)
}

// detectTerminalTheme detects whether the terminal is in dark or light mode.
// It uses multiple detection methods for reliability:
// 1. COLORFGBG environment variable
// 2. OSC 11 query via lipgloss (with improved timeout handling)
// 3. Falls back to dark mode as default
func detectTerminalTheme() bool {
	// Method 1: Check COLORFGBG environment variable
	// Format: "fg;bg" where bg > 6 typically means light background
	if colorfgbg := os.Getenv("COLORFGBG"); colorfgbg != "" {
		parts := strings.Split(colorfgbg, ";")
		if len(parts) >= 2 {
			if bg, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
				// Background colors: 0-6 are dark, 7-15 and higher are light
				// This is based on standard 16-color terminal palette
				isDark := bg <= 6 || (bg >= 8 && bg <= 14)
				logging.Debug("Theme detection via COLORFGBG=%s: isDark=%v (bg=%d)", colorfgbg, isDark, bg)
				return isDark
			}
		}
	}

	// Method 2: Try OSC 11 query via lipgloss with improved handling
	// Flush stdout and give terminal time to settle
	_ = os.Stdout.Sync()
	time.Sleep(10 * time.Millisecond)

	// Run detection with multiple attempts for reliability
	// In case of timeout or failure, some terminals don't respond to OSC queries
	const attempts = 5
	darkCount := 0
	validCount := 0

	for i := 0; i < attempts; i++ {
		// Use a goroutine with timeout to avoid hanging
		resultCh := make(chan bool, 1)
		go func() {
			resultCh <- lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
		}()

		select {
		case result := <-resultCh:
			validCount++
			if result {
				darkCount++
			}
		case <-time.After(50 * time.Millisecond):
			// Timeout - detection failed for this attempt
			logging.Trace("Theme detection attempt %d timed out", i+1)
		}

		// Small delay between attempts
		if i < attempts-1 {
			time.Sleep(15 * time.Millisecond)
		}
	}

	// If we got at least some valid responses, use majority vote
	if validCount >= 2 {
		isDark := darkCount > validCount/2
		logging.Debug("Theme detection via OSC: isDark=%v (dark=%d/%d)", isDark, darkCount, validCount)
		return isDark
	}

	// Method 3: Check if we're in a tmux session and try to query the terminal
	if os.Getenv("TMUX") != "" {
		// Inside tmux, OSC queries might not work reliably
		// Check if parent terminal type suggests light/dark
		term := os.Getenv("TERM_PROGRAM")
		if strings.Contains(strings.ToLower(term), "apple_terminal") {
			// Apple Terminal default is light
			logging.Debug("Theme detection: Apple Terminal detected, assuming light mode")
			return false
		}
	}

	// Default fallback: assume dark mode (most common for terminal users)
	logging.Debug("Theme detection fallback: assuming dark mode")
	return true
}

// resolveThemePreset resolves the theme preset, auto-detecting if necessary.
// If preset is "auto", it detects the terminal theme and returns the appropriate preset.
func resolveThemePreset(preset ThemePreset) ThemePreset {
	if preset == ThemeAuto || preset == "" {
		if detectTerminalTheme() {
			return ThemeDark
		}
		return ThemeLight
	}
	return preset
}

// IsDarkTheme returns whether the given preset is a dark theme.
func IsDarkTheme(preset ThemePreset) bool {
	switch preset {
	case ThemeDark, ThemeDarkBlue, ThemeDarkGreen, ThemeDarkPurple, ThemeDarkWarm, ThemeDarkMono:
		return true
	case ThemeLight, ThemeLightBlue, ThemeLightGreen, ThemeLightPurple, ThemeLightWarm, ThemeLightMono:
		return false
	case ThemeAuto:
		return detectTerminalTheme()
	default:
		return true // Default to dark
	}
}
