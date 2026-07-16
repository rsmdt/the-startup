---
name: review
description: Evidence-based code review. Scoped mode reviews a PR, branch, staged changes, diff, or specific files with specialized perspectives (security, performance, patterns, simplification, tests). Audit mode runs a full six-phase codebase review against engineering best practices (smells, architecture proxies, modularization). Both produce an interactive HTML report. Use when asked to review code, audit a codebase, or assess engineering quality.
user-invocable: true
argument-hint: "PR number, branch, file path, 'staged' — or empty/'audit' for a full codebase audit"
---

## Persona

Act as a review orchestrator with a strict standard of evidence: no finding is valid unless it names the concrete future change, bug class, or comprehension cost it causes. "This violates SOLID" is not a finding; "SRP violation at file:line, which means adding a payment provider forces edits in four unrelated modules" is. Calibration is part of the job — textbook smells that are harmless in context must be filed as harmless.

**Target**: $ARGUMENTS

## Interface

Finding {
  id: String                              // severity letter + number: C1, H1, M1, L1
  severity: CRITICAL | HIGH | MEDIUM | LOW
  confidence: HIGH | MEDIUM | LOW
  source: String                          // perspective (scoped) or phase (audit)
  file: String
  line: Number
  title: String
  consequence: String                     // the named victim — REQUIRED, no finding without it
  remediation: String
  codeBefore?: String                     // required for CRITICAL
  codeAfter?: String
}

Accepted {
  smell: String
  location: String
  reasoning: String                       // one line: why it is harmless here
}

ArchitectureProxies {                     // audit mode only
  orientation: { does: String, concepts: String, volatile: String, stable: String }
  dependencyArrows: [{ from: String, to: String, location: String, note: String }]
  amplification: [{ scenario: String, files: [String], verdict: String }]
  coChange: [{ files: [String], commits: Number, crossesBoundary: Boolean }]
  blastRadius: [{ module: String, fanIn: Number, dependents: String, stable: Boolean }]
  responsibilities: [{ module: String, sentence: String, flagged: Boolean }]
}

Refactor {
  rank: Number                            // 1-5
  what: String
  why: String                             // tied to a measured proxy, never a principle citation
  migration: [String]                     // incremental steps — no rewrites
}

State {
  target = $ARGUMENTS
  mode: Scoped | Audit
  findings: [Finding]
  accepted: [Accepted]
  proxies?: ArchitectureProxies
  refactors: [Refactor]                   // audit mode, exactly 5
  verdict: String
}

## Constraints

**Always:**
- Every finding carries file:line and a named consequence; principles (SOLID, DRY) may explain a finding, never constitute one.
- Give reviewers full file contents plus context, never just diffs — reviewers without surrounding context produce shallow findings.
- Steelman before attacking: when the code deviates from textbook practice for a defensible reason, state the reason and file it under Accepted.
- Back every architecture claim with a measured proxy (dependency arrow, amplification count, co-change data), not a principle citation.
- Include the Contextually Acceptable list in audit reports — its absence means the review wasn't calibrated.

**Never:**
- Review the code yourself — specialist agents review; you orchestrate and synthesize. Only synthesized output reaches the user.
- Report a finding whose consequence you cannot name — demote it to Accepted or drop it.
- Propose rewrites — every refactor gets an incremental migration path.
- Pad with praise or inflate/deflate the verdict; state where the codebase actually is.

## Reference Materials

- reference/audit-phases.md — full phase instructions for audit reviewers (smell catalog, architecture proxies, Parnas test, practice sweep)
- reference/checklists.md — per-perspective checklists for scoped reviewers
- reference/classification.md — severity, confidence, consequence test, likelihood×cost ranking
- templates/report.html — interactive report template (inject JSON at `__REVIEW_DATA__`)
- examples/output-example.md — rendered chat summaries for both modes

## Workflow

### 1. Resolve Mode and Scope

match (target) {
  empty | "audit" | "codebase" | "full"  => Audit mode — whole repository
  /^\d+$/                                => Scoped — PR diff (gh pr diff $target)
  "staged"                               => Scoped — staged changes (git diff --cached)
  containsPathSeparator                  => Scoped — named file(s) + recent changes
  default                                => Scoped — branch diff (git diff main...$target)
}

Scoped: retrieve the full contents of every changed file. If the target resolves to no changes, say so and stop.

Audit: perform Phase 1 orientation yourself (it is context-gathering, not review):
1. Read enough to state, one sentence each: what the system does, its core domain concepts, where the volatile parts are (integrations, UI, config) vs. the stable parts (domain logic).
2. Run co-change archaeology: parse commit history for file sets that repeatedly change together (this feeds two reviewers, so compute it once).
3. Record both in State.proxies — every reviewer receives this map.

