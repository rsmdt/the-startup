// implement-workflow: v1
//
// The Agentic Startup — single implementation loop.
//
// One reusable dynamic-workflow script. It is NOT regenerated per spec: /start:implement
// installs this file to <project>/.claude/workflows/implement.js and launches it via the
// Workflow tool, passing the spec's plan as `args`. A 1-BEAN plan and a 40-BEAN plan run
// this same script — scale is emergent from the number of BEANs, not from a strategy choice.
//
// A BEAN (Bundled EARS + Acceptance Node) is one vertical slice of the plan. Each BEAN is
// driven through a barriered code→eval retry loop:
//   • CODE agent  — builds the BEAN via TDD. Sees the BEAN spec + filtered retry summaries.
//                   NEVER sees the holdout scenarios.
//   • EVAL agent  — proves the BEAN against its holdout scenarios through external interfaces.
//                   NEVER sees the source, the BEAN spec, or the code agent's output.
// The information barrier is STRUCTURAL: this script threads BEAN-spec content only into the
// code prompt and scenario content only into the eval prompt. The two never cross in-script.
// The prompts additionally carry belt-and-suspenders "do not read" instructions.
//
// Independent BEANs (no unmet dependency) run concurrently via parallel(); dependent BEANs
// run in order. The orchestrator (/start:implement) topologically sorts them into groups and
// passes that ordering in `args.groups`.
//
// Runtime note: the script itself has no filesystem/shell access — only agents do. So each
// EVAL agent owns the full service lifecycle (start → health-check → run 3× → tear down).

export const meta = {
  name: "implement",
  description: "Build a spec by looping each BEAN through a barriered code→eval retry loop",
  phases: [
    { title: "Build", detail: "code agent implements each BEAN via TDD (blind to scenarios)" },
    { title: "Evaluate", detail: "eval agent proves each BEAN against holdout scenarios (blind to code)" },
  ],
};

// ── Input contract (args, supplied by /start:implement from plan.mdx) ─────────────────
// {
//   title:         string,
//   threshold:     number,          // satisfaction ratio for "done" (e.g. 0.90)
//   maxIterations: number,          // retry limit per BEAN (e.g. 5)
//   servicePort:   number,          // discovered service port
//   startCommand:  string,          // how the eval agent brings the service up
//   healthPath:    string,          // e.g. "/health" (eval agent polls this)
//   beans: [ {
//     id:          string,          // internal orchestration label (never leaks into code)
//     title:       string,
//     specContent: string,          // Features (ACs) + EARS behavior + quality gates — the code agent's brief
//     scenarioContent: string,      // concatenated holdout scenario text — the eval agent's brief
//   } ],
//   groups: [ { mode: "parallel" | "sequential", beanIds: string[] } ],  // topo-sorted
// }

const {
  title = "implementation",
  threshold = 0.9,
  maxIterations = 5,
  servicePort,
  startCommand,
  healthPath = "/health",
  beans = [],
  groups = [{ mode: "sequential", beanIds: beans.map((b) => b.id) }],
} = args ?? {};

const byId = new Map(beans.map((b) => [b.id, b]));

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

function codePrompt(bean, iteration, failures) {
  const retry =
    iteration > 0 && failures.length
      ? `\nA previous attempt did not satisfy the acceptance evidence ` +
        `(attempt ${iteration + 1} of ${maxIterations}). Observed symptoms:\n` +
        failures.map((f) => `- ${f}`).join("\n") +
        `\nInvestigate the codebase to find and fix the root causes.\n`
      : "";

  return `You are implementing the following slice of a specification using Test-Driven Development (TDD).

${bean.specContent}

Codebase context: read AGENTS.md (or CLAUDE.md) for project orientation and conventions.

DO NOT read or access any files under a scenarios/ directory, or under .start/specs/*/scenarios/.
Those hold the holdout acceptance evidence and must not influence how you build.

## File naming
Name test files after the feature or component under test — never after an internal id.
  ✅ tests/features/orders/test_order_repository.py
  ❌ tests/features/orders/test_bean3_repository.py

## TDD process (mandatory)
For EACH acceptance criterion and EARS behavior above, follow red-green-refactor:
1. RED: write a failing test for the behavior. Run the suite — confirm it fails.
2. GREEN: write the minimal code to make it pass. Run the suite — confirm it passes.
3. REFACTOR: improve the code while keeping every test green.
Do not write implementation code without a failing test first.
${retry}
When complete:
1. All existing tests pass.
2. New tests cover every acceptance criterion and EARS behavior above (written BEFORE implementation).
3. Run the full test suite and report results, plus how many red-green-refactor cycles you completed.`;
}

