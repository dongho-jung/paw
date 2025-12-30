# TAW (Tmux + Agent + Worktree)

Claude Code-based autonomous agent work environment

## Keyboard Shortcuts

### Mouse
  Click           Select pane
  Drag            Select text (copy mode)
  Scroll          Scroll pane
  Border drag     Resize pane

### Navigation
  ⌥Tab        Move to next pane (cycle)
  ⌥←/→        Move to previous/next window
  ⌃C ⌃C       Exit session (press twice within 1 second)
  ⌃D ⌃D       Exit session (press twice within 1 second)

### Command Palette
  ⌃R          Open command palette (fuzzy search)

Available commands in palette:
  new-task      Create a new task
  end-task      End current task
  show-tasks    Show task list (active + done)
  show-log      Show log viewer
  show-shell    Open shell pane
  show-help     Show this help
  add-queue     Add task to queue
  merge-all     Merge all completed tasks
  detach        Exit session

## Slash Commands (for agents)

  /commit     Smart commit (auto-generate message from diff analysis)
  /test       Auto-detect and run project tests
  /pr         Auto-create PR and open browser
  /merge      Merge worktree branch to project branch

## Directory Structure

  .taw/
  ├── config                 Project configuration file
  ├── PROMPT.md              Project-specific agent instructions
  ├── memory                 Shared project memory (YAML)
  ├── log                    Unified log file
  ├── .queue/                Quick task queue (add with add-queue)
  ├── history/               Completed task history
  │   └── YYMMDD_HHMMSS_name Task content + work capture
  └── agents/{task-name}/
      ├── task               Task content
      ├── origin/            Project root (symlink)
      └── worktree/          git worktree (auto-created)

## Window Status Icons

  ⭐️  New task input window
  🤖  Agent working
  💬  Waiting for user input
  ✅  Task completed
  ⚠️  Warning (merge failed, needs manual resolution)

## Task List Viewer (show-tasks)

View all active and completed tasks with preview panel.

### Navigation
  ↑/↓         Navigate tasks
  PgUp/PgDn   Scroll preview panel
  ⏎/Space     Focus selected task window
  q/Esc       Close the task list

### Actions
  c           Cancel task (active tasks only)
  m           Merge task (triggers end-task flow)
  p           Push branch to remote
  r           Resume task (history items only, creates new task)

### Status Icons
  🤖  Working (agent active)
  💬  Waiting (needs user input)
  ✅  Done (ready to merge)
  📁  History (completed, from history)

## Log Viewer (show-log)

  ↑/↓         Scroll vertically
  ←/→         Scroll horizontally (when word wrap is off)
  g           Jump to top
  G           Jump to bottom
  PgUp/PgDn   Page scroll
  s           Toggle tail mode (follow new logs)
  w           Toggle word wrap
  l           Cycle log level filter (L0+ → L1+ → ... → L5 only)
  q/Esc       Close the log viewer

## Environment Variables (for agents)

  TASK_NAME     Task name
  TAW_DIR       .taw directory path
  PROJECT_DIR   Project root path
  WORKTREE_DIR  Worktree path
  WINDOW_ID     tmux window ID (for status updates)

---
Press q to exit
