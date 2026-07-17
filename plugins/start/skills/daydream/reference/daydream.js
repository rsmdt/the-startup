// daydream-workflow: v4
//
// The Agentic Startup — daydream workflow (self-improvement from session logs, tracker-backed).
//
// One reusable dynamic-workflow script. /start:daydream installs this to
// <project>/.claude/workflows/daydream.js and launches it via the Workflow tool, passing an
// approved daydream plan as `args`. It is NOT regenerated per run.
//
// It reads the Claude Code session logs over a date range, finds where the agentic system
// underperformed, and PROPOSES improvements to this repo's skills/agents/workflows/templates —
// each proposal a concrete diff + the cited log evidence that motivates it. It is ADVISORY by
// default: v1 never edits a real artifact — it writes proposals to the tracker and a report a
// human reviews. The three-tier safe-apply ladder (advisory → bounded-patch → eval-gated) is the
// design roadmap (see .start/analysis/001-daydream-self-improvement/report.md §5.8); higher tiers
// are gated behind `applyTier` and are NOT implemented in v1.
//
// Pipeline (grounded in the analysis report — self-improvement literature + shipped analogues):
//   MINE     — fan out one agent per signal angle (tool-errors, user-corrections, retries, cost).
//              Each walks the session JSONL over the date range, does the two-pass join
//              (tool_use_id → {tool, attributionSkill}; then is_error results + interruptions,
//              grouped by skill), and writes CITED findings into its task via task.js.
//   CLUSTER  — ONE agent reads all mined findings, dedups/ranks them into distinct improvement
//              OPPORTUNITIES (each names one editable target + its evidence), and creates a
//              proposal task per opportunity via task.js. Returns their ids.
//   PROPOSE  — per opportunity, draft a concrete edit (a diff to the target artifact) + rationale,
//              and classify which safety tier the edit qualifies for. Records it into the task.
//   VALIDATE — adversarial gate: confirm the cited log evidence ACTUALLY supports the proposal
//              (presumption of failure). Unsupported proposals are dropped, not deferred.
//   SYNTHESIZE — ONCE, at the end: write the daydream report — ranked, validated proposals with
//              their diffs, evidence, and tier — grouped by target. A human applies them.
//
// Runtime note: the script has no shell/filesystem access — only its agents do. Every task.js call
// and every file read is made BY an agent; the script only orchestrates and branches on the
// structured verdicts agents return.
//
// Sandbox note: the session-log archive lives OUTSIDE the project root (under the user's Claude
// home, e.g. ~/.claude/projects/<slug>/). The command sandbox denies reads there by default, so a
// mine agent that shells into the archive may hit "Operation not permitted" and must retry with the
// sandbox disabled. Reading via the file-read tool is unaffected. Mine prompts say so explicitly.

export const meta = {
  name: "daydream",
  description: "Self-improvement from session logs: mine underperformance, propose cited diffs to skills/agents/workflows, validate the evidence, report",
  phases: [
    { title: "Mine", detail: "parallel agents mine the session JSONL per signal angle → cited findings" },
    { title: "Cluster", detail: "dedup/rank findings into improvement opportunities, one proposal task each" },
    { title: "Propose", detail: "draft a concrete diff + safety tier per opportunity" },
    { title: "Validate", detail: "adversarially confirm the log evidence supports each proposal" },
    { title: "Synthesize", detail: "write the proposals report; revise it to clear the grade" },
    { title: "Grade", detail: "grade each proposal on behaviour / improvement / proof; block and revise on failures" },
  ],
};

// ── Input contract (args, supplied by /start:daydream from the approved plan) ──────────────
// {
//   name:            string,          // run slug (e.g. "2026-07-16")
//   home:            string,          // abs path to .start/daydream/<slug> (tracker home)
//   taskScript:      string,          // abs path to specify-meta/task.js
//   description:     string,          // terse purpose
//   dateRange:       { since, until, label },  // ISO dates + a human label (e.g. "yesterday")
//   logRoot:         string,          // abs path to this project's session-log dir (~/.claude/projects/<slug>/)
//   editTargetsRoot: string,          // abs path to the git-tracked artifacts root (this repo's plugins/)
//   angles:          [{ id, name, kind? }],  // seeded mine task ids. kind: "friction" (default) mines
//                                             // structured underperformance; "router" mines model-routing
//                                             // (over/under-provisioned models per role); "conversation"
//                                             // reads dialogue content for semantic inefficiency (opus).
//   reportPath:      string,          // where SYNTHESIZE writes the report
//   applyTier:       string,          // safety ceiling: "advisory" (v1 default) | "bounded" | "eval"
//   maxProposals:    number,          // cap on proposals carried to the report (default 12)
// }

