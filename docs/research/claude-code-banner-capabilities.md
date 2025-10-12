# Claude Code Plugin Banner Capabilities Research

**Research Date:** 2025-10-10
**Context:** Investigating whether SessionStart hooks can display ASCII art banners and styled terminal output

## Executive Summary

### Can SessionStart hooks display ASCII art banners?
**PARTIALLY - with significant limitations**

### Can SessionStart hooks use ANSI colors/formatting?
**YES - but only in context, not directly visible to users**

---

## Key Findings

### 1. SessionStart Hook Output Behavior

**Primary Finding:** SessionStart hook stdout is **added to Claude's context**, NOT displayed directly to users.

From official documentation:
> "For `SessionStart` hooks, stdout is added as context for Claude"

This means:
- ✅ The hook can output ANSI-formatted text
- ✅ Multi-line output including ASCII art is supported
- ❌ Output is NOT shown in the UI/terminal to users
- ❌ Users will NOT see banners directly

### 2. ANSI Color Code Support

**Evidence:** Production hooks successfully use ANSI escape codes.

Example from `disler/claude-code-hooks-mastery` (status_line.py):
```python
# ANSI color codes used in production
CYAN = '\033[36m'
BLUE = '\033[34m'
GREEN = '\033[32m'
GRAY = '\033[90m'
RED = '\033[31m'

# Output example
print(f"{CYAN}[{model_name}]{RESET} | 📁 {BLUE}{dir_name}{RESET} | 🌿 {GREEN}{branch}{RESET}")
```

**Confirmed:**
- ✅ ANSI escape codes are NOT stripped by Claude Code
- ✅ Color formatting works in hook scripts
- ✅ Unicode emojis are supported

### 3. Output Display Mechanisms

There are **THREE** ways hooks can communicate with users:

#### Option 1: stdout (Context Only - SessionStart)
```bash
#!/bin/bash
# This goes into Claude's context, NOT shown to user
echo -e "\033[1;34m╔═══════════════════════════════════════╗\033[0m"
echo -e "\033[1;34m║     THE AGENTIC STARTUP               ║\033[0m"
echo -e "\033[1;34m╚═══════════════════════════════════════╝\033[0m"
```

**Result:** Claude sees the formatted banner in context, user does NOT

#### Option 2: systemMessage (JSON Field)
```bash
#!/bin/bash
cat << 'EOF'
{
  "systemMessage": "⚠️ The Agentic Startup plugin loaded successfully",
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": "Project initialized with The Agentic Startup patterns"
  }
}
EOF
```

**Result:**
- User sees: "⚠️ The Agentic Startup plugin loaded successfully" (plain text warning)
- Claude receives: "Project initialized with The Agentic Startup patterns" (context)

**Limitations:**
- systemMessage is described as "warning message" - likely plain text
- No documented support for ANSI codes or multi-line ASCII art
- Format/styling unknown - documentation does not specify

#### Option 3: Transcript Mode (CTRL-R)
- SessionStart hooks do NOT show stdout in transcript mode
- Only certain hooks (not SessionStart) display in transcript

### 4. Alternative Hook Types

**Other hooks that show output to users:**

