# Example: File Coordination for Parallel Agents

This example shows how to prevent file path collisions when multiple agents create documentation in parallel.

## Scenario

**User Request:** "Document our codebase patterns for authentication, caching, and error handling"

## Initial Decomposition (Incorrect)

```
❌ WRONG APPROACH

Activities:
1. Document authentication patterns
2. Document caching patterns
3. Document error handling patterns

All parallel, all OUTPUT: "docs/patterns/[pattern-name].md"

Problem: What if agents choose same filename?
- Agent 1 might create: docs/patterns/auth.md
- Agent 2 might create: docs/patterns/cache.md
- Agent 3 might create: docs/patterns/error.md

OR worse:
- Agent 1: docs/patterns/authentication.md
- Agent 2: docs/patterns/authentication-patterns.md
- Both trying to document auth? Collision!

Result: Ambiguous, potential collisions, inconsistent naming
```

## Correct Decomposition (File Coordination)

```
✅ CORRECT APPROACH

Activities:
1. Document authentication patterns
   - OUTPUT: docs/patterns/authentication-flow.md (EXPLICIT PATH)

2. Document caching patterns
   - OUTPUT: docs/patterns/caching-strategy.md (EXPLICIT PATH)

3. Document error handling patterns
   - OUTPUT: docs/patterns/error-handling.md (EXPLICIT PATH)

File Coordination Check:
✅ All paths explicit and unique
✅ No ambiguity in naming
✅ No collision risk

Status: SAFE FOR PARALLEL EXECUTION
```

## Agent Prompts with File Coordination

### Agent 1: Authentication Pattern

```
DISCOVERY_FIRST: Before starting, check existing documentation:
    - List docs/patterns/ directory
    - Search for existing auth-related files
    - Note naming conventions used
    - Check if authentication-flow.md already exists

FOCUS: Document authentication patterns discovered in codebase
    - JWT token generation and validation
    - Password hashing (bcrypt usage)
    - Session management approach
    - Protected route patterns
    - Error responses for auth failures

EXCLUDE:
    - Don't document caching (Agent 2 handles this)
    - Don't document error handling generally (Agent 3 handles this)
    - Don't create multiple files (single document)
    - Don't modify existing authentication files if they exist

CONTEXT: Documenting patterns found in src/middleware/auth.*, src/routes/auth.*
    - Focus on how authentication works, not implementation details
    - Use pattern template: docs/templates/pattern-template.md
    - Follow existing documentation style

OUTPUT: EXACTLY this path: docs/patterns/authentication-flow.md
    - If file exists: STOP and report (don't overwrite)
    - If file doesn't exist: Create new
    - Use pattern template structure
    - Include code examples from codebase

SUCCESS: Authentication patterns documented
    - File created at exact path specified
    - Follows pattern template
    - Includes JWT, bcrypt, and session patterns
    - Code examples are accurate
    - No collision with other agents

TERMINATION:
    - Documentation complete at specified path
    - File already exists (report to user)
    - No authentication patterns found in codebase
```

### Agent 2: Caching Pattern

```
DISCOVERY_FIRST: Before starting, check existing documentation:
    - List docs/patterns/ directory
    - Search for existing cache-related files
    - Note naming conventions
    - Check if caching-strategy.md already exists

FOCUS: Document caching patterns discovered in codebase
    - Redis usage patterns
    - Cache key naming conventions
    - TTL (time-to-live) strategies
    - Cache invalidation approaches
    - What gets cached and why

EXCLUDE:
    - Don't document authentication (Agent 1 handles this)
    - Don't document error handling (Agent 3 handles this)
    - Don't create multiple files (single document)
    - Don't overlap with Agent 1's work

CONTEXT: Documenting patterns found in src/cache/*, src/services/*
    - Focus on caching strategy, not Redis API details
    - Use pattern template
    - Follow existing documentation style

OUTPUT: EXACTLY this path: docs/patterns/caching-strategy.md
    - If file exists: STOP and report (don't overwrite)
    - If file doesn't exist: Create new
    - Use pattern template structure
    - Include code examples from codebase

SUCCESS: Caching patterns documented
    - File created at exact path specified
    - Follows pattern template
    - Includes Redis patterns and cache strategies
    - No collision with other agents

TERMINATION:
    - Documentation complete at specified path
    - File already exists (report)
    - No caching patterns found
```

### Agent 3: Error Handling Pattern

