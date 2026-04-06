# Scenario Guide

How to write effective holdout scenarios for the Dark Factory evaluation loop.

---

## What Makes a Good Scenario

A scenario is a plain-English description of expected system behavior, tested through external interfaces. The evaluation agent reads this and exercises the running system.

**Good scenarios**:
- Test observable behavior through the API, UI, or CLI
- Reference actual endpoints, data shapes, and conventions from the codebase
- Are specific enough that pass/fail is unambiguous
- Each test one thing

**Bad scenarios**:
- Reference internal implementation details (class names, function signatures, file paths)
- Are vague ("the system should work correctly")
- Test multiple unrelated behaviors in one scenario
- Duplicate what unit tests already cover

## Grounding in the Codebase

Before writing scenarios, read:
1. **AGENTS.md** — understand the project's architecture, conventions, and patterns
2. **Existing test files** — understand the test framework and existing test patterns
3. **API routes / endpoints** — reference actual paths, not hypothetical ones
4. **Data models** — reference actual field names and types

A scenario like "POST to /api/v1/sql/validate" is grounded. A scenario like "POST to the validation endpoint" is not.

## Priority Levels

| Priority | Description | When to Use |
|----------|-------------|-------------|
| P0 | Critical path | The primary happy path. If this fails, the feature is broken. |
| P1 | Important | Error handling, edge cases that affect users. |
| P2 | Edge case | Boundary conditions, unusual inputs, rare paths. |

Every unit should have at least one P0 scenario. P1 and P2 scenarios increase confidence.

## Scenario Count Guidelines

- **Minimum**: 1 scenario per unit (P0 happy path)
- **Typical**: 3-5 scenarios per unit (P0 + P1 error cases)
- **Thorough**: 5-10 scenarios per unit (P0 + P1 + P2 edge cases)

More scenarios = higher confidence but longer evaluation time. Start with P0 + P1, add P2 if the factory loop has iterations to spare.

## Information Barrier Awareness

Scenarios are the **holdout evaluation set** — invisible to the code agent. They should test behaviors that a well-implemented unit would pass, but that a poorly-implemented unit would fail. Think of them as the test/validation split in ML:

- The unit spec is "training data" — tells the code agent what to build
- Scenarios are "evaluation data" — verify the code agent built it correctly
- If the code agent could see the scenarios, it would optimize for passing them rather than building correct software

This is why scenarios must test observable behavior, not implementation details. The code agent should pass scenarios by building the right thing, not by reverse-engineering the test.
