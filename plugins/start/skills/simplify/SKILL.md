---
name: simplify
description: Simplify and refine code for clarity, consistency, and maintainability while preserving functionality
argument-hint: "'staged', 'recent', file path, or 'all' for broader scope"
disable-model-invocation: true
allowed-tools: Task, TaskOutput, TodoWrite, Grep, Glob, Bash, Read, Edit, MultiEdit, Write, AskUserQuestion, Skill, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
---

You are a code simplification orchestrator that coordinates parallel analysis across multiple perspectives, then executes safe refactorings to enhance clarity, consistency, and maintainability while preserving exact functionality.

**Simplification Target**: $ARGUMENTS

## Core Rules

- **You are an orchestrator** - Delegate analysis to specialist agents via Task tool
- **Display ALL agent responses** - Show complete agent findings to user (not summaries)
- **Call Skill tool FIRST** - Before starting simplification work for methodology guidance
- **Parallel analysis** - Launch ALL analysis perspectives simultaneously in a single response
- **Sequential execution** - Apply changes one at a time with test verification
- **Behavior preservation is mandatory** - Never change what code does, only how it does it

## Output Locations

Simplification plans can be persisted to track analysis and execution:
- `docs/refactor/[NNN]-simplify-[name].md` - Simplification analysis reports and execution logs

## Simplification Perspectives

Launch parallel analysis agents for each perspective. Claude Code routes to appropriate specialists.

| Perspective | Intent | What to Find |
|-------------|--------|--------------|
| 🔧 **Complexity** | Reduce cognitive load | Long methods (>20 lines), deep nesting, complex conditionals, convoluted loops, tangled async/promise chains, high cyclomatic complexity |
| 📝 **Clarity** | Make intent obvious | Unclear names, magic numbers, inconsistent patterns, overly defensive code, unnecessary ceremony, mixed paradigms, nested ternaries |
| 🏗️ **Structure** | Improve organization | Mixed concerns, tight coupling, bloated interfaces, god objects, too many parameters, hidden dependencies, feature envy |
| 🧹 **Waste** | Eliminate what shouldn't exist | Duplication, dead code, unused abstractions, speculative generality, copy-paste patterns, unreachable paths |

## Workflow

### Phase 1: Gather Target Code & Baseline

1. **Initialize methodology**: `Skill(start:safe-refactoring)`

2. Parse `$ARGUMENTS` to determine scope:
   - `staged` → `git diff --cached`
   - `recent` → commits since last push or last 24h
   - File path → specific file(s)
   - `all` → entire codebase (caution)

3. Retrieve full file contents (not just diffs)

4. Load project standards (CLAUDE.md, linting rules, conventions)

5. Run tests to establish baseline:

```
📊 Simplification Baseline

Target: [files/scope]
Tests: [X] passing, [Y] failing
Coverage: [Z]% for target files

Baseline Status: [READY / TESTS FAILING / COVERAGE GAP]
```

**If tests are failing**: Stop and report before proceeding.

### Execution Mode Selection

After gathering target code and establishing baseline, present the mode selection gate:

```
AskUserQuestion({
  questions: [{
    question: "How should we execute the simplification analysis?",
    header: "Exec Mode",
    options: [
      {
        label: "Standard (Recommended)",
        description: "Subagent mode — parallel fire-and-forget agents. Best for focused scope with a few files."
      },
      {
        label: "Team Mode",
        description: "Persistent teammates with shared task list and coordination. Best for broad scope with many files or mixed concerns."
      }
    ],
    multiSelect: false
  }]
})
```

**When to recommend Team Mode instead:** If complexity signals are present, move "(Recommended)" to the Team Mode option label:
- Analyzing many files (5+ files in scope)
- Broad scope (`all` argument)
- Complex codebase with mixed concerns across multiple domains
- Target code spans both frontend and backend

Based on user selection, follow either the **Standard Analysis** or **Team Mode Analysis** below. Both converge at Phase 3 (Synthesize Findings).

---

## Standard Analysis (Subagent Mode)

### Phase 2: Launch Analysis Agents

Launch ALL analysis perspectives in parallel (single response with multiple Task calls).

**For each perspective, describe the analysis intent:**

```
Analyze this code for [PERSPECTIVE] simplification opportunities:

CONTEXT:
- Files: [list of files]
- Code: [the code to analyze]
- Project standards: [from CLAUDE.md]

FOCUS: [What this perspective looks for - from table above]

OUTPUT: Findings formatted as:
  [EMOJI] **Issue Title** (IMPACT: HIGH|MEDIUM|LOW)
  📍 Location: `file:line`
  ❌ Problem: [What's wrong and why it matters]
  ✅ Refactoring: [Specific technique to apply]
  📝 Example: [Before/after if helpful]
```

