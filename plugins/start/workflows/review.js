// The Agentic Startup — code review workflow (dynamic dimensions → adversarial verify → rank → report).
//
// One reusable dynamic-workflow script, shipped with the `start` plugin at workflows/review.js. The
// /start:review SKILL is the entry point: it launches this script in place via the Workflow tool,
// passing a resolved review target as `args`.
// It ships in workflows/, so it is discoverable and namespaced as a command — but it is NOT a user
// entry point: it needs an `args` object only the skill can build (see the input contract below).
// Invoking it bare hits the fail-fast guard below and aborts with a pointer to the right command.
//
// It generalizes the solution-design review (solution.js) to code: same pipeline — fan out lenses,
// adversarially verify each finding, synthesize, report — with ONE addition, a triage stage, because a
// code diff's concerns aren't known up front the way an SDD's are.
//
// Pipeline:
//   TRIAGE   — (scoped only) read the diff and choose the review dimensions dynamically: an always-active
//              core that can't be switched off, plus conditional dimensions the change's signals fire, plus
//              a per-dimension reviewer count + effort scaled to blast radius. Audit mode skips triage — its
//              dimensions are the four fixed phases.
//   REVIEW   — fan out one agent per (dimension × file-group). Each reads the full file contents + the diff,
//              reviews ONLY through its dimension against the evidence standard, and returns findings.
//   VERIFY   — adversarially confirm each finding against the code (default to disbelief); drop false
//              positives and demote unconsequenced ones to Accepted.
//   SYNTHESIZE — one agent dedups by location, assigns ids, ranks, and (audit) distills the top refactors.
//   REPORT   — write a MARKDOWN report. The verdict is computed in the script from severity counts.
//
// Runtime note: the script has no shell/filesystem access — only agents do. The triage/review/verify/
// report agents run git, read files, and write the report; the script orchestrates and branches on
// structured verdicts. The information barrier here is soft (review is not a build/eval loop), but the
// VERIFY agent is a DIFFERENT agent from the reviewer, so findings are confirmed independently.

export const meta = {
  name: "run-review",
  description: "Code review: pick dimensions from scope, fan out reviewers, adversarially verify findings, rank, report",
  phases: [
    { title: "Triage", detail: "read the diff and choose review dimensions + fan-out scale (scoped mode)" },
    { title: "Review", detail: "parallel reviewers, one per dimension × file-group, against the evidence standard" },
    { title: "Verify", detail: "adversarially confirm each finding against the code; drop false positives" },
    { title: "Synthesize", detail: "dedup, rank, assign ids, distill refactors (audit)" },
    { title: "Report", detail: "write the markdown review report; verdict derived from severity counts" },
  ],
};

// ── Input contract (args, supplied by /start:review) ───────────────────────────────────────
// {
//   name:        string,          // target label for the report header
//   mode:        "scoped" | "audit",
//   diffCommand: string,          // scoped: how an agent gets the diff (e.g. "git diff main...feat")
//   changedFiles:string[],        // scoped: files to read in full
//   hasConstitution: boolean,     // scoped: project has CONSTITUTION.md → Constitution dimension fires
//   orientation: string,          // audit: one-paragraph system orientation (skill computed it)
//   coChange:    string,          // audit: co-change archaeology summary (skill computed it once)
//   reportPath:  string,          // abs path for the MARKDOWN report (skill bakes the date into it)
//   classificationPath: string,   // skill's reference/classification.md (severity, consequence test, ranking)
//   checklistsPath:     string,   // skill's reference/checklists.md (per-dimension scoped checklists)
//   auditPhasesPath:    string,   // skill's reference/audit-phases.md (audit phase instructions)
// }

const parsedArgs = typeof args === "string" ? JSON.parse(args) : args;

const {
  name = "review",
  mode = "scoped",
  diffCommand = "git diff",
  changedFiles = [],
  hasConstitution = false,
  orientation = "",
  coChange = "",
  reportPath,
  classificationPath = "",
  checklistsPath = "",
  auditPhasesPath = "",
} = parsedArgs ?? {};

