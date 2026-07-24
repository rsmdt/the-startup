// implement-workflow: v8
//
// v8 — schema application in the loop. When a code agent's fix is a database migration, the eval
// tests the LIVE service — whose DB does not auto-apply new migration files. So a schema-dependent
// fix (a new run_type, a new column, a constraint) sits unapplied, the eval keeps hitting the stale
// schema, and the scenario fails EVERY iteration no matter how correct the code is — the loop can
// never converge on anything requiring a migration. v8 adds a deterministic migrate step between the
// code and eval agents (a tiny agent runs `migrateCommand`, e.g. `supabase migration up`), so the DB
// the eval tests always matches the code under test. Gated on `migrateCommand` being supplied; a
// no-op when there is nothing pending. This is the same lesson as v7's env pinning: do not rely on an
// agent REMEMBERING a setup step — make it a guaranteed stage in the loop.
//
// v7 — eval harness determinism. The service is now PINNED: the launcher brings it up once,
// already configured, and every eval agent evaluates against it read-only. Eval agents must not
// start, restart, stop, or reconfigure the service (no editing .env / compose / feature flags),
// and must isolate their own test data (own workspace/tenant) so parallel evals cannot collide.
// This removes two non-determinism sources that let environment state masquerade as code failure:
//   (1) agents mutating a shared .env per iteration and restoring it — a step one agent could skip,
//       silently disabling a feature and zeroing an otherwise-passing slice;
//   (2) parallel eval agents each owning the lifecycle of ONE shared stack on ONE port — one
//       agent's teardown killing the service another is mid-evaluation against.
// A service that is down or misconfigured is now REPORTED as an environment failure, never
// "fixed" by the eval agent — so a harness problem can never be scored as a code defect.
//
// The Agentic Startup — single implementation loop (unified task model).
//
// One reusable dynamic-workflow script. /start:implement installs this file to
// <project>/.claude/workflows/implement.js and launches it via the Workflow tool, passing the
// spec's implementation tasks as `args`. A 1-task plan and a 40-task plan run this same script —
// scale is emergent from the number of tasks.
//
// Everything is a task (see specify-meta/task.js). For a spec's implement phase the tracker holds:
//   • a root phase task,
//   • one SLICE task per unit of work — its body is the code brief (features + EARS + gates),
//   • SCENARIO child tasks under each slice — their bodies are the holdout acceptance evidence.
// specify-plan authors these; /start:implement reads them and passes the slice task ids here.
//
// Each slice task is driven through a barriered code→eval retry loop:
//   • COVERAGE auditor — a spec-completeness pre-flight: reads BOTH the slice brief and its scenarios
//                  and reports any feature / EARS behavior that NO scenario covers. A pass-rate can't
//                  see an absent test, so this is what stops acceptance surfaces being silently missed.
//   • CODE agent — builds the slice via TDD. Reads the SLICE task (its brief). NEVER reads the
//                  slice's scenario child tasks.
//   • EVAL agent — proves the slice against its SCENARIO child tasks through external interfaces.
//                  NEVER reads the slice task, the source, or spec files. Scores TWO acceptance
//                  surfaces separately — business (requirement/feature) and technical (solution/EARS) —
//                  and a slice is done only when BOTH clear the threshold (a flat blended ratio would
//                  let feature passes mask EARS-behavior failures).
// The information barrier is "which task each agent reads": the code agent reads the slice, the eval
// agent reads the scenarios — they never overlap (belt-and-suspenders "do not read" instructions too).
// The coverage auditor reads both, but it neither builds nor scores the implementation and its output
// toward the build is only "surface X has no test" — so the barrier survives.
//
// Independent slices run concurrently via parallel(); dependent slices run in order (the skill
// topo-sorts them into `groups`).
//
// Runtime note: the script has no shell/filesystem access — only agents do. So every task.js call
// (read/list/update/create) is made BY an agent; the script only orchestrates and branches on the
// structured verdicts agents return.

export const meta = {
  name: "implement",
  description: "Build a spec by looping each slice task through a barriered code→eval retry loop",
  phases: [
    { title: "Coverage", detail: "audit that every feature + EARS behavior in the brief maps to a scenario (barrier-safe)" },
    { title: "Build", detail: "code agent implements each slice via TDD (reads the slice task, blind to its scenarios)" },
    { title: "Evaluate", detail: "eval agent scores business + technical surfaces separately (blind to the slice + source)" },
  ],
};

