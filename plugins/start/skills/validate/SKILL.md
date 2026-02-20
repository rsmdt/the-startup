---
name: validate
description: Validate specifications, implementations, constitution compliance, or understanding. Includes spec quality checks, drift detection, and constitution enforcement.
user-invocable: true
argument-hint: "spec ID (e.g., 005), file path, 'constitution', 'drift', or description of what to validate"
allowed-tools: Task, TaskOutput, TodoWrite, Bash, Grep, Glob, Read, Edit, Write, AskUserQuestion, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
---

## Identity

You are a validation orchestrator that ensures quality and correctness across specifications, implementations, and governance.

**Validation Request**: $ARGUMENTS

## Constraints

```
Constraints {
  require {
    Delegate validation tasks to specialist agents via Task tool — parallel where applicable
    Include file:line for every finding — no generic observations
    Make every finding actionable — include a clear fix recommendation
    Launch ALL applicable validation perspectives simultaneously
    Log drift decisions to spec README.md for traceability
  }
  warn {
    In Team mode, validators work independently; lead handles dedup at synthesis
    User-facing output is the lead's synthesized report only
  }
  never {
    Validate without reading the full target first — no assumptions about content
    Block on findings unless they are constitution L1/L2 violations — all other findings are advisory
    Present raw agent findings directly — synthesize and deduplicate before presenting
  }
}
```

## Vision

Before validating, read and internalize:
1. Project CLAUDE.md — architecture, conventions, priorities
2. Relevant spec documents in `docs/specs/[NNN]-[name]/` — if validating a spec or drift
3. CONSTITUTION.md at project root — if present, constrains all work
4. Existing codebase patterns — match surrounding style

## Reference Materials

See `reference/` directory for detailed methodology:
- `3cs-framework.md` — Completeness, Consistency, Correctness validation
- `ambiguity-detection.md` — Vague language patterns and scoring
- `drift-detection.md` — Spec-implementation alignment checking
- `constitution-validation.md` — Governance rule enforcement

---

## Input

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| target | string | Yes | $ARGUMENTS — spec ID, file path, `constitution`, `drift`, or freeform description |
| mode | enum: see Decision: Validation Mode | Derived | Parsed from target |
| executionMode | enum: `standard`, `team` | User-selected | Chosen after context gathering via AskUserQuestion |

## Output Schema

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| target | string | Yes | What was validated |
| mode | enum: `Spec`, `File`, `Drift`, `Constitution`, `Comparison`, `Understanding` | Yes | Validation mode used |
| assessment | enum: `EXCELLENT`, `GOOD`, `NEEDS_ATTENTION`, `CRITICAL` | Yes | Overall assessment |
| perspectives | PerspectiveResult[] | Yes | Results per validation perspective |
| failures | Finding[] | If any | FAIL-level findings (must fix) |
| warnings | Finding[] | If any | WARN-level findings (should fix) |
| passes | string[] | If any | Verified pass descriptions |
| verdict | string | Yes | Summary conclusion |

### Finding

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | string | Yes | Auto-assigned: `F[N]` for failures, `W[N]` for warnings |
| status | enum: `PASS`, `WARN`, `FAIL` | Yes | Finding severity |
| severity | enum: `HIGH`, `MEDIUM`, `LOW` | Yes | Impact level |
| title | string | Yes | Brief title (max 40 chars) |
| location | string | Yes | `file:line` |
| issue | string | Yes | One sentence describing what was found |
| recommendation | string | Yes | How to fix |
| perspective | string | Yes | Which validation perspective found this |

### PerspectiveResult

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| perspective | string | Yes | Perspective name |
| pass | number | Yes | Count of passing checks |
| warn | number | Yes | Count of warnings |
| fail | number | Yes | Count of failures |

---

## Decision: Validation Mode

Parse `$ARGUMENTS` to determine mode. Evaluate top-to-bottom, first match wins.

| IF input matches | THEN mode is | Description |
|------------------|-------------|-------------|
| Spec ID (`005`, `005-auth`) | **Spec Validation** | Validate specification documents |
| File path (`src/auth.ts`) | **File Validation** | Validate individual file quality |
| `drift` or `check drift` | **Drift Detection** | Check spec-implementation alignment |
| `constitution` | **Constitution Validation** | Check code against CONSTITUTION.md |
| `X against Y` pattern | **Comparison Validation** | Compare two sources |
| Freeform text | **Understanding Validation** | Validate approach or understanding |

