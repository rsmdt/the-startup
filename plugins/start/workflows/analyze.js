// The Agentic Startup — analysis workflow (deep-research style, tracker-backed).
//
// One reusable dynamic-workflow script, shipped with the `start` plugin at workflows/analyze.js. The
// /start:analyze SKILL is the entry point: it gets the plan approved, seeds the tracker, and launches
// this script in place via the Workflow tool with an `args` object. It is NOT regenerated per analysis.
// It ships in workflows/, so it is discoverable and namespaced as a command — but it is NOT a user
// entry point: it needs an `args` object only the skill can build (see the input contract below).
// Invoking it bare hits the fail-fast guard below and aborts with a pointer to the right command.
//
// Pipeline (grounded in deep-research literature — see SKILL.md "Synthesizer placement"):
//   RESEARCH  — fan out one agent per open research task; each gathers CITED findings and
//               writes them into its task via task.js. Never synthesizes a report.
//   VALIDATE  — one validator scores the findings against the `goals` rubric (coverage gate,
//               presumption-of-failure). Unmet goals → it creates gap research tasks via
//               task.js and returns their ids; loop back to RESEARCH.
//   SYNTHESIZE — after the loop converges, write the report from the validated findings.
//               (Evidence-first: validate coverage BEFORE writing.)
//   GRADE      — grade the written report against a quality rubric (citations, faithfulness, goal
//               coverage, mechanism depth) ONCE. There is no revision loop: a failing dimension is
//               returned to the skill for the human to act on, not fed back into another synthesis.
//
// Runtime note: the script has no shell/filesystem access — only its agents do. So every
// task.js call (read/update/create/graph) is made BY an agent; the script only orchestrates
// and branches on the structured verdicts agents return.

export const meta = {
  name: "run-analyze",
  description: "Deep-research analysis: fan out cited research, validate against a goals rubric, synthesize once at the end",
  phases: [
    { title: "Research", detail: "parallel agents gather cited findings per topic (code + online)" },
    { title: "Validate", detail: "score findings against the goals rubric; spawn gap research if unmet" },
    { title: "Synthesize", detail: "write the report from validated findings" },
    { title: "Grade", detail: "grade the report against a quality rubric; failures are reported, not revised" },
  ],
};

// ── Input contract (args, supplied by /start:analyze from the approved plan) ───────────
// {
//   name:        string,          // analysis slug
//   home:        string,          // abs path to .start/analysis/<NNN>-<name> (tracker home)
//   taskScript:  string,          // abs path to specify-meta/task.js
//   description: string,          // terse purpose
//   goals:       string[],        // the rubric — objectives that make the analysis review-worthy
//   topics:      [{ id, name }],  // initial research task ids, seeded by the skill
//   reportPath:  string,          // where SYNTHESIZE writes the report
//   maxRounds:   number,          // research↔validate loop cap (default 2)
// }

// Background workflow launches may serialize `args`; on some runtimes it is revived as a
// JSON string rather than an object. Parse defensively so destructuring never yields all-undefined.
const parsedArgs = typeof args === "string" ? JSON.parse(args) : args;

const {
  name = "analysis",
  home,
  taskScript,
  description = "",
  goals = [],
  topics = [],
  reportPath,
  maxRounds = 2,
} = parsedArgs ?? {};

// FAIL FAST on a mis-wired launch. This script is not a user entry point: it is launched by the
// /start:analyze skill with an `args` object it builds. Invoked bare (e.g. by typing the command), every
// prompt below would template literal "undefined" paths ("node undefined --home undefined"), research
// nothing, and synthesize an empty report. Abort before any agent runs instead.
{
  const required = { home, taskScript, reportPath };
  const missing = Object.entries(required)
    .filter(([, v]) => typeof v !== "string" || !v.trim() || v === "undefined")
    .map(([k]) => k);
  if (!goals.length) missing.push("goals");
  if (!topics.some((t) => t && t.id)) missing.push("topics");
  if (missing.length) {
    throw new Error(
      `run-analyze: launch args missing or unresolved: ${missing.join(", ")}. ` +
        `This workflow is launched BY the /start:analyze skill, not invoked directly — the skill gets the ` +
        `goals rubric approved and seeds the research tasks first. Run /start:analyze instead.`,
    );
  }
}

