// Package constants defines shared constants used throughout the PAW application.
package constants

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"time"
)

// Window status emojis
const (
	EmojiWorking = "🤖"
	EmojiWaiting = "💬"
	EmojiDone    = "✅"
	EmojiWarning = "⚠️"
	EmojiNew     = "⭐️"
)

// TaskEmojis contains all emojis used for task windows.
var TaskEmojis = []string{
	EmojiWorking,
	EmojiWaiting,
	EmojiDone,
	EmojiWarning,
}

// IsTaskWindow returns true if the window name has a task emoji prefix.
func IsTaskWindow(windowName string) bool {
	for _, emoji := range TaskEmojis {
		if strings.HasPrefix(windowName, emoji) {
			return true
		}
	}
	return false
}

// ExtractTaskName extracts the window token from a window name by removing the emoji prefix.
// Returns the token and true if a task emoji was found, or empty string and false otherwise.
func ExtractTaskName(windowName string) (string, bool) {
	for _, emoji := range TaskEmojis {
		if strings.HasPrefix(windowName, emoji) {
			return strings.TrimPrefix(windowName, emoji), true
		}
	}
	return "", false
}

const (
	WindowTokenSep = "~"
	WindowIDLen    = 4
)

// TruncateForWindowName returns a stable window token for a task name.
// The token includes a short ID suffix to avoid collisions.
func TruncateForWindowName(name string) string {
	return WindowToken(name)
}

// LegacyTruncateForWindowName truncates a task name without an ID suffix.
// This preserves backward compatibility with older window names.
func LegacyTruncateForWindowName(name string) string {
	if len(name) > MaxWindowNameLen {
		return name[:MaxWindowNameLen]
	}
	return name
}

// WindowToken builds a window-safe token with a short stable ID suffix.
func WindowToken(name string) string {
	id := ShortTaskID(name)
	suffix := WindowTokenSep + id
	maxBase := MaxWindowNameLen - len(suffix)
	if maxBase < 1 {
		maxBase = 1
	}
	base := name
	if len(base) > maxBase {
		base = base[:maxBase]
	}
	return base + suffix
}

// ShortTaskID returns a stable short ID for a task name.
func ShortTaskID(name string) string {
	sum := sha1.Sum([]byte(name))
	return hex.EncodeToString(sum[:])[:WindowIDLen]
}

// MatchesWindowToken returns true if the extracted window token matches the task name.
func MatchesWindowToken(extracted, taskName string) bool {
	return extracted == TruncateForWindowName(taskName) || extracted == LegacyTruncateForWindowName(taskName)
}

// Display limits
const (
	MaxDisplayNameLen = 32
	MaxTaskNameLen    = 32
	MinTaskNameLen    = 8
	MaxWindowNameLen  = 12 // Max task name length in tmux window names
)

// Claude interaction timeouts
const (
	ClaudeReadyMaxAttempts  = 60
	ClaudeReadyPollInterval = 500 * time.Millisecond
	ClaudeNameGenTimeout1   = 1 * time.Minute // haiku
	ClaudeNameGenTimeout2   = 2 * time.Minute // sonnet
	ClaudeNameGenTimeout3   = 3 * time.Minute // opus
	ClaudeNameGenTimeout4   = 4 * time.Minute // opus with thinking
)

// Git/Worktree timeouts
const (
	WorktreeTimeout       = 30 * time.Second
	WindowCreationTimeout = 30 * time.Second
)

// Tmux command timeout
const (
	TmuxCommandTimeout = 10 * time.Second
)

// Default configuration values
const (
	DefaultMainBranch = "main"
	DefaultWorkMode   = "worktree"
	DefaultOnComplete = "confirm"
)

// Directory and file names
const (
	PawDirName       = ".paw"
	AgentsDirName    = "agents"
	HistoryDirName   = "history"
	WindowMapFileName = "window-map.json"
	ConfigFileName   = "config"
	LogFileName      = "log"
	MemoryFileName   = "memory"
	PromptFileName   = "PROMPT.md"
	TaskFileName     = "task"
	TabLockDirName   = ".tab-lock"
	WindowIDFileName = "window_id"
	PRFileName       = ".pr"
	GitRepoMarker    = ".is-git-repo"
	GlobalPromptLink = ".global-prompt"
	ClaudeLink       = ".claude"
)

// Tmux related constants
const (
	TmuxSocketPrefix = "paw-"
	NewWindowName    = EmojiNew + "main"
)

// Pane capture settings
const (
	PaneCaptureLines = 10000 // Number of lines to capture from pane history
	SummaryMaxLen    = 8000  // Max characters to send for summary generation
)

// Merge lock settings
const (
	MergeLockMaxRetries    = 30              // Maximum retries to acquire merge lock
	MergeLockRetryInterval = 1 * time.Second // Interval between lock retries
)

// Commit message templates
const (
	CommitMessageMerge           = "feat: %s" // Format string for merge commits
	CommitMessageAutoCommit      = "chore: auto-commit on task end\n\n%s"
	CommitMessageAutoCommitMerge = "chore: auto-commit before merge\n\n%s"
	CommitMessageAutoCommitPush  = "chore: auto-commit before push"
)

// Double-press detection
const (
	DoublePressIntervalSec = 2 // Seconds to wait for second keypress
)

// Task window handling
const (
	WindowIDWaitMaxAttempts = 60                     // Max attempts to wait for window ID file
	WindowIDWaitInterval    = 500 * time.Millisecond // Interval between checks
)
