# Audit Phase Instructions

Full instructions for the four audit reviewers. Each reviewer receives the orientation map (what the system does, core domain concepts, volatile vs. stable parts) and the co-change data from the orchestrator, plus its own phase section below. All judgments are relative to that map.

Standard of evidence for every phase: a finding must name the concrete future change, bug class, or comprehension cost it causes, with file:line. If you cannot name a victim, file the observation under CONTEXTUALLY ACCEPTABLE with one line of reasoning.

---

## Phase 2 — Smell Detection with Consequence Test (Smell Auditor)

Sweep for the classic catalog:

- God classes/modules
- Long parameter lists
- Feature envy
- Shotgun surgery
- Divergent change
- Primitive obsession
- Data clumps
- Deep inheritance
- Boolean-flag parameters
- Comments that apologize for the code
- Dead code
- Error swallowing
- Mixed abstraction levels within a function

For every hit, run the consequence test before reporting:

1. What realistic modification does this make harder, or what bug class does it invite? Name it specifically.
2. If you cannot name a victim, file it under CONTEXTUALLY ACCEPTABLE with one line of reasoning, not under findings.

Also hunt contextual smells the catalog misses:

- Functions that mix policy with I/O
- Names that lie about behavior
- Invariants enforced in some call paths but not others
- Error handling that differs across similar operations for no stated reason

**Output**: findings (each with file:line, consequence, remediation) + the Accepted list. Both are mandatory — an empty Accepted list means the sweep wasn't calibrated.

---

## Phase 3 — Architecture Assessment via Operational Proxies (Architecture Auditor)

Do not opine on architecture style. Measure:

### 3.1 Dependency Direction

Build the import graph between major modules. Does domain logic import frameworks/IO/vendors, or the reverse? Volatile → stable is expected; stable → volatile is a finding. List every wrong-way arrow with file:line and what change it makes risky.

### 3.2 Change Amplification

Choose 3 realistic change scenarios grounded in this domain (e.g., add a field to core entity X, swap vendor Y, add a new consumer of capability Z). For each, list every file that must change. More than 5 scattered locations for a conceptually local change is a finding — name the coupling that causes it.

### 3.3 Co-Change Archaeology

The orchestrator supplies file sets that repeatedly change in the same commits. Where co-changing files live in different modules, the drawn boundary contradicts the real one — report the mismatch. (You interpret; the raw data is already computed.)

### 3.4 Blast Radius

For the 5 highest fan-in modules: state what depends on them and whether their interfaces are stable enough to justify that fan-in. Unstable interface + high fan-in = a finding with the change that would ripple.

### 3.5 Explainability

Write a one-sentence responsibility for each major module. Every module needing "and", or whose name doesn't match its contents, is a finding — the consequence is comprehension cost, name where it bites (onboarding, code placement decisions, review misrouting).

**Output**: the five proxy datasets (structured per the ArchitectureProxies interface) + findings derived from them. Every architecture claim must cite a measurement, never a principle.

---

## Phase 4 — Modularization Audit / Parnas Test (Modularization Auditor)

For every module boundary:

### What Does It Hide?

Name the decision (data format, algorithm, vendor, protocol) concealed behind it. A boundary hiding nothing is organizational, not architectural — flag it.

### Does the Hiding Hold?

Find leaks, each cited with the caller that exploits it:

- Internal types in public signatures
- Callers reaching through the module (a.b.c.d chains)
- Config passthrough
- Escape hatches exposing internals

### Interface/Implementation Ratio

Modules whose public surface approaches their internal surface are folders pretending to be modules.

### Wrong-Boundary Evidence

Cross-reference the co-change data: co-change across a boundary plus leaks equals a boundary to redraw. Propose the redraw in one line — what moves where, what the new hidden decision is.

### The Inverse

Check for things bundled together that change for different reasons and should be split.

**Output**: per-boundary verdicts (hides / leaks / ratio), proposed redraws (one line each), findings with file:line + consequence.

---

## Phase 5 — Best-Practice Sweep (Practices Auditor)

Mechanical and terse. One line per finding: file:line, consequence.

- Error handling consistency
- Resource cleanup
- Input validation at trust boundaries
- Logging that would actually support debugging an incident
- Concurrency hazards: shared mutable state, unguarded caches
- Configuration handling
- Dependency freshness — only where a known CVE or EOL applies

Skip style-guide trivia a formatter or linter would catch — recommend the tool once instead.

**Output**: one-line findings + a single tooling recommendation if applicable.
