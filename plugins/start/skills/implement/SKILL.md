---
name: implement
description: Implementation entry point. Use to execute a completed specification.
user-invocable: true
argument-hint: "spec ID to implement (e.g., 002), or file path"
---

## Persona

Act as the implementation entry-point dispatcher. Resolve the spec, detect which decomposition artifacts are present, hand off to the matching execution sub-skill. You do not orchestrate implementation directly — every loop body lives in a sub-skill.

**Implementation Target**: $ARGUMENTS

## Interface

DispatchTarget {
  tier: Direct | Standard | Factory
  skill: implement-direct | implement-standard | implement-factory
  artifact: string                // path that triggered the dispatch
}

State {
  target = $ARGUMENTS
  specDirectory: string           // resolved .start/specs/NNN-name/ path
  artifacts: {
    plan?: string                 // path to plan/README.md if present
    manifest?: string             // path to manifest.md if present
    requirements?: string         // path to requirements.md
    solution?: string             // path to solution.md
  }
  dispatch: DispatchTarget
}

## Constraints

**Always:**
- Resolve the spec via Skill(start:specify-meta) before inspecting artifacts.
- Detect artifacts before dispatching — do not assume a tier.
- Pass `$ARGUMENTS` through unchanged to the sub-skill.
- When multiple decomposition artifacts coexist (both `plan/` and `manifest.md`), ask the user which to run.
- Surface what was detected so the user can confirm the dispatch target.

**Never:**
- Implement code directly — sub-skills own all execution logic.
- Run two sub-skills in parallel — pick one tier per `/start:implement` invocation.
- Modify spec artifacts during dispatch — that is each sub-skill's responsibility.

## Reference Materials

This skill is a thin dispatcher and has no reference materials of its own.
Each sub-skill owns its references, examples, and templates:

- [implement-direct](../implement-direct/SKILL.md) — phase-less orchestrator for low-complexity work
- [implement-standard](../implement-standard/SKILL.md) — linear phase loop for single-feature plans
- [implement-factory](../implement-factory/SKILL.md) — factory loop with information barriers and retry

## Workflow

### 1. Resolve Spec

Invoke `Skill(start:specify-meta)` with `$ARGUMENTS` to resolve `specDirectory`. The `--read` output reveals which artifacts exist via the `plan_dir`, `plan`, `manifest`, `units`, and `scenarios_*` keys.

If `$ARGUMENTS` is a file path or freeform brief (no spec ID), skip spec-meta resolution. Treat this as a direct-mode invocation and route to `implement-direct`.

### 2. Detect Artifacts

Inspect the resolved spec directory for decomposition artifacts:

```
manifest_md    = exists(specDirectory / "manifest.md")
plan_readme    = exists(specDirectory / "plan" / "README.md")
requirements   = exists(specDirectory / "requirements.md")
solution       = exists(specDirectory / "solution.md")
```

### 3. Dispatch

Apply the routing rules:

```
match (artifacts) {
  manifest_md exists AND NOT plan_readme   => Skill(start:implement-factory)
  plan_readme exists AND NOT manifest_md   => Skill(start:implement-standard)
  manifest_md exists AND plan_readme       => AskUserQuestion (header "Pipeline"):
                                                Factory   — run factory loop on manifest.md
                                                Standard  — run phase loop on plan/
                                                Abort     — exit without action
  none of the above
    AND (requirements OR solution exists)  => Skill(start:implement-direct)
  none of the above AND no specs           => Error: no specification artifacts found.
                                              Run /start:specify first, or pass a brief.
}
```

Before invoking the sub-skill, present a one-line dispatch summary to the user:

```
Detected: manifest.md → routing to implement-factory
```

```
Detected: plan/README.md → routing to implement-standard
```

```
Detected: only requirements.md and solution.md → routing to implement-direct
```

### 4. Hand Off

Invoke the chosen sub-skill via `Skill(start:implement-{tier})` with the same `$ARGUMENTS`.

Each sub-skill is self-contained: it reads its own artifacts, runs its own loop, and reports its own completion summary. The dispatcher does not post-process sub-skill output.

### Notes on tier mismatch

If the spec README decision log records a tier (e.g., "Decomposition tier: Standard") but the corresponding artifact is missing (e.g., no `plan/` directory), report the mismatch to the user before falling through to the next available tier. This typically indicates an interrupted `specify` run — the user should re-run `/start:specify` for that spec to complete decomposition, or explicitly choose a different tier.
