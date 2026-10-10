# Loop Contract

The workflow script you author must implement this contract exactly. The script has no filesystem or
shell access: it orchestrates and branches on the structured results its agents return. Every agent call
passes a `schema`, so the loop branches on data, never on prose.

## Script rules

- `meta` is a pure literal: `name`, `description`, `phases` as plain strings — no concatenation, no
  variables, no template interpolation. A non-literal `meta` fails to register.
- Phases: `Preflight`, `Coverage`, `Build`, `Evaluate`.
- Inputs (via `args`): `specId`, `taskScript` (absolute path), `threshold`, `maxIterations`,
  `startCommand`, `servicePort`, `healthPath`, `migrateCommand` (optional), `requiredEnv` (optional),
  `groups` (ordered; each `{ parallel: boolean, sliceIds: string[] }`).
- Every task.js command in a prompt is `node {taskScript} <verb> ... --spec {specId}`.

## Harness abort

An agent that returns `null`, or an object missing its stage's key field — code `approach`, preflight
`ok`, eval `business`, coverage `slices` — did no work (quota, session limit, spawn failure). This is
never a verdict:

- latch a run-wide abort flag and spawn no further agents;
- charge no iteration and produce no score for the slice in progress;
- that slice's result is `{ id, status: "harness-abort", abortedAt: <stage label> }`.

## Stages

### 1. Run preflight (once)

Preflight agent with migrations **not** applied. Schema:

```
{ ok: boolean (required), faults: string[] }
```

`ok: false` or no result → return `{ specId, threshold, status: "blocked-environment", environment: faults, results: [] }`.

### 2. Coverage audit (once, all slices)

One coverage agent over every slice id. Schema:

```
{ slices: [ { sliceId: string, covered: boolean, gaps: [ { kind: "business"|"technical", item: string, why: string } ] } ] (required) }
```

Record gaps per slice. A slice the auditor omitted is logged as **unaudited** — never treated as covered.
A harness abort here → return `{ specId, threshold, status: "harness-abort", results: [] }`.

### 3. Build each slice

Run the groups in order. Within a parallel group, run its slices concurrently; otherwise one at a time.
Each slice runs this loop, starting with `iteration = 0`, empty `failures`, empty `history`, `best = -1`:

1. **Code agent** — prompt from code-agent.md with the last three `history` entries and the current
   `failures`. Schema:
   ```
   { approach: string (required), wroteCode: boolean, wroteMigration: boolean }
   ```
2. **Preflight** — only when `migrateCommand` is set and `wroteMigration` is not explicitly `false`.
   Preflight agent with migrations applied. `ok: false` → slice result `blocked-environment` with its
   faults; iteration **not** charged.
3. **Eval agent** — prompt from eval-agent.md. Schema:
   ```
   surface = { satisfaction: number (required), failures: string[], notExercisable: [ { scenario: string, missingLever: string } ] }
   { business: surface (required), technical: surface (required), environment: string[] }
   ```
4. **Decide:**
   - `environment` non-empty → `blocked-environment` with those faults; iteration **not** charged.
   - `business.satisfaction ≥ threshold` **and** `technical.satisfaction ≥ threshold` → `completed`, or
     `incomplete-coverage` when the slice has coverage gaps or any `notExercisable` entries. Both
     surfaces must clear the threshold independently; never blend them.
   - Otherwise: `iteration += 1`; `failures` = both surfaces' failures, first six; append
     `{ iteration, business, technical, failures, approach }` to `history`.
     - `iteration ≥ maxIterations` → stop with status `failed`.
     - `score = business + technical`. If `score > best`, set `best = score` and loop to step 1;
       otherwise stop with status `stalled` — retries remain, but the loop is not converging.

Slice result:

```
{ id, status, business, technical, iterations, coverageGaps, notExercisable, failures?, history?, environment? }
```

## Result

```
{ specId, threshold, results: [ slice result, ... ], harnessAbort?: { label, detail } }
```

A run-level `status` (`blocked-environment` or `harness-abort`) appears only when nothing was built.

Log one line per decision (slice, scores, retry n/max, or why it stopped) and a closing summary:
satisfied count, slices with unproven surfaces, stalled slices, blocked slices, aborted slices.
