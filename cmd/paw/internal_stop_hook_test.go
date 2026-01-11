package main

import (
	"testing"

	"github.com/dongho-jung/paw/internal/task"
)

func TestParseStopHookDecision(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   task.Status
		ok     bool
	}{
		{name: "working exact", output: "WORKING", want: task.StatusWorking, ok: true},
		{name: "working lowercase", output: "working", want: task.StatusWorking, ok: true},
		{name: "waiting exact maps to working", output: "WAITING", want: task.StatusWorking, ok: true}, // WAITING -> WORKING (watch-wait handles it)
		{name: "done lowercase", output: "done", want: task.StatusDone, ok: true},
		{name: "warning exact", output: "WARNING", want: task.StatusWaiting, ok: true}, // WARNING -> WAITING (removed from UI)
		{name: "warning prefix", output: "warn", want: task.StatusWaiting, ok: true},  // WARNING -> WAITING (removed from UI)
		{name: "contains working", output: "Status: WORKING", want: task.StatusWorking, ok: true},
		{name: "contains waiting maps to working", output: "Result: WAITING", want: task.StatusWorking, ok: true}, // WAITING -> WORKING
		{name: "unknown", output: "maybe", want: "", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseStopHookDecision(tt.output)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("status = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHasDoneMarker(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "marker at end",
			content: "Some output\nVerification complete\nPAW_DONE\n",
			want:    true,
		},
		{
			name:    "marker with trailing whitespace",
			content: "Some output\n  PAW_DONE  \n",
			want:    true,
		},
		{
			name:    "marker in middle (within last 20 lines)",
			content: "Line 1\nPAW_DONE\nReady for review\n",
			want:    true,
		},
		{
			name:    "marker with Claude Code prefix",
			content: "Some output\n⏺ PAW_DONE\nReady for review\n",
			want:    true,
		},
		{
			name:    "no marker",
			content: "Some output\nReady for review\n",
			want:    false,
		},
		{
			name:    "partial marker",
			content: "PAW_DONE_EXTRA\n",
			want:    false,
		},
		{
			name:    "marker embedded in text",
			content: "Text PAW_DONE text\n",
			want:    false,
		},
		{
			name:    "empty content",
			content: "",
			want:    false,
		},
		{
			name:    "done marker in last segment",
			content: "⏺ First response\nPAW_DONE\nReady.\n⏺ Second response\nWorking...\nPAW_DONE\nDone again.\n",
			want:    true,
		},
		{
			name:    "done marker only in previous segment (not last)",
			content: "⏺ First response\nAll done!\nPAW_DONE\nReady.\n⏺ New task started\nWorking on the new task...\n",
			want:    false,
		},
		{
			name:    "done marker with multiple segments",
			content: "⏺ First\nPAW_DONE\n⏺ Second\nPAW_DONE\n⏺ Third (new task)\nWorking...\n",
			want:    false,
		},
		{
			name:    "done marker without segment markers (backward compat)",
			content: "Some output\nTask completed.\nPAW_DONE\nReady for review.\n",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasDoneMarker(tt.content)
			if got != tt.want {
				t.Fatalf("hasDoneMarker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasWaitingMarker(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "marker at end",
			content: "Some output\nWorking on it...\nPAW_WAITING\n",
			want:    true,
		},
		{
			name:    "marker with trailing whitespace",
			content: "Some output\n  PAW_WAITING  \n",
			want:    true,
		},
		{
			name:    "marker with UI after (within distance)",
			content: "Some output\nPAW_WAITING\n🔒 Plan\n> 1. Option A\n  Description\n2. Option B\n  Description\nEnter to select\n",
			want:    true,
		},
		{
			name:    "marker with Claude Code prefix",
			content: "Some output\n⏺ PAW_WAITING\nWaiting for input...\n",
			want:    true,
		},
		{
			name:    "no marker",
			content: "Some output\nStill working...\n",
			want:    false,
		},
		{
			name:    "partial marker",
			content: "PAW_WAITING_EXTRA\n",
			want:    false,
		},
		{
			name:    "marker embedded in text",
			content: "Text PAW_WAITING text\n",
			want:    false,
		},
		{
			name:    "empty content",
			content: "",
			want:    false,
		},
		{
			name:    "marker in last segment",
			content: "⏺ First response\nPAW_WAITING\nUI here.\n⏺ Second response\nWorking...\nPAW_WAITING\nMore UI.\n",
			want:    true,
		},
		{
			name:    "marker only in previous segment (not last)",
			content: "⏺ First response\nPAW_WAITING\nUI here.\n⏺ New task started\nWorking on the new task...\n",
			want:    false,
		},
		{
			name:    "waiting after done (new work started)",
			content: "⏺ First response\nPAW_DONE\nReady.\n⏺ New question response\nPAW_WAITING\n> 1. Option\n",
			want:    true,
		},
		{
			name:    "waiting after done without segment marker (real bug scenario)",
			content: "⏺ PAW_DONE\n⏺ New response\n  PAW_WAITING\n☐ Notify\nEnter to select\n",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasWaitingMarker(tt.content)
			if got != tt.want {
				t.Fatalf("hasWaitingMarker() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestWaitingPriorityOverDone tests that PAW_WAITING takes priority over PAW_DONE
// when both markers exist. This is the real bug scenario where:
// 1. Task outputs PAW_DONE (Done state)
// 2. User asks new question
// 3. Agent outputs PAW_WAITING (should become Waiting, not stay Done)
func TestWaitingPriorityOverDone(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		hasWaiting    bool
		hasDone       bool
		hasAskUser    bool
		expectedPrio  string // "waiting" or "done" or "classify"
	}{
		{
			name:         "only done marker",
			content:      "⏺ First response\nAll done!\nPAW_DONE\nReady.\n",
			hasWaiting:   false,
			hasDone:      true,
			hasAskUser:   false,
			expectedPrio: "done",
		},
		{
			name:         "done then waiting (new work)",
			content:      "⏺ First response\nPAW_DONE\n⏺ New response\nPAW_WAITING\nEnter to select\n",
			hasWaiting:   true,
			hasDone:      false, // hasDoneMarker correctly ignores previous segment
			hasAskUser:   false,
			expectedPrio: "waiting",
		},
		{
			name:         "done then AskUserQuestion",
			content:      "⏺ First response\nPAW_DONE\n⏺ New response\nAskUserQuestion:\n  - question: Ready?\n",
			hasWaiting:   false,
			hasDone:      false, // hasDoneMarker correctly ignores previous segment
			hasAskUser:   true,
			expectedPrio: "waiting",
		},
		{
			// Real bug scenario: agent outputs ⏺ PAW_DONE then user asks question
			// Without a new ⏺ marker in the new response
			name:         "done with waiting in same segment (no new segment marker)",
			content:      "⏺ PAW_DONE\nReady.\n\n...user input...\n  PAW_WAITING\n☐ Question\n",
			hasWaiting:   true,
			hasDone:      true, // Both in same segment, both detected
			hasAskUser:   false,
			expectedPrio: "waiting", // But waiting should win
		},
		{
			name:         "no markers",
			content:      "⏺ Response\nWorking on task...\n",
			hasWaiting:   false,
			hasDone:      false,
			hasAskUser:   false,
			expectedPrio: "classify",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWaiting := hasWaitingMarker(tt.content)
			gotDone := hasDoneMarker(tt.content)
			gotAskUser := hasAskUserQuestionInLastSegment(tt.content)

			if gotWaiting != tt.hasWaiting {
				t.Errorf("hasWaitingMarker() = %v, want %v", gotWaiting, tt.hasWaiting)
			}
			if gotDone != tt.hasDone {
				t.Errorf("hasDoneMarker() = %v, want %v", gotDone, tt.hasDone)
			}
			if gotAskUser != tt.hasAskUser {
				t.Errorf("hasAskUserQuestionInLastSegment() = %v, want %v", gotAskUser, tt.hasAskUser)
			}

			// Verify priority logic: waiting/ask > done > classify
			var gotPrio string
			if gotWaiting || gotAskUser {
				gotPrio = "waiting"
			} else if gotDone {
				gotPrio = "done"
			} else {
				gotPrio = "classify"
			}

			if gotPrio != tt.expectedPrio {
				t.Errorf("priority = %v, want %v", gotPrio, tt.expectedPrio)
			}
		})
	}
}

func TestHasAskUserQuestionInLastSegment(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "AskUserQuestion at end of last segment",
			content: "⏺ Working on task\nDoing work...\nAskUserQuestion:\n  - question: How?\n",
			want:    true,
		},
		{
			name:    "AskUserQuestion with options",
			content: "⏺ Response\nAskUserQuestion:\n  questions:\n    - question: Which one?\n      options:\n        - Option A\n        - Option B\n",
			want:    true,
		},
		{
			name:    "AskUserQuestion in previous segment only",
			content: "⏺ First response\nAskUserQuestion:\n  - question: Done?\n⏺ New response\nWorking on changes...\n",
			want:    false,
		},
		{
			name:    "no AskUserQuestion",
			content: "⏺ Response\nAll done!\nPAW_DONE\n",
			want:    false,
		},
		{
			name:    "empty content",
			content: "",
			want:    false,
		},
		{
			name:    "AskUserQuestion without segment marker",
			content: "Working on task...\nAskUserQuestion:\n  - question: Ready?\n",
			want:    true,
		},
		{
			name:    "AskUserQuestion mentioned in text (not tool call)",
			content: "⏺ Response\nI will use AskUserQuestion to ask you\n",
			want:    false,
		},
		{
			name:    "AskUserQuestion tool invocation format",
			content: "⏺ Response\n  AskUserQuestion:\n    questions:\n",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasAskUserQuestionInLastSegment(tt.content)
			if got != tt.want {
				t.Fatalf("hasAskUserQuestionInLastSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}