// Background launches may serialize `args`; on some runtimes it is revived as a JSON string.
const parsedArgs = typeof args === "string" ? JSON.parse(args) : args;

const {
  name = "daydream",
  home,
  taskScript,
  description = "",
  dateRange = {},
  logRoot,
  editTargetsRoot,
  angles = [],
  reportPath,
  applyTier = "advisory",
  maxProposals = 12,
} = parsedArgs ?? {};

const tj = `node ${taskScript}`; // task.js invocation prefix agents use
const range = `${dateRange.label ?? ""} (${dateRange.since ?? "?"} → ${dateRange.until ?? "?"})`;

// v1 ships advisory only. Guard loudly rather than silently pretending to auto-apply.
if (applyTier !== "advisory") {
  log(`daydream v1 is advisory-only; requested applyTier="${applyTier}" is not implemented — proceeding advisory (proposals are reported, not applied).`);
}

// ── Structured verdicts, so the loop branches on data, not prose ──────────────────────────

const CLUSTER = {
  type: "object",
  properties: {
    // ids of the proposal tasks the cluster agent CREATED via task.js (one per opportunity)
    proposalIds: { type: "array", items: { type: "string" } },
    droppedCount: { type: "number", description: "opportunities discarded as thin/duplicate" },
  },
  required: ["proposalIds"],
};

const PROPOSAL = {
  type: "object",
  properties: {
    target: { type: "string", description: "repo-relative path of the artifact to edit" },
    targetKind: { type: "string", description: "skill | agent | workflow | template | output-style" },
    tier: { type: "string", description: "advisory | bounded | eval — the highest tier this edit qualifies for" },
    hasEvalSuite: { type: "boolean", description: "true iff the target ships an evals/evals.json" },
    summary: { type: "string", description: "one line: what the edit changes and why" },
    // whether a concrete diff was written into the task body
    diffWritten: { type: "boolean" },
  },
  required: ["target", "summary"],
};

const VERDICT = {
  type: "object",
  properties: {
    supported: { type: "boolean", description: "true iff the cited log evidence actually motivates this edit" },
    reason: { type: "string" },
  },
  required: ["supported"],
};

// Report-quality grade — a rubric of per-dimension scores (fail | partial | pass). Daydream's rubric:
// every proposal must explain the BEHAVIOUR, the IMPROVEMENT, and the PROOF. The gate is per-dimension:
// any dimension scored `fail` drives a revision pass; `partial` ships with a note. `issues` localize the
// fixes for the reviser. The overall scalar is DERIVED from the scores by the script below — not declared
// by the model, so there is one source of truth. `dimensions` is extensible with problems the grader finds.
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
          detail: { type: "string", description: "the exact proposal at fault" },
          fix: { type: "string", description: "the concrete change that resolves it" },
        },
        required: ["dimension", "detail"],
      },
    },
  },
  required: ["dimensions"],
};

// ── Prompt builders ───────────────────────────────────────────────────────────────────────

