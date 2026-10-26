# Review Workflow Contract

The workflow script you author must implement this contract exactly. The script has no filesystem
access: its agents read the spec document; the script orchestrates and grades. Every agent call except the
report passes a `schema`.

## Script rules

- `meta` is a pure literal: `name`, `description`, `phases` as plain strings — no concatenation, no
  variables, no template interpolation. A non-literal `meta` fails to register.
- Phases: `Review`, `Verify`, `Report`.
- Inputs (via `args`): `name`, `specId`, `specPath`, `reportPath` (both paths absolute). In `spec.md`, the
  `## Requirements` part is the PRD and the `## Solution` part is the design under review.
- `MAX_FINDINGS_PER_LENS = 8`.

## Lenses

| key | title | reads PRD | focus |
|-----|-------|-----------|-------|
| `requirement-coverage` | Requirement coverage | yes | Every PRD requirement maps to exactly one SDD component that owns it. Report any requirement with no owner or more than one. |
| `ac-ears-traceability` | AC → EARS traceability | yes | Every edge and error state a PRD acceptance criterion implies is stated, either by the criterion itself or by an EARS criterion in the SDD that refines it. Report any implied state stated by neither, and any EARS criterion inventing behavior the PRD does not ask for. An AC that already states its behavior needs no EARS restatement. |
| `mece` | MECE (boundaries) | no | No two components share domain logic; no responsibility is unowned; no interface duplicates a consumer→provider path; each business datum has one owning entity. Report overlaps and gaps. |
| `design-risk` | Design risk | yes | Three concerns, each given a deliberate pass: **security** (trust boundaries, authN/authZ, input handling, data protection, secrets), **robustness** (needless complexity, brittle abstractions, unsafe concurrency, unhandled failure modes), **compatibility** (breaking changes to existing interfaces, schemas, or contracts without a migration path). Prefix each finding with `security:`, `robustness:`, or `compatibility:`. A finding is `blocking` only when the risk defeats a PRD requirement or a stated constraint; hardening the PRD does not call for is `advisory`. |
| `simplicity` | Simplicity (YAGNI) | yes | The design solves the PRD's problem and no other. Report every component, interface, data field, ADR, or quality mechanism that no PRD requirement or acceptance criterion needs (speculative hardening, extra endpoints, generality, abstraction layers with one implementation), and every place an existing structure in the codebase or platform would replace new design. Each finding names what to remove or replace. Mark it `blocking` when the unneeded part adds a component, interface, or slice of work; `advisory` otherwise. |

## Stages

Run each lens as a pipeline — review, then verify that lens's findings as soon as it returns.

### 1. Review (one agent per lens)

Schema:

```
{ findings: [ { severity: "blocking"|"advisory" (required), detail: string (required), location: string, fix: string } ] (required) }
```

Prompt:

> You are reviewing a drafted solution design for the spec "{name}" through one lens: {title}.
> Read {the requirements — the `## Requirements` part — and} the design under review — the `## Solution`
> part — of {specPath}.
>
> ## Your lens
> {focus}
>
> Read the design in full and judge only through your lens. For every issue, cite the exact SDD location
> (section, component, interface, entity) and the concrete change that resolves it. Mark it `blocking`
> if the design is not sound to build on without the fix, `advisory` if it can ship with a note. Return
> an empty list if the design is clean on your lens.

A null result marks the lens **unreviewed** — never clean.

### 2. Verify (one agent per lens with findings)

Sort the lens's findings blocking first and verify the first `MAX_FINDINGS_PER_LENS`; the rest are
carried to the report marked **unverified**, and the cut is logged. Schema:

```
{ verdicts: [ { index: number (required), real: boolean (required), reason: string } ] (required) }
```

Prompt:

> You are the adversarial verifier for the {count} finding(s) below on the solution design "{name}"
> (lens: {title}). Assume each finding is wrong until the design forces the opposite conclusion.
>
> {numbered list: `i. [severity] detail — location`}
>
> Open the design yourself: the `## Solution` part {and the `## Requirements` part} of {specPath}. For each
> finding, confirm the design actually has that flaw at the cited location. A finding that misreads the design, duplicates a concern
> handled elsewhere, or describes a non-issue is not real. When in doubt, rule not real.
>
> Judge each finding on its own evidence; a run of identical verdicts means you stopped opening the
> design. Return one verdict per finding, `index` matching the numbering above.

A finding with `real: true` is confirmed. A finding with no verdict, or a null verifier result, is
carried as **unverified** — never dropped silently.

### 3. Grade (in code — no agent)

Per lens:
- **unreviewed** lens → `fail`, note "lens did not run".
- no confirmed and no unverified findings → `pass`.
- any confirmed `blocking` finding, or any confirmed finding on `requirement-coverage` or
  `ac-ears-traceability` → `fail`.
- otherwise (only advisory confirmed findings, or only unverified ones) → `partial`.

`overall` = points (fail 0, partial 1, pass 2) as `<points>/<2 × lenses>`. `issues` lists every confirmed
and unverified finding with its lens, detail, location, and fix.

### 4. Report (one agent)

> Write the solution-design review report for "{name}" to {reportPath}:
> `# Solution Design Review — {name} ({specId})`, a **Grade** section (overall and each lens's score and
> note), **Findings by dimension** (one section per lens, each finding as `[severity] detail — location →
> fix`, unverified ones marked `(unverified)`, `(clean)` when empty), and **How to act**: the design is
> human-owned; resolve each blocking finding with the user (confirming any architecture decision), then
> re-run the review. Advisory findings may ship with a note.

## Result

```
{ name, specId, reportPath, confirmed, unverified, grade: { dimensions, issues, overall },
  failingDimensions: string[], reviewReady: boolean }
```

`reviewReady` is true only when no lens scored `fail`.
