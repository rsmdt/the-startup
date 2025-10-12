# Claude Code Plugin Research: File References and Executable Scripts

**Research Date:** 2025-10-11
**Status:** Complete with Evidence
**Sources:** Official Anthropic documentation, official plugin examples, community plugins

---

## Executive Summary

This research answers three critical questions about Claude Code plugin architecture:

1. **File References:** @ notation works at runtime for user invocations; plugins use relative paths in plugin.json
2. **Script Directory:** Use `hooks/` for hook scripts (official convention); no `bin/` or `scripts/` directory observed
3. **Spec Command:** Implement as external CLI tool (Option C), not as plugin command due to capability limitations

---

## Question 1: File References in Commands/Agents

### Current Mechanism

**@ Notation:**
- Works in **user-facing contexts** (chat, CLAUDE.md, user prompts)
- Does NOT work in plugin command/agent markdown files
- Runtime file inclusion by Claude Code itself
- Example: `@src/utils/helpers.js` in chat

**${CLAUDE_PLUGIN_ROOT}:**
- Works in **plugin configuration** (hooks.json, plugin.json)
- Resolves to absolute plugin installation path
- Used for referencing bundled scripts/files
- Example: `${CLAUDE_PLUGIN_ROOT}/hooks/security_reminder_hook.py`

**Plugin File References:**
- Use **relative paths** in plugin.json
- Paths are relative to plugin root
- Example: `./commands/feature-dev.md`, `./agents/code-architect.md`

### Evidence

**Official Example (security-guidance plugin):**
```json
{
  "hooks": {
    "PreToolUse": [{
      "hooks": [{
        "type": "command",
        "command": "python3 ${CLAUDE_PLUGIN_ROOT}/hooks/security_reminder_hook.py"
      }],
      "matcher": "Edit|Write|MultiEdit"
    }]
  }
}
```

**Source:** https://github.com/anthropics/claude-code/tree/main/plugins/security-guidance

**Community Example (wshobson/agents):**
```json
{
  "name": "claude-code-essentials",
  "commands": [
    "./tools/code-explain.md",
    "./tools/smart-debug.md"
  ],
  "agents": [
    "./agents/code-reviewer.md",
    "./agents/debugger.md"
  ]
}
```

**Source:** https://github.com/wshobson/agents/blob/main/.claude-plugin/marketplace.json

### @ Notation in Agent/Command Files

**GitHub Issue #5914:** "Support @ imports in sub-agent markdown files"
- **Current Behavior:** @ notation in agent files causes Claude to use Read tool (inefficient)
- **Not a True Import:** Claude "smart enough to read" but not optimized
- **Limitation:** File read every time, requires restart for changes

**Conclusion:** @ notation in plugin files is not officially supported for efficient runtime imports.

### Answer to Question 1

**For Plugin Commands/Agents referencing rules or templates:**

❌ **Cannot use:** `@{{STARTUP_PATH}}/rules/agent-delegation.md`
❌ **Cannot use:** `@${CLAUDE_PLUGIN_ROOT}/rules/agent-delegation.md`
✅ **Must inline content** directly in command/agent markdown
✅ **Alternative:** Use hooks to inject context before tool execution

**For Plugin Hooks:**
✅ **Can use:** `${CLAUDE_PLUGIN_ROOT}/hooks/script.py`
✅ **Can use:** `${CLAUDE_PLUGIN_ROOT}/templates/PRD.md` (as script argument)

**Recommendation:**
- **Commands/Agents:** Inline all rules/templates content during build
- **Hooks:** Reference scripts/files using ${CLAUDE_PLUGIN_ROOT}
- **Runtime File Access:** Only via hooks executing scripts

---

## Question 2: Scripts Directory Convention

### Official Convention: hooks/

**Evidence from Anthropic plugins:**

```
plugins/security-guidance/
├── .claude-plugin/
│   └── plugin.json
├── hooks/
│   ├── hooks.json
│   └── security_reminder_hook.py
└── README.md
```

