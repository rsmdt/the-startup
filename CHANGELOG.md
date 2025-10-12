# Changelog

All notable changes to The Agentic Startup will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2025-10-12

### Changed
- **BREAKING:** Complete migration from npm CLI package to Claude Code plugin architecture
- Installation now uses `/plugin install` instead of `npx the-agentic-startup install`
- Removed Ink-based TUI installer (no longer needed with plugin system)
- Simplified installation process - one command installs everything

### Added
- **Hooks System**: SessionStart and UserPromptSubmit hooks
  - Welcome banner on first plugin session
  - Git branch statusline integration
- **Plugin Manifest**: `.claude-plugin/plugin.json` for plugin discovery
- **Scripts Directory**: `scripts/spec.sh` for spec generation
- **Spec Command**: `/s:spec` for creating numbered specification directories
- Auto-incrementing spec IDs (001, 002, 003...)
- TOML output format for spec metadata reading
- Template generation support via `--add` flag

### Improved
- **File References**: Commands now use @ notation (`@rules/agent-delegation.md`) instead of placeholders
- **Component Discovery**: All components (agents, commands, hooks) auto-discovered by Claude Code
- **Directory Structure**: Flattened structure with all components at repository root
- **Documentation**: Updated README for plugin installation and usage
- **Agent Access**: All 50 agents immediately available after installation

### Removed
- npm package installation workflow
- Interactive TUI installer (Ink components)
- Lock file management system
- Settings.json merger and backup/restore
- CLI-specific source code (`src/cli/`, `src/ui/`, `src/core/installer/`)
- Build-time placeholder replacement

### Technical
- Plugin structure follows official Claude Code specifications
- Hooks use `${CLAUDE_PLUGIN_ROOT}` for script paths
- Commands use @ notation for runtime file references
- No build step required - files used as committed to Git
- Cross-platform statusline support (bash/PowerShell)

### Migration Guide

**From 1.x (npm) to 2.x (plugin):**

1. Uninstall npm package:
   ```bash
   npx the-agentic-startup uninstall
   npm uninstall -g the-agentic-startup
   ```

2. Install plugin:
   ```bash
   /plugin install irudiperera/the-startup
   ```

3. Output style (manual installation):
   - Copy `assets/claude/output-styles/the-startup.md` to `~/.claude/output-styles/`
   - Activate: `/settings add "outputStyle": "the-startup"`

**What stays the same:**
- All 50 agents work identically
- All slash commands work identically
- Specification workflow unchanged
- Documentation structure unchanged
- Agent delegation rules unchanged

**What's better:**
- Simpler installation (one command)
- Automatic updates via plugin system
- Welcome banner on first use
- Git statusline integration

## [1.0.0] - 2025-09-13

### Added
- Initial release as npm CLI package
- Interactive installation via Ink-based TUI
- 50 specialized agents across 9 professional roles
- 5 slash commands: `/s:specify`, `/s:analyze`, `/s:implement`, `/s:refactor`, `/s:init`
- The Startup output style
- Statusline integration (manual configuration)
- Agent delegation rules and cycle patterns
- Template system for PRD, SDD, PLAN, DOR, DOD, TASK-DOD
- Lock file system for tracking installed components
- Settings.json deep merge with backup/restore
- Rollback mechanism for failed installations
- Component selection during installation
- Cross-platform support (macOS, Linux, Windows)

### Technical
- Built with TypeScript
- CLI using Commander.js
- TUI using Ink (React for CLI)
- Published to npm registry
- Installable via `npx` or `npm install -g`

---

[2.0.0]: https://github.com/irudiperera/the-startup/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/irudiperera/the-startup/releases/tag/v1.0.0