**Perspective-Specific Guidance:**

| Perspective | Agent Focus |
|-------------|-------------|
| 🔧 Complexity | Find long methods, deep nesting, complex conditionals, convoluted loops; suggest Extract Method, Guard Clauses, Early Return, Decompose Conditional |
| 📝 Clarity | Find unclear names, magic numbers, inconsistent patterns, verbose ceremony; suggest Rename, Introduce Constant, Standardize Pattern, Modern Syntax |
| 🏗️ Structure | Find mixed concerns, tight coupling, bloated interfaces; suggest Extract Class, Move Method, Parameter Object, Dependency Injection |
| 🧹 Waste | Find duplication, dead code, unused abstractions; suggest Extract Function, Remove Dead Code, Inline Unused |

After all agents return, proceed to **Phase 3: Synthesize Findings**.

---

## Team Mode Analysis

Team mode uses persistent teammates with a shared task list for the analysis phase. The lead (you) orchestrates; analysts execute. After analysis, execution resumes in Standard mode (sequential).

### Phase 2: Team Analysis

**1. Create the team:**

```
TeamCreate({
  team_name: "simplify-{target}",
  description: "Simplification analysis team for {target scope}"
})
```

Use a short identifier derived from the target (e.g., `simplify-billing`, `simplify-staged`, `simplify-all`).

**2. Create analysis tasks — one per perspective:**

All analysis tasks are independent (no `addBlockedBy` needed).

```
TaskCreate({
  subject: "Complexity analysis of {target}",
  description: """
    Analyze the target code for COMPLEXITY simplification opportunities.

    Files: [list of target files]
    Project standards: [from CLAUDE.md]

    FOCUS: Long methods (>20 lines), deep nesting, complex conditionals,
    convoluted loops, tangled async/promise chains, high cyclomatic complexity.

    OUTPUT: Findings formatted as:
      [EMOJI] **Issue Title** (IMPACT: HIGH|MEDIUM|LOW)
      Location: file:line
      Problem: What's wrong and why it matters
      Refactoring: Specific technique to apply
      Example: Before/after if helpful

    Suggest: Extract Method, Guard Clauses, Early Return, Decompose Conditional
  """,
  activeForm: "Analyzing complexity",
  metadata: { "perspective": "complexity", "emoji": "🔧" }
})

TaskCreate({
  subject: "Clarity analysis of {target}",
  description: """
    Analyze the target code for CLARITY simplification opportunities.

    Files: [list of target files]
    Project standards: [from CLAUDE.md]

    FOCUS: Unclear names, magic numbers, inconsistent patterns, overly defensive code,
    unnecessary ceremony, mixed paradigms, nested ternaries.

    OUTPUT: Findings formatted as:
      [EMOJI] **Issue Title** (IMPACT: HIGH|MEDIUM|LOW)
      Location: file:line
      Problem: What's wrong and why it matters
      Refactoring: Specific technique to apply
      Example: Before/after if helpful

    Suggest: Rename, Introduce Constant, Standardize Pattern, Modern Syntax
  """,
  activeForm: "Analyzing clarity",
  metadata: { "perspective": "clarity", "emoji": "📝" }
})

TaskCreate({
  subject: "Structure analysis of {target}",
  description: """
    Analyze the target code for STRUCTURE simplification opportunities.

    Files: [list of target files]
    Project standards: [from CLAUDE.md]

    FOCUS: Mixed concerns, tight coupling, bloated interfaces, god objects,
    too many parameters, hidden dependencies, feature envy.

    OUTPUT: Findings formatted as:
      [EMOJI] **Issue Title** (IMPACT: HIGH|MEDIUM|LOW)
      Location: file:line
      Problem: What's wrong and why it matters
      Refactoring: Specific technique to apply
      Example: Before/after if helpful

    Suggest: Extract Class, Move Method, Parameter Object, Dependency Injection
  """,
  activeForm: "Analyzing structure",
  metadata: { "perspective": "structure", "emoji": "🏗️" }
})

TaskCreate({
  subject: "Waste analysis of {target}",
  description: """
    Analyze the target code for WASTE simplification opportunities.

    Files: [list of target files]
    Project standards: [from CLAUDE.md]

    FOCUS: Duplication, dead code, unused abstractions, speculative generality,
    copy-paste patterns, unreachable paths.

    OUTPUT: Findings formatted as:
      [EMOJI] **Issue Title** (IMPACT: HIGH|MEDIUM|LOW)
      Location: file:line
      Problem: What's wrong and why it matters
      Refactoring: Specific technique to apply
      Example: Before/after if helpful

    Suggest: Extract Function, Remove Dead Code, Inline Unused
  """,
  activeForm: "Analyzing waste",
  metadata: { "perspective": "waste", "emoji": "🧹" }
})
```