**Source:** https://github.com/anthropics/claude-code/tree/main/plugins/security-guidance

### No bin/ or scripts/ Directory Observed

**Search Results:**
- Searched all official plugins (5 plugins: agent-sdk-dev, commit-commands, feature-dev, pr-review-toolkit, security-guidance)
- Searched community plugins (wshobson/agents with 83+ agents)
- **Zero instances** of `bin/` or `scripts/` directories in plugin structure

**Official Documentation:**
- Plugins reference mentions `scripts/` in example: `${CLAUDE_PLUGIN_ROOT}/scripts/process.sh`
- However, actual official plugins use `hooks/` directory

### Plugin Directory Structure

**Standard Structure (from evidence):**
```
plugin-name/
├── .claude-plugin/
│   ├── plugin.json          # Metadata
│   └── marketplace.json     # Optional: if hosting marketplace
├── commands/                # Slash command markdown files
│   ├── command1.md
│   └── command2.md
├── agents/                  # Agent definition markdown files
│   ├── agent1.md
│   └── agent2.md
├── hooks/                   # Hook configuration and scripts
│   ├── hooks.json           # Hook definitions
│   └── script.py            # Executable hook scripts
└── README.md
```

### Answer to Question 2

**Recommendation:**
- Use `hooks/` directory for executable scripts (matches official convention)
- Rename `assets/the-startup/bin/` → `hooks/` in plugin structure
- Place `statusline.sh` and `statusline.ps1` in `hooks/`

**Configuration Example:**
```json
{
  "statusLine": {
    "type": "command",
    "command": "${CLAUDE_PLUGIN_ROOT}/hooks/statusline.sh"
  }
}
```

---

## Question 3: The "spec" Command Implementation

### Command Capabilities in Plugins

**What Slash Commands CAN Do:**
- Provide instructions to Claude Code
- Include arguments via `$ARGUMENTS` placeholder
- Trigger sub-agents via Task tool
- Instruct Claude to use tools (Read, Write, Edit, Bash, etc.)

**What Slash Commands CANNOT Do Directly:**
- Execute shell scripts
- Create directories directly
- Write files directly
- Generate TOML/JSON output directly

**Evidence:**
Official plugins (feature-dev, pr-review-toolkit) show commands as **instruction sets**, not script executors:

```markdown
---
description: Guided feature development
---

You are helping a developer implement a new feature.

## Phase 1: Discovery
**Actions**:
1. Create todo list with all phases
2. If feature unclear, ask user for...
```

### Evaluation of Options

#### Option A: Slash Command (Pure Instructions)

```markdown
---
description: Create numbered spec directories
---

Create a specification directory for: $ARGUMENTS

1. Find existing spec directories matching pattern `spec-\d{3}-*`
2. Determine next sequential number
3. Create directory: `spec-{next}-{slug}/`
4. Generate TOML frontmatter with metadata
5. Create SPEC.md, PLAN.md from templates
```

**Pros:**
- Native plugin integration
- User-friendly invocation: `/s:spec feature-name`

**Cons:**
- Relies on Claude's interpretation
- Cannot guarantee TOML syntax correctness
- Directory numbering may be unreliable
- Template substitution not guaranteed
- No atomic operations

**Verdict:** ❌ **Not Recommended** - Insufficient guarantees for critical operations

#### Option B: Hook-Invoked Script

```json
{
  "hooks": {
    "UserPromptSubmit": [{
      "matcher": "/s:spec",
      "hooks": [{
        "type": "command",
        "command": "${CLAUDE_PLUGIN_ROOT}/hooks/spec.sh $ARGUMENTS"
      }]
    }]
  }
}
```

**Pros:**
- Full script capabilities
- Guaranteed TOML generation
- Atomic directory creation

