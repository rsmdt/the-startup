// solution-workflow: v1
//
// The Agentic Startup — solution-design review workflow (adversarial architect panel + grade).
//
// One reusable dynamic-workflow script. /start:specify (solution phase) installs this to
// <project>/.claude/workflows/solution.js and launches it via the Workflow tool, passing the paths of
// a DRAFTED solution design (SDD) and its requirements (PRD) as `args`.
//
// It does NOT author the SDD. Authoring stays interactive in the skill — architecture decisions (ADRs)
// require user confirmation, so the design is human-owned. This workflow is the VERIFIER: it reviews a
// drafted SDD across independent lenses, adversarially confirms each finding, and grades the design
// against a rubric. The revise loop lives in the human review cycle (revise the SDD → re-run → re-grade),
// not inside the workflow — the workflow never edits the human-owned design.
//
// Pipeline (the canonical review shape — dimensions → find → verify → grade):
//   REVIEW  — fan out one agent per lens (requirement-coverage, AC→EARS traceability, MECE, security,
//             robustness, compatibility). Each reads the drafted SDD (+ the PRD for the traceability
//             lenses) and returns FINDINGS only — it does not score itself.
//   VERIFY  — adversarially confirm each finding against the SDD (default to disbelief); kill false
//             positives so the human acts only on real issues.
//   GRADE   — ONE grader scores each dimension fail|partial|pass from the confirmed findings + the
//             traceability tallies, and lists the issues. The overall value is derived from the scores.
//   REPORT  — write the review report (findings grouped by dimension + the grade) for the human.
//
// Runtime note: the script has no shell/filesystem access — only its agents do. The lens/verify/grade
// agents read the SDD and PRD files; the script only orchestrates and branches on structured verdicts.

export const meta = {
  name: "solution",
  description: "Review a drafted solution design across independent lenses, verify each finding, grade the design",
  phases: [
    { title: "Review", detail: "parallel lenses find issues in the drafted SDD (coverage, traceability, MECE, security, robustness, compatibility)" },
    { title: "Verify", detail: "adversarially confirm each finding against the SDD; drop false positives" },
    { title: "Grade", detail: "score each dimension fail|partial|pass from confirmed findings + traceability" },
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
    key: "security",
    title: "Security",
    needsPrd: false,
    focus:
      "The design's trust boundaries, authN/authZ, input handling, data protection, and secret/credential " +
      "flow. Report where the architecture leaves a sensitive path unprotected or a boundary unclear.",
  },
  {
    key: "robustness",
    title: "Robustness",
    needsPrd: false,
    focus:
      "Unnecessary complexity, brittle abstractions, unsafe concurrency, and unhandled failure modes in the " +
      "design. Report where the design is harder to reason about than the problem requires, or where a " +
      "failure/ordering/race case is unaddressed.",
  },
  {
    key: "compatibility",
    title: "Compatibility",
    needsPrd: false,
    focus:
      "Breaking changes to existing public interfaces, schemas, or contracts this design implies, and whether " +
      "a migration path exists. Report interface/version changes with no stated migration or consumer impact.",
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

// The adversarial verifier's verdict on a single finding.
const VERIFY = {
  type: "object",
  properties: {
    real: { type: "boolean", description: "true iff the SDD actually has this flaw" },
    reason: { type: "string" },
  },
  required: ["real"],
};

// Report-quality grade — a rubric of per-dimension scores (fail | partial | pass). The gate is
// per-dimension: any dimension scored `fail` means the design is not review-ready. `issues` localize the
// fixes. The overall scalar is DERIVED from the scores by the script — one source of truth.
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
          detail: { type: "string" },
          fix: { type: "string" },
        },
        required: ["dimension", "detail"],
      },
    },
  },
  required: ["dimensions"],
};

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

function verifyPrompt(lens, finding) {
  return `You are the ADVERSARIAL verifier for one review finding on the solution design "${name}"
(lens: ${lens.title}). Default to disbelief — assume the finding is WRONG until the design forces the
opposite conclusion.

The finding:
  detail:   ${finding.detail}
  location: ${finding.location ?? "(unspecified)"}
  severity: ${finding.severity}

Open the design yourself and check it:
  SDD: ${sddPath}${lens.needsPrd ? `\n  PRD: ${prdPath}` : ""}
Confirm the design ACTUALLY has this flaw at the cited location. A finding that misreads the design,
duplicates a concern already handled elsewhere, or describes a non-issue is NOT real. When in doubt, rule
NOT real — a false finding wastes the author's revision.

Return { real, reason }.`;
}

function gradePrompt(confirmedByDim) {
  const blocks = panel
    .map((l) => {
      const items = confirmedByDim[l.key] ?? [];
      const body = items.length
        ? items.map((f) => `  - [${f.severity}] ${f.detail}${f.location ? ` (${f.location})` : ""}`).join("\n")
        : "  (no confirmed findings)";
      return `### ${l.title} [${l.key}]\n${body}`;
    })
    .join("\n");
  return `You are the grader for the solution design "${name}". Score the design on each dimension below
from its CONFIRMED review findings. You do not re-review — you score what the panel confirmed.

Score each dimension: \`pass\` (no real problem), \`partial\` (only advisory findings — ships with a note),
or \`fail\` (one or more blocking findings). For the two coverage dimensions (requirement-coverage,
ac-ears-traceability) any confirmed gap is \`fail\` — a missed requirement or an untraced acceptance
criterion is not something the design can ship without.

## Confirmed findings, by dimension
${blocks}

## Return
- dimensions: score EACH dimension { name, score ("fail" | "partial" | "pass"), note } — use the lens key as name.
- issues: for every dimension scored below \`pass\`, one entry per problem { dimension, detail, fix } — carry
  the confirmed finding's location and the concrete change. Be specific.`;
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
  ({ lens, findings }) =>
    parallel(
      findings.map((f) => () =>
        agent(verifyPrompt(lens, f), { phase: "Verify", label: `verify:${lens.key}`, schema: VERIFY }).then((v) => ({ ...f, real: v?.real })),
      ),
    ).then((verified) => ({ lens, confirmed: verified.filter(Boolean).filter((f) => f.real) })),
);

// Group confirmed findings by lens key for the grader and the report.
const confirmedByDim = {};
let confirmedCount = 0;
for (const entry of perLens.filter(Boolean)) {
  confirmedByDim[entry.lens.key] = entry.confirmed;
  confirmedCount += entry.confirmed.length;
}
log(`solution review "${name}": ${confirmedCount} confirmed finding(s) across ${panel.length} lens(es)`);

const grade = await agent(gradePrompt(confirmedByDim), { phase: "Grade", label: "grade", schema: GRADE });

// Derived rubric value — single source of truth is the per-dimension scores (fail=0, partial=1, pass=2).
if (grade?.dimensions?.length) {
  const pts = grade.dimensions.reduce((s, d) => s + (({ fail: 0, partial: 1, pass: 2 })[d.score] ?? 0), 0);
  grade.overall = `${pts}/${grade.dimensions.length * 2}`;
}

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
