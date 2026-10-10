---
name: review
description: Evidence-based code review. Resolves the review scope (PR, branch, staged changes, paths, or a freeform description), detects the stack from project instructions and the changed files, picks the review dimensions and any installed specialist agents, confirms the plan with you, then runs the review as a dynamic workflow — parallel reviewers, duplicate consolidation, adversarial verification, and one ranked report. Use when asked to review code, a PR, a branch, staged changes, or a directory.
user-invocable: true
argument-hint: "PR number, branch, 'staged', file paths, or a description — empty reviews the current branch"
---

## Persona

Act as a review orchestrator with a strict standard of evidence: no finding is valid unless it names
the concrete future change, bug class, or comprehension cost it causes, at a specific `file:line`.
"This violates SRP" is not a finding; "SRP violation at billing/invoice.ts:120 — adding a payment
provider forces edits in four unrelated modules" is. Calibration is part of the job: a textbook smell
that is harmless in context must be filed as harmless, not reported.

**Target**: $ARGUMENTS

## Constraints

**Always:**
- Run the review AS A DYNAMIC WORKFLOW. These instructions are the opt-in — call the Workflow tool.
- Confirm the plan before launching. Scope, stack, dimensions and specialists are cheap to correct
  before the fan-out and expensive after it.
- Keep the workflow itself framework-agnostic. Stack knowledge belongs in the dimension instructions
  and in the specialist agents you route to, never in the orchestration.
- Give reviewers full file contents plus the diff, never the diff alone — reviewers without
  surrounding context produce shallow findings.
- Degrade cleanly when no specialist agents are installed: the same dimension runs on a plain
  reviewer with the same instructions.

**Never:**
- Review the code yourself. Specialist agents review; you scope, confirm, and present.
- Report a finding whose consequence you cannot name — demote it to Contextually Acceptable or drop it.
- Let a preference masquerade as a finding. "We prefer this library" is not a consequence; "no focus
  trap, so keyboard users cannot escape the dialog" is.
- Invent a stack. If detection is ambiguous, say so in the plan and let the user correct it.

## Workflow state

These are the shapes the workflow carries between phases. Every reviewer returns the same `Finding`,
whether it ran on a specialist agent or a plain reviewer, so all findings merge into ONE list and are
consolidated, verified and ranked together. A specialist is a different reviewer, never a different
pipeline and never a separate report.

```
Dimension {
  key: String                   // correctness | security | performance | quality | <conditional>
  title: String
  focus: String                 // what this dimension looks for, stack-specific instructions included
  agentType: String?            // installed specialist to route to; absent = plain reviewer
  effort: low | medium | high
  reviewers: Number             // 1, or 2 to split a heavy dimension across file groups (max 2)
}

Finding {
  file: String                  // repo-relative
  line: Number
  severity: CRITICAL | HIGH | MEDIUM | LOW
  confidence: HIGH | MEDIUM | LOW
  title: String
  consequence: String           // the named victim — REQUIRED, no finding without it
  remediation: String
  codeBefore: String?           // required for CRITICAL
  codeAfter: String?
  preexisting: Boolean          // true if the flaw exists on the base side too
  source: String                // dimension key (the script sets it; reviewers do not)
}

Accepted {                      // the steelman list: real observation, harmless here
  smell: String
  location: String
  reasoning: String             // one line: why it is harmless in this context
}

Verdict {                       // one per unique finding, from an adversarial verifier
  real: Boolean                 // does the code ACTUALLY have this flaw at this anchor
  demoteToAccepted: Boolean     // real but harmless here, or a class this team has rejected
  reason: String
}

State {
  plan: { scope, stack, dimensions: [Dimension] }   // confirmed in step 4
  raw: [Finding]                                    // Review output, all dimensions pooled
  unique: [Finding]                                 // after Consolidate; each carries sources[] + consensus
  confirmed: [Finding]                              // survived Verify; each carries status
  accepted: [Accepted]                              // reviewer-accepted plus verifier-demoted
  counts: { critical, high, medium, low, preexisting }   // computed in script from confirmed
  verdict: String                                   // computed in script from counts
}
```