// Gap research per round. Round 1 fans out one agent per seeded topic; every later round batches its gap
// briefs into ONE agent. The gap fan-out used to consume whatever the validator returned, uncapped — the
// only term in this file with no ceiling at all.
const MAX_GAPS_PER_ROUND = 3;

const goalList = goals.map((g, i) => `${i + 1}. ${g}`).join("\n");
const tj = `node ${taskScript}`; // task.js invocation prefix agents use

// Structured verdict the validator returns, so the loop branches on data, not prose.
const VALIDATOR = {
  type: "object",
  properties: {
    satisfied: { type: "boolean", description: "true iff every goal is met by cited evidence" },
    goalResults: {
      type: "array",
      items: {
        type: "object",
        properties: { goal: { type: "string" }, met: { type: "boolean" }, note: { type: "string" } },
        required: ["goal", "met"],
      },
    },
    // ids of the gap research tasks the validator CREATED via task.js this round (empty if satisfied)
    nextTaskIds: { type: "array", items: { type: "string" } },
  },
  required: ["satisfied"],
};

// Report-quality grade — a rubric of per-dimension scores (fail | partial | pass). The gate is
// per-dimension: any dimension scored `fail` (must-fix) is surfaced to the human; `partial` ships with a
// note. `issues` localize the fixes for the human who reads the grade. The overall scalar is DERIVED from the scores by
// the script below — not declared by the model, so there is a single source of truth. `dimensions` is a
// rubric the grader may extend with problems it finds.
const GRADE = {
  type: "object",
  properties: {
    dimensions: {
      type: "array",
      items: {
        type: "object",
        properties: {
          name: { type: "string" },
          score: { type: "string", description: "fail | partial | pass" },
          note: { type: "string" },
        },
        required: ["name", "score"],
      },
    },
    issues: {
      type: "array",
      items: {
        type: "object",
        properties: {
          dimension: { type: "string" },
          detail: { type: "string", description: "the exact claim/section at fault" },
          fix: { type: "string", description: "the concrete change that resolves it" },
        },
        required: ["dimension", "detail"],
      },
    },
  },
  required: ["dimensions"],
};

// ── Prompt builders ────────────────────────────────────────────────────────────────────

// Takes a LIST of task ids. Round 1 calls this once per seeded topic (one id each, genuinely parallel
// angles); later rounds call it ONCE with every gap id, because gap briefs are narrow follow-ups that do
// not each warrant their own context — and fanning them out was the uncapped term in this file.
function researchPrompt(taskIds) {
  const many = taskIds.length > 1;
  return `You are a research agent for the analysis "${name}".
Purpose: ${description}

${many
      ? `You have ${taskIds.length} assignments, each in a tracker task. Read them all first:
${taskIds.map((id) => `  ${tj} read ${id} --home ${home}`).join("\n")}
Work them in order and give each its own genuine investigation — they are separate questions, not one
question in three parts. Record findings back into the task they belong to, separately.`
      : `Your assignment lives in a tracker task. Read it first:
  ${tj} read ${taskIds[0]} --home ${home}`}
It names your topic and may carry a brief from a prior validation round telling you what to dig into.

## How to research
- Code/repository questions: read the actual files. Cite as \`file:line\` (e.g. \`auth/config.ts:12\`).
- External/current questions: use web search + fetch. Cite the source URL.
- Investigate mechanisms, not surface facts. Prefer primary sources.

## Record findings (CITED) back into your task
Write your findings as markdown with an inline source on EVERY claim, then store them${many ? " — one file per task, into that task" : ""}:
  1. write them to a file, e.g. findings.md
  2. ${tj} update <task-id> --home ${home} --body-file findings.md
Before you mark ${many ? "a" : "the"} task passed, RE-READ your findings and confirm EVERY claim carries a source
(\`file:line\` or URL). If any claim is uncited, add the source or delete the claim — an uncited
claim is not a finding.
When done, mark ${many ? "each" : "the"} task researched:
  ${tj} update <task-id> --home ${home} --status passed

Do NOT write the final report — only gather cited findings. A claim with no source is not done.`;
}