function minePrompt(taskId) {
  return `You are a session-log mining agent for the daydream run "${name}".
Purpose: ${description}
Date range: ${range}

Your assignment (a signal angle) lives in a tracker task. Read it first:
  ${tj} read ${taskId} --home ${home}

## Where the data is
Session transcripts are JSONL under: ${logRoot}
- Top-level sessions:  <session-uuid>.jsonl
- Subagent transcripts (SEPARATE, must also be walked):  <session-uuid>/subagents/agent-<id>.jsonl
Reading this archive may require the shell sandbox to be DISABLED — it lives outside the project
root, so a plain read can fail with "Operation not permitted". If that happens, retry the command
with the sandbox disabled. Filter records to the date range by the ISO \`timestamp\` on each record.

## How to mine (two-pass join — highest-yield extraction)
1. Pass 1: from \`assistant\` records, build \`tool_use_id → { tool, attributionSkill }\`.
2. Pass 2: collect the signal for YOUR angle — e.g. \`tool_result\` blocks with \`is_error:true\`;
   \`[Request interrupted by user]\` markers; user-rejected tool uses; repeated identical tool_use
   inputs after an error in the \`parentUuid\` chain (retry loops); expensive turns from \`usage\`.
   Group by \`attributionSkill\` so each finding BLAMES a concrete skill/plugin.

## Record findings (CITED) back into your task
Every finding must cite its source — the session file + record (e.g.
\`<session-uuid>.jsonl record uuid=… / tool_use_id=…\`) and, where the fix points at code, the repo
\`file:line\`. Compose the findings as markdown and store them INLINE — do NOT create a findings.md (or
any report-shaped file); the subagent guard blocks report-file writes and a stray file pollutes the repo.
Pipe the markdown straight in via a shell heredoc:
  ${tj} update ${taskId} --home ${home} --body-file /dev/stdin <<'FINDINGS'
  ...your cited markdown findings here...
  FINDINGS
Then mark the task mined:  ${tj} update ${taskId} --home ${home} --status passed

Report counts (e.g. "skill X caused N errors of category C over the range"), not narratives. A finding
with no source is not a finding. Do NOT propose fixes here — only gather cited signal.`;
}

function conversationPrompt(taskId) {
  return `You are the CONVERSATION-QUALITY mining agent for the daydream run "${name}".
Purpose: ${description}
Date range: ${range}

Read your assignment:  ${tj} read ${taskId} --home ${home}
You read the actual dialogue — user and assistant message CONTENT — and judge where the interaction was
inefficient, including inefficiencies that leave NO structured trace (nothing errored, nothing was
interrupted — the exchange was simply less efficient than it could have been). The goal is general
conversational efficiency, not a fixed checklist.

## Where the data is
Session transcripts (JSONL) under: ${logRoot} — top-level <session-uuid>.jsonl and
<session-uuid>/subagents/agent-<id>.jsonl. Read \`assistant\` records' \`message.content\` (text and
thinking blocks) and \`user\` records' typed prompts (string \`message.content\`). Filter by ISO
\`timestamp\` to the range. Reading this archive may need the sandbox DISABLED — retry unsandboxed on
"Operation not permitted".

## Scope by SAMPLING — do not read every session in full
There may be many sessions; reading all fully is prohibitively expensive. Prioritize the highest-signal
ones (most turns, longest, or those with visible back-and-forth) and take a representative sample of the
rest. REPORT exactly which sessions you sampled and how many you skipped — never imply full coverage.

## What to hunt (these are ILLUSTRATIONS, not an exhaustive list)
Surface ANY recurring pattern where the exchange was less efficient than it could have been, e.g.:
- Under-delivery → follow-up: the user had to ask for detail/scope the assistant should have supplied
  up front ("actually also…", "you forgot…", "give me the full X").
- Flip-flopping: the assistant reversed a recommendation with little new information, confusing the user.
- Over-asking / ceremony: menus or clarifying questions where a sensible default existed; verbosity the
  user pushed back on ("keep it terse", "no ceremony").
- Ignored-then-restated: a constraint the user had to repeat because it was dropped.
- Re-explanation loops: the assistant needed several tries to land a point.
- Premature confidence or unverified claims the user then had to challenge.
Invent categories the examples don't cover if the dialogue shows them — you are looking for GENERAL
inefficiency, not matching a template.

## What makes it a finding
A single awkward turn is noise. A PATTERN that recurs across sessions (or a costly single instance) is a
finding. For each: name the pattern, quote the MINIMAL telltale exchange (a few words, not whole turns —
these are real conversations, keep quotes short), cite every instance (session file + record uuid), and
name the concrete target whose instructions would prevent it — usually a skill's Constraints, a team
agent, or an output-style under ${editTargetsRoot}.

## Record findings (CITED) inline — do NOT create a file
  ${tj} update ${taskId} --home ${home} --body-file /dev/stdin <<'FINDINGS'
  ...your cited conversation-quality patterns + target per pattern here...
  FINDINGS
Then:  ${tj} update ${taskId} --home ${home} --status passed
Diagnose only; do NOT edit any file.`;
}

