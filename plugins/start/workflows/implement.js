// v9 — the loop learns to stop. Through v8 the retry loop could not tell the difference between "the
// code is wrong", "the environment is broken", "the harness went away" and "no interface can prove
// this" — every one of them arrived as a low satisfaction score, burned a retry, and re-spawned a code
// agent. v9 separates those channels and gives each its own exit:
//
//   (a) MEMORY + CONVERGENCE — each retry was amnesiac: the code agent saw six one-line symptoms and
//       nothing about what earlier attempts had already tried or already written to disk, so
//       iterations re-derived the same wrong hypothesis or spent themselves rediscovering that the
//       slice was already built. The code agent now returns a one-line `approach`, accumulated into a
//       per-slice history replayed into every later brief as "already tried, do not repeat"; the brief
//       opens with a state-on-disk inventory; and the loop halts as `stalled` once `stallLimit`
//       consecutive iterations fail to improve the score, or one symptom survives three attempts.
//       The retry budget is capacity, not a quota — it is only worth spending while the loop moves.
//
//   (b) A REAL ENVIRONMENT BARRIER — v7 told the eval agent to REPORT an environment fault and v8
//       added a migrate step, but neither was ENFORCED: the migrate agent trusted its own exit string
//       (an `npx …` command that never resolved in the agent shell still reported "none pending"),
//       nothing asserted the stack was alive before an expensive eval was spawned, and an
//       "ENVIRONMENT:" verdict was merely a failure STRING — so it scored 0 and charged the code
//       agent's retry budget. v9 replaces the migrate step with a PREFLIGHT returning a structured
//       verdict the loop branches on, requires `migrateCommand` to be directly executable (no shim
//       indirection), runs the preflight once before the whole run and again before each eval, and
//       makes an environment fault abort the slice as `blocked-environment` WITHOUT decrementing
//       maxIterations. A harness fault is not a retry.
//
//   (c) HARNESS ABORTS ARE NOT VERDICTS — a subagent can return having done no work at all (quota or
//       session limit, spawn failure). v8 consumed that emptiness as data: `verdict?.business?.
//       satisfaction ?? 0` scored an aborted eval as 0% on both surfaces, burned a retry, and handed
//       the next code agent an empty failures[]. In one run 26 of 48 subagents returned nothing but
//       "You've hit your session limit" and all three slices hit the iteration cap — a quota outage
//       recorded as a code-convergence failure. Every agent call now goes through `runAgent`, which
//       detects the zero-work shape, latches, stops spawning, charges no iteration, and surfaces the
//       run as `harness-abort`.
//
//   (d) PER-ROLE MODEL ROUTING — every agent previously inherited the parent model, so a deterministic
//       preflight and a judgment gate were billed identically despite a ~45x spread in output per run,
//       and one EVAL agent burned 1,281 turns with no ceiling. v9 pins each role explicitly (see OPTS)
//       and gives the EVAL agent a budget/stop clause, because the runtime exposes no per-agent turn
//       cap — only `model` and `effort` — so a runaway evaluation can only be bounded from the prompt.
//
//   (e) A HONEST BUG-FILING CONTRACT — the eval agent filed a bug per failed scenario with nothing to
//       reconcile against what was already filed, so every retry round re-filed the same symptom in
//       new words (one spec ended with 58 pending bug tasks under 3 slices, clustering 8/8/7/5/4/4
//       restatements of six problems); and with only pass/fail available, a missing TEST AFFORDANCE
//       was scored as a product defect and re-litigated every iteration (19 of those 58). v9 adds a
//       dedupe read before filing, a distinct `gap` task type + `notExercisable` verdict channel (out
//       of the satisfaction denominator so the loop stops burning retries on it, but blocking
//       "completed" exactly as a coverage gap does), and makes the environment channel exclusive — an
//       environment fault never becomes a child task. The tracker is the durable artifact a human
//       reads after the run; it must state real work only.
//
// v8 — schema application in the loop. When a code agent's fix is a database migration, the eval
// tests the LIVE service — whose DB does not auto-apply new migration files. So a schema-dependent
// fix (a new run_type, a new column, a constraint) sits unapplied, the eval keeps hitting the stale
// schema, and the scenario fails EVERY iteration no matter how correct the code is — the loop can
// never converge on anything requiring a migration. v8 adds a deterministic migrate step between the
// code and eval agents (a tiny agent runs `migrateCommand`, e.g. `supabase migration up`), so the DB
// the eval tests always matches the code under test. Gated on `migrateCommand` being supplied; a
// no-op when there is nothing pending. This is the same lesson as v7's env pinning: do not rely on an
// agent REMEMBERING a setup step — make it a guaranteed stage in the loop. (v9 folds this step into
// the preflight, where its result is asserted rather than trusted.)
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
// One reusable dynamic-workflow script, shipped with the `start` plugin at workflows/implement.js. The
// /start:implement SKILL is the entry point: it resolves the spec, reads and topo-sorts the tracker
// tasks, discovers the service config, and launches THIS script in place via the Workflow tool with an
// `args` object. It lives outside workflows/ deliberately — a script in workflows/ registers a slash
// command, and this one cannot be invoked directly: the runtime gives a workflow script no filesystem
// or shell access, so it can never resolve its own `tasks`/`groups`/service config. A 1-task plan and a 40-task plan run this same script —
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
//   • PREFLIGHT — deterministic harness assertion between code and eval: applies pending migrations
//                  and proves the service answers, so an eval never scores a broken environment.
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
  name: "run-implement",
  description: "Build a spec by looping each slice task through a barriered code→eval retry loop",
  phases: [
    { title: "Coverage", detail: "one audit that every feature + EARS behavior in every brief maps to a scenario (barrier-safe)" },
    { title: "Build", detail: "code agent implements each slice via TDD (reads the slice task, blind to its scenarios)" },
    { title: "Preflight", detail: "assert the harness once per run, and re-apply migrations whenever a slice writes one" },
    { title: "Evaluate", detail: "eval agent scores business + technical surfaces separately (blind to the slice + source)" },
  ],
};