// FAIL FAST on a mis-wired launch. This script is not a user entry point: it is launched by the
// /start:review skill, which resolves the mode, the diff command, and the changed-file list. Invoked
// bare, the reviewers would be pointed at nothing and the report path would be "undefined".
{
  const required = { reportPath };
  const missing = Object.entries(required)
    .filter(([, v]) => typeof v !== "string" || !v.trim() || v === "undefined")
    .map(([k]) => k);
  if (mode === "scoped" && !changedFiles.length && !diffCommand.trim()) missing.push("diffCommand/changedFiles");
  if (mode === "audit" && !auditPhasesPath.trim()) missing.push("auditPhasesPath");
  if (missing.length) {
    throw new Error(
      `run-review: launch args missing or unresolved: ${missing.join(", ")}. ` +
        `This workflow is launched BY the /start:review skill, not invoked directly — the skill resolves ` +
        `the review mode and scope first. Run /start:review instead.`,
    );
  }
}

// The evidence standard — carried verbatim into every reviewer and verifier so the review stays
// calibrated (no finding without a named consequence; steelman defensible deviations into Accepted).
const EVIDENCE = `Evidence standard: a finding is valid ONLY if it names the concrete consequence — the
future change, bug class, or comprehension cost it causes — at a specific \`file:line\`. A principle
(SOLID, DRY) may EXPLAIN a finding, never CONSTITUTE one. If the code deviates from textbook practice for
a defensible reason, STEELMAN it: state the reason and file it under Accepted (harmless here) rather than
reporting it. No finding without a named consequence.`;

// ── Structured verdicts ─────────────────────────────────────────────────────────────────────

// Triage picks the scoped review panel dynamically. The core dimensions are always present; triage adds
// conditional ones and assigns per-dimension reviewer count + effort by blast radius.
const TRIAGE = {
  type: "object",
  properties: {
    dimensions: {
      type: "array",
      items: {
        type: "object",
        properties: {
          key: { type: "string" },
          title: { type: "string" },
          focus: { type: "string", description: "what this dimension looks for, one or two sentences" },
          reviewers: { type: "number", description: "parallel reviewers for this dimension: 1, or 2 for a heavy dimension on a large diff. The script clamps anything higher." },
          effort: { type: "string", description: "low | medium | high — reasoning effort for this dimension's reviewers" },
          rationale: { type: "string", description: "why this dimension fired (or its scale), for the audit trail" },
        },
        required: ["key", "title", "focus"],
      },
    },
  },
  required: ["dimensions"],
};

const FINDINGS = {
  type: "object",
  properties: {
    findings: {
      type: "array",
      items: {
        type: "object",
        properties: {
          severity: { type: "string", description: "CRITICAL | HIGH | MEDIUM | LOW" },
          confidence: { type: "string", description: "HIGH | MEDIUM | LOW" },
          file: { type: "string" },
          line: { type: "number" },
          title: { type: "string" },
          consequence: { type: "string", description: "the named victim — REQUIRED" },
          remediation: { type: "string" },
          codeBefore: { type: "string" },
          codeAfter: { type: "string" },
        },
        required: ["severity", "file", "title", "consequence"],
      },
    },
    accepted: {
      type: "array",
      items: {
        type: "object",
        properties: { smell: { type: "string" }, location: { type: "string" }, reasoning: { type: "string" } },
        required: ["smell", "reasoning"],
      },
    },
  },
  required: ["findings"],
};

// One verifier judges ALL of a reviewer's findings in a single pass. Per-finding verification was
// multiplicative (reviewers × findings, both LLM-chosen and uncapped) and drove the worst case past the
// runtime's own agent ceiling. The adversarial framing is preserved verbatim in the prompt; what is lost
// is per-finding context isolation, which the prompt counters with an explicit anti-anchoring rule.
const VERIFY_BATCH = {
  type: "object",
  properties: {
    verdicts: {
      type: "array",
      items: {
        type: "object",
        properties: {
          index: { type: "number", description: "the finding's index in the numbered list you were given" },
          real: { type: "boolean", description: "true iff the code actually has this flaw AND its consequence is named" },
          demoteToAccepted: { type: "boolean", description: "true if it is a real observation but harmless here" },
          reason: { type: "string" },
        },
        required: ["index", "real"],
      },
    },
  },
  required: ["verdicts"],
};