function routerPrompt(taskId) {
  return `You are the MODEL-ROUTING mining agent for the daydream run "${name}".
Purpose: ${description}
Date range: ${range}

Read your assignment:  ${tj} read ${taskId} --home ${home}
You mine the session logs for where the wrong-sized LLM model was used — over-provisioned roles paying
for capability they didn't need, and under-provisioned roles that thrashed for lack of it.

## Where the data is
Session transcripts (JSONL) under: ${logRoot} — top-level <session-uuid>.jsonl AND
<session-uuid>/subagents/agent-<id>.jsonl. Subagent role is in the sidecar agent-<id>.meta.json
(\`agentType\`, \`description\` — e.g. a workflow's mine/validate/code/eval agent). Filter by the ISO
\`timestamp\` to the date range. Reading this archive may need the shell sandbox DISABLED — retry
unsandboxed on "Operation not permitted".

## What to extract (per ROLE, not per session)
For each assistant record: \`message.model\` (which model ran), \`message.usage\` (token spend), and the
\`attributionSkill\`/subagent role it belongs to. Group by role. Then cross-reference the SAME role's
outcome signal in that session:
- OVER-provisioned → a strong model (e.g. opus) doing mechanical, schema-filled, or high-fan-out work
  that was never corrected, retried, or eval-failed. Candidate to DOWN-tier (opus→sonnet→haiku) or
  lower \`effort\`.
- UNDER-provisioned → a cheaper model (haiku/sonnet) that THRASHED: repeated retries after error in the
  \`parentUuid\` chain, user corrections, or eval-loop failures. Candidate to UP-tier / raise \`effort\`.
A verifier/gate role must never be recommended BELOW the generator it checks — call that out if seen.

## Record findings (CITED) inline
Every recommendation cites: the session file + record uuid, the \`message.model\` and \`usage\` numbers,
and the outcome signal (retry/correction/eval-fail record, or its absence). Name the concrete target a
proposal would edit — an \`agent(..., { model, effort })\` call in a workflow golden copy under
${editTargetsRoot} (\`plugins/*/skills/*/reference/*.js\`) or a \`model:\` frontmatter in a team agent
(\`plugins/team/agents/*/*.md\`). Store findings INLINE — do NOT create a file:
  ${tj} update ${taskId} --home ${home} --body-file /dev/stdin <<'FINDINGS'
  ...your cited routing findings + per-role model/effort recommendation here...
  FINDINGS
Then:  ${tj} update ${taskId} --home ${home} --status passed
Recommend routing only; do NOT edit any file.`;
}

function clusterPrompt() {
  return `You are the clustering gate for the daydream run "${name}".
Purpose: ${description}

## Read every mined finding
  ${tj} graph --home ${home}
  ${tj} list --home ${home} --type mine
Read each mine task's findings:  ${tj} read <id> --home ${home}

## Turn signal into distinct improvement OPPORTUNITIES
- Dedup: collapse findings that point at the same root cause / same target artifact.
- Rank by (frequency × severity): a category that recurred often and blocked work outranks a one-off.
- Each opportunity must name ONE editable target under ${editTargetsRoot} (a skill/agent/workflow/
  template/output-style file) and carry the cited evidence that motivates it.
- Drop thin or unactionable signal (transient env errors, one-offs with no pattern). Count what you drop.
- Keep at most ${maxProposals} opportunities.

## Create one proposal task per opportunity
  ${tj} create --home ${home} --type proposal --name "<target>: <one-line opportunity>" --body "<the cited evidence + why this target>"
Collect the ids you create — return them as proposalIds. Also return droppedCount.`;
}

function proposePrompt(taskId) {
  return `You are a proposal agent for the daydream run "${name}".
Read your opportunity task:  ${tj} read ${taskId} --home ${home}
It names ONE target artifact under ${editTargetsRoot} and the cited log evidence motivating a change.

## Draft a concrete edit
- Read the target artifact in full first.
- Write a MINIMAL, surgical diff (unified-diff style, or clearly-marked before/after) that addresses
  the evidenced problem — nothing speculative, nothing beyond what the evidence supports.
- Prefer a small patch to an existing section over a rewrite (bounded edits are safer to land).

## Classify the safety tier (does NOT auto-apply — this is advisory metadata for the human)
- \`eval\`   — the target directory ships an evals/evals.json that could gate the edit.
- \`bounded\` — no eval, but the edit is small and structurally checkable (e.g. a frontmatter/section patch).
- \`advisory\` — anything larger or judgement-heavy; a human must apply and verify it.
Report hasEvalSuite (true iff evals/evals.json exists next to the target).

## Store the proposal
Store the diff + rationale (keeping every evidence citation) INLINE — do NOT create a proposal.md or
any scratch file (it pollutes the repo). Pipe it straight in via a shell heredoc:
  ${tj} update ${taskId} --home ${home} --body-file /dev/stdin <<'PROPOSAL'
  ...your diff + rationale + citations here...
  PROPOSAL
Return { target, targetKind, tier, hasEvalSuite, summary, diffWritten }.
Do NOT edit the real artifact — only propose.`;
}

