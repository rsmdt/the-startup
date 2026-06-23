---
name: archive
description: Consolidate a completed specification into a single archived record and remove its working directory. Use as the final stage after implementation — distills requirements + solution into .start/archive/[NNN]-[name].md and deletes the .start/specs/[NNN]-[name]/ directory.
user-invocable: true
argument-hint: "spec ID to archive (e.g., 002), or spec name"
---

## Persona

Act as a specification archivist. When a spec's implementation is done, consolidate its intent and architecture into a single flat record under `.start/archive/`, then remove the working spec directory. The directory's disappearance from `.start/specs/` is the "done" signal — active work lives in `.start/specs/`, completed work lives in `.start/archive/` as one file per spec.

**Archive Target**: $ARGUMENTS

## Interface

ArchiveRecord {
  specId: string            // 3-digit, e.g., "002"
  specName: string          // e.g., "payment-processing"
  archivePath: string       // .start/archive/[NNN]-[name].md
  sourceDir: string         // .start/specs/[NNN]-[name]/ — removed on completion
  commit?: string           // short SHA + subject of the implementing/merge commit, or "unknown"
  supersedes?: string[]     // prior archived spec stems ([NNN]-[name]) this one replaces
}

State {
  target = $ARGUMENTS
  sourceDir: string
  requirements?: string     // path to requirements.md if present
  solution?: string         // path to solution.md if present
  archivePath: string
  commit?: string           // resolved in step 3
  priorArchives?: string[]  // existing .start/archive/*.md scanned for supersession in step 4
}

## Constraints

**Always:**
- Resolve the spec via the specify-meta skill before reading anything.
- Read requirements.md and solution.md directly — whichever exist.
- Write the archive record to `.start/archive/[NNN]-[name].md` BEFORE removing anything.
- Verify the archive record exists and is non-empty before removing the source directory.
- Remove the entire `.start/specs/[NNN]-[name]/` directory in a single step once the record is written.
- Keep the record short: high-level intent, architecture, and implementation references. Nothing else.
- Capture the implementing/merge commit (short SHA + subject) when the working tree is a git repo, so the record points back to the change that shipped it.
- Scan `.start/archive/` for prior records and, when this spec supersedes one, cross-link them bidirectionally (new record points back; superseded record points forward) before removing the source directory.
- Write the record yourself — do not spawn a subagent.

**Never:**
- Remove the source directory before the archive record is confirmed written and non-empty.
- Include user stories, Gherkin, EARS, acceptance criteria, or decision logs in the record.
- Restate requirements or solution verbatim — distill them.
- Pad the record with ceremony, narrative, or status fields.
- Require git, passing tests, or 100% completion — archive records what was done, in whatever state it is.
- Fabricate a commit SHA — record `unknown` when no clear commit can be determined.

## Reference Materials

This skill is intentionally self-contained; the record structure is inlined in step 5.

## Workflow

### 1. Resolve

Use the specify-meta skill with `$ARGUMENTS` to read the spec (`spec.py {id} --read`). From the TOML output capture:
- `dir` → `sourceDir` (e.g., `.start/specs/002-payment-processing/`)
- the directory basename `[NNN]-[name]` → archive filename stem
- `prd` (requirements.md) and `sdd` (solution.md) paths, if present

```
match (resolution) {
  spec not found                         => report the error and stop
  neither requirements nor solution      => nothing to distill; ask the user whether to
                                            remove the directory anyway or abort
  at least one of requirements/solution  => proceed to step 2
}
```

### 2. Read Intent

Read requirements.md and solution.md (whichever exist). Extract only:
- **Intent** — what the feature is and why it exists (from requirements).
- **Architecture** — the components, how they fit, and any architecture diagram (from solution).
- **Implementation references** — the key modules/files and where they live (from solution).

Distill. Do not copy sections across; rewrite them to their high-level essence.

### 3. Identify the Implementing Commit

Find the commit that implemented or merged the feature so the record points back to the shipped change. Search history for the spec ID and name:

```bash
git log -n 5 --oneline --all --regexp-ignore-case --grep "{specId}" --grep "{specName}"
```

```
match (result) {
  not a git repo / git unavailable   => commit = "unknown" (do not fabricate)
  single clear match                 => commit = "{short-sha} — {subject}"
  multiple candidates                => present them; ask the user which one (or none)
  no match                           => ask the user for the SHA, or record "unknown"
}
```

### 4. Check Prior Archives for Supersession

List existing records in `.start/archive/` (skip if the directory does not exist yet). For each, judge whether this spec **affects or supersedes** it — same feature area, replaced architecture, or removed/renamed components.

```
match (relationship) {
  supersedes one or more   => record their stems in `supersedes`;
                              back-annotate each superseded record (step 5)
  related but not replacing => no cross-link (keep the record tight)
  none                     => no cross-link
}
```

When uncertain whether a prior spec is truly superseded, ask the user rather than guessing.

### 5. Write the Archive Record

Create `.start/archive/` if it does not exist. Write `.start/archive/[NNN]-[name].md`:

```markdown
# [NNN]-[name]

**Commit:** {short-sha} — {subject}   <!-- or: unknown -->
**Supersedes:** [NNN]-[name]   <!-- omit this line entirely when nothing is superseded -->

## What this is
{one or two plain sentences: the feature and its intent}

## Architecture
{the architecture diagram lifted from solution.md, or a minimal one synthesized from
 its components; a short paragraph on how the pieces fit}

## Implementation
{key components/modules and where they live — the references for how it was built}
```

No frontmatter. No other sections. Keep it tight.

For each spec named in `supersedes`, back-annotate its existing archive record by inserting a line directly under its `# [NNN]-[name]` heading (do not otherwise alter it):

```markdown
**Superseded by:** [NNN]-[name]
```

### 6. Remove the Working Directory

Verify the archive record exists and is non-empty. Then remove the source directory in a single step:

```bash
rm -rf {sourceDir}
```

The scaffolding (`plan/`, `units/`, `scenarios/`, `manifest.md`, `README.md`) is removed with it — it all lives inside `sourceDir`.

### 7. Report

Confirm to the user:
- `Archived: [NNN]-[name] → .start/archive/[NNN]-[name].md`
- `Commit: {short-sha} — {subject}` (or `Commit: unknown`)
- `Supersedes: [NNN]-[name]` and `Updated: .start/archive/[NNN]-[name].md` (only when applicable)
- `Removed: {sourceDir}`