function validatePrompt(round) {
  return `You are the ADVERSARIAL validation gate for the analysis "${name}" (round ${round + 1}).
Purpose: ${description}

## Adversarial stance (default to disbelief)
Assume NOTHING is true. Every goal is unmet and every claim is false until the evidence forces the
opposite conclusion. The burden of proof is on the findings, never on you. A confident sentence is not
evidence — only a source you have checked yourself is. Try to REFUTE each claim; keep only what survives.

## The goals rubric — the analysis is review-worthy only when ALL of these are met by verified evidence:
${goalList}

## Read the research gathered so far
  ${tj} graph --home ${home}                 (see all tasks + status)
  ${tj} list --home ${home} --type research   (the research tasks)
Then read each research task's findings:  ${tj} read <id> --home ${home}

## Verify each goal against the evidence — do not take findings on trust
For each goal, and each claim a finding makes toward it:
- OPEN the cited source yourself — read the \`file:line\` in the repo, or fetch the URL — and confirm it
  ACTUALLY substantiates the claim. A citation that does not support the claim counts as UNVERIFIED.
- A goal is met ONLY when every claim it rests on is independently verified against its source.
  Uncited, thin, unverifiable, or "the source doesn't actually say that" → NOT met.
- Rule NOT met only when you OPENED the source and it failed to substantiate the claim. "I did not check
  it" is not the same as "it is unsupported" — go and check it, then rule. Do not rule NOT met on a hunch:
  every unmet goal costs a full research round, so an unexamined doubt is not grounds for one.

## Citation integrity — check EVERY research task, not just goal coverage
Scan each research task's findings. If ANY claim lacks a source (\`file:line\` or URL), the analysis is
NOT satisfied — even when the goal is otherwise covered. Uncited findings are fixed in research, never
deferred to the report: open a gap task naming exactly which claims need a citation (or removal).

## If any goal is unmet OR any finding is uncited, open gap research (so the next round can close it)
For each gap, create a research task with a precise brief of what to investigate:
  ${tj} create --home ${home} --type research --name "<short topic>" --body "<what to research and why>"
Collect the ids you create — return them as nextTaskIds so the loop dispatches them.
Mark research tasks that fully satisfied their goal as passed if not already.

## Return the structured verdict
- satisfied: true only if every goal is met
- goalResults: per goal { goal, met, note }
- nextTaskIds: ids of the gap research tasks you just created (empty if satisfied)`;
}

function synthesizePrompt() {
  return `You are the synthesizer for the analysis "${name}". The research is complete and validated.
Purpose: ${description}

## Gather the validated findings
  ${tj} graph --home ${home}
  Read every research task's findings:  ${tj} read <id> --home ${home}

## Write the report to: ${reportPath}
- Structure it around the goals it must satisfy:
${goalList}
- Preserve EVERY citation from the findings — each claim keeps its source (code \`file:line\` or URL).
- Explain mechanisms and cross-cutting connections; do not just list facts.
- Do NOT introduce claims that have no source in the findings.

Write the file. Report the path and a one-line summary.`;
}

