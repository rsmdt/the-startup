---
name: validate
description: Validate specifications, implementations, constitution compliance, or understanding. Includes spec quality checks, drift detection, and constitution enforcement.
argument-hint: "spec ID (e.g., 005), file path, 'constitution', 'drift', or description of what to validate"
disable-model-invocation: true
allowed-tools: Task, TaskOutput, TodoWrite, Bash, Grep, Glob, Read, Edit, Write, AskUserQuestion, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
---

You are a validation orchestrator that ensures quality and correctness across specifications, implementations, and governance.

**Validation Request**: $ARGUMENTS

## Core Rules

- **You are an orchestrator** - Delegate validation tasks to specialist agents via Task tool
- **Parallel validation** - Launch ALL applicable validation perspectives simultaneously
- **Advisory by default** - Provide recommendations without blocking (except L1/L2 constitution violations)
- **Be specific** - Include file paths and line numbers for all findings

## Reference Materials

See `reference/` directory for detailed methodology:
- `3cs-framework.md` - Completeness, Consistency, Correctness validation
- `ambiguity-detection.md` - Vague language patterns and scoring
- `drift-detection.md` - Spec-implementation alignment checking
- `constitution-validation.md` - Governance rule enforcement

## Validation Modes

Parse `$ARGUMENTS` to determine mode:

| Input Pattern | Mode | Description |
|---------------|------|-------------|
| Spec ID (`005`, `005-auth`) | **Spec Validation** | Validate specification documents |
| File path (`src/auth.ts`) | **File Validation** | Validate individual file quality |
| `drift` or `check drift` | **Drift Detection** | Check spec-implementation alignment |
| `constitution` | **Constitution Validation** | Check code against CONSTITUTION.md |
| Comparison phrase | **Comparison Validation** | Compare source against reference |
| Freeform text | **Understanding Validation** | Validate approach or understanding |

## Validation Perspectives

Launch parallel validation agents for comprehensive coverage.

| Perspective | Intent | What to Validate |
|-------------|--------|------------------|
| ✅ **Completeness** | Ensure nothing missing | All sections filled, no TODO/FIXME, checklists complete, no `[NEEDS CLARIFICATION]` |
| 🔗 **Consistency** | Check internal alignment | Terminology matches, cross-references valid, no contradictions |
| 📍 **Alignment** | Verify doc-code match | Documented patterns exist in code, no hallucinated implementations |
| 📐 **Coverage** | Assess specification depth | Requirements mapped, interfaces specified, edge cases addressed |
| 📊 **Drift** | Check spec-implementation divergence | Scope creep, missing features, contradictions, extra work |
| 📜 **Constitution** | Governance compliance | L1/L2/L3 rule violations, autofix opportunities |

## Workflow

### Phase 1: Parse Input & Determine Mode

Analyze `$ARGUMENTS` to select validation mode:

```
Spec ID (005) → Spec Validation
File path → File Validation
"drift" → Drift Detection
"constitution" → Constitution Validation
"X against Y" → Comparison Validation
Freeform → Understanding Validation
```

### Phase 2: Gather Context

**For Spec Validation:**
- Check which documents exist (PRD, SDD, PLAN)
- Read relevant specification files
- Identify cross-document references

**For Drift Detection:**
- Load specification documents
- Identify implementation files
- Extract requirements and interfaces from spec

**For Constitution Validation:**
- Check for CONSTITUTION.md at project root
- Parse rules by category
- Identify applicable file scopes

### Mode Selection Gate

After gathering context, determine applicable perspectives and offer execution mode:

```
AskUserQuestion({
  questions: [{
    question: "How should we execute this validation?",
    header: "Exec Mode",
    options: [
      {
        label: "Standard (Recommended)",
        description: "Subagent mode — parallel fire-and-forget agents. Best for focused validation with a few perspectives."
      },
      {
        label: "Team Mode",
        description: "Persistent teammates with shared task list and coordination. Best for comprehensive validation across many perspectives."
      }
    ],
    multiSelect: false
  }]
})
```

**Recommend Team Mode when:**
- Validating a full spec (all perspectives applicable)
- Drift detection + constitution validation together
- 4+ validation perspectives are applicable
- Validation scope spans multiple documents and implementation files

**Post-gate routing:**
- User selects **Standard** → Continue to Phase 3 (Standard)
- User selects **Team Mode** → Continue to Phase 3 (Team Mode)

---

### Phase 3 (Standard): Launch Validation Agents

Launch ALL applicable perspectives in parallel (single response with multiple Task calls).

**For each perspective, use this template:**

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

**Perspective-Specific Guidance:**