const SYNTHESIS = {
  type: "object",
  properties: {
    findings: {
      type: "array",
      items: {
        type: "object",
        properties: {
          id: { type: "string", description: "severity letter + number: C1, H1, M1, L1" },
          severity: { type: "string", description: "CRITICAL | HIGH | MEDIUM | LOW" },
          confidence: { type: "string" },
          source: { type: "string", description: "the dimension/phase it came from" },
          file: { type: "string" },
          line: { type: "number" },
          title: { type: "string" },
          consequence: { type: "string" },
          remediation: { type: "string" },
          codeBefore: { type: "string" },
          codeAfter: { type: "string" },
        },
        required: ["id", "severity", "file", "title", "consequence"],
      },
    },
    accepted: { type: "array", items: { type: "object", properties: { smell: { type: "string" }, location: { type: "string" }, reasoning: { type: "string" } } } },
    refactors: {
      type: "array",
      description: "audit mode only — the highest-leverage refactors, each tied to a measured finding",
      items: {
        type: "object",
        properties: { rank: { type: "number" }, what: { type: "string" }, why: { type: "string" }, migration: { type: "array", items: { type: "string" } } },
        required: ["what", "why"],
      },
    },
    strengths: { type: "array", items: { type: "string" } },
  },
  required: ["findings"],
};

// ── Prompt builders ─────────────────────────────────────────────────────────────────────────

function triagePrompt() {
  return `You are the triage stage of a code review of "${name}". Read the change and choose the review
dimensions dynamically — do not review the code yourself.

Get the diff and read the changed files in full:
  diff:  run \`${diffCommand}\`
  files: ${changedFiles.length ? changedFiles.join(", ") : "(from the diff)"}

## Choose the dimensions (Dial A — which fire)
ALWAYS include this core (it fires on every change; never drop one):
  correctness (logic, edge cases, error paths, does it do what the change claims, AND whether the new
    paths are actually covered by tests),
  security,
  performance,
  quality (naming, patterns, smells, AND YAGNI / over-engineering — code that does more than it needs to).
ADD conditional dimensions when the change's signals fire them:
  concurrency (async / shared state), dependencies (manifest changes), compatibility (public API / schema /
  config change), data-migration (migration or schema files), accessibility (UI changes)${hasConstitution ? `,
  constitution (this project HAS a CONSTITUTION.md — include it)` : ""}.
Fire a conditional dimension on ANY hard signal (a manifest change fires dependencies even if it "looks"
trivial). If a change looks cosmetic but touches a risky path (auth, money, data), fire the matching
dimension anyway and say why in the rationale.

## Scale each dimension (Dial B — how much)
Set reviewers (1 or 2) and effort per dimension by blast radius, NOT file type: a tiny one-file diff → 1
reviewer, low effort each; a large multi-subsystem diff → split a heavy dimension across two file-groups
(reviewers: 2) and raise effort ("high") on the risky ones (security on an auth change, correctness on
core logic). Effort is your main dial, not reviewer count — a "high"-effort single reviewer beats two
shallow ones. ${MAX_REVIEWERS_PER_DIM} reviewers per dimension is the ceiling and the script enforces it,
so asking for more only gets clamped.

Return { dimensions: [{ key, title, focus, reviewers, effort, rationale }] }.`;
}

function reviewPrompt(task) {
  const scope =
    mode === "audit"
      ? `## Scope — whole codebase (audit)
System orientation (given): ${orientation || "(none provided)"}
Co-change archaeology (given, do not recompute): ${coChange || "(none provided)"}
Read the phase instructions for your phase:  ${auditPhasesPath}`
      : `## Scope — the change under review
Get the diff:  run \`${diffCommand}\`
Read these changed files IN FULL (not just the diff — context matters):  ${(task.files ?? changedFiles).join(", ") || "(from the diff)"}
Read your dimension's checklist:  ${checklistsPath}`;
  return `You are a code reviewer for "${name}", reviewing through ONE dimension: ${task.title}.

## Your dimension
${task.focus}

${scope}

Read the severity / confidence / consequence rules:  ${classificationPath}

${EVIDENCE}

Review ONLY through your dimension — other dimensions are covered by other reviewers. Return every real
finding with its \`file:line\`, severity, confidence, named consequence, and remediation (code before/after
for CRITICAL). Put defensible-deviation smells under \`accepted\` with the one-line reason they are harmless
here. Return an empty findings list if your dimension is clean — do not invent findings to fill a quota.

Return { findings: [...], accepted: [...] }.`;
}

