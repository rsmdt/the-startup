# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

### Development Mode (No Build Required)

Run commands directly from TypeScript source using tsx:

```bash
# Run installer in development mode
npm run dev:install

# Run other commands
npm run dev:uninstall
npm run dev:init
npm run dev:spec

# Run with custom arguments
npm run dev:run -- install --yes
npm run dev:run -- spec "Add user authentication"

# Or use tsx directly
npx tsx src/index.ts install
```

### Build and Run

```bash
# Build the project
npm run build

# Run the built CLI directly
node dist/index.js install
node dist/index.js --help

# Or link globally for easier testing (simulates npx behavior)
npm link
the-agentic-startup install
the-agentic-startup --help

# Unlink when done testing
npm unlink -g the-agentic-startup

# Watch mode (rebuilds on file changes)
npm run dev
# Then in another terminal, run:
node dist/index.js install
```

### Testing
```bash
# Run all tests
npm test

# Run tests with coverage
npm run test:coverage

# Run tests in watch mode
npm run test:watch

# Run specific test file
npm test tests/core/installer/Installer.test.ts

# Run tests with UI (if available)
npm run test:ui
```

### Development Workflow
```bash
# Format code
npm run format

# Lint code
npm run lint

# Type check
npm run typecheck

# Clean build artifacts
npm run clean

# Full check (lint + typecheck + test)
npm run check
```

## Architecture Overview

### Package Structure

The project follows a standard TypeScript/Node.js layout with clear separation of concerns:

- **`src/index.ts`**: Entry point that exports CLI runner
- **`src/cli/`**: CLI command implementations using Commander.js
  - `index.ts`: CLI setup with Commander.js, registers all commands
  - `install.ts`: Installation command that launches Ink-based TUI
  - `uninstall.ts`: Uninstall command with lock file reading
  - `init.ts`: Initialize DOR/DOD/TASK-DOD templates
  - `spec.ts`: Create numbered spec directories with TOML output
  - `statusline.ts`: Cross-platform stdio passthrough to shell scripts

- **`src/core/`**: Core business logic modules
  - `installer/`: Installation logic and file management
    - `Installer.ts`: Main installer with rollback mechanism
    - `LockManager.ts`: Lock file creation, reading, checksum management
    - `SettingsMerger.ts`: Deep merge settings with backup/restore
  - `init/`: Template initialization logic
    - `Initializer.ts`: DOR/DOD/TASK-DOD template processing
  - `spec/`: Specification directory management
    - `SpecGenerator.ts`: Auto-incrementing IDs, TOML output, template generation
  - `types/`: TypeScript type definitions
    - `config.ts`: Configuration types for all commands
    - `settings.ts`: Claude settings.json types
    - `lock.ts`: Lock file format types with v1→v2 migration support

- **`src/ui/`**: Ink-based interactive UI components (React)
  - `install/`: Installation wizard components
    - `InstallWizard.tsx`: Main wizard with state machine (6 states)
    - `ChoiceSelector.tsx`, `FinalConfirmation.tsx`, `Complete.tsx`: UI components
  - `uninstall/`: Uninstall wizard
    - `UninstallWizard.tsx`: Lock file reading and confirmation
  - `shared/`: Reusable UI components
    - `theme.ts`, `Spinner.tsx`, `ErrorDisplay.tsx`, `Banner.tsx`

- **`tests/`**: Comprehensive test suite
  - `core/`: Unit tests for all core modules
  - `ui/`: UI component tests using ink-testing-library
  - `integration/`: Integration tests for full workflows

- **`docs/`**: Project documentation structure
  - `domain/`: Business rules, workflows, and domain patterns
  - `patterns/`: Technical patterns and architectural solutions
  - `interfaces/`: API contracts and service integrations
  - `specs/`: Feature specifications (PRD, SDD, PLAN documents)

### Embedded Assets

The application includes all assets as separate files distributed with the npm package:
- `assets/claude/agents/**/*.md`: Agent definitions (11 roles, 39 activities)
- `assets/claude/commands/**/*.md`: Slash command definitions (5 commands)
- `assets/claude/output-styles/the-startup.md`: Custom output style
- `assets/the-startup/templates/*`: Template files (PRD, BRD, SDD, PLAN, DOR, DOD, TASK-DOD)
- `assets/the-startup/rules/*`: Agent delegation and cycle pattern rules

### UI Architecture (Ink/React)

The installer uses React components with Ink for terminal UI:

1. **InstallWizard**: Main component orchestrating installation flow
2. **State Machine**: Manages transitions between installation steps
   - Intro → StartupPath → ClaudePath → FileSelection → Installing → Complete
3. **React Components**: Each step is a separate React component
   - ChoiceSelector: Arrow key navigation menu for path selection
   - FinalConfirmation: Final confirmation screen with file tree display
   - Complete: Shows installation success summary
   - Banner: ASCII art banner display
   - ErrorDisplay: Handles error display with recovery options

### Installation Flow

1. User runs `the-agentic-startup install` or `npm run dev install`
2. Ink TUI launches with path selection
3. Files are selected using interactive tree (or --yes flag for all)
4. Assets are copied to:
   - `.claude/agents/` and `.claude/commands/`: Agent and command definitions
   - `.the-startup/templates/`: Template files
   - `.the-startup/rules/`: Agent delegation rules
5. Settings.json is merged with hooks (backup created)
6. Lock file (v2 format) is created with checksums
7. Rollback mechanism ensures clean state on any failure

### Lock File Format

The project supports two lock file versions with automatic migration:

**v1 (deprecated)**: Flat map of files
```json
{
  "version": "1.0.0",
  "files": {
    "agents/the-chief.md": { "size": 2048, "checksum": "sha256:..." }
  }
}
```

**v2 (current)**: Categorized arrays with metadata
```json
{
  "version": "2.0.0",
  "install_date": "2025-10-06T10:00:00Z",
  "categories": {
    "agents": [
      { "path": "agents/the-chief.md", "size": 2048, "checksum": "sha256:..." }
    ]
  }
}
```

Migration happens automatically on read via `LockManager.migrateLockFile()`.

## Key Implementation Details

### File Path Handling
- Installation paths support `~` expansion for home directory
- Project-local installation uses `.the-startup` directory
- Claude configuration expected at `~/.claude`
- All paths use `path.join()` for cross-platform compatibility

### Placeholder Replacement
Templates use placeholders that are replaced during installation:
- `{{STARTUP_PATH}}`: Installation directory path
- `{{CLAUDE_PATH}}`: Claude configuration directory

### Rollback Mechanism
The installer implements atomic operations with full rollback:
- Tracks all installed files during installation
- Creates backup of settings.json before modification
- On any failure: deletes all installed files, restores settings backup
- Ensures system is in clean state after failed installation

### Error Handling
- Custom error types with specific messages (ENOENT, EACCES, ENOSPC)
- Installation validates paths and provides clear error messages
- Settings merger handles JSON parse errors gracefully
- Lock file migration handles both v1 and v2 formats

### Testing Strategy
- Unit tests for all core modules (Installer, LockManager, SettingsMerger, etc.)
- Integration tests for full installation/uninstall workflows
- UI tests using ink-testing-library for React components
- Migration test from v1 to v2 lock file format
- Real-file integration tests using actual assets for verification

## Distribution

The project is distributed as an npm package:
- Published to npm registry
- Installed globally: `npm install -g the-agentic-startup`
- Or used via npx: `npx the-agentic-startup install`
- Assets are included in the npm package (not embedded at build time)
- Entry point: `dist/index.js` (built from TypeScript)