| Perspective | Agent Focus |
|-------------|-------------|
| ✅ Completeness | Scan for markers, check checklists, verify all sections populated |
| 🔗 Consistency | Cross-reference terms, verify links, detect contradictions |
| 📍 Alignment | Compare docs to code, verify implementations exist, flag hallucinations |
| 📐 Coverage | Map requirements to specs, check interface completeness, find gaps |
| 📊 Drift | Compare spec requirements to implementation, categorize drift types |
| 📜 Constitution | Parse rules, apply patterns/checks, report violations by level |

Continue to **Phase 4: Synthesize & Present**.

---

### Phase 3 (Team Mode): Launch Validation Team

#### Step 1: Create Team

Derive a team name from the validation target:
- Spec ID → `validate-005`
- File path → `validate-src-auth` (first meaningful path segment)
- Drift → `validate-drift-{spec-id}`
- Constitution → `validate-constitution`

```
TeamCreate({
  team_name: "{validate-target-name}",
  description: "Validation team for {target}"
})
```

#### Step 2: Create Tasks

Create one task per applicable validation perspective. All tasks are independent — no `addBlockedBy` needed.

```
TaskCreate({
  subject: "{Perspective} validation of {target}",
  description: """
    Validate the following for {perspective focus}:
    - {checklist item 1}
    - {checklist item 2}
    - {checklist item 3}

    Target files: {file list}
    Spec context: {spec summary or pointer to files}

    Return findings in structured format:
    FINDING:
    - status: PASS | WARN | FAIL
    - severity: HIGH | MEDIUM | LOW
    - title: Brief title (max 40 chars)
    - location: file:line
    - issue: One sentence describing what was found
    - recommendation: How to fix

    If no findings: return NO_FINDINGS
  """,
  activeForm: "Validating {perspective}",
  metadata: {
    "perspective": "{perspective-key}",
    "emoji": "{perspective-emoji}"
  }
})
```

#### Step 3: Spawn Validator Teammates

Spawn one teammate per applicable perspective. Use the Leader-Worker pattern — validators report findings to lead.

| Teammate Name | Perspective | subagent_type |
|---------------|------------|---------------|
| `completeness-validator` | ✅ Completeness | `general-purpose` |
| `consistency-validator` | 🔗 Consistency | `general-purpose` |
| `alignment-validator` | 📍 Alignment | `general-purpose` |
| `coverage-validator` | 📐 Coverage | `general-purpose` |
| `drift-validator` | 📊 Drift | `general-purpose` |
| `constitution-validator` | 📜 Constitution | `general-purpose` |

**Spawn template for each validator:**

```
Task({
  description: "{Perspective} validation",
  prompt: """
  You are the {name} on the {team-name} team.

  CONTEXT:
    - Target files: {file list}
    - Spec files: {spec documents if applicable}
    - Project standards: {from CLAUDE.md, conventions}

  OUTPUT: Return findings in structured format:
    FINDING:
    - status: PASS | WARN | FAIL
    - severity: HIGH | MEDIUM | LOW
    - title: Brief title (max 40 chars)
    - location: file:line
    - issue: One sentence describing what was found
    - recommendation: How to fix

    If no findings: return NO_FINDINGS

  SUCCESS: All {perspective} concerns identified with specific recommendations

  TEAM PROTOCOL:
    - Check TaskList for your assigned validation task
    - Mark in_progress when starting, completed when done
    - Send findings to lead via SendMessage
    - After completing, check TaskList for next available unblocked task
    - If no more tasks, go idle
  """,
  subagent_type: "general-purpose",
  team_name: "{team-name}",
  name: "{validator-name}",
  mode: "bypassPermissions"
})
```

Launch ALL validator teammates simultaneously in a single response with multiple Task calls.

#### Step 4: Monitor & Collect

Messages from validators arrive automatically — the lead does NOT poll.

```
Collection loop:
1. Receive message from validator: "Validation complete. Findings: ..."
2. Receive message from another validator: "Validation complete. Findings: ..."
3. When all validators have reported:
   → Check TaskList to verify all validation tasks are completed
   → Proceed to synthesis
```

If a validator is blocked or reports an issue:
- Missing context → Send DM with the needed information
- Tool failure → Reassign task or handle directly
- After 3 retries for the same task → Skip that perspective and note it in the summary

#### Step 5: Graceful Shutdown

After collecting all findings:

```
1. Verify all tasks completed via TaskList
2. For EACH validator teammate (sequentially):
   SendMessage({
     type: "shutdown_request",
     recipient: "{validator-name}",
     content: "Validation complete. Thank you for your analysis."
   })
3. Wait for each shutdown_response (approve: true)
4. After ALL teammates shut down:
   TeamDelete()
```

