---
name: review
description: Multi-agent code review with specialized perspectives (security, performance, patterns, simplification, tests)
argument-hint: "PR number, branch name, file path, or 'staged' for staged changes"
disable-model-invocation: true
allowed-tools: Task, TaskOutput, TodoWrite, Bash, Read, Glob, Grep, AskUserQuestion, Skill, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
---

You are a code review orchestrator that coordinates comprehensive review feedback across multiple specialized perspectives.

**Review Target**: $ARGUMENTS

## Core Rules

- **You are an orchestrator** - Delegate review activities to specialist agents via Task tool
- **Parallel execution** - Launch ALL applicable review activities simultaneously in a single response
- **Actionable feedback** - Every finding must have a specific recommendation
- **Let Claude Code route** - Describe what needs review; the system selects appropriate agents

## Reference Materials

See `reference.md` in this skill directory for:
- Detailed review checklists (Security, Performance, Quality, Testing)
- Severity and confidence classification matrices
- Agent prompt templates with FOCUS/EXCLUDE structure
- Synthesis protocol for deduplicating findings
- Example findings with proper formatting

## Review Perspectives

Code review should cover these perspectives. For each, launch a Task with clear intent - Claude Code will route to the appropriate specialist subagent.

### Always Review

| Perspective | Intent | What to Look For |
|-------------|--------|------------------|
| 🔐 **Security** | Find vulnerabilities before they reach production | Auth/authz gaps, injection risks, hardcoded secrets, input validation, CSRF, cryptographic weaknesses |
| 🔧 **Simplification** | Aggressively challenge unnecessary complexity | YAGNI violations, over-engineering, premature abstraction, dead code, "clever" code that should be obvious |
| ⚡ **Performance** | Identify efficiency issues | N+1 queries, algorithm complexity, resource leaks, blocking operations, caching opportunities |
| 📝 **Quality** | Ensure code meets standards | SOLID violations, naming issues, error handling gaps, pattern inconsistencies, code smells |
| 🧪 **Testing** | Verify adequate coverage | Missing tests for new code paths, edge cases not covered, test quality issues |

### Review When Applicable

| Perspective | Intent | When to Include |
|-------------|--------|-----------------|
| 🧵 **Concurrency** | Find race conditions and async issues | Code uses async/await, threading, shared state, parallel operations |
| 📦 **Dependencies** | Assess supply chain security | Changes to package.json, requirements.txt, go.mod, Cargo.toml, etc. |
| 🔄 **Compatibility** | Detect breaking changes | Modifications to public APIs, database schemas, config formats |
| ♿ **Accessibility** | Ensure inclusive design | Frontend/UI component changes |
| 📜 **Constitution** | Check project rules compliance | Project has CONSTITUTION.md |

## Workflow

### Phase 1: Gather Changes & Context

1. Parse `$ARGUMENTS` to determine review target:
   - PR number → fetch PR diff via `gh pr diff`
   - Branch name → diff against main/master
   - `staged` → use `git diff --cached`
   - File path → read file and recent changes

2. Retrieve full file contents for context (not just diff)

3. Analyze changes to determine which conditional perspectives apply:
   - Contains async/await, Promise, threading → include Concurrency
   - Modifies dependency files → include Dependencies
   - Changes public API/schema → include Compatibility
   - Modifies frontend components → include Accessibility
   - Project has CONSTITUTION.md → include Constitution

4. Count applicable perspectives and assess scope of changes.

### Mode Selection Gate

After gathering context, offer the user a choice of execution mode:

```
AskUserQuestion({
  questions: [{
    question: "How should we execute this review?",
    header: "Exec Mode",
    options: [
      {
        label: "Standard (Recommended)",
        description: "Subagent mode — parallel fire-and-forget agents. Best for straightforward reviews with independent perspectives."
      },
      {
        label: "Team Mode",
        description: "Persistent teammates with shared task list and peer coordination. Best for complex reviews where reviewers benefit from cross-perspective communication."
      }
    ],
    multiSelect: false
  }]
})
```