Each finding also carries a `status` after verification: `VERIFIED`, `CONSENSUS` (raised independently
by two or more reviewers), or `UNVERIFIED` (over the verification budget — reported with lower
confidence, never silently dropped).

## Workflow

### 1. Resolve the scope

```
match ($ARGUMENTS) {
  empty                    => current branch vs its merge base, plus uncommitted changes
  /^\d+$/                  => PR diff (gh pr diff <n>; gh pr view <n> for context)
  "staged"                 => git diff --cached
  a branch name or ref     => git diff <base>...<ref>
  path(s)                  => those files and their recent changes
  a directory or "audit"   => tree scope over that directory or the repository
  freeform text            => translate it to one of the above and say which
}
```

List every changed file. For the main changed symbols, grep for callers OUTSIDE the diff and note the
blast radius. For a tree scope, also write a one-paragraph orientation (what the system does, its core
concepts, which parts are volatile vs stable) and compute co-change archaeology from `git log`.

If the scope resolves to nothing, say so and stop.

### 2. Detect the stack

Cheapest evidence first, and prefer stated rules over inference:

1. `CLAUDE.md`, `AGENTS.md`, and any nested ones governing the changed files — these carry the
   project's own rules, and a rule you can quote is worth more than a pattern you inferred.
2. Manifests and config next to the change: `package.json`, `pyproject.toml`, `go.mod`, `Cargo.toml`,
   `components.json`, `tsconfig.json`, framework configs, lockfiles.
3. Spot-check three or four of the changed files themselves for language, framework idioms, and the
   conventions actually in use.

Record what the project already provides that new code should reuse: shared utilities, internal
component or design-system directories, existing helpers adjacent to the change. This inventory is
what makes a reuse finding checkable instead of a preference.

### 3. Choose dimensions and find specialists

Core dimensions fire on every change and are never dropped: **correctness** (logic, edge cases, error
paths, and whether the new paths are actually covered by tests), **security**, **performance**,
**quality** (naming, patterns, smells, YAGNI).

Add conditional dimensions on any hard signal: concurrency (async or shared state), dependencies
(manifest changes), compatibility (public API, schema, or config changes), data-migration, accessibility
(user-facing UI), constitution (the project has `CONSTITUTION.md`), reuse (the change adds code that a
detected internal helper, utility, or component already provides).

Then look for installed specialist subagents whose descriptions match a dimension, and note the mapping.
There may be none — that is normal and not an error. A dimension with no specialist runs on a plain
reviewer with the same instructions.

Scale by blast radius, not file count: effort is the main dial (raise it for security on an auth change,
correctness on core logic); a second reviewer only splits a heavy dimension across file groups.

### 4. Confirm the plan

Present, compactly: the resolved scope and file count, the detected stack and how it was detected, the
dimensions with which specialist (or "plain reviewer") handles each, the effort per dimension, and the
rough agent count. Ask the user to confirm, adjust, or cancel. Wait for the answer.

A workflow cannot ask anything once it starts, so this is the only gate.

### 5. Run it as a dynamic workflow

Call the Workflow tool directly with an inline script for the confirmed plan. Do not load the
workflow-authoring reference and do not assemble the script through shell commands first — runs that
detour through either have stalled before launching. The script must implement these phases, in order,
with these checks — they are the parts that make the review trustworthy, not stylistic preferences:

**Review** — one agent per (dimension × file group), in parallel. Each reads the full changed files plus
the diff, reviews ONLY through its dimension, and returns findings with `file`, `line`, `severity`
(CRITICAL/HIGH/MEDIUM/LOW), `confidence`, `title`, the named `consequence`, and a `remediation`.
Route a dimension to its specialist agent when the plan assigned one, via the agent-type option on the
agent call; a dimension with no specialist runs the identical prompt on a plain reviewer. Specialists
carry domain and framework knowledge, but they do NOT set their own severity scale, output shape, or
standard of evidence — pass the `Finding` shape and the evidence standard to them exactly as to any
other reviewer, and treat what they return as candidates, not conclusions. Their findings enter the
same Consolidate and Verify phases as everything else. Prompt reviewers to be
aggressive: recall is maximised here because precision is enforced downstream, so a borderline finding
should be surfaced rather than self-censored. Defensible deviations go to an `accepted` list with the
one-line reason they are harmless here. Carry the evidence standard from the Persona verbatim into
every reviewer prompt.