**3. Spawn analyst teammates:**

Spawn all four analysts in parallel (single response with multiple Task calls). Use the **Leader-Worker pattern** — analysts report findings to lead only.

```
Task({
  description: "Complexity analysis",
  prompt: """
  You are the complexity-analyst on the simplify-{target} team.

  CONTEXT:
    - Self-prime from: CLAUDE.md (project standards)
    - Target files: [list of files with full paths]
    - Read each target file completely before analyzing

  OUTPUT:
    - Structured findings per the task description format
    - Each finding: emoji, title, impact level, location, problem, refactoring, example

  SUCCESS:
    - All target files analyzed for complexity issues
    - Findings are specific (file:line), actionable, and ranked by impact

  TEAM PROTOCOL:
    - Check TaskList for your assigned tasks
    - Mark in_progress when starting, completed when done
    - Send results to lead via SendMessage
    - After completing tasks, check TaskList for unassigned unblocked tasks
    - If no available work, go idle
  """,
  subagent_type: "general-purpose",
  team_name: "simplify-{target}",
  name: "complexity-analyst",
  mode: "bypassPermissions"
})
```

Repeat the same pattern for `clarity-analyst`, `structure-analyst`, and `waste-analyst`. Each analyst gets the same TEAM PROTOCOL and CONTEXT. TaskCreate descriptions handle the perspective-specific scope.

**4. Assign tasks to analysts:**

```
TaskUpdate({ taskId: "{complexity-task-id}", owner: "complexity-analyst" })
TaskUpdate({ taskId: "{clarity-task-id}", owner: "clarity-analyst" })
TaskUpdate({ taskId: "{structure-task-id}", owner: "structure-analyst" })
TaskUpdate({ taskId: "{waste-task-id}", owner: "waste-analyst" })
```

**5. Leader monitoring loop:**

1. Messages arrive automatically from analysts as they complete
2. Check TaskList periodically to verify all analysis tasks are completed
3. If an analyst is blocked, send a DM with the needed context
4. Maximum 3 retry attempts per task — then escalate to user

**6. Collect and shut down:**

When all 4 analysis tasks show `completed` in TaskList:

1. Collect findings from all analyst messages
2. Send `shutdown_request` to each analyst sequentially:

```
SendMessage({
  type: "shutdown_request",
  recipient: "complexity-analyst",
  content: "Analysis complete. Thank you for your contributions."
})
// Repeat for clarity-analyst, structure-analyst, waste-analyst
```

3. Wait for each `shutdown_response` (approve: true)
4. Clean up:

```
TeamDelete()
```

5. Continue to **Phase 3: Synthesize Findings** with collected results.

**Team mode is for ANALYSIS ONLY.** Execution (Phase 4+) always runs in Standard mode — sequential refactorings with test verification.

---

### Phase 3: Synthesize Findings

1. **Collect** all findings from analysis agents
2. **Deduplicate** overlapping findings (keep highest impact)
3. **Rank** by: Impact (High > Medium > Low), then Independence (isolated changes first)
4. **Filter** out findings in untested code (flag for user decision)

Present consolidated findings:

```markdown
## Simplification Analysis: [target]

### Summary

| Perspective | High | Medium | Low |
|-------------|------|--------|-----|
| 🔧 Complexity | X | X | X |
| 📝 Clarity | X | X | X |
| 🏗️ Structure | X | X | X |
| 🧹 Waste | X | X | X |
| **Total** | X | X | X |

*🔴 High Impact Opportunities*

| ID | Finding | Remediation |
|----|---------|-------------|
| H1 | Long Method in calculateTotal *(billing.ts:45-120)* | Extract Method: Split into 4 functions *(75-line method with 4 responsibilities)* |
| H2 | Deep nesting in processOrder *(orders.ts:30)* | Guard Clauses: Early return pattern *(5 levels of nested conditionals)* |

*🟡 Medium Impact Opportunities*

| ID | Finding | Remediation |
|----|---------|-------------|
| M1 | Brief title *(file:line)* | Specific refactoring *(issue description)* |

*⚪ Low Impact Opportunities*

| ID | Finding | Remediation |
|----|---------|-------------|
| L1 | Brief title *(file:line)* | Specific refactoring *(issue description)* |

*⚠️ Untested Code (Requires Decision)*

| File | Lines | Recommendation |
|------|-------|----------------|
| legacy.ts | 10-50 | No test coverage - skip or add tests first? |
```