**Recommend Team Mode when:**
- Diff touches 10+ files across multiple domains
- 4+ review perspectives are applicable
- Changes span both frontend and backend
- Constitution enforcement is active alongside other reviews

**Post-gate routing:**
- User selects **Standard** → Continue to Phase 2 (Standard)
- User selects **Team Mode** → Continue to Phase 2 (Team Mode)

---

### Phase 2 (Standard): Launch Review Activities

Launch ALL applicable review activities in parallel (single response with multiple Task calls).

**For each perspective, describe the review intent:**

```
Review this code for [PERSPECTIVE]:

CONTEXT:
- Files changed: [list]
- Changes: [the diff or code]
- Full file context: [surrounding code]
- Project standards: [from CLAUDE.md, .editorconfig, etc.]

FOCUS: [What this perspective looks for - from table above]

OUTPUT: Return findings as a structured list, one per finding:

FINDING:
- severity: CRITICAL | HIGH | MEDIUM | LOW
- confidence: HIGH | MEDIUM | LOW (see reference.md for classification matrix)
- title: Brief title (max 40 chars, e.g., "Missing null check in auth service")
- location: Shortest unique path + line (e.g., "auth/service.ts:42-45")
- issue: One sentence describing what's wrong (e.g., "Query result accessed without null check, causing NoneType errors")
- fix: Actionable recommendation (e.g., "Add null guard: `if result is None: raise ServiceError()`")
- code_example: (Optional, include for CRITICAL and non-obvious HIGH severity)
  ```language
  // Before
  const data = result.data;

  // After
  if (!result) throw new Error('No result');
  const data = result.data;
  ```

Confidence Guidelines:
- HIGH: Clear violation of established pattern or security rule
- MEDIUM: Likely issue but context-dependent
- LOW: Potential improvement, may not be applicable

If no findings for this perspective, return: NO_FINDINGS
```

Continue to **Phase 3: Synthesize & Present**.

---

### Phase 2 (Team Mode): Launch Review Team

#### Step 1: Create Team

Derive a team name from the review target:
- PR number → `review-pr-123`
- Branch → `review-feature-auth`
- Staged → `review-staged`
- File path → `review-src-auth` (first meaningful path segment)

```
TeamCreate({
  team_name: "{review-target-name}",
  description: "Code review team for {target}"
})
```

#### Step 2: Create Tasks

Create one task per applicable review perspective. All tasks are independent — no `addBlockedBy` needed.

```
TaskCreate({
  subject: "{Perspective} review of {target}",
  description: """
    Review the following changes for {perspective focus}:
    - {checklist item 1}
    - {checklist item 2}
    - {checklist item 3}

    Files changed: {file list}
    Diff context: {diff summary or pointer to files}

    Return findings in structured format:
    FINDING:
    - severity: CRITICAL | HIGH | MEDIUM | LOW
    - confidence: HIGH | MEDIUM | LOW
    - title: Brief title (max 40 chars)
    - location: file:line
    - issue: One sentence
    - fix: Actionable recommendation
    - code_example: (Optional, for CRITICAL/HIGH)

    If no findings: return NO_FINDINGS
  """,
  activeForm: "Reviewing {perspective}",
  metadata: {
    "perspective": "{perspective-key}",
    "emoji": "{perspective-emoji}"
  }
})
```

#### Step 3: Spawn Reviewer Teammates

Spawn one teammate per applicable perspective. Use specialized agent types from the team plugin where available.