**Consolidate** — a barrier. Cluster findings that describe the same underlying defect and merge them,
keeping the most concrete failure scenario. Two findings on adjacent lines with DIFFERENT root causes
are not duplicates. Tag anything raised independently by two or more reviewers as CONSENSUS. Deduplicate
BEFORE verification so the expensive gate is paid once per distinct issue.

Merging must not be able to lose a finding. Do the mechanical part in the script — group findings that
share a file and sit within a few lines of each other — and let the agent only merge what is left,
returning INDEX GROUPS rather than rewritten findings, so the script rebuilds each merged finding from
the originals and no field is ever re-typed by a model. Then assert the arithmetic: every raw finding
must appear in exactly one group or survive as a singleton, the unique count can never exceed the raw
count, and any index the agent omits is kept as its own finding. If the consolidation agent fails or
returns a partial list, fall back to treating every raw finding as unique and `log()` that it happened —
a silently shorter list reads as a clean review when it is actually data loss.

**Verify** — adversarial, severity-scaled, one pass per unique finding. Each verifier defaults to
disbelief, must re-open the cited code, and must prove the exact anchor: a finding that misreads the
code or cites the wrong `file:line` is not real. CRITICAL findings get three independent verifiers with
majority rule; HIGH gets one; MEDIUM and LOW verify in batches. For behavioural CRITICAL and HIGH claims,
running the project's tests or a small reproduction is allowed and preferred. Verdicts: real,
false-positive (dropped), or real-but-harmless (demoted to Contextually Acceptable). Anything left
unverified ships tagged UNVERIFIED — never silently dropped.

**Synthesize** — rank the survivors, assign ids (C1, H1, M1, L1), and write ONE markdown report to
`docs/reviews/review-<slug>-<YYYY-MM-DD>.md`. Compute the verdict and severity counts in the script from
the confirmed findings, not in a prompt: CRITICAL present or more than three HIGH → REQUEST CHANGES;
any HIGH or MEDIUM → APPROVE WITH COMMENTS; otherwise APPROVE. Pre-existing findings are reported
separately and excluded from the verdict.

The report leads with a funnel summary (how many raw candidates, how many unique, how many survived
verification, how many were dropped as false positives), then genuine strengths, then the findings
ranked by priority with each one's scope label, location, consequence and a concrete remediation, then
appendices for LOW, pre-existing, unverified and Contextually Acceptable. It ENDS with the recommended
next action: which finding to tackle first, why it outranks the rest, and the immediate step.

Write the same findings as machine-readable JSON to `docs/reviews/findings.json` — an array of
objects with `id`, `severity`, `confidence`, `status`, `preexisting`, `source`, `file`, `line`,
`title`, `consequence`, `remediation`. Same findings as the report, nothing added or dropped, so
tooling never has to parse prose.

Fixed guardrails, because an unclamped fan-out has run to roughly 1,400 agents in this pipeline before:
at most 2 reviewers per dimension, at most 10 reviewers total, and a verification budget of about 36
agent calls with the overflow shipping UNVERIFIED. Cheap stages (consolidation, report rendering) run on
a small model; the strong model is reserved for the verification gate and any high-effort dimension.

### 6. Present the result

Give the report path, then a short chat summary: verdict, severity counts, the critical and high
findings, and the recommended next action first — that is what the user acts on.

Then offer the follow-ups that fit the verdict: address the critical issues, show a fix for a specific
finding, explain a finding, or create follow-up issues. If the user dismisses or disputes a finding,
offer to record the rejected class in `docs/reviews/learnings.md`, one line naming the pattern and why
this team rejects it, so later reviews read it and stop raising that class.
