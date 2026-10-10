# Workflow Contract

The workflow script you author must implement this contract exactly. The script has no filesystem or
shell access: every task.js call and file read is made by an agent; the script orchestrates and branches
on the structured results agents return. Every agent call except Synthesize passes a `schema`.

## Script rules

- `meta` is a pure literal: `name`, `description`, `phases` as plain strings — no concatenation, no
  variables, no template interpolation. A non-literal `meta` fails to register.
- Phases: `Mine`, `Cluster`, `Propose`, `Validate`, `Synthesize`, `Grade`.
- Inputs (via `args`): `name`, `home` (absolute `.start/reflect/<name>`), `taskScript` (absolute path to
  task.js), `description`, `dateRange` `{ since, until, label }`, `logRoot`, `projectRoot`, `pluginsRoot`
  (absolute path to the repo's `plugins/`), `lenses` (`[{ id, kind }]`, kind one of `friction`,
  `conversation`, `router`, `capability`, `guardrail`, `instruction`), `reportPath`
  (`home/report.md`), `maxProposals` (default 8), `maxMineRounds` (default 1).
- Every task.js command in a prompt is `node {taskScript} <verb> ... --home {home}`.

## Stages

### 1. Mine

Run every lens concurrently with the prompt from mine-agent.md (the shared block plus that lens's
section). Schema:

```
{ freshFindings: number (required), examinedSessions: string[], unexaminedRemain: boolean }
```

Track the examined sessions per lens and sum `freshFindings`. While rounds remain (`maxMineRounds`),
re-run only the lenses that returned `unexaminedRemain: true`, passing each the sessions it already
examined. A null result ends that lens.

### 2. Cluster

One cluster agent over all mine tasks. Schema:

```
{ proposalIds: string[] (required), droppedCount: number }
```

Keep at most `maxProposals` ids; log how many were cut.

### 3. Propose

One propose agent per proposal id, concurrently. Schema:

```
{ target: string (required), targetKind: string, action: "improve"|"create" (required),
  tier: "advisory"|"bounded"|"eval", hasEvalSuite: boolean, summary: string (required), diffWritten: boolean }
```

A null result is dropped and counted — it never reaches validation.

### 4. Validate

One validate agent over every drafted proposal id. Schema:

```
{ verdicts: [ { id: string, supported: boolean, reason: string } ] (required) }
```

Keep a proposal only when its verdict says `supported: true`. A proposal with no verdict was not audited
and is dropped.

### 5. Synthesize

One synthesize agent with the kept proposals; it writes `reportPath`.

### 6. Grade

Only when at least one proposal was kept: one grade agent. Schema:

```
{ dimensions: [ { name: string, score: "fail"|"partial"|"pass", note: string } ] (required),
  issues: [ { dimension: string, detail: string, fix: string } ] }
```

Compute `overall` in code: fail = 0, partial = 1, pass = 2, reported as `<points>/<2 × dimensions>`.
The report is not revised.

## Result

```
{ name, dateRange, mined, findings, opportunities, proposals, supported, dropped, reportPath,
  grade: { dimensions, issues, overall } | null, failingDimensions: string[] }
```

Log the mine rounds, the counts at each stage, and any failing grade dimension.