```
DISCOVERY_FIRST: Before starting, check existing documentation:
    - List docs/patterns/ directory
    - Search for existing error-related files
    - Note naming conventions
    - Check if error-handling.md already exists

FOCUS: Document error handling patterns discovered in codebase
    - Error class hierarchy
    - HTTP status code mapping
    - Error response format
    - Logging strategy for errors
    - Recovery patterns

EXCLUDE:
    - Don't document authentication (Agent 1's domain)
    - Don't document caching (Agent 2's domain)
    - Don't create multiple files
    - Don't overlap with other agents

CONTEXT: Documenting patterns found in src/errors/*, src/middleware/error.*
    - Focus on error handling strategy and patterns
    - Use pattern template
    - Follow existing documentation style

OUTPUT: EXACTLY this path: docs/patterns/error-handling.md
    - If file exists: STOP and report (don't overwrite)
    - If file doesn't exist: Create new
    - Use pattern template structure
    - Include code examples

SUCCESS: Error handling patterns documented
    - File created at exact path specified
    - Follows pattern template
    - Includes error classes and response patterns
    - No collision with other agents

TERMINATION:
    - Documentation complete at specified path
    - File already exists (report)
    - No error handling patterns found
```

## File Coordination Checklist

Before launching these agents:

```
📋 File Coordination Pre-Flight Check

Agent 1 OUTPUT: docs/patterns/authentication-flow.md
Agent 2 OUTPUT: docs/patterns/caching-strategy.md
Agent 3 OUTPUT: docs/patterns/error-handling.md

✅ All paths are explicit (no ambiguity)
✅ All paths are unique (no two agents write same file)
✅ All paths follow convention (docs/patterns/[name].md)
✅ All agents instructed to check if file exists first
✅ All agents instructed to STOP if collision detected

File Collision Risk: NONE
Safe to launch in parallel: YES
```

## Execution Flow

### Launch All Three Agents in Parallel

```
🚀 Launching 3 parallel documentation agents

Agent 1: Authentication Pattern → RUNNING
  TARGET: docs/patterns/authentication-flow.md

Agent 2: Caching Pattern → RUNNING
  TARGET: docs/patterns/caching-strategy.md

Agent 3: Error Handling Pattern → RUNNING
  TARGET: docs/patterns/error-handling.md

File coordination: ✅ All unique paths
Collision risk: ✅ None
Parallel safety: ✅ Confirmed
```

### Monitoring for Collisions

```
⏳ Agents running...

[Agent 1] Checking: docs/patterns/authentication-flow.md → NOT EXISTS
[Agent 1] Safe to create → PROCEEDING

[Agent 2] Checking: docs/patterns/caching-strategy.md → NOT EXISTS
[Agent 2] Safe to create → PROCEEDING

[Agent 3] Checking: docs/patterns/error-handling.md → NOT EXISTS
[Agent 3] Safe to create → PROCEEDING

No collisions detected. All agents proceeding independently.
```

### Completion

```
Agent 1: COMPLETE ✅ (22 minutes)
  Created: docs/patterns/authentication-flow.md (3.2 KB)

Agent 2: COMPLETE ✅ (18 minutes)
  Created: docs/patterns/caching-strategy.md (2.8 KB)

Agent 3: COMPLETE ✅ (25 minutes)
  Created: docs/patterns/error-handling.md (4.1 KB)

All agents complete. No collisions occurred.
```

## Results

```
📁 docs/patterns/
├── authentication-flow.md    ✅ Created by Agent 1
├── caching-strategy.md        ✅ Created by Agent 2
└── error-handling.md          ✅ Created by Agent 3
```

**Total time:** 25 minutes (parallel)
**Sequential would take:** 65 minutes
**Time saved:** 40 minutes (61% faster)

## Alternative Coordination Strategies

### Strategy 1: Directory Separation

If agents might create multiple files:

```
Agent 1 OUTPUT: docs/patterns/authentication/
  - flow.md
  - jwt-tokens.md
  - session-management.md

Agent 2 OUTPUT: docs/patterns/caching/
  - redis-usage.md
  - invalidation.md
  - key-naming.md

Agent 3 OUTPUT: docs/patterns/error-handling/
  - error-classes.md
  - response-format.md
  - recovery-patterns.md
```

**Result:** Each agent owns a directory, no file collisions possible.

### Strategy 2: Timestamp-Based Naming

For logs or reports that accumulate:

```
Agent 1 OUTPUT: logs/auth-analysis-2025-01-24-10-30-00.md
Agent 2 OUTPUT: logs/cache-analysis-2025-01-24-10-30-00.md
Agent 3 OUTPUT: logs/error-analysis-2025-01-24-10-30-00.md
```

**Result:** Timestamps ensure uniqueness even if same topic.

