# Canonical Skill Example: review

A complete skill demonstrating all gold-standard conventions, annotated with convention callouts.

---

```yaml
---
name: review                          # kebab-case, matches entry-point function
description: Multi-agent code review  # trigger-focused, no workflow details
  with specialized perspectives
user-invocable: true
argument-hint: "PR number, branch     # shown in / menu
  name, file path, or 'staged'"
allowed-tools: Task, TaskOutput, ...  # tools used without permission prompts
---
```

## Persona                             ← PICS: P

Act as a code review orchestrator...   ← role + expertise frame

**Review Target**: $ARGUMENTS          ← binds argument to context

## Interface                           ← PICS: I — data shapes, then fn signatures

Finding {                              ← data shape with inlined enums (no type aliases)
  severity: CRITICAL | HIGH | MEDIUM | LOW
  confidence: HIGH | MEDIUM | LOW
  title: String
  location: String
  issue: String
  fix: String
  code_example?: String                ← ? for optional
}

fn gatherContext(target)               ← forward declarations (fn = define, not execute)
fn selectMode()
fn launchReviews(mode)
fn synthesize(findings)
fn nextSteps(verdict)

## Constraints                         ← PICS: C — require/never split

Constraints {
  require {
    Launch ALL applicable review activities simultaneously.
    Every finding must have a specific, implementable fix.
  }
  never {
    Review code yourself — always delegate.     ← enforcement rule moved from Persona
    Present findings without actionable fix.
  }
}

## State                               ← PICS: S — concrete defaults, no infer()

State {
  target = $ARGUMENTS
  perspectives = []                    ← comment explains origin: populated by gatherContext
  mode: Standard | Team                ← chosen by user in selectMode
  findings: [Finding]                  ← collected from agents
}

## Reference Materials                 ← optional, progressive disclosure

See `reference/` directory:
- [Perspectives](reference/perspectives.md)
- [Output Format](reference/output-format.md)
- [Checklist](reference/checklists.md)
- [Classification](reference/classification.md)

## Workflow                            ← fn definitions + entry-point pipe chain

fn gatherContext(target) {             ← fn = definition, not execution
  match (target) {                     ← match for routing
    /^\d+$/       => gh pr diff $target
    "staged"      => git diff --cached
    default       => git diff main...$target
  }

  match (changes) {                    ← conditional perspective activation
    async/await | Promise => +Concurrency
    dependency changes    => +Dependencies
  }
}

fn selectMode() { ... }
fn launchReviews(mode) { ... }

fn synthesize(findings) {
  findings
    |> deduplicate                     ← pipe operator for data pipeline
    |> sort(by: [severity desc])
    |> assignIds(pattern: "$severityLetter$number")
    |> buildSummaryTable

  verdict = match (criticalCount, highCount, mediumCount) {
    (> 0, _, _)  => 🔴 REQUEST CHANGES
    (0, 0, 0)    => ✅ APPROVE
  }
}

fn nextSteps(verdict) { ... }

review(target) {                       ← entry point: NO fn keyword, name matches skill
  gatherContext(target) |> selectMode |> launchReviews |> synthesize |> nextSteps
}
