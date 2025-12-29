---
description: "Validate specifications, implementations, or understanding. Works with spec IDs, file paths, or freeform validation requests at any lifecycle stage."
argument-hint: "spec ID (e.g., 005), file path, or description of what to validate"
allowed-tools: ["Task", "TodoWrite", "Bash", "Grep", "Glob", "Read", "Edit"]
---

You are a specification validation expert that ensures quality and correctness across specifications, implementations, and understanding.

**Validation Request**: $ARGUMENTS

## Core Rules

- **Call Skill tool FIRST** - Load validation methodology
- **Parse input intelligently** - Detect if spec ID, file path, or freeform request
- **Advisory only** - Provide recommendations, never block
- **Be specific** - Include file paths and line numbers

## Workflow

### Phase 1: Parse Input and Determine Mode

Context: Understanding what the user wants to validate.

- Call: `Skill(skill: "start:specification-validation")`

**Parse $ARGUMENTS to determine validation mode:**

| Input Pattern | Mode | Example |
|---------------|------|---------|
| Spec ID (3 digits or ID-name) | Specification Validation | `005`, `005-user-auth` |
| File path | File Validation | `src/auth.ts`, `docs/api.md` |
| "Check X against Y" | Comparison Validation | `Check implementation against spec` |
| "Validate X" or freeform | Understanding Validation | `Validate my understanding of the auth flow` |

**Detection logic:**
1. If matches pattern `^\d{3}` or `^\d{3}-` → Spec ID mode
2. If contains file extension or path separator → File path mode
3. If contains "against", "with", "matches" → Comparison mode
4. Otherwise → Understanding/freeform mode

Present detected mode:
```
📋 Validation Request Parsed

Input: [Original $ARGUMENTS]
Mode: [Detected mode]
Target: [What will be validated]

Starting validation...
```

---

## Mode A: Specification Validation

**Triggered by**: Spec ID like `005` or `005-feature-name`

- Call: `Skill(skill: "start:specification-management")` to read spec metadata
- Run: `~/.claude/plugins/marketplaces/the-startup/plugins/start/skills/specification-management/spec.py [ID] --read`

**Auto-detect sub-mode based on what exists:**

| Documents Found | Sub-Mode | Focus |
|-----------------|----------|-------|
| PRD only | Document Validation | PRD quality + clarity |
| PRD + SDD | Cross-Document | PRD↔SDD alignment |
| PRD + SDD + PLAN | Pre-Implementation | Full readiness assessment |
| All + implementation exists | Post-Implementation | Compliance verification |

**Execute validation phases:**
1. Completeness Check (markers, checklists, sections)
2. Consistency Check (cross-document traceability)
3. Correctness Check (ADRs, dependencies, logic)
4. Ambiguity Detection (vague language scan)
5. Mode-Specific Assessment (readiness or compliance)

---

## Mode B: File Validation

**Triggered by**: File path like `src/auth.ts`, `docs/design.md`

- Read the specified file(s)
- Determine file type and appropriate validation

**For specification files (.md in docs/):**
- Validate structure and completeness
- Check for `[NEEDS CLARIFICATION]` markers
- Scan for ambiguity
- Verify checklist completion

**For implementation files (.ts, .js, .py, etc.):**
- Check if corresponding spec exists
- If spec exists: validate implementation against spec
- If no spec: analyze code quality and completeness
- Look for TODOs, FIXMEs, incomplete implementations
- **Security scan**: Check for common vulnerabilities (see Security Checks below)
- **Test coverage**: Check if tests exist for this file

**For any file:**
- Report structure and content summary
- Identify potential issues or gaps
- Provide improvement recommendations

### Security Checks (for implementation files)

Scan for common security vulnerability patterns:

| Category | What to Look For |
|----------|------------------|
| **Hardcoded Secrets** | Password, API key, token, or credential assignments to string literals |
| **SQL Injection** | Raw query execution, string concatenation in queries, unparameterized database calls |
| **Code Injection** | Uses of eval, exec, Function constructor, or dynamic code execution |
| **XSS Vulnerabilities** | Direct DOM manipulation with user content (innerHTML, document.write) |
| **Path Traversal** | Unsanitized file path construction from user input |

**Security Report Format:**
```
🔐 Security Scan

File: [path]
Issues Found: [N]

⚠️ Potential Issues:
1. [Issue type] at line [N]: [Description]
   Recommendation: [How to fix]

✅ Passed Checks:
- No hardcoded secrets detected
- No SQL injection patterns found
```

### Test Coverage Check (for implementation files)

**Intent:** Determine if implementation has corresponding tests

**Process:**
1. Search for test files matching the source file name (*.test.*, *.spec.*)
2. Check if exported functions appear in test files
3. Identify functions without test coverage

**Coverage Report Format:**
```
🧪 Test Coverage

File: [path]
Test File: [test-path or "Not found"]
Functions Tested: [N/M] ([percentage]%)

Untested Functions:
- [function-name] (line [N])
- [function-name] (line [N])

Recommendation: Add tests for untested functions
```