// ── Input contract (args, supplied by /start:implement from the tracker) ───────────────
// {
//   specId:        string,          // for task.js --spec
//   taskScript:    string,          // abs path to specify-meta/task.js
//   threshold:     number,          // satisfaction ratio for "done" (e.g. 0.90)
//   maxIterations: number,          // retry limit per slice (e.g. 5)
//   stallLimit:    number|undefined,// consecutive non-improving iterations before halting a slice as
//                                   //   "stalled" (default 2). The retry budget is capacity, not a
//                                   //   quota — stop early when the loop stops converging.
//   servicePort:   number,          // port the (already-running) service listens on
//   startCommand:  string,          // how the LAUNCHER started the service once, before the run —
//                                   //   documented for humans; eval agents must NOT run it themselves
//   healthPath:    string,          // e.g. "/health" (eval agent polls this to VERIFY readiness only)
//   migrateCommand:string|undefined,// optional: applies pending DB migrations to the running service's
//                                   //   DB before each eval (e.g. "supabase migration up --local").
//                                   //   MUST be idempotent AND a directly-executable binary — no `npx`,
//                                   //   `pnpm dlx`, or version-manager shim: agent shells routinely
//                                   //   cannot resolve them and the step then silently no-ops.
//   requiredEnv:   string[],        // optional: env vars the scenarios depend on (e.g. feature flags).
//                                   //   Preflight asserts each is set in the env the SERVICE reads.
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
  // 3, not 5. The stall/repeat detectors below already halt a non-moving slice at iteration 3, so
  // attempts 4-5 were only ever reachable by a slice improving its score every single round — rare, and
  // when it happens the halt reports `failed` with the symptoms so a human can relaunch. The retry budget
  // is still capacity, not a quota: a slice that clears the threshold on attempt 1 spends nothing.
  maxIterations = 3,
  stallLimit = 2,
  servicePort,
  startCommand,
  healthPath = "/health",
  migrateCommand,
  requiredEnv = [],
  tasks = [],
  groups = [{ mode: "sequential", taskIds: tasks.map((t) => t.id) }],
} = parsedArgs ?? {};

// FAIL FAST on a mis-wired launch. This script is not a user entry point: it is launched by the
// /start:implement skill, which resolves the spec, topo-sorts the tracker tasks, and discovers the
// service config. The runtime gives a workflow script no filesystem or shell access, so it can never
// resolve those itself. Invoked bare, every prompt would carry "undefined" and no slice could be built.
{
  const required = { specId, taskScript, startCommand };
  const missing = Object.entries(required)
    .filter(([, v]) => typeof v !== "string" || !v.trim() || v === "undefined")
    .map(([k]) => k);
  if (!tasks.some((t) => t && t.id)) missing.push("tasks");
  if (!groups.some((g) => g && (g.taskIds ?? []).length)) missing.push("groups");
  if (missing.length) {
    throw new Error(
      `run-implement: launch args missing or unresolved: ${missing.join(", ")}. ` +
        `This workflow is launched BY the /start:implement skill, not invoked directly — the skill reads ` +
        `the spec's tracker tasks and service config first. Run /start:implement <specId> instead.`,
    );
  }
}

