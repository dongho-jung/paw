package github

import (
	"encoding/json"
	"testing"
)

func TestNew(t *testing.T) {
	client := New()
	if client == nil {
		t.Fatal("New() returned nil")
	}
}

func TestPRStatusJSON(t *testing.T) {
	// Test that PRStatus can be properly unmarshaled from JSON
	jsonData := `{"number": 123, "state": "open", "merged": false, "url": "https://github.com/owner/repo/pull/123"}`

	var status PRStatus
	if err := json.Unmarshal([]byte(jsonData), &status); err != nil {
		t.Fatalf("Failed to unmarshal PRStatus: %v", err)
	}

	if status.Number != 123 {
		t.Errorf("Number = %d, want 123", status.Number)
	}
	if status.State != "open" {
		t.Errorf("State = %q, want %q", status.State, "open")
	}
	if status.Merged != false {
		t.Errorf("Merged = %v, want false", status.Merged)
	}
	if status.URL != "https://github.com/owner/repo/pull/123" {
		t.Errorf("URL = %q, want %q", status.URL, "https://github.com/owner/repo/pull/123")
	}
}

func TestPRStatusJSONMerged(t *testing.T) {
	// Test merged state
	jsonData := `{"number": 456, "state": "closed", "merged": true, "url": "https://github.com/owner/repo/pull/456"}`

	var status PRStatus
	if err := json.Unmarshal([]byte(jsonData), &status); err != nil {
		t.Fatalf("Failed to unmarshal PRStatus: %v", err)
	}

	if status.Number != 456 {
		t.Errorf("Number = %d, want 456", status.Number)
	}
	if status.State != "closed" {
		t.Errorf("State = %q, want %q", status.State, "closed")
	}
	if status.Merged != true {
		t.Errorf("Merged = %v, want true", status.Merged)
	}
}

func TestIsInstalled(t *testing.T) {
	client := New()
	// Just test that it doesn't panic - the result depends on the environment
	_ = client.IsInstalled()
}
