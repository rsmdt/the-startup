# Scenario Guide

How to write effective holdout acceptance scenarios. An acceptance scenario is a plain-English description of expected system behavior, tested through external interfaces — typically verifying one EARS statement from its slice. During implementation, an evaluation agent reads it and exercises the running system to decide pass/fail. Each scenario is created as a **scenario child task** (`task.js --type scenario`) under its slice.

---

## What Makes a Good Scenario

**Good scenarios**:
- Test observable behavior through the API, UI, or CLI
- Exercise **one slice's own** external interface, so it is verifiable independently
- Reference actual endpoints, data shapes, and conventions from the codebase
- Are specific enough that pass/fail is unambiguous
- Each test one thing

**Bad scenarios**:
- Reference internal implementation details (class names, function signatures, file paths)
- Are vague ("the system should work correctly")
- Test multiple unrelated behaviors at once
- Duplicate what unit tests already cover

## Grounding in the Codebase

Before writing scenarios, read:
1. **Project instructions** (AGENTS.md / CLAUDE.md) — architecture, conventions, patterns
2. **Existing test files** — the framework and existing test patterns
3. **API routes / endpoints** — reference actual paths, not hypothetical ones
4. **Data models** — reference actual field names and types

A scenario like "POST to /api/v1/sql/validate" is grounded. "POST to the validation endpoint" is not.

## Priority Levels

| Priority | Description | When to Use |
|----------|-------------|-------------|
| P0 | Critical path | The primary happy path. If this fails, the feature is broken. |
| P1 | Important | Error handling, edge cases that affect users. |
| P2 | Edge case | Boundary conditions, unusual inputs, rare paths. |

Every slice should have at least one P0 scenario. P1 and P2 increase confidence.

## Scenario Count Guidelines

- **Minimum**: 1 scenario per slice (P0 happy path)
- **Typical**: 3-5 scenarios per slice (P0 + P1 error cases)
- **Thorough**: 5-10 scenarios per slice (P0 + P1 + P2 edge cases)

More scenarios = higher confidence but longer evaluation. Start with P0 + P1, add P2 if useful.

## Why Holdout — the Information Barrier

Acceptance scenarios are a **holdout evaluation set**. When `/implement` builds a slice, the code agent that writes the implementation never reads its scenario child tasks — only the evaluation agent does. Think of it as the train/test split in ML:

- The requirements + solution + slice brief are "training data" — they tell the implementer what to build.
- Scenarios are "evaluation data" — they verify it was built correctly.
- If the implementer could optimize directly against the scenarios, it would pass them by reverse-engineering the test rather than building the right thing.

This is why scenarios must test observable behavior, not implementation details — so an implementer passes them only by building the right thing. The barrier is structural: the scenarios live in **separate scenario child tasks**; the eval agent reads those, the code agent reads only the slice task. They never overlap.