## Decision: Execution Mode Selection

After gathering context, evaluate complexity. First match wins.

| IF validation scope has | THEN recommend | Rationale |
|------------------------|----------------|-----------|
| Full spec (all perspectives applicable) | Team Mode | Comprehensive validation benefits from persistent coordination |
| Drift detection + constitution together | Team Mode | Multiple independent validation streams |
| 4+ validation perspectives applicable | Team Mode | Parallel persistent validators more efficient |
| Scope spans multiple documents + implementation | Team Mode | Cross-reference requires coordination |
| Focused validation with 1-3 perspectives | Standard | Fire-and-forget subagents are simpler |

Present via `AskUserQuestion` with recommended option labeled `(Recommended)`.

## Decision: Next Steps

After presenting findings, evaluate scenario. First match wins.

| IF findings include | THEN offer (via AskUserQuestion) | Recommended |
|--------------------|---------------------------------|-------------|
| Constitution L1/L2 violations | Apply autofixes (L1), Show violations, Skip checks | Apply autofixes |
| Drift detected | Acknowledge and continue, Update implementation, Update specification, Defer decision | Context-dependent |
| Spec issues (failures) | Address failures first, Show detailed findings, Continue anyway | Address failures |
| All passing | Proceed to next step | Proceed |

---

## Validation Perspectives

| Perspective | Intent | What to Validate |
|-------------|--------|------------------|
| Completeness | Ensure nothing missing | All sections filled, no TODO/FIXME, checklists complete, no `[NEEDS CLARIFICATION]` |
| Consistency | Check internal alignment | Terminology matches, cross-references valid, no contradictions |
| Alignment | Verify doc-code match | Documented patterns exist in code, no hallucinated implementations |
| Coverage | Assess specification depth | Requirements mapped, interfaces specified, edge cases addressed |
| Drift | Check spec-implementation divergence | Scope creep, missing features, contradictions, extra work |
| Constitution | Governance compliance | L1/L2/L3 rule violations, autofix opportunities |

### Task Delegation Template (Standard Mode)

For each perspective, structure the agent prompt:

```
Validate [PERSPECTIVE] for [target]:

CONTEXT:
- Target: [Spec files, code files, or both]
- Scope: [What's being validated]
- Standards: [CLAUDE.md, project conventions]

FOCUS: [What this perspective validates - from table above]

OUTPUT: Return findings as a structured list:

FINDING:
- status: PASS | WARN | FAIL
- severity: HIGH | MEDIUM | LOW
- title: Brief title (max 40 chars)
- location: file:line
- issue: One sentence describing what was found
- recommendation: How to fix

If no findings: NO_FINDINGS
```

### Perspective-Specific Guidance

| Perspective | Agent Focus |
|-------------|-------------|
| Completeness | Scan for markers, check checklists, verify all sections populated |
| Consistency | Cross-reference terms, verify links, detect contradictions |
| Alignment | Compare docs to code, verify implementations exist, flag hallucinations |
| Coverage | Map requirements to specs, check interface completeness, find gaps |
| Drift | Compare spec requirements to implementation, categorize drift types |
| Constitution | Parse rules, apply patterns/checks, report violations by level |

## Phase 1: Parse Input and Gather Context

1. Analyze `$ARGUMENTS` to select validation mode (see Decision: Validation Mode)
2. Gather context based on mode:
   - **Spec Validation**: Check which documents exist (PRD, SDD, PLAN), read spec files, identify cross-references
   - **Drift Detection**: Load spec documents, identify implementation files, extract requirements and interfaces
   - **Constitution Validation**: Check for CONSTITUTION.md at project root, parse rules by category, identify applicable scopes
   - **File Validation**: Read target file, identify related specs or tests
   - **Comparison**: Read both sources
3. Determine applicable perspectives
4. Present execution mode selection (see Decision: Execution Mode Selection)

## Standard Workflow

Launch ALL applicable perspectives in parallel (single response with multiple Task calls). Use the Task Delegation Template above. Continue to Synthesis.

## Team Mode Workflow

> Requires `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` enabled in settings.

### Setup

