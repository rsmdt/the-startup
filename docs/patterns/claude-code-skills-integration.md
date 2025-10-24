# Claude Code Skills Integration Pattern

> **Category:** Technical Pattern
> **Last Updated:** 2025-10-24
> **Status:** Active

## Purpose

This document describes how Claude Code Skills are integrated into The Agentic Startup framework, replacing explicit documentation instructions with autonomous skill activation.

## Context

**When to use this pattern:**
- Building multi-agent orchestration systems
- Creating reusable capabilities across commands
- Implementing autonomous behavior based on context
- Reducing duplication of instructions across multiple workflows

**When NOT to use this pattern:**
- User-invoked operations (use slash commands instead)
- Simple one-off tasks that won't be reused
- Features that require explicit user control

## What Are Claude Code Skills?

Skills are **model-invoked** capabilities that Claude autonomously activates based on:
- Description relevance to the current request
- Presence of trigger terms in user input or agent context
- Context matching with the skill's purpose

**Key Difference from Slash Commands:**
- **Slash Commands:** User explicitly invokes (e.g., `/start:analyze`)
- **Skills:** Claude autonomously decides to activate based on description

## SKILL.md Structure

Every skill requires this exact YAML frontmatter:

```yaml
---
name: skill-name-here
description: What the Skill does and when Claude should use it (max 1024 chars)
allowed-tools: Read, Write, Edit, Grep, Glob (optional)
---

[Markdown instructions for Claude when skill is activated]
```

### Validation Rules

- `name`: lowercase + hyphens + numbers only, max 64 characters
- `description`: **MUST include WHAT + WHEN** (functionality + activation triggers)
- `allowed-tools`: Optional comma-separated list of permitted tools

### Description Best Practices

❌ **Bad:** "Helps with documentation"
✅ **Good:** "Document business rules, technical patterns, and service interfaces discovered during analysis. Use when you find reusable patterns, external integrations, domain rules, or when checking existing documentation before creating new files."

**Formula:** `[What it does] + [When to use it] + [Trigger terms]`

## File Organization

```
skill-name/
├── SKILL.md           # Required: frontmatter + instructions
├── reference.md       # Optional: detailed documentation
├── examples.md        # Optional: usage examples
├── scripts/           # Optional: helper utilities
└── templates/         # Optional: file templates
```

**Progressive Disclosure:**
- SKILL.md: Always loaded when skill activates
- Other files: Loaded only when contextually relevant (token efficiency)

## Storage Locations

Skills can be stored in three locations:

1. **Plugin Skills:** `plugins/start/skills/` → Bundled with plugin distribution
2. **Project Skills:** `.claude/skills/` → Team-shared via git
3. **Personal Skills:** `~/.claude/skills/` → Individual workflows

Our implementation uses plugin skills for framework-wide capabilities.

## Documentation Skill Implementation

### Created Files

```
plugins/start/skills/documentation/
├── SKILL.md                           # Main skill with activation logic
├── reference.md                       # Detailed protocols and edge cases
└── templates/
    ├── pattern-template.md            # For technical patterns
    ├── interface-template.md          # For external integrations
    └── domain-template.md             # For business rules
```

### Trigger Terms

The skill activates when Claude encounters:
- "pattern", "reusable pattern", "technical pattern"
- "interface", "integration", "external service", "API"
- "domain rule", "business rule", "workflow"
- "document", "documentation"
- "check existing", "avoid duplication"

### What It Does

1. **Checks for existing documentation** (deduplication)
2. **Categorizes correctly** (domain/patterns/interfaces)
3. **Uses appropriate template**
4. **Applies naming conventions** (descriptive, searchable)
5. **Creates cross-references** between related docs
6. **Reports what was created/updated**

## Command Updates for Skill Activation

### Before: Explicit Instructions

```markdown
**Documentation Phase**:
- Create/update docs/domain/ for business rules and domain patterns discovered
- Create/update docs/patterns/ for reusable solution patterns identified
- Create/update docs/interfaces/ for external service contracts needed
```

### After: Trigger Language

```markdown
**Documentation Phase**:
- Document any discovered patterns, interfaces, or domain rules for future reference
```

### Why This Works

The natural language "discovered patterns", "interfaces", and "domain rules" contains the trigger terms that activate the documentation skill autonomously.

## Updated Command Files

1. **`/start:specify`** (plugins/start/commands/specify.md)
   - Removed explicit docs/ structure instructions
   - Changed to trigger language: "discovered domain rules, patterns, or external integrations"
   - Simplified Supporting Documentation section

2. **`/start:analyze`** (plugins/start/commands/analyze.md)
   - Simplified Documentation Structure section
   - Reduced Important Notes to trigger language

3. **`cycle-pattern.md`** (plugins/start/rules/cycle-pattern.md)
   - Updated Documentation Phase to use triggers
   - Updated Review Phase to report documented items
   - Updated "Ask yourself" checklist

## Benefits of Skills Pattern