function verifyBatchPrompt(task, findings) {
  const list = findings
    .map((f, i) => `${i}. ${f.severity} — ${f.title}\n   at ${f.file}:${f.line ?? "?"}\n   consequence claimed: ${f.consequence}`)
    .join("\n");
  return `You are the ADVERSARIAL verifier for the ${findings.length} review finding(s) below on "${name}"
(dimension: ${task.title}). Default to disbelief — assume each finding is WRONG until the code forces the
opposite conclusion.

## The findings
${list}

Open the code yourself (read each cited file, and run \`${diffCommand}\` if you need the diff) and check
EACH finding independently:
1. Does the code ACTUALLY have this flaw at that location? A finding that misreads the code is NOT real.
2. Is the consequence genuinely named and plausible? ${EVIDENCE}
If it is a real observation but harmless in this context, set demoteToAccepted. When in doubt, rule NOT
real — a false finding destroys the review's credibility.

Judge each finding ON ITS OWN EVIDENCE. The verdicts are independent: having just confirmed three findings
is no reason to wave the fourth through, and rejecting one is no reason to reject the next. A run of
identical verdicts means you have stopped opening the code — go back and open it.

Return { verdicts: [{ index, real, demoteToAccepted, reason }] } — one entry per finding, \`index\` matching
the numbering above.`;
}

function synthesizePrompt(confirmed, accepted) {
  const list = confirmed
    .map((f) => `- [${f.severity}/${f.confidence ?? "?"}]${f.unverified ? "[UNVERIFIED]" : ""} ${f.source}: ${f.title} @ ${f.file}:${f.line ?? "?"} — ${f.consequence}`)
    .join("\n");
  const unverifiedCount = confirmed.filter((f) => f.unverified).length;
  return `You are the synthesizer for the code review of "${name}" (${mode} mode). The findings below are
already verified — do not re-review. Consolidate them.${unverifiedCount
      ? `\n\nNOTE: ${unverifiedCount} finding(s) are tagged [UNVERIFIED] — their reviewer returned more findings
than the verification cap allowed, so no adversarial verifier checked them. Keep them, but mark each one
"unverified" in its entry so the reader knows its confidence is lower.`
      : ""}

## Confirmed findings
${list || "(none)"}

## Do
1. Deduplicate by location (same file within ~5 lines) — keep the highest severity, merge complementary detail.
2. Assign ids: C1, H1, M1, L1 … ${mode === "audit"
      ? "Rank by likelihood the consequence materializes × cost when it does (see " + classificationPath + ")."
      : "Sort by severity, then confidence."}
3. Carry each finding's source (dimension), file:line, consequence, remediation, and code before/after if present.
${mode === "audit" ? "4. Distill the FIVE highest-leverage refactors, each tied to a specific finding, each with an incremental migration path (no rewrites).\n5. List genuine strengths." : "4. List genuine strengths of the change."}

Return { findings: [ranked, with ids], accepted, ${mode === "audit" ? "refactors, " : ""}strengths }.
(The Accepted list already gathered: ${accepted.length} item(s) — preserve them.)`;
}

function reportPrompt(synth, verdict, counts) {
  const data = {
    findings: synth?.findings ?? [],
    accepted: synth?.accepted ?? [],
    ...(mode === "audit" ? { refactors: synth?.refactors ?? [] } : {}),
    strengths: synth?.strengths ?? [],
  };
  return `You are writing the code review report for "${name}" (${mode} mode) to: ${reportPath}
Write MARKDOWN (not HTML). Render ALL of the data below — preserve every \`file:line\` and consequence,
invent nothing, drop nothing (LOW findings go in the table too).

## Data to render (JSON)
${JSON.stringify(data, null, 2)}

## Structure
# Code Review — ${name}

**Verdict:** ${verdict}
**Findings:** ${counts.critical} critical · ${counts.high} high · ${counts.medium} medium · ${counts.low} low

## Findings
A table — columns: ID · Severity · Confidence · Location (\`file:line\`) · Dimension (source) · Title ·
Consequence — then, for each CRITICAL/HIGH, a subsection with its consequence, remediation, and code
before/after if present.

## Contextually Acceptable
The steelman list (accepted) — smells harmless here, each with its one-line reason.${mode === "audit" ? " Its presence proves the review was calibrated." : ""}
${mode === "audit" ? "\n## Top refactors\nThe ranked refactors — what / why (tied to a finding) / incremental migration steps.\n\n## Verdict\nA paragraph stating where the codebase actually is — real risks and real strengths, no grade inflation.\n" : "\n## Strengths\nWhat the change does well.\n"}
Write the file, then report the path and a one-line summary.`;
}

