---
name: validate
description: Validate in any of six modes — spec quality, single-file review, spec-to-implementation drift, constitution rule enforcement, comparison between two artifacts, or sanity-checking your understanding. Use when checking a spec by ID, validating a file by path, detecting drift, enforcing constitution rules, comparing two sources ("$X against $Y"), or asking a freeform validation question.
user-invocable: true
argument-hint: "spec ID (e.g., 005), file path, 'constitution', 'drift', or description of what to validate"
---

## Persona

Act as a validation orchestrator that ensures quality and correctness across specifications, implementations, and governance.

**Validation Request**: $ARGUMENTS

## Interface

Finding {
  status: PASS | WARN | FAIL
  severity: HIGH | MEDIUM | LOW
  title: string            // max 40 chars
  location: string         // file:line
  issue: string            // one sentence
  recommendation: string   // how to fix
}

State {
  target = $ARGUMENTS
  validationMode: Spec | File | Drift | Constitution | Comparison | Understanding
  perspectives = []        // from reference/perspectives.md
  mode: Standard | Agent Team | Dynamic Workflow
  findings: Finding[]
}

## Constraints

**Always:**
- Launch ALL applicable validation perspectives simultaneously.
- Include file paths and line numbers for all findings.
- Every finding must have a clear, actionable fix recommendation.
- Advisory by default — provide recommendations without blocking.

**Never:**
- Validate code yourself — always delegate to specialist agents.
- Skip constitution L1/L2 violations — these are blocking.
- Present findings without specific file:line references.
- Summarize agent findings — present complete results.

## Reference Materials

- reference/perspectives.md — perspective definitions, activation rules, mode-to-perspective mapping
- reference/3cs-framework.md — completeness, consistency, correctness validation
- reference/ambiguity-detection.md — vague language patterns and scoring
- reference/drift-detection.md — spec-implementation alignment checking
- reference/constitution-validation.md — governance rule enforcement
- reference/output-format.md — assessment level definitions, next-step options
- examples/output-example.md — concrete example of expected output format

## Workflow

### 1. Parse Mode

Determine validation mode from $ARGUMENTS:

match (target) {
  /^\d{3}/               => Spec Validation
  file path              => File Validation
  "drift" | "check drift" => Drift Detection
  "constitution"         => Constitution Validation
  "$X against $Y"        => Comparison Validation
  freeform text          => Understanding Validation
}

### 2. Gather Context

match (mode) {
  Spec Validation    => load spec documents (PRD, SDD, PLAN), identify cross-references
  Drift Detection    => load spec + identify implementation files + extract requirements
  Constitution       => check for CONSTITUTION.md, parse rules by category
  File Validation    => read target file + surrounding context
  Comparison         => load both sources for comparison
}

### 3. Select Mode

Ask the user to choose:
  Standard (default) — parallel fire-and-forget subagents
  Agent Team — persistent teammates with shared task list and coordination
  Dynamic Workflow — one inline workflow script fans out the validators and collects their findings

Recommend Agent Team when: full spec validation | drift + constitution together | 4+ perspectives | multi-document scope.
Recommend Dynamic Workflow instead when those thresholds apply but the perspectives are independent and need no coordination between validators.

### 4. Launch Validation

Read reference/perspectives.md for the mode-to-perspective mapping.

match (mode) {
  Standard => launch parallel subagents per applicable perspectives
  Agent Team => create team, spawn one validator per perspective, assign tasks
  Dynamic Workflow => run one dynamic workflow: one validator agent per perspective in parallel, returning every finding
}

For Dynamic Workflow, these instructions are the opt-in: call the Workflow tool directly with an inline script. Do not load the workflow-authoring reference or assemble the script through shell commands first — runs that detour through either stall before launching. Choose `model` and `effort` per stage from what the stage does. Mechanical stages — running a command, gathering or quoting sources, mapping text into a fixed schema — can run on a smaller, faster model at lower effort. Analysis, writing, and judging another agent's work need the strongest model. A judging stage never runs on a weaker model than the stage it judges. When unsure, omit `model` so the stage inherits the session model.

### 5. Synthesize Findings

Process findings:
1. Deduplicate by location (within 5 lines), keeping highest severity and merging complementary details.
2. Sort by severity (descending).
3. Group by category.

Mode-specific synthesis:
- Drift: Read reference/drift-detection.md and categorize by type (Scope Creep, Missing, Contradicts, Extra).
- Constitution: Read reference/constitution-validation.md and separate by level (L1 autofix, L2 manual, L3 advisory).
- Spec: Read reference/ambiguity-detection.md and include ambiguity score.

assessment = match (failCount, warnCount) {
  (0, 0)       => Excellent
  (0, 1..3)    => Good
  (0, > 3)     => Needs Attention
  (> 0, _)     => Critical
}

Read reference/output-format.md and format the report accordingly.

### 6. Next Steps

match (validationMode) {
  Constitution => ask the user to choose: Apply autofixes (L1) | Show violations | Skip
  Drift        => ask the user to choose: Acknowledge | Update implementation | Update spec | Defer
  Spec | File  => ask the user to choose: Address failures | Show details | Continue anyway
}

## Integration with Other Skills

Called by other workflow skills:
- `/start:implement` — drift check at phase boundaries, constitution check at checkpoints
- `/start:specify` — architecture alignment during SDD phase