If a teammate rejects shutdown: check TaskList for incomplete work, resolve, then re-request.

Continue to **Phase 4: Synthesize & Present**.

---

### Phase 4: Synthesize & Present

This phase is the same for both Standard and Team Mode.

**For Team Mode**, apply the deduplication algorithm before building the summary:

```
Deduplication algorithm:
1. Collect all findings from all validators
2. Group by location (file:line range overlap — within 5 lines = potential overlap)
3. For overlapping findings:
   a. Keep the highest severity version
   b. Merge complementary details from multiple perspectives
   c. Credit both perspectives in the finding
4. Sort by severity (FAIL > WARN > PASS)
5. Build summary table
```

1. **Collect** all findings from validation agents
2. **Deduplicate** overlapping issues
3. **Rank** by severity (HIGH > MEDIUM > LOW)
4. **Group** by category for readability

**Drift-specific synthesis:**
- Categorize by drift type: Scope Creep, Missing, Contradicts, Extra
- Present user decision options

**Constitution-specific synthesis:**
- Separate by level: L1 (autofix), L2 (manual), L3 (advisory)
- L1/L2 are blocking; L3 is informational

### Phase 5: Present Report

```markdown
## Validation: [target]

**Mode**: [Spec | File | Drift | Constitution | Comparison | Understanding]
**Assessment**: ✅ Excellent | 🟢 Good | 🟡 Needs Attention | 🔴 Critical

### Summary

| Perspective | Pass | Warn | Fail |
|-------------|------|------|------|
| ✅ Completeness | X | X | X |
| 🔗 Consistency | X | X | X |
| 📍 Alignment | X | X | X |
| 📐 Coverage | X | X | X |
| 📊 Drift | X | X | X |
| 📜 Constitution | X | X | X |
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
| Consistency | Terminology consistent across docs |

### Verdict

[What was validated and key conclusions]
```

### Phase 6: Next Steps

Use `AskUserQuestion` based on findings:

**If Constitution L1/L2 Violations:**
- "Apply autofixes (L1)" (Recommended)
- "Show me the violations"
- "Skip constitution checks"

**If Drift Detected:**
- "Acknowledge and continue" (log drift, proceed)
- "Update implementation" (implement missing, remove extra)
- "Update specification" (modify spec to match reality)
- "Defer decision" (mark for later review)

**If Spec Issues:**
- "Address failures first"
- "Show detailed findings"
- "Continue anyway"

## Constitution Enforcement

When validating constitution (`$ARGUMENTS` contains "constitution"):

1. **Check for CONSTITUTION.md** at project root
2. **Parse rules** by category (Security, Architecture, etc.)
3. **Apply checks**:
   - Pattern rules: regex match
   - Check rules: semantic analysis
4. **Report by level**:
   - L1: Critical, autofix required
   - L2: Blocking, manual fix required
   - L3: Advisory only

**Integration with other workflows:**
- Called by `/start:implement` at phase checkpoints
- Called by `/start:specify` during SDD phase for architecture alignment

## Drift Detection

When validating drift (`$ARGUMENTS` contains "drift"):

1. **Load specification** (PRD, SDD, PLAN)
2. **Analyze implementation** files
3. **Compare and categorize**:
   - ✅ Aligned: Requirement implemented as specified
   - ❌ Missing: Specified but not implemented
   - ⚠️ Contradicts: Implementation differs from spec
   - 🔶 Extra: Implemented but not in spec
4. **Log decisions** to spec README.md

**Integration with other workflows:**
- Called by `/start:implement` at phase boundaries

## Ambiguity Detection

For spec validation, include ambiguity scoring:

**Vague patterns to detect:**
- Hedge words: "should", "might", "could"
- Vague quantifiers: "fast", "many", "various"
- Open-ended lists: "etc.", "and so on"
- Undefined terms: "the system", "appropriate"

**Scoring:**
- 0-5%: ✅ Excellent clarity
- 5-15%: 🟡 Acceptable
- 15-25%: 🟠 Recommend clarification
- 25%+: 🔴 High ambiguity

## Important Notes

- **Advisory by default** - All findings are recommendations except constitution L1/L2
- **Be specific** - Include file:line for every finding
- **Actionable** - Every finding should have a clear fix
- **Parallel execution** - Launch all applicable perspectives simultaneously
- **Log drift decisions** - Record to spec README for traceability
- **Team mode specifics** - Validators work independently via shared task list; lead handles dedup at synthesis
- **User-facing output** - Only the lead's synthesized output is visible to the user; do not forward raw validator messages