const tj = `node ${taskScript} `; // task.js prefix (with --spec appended per call)
const spec = `--spec ${specId}`;
const pct = (n) => `${(n * 100).toFixed(0)}%`;

// ── Per-role model routing (explicit — never inherited) ────────────────────────────────
// The roles below differ by ~45x in output per run, so inheriting one model bills a "run exactly this
// command" step at judgment-agent rates. Pin each role to the smallest model that can still be trusted
// with the verdict it returns:
//   CODE / EVAL → opus + high. EVAL is the gate OVER code; a gate must never run below the thing it
//                 judges, or a cheaper judge rubber-stamps work a stronger model produced. EVAL is the
//                 single most tempting down-tier target on raw cost — that is exactly why it is pinned.
//   COVERAGE    → sonnet + medium. It audits the SPEC (brief ↔ scenarios, a mechanical mapping into a
//                 fixed schema) and never judges the code agent, so this places no gate below its
//                 generator. NOT haiku: a false `covered: true` silently green-lights an untested
//                 acceptance surface and nothing downstream can detect it.
//   PREFLIGHT   → haiku + low. Fully deterministic: run a command, curl a URL, check a variable. Safe
//                 to down-tier because it never judges — it reports what a command actually returned,
//                 and /start:implement validates `migrateCommand` resolves before the run starts.
const OPTS = {
  coverage:  { model: "sonnet", effort: "medium" },
  code:      { model: "opus",   effort: "high" },
  eval:      { model: "opus",   effort: "high" },
  preflight: { model: "haiku",  effort: "low" },
};

// Structured verdict the eval agent must return, so the loop branches on data, not prose.
// Two acceptance surfaces scored SEPARATELY — business (requirement/feature) and technical (solution/
// EARS behavior) — each gated independently so one surface's passes can't mask the other's failures.
const SURFACE = {
  type: "object",
  properties: {
    satisfaction: {
      type: "number",
      description: "passed / (total - notExercisable) for this surface; 1.0 if no scorable scenarios remain",
    },
    // one-line observable symptoms only — these become the code agent's retry guidance
    failures: { type: "array", items: { type: "string" } },
    // Scenarios no external interface can induce: a SPEC/testability gap, never a code defect. Kept out
    // of the satisfaction denominator so the loop cannot spend every retry on something code can't fix —
    // and surfaced separately so it still blocks "completed" (see runSlice).
    notExercisable: {
      type: "array",
      items: {
        type: "object",
        properties: {
          scenario: { type: "string", description: "the scenario task id" },
          missingLever: { type: "string", description: "the external affordance that does not exist" },
        },
        required: ["scenario", "missingLever"],
      },
    },
  },
  required: ["satisfaction"],
};
const VERDICT = {
  type: "object",
  properties: {
    business: SURFACE,
    technical: SURFACE,
    // Environment faults are a SEPARATE channel, not a failure string: the loop branches on this and
    // aborts for a human instead of charging the fault to the code agent's retry budget.
    environment: {
      type: "array",
      items: { type: "string" },
      description: "one line per harness fault (service down, schema missing, flag unset); EMPTY if the environment was valid",
    },
  },
  required: ["business", "technical"],
};

// Deterministic harness assertion — cheap, runs before the expensive eval.
const PREFLIGHT = {
  type: "object",
  properties: {
    ok: { type: "boolean", description: "true iff EVERY check passed" },
    faults: { type: "array", items: { type: "string" }, description: "one line per failed check: the exact command/URL/variable and what it actually returned" },
  },
  required: ["ok"],
};

// What the code agent reports back about the attempt it just made. This is the loop's memory: the
// `approach` line is replayed verbatim into every later retry brief for this slice, so attempt N+1
// can see — and avoid — the hypotheses attempts 1..N already spent themselves on.
const ATTEMPT = {
  type: "object",
  properties: {
    approach: { type: "string", description: "one or two lines: the root cause diagnosed this iteration and the change made for it" },
    wroteCode: { type: "boolean", description: "false if the slice was found already implemented and nothing was changed" },
    wroteMigration: {
      type: "boolean",
      description:
        "true if this attempt added or changed a database migration / schema file. Say true if unsure — " +
        "a false negative means the migration is never applied and the eval runs against a stale schema.",
    },
  },
  required: ["approach"],
};

