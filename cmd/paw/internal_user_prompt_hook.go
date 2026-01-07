package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/dongho-jung/paw/internal/constants"
	"github.com/dongho-jung/paw/internal/logging"
	"github.com/dongho-jung/paw/internal/tmux"
)

var userPromptSubmitHookCmd = &cobra.Command{
	Use:   "user-prompt-submit-hook",
	Short: "Handle Claude UserPromptSubmit hook to set working status",
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionName := os.Getenv("SESSION_NAME")
		windowID := os.Getenv("WINDOW_ID")
		taskName := os.Getenv("TASK_NAME")
		if sessionName == "" || windowID == "" || taskName == "" {
			return nil
		}

		if pawDir := os.Getenv("PAW_DIR"); pawDir != "" {
			logger, _ := logging.New(filepath.Join(pawDir, constants.LogFileName), os.Getenv("PAW_DEBUG") == "1")
			if logger != nil {
				defer func() { _ = logger.Close() }()
				logger.SetScript("user-prompt-submit-hook")
				logger.SetTask(taskName)
				logging.SetGlobal(logger)
			}
		}

		logging.Trace("userPromptSubmitHookCmd: start session=%s windowID=%s task=%s", sessionName, windowID, taskName)
		defer logging.Trace("userPromptSubmitHookCmd: end")

		tm := tmux.New(sessionName)
		paneID := windowID + ".0"
		if !tm.HasPane(paneID) {
			logging.Debug("userPromptSubmitHookCmd: pane %s not found, skipping", paneID)
			return nil
		}

		// Set window to working state (user submitted a prompt, agent is now working)
		newName := constants.EmojiWorking + constants.TruncateForWindowName(taskName)
		if err := renameWindowCmd.RunE(renameWindowCmd, []string{windowID, newName}); err != nil {
			logging.Warn("userPromptSubmitHookCmd: failed to rename window: %v", err)
			return nil
		}

		logging.Debug("userPromptSubmitHookCmd: status updated to working")
		return nil
	},
}