// ── Caps ─────────────────────────────────────────────────────────────────────────────────────
//
// Every fan-out in this file gets a NUMBER. Before these existed, each multiplier was chosen by an LLM
// (the triage agent's `reviewers`, each reviewer's finding count) with nothing clamping it, so the agent
// count was unbounded in principle and ~1400 in practice on a large diff — past the runtime's own 1000
// ceiling. Caps never drop work silently: over-cap findings ship to the report tagged unverified, and a
// truncated reviewer split is logged.
const MAX_REVIEWERS_PER_DIM = 2; // a heavy dimension may split across file-groups, but only so far
const MAX_REVIEWERS_TOTAL = 10; // ceiling across all dimensions; past it, one reviewer per dimension
const MAX_FINDINGS_PER_REVIEWER = 8; // verified per reviewer; the rest are reported unverified

const SEVERITY_ORDER = { CRITICAL: 3, HIGH: 2, MEDIUM: 1, LOW: 0 };
const sevRank = (s) => SEVERITY_ORDER[String(s ?? "").toUpperCase()] ?? 0;

// ── Helpers ──────────────────────────────────────────────────────────────────────────────────

function splitFiles(files, n) {
  if (!files || n <= 1 || files.length <= 1) return [files ?? null];
  const groups = Array.from({ length: Math.min(n, files.length) }, () => []);
  files.forEach((f, i) => groups[i % groups.length].push(f));
  return groups;
}

function scopedVerdict(counts) {
  if (counts.critical > 0) return "REQUEST CHANGES";
  if (counts.high > 3) return "REQUEST CHANGES";
  if (counts.high >= 1) return "APPROVE WITH COMMENTS";
  if (counts.medium > 0) return "APPROVE WITH COMMENTS";
  return "APPROVE";
}

const AUDIT_PHASES = [
  { key: "smells", title: "Smell auditor", effort: "medium", focus: "Phase 2 — smells that pass the consequence test, plus the Accepted list. Read your phase section." },
  { key: "architecture", title: "Architecture auditor", effort: "high", focus: "Phase 3 — dependency arrows, amplification, blast radius, module responsibilities, backed by measured proxies." },
  { key: "modularization", title: "Modularization auditor", effort: "high", focus: "Phase 4 — the Parnas test: what each boundary hides, leaks with exploiting callers, boundaries to redraw (use the co-change data)." },
  { key: "practices", title: "Practices auditor", effort: "low", focus: "Phase 5 — one-line mechanical findings: file:line + consequence." },
];

// ── Orchestration ─────────────────────────────────────────────────────────────────────────────

// 1. Build the review tasks. Scoped → triage picks them dynamically; audit → the four fixed phases.
let reviewTasks = [];
if (mode === "audit") {
  reviewTasks = AUDIT_PHASES.map((p) => ({ ...p, files: null }));
} else {
  const triage = await agent(triagePrompt(), { phase: "Triage", label: "triage", schema: TRIAGE });
  const dimensions = (triage?.dimensions ?? []).filter((x) => x && x.key);
  for (const d of dimensions) {
    // The triage agent proposes the split; the script clamps it. An LLM-chosen integer never sizes a fan-out.
    const n = Math.min(MAX_REVIEWERS_PER_DIM, Math.max(1, Math.floor(d.reviewers ?? 1)));
    for (const group of splitFiles(changedFiles, n)) {
      reviewTasks.push({ key: d.key, title: d.title ?? d.key, focus: d.focus ?? "", effort: d.effort, files: group });
    }
  }
  if (reviewTasks.length > MAX_REVIEWERS_TOTAL) {
    log(`review "${name}": ${reviewTasks.length} reviewers exceeds the cap of ${MAX_REVIEWERS_TOTAL} — collapsing to one reviewer per dimension (each now reads the full change).`);
    reviewTasks = dimensions.map((d) => ({ key: d.key, title: d.title ?? d.key, focus: d.focus ?? "", effort: d.effort, files: null }));
  }
  log(`review "${name}": triaged ${new Set(reviewTasks.map((t) => t.key)).size} dimension(s) → ${reviewTasks.length} reviewer(s)`);
}