// ── Input contract (args, supplied by /start:implement from the tracker) ───────────────
// {
//   specId:        string,          // for task.js --spec
//   taskScript:    string,          // abs path to specify-meta/task.js
//   threshold:     number,          // satisfaction ratio for "done" (e.g. 0.90)
//   maxIterations: number,          // retry limit per slice (e.g. 5)
//   servicePort:   number,          // port the (already-running) service listens on
//   startCommand:  string,          // how the LAUNCHER started the service once, before the run —
//                                   //   documented for humans; eval agents must NOT run it themselves
//   healthPath:    string,          // e.g. "/health" (eval agent polls this to VERIFY readiness only)
//   migrateCommand:string|undefined,// optional: applies pending DB migrations to the running service's
//                                   //   DB before each eval (e.g. "npx supabase migration up --local").
//                                   //   Omit for projects with no migration step. Must be idempotent.
//   tasks:  [ { id: string, dependsOn: string[] } ],           // slice task ids + deps (tracker ids)
//   groups: [ { mode: "parallel" | "sequential", taskIds: string[] } ],  // topo-sorted
// }

// Background workflow launches may serialize `args`; on some runtimes it is revived as a
// JSON string rather than an object. Parse defensively so destructuring never yields all-undefined.
const parsedArgs = typeof args === "string" ? JSON.parse(args) : args;

const {
  specId,
  taskScript,
  threshold = 0.9,
  maxIterations = 5,
  servicePort,
  startCommand,
  healthPath = "/health",
  migrateCommand,
  tasks = [],
  groups = [{ mode: "sequential", taskIds: tasks.map((t) => t.id) }],
} = parsedArgs ?? {};

const tj = `node ${taskScript} `; // task.js prefix (with --spec appended per call)
const spec = `--spec ${specId}`;

// Structured verdict the eval agent must return, so the loop branches on data, not prose.
// Two acceptance surfaces scored SEPARATELY — business (requirement/feature) and technical (solution/
// EARS behavior) — each gated independently so one surface's passes can't mask the other's failures.
const SURFACE = {
  type: "object",
  properties: {
    satisfaction: { type: "number", description: "passed/total for this surface; 1.0 if it has no scenarios" },
    // one-line observable symptoms only — these become the code agent's retry guidance
    failures: { type: "array", items: { type: "string" } },
  },
  required: ["satisfaction"],
};
const VERDICT = {
  type: "object",
  properties: { business: SURFACE, technical: SURFACE },
  required: ["business", "technical"],
};

// Spec-completeness verdict from the coverage auditor: is every acceptance surface backed by a scenario?
const COVERAGE = {
  type: "object",
  properties: {
    covered: { type: "boolean", description: "true iff every feature and every EARS behavior maps to >= 1 scenario" },
    gaps: {
      type: "array",
      items: {
        type: "object",
        properties: {
          kind: { type: "string", description: "business | technical" },
          item: { type: "string", description: "the exact feature or EARS behavior no scenario covers" },
          why: { type: "string" },
        },
        required: ["kind", "item"],
      },
    },
  },
  required: ["covered"],
};

// ── Prompt builders — the two halves of the information barrier + the coverage pre-flight ─

function coveragePrompt(taskId) {
  return `You are the coverage auditor for one slice of a specification. You confirm that every
acceptance surface in the slice brief is backed by at least one scenario — so nothing is left untested.
You audit spec completeness ONLY: you do not build, run, test, or edit anything.

Read BOTH:
  the slice brief (business features / acceptance criteria AND technical EARS behavior):
    ${tj}read ${taskId} ${spec}
  its scenario child tasks (the holdout acceptance evidence):
    ${tj}list ${spec} --parent ${taskId} --type scenario   then read each:  ${tj}read <scenarioId> ${spec}

Map every business FEATURE / acceptance criterion and every technical EARS behavior in the brief to the
scenario(s) that would exercise it. EARS behaviors include edge and error states (WHEN/WHILE/IF … the
system SHALL …) — those are the ones most often left uncovered. Report anything NO scenario covers.

## Return
- covered: true ONLY if every feature and every EARS behavior maps to >= 1 scenario.
- gaps: for each uncovered item { kind ("business" | "technical"), item, why } — name the exact unmapped
  feature or EARS behavior.`;
}

