package tui

import "testing"

func TestSetVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "dev version unchanged",
			input:    "dev",
			expected: "dev",
		},
		{
			name:     "tag only unchanged",
			input:    "v0.3.0",
			expected: "v0.3.0",
		},
		{
			name:     "tag with commits unchanged",
			input:    "v0.3.0-32",
			expected: "v0.3.0-32",
		},
		{
			name:     "strip git hash",
			input:    "v0.3.0-32-gabcdef1",
			expected: "v0.3.0-32",
		},
		{
			name:     "strip git hash with dirty",
			input:    "v0.3.0-32-gabcdef1-dirty",
			expected: "v0.3.0-32",
		},
		{
			name:     "strip long git hash",
			input:    "v0.3.0-100-g1234567890abcdef",
			expected: "v0.3.0-100",
		},
		{
			name:     "unknown version unchanged",
			input:    "unknown",
			expected: "unknown",
		},
		{
			name:     "short commit hash only",
			input:    "gabcdef1",
			expected: "gabcdef1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetVersion(tt.input)
			if Version != tt.expected {
				t.Errorf("SetVersion(%q) = %q, want %q", tt.input, Version, tt.expected)
			}
		})
	}
}

func TestGetTip(t *testing.T) {
	// GetTip should return a non-empty string from the tips slice
	tip := GetTip()
	if tip == "" {
		t.Error("GetTip() returned empty string")
	}

	// Verify the tip is in the tips slice
	found := false
	for _, validTip := range tips {
		if tip == validTip {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("GetTip() returned %q which is not in tips slice", tip)
	}
}

func TestGetTipRandomness(t *testing.T) {
	// Call GetTip multiple times and verify we get different tips
	// (with 20 tips, calling 50 times should almost certainly get different results)
	tipCounts := make(map[string]int)
	for i := 0; i < 50; i++ {
		tip := GetTip()
		tipCounts[tip]++
	}

	// We should see at least 2 different tips (very high probability with random selection)
	if len(tipCounts) < 2 {
		t.Errorf("GetTip() returned the same tip all 50 times, expected randomness")
	}
}