| Teammate Name | Perspective | subagent_type |
|---------------|------------|---------------|
| `security-reviewer` | 🔐 Security | `team:the-architect:review-security` |
| `simplification-reviewer` | 🔧 Simplification | `team:the-architect:review-complexity` |
| `performance-reviewer` | ⚡ Performance | `team:the-developer:optimize-performance` |
| `quality-reviewer` | 📝 Quality | `general-purpose` |
| `test-reviewer` | 🧪 Testing | `team:the-tester:test-quality` |
| `concurrency-reviewer` | 🧵 Concurrency | `team:the-developer:review-concurrency` |
| `dependency-reviewer` | 📦 Dependencies | `team:the-devops:review-dependency` |
| `compatibility-reviewer` | 🔄 Compatibility | `team:the-architect:review-compatibility` |
| `accessibility-reviewer` | ♿ Accessibility | `team:the-designer:build-accessibility` |

**Spawn template for each reviewer:**

```
Task({
  description: "{Perspective} code review",
  prompt: """
  You are the {name} on the {team-name} team.

  CONTEXT:
    - Files changed: {file list}
    - Changes: {the diff or code}
    - Full file context: {surrounding code for each changed file}
    - Project standards: {from CLAUDE.md, .editorconfig, etc.}

  OUTPUT: Return findings in structured format:
    FINDING:
    - severity: CRITICAL | HIGH | MEDIUM | LOW
    - confidence: HIGH | MEDIUM | LOW
    - title: Brief title (max 40 chars)
    - location: file:line
    - issue: One sentence describing what's wrong
    - fix: Actionable recommendation
    - code_example: (Optional, for CRITICAL/HIGH)

    If no findings: return NO_FINDINGS

  SUCCESS: All {perspective} concerns identified with remediation steps

  TEAM PROTOCOL:
    - Check TaskList for your assigned review task
    - Mark in_progress when starting, completed when done
    - Send findings to lead via SendMessage
    - Discover teammates via ~/.claude/teams/{team-name}/config.json
    - If you find an issue overlapping another reviewer's domain,
      DM them: "FYI: Found {issue} at {location} — relates to your review"
    - Do NOT wait for peer responses — send findings to lead regardless
  """,
  subagent_type: "{subagent_type from table}",
  team_name: "{team-name}",
  name: "{teammate-name}",
  mode: "bypassPermissions"
})
```

Launch ALL reviewer teammates simultaneously in a single response with multiple Task calls.

#### Step 4: Monitor & Collect

Messages from reviewers arrive automatically — the lead does NOT poll.

```
Collection loop:
1. Receive message from reviewer: "Review complete. Findings: ..."
2. Receive message from another reviewer: "Review complete. Findings: ..."
3. When all reviewers have reported:
   → Check TaskList to verify all review tasks are completed
   → Proceed to synthesis
```

If a reviewer is blocked or reports an issue:
- Missing context → Send DM with the needed information
- Tool failure → Reassign task or handle directly
- After 3 retries for the same task → Skip that perspective and note it in the summary

#### Step 5: Graceful Shutdown

After collecting all findings:

```
1. Verify all tasks completed via TaskList
2. For EACH reviewer teammate (sequentially):
   SendMessage({
     type: "shutdown_request",
     recipient: "{reviewer-name}",
     content: "Review complete. Thank you for your analysis."
   })
3. Wait for each shutdown_response (approve: true)
4. After ALL teammates shut down:
   TeamDelete()
```

If a teammate rejects shutdown: check TaskList for incomplete work, resolve, then re-request.

Continue to **Phase 3: Synthesize & Present**.

---

### Phase 3: Synthesize & Present

This phase is the same for both Standard and Team Mode.

**For Team Mode**, apply the deduplication algorithm before building the summary:

```
Deduplication algorithm:
1. Collect all findings from all reviewers
2. Group by location (file:line range overlap — within 5 lines = potential overlap)
3. For overlapping findings:
   a. Keep the highest severity version
   b. Merge complementary details (e.g., security + quality insights)
   c. Credit both perspectives in the finding
4. Sort by severity (Critical > High > Medium > Low) then confidence
5. Assign finding IDs (C1, C2, H1, H2, M1, etc.)
6. Build summary table
```

1. **Collect** all findings from review activities
2. **Deduplicate** overlapping findings (keep highest severity)
3. **Rank** by severity (Critical > High > Medium > Low) then confidence
4. **Group** by category for readability