function evalPrompt(bean) {
  return `You are evaluating whether an implementation satisfies the following acceptance scenarios.

${bean.scenarioContent}

## Service lifecycle (you own it)
1. Start the service:  ${startCommand}
2. Poll http://localhost:${servicePort}${healthPath} until it responds (retry with backoff, ~5 attempts).
   If the service has no health endpoint, adapt using the project's conventions.
3. Run your evaluation against http://localhost:${servicePort}.
4. Tear the service down when you are done (kill the process you started).

## Evaluation method (use the first available, in priority order)
1. E2E test automation — write and run automated end-to-end tests for the scenarios, using the project's test framework.
2. Browser automation — drive the UI/API through a real browser.
3. CLI fallback — exercise endpoints directly with curl/httpie.

## Restrictions
DO NOT read source code, specification files, unit/BEAN specs, or any implementation details.
DO NOT access files under src/, plugins/, or .start/specs/*/ other than the scenarios you were given.
Test only through the external interface. Name any test files after the feature under test, never an internal id.

## Reporting (return the structured verdict)
Run each scenario 3 times; a scenario passes if it passes at least 2 of 3.
Return:
- satisfaction: passed_scenarios / total_scenarios (0.0–1.0)
- passed: names of scenarios that passed
- failures: for each failed scenario, ONE line stating the observable symptom (no internal detail)`;
}

// ── The BEAN loop — build → evaluate → (retry with filtered summaries) ─────────────────

async function buildBean(bean) {
  let iteration = 0;
  let failures = [];

  while (true) {
    // CODE agent: sees the BEAN spec + one-line retry summaries. Never the scenarios.
    await agent(codePrompt(bean, iteration, failures), {
      phase: "Build",
      label: `build:${bean.id}`,
    });

    // EVAL agent: sees the scenarios + service config. Never the source or the spec.
    const verdict = await agent(evalPrompt(bean), {
      phase: "Evaluate",
      label: `eval:${bean.id}`,
      schema: VERDICT,
    });

    const satisfaction = verdict?.satisfaction ?? 0;

    if (satisfaction >= threshold) {
      return { id: bean.id, title: bean.title, status: "completed", satisfaction, iterations: iteration + 1 };
    }

    iteration += 1;

    // Filter eval output down to one-line observable symptoms before it reaches the next code agent.
    failures = (verdict?.failures ?? []).slice(0, 6);

    if (iteration >= maxIterations) {
      return { id: bean.id, title: bean.title, status: "failed", satisfaction, iterations: iteration, failures };
    }

    log(`BEAN ${bean.id} at ${(satisfaction * 100).toFixed(0)}% (< ${(threshold * 100).toFixed(0)}%) — retry ${iteration}/${maxIterations}`);
  }
}

// ── Execution: run groups in order; parallelize independent BEANs within a group ───────

const results = [];

for (const group of groups) {
  const groupBeans = group.beanIds.map((id) => byId.get(id)).filter(Boolean);
  if (!groupBeans.length) continue;

  if (group.mode === "parallel") {
    const groupResults = await parallel(groupBeans.map((b) => () => buildBean(b)));
    results.push(...groupResults.filter(Boolean));
  } else {
    for (const b of groupBeans) {
      results.push(await buildBean(b));
    }
  }
}

const completed = results.filter((r) => r.status === "completed").length;
log(`${completed}/${results.length} BEANs satisfied for "${title}"`);

return { title, threshold, results };
