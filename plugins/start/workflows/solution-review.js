// The Agentic Startup — solution-design review workflow (adversarial architect panel + grade).
//
// One reusable dynamic-workflow script, shipped with the `start` plugin at workflows/solution-review.js. The
// /start:specify SKILL (solution phase) is the entry point: it launches this script in place via the
// Workflow tool, passing the paths of a DRAFTED solution design (SDD) and its requirements (PRD) as `args`.
// It ships in workflows/, so it is discoverable and namespaced as a command — but it is NOT a user
// entry point: it needs an `args` object only the skill can build (see the input contract below).
// Invoking it bare hits the fail-fast guard below and aborts with a pointer to the right command.
//
// It does NOT author the SDD. Authoring stays interactive in the skill — architecture decisions (ADRs)
// require user confirmation, so the design is human-owned. This workflow is the VERIFIER: it reviews a
// drafted SDD across independent lenses, adversarially confirms each finding, and grades the design
// against a rubric. The revise loop lives in the human review cycle (revise the SDD → re-run → re-grade),
// not inside the workflow — the workflow never edits the human-owned design.
//
// Pipeline (the canonical review shape — dimensions → find → verify → grade):
//   REVIEW  — fan out one agent per lens (requirement-coverage, AC→EARS traceability, MECE, and
//             design-risk covering security/robustness/compatibility). Each reads the drafted SDD
//             (+ the PRD for the traceability lenses) and returns FINDINGS only — it does not score itself.
//   VERIFY  — ONE adversarial verifier per lens confirms that lens's findings against the SDD (default to
//             disbelief); kill false positives so the human acts only on real issues.
//   GRADE   — computed IN-SCRIPT from the confirmed findings (see computeGrade). No agent: the rule is
//             deterministic, so it is code, and the overall value derives from the per-dimension scores.
//   REPORT  — write the review report (findings grouped by dimension + the grade) for the human.
//
// Runtime note: the script has no shell/filesystem access — only its agents do. The lens/verify agents
// read the SDD and PRD files; the script only orchestrates and branches on structured verdicts.

export const meta = {
  name: "run-solution-review",
  description: "Review a drafted solution design across independent lenses, verify each finding, grade the design",
  phases: [
    { title: "Review", detail: "parallel lenses find issues in the drafted SDD (coverage, traceability, MECE, design risk)" },
    { title: "Verify", detail: "one adversarial verifier per lens confirms that lens's findings; drop false positives" },
    { title: "Report", detail: "write the review report — confirmed findings grouped by dimension + the grade" },
  ],
};

// ── Input contract (args, supplied by /start:specify from the solution phase) ──────────────
// {
//   name:       string,          // spec slug / title (for labels + report header)
//   specId:     string,          // spec id (context only)
//   prdPath:    string,          // abs path to requirements.mdx (WHY — features + acceptance criteria)
//   sddPath:    string,          // abs path to the drafted solution.mdx (WHAT — the design under review)
//   reportPath: string,          // abs path where REPORT writes the review
//   lenses:     [{ key, title, focus, needsPrd }],  // optional override of the default lens panel
// }

// Background launches may serialize `args`; on some runtimes it is revived as a JSON string.
const parsedArgs = typeof args === "string" ? JSON.parse(args) : args;

const {
  name = "solution",
  specId = "",
  prdPath,
  sddPath,
  reportPath,
  lenses,
} = parsedArgs ?? {};

// FAIL FAST on a mis-wired launch. This script is not a user entry point: it is launched by the
// /start:specify skill (solution phase) with the paths of a DRAFTED design and its requirements.
// Invoked bare, every lens would be asked to read "undefined".
{
  const required = { prdPath, sddPath, reportPath };
  const missing = Object.entries(required)
    .filter(([, v]) => typeof v !== "string" || !v.trim() || v === "undefined")
    .map(([k]) => k);
  if (missing.length) {
    throw new Error(
      `run-solution-review: launch args missing or unresolved: ${missing.join(", ")}. ` +
        `This workflow is launched BY the /start:specify skill's solution phase, not invoked directly — ` +
        `it reviews an already-drafted solution design. Run /start:specify instead.`,
    );
  }
}