function validatePrompt(taskId) {
  return `You are the ADVERSARIAL validation gate for a daydream proposal (default to disbelief).
Read the proposal:  ${tj} read ${taskId} --home ${home}

Assume the proposed edit is UNJUSTIFIED until the cited log evidence forces the opposite conclusion.
For each citation the proposal rests on:
- OPEN the cited source yourself — the session-log record it points to (disable the sandbox if the
  read is denied) and any repo \`file:line\` — and confirm it ACTUALLY shows the failure/friction the
  proposal claims. A citation that does not substantiate the claim counts as UNSUPPORTED.
- The proposal is supported ONLY when the evidence genuinely motivates THIS edit (not merely "an error
  happened nearby"). The edit must plausibly reduce the evidenced friction.
When in doubt, rule NOT supported — a dropped weak proposal is better than a wrong edit suggestion.

Return { supported, reason }.`;
}

function gradePrompt() {
  return `You are the report grader for the daydream run "${name}". Grade the proposals report at
${reportPath} against the rubric below. You do NOT fix anything — you assess, and list what must change.
Every proposal must stand on three legs: the BEHAVIOUR it addresses, the IMPROVEMENT it makes, and the
PROOF behind it.

Score each dimension: \`pass\` (every proposal meets it), \`partial\` (mostly, with a fixable gap), or
\`fail\` (one or more proposals miss it). \`fail\` forces a revision; \`partial\` ships with a note.

## Rubric — across all proposals in the report
1. Behaviour — each proposal names the concrete behaviour/pattern it fixes: what the system actually did
   that was inefficient or wrong, not merely "an error occurred nearby". \`fail\` if any is vague/missing.
2. Improvement — each carries a concrete, minimal edit (a diff or clearly-marked before/after) to a real
   target artifact, not a vague suggestion. \`fail\` if any lacks an actionable change.
3. Proof — each behaviour is backed by traceable evidence (session-log record uuid and/or repo
   \`file:line\`) a reader can open, that genuinely motivates THIS change. \`fail\` if any is unsourced or
   the evidence does not support the change.
4. Faithfulness — the report contains no proposal absent from the validated tracker; nothing invented.
   \`fail\` if any is invented.
If a proposal shows a quality problem none of these name, add it as your own dimension and score it.

## Return
- dimensions: score EACH dimension { name, score ("fail" | "partial" | "pass"), note }.
- issues: for every dimension scored below \`pass\`, one entry per concrete problem { dimension, detail,
  fix } — name the exact proposal and the change that resolves it. Be specific.`;
}

function synthesizePrompt(kept, issues) {
  if (issues && issues.length) {
    const fixes = issues
      .map((x) => `- ${x.dimension}: ${x.detail}${x.fix ? ` → ${x.fix}` : ""}`)
      .join("\n");
    return `You are revising the daydream proposals report at ${reportPath} to fix specific grading issues.
Read the current report first.
Fix EACH issue below, changing only what it requires — preserve every other proposal, diff, and citation:
${fixes}
If closing an issue needs evidence or a diff, re-open the proposal task it refers to:
  ${tj} list --home ${home} --type proposal
  ${tj} read <id> --home ${home}
Do NOT invent a proposal that is not in the tracker, and do NOT weaken any citation.
Rewrite the file in place. Report the path and a one-line summary of what you changed.`;
  }
  const list = kept
    .map((p, i) => `${i + 1}. [${p.tier ?? "advisory"}] ${p.target} — ${p.summary} (task ${p.id})`)
    .join("\n");
  return `You are the synthesizer for the daydream run "${name}". Proposals are drafted and validated.
Purpose: ${description}
Date range: ${range}

## The validated proposals to include (read each task for its diff + evidence)
${list || "(none survived validation)"}
Read each:  ${tj} read <id> --home ${home}

## Write the report to: ${reportPath}
- Open with a short summary: what range was mined, how many findings/opportunities, how many proposals
  survived validation, and the per-skill friction ranking.
- Then one section per proposal, GROUPED by target artifact and ordered by tier then rank:
  the concrete diff, the one-line rationale, the cited log evidence, and the safety tier.
- Preserve EVERY citation. Do NOT invent proposals not in the tracker.
- Close with a "How to apply" note: daydream is advisory — a human reviews each diff and applies it on
  a branch/PR; note which proposals ship an eval suite (tier \`eval\`) as the safest to automate later.

Write the file. Report the path and a one-line summary.`;
}

