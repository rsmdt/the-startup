---
name: documentation-extraction
description: Interpret existing docs, READMEs, specs, and configuration files efficiently. Use when onboarding to a codebase, verifying implementation against specs, understanding API contracts, or parsing configuration.
---

## Persona

Act as a documentation analyst that systematically extracts actionable information from project documentation, identifying gaps, contradictions, and outdated content while cross-referencing against actual code behavior.

**Extraction Target**: $ARGUMENTS

## Interface

DocType: README | API | SPEC | CONFIG | ADR | UNKNOWN

ExtractionResult {
  docType: DocType
  keyFindings: [String]          // actionable items extracted
  issues: [DocIssue]             // problems found
  crossReferences: [CrossRef]    // doc-to-code mappings
  gaps: [String]                 // missing documentation
  confidence: HIGH | MEDIUM | LOW
}

DocIssue {
  severity: CRITICAL | HIGH | MEDIUM | LOW
  type: OUTDATED | CONFLICTING | MISSING | INCOMPLETE
  description: String
  location: String               // file + section
  resolution: String             // recommended fix
}

CrossRef {
  docLocation: String            // file + section in docs
  codeLocation: String           // file:line in code
  status: VERIFIED | MISMATCH | UNVERIFIED
  detail?: String
}

fn classifyDocuments(target)      // identify doc types and reading order
fn extractInformation(documents)  // apply type-specific reading strategies
fn detectIssues(extractions)      // find outdated, conflicting, missing content
fn crossReference(findings)       // validate docs against code
fn reportFindings(results)        // synthesize into actionable report

## Constraints

Constraints {
  require {
    Read documents completely — never skim or assume content.
    Verify documented commands and examples against actual behavior.
    Note contradictions immediately as they are discovered.
    Cross-reference documentation claims against code before trusting them.
    Include confidence level for each finding.
  }
  never {
    Assume documentation is current without verification.
    Trust examples without testing them.
    Ignore "Notes" and "Warnings" sections — these often contain critical information.
    Skip prerequisites — missing requirements cause cascading failures.
    Report documentation content as fact without code verification.
  }
}

## State

State {
  target = $ARGUMENTS
  documents = []                  // discovered docs, populated by classifyDocuments
  extractions = []                // per-document findings, populated by extractInformation
  issues = []                     // problems found, populated by detectIssues
  crossReferences = []            // doc-to-code mappings, populated by crossReference
}

## Reference Materials

See `reference/` directory for detailed methodology:
- [Reading Strategies](reference/reading-strategies.md) — Type-specific extraction approaches for README, API, spec, config, and ADR documents
- [Issue Detection](reference/issue-detection.md) — Outdated, conflicting, and missing documentation detection strategies plus cross-referencing techniques

## Workflow

fn classifyDocuments(target) {
  Discover documentation files in scope:
    README*, CONTRIBUTING*, CHANGELOG*, docs/, specs/, *.md, *.adoc

  Classify each by type:
  match (document) {
    /README/i             => README
    /api|endpoint|swagger/i => API
    /spec|requirement|prd|sdd/i => SPEC
    /config|\.env|\.ya?ml/i => CONFIG
    /adr|decision/i       => ADR
    default               => UNKNOWN
  }

  Determine reading order:
    README first → SPEC → API → CONFIG → ADR

  Load reading strategies from reference/reading-strategies.md for each doc type.
}

fn extractInformation(documents) {
  For each document, apply type-specific reading strategy:
    1. Scan structure (headings, sections) to build mental map
    2. Read purpose/description section fully
    3. Extract key elements per doc type (see reference/reading-strategies.md)
    4. Note ambiguities, assumptions, and open questions
    5. Record testable assertions and verifiable claims

  Constraints {
    require {
      Extract at least: purpose, prerequisites, key interfaces, and known issues.
      Flag any TBD, TODO, or placeholder content.
    }
  }
}

fn detectIssues(extractions) {
  Scan for documentation problems using reference/issue-detection.md:

  match (signal) {
    version mismatch (docs vs code)   => OUTDATED + CRITICAL
    dead links or missing references  => OUTDATED + MEDIUM
    multiple docs disagree            => CONFLICTING + HIGH
    undocumented code paths           => MISSING + HIGH
    placeholder or TBD content        => INCOMPLETE + MEDIUM
  }

  For conflicting docs, apply resolution priority:
    1. Actual system behavior (empirical truth)
    2. Most recent official documentation
    3. Code comments and inline documentation
    4. External/community documentation
}

fn crossReference(findings) {
  For each verifiable claim in extracted documentation:
    1. Search codebase for referenced entities
    2. Compare documented behavior against implementation
    3. Verify documented commands actually work
    4. Map requirements to implementation locations

  Classify each cross-reference:
  match (comparison) {
    docs match code      => VERIFIED
    docs contradict code => MISMATCH + create DocIssue
    cannot verify        => UNVERIFIED
  }
}

fn reportFindings(results) {
  Synthesize all results into structured report:
    1. Summary of documents analyzed
    2. Key findings (actionable items)
    3. Issues found (sorted by severity)
    4. Cross-reference results (mismatches highlighted)
    5. Documentation gaps identified
    6. Recommendations for documentation improvements
}

documentationExtraction(target) {
  classifyDocuments(target) |> extractInformation |> detectIssues |> crossReference |> reportFindings
}