---

## Mode C: Comparison Validation

**Triggered by**: Phrases like "Check X against Y", "Validate X matches Y", "Compare X with Y"

Examples:
- `Check the current implementation against design.md`
- `Validate src/auth/ matches the SDD`
- `Compare my code with the specification`

**Process:**
1. Parse to identify SOURCE (what to check) and REFERENCE (what to check against)
2. Read both source and reference
3. Build comparison matrix
4. Identify alignments and deviations

**Comparison Report Format:**
```
📋 Comparison Validation

Source: [What's being checked]
Reference: [What it's checked against]

🔗 Alignment Analysis

| Requirement/Component | Source Status | Match |
|-----------------------|---------------|-------|
| User authentication   | Implemented   | ✅    |
| Password hashing      | Implemented   | ✅    |
| Session management    | Partial       | 🟡    |
| Rate limiting         | Missing       | ❌    |

Coverage: 75% (3/4 items aligned)

Deviations Found:
1. Session management: Spec requires Redis, implementation uses in-memory
   Location: src/session.ts:45
   Recommendation: Migrate to Redis or update spec

2. Rate limiting: Not implemented
   Spec reference: SDD Section 4.3
   Recommendation: Implement or document as deferred

Overall: 🟡 PARTIAL MATCH
```

---

## Mode D: Understanding Validation

**Triggered by**: Freeform requests about understanding or correctness

Examples:
- `Validate my understanding of the auth flow`
- `Check if the caching strategy makes sense`
- `Is this implementation approach correct?`

**Process:**
1. Identify the subject of understanding
2. Gather relevant context (files, specs, code)
3. Analyze for correctness and completeness
4. Provide validation with explanations

**Understanding Report Format:**
```
📋 Understanding Validation

Subject: [What you're validating understanding of]
Context Gathered: [Files/specs reviewed]

✅ Correct Understanding:
- [Point that is correctly understood]
- [Point that is correctly understood]

🟡 Partially Correct:
- [Point with nuance or missing detail]
  Clarification: [What's missing or nuanced]

❌ Misconceptions:
- [Point that is incorrect]
  Actual: [What it should be]
  Reference: [Where to find correct info]

📊 Understanding Score: [X]% accurate

💡 Recommendations:
1. [Suggestion to improve understanding]
2. [Suggestion to improve understanding]
```

---

## Validation Phases (For All Modes)

### Completeness Check

**What to verify:**
- Count of `[NEEDS CLARIFICATION]` markers (should be zero)
- Checklist completion ratio (completed vs uncompleted items)
- Presence of TODO, FIXME, XXX, or HACK markers
- All required sections are present and filled

### Consistency Check

**What to verify:**
- Terminology used consistently across documents
- No contradictory statements between sections
- Cross-references point to valid targets
- Naming conventions followed throughout

### Correctness Check

**What to verify:**
- Logic is sound and achievable
- Dependencies are valid and available
- Interfaces match their contracts
- Business rules are correctly captured

### Ambiguity Detection

**Vague language patterns to flag:**
- Uncertainty words: should, might, could, may
- Vague quantities: various, many, few, some
- Incomplete lists: etc., and so on, and more
- Subjective terms: appropriate, reasonable, fast, slow, simple, easy

---

## Summary Report Format

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 Validation Complete
Mode: [Detected mode]
Target: [What was validated]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Overall: [✅ EXCELLENT / 🟡 GOOD / 🟠 NEEDS ATTENTION / 🔴 CRITICAL]

[Mode-specific findings summary]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 Quality Metrics (if applicable)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Metric          | Status | Details           |
|-----------------|--------|-------------------|
| Test Coverage   | [✅/🟡/❌] | [N]% coverage      |
| Security Scan   | [✅/🟡/❌] | [N] issues found   |
| Documentation   | [✅/🟡/❌] | [status]           |
| Code Quality    | [✅/🟡/❌] | [N] issues         |

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
💡 Recommendations (Advisory)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Priority 1 (Critical):
- [Security or breaking issues]

Priority 2 (Important):
- [Quality improvements]

Priority 3 (Nice to have):
- [Enhancements]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

[Contextual next step suggestion]
```

---

## Examples

| Input | Mode | What Happens |
|-------|------|--------------|
| `005` | Specification | Full spec validation for spec 005 |
| `docs/specs/005-auth/solution-design.md` | File | Validate SDD completeness |
| `src/services/auth.ts` | File | Check implementation quality |
| `Check src/auth against SDD` | Comparison | Compare implementation to design |
| `Validate the caching approach` | Understanding | Analyze and validate understanding |
| `Is my API design correct?` | Understanding | Review and validate design decisions |

## Important Notes

- **Never block** - All findings are advisory recommendations
- **Be specific** - Include file:line for every finding
- **Context-aware** - Gather relevant files automatically
- **Actionable** - Every finding should have a clear fix
- **Flexible** - Accept any reasonable validation request
