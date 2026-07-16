---
name: implement
description: Implementation entry point. Use to execute a completed specification. Reads the spec (requirements, solution, plan) and selects an execution strategy — Direct, Iterative, or Decompose — then hands off to the matching engine. Also runs legacy specs (plan/ or manifest.md) by back-compatible detection.
user-invocable: true
argument-hint: "spec ID to implement (e.g., 002), or file path"
---

## Persona

Act as the implementation entry point. You read the specification, decide **how** to execute it (Direct, Iterative, or Decompose), record that decision, and hand off to the matching engine skill. The specification tells you the target state and what "done" looks like; you choose the execution shape. You do not implement code directly — every loop body lives in an engine skill.

**Implementation Target**: $ARGUMENTS

## Interface

Strategy = Direct | Iterative | Decompose

DispatchTarget {
  strategy: Strategy
  skill: implement-direct | implement-incremental | implement-factory   // incremental=Iterative, factory=Decompose
  reason: string                  // why this strategy (selector rationale or legacy detection)
}

State {
  target = $ARGUMENTS
  specDirectory: string           // resolved .start/specs/NNN-name/ path
  artifacts: {
    requirements?: string         // requirements.md(x)
    solution?: string             // solution.md(x)
    plan?: string                 // plan.mdx (new) — completion contract
    scenarios?: string            // scenarios/phase-N/ (new) — holdout validation scenarios
    legacyPlanDir?: string        // plan/README.md (legacy incremental)
    legacyManifest?: string       // manifest.md (legacy factory)
  }
  dispatch: DispatchTarget
}

## Constraints

**Always:**
- Resolve the spec via the specify-meta skill before inspecting artifacts.
- Detect artifacts before selecting — legacy artifacts route by back-compat; new specs route by the strategy selector.
- Run the strategy selector (reference/strategy-selector.md) on requirements + solution + plan when no legacy decomposition artifact exists.
- Record the chosen strategy in `.impl/strategy.md` and the spec README decisions log before handing off.
- Surface what was detected and the recommended strategy so the user can confirm or override.
- Pass `$ARGUMENTS` and the chosen strategy through to the engine skill.

**Never:**
- Implement code directly — engine skills own all execution logic.
- Run two engines in parallel — pick one strategy per `/start:implement` invocation.
- Author or modify specification artifacts (requirements/solution/plan/scenarios) — those are the spec's, and the engine writes only working-tree code plus `.impl/`.
- Choose a strategy without showing the user the rationale.

## Reference Materials

- [Strategy Selector](reference/strategy-selector.md) — Heuristic for choosing Direct, Iterative, or Decompose

Each engine skill owns its own loop, references, and templates:
- [implement-direct](../implement-direct/SKILL.md) — Direct: small, phase-less work
- [implement-incremental](../implement-incremental/SKILL.md) — Iterative: linear phase loop (generated at exec time)
- [implement-factory](../implement-factory/SKILL.md) — Decompose: parallel units + holdout-scenario retry loop (generated at exec time)

## Workflow

### 1. Resolve Spec

Use the specify-meta skill with `$ARGUMENTS` to resolve `specDirectory`.

If `$ARGUMENTS` is a file path or freeform brief (no spec ID), skip spec-meta resolution, set strategy = Direct, and route to implement-direct.

### 2. Detect Artifacts

Inspect the resolved spec directory:

```
# New layout (completion contract)
plan_mdx       = exists(specDirectory / "plan.mdx")
scenarios_dir  = exists(specDirectory / "scenarios")
requirements   = exists(specDirectory / "requirements.md") or "requirements.mdx"
solution       = exists(specDirectory / "solution.md")     or "solution.mdx"

# Legacy layout (pre-built decomposition)
legacy_plan    = exists(specDirectory / "plan" / "README.md")
legacy_manifest= exists(specDirectory / "manifest.md")
```

### 3. Select Strategy

```
match (artifacts) {
  legacy_manifest exists AND NOT legacy_plan  => Decompose (back-compat: run factory loop on manifest.md)
  legacy_plan exists AND NOT legacy_manifest  => Iterative (back-compat: run phase loop on plan/)
  legacy_manifest AND legacy_plan             => ask the user under header "Strategy":
                                                   Decompose — run factory loop on manifest.md
                                                   Iterative — run phase loop on plan/
                                                   Abort     — exit without action

  no legacy artifacts                         => run the Strategy Selector:
       Read reference/strategy-selector.md and apply it to plan.mdx (the phase baseline),
       with solution + requirements for context.
       Surface the recommendation with its signals (phase count, AC count, component
       count, change type, parallel markers).
       Ask the user to choose under header "Strategy" — Direct, Iterative, or Decompose —
       with the recommendation highlighted. The user may override freely.
       The choice decides how to consume the plan's phases: Direct ignores them and builds
       to the gates + scenarios, Iterative follows them, Decompose splits them into units.

  no specs at all                             => Error: no specification found.
                                                 Run the specify skill first, or pass a brief.
}
```

When a new-layout spec has a `plan.mdx`, it already carries a **phase baseline** — vertical-slice phases, each correlating features (ACs), EARS behavior, quality gates, and validation scenarios. The plan's `scenarios/phase-N/` are the holdout set; all engines validate against them and never re-author them at implement time. Iterative follows the phases, Decompose splits them into finer units, Direct builds straight to the gates + scenarios.

### 4. Record the Decision

Write `.impl/strategy.md` capturing: chosen strategy, the selector rationale (or legacy-detection note), and the source artifacts. Append the decision to the spec README decisions log per reference/strategy-selector.md.

The `.impl/` directory holds all implementation-time artifacts (strategy, generated phases/units, manifest) — kept separate from the specification, and resumable across sessions.

### 5. Hand Off

Present a one-line dispatch summary, then invoke the engine:

```
Strategy: Iterative (2 phases, 2 components) → implement-incremental
Strategy: Decompose (4 components, parallel flagged) → implement-factory
Strategy: Direct (fix, existing surface) → implement-direct
Detected legacy manifest.md → Decompose → implement-factory
```

Use the chosen engine skill with the same `$ARGUMENTS` and the selected strategy. Each engine is self-contained: for new-layout specs it generates its own decomposition into `.impl/` from requirements + solution + plan; for legacy specs it reads the pre-built `plan/` or `manifest.md`. The engine reads its artifacts, runs its loop, and reports its own completion summary.