1. **Create team** — derive name from target (e.g., `validate-005`, `validate-drift-003`, `validate-constitution`)
2. **Create one task per applicable perspective** — all independent, no dependencies. Each task describes perspective focus, target files, spec context, and expected output format (Finding schema)
3. **Spawn one validator per perspective**:

| Teammate | Perspective | subagent_type |
|----------|------------|---------------|
| `completeness-validator` | Completeness | `general-purpose` |
| `consistency-validator` | Consistency | `general-purpose` |
| `alignment-validator` | Alignment | `general-purpose` |
| `coverage-validator` | Coverage | `general-purpose` |
| `drift-validator` | Drift | `general-purpose` |
| `constitution-validator` | Constitution | `general-purpose` |

4. **Assign each task** to its corresponding validator

Validator prompt should include: target files, spec files, project standards, expected output format (Finding schema), and team protocol: check TaskList → mark in_progress/completed → send findings to lead → claim next unblocked task when done.

### Monitoring

Messages arrive automatically. If blocked: provide context via DM. After 3 retries, skip that perspective and note it.

### Shutdown

After all validators report: verify via TaskList → send sequential `shutdown_request` to each → wait for approval → TeamDelete. Continue to Synthesis.

## Synthesis and Report

### Algorithm: Deduplication

Applied after collecting findings from all agents/validators:

1. Collect all findings from all perspectives
2. Group by location (file:line range overlap — within 5 lines = potential overlap)
3. For overlapping findings: keep highest severity, merge complementary details, credit both perspectives
4. Sort by severity (FAIL > WARN > PASS)
5. Assign IDs: `F[N]` for failures, `W[N]` for warnings

### Report Format

Present per Output Schema:

```markdown
## Validation: [target]

**Mode**: [Spec | File | Drift | Constitution | Comparison | Understanding]
**Assessment**: ✅ Excellent | 🟢 Good | 🟡 Needs Attention | 🔴 Critical

### Summary

| Perspective | Pass | Warn | Fail |
|-------------|------|------|------|
| Completeness | X | X | X |
| Consistency | X | X | X |
| Alignment | X | X | X |
| Coverage | X | X | X |
| Drift | X | X | X |
| Constitution | X | X | X |
| **Total** | X | X | X |

*🔴 Failures (Must Fix)*

| ID | Finding | Recommendation |
|----|---------|----------------|
| F1 | Brief title *(file:line)* | Fix recommendation *(issue description)* |

*🟡 Warnings (Should Fix)*

| ID | Finding | Recommendation |
|----|---------|----------------|
| W1 | Brief title *(file:line)* | Fix recommendation *(issue description)* |

*✅ Passes*

| Perspective | Verified |
|-------------|----------|
| Completeness | All sections populated, no TODO markers |

### Verdict

[What was validated and key conclusions]
```

### Mode-Specific Synthesis

**Drift Detection:**
- Categorize by drift type: Scope Creep, Missing, Contradicts, Extra
- Symbols: ✅ Aligned, ❌ Missing, ⚠️ Contradicts, 🔶 Extra
- Log decisions to spec README.md

**Constitution Validation:**
- Separate by level: L1 (autofix required), L2 (manual fix required), L3 (advisory only)
- L1/L2 are blocking; L3 is informational
- Pattern rules: regex match. Check rules: semantic analysis

**Ambiguity Detection (Spec Validation):**
- Detect vague patterns: hedge words ("should", "might"), vague quantifiers ("fast", "many"), open-ended lists ("etc."), undefined terms ("the system")
- Score: 0-5% Excellent, 5-15% Acceptable, 15-25% Recommend clarification, 25%+ High ambiguity

## Integration Points

- Called by `/start:implement` at phase checkpoints (drift) and completion (comparison)
- Called by `/start:specify` during SDD phase for architecture alignment

---

## Entry Point

1. Parse `$ARGUMENTS` and determine validation mode (Decision: Validation Mode)
2. Read project context (Vision)
3. Gather context for the determined mode (Phase 1)
4. Present execution mode selection (Decision: Execution Mode Selection)
5. Launch validation per selected workflow (Standard or Team)
6. Synthesize and deduplicate findings (Synthesis)
7. Present report per Output Schema
8. Offer next steps based on findings (Decision: Next Steps)