| Hook Type | Output Visibility | Use Case |
|-----------|------------------|----------|
| PreToolUse | ❌ Cannot display (Issue #4084) | Before tool execution |
| PostToolUse | ✅ Shows in transcript (CTRL-R) | After tool execution |
| Stop | ✅ Shows in transcript (CTRL-R) | When Claude finishes |
| UserPromptSubmit | ❌ Cannot display (Issue #4084) | Before prompt processing |

**Known Issue:** GitHub Issue #4084 reports that some hooks cannot display output to users despite documentation suggesting they should.

---

## Evidence-Based Conclusions

### What WORKS:
1. ✅ **ANSI codes in hook scripts** - Proven in production (status_line.py)
2. ✅ **Multi-line output** - No restrictions documented
3. ✅ **Unicode/emoji** - Used in production hooks
4. ✅ **Structured JSON output** - systemMessage field for user warnings

### What DOES NOT WORK:
1. ❌ **Direct banner display to users from SessionStart** - stdout goes to context only
2. ❌ **ASCII art visible in UI** - SessionStart output not shown to users
3. ❌ **Formatted terminal output in SessionStart** - context injection only

### What is UNKNOWN:
1. ❓ **systemMessage formatting support** - Documentation does not specify:
   - Can it contain ANSI codes? (Unknown)
   - Can it be multi-line? (Unknown)
   - How is it displayed in UI? (Unknown - described as "warning")
   - Examples? (None found)

---

## Practical Recommendations

### For Banner Display in SessionStart Hook:

**Option A: Context-Only Banner (Works but invisible to users)**
```bash
#!/bin/bash
# Banner will be in Claude's context, helpful for Claude's awareness
cat << 'EOF'
╔═══════════════════════════════════════╗
║     THE AGENTIC STARTUP               ║
║     Enterprise AI Development         ║
╚═══════════════════════════════════════╝

Project patterns loaded:
- Agent delegation rules
- Specification workflows
- Development templates
EOF
```

**Use Case:** Claude becomes aware of plugin initialization, can mention it in responses

**Option B: systemMessage for User Notification (Plain text)**
```bash
#!/bin/bash
cat << 'EOF'
{
  "systemMessage": "✓ The Agentic Startup plugin loaded successfully",
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": "╔═══════════════════════════════════════╗\n║     THE AGENTIC STARTUP               ║\n║     Enterprise AI Development         ║\n╚═══════════════════════════════════════╝\n\nProject patterns active:\n- Agent delegation framework\n- Specification-driven development\n- Enterprise templates"
  }
}
EOF
```

**Use Case:**
- User sees simple confirmation message
- Claude receives formatted banner in context
- Best of both worlds (but banner not directly visible)

**Option C: Welcome Message in First Response**
```bash
#!/bin/bash
# Context tells Claude to show banner in first response
cat << 'EOF'
{
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": "IMPORTANT: When starting this session, display the following banner to the user in your first response:

╔═══════════════════════════════════════╗
║     THE AGENTIC STARTUP               ║
║     Enterprise AI Development         ║
╚═══════════════════════════════════════╝

Then inform them that The Agentic Startup plugin is active with agent delegation, specification workflows, and enterprise templates available."
  }
}
EOF
```

**Use Case:**
- Claude displays banner in first message to user
- User sees formatted banner (via Claude, not hook directly)
- Requires Claude cooperation

---

## Testing Recommendations

To definitively answer unknown questions:

### Test 1: systemMessage with Multi-line
```bash
#!/bin/bash
cat << 'EOF'
{
  "systemMessage": "Line 1\nLine 2\nLine 3"
}
EOF
```

**Question:** Does systemMessage support `\n` newlines?

### Test 2: systemMessage with ANSI Codes
```bash
#!/bin/bash
cat << 'EOF'
{
  "systemMessage": "\033[1;34mBlue Bold Text\033[0m"
}
EOF
```

**Question:** Does systemMessage render ANSI codes or show them raw?

### Test 3: systemMessage with ASCII Art
```bash
#!/bin/bash
cat << 'EOF'
{
  "systemMessage": "╔═══╗\n║ A ║\n╚═══╝"
}
EOF
```

**Question:** Can systemMessage display box-drawing characters?

---

## References

### Official Documentation
- [Claude Code Hooks Reference](https://docs.claude.com/en/docs/claude-code/hooks)
- Hook output: "For SessionStart hooks, stdout is added as context for Claude"
- systemMessage: "Optional warning message shown to the user"

### Production Examples
- [disler/claude-code-hooks-mastery](https://github.com/disler/claude-code-hooks-mastery)
  - Uses ANSI color codes successfully
  - status_line.py with `\033[36m` cyan, `\033[34m` blue, etc.
- [EvanL1/claude-code-hooks](https://github.com/EvanL1/claude-code-hooks)
  - terminal-ui.sh for terminal enhancement

### Known Issues
- [Issue #4084](https://github.com/anthropics/claude-code/issues/4084) - Hook Output Visibility Blocked
- [Issue #4318](https://github.com/anthropics/claude-code/issues/4318) - SessionStart/SessionEnd Feature Request

---

## Final Answer

### Can SessionStart hooks display ASCII art banners?
**NO - Not directly to users.**

SessionStart stdout goes to Claude's context only. Users will NOT see the banner unless:
1. Claude mentions it in a response (Option C above)
2. You use systemMessage (unknown if ASCII art works - needs testing)

### Can SessionStart hooks use ANSI colors/formatting?
**YES - But only in context, not in user-visible output.**

ANSI codes work in hook scripts and are NOT stripped. However:
- SessionStart output → Claude's context (invisible to users)
- systemMessage field → User warning (format support unknown)

### Recommended Approach
**Use Option C:** Instruct Claude via context to display the banner in its first response. This is the most reliable way to show a banner to users while maintaining proper formatting.

```bash
#!/bin/bash
cat << 'EOF'
{
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": "Display this banner in your first response:\n\n╔═══════════════════════════════════════╗\n║     THE AGENTIC STARTUP               ║\n║     Enterprise AI Development         ║\n╚═══════════════════════════════════════╝\n\nInform the user that The Agentic Startup plugin is active."
  }
}
EOF
```

This delegates banner display to Claude, ensuring users see it while working within documented SessionStart behavior.
