package fileutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.txt")

	if err := WriteFileAtomic(path, []byte("first"), 0644); err != nil {
		t.Fatalf("WriteFileAtomic first write failed: %v", err)
	}
	if got, err := os.ReadFile(path); err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	} else if string(got) != "first" {
		t.Fatalf("ReadFile = %q, want %q", string(got), "first")
	}

	if err := WriteFileAtomic(path, []byte("second"), 0644); err != nil {
		t.Fatalf("WriteFileAtomic overwrite failed: %v", err)
	}
	if got, err := os.ReadFile(path); err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	} else if string(got) != "second" {
		t.Fatalf("ReadFile = %q, want %q", string(got), "second")
	}

	if info, err := os.Stat(path); err != nil {
		t.Fatalf("Stat failed: %v", err)
	} else if info.Mode().Perm() != 0644 {
		t.Fatalf("File permissions = %v, want %v", info.Mode().Perm(), os.FileMode(0644))
	}

	if matches, _ := filepath.Glob(filepath.Join(dir, "data.txt.tmp-*")); len(matches) != 0 {
		t.Fatalf("Expected no temp files, found %d", len(matches))
	}
}

func TestWriteFileAtomicMissingDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing", "data.txt")
	if err := WriteFileAtomic(path, []byte("data"), 0644); err == nil {
		t.Fatal("WriteFileAtomic expected error for missing directory")
	}
}

func TestBackupCorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")

	if err := BackupCorruptFile(path); err != nil {
		t.Fatalf("BackupCorruptFile missing file returned error: %v", err)
	}

	if err := os.WriteFile(path, []byte("bad"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	if err := BackupCorruptFile(path); err != nil {
		t.Fatalf("BackupCorruptFile failed: %v", err)
	}
	if _, err := os.Stat(path); err == nil || !os.IsNotExist(err) {
		t.Fatal("Expected original file to be removed")
	}
	backupPath := path + ".corrupt"
	if got, err := os.ReadFile(backupPath); err != nil {
		t.Fatalf("ReadFile backup failed: %v", err)
	} else if string(got) != "bad" {
		t.Fatalf("Backup contents = %q, want %q", string(got), "bad")
	}

	if err := os.WriteFile(path, []byte("new"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	if err := os.WriteFile(backupPath, []byte("old"), 0644); err != nil {
		t.Fatalf("WriteFile backup failed: %v", err)
	}
	if err := BackupCorruptFile(path); err != nil {
		t.Fatalf("BackupCorruptFile with existing backup failed: %v", err)
	}
	if _, err := os.Stat(path); err == nil || !os.IsNotExist(err) {
		t.Fatal("Expected original file to be removed after backup rotation")
	}
	if got, err := os.ReadFile(backupPath); err != nil {
		t.Fatalf("ReadFile existing backup failed: %v", err)
	} else if string(got) != "old" {
		t.Fatalf("Existing backup contents = %q, want %q", string(got), "old")
	}

	matches, _ := filepath.Glob(path + ".*.corrupt")
	if len(matches) == 0 {
		t.Fatal("Expected timestamped backup file")
	}
	foundNew := false
	for _, match := range matches {
		data, err := os.ReadFile(match)
		if err != nil {
			t.Fatalf("ReadFile rotated backup failed: %v", err)
		}
		if string(data) == "new" {
			foundNew = true
			break
		}
	}
	if !foundNew {
		t.Fatal("Expected rotated backup to contain new content")
	}
}
