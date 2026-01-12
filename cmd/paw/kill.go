package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/dongho-jung/paw/internal/tmux"
)

const (
	gracefulShutdownTimeout = 15 * time.Second
	shutdownPollInterval    = 500 * time.Millisecond
)

var killCmd = &cobra.Command{
	Use:   "kill [session]",
	Short: "Kill PAW tmux sessions",
	Long: `Kill a PAW tmux session without removing .paw directory.

If a session name is provided, kills that specific session.
If no session name is provided, lists available sessions to choose from.

This command sends SIGINT (Ctrl+C) to all processes in the session first,
waits up to 15 seconds for graceful shutdown, then force kills if needed.

Unlike 'paw clean', this preserves the .paw directory, worktrees, and branches.

Examples:
  paw kill             # List and select a session to kill
  paw kill myproject   # Kill the 'myproject' session directly

See also: paw kill-all (to kill all sessions at once)`,
	Args: cobra.MaximumNArgs(1),
	RunE: runKill,
}

var killAllCmd = &cobra.Command{
	Use:   "kill-all",
	Short: "Kill all running PAW sessions",
	Long: `Kill all running PAW tmux sessions without removing .paw directories.

This command finds all PAW sessions and kills them gracefully:
1. Sends SIGINT (Ctrl+C) to all processes
2. Waits up to 15 seconds for graceful shutdown
3. Force kills remaining sessions

Unlike 'paw clean', this preserves .paw directories, worktrees, and branches.`,
	RunE: runKillAll,
}

// Note: killAllCmd is registered in main.go as a root command

func runKill(cmd *cobra.Command, args []string) error {
	sessions, err := findPawSessions()
	if err != nil {
		return fmt.Errorf("failed to find PAW sessions: %w", err)
	}

	if len(sessions) == 0 {
		fmt.Println("No running PAW sessions found.")
		return nil
	}

	var targetSession pawSession

	if len(args) == 1 {
		// Direct kill of specified session
		sessionName := args[0]
		for _, s := range sessions {
			if s.Name == sessionName {
				targetSession = s
				break
			}
		}
		if targetSession.Name == "" {
			// Try partial match
			var matches []pawSession
			for _, s := range sessions {
				if strings.Contains(s.Name, sessionName) {
					matches = append(matches, s)
				}
			}
			if len(matches) == 1 {
				targetSession = matches[0]
			} else if len(matches) > 1 {
				fmt.Printf("Multiple sessions match '%s':\n", sessionName)
				for _, m := range matches {
					fmt.Printf("  - %s\n", m.Name)
				}
				return fmt.Errorf("please specify a unique session name")
			} else {
				fmt.Printf("Session '%s' not found.\n\n", sessionName)
				fmt.Println("Available sessions:")
				for _, s := range sessions {
					fmt.Printf("  - %s\n", s.Name)
				}
				return fmt.Errorf("session not found")
			}
		}
	} else if len(sessions) == 1 {
		// Only one session, confirm and kill
		targetSession = sessions[0]
		fmt.Printf("Found session: %s\n", targetSession.Name)
	} else {
		// Multiple sessions, prompt for selection
		fmt.Println("Running PAW sessions:")
		fmt.Println()
		for i, s := range sessions {
			fmt.Printf("  %d. %s\n", i+1, s.Name)
		}
		fmt.Println()
		fmt.Print("Select session to kill [1]: ")

		var input string
		_, _ = fmt.Scanln(&input)
		input = strings.TrimSpace(input)

		// Default to first session
		idx := 0
		if input != "" {
			var n int
			if _, err := fmt.Sscanf(input, "%d", &n); err != nil || n < 1 || n > len(sessions) {
				return fmt.Errorf("invalid selection: %s", input)
			}
			idx = n - 1
		}
		targetSession = sessions[idx]
	}

	return killSession(targetSession, true)
}

func runKillAll(cmd *cobra.Command, args []string) error {
	sessions, err := findPawSessions()
	if err != nil {
		return fmt.Errorf("failed to find PAW sessions: %w", err)
	}

	if len(sessions) == 0 {
		fmt.Println("No running PAW sessions found.")
		return nil
	}

	fmt.Printf("Found %d PAW session(s) to kill:\n", len(sessions))
	for _, s := range sessions {
		fmt.Printf("  - %s\n", s.Name)
	}
	fmt.Println()

	var failed []string
	for _, s := range sessions {
		if err := killSession(s, false); err != nil {
			failed = append(failed, fmt.Sprintf("%s: %v", s.Name, err))
		}
	}

	if len(failed) > 0 {
		fmt.Println("\nFailed to kill some sessions:")
		for _, f := range failed {
			fmt.Printf("  - %s\n", f)
		}
		return fmt.Errorf("failed to kill %d session(s)", len(failed))
	}

	fmt.Println("\nAll sessions killed successfully.")
	return nil
}

// killSession kills a PAW session with graceful shutdown.
// If verbose is true, prints progress messages.
func killSession(session pawSession, verbose bool) error {
	tm := tmux.New(session.Name)

	// Check if session still exists
	if !tm.HasSession(session.Name) {
		if verbose {
			fmt.Printf("Session '%s' is no longer running.\n", session.Name)
		}
		return nil
	}

	if verbose {
		fmt.Printf("Killing session '%s'...\n", session.Name)
	}

	// Step 1: Send Ctrl+C to all panes for graceful shutdown
	if verbose {
		fmt.Print("  Sending interrupt signal to processes...")
	}
	sendInterruptToAllPanes(tm, session.Name)
	if verbose {
		fmt.Println(" done")
	}

	// Step 2: Wait for graceful shutdown (poll for session termination)
	if verbose {
		fmt.Printf("  Waiting up to %s for graceful shutdown...", gracefulShutdownTimeout)
	}

	deadline := time.Now().Add(gracefulShutdownTimeout)
	sessionDied := false

	for time.Now().Before(deadline) {
		if !tm.HasSession(session.Name) {
			sessionDied = true
			break
		}
		time.Sleep(shutdownPollInterval)
	}

	if sessionDied {
		if verbose {
			fmt.Println(" processes exited gracefully")
			fmt.Println("  Session terminated.")
		} else {
			fmt.Printf("Killed: %s (graceful)\n", session.Name)
		}
		return nil
	}

	// Step 3: Force kill the session
	if verbose {
		fmt.Println(" timeout, force killing")
		fmt.Print("  Force killing session...")
	}

	if err := tm.KillSession(session.Name); err != nil {
		if verbose {
			fmt.Println(" failed")
		}
		return fmt.Errorf("failed to kill session: %w", err)
	}

	if verbose {
		fmt.Println(" done")
		fmt.Println("  Session terminated.")
	} else {
		fmt.Printf("Killed: %s (forced)\n", session.Name)
	}

	return nil
}

// sendInterruptToAllPanes sends Ctrl+C to all panes in a session.
func sendInterruptToAllPanes(tm tmux.Client, sessionName string) {
	// Get all windows in the session
	windows, err := tm.ListWindows()
	if err != nil {
		return
	}

	// Send Ctrl+C to each window's active pane
	// Note: We target the window which sends to its active pane
	for _, w := range windows {
		// Use window ID to target the pane
		target := sessionName + ":" + w.Name
		_ = tm.SendKeys(target, "C-c")
	}
}

