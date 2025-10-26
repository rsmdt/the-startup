# Migration Guide: v1.x to v2.0

This guide helps you migrate from The Agentic Startup v1.x (bash script installer) to v2.0 (Claude Code marketplace plugin).

## Overview

**The Agentic Startup v1.x** was distributed via bash script installer and manually placed files in your Claude Code configuration directory.

**Version 2.0** is distributed through the Claude Code marketplace plugin system with automated installation and management.

## Migration Steps

### Step 1: Uninstall v1.x

The old version installed files to your Claude Code configuration directory. Remove them before installing v2.0.

#### Remove Installed Files

```bash
# Remove agents (if installed globally)
rm -rf ~/.claude/agents/the-*.md

# Remove commands (if installed globally)
rm -rf ~/.claude/commands/s-*.md

# Remove local installation (if used --local flag)
rm -rf ./.the-startup

# Remove output style
rm -rf ~/.claude/output-styles/the-startup.md

# Remove hooks (if any were installed)
rm -rf ~/.claude/hooks/the-startup-*
```

#### Clean Up settings.json

The old installer may have modified `~/.claude/settings.json`. Check for and remove The Agentic Startup references:

**Before:**
```json
{
  "hooks": {
    "sessionStart": "~/.claude/hooks/the-startup-session.sh",
    "statuslineComplete": "~/.claude/hooks/the-startup-statusline.sh"
  },
  "outputStyle": "The Startup"
}
```

**After:**
```json
{
  "hooks": {},
  "outputStyle": "Default"
}
```

**Note:** Remove only The Agentic Startup-related hooks and settings. Keep any other customizations you've made.

#### Verify Complete Removal

```bash
# Check for any remaining files
find ~/.claude -name "*startup*" -o -name "*the-startup*"
```

If this returns no results, cleanup is complete.

### Step 2: Install v2.0

After uninstalling v1.x, install the new marketplace version:

```bash
# Add The Agentic Startup marketplace
/plugin marketplace add rsmdt/the-startup

# Install the Start plugin (required - core workflows)
/plugin install start@the-startup

# (Optional) Install the Team plugin for specialized agents
/plugin install team@the-startup
```

### Step 3: Initialize Configuration

Run the initialization command to set up output style and statusline:

```bash
/start:init
```

This will:
- Configure "The Startup" output style
- Set up git-aware statusline
- Add necessary hooks automatically
- Prompt for your preferences

## Key Differences

| Feature | v1.x (Deprecated) | v2.0 (Current) |
|---------|-------------------|----------------|
| **Distribution** | Bash script installer | Claude Code marketplace |
| **Installation** | Manual file placement | Automatic plugin installation |
| **Scope** | Global or local | Plugin-based (isolated) |
| **Agents** | Static files in `.claude/agents/` | Dynamic skills system |
| **Commands** | `/s:*` prefix | `/start:*` prefix |
| **Configuration** | Manual `settings.json` editing | `/start:init` wizard |
| **Updates** | Re-run install script | `/plugin update` |
| **Uninstall** | Manual file deletion | `/plugin uninstall start@the-startup` |

## Breaking Changes

### Command Prefix Change

**v1.x:**
```bash
/s:specify Build a feature
/s:implement 001
/s:analyze patterns
/s:refactor code
```

**v2.0:**
```bash
/start:specify Build a feature
/start:implement 001
/start:analyze patterns
/start:refactor code
```

Update any documentation or scripts that reference the old command names.

### Agent Files Removed

v1.x used static agent definition files:
- `~/.claude/agents/the-chief.md`
- `~/.claude/agents/the-architect.md`
- etc.

v2.0 uses a **skills system** instead:
- `documentation` skill - Auto-documents patterns/interfaces
- `agent-delegation` skill - Coordinates specialist agents

**Migration:** No action needed. The new skills system provides enhanced functionality automatically.

### Hooks Management

**v1.x:** Manual hook script files
**v2.0:** Hooks managed by plugin system

Run `/start:init` to configure hooks automatically. The plugin handles hook lifecycle.

### Settings Configuration

**v1.x:** Manual JSON editing
**v2.0:** Interactive wizard via `/start:init`

The new approach prevents configuration errors and provides validation.

## Verification

After migration, verify the installation:

```bash
# Check installed plugins
/plugin list

# Verify commands are available
/start:init

# Confirm output style is configured
# Check bottom of Claude Code for statusline
```

You should see:
- `start@the-startup` in plugin list
- `/start:*` commands in command palette
- Git branch in statusline (after running `/start:init`)

## Troubleshooting

### Commands Not Available

**Issue:** `/start:*` commands don't appear

**Solution:**
1. Verify plugin installation: `/plugin list`
2. Reinstall if needed: `/plugin install start@the-startup`
3. Restart Claude Code

### Statusline Not Showing

**Issue:** Git branch not appearing in statusline

**Solution:**
1. Run `/start:init`
2. Confirm statusline setup when prompted
3. Check `~/.claude/settings.json` for hooks configuration

### Old Commands Still Present

**Issue:** Both `/s:*` and `/start:*` commands appear

**Solution:**
1. Complete Step 1 (Uninstall v1.x) thoroughly
2. Verify cleanup: `find ~/.claude -name "s-*.md"`
3. Remove any remaining `/s:*` command files

## Getting Help

If you encounter issues during migration:

1. **Check Documentation:** [README.md](README.md)
2. **Report Issues:** [GitHub Issues](https://github.com/rsmdt/the-startup/issues)
3. **Verify Installation:** Ensure Claude Code v2.0+ is installed

## Additional Resources

- [Claude Code Documentation](https://docs.claude.com/claude-code)
- [Plugin System Guide](https://docs.claude.com/claude-code/plugins)
- [The Agentic Startup GitHub](https://github.com/rsmdt/the-startup)
