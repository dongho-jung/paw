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
		{name: "waiting exact", output: "WAITING", want: task.StatusWaiting, ok: true},
		{name: "done lowercase", output: "done", want: task.StatusDone, ok: true},
		{name: "warning exact", output: "WARNING", want: task.StatusCorrupted, ok: true},
		{name: "warning prefix", output: "warn", want: task.StatusCorrupted, ok: true},
		{name: "contains waiting", output: "Result: WAITING", want: task.StatusWaiting, ok: true},
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