// Spec-completeness verdict from the coverage auditor: is every acceptance surface backed by a scenario?
const COVERAGE = {
  type: "object",
  properties: {
    slices: {
      type: "array",
      items: {
        type: "object",
        properties: {
          sliceId: { type: "string", description: "the slice task id this verdict is for" },
          covered: { type: "boolean", description: "true iff every feature and every EARS behavior in THIS slice maps to >= 1 scenario" },
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
        required: ["sliceId", "covered"],
      },
    },
  },
  required: ["slices"],
};

// ── Harness aborts — an agent that never ran is not a failing verdict ───────────────────
// The script sees only what agent() returns, so "did no work" has to be inferred from the result's
// shape. Every call site in this workflow is schema'd, so the test is structural and needs no
// string-matching against the harness's banner wording: a schema'd agent that produced no object, or
// omitted the required field, never answered.
class HarnessAbort extends Error {
  constructor(label, detail) {
    super(`harness abort in ${label}: ${detail}`);
    this.label = label;
    this.detail = detail;
  }
}

// Set once, then read by every subsequent runAgent — the run must not keep fanning agents into a
// harness that is already refusing to serve them.
let harnessAbort = null;

// Every agent() call goes through this. `requiredKey` names the schema field a real result must carry.
async function runAgent(prompt, opts, requiredKey) {
  if (harnessAbort) throw harnessAbort;

  const result = await agent(prompt, opts);

  if (result == null || typeof result !== "object" || result[requiredKey] === undefined) {
    const detail = `returned no "${requiredKey}"`;
    harnessAbort = new HarnessAbort(opts.label, detail);
    log(`HARNESS ABORT — ${opts.label} did no work (${detail}). This is an environment/quota failure, `
      + `NOT a verdict: no iteration is charged and nothing further is scored.`);
    throw harnessAbort;
  }
  return result;
}

// ── Prompt builders — the two halves of the information barrier + the coverage pre-flight ─

// Takes ALL slice ids. This used to run once per slice inside runSlice, but it is a static brief↔scenario
// mapping with no per-slice state and no dependency on anything the loop produces — so one auditor does
// the whole spec before any building starts, instead of S auditors interleaved with the retry loops.
function coveragePrompt(taskIds) {
  return `You are the coverage auditor for a specification. You confirm that every acceptance surface in
each slice brief is backed by at least one scenario — so nothing is left untested. You audit spec
completeness ONLY: you do not build, run, test, or edit anything.

Audit these ${taskIds.length} slice(s), one at a time:
${taskIds.map((id) => `  - ${id}`).join("\n")}

For EACH slice id above, read BOTH:
  the slice brief (business features / acceptance criteria AND technical EARS behavior):
    ${tj}read <sliceId> ${spec}
  its scenario child tasks (the holdout acceptance evidence):
    ${tj}list ${spec} --parent <sliceId> --type scenario   then read each:  ${tj}read <scenarioId> ${spec}

Map every business FEATURE / acceptance criterion and every technical EARS behavior in the brief to the
scenario(s) that would exercise it. EARS behaviors include edge and error states (WHEN/WHILE/IF … the
system SHALL …) — those are the ones most often left uncovered. Report anything NO scenario covers.

Give every slice the same scrutiny — a slice audited late is not audited less. Attribute each gap to the
slice it belongs to; a gap on the wrong slice is worse than no gap at all.

## Return
- slices: one entry per slice id above { sliceId, covered, gaps }.
  - covered: true ONLY if every feature and every EARS behavior in THAT slice maps to >= 1 scenario.
  - gaps: for each uncovered item { kind ("business" | "technical"), item, why } — name the exact unmapped
    feature or EARS behavior.`;
}

function codePrompt(taskId, iteration, failures, history) {
  const priorAttempts = history.length
    ? `\n## Attempts already made on this slice — do NOT repeat them\n` +
      history
        .map(
          (h) =>
            `- attempt ${h.iteration} (scored business ${pct(h.business)} / technical ${pct(h.technical)}): ${h.approach}\n` +
            `    symptoms that SURVIVED it: ${h.failures.length ? h.failures.join(" | ") : "(none reported)"}`,
        )
        .join("\n") +
      `\nEvery approach above has already been tried and did not clear the acceptance evidence. Form a
DIFFERENT hypothesis. If a symptom has now survived two or more distinct fixes, stop patching the same
code path — re-derive the root cause from observed behavior. And if you conclude the failure is not in
the code this slice owns (environment, a missing migration, the harness, an unowned module), say exactly
that in your returned \`approach\`: that is a useful answer and it ends the loop, another blind fix is not.\n`
    : "";

  const retry =
    iteration > 0 && failures.length
      ? `\nThe latest attempt did not satisfy the acceptance evidence ` +
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

## State on disk — start from the delta, not from zero
This working tree may ALREADY contain work for this slice: from an earlier iteration of this same loop,
or from a previous /implement run that stopped part-way. Inventory it before you plan anything:
  git status --short          # staged + untracked files an earlier attempt left behind
  git diff --stat HEAD        # everything changed since the last commit
Read whatever touches your brief and treat it as your starting point. If the slice is already built, do
NOT rewrite it and do NOT spend the iteration re-verifying it end to end: run the test suite, locate the
specific gap named in the symptoms below, fix only that, and report \`wroteCode: false\` if there was
genuinely nothing to change. Rediscovering work that is already present is the most common way an
iteration of this loop finishes having produced nothing.

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
${priorAttempts}${retry}
When complete:
1. All existing tests pass.
2. New tests cover every acceptance criterion and EARS behavior in your brief.
3. Run the full test suite and confirm it is green.

## Return
- approach: one or two lines — the root cause you diagnosed this iteration and the change you made for
  it. This is handed verbatim to the next attempt, so name the HYPOTHESIS, not the file list.
- wroteCode: false if you found the slice already implemented and changed nothing.`;
}

function preflightPrompt(applyMigrations) {
  const doMigrate = applyMigrations && Boolean(migrateCommand);
  const migrate = doMigrate
    ? `1. Apply pending database migrations — run exactly this, from the project root:
     ${migrateCommand}
   Then run it a SECOND time and confirm it reports nothing left to apply. Do NOT trust the first
   run's exit message: if the command was not found, or was proxied/shimmed and never actually
   executed, that is a FAULT — quote the stderr verbatim. Do not author migrations or edit files.
`
    : "";
  const n = (i) => String(i + (doMigrate ? 1 : 0));
  const env = requiredEnv.length
    ? `${n(2)}. Confirm these variables are set — non-empty, not commented out — in the environment the
   SERVICE actually reads (its .env / compose env, not just your shell): ${requiredEnv.join(", ")}.
`
    : "";

  return `You are the environment preflight for an acceptance-evaluation loop. Your only job is to assert
that the harness is VALID before an expensive evaluation runs. You do not build, fix, reconfigure, or
restart anything — you observe and report. Prefer direct executables: a version-manager shim (mise, nvm,
asdf) or an \`npx\`-style proxy may not resolve in this shell, and a command that never ran is a fault,
not a pass.

${migrate}${n(1)}. Confirm the service answers: curl http://localhost:${servicePort}${healthPath} (retry ~3 times, with backoff).
${env}
## Return
- ok: true ONLY if EVERY check above passed.
- faults: one line per failed check — name the exact command, URL, or variable and what it actually
  returned. Never repair the environment: a fault is reported to a human, not fixed here.`;
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
   scenarios need is absent — a table, a migration, a feature flag, a stale worker process), STOP
   IMMEDIATELY and report it in the \`environment\` array of your verdict: one line naming exactly what
   was unreachable or missing. The environment is a PRECONDITION of your evaluation, not part of what
   you score. Therefore, when \`environment\` is non-empty:
   - do NOT run the remaining scenarios, and do NOT lower a surface's satisfaction because of it;
   - do NOT file the fault as a \`--type bug\` or \`--type gap\` child task — those mean a CODE defect
     and a testability gap respectively, and an environment fault filed as either sends the next code
     agent hunting something that does not exist;
   - do NOT reconfigure, restart, re-migrate, or otherwise repair the environment to make a scenario
     pass. The run aborts and a human is told. That is the correct outcome.
   Leave \`environment\` EMPTY whenever the harness was valid — an empty array is what lets your scores
   be trusted.

## Isolate your test data (parallel-safe)
Create your own workspace/tenant (and any fixtures) for this evaluation, scoped to a unique id, so a
concurrent eval's data can never be mistaken for yours. Do not assume an empty database or reuse shared
records. Clean up only the data you created — never global state.

## Evaluation method (use the first available, in priority order)
1. E2E test automation — write and run automated end-to-end tests for the scenarios.
2. Browser automation — drive the UI/API through a real browser.
3. CLI fallback — exercise endpoints directly with curl/httpie.

## Record the outcome (bookkeeping — does NOT replace your returned verdict)
Run each scenario 3 times; a scenario passes if it passes at least 2 of 3. Classify every non-passing
scenario FIRST — the three outcomes are recorded differently:
- FAILED — the interface exists, you exercised it, the behavior was wrong. A code defect.
- NOT EXERCISABLE — the interface offers no lever to induce the condition (nothing can crash the worker
  mid-run; no way to interrupt a compile; the scale the scenario needs exceeds the eval window). This is
  a SPEC/testability gap, NOT a code defect — no amount of retrying can clear it, so it is never a failure.
- ENVIRONMENT — service unreachable or misconfigured (see above). Verdict only; file nothing.

Then, BEFORE filing anything, read what is already filed. You are one of several retry rounds against
this same slice; the same problem must never be filed twice:
  ${tj}list ${spec} --parent ${taskId} --type bug --format json
  ${tj}list ${spec} --parent ${taskId} --type gap --format json
For each failed / not-exercisable scenario, match it against those names:
- an existing task describes the SAME underlying problem (different wording still counts as the same) →
  append your evidence to it; do NOT create a task:
      ${tj}update <existingId> ${spec} --append-note "<what you observed this round>"
- nothing matches → create exactly ONE task:
      failed:           ${tj}create ${spec} --parent ${taskId} --type bug --name "<observable symptom>"
      not exercisable:  ${tj}create ${spec} --parent ${taskId} --type gap --name "<scenario>: <the external lever that does not exist>"
Finally set the slice status:
- every scenario passed:  ${tj}update ${taskId} ${spec} --status passed
- otherwise:              ${tj}update ${taskId} ${spec} --status failed

## Budget — you must terminate
You are the gate, not a second code agent. Do NOT debug, repair, or root-cause the implementation, and
do NOT re-run a scenario beyond the 3 attempts specified above hoping it flips. Once every scenario has
its 3 attempts, return the verdict you have.
Ceiling: if you exceed ~60 tool calls, STOP immediately and return your current tallies, counting every
scenario you could not resolve as a failure with a one-line symptom. A returned partial verdict is
correct behaviour — the loop retries the slice. Continuing to probe is not.

## Return the structured verdict — PER SURFACE
- business:  { satisfaction, failures[], notExercisable[] } over the BUSINESS scenarios
- technical: { satisfaction, failures[], notExercisable[] } over the TECHNICAL scenarios
- environment: [] when the harness was valid; otherwise one line per harness fault (see above)
- satisfaction (per surface): passed / (total - not-exercisable); 1.0 if no scorable scenario remains.
  A not-exercisable scenario is neither a pass nor a fail — it leaves the denominator entirely.
- failures (within each surface): for each failed scenario, ONE line stating the observable symptom
  (no internal detail)
- notExercisable (within each surface): { scenario, missingLever } — name the missing external
  affordance, not a guess at the code. Do NOT use this channel for a behavior you COULD exercise and
  that was simply wrong: that is a failure. Marking a scenario not-exercisable does not make the slice
  done — it records that the spec gave no way to prove it.`;
}

