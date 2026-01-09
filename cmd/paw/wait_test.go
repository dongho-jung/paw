package main

import (
	"strings"
	"testing"
)

func TestDetectWaitInContentAskUserQuestionUI(t *testing.T) {
	content := strings.Join([]string{
		"Which fruit would you like to pick?",
		"> 1. Orange",
		"  A citrus fruit",
		"2. Apple",
		"  A classic fruit",
		"3. Type something.",
		"",
		"Enter to select - Tab/Arrow keys to navigate - Esc to cancel",
	}, "\n")

	waitDetected, reason := detectWaitInContent(content)
	if !waitDetected {
		t.Fatalf("expected wait to be detected for AskUserQuestion UI")
	}
	if reason != "AskUserQuestionUI" {
		t.Fatalf("expected reason AskUserQuestionUI, got %q", reason)
	}
}

func TestDetectWaitInContentMarkerWithoutPrompt(t *testing.T) {
	content := strings.Join([]string{
		"Working on it...",
		"PAW_WAITING",
	}, "\n")

	waitDetected, reason := detectWaitInContent(content)
	if !waitDetected {
		t.Fatalf("expected wait to be detected for PAW_WAITING marker")
	}
	if reason != "marker" {
		t.Fatalf("expected reason marker, got %q", reason)
	}
}

func TestDetectDoneInContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name: "done marker at end",
			content: strings.Join([]string{
				"All tests passed.",
				"PAW_DONE",
				"Ready for review. Press ⌃F to finish.",
			}, "\n"),
			expected: true,
		},
		{
			name: "done marker only",
			content: strings.Join([]string{
				"Task completed.",
				"PAW_DONE",
			}, "\n"),
			expected: true,
		},
		{
			name: "done marker with whitespace",
			content: strings.Join([]string{
				"Task completed.",
				"  PAW_DONE  ",
			}, "\n"),
			expected: true,
		},
		{
			name: "done marker with Claude Code prefix",
			content: strings.Join([]string{
				"Task completed.",
				"⏺ PAW_DONE",
				"Ready for review.",
			}, "\n"),
			expected: true,
		},
		{
			name: "no done marker",
			content: strings.Join([]string{
				"Still working on it...",
				"Running tests...",
			}, "\n"),
			expected: false,
		},
		{
			name: "done marker too far from end",
			content: func() string {
				lines := []string{"PAW_DONE"}
				// Add more than doneMarkerMaxDistance lines after
				for i := 0; i < doneMarkerMaxDistance+10; i++ {
					lines = append(lines, "more output...")
				}
				return strings.Join(lines, "\n")
			}(),
			expected: false,
		},
		{
			name:     "empty content",
			content:  "",
			expected: false,
		},
		{
			name: "partial match should not detect",
			content: strings.Join([]string{
				"PAW_DONE_NOT",
				"still working",
			}, "\n"),
			expected: false,
		},
		{
			name: "done marker in last segment",
			content: strings.Join([]string{
				"⏺ First response",
				"PAW_DONE",
				"Ready for review.",
				"⏺ Second response after new task",
				"Working on it...",
				"PAW_DONE",
				"Done again.",
			}, "\n"),
			expected: true,
		},
		{
			name: "done marker only in previous segment (not last)",
			content: strings.Join([]string{
				"⏺ First response",
				"All done!",
				"PAW_DONE",
				"Ready for review.",
				"⏺ New task started",
				"Working on the new task...",
				"Still processing...",
			}, "\n"),
			expected: false,
		},
		{
			name: "done marker with multiple segments",
			content: strings.Join([]string{
				"⏺ First response",
				"PAW_DONE",
				"⏺ Second response",
				"PAW_DONE",
				"⏺ Third response (new task)",
				"Working on new task...",
			}, "\n"),
			expected: false,
		},
		{
			name: "done marker without any segment markers (backward compat)",
			content: strings.Join([]string{
				"Some output without segment markers",
				"Task completed.",
				"PAW_DONE",
				"Ready for review.",
			}, "\n"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectDoneInContent(tt.content)
			if result != tt.expected {
				t.Errorf("detectDoneInContent() = %v, want %v", result, tt.expected)
			}
		})
	}
}
