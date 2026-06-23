---
name: review
description: Multi-agent code review with specialized perspectives (security, performance, patterns, simplification, tests). Use when asked to review a PR, a branch, staged changes, a diff, or specific files before merge.
user-invocable: true
argument-hint: "PR number, branch name, file path, or 'staged' for staged changes"
---

Orchestrate a multi-perspective review of **$ARGUMENTS**. Specialist agents do the reviewing — never review the code yourself. Only your synthesized report reaches the user, never raw reviewer output.

## Workflow

### 1. Resolve Target

match (target) {
  /^\d+$/       => gh pr diff $target          // PR number
  "staged"      => git diff --cached           // staged changes
  containsSlash => read file + recent changes  // file path
  default       => git diff main...$target     // branch name
}

Retrieve the **full contents of every changed file**, not just the diff — reviewers without surrounding context produce shallow findings. If the target resolves to no changes, say so and stop.

### 2. Activate Perspectives

Always active:

| Perspective | Intent | What to Look For |
|-------------|--------|------------------|
| 🔐 Security | Find vulnerabilities before production | Auth/authz gaps, injection risks, hardcoded secrets, input validation, CSRF, cryptographic weaknesses |
| 🔧 Simplification | Challenge unnecessary complexity | YAGNI violations, over-engineering, premature abstraction, dead code, "clever" code that should be obvious |
| ⚡ Performance | Identify efficiency issues | N+1 queries, algorithm complexity, resource leaks, blocking operations, caching opportunities |
| 📝 Quality | Ensure code meets standards | SOLID violations, naming issues, error handling gaps, pattern inconsistencies, code smells |
| 🧪 Testing | Verify adequate coverage | Missing tests for new code paths, uncovered edge cases, test quality issues |

Add when the changes match:

| Perspective | Activation Rule |
|-------------|-----------------|
| 🧵 Concurrency | Any use of async/await, Promises, threading, shared state, or parallel operations — one match suffices |
| 📦 Dependencies | Changes to package.json, requirements.txt, go.mod, Cargo.toml, etc. |
| 🔄 Compatibility | Public API, database schema, config format, or migration changes |
| ♿ Accessibility | Frontend/UI component changes |
| 📜 Constitution | Project has CONSTITUTION.md |

### 3. Launch Reviews

Ask the user to choose:
  Standard (default) — parallel fire-and-forget subagents
  Agent Team — persistent teammates with peer coordination (recommend when files > 10, perspectives >= 4, cross-domain, or constitution active)

Launch ALL applicable perspectives **simultaneously in a single response**. Each reviewer receives the full file contents, the diff, and its perspective row (intent + what to look for). For a reviewer's domain-specific checklist, point it at reference/checklists.md.

### 4. Synthesize Findings

1. Deduplicate by location (same file, within 5 lines) — keep the highest severity, merge complementary details.
2. Sort by severity descending, then confidence descending.
3. Assign IDs: severity letter + number (C1, C2, H1, M1, L1...).
4. Determine the verdict:

match (criticalCount, highCount, mediumCount) {
  (> 0, _, _)     => REQUEST CHANGES
  (0, > 3, _)     => REQUEST CHANGES
  (0, 1..3, _)    => APPROVE WITH COMMENTS
  (0, 0, > 0)     => APPROVE WITH COMMENTS
  (0, 0, 0)       => APPROVE
}

Report with: verdict, findings table, code examples, strengths. One example row sets the format:

| ID | Severity | Confidence | Finding | Remediation |
|----|----------|------------|---------|-------------|
| C1 | CRITICAL | HIGH | SQL injection *(api/users.ts:42)* | Use parameterized query *(user input concatenated into SQL)* |

- Every finding needs an actionable remediation — a finding without a fix is noise.
- Code examples (before/after): REQUIRED for Critical, include for High when the fix is non-obvious, omit for Medium/Low.
- Always close with a **Strengths** section highlighting what the code does well.
- When unsure how to grade a finding, read reference/classification.md.

### 5. Next Steps

Ask the user to choose, based on verdict:

match (verdict) {
  REQUEST CHANGES       => Address critical issues first | Show fixes for [issue] | Explain [finding] in detail
  APPROVE WITH COMMENTS => Apply suggested fixes | Create follow-up issues | Proceed without changes
  APPROVE               => Add to PR comments (if PR review) | Done
}

## Deeper Reference (load only when needed)

- reference/checklists.md — per-domain checklists (security, performance, quality, test coverage) for reviewers that need them
- reference/classification.md — severity/confidence definitions and edge-case classification matrix
- examples/output-example.md — a full rendered report
