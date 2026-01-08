package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dongho-jung/paw/internal/app"
	"github.com/dongho-jung/paw/internal/config"
	"github.com/dongho-jung/paw/internal/constants"
	"github.com/dongho-jung/paw/internal/embed"
	"github.com/dongho-jung/paw/internal/git"
	"github.com/dongho-jung/paw/internal/task"
	"github.com/dongho-jung/paw/internal/tmux"
)

type doctorCheck struct {
	name     string
	ok       bool
	message  string
	required bool
	fix      func() error
}

var doctorFix bool

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose project and session health",
	RunE:  runDoctor,
}

func init() {
	doctorCmd.Flags().BoolVar(&doctorFix, "fix", false, "Attempt safe fixes where possible")
}

func runDoctor(cmd *cobra.Command, args []string) error {
	application, err := buildAppFromCwd()
	if err != nil {
		return err
	}

	results := doctorChecks(application)
	printDoctorResults(results)

	if doctorFix {
		fixResults := applyDoctorFixes(results)
		if len(fixResults) > 0 {
			fmt.Println()
			fmt.Println("Fixes:")
			printDoctorResults(fixResults)
		}

		results = doctorChecks(application)
		if len(results) > 0 {
			fmt.Println()
			fmt.Println("Recheck:")
			printDoctorResults(results)
		}
	}

	hasErrors := false
	for _, r := range results {
		if r.required && !r.ok {
			hasErrors = true
			break
		}
	}
	if hasErrors {
		return fmt.Errorf("doctor found required issues")
	}
	return nil
}

func buildAppFromCwd() (*app.App, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	gitClient := git.New()
	isGitRepo := gitClient.IsGitRepo(cwd)
	projectDir := cwd
	if isGitRepo {
		if repoRoot, err := gitClient.GetRepoRoot(cwd); err == nil {
			projectDir = repoRoot
		}
	}

	application, err := app.New(projectDir)
	if err != nil {
		return nil, err
	}

	pawHome, _ := getPawHome()
	application.SetPawHome(pawHome)
	application.SetGitRepo(isGitRepo)

	return application, nil
}

func doctorChecks(appCtx *app.App) []doctorCheck {
	pawDirExists := pathExists(appCtx.PawDir)

	results := []doctorCheck{
		{
			name:     ".paw directory",
			ok:       pawDirExists,
			required: true,
			message:  boolMessage(pawDirExists, appCtx.PawDir, "missing (run paw to initialize)"),
			fix: func() error {
				return appCtx.Initialize()
			},
		},
	}

	if !pawDirExists {
		return results
	}

	configPath := filepath.Join(appCtx.PawDir, constants.ConfigFileName)
	configExists := pathExists(configPath)
	results = append(results, doctorCheck{
		name:     "config",
		ok:       configExists,
		required: false,
		message:  boolMessage(configExists, configPath, "missing (run paw setup)"),
		fix: func() error {
			cfg := config.DefaultConfig()
			return cfg.Save(appCtx.PawDir)
		},
	})

	claudeDir := filepath.Join(appCtx.PawDir, constants.ClaudeLink)
	claudeExists := pathExists(claudeDir)
	results = append(results, doctorCheck{
		name:     "claude settings",
		ok:       claudeExists,
		required: false,
		message:  boolMessage(claudeExists, claudeDir, "missing"),
		fix: func() error {
			return embed.WriteClaudeFiles(claudeDir)
		},
	})

	historyDir := filepath.Join(appCtx.PawDir, constants.HistoryDirName)
	historyExists := pathExists(historyDir)
	results = append(results, doctorCheck{
		name:     "history directory",
		ok:       historyExists,
		required: false,
		message:  boolMessage(historyExists, historyDir, "missing"),
		fix: func() error {
			return os.MkdirAll(historyDir, 0755)
		},
	})

	agentsDir := filepath.Join(appCtx.PawDir, constants.AgentsDirName)
	agentsExists := pathExists(agentsDir)
	results = append(results, doctorCheck{
		name:     "agents directory",
		ok:       agentsExists,
		required: false,
		message:  boolMessage(agentsExists, agentsDir, "missing"),
		fix: func() error {
			return os.MkdirAll(agentsDir, 0755)
		},
	})

	if configExists {
		cfg, err := config.Load(appCtx.PawDir)
		if err == nil {
			appCtx.Config = cfg
			warnings := cfg.Normalize()
			if len(warnings) > 0 {
				results = append(results, doctorCheck{
					name:     "config values",
					ok:       false,
					required: false,
					message:  fmt.Sprintf("normalize warnings: %s", stringsJoin(warnings)),
				})
			} else {
				results = append(results, doctorCheck{
					name:     "config values",
					ok:       true,
					required: false,
					message:  "ok",
				})
			}
		}
	}

	results = append(results, worktreeChecks(appCtx)...)
	results = append(results, sessionChecks(appCtx)...)

	return results
}

