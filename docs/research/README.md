# Research Documentation

This directory contains research findings and proof-of-concept implementations for The Agentic Startup plugin.

## Contents

### 1. Banner Capabilities Research
**File:** `claude-code-banner-capabilities.md`

Comprehensive research on Claude Code plugin capabilities for displaying banners and styled terminal output in SessionStart hooks.

**Key Findings:**
- ✅ ANSI codes supported in hook scripts
- ❌ SessionStart stdout NOT shown to users (goes to context only)
- ✅ Can instruct Claude to display banner in first response
- ❓ systemMessage formatting support unknown (needs testing)

**Evidence Sources:**
- Official Claude Code documentation
- Production hook examples (disler/claude-code-hooks-mastery)
- GitHub issues and community examples

### 2. SessionStart Hook Implementations

#### Basic Implementation
**File:** `sessionstart-hook-implementation.sh`

Simple SessionStart hook that instructs Claude to display The Agentic Startup banner in its first response.

**Features:**
- ASCII art banner
- Plugin capability summary
- Available agents and commands list
- Simple systemMessage notification

**Usage:**
```bash
./sessionstart-hook-implementation.sh
```

#### Advanced Implementation
**File:** `sessionstart-hook-advanced.sh`

Context-aware SessionStart hook that detects project installation state and provides appropriate messaging.

**Features:**
- Installation state detection (checks for `.the-startup/.lock`)
- Different messages for installed vs. available states
- Comprehensive capability listing
- Conditional banner display

**Usage:**
```bash
./sessionstart-hook-advanced.sh
```

## Implementation Approach

### Why Claude Displays the Banner

SessionStart hooks cannot directly show output to users. Instead, we use a clever workaround:

1. **Hook outputs JSON** with `hookSpecificOutput.additionalContext`
2. **Context includes instruction** for Claude to display banner
3. **Claude reads context** when session starts
4. **Claude displays banner** in first response to user

This approach:
- ✅ Works within documented SessionStart behavior
- ✅ Ensures users see formatted banner
- ✅ Provides full ANSI/ASCII art support (via Claude's rendering)
- ✅ No undocumented features or hacks

### Integration with The Agentic Startup

To use these hooks in the plugin:

1. **Include hook script in assets:**
   ```
   assets/claude/hooks/SessionStart.sh
   ```

2. **Install hook during plugin installation:**
   - Copy to `~/.claude/hooks/SessionStart.sh`
   - Or to `.claude/hooks/SessionStart.sh` for project-local

3. **Configure in settings.json:**
   ```json
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

## Testing Strategy

### Manual Testing

1. **Test basic hook:**
   ```bash
   cd /Users/irudi/Code/personal/the-startup/docs/research
   ./sessionstart-hook-implementation.sh
   ```

2. **Verify JSON output:**
   - Should output valid JSON
   - Check `systemMessage` field
   - Check `hookSpecificOutput.additionalContext` contains banner instruction

3. **Test in Claude Code:**
   - Configure hook in `.claude/settings.json`
   - Start new session
   - Verify systemMessage appears
   - Verify Claude displays banner in first response

### Automated Testing

Consider adding tests for:
- JSON validity
- Banner format preservation
- Context message generation
- Installation state detection

## Open Questions

### systemMessage Formatting

**Unknown:** What formatting does systemMessage support?

**Tests to run:**
1. Multi-line with `\n`: `"Line 1\nLine 2"`
2. ANSI codes: `"\033[1;34mBlue\033[0m"`
3. ASCII art: `"╔═══╗\n║ A ║"`

**How to test:**
```bash
# Create test hook
cat > test-systemmsg.sh << 'EOF'
#!/bin/bash
echo '{"systemMessage": "Line 1\nLine 2\nLine 3"}'
EOF

# Configure in settings.json and observe output
```

### Alternative Display Methods

**Explore:**
- Can PostToolUse hooks display banners after first tool execution?
- Can Stop hooks show session summary with banner?
- Does verbose mode (CTRL-R) show SessionStart output?

## Next Steps

1. **Decide on implementation:**
   - Basic (simple, reliable)
   - Advanced (context-aware, feature-rich)

2. **Integrate into installer:**
   - Add hook script to assets
   - Copy during installation
   - Configure in settings.json merge

3. **Test with real Claude Code session:**
   - Verify banner display
   - Confirm context injection
   - Validate user experience

4. **Document for users:**
   - Explain banner behavior
   - Provide customization options
   - Show how to disable if desired

## References

- [Claude Code Hooks Documentation](https://docs.claude.com/en/docs/claude-code/hooks)
- [disler/claude-code-hooks-mastery](https://github.com/disler/claude-code-hooks-mastery) - Production examples
- [Issue #4084](https://github.com/anthropics/claude-code/issues/4084) - Hook output visibility
- [Issue #4318](https://github.com/anthropics/claude-code/issues/4318) - SessionStart feature request

---

**Research Date:** 2025-10-10
**Status:** Complete - Ready for implementation
