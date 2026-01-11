// Package notify provides simple desktop notifications.
package notify

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/dongho-jung/paw/internal/logging"
)

// SoundType represents different notification sounds.
type SoundType string

const (
	// SoundTaskCreated is played when a task window is created.
	SoundTaskCreated SoundType = "Glass"
	// SoundTaskCompleted is played when a task is completed successfully.
	SoundTaskCompleted SoundType = "Hero"
	// SoundNeedInput is played when user intervention is needed.
	SoundNeedInput SoundType = "Funk"
	// SoundError is played when an error or problem occurs.
	SoundError SoundType = "Basso"
	// SoundCancelPending is played when waiting for second Ctrl+C to cancel.
	SoundCancelPending SoundType = "Tink"
)

// Send shows a desktop notification using AppleScript (macOS only).
func Send(title, message string) error {
	logging.Debug("-> Send(title=%q, message=%q)", title, message)
	defer logging.Debug("<- Send")

	if runtime.GOOS != "darwin" {
		return nil
	}

	script := fmt.Sprintf(`display notification %q with title %q`, message, title)
	cmd := appleScriptCommand("-e", script)
	if err := cmd.Run(); err != nil {
		fallbackErr := exec.Command("osascript", "-e", script).Run()
		if fallbackErr == nil {
			return nil
		}
		return err
	}
	return nil
}

// PlaySound plays a system sound (macOS only).
// It runs in the background and does not block.
// Uses nohup to ensure the sound plays even if parent process exits.
func PlaySound(soundType SoundType) {
	logging.Debug("-> PlaySound(soundType=%s)", soundType)
	defer logging.Debug("<- PlaySound")

	if runtime.GOOS != "darwin" {
		return
	}

	soundPath := fmt.Sprintf("/System/Library/Sounds/%s.aiff", soundType)

	// Check if sound file exists
	if _, err := os.Stat(soundPath); os.IsNotExist(err) {
		logging.Trace("PlaySound: sound file not found path=%s", soundPath)
		return
	}

	// Run afplay via nohup to ensure it survives parent process exit
	// Redirect output to /dev/null to fully detach
	cmd := exec.Command("nohup", "afplay", soundPath)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	if err := cmd.Start(); err != nil {
		logging.Warn("PlaySound: failed to start afplay err=%v", err)
	}
}

func appleScriptCommand(args ...string) *exec.Cmd {
	if runtime.GOOS != "darwin" {
		return exec.Command("osascript", args...)
	}
	uid := os.Getuid()
	if uid > 0 {
		cmdArgs := append([]string{"asuser", fmt.Sprintf("%d", uid), "osascript"}, args...)
		return exec.Command("launchctl", cmdArgs...)
	}
	return exec.Command("osascript", args...)
}

// SendWithActions shows a notification. Action buttons are not supported
// without the native notification helper, so this falls back to a simple notification.
// Always returns -1 (no action selected) since action buttons are not available.
func SendWithActions(title, message, iconPath string, actions []string, timeoutSec int) (int, error) {
	logging.Debug("-> SendWithActions(title=%q, actions=%v, timeout=%d)", title, actions, timeoutSec)
	defer logging.Debug("<- SendWithActions")

	if runtime.GOOS != "darwin" {
		return -1, nil
	}

	// Action buttons require a native macOS app helper.
	// Fall back to simple notification.
	logging.Debug("SendWithActions: action buttons not available, sending simple notification")
	if err := Send(title, message); err != nil {
		return -1, err
	}
	return -1, nil
}