// ── The slice loop — build → preflight → evaluate → (retry with filtered summaries) ─────

// sliceId -> gaps[], filled by the single coverage audit that runs before any building starts.
const coverageGapsById = new Map();

async function runSlice(taskId) {
  // Coverage was audited up-front for every slice (barrier-safe third role): every feature + EARS behavior
  // in the brief maps to a scenario, checked BEFORE we trust any eval's pass-rate. A pass-rate can't see
  // an absent test, so an uncovered surface would otherwise sail through as "done".
  const gaps = coverageGapsById.get(taskId) ?? [];
  if (gaps.length) log(`slice ${taskId}: ${gaps.length} acceptance surface(s) have NO scenario — coverage gap (spec incomplete)`);

  let iteration = 0;
  let failures = [];
  const history = [];  // one entry per failed attempt — the loop's memory, replayed into retry briefs
  let best = -1;       // best combined (business + technical) score seen for this slice
  let stalled = 0;     // consecutive iterations that failed to beat `best`

  while (true) {
    // CODE agent: reads the slice task + the retry symptoms + what earlier attempts already tried.
    // Never the scenario tasks. Returns its approach so the next attempt inherits it.
    const attempt = await runAgent(codePrompt(taskId, iteration, failures, history.slice(-3)), {
      phase: "Build", label: `build:${taskId}`, schema: ATTEMPT, ...OPTS.code,
    }, "approach");

    // PREFLIGHT (deterministic, cheap): apply the migrations the code agent just wrote, then ASSERT the
    // harness before paying for an eval. A broken environment aborts here — it is a human's problem,
    // never a code agent's retry.
    //
    // This ran unconditionally on every iteration of every slice. It only has work to do when there are
    // migrations to apply: the run-level gate below already asserted the harness once before anything was
    // built, and if the harness breaks mid-run the eval agent's `environment` channel still catches it
    // and aborts WITHOUT spending a retry. So it now runs only when this project HAS migrations.
    //
    // Fail-safe by construction: a missing or undefined `wroteMigration` runs the preflight. Only an
    // explicit `false` from the code agent skips it, and the schema tells the agent to say true if unsure.
    // A false negative would let an eval run against a stale schema, so the default must never be "skip".
    const mayHaveMigrated = Boolean(migrateCommand) && attempt?.wroteMigration !== false;
    if (mayHaveMigrated) {
      const pre = await runAgent(preflightPrompt(true), {
        phase: "Preflight", label: `preflight:${taskId}`, schema: PREFLIGHT, ...OPTS.preflight,
      }, "ok");
      if (pre.ok === false) {
        return { id: taskId, status: "blocked-environment", environment: pre.faults ?? [], coverageGaps: gaps, iterations: iteration };
      }
    }

    // EVAL agent: reads the scenario child tasks + service config. Never the slice or the source.
    const verdict = await runAgent(evalPrompt(taskId), { phase: "Evaluate", label: `eval:${taskId}`, schema: VERDICT, ...OPTS.eval }, "business");

    // An ENVIRONMENT verdict is not a code defect. Abort for a human WITHOUT spending a retry slot —
    // `iteration` is deliberately NOT incremented here.
    const environment = verdict?.environment ?? [];
    if (environment.length) {
      return { id: taskId, status: "blocked-environment", environment, coverageGaps: gaps, iterations: iteration };
    }

    const business = verdict?.business?.satisfaction ?? 0;
    const technical = verdict?.technical?.satisfaction ?? 0;
    // Scenarios no interface can induce: excluded from satisfaction (retrying can't fix them), so they
    // must not silently pass as "done" either — same treatment as a coverage gap.
    const notExercisable = [
      ...(verdict?.business?.notExercisable ?? []),
      ...(verdict?.technical?.notExercisable ?? []),
    ];
    const surfacesPass = business >= threshold && technical >= threshold; // BOTH must clear — no blending

    if (surfacesPass) {
      // Both surfaces green, but an UNPROVEN acceptance surface — no scenario covers it (coverage gap),
      // or no interface can induce the one that does — is NOT "done".
      const status = gaps.length || notExercisable.length ? "incomplete-coverage" : "completed";
      return { id: taskId, status, business, technical, coverageGaps: gaps, notExercisable, iterations: iteration + 1 };
    }

    iteration += 1;

    // Merge both surfaces' one-line symptoms; filter before they reach the next code agent.
    failures = [...(verdict?.business?.failures ?? []), ...(verdict?.technical?.failures ?? [])].slice(0, 6);

    history.push({
      iteration,
      business,
      technical,
      failures,
      approach: attempt?.approach ?? "(not reported)",
      wroteCode: attempt?.wroteCode !== false,
    });

    // Convergence check — is this loop still moving?
    const score = business + technical;
    if (score > best) { best = score; stalled = 0; } else { stalled += 1; }
    const repeated = failures.some((f) => history.filter((h) => h.failures.includes(f)).length >= 3);

    if (iteration >= maxIterations || stalled >= stallLimit || repeated) {
      const status = iteration >= maxIterations ? "failed" : "stalled";
      if (status === "stalled") {
        log(`slice ${taskId}: HALTING at ${iteration}/${maxIterations} — ${repeated ? "the same symptom survived 3 attempts" : `${stalled} iterations with no score improvement`}. ${maxIterations - iteration} retries left unspent; this needs a human, not another attempt.`);
      }
      return { id: taskId, status, business, technical, coverageGaps: gaps, notExercisable, iterations: iteration, failures, history };
    }

    log(`slice ${taskId}: business ${pct(business)} / technical ${pct(technical)} (< ${pct(threshold)}) — retry ${iteration}/${maxIterations}`);
  }
}