// 2. REVIEW → VERIFY as one pipeline per task: a reviewer's findings verify as soon as it returns.
const perTask = await pipeline(
  reviewTasks,
  (task) =>
    agent(reviewPrompt(task), { phase: "Review", label: `review:${task.key}`, effort: task.effort })
      .then((r) => ({ task, findings: r?.findings ?? [], accepted: r?.accepted ?? [] })),
  ({ task, findings, accepted }) => {
    if (!findings.length) return { task, confirmed: [], accepted };
    // Severity-first, so if the cap truncates it drops the least consequential findings, not arbitrary ones.
    const ranked = [...findings].sort((a, b) => sevRank(b.severity) - sevRank(a.severity));
    const toVerify = ranked.slice(0, MAX_FINDINGS_PER_REVIEWER);
    const overflow = ranked.slice(MAX_FINDINGS_PER_REVIEWER);
    if (overflow.length) {
      log(`review "${name}": ${task.key} returned ${findings.length} findings — verifying the top ${MAX_FINDINGS_PER_REVIEWER} by severity; ${overflow.length} ship unverified.`);
    }
    return agent(verifyBatchPrompt(task, toVerify), { phase: "Verify", label: `verify:${task.key}`, schema: VERIFY_BATCH }).then((res) => {
      const confirmed = [];
      const demoted = [];
      for (const v of res?.verdicts ?? []) {
        const f = toVerify[v.index];
        if (!f) continue; // an out-of-range index means the verifier lost track — drop, do not guess
        if (v.real && !v.demoteToAccepted) confirmed.push({ ...f, source: task.key });
        else if (v.real && v.demoteToAccepted) demoted.push({ smell: f.title, location: `${f.file}:${f.line ?? "?"}`, reasoning: v.reason ?? "harmless here" });
      }
      // Over-cap findings reach the report tagged, rather than disappearing without trace.
      for (const f of overflow) confirmed.push({ ...f, source: task.key, unverified: true });
      return { task, confirmed, accepted: [...accepted, ...demoted] };
    });
  },
);

// 3. Flatten confirmed findings + accepted.
const confirmed = [];
const accepted = [];
for (const entry of perTask.filter(Boolean)) {
  confirmed.push(...entry.confirmed);
  accepted.push(...entry.accepted);
}
log(`review "${name}": ${confirmed.length} confirmed finding(s), ${accepted.length} accepted (harmless)`);

// 4. Synthesize (dedup / rank / ids / refactors).
const synth = await agent(synthesizePrompt(confirmed, accepted), { phase: "Synthesize", label: "synthesize", schema: SYNTHESIS });

// 5. Verdict — derived in the script from severity counts (single source of truth).
const finalFindings = synth?.findings ?? [];
const counts = { critical: 0, high: 0, medium: 0, low: 0 };
for (const f of finalFindings) {
  const s = (f.severity ?? "").toUpperCase();
  if (s === "CRITICAL") counts.critical++;
  else if (s === "HIGH") counts.high++;
  else if (s === "MEDIUM") counts.medium++;
  else if (s === "LOW") counts.low++;
}
const verdict = mode === "audit" ? "SEE VERDICT PARAGRAPH" : scopedVerdict(counts);

// 6. Report (markdown).
await agent(reportPrompt(synth, verdict, counts), { phase: "Report", label: "report" });

return {
  name,
  mode,
  reportPath,
  verdict,
  counts,
  findings: finalFindings.length,
  accepted: (synth?.accepted ?? accepted).length,
  refactors: (synth?.refactors ?? []).length,
};