### 1. DRY (Don't Repeat Yourself)
- Documentation logic defined once in skill
- No repeating instructions across multiple commands
- Single source of truth for protocols

### 2. Autonomous Activation
- Claude decides when documentation is needed
- Works in ANY context, not just specific commands
- Natural language triggers instead of explicit calls

### 3. Consistency
- Same deduplication checks every time
- Same categorization logic
- Same template usage

### 4. Extensibility
- Works in ad-hoc conversations
- Works with future commands automatically
- Works for ANY agent, not just orchestrators

### 5. Context Efficiency
- Progressive disclosure loads details only when needed
- Keeps command files focused on orchestration
- Reduces token usage in common cases

## Future Skills Candidates

Based on the analysis, these capabilities would benefit from Skills extraction:

### 1. Specification Review Skill ⭐⭐

**Purpose:** Validate implementation against spec documents (PRD/SDD/PLAN)

**Trigger terms:** "check spec", "verify requirements", "SDD compliance", "context drift"

**Why a skill:** Used by both `/start:specify` and `/start:implement`, plus ad-hoc validation

### 2. Agent Delegation Skill ⭐

**Purpose:** Intelligent task decomposition and parallel execution patterns

**Trigger terms:** "delegate", "launch agents", "parallel execution", "break down task"

**Why a skill:** Framework-wide pattern used across ALL orchestration commands

### 3. Quality Gates Skill ⭐

**Purpose:** Execute DOR/DOD/TASK-DOD validations

**Trigger terms:** "validate", "quality gate", "definition of done", "readiness check"

**Why a skill:** Flexible invocation (used when gates exist, skipped when they don't)

## Testing Strategy

### Activation Testing

1. **Direct trigger:** Ask questions with trigger terms
   - "Should I document this pattern?"
   - "I found a reusable caching strategy"
   - "This external API needs documentation"

2. **Contextual trigger:** Work on tasks that naturally involve documentation
   - Analyzing codebase patterns
   - Designing integrations
   - Discovering business rules

3. **Discovery verification:** Ask Claude
   - "What skills are available?"
   - Should list "documentation" skill

### Validation Testing

1. **Deduplication:** Ensure existing docs are checked first
2. **Categorization:** Verify correct docs/ category selection
3. **Template usage:** Confirm appropriate template is used
4. **Cross-references:** Check links between related docs

## Integration with Plugin Installation

The documentation skill is distributed with the `@the-startup` plugin:

```
plugins/start/
├── skills/documentation/     # Skill files
├── commands/                 # Command files (updated to trigger skill)
└── rules/                    # Rule files (updated to trigger skill)
```

When users install the plugin:
1. Skill copies to `.claude/skills/documentation/`
2. Commands copy to `.claude/commands/`
3. Skill becomes available project-wide
4. Commands naturally trigger skill via language

## Common Patterns

### Pattern 1: Discovery → Documentation

```
Agent discovers pattern
  ↓
Claude recognizes "reusable pattern" (trigger)
  ↓
Documentation skill activates
  ↓
Checks existing docs
  ↓
Categorizes (domain/patterns/interfaces)
  ↓
Creates or updates documentation
  ↓
Reports what was done
```

### Pattern 2: Explicit Check

```
User or orchestrator mentions "check existing documentation"
  ↓
Documentation skill activates
  ↓
Searches docs/ hierarchy
  ↓
Reports findings
  ↓
Suggests create vs update
```

## Edge Cases

### What if Multiple Skills Match?

Claude chooses based on:
- Description relevance to current context
- Specificity of trigger terms
- Most recent skill definitions

**Best practice:** Use distinct trigger terms for different skills

### What if Skill Doesn't Activate?

**Debugging steps:**
1. Check description has clear trigger terms
2. Verify YAML frontmatter is valid
3. Ensure file is in correct location
4. Try more explicit trigger language
5. Use `claude --debug` (if available)

### What if Commands Conflict with Skill?

**Resolution:**
- Skills are supplementary to commands
- Commands orchestrate workflow
- Skills provide specialized capabilities
- If conflict occurs, command takes precedence

## Performance Considerations

### Token Efficiency

- SKILL.md loaded on activation (~2KB)
- reference.md loaded only if needed (~8KB)
- templates/ loaded only if needed (~6KB)
- Total skill: ~16KB maximum, ~2KB typical

### Activation Speed

- Skill discovery: Automatic (no performance cost)
- Skill activation: Near-instantaneous
- Progressive disclosure: Loads details on-demand

## Related Documentation

- **Skills Official Docs:** https://docs.claude.com/en/docs/claude-code/skills
- **Plugin Structure:** See plugins/start/README.md
- **Command Architecture:** See docs/patterns/slash-commands.md (if exists)

## References

- [Claude Code Skills Documentation](https://docs.claude.com/en/docs/claude-code/skills)
- [The Agentic Startup Framework](../README.md)

## Version History

| Date | Change | Author |
|------|--------|--------|
| 2025-10-24 | Initial pattern documentation | Claude (Documentation Skill Analysis) |