// The default review panel. Each lens is one independent angle over the drafted SDD. The traceability
// lenses (needsPrd) also read the PRD — they are the upstream coverage gate (requirement → component,
// acceptance criterion → EARS behavior) that keeps the WHY→WHAT chain from dropping anything.
const DEFAULT_LENSES = [
  {
    key: "requirement-coverage",
    title: "Requirement coverage",
    needsPrd: true,
    focus:
      "Every PRD requirement maps to at least one SDD component that owns it, and exactly one owner (no " +
      "orphan requirement, no requirement split across components). Report any requirement with no owner " +
      "or more than one.",
  },
  {
    key: "ac-ears-traceability",
    title: "AC → EARS traceability",
    needsPrd: true,
    focus:
      "Every PRD acceptance criterion has a corresponding EARS criterion in the SDD (WHEN/WHILE/IF … the " +
      "system SHALL …), including the edge and error states the criterion implies. Report any acceptance " +
      "criterion with no EARS criterion, and any EARS criterion that invents behavior the PRD does not ask for.",
  },
  {
    key: "mece",
    title: "MECE (boundaries)",
    needsPrd: false,
    focus:
      "Component responsibilities are mutually exclusive (no two components share domain logic) and " +
      "collectively exhaustive (no responsibility is unowned). Interfaces don't duplicate a consumer→provider " +
      "path; each business datum has a single owning entity. Report overlaps and gaps.",
  },
  {
    // Security, robustness and compatibility were three separate lenses. They are one now: all three ask
    // "what does this design get wrong about the world outside the happy path", and three agents each read
    // the same SDD in full to answer a third of that question. Tag each finding with which concern it hits
    // so the report can still group them.
    key: "design-risk",
    title: "Design risk (security / robustness / compatibility)",
    needsPrd: false,
    focus:
      "Three concerns, weighted equally — give each a deliberate pass before you report, and do not stop " +
      "at the first one that yields findings:\n" +
      "  (a) SECURITY — trust boundaries, authN/authZ, input handling, data protection, secret/credential " +
      "flow. Where does the architecture leave a sensitive path unprotected or a boundary unclear?\n" +
      "  (b) ROBUSTNESS — unnecessary complexity, brittle abstractions, unsafe concurrency, unhandled " +
      "failure modes. Where is the design harder to reason about than the problem requires, or a " +
      "failure/ordering/race case unaddressed?\n" +
      "  (c) COMPATIBILITY — breaking changes to existing public interfaces, schemas or contracts this " +
      "design implies, and whether a migration path exists. Which interface/version changes state no " +
      "migration or consumer impact?\n" +
      "Prefix every finding's detail with the concern it belongs to: \"security:\", \"robustness:\" or " +
      "\"compatibility:\".",
  },
];

const panel = (Array.isArray(lenses) && lenses.length ? lenses : DEFAULT_LENSES).filter((l) => l && l.key);

// ── Structured verdicts, so the loop branches on data, not prose ──────────────────────────

// Each lens returns findings only — generators find problems; they do not score themselves.
const FINDINGS = {
  type: "object",
  properties: {
    findings: {
      type: "array",
      items: {
        type: "object",
        properties: {
          severity: { type: "string", description: "blocking | advisory" },
          detail: { type: "string", description: "the specific flaw" },
          location: { type: "string", description: "the SDD section / component / interface it is in" },
          fix: { type: "string", description: "the concrete change that resolves it" },
        },
        required: ["severity", "detail"],
      },
    },
  },
  required: ["findings"],
};

// One verifier judges ALL of a lens's findings in a single pass. Per-finding verification was
// multiplicative (lenses × findings) with no cap on the finding count, so the agent total was set by how
// chatty the lenses happened to be. The adversarial framing survives verbatim; the anti-anchoring rule in
// the prompt is what replaces per-finding context isolation.
const VERIFY_BATCH = {
  type: "object",
  properties: {
    verdicts: {
      type: "array",
      items: {
        type: "object",
        properties: {
          index: { type: "number", description: "the finding's index in the numbered list you were given" },
          real: { type: "boolean", description: "true iff the SDD actually has this flaw" },
          reason: { type: "string" },
        },
        required: ["index", "real"],
      },
    },
  },
  required: ["verdicts"],
};

