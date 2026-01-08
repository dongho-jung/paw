package constants

import (
	"strings"
	"testing"
)

func TestExtractTaskName(t *testing.T) {
	tests := []struct {
		name         string
		windowName   string
		wantTaskName string
		wantFound    bool
	}{
		{
			name:         "working emoji prefix",
			windowName:   "🤖" + WindowToken("my-task"),
			wantTaskName: WindowToken("my-task"),
			wantFound:    true,
		},
		{
			name:         "waiting emoji prefix",
			windowName:   "💬" + WindowToken("my-task"),
			wantTaskName: WindowToken("my-task"),
			wantFound:    true,
		},
		{
			name:         "done emoji prefix",
			windowName:   "✅" + WindowToken("my-task"),
			wantTaskName: WindowToken("my-task"),
			wantFound:    true,
		},
		{
			name:         "warning emoji prefix",
			windowName:   "⚠️" + WindowToken("my-task"),
			wantTaskName: WindowToken("my-task"),
			wantFound:    true,
		},
		{
			name:         "no emoji prefix",
			windowName:   "my-task",
			wantTaskName: "",
			wantFound:    false,
		},
		{
			name:         "different emoji",
			windowName:   "🚀my-task",
			wantTaskName: "",
			wantFound:    false,
		},
		{
			name:         "empty string",
			windowName:   "",
			wantTaskName: "",
			wantFound:    false,
		},
		{
			name:         "emoji only",
			windowName:   "🤖",
			wantTaskName: "",
			wantFound:    true,
		},
		{
			name:         "task with spaces",
			windowName:   "🤖task with spaces",
			wantTaskName: "task with spaces",
			wantFound:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTaskName, gotFound := ExtractTaskName(tt.windowName)
			if gotTaskName != tt.wantTaskName {
				t.Errorf("ExtractTaskName(%q) taskName = %q, want %q", tt.windowName, gotTaskName, tt.wantTaskName)
			}
			if gotFound != tt.wantFound {
				t.Errorf("ExtractTaskName(%q) found = %v, want %v", tt.windowName, gotFound, tt.wantFound)
			}
		})
	}
}

func TestIsTaskWindow(t *testing.T) {
	tests := []struct {
		name       string
		windowName string
		want       bool
	}{
		{
			name:       "working emoji prefix",
			windowName: "🤖my-task",
			want:       true,
		},
		{
			name:       "waiting emoji prefix",
			windowName: "💬my-task",
			want:       true,
		},
		{
			name:       "done emoji prefix",
			windowName: "✅my-task",
			want:       true,
		},
		{
			name:       "warning emoji prefix",
			windowName: "⚠️my-task",
			want:       true,
		},
		{
			name:       "no emoji prefix",
			windowName: "my-task",
			want:       false,
		},
		{
			name:       "new window emoji",
			windowName: "⭐️main",
			want:       false, // EmojiNew is not a task emoji
		},
		{
			name:       "empty string",
			windowName: "",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTaskWindow(tt.windowName); got != tt.want {
				t.Errorf("IsTaskWindow(%q) = %v, want %v", tt.windowName, got, tt.want)
			}
		})
	}
}

func TestTaskEmojis(t *testing.T) {
	expectedEmojis := []string{
		EmojiWorking,
		EmojiWaiting,
		EmojiDone,
		EmojiWarning,
	}

	if len(TaskEmojis) != len(expectedEmojis) {
		t.Errorf("TaskEmojis length = %d, want %d", len(TaskEmojis), len(expectedEmojis))
	}

	for i, emoji := range expectedEmojis {
		if TaskEmojis[i] != emoji {
			t.Errorf("TaskEmojis[%d] = %q, want %q", i, TaskEmojis[i], emoji)
		}
	}
}

