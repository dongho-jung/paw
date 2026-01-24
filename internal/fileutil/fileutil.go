// Package fileutil provides helpers for working with files safely.
package fileutil

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// WriteFileAtomic writes data to path by renaming a temp file into place.
func WriteFileAtomic(path string, data []byte, perm fs.FileMode) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	tmpFile, err := os.CreateTemp(dir, base+".tmp-*")
	if err != nil {
		return err
	}

	tmpName := tmpFile.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpName)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Chmod(perm); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpName, path); err != nil {
		if errors.Is(err, fs.ErrExist) || os.IsExist(err) {
			if removeErr := os.Remove(path); removeErr != nil {
				return err
			}
			if err := os.Rename(tmpName, path); err != nil {
				return err
			}
			cleanup = false
			return nil
		}
		return err
	}

	cleanup = false
	return nil
}

// BackupCorruptFile renames path to a .corrupt backup when possible.
func BackupCorruptFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	backupPath := path + ".corrupt"
	if _, err := os.Stat(backupPath); err == nil {
		backupPath = fmt.Sprintf("%s.%d.corrupt", path, time.Now().UnixNano())
	}

	return os.Rename(path, backupPath)
}