// Findings verified per lens. Past this, findings ship to the report tagged unverified rather than
// silently dropped — a cap must never make work disappear without trace.
const MAX_FINDINGS_PER_LENS = 8;

// Report-quality grade — a rubric of per-dimension scores (fail | partial | pass). The gate is
// per-dimension: any dimension scored `fail` means the design is not review-ready. `issues` localize the
// fixes.
//
// There is no grader agent. The rule below was previously written out in a prompt and handed to an LLM,
// which then applied it to a list the script already held — a whole agent spent re-deriving something
// deterministic. It is now `computeGrade`, and the scoring is identical by construction rather than by
// the grader's good behaviour. Shape returned to the caller is unchanged:
//   { dimensions: [{ name, score, note }], issues: [{ dimension, detail, fix }], overall: "n/m" }
const COVERAGE_DIMENSIONS = new Set(["requirement-coverage", "ac-ears-traceability"]);

function computeGrade(confirmedByDim) {
  const dimensions = panel.map((l) => {
    const items = confirmedByDim[l.key] ?? [];
    const blocking = items.filter((f) => String(f.severity).toLowerCase() === "blocking");
    // A missed requirement or an untraced acceptance criterion is not something the design can ship
    // without, so on the coverage lenses ANY confirmed gap fails — severity does not soften it.
    const score = !items.length
      ? "pass"
      : blocking.length || COVERAGE_DIMENSIONS.has(l.key)
        ? "fail"
        : "partial";
    const note =
      score === "pass"
        ? "no confirmed findings"
        : score === "fail"
          ? `${blocking.length || items.length} blocking finding(s)${COVERAGE_DIMENSIONS.has(l.key) && !blocking.length ? " — any confirmed gap fails a coverage lens" : ""}`
          : `${items.length} advisory finding(s) — ships with a note`;
    return { name: l.key, score, note };
  });

  const issues = [];
  for (const l of panel) {
    for (const f of confirmedByDim[l.key] ?? []) {
      issues.push({ dimension: l.key, detail: f.location ? `${f.detail} (${f.location})` : f.detail, fix: f.fix });
    }
  }

  const pts = dimensions.reduce((s, d) => s + ({ fail: 0, partial: 1, pass: 2 })[d.score], 0);
  return { dimensions, issues, overall: `${pts}/${dimensions.length * 2}` };
}

// ── Prompt builders ───────────────────────────────────────────────────────────────────────

function reviewPrompt(lens) {
  const sources = lens.needsPrd
    ? `Read BOTH the requirements and the design:\n  PRD (requirements): ${prdPath}\n  SDD (design under review): ${sddPath}`
    : `Read the design under review:\n  SDD: ${sddPath}`;
  return `You are reviewing a drafted solution design for the spec "${name}" through ONE lens: ${lens.title}.
${sources}

## Your lens
${lens.focus}

## How to review
Read the design in full. Judge ONLY through your lens — other lenses cover other concerns. For every
issue, cite the exact SDD location (section / component / interface / data entity) it lives in, and give
the concrete change that resolves it. Mark an issue \`blocking\` if it makes the design not sound to build
on, \`advisory\` if it is a quality concern that can ship with a note. Do not invent issues to fill a
quota — return an empty list if the design is clean on your lens.

Return { findings: [{ severity, detail, location, fix }] }.`;
}

function verifyBatchPrompt(lens, findings) {
  const list = findings
    .map((f, i) => `${i}. [${f.severity}] ${f.detail}\n   location: ${f.location ?? "(unspecified)"}`)
    .join("\n");
  return `You are the ADVERSARIAL verifier for the ${findings.length} review finding(s) below on the solution
design "${name}" (lens: ${lens.title}). Default to disbelief — assume each finding is WRONG until the
design forces the opposite conclusion.

## The findings
${list}

Open the design yourself and check them:
  SDD: ${sddPath}${lens.needsPrd ? `\n  PRD: ${prdPath}` : ""}
For EACH finding, confirm the design ACTUALLY has that flaw at the cited location. A finding that misreads
the design, duplicates a concern already handled elsewhere, or describes a non-issue is NOT real. When in
doubt, rule NOT real — a false finding wastes the author's revision.

Judge each finding ON ITS OWN EVIDENCE. The verdicts are independent: having just confirmed three findings
is no reason to wave the fourth through, and rejecting one is no reason to reject the next. A run of
identical verdicts means you have stopped opening the design — go back and open it.

Return { verdicts: [{ index, real, reason }] } — one entry per finding, \`index\` matching the numbering above.`;
}