### Strategy 3: Agent ID Namespacing

When dynamic number of agents:

```
For each module in [moduleA, moduleB, moduleC, moduleD]:
  Launch agent with OUTPUT: analysis/module-${MODULE_NAME}.md

Results:
- analysis/module-moduleA.md
- analysis/module-moduleB.md
- analysis/module-moduleC.md
- analysis/module-moduleD.md
```

**Result:** Template-based naming prevents collisions.

## What Could Go Wrong (Anti-Patterns)

### ❌ Anti-Pattern 1: Ambiguous Paths

```
BAD:
Agent 1 OUTPUT: "Create pattern documentation"
Agent 2 OUTPUT: "Document caching patterns"

Problem: Where exactly? What filename?
Result: Agents might choose same name or wrong location
```

### ❌ Anti-Pattern 2: Overlapping Domains

```
BAD:
Agent 1 FOCUS: "Document authentication and security"
Agent 2 FOCUS: "Document security patterns"

Problem: Both might document auth security!
Result: Duplicate or conflicting documentation
```

**Fix:** Clear FOCUS boundaries, explicit EXCLUDE

### ❌ Anti-Pattern 3: No Existence Check

```
BAD:
OUTPUT: docs/patterns/auth.md
(No instruction to check if exists)

Problem: If file exists, agent might overwrite
Result: Lost documentation
```

**Fix:** Always include existence check in FOCUS or OUTPUT

### ❌ Anti-Pattern 4: Generic Filenames

```
BAD:
Agent 1 OUTPUT: docs/patterns/pattern.md
Agent 2 OUTPUT: docs/patterns/patterns.md
Agent 3 OUTPUT: docs/patterns/pattern-doc.md

Problem: All similar, confusing, might collide
Result: Unclear which agent created what
```

**Fix:** Descriptive, specific filenames

## File Coordination Best Practices

### 1. Explicit Paths

✅ **Always specify exact OUTPUT path:**
```
OUTPUT: docs/patterns/authentication-flow.md
NOT: "Document authentication patterns"
```

### 2. Unique Names

✅ **Ensure all agents have unique filenames:**
```
Before launching:
- Agent 1: authentication-flow.md
- Agent 2: caching-strategy.md
- Agent 3: error-handling.md
✅ All unique → SAFE
```

### 3. Existence Checks

✅ **Instruct agents to check before creating:**
```
OUTPUT: docs/patterns/authentication-flow.md
- If file exists: STOP and report (don't overwrite)
- If file doesn't exist: Create new
```

### 4. Clear Boundaries

✅ **Use EXCLUDE to prevent overlap:**
```
Agent 1:
FOCUS: Document authentication
EXCLUDE: Don't document caching (Agent 2) or errors (Agent 3)

Agent 2:
FOCUS: Document caching
EXCLUDE: Don't document auth (Agent 1) or errors (Agent 3)
```

### 5. Validation

✅ **Check coordination before launching:**
```
Checklist:
- [ ] All paths explicit
- [ ] All paths unique
- [ ] All agents have existence checks
- [ ] Clear FOCUS boundaries
- [ ] No overlap in domains
```

## Integration with Documentation Skill

When these agents complete, documentation skill may activate:

```
Agent 1 completes → Creates authentication-flow.md
  ↓
Documentation skill notices "pattern" created
  ↓
Checks: Is this in correct location? (docs/patterns/ ✅)
Checks: Does it follow template? (✅)
Checks: Should it be cross-referenced? (Yes)
  ↓
Documentation skill adds cross-references:
- Links to related authentication interfaces
- Links to domain rules about user permissions
```

**Coordination:** Agent-delegation ensures unique paths, documentation skill ensures quality and linking.

## Lessons Learned

### What Worked

✅ **Explicit paths:** Zero ambiguity, zero collisions
✅ **Existence checks:** Prevented accidental overwrites
✅ **Clear boundaries:** No domain overlap
✅ **Parallel execution:** 61% time savings

### What to Watch For

⚠️ **Naming conventions:** Ensure consistency (kebab-case vs snake_case)
⚠️ **Template usage:** All agents should use same template
⚠️ **Directory structure:** Verify docs/patterns/ exists before launching

## Reusable Coordination Template

For any parallel file creation:

```
1. List all files to be created
2. Verify all paths are unique
3. Add existence checks to OUTPUT
4. Use EXCLUDE to prevent overlap
5. Launch in parallel with confidence
```

**This prevents:**
- File path collisions
- Accidental overwrites
- Domain overlap
- Naming inconsistencies
