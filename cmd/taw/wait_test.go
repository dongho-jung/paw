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
		"TAW_WAITING",
	}, "\n")

	waitDetected, reason := detectWaitInContent(content)
	if !waitDetected {
		t.Fatalf("expected wait to be detected for TAW_WAITING marker")
	}
	if reason != "marker" {
		t.Fatalf("expected reason marker, got %q", reason)
	}
}

func TestBuildAppleScriptUsesSystemEvents(t *testing.T) {
	script := buildAppleScript("Pick one?", []string{"Yes", "No"})
	if !strings.Contains(script, "display dialog") {
		t.Fatalf("expected display dialog in AppleScript")
	}
}

func TestBuildWaitDialogScriptUsesSystemEvents(t *testing.T) {
	script := buildWaitDialogScript("Need input")
	if !strings.Contains(script, "display dialog") {
		t.Fatalf("expected display dialog in wait dialog script")
	}
}

func TestBuildWaitDialogJXAScript(t *testing.T) {
	script := buildWaitDialogJXAScript("Need input")
	if !strings.Contains(script, "NSAlert") {
		t.Fatalf("expected NSAlert usage in JXA wait dialog script")
	}
}