function reportPrompt(confirmedByDim, grade) {
  const blocks = panel
    .map((l) => {
      const items = confirmedByDim[l.key] ?? [];
      const body = items.length
        ? items.map((f) => `- [${f.severity}] ${f.detail}${f.location ? ` — ${f.location}` : ""}${f.fix ? ` → ${f.fix}` : ""}`).join("\n")
        : "- (clean)";
      return `## ${l.title}\n${body}`;
    })
    .join("\n\n");
  const scores = (grade?.dimensions ?? []).map((d) => `- ${d.name}: ${d.score}${d.note ? ` — ${d.note}` : ""}`).join("\n");
  return `You are writing the solution-design review report for the spec "${name}" to: ${reportPath}

Write a clear markdown report with this structure:
# Solution Design Review — ${name}${specId ? ` (${specId})` : ""}

## Grade
Overall: ${grade?.overall ?? "(n/a)"}
${scores || "(no dimensions)"}

## Findings by dimension
${blocks}

## How to act
The solution design is human-owned — these are review findings, not edits. Resolve each blocking finding
(confirming any architecture decision with the user), then re-run this review to re-grade. Advisory
findings may ship with a note.

Write the file exactly as structured above (you may tighten wording), then report the path and a one-line summary.`;
}

// ── Orchestration: review (fan-out) → verify (pipeline) → grade → report ─────────────────

// REVIEW + VERIFY as one pipeline per lens: a lens's findings verify as soon as that lens returns, so a
// slow lens never blocks a fast one. Each stage is pinned to its phase to avoid racing the global phase().
const perLens = await pipeline(
  panel,
  (lens) => agent(reviewPrompt(lens), { phase: "Review", label: `review:${lens.key}`, schema: FINDINGS }).then((r) => ({ lens, findings: r?.findings ?? [] })),
  ({ lens, findings }) => {
    if (!findings.length) return { lens, confirmed: [] };
    // Blocking first, so if the cap truncates it drops advisory findings rather than the ones that gate.
    const ranked = [...findings].sort(
      (a, b) => (String(b.severity).toLowerCase() === "blocking") - (String(a.severity).toLowerCase() === "blocking"),
    );
    const toVerify = ranked.slice(0, MAX_FINDINGS_PER_LENS);
    const overflow = ranked.slice(MAX_FINDINGS_PER_LENS);
    if (overflow.length) {
      log(`solution review "${name}": lens ${lens.key} returned ${findings.length} findings — verifying the top ${MAX_FINDINGS_PER_LENS}; ${overflow.length} ship unverified.`);
    }
    return agent(verifyBatchPrompt(lens, toVerify), { phase: "Verify", label: `verify:${lens.key}`, schema: VERIFY_BATCH }).then((res) => {
      const confirmed = [];
      for (const v of res?.verdicts ?? []) {
        const f = toVerify[v.index];
        if (f && v.real) confirmed.push(f);
      }
      for (const f of overflow) confirmed.push({ ...f, unverified: true });
      return { lens, confirmed };
    });
  },
);

// Group confirmed findings by lens key for the grader and the report.
const confirmedByDim = {};
let confirmedCount = 0;
for (const entry of perLens.filter(Boolean)) {
  confirmedByDim[entry.lens.key] = entry.confirmed;
  confirmedCount += entry.confirmed.length;
}
log(`solution review "${name}": ${confirmedCount} confirmed finding(s) across ${panel.length} lens(es)`);

// Scored in-script from the confirmed findings — see computeGrade. No agent.
const grade = computeGrade(confirmedByDim);

await agent(reportPrompt(confirmedByDim, grade), { phase: "Report", label: "report" });

const failing = (grade?.dimensions ?? []).filter((d) => d.score === "fail").map((d) => d.name);

return {
  name,
  specId,
  reportPath,
  confirmed: confirmedCount,
  grade,
  failingDimensions: failing,
  reviewReady: failing.length === 0,
};
