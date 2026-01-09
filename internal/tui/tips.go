// Package tui provides terminal user interface components for PAW.
package tui

import (
	"math/rand/v2"
	"regexp"
)

// Version is the PAW version string, set from main package.
var Version = "dev"

// tips contains usage tips shown to users.
var tips = []string{
	"Press ⌃R to search task history",
	"Press ⌥Tab to switch between panels",
	"Press ⌃T to open template selector",
	"Press ⌃O to view logs",
	"Press ⌃B to toggle bottom shell",
	"Press ⌃/ for help",
	"Press ⌃F to finish a completed task",
	"Press ⌃K to cancel a running task",
	"Press ⌃↓ to sync with main branch",
	"Use mouse to select and copy text",
	"Drag to select text, then ⌃C to copy",
	"Scroll with mouse wheel in kanban",
	"Each task runs in its own git worktree",
	"Tasks can run in parallel without conflicts",
	"Press Esc twice quickly to cancel input",
	"Configure notifications in .paw/config",
	"Use templates (⌃T) for repeated tasks",
	"Task history is saved for easy reuse",
	"Worktree mode keeps main branch clean",
	"Run 'paw check' to verify dependencies",
}

// versionHashRegex matches the git hash suffix in version strings.
// Pattern: -g[0-9a-f]+ or -g[0-9a-f]+-dirty at the end
var versionHashRegex = regexp.MustCompile(`-g[0-9a-f]+(-dirty)?$`)

// GetTip returns a random usage tip.
// Each call returns a different random tip.
func GetTip() string {
	return tips[rand.IntN(len(tips))]
}

// SetVersion sets the PAW version string, stripping the git hash suffix.
// e.g., "v0.3.0-32-gabcdef" -> "v0.3.0-32"
func SetVersion(v string) {
	Version = versionHashRegex.ReplaceAllString(v, "")
}
