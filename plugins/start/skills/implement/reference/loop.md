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
  `phases` (ordered; each `{ phaseId: string, groups: [ { parallel: boolean, sliceIds: string[] } ] }`,
  groups ordered within the phase).
- Every task.js command in a prompt is `node {taskScript} <verb> ... --spec {specId}`.

## Harness abort

An agent that returns `null`, or an object missing its stage's key field — code `approach`, preflight
and service `ok`, eval `slices`, coverage `slices` — did no work (quota, session limit, spawn failure). This is
never a verdict:

- latch a run-wide abort flag and spawn no further agents;
- charge no iteration and produce no score for the slice in progress;
- every slice of the phase in progress that is not yet decided gets the result
  `{ id, status: "harness-abort", abortedAt: <stage label> }`.

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

### 3. Build and evaluate each phase

Run the phases in order. Each slice keeps its own state: `iteration = 0`, empty `failures`, empty
`history`, `best = -1`. A phase repeats these rounds until every one of its slices is decided:

1. **Code agents** — for each undecided slice, a code agent with the prompt from code-agent.md, the last
   three `history` entries, and the current `failures`. Run the phase's groups in order; within a
   parallel group, run its slices concurrently. Schema:
   ```
   { approach: string (required), wroteCode: boolean, wroteMigration: boolean }
   ```
2. **Preflight** — only when `migrateCommand` is set and some code agent's `wroteMigration` is not
   explicitly `false`. Preflight agent with migrations applied. `ok: false` → every undecided slice is
   `blocked-environment` with its faults; iteration **not** charged.
3. **Service** — one service agent (prompt from service-agent.md) restarts the service from the current
   code. `ok: false` → every undecided slice is `blocked-environment` with its faults; iteration **not**
   charged.
4. **Eval agent** — one eval agent for the phase (prompt from eval-agent.md) over every undecided slice.
   Schema:
   ```
   surface = { satisfaction: number (required), failures: string[], notExercisable: [ { scenario: string, missingLever: string } ] }
   { slices: [ { sliceId: string (required), business: surface (required), technical: surface (required) } ] (required), environment: string[] }
   ```
5. **Decide, per slice:**
   - `environment` non-empty → every undecided slice is `blocked-environment` with those faults;
     iteration **not** charged.
   - A slice missing from `slices` was not judged → `harness-abort` for that slice.
   - `business.satisfaction ≥ threshold` **and** `technical.satisfaction ≥ threshold` → `completed`, or
     `incomplete-coverage` when the slice has coverage gaps or any `notExercisable` entries. Both
     surfaces must clear the threshold independently; never blend them.
   - Otherwise: `iteration += 1`; `failures` = both surfaces' failures, first six; append
     `{ iteration, business, technical, failures, approach }` to `history`.
     - `iteration ≥ maxIterations` → stop with status `failed`.
     - `score = business + technical`. If `score > best`, set `best = score` and the slice stays
       undecided for the next round; otherwise stop with status `stalled` — retries remain, but the
       loop is not converging.

The next round rebuilds only the undecided slices; slices already decided are not rebuilt or
re-evaluated.

Slice result:

```
{ id, status, business, technical, iterations, coverageGaps, notExercisable, failures?, history?, environment? }
```

## Result

```
{ specId, threshold, results: [ slice result, ... ], harnessAbort?: { label, detail } }
```

A run-level `status` (`blocked-environment` or `harness-abort`) appears only when nothing was built.

Log one line per decision (phase, slice, scores, retry n/max, or why it stopped) and a closing summary:
satisfied count, slices with unproven surfaces, stalled slices, blocked slices, aborted slices.