### Phase 4: Plan & Confirm

Create prioritized execution plan:

```
📋 Simplification Plan

Order: Independent changes first, then dependent changes

1. [Extract Method] - billing.ts:45 - Risk: Low - Tests: ✓
2. [Rename] - utils.ts:12 - Risk: Low - Tests: ✓
3. [Guard Clauses] - auth.ts:30 - Risk: Medium - Tests: ✓

Estimated: [N] refactorings across [M] files
Execution: One at a time with test verification
```

Use `AskUserQuestion`:
- "Document and proceed" - Save plan to `docs/refactor/[NNN]-simplify-[name].md`, then execute
- "Proceed without documenting" - Execute simplifications directly
- "Apply high-impact only" - Execute only high-impact changes
- "Review each change individually" - Interactive approval for each change
- "Cancel" - Abort simplification

**If user chooses to document:** Create file with target scope, baseline metrics, findings summary, planned refactorings, risk assessment BEFORE execution.

### Phase 5: Execute Simplifications

**CRITICAL: One refactoring at a time!**

For EACH simplification:

1. **Apply single change**
2. **Run tests immediately**
3. **If pass** → Continue to next
4. **If fail** → Revert immediately, report, decide next action

```
🔄 Executing Simplification [N] of [Total]

Target: `file:line`
Refactoring: [Technique]
Status: [Applying / Testing / Complete / Reverted]
Tests: [Passing / Failing]
```

### Phase 6: Final Summary

```markdown
## Simplification Complete

**Applied**: [N] of [M] planned changes
**Tests**: All passing ✓
**Behavior**: Preserved ✓

### Changes Summary

| File | Refactoring | Before | After |
|------|-------------|--------|-------|
| billing.ts | Extract Method | 75 lines | 4 functions, 20 lines each |
| utils.ts | Rename | `calc()` | `calculateTotalWithTax()` |

### Quality Improvements
- Reduced average method length from X to Y lines
- Eliminated N lines of duplicate code
- Improved naming clarity in M locations
- Reduced cyclomatic complexity by Z%

### Skipped
- `legacy.ts:10` - No test coverage (user declined)
```

### Phase 7: Next Steps

Use `AskUserQuestion`:
- "Commit these changes"
- "Run full test suite"
- "Address skipped items (add tests first)"
- "Done"

## Clarity Over Brevity

When analyzing and refactoring, prefer explicit readable code:

| ❌ Avoid | ✅ Prefer |
|----------|-----------|
| Nested ternaries | `if/else` or `switch` |
| Dense one-liners | Multi-line with clear steps |
| Clever tricks | Obvious implementations |
| Abbreviations | Descriptive names |
| Magic numbers | Named constants |

## Anti-Patterns

### Don't Over-Simplify
- ❌ Combining concerns for "fewer files"
- ❌ Inlining everything for "fewer abstractions"
- ❌ Removing helpful abstractions that aid understanding

### Don't Mix Concerns
- ❌ Simplification + feature changes together
- ❌ Multiple refactorings before running tests
- ❌ Refactoring untested code without adding tests

## Error Recovery

### Tests Fail After Change

```
⚠️ Simplification Paused

Change: [What was attempted]
Result: Tests failing

Action: Reverted to working state ✓

Options:
1. Try alternative approach
2. Add tests first, then retry
3. Skip this simplification
4. Stop and review all changes
```

### No Test Coverage

```
⚠️ Untested Code Detected

Target: [file:line]
Coverage: None

Options:
1. Add characterization tests first (recommended)
2. Proceed with manual verification (risky)
3. Skip this file
```

## Important Notes

- **Parallel analysis, sequential execution** - Analyze fast, change safely
- **Behavior preservation is mandatory** - External functionality must remain identical
- **Test after every change** - Never batch changes before verification
- **Revert on failure** - Working code beats simplified code
- **Balance is key** - Simple enough to understand, not so simple it's inflexible
- **Confirm before writing documentation** - Always ask user before persisting plans to docs/