// ── Orchestration: mine → cluster → (propose → validate) pipeline → synthesize ─────────────

const mineAngles = angles.filter((a) => a && a.id);
const mineIds = mineAngles.map((a) => a.id);

if (mineAngles.length) {
  await parallel(
    mineAngles.map((a) => () => {
      // Semantic conversation mining is judgment-heavy → strong model; the rest stay on the inherited
      // model (mechanical extraction). This is the router principle applied to daydream's own agents.
      if (a.kind === "conversation") {
        return agent(conversationPrompt(a.id), { phase: "Mine", label: `convo:${a.id}`, model: "opus", effort: "high" });
      }
      if (a.kind === "router") {
        return agent(routerPrompt(a.id), { phase: "Mine", label: `route:${a.id}` });
      }
      return agent(minePrompt(a.id), { phase: "Mine", label: `mine:${a.id}` });
    }),
  );
}

const cluster = await agent(clusterPrompt(), { phase: "Cluster", label: "cluster", schema: CLUSTER });
const proposalIds = (cluster?.proposalIds ?? []).filter(Boolean).slice(0, maxProposals);

let proposals = [];
if (proposalIds.length) {
  proposals = await pipeline(
    proposalIds,
    (id) => agent(proposePrompt(id), { phase: "Propose", label: `propose:${id}`, schema: PROPOSAL }).then((p) => ({ ...p, id })),
    (p, id) =>
      agent(validatePrompt(id), { phase: "Validate", label: `validate:${id}`, schema: VERDICT }).then((v) => ({ ...p, id, verdict: v })),
  );
}

const kept = proposals.filter(Boolean).filter((p) => p.verdict?.supported);
const dropped = proposals.filter(Boolean).length - kept.length;
log(`daydream "${name}": ${mineIds.length} angle(s) mined → ${proposalIds.length} opportunit(ies) → ${kept.length} proposal(s) validated (${dropped} dropped).`);

// Synthesize the proposals report, then grade it; blocking issues drive a revision pass. Loop until the
// report clears the grade or the cap — the grade gate has teeth, it does not just attach a warning.
const maxGradeRounds = 3; // initial synthesis + up to 2 revisions
let grade = null;

for (let g = 0; ; g++) {
  await agent(synthesizePrompt(kept, g === 0 ? null : grade?.issues ?? []), {
    phase: "Synthesize",
    label: g === 0 ? "synthesize" : `revise:r${g}`,
  });

  grade = await agent(gradePrompt(), { phase: "Grade", label: `grade:r${g + 1}`, schema: GRADE });

  const failing = (grade?.dimensions ?? []).filter((d) => d.score === "fail");
  if (!failing.length || g + 1 >= maxGradeRounds) {
    if (failing.length) log(`daydream "${name}": report graded with ${failing.length} failing dimension(s) after ${g + 1} pass(es).`);
    break;
  }
  log(`grade round ${g + 1}: ${failing.length} failing dimension(s) → revising report`);
}

// Derived rubric value — single source of truth is the per-dimension scores (fail=0, partial=1, pass=2).
if (grade?.dimensions?.length) {
  const pts = grade.dimensions.reduce((s, d) => s + (({ fail: 0, partial: 1, pass: 2 })[d.score] ?? 0), 0);
  grade.overall = `${pts}/${grade.dimensions.length * 2}`;
}

return {
  name,
  dateRange,
  mined: mineIds.length,
  opportunities: proposalIds.length,
  proposals: proposals.filter(Boolean).length,
  supported: kept.length,
  dropped,
  applyTier: "advisory",
  reportPath,
  grade,
};