func worktreeChecks(appCtx *app.App) []doctorCheck {
	if !appCtx.IsGitRepo || appCtx.Config == nil || appCtx.Config.WorkMode != config.WorkModeWorktree {
		return nil
	}

	mgr := task.NewManager(appCtx.AgentsDir, appCtx.ProjectDir, appCtx.PawDir, appCtx.IsGitRepo, appCtx.Config)
	corrupted, err := mgr.FindCorruptedTasks()
	if err != nil {
		return []doctorCheck{
			{
				name:     "worktree health",
				ok:       false,
				required: false,
				message:  fmt.Sprintf("error: %v", err),
			},
		}
	}

	if len(corrupted) == 0 {
		return []doctorCheck{
			{
				name:     "worktree health",
				ok:       true,
				required: false,
				message:  "ok",
			},
		}
	}

	return []doctorCheck{
		{
			name:     "worktree health",
			ok:       false,
			required: false,
			message:  fmt.Sprintf("corrupted worktrees: %d", len(corrupted)),
		},
	}
}

func sessionChecks(appCtx *app.App) []doctorCheck {
	if _, err := os.Stat(appCtx.PawDir); err != nil {
		return nil
	}

	if _, err := exec.LookPath("tmux"); err != nil {
		return nil
	}

	tm := tmux.New(appCtx.SessionName)
	if !tm.HasSession(appCtx.SessionName) {
		return []doctorCheck{
			{
				name:     "tmux session",
				ok:       false,
				required: false,
				message:  "not running",
			},
		}
	}

	mgr := task.NewManager(appCtx.AgentsDir, appCtx.ProjectDir, appCtx.PawDir, appCtx.IsGitRepo, appCtx.Config)
	mgr.SetTmuxClient(tm)

	orphaned, err := mgr.FindOrphanedWindows()
	if err != nil {
		return []doctorCheck{
			{
				name:     "tmux session",
				ok:       false,
				required: false,
				message:  fmt.Sprintf("session check failed: %v", err),
			},
		}
	}

	stopped, err := mgr.FindStoppedTasks()
	if err != nil {
		return []doctorCheck{
			{
				name:     "tmux session",
				ok:       false,
				required: false,
				message:  fmt.Sprintf("session check failed: %v", err),
			},
		}
	}

	incomplete, err := mgr.FindIncompleteTasks(appCtx.SessionName)
	if err != nil {
		return []doctorCheck{
			{
				name:     "tmux session",
				ok:       false,
				required: false,
				message:  fmt.Sprintf("session check failed: %v", err),
			},
		}
	}

	msg := fmt.Sprintf("orphaned=%d stopped=%d incomplete=%d", len(orphaned), len(stopped), len(incomplete))
	ok := len(orphaned) == 0 && len(stopped) == 0 && len(incomplete) == 0

	return []doctorCheck{
		{
			name:     "tmux session",
			ok:       ok,
			required: false,
			message:  msg,
		},
	}
}

func printDoctorResults(results []doctorCheck) {
	for _, r := range results {
		printDoctorResult(r)
	}
}

func printDoctorResult(r doctorCheck) {
	var icon string
	if r.ok {
		icon = "[OK]"
	} else if r.required {
		icon = "[ERR]"
	} else {
		icon = "[WARN]"
	}

	optionalSuffix := ""
	if !r.required && !r.ok {
		optionalSuffix = " (optional)"
	}

	fmt.Printf("%s %s: %s%s\n", icon, r.name, r.message, optionalSuffix)
}

func applyDoctorFixes(results []doctorCheck) []doctorCheck {
	var fixes []doctorCheck
	for _, r := range results {
		if r.ok || r.fix == nil {
			continue
		}
		fixCheck := doctorCheck{
			name:     r.name + " fix",
			required: false,
		}
		if err := r.fix(); err != nil {
			fixCheck.ok = false
			fixCheck.message = err.Error()
		} else {
			fixCheck.ok = true
			fixCheck.message = "applied"
		}
		fixes = append(fixes, fixCheck)
	}
	return fixes
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func boolMessage(ok bool, okMessage, badMessage string) string {
	if ok {
		return okMessage
	}
	return badMessage
}

func stringsJoin(values []string) string {
	return strings.Join(values, "; ")
}