function codePrompt(taskId, iteration, failures) {
  const retry =
    iteration > 0 && failures.length
      ? `\nA previous attempt did not satisfy the acceptance evidence ` +
        `(attempt ${iteration + 1} of ${maxIterations}). Observed symptoms:\n` +
        failures.map((f) => `- ${f}`).join("\n") +
        `\nInvestigate the codebase to find and fix the root causes.\n`
      : "";

  return `You are implementing one slice of a specification using Test-Driven Development (TDD).

Read your task — it is your code brief (features / acceptance criteria, EARS behavior, quality gates):
  ${tj}read ${taskId} ${spec}
Mark it in progress before you start (bookkeeping only):
  ${tj}update ${taskId} ${spec} --status building

Codebase context: read AGENTS.md (or CLAUDE.md) for project orientation and conventions.

DO NOT read this slice's scenario child tasks (type: scenario), nor any files under a scenarios/
directory. Those hold the holdout acceptance evidence and must not influence how you build.

## File naming
Name test files after the feature or component under test — never after an internal task id.
  ✅ tests/features/orders/test_order_repository.py
  ❌ tests/features/orders/test_${taskId}.py

## TDD process (mandatory)
For EACH acceptance criterion and EARS behavior in your brief, follow red-green-refactor:
1. RED: write a failing test for the behavior. Run the suite — confirm it fails.
2. GREEN: write the minimal code to make it pass. Run the suite — confirm it passes.
3. REFACTOR: improve the code while keeping every test green.
Do not write implementation code without a failing test first.
${retry}
When complete:
1. All existing tests pass.
2. New tests cover every acceptance criterion and EARS behavior in your brief.
3. Run the full test suite and confirm it is green.`;
}

function migratePrompt() {
  return `Apply any pending database migrations to the running service's database so the acceptance
evaluation tests the CURRENT schema. Run exactly this command from the project root:

  ${migrateCommand}

This is idempotent — if nothing is pending it is a harmless no-op. Do not write code, edit files, or
author new migrations; only apply what already exists on disk. If the command fails, report the error
verbatim as your final message. Otherwise report the migrations applied (or "none pending").`;
}

function evalPrompt(taskId) {
  return `You are evaluating whether an implementation satisfies a slice's acceptance scenarios.

Read the acceptance scenarios — they are the scenario child tasks of this slice:
  ${tj}list ${spec} --parent ${taskId} --type scenario
  then read each one:  ${tj}read <scenarioId> ${spec}

DO NOT read the slice task itself, its source code, spec files, or unit tests. Test ONLY through
the external interface — the scenarios are all you are given.

## Two acceptance surfaces — score them SEPARATELY
Each scenario proves either a BUSINESS check (a requirement / feature acceptance — does it do what the
user needs?) or a TECHNICAL check (a solution EARS behavior — the specified response to a condition or
event, including edge and error states). Determine which from the scenario itself, or from its declared
kind if it carries one. Keep the two tallies apart: a slice is not done until BOTH surfaces pass, so
feature passes must never average away EARS-behavior failures.

## Service (already running and configured — do NOT touch its lifecycle or config)
The service was started ONCE by the launcher (\`${startCommand}\`) and is already configured. Other
eval agents may be evaluating against this same instance concurrently. Therefore:
1. VERIFY readiness only: poll http://localhost:${servicePort}${healthPath} (retry with backoff, ~5
   attempts). If it has no health endpoint, adapt using the project's conventions.
2. Do NOT start, restart, or stop the service, and do NOT kill any process — you did not start it and
   others depend on it. Do NOT edit .env, compose files, feature flags, or any configuration; the
   environment is pinned. Evaluate against http://localhost:${servicePort} exactly as it stands.
3. If the service is unreachable after your retries, or is plainly misconfigured (a core capability the
   scenarios need is absent), STOP and return a verdict that records this as an ENVIRONMENT failure:
   set the affected surface's satisfaction to 0 and put a failure line beginning "ENVIRONMENT:" naming
   what was unreachable or missing. Never reconfigure the service to make a scenario pass — a harness
   problem must be reported as such, not scored as a code defect and not silently worked around.

## Isolate your test data (parallel-safe)
Create your own workspace/tenant (and any fixtures) for this evaluation, scoped to a unique id, so a
concurrent eval's data can never be mistaken for yours. Do not assume an empty database or reuse shared
records. Clean up only the data you created — never global state.

## Evaluation method (use the first available, in priority order)
1. E2E test automation — write and run automated end-to-end tests for the scenarios.
2. Browser automation — drive the UI/API through a real browser.
3. CLI fallback — exercise endpoints directly with curl/httpie.

## Record the outcome (bookkeeping — does NOT replace your returned verdict)
Run each scenario 3 times; a scenario passes if it passes at least 2 of 3. Then:
- if every scenario passed:  ${tj}update ${taskId} ${spec} --status passed
- otherwise:                 ${tj}update ${taskId} ${spec} --status failed
  and, for EACH failed scenario, file a bug as a child task:
                             ${tj}create ${spec} --parent ${taskId} --type bug --name "<observable symptom>"

## Return the structured verdict — PER SURFACE
- business:  { satisfaction (passed/total of BUSINESS scenarios; 1.0 if there are none), failures[] }
- technical: { satisfaction (passed/total of TECHNICAL scenarios; 1.0 if there are none), failures[] }
- failures (within each surface): for each failed scenario, ONE line stating the observable symptom
  (no internal detail)`;
}

