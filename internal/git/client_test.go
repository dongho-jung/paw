package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	client := New()
	if client == nil {
		t.Fatal("New() returned nil")
	}
}

func TestCopyUntrackedFiles(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	dstDir := filepath.Join(tempDir, "dst")

	// Create source directory structure
	if err := os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755); err != nil {
		t.Fatalf("Failed to create src dir: %v", err)
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		t.Fatalf("Failed to create dst dir: %v", err)
	}

	// Create source files
	files := map[string]string{
		"file1.txt":        "content 1",
		"file2.txt":        "content 2",
		"subdir/file3.txt": "content 3",
	}

	for path, content := range files {
		fullPath := filepath.Join(srcDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir for %s: %v", path, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", path, err)
		}
	}

	// Copy files
	fileList := []string{"file1.txt", "file2.txt", "subdir/file3.txt"}
	if err := CopyUntrackedFiles(fileList, srcDir, dstDir); err != nil {
		t.Fatalf("CopyUntrackedFiles() error = %v", err)
	}

	// Verify copied files
	for path, wantContent := range files {
		dstPath := filepath.Join(dstDir, path)
		data, err := os.ReadFile(dstPath)
		if err != nil {
			t.Errorf("Failed to read copied file %s: %v", path, err)
			continue
		}
		if string(data) != wantContent {
			t.Errorf("Copied file %s content = %q, want %q", path, string(data), wantContent)
		}
	}
}

func TestCopyUntrackedFilesPreservesMode(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	dstDir := filepath.Join(tempDir, "dst")

	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create src dir: %v", err)
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		t.Fatalf("Failed to create dst dir: %v", err)
	}

	// Create executable file
	srcPath := filepath.Join(srcDir, "script.sh")
	if err := os.WriteFile(srcPath, []byte("#!/bin/bash\necho hello"), 0755); err != nil {
		t.Fatalf("Failed to write script: %v", err)
	}

	// Copy file
	if err := CopyUntrackedFiles([]string{"script.sh"}, srcDir, dstDir); err != nil {
		t.Fatalf("CopyUntrackedFiles() error = %v", err)
	}

	// Verify mode
	dstPath := filepath.Join(dstDir, "script.sh")
	info, err := os.Stat(dstPath)
	if err != nil {
		t.Fatalf("Failed to stat copied file: %v", err)
	}

	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		t.Fatalf("Failed to stat source file: %v", err)
	}

	if info.Mode() != srcInfo.Mode() {
		t.Errorf("Copied file mode = %v, want %v", info.Mode(), srcInfo.Mode())
	}
}

func TestCopyUntrackedFilesEmptyList(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	dstDir := filepath.Join(tempDir, "dst")

	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create src dir: %v", err)
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		t.Fatalf("Failed to create dst dir: %v", err)
	}

	// Copy empty list
	if err := CopyUntrackedFiles([]string{}, srcDir, dstDir); err != nil {
		t.Fatalf("CopyUntrackedFiles() with empty list error = %v", err)
	}
}

func TestCopyUntrackedFilesNonexistentSource(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	dstDir := filepath.Join(tempDir, "dst")

	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create src dir: %v", err)
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		t.Fatalf("Failed to create dst dir: %v", err)
	}

	// Try to copy non-existent file
	err := CopyUntrackedFiles([]string{"nonexistent.txt"}, srcDir, dstDir)
	if err == nil {
		t.Error("CopyUntrackedFiles() should return error for non-existent source")
	}
}

func TestCopyUntrackedFilesCreatesDirectories(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	dstDir := filepath.Join(tempDir, "dst")

	// Create source file in nested directory
	nestedDir := filepath.Join(srcDir, "a", "b", "c")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested src dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nestedDir, "file.txt"), []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Create empty dst dir (without nested structure)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		t.Fatalf("Failed to create dst dir: %v", err)
	}

	// Copy file - should create nested directories
	if err := CopyUntrackedFiles([]string{"a/b/c/file.txt"}, srcDir, dstDir); err != nil {
		t.Fatalf("CopyUntrackedFiles() error = %v", err)
	}

	// Verify file exists
	dstPath := filepath.Join(dstDir, "a", "b", "c", "file.txt")
	if _, err := os.Stat(dstPath); err != nil {
		t.Errorf("Copied file not found: %v", err)
	}
}

func TestGenerateMergeCommitMessage(t *testing.T) {
	tests := []struct {
		name            string
		taskName        string
		commits         []CommitInfo
		wantContains    []string
		wantNotContains []string
	}{
		{
			name:     "fix task with no commits",
			taskName: "fix-kanban-drag-select",
			commits:  nil,
			wantContains: []string{
				"fix: kanban drag select",
			},
			wantNotContains: []string{
				"Changes:",
			},
		},
		{
			name:     "fix task with commits",
			taskName: "fix-login-bug",
			commits: []CommitInfo{
				{Hash: "abc123", Subject: "Fix null pointer exception"},
				{Hash: "def456", Subject: "Add error handling"},
			},
			wantContains: []string{
				"fix: login bug",
				"Changes:",
				"- Fix null pointer exception",
				"- Add error handling",
			},
		},
		{
			name:     "feature task",
			taskName: "add-dark-mode",
			commits: []CommitInfo{
				{Hash: "abc123", Subject: "Implement dark mode toggle"},
			},
			wantContains: []string{
				"feat: dark mode",
				"Changes:",
				"- Implement dark mode toggle",
			},
		},
		{
			name:     "refactor task with improve keyword",
			taskName: "improve-commit-messages",
			commits: []CommitInfo{
				{Hash: "abc123", Subject: "Add commit type inference"},
				{Hash: "def456", Subject: "Add commit body generation"},
			},
			wantContains: []string{
				"refactor: improve commit messages",
				"Changes:",
				"- Add commit type inference",
				"- Add commit body generation",
			},
		},
		{
			name:     "docs task",
			taskName: "docs-update-readme",
			commits: []CommitInfo{
				{Hash: "abc123", Subject: "Update installation instructions"},
			},
			wantContains: []string{
				"docs: update readme",
				"Changes:",
				"- Update installation instructions",
			},
		},
		{
			name:     "long commit subject truncated",
			taskName: "fix-bug",
			commits: []CommitInfo{
				{Hash: "abc123", Subject: "This is a very long commit message that exceeds the seventy-two character limit and should be truncated"},
			},
			wantContains: []string{
				"fix: bug",
				"...",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateMergeCommitMessage(tt.taskName, tt.commits)

			for _, want := range tt.wantContains {
				if !contains(result, want) {
					t.Errorf("GenerateMergeCommitMessage() result does not contain %q\nGot:\n%s", want, result)
				}
			}

			for _, notWant := range tt.wantNotContains {
				if contains(result, notWant) {
					t.Errorf("GenerateMergeCommitMessage() result should not contain %q\nGot:\n%s", notWant, result)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