Present in this format:

```markdown
## Code Review: [target]

**Verdict**: 🔴 REQUEST CHANGES | 🟡 APPROVE WITH COMMENTS | ✅ APPROVE

### Summary

| Category | Critical | High | Medium | Low |
|----------|----------|------|--------|-----|
| 🔐 Security | X | X | X | X |
| 🔧 Simplification | X | X | X | X |
| ⚡ Performance | X | X | X | X |
| 📝 Quality | X | X | X | X |
| 🧪 Testing | X | X | X | X |
| **Total** | X | X | X | X |

*🔴 Critical & High Findings (Must Address)*

| ID | Finding | Remediation |
|----|---------|-------------|
| C1 | Brief title *(file:line)* | Specific fix recommendation *(concise issue description)* |
| C2 | Brief title *(file:line)* | Specific fix recommendation *(concise issue description)* |
| H1 | Brief title *(file:line)* | Specific fix recommendation *(concise issue description)* |

#### Code Examples for Critical Fixes

**[C1] Title**
```language
// Before
old code

// After
new code
```

**[C2] Title**
```language
// Before
old code

// After
new code
```

*🟡 Medium Findings (Should Address)*

| ID | Finding | Remediation |
|----|---------|-------------|
| M1 | Brief title *(file:line)* | Specific fix recommendation *(concise issue description)* |
| M2 | Brief title *(file:line)* | Specific fix recommendation *(concise issue description)* |

*⚪ Low Findings (Consider)*

| ID | Finding | Remediation |
|----|---------|-------------|
| L1 | Brief title *(file:line)* | Specific fix recommendation *(concise issue description)* |

### Strengths

- ✅ [Positive observation with specific code reference]
- ✅ [Good patterns noticed]

### Verdict Reasoning

[Why this verdict was chosen based on findings]
```

**Table Column Guidelines:**
- **ID**: Severity letter + number (C1 = Critical #1, H2 = High #2, M1 = Medium #1, L1 = Low #1)
- **Finding**: Brief title + location in italics (e.g., `Missing null check *(auth/service.ts:42)*`)
- **Remediation**: Fix recommendation + issue context in italics (e.g., `Add null guard *(query result accessed without check)*`)

**Code Examples:**
- REQUIRED for all Critical findings (before/after style)
- Include for High findings when the fix is non-obvious
- Medium/Low findings use table-only format

### Phase 4: Next Steps

Use `AskUserQuestion` with options based on verdict:

**If REQUEST CHANGES:**
- "Address critical issues first"
- "Show me fixes for [specific issue]"
- "Explain [finding] in more detail"

**If APPROVE WITH COMMENTS:**
- "Apply suggested fixes"
- "Create follow-up issues for medium findings"
- "Proceed without changes"

**If APPROVE:**
- "Add to PR comments (if PR review)"
- "Done"

## Verdict Decision Matrix

| Critical | High | Decision |
|----------|------|----------|
| > 0 | Any | 🔴 REQUEST CHANGES |
| 0 | > 3 | 🔴 REQUEST CHANGES |
| 0 | 1-3 | 🟡 APPROVE WITH COMMENTS |
| 0 | 0 (Medium > 0) | 🟡 APPROVE WITH COMMENTS |
| 0 | 0 (Low only) | ✅ APPROVE |

## Important Notes

- **Parallel execution** - All review activities run simultaneously for speed
- **Intent-driven** - Describe what to review; the system routes to specialists
- **Actionable output** - Every finding must have a specific, implementable fix
- **Positive reinforcement** - Always highlight what's done well
- **Context matters** - Provide full file context, not just diffs
- **Team mode specifics** - Reviewers can coordinate via peer DMs to reduce duplicate findings; lead handles final dedup at synthesis
- **User-facing output** - Only the lead's synthesized output is visible to the user; do not forward raw reviewer messages