function gradePrompt() {
  return `You are the report grader for the analysis "${name}". Grade the report at ${reportPath} against
the rubric below. You do NOT fix anything — you score each dimension and list exactly what must change.

Score each dimension: \`pass\` (fully meets the bar), \`partial\` (mostly there with a fixable gap), or
\`fail\` (does not meet it). \`fail\` forces a revision; \`partial\` ships with a note.

## Rubric
1. Citations — every substantive claim carries a traceable source (code \`file:line\` or a URL).
   \`fail\` if ANY claim is uncited.
2. Faithfulness — no claim unsupported by the findings; nothing invented in the write-up.
   \`fail\` if any claim is unsupported.
3. Goal coverage — every goal below is addressed by the report's prose:
${goalList}
   \`fail\` if any goal is left unaddressed.
4. Mechanism depth — explains mechanisms and cross-cutting connections, not just lists facts.
   \`fail\` only if shallowness leaves a goal effectively unanswered; otherwise \`partial\`/\`pass\`.
If the report shows a quality problem none of these name, add it as your own dimension and score it.

## Return
- dimensions: score EACH dimension { name, score ("fail" | "partial" | "pass"), note }.
- issues: for every dimension scored below \`pass\`, one entry per concrete problem
  { dimension, detail, fix } — name the exact claim or section and the change that resolves it. Be specific.`;
}

// ── Orchestration: research → validate (loop) → synthesize → grade ─────────────────────

let round = 0;
let researchIds = topics.map((t) => t.id).filter(Boolean);
let verdict = null;

while (true) {
  if (researchIds.length) {
    if (round === 0) {
      // Round 1 only: the seeded topics are independent angles, so they are worth separate contexts.
      await parallel(researchIds.map((id) => () => agent(researchPrompt([id]), { phase: "Research", label: `research:${id}` })));
    } else {
      // Later rounds: gap briefs are narrow follow-ups. One agent takes them all.
      await agent(researchPrompt(researchIds), { phase: "Research", label: `research:gaps-r${round}` });
    }
  }

  verdict = await agent(validatePrompt(round), { phase: "Validate", label: `validate:r${round + 1}`, schema: VALIDATOR });

  if (verdict?.satisfied) break;

  round += 1;
  const proposed = (verdict?.nextTaskIds ?? []).filter(Boolean);
  researchIds = proposed.slice(0, MAX_GAPS_PER_ROUND);
  if (proposed.length > researchIds.length) {
    log(`round ${round}: validator opened ${proposed.length} gap task(s) — researching the first ${MAX_GAPS_PER_ROUND}; the rest stay open in the tracker.`);
  }

  if (round >= maxRounds || researchIds.length === 0) {
    log(`analysis "${name}": stopping after ${round} round(s) — ${researchIds.length ? "round cap reached" : "no further research proposed"}, goals not fully met`);
    break;
  }
  const unmet = (verdict?.goalResults ?? []).filter((g) => !g.met).length;
  log(`round ${round}: ${unmet} goal(s) unmet → ${researchIds.length} gap research task(s)`);
}

// Synthesize once, then grade once. There is no revision loop: the grader re-checked rules the
// synthesizer had already been given, and on cap exhaustion a failing grade only logged — a gate with no
// teeth costing up to four extra agents. The grade now travels in the return value instead, so the skill
// surfaces a failing dimension to the human, who is a better reviser than a third synthesis pass.
await agent(synthesizePrompt(), { phase: "Synthesize", label: "synthesize" });

const grade = await agent(gradePrompt(), { phase: "Grade", label: "grade", schema: GRADE });

const failingDimensions = (grade?.dimensions ?? []).filter((d) => String(d.score).toLowerCase() === "fail");
if (failingDimensions.length) {
  log(`analysis "${name}": report graded with ${failingDimensions.length} failing dimension(s) — see the grade in the result; the report ships as-is.`);
}

// Derived rubric value — single source of truth is the per-dimension scores (fail=0, partial=1, pass=2).
if (grade?.dimensions?.length) {
  const pts = grade.dimensions.reduce((s, d) => s + (({ fail: 0, partial: 1, pass: 2 })[d.score] ?? 0), 0);
  grade.overall = `${pts}/${grade.dimensions.length * 2}`;
}

return {
  name,
  satisfied: !!verdict?.satisfied,
  rounds: round + 1,
  reportPath,
  goalResults: verdict?.goalResults ?? [],
  grade,
  // Named explicitly so the skill can surface them: the report is no longer revised to clear the grade.
  failingDimensions: failingDimensions.map((d) => d.name),
};
