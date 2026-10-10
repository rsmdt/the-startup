# Review Workflow Contract

The workflow script you author must implement this contract exactly. The script has no filesystem
access: its agents read the PRD and SDD; the script orchestrates and grades. Every agent call except the
report passes a `schema`.

## Script rules

- `meta` is a pure literal: `name`, `description`, `phases` as plain strings — no concatenation, no
  variables, no template interpolation. A non-literal `meta` fails to register.
- Phases: `Review`, `Verify`, `Report`.
- Inputs (via `args`): `name`, `specId`, `prdPath`, `sddPath`, `reportPath` (all paths absolute).
- `MAX_FINDINGS_PER_LENS = 8`.

## Lenses

| key | title | reads PRD | focus |
|-----|-------|-----------|-------|
| `requirement-coverage` | Requirement coverage | yes | Every PRD requirement maps to exactly one SDD component that owns it. Report any requirement with no owner or more than one. |
| `ac-ears-traceability` | AC → EARS traceability | yes | Every PRD acceptance criterion has an EARS criterion in the SDD, including the edge and error states it implies. Report any AC with none, and any EARS criterion inventing behavior the PRD does not ask for. |
| `mece` | MECE (boundaries) | no | No two components share domain logic; no responsibility is unowned; no interface duplicates a consumer→provider path; each business datum has one owning entity. Report overlaps and gaps. |
| `design-risk` | Design risk | no | Three concerns, each given a deliberate pass: **security** (trust boundaries, authN/authZ, input handling, data protection, secrets), **robustness** (needless complexity, brittle abstractions, unsafe concurrency, unhandled failure modes), **compatibility** (breaking changes to existing interfaces, schemas, or contracts without a migration path). Prefix each finding with `security:`, `robustness:`, or `compatibility:`. |

## Stages

Run each lens as a pipeline — review, then verify that lens's findings as soon as it returns.

### 1. Review (one agent per lens)

Schema:

```
{ findings: [ { severity: "blocking"|"advisory" (required), detail: string (required), location: string, fix: string } ] (required) }
```

Prompt:

> You are reviewing a drafted solution design for the spec "{name}" through one lens: {title}.
> Read {the PRD at {prdPath} and} the design under review at {sddPath}.
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
> Open the design yourself: {sddPath} {and the PRD at {prdPath}}. For each finding, confirm the design
> actually has that flaw at the cited location. A finding that misreads the design, duplicates a concern
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
