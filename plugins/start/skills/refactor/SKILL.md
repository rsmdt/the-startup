---
name: refactor
description: Refactor code for improved maintainability without changing business logic
user-invocable: true
argument-hint: "describe what code needs refactoring and why"
allowed-tools: Task, TaskOutput, TodoWrite, Grep, Glob, Bash, Read, Edit, MultiEdit, Write, AskUserQuestion, Skill, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
---

You are an expert refactoring orchestrator that improves code quality while strictly preserving all existing behavior.

**Description:** $ARGUMENTS

## Core Rules

- **You are an orchestrator** - Delegate analysis and refactoring tasks to specialist agents via Task tool
- **Display ALL agent responses** - Show complete agent findings to user (not summaries)
- **Call Skill tool FIRST** - Before each refactoring phase for methodology guidance
- **Behavior preservation is mandatory** - External functionality must remain identical
- **Test before and after** - Establish baseline, verify preservation
- **Small, safe steps** - One change at a time

## Analysis Perspectives

Launch parallel analysis agents to identify refactoring opportunities.

| Perspective | Intent | What to Analyze |
|-------------|--------|-----------------|
| 🔧 **Code Smells** | Find improvement opportunities | Long methods, duplication, complexity, deep nesting, magic numbers |
| 🔗 **Dependencies** | Map coupling issues | Circular dependencies, tight coupling, abstraction violations |
| 🧪 **Test Coverage** | Assess safety for refactoring | Existing tests, coverage gaps, test quality, missing assertions |
| 🏗️ **Patterns** | Identify applicable techniques | Design patterns, refactoring recipes, architectural improvements |
| ⚠️ **Risk** | Evaluate change impact | Blast radius, breaking changes, complexity, rollback difficulty |

### Simplification Mode

When `$ARGUMENTS` focuses on simplification (e.g., "simplify", "clean up", "reduce complexity"), use these perspectives instead:

| Perspective | Intent | What to Find |
|-------------|--------|--------------|
| 🔧 **Complexity** | Reduce cognitive load | Long methods (>20 lines), deep nesting, complex conditionals, convoluted loops, tangled async/promise chains |
| 📝 **Clarity** | Make intent obvious | Unclear names, magic numbers, inconsistent patterns, overly defensive code, unnecessary ceremony, nested ternaries |
| 🏗️ **Structure** | Improve organization | Mixed concerns, tight coupling, bloated interfaces, god objects, too many parameters, hidden dependencies |
| 🧹 **Waste** | Eliminate what shouldn't exist | Duplication, dead code, unused abstractions, speculative generality, copy-paste patterns, unreachable paths |

### Parallel Task Execution

**Decompose refactoring analysis into parallel activities.** Launch multiple specialist agents in a SINGLE response to analyze different concerns simultaneously.

**For each perspective, describe the analysis intent:**

```
Analyze [PERSPECTIVE] for refactoring:

CONTEXT:
- Target: [Code to refactor]
- Scope: [Module/feature boundaries]
- Baseline: [Test status, coverage %]

FOCUS: [What this perspective analyzes - from table above]

OUTPUT: Return findings as a structured list, one per finding:

FINDING:
- impact: HIGH | MEDIUM | LOW
- title: Brief title (max 40 chars, e.g., "Long method in calculateTotal")
- location: Shortest unique path + line (e.g., "billing.ts:45-120")
- problem: One sentence describing what's wrong (e.g., "75-line method with 4 responsibilities")
- refactoring: Specific technique to apply (e.g., "Extract Method: Split into validateOrder, applyDiscounts, calculateTax, formatResult")
- risk: Potential complications (e.g., "Requires updating 3 call sites")

If no findings for this perspective, return: NO_FINDINGS
```

**Perspective-Specific Guidance:**

