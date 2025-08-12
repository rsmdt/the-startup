# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

### Build and Run
```bash
# Build the binary
go build -o the-startup

# Run directly without building
go run . install
go run . log --assistant < input.json
go run . log --user < input.json

# Run the compiled binary
./the-startup install
./the-startup --help
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/log/...
go test ./internal/ui/...
go test ./internal/installer/...

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Development Workflow
```bash
# Format code
go fmt ./...

# Check for issues
go vet ./...

# Update dependencies
go mod tidy
```

## Architecture Overview

### Package Structure

The project follows a standard Go layout with clear separation of concerns:

- **`main.go`**: Entry point that sets up Cobra commands and embeds asset files
- **`cmd/`**: Command implementations using Cobra framework
  - `install.go`: Installation command that launches BubbleTea TUI
  - `log.go`: Processes hook data from Claude Code via stdin
  - `commands.go`: Other commands (update, validate, hooks)
  
- **`internal/`**: Core application logic
  - `installer/`: Installation logic and file management
    - Handles copying embedded assets to appropriate directories
    - Manages lock files and configuration updates
  - `ui/`: BubbleTea-based interactive TUI components
    - Composable model pattern with state machine
    - Separate models for each installation step
  - `log/`: JSONL logging for agent interactions
    - Processes PreToolUse and PostToolUse hook data
    - Manages session-based and global log files
  - `config/`: Configuration structures (lock files)
  - `assets/`: Embedded filesystem management

### Embedded Assets

The application embeds all assets at compile time using Go's embed package:
- `assets/agents/*.md`: Agent definitions
- `assets/commands/**/*.md`: Command definitions with nested structure
- `assets/templates/*`: Template files (BRD, PRD, SDD, PLAN)
- `assets/settings.json`: Claude Code settings template


### UI Architecture (BubbleTea)

The installer uses a composable model pattern with clear state transitions:

1. **MainModel**: Orchestrates the overall flow by composing sub-models
2. **State Machine**: Manages transitions between installation steps
   - StateStartupPath → StateClaudePath → StateFileSelection → StateComplete
3. **Sub-models**: Each handles a specific step with its own Update/View logic
   - StartupPathModel: Selects installation directory
   - ClaudePathModel: Selects Claude configuration directory  
   - FileSelectionModel: Interactive tree selector for choosing files
   - CompleteModel: Shows installation success
   - ErrorModel: Handles error display

### Installation Flow

1. User runs `the-startup install`
2. TUI launches with path selection
3. Files are selected using tree navigation
4. Assets are copied to:
   - `.claude/agents/` and `.claude/commands/`: Agent and command definitions
   - `.the-startup/templates/`: Template files
   - `.the-startup/bin/`: The startup binary
5. Settings.json is updated with hooks configuration
6. Lock file is created to track installation

### Hook Processing

The `log` command processes Claude Code hook data:
- Reads JSON from stdin (PreToolUse or PostToolUse events)
- Filters for Task tool with specific subagent types
- Writes JSONL logs to `.the-startup/<session-id>/agent-instructions.jsonl`
- Maintains a global log at `.the-startup/all-agent-instructions.jsonl`

## Key Implementation Details

### File Path Handling
- Installation paths support `~` expansion for home directory
- Project-local installation uses `.the-startup` directory
- Claude configuration expected at `~/.claude`

### Placeholder Replacement
Templates use placeholders that are replaced during installation:
- `{{STARTUP_PATH}}`: Installation directory path
- `{{CLAUDE_PATH}}`: Claude configuration directory

### Session Management
Hook logging automatically detects Claude Code session IDs from:
- `CLAUDE_SESSION_ID` environment variable
- `claude_session_id` in JSON input
- Falls back to finding latest `dev-*` directory

### Error Handling
- Hook commands exit silently (code 0) on errors to not disrupt Claude Code
- Debug output available via `DEBUG_HOOKS=1` environment variable
- Installation validates paths and provides clear error messages