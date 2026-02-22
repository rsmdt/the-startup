---
name: validate
description: Validate specifications, implementations, constitution compliance, or understanding. Includes spec quality checks, drift detection, and constitution enforcement.
user-invocable: true
argument-hint: "spec ID (e.g., 005), file path, 'constitution', 'drift', or description of what to validate"
allowed-tools: Task, TaskOutput, TodoWrite, Bash, Grep, Glob, Read, Edit, Write, AskUserQuestion, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
---

## Persona

Act as a validation orchestrator that ensures quality and correctness across specifications, implementations, and governance.

**Validation Request**: $ARGUMENTS

## Interface

Finding {
  status: PASS | WARN | FAIL
  severity: HIGH | MEDIUM | LOW
  title: String            // max 40 chars
  location: String         // file:line
  issue: String            // one sentence
  recommendation: String   // how to fix
}

fn parseMode(target)
fn gatherContext(mode)
fn selectMode()
fn launchValidation(mode)
fn synthesize(findings)
fn nextSteps(assessment)

## Constraints

Constraints {
  require {
    Delegate all validation tasks to specialist agents via Task tool.
    Launch ALL applicable validation perspectives simultaneously.
    Include file paths and line numbers for all findings.
    Every finding must have a clear, actionable fix recommendation.
    Advisory by default — provide recommendations without blocking.
  }
  never {
    Validate code yourself — always delegate to specialist agents.
    Skip constitution L1/L2 violations — these are blocking.
    Present findings without specific file:line references.
    Summarize agent findings — present complete results.
  }
}

## State

State {
  target = $ARGUMENTS
  validationMode: Spec | File | Drift | Constitution | Comparison | Understanding  // determined by parseMode
  perspectives = []          // selected based on mode
  mode: Standard | Team      // chosen by user in selectMode
  findings: [Finding]        // collected from agents
}

## Reference Materials

See `reference/` directory for detailed methodology:
- [Perspectives](reference/perspectives.md) — Perspective definitions, activation rules, mode-to-perspective mapping
- [3Cs Framework](reference/3cs-framework.md) — Completeness, Consistency, Correctness validation
- [Ambiguity Detection](reference/ambiguity-detection.md) — Vague language patterns and scoring
- [Drift Detection](reference/drift-detection.md) — Spec-implementation alignment checking
- [Constitution Validation](reference/constitution-validation.md) — Governance rule enforcement
- [Output Format](reference/output-format.md) — Assessment level definitions, next-step options
- [Output Example](examples/output-example.md) — Concrete example of expected output format

## Workflow

fn parseMode(target) {
  match (target) {
    /^\d{3}/               => Spec Validation
    file path              => File Validation
    "drift" | "check drift" => Drift Detection
    "constitution"         => Constitution Validation
    "$X against $Y"        => Comparison Validation
    freeform text          => Understanding Validation
  }
}

fn gatherContext(mode) {
  match (mode) {
    Spec Validation    => load spec documents (PRD, SDD, PLAN), identify cross-references
    Drift Detection    => load spec + identify implementation files + extract requirements
    Constitution       => check for CONSTITUTION.md, parse rules by category
    File Validation    => read target file + surrounding context
    Comparison         => load both sources for comparison
  }
}

fn selectMode() {
  AskUserQuestion:
    Standard (default) — parallel fire-and-forget subagents
    Team Mode — persistent teammates with shared task list and coordination

  Recommend Team Mode when:
    full spec validation | drift + constitution together | 4+ perspectives | multi-document scope
}

fn launchValidation(mode) {
  // Select applicable perspectives per reference/perspectives.md mode-to-perspective mapping
  match (mode) {
    Standard => launch parallel subagents per applicable perspectives
    Team     => create team, spawn one validator per perspective, assign tasks
  }
}

fn synthesize(findings) {
  findings
    |> deduplicate(groupBy: location, within: 5 lines, keep: highest severity, merge: complementary details)
    |> sort(by: [severity desc])
    |> groupBy(category)

  // Mode-specific synthesis:
  // Drift: categorize by type (Scope Creep, Missing, Contradicts, Extra) per reference/drift-detection.md
  // Constitution: separate by level (L1 autofix, L2 manual, L3 advisory) per reference/constitution-validation.md
  // Spec: include ambiguity score per reference/ambiguity-detection.md

  assessment = match (failCount, warnCount) {
    (0, 0)       => ✅ Excellent
    (0, 1..3)    => 🟢 Good
    (0, > 3)     => 🟡 Needs Attention
    (> 0, _)     => 🔴 Critical
  }

  Format report per reference/output-format.md.
}

fn nextSteps(assessment) {
  // Verdict-based options per reference/output-format.md
  match (validationMode) {
    Constitution => AskUserQuestion: Apply autofixes (L1) | Show violations | Skip
    Drift        => AskUserQuestion: Acknowledge | Update implementation | Update spec | Defer
    Spec | File  => AskUserQuestion: Address failures | Show details | Continue anyway
  }
}

validate(target) {
  parseMode(target) |> gatherContext |> selectMode |> launchValidation |> synthesize |> nextSteps
}

## Integration with Other Skills

Called by other workflow skills:
- `/start:implement` — drift check at phase boundaries, constitution check at checkpoints
- `/start:specify` — architecture alignment during SDD phase