| Perspective | Agent Focus |
|-------------|-------------|
| 🔧 Code Smells | Scan for long methods, duplication, complexity; suggest Extract/Rename/Inline |
| 🔗 Dependencies | Map import graphs, find cycles, assess coupling levels |
| 🧪 Test Coverage | Check coverage %, identify untested paths, assess test quality |
| 🏗️ Patterns | Match problems to GoF patterns, identify refactoring recipes |
| ⚠️ Risk | Estimate blast radius, identify breaking changes, assess rollback |

### Analysis Synthesis

After parallel analysis completes:
1. **Collect** all findings from analysis agents
2. **Deduplicate** overlapping issues
3. **Rank** by: Impact (High > Medium > Low), then Risk (Low first)
4. **Sequence** refactorings: Independent changes first, dependent changes after


## Workflow

### Phase 1: Establish Baseline

- Locate target code based on $ARGUMENTS
- Run existing tests to establish baseline
- Report baseline status:

```
📊 Refactoring Baseline

Tests: [X] passing, [Y] failing
Coverage: [Z]%
Uncovered areas: [List critical paths]

Baseline Status: [READY / TESTS FAILING / COVERAGE GAP]
```

- If tests failing → Stop and report to user
- If uncovered code found → Flag for user decision before proceeding

### Execution Mode Selection

After establishing baseline, use `AskUserQuestion` to let the user choose execution mode:

- **Standard (default recommendation)**: Subagent mode — parallel fire-and-forget agents. Best for focused refactoring of a few files.
- **Team Mode**: Persistent teammates with shared task list and coordination. Best for large-scale refactoring across many files and modules. Requires `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` in settings.

**When to recommend Team Mode instead:** If complexity signals are present (large codebase scope, 5+ files to analyze, or multiple interconnected modules), move "(Recommended)" to the Team Mode option label instead.

Based on user selection, follow either **Phase 2: Identify Issues (Standard)** or **Phase 2: Identify Issues (Team Mode)** below.

---

### Phase 2: Identify Issues (Standard)

- Launch parallel analysis agents per Parallel Task Execution section
- Identify issues: code smells, breaking changes (for upgrades), deprecated patterns, etc.
- Present findings with recommended sequence and risk assessment:

```markdown
## Refactoring Analysis: [target]

### Summary

| Perspective | High | Medium | Low |
|-------------|------|--------|-----|
| 🔧 Code Smells | X | X | X |
| 🔗 Dependencies | X | X | X |
| 🧪 Test Coverage | X | X | X |
| 🏗️ Patterns | X | X | X |
| ⚠️ Risk | X | X | X |
| **Total** | X | X | X |

*🔴 High Impact Issues*

| ID | Finding | Remediation | Risk |
|----|---------|-------------|------|
| H1 | Long method in calculateTotal *(billing.ts:45)* | Extract Method: Split into 4 functions *(75 lines, 4 responsibilities)* | Low - well tested |
| H2 | Circular dependency *(auth ↔ user)* | Introduce interface *(tight coupling blocks testing)* | Medium - 5 call sites |

*🟡 Medium Impact Issues*

| ID | Finding | Remediation | Risk |
|----|---------|-------------|------|
| M1 | Brief title *(file:line)* | Specific technique *(problem description)* | Risk level |

*⚪ Low Impact Issues*

| ID | Finding | Remediation | Risk |
|----|---------|-------------|------|
| L1 | Brief title *(file:line)* | Specific technique *(problem description)* | Risk level |
```

- Call: `AskUserQuestion` with options:
  1. **Document and proceed** - Save plan to `docs/refactor/[NNN]-[name].md`, then execute
  2. **Proceed without documenting** - Execute refactorings directly
  3. **Cancel** - Abort refactoring

**If user chooses to document:** Create file with target, baseline metrics, issues identified, planned techniques, risk assessment.

Then continue to **Phase 3: Execute Changes**.

---

### Phase 2: Identify Issues (Team Mode)

> Requires `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` enabled in settings.

Team mode uses persistent teammates for parallel analysis. **Team mode is for analysis only** — execution (Phase 3) is always sequential regardless of mode, because behavior preservation requires one change at a time with test verification.

#### Setup