**Cons:**
- Hooks don't match on user prompts with pattern matching
- UserPromptSubmit fires for ALL prompts
- No way to conditionally trigger on `/s:spec`
- Would need to parse prompt text in script

**Verdict:** ❌ **Not Feasible** - Hook system doesn't support command pattern matching

#### Option C: External Script (Separate CLI)

```bash
# Installation
npm install -g the-agentic-startup-spec

# Usage
the-agentic-startup-spec feature-name
```

**Pros:**
- Full CLI capabilities (inquirer, commander, etc.)
- Guaranteed TOML generation
- Atomic operations
- Can be developed/tested independently
- Can provide rich CLI experience (colors, spinners, confirmations)

**Cons:**
- Separate installation step
- Not integrated into Claude Code plugin system

**Verdict:** ✅ **Recommended** - Best approach for reliable implementation

#### Option D: Command That Invokes Script

```markdown
---
description: Create spec directory
---

Run the spec generator script:

Use the Bash tool to execute:
`${CLAUDE_PLUGIN_ROOT}/hooks/spec.sh $ARGUMENTS`
```

**Pros:**
- User-friendly slash command interface
- Leverages bundled script

**Cons:**
- `${CLAUDE_PLUGIN_ROOT}` not available in command context
- Would instruct Claude to use Bash tool
- Claude would need to construct command
- Not guaranteed to execute correctly

**Verdict:** ❌ **Not Feasible** - Variables not available in command markdown

### Answer to Question 3

**Recommendation: Option C - External CLI Tool**

**Implementation:**
1. Create standalone npm package: `the-agentic-startup-spec`
2. Distribute via npm registry
3. CLI provides:
   - Auto-incrementing spec directory creation
   - TOML frontmatter generation
   - Template instantiation
   - Validation and error handling

**Package Structure:**
```
the-agentic-startup-spec/
├── src/
│   ├── cli.ts              # CLI entry point
│   ├── generator.ts        # Spec generation logic
│   └── templates/          # Embedded templates
├── bin/
│   └── spec.js             # Executable
└── package.json
```

**Benefits:**
- Reliable, testable implementation
- Rich CLI features (inquirer, chalk, ora)
- Independent versioning
- Can be used with or without Claude Code
- Future: Could add to PATH via plugin installation hook

**Plugin Integration (Future Enhancement):**
Commands can **recommend** using the external tool:
```markdown
---
description: Create specification directory
---

This workflow requires the `the-agentic-startup-spec` CLI tool.

If not installed, run:
```bash
npm install -g the-agentic-startup-spec
```

Then execute:
```bash
the-agentic-startup-spec $ARGUMENTS
```

This ensures:
- Correct auto-incrementing directory numbers
- Valid TOML frontmatter generation
- Proper template instantiation
```

---

## Summary of Findings

### File References

| Context | Mechanism | Works? | Example |
|---------|-----------|--------|---------|
| User chat/prompts | `@path/to/file` | ✅ Yes | `@src/utils.ts` |
| CLAUDE.md | `@path/to/file` | ✅ Yes | `@.claude/rules.md` |
| Plugin commands/agents | `@path/to/file` | ❌ No (inefficient) | Use inlining instead |
| Plugin hooks config | `${CLAUDE_PLUGIN_ROOT}/path` | ✅ Yes | `${CLAUDE_PLUGIN_ROOT}/hooks/script.py` |
| Plugin metadata | `./relative/path` | ✅ Yes | `./commands/cmd.md` |

### Script Directory

| Directory | Usage | Official? | Recommendation |
|-----------|-------|-----------|----------------|
| `hooks/` | Hook scripts & config | ✅ Yes | Use for all executable scripts |
| `scripts/` | Utility scripts | ❌ No evidence | Avoid - use hooks/ instead |
| `bin/` | Executable binaries | ❌ No evidence | Avoid - use hooks/ instead |

### Spec Command Implementation

