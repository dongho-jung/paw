package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/dongho-jung/paw/internal/logging"
	"github.com/dongho-jung/paw/internal/tui"
)

func readImageAttachments(path string) ([]tui.ImageAttachment, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	var attachments []tui.ImageAttachment
	if err := json.Unmarshal(data, &attachments); err != nil {
		return nil, err
	}
	return attachments, nil
}

func writeImageAttachments(path string, attachments []tui.ImageAttachment) error {
	if len(attachments) == 0 {
		return nil
	}
	data, err := json.Marshal(attachments)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func moveImageAttachments(attachments []tui.ImageAttachment, destDir string) ([]tui.ImageAttachment, error) {
	if len(attachments) == 0 {
		return nil, nil
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		logging.Warn("moveImageAttachments: failed to create dir %s: %v", destDir, err)
		return attachments, nil
	}

	moved := make([]tui.ImageAttachment, 0, len(attachments))
	for _, attachment := range attachments {
		base := filepath.Base(attachment.Path)
		destPath := filepath.Join(destDir, base)
		if err := moveFile(attachment.Path, destPath); err != nil {
			logging.Warn("moveImageAttachments: failed to move %s: %v", attachment.Path, err)
			if _, statErr := os.Stat(attachment.Path); statErr == nil {
				moved = append(moved, attachment)
			}
			continue
		}
		if absPath, err := filepath.Abs(destPath); err == nil {
			attachment.Path = absPath
		} else {
			attachment.Path = destPath
		}
		moved = append(moved, attachment)
	}

	return moved, nil
}

func removeImageAttachments(attachments []tui.ImageAttachment) {
	for _, attachment := range attachments {
		if err := os.Remove(attachment.Path); err != nil && !os.IsNotExist(err) {
			logging.Debug("removeImageAttachments: failed to remove %s: %v", attachment.Path, err)
		}
	}
}

func moveFile(src, dest string) error {
	if err := os.Rename(src, dest); err == nil {
		return nil
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	if err := out.Sync(); err != nil {
		_ = out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}

	return os.Remove(src)
}
