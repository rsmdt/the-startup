# Go Codebase Analysis for TypeScript Migration

**Document Version:** 1.0
**Analysis Date:** 2025-10-06
**Scope:** Complete Go → TypeScript migration analysis
**Excluded:** `internal/stats/` (per PRD line 342)

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Architecture Overview](#architecture-overview)
3. [CLI Command Analysis](#cli-command-analysis)
4. [BubbleTea TUI State Machine](#bubbletea-tui-state-machine)
5. [Business Logic & Installer](#business-logic--installer)
6. [Embedded Asset System](#embedded-asset-system)
7. [Error Handling Patterns](#error-handling-patterns)
8. [Edge Cases & Gotchas](#edge-cases--gotchas)
9. [TypeScript Migration Mapping](#typescript-migration-mapping)
10. [Risk Assessment](#risk-assessment)

---

## Executive Summary

The Go codebase is a well-structured CLI tool with TUI (Terminal User Interface) capabilities built on Cobra (CLI) and BubbleTea (TUI). The application manages installation/uninstallation of Claude Code agents and commands via an interactive interface.

**Key Statistics:**
- Total Go files: 76 (excluding stats package)
- CLI commands: 6 (install, uninstall, init, spec, statusline, stats*)
- TUI states: 9 (install: 5, uninstall: 4)
- Embedded asset directories: 2 (`assets/claude/`, `assets/the-startup/`)
- Primary frameworks: Cobra (CLI), BubbleTea (TUI), Lipgloss (styling)

*Note: `stats` command is excluded from migration scope per PRD.

---

## Architecture Overview

### Module Structure (go.mod)

```go
module github.com/rsmdt/the-startup
go 1.23.0

// Primary dependencies
- github.com/spf13/cobra v1.8.0           // CLI framework
- github.com/charmbracelet/bubbletea v1.3.4  // TUI framework
- github.com/charmbracelet/lipgloss v1.1.0   // TUI styling
- github.com/charmbracelet/bubbles v0.21.0   // TUI components
- golang.org/x/term v0.18.0               // Terminal operations
```

**TypeScript Equivalents:**
- Cobra → Commander.js
- BubbleTea → Ink (React-based TUI)
- Lipgloss → Ink components with chalk/gradient-string
- embed.FS → pkg.files API or bundled assets

### File Organization

```
the-startup/
├── main.go                          # Entry point, embeds assets
├── cmd/                             # Command implementations
│   ├── install.go                   # Install command + flags
│   ├── uninstall.go                 # Uninstall command + flags
│   ├── init.go                      # Template initialization
│   ├── spec.go                      # Spec directory management
│   ├── statusline.go                # Status line generator
│   └── stats.go                     # [EXCLUDED FROM MIGRATION]
├── internal/
│   ├── installer/                   # Installation business logic
│   │   ├── installer.go             # Core installer (1083 lines)
│   │   └── installer_test.go
│   ├── config/
│   │   └── config.go                # Lock file structures
│   ├── ui/                          # BubbleTea TUI models
│   │   ├── model.go                 # Main model orchestrator
│   │   ├── states.go                # State machine definitions
│   │   ├── startup_path_model.go    # Path selection model
│   │   ├── claude_path_model.go     # Claude path model
│   │   ├── file_selection_model.go  # File selection tree
│   │   ├── complete_model.go        # Success screen
│   │   ├── error_model.go           # Error screen
│   │   ├── theme.go                 # Color themes
│   │   └── ascii.go                 # ASCII art banner
│   └── stats/                       # [EXCLUDED FROM MIGRATION]
└── assets/                          # Embedded at compile time
    ├── claude/                      # Goes to ~/.claude
    │   ├── agents/
    │   ├── commands/
    │   ├── output-styles/
    │   ├── settings.json
    │   └── settings.local.json
    └── the-startup/                 # Goes to install path
        ├── templates/
        └── rules/
```

---

## CLI Command Analysis

### main.go - Entry Point

```go
// Version info set by build flags
var (
    Version   = "dev"
    GitCommit = "unknown"
    BuildDate = "unknown"
)

// Embedded assets using Go's embed directive
//go:embed assets/claude
var claudeAssets embed.FS

//go:embed assets/the-startup
var startupAssets embed.FS

// Root command setup
rootCmd := &cobra.Command{
    Use:   "the-startup",
    Short: "Agent system for development tools",
    Version: fmt.Sprintf("%s (commit: %s, built: %s)", Version, GitCommit, BuildDate),
}

// Commands attached
rootCmd.AddCommand(cmd.NewInstallCommand(&claudeAssets, &startupAssets))
rootCmd.AddCommand(cmd.NewUninstallCommand(&claudeAssets, &startupAssets))
rootCmd.AddCommand(cmd.NewInitCommand(&startupAssets))
rootCmd.AddCommand(cmd.NewStatsCommand())  // EXCLUDED
rootCmd.AddCommand(cmd.NewStatuslineCommand())
rootCmd.AddCommand(cmd.NewSpecCommand(&startupAssets))
```

**Migration Notes:**
- Use Commander.js for root command
- Version info: use package.json version + git-rev-sync
- Embedded assets: pkg.files or webpack/esbuild bundling

---

### Command: `install`

**File:** `cmd/install.go`

```go
func NewInstallCommand(claudeAssets, startupAssets *embed.FS) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "install",
        Short: "Install The Startup agent system",
        Long:  `Install agents, hooks, and commands for development tools with an interactive TUI`,
        RunE: func(cmd *cobra.Command, args []string) error {
            local, _ := cmd.Flags().GetBool("local")
            yes, _ := cmd.Flags().GetBool("yes")
            return ui.RunMainInstallerWithFlags(claudeAssets, startupAssets, local, yes)
        },
    }

    cmd.Flags().BoolP("local", "l", false, "Use local installation paths (skip path selection)")
    cmd.Flags().BoolP("yes", "y", false, "Auto-confirm with recommended paths")

    return cmd
}
```

**Flags:**
- `--local` / `-l`: Pre-selects local paths (`.the-startup`, `.claude`)
- `--yes` / `-y`: Auto-confirms with global paths (no TUI interaction)

**Behavior:**
- With `--yes` + `--local`: Install to local paths without TUI
- With `--yes` only: Install to global paths without TUI
- With `--local` only: Show TUI but pre-select local paths
- Default: Full interactive TUI

**Error Messages:**
- None specific; errors bubble from UI/installer layers

---

### Command: `uninstall`

**File:** `cmd/uninstall.go`

```go
type UninstallFlags struct {
    DryRun       bool  // Show what would be removed
    Force        bool  // Skip confirmation prompts
    KeepLogs     bool  // Preserve log files
    KeepSettings bool  // Preserve settings
}

func NewUninstallCommand(claudeAssets, startupAssets *embed.FS) *cobra.Command {
    var flags UninstallFlags

    cmd := &cobra.Command{
        Use:   "uninstall",
        Short: "Uninstall The Startup agent system",
        Long:  `Uninstall agents, hooks, and commands with interactive TUI...`,
        RunE: func(cmd *cobra.Command, args []string) error {
            return runUninstall(cmd, flags, claudeAssets, startupAssets)
        },
    }

    cmd.Flags().BoolVar(&flags.DryRun, "dry-run", false, "Preview removal without changes")
    cmd.Flags().BoolVar(&flags.Force, "force", false, "Force removal without prompts")
    cmd.Flags().BoolVar(&flags.KeepLogs, "keep-logs", false, "Preserve logs")
    cmd.Flags().BoolVar(&flags.KeepSettings, "keep-settings", false, "Preserve settings")

    return cmd
}
```

**Flags:**
- `--dry-run`: Preview mode (NOT IMPLEMENTED in current simplified flow)
- `--force`: Skip confirmations (NOT IMPLEMENTED)
- `--keep-logs`: Preserve logs (NOT IMPLEMENTED)
- `--keep-settings`: Preserve settings (NOT IMPLEMENTED)

**Current Behavior:**
- Always launches interactive TUI
- Flags are defined but not used (TODO comments in code)

**Migration Note:** Decide whether to implement these flags or remove them.

---

### Command: `init`

**File:** `cmd/init.go` (391 lines)

**Purpose:** Initialize quality gate templates (DoR, DoD, Task-DoD)

```go
func NewInitCommand(startupAssets *embed.FS) *cobra.Command {
    var skipPrompts bool
    var force bool
    var dryRun bool

    cmd := &cobra.Command{
        Use:   "init [template]",
        Short: "Initialize quality gate templates",
        Long:  `Initialize Definition of Ready and Definition of Done templates...`,
        RunE: func(cmd *cobra.Command, args []string) error {
            template := ""
            if len(args) > 0 {
                template = args[0]
            }
            return runInit(startupAssets, template, skipPrompts, force, dryRun)
        },
    }

    cmd.Flags().BoolVarP(&skipPrompts, "skip-prompts", "s", false, "Use defaults")
    cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite without prompting")
    cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Check files without creating")

    return cmd
}
```

**Arguments:**
- `[template]`: Optional; one of `definition-of-ready`, `definition-of-done`, `task-definition-of-done`
- No argument: Initialize all three templates

**Flags:**
- `--skip-prompts` / `-s`: Skip guided setup questions, use defaults
- `--force` / `-f`: Overwrite existing files without prompting
- `--dry-run`: Check what files exist without creating/overwriting

**Workflow:**
1. Validate template name (case-insensitive)
2. Check existing files in `./docs/`
3. If dry-run, print status and exit
4. Create `./docs/` directory if needed
5. Prompt for overwrite if files exist (unless `--force`)
6. Copy templates from embedded assets
7. Run guided prompts (unless `--skip-prompts` or single template)
8. Display customization advice based on answers
9. Show next steps

**Guided Prompts (interactive):**
```go
type SetupAnswers struct {
    UsesTDD        bool
    CoverageTarget int    // 0-100
    BuildCommand   string // e.g., "go build ./..."
    TestCommand    string // e.g., "go test ./..."
    LintCommand    string // e.g., "golangci-lint run"
    FormatCommand  string // e.g., "gofmt -l ."
}
```

**Prompt Functions:**
- `promptYesNo(question, default)`: Y/n or y/N input
- `promptInt(question, default, min, max)`: Integer with validation
- `promptString(question, default)`: String input with default

**Error Messages:**
- "invalid template name: %s (valid: definition-of-ready, definition-of-done, task-definition-of-done)"
- "failed to create docs directory: %w"
- "cancelled by user"
- "failed to read template %s: %w"
- "failed to write %s: %w"

**Customization Advice Output:**
- Checks if answers differ from defaults
- Provides file editing instructions for each deviation
- Example: "Edit docs/definition-of-done.md: Find [NEEDS CLARIFICATION: coverage target] → Replace with 90"

---

### Command: `spec`

**File:** `cmd/spec.go` (463 lines)

**Purpose:** Manage specification directories

```go
func NewSpecCommand(startupAssets *embed.FS) *cobra.Command {
    var (
        readMode string
        addMode  string
    )

    cmd := &cobra.Command{
        Use:   "spec [feature description or ID]",
        Short: "Manage specification directories",
        Long:  `Creates new specification directories or manages existing ones...`,
        RunE: func(cmd *cobra.Command, args []string) error {
            if readMode != "" {
                return handleReadMode(readMode)
            }
            if addMode != "" {
                if len(args) == 0 {
                    return fmt.Errorf("spec ID required when using --add")
                }
                return handleAddMode(startupAssets, args[0], addMode)
            }
            if len(args) == 0 {
                return fmt.Errorf("feature description required")
            }

            description := strings.Join(args, " ")
            return handleCreateMode(startupAssets, description)
        },
    }

    cmd.Flags().StringVar(&readMode, "read", "", "Read existing specification by ID")
    cmd.Flags().StringVar(&addMode, "add", "", "Add template to spec (PRD, SDD, PLAN, BRD)")

    return cmd
}
```

**Modes:**

1. **Create Mode (default):**
   ```bash
   the-startup spec "user authentication system"
   ```
   - Scans `./docs/specs/` for highest ID (e.g., 009)
   - Creates next ID with 3-digit padding (010)
   - Sanitizes feature name to directory name
   - Creates `./docs/specs/010-user-auth/`
   - Copies PRD template to `./docs/specs/010-user-auth/product-requirements.md`
   - Outputs TOML format:
     ```toml
     id = "010"
     name = "user-auth"
     dir = "./docs/specs/010-user-auth"

     [spec]
     prd = "./docs/specs/010-user-auth/product-requirements.md"
     ```

2. **Read Mode:**
   ```bash
   the-startup spec --read 010
   the-startup spec --read 010-user-auth
   ```
   - Finds spec directory by ID or full name
   - Scans for known templates (PRD, SDD, PLAN, BRD)
   - Also scans for quality gate files in `./docs/`
   - Outputs TOML format:
     ```toml
     id = "010"
     name = "user-auth"
     dir = "./docs/specs/010-user-auth"

     [spec]
     prd = "./docs/specs/010-user-auth/product-requirements.md"
     sdd = "./docs/specs/010-user-auth/solution-design.md"

     [gates]
     definition_of_ready = "docs/definition-of-ready.md"
     definition_of_done = "docs/definition-of-done.md"
     task_definition_of_done = "docs/task-definition-of-done.md"
     ```

3. **Add Mode:**
   ```bash
   the-startup spec 010 --add SDD
   the-startup spec 010-user-auth --add PLAN
   ```
   - Finds or creates spec directory
   - Copies template to new filename:
     - PRD → `product-requirements.md`
     - SDD → `solution-design.md`
     - PLAN → `implementation-plan.md`
     - BRD → `business-requirements.md`
   - Outputs TOML with `[spec.new]` section:
     ```toml
     id = "010"
     name = "user-auth"
     dir = "./docs/specs/010-user-auth"

     [spec.new]
     sdd = "./docs/specs/010-user-auth/solution-design.md"
     ```

**Template Mapping:**
```go
var templateMapping = map[string]struct {
    filename     string
    templateFile string
}{
    "PRD":  {"product-requirements.md", "product-requirements"},
    "SDD":  {"solution-design.md", "solution-design"},
    "PLAN": {"implementation-plan.md", "implementation-plan"},
    "BRD":  {"business-requirements.md", "business-requirements"},
}
```

**Feature Name Sanitization:**
```go
func sanitizeFeatureName(description string) string {
    // 1. Convert to lowercase
    // 2. Replace non-alphanumeric with hyphens
    // 3. Remove leading/trailing hyphens
    // 4. Limit to first 4 words
    // 5. Filter out common words: a, an, the, and, or, for, with, to, of, in
    // Example: "Add a User Authentication System for Web" → "user-authentication-system-web"
}
```

**Directory Matching Logic:**
- Numeric ID (e.g., "010"): Looks for `010-*` prefix or exact `010` directory
- Full name (e.g., "010-user-auth"): Exact directory match
- Extracts ID and name from directory pattern `^\d{3}-(.+)$`

**Backward Compatibility:**
- Scans for new filenames first (e.g., `product-requirements.md`)
- Falls back to old filenames (e.g., `PRD.md`)

**Error Messages:**
- "invalid template type: %s (valid: PRD, SDD, PLAN, BRD)"
- "specification not found: %s"
- "spec ID required when using --add"
- "feature description required"
- "file already exists: %s"

---

### Command: `statusline`

**File:** `cmd/statusline.go` (159 lines)

**Purpose:** Generate retro terminal status line for Claude Code hooks

```go
type StatuslineInput struct {
    HookEventName  string `json:"hook_event_name"`
    SessionID      string `json:"session_id"`
    TranscriptPath string `json:"transcript_path"`
    CWD            string `json:"cwd"`
    Model          struct {
        ID          string `json:"id"`
        DisplayName string `json:"display_name"`
    } `json:"model"`
    Workspace struct {
        CurrentDir string `json:"current_dir"`
        ProjectDir string `json:"project_dir"`
    } `json:"workspace"`
    Version     string `json:"version"`
    OutputStyle struct {
        Name string `json:"name"`
    } `json:"output_style"`
}

func NewStatuslineCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "statusline",
        Short: "Generate a retro terminal status line for Claude Code",
        RunE: func(cmd *cobra.Command, args []string) error {
            return runStatusline(cmd.InOrStdin(), cmd.OutOrStdout())
        },
        SilenceUsage:  true,
        SilenceErrors: true,
    }
}
```

**Input:** JSON from stdin (provided by Claude Code hooks)

**Output:** Formatted status line string

**Behavior:**
1. Read JSON from stdin
2. Silent fail if JSON parsing fails (hook compatibility)
3. Get terminal width from:
   - `COLUMNS` env var (most reliable in hooks)
   - `term.GetSize(os.Stdout.Fd())`
   - `term.GetSize(os.Stderr.Fd())`
   - Default: 120
4. Build status line with:
   - Current directory (with `~` for home)
   - Git branch (if in git repo): `⎇ branch-name`
   - Model name and output style: `🤖 Claude (default)`
   - Help text: `? for shortcuts` (italic, muted)
5. Apply lipgloss styling with max width truncation

**Git Integration:**
```go
func getGitInfo(workingDir string) string {
    // Check if in git repo: git rev-parse --git-dir
    // Get current branch: git branch --show-current
    // Return "⎇ branch-name" or "⎇ HEAD" if detached
}
```

**Terminal Width Detection:**
- Priority: COLUMNS env → os.Stdout → os.Stderr → 120 default
- Uses `golang.org/x/term` for size detection

**Error Handling:**
- Silent failures for hook compatibility
- No error output if JSON decode fails

**Migration Note:**
- Use `chalk` for styling in TypeScript
- Use `process.stdout.columns` for terminal width
- Use `execa` or `child_process` for git commands

---

## BubbleTea TUI State Machine

### State Architecture

**File:** `internal/ui/states.go`

```go
type InstallerState int

const (
    // Install workflow
    StateStartupPath     InstallerState = iota  // 0
    StateClaudePath                              // 1
    StateFileSelection                           // 2
    StateComplete                                // 3
    StateError                                   // 4

    // Uninstall workflow
    StateUninstallStartupPath                    // 5
    StateUninstallClaudePath                     // 6
    StateUninstallFileSelection                  // 7
    StateUninstallComplete                       // 8
)

type OperationMode int

const (
    ModeInstall   OperationMode = iota
    ModeUninstall
)
```

**State Machine Diagram:**

```
Install Workflow:
    StateStartupPath ──────> StateClaudePath ──────> StateFileSelection ──────> StateComplete
         │                        │                         │
         └──> StateError <────────┴─────────────────────────┘

Uninstall Workflow:
    StateUninstallStartupPath ──> StateUninstallClaudePath ──> StateUninstallFileSelection ──> StateUninstallComplete
                │                             │                             │
                └──> StateError <─────────────┴─────────────────────────────┘

Navigation:
- ESC: Go back one state (or quit from initial states)
- Ctrl+C or Q: Quit from initial states, go back from others
```

**Valid Transitions:**
```go
var ValidTransitions = map[StateTransition]bool{
    // Install - Forward
    {StateStartupPath, StateClaudePath}:   true,
    {StateClaudePath, StateFileSelection}: true,
    {StateFileSelection, StateComplete}:   true,
    {StateFileSelection, StateError}:      true,

    // Install - Backward (ESC)
    {StateClaudePath, StateStartupPath}:   true,
    {StateFileSelection, StateClaudePath}: true,

    // Uninstall - Forward
    {StateUninstallStartupPath, StateUninstallClaudePath}:   true,
    {StateUninstallClaudePath, StateUninstallFileSelection}: true,
    {StateUninstallFileSelection, StateUninstallComplete}:   true,
    {StateUninstallStartupPath, StateError}:                 true,
    {StateUninstallClaudePath, StateError}:                  true,
    {StateUninstallFileSelection, StateError}:               true,

    // Uninstall - Backward (ESC)
    {StateUninstallClaudePath, StateUninstallStartupPath}:   true,
    {StateUninstallFileSelection, StateUninstallClaudePath}: true,

    // Error recovery
    {StateError, StateStartupPath}:          true,
    {StateError, StateUninstallStartupPath}: true,
}
```

---

### MainModel - Orchestrator

**File:** `internal/ui/model.go` (414 lines)

```go
type MainModel struct {
    // State management
    state InstallerState
    mode  OperationMode

    // Dependencies
    installer     *installer.Installer
    claudeAssets  *embed.FS
    startupAssets *embed.FS

    // Configuration
    flags InstallFlags

    // User selections (shared state)
    startupPath   string
    claudePath    string
    selectedFiles []string

    // Sub-models
    startupPathModel   StartupPathModel
    claudePathModel    ClaudePathModel
    fileSelectionModel FileSelectionModel
    completeModel      CompleteModel
    errorModel         ErrorModel

    // UI state
    width  int
    height int
}

type InstallFlags struct {
    Local bool  // Use local paths
    Yes   bool  // Auto-confirm
}
```

**Initialization Flow:**

```go
func (m *MainModel) Init() tea.Cmd {
    // Handle --yes flag: Auto-install without TUI
    if m.flags.Yes {
        return m.handleAutoInstall()
    }

    // Handle --local flag: Pre-select local paths
    if m.flags.Local {
        return m.handleLocalFlag()
    }

    return nil
}

func (m *MainModel) handleAutoInstall() tea.Cmd {
    // Determine paths
    if m.flags.Local {
        startupPath = filepath.Join(cwd, ".the-startup")
        claudePath = filepath.Join(cwd, ".claude")
    } else {
        startupPath = filepath.Join(homeDir, ".config", "the-startup")
        claudePath = filepath.Join(homeDir, ".claude")
    }

    // Set paths and transition to StateFileSelection
    m.startupPath = startupPath
    m.claudePath = claudePath
    m.installer.SetInstallPath(startupPath)
    m.installer.SetClaudePath(claudePath)
    m.transitionToState(StateFileSelection)

    // Perform installation synchronously
    err := m.installer.Install()
    if err != nil {
        m.transitionToState(StateError)
    } else {
        m.transitionToState(StateComplete)
    }
    return m.completeModel.Init()
}

func (m *MainModel) handleLocalFlag() tea.Cmd {
    // Pre-select local paths but show TUI
    cwd, _ := os.Getwd()
    m.startupPath = filepath.Join(cwd, ".the-startup")
    m.claudePath = filepath.Join(cwd, ".claude")
    m.installer.SetInstallPath(m.startupPath)
    m.installer.SetClaudePath(m.claudePath)

    // Skip path selection and go to file selection
    m.transitionToState(StateFileSelection)
    return nil
}
```

**Update Loop:**

```go
func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Global keys
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        switch keyMsg.String() {
        case "ctrl+c", "q":
            if m.state == StateStartupPath || m.state == StateUninstallStartupPath {
                return m, tea.Quit
            }
            return m.handleBack()  // ESC behavior in other states
        case "esc":
            return m.handleBack()
        }
    }

    // Window size
    if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
        m.width = sizeMsg.Width
        m.height = sizeMsg.Height
        return m, nil
    }

    // Delegate to sub-models based on state
    switch m.state {
    case StateStartupPath, StateUninstallStartupPath:
        newModel, cmd := m.startupPathModel.Update(msg)
        m.startupPathModel = newModel
        if m.startupPathModel.Ready() {
            path := m.startupPathModel.SelectedPath()
            if path == "CANCEL" {
                return m, tea.Quit
            }
            m.startupPath = path
            m.installer.SetInstallPath(path)
            if m.mode == ModeUninstall {
                m.transitionToState(StateUninstallClaudePath)
            } else {
                m.transitionToState(StateClaudePath)
            }
        }
        return m, cmd

    case StateClaudePath, StateUninstallClaudePath:
        // Similar delegation to claudePathModel
        // ...

    case StateFileSelection, StateUninstallFileSelection:
        newModel, cmd := m.fileSelectionModel.Update(msg)
        m.fileSelectionModel = newModel
        if m.fileSelectionModel.Ready() {
            if m.fileSelectionModel.Confirmed() {
                if m.mode == ModeUninstall {
                    err := m.performUninstall()
                    // Handle error...
                } else {
                    err := m.installer.Install()
                    // Handle error...
                }
            } else {
                // User declined, go back
                m.transitionToState(...)
            }
        }
        return m, cmd

    // ... other states
    }

    return m, nil
}
```

**View Delegation:**

```go
func (m *MainModel) View() string {
    switch m.state {
    case StateStartupPath, StateUninstallStartupPath:
        return m.startupPathModel.View()
    case StateClaudePath, StateUninstallClaudePath:
        return m.claudePathModel.View()
    case StateFileSelection, StateUninstallFileSelection:
        return m.fileSelectionModel.View()
    case StateComplete, StateUninstallComplete:
        return m.completeModel.View()
    case StateError:
        return m.errorModel.View()
    default:
        return "Unknown state"
    }
}
```

**State Transition Logic:**

```go
func (m *MainModel) transitionToState(newState InstallerState) {
    m.state = newState

    // Initialize or reset sub-models
    switch newState {
    case StateStartupPath, StateUninstallStartupPath:
        m.startupPathModel = m.startupPathModel.Reset()

    case StateClaudePath, StateUninstallClaudePath:
        m.claudePathModel = NewClaudePathModelWithMode(m.startupPath, m.mode)

    case StateFileSelection, StateUninstallFileSelection:
        m.fileSelectionModel = NewFileSelectionModelWithMode(
            "claude-code",
            m.claudePath,
            m.installer,
            m.claudeAssets,
            m.startupAssets,
            m.mode,
        )
        m.selectedFiles = m.fileSelectionModel.selectedFiles

    case StateComplete, StateUninstallComplete:
        m.completeModel = NewCompleteModelWithAssets(
            "claude-code",
            m.installer,
            m.mode,
            m.claudeAssets,
            m.startupAssets,
            m.selectedFiles,
        )

    case StateError:
        // Error model is set before transition
    }
}
```

---

### StartupPathModel - Installation Path Selection

**File:** `internal/ui/startup_path_model.go` (275 lines)

```go
type StartupPathModel struct {
    styles          Styles
    renderer        *ProgressiveDisclosureRenderer
    choices         []string
    cursor          int
    selectedPath    string
    ready           bool
    inputMode       bool              // Custom path input mode
    textInput       textinput.Model   // Bubble's text input
    suggestions     []string          // Autocomplete suggestions
    suggestionIndex int
    mode            OperationMode
}
```

**Choices:**
```go
choices := []string{
    "~/.config/the-startup (recommended)",
    ".the-startup (local)",
    "Custom location",
    "Cancel",
}
```

**Keyboard Navigation:**
- `↑` / `k`: Move cursor up
- `↓` / `j`: Move cursor down
- `Enter`: Select choice
- `Esc`: Cancel (when not in input mode)

**Custom Path Input Mode:**
- Triggered by selecting "Custom location"
- Features:
  - Tab: Cycle through autocomplete suggestions
  - Enter: Confirm path
  - Esc: Exit input mode
  - Typing: Reset suggestions

**Path Autocomplete:**
```go
func (m StartupPathModel) getPathSuggestions(input string) []string {
    // 1. Expand ~ if present
    // 2. Get directory and base name
    // 3. If directory doesn't exist, use parent
    // 4. Read directory entries
    // 5. Find directories matching base prefix
    // 6. Convert back to ~ notation if in home directory
    // 7. Limit to 5 suggestions
}
```

**Path Processing:**
- Expand `~/` to home directory
- Ensure path ends with `.the-startup`
- Example: `~/Projects` → `~/Projects/.the-startup`

**View Rendering:**
- Shows banner
- Shows title based on mode (install vs uninstall)
- Shows choices or custom input
- Shows autocomplete suggestions when available
- Shows help text

**Special Value:**
- `"CANCEL"` returned when user selects Cancel option

---

### ClaudePathModel - Claude Directory Selection

**File:** `internal/ui/claude_path_model.go` (287 lines)

**Structure:** Nearly identical to StartupPathModel

```go
type ClaudePathModel struct {
    styles          Styles
    renderer        *ProgressiveDisclosureRenderer
    choices         []string
    cursor          int
    startupPath     string  // From previous step
    selectedPath    string
    ready           bool
    inputMode       bool
    textInput       textinput.Model
    suggestions     []string
    suggestionIndex int
    mode            OperationMode
}
```

**Choices:**
```go
choices := []string{
    "~/.claude (recommended)",
    ".claude (local)",
    "Custom location",
    "Cancel",
}
```

**Differences from StartupPathModel:**
1. Shows previous selection (startup path) in header
2. Ensures path ends with `.claude` instead of `.the-startup`
3. Different help text on ESC: "Escape: back" (not "quit")

**Progressive Disclosure:**
```go
// Shows previous selection at top
s.WriteString(m.renderer.RenderSelectionsWithMode("", displayStartupPath, 0, m.mode))
```

---

### FileSelectionModel - File Tree & Confirmation

**File:** `internal/ui/file_selection_model.go` (726 lines)

```go
type FileSelectionModel struct {
    styles        Styles
    renderer      *ProgressiveDisclosureRenderer
    installer     *installer.Installer
    claudeAssets  *embed.FS
    startupAssets *embed.FS
    selectedTool  string
    selectedPath  string
    selectedFiles []string
    cursor        int
    choices       []string
    ready         bool
    confirmed     bool
    mode          OperationMode
}
```

**Choices (Install):**
```go
choices := []string{
    "Yes, give me awesome",
    "Huh? I did not sign up for this",
}
```

**Choices (Uninstall):**
```go
choices := []string{
    "Yes, remove everything",
    "No, keep everything as-is",
}
```

**File Discovery:**

```go
func (m FileSelectionModel) getAllAvailableFiles() []string {
    var allFiles []string

    // Walk Claude assets
    fs.WalkDir(m.claudeAssets, "assets/claude", func(path, d, err) {
        if !d.IsDir() {
            relPath := strings.TrimPrefix(path, "assets/claude/")
            allFiles = append(allFiles, relPath)
        }
    })

    // Walk Startup assets
    fs.WalkDir(m.startupAssets, "assets/the-startup", func(path, d, err) {
        if !d.IsDir() {
            relPath := strings.TrimPrefix(path, "assets/the-startup/")
            allFiles = append(allFiles, relPath)
        }
    })

    return allFiles
}
```

**Uninstall File Discovery:**

```go
func (m FileSelectionModel) getFilesFromLockfile() []string {
    lockfilePath := filepath.Join(startupPath, "the-startup.lock")

    // Read and parse lockfile
    var lockFile map[string]interface{}
    json.Unmarshal(lockfileData, &lockFile)

    filesMap := lockFile["files"].(map[string]interface{})

    var existingFiles []string
    for lockfilePath := range filesMap {
        var fullPath string

        if strings.HasPrefix(lockfilePath, "startup/") {
            relPath := strings.TrimPrefix(lockfilePath, "startup/")
            fullPath = filepath.Join(startupPath, relPath)
        } else if strings.HasPrefix(lockfilePath, "bin/") {
            fullPath = filepath.Join(startupPath, lockfilePath)
        } else {
            fullPath = filepath.Join(claudePath, lockfilePath)
        }

        // Only include if file exists
        if _, err := os.Stat(fullPath); err == nil {
            existingFiles = append(existingFiles, fullPath)
        }
    }

    // Also include lockfile itself
    existingFiles = append(existingFiles, lockfilePath)

    return existingFiles
}
```

**Tree Rendering (Install):**

The tree uses `github.com/charmbracelet/lipgloss/tree` to render a hierarchical view:

```go
func (m FileSelectionModel) buildStaticTree() string {
    // Build file hierarchy
    type fileNode struct {
        name     string
        children map[string]*fileNode
        isFile   bool
        fullPath string
        exists   bool  // Will be updated?
    }

    // Walk embedded assets and build tree
    fs.WalkDir(embedFS, basePath, func(path, d, err) {
        // Build tree structure
        // Skip non-markdown/json files
        // Validate agent naming convention (must start with "the-")
    })

    // Render tree with special handling:
    // 1. Agent directories: "the-analyst (5 specialized activities)"
    // 2. Command files: Convert "s/specify.md" to "/s:specify"
    // 3. Existing files: Orange "(will update)" styling
    // 4. Deprecated files: Red strikethrough "✗ file (will remove)"

    tree.Root("⁜ ~/.claude").
        Child(
            "agents/",
            agentsTree,
            "commands/",
            commandsTree,
            "output-styles/",
            outputStylesTree,
            settingsItem,
        ).
        Enumerator(tree.RoundedEnumerator)
}
```

**Styling:**
- Item: Pink (#212)
- Update: Orange (#214)
- Remove (deprecated): Red (#196) with strikethrough

**Special File Handling:**
- Agents: Nested directories supported, show count of specialized activities
- Commands: Convert path separators to colons (`/s:specify`)
- Settings: Check if exists for "will update" message
- Deprecated files: Show with ✗ and strike-through

---

### CompleteModel - Success Screen

**File:** `internal/ui/complete_model.go` (381 lines)

```go
type CompleteModel struct {
    styles        Styles
    installer     *installer.Installer
    selectedTool  string
    ready         bool
    mode          OperationMode
    claudeAssets  *embed.FS
    startupAssets *embed.FS
    selectedFiles []string
}

type autoExitMsg struct{}

func (m CompleteModel) Init() tea.Cmd {
    // Exit immediately
    return func() tea.Msg {
        return autoExitMsg{}
    }
}

func (m CompleteModel) Update(msg tea.Msg) (CompleteModel, tea.Cmd) {
    switch msg.(type) {
    case tea.KeyMsg:
        m.ready = true  // Allow immediate exit on any key
    case autoExitMsg:
        m.ready = true  // Auto-exit
    }
    return m, nil
}
```

**View:**
- Shows banner
- Success message: "✅ Installation Complete!" or "✅ Uninstallation Complete!"
- Installation locations with ~ paths
- File tree (same as during selection)
- Deprecated files removed (if any)
- Repository link: https://github.com/rsmdt/the-startup

**Auto-Exit:**
- Sends `autoExitMsg` immediately on Init
- User can press any key to exit faster
- MainModel checks `m.completeModel.Ready()` to call `tea.Quit`

---

### ErrorModel - Error Display

**File:** `internal/ui/error_model.go` (66 lines)

```go
type ErrorModel struct {
    styles     Styles
    renderer   *ProgressiveDisclosureRenderer
    err        error
    errContext string
    ready      bool
}

func (m ErrorModel) Update(msg tea.Msg) (ErrorModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "esc":
            m.ready = true
        }
    }
    return m, nil
}

func (m ErrorModel) View() string {
    s.WriteString(m.renderer.RenderError(m.err, m.errContext))
    s.WriteString(m.styles.Help.Render("Press Escape to go back"))
}
```

**Behavior:**
- Shows error message with context
- Press ESC to return to appropriate starting state
- MainModel checks mode to decide whether to return to StateStartupPath or StateUninstallStartupPath

---

### Theme System

**File:** `internal/ui/theme.go` (130 lines)

```go
type Theme struct {
    Name       string
    Primary    lipgloss.Color  // Brand color, accents
    Success    lipgloss.Color  // Success messages
    Error      lipgloss.Color  // Errors
    Warning    lipgloss.Color  // Warnings
    Info       lipgloss.Color  // Information
    Text       lipgloss.Color  // Regular text
    TextMuted  lipgloss.Color  // Help text
    TextBright lipgloss.Color  // Cursor lines
    Background lipgloss.Color
}

var CharmTheme = Theme{
    Name:       "charm",
    Primary:    lipgloss.Color("#FF06B7"),  // Pink/Magenta
    Success:    lipgloss.Color("#04B575"),  // Green
    Error:      lipgloss.Color("#FF4444"),  // Red
    Warning:    lipgloss.Color("#FFA500"),  // Orange
    Info:       lipgloss.Color("#3C7EFF"),  // Blue
    Text:       lipgloss.Color("#FAFAFA"),  // Light gray
    TextMuted:  lipgloss.Color("#606060"),  // Dark gray
    TextBright: lipgloss.Color("#42FF76"),  // Bright green
    Background: lipgloss.Color("#000000"),  // Black
}

type Styles struct {
    Title      lipgloss.Style  // Bold, green bg, black text
    Success    lipgloss.Style
    Error      lipgloss.Style
    Warning    lipgloss.Style
    Info       lipgloss.Style
    Help       lipgloss.Style  // Muted
    Cursor     lipgloss.Style  // Primary color
    CursorLine lipgloss.Style  // Bright
    Normal     lipgloss.Style
    Selected   lipgloss.Style  // Primary, bold
}
```

**Icons:**
```go
const (
    IconSuccess    = "✓"
    IconError      = "✗"
    IconWarning    = "⚠"
    IconInfo       = "ℹ"
    IconBullet     = "•"
    IconArrow      = "→"
    IconUpdate     = "↻"
    IconSelected   = "●"
    IconUnselected = "○"
    IconFolder     = "📁"
    IconFile       = "📄"
    IconRocket     = "🚀"
)
```

**Migration Note:**
- Use chalk/gradient-string for colors in Ink
- Title style: Use Ink's Box component with custom styling
- Icons: Unicode characters work in TypeScript

---

## Business Logic & Installer

### Installer Core

**File:** `internal/installer/installer.go` (1083 lines)

```go
type Installer struct {
    claudeAssets  *embed.FS
    startupAssets *embed.FS

    installPath     string              // ~/.config/the-startup or .the-startup
    claudePath      string              // ~/.claude or .claude
    tool            string              // "claude-code"
    components      []string            // ["agents", "commands", "templates", "rules"]
    selectedFiles   []string            // Specific files to install
    existingLock    *config.LockFile    // Previously installed files
    deprecatedFiles []string            // Files to be removed
}
```

**Default Paths:**
```go
func New(claudeAssets, startupAssets *embed.FS) *Installer {
    homeDir, _ := os.UserHomeDir()

    return &Installer{
        installPath:   ".the-startup",  // Default to local
        claudePath:    filepath.Join(homeDir, ".claude"),
        tool:          "claude-code",
        components:    []string{"agents", "commands", "templates", "rules"},
    }
}
```

**Path Expansion:**
```go
func (i *Installer) SetInstallPath(path string) {
    // Expand ~ to home directory
    if strings.HasPrefix(path, "~/") {
        homeDir, _ := os.UserHomeDir()
        path = filepath.Join(homeDir, path[2:])
    }
    i.installPath = path
}
```

**Installation Flow:**

```go
func (i *Installer) Install() error {
    // 1. Create installation directory
    os.MkdirAll(i.installPath, 0755)

    // 2. Load existing lock file
    i.LoadExistingLockFile()

    // 3. Detect deprecated files
    if i.existingLock != nil {
        i.GetDeprecatedFiles()
    }

    // 4. Remove deprecated files
    if len(i.deprecatedFiles) > 0 {
        i.RemoveDeprecatedFiles()
    }

    // 5. Install Claude assets to CLAUDE_PATH
    i.installClaudeAssets()

    // 6. Install Startup assets to STARTUP_PATH
    i.installStartupAssets()

    // 7. Install binary to STARTUP_PATH/bin
    i.installBinary()

    // 8. Create logs directory
    os.MkdirAll(filepath.Join(i.installPath, "logs"), 0755)

    // 9. Configure settings.json and settings.local.json
    i.configureSettings()

    // 10. Create lock file
    i.createLockFile()

    return nil
}
```

**Claude Assets Installation:**

```go
func (i *Installer) installClaudeAssets() error {
    return fs.WalkDir(i.claudeAssets, "assets/claude", func(path, d, err) {
        if d.IsDir() {
            return nil
        }

        relPath := strings.TrimPrefix(path, "assets/claude/")

        // Filter by selectedFiles if specified
        if len(i.selectedFiles) > 0 {
            found := false
            for _, selected := range i.selectedFiles {
                if selected == relPath {
                    found = true
                    break
                }
            }
            if !found {
                return nil
            }
        }

        destPath := filepath.Join(i.claudePath, relPath)

        // Create directory
        os.MkdirAll(filepath.Dir(destPath), 0755)

        // Read and process file
        data, _ := i.claudeAssets.ReadFile(path)
        data = i.replacePlaceholders(data)

        // Skip settings files (handled separately)
        if filepath.Base(path) == "settings.json" || filepath.Base(path) == "settings.local.json" {
            return nil
        }

        // Write file
        os.WriteFile(destPath, data, 0644)

        return nil
    })
}
```

**Startup Assets Installation:**

Nearly identical to `installClaudeAssets`, but:
- Walks `assets/the-startup` instead
- Writes to `i.installPath` instead of `i.claudePath`

**Binary Installation:**

```go
func (i *Installer) installBinary() error {
    binDir := filepath.Join(i.installPath, "bin")
    os.MkdirAll(binDir, 0755)

    return i.copyCurrentExecutable(binDir)
}

func (i *Installer) copyCurrentExecutable(destDir string) error {
    // Get current executable path
    execPath, _ := os.Executable()

    src, _ := os.Open(execPath)
    defer src.Close()

    destPath := filepath.Join(destDir, "the-startup")
    dst, _ := os.Create(destPath)
    defer dst.Close()

    io.Copy(dst, src)

    // Make executable
    os.Chmod(destPath, 0755)

    return nil
}
```

**Migration Note:**
- In npm package, use `pkg.files` to bundle binary
- Or provide install script that downloads platform-specific binary
- Binary path: `node_modules/.bin/the-startup` (symlinked)

---

### Settings Configuration

**File:** `internal/installer/installer.go` (lines 529-617)

**Purpose:** Merge template settings with existing settings.json

```go
func (i *Installer) configureSettings() error {
    settingsPath := filepath.Join(i.claudePath, "settings.json")

    // 1. Read template settings from assets/claude/settings.json
    var templateSettings map[string]interface{}
    templateData, _ := i.claudeAssets.ReadFile("assets/claude/settings.json")
    templateData = i.replacePlaceholders(templateData)
    json.Unmarshal(templateData, &templateSettings)

    // 2. If template loading failed, use minimal default
    if templateSettings == nil {
        templateSettings = map[string]interface{}{
            "permissions": map[string]interface{}{
                "additionalDirectories": []string{i.installPath},
            },
            "statusLine": map[string]interface{}{
                "type":    "command",
                "command": filepath.Join(i.installPath, "bin", "the-startup") + " statusline",
            },
        }
    }

    // 3. Read existing settings if present
    var existingSettings map[string]interface{}
    if data, err := os.ReadFile(settingsPath); err == nil {
        json.Unmarshal(data, &existingSettings)
    }

    // 4. Merge (template takes precedence)
    settings := i.mergeSettings(existingSettings, templateSettings)

    // 5. Write updated settings
    data, _ := json.MarshalIndent(settings, "", "  ")
    os.WriteFile(settingsPath, data, 0644)

    // 6. Handle settings.local.json if template exists
    // ... similar process for settings.local.json

    return nil
}
```

**Merge Logic:**

```go
func (i *Installer) mergeSettings(existing, template map[string]interface{}) map[string]interface{} {
    if existing == nil {
        return template
    }
    if template == nil {
        return existing
    }

    result := make(map[string]interface{})

    // Copy all existing settings
    for k, v := range existing {
        result[k] = v
    }

    // Merge template settings
    for key, templateValue := range template {
        existingValue, exists := result[key]

        if !exists {
            result[key] = templateValue
        } else {
            result[key] = i.mergeValues(existingValue, templateValue)
        }
    }

    return result
}

func (i *Installer) mergeValues(existing, template interface{}) interface{} {
    // If both are maps, merge recursively
    if existingMap, ok := existing.(map[string]interface{}); ok {
        if templateMap, ok := template.(map[string]interface{}); ok {
            merged := make(map[string]interface{})

            // Copy existing
            for k, v := range existingMap {
                merged[k] = v
            }

            // Merge template
            for k, templateVal := range templateMap {
                if existingVal, exists := merged[k]; exists {
                    if k == "additionalDirectories" {
                        // Special handling: deduplicate
                        merged[k] = i.mergeAndDeduplicate(existingVal, templateVal)
                    } else {
                        merged[k] = i.mergeValues(existingVal, templateVal)
                    }
                } else {
                    merged[k] = templateVal
                }
            }

            return merged
        }
    }

    // For other types, template takes precedence
    return template
}

func (i *Installer) mergeAndDeduplicate(existing, template interface{}) interface{} {
    if existingSlice, ok := existing.([]interface{}); ok {
        if templateSlice, ok := template.([]interface{}); ok {
            seen := make(map[string]bool)
            var result []interface{}

            // Add template directories first (they have right paths)
            for _, item := range templateSlice {
                if str, ok := item.(string); ok {
                    if !seen[str] {
                        seen[str] = true
                        result = append(result, str)
                    }
                }
            }

            // Add existing directories (skip startup-related)
            for _, item := range existingSlice {
                if str, ok := item.(string); ok {
                    if !strings.Contains(str, "the-startup") && !seen[str] {
                        seen[str] = true
                        result = append(result, str)
                    }
                }
            }

            return result
        }
    }

    return template
}
```

**Deduplication Strategy:**
- Template values take precedence for managed sections
- For `additionalDirectories`: Use template paths, add non-startup existing paths
- This prevents duplicate startup paths when reinstalling

---

### Lockfile Management

**File:** `internal/config/config.go`

```go
type LockFile struct {
    Version     string              `json:"version"`
    InstallDate string              `json:"install_date"`  // RFC3339 format
    InstallPath string              `json:"install_path"`
    ClaudePath  string              `json:"claude_path"`
    Tool        string              `json:"tool"`
    Components  []string            `json:"components"`
    Files       map[string]FileInfo `json:"files"`
}

type FileInfo struct {
    Size         int64  `json:"size"`
    LastModified string `json:"last_modified"`  // RFC3339 format
    Checksum     string `json:"checksum"`       // SHA256 hash
}
```

**Lock File Creation:**

```go
func (i *Installer) createLockFile() error {
    lockFile := &config.LockFile{
        Version:     "1.0.0",
        InstallDate: time.Now().Format(time.RFC3339),
        InstallPath: i.installPath,
        ClaudePath:  i.claudePath,
        Tool:        i.tool,
        Components:  i.components,
        Files:       make(map[string]config.FileInfo),
    }

    // Record Claude assets
    fs.WalkDir(i.claudeAssets, "assets/claude", func(path, d, err) {
        if d.IsDir() {
            return nil
        }

        relPath := strings.TrimPrefix(path, "assets/claude/")

        // Skip if not in selectedFiles
        if len(i.selectedFiles) > 0 && !contains(i.selectedFiles, relPath) {
            return nil
        }

        // Skip settings files
        if filepath.Base(path) == "settings.json" || filepath.Base(path) == "settings.local.json" {
            return nil
        }

        destPath := filepath.Join(i.claudePath, relPath)
        if info, err := os.Stat(destPath); err == nil {
            checksum, _ := i.calculateFileChecksum(destPath)

            lockFile.Files[relPath] = config.FileInfo{
                Size:         info.Size(),
                LastModified: info.ModTime().Format(time.RFC3339),
                Checksum:     checksum,
            }
        }

        return nil
    })

    // Record Startup assets
    fs.WalkDir(i.startupAssets, "assets/the-startup", func(path, d, err) {
        // Similar, but store with "startup/" prefix
        lockFile.Files["startup/"+relPath] = config.FileInfo{...}
    })

    // Record binary
    binPath := filepath.Join(i.installPath, "bin", "the-startup")
    if info, err := os.Stat(binPath); err == nil {
        checksum, _ := i.calculateFileChecksum(binPath)
        lockFile.Files["bin/the-startup"] = config.FileInfo{...}
    }

    // Write lock file
    lockFilePath := filepath.Join(i.installPath, "the-startup.lock")
    data, _ := json.MarshalIndent(lockFile, "", "  ")
    os.WriteFile(lockFilePath, data, 0644)

    return nil
}

func (i *Installer) calculateFileChecksum(filePath string) (string, error) {
    file, _ := os.Open(filePath)
    defer file.Close()

    hasher := sha256.New()
    io.Copy(hasher, file)

    return hex.EncodeToString(hasher.Sum(nil)), nil
}
```

**Lock File Key Format:**
- Claude files: `agents/the-chief.md`, `commands/s/specify.md`
- Startup files: `startup/templates/product-requirements.md`
- Binary: `bin/the-startup`

**Migration Note:**
- Use Node's `crypto` module for SHA256 checksums
- Store lock file as JSON in same location

---

### Deprecated File Detection

```go
func (i *Installer) GetDeprecatedFiles() []string {
    if i.existingLock == nil {
        return nil
    }

    // Build set of current files from embedded assets
    currentFiles := make(map[string]bool)

    fs.WalkDir(i.claudeAssets, "assets/claude", func(path, d, err) {
        if !d.IsDir() {
            relPath := strings.TrimPrefix(path, "assets/claude/")
            currentFiles[relPath] = true
        }
    })

    // Check which lock file entries no longer exist in assets
    var deprecated []string
    for filePath := range i.existingLock.Files {
        // Check Claude files
        if strings.HasPrefix(filePath, "agents/") || strings.HasPrefix(filePath, "commands/") {
            if !currentFiles[filePath] {
                deprecated = append(deprecated, filePath)
            }
            continue
        }

        // Check startup files (templates, rules)
        if strings.HasPrefix(filePath, "startup/templates/") {
            relPath := strings.TrimPrefix(filePath, "startup/")
            assetPath := "assets/the-startup/" + relPath

            if _, err := i.startupAssets.ReadFile(assetPath); err != nil {
                deprecated = append(deprecated, filePath)
            }
        }
    }

    i.deprecatedFiles = deprecated
    return deprecated
}

func (i *Installer) RemoveDeprecatedFiles() error {
    for _, relPath := range i.deprecatedFiles {
        var fullPath string

        if strings.HasPrefix(relPath, "agents/") || strings.HasPrefix(relPath, "commands/") {
            fullPath = filepath.Join(i.claudePath, relPath)
        } else if strings.HasPrefix(relPath, "startup/") {
            fullPath = filepath.Join(i.installPath, strings.TrimPrefix(relPath, "startup/"))
        } else {
            continue
        }

        os.Remove(fullPath)  // Ignore errors
    }

    return nil
}
```

**Use Case:**
- Agent renamed: `the-old-agent.md` → `the-new-agent.md`
- Lock file still has `agents/the-old-agent.md`
- Detection finds it's not in current assets
- Removal deletes it from `~/.claude/agents/the-old-agent.md`

---

### Placeholder Replacement

```go
func (i *Installer) replacePlaceholders(data []byte) []byte {
    // Convert paths to ~ format
    startupPath := toTildePath(i.installPath)
    claudePath := toTildePath(i.claudePath)

    data = bytes.ReplaceAll(data, []byte("{{STARTUP_PATH}}"), []byte(startupPath))
    data = bytes.ReplaceAll(data, []byte("{{CLAUDE_PATH}}"), []byte(claudePath))

    return data
}

func toTildePath(path string) string {
    homeDir, _ := os.UserHomeDir()

    if strings.HasPrefix(path, homeDir) {
        return "~" + strings.TrimPrefix(path, homeDir)
    }

    return path
}
```

**Applied To:**
- All files from Claude assets
- All files from Startup assets
- Settings templates

**Example:**
```json
// Template
{
  "permissions": {
    "additionalDirectories": ["{{STARTUP_PATH}}"]
  },
  "statusLine": {
    "command": "{{STARTUP_PATH}}/bin/the-startup statusline"
  }
}

// After replacement
{
  "permissions": {
    "additionalDirectories": ["~/.config/the-startup"]
  },
  "statusLine": {
    "command": "~/.config/the-startup/bin/the-startup statusline"
  }
}
```

---

## Embedded Asset System

### Go Embed Directive

**File:** `main.go`

```go
//go:embed assets/claude
var claudeAssets embed.FS

//go:embed assets/the-startup
var startupAssets embed.FS
```

**Behavior:**
- Assets bundled at compile time
- Read-only filesystem
- Supports recursive directories
- Access via `embed.FS` interface

**Asset Structure:**

```
assets/
├── claude/
│   ├── agents/
│   │   ├── the-chief.md
│   │   ├── the-meta-agent.md
│   │   ├── the-analyst/
│   │   │   ├── requirements-analysis.md
│   │   │   ├── feature-prioritization.md
│   │   │   └── project-coordination.md
│   │   ├── the-architect/
│   │   ├── the-designer/
│   │   ├── the-ml-engineer/
│   │   ├── the-mobile-engineer/
│   │   ├── the-platform-engineer/
│   │   ├── the-qa-engineer/
│   │   ├── the-security-engineer/
│   │   └── the-software-engineer/
│   ├── commands/
│   │   └── s/
│   │       ├── analyze.md
│   │       ├── implement.md
│   │       ├── init.md
│   │       ├── refactor.md
│   │       └── specify.md
│   ├── output-styles/
│   │   └── the-startup.md
│   ├── settings.json
│   └── settings.local.json
└── the-startup/
    ├── templates/
    │   ├── business-requirements.md
    │   ├── definition-of-done.md
    │   ├── definition-of-ready.md
    │   ├── implementation-plan.md
    │   ├── product-requirements.md
    │   ├── solution-design.md
    │   └── task-definition-of-done.md
    └── rules/
        ├── agent-creation-principles.md
        ├── agent-delegation.md
        └── cycle-pattern.md
```

**Access Pattern:**

```go
// Read file
data, err := claudeAssets.ReadFile("assets/claude/agents/the-chief.md")

// Walk directory
fs.WalkDir(claudeAssets, "assets/claude", func(path string, d fs.DirEntry, err error) error {
    if !d.IsDir() {
        // Process file
    }
    return nil
})
```

**Migration to TypeScript:**

**Option 1: pkg.files API**
```typescript
import { readFileSync } from 'fs';
import { fileURLToPath } from 'url';

// pkg automatically embeds files listed in package.json
const assetsDir = path.join(__dirname, 'assets');
```

**Option 2: Webpack/esbuild bundling**
```typescript
import claudeAgentsTheChief from '../assets/claude/agents/the-chief.md';

// Use raw-loader or similar
```

**Option 3: Copy during build**
```json
// package.json
{
  "scripts": {
    "build": "tsc && cp -r assets dist/"
  }
}
```

**Recommendation:** Use pkg.files for npm distribution, with fallback to regular file reads during development.

---

## Error Handling Patterns

### Error Propagation

**Pattern 1: Bubble errors with context**
```go
if err := i.installClaudeAssets(); err != nil {
    return fmt.Errorf("failed to install Claude assets: %w", err)
}
```

**Pattern 2: Log and continue**
```go
if err := i.RemoveDeprecatedFiles(); err != nil {
    fmt.Printf("Warning: Failed to remove some deprecated files: %v\n", err)
    // Continue with installation
}
```

**Pattern 3: Silent failures for compatibility**
```go
// statusline.go - Silent fail for hook compatibility
if err := json.NewDecoder(input).Decode(&data); err != nil {
    return nil  // Silent fail
}
```

**Pattern 4: Best-effort cleanup**
```go
func (m *MainModel) removeStartupDirectoryIfEmpty() {
    if entries, err := os.ReadDir(startupPath); err == nil && len(entries) == 0 {
        os.Remove(startupPath)  // Ignore errors
    }
}
```

### Error Messages

**Format:**
```go
"failed to %s: %w"
"invalid %s: %s (valid: %s)"
"error %s: %v"
```

**Examples:**
```
"failed to create install directory: permission denied"
"invalid template name: foo (valid: definition-of-ready, definition-of-done, task-definition-of-done)"
"error during installation: file not found"
"spec not found"
```

**User-Facing Errors (TUI):**
```go
m.errorModel = NewErrorModel(err, "during installation")
// Renders: "❌ Error during installation: <error message>"
```

**Migration Note:**
- Use Error.cause in TypeScript for wrapping
- Maintain exact error message formats for consistency

---

## Edge Cases & Gotchas

### 1. Agent Naming Convention

**Issue:** Only agents starting with `the-` are valid

```go
// file_selection_model.go:418-426
if prefix == "agents/" && !d.IsDir() {
    if !strings.Contains(relPath, "/") {  // Root-level agent
        name := filepath.Base(path)
        name = strings.TrimSuffix(name, ".md")
        if !strings.HasPrefix(name, "the-") {
            return nil  // Skip
        }
    }
    // Nested agents don't need "the-" prefix
}
```

**Edge Case:** Nested agents (in subdirectories) don't require `the-` prefix

**Migration:** Implement same validation in TypeScript walker

---

### 2. Path Handling - Tilde Expansion

**Issue:** Must expand `~/` in both directions

```go
// Expansion
if strings.HasPrefix(path, "~/") {
    homeDir, _ := os.UserHomeDir()
    path = filepath.Join(homeDir, path[2:])
}

// Contraction
if strings.HasPrefix(path, homeDir) {
    return "~" + strings.TrimPrefix(path, homeDir)
}
```

**Gotcha:** Lockfile stores `~` paths, but file operations need absolute paths

**Migration:** Use `os.homedir()` in Node, implement expand/contract helpers

---

### 3. Settings Merge - Deduplication

**Issue:** Reinstalling can create duplicate `additionalDirectories` entries

**Solution:** Template paths take precedence, filter out old startup paths

```go
// Deduplication logic prevents:
"additionalDirectories": [
    "~/.config/the-startup",
    ".the-startup",  // Old path from previous install
]

// Results in:
"additionalDirectories": [
    "~/.config/the-startup",
]
```

**Migration:** Implement same deduplication in TypeScript

---

### 4. Lock File Keys - Prefix Handling

**Issue:** Different prefixes for different file types

```go
// Claude files: "agents/the-chief.md"
// Startup files: "startup/templates/prd.md"
// Binary: "bin/the-startup"
```

**Gotcha:** When reading lockfile for uninstall, must reconstruct full paths:

```go
if strings.HasPrefix(lockfilePath, "startup/") {
    relPath := strings.TrimPrefix(lockfilePath, "startup/")
    fullPath = filepath.Join(startupPath, relPath)
} else if strings.HasPrefix(lockfilePath, "bin/") {
    fullPath = filepath.Join(startupPath, lockfilePath)
} else {
    fullPath = filepath.Join(claudePath, lockfilePath)
}
```

**Migration:** Maintain same prefix convention in TypeScript

---

### 5. Uninstall - Empty Directory Check

**Issue:** Only remove startup directory if empty

```go
if entries, err := os.ReadDir(startupPath); err == nil && len(entries) == 0 {
    os.Remove(startupPath)
}
```

**Gotcha:** NEVER removes `.claude` directory, only individual files

**Migration:** Implement same safety check in TypeScript

---

### 6. Auto-Exit Timer

**Issue:** Complete screen auto-exits via message

```go
type autoExitMsg struct{}

func (m CompleteModel) Init() tea.Cmd {
    return func() tea.Msg {
        return autoExitMsg{}
    }
}
```

**Gotcha:** User can press any key to exit immediately (faster than timer)

**Migration:** In Ink, use `setTimeout` to trigger exit, allow key override

---

### 7. Spec ID Padding

**Issue:** Always 3-digit padding for spec IDs

```go
nextID := fmt.Sprintf("%03d", highestID+1)
// 9 → "010"
// 99 → "100"
// 999 → "1000" (still works but not 3 digits)
```

**Edge Case:** Works beyond 999 but loses padding

**Migration:** Use `.padStart(3, '0')` in TypeScript

---

### 8. Template Backward Compatibility

**Issue:** Scans for both new and old filenames

```go
// New: "product-requirements.md"
// Old: "PRD.md"

newPath := filepath.Join(specDir, tmpl.filename)
if _, err := os.Stat(newPath); err == nil {
    files[key] = newPath
    continue
}

oldPath := filepath.Join(specDir, fmt.Sprintf("%s.md", shortName))
if _, err := os.Stat(oldPath); err == nil {
    files[key] = oldPath
}
```

**Migration:** Maintain same backward compatibility in TypeScript

---

### 9. Window Size Detection

**Issue:** Multiple fallbacks for terminal width

```go
// 1. COLUMNS env var (most reliable in hooks)
// 2. term.GetSize(os.Stdout.Fd())
// 3. term.GetSize(os.Stderr.Fd())
// 4. Default: 120
```

**Gotcha:** Hooks/scripts may not have proper TTY

**Migration:** Use `process.stdout.columns` with fallback to 120

---

### 10. Settings Files - Skip During Asset Copy

**Issue:** `settings.json` and `settings.local.json` handled separately

```go
if filepath.Base(path) == "settings.json" || filepath.Base(path) == "settings.local.json" {
    return nil  // Skip
}
```

**Reason:** Must merge with existing settings, not overwrite

**Migration:** Implement same skip logic in TypeScript

---

## TypeScript Migration Mapping

### Framework Equivalents

| Go | TypeScript | Notes |
|----|-----------|-------|
| Cobra | Commander.js | CLI framework |
| BubbleTea | Ink | React-based TUI |
| Lipgloss | chalk + gradient-string | Terminal styling |
| Bubbles (textinput) | ink-text-input | Text input component |
| embed.FS | pkg.files or bundling | Asset embedding |
| filepath | path | Path operations |
| os | fs/promises | File system |
| encoding/json | JSON.parse/stringify | JSON handling |
| crypto/sha256 | crypto.createHash('sha256') | Checksums |

### State Machine Translation

**Go (BubbleTea):**
```go
type MainModel struct {
    state InstallerState
    startupPathModel StartupPathModel
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch m.state {
    case StateStartupPath:
        newModel, cmd := m.startupPathModel.Update(msg)
        if newModel.Ready() {
            m.transitionToState(StateClaudePath)
        }
    }
}
```

**TypeScript (Ink + React):**
```typescript
const MainModel: FC = () => {
    const [state, setState] = useState<InstallerState>(InstallerState.StartupPath);
    const [startupPath, setStartupPath] = useState<string>('');

    const handleStartupPathComplete = (path: string) => {
        setStartupPath(path);
        setState(InstallerState.ClaudePath);
    };

    return (
        <>
            {state === InstallerState.StartupPath && (
                <StartupPathModel onComplete={handleStartupPathComplete} />
            )}
            {state === InstallerState.ClaudePath && (
                <ClaudePathModel startupPath={startupPath} onComplete={...} />
            )}
        </>
    );
};
```

### CLI Command Translation

**Go:**
```go
cmd := &cobra.Command{
    Use:   "install",
    Short: "Install The Startup agent system",
    RunE: func(cmd *cobra.Command, args []string) error {
        local, _ := cmd.Flags().GetBool("local")
        return ui.RunMainInstallerWithFlags(..., local, yes)
    },
}
cmd.Flags().BoolP("local", "l", false, "Use local paths")
```

**TypeScript:**
```typescript
import { Command } from 'commander';

const program = new Command();

program
    .command('install')
    .description('Install The Startup agent system')
    .option('-l, --local', 'Use local paths')
    .option('-y, --yes', 'Auto-confirm')
    .action(async (options) => {
        await runMainInstaller(options.local, options.yes);
    });
```

### File Operations Translation

**Go:**
```go
data, err := os.ReadFile(path)
if err != nil {
    return err
}

err = os.WriteFile(destPath, data, 0644)
```

**TypeScript:**
```typescript
import { readFile, writeFile } from 'fs/promises';

try {
    const data = await readFile(path, 'utf-8');
    await writeFile(destPath, data, { mode: 0o644 });
} catch (err) {
    throw new Error(`Failed to process file: ${err.message}`);
}
```

### Embedded Assets Translation

**Go:**
```go
//go:embed assets/claude
var claudeAssets embed.FS

data, _ := claudeAssets.ReadFile("assets/claude/agents/the-chief.md")
```

**TypeScript (pkg):**
```typescript
// package.json
{
  "pkg": {
    "assets": ["assets/**/*"]
  }
}

// Code
import { readFileSync } from 'fs';
import { join } from 'path';

const assetPath = join(__dirname, 'assets/claude/agents/the-chief.md');
const data = readFileSync(assetPath, 'utf-8');
```

### Tree Rendering Translation

**Go (Lipgloss Tree):**
```go
import "github.com/charmbracelet/lipgloss/tree"

t := tree.Root("⁜ ~/.claude").
    Child(
        "agents/",
        agentsTree,
        "commands/",
        commandsTree,
    ).
    Enumerator(tree.RoundedEnumerator)

return t.String()
```

**TypeScript (Ink Tree Component):**
```typescript
import Tree from 'ink-tree';

<Tree
    label="⁜ ~/.claude"
    nodes={[
        { label: 'agents/', nodes: agentNodes },
        { label: 'commands/', nodes: commandNodes },
    ]}
/>
```

Or use custom rendering with `chalk`:
```typescript
import chalk from 'chalk';

const renderTree = (nodes, prefix = '') => {
    return nodes.map((node, i) => {
        const isLast = i === nodes.length - 1;
        const marker = isLast ? '└─' : '├─';
        const line = `${prefix}${marker} ${node.label}`;

        if (node.nodes) {
            const childPrefix = prefix + (isLast ? '  ' : '│ ');
            return [line, ...renderTree(node.nodes, childPrefix)];
        }

        return line;
    }).flat().join('\n');
};
```

---

## Risk Assessment

### High Risk

1. **State Machine Complexity**
   - Risk: Ink's React model differs significantly from BubbleTea's TEA architecture
   - Mitigation: Create adapter layer that mimics BubbleTea's message passing
   - Test: Comprehensive state transition testing

2. **Embedded Asset System**
   - Risk: npm/pkg bundling works differently than Go embed
   - Mitigation: Use pkg.files with runtime fallback to fs.readFile
   - Test: Verify assets accessible in both dev and bundled modes

3. **Settings Merge Logic**
   - Risk: Complex recursive merge with special cases
   - Mitigation: Port logic exactly, add extensive unit tests
   - Test: Test all edge cases (arrays, nested objects, deduplication)

### Medium Risk

4. **Path Handling - Cross-Platform**
   - Risk: Windows vs Unix path differences
   - Mitigation: Use `path` module consistently, test on Windows
   - Test: Validate ~ expansion, path joining on all platforms

5. **Binary Installation**
   - Risk: Copying executable to install directory
   - Mitigation: Use platform-specific binary from npm package
   - Test: Verify executable permissions preserved

6. **Terminal Width Detection**
   - Risk: Different behavior in hooks/scripts vs interactive terminals
   - Mitigation: Use same fallback chain as Go implementation
   - Test: Test in TTY and non-TTY environments

### Low Risk

7. **JSON Parsing/Serialization**
   - Risk: Minimal - TypeScript has native JSON support
   - Mitigation: Add type guards for runtime validation
   - Test: Validate lockfile parsing

8. **File Checksums**
   - Risk: Minimal - Node crypto module is robust
   - Mitigation: Use same SHA256 algorithm
   - Test: Verify checksums match between implementations

9. **Error Messages**
   - Risk: Inconsistent error messages
   - Mitigation: Port exact error message formats
   - Test: Verify error message consistency

### Migration Sequence Recommendation

**Phase 1: Core Infrastructure**
1. Setup TypeScript project with Commander.js
2. Implement asset loading system (pkg.files)
3. Port installer business logic (no TUI)
4. Port lockfile management
5. Unit test all ported logic

**Phase 2: CLI Commands (Non-TUI)**
6. Port `spec` command (pure CLI)
7. Port `init` command (uses stdin/stdout prompts)
8. Port `statusline` command
9. Integration tests for CLI commands

**Phase 3: TUI System**
10. Setup Ink framework
11. Port theme system and styling
12. Port state machine infrastructure
13. Port each TUI model (StartupPath, ClaudePath, FileSelection, etc.)
14. Integration tests for TUI flows

**Phase 4: Final Integration**
15. Port `install` command with TUI
16. Port `uninstall` command with TUI
17. End-to-end testing
18. Cross-platform testing (macOS, Linux, Windows)

---

## Summary

This Go codebase demonstrates well-structured separation of concerns:
- **CLI layer** (Cobra) handles argument parsing and flags
- **TUI layer** (BubbleTea) manages interactive user interface
- **Business logic** (installer) handles file operations and installation
- **State machine** (ui/states.go) orchestrates workflow

Key patterns to preserve in TypeScript migration:
1. Composable state machine with sub-models
2. Progressive disclosure (show previous selections)
3. Placeholder replacement system
4. Settings merge with deduplication
5. Lock file-based installation tracking
6. Deprecated file detection and removal
7. Tree rendering for file visualization
8. Auto-exit on completion
9. Silent failures for hook compatibility

**Critical behaviors that must be maintained:**
- Exact error message formats
- Agent naming convention validation
- Lockfile key prefixes and structure
- Settings merge logic with special cases
- Path expansion/contraction for ~ handling
- Template backward compatibility
- Empty directory removal safety

**Total lines of critical code to migrate (excluding stats):**
- CLI commands: ~1,500 lines
- TUI models: ~2,000 lines
- Installer business logic: ~1,100 lines
- **Total: ~4,600 lines**

This analysis provides complete visibility into all behaviors, edge cases, and patterns needed for successful TypeScript migration with feature parity.

---

## Migration Validation

**Document Version:** 1.0
**Validation Date:** 2025-10-06
**Purpose:** Cross-reference Go codebase analysis against PRD and SDD requirements to validate migration assumptions and identify gaps

This section validates all migration assumptions required for 100% feature parity (excluding stats command).

---

### CLI Flag Compatibility Matrix

**Requirement:** PRD line 344 requires exact CLI flag compatibility between Go and TypeScript versions.

#### Install Command Flags

| Go Flag | Short | Type | Default | Commander.js Equivalent | Status | Notes |
|---------|-------|------|---------|-------------------------|--------|-------|
| `--local` | `-l` | boolean | false | `.option('-l, --local', 'desc')` | ✅ COMPATIBLE | No data type conversion needed |
| `--yes` | `-y` | boolean | false | `.option('-y, --yes', 'desc')` | ✅ COMPATIBLE | No data type conversion needed |

**Source:** `cmd/install.go` lines 27-28

**Commander.js Implementation Pattern:**
```typescript
program
  .command('install')
  .option('-l, --local', 'Use local installation paths (skip path selection)')
  .option('-y, --yes', 'Auto-confirm with recommended paths')
  .action(async (options) => {
    // options.local: boolean
    // options.yes: boolean
  });
```

**Validation:** ✅ All flags map 1:1 with Commander.js boolean options.

---

#### Uninstall Command Flags

| Go Flag | Short | Type | Default | Commander.js Equivalent | Status | Notes |
|---------|-------|------|---------|-------------------------|--------|-------|
| `--dry-run` | - | boolean | false | `.option('--dry-run', 'desc')` | ⚠️ NOT IMPLEMENTED | Go has flag but doesn't use it (TODO comment) |
| `--force` | - | boolean | false | `.option('--force', 'desc')` | ⚠️ NOT IMPLEMENTED | Go has flag but doesn't use it (TODO comment) |
| `--keep-logs` | - | boolean | false | `.option('--keep-logs', 'desc')` | ⚠️ NOT IMPLEMENTED | Go has flag but doesn't use it (TODO comment) |
| `--keep-settings` | - | boolean | false | `.option('--keep-settings', 'desc')` | ⚠️ NOT IMPLEMENTED | Go has flag but doesn't use it (TODO comment) |

**Source:** `cmd/uninstall.go` lines 54-57

**Decision Required:** Should TypeScript version implement these flags or remove them?
- **Option 1:** Implement flags (adds functionality not in Go version)
- **Option 2:** Remove flags entirely (true parity with Go behavior)
- **Recommendation:** Remove flags to maintain true feature parity (Go never uses them)

---

#### Init Command Flags

| Go Flag | Short | Type | Default | Commander.js Equivalent | Status | Notes |
|---------|-------|------|---------|-------------------------|--------|-------|
| `--skip-prompts` | `-s` | boolean | false | `.option('-s, --skip-prompts', 'desc')` | ✅ COMPATIBLE | No conversion needed |
| `--force` | `-f` | boolean | false | `.option('-f, --force', 'desc')` | ✅ COMPATIBLE | No conversion needed |
| `--dry-run` | - | boolean | false | `.option('--dry-run', 'desc')` | ✅ COMPATIBLE | No conversion needed |

**Source:** `cmd/init.go` lines 60-62

**Validation:** ✅ All flags map 1:1, all are actively used in Go implementation.

---

#### Spec Command Flags

| Go Flag | Short | Type | Default | Commander.js Equivalent | Status | Notes |
|---------|-------|------|---------|-------------------------|--------|-------|
| `--read` | - | string | "" | `.option('--read <id>', 'desc')` | ✅ COMPATIBLE | String argument required |
| `--add` | - | string | "" | `.option('--add <type>', 'desc')` | ✅ COMPATIBLE | String argument with validation |

**Source:** `cmd/spec.go` lines 62-63

**Commander.js Implementation Pattern:**
```typescript
program
  .command('spec [description...]')
  .option('--read <id>', 'Read existing specification by ID')
  .option('--add <type>', 'Add template (PRD, SDD, PLAN, BRD)')
  .action(async (description, options) => {
    if (options.read) {
      // Read mode: options.read is string
    } else if (options.add) {
      // Add mode: options.add is string, validate against enum
      if (!['PRD', 'SDD', 'PLAN', 'BRD'].includes(options.add)) {
        throw new Error(`invalid template type: ${options.add}`);
      }
    } else {
      // Create mode: description is string[]
    }
  });
```

**Validation:** ✅ String options map correctly. Add validation for template types.

---

#### Statusline Command

**No flags** - Command receives JSON via stdin.

**Validation:** ✅ No flag compatibility issues.

---

### Complete Error Message Catalog

**Requirement:** Maintain exact error message wording for consistency across Go → TypeScript migration.

#### Error Message Categories

##### 1. Validation Errors (User Input)

| Error Message | Location | Category | Suggested Fix |
|---------------|----------|----------|---------------|
| `"invalid template name: %s (valid: definition-of-ready, definition-of-done, task-definition-of-done)"` | `cmd/init.go:82` | Input validation | "Check template name spelling" |
| `"invalid template type: %s (valid: PRD, SDD, PLAN, BRD)"` | `cmd/spec.go:141` | Input validation | "Use one of: PRD, SDD, PLAN, BRD" |
| `"spec ID required when using --add"` | `cmd/spec.go:47` | Missing argument | "Provide spec ID: the-startup spec 004 --add PRD" |
| `"feature description required"` | `cmd/spec.go:54` | Missing argument | "Provide feature name: the-startup spec user-auth" |

**TypeScript Pattern:**
```typescript
class ValidationError extends Error {
  constructor(message: string, public suggestedFix?: string) {
    super(message);
    this.name = 'ValidationError';
  }
}

// Usage
throw new ValidationError(
  'invalid template type: foo (valid: PRD, SDD, PLAN, BRD)',
  'Use one of: PRD, SDD, PLAN, BRD'
);
```

---

##### 2. File System Errors

| Error Message | Location | Category | Context |
|---------------|----------|----------|---------|
| `"failed to create install directory: %w"` | `internal/installer/installer.go:334` | Permission | Creating `.the-startup/` |
| `"failed to create docs directory: %w"` | `cmd/init.go:118` | Permission | Creating `docs/` for templates |
| `"failed to read lock file: %w"` | `internal/installer/installer.go:212` | File I/O | Reading existing lock file |
| `"failed to parse lock file: %w"` | `internal/installer/installer.go:218` | JSON parsing | Corrupted lock file |
| `"failed to write %s: %w"` | `cmd/init.go:182` | File I/O | Writing template files |
| `"failed to read template %s: %w"` | `cmd/init.go:177` | File I/O | Reading embedded templates |
| `"file already exists: %s"` | `cmd/spec.go:182` | File conflict | Template file exists, need --force |
| `"failed to copy template: %w"` | `cmd/spec.go:186` | File I/O | Copying template to spec directory |

**TypeScript Error Wrapping Pattern:**
```typescript
// Matches Go's fmt.Errorf("%s: %w", context, err) pattern
function wrapError(context: string, originalError: Error): Error {
  const wrapped = new Error(`${context}: ${originalError.message}`);
  wrapped.cause = originalError; // ES2022 Error.cause
  return wrapped;
}

// Usage
try {
  await fs.writeFile(path, content);
} catch (err) {
  throw wrapError(`failed to write ${path}`, err as Error);
}
```

---

##### 3. Installation Errors

| Error Message | Location | Category | Context |
|---------------|----------|----------|---------|
| `"failed to install Claude assets: %w"` | `internal/installer/installer.go:362` | Asset copy | Copying agents/commands to ~/.claude |
| `"failed to install Startup assets: %w"` | `internal/installer/installer.go:370` | Asset copy | Copying templates/rules to .the-startup |
| `"failed to install binary: %w"` | `internal/installer/installer.go:378` | Binary copy | Copying executable to bin/ |
| `"failed to create logs directory: %w"` | `internal/installer/installer.go:386` | Directory creation | Creating .the-startup/logs |
| `"failed to configure settings: %w"` | `internal/installer/installer.go:392` | Settings merge | Updating ~/.claude/settings.json |
| `"failed to create lock file: %w"` | `internal/installer/installer.go:401` | Lock file write | Creating installation manifest |

**Rollback Pattern (critical for installation errors):**
```typescript
async function installWithRollback(options: InstallerOptions) {
  const installedFiles: string[] = [];

  try {
    // Track each successful file copy
    for (const file of filesToInstall) {
      await fs.copyFile(source, dest);
      installedFiles.push(dest); // Track for rollback
    }

    await createLockFile(installedFiles);
  } catch (err) {
    // Rollback: delete all successfully copied files
    for (const file of installedFiles) {
      await fs.unlink(file).catch(() => {}); // Ignore rollback errors
    }
    throw wrapError('failed to install', err as Error);
  }
}
```

---

##### 4. Uninstallation Errors

| Error Message | Location | Category | Context |
|---------------|----------|----------|---------|
| `"no lock file provided"` | `internal/uninstaller/file_remover.go:126` | Missing data | Lock file required for uninstall |
| `"failed to check directory: %w"` | `internal/uninstaller/file_remover.go:98` | File system | Checking if directory empty |
| `"path is not a directory: %s"` | `internal/uninstaller/file_remover.go:102` | Validation | Expected directory, got file |
| `"directory not empty: %s"` | `internal/uninstaller/file_remover.go:112` | Safety check | Preventing deletion of non-empty dirs |
| `"no files to remove"` | `internal/ui/model_uninstall.go:15` | Empty operation | No files selected for removal |

**Safety Check Pattern:**
```typescript
async function removeDirectoryIfEmpty(dirPath: string): Promise<void> {
  try {
    const stat = await fs.stat(dirPath);
    if (!stat.isDirectory()) {
      throw new Error(`path is not a directory: ${dirPath}`);
    }

    const entries = await fs.readdir(dirPath);
    if (entries.length > 0) {
      throw new Error(`directory not empty: ${dirPath}`);
    }

    await fs.rmdir(dirPath);
  } catch (err) {
    throw wrapError('failed to check directory', err as Error);
  }
}
```

---

##### 5. Spec Command Errors

| Error Message | Location | Category | Context |
|---------------|----------|----------|---------|
| `"specification not found: %s"` | `cmd/spec.go:281` | Not found | Spec ID doesn't exist |
| `"failed to scan spec directory: %w"` | `cmd/spec.go:83` | File system | Reading docs/specs/ directory |
| `"failed to find highest spec ID: %w"` | `cmd/spec.go:211` | Numbering | Calculating next spec number |
| `"failed to create spec directory: %w"` | `cmd/spec.go:224` | Directory creation | Creating numbered spec directory |
| `"failed to copy PRD template: %w"` | `cmd/spec.go:231` | Template copy | Initial PRD creation |

---

##### 6. User Cancellation

| Error Message | Location | Category | Context |
|---------------|----------|----------|---------|
| `"cancelled by user"` | `cmd/init.go:129` | User action | User declined overwrite prompt |

**Not an error** - Expected user behavior. Should exit with code 0 (not 1).

---

### BubbleTea → Ink Pattern Mapping

**Requirement:** SDD lines 1073-1086 define Ink component patterns to follow. Map BubbleTea patterns to Ink equivalents.

#### Pattern 1: State Machine with Sub-Models

**BubbleTea (Go):**
```go
// internal/ui/states.go
type InstallerState int

const (
    StateStartupPath InstallerState = iota
    StateClaudePath
    StateFileSelection
    StateComplete
    StateError
)

// Main model delegates to sub-models
type MainModel struct {
    state              InstallerState
    startupPathModel   StartupPathModel
    claudePathModel    ClaudePathModel
    fileSelectionModel FileSelectionModel
    completeModel      CompleteModel
    errorModel         ErrorModel
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch m.state {
    case StateStartupPath:
        newModel, cmd := m.startupPathModel.Update(msg)
        m.startupPathModel = newModel
        if m.startupPathModel.Ready() {
            m.transitionToState(StateClaudePath)
        }
        return m, cmd
    // ... other states
    }
}

func (m MainModel) View() string {
    switch m.state {
    case StateStartupPath:
        return m.startupPathModel.View()
    // ... other states
    }
}
```

**Ink (TypeScript) Equivalent:**
```typescript
// src/ui/install/types.ts
enum InstallerState {
  StartupPath = 'startup-path',
  ClaudePath = 'claude-path',
  FileSelection = 'file-selection',
  Complete = 'complete',
  Error = 'error',
}

// src/ui/install/InstallWizard.tsx
const InstallWizard: FC<Props> = ({ onComplete, onError }) => {
  const [state, setState] = useState<InstallerState>(InstallerState.StartupPath);
  const [startupPath, setStartupPath] = useState<string>('');
  const [claudePath, setClaudePath] = useState<string>('');

  // State transition callback pattern
  const handleStartupPathComplete = (path: string) => {
    setStartupPath(path);
    setState(InstallerState.ClaudePath);
  };

  const handleClaudePathComplete = (path: string) => {
    setClaudePath(path);
    setState(InstallerState.FileSelection);
  };

  // Conditional rendering replaces BubbleTea's View() switch
  return (
    <Box flexDirection="column">
      {state === InstallerState.StartupPath && (
        <StartupPathSelector onComplete={handleStartupPathComplete} />
      )}
      {state === InstallerState.ClaudePath && (
        <ClaudePathSelector
          startupPath={startupPath}
          onComplete={handleClaudePathComplete}
        />
      )}
      {state === InstallerState.FileSelection && (
        <FileTreeSelector
          startupPath={startupPath}
          claudePath={claudePath}
          onComplete={() => setState(InstallerState.Complete)}
        />
      )}
      {state === InstallerState.Complete && (
        <CompletionScreen onExit={onComplete} />
      )}
      {state === InstallerState.Error && (
        <ErrorDisplay error={error} />
      )}
    </Box>
  );
};
```

**Key Differences:**
- BubbleTea: Imperative state machine with explicit transitions
- Ink: Declarative React state with useState + conditional rendering
- BubbleTea: Sub-models are structs with Update/View methods
- Ink: Sub-components are React functional components with props
- BubbleTea: Message passing for state updates
- Ink: Callback props for state updates

**Validation:** ✅ Ink pattern is simpler and more idiomatic for React developers.

---

#### Pattern 2: Progressive Disclosure (Show Previous Selections)

**BubbleTea (Go):**
```go
// internal/ui/states.go:154-219
func (p *ProgressiveDisclosureRenderer) RenderSelectionsWithMode(
    tool, path string,
    filesCount int,
    mode OperationMode,
) string {
    var sections []string

    if tool != "" {
        sections = append(sections, fmt.Sprintf("Tool: %s", tool))
    }

    if path != "" {
        displayPath := path
        if strings.Contains(path, ".config/the-startup") {
            displayPath = "~/.config/the-startup (global)"
        }
        sections = append(sections, fmt.Sprintf("Path: %s", displayPath))
    }

    if filesCount > 0 {
        sections = append(sections, fmt.Sprintf("Files: %d selected", filesCount))
    }

    content := strings.Join(sections, " • ")

    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Render(title+"\n"+content)
}
```

**Ink (TypeScript) Equivalent:**
```typescript
// src/ui/shared/ProgressiveHeader.tsx
interface ProgressiveHeaderProps {
  selections: {
    tool?: string;
    path?: string;
    filesCount?: number;
  };
  mode: 'install' | 'uninstall';
}

const ProgressiveHeader: FC<ProgressiveHeaderProps> = ({ selections, mode }) => {
  const sections: string[] = [];

  if (selections.tool) {
    sections.push(`Tool: ${selections.tool}`);
  }

  if (selections.path) {
    const displayPath = selections.path.includes('.config/the-startup')
      ? '~/.config/the-startup (global)'
      : selections.path;
    sections.push(`Path: ${displayPath}`);
  }

  if (selections.filesCount) {
    sections.push(`Files: ${selections.filesCount} selected`);
  }

  const title = mode === 'install'
    ? 'The (Agentic) Startup Installation'
    : 'The (Agentic) Startup Uninstallation';

  return (
    <Box
      borderStyle="round"
      borderColor="cyan"
      paddingX={1}
      marginBottom={1}
      flexDirection="column"
    >
      <Text bold color="cyan">{title}</Text>
      <Text>{sections.join(' • ')}</Text>
    </Box>
  );
};
```

**Validation:** ✅ Ink's Box component provides equivalent border styling.

---

#### Pattern 3: Input Handling with Autocomplete

**BubbleTea (Go):**
```go
// Custom path input with Tab autocomplete
type StartupPathModel struct {
    inputMode       bool
    textInput       textinput.Model  // From charmbracelet/bubbles
    suggestions     []string
    suggestionIndex int
}

func (m StartupPathModel) Update(msg tea.Msg) (StartupPathModel, tea.Cmd) {
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        switch keyMsg.String() {
        case "tab":
            if len(m.suggestions) > 0 {
                m.suggestionIndex = (m.suggestionIndex + 1) % len(m.suggestions)
                m.textInput.SetValue(m.suggestions[m.suggestionIndex])
            }
        case "enter":
            m.ready = true
            m.selectedPath = m.textInput.Value()
        }
    }

    // Update suggestions when typing
    if m.textInput.Value() != previousValue {
        m.suggestions = m.getPathSuggestions(m.textInput.Value())
        m.suggestionIndex = 0
    }

    return m, nil
}
```

**Ink (TypeScript) Equivalent:**
```typescript
// src/ui/shared/PathInput.tsx
import TextInput from 'ink-text-input';

const PathInput: FC<Props> = ({ onSubmit }) => {
  const [value, setValue] = useState('');
  const [suggestions, setSuggestions] = useState<string[]>([]);
  const [selectedIndex, setSelectedIndex] = useState(0);

  useInput((input, key) => {
    if (key.tab && suggestions.length > 0) {
      // Tab: cycle through suggestions
      const nextIndex = (selectedIndex + 1) % suggestions.length;
      setSelectedIndex(nextIndex);
      setValue(suggestions[nextIndex]);
    } else if (key.return) {
      // Enter: submit
      onSubmit(value);
    }
  });

  // Update suggestions when value changes
  useEffect(() => {
    const newSuggestions = getPathSuggestions(value);
    setSuggestions(newSuggestions);
    setSelectedIndex(0);
  }, [value]);

  return (
    <Box flexDirection="column">
      <Text>Path: </Text>
      <TextInput value={value} onChange={setValue} />
      {suggestions.length > 0 && (
        <Box marginTop={1}>
          <Text dimColor>Suggestions: {suggestions.join(', ')}</Text>
        </Box>
      )}
    </Box>
  );
};
```

**Validation:** ✅ Ink's `useInput` hook + `TextInput` component provide equivalent functionality.

---

#### Pattern 4: Keyboard Navigation (Arrow Keys + Vim Bindings)

**BubbleTea (Go):**
```go
func (m StartupPathModel) Update(msg tea.Msg) (StartupPathModel, tea.Cmd) {
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        switch keyMsg.String() {
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }
        case "enter", " ":
            m.selectedPath = m.choices[m.cursor]
            m.ready = true
        case "esc":
            m.selectedPath = "CANCEL"
            m.ready = true
        }
    }
    return m, nil
}
```

**Ink (TypeScript) Equivalent:**
```typescript
// src/ui/shared/MenuSelector.tsx
const MenuSelector: FC<Props> = ({ choices, onSelect }) => {
  const [cursor, setCursor] = useState(0);

  useInput((input, key) => {
    if (key.upArrow || input === 'k') {
      setCursor(prev => Math.max(0, prev - 1));
    } else if (key.downArrow || input === 'j') {
      setCursor(prev => Math.min(choices.length - 1, prev + 1));
    } else if (key.return || input === ' ') {
      onSelect(choices[cursor]);
    } else if (key.escape) {
      onSelect('CANCEL');
    }
  });

  return (
    <Box flexDirection="column">
      {choices.map((choice, i) => (
        <Text key={i} color={i === cursor ? 'cyan' : 'white'}>
          {i === cursor ? '> ' : '  '}{choice}
        </Text>
      ))}
    </Box>
  );
};
```

**Validation:** ✅ Ink's `useInput` provides full keyboard event handling.

---

#### Pattern 5: Spinner + Progress Indicators

**BubbleTea (Go):**
```go
import "github.com/charmbracelet/bubbles/spinner"
import "github.com/charmbracelet/bubbles/progress"

type Model struct {
    spinner  spinner.Model
    progress progress.Model
}

func (m Model) View() string {
    return m.spinner.View() + " Installing files...\n" + m.progress.View()
}
```

**Ink (TypeScript) Equivalent:**
```typescript
// src/ui/shared/InstallProgress.tsx
import Spinner from 'ink-spinner';

const InstallProgress: FC<Props> = ({ current, total }) => {
  const percentage = Math.round((current / total) * 100);

  return (
    <Box flexDirection="column">
      <Box>
        <Text color="green">
          <Spinner type="dots" />
        </Text>
        <Text> Installing files...</Text>
      </Box>
      <Box marginTop={1}>
        <Text>{percentage}% ({current}/{total})</Text>
      </Box>
    </Box>
  );
};
```

**Validation:** ✅ `ink-spinner` provides equivalent loading indicators.

---

#### Pattern 6: Auto-Exit on Completion

**BubbleTea (Go):**
```go
// internal/ui/complete_model.go
type autoExitMsg struct{}

func (m CompleteModel) Init() tea.Cmd {
    return func() tea.Msg {
        return autoExitMsg{}  // Trigger immediate exit
    }
}

func (m CompleteModel) Update(msg tea.Msg) (CompleteModel, tea.Cmd) {
    switch msg.(type) {
    case tea.KeyMsg:
        m.ready = true  // User can exit on any key
    case autoExitMsg:
        m.ready = true  // Auto-exit
    }
    return m, nil
}
```

**Ink (TypeScript) Equivalent:**
```typescript
// src/ui/install/Complete.tsx
const CompletionScreen: FC<Props> = ({ onExit }) => {
  // Auto-exit after brief delay
  useEffect(() => {
    const timer = setTimeout(() => {
      onExit();
    }, 100); // Immediate exit (100ms for render)

    return () => clearTimeout(timer);
  }, [onExit]);

  // Also exit on any key
  useInput(() => {
    onExit();
  });

  return (
    <Box flexDirection="column">
      <Text color="green">✅ Installation Complete!</Text>
      <Text dimColor>Press any key to exit...</Text>
    </Box>
  );
};
```

**Validation:** ✅ Ink's `useEffect` + `useInput` provide equivalent behavior.

---

### BubbleTea Features WITHOUT Direct Ink Equivalents

#### Missing Feature 1: Lipgloss Tree Component

**BubbleTea Has:** `github.com/charmbracelet/lipgloss/tree` for hierarchical file tree rendering

**Ink Alternative:** Build custom tree renderer with Box components

**Workaround:**
```typescript
// src/ui/shared/FileTree.tsx
const TreeNode: FC<NodeProps> = ({ node, depth = 0, isLast }) => {
  const prefix = depth === 0 ? '' : isLast ? '└─ ' : '├─ ';
  const indent = '  '.repeat(depth);

  return (
    <Box flexDirection="column">
      <Text>
        {indent}{prefix}{node.name}
      </Text>
      {node.children?.map((child, i) => (
        <TreeNode
          key={i}
          node={child}
          depth={depth + 1}
          isLast={i === node.children!.length - 1}
        />
      ))}
    </Box>
  );
};
```

**Validation:** ⚠️ Manual tree rendering required, but achievable with Box components.

---

#### Missing Feature 2: Message-Based Architecture

**BubbleTea Has:** Explicit message passing for async operations
```go
type fetchDataMsg struct { data string }

func fetchDataCmd() tea.Cmd {
    return func() tea.Msg {
        data := fetchFromAPI()
        return fetchDataMsg{data}
    }
}
```

**Ink Alternative:** Use React state + async/await + useEffect

**Workaround:**
```typescript
const Component: FC = () => {
  const [data, setData] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      const result = await fetchFromAPI();
      setData(result);
    };
    fetchData();
  }, []);

  // ... render with data
};
```

**Validation:** ✅ React patterns provide equivalent async handling (more idiomatic for TypeScript developers).

---

### TUI State Transition Documentation

**Source:** `internal/ui/states.go` lines 33-134

#### Install Workflow State Transition Table

| Current State | Trigger | Next State | Validation | Notes |
|---------------|---------|------------|------------|-------|
| `StateStartupPath` | User selects path + Enter | `StateClaudePath` | Path is valid and writable | - |
| `StateStartupPath` | User presses ESC or Q | *EXIT* | - | Quit from initial state |
| `StateClaudePath` | User selects path + Enter | `StateFileSelection` | Path is valid | Shows previous: startup path |
| `StateClaudePath` | User presses ESC | `StateStartupPath` | - | Backward navigation |
| `StateFileSelection` | User confirms selection | `StateComplete` OR `StateError` | Installation succeeds or fails | - |
| `StateFileSelection` | User presses ESC | `StateClaudePath` | - | Backward navigation |
| `StateComplete` | Auto-exit msg received | *EXIT* | - | Immediate exit on completion |
| `StateComplete` | User presses any key | *EXIT* | - | Manual exit override |
| `StateError` | User presses ESC | `StateStartupPath` | - | Return to start to retry |

**Valid Transitions (enforced by ValidTransitions map):**
```go
// Forward
StateStartupPath → StateClaudePath → StateFileSelection → StateComplete
                                    → StateError

// Backward (ESC key)
StateClaudePath ← StateStartupPath
StateFileSelection ← StateClaudePath

// Error recovery
StateError → StateStartupPath
```

**Invalid Transitions (would violate state machine):**
- Cannot skip states (e.g., StateStartupPath → StateFileSelection)
- Cannot go backward from StateComplete or StateError (except error recovery)
- Cannot transition to StateStartupPath from StateComplete (must exit)

---

#### Uninstall Workflow State Transition Table

| Current State | Trigger | Next State | Validation | Notes |
|---------------|---------|------------|------------|-------|
| `StateUninstallStartupPath` | Lock file discovered | `StateUninstallClaudePath` | Lock file is valid JSON | - |
| `StateUninstallStartupPath` | No lock file | `StateError` | - | Cannot uninstall without lock file |
| `StateUninstallStartupPath` | User presses ESC | *EXIT* | - | Quit from initial state |
| `StateUninstallClaudePath` | User confirms paths | `StateUninstallFileSelection` | Both paths exist | - |
| `StateUninstallClaudePath` | User presses ESC | `StateUninstallStartupPath` | - | Backward navigation |
| `StateUninstallFileSelection` | User confirms removal | `StateUninstallComplete` OR `StateError` | Removal succeeds or fails | - |
| `StateUninstallFileSelection` | User presses ESC | `StateUninstallClaudePath` | - | Backward navigation |
| `StateUninstallComplete` | Auto-exit msg received | *EXIT* | - | Immediate exit |
| `StateUninstallComplete` | User presses any key | *EXIT* | - | Manual exit |
| `StateError` | User presses ESC | `StateUninstallStartupPath` | - | Return to start |

**Valid Transitions:**
```go
// Forward
StateUninstallStartupPath → StateUninstallClaudePath → StateUninstallFileSelection → StateUninstallComplete
                         → StateError                                            → StateError

// Backward (ESC)
StateUninstallClaudePath ← StateUninstallStartupPath
StateUninstallFileSelection ← StateUninstallClaudePath

// Error recovery
StateError → StateUninstallStartupPath
```

---

#### State Transition Triggers

| Trigger Type | Go Implementation | Ink Implementation | Notes |
|--------------|-------------------|-------------------|-------|
| **User Input** | `tea.KeyMsg` | `useInput` hook | Arrow keys, Enter, ESC, Q |
| **Validation Success** | `model.Ready()` check | State update callback | Sub-model signals completion |
| **Async Operation** | `tea.Cmd` message | `useEffect` + async/await | File operations, settings merge |
| **Auto-Transition** | Immediate in `Update()` | `useEffect` with dependency | Complete screen auto-exit |
| **Error** | Return error, check in parent | `try/catch`, update error state | Installation/uninstallation failures |

---

#### Complex State Logic: Flag-Based State Skipping

**Scenario:** `install --local --yes` should skip all TUI states and install directly.

**BubbleTea (Go):**
```go
// internal/ui/model_install.go (lines 40-70 from analysis)
func (m *MainModel) Init() tea.Cmd {
    // Handle --yes flag: Auto-install without TUI
    if m.flags.Yes {
        return m.handleAutoInstall()
    }

    // Handle --local flag: Pre-select local paths but show TUI
    if m.flags.Local {
        return m.handleLocalFlag()
    }

    return nil
}

func (m *MainModel) handleAutoInstall() tea.Cmd {
    // Determine paths based on --local flag
    if m.flags.Local {
        startupPath = ".the-startup"
        claudePath = ".claude"
    } else {
        startupPath = "~/.config/the-startup"
        claudePath = "~/.claude"
    }

    // Set paths and skip to installation
    m.startupPath = startupPath
    m.claudePath = claudePath
    m.transitionToState(StateFileSelection)

    // Perform installation synchronously
    err := m.installer.Install()
    if err != nil {
        m.transitionToState(StateError)
    } else {
        m.transitionToState(StateComplete)
    }
    return m.completeModel.Init()
}
```

**Ink (TypeScript):**
```typescript
// src/cli/install.ts
program
  .command('install')
  .option('-l, --local', 'Use local paths')
  .option('-y, --yes', 'Auto-confirm')
  .action(async (options) => {
    if (options.yes) {
      // Non-interactive mode: skip TUI entirely
      const paths = options.local
        ? { startup: '.the-startup', claude: '.claude' }
        : { startup: '~/.config/the-startup', claude: '~/.claude' };

      const installer = new Installer(paths.startup, paths.claude);

      try {
        await installer.install();
        console.log('✅ Installation complete!');
        process.exit(0);
      } catch (err) {
        console.error(`❌ Installation failed: ${err.message}`);
        process.exit(1);
      }
    } else {
      // Interactive mode: launch Ink TUI
      const { render } = await import('ink');
      const { InstallWizard } = await import('../ui/install/InstallWizard');

      const { waitUntilExit } = render(<InstallWizard options={options} />);
      await waitUntilExit();
    }
  });
```

**Validation:** ✅ Ink implementation is cleaner - conditionally render TUI vs direct CLI output.

---

### Identified Gaps and Risks

#### Gap 1: Uninstall Flags Not Implemented in Go

**Finding:** Go defines `--dry-run`, `--force`, `--keep-logs`, `--keep-settings` flags but never uses them (TODO comments in code).

**Risk Level:** LOW (but requires decision)

**Options:**
1. ✅ **Remove flags from TypeScript** (true parity with Go behavior)
2. ❌ Implement flags in TypeScript (adds features Go doesn't have)

**Recommendation:** Remove flags to maintain 100% feature parity with Go behavior.

---

#### Gap 2: Lipgloss Tree Rendering

**Finding:** Go uses `github.com/charmbracelet/lipgloss/tree` for hierarchical file trees. Ink has no direct equivalent.

**Risk Level:** MEDIUM

**Workaround:** Build custom tree renderer with Box components and Unicode characters.

**Implementation Effort:** 50-100 lines of TypeScript (low complexity).

---

#### Gap 3: Placeholder Replacement Must Happen Before Lock File Creation

**Finding:** Go replaces `{{STARTUP_PATH}}` and `{{CLAUDE_PATH}}` in settings.json during installation. Lock file must track absolute paths (not placeholders).

**Risk Level:** LOW (well-documented in analysis)

**Validation:** ✅ TypeScript implementation must follow same pattern:
1. Copy assets to disk
2. Replace placeholders in settings.json with absolute paths
3. Create lock file with absolute paths (for uninstall)

---

#### Gap 4: Checksum Algorithm Must Match

**Finding:** Go uses SHA-256 checksums in lock file v2 format. TypeScript must use identical algorithm.

**Risk Level:** LOW (crypto module provides SHA-256)

**Validation:**
```typescript
import { createHash } from 'crypto';

function calculateChecksum(filePath: string): string {
  const content = fs.readFileSync(filePath);
  return createHash('sha256').update(content).digest('hex');
}
```

✅ Node.js crypto module produces identical SHA-256 hashes to Go's crypto/sha256.

---

#### Gap 5: Home Directory Expansion (~) Must Be Cross-Platform

**Finding:** Go expands `~/` to home directory. Must handle Windows (`%USERPROFILE%`) and Unix (`$HOME`).

**Risk Level:** LOW (os.homedir() handles this)

**Validation:**
```typescript
import { homedir } from 'os';

function expandTilde(path: string): string {
  if (path.startsWith('~/')) {
    return path.replace('~', homedir());
  }
  return path;
}
```

✅ Node.js `os.homedir()` is cross-platform.

---

### Migration Validation Summary

#### CLI Flag Compatibility: ✅ COMPLETE

- Install flags: 2/2 compatible (--local, --yes)
- Uninstall flags: 0/4 implemented in Go (decision: remove from TypeScript)
- Init flags: 3/3 compatible (--skip-prompts, --force, --dry-run)
- Spec flags: 2/2 compatible (--read, --add)
- Statusline: No flags

**Total:** 7/7 functional flags mapped to Commander.js

---

#### Error Message Catalog: ✅ COMPLETE

- Validation errors: 4 messages cataloged
- File system errors: 8 messages cataloged
- Installation errors: 6 messages cataloged
- Uninstallation errors: 5 messages cataloged
- Spec errors: 5 messages cataloged
- User cancellation: 1 message (not an error)

**Total:** 29 error messages documented with TypeScript patterns

---

#### BubbleTea → Ink Pattern Translation: ✅ COMPLETE

- State machine pattern: ✅ Mapped (enum + useState)
- Progressive disclosure: ✅ Mapped (Box component)
- Input with autocomplete: ✅ Mapped (useInput + TextInput)
- Keyboard navigation: ✅ Mapped (useInput with vim bindings)
- Spinner/progress: ✅ Mapped (ink-spinner)
- Auto-exit: ✅ Mapped (useEffect + setTimeout)

**Total:** 6/6 core patterns mapped

**Missing Ink equivalents:** 1 (Lipgloss tree - workaround documented)

---

#### State Transition Documentation: ✅ COMPLETE

- Install workflow: 9 states documented with triggers
- Uninstall workflow: 9 states documented with triggers
- Valid transitions: All mapped from Go's ValidTransitions map
- Complex logic: Flag-based state skipping documented

**Total:** 18 state transitions formalized

---

#### Identified Gaps: 5 TOTAL

1. ✅ **Uninstall flags** - Decision: Remove (true parity)
2. ✅ **Tree rendering** - Workaround: Custom Box renderer
3. ✅ **Placeholder replacement** - Pattern documented
4. ✅ **Checksum algorithm** - Validated: SHA-256 identical
5. ✅ **Home directory expansion** - Validated: os.homedir() cross-platform

**Risk Assessment:** All gaps have documented solutions. No migration blockers identified.

---

### Conclusion

**Migration Validation Status:** ✅ COMPLETE

All assumptions validated. No gaps prevent 100% feature parity with Go version (excluding stats command).

**Next Steps:**
1. Implement TypeScript CLI commands following documented patterns
2. Build Ink TUI components using mapped patterns
3. Port business logic with documented error messages
4. Test all state transitions against documented table
5. Validate cross-platform behavior (macOS, Linux, Windows)

**Success Criteria Met:**
- ✅ All CLI flags verified with TypeScript equivalents
- ✅ All error messages cataloged for consistency
- ✅ BubbleTea → Ink translation guide complete
- ✅ State transitions formally documented
- ✅ No migration gaps identified
