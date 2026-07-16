// implement-workflow: v4
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
//   • CODE agent — builds the slice via TDD. Reads the SLICE task (its brief). NEVER reads the
//                  slice's scenario child tasks.
//   • EVAL agent — proves the slice against its SCENARIO child tasks through external interfaces.
//                  NEVER reads the slice task, the source, or spec files.
// The information barrier is "which task each agent reads": the script tells the code agent to read
// the slice task and the eval agent to read the scenario child tasks — they never overlap. The
// prompts carry belt-and-suspenders "do not read" instructions too.
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
    { title: "Build", detail: "code agent implements each slice via TDD (reads the slice task, blind to its scenarios)" },
    { title: "Evaluate", detail: "eval agent proves each slice against its scenario tasks (blind to the slice + source)" },
  ],
};

// ── Input contract (args, supplied by /start:implement from the tracker) ───────────────
// {
//   specId:        string,          // for task.js --spec
//   taskScript:    string,          // abs path to specify-meta/task.js
//   threshold:     number,          // satisfaction ratio for "done" (e.g. 0.90)
//   maxIterations: number,          // retry limit per slice (e.g. 5)
//   servicePort:   number,          // discovered service port
//   startCommand:  string,          // how the eval agent brings the service up
//   healthPath:    string,          // e.g. "/health" (eval agent polls this)
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
  tasks = [],
  groups = [{ mode: "sequential", taskIds: tasks.map((t) => t.id) }],
} = parsedArgs ?? {};

const tj = `node ${taskScript} `; // task.js prefix (with --spec appended per call)
const spec = `--spec ${specId}`;

// Structured verdict the eval agent must return, so the loop branches on data, not prose.
const VERDICT = {
  type: "object",
  properties: {
    satisfaction: { type: "number", description: "ratio 0.0–1.0 of scenarios passed" },
    passed: { type: "array", items: { type: "string" } },
    // one-line observable symptoms only — these become the code agent's retry guidance
    failures: { type: "array", items: { type: "string" } },
  },
  required: ["satisfaction"],
};

// ── Prompt builders — the two halves of the information barrier ────────────────────────

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
2. New tests cover every acceptance criterion and EARS behavior (written BEFORE implementation).
3. Run the full test suite and report results, plus how many red-green-refactor cycles you completed.`;
}

function evalPrompt(taskId) {
  return `You are evaluating whether an implementation satisfies a slice's acceptance scenarios.

Read the acceptance scenarios — they are the scenario child tasks of this slice:
  ${tj}list ${spec} --parent ${taskId} --type scenario
  then read each one:  ${tj}read <scenarioId> ${spec}

DO NOT read the slice task itself, its source code, spec files, or unit tests. Test ONLY through
the external interface — the scenarios are all you are given.

## Service lifecycle (you own it)
1. Start the service:  ${startCommand}
2. Poll http://localhost:${servicePort}${healthPath} until it responds (retry with backoff, ~5 attempts).
   If the service has no health endpoint, adapt using the project's conventions.
3. Run your evaluation against http://localhost:${servicePort}.
4. Tear the service down when you are done (kill the process you started).

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

## Return the structured verdict
- satisfaction: passed_scenarios / total_scenarios (0.0–1.0)
- passed: names of scenarios that passed
- failures: for each failed scenario, ONE line stating the observable symptom (no internal detail)`;
}

// ── The slice loop — build → evaluate → (retry with filtered summaries) ─────────────────

async function buildTask(taskId) {
  let iteration = 0;
  let failures = [];

  while (true) {
    // CODE agent: reads the slice task + one-line retry summaries. Never the scenario tasks.
    await agent(codePrompt(taskId, iteration, failures), { phase: "Build", label: `build:${taskId}` });

    // EVAL agent: reads the scenario child tasks + service config. Never the slice or the source.
    const verdict = await agent(evalPrompt(taskId), { phase: "Evaluate", label: `eval:${taskId}`, schema: VERDICT });

    const satisfaction = verdict?.satisfaction ?? 0;

    if (satisfaction >= threshold) {
      return { id: taskId, status: "completed", satisfaction, iterations: iteration + 1 };
    }

    iteration += 1;

    // Filter eval output down to one-line observable symptoms before it reaches the next code agent.
    failures = (verdict?.failures ?? []).slice(0, 6);

    if (iteration >= maxIterations) {
      return { id: taskId, status: "failed", satisfaction, iterations: iteration, failures };
    }

    log(`slice ${taskId} at ${(satisfaction * 100).toFixed(0)}% (< ${(threshold * 100).toFixed(0)}%) — retry ${iteration}/${maxIterations}`);
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
log(`${completed}/${results.length} slices satisfied`);

return { specId, threshold, results };
