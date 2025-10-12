# Quick Reference: Claude Code Banner Display

## TL;DR

**Q: Can SessionStart hooks display ASCII art banners?**
**A: NO - not directly. But YES - via Claude displaying it for you.**

## The Solution

SessionStart hooks **cannot** show output directly to users. But they **can** instruct Claude to display a banner in its first response.

## Working Implementation

```bash
#!/bin/bash
# SessionStart hook script

BANNER="╔═══════════════════════════════════════╗
║     THE AGENTIC STARTUP               ║
║     Enterprise AI Development         ║
╚═══════════════════════════════════════╝"

cat << EOF
{
  "systemMessage": "✓ The Agentic Startup plugin loaded",
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": "Display this banner in your first response:\n\n${BANNER}\n\nThen inform the user that The Agentic Startup plugin is active."
  }
}
EOF
```

## How It Works

1. **Hook outputs JSON** → Goes to Claude (not user)
2. **additionalContext contains instruction** → Claude reads it
3. **Claude displays banner** → User sees it in first response

## What Works vs. What Doesn't

### ✅ WORKS
- ANSI color codes in hooks (proven in production)
- Multi-line output and ASCII art
- Unicode characters and emojis
- Instructing Claude to display banner
- systemMessage for simple notifications

### ❌ DOESN'T WORK
- Direct banner display to users from SessionStart
- SessionStart stdout visible in UI (goes to context only)
- Formatted terminal output from SessionStart (context only)

### ❓ UNKNOWN (Needs Testing)
- Does systemMessage support ANSI codes?
- Does systemMessage support multi-line text?
- Can systemMessage display ASCII art?

## Evidence

**Production Example:** [disler/claude-code-hooks-mastery](https://github.com/disler/claude-code-hooks-mastery)
```python
# Uses ANSI codes successfully
CYAN = '\033[36m'
BLUE = '\033[34m'
print(f"{CYAN}[{model}]{RESET} | {BLUE}{dir}{RESET}")
```

**Official Docs:** [Claude Code Hooks](https://docs.claude.com/en/docs/claude-code/hooks)
> "For SessionStart hooks, stdout is added as context for Claude"

## Files in This Research

1. **claude-code-banner-capabilities.md** - Full research findings
2. **sessionstart-hook-implementation.sh** - Basic working implementation
3. **sessionstart-hook-advanced.sh** - Advanced context-aware version
4. **README.md** - Complete documentation
5. **QUICK-REFERENCE.md** - This file

## Recommendation

**Use the "Claude displays it" approach:**
- Most reliable
- Works within documented behavior
- Supports full formatting
- Best user experience

**Implementation:**
```json
// .claude/settings.json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "/path/to/SessionStart.sh"
          }
        ]
      }
    ]
  }
}
```

## Alternative: systemMessage Testing

Want to test if systemMessage supports formatting? Try:

```bash
# Test 1: Multi-line
echo '{"systemMessage": "Line 1\nLine 2\nLine 3"}'

# Test 2: ANSI colors
echo '{"systemMessage": "\033[1;34mBlue Bold\033[0m"}'

# Test 3: ASCII art
echo '{"systemMessage": "╔═══╗\n║ A ║\n╚═══╝"}'
```

Configure hook in settings.json and observe the output.

---

**Bottom Line:** Don't fight the system. Let Claude display your banner. It works perfectly.