// A harness abort unwinds the slice WITHOUT charging the iteration (the throw happens before
// `iteration += 1`) and without producing a satisfaction score — the slice is unjudged, not failed.
async function buildTask(taskId) {
  try {
    return await runSlice(taskId);
  } catch (e) {
    if (e instanceof HarnessAbort) {
      return { id: taskId, status: "harness-abort", abortedAt: e.label, detail: e.detail };
    }
    throw e;
  }
}

// ── Execution: run groups in order; parallelize independent slices within a group ───────

// One cheap gate before any expensive agent runs: a dead or unconfigured stack must never be
// discovered — five times over — by the eval agents themselves. Called directly rather than through
// runAgent: there is no slice to attribute an abort to yet, and a null result here is treated as
// blocked anyway, which also covers a harness that is already refusing to serve agents.
const gate = await agent(preflightPrompt(false), {
  phase: "Preflight", label: "preflight:run", schema: PREFLIGHT, ...OPTS.preflight,
});
if (!gate || gate.ok === false) {
  const faults = gate?.faults ?? ["preflight agent returned no verdict (harness unavailable?)"];
  log(`environment preflight FAILED — nothing built. ${faults.join("; ")}`);
  return { specId, threshold, status: "blocked-environment", environment: faults, results: [] };
}

