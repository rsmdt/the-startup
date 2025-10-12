# Claude Code Plugin Post-Installation Messaging Research

**Research Date:** 2025-10-10
**Status:** COMPLETE
**Confidence:** HIGH (based on official documentation and examples)

## Executive Summary

Claude Code plugins **DO NOT have native post-installation messaging** mechanisms like completion screens or automatic README display. However, plugins **CAN use SessionStart hooks** to display welcome messages when Claude Code starts a new session after plugin installation.

## Key Findings

### 1. No Native Post-Installation Hooks

**Evidence:**
- Official plugin.json schema has no fields for installation messages, welcome text, or post-install notifications
- Official example plugins (security-guidance, feature-dev) contain no installation message mechanisms
- Plugin documentation does not mention any post-installation display system
- No README.md automatic display after `/plugin install`

**Source:**
- https://docs.claude.com/en/docs/claude-code/plugins-reference
- https://github.com/anthropics/claude-code/tree/main/plugins/security-guidance
- https://github.com/anthropics/claude-code/tree/main/plugins/feature-dev

### 2. SessionStart Hook Available (RECOMMENDED APPROACH)

**What it is:**
- Fires when Claude Code starts a new session or resumes existing session
- Runs after plugin installation on first session start
- Can execute commands that output to stdout
- Stdout content is automatically added to context (visible to user)

**Sources:**
- "startup": Fresh launch of Claude Code
- "resume": Resuming existing session
- "clear": After clearing session
- "compact": After context compaction

**Configuration:**
```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "echo 'Welcome to The Agentic Startup!'"
          }
        ]
      }
    ]
  }
}
```

**Evidence:**
- Hooks reference: https://docs.claude.com/en/docs/claude-code/hooks
- SessionStart hook guide: https://docs.claude.com/en/docs/claude-code/hooks-guide
- Example implementation: https://github.com/listfold/claude-git (uses SessionStart for initialization)
- Blog post: https://blog.gitbutler.com/automate-your-ai-workflows-with-claude-code-hooks

**Key characteristics:**
- SessionStart hooks do not use matchers (unlike PreToolUse/PostToolUse)
- Multiple commands can be chained
- Stdout is automatically added to context and shown to user
- Perfect for displaying setup instructions, available commands, and first-run guidance

### 3. Plugin Installation UX

**What users see:**
1. Run `/plugin install plugin-name@marketplace-name`
2. Plugin installs immediately (no restart required as of latest version)
3. Plugin components (commands, agents, hooks) become active
4. No automatic completion message or README display
5. Users can verify with `/help` to see new commands
6. **SessionStart hook fires on next session start** (after install completes)

**Evidence:**
- Plugin documentation: https://docs.claude.com/en/docs/claude-code/plugins
- No restart required: https://www.anthropic.com/news/claude-code-plugins

### 4. plugin.json Schema

**Available fields:**
- `name`, `version`, `description` (shown in plugin browser)
- `author`, `homepage`, `repository`, `license`, `keywords`
- `commands`, `agents`, `hooks`, `mcpServers` (paths to assets)

**NOT available:**
- No `welcomeMessage` field
- No `postInstall` field
- No `installationMessage` field
- No `notification` field

**Evidence:**
https://docs.claude.com/en/docs/claude-code/plugins-reference

### 5. Alternative Approaches Considered

#### Notification Hook
- **Purpose:** Triggers when Claude sends notifications or waits for input
- **Verdict:** NOT suitable - only fires during permission requests, not installation
- **Source:** https://docs.claude.com/en/docs/claude-code/hooks

#### README.md Display
- **Verdict:** NOT supported - no automatic display mechanism
- **Reality:** READMEs exist for documentation but are not shown to users after install

#### Command Description
- **What it is:** Description field in slash command definitions
- **Verdict:** Partial solution - shown in `/help` but not proactive
- **Limitation:** Users must run `/help` to discover commands

## Recommended Implementation

### Option 1: SessionStart Hook with Welcome Message (RECOMMENDED)

Create a SessionStart hook that displays a formatted welcome message:

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "cat ~/.the-startup/welcome.txt",
            "description": "Show welcome message after installation"
          }
        ]
      }
    ]
  }
}
```

**Advantages:**
- Official, documented mechanism
- Runs automatically on first session after install
- Stdout added to context (visible to user)
- Can display formatted text, ASCII art, instructions
- Works across terminal and VS Code

**Disadvantages:**
- Fires on EVERY session start (not just first run)
- Need to manage when to show (could check for flag file)
- Less control over timing than CLI completion screen

### Option 2: First-Run Detection Pattern

Combine SessionStart with state tracking:

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "/path/to/first-run-check.sh",
            "description": "Check if first run and show welcome"
          }
        ]
      }
    ]
  }
}
```

Script checks for `.the-startup/.initialized` flag:
- If missing: Show welcome message, create flag
- If present: Skip welcome message

**Advantages:**
- Shows welcome only once
- Clean UX (no repeated messages)
- Simple state management

**Disadvantages:**
- Requires script execution
- More complex than simple echo
- Need to handle flag file lifecycle

### Option 3: Welcome Slash Command

Create a `/startup:welcome` command users can run:

```markdown
# .claude/commands/startup-welcome.md

Display welcome message and getting started guide for The Agentic Startup.

Show the following information:
1. Installation complete
2. Available commands: /prd:create, /prd:execute, /s:specify, /s:implement, etc.
3. How to activate output style: /settings add "outputStyle": "the-startup"
4. First steps: Run /s:init to initialize templates
5. Documentation: Point to .the-startup/templates/ directory

Format with clear headings and bullet points.
```