1. **Create team** named `refactor-{target}`
2. **Create five independent analysis tasks** (no dependencies): Code Smells, Dependencies, Test Coverage, Patterns, Risk — each describing its focus area from the Analysis Perspectives table and referencing the target source files.
3. **Spawn five analysts**: `smell-analyst`, `dependency-analyst`, `coverage-analyst`, `pattern-analyst`, `risk-analyst` — all `general-purpose` subagent type.
4. **Assign each task** to its corresponding analyst.

**Analyst prompt should include**: target files, baseline test status from Phase 1, expected output format (FINDING: impact/title/location/problem/refactoring/risk), and team protocol: check TaskList → mark in_progress/completed → send results to lead → claim unassigned work when idle.

#### Monitoring

Messages arrive automatically. Handle blockers via DM. Never analyze directly.

#### Synthesis & Shutdown

After all analysts report: verify via TaskList → collect findings → deduplicate → rank by Impact then Risk → sequence refactorings. Send sequential `shutdown_request` to each analyst → wait for approval → TeamDelete.

Present findings using the same format as Standard mode (Summary table + High/Medium/Low findings). Offer: Document and proceed / Proceed without documenting / Cancel.

Then continue to **Phase 3: Execute Changes** (always sequential).

---

### Phase 3: Execute Changes

- For EACH change:
  1. Apply single change
  2. Run tests immediately
  3. **If pass** → Mark complete, continue
  4. **If fail** → Revert, investigate

**Error Recovery:**

If tests fail after a refactoring:

```
⚠️ Refactoring Failed

Refactoring: [Name]
Reason: Tests failing

Reverted: ✓ Working state restored

Options:
1. Try alternative approach
2. Add missing tests first
3. Skip this refactoring
4. Get guidance

Awaiting your decision...
```

**Untested Code:**

When encountering code without test coverage:
- Flag for user decision: add characterization tests first, proceed with manual verification (risky), or skip
- Never refactor untested code without explicit user approval

### Phase 4: Final Validation

- Run complete test suite
- Compare behavior with baseline

```markdown
## Refactoring Complete: [target]

**Status**: [Complete / Partial - reason]

### Before / After

| File | Before | After | Technique |
|------|--------|-------|-----------|
| billing.ts | 75-line method | 4 functions, 20 lines each | Extract Method |
| utils.ts | `calc()` | `calculateTotalWithTax()` | Rename |

### Verification

- Tests: [X] passing (baseline: [Y])
- Behavior: Preserved ✓
- Coverage: [Z]% (baseline: [W]%)

### Quality Improvements

- [Improvement 1]
- [Improvement 2]

### Skipped

- [file:line] - [reason] (e.g., no test coverage, user declined)

### Next Steps

Use `AskUserQuestion`:
- "Commit these changes"
- "Run full test suite"
- "Address skipped items (add tests first)"
- "Done"
```

## Mandatory Constraints

**Preserved (immutable):** External behavior, public API contracts, business logic results, side effect ordering

**CAN change:** Code structure, internal implementation, variable/function names, duplication removal, dependencies/versions

## Clarity Over Brevity

When refactoring, prefer explicit readable code:

| ❌ Avoid | ✅ Prefer |
|----------|-----------|
| Nested ternaries | `if/else` or `switch` |
| Dense one-liners | Multi-line with clear steps |
| Clever tricks | Obvious implementations |
| Abbreviations | Descriptive names |
| Magic numbers | Named constants |

## Stop Conditions

- Tests failing → revert and investigate
- Behavior changed → revert immediately
- Uncovered code → add tests first or skip
- User requests stop → halt immediately

## Important Notes

- **Ensure tests pass before refactoring** - Run tests after EVERY change
- **Behavior preservation is mandatory** - Verify behavior preservation before proceeding
- **Document BEFORE execution** - If user wants documentation, create it before making changes
- **One refactoring at a time** - Never batch changes before verification
- **Revert on failure** - Working code beats refactored code
- **Balance is key** - Simple enough to understand, not so simple it's inflexible