const byId = new Map(tasks.map((t) => [t.id, t]));
const results = [];

// ONE coverage audit for the whole spec, before any slice is built. It reads briefs and scenarios only —
// never source — so it stays on the safe side of the code/eval barrier exactly as it did per-slice.
{
  const auditIds = groups.flatMap((g) => g.taskIds).filter((id) => byId.has(id));
  if (auditIds.length) {
    try {
      const coverage = await runAgent(
        coveragePrompt(auditIds),
        { phase: "Coverage", label: "coverage", schema: COVERAGE, ...OPTS.coverage },
        "slices",
      );
      for (const s of coverage?.slices ?? []) {
        if (s?.sliceId) coverageGapsById.set(s.sliceId, s.covered ? [] : (s.gaps ?? []));
      }
      // A slice the auditor silently omitted has no verdict. Say so rather than recording "no gaps",
      // which would read as a clean audit it never actually performed.
      const unaudited = auditIds.filter((id) => !coverageGapsById.has(id));
      if (unaudited.length) {
        log(`coverage audit returned no verdict for ${unaudited.length} slice(s): ${unaudited.join(", ")} — treated as unaudited, not as covered.`);
      }
    } catch (e) {
      // runAgent throws HarnessAbort on a quota/environment fault. Per-slice, buildTask caught this; now
      // that the audit is hoisted, it must be caught here or one fault would crash the whole run instead
      // of returning the structured result the skill reads.
      if (!(e instanceof HarnessAbort)) throw e;
      log(`coverage audit hit a harness fault (${e.detail}) — nothing built.`);
      return {
        specId,
        threshold,
        status: "harness-abort",
        results: [],
        harnessAbort: { label: e.label, detail: e.detail },
      };
    }
  }
}

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
const withGaps = results.filter((r) => (r.coverageGaps ?? []).length || (r.notExercisable ?? []).length);
const halted = results.filter((r) => r.status === "stalled");
const blocked = results.filter((r) => r.status === "blocked-environment");
const aborted = results.filter((r) => r.status === "harness-abort");

log(`${completed}/${results.length} slices satisfied`
  + `${withGaps.length ? `; ${withGaps.length} with unproven acceptance surfaces (no scenario, or no interface lever to induce one)` : ""}`
  + `${halted.length ? `; ${halted.length} halted early as not converging (${halted.map((r) => r.id).join(", ")}) — needs a human` : ""}`
  + `${blocked.length ? `; ${blocked.length} BLOCKED on the environment (no retry spent — a human must fix the harness)` : ""}`);

if (aborted.length) {
  log(`HARNESS ABORT — ${aborted.length} slice(s) stopped without a verdict (first: ${aborted[0].abortedAt}). `
    + `The harness stopped serving agents; this is NOT a code failure and these slices were NOT judged. `
    + `Their retry budget is intact — re-launch the same slice ids once the harness is available.`);
}

return { specId, threshold, results, harnessAbort: harnessAbort && { label: harnessAbort.label, detail: harnessAbort.detail } };