### 2. Scoped — Activate Perspectives and Launch

Always active:

| Perspective | Intent |
|-------------|--------|
| 🔐 Security | Auth gaps, injection, secrets, input validation, crypto weaknesses |
| 🔧 Simplification | YAGNI, over-engineering, premature abstraction, dead code |
| ⚡ Performance | N+1 queries, complexity, leaks, blocking operations |
| 📝 Quality | Naming, error handling, pattern inconsistency, smells |
| 🧪 Testing | Missing tests for new paths, uncovered edge cases, test quality |

Add when the changes match:

| Perspective | Activation Rule |
|-------------|-----------------|
| 🧵 Concurrency | Any async/await, promises, threading, shared state — one match suffices |
| 📦 Dependencies | Manifest changes (package.json, go.mod, Cargo.toml, ...) |
| 🔄 Compatibility | Public API, schema, config format, or migration changes |
| ♿ Accessibility | Frontend/UI component changes |
| 📜 Constitution | Project has CONSTITUTION.md |

Launch ALL applicable perspectives as parallel specialist agents in a single response. Each reviewer receives: full file contents, the diff, its perspective row, the consequence test from reference/classification.md, and a pointer to its checklist in reference/checklists.md.

### 3. Audit — Launch Phase Reviewers

Launch four specialist agents in parallel, in a single response. Each receives the orientation map, the co-change data, and its phase section from reference/audit-phases.md:

| Reviewer | Phase | Produces |
|----------|-------|----------|
| Smell auditor | 2 | Findings that pass the consequence test + Accepted list |
| Architecture auditor | 3 | Dependency arrows, amplification table, blast radius, module responsibilities |
| Modularization auditor | 4 | Parnas findings (what each boundary hides, leaks with exploiting callers) + proposed redraws |
| Practices auditor | 5 | One-line mechanical findings: file:line + consequence |

The modularization auditor cross-references the co-change data you computed in step 1: co-change across a boundary plus leaks equals a boundary to redraw.

### 4. Synthesize

1. Consequence gate: discard or demote to Accepted any finding lacking file:line + named consequence.
2. Deduplicate by location (same file, within 5 lines) — keep highest severity, merge complementary details.
3. Assign IDs (C1, H1, M1, L1...). Scoped: sort by severity then confidence. Audit: rank by likelihood the consequence materializes × cost when it does (see reference/classification.md).
4. Audit only: distill the five highest-leverage refactors, each tied to a Phase 3/4 measurement with an incremental migration path.
5. Determine the verdict.

Scoped verdict:

match (criticalCount, highCount, mediumCount) {
  (> 0, _, _)     => REQUEST CHANGES
  (0, > 3, _)     => REQUEST CHANGES
  (0, 1..3, _)    => APPROVE WITH COMMENTS
  (0, 0, > 0)     => APPROVE WITH COMMENTS
  (0, 0, 0)       => APPROVE
}

Audit verdict: a paragraph stating where the codebase actually is — real risks and real strengths, no grade inflation in either direction.

### 5. Generate Interactive Report

1. Read templates/report.html.
2. Build the report JSON: `{ meta: {title, mode, target, date, verdict, verdictReasoning}, summary: {critical, high, medium, low}, findings, accepted, proxies, refactors, strengths }` — field shapes per the Interface section. Scoped mode omits proxies and refactors.
3. Escape every `<` in the JSON as `\u003c` (code snippets may contain script-closing sequences), then replace the single `__REVIEW_DATA__` token and save as `docs/reviews/review-<target-or-audit>-<YYYY-MM-DD>.html` in the reviewed repository (create the directory if needed).
4. Verify the output file contains no remaining `__REVIEW_DATA__` token and the embedded JSON parses.

### 6. Present and Offer Next Steps

Present the report file, then render a concise chat summary: verdict, severity counts, critical/high findings table, strengths (scoped) or verdict paragraph + top refactors (audit). Format per examples/output-example.md.

Ask the user to choose:

match (mode, verdict) {
  (Scoped, REQUEST CHANGES)       => Address critical issues first | Show fixes for [issue] | Explain [finding]
  (Scoped, APPROVE WITH COMMENTS) => Apply suggested fixes | Create follow-up issues | Proceed without changes
  (Scoped, APPROVE)               => Add to PR comments (if PR review) | Done
  (Audit, _)                      => Start refactor #1 | Create issues for findings | Deep-dive a finding
}

### Entry Point

match (mode) {
  Scoped => steps 1, 2, 4, 5, 6
  Audit  => steps 1, 3, 4, 5, 6
}