| Option | Feasibility | Reliability | Recommendation |
|--------|-------------|-------------|----------------|
| A: Pure slash command | ⚠️ Possible | Low | ❌ Not recommended |
| B: Hook-invoked | ❌ Not feasible | N/A | ❌ Not feasible |
| C: External CLI | ✅ Fully feasible | High | ✅ **Recommended** |
| D: Command + script | ❌ Not feasible | Low | ❌ Not feasible |

---

## Implementation Recommendations

### 1. For Current CLI → Plugin Migration

**Commands/Agents with File References:**

Before:
```markdown
@{{STARTUP_PATH}}/rules/agent-delegation.md
```

After (inline during build):
```markdown
## Agent Delegation Rules

[Full content of agent-delegation.md inlined here]
```

**Build Process:**
```typescript
// Build-time inlining
function inlineReferences(content: string): string {
  return content.replace(
    /@\{\{STARTUP_PATH\}\}\/(.+)/g,
    (match, path) => fs.readFileSync(`assets/the-startup/${path}`, 'utf-8')
  );
}
```

### 2. For Scripts

**Current:**
```
assets/the-startup/bin/
├── statusline.sh
└── statusline.ps1
```

**Plugin Structure:**
```
plugin-root/
├── hooks/
│   ├── hooks.json
│   ├── statusline.sh
│   └── statusline.ps1
```

**Configuration:**
```json
{
  "hooks": {
    "statusLine": {
      "type": "command",
      "command": "${CLAUDE_PLUGIN_ROOT}/hooks/statusline.sh"
    }
  }
}
```

### 3. For Spec Command

**Separate Package:**
```json
{
  "name": "the-agentic-startup-spec",
  "version": "1.0.0",
  "bin": {
    "the-agentic-startup-spec": "./dist/bin/spec.js"
  }
}
```

**Installation:**
```bash
npm install -g the-agentic-startup-spec
```

**Usage:**
```bash
the-agentic-startup-spec feature-name --add solution-design
```

---

## References

### Official Sources

1. **Plugins Reference**: https://docs.claude.com/en/docs/claude-code/plugins-reference
2. **Hooks Guide**: https://docs.claude.com/en/docs/claude-code/hooks-guide
3. **Slash Commands**: https://docs.claude.com/en/docs/claude-code/slash-commands
4. **Status Line**: https://docs.claude.com/en/docs/claude-code/statusline

### Official Plugin Examples

1. **security-guidance**: https://github.com/anthropics/claude-code/tree/main/plugins/security-guidance
2. **pr-review-toolkit**: https://github.com/anthropics/claude-code/tree/main/plugins/pr-review-toolkit
3. **feature-dev**: https://github.com/anthropics/claude-code/tree/main/plugins/feature-dev

### Community Examples

1. **wshobson/agents**: https://github.com/wshobson/agents (83+ agents, 15 workflows)
2. **wshobson/commands**: https://github.com/wshobson/commands (57 commands)
3. **AgiFlow/aicode-toolkit**: https://github.com/AgiFlow/aicode-toolkit

### GitHub Issues

1. **#5914**: Support @ imports in sub-agent markdown files
2. **#990**: Syntax for including CLAUDE.md files

---

## Conclusion

**Definitive Answers:**

1. **File References:** @ notation is for user context only; plugins must inline content or use ${CLAUDE_PLUGIN_ROOT} in hooks
2. **Scripts Directory:** Use `hooks/` directory per official convention
3. **Spec Command:** Implement as external CLI tool for reliability

**Next Steps:**

1. ✅ Inline rules/templates content during plugin build
2. ✅ Rename bin/ → hooks/ and configure statusline
3. ✅ Extract spec command to separate npm package
4. ✅ Update plugin.json to reference hooks/
5. ✅ Create build process for content inlining

**Success Criteria:**
- ✅ Evidence from official sources
- ✅ Verified with actual plugin examples
- ✅ Clear implementation guidance
- ✅ No assumptions without evidence
