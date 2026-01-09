// Package tui provides terminal user interface components for PAW.
package tui

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss/v2"

	"github.com/dongho-jung/paw/internal/config"
	"github.com/dongho-jung/paw/internal/constants"
)

// DetectDarkMode returns whether the terminal is in dark mode.
// It checks the theme config setting first:
//   - "light": always returns false (dark mode = off)
//   - "dark": always returns true (dark mode = on)
//   - "auto" or empty: uses lipgloss.HasDarkBackground() to auto-detect
//
// This function should be called BEFORE bubbletea starts, as
// lipgloss.HasDarkBackground() reads from stdin.
func DetectDarkMode() bool {
	// Try to load theme from config
	theme := loadThemeFromConfig()

	switch theme {
	case config.ThemeLight:
		return false
	case config.ThemeDark:
		return true
	default:
		// Auto-detect
		return lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	}
}

// loadThemeFromConfig attempts to load the theme setting from .paw/config.
// Returns ThemeAuto if the config cannot be loaded.
func loadThemeFromConfig() config.Theme {
	// Find .paw directory
	pawDir := findPawDir()
	if pawDir == "" {
		return config.ThemeAuto
	}

	cfg, err := config.Load(pawDir)
	if err != nil {
		return config.ThemeAuto
	}

	if cfg.Theme == "" {
		return config.ThemeAuto
	}

	return cfg.Theme
}

// findPawDir looks for .paw directory starting from current dir up to root.
func findPawDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		pawDir := filepath.Join(dir, constants.PawDirName)
		if info, err := os.Stat(pawDir); err == nil && info.IsDir() {
			return pawDir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}