**Advantages:**
- User-controlled (no automatic trigger)
- Available anytime users need help
- Simple implementation (just a command file)

**Disadvantages:**
- NOT proactive (users must know to run it)
- Defeats purpose of post-install guidance
- Less discoverable than automatic message

## Comparison to CLI Installation

### Current CLI UX (the-agentic-startup install)
```
✓ Installation complete!

Next steps:
1. Activate the output style in Claude settings
2. Run /s:init to initialize templates
3. Use /s:specify to create your first spec
4. Check .the-startup/templates/ for available templates
```

### Plugin UX (with SessionStart hook)
```bash
# After /plugin install the-agentic-startup@marketplace
# On next session start:

Welcome to The Agentic Startup!

Installation complete. Next steps:
1. Activate output style: /settings add "outputStyle": "the-startup"
2. Initialize templates: /s:init
3. Create your first spec: /s:specify "feature description"
4. Available commands: /prd:create, /prd:execute, /s:specify, /s:implement

For more information, run /startup:welcome anytime.
```

### Key Differences
- CLI: Shows immediately after installation completes
- Plugin: Shows on next session start (slight delay)
- CLI: One-time message (script ends)
- Plugin: Needs first-run detection to avoid repetition

## Recommendations

### For The Agentic Startup Plugin

**Implement both approaches:**

1. **SessionStart hook with first-run detection**
   - Show welcome message once after installation
   - Use `.the-startup/.plugin-initialized` flag
   - Display essential setup steps and available commands

2. **Permanent /startup:welcome command**
   - Available anytime users need guidance
   - Listed in `/help` for discoverability
   - Detailed reference version of welcome message

3. **Output style reminder**
   - Include clear instructions in welcome message
   - This is the biggest gap vs CLI (no auto-merge of settings.json)

### Welcome Message Content

```
========================================
  THE AGENTIC STARTUP - READY TO USE
========================================

Installation complete! Your plugin is active.

IMPORTANT: Activate the output style
  Run: /settings add "outputStyle": "the-startup"

Available Commands:
  /s:specify      - Create feature specifications
  /s:implement    - Execute specification plans
  /s:init         - Initialize DOR/DOD templates
  /s:analyze      - Analyze business rules
  /s:refactor     - Refactor code safely
  /prd:create     - Create Product Requirements
  /prd:execute    - Execute PRD implementation

First Steps:
  1. Initialize templates: /s:init
  2. Create your first spec: /s:specify "your feature idea"
  3. Review templates in .the-startup/templates/

Need help? Run /startup:welcome anytime.
========================================
```

## Implementation Details

### File Structure
```
.claude-plugin/
  plugin.json         # Plugin metadata
hooks/
  hooks.json          # SessionStart hook configuration
  welcome.sh          # First-run check and message display
scripts/
  welcome-message.txt # Formatted welcome text
commands/
  startup-welcome.md  # Manual welcome command
```

### hooks/hooks.json
```json
{
  "SessionStart": [
    {
      "hooks": [
        {
          "type": "command",
          "command": "./hooks/welcome.sh",
          "description": "Show welcome message on first run"
        }
      ]
    }
  ]
}
```

### hooks/welcome.sh
```bash
#!/bin/bash
FLAG_FILE="$HOME/.the-startup/.plugin-initialized"

if [ ! -f "$FLAG_FILE" ]; then
  cat "$(dirname "$0")/../scripts/welcome-message.txt"
  mkdir -p "$(dirname "$FLAG_FILE")"
  touch "$FLAG_FILE"
fi
```

## Limitations and Constraints

### What Plugins CANNOT Do
1. Show messages during `/plugin install` execution
2. Display README automatically after installation
3. Modify user settings.json automatically (unlike CLI)
4. Show one-time completion screens (without state tracking)
5. Trigger notifications to user's system

### What Plugins CAN Do
1. Run commands on session start via SessionStart hook
2. Display text via stdout (added to context)
3. Execute scripts for state management
4. Provide slash commands for on-demand information
5. Include hooks in plugin configuration

## Conclusion

**Answer:** Claude Code plugins support post-installation messaging through **SessionStart hooks**, not through native installation completion screens.

**Implementation:** Use SessionStart hook with first-run detection to display welcome message once, plus a permanent `/startup:welcome` command for on-demand access.

**Key Difference from CLI:** Plugins cannot show completion messages immediately after `/plugin install` - there's always a slight delay until next session start. This is an acceptable trade-off for the benefits of the plugin system.

**Evidence Level:** HIGH - based on official documentation, example plugins, and community implementations.

## Sources

### Official Documentation
- Plugin Reference: https://docs.claude.com/en/docs/claude-code/plugins-reference
- Hooks Reference: https://docs.claude.com/en/docs/claude-code/hooks
- Hooks Guide: https://docs.claude.com/en/docs/claude-code/hooks-guide
- Plugin Announcement: https://www.anthropic.com/news/claude-code-plugins

### Example Implementations
- security-guidance plugin: https://github.com/anthropics/claude-code/tree/main/plugins/security-guidance
- feature-dev plugin: https://github.com/anthropics/claude-code/tree/main/plugins/feature-dev
- claude-git (SessionStart example): https://github.com/listfold/claude-git

### Community Resources
- GitButler blog: https://blog.gitbutler.com/automate-your-ai-workflows-with-claude-code-hooks
- SessionStart feature request: https://github.com/anthropics/claude-code/issues/4318 (implemented)

## Related Research
- [Plugin Migration Strategy](./plugin-migration-strategy.md) - Overall plugin conversion approach
- [Settings Merge Alternatives](./settings-merge-alternatives.md) - How to handle outputStyle activation