// ── The slice loop — build → evaluate → (retry with filtered summaries) ─────────────────

async function buildTask(taskId) {
  // Coverage pre-flight (barrier-safe third role): confirm every feature + EARS behavior in the brief
  // maps to a scenario BEFORE we trust the eval's pass-rate. A pass-rate can't see an absent test, so
  // an uncovered surface would otherwise sail through as "done".
  const coverage = await agent(coveragePrompt(taskId), { phase: "Coverage", label: `coverage:${taskId}`, schema: COVERAGE });
  const gaps = coverage?.covered ? [] : (coverage?.gaps ?? []);
  if (gaps.length) log(`slice ${taskId}: ${gaps.length} acceptance surface(s) have NO scenario — coverage gap (spec incomplete)`);

  let iteration = 0;
  let failures = [];

  while (true) {
    // CODE agent: reads the slice task + one-line retry summaries. Never the scenario tasks.
    await agent(codePrompt(taskId, iteration, failures), { phase: "Build", label: `build:${taskId}` });

    // MIGRATE step (deterministic): apply any migrations the code agent just wrote to the running
    // service's DB, so the eval tests the current schema — not a stale one that fails every iteration.
    if (migrateCommand) {
      await agent(migratePrompt(), { phase: "Build", label: `migrate:${taskId}` });
    }

    // EVAL agent: reads the scenario child tasks + service config. Never the slice or the source.
    const verdict = await agent(evalPrompt(taskId), { phase: "Evaluate", label: `eval:${taskId}`, schema: VERDICT });

    const business = verdict?.business?.satisfaction ?? 0;
    const technical = verdict?.technical?.satisfaction ?? 0;
    const surfacesPass = business >= threshold && technical >= threshold; // BOTH must clear — no blending

    if (surfacesPass) {
      // Both surfaces green, but a coverage gap means an untested surface — that is NOT "done".
      const status = gaps.length ? "incomplete-coverage" : "completed";
      return { id: taskId, status, business, technical, coverageGaps: gaps, iterations: iteration + 1 };
    }

    iteration += 1;

    // Merge both surfaces' one-line symptoms; filter before they reach the next code agent.
    failures = [...(verdict?.business?.failures ?? []), ...(verdict?.technical?.failures ?? [])].slice(0, 6);

    if (iteration >= maxIterations) {
      return { id: taskId, status: "failed", business, technical, coverageGaps: gaps, iterations: iteration, failures };
    }

    log(`slice ${taskId}: business ${(business * 100).toFixed(0)}% / technical ${(technical * 100).toFixed(0)}% (< ${(threshold * 100).toFixed(0)}%) — retry ${iteration}/${maxIterations}`);
  }
}

// ── Execution: run groups in order; parallelize independent slices within a group ───────

const byId = new Map(tasks.map((t) => [t.id, t]));
const results = [];

for (const group of groups) {
  const ids = group.taskIds.filter((id) => byId.has(id));
  if (!ids.length) continue;

  if (group.mode === "parallel") {
    const groupResults = await parallel(ids.map((id) => () => buildTask(id)));
    results.push(...groupResults.filter(Boolean));
  } else {
    for (const id of ids) results.push(await buildTask(id));
  }
}

const completed = results.filter((r) => r.status === "completed").length;
const withGaps = results.filter((r) => (r.coverageGaps ?? []).length);
log(`${completed}/${results.length} slices satisfied${withGaps.length ? `; ${withGaps.length} with coverage gaps (untested acceptance surfaces)` : ""}`);

return { specId, threshold, results };
