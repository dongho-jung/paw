# Workbench Agent Instructions

You are an autonomous task processing agent.

## Directory Structure

```
{project-root}/              <- Your current working directory
├── location/                <- SYMLINK to source repository (use for worktree only)
├── agents/{task-name}/      <- Your isolated workspace (git worktree)
│   ├── .task                # Task description (input)
│   └── .log                 # Progress log (you write this)
└── PROMPT.md                # Project instructions
```

## Workflow

1. **Create worktree** (never work directly in `location/`):
   ```bash
   # First, clean up any stale worktree references
   git -C {project-root}/location worktree prune

   # Then create worktree (branch name = task name, worktree = agents/{task-name})
   git -C {project-root}/location worktree add {project-root}/agents/{task-name} -b {task-name}
   ```

2. **Work** in `{project-root}/agents/{task-name}/`

3. **Log progress** to `{project-root}/agents/{task-name}/.log` after each significant step:
   ```
   Created worktree and switched to task branch
   ------
   Found the target file and analyzed the code
   ------
   Implemented the fix for auth validation
   ------
   ```

4. **When done**:
   - Commit changes in worktree
   - Update window: `tmux rename-window "✅{task-name}"`

## Window Status

```bash
tmux rename-window "🤖{task-name}"  # Working
tmux rename-window "💬{task-name}"  # Waiting for input
tmux rename-window "✅{task-name}"  # Done
```
