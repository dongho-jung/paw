package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dongho-jung/paw/internal/constants"
	"github.com/dongho-jung/paw/internal/logging"
)

// ImageAttachment represents an image pasted into the task input.
type ImageAttachment struct {
	Label string `json:"label"`
	Path  string `json:"path"`
}

func imagePlaceholder(label string) string {
	return "[" + label + "]"
}

func (m *TaskInput) handleImagePaste() bool {
	img, err := readClipboardImage()
	if err != nil {
		logging.Debug("handleImagePaste: clipboard read failed: %v", err)
	}
	if img == nil || len(img.data) == 0 {
		return false
	}

	pawDir := m.pawDirPath()
	if pawDir == "" {
		logging.Warn("handleImagePaste: pawDir not found")
		return false
	}

	attachmentsDir := filepath.Join(pawDir, constants.InputImageDirName)
	if err := os.MkdirAll(attachmentsDir, 0755); err != nil {
		logging.Warn("handleImagePaste: failed to create attachments dir: %v", err)
		return false
	}

	nextIndex := len(m.imageAttachments) + 1
	label := fmt.Sprintf("Image %d", nextIndex)
	ext := strings.ToLower(img.ext)
	if ext == "" {
		ext = "png"
	}

	now := time.Now().UTC()
	filename := fmt.Sprintf("img-%s-%d-%02d.%s", now.Format("20060102-150405"), now.UnixNano(), nextIndex, ext)
	path := filepath.Join(attachmentsDir, filename)

	if err := os.WriteFile(path, img.data, 0644); err != nil {
		logging.Warn("handleImagePaste: failed to write image: %v", err)
		return false
	}

	m.imageAttachments = append(m.imageAttachments, ImageAttachment{
		Label: label,
		Path:  path,
	})

	m.textarea.InsertString(imagePlaceholder(label))
	m.updateTextareaHeight()

	return true
}

func (m *TaskInput) selectedImageAttachments(content string) []ImageAttachment {
	if len(m.imageAttachments) == 0 {
		return nil
	}

	selected := make([]ImageAttachment, 0, len(m.imageAttachments))
	for _, attachment := range m.imageAttachments {
		if strings.Contains(content, imagePlaceholder(attachment.Label)) {
			selected = append(selected, attachment)
		}
	}

	return selected
}

func (m *TaskInput) cleanupImageAttachments(keep []ImageAttachment) {
	if len(m.imageAttachments) == 0 {
		return
	}

	keepPaths := make(map[string]struct{}, len(keep))
	for _, attachment := range keep {
		keepPaths[attachment.Path] = struct{}{}
	}

	for _, attachment := range m.imageAttachments {
		if _, ok := keepPaths[attachment.Path]; ok {
			continue
		}
		if err := os.Remove(attachment.Path); err != nil && !os.IsNotExist(err) {
			logging.Debug("cleanupImageAttachments: failed to remove %s: %v", attachment.Path, err)
		}
	}
}