func TestConstants(t *testing.T) {
	// Test that constants have expected values
	if PawDirName != ".paw" {
		t.Errorf("PawDirName = %q, want %q", PawDirName, ".paw")
	}
	if AgentsDirName != "agents" {
		t.Errorf("AgentsDirName = %q, want %q", AgentsDirName, "agents")
	}
	if DefaultMainBranch != "main" {
		t.Errorf("DefaultMainBranch = %q, want %q", DefaultMainBranch, "main")
	}
	if NewWindowName != EmojiNew+"main" {
		t.Errorf("NewWindowName = %q, want %q", NewWindowName, EmojiNew+"main")
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "kebab-case to camelCase",
			input:    "cancel-task-twice",
			expected: "cancelTaskTwice",
		},
		{
			name:     "snake_case to camelCase",
			input:    "my_task_name",
			expected: "myTaskName",
		},
		{
			name:     "mixed separators",
			input:    "my-task_name",
			expected: "myTaskName",
		},
		{
			name:     "single word unchanged",
			input:    "task",
			expected: "task",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "already camelCase",
			input:    "myTaskName",
			expected: "myTaskName",
		},
		{
			name:     "consecutive separators",
			input:    "my--task",
			expected: "myTask",
		},
		{
			name:     "separator at start",
			input:    "-my-task",
			expected: "myTask",
		},
		{
			name:     "separator at end",
			input:    "my-task-",
			expected: "myTask",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTruncateWithWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "short name fits",
			input:    "my-task",
			maxLen:   20,
			expected: "myTask", // camelCase conversion
		},
		{
			name:     "long name truncated",
			input:    "cancel-task-twice",
			maxLen:   10,
			expected: "cancelTas…",
		},
		{
			name:     "exact fit",
			input:    "my-task",
			maxLen:   6,
			expected: "myTask",
		},
		{
			name:     "very short width",
			input:    "cancel-task-twice",
			maxLen:   1,
			expected: "…",
		},
		{
			name:     "zero width",
			input:    "cancel-task-twice",
			maxLen:   0,
			expected: "",
		},
		{
			name:     "empty input",
			input:    "",
			maxLen:   10,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateWithWidth(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("TruncateWithWidth(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestTruncateForWindowName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
	}{
		{
			name:  "short name unchanged",
			input: "my-task",
		},
		{
			name:  "exact length unchanged",
			input: "exactly12chr",
		},
		{
			name:  "long name truncated",
			input: "this-is-a-very-long-task-name",
		},
		{
			name:  "empty string unchanged",
			input: "",
		},
		{
			name:  "unicode name truncated by bytes",
			input: "한글태스크",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateForWindowName(tt.input)
			if len(result) > MaxWindowNameLen {
				t.Errorf("TruncateForWindowName(%q) length = %d, want <= %d", tt.input, len(result), MaxWindowNameLen)
			}
			if !strings.Contains(result, WindowTokenSep) {
				t.Errorf("TruncateForWindowName(%q) missing token separator", tt.input)
			}
			suffix := WindowTokenSep + ShortTaskID(tt.input)
			if !strings.HasSuffix(result, suffix) {
				t.Errorf("TruncateForWindowName(%q) = %q, want suffix %q", tt.input, result, suffix)
			}
		})
	}
}

func TestTruncateForWindowNameUsesCamelCase(t *testing.T) {
	// Verify that TruncateForWindowName converts to camelCase
	input := "cancel-task-twice"
	result := TruncateForWindowName(input)

	// The result should contain camelCase base (cancelTaskTwice or truncated)
	// Since MaxWindowNameLen is 20 and suffix is 5 chars (~xxxx), base can be 15 chars
	// "cancelTaskTwice" is 15 chars, so it fits exactly

	if !strings.HasPrefix(result, "cancelTaskTwice~") {
		t.Errorf("TruncateForWindowName(%q) = %q, expected camelCase prefix 'cancelTaskTwice~'", input, result)
	}
}

func TestDisplayLimits(t *testing.T) {
	// Verify limit constants have sensible values
	if MaxDisplayNameLen <= 0 {
		t.Errorf("MaxDisplayNameLen should be positive, got %d", MaxDisplayNameLen)
	}
	if MaxTaskNameLen <= 0 {
		t.Errorf("MaxTaskNameLen should be positive, got %d", MaxTaskNameLen)
	}
	if MinTaskNameLen <= 0 {
		t.Errorf("MinTaskNameLen should be positive, got %d", MinTaskNameLen)
	}
	if MinTaskNameLen > MaxTaskNameLen {
		t.Errorf("MinTaskNameLen (%d) should be <= MaxTaskNameLen (%d)", MinTaskNameLen, MaxTaskNameLen)
	}
	if MaxWindowNameLen <= 0 {
		t.Errorf("MaxWindowNameLen should be positive, got %d", MaxWindowNameLen)
	}
}

func TestTimeoutConstants(t *testing.T) {
	// Verify timeout constants have sensible values
	if ClaudeReadyMaxAttempts <= 0 {
		t.Errorf("ClaudeReadyMaxAttempts should be positive, got %d", ClaudeReadyMaxAttempts)
	}
	if ClaudeReadyPollInterval <= 0 {
		t.Errorf("ClaudeReadyPollInterval should be positive, got %v", ClaudeReadyPollInterval)
	}
	if WorktreeTimeout <= 0 {
		t.Errorf("WorktreeTimeout should be positive, got %v", WorktreeTimeout)
	}
	if TmuxCommandTimeout <= 0 {
		t.Errorf("TmuxCommandTimeout should be positive, got %v", TmuxCommandTimeout)
	}
}

func TestFileAndDirNames(t *testing.T) {
	// Verify file/dir name constants are non-empty
	names := map[string]string{
		"PawDirName":       PawDirName,
		"AgentsDirName":    AgentsDirName,
		"HistoryDirName":   HistoryDirName,
		"ConfigFileName":   ConfigFileName,
		"LogFileName":      LogFileName,
		"MemoryFileName":   MemoryFileName,
		"PromptFileName":   PromptFileName,
		"TaskFileName":     TaskFileName,
		"TabLockDirName":   TabLockDirName,
		"WindowIDFileName": WindowIDFileName,
		"PRFileName":       PRFileName,
		"GitRepoMarker":    GitRepoMarker,
		"GlobalPromptLink": GlobalPromptLink,
		"ClaudeLink":       ClaudeLink,
	}

	for name, value := range names {
		if value == "" {
			t.Errorf("%s should not be empty", name)
		}
	}
}

func TestEmojiConstants(t *testing.T) {
	// Verify emoji constants are non-empty
	emojis := map[string]string{
		"EmojiWorking": EmojiWorking,
		"EmojiWaiting": EmojiWaiting,
		"EmojiDone":    EmojiDone,
		"EmojiWarning": EmojiWarning,
		"EmojiNew":     EmojiNew,
	}

	for name, value := range emojis {
		if value == "" {
			t.Errorf("%s should not be empty", name)
		}
	}
}
